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

package tests

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var (
	mockedConsumedResourceWithType = func() pmetric.Metrics {
		md := pmetric.NewMetrics()
		rm := md.ResourceMetrics().AppendEmpty()
		rm.Resource().Attributes().UpsertString("opencensus.resourcetype", "host")
		rm.Resource().Attributes().UpsertString("label-key", "label-value")
		m := rm.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
		m.SetName("metric-name")
		m.SetDescription("metric-description")
		m.SetUnit("metric-unit")
		m.SetDataType(pmetric.MetricDataTypeGauge)
		m.Gauge().DataPoints().AppendEmpty().SetIntVal(0)
		return md
	}()

	mockedConsumedResourceEmpty = func() pmetric.Metrics {
		md := pmetric.NewMetrics()
		rm := md.ResourceMetrics().AppendEmpty()
		m := rm.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
		m.SetName("metric-name")
		m.SetDescription("metric-description")
		m.SetUnit("metric-unit")
		m.SetDataType(pmetric.MetricDataTypeGauge)
		m.Gauge().DataPoints().AppendEmpty().SetIntVal(0)
		return md
	}()
)

type resourceProcessorTestCase struct {
	name                    string
	resourceProcessorConfig string
	mockedConsumedMetrics   pmetric.Metrics
	expectedMetrics         pmetric.Metrics
}

func getResourceProcessorTestCases() []resourceProcessorTestCase {

	tests := []resourceProcessorTestCase{
		{
			name: "update_and_rename_existing_attributes",
			resourceProcessorConfig: `
  resource:
    attributes:
    - key: label-key
      value: new-label-value
      action: update
    - key: resource-type
      from_attribute: opencensus.resourcetype
      action: upsert
    - key: opencensus.resourcetype
      action: delete
`,
			mockedConsumedMetrics: mockedConsumedResourceWithType,
			expectedMetrics: func() pmetric.Metrics {
				md := pmetric.NewMetrics()
				rm := md.ResourceMetrics().AppendEmpty()
				rm.Resource().Attributes().UpsertString("resource-type", "host")
				rm.Resource().Attributes().UpsertString("label-key", "new-label-value")
				return md
			}(),
		},
		{
			name: "set_attribute_on_empty_resource",
			resourceProcessorConfig: `
  resource:
    attributes:
    - key: additional-label-key
      value: additional-label-value
      action: insert

`,
			mockedConsumedMetrics: mockedConsumedResourceEmpty,
			expectedMetrics: func() pmetric.Metrics {
				md := pmetric.NewMetrics()
				rm := md.ResourceMetrics().AppendEmpty()
				rm.Resource().Attributes().UpsertString("additional-label-key", "additional-label-value")
				return md
			}(),
		},
	}

	return tests
}

