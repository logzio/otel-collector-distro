// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package jsonlogexporter

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	_ "github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func exportLogs(ld plog.Logs, t *testing.T, cfg *Config) {
	params := componenttest.NewNopExporterCreateSettings()
	exporter, err := createLogsExporter(context.Background(), params, cfg)
	require.NotNil(t, exporter)
	require.NoError(t, err)
	ctx := context.Background()
	err = exporter.ConsumeLogs(ctx, ld)
	require.NoError(t, err)
	err = exporter.Shutdown(ctx)
	require.NoError(t, err)
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

func TestNullTracesExporterConfig(tester *testing.T) {
	params := componenttest.NewNopExporterCreateSettings()
	_, err := newJsonLogExporter(nil, params)
	assert.Error(tester, err, "Null exporter config should produce error")
}

func TestConsumeLogs(t *testing.T) {
	var lock sync.Mutex
	var recordedRequests []byte
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		recordedRequests, _ = ioutil.ReadAll(req.Body)
		rw.WriteHeader(http.StatusOK)
	}))
	cfg := Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		Region:           "us",
		Token:            "",
		CustomEndpoint:   server.URL,
		DrainInterval:    100000,
		QueueMaxLength:   500000,
		QueueCapacity:    20 * 1024 * 1024,
	}
	defer server.Close()
	var log map[string]interface{}
	ld := GenerateLogsManyLogRecordsSameResource(5)
	ld.ResourceLogs().At(0).Resource().Attributes().InsertString("att", "test_att")
	lock.Lock()
	exportLogs(ld, t, &cfg)
	lock.Unlock()
	decoded, _ := gUnzipData(recordedRequests)
	requests := strings.Split(string(decoded), "\n")
	require.Equal(t, 6, len(requests))
	assert.NoError(t, json.Unmarshal([]byte(requests[0]), &log))
	require.Equal(t, "test_att", log["att"])

	ld = GenerateLogsTwoLogRecordsSameResource()
	ld.ResourceLogs().At(0).Resource().Attributes().InsertString("test", "test_value")
	exportLogs(ld, t, &cfg)
	decoded, _ = gUnzipData(recordedRequests)
	requests = strings.Split(string(decoded), "\n")
	require.Equal(t, 3, len(requests))
	assert.NoError(t, json.Unmarshal([]byte(requests[0]), &log))
	require.Equal(t, "test_value", log["test"])

}

func TestConvertLogRecordToJson(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg := factory.CreateDefaultConfig()
	params := componenttest.NewNopExporterCreateSettings()
	exporter, err := newJsonLogExporter(cfg.(*Config), params)
	require.NotNil(t, exporter)
	require.NoError(t, err)

	type convertLogRecordToJsonTest struct {
		log      plog.LogRecord
		resource pcommon.Resource
		expected map[string]interface{}
	}

	var convertLogRecordToJsonTests = []convertLogRecordToJsonTest{
		{GenerateLogRecordWithNestedBody(),
			pcommon.NewResource(),
			map[string]interface{}{
				"23":           float64(45),
				"app":          "server",
				"foo":          "bar",
				"instance_num": float64(1),
				"level":        "Info",
				"message":      "hello there",
				"nested":       map[string]interface{}{"number": float64(499), "string": "v1"},
				"spanID":       "0102040800000000",
				"traceID":      "08040201000000000000000000000000",
			},
		},
		{GenerateLogRecordWithMultiTypeValues(),
			pcommon.NewResource(),
			map[string]interface{}{
				"bool":     true,
				"customer": "acme",
				"env":      "dev",
				"level":    "Info",
				"message":  "something happened",
				"number":   float64(64),
			},
		},
	}
	for _, test := range convertLogRecordToJsonTests {
		output := exporter.ConvertLogRecordToJson(test.log, test.resource)
		require.Equal(t, output, test.expected)
	}
}

func TestGetListenerUrl(t *testing.T) {
	var testLogger = hclog2ZapLogger{
		Zap:  zap.NewExample(),
		name: loggerName,
	}
	type getListenerUrlTest struct {
		arg1     string
		arg2     hclog2ZapLogger
		expected string
	}
	var getListenerUrlTests = []getListenerUrlTest{
		{"us", testLogger, "https://listener.logz.io:8071"},
		{"eu", testLogger, "https://listener-eu.logz.io:8071"},
		{"au", testLogger, "https://listener-au.logz.io:8071"},
		{"ca", testLogger, "https://listener-ca.logz.io:8071"},
		{"nl", testLogger, "https://listener-nl.logz.io:8071"},
		{"uk", testLogger, "https://listener-uk.logz.io:8071"},
		{"wa", testLogger, "https://listener-wa.logz.io:8071"},
		{"not-valid", testLogger, "https://listener.logz.io:8071"},
		{"", testLogger, "https://listener.logz.io:8071"},
	}
	for _, test := range getListenerUrlTests {
		output := getListenerUrl(test.arg1, &test.arg2)
		require.Equal(t, output, test.expected)
	}
}
