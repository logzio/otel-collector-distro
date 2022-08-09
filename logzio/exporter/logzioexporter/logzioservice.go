// Copyright  The OpenTelemetry Authors
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

package logzioexporter // import "github.com/logzio/otel-collector-distro/logzio/exporter/logzioexporter"

import (
	"fmt"
	"hash/fnv"

	"github.com/jaegertracing/jaeger/model"
)

const serviceLogType = "jaegerService"

// LogzioService type, for query purposes
type LogzioService struct {
	OperationName string `json:"operationName"`
	ServiceName   string `json:"serviceName"`
	Type          string `json:"type"`
}

// newLogzioService creates a new logzio service from a span
func newLogzioService(span *model.Span) LogzioService {
	service := LogzioService{
		ServiceName:   span.Process.ServiceName,
		OperationName: span.OperationName,
		Type:          serviceLogType,
	}
	return service
}

// HashCode receives a logzio service and returns a hash representation of it's service name and operation name.
func (service *LogzioService) HashCode() (string, error) {
	hash := fnv.New64a()
	_, err := hash.Write(append([]byte(service.ServiceName), []byte(service.OperationName)...))
	return fmt.Sprintf("%x", hash.Sum64()), err
}
