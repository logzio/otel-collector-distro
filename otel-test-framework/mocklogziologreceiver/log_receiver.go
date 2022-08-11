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

package mocklogziologreceiver // Package mocklogziologreceiver import "github.com/logzio/otel-collector-distro/otel-test-framework/mocklogziologreceiver"

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/plog"
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
	"go.uber.org/zap"
)

// MockLogzioLogReceiver type is used to handle spans received in the logzio span format.
type MockLogzioLogReceiver struct {
	mu     sync.Mutex
	logger *zap.Logger

	config *Config
	server *http.Server

	nextConsumer consumer.Logs
}

// New creates a new MockLogzioLogReceiver reference.
func New(
	nextConsumer consumer.Logs,
	params component.ReceiverCreateSettings,
	config *Config) (*MockLogzioLogReceiver, error) {
	if nextConsumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	lr := &MockLogzioLogReceiver{
		logger:       params.Logger,
		config:       config,
		nextConsumer: nextConsumer,
	}
	return lr, nil
}

// Start spins up the receiver's HTTP server and makes the receiver start its processing.
func (lr *MockLogzioLogReceiver) Start(_ context.Context, host component.Host) error {
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
	nr.HandleFunc("/Log", lr.HTTPHandlerFunc)

	// create a server with the handler
	lr.server = &http.Server{Handler: nr}

	// run the server on a routine
	go func() {
		host.ReportFatalError(lr.server.Serve(ln))
	}()
	return nil
}

// handleRequest parses an http request containing aws json request and passes the count of the logs to next consumer
func (lr *MockLogzioLogReceiver) handleRequest(req *http.Request) error {
	transport := "http"
	receiverCreateSettings := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger: zap.L(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	obsrecv := obsreport.NewReceiver(obsreport.ReceiverSettings{ReceiverID: lr.config.ID(), Transport: transport, ReceiverCreateSettings: receiverCreateSettings})
	ctx := obsrecv.StartLogsOp(req.Context())
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	decompressedBody, err := gUnzipData(body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyString := string(decompressedBody)
	logSplit := strings.Split(bodyString, "\n")
	logCount := len(logSplit) - 1

	// Generate log data for next consumer
	logData := plog.NewLogs()
	rlog := logData.ResourceLogs().AppendEmpty()
	ils := rlog.ScopeLogs().AppendEmpty()
	ils.LogRecords().EnsureCapacity(logCount)

	for i := 0; i < logCount; i++ {
		ils.LogRecords().AppendEmpty()
	}
	err = lr.nextConsumer.ConsumeLogs(ctx, logData)
	obsrecv.EndLogsOp(ctx, typeStr, logData.LogRecordCount(), err)
	return err
}

// HTTPHandlerFunc returns a http.HandlerFunc that handles requests
func (lr *MockLogzioLogReceiver) HTTPHandlerFunc(rw http.ResponseWriter, req *http.Request) {
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
func (lr *MockLogzioLogReceiver) Shutdown(context.Context) error {
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
