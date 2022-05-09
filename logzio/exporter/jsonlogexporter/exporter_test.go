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
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"internal/core/testdata"
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonLogLogsExporter(t *testing.T) {
	cfg := Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		Region:           "us",
		Token:            "",
		DrainInterval:    3,
		QueueMaxLength:   500000,
		QueueCapacity:    20 * 1024 * 1024,
	}
	params := componenttest.NewNopExporterCreateSettings()
	exporter, err := createLogsExporter(context.Background(), params, &cfg)
	require.NotNil(t, exporter)
	require.NoError(t, err)
	ctx := context.Background()
	ld := testdata.GenerateLogsTwoLogRecordsSameResource()
	assert.NoError(t, exporter.ConsumeLogs(ctx, ld))

}

// tempFileName provides a temporary file name for testing.
func tempFileName(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("", "*.json")
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())
	socket := tmpfile.Name()
	require.NoError(t, os.Remove(socket))
	return socket
}

// errorWriter is an io.Writer that will return an error all ways
type errorWriter struct {
}

func (e errorWriter) Write([]byte) (n int, err error) {
	return 0, errors.New("all ways return error")
}

func (e *errorWriter) Close() error {
	return nil
}
