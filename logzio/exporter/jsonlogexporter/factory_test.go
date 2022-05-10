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
	"context"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configtest.CheckConfigStruct(cfg))
}
func TestCreateLogsExporter(t *testing.T) {
	factories, err := componenttest.NopFactories()
	require.NoError(t, err)
	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	require.NoError(t, err)

	params := componenttest.NewNopExporterCreateSettings()
	conf := cfg.Exporters[config.NewComponentIDWithName(typeStr, "")]
	exporter, err := factory.CreateLogsExporter(context.Background(), params, conf)
	assert.Nil(t, err)
	assert.NotNil(t, exporter)
}

func TestCreateDefaultLogsExporter(t *testing.T) {
	cfg := createDefaultConfig()
	exp, err := createLogsExporter(
		context.Background(),
		componenttest.NewNopExporterCreateSettings(),
		cfg)
	assert.NoError(t, err)
	require.NotNil(t, exp)
}
