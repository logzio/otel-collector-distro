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
	"context"
	"go.opentelemetry.io/collector/config/confighttp"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configtest.CheckConfigStruct(cfg))
}

func TestCreateTracesExporter(t *testing.T) {

	factories, err := componenttest.NopFactories()
	require.NoError(t, err)
	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	require.NoError(t, err)

	params := componenttest.NewNopExporterCreateSettings()
	exporter, err := factory.CreateTracesExporter(context.Background(), params, cfg.Exporters[config.NewComponentIDWithName(typeStr, "2")])
	assert.Nil(t, err)
	assert.NotNil(t, exporter)
}

func TestGenerateUrl(t *testing.T) {
	type generateUrlTest struct {
		endpoint string
		region   string
		expected string
	}
	var generateUrlTests = []generateUrlTest{
		{"", "us", "https://listener.logz.io:8071/?token=token"},
		{"", "", "https://listener.logz.io:8071/?token=token"},
		{"https://doesnotexist.com", "", "https://doesnotexist.com"},
		{"https://doesnotexist.com", "us", "https://doesnotexist.com"},
		{"https://doesnotexist.com", "not-valid", "https://doesnotexist.com"},
		{"", "not-valid", "https://listener.logz.io:8071/?token=token"},
		{"", "US", "https://listener.logz.io:8071/?token=token"},
		{"", "Us", "https://listener.logz.io:8071/?token=token"},
		{"", "EU", "https://listener-eu.logz.io:8071/?token=token"},
	}
	for _, test := range generateUrlTests {
		cfg := &Config{
			Region:           test.region,
			Token:            "token",
			ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
			HTTPClientSettings: confighttp.HTTPClientSettings{
				Endpoint: test.endpoint,
			},
		}
		output, _ := generateEndpoint(cfg)
		require.Equal(t, test.expected, output)
	}
}

func TestGetListenerUrl(t *testing.T) {
	type getListenerUrlTest struct {
		arg1     string
		expected string
	}
	var getListenerUrlTests = []getListenerUrlTest{
		{"us", "https://listener.logz.io:8071"},
		{"eu", "https://listener-eu.logz.io:8071"},
		{"au", "https://listener-au.logz.io:8071"},
		{"ca", "https://listener-ca.logz.io:8071"},
		{"nl", "https://listener-nl.logz.io:8071"},
		{"uk", "https://listener-uk.logz.io:8071"},
		{"wa", "https://listener-wa.logz.io:8071"},
		{"not-valid", "https://listener.logz.io:8071"},
		{"", "https://listener.logz.io:8071"},
		{"US", "https://listener.logz.io:8071"},
		{"Us", "https://listener.logz.io:8071"},
	}
	for _, test := range getListenerUrlTests {
		output := getListenerUrl(test.arg1)
		require.Equal(t, output, test.expected)
	}
}
