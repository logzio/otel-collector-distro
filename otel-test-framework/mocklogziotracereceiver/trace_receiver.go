// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocklogziotracereceiver // Package mocklogziotracereceiver import "github.com/logzio/otel-collector-distro/otel-test-framework/mocklogziotracereceiver"

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// MockLogzioTraceReceiver type is used to handle spans received in the logzio span format.
type MockLogzioTraceReceiver struct {
	mu     sync.Mutex
	logger *zap.Logger

	config *Config
	server *http.Server

	nextConsumer consumer.Traces
}

// New creates a new MockLogzioTraceReceiver reference.
func New(
	nextConsumer consumer.Traces,
	params component.ReceiverCreateSettings,
	config *Config) (*MockLogzioTraceReceiver, error) {
	if nextConsumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	lr := &MockLogzioTraceReceiver{
		logger:       params.Logger,
		config:       config,
		nextConsumer: nextConsumer,
	}
	return lr, nil
}

// Start spins up the receiver's HTTP server and makes the receiver start its processing.
func (lr *MockLogzioTraceReceiver) Start(_ context.Context, host component.Host) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	// set up the listener
	ln, err := net.Listen("tcp", lr.config.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", lr.config.Endpoint, err)
	}
	lr.logger.Info(fmt.Sprintf("listen to address %s", lr.config.Endpoint))

	// use gorilla mux to create a router/handler
	nr := mux.NewRouter()
	nr.HandleFunc("/Trace", lr.HTTPHandlerFunc)

	// create a server with the handler
	lr.server = &http.Server{Handler: nr}

	// run the server on a routine
	go func() {
		host.ReportFatalError(lr.server.Serve(ln))
	}()
	return nil
}

// handleRequest parses an http request containing aws json request and passes the count of the traces to next consumer
func (lr *MockLogzioTraceReceiver) handleRequest(req *http.Request) error {
	transport := "http"
	receiverCreateSettings := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			TracerProvider: trace.NewNoopTracerProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	obsrecv := obsreport.NewReceiver(obsreport.ReceiverSettings{ReceiverID: lr.config.ID(), Transport: transport, ReceiverCreateSettings: receiverCreateSettings})
	ctx := obsrecv.StartTracesOp(req.Context())
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	decompressedBody, err := gUnzipData(body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyString := string(decompressedBody)
	spansAndServicesSplit := strings.Split(bodyString, "\n")
	spanCount := len(spansAndServicesSplit) - 1

	var span map[string]interface{}
	for i := 0; i < len(spansAndServicesSplit); i++ {
		if err = json.Unmarshal([]byte(spansAndServicesSplit[i]), &span); err == nil {
			if span["type"] != "jaegerSpan" {
				// Remove jaeger service for correct count
				spanCount = spanCount - 1
			}
		}
	}

	// Generate trace data for next consumer
	traceData := ptrace.NewTraces()
	rspan := traceData.ResourceSpans().AppendEmpty()
	ils := rspan.ScopeSpans().AppendEmpty()
	ils.Spans().EnsureCapacity(spanCount)

	for i := 0; i < spanCount; i++ {
		ils.Spans().AppendEmpty()
	}
	err = lr.nextConsumer.ConsumeTraces(ctx, traceData)
	obsrecv.EndTracesOp(ctx, typeStr, traceData.SpanCount(), err)
	return err
}

// HTTPHandlerFunc returns a http.HandlerFunc that handles requests
func (lr *MockLogzioTraceReceiver) HTTPHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	// handle the request payload
	err := lr.handleRequest(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Shutdown tells the receiver that should stop reception,
// giving it a chance to perform any necessary clean-up and shutting down
// its HTTP server.
func (lr *MockLogzioTraceReceiver) Shutdown(context.Context) error {
	return lr.server.Close()
}

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}
