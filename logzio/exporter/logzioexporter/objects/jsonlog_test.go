package objects

import (
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"testing"
	"time"
)

var (
	TestLogTime      = time.Now()
	TestLogTimestamp = pcommon.NewTimestampFromTime(TestLogTime)
)

// Logs
func GenerateLogRecordWithNestedBody() plog.LogRecord {
	lr := plog.NewLogRecord()
	fillLogOne(lr)
	return lr
}
func GenerateLogRecordWithMultiTypeValues() plog.LogRecord {
	lr := plog.NewLogRecord()
	fillLogTwo(lr)
	return lr
}

func fillLogOne(log plog.LogRecord) {
	log.SetTimestamp(TestLogTimestamp)
	log.SetDroppedAttributesCount(1)
	log.SetSeverityNumber(plog.SeverityNumberINFO)
	log.SetSeverityText("Info")
	log.SetSpanID(pcommon.NewSpanID([8]byte{0x01, 0x02, 0x04, 0x08}))
	log.SetTraceID(pcommon.NewTraceID([16]byte{0x08, 0x04, 0x02, 0x01}))

	attrs := log.Attributes()
	attrs.InsertString("app", "server")
	attrs.InsertDouble("instance_num", 1)

	// nested body map
	attVal := pcommon.NewValueMap()
	attNestedVal := pcommon.NewValueMap()

	attMap := attVal.MapVal()
	attMap.InsertDouble("23", 45)
	attMap.InsertString("foo", "bar")
	attMap.InsertString("message", "hello there")
	attNestedMap := attNestedVal.MapVal()
	attNestedMap.InsertString("string", "v1")
	attNestedMap.InsertDouble("number", 499)
	attMap.Insert("nested", attNestedVal)
	attVal.CopyTo(log.Body())

}

func fillLogTwo(log plog.LogRecord) {
	log.SetTimestamp(TestLogTimestamp)
	log.SetDroppedAttributesCount(1)
	log.SetSeverityNumber(plog.SeverityNumberINFO)
	log.SetSeverityText("Info")

	attrs := log.Attributes()
	attrs.InsertString("customer", "acme")
	attrs.InsertDouble("number", 64)
	attrs.InsertBool("bool", true)
	attrs.InsertString("env", "dev")
	log.Body().SetStringVal("something happened")
}

func TestConvertLogRecordToJson(t *testing.T) {
	logger := hclog.NewNullLogger()
	type convertLogRecordToJsonTest struct {
		log      plog.LogRecord
		resource pcommon.Resource
		expected map[string]interface{}
	}

	var convertLogRecordToJsonTests = []convertLogRecordToJsonTest{
		{GenerateLogRecordWithNestedBody(),
			pcommon.NewResource(),
			map[string]interface{}{
				"23":           float64(45),
				"app":          "server",
				"foo":          "bar",
				"instance_num": float64(1),
				"level":        "Info",
				"message":      "hello there",
				"nested":       map[string]interface{}{"number": float64(499), "string": "v1"},
				"spanID":       "0102040800000000",
				"traceID":      "08040201000000000000000000000000",
			},
		},
		{GenerateLogRecordWithMultiTypeValues(),
			pcommon.NewResource(),
			map[string]interface{}{
				"bool":     true,
				"customer": "acme",
				"env":      "dev",
				"level":    "Info",
				"message":  "something happened",
				"number":   float64(64),
			},
		},
	}
	for _, test := range convertLogRecordToJsonTests {
		output := ConvertLogRecordToJson(test.log, test.resource, logger)
		require.Equal(t, output, test.expected)
	}
}
