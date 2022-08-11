// Copyright 2020 OpenTelemetry Authors
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

// nolint:gocritic
package customdatareceivers // import "github.com/logzio/otel-collector-distro/otel-test-framework/customdatareceivers"

import (
	"context"
	"fmt"
	"github.com/logzio/otel-collector-distro/otel-test-framework/mocklogziotracereceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/testbed/testbed"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
)

// MockLogzioTraceDataReceiver implements logzio span format receiver.
type MockLogzioTraceDataReceiver struct {
	testbed.DataReceiverBase
	receiver component.TracesReceiver
}

// NewMockLogzioTraceDataReceiver creates a new  MockDataReceiver
func NewMockLogzioTraceDataReceiver(port int) *MockLogzioTraceDataReceiver {
	return &MockLogzioTraceDataReceiver{DataReceiverBase: testbed.DataReceiverBase{Port: port}}
}

//Start listening on the specified port
func (ar *MockLogzioTraceDataReceiver) Start(tc consumer.Traces, _ consumer.Metrics, _ consumer.Logs) error {
	var err error
	mockDataReceiverCFG := mocklogziotracereceiver.Config{
		Endpoint: fmt.Sprintf("localhost:%d", ar.Port),
	}
	ar.receiver, err = mocklogziotracereceiver.New(tc, componenttest.NewNopReceiverCreateSettings(), &mockDataReceiverCFG)

	if err != nil {
		return err
	}

	return ar.receiver.Start(context.Background(), componenttest.NewNopHost())
}

func (ar *MockLogzioTraceDataReceiver) Stop() error {
	return ar.receiver.Shutdown(context.Background())
}

func (ar *MockLogzioTraceDataReceiver) GenConfigYAMLStr() string {
	// Note that this generates an exporter config for agent.
	return fmt.Sprintf(`
  logzio:
    account_token: token
    endpoint: http://localhost:%d/Trace`, ar.Port)
}

func (ar *MockLogzioTraceDataReceiver) ProtocolName() string {
	return "logzio"
}
