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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testdata "internal/core/testdata"
)

func TestFileLogsExporter(t *testing.T) {
	jle := &jsonlogexporter{token: tempFileName(t)}
	require.NotNil(t, jle)

	ld := testdata.GenerateLogsTwoLogRecordsSameResource()
	assert.NoError(t, jle.ConsumeLogs(context.Background(), ld))

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
