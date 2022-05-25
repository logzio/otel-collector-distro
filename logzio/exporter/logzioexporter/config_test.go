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
	"go.opentelemetry.io/collector/config/configcompression"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.Nil(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 2, len(cfg.Exporters))

	cfgExp := cfg.Exporters[config.NewComponentIDWithName(typeStr, "2")]
	expected := &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentIDWithName(typeStr, "2")),
		TracesToken:      "token",
		Region:           "eu",
	}
	expected.RetrySettings = exporterhelper.NewDefaultRetrySettings()
	expected.RetrySettings.MaxInterval = 5 * time.Second
	expected.QueueSettings = exporterhelper.NewDefaultQueueSettings()
	expected.QueueSettings.Enabled = false
	expected.HTTPClientSettings = confighttp.HTTPClientSettings{
		Endpoint: "",
		Timeout:  30 * time.Second,
		Headers:  map[string]string{},
		// Default to gzip compression
		Compression: configcompression.Gzip,
		// We almost read 0 bytes, so no need to tune ReadBufferSize.
		WriteBufferSize: 512 * 1024,
	}
	assert.Equal(t, expected, cfgExp)
}

func TestDefaultLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.Nil(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "configd.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 2, len(cfg.Exporters))

	cfgExp := cfg.Exporters[config.NewComponentIDWithName(typeStr, "2")]
	expected := &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentIDWithName(typeStr, "2")),
		TracesToken:      "logzioTESTtoken",
	}
	expected.RetrySettings = exporterhelper.NewDefaultRetrySettings()
	expected.QueueSettings = exporterhelper.NewDefaultQueueSettings()
	expected.HTTPClientSettings = confighttp.HTTPClientSettings{
		Endpoint: "",
		Timeout:  30 * time.Second,
		Headers:  map[string]string{},
		// Default to gzip compression
		Compression: configcompression.Gzip,
		// We almost read 0 bytes, so no need to tune ReadBufferSize.
		WriteBufferSize: 512 * 1024,
	}
	assert.Equal(t, expected, cfgExp)
}
