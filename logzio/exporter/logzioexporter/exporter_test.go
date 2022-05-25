// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logzioexporter

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter/objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.6.1"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testService   = "testService"
	testHost      = "testHost"
	testOperation = "testOperation"
)

func newTestTracesWithAttributes() ptrace.Traces {
	td := ptrace.NewTraces()
	for i := 0; i < 10; i++ {
		s := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
		s.SetName(fmt.Sprintf("%s-%d", testOperation, i))
		s.SetTraceID(pcommon.NewTraceID([16]byte{byte(i), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}))
		s.SetSpanID(pcommon.NewSpanID([8]byte{byte(i), 0, 0, 0, 0, 0, 0, 2}))
		for j := 0; j < 5; j++ {
			s.Attributes().Insert(fmt.Sprintf("k%d", j), pcommon.NewValueString(fmt.Sprintf("v%d", j)))
		}
		s.SetKind(ptrace.SpanKindServer)
	}
	return td
}

func newTestTraces() ptrace.Traces {
	td := ptrace.NewTraces()
	s := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	s.SetName(testOperation)
	s.SetTraceID(pcommon.NewTraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}))
	s.SetSpanID(pcommon.NewSpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 2}))
	s.SetKind(ptrace.SpanKindServer)
	return td
}

func testTracesExporter(td ptrace.Traces, t *testing.T, cfg *Config) {
	params := componenttest.NewNopExporterCreateSettings()
	exporter, err := CreateTracesExporter(context.Background(), params, cfg)
	err = exporter.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	ctx := context.Background()
	err = exporter.ConsumeTraces(ctx, td)
	require.NoError(t, err)
	err = exporter.Shutdown(ctx)
	require.NoError(t, err)
}

func TestNullTracesExporterConfig(tester *testing.T) {
	params := componenttest.NewNopExporterCreateSettings()
	_, err := newLogzioTracesExporter(nil, params)
	assert.Error(tester, err, "Null exporter config should produce error")
}

func TestNullExporterConfig(tester *testing.T) {
	params := componenttest.NewNopExporterCreateSettings()
	_, err := newLogzioExporter(nil, params)
	assert.Error(tester, err, "Null exporter config should produce error")
}

func TestNullTokenConfig(tester *testing.T) {
	cfg := Config{
		Region: "eu",
	}
	params := componenttest.NewNopExporterCreateSettings()
	_, err := CreateTracesExporter(context.Background(), params, &cfg)
	assert.Error(tester, err, "Empty token should produce error")
}

func TestEmptyNode(tester *testing.T) {
	cfg := Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		TracesToken:      "test",
		Region:           "eu",
	}
	testTracesExporter(ptrace.NewTraces(), tester, &cfg)
}

//func TestConversionTraceError(tester *testing.T) {
//	cfg := Config{
//		TracesToken: "test",
//		Region:      "eu",
//	}
//	params := componenttest.NewNopExporterCreateSettings()
//	exporter, _ := newLogzioExporter(&cfg, params)
//	oldFunc := exporter.InternalTracesToJaegerTraces
//	defer func() { exporter.InternalTracesToJaegerTraces = oldFunc }()
//	exporter.InternalTracesToJaegerTraces = func(td ptrace.Traces) ([]*model.Batch, error) {
//		return nil, errors.New("fail")
//	}
//	err := exporter.pushTraceData(context.Background(), newTestTraces())
//	assert.Error(tester, err)
//}

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

func TestPushTraceData(tester *testing.T) {
	var recordedRequests []byte
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		recordedRequests, _ = ioutil.ReadAll(req.Body)
		rw.WriteHeader(http.StatusOK)
	}))
	cfg := Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		TracesToken:      "",
		Region:           "us",
	}
	defer server.Close()
	td := newTestTraces()
	res := td.ResourceSpans().At(0).Resource()
	res.Attributes().UpsertString(conventions.AttributeServiceName, testService)
	res.Attributes().UpsertString(conventions.AttributeHostName, testHost)
	testTracesExporter(td, tester, &cfg)

	var logzioSpan objects.LogzioSpan
	decoded, _ := gUnzipData(recordedRequests)
	requests := strings.Split(string(decoded), "\n")
	assert.NoError(tester, json.Unmarshal([]byte(requests[0]), &logzioSpan))
	assert.Equal(tester, testOperation, logzioSpan.OperationName)
	assert.Equal(tester, testService, logzioSpan.Process.ServiceName)

	var logzioService objects.LogzioService
	assert.NoError(tester, json.Unmarshal([]byte(requests[1]), &logzioService))

	assert.Equal(tester, testOperation, logzioService.OperationName)
	assert.Equal(tester, testService, logzioService.ServiceName)

}
