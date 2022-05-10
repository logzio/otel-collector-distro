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

package testdata

import "go.opentelemetry.io/collector/pdata/pcommon"

var (
	resourceAttributes1 = map[string]interface{}{"resource-attr": "resource-attr-val-1"}
	resourceAttributes2 = map[string]interface{}{"resource-attr": "resource-attr-val-2"}
	spanEventAttributes = map[string]interface{}{"span-event-attr": "span-event-attr-val"}
	spanLinkAttributes  = map[string]interface{}{"span-link-attr": "span-link-attr-val"}
	spanAttributes      = map[string]interface{}{"span-attr": "span-attr-val"}
)

func initResourceAttributes1(dest pcommon.Map) {
	pcommon.NewMapFromRaw(resourceAttributes1).CopyTo(dest)
}

func initResourceAttributes2(dest pcommon.Map) {
	pcommon.NewMapFromRaw(resourceAttributes2).CopyTo(dest)
}
