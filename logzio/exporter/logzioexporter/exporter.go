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

package logzioexporter // Package logzioexporter import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter"

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"google.golang.org/genproto/googleapis/rpc/status"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	loggerName               = "logzio-exporter"
	headerRetryAfter         = "Retry-After"
	maxHTTPResponseReadBytes = 64 * 1024
)

// logzioExporter implements an OpenTelemetry trace exporter that exports all spans to Logz.io
type logzioExporter struct {
	config   *Config
	client   *http.Client
	logger   hclog.Logger
	settings component.TelemetrySettings
}

func newLogzioExporter(cfg *Config, params component.ExporterCreateSettings) (*logzioExporter, error) {
	logger := hclog2ZapLogger{
		Zap:  params.Logger,
		name: loggerName,
	}
	if cfg == nil {
		return nil, errors.New("exporter config can't be null")
	}
	return &logzioExporter{
		config:   cfg,
		logger:   &logger,
		settings: params.TelemetrySettings,
	}, nil
}

func newLogzioTracesExporter(config *Config, set component.ExporterCreateSettings) (component.TracesExporter, error) {
	exporter, err := newLogzioExporter(config, set)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	exporter.config.Endpoint, err = generateEndpoint(config, config.Region)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewTracesExporter(
		config,
		set,
		exporter.pushTraceData,
		exporterhelper.WithStart(exporter.start),
		exporterhelper.WithQueue(config.QueueSettings),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

func (exporter *logzioExporter) start(_ context.Context, host component.Host) error {
	client, err := exporter.config.HTTPClientSettings.ToClient(host.GetExtensions(), exporter.settings)
	if err != nil {
		return err
	}
	exporter.client = client
	return nil
}

func (exporter *logzioExporter) pushTraceData(ctx context.Context, traces ptrace.Traces) error {
	var dataBuffer bytes.Buffer
	batches, err := jaeger.ProtoFromTraces(traces)
	if err != nil {
		return err
	}
	for _, batch := range batches {
		for _, span := range batch.Spans {
			span.Process = batch.Process
			logzioSpan, err := TransformToLogzioSpanBytes(span)
			dataBuffer.Write(append(logzioSpan, '\n'))
			if err != nil {
				return err
			}
		}
	}
	err = exporter.export(ctx, exporter.config.Endpoint, dataBuffer.Bytes())
	dataBuffer.Reset()
	return err
}

func (exporter *logzioExporter) export(ctx context.Context, url string, request []byte) error {
	exporter.logger.Debug("Preparing to make HTTP request")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(request))
	if err != nil {
		return consumererror.NewPermanent(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := exporter.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make an HTTP request: %w", err)
	}

	defer func() {
		// Discard any remaining response body when we are done reading.
		io.CopyN(ioutil.Discard, resp.Body, maxHTTPResponseReadBytes) // nolint:errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Request is successful.
		exporter.logger.Debug("Request is successful")
		return nil
	}
	respStatus := readResponse(resp)

	// Format the error message. Use the status if it is present in the response.
	var formattedErr error
	if respStatus != nil {
		formattedErr = fmt.Errorf(
			"error exporting items, request to %s responded with HTTP Status Code %d, Message=%s, Details=%v",
			url, resp.StatusCode, respStatus.Message, respStatus.Details)
	} else {
		formattedErr = fmt.Errorf(
			"error exporting items, request to %s responded with HTTP Status Code %d",
			url, resp.StatusCode)
	}

	// Check if the server is overwhelmed.
	// See spec https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md#throttling-1
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
		// Fallback to 0 if the Retry-After header is not present. This will trigger the
		// default backoff policy by our caller (retry handler).
		retryAfter := 0
		if val := resp.Header.Get(headerRetryAfter); val != "" {
			if seconds, err2 := strconv.Atoi(val); err2 == nil {
				retryAfter = seconds
			}
		}
		// Indicate to our caller to pause for the specified number of seconds.
		return exporterhelper.NewThrottleRetry(formattedErr, time.Duration(retryAfter)*time.Second)
	}

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return consumererror.NewPermanent(formattedErr)
	}

	// All other errors are retryable, so don't wrap them in consumererror.NewPermanent().
	return formattedErr
}

// Read the response and decode the status.Status from the body.
// Returns nil if the response is empty or cannot be decoded.
func readResponse(resp *http.Response) *status.Status {
	var respStatus *status.Status
	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		// Request failed. Read the body. OTLP spec says:
		// "Response body for all HTTP 4xx and HTTP 5xx responses MUST be a
		// Protobuf-encoded Status message that describes the problem."
		maxRead := resp.ContentLength
		if maxRead == -1 || maxRead > maxHTTPResponseReadBytes {
			maxRead = maxHTTPResponseReadBytes
		}
		respBytes := make([]byte, maxRead)
		n, err := io.ReadFull(resp.Body, respBytes)
		if err == nil && n > 0 {
			// Decode it as Status struct. See https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md#failures
			respStatus = &status.Status{}
			err = proto.Unmarshal(respBytes, respStatus)
			if err != nil {
				respStatus = nil
			}
		}
	}

	return respStatus
}
