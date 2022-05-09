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

package jsonlogexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jsonlogexporter"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logzio/logzio-go"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	_ "go.opentelemetry.io/collector/model/otlp"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"log"
	"os"
)

// jsonlogexporter is the implementation of file exporter that writes telemetry data to a file
// in Protobuf-JSON format.
type jsonlogexporter struct {
	token  string
	sender logzio.LogzioSender
}

func newJsonLogExporter(config *Config, params component.ExporterCreateSettings) (*jsonlogexporter, error) {
	l, err := logzio.New(
		config.Token,
		logzio.SetDebug(os.Stderr),
		logzio.SetUrl(getListenerUrl(config.Region)),
		logzio.SetInMemoryQueue(true),
		logzio.SetCompress(true),
		logzio.SetlogCountLimit(config.QueueMaxLength),
		logzio.SetinMemoryCapacity(uint64(config.QueueCapacity)),
	)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("exporter config can't be null")
	}
	if err != nil {
		return nil, err
	}

	return &jsonlogexporter{
		token:  config.Token,
		sender: *l,
	}, nil
}

func newJsonLogLogsExporter(config *Config, set component.ExporterCreateSettings) (component.LogsExporter, error) {
	exporter, err := newJsonLogExporter(config, set)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return exporterhelper.NewLogsExporter(
		config,
		set,
		exporter.ConsumeLogs,
		exporterhelper.WithShutdown(exporter.Shutdown))
}

func getListenerUrl(region string) string {
	var url string
	switch region {
	case "us":
		url = "https://listener.logz.io:8071"
	case "ca":
		url = "https://listener-ca.logz.io:8071"
	case "eu":
		url = "https://listener-eu.logz.io:8071"
	case "uk":
		url = "https://listener-uk.logz.io:8071"
	case "au":
		url = "https://listener-au.logz.io:8071"
	case "nl":
		url = "https://listener-nl.logz.io:8071"
	case "wa":
		url = "https://listener-wa.logz.io:8071"
	default:
		log.Printf("Region '%s' is not supported yet, setting url to default value", region)
		url = "https://listener.logz.io:8071"
	}
	log.Printf("Setting logzio listener url to: %s", url)
	return url
}

func convertAttributeValue(value pcommon.Value) interface{} {
	switch value.Type() {
	case pcommon.ValueTypeInt:
		return value.IntVal()
	case pcommon.ValueTypeBool:
		return value.BoolVal()
	case pcommon.ValueTypeDouble:
		return value.DoubleVal()
	case pcommon.ValueTypeString:
		return value.StringVal()
	case pcommon.ValueTypeMap:
		values := map[string]interface{}{}
		value.MapVal().Range(func(k string, v pcommon.Value) bool {
			values[k] = convertAttributeValue(v)
			return true
		})
		return values
	case pcommon.ValueTypeSlice:
		arrayVal := value.SliceVal()
		values := make([]interface{}, arrayVal.Len())
		for i := 0; i < arrayVal.Len(); i++ {
			values[i] = convertAttributeValue(arrayVal.At(i))
		}
		return values
	case pcommon.ValueTypeEmpty:
		return nil
	default:
		return value
	}
}

func (e *jsonlogexporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConvertLogRecordToJson Takes `plog.LogRecord` and `pcommon.Resource` input, outputs byte array that represents the log record as json string
func (e *jsonlogexporter) ConvertLogRecordToJson(log plog.LogRecord, resource pcommon.Resource) ([]byte, error) {
	jsonLog := map[string]interface{}{}
	if spanID := log.SpanID().HexString(); spanID != "" {
		jsonLog["spanID"] = spanID
	}
	if traceID := log.TraceID().HexString(); traceID != "" {
		jsonLog["traceID"] = traceID
	}
	if log.SeverityText() != "" {
		jsonLog["level"] = log.SeverityText()
	}
	jsonLog["@timestamp"] = log.Timestamp().AsTime()

	// add resource attributes to each json log
	resource.Attributes().Range(func(k string, v pcommon.Value) bool {
		jsonLog[k] = convertAttributeValue(v)
		return true
	})
	// add log record attributes to each json log
	log.Attributes().Range(func(k string, v pcommon.Value) bool {
		jsonLog[k] = convertAttributeValue(v)
		return true
	})

	switch log.Body().Type() {
	case pcommon.ValueTypeString:
		jsonLog["message"] = log.Body().StringVal()
	case pcommon.ValueTypeMap:
		bodyFieldsMap := convertAttributeValue(log.Body()).(map[string]interface{})
		for key, value := range bodyFieldsMap {
			jsonLog[key] = value
		}
	}
	buf, err := json.Marshal(jsonLog)
	fmt.Printf("json data: %s\n", buf)
	return buf, err
}

func (e *jsonlogexporter) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	resourceLogs := ld.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		resource := resourceLogs.At(i).Resource()
		scopeLogs := resourceLogs.At(i).ScopeLogs()
		for j := 0; j < scopeLogs.Len(); j++ {
			logRecords := scopeLogs.At(j).LogRecords()
			for k := 0; k < logRecords.Len(); k++ {
				log := logRecords.At(k)
				buf, err := e.ConvertLogRecordToJson(log, resource)
				if err != nil {
					return err
				}
				e.sender.Send(buf)
			}
		}
	}
	e.sender.Drain()
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
func (e *jsonlogexporter) Shutdown(context.Context) error {
	e.sender.Stop()
	return nil
}
