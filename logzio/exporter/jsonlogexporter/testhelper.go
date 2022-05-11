package jsonlogexporter

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"time"
)

var (
	TestLogTime         = time.Now()
	TestLogTimestamp    = pcommon.NewTimestampFromTime(TestLogTime)
	resourceAttributes1 = map[string]interface{}{"resource-attr": "resource-attr-val-1"}
	resourceAttributes2 = map[string]interface{}{"resource-attr": "resource-attr-val-2"}
	spanEventAttributes = map[string]interface{}{"span-event-attr": "span-event-attr-val"}
	spanLinkAttributes  = map[string]interface{}{"span-link-attr": "span-link-attr-val"}
	spanAttributes      = map[string]interface{}{"span-attr": "span-attr-val"}
)

// Resource Attributes
func initResourceAttributes1(dest pcommon.Map) {
	pcommon.NewMapFromRaw(resourceAttributes1).CopyTo(dest)
}

// Resources
func initResource1(r pcommon.Resource) {
	initResourceAttributes1(r.Attributes())
}

// Logs

func GenerateLogsOneEmptyResourceLogs() plog.Logs {
	ld := plog.NewLogs()
	ld.ResourceLogs().AppendEmpty()
	return ld
}

func GenerateLogsNoLogRecords() plog.Logs {
	ld := GenerateLogsOneEmptyResourceLogs()
	initResource1(ld.ResourceLogs().At(0).Resource())
	return ld
}

func GenerateLogsOneEmptyLogRecord() plog.Logs {
	ld := GenerateLogsNoLogRecords()
	rs0 := ld.ResourceLogs().At(0)
	rs0.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	return ld
}

func GenerateLogsTwoLogRecordsSameResource() plog.Logs {
	ld := GenerateLogsOneEmptyLogRecord()
	logs := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	fillLogOne(logs.At(0))
	fillLogTwo(logs.AppendEmpty())
	return ld
}

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

func GenerateLogsManyLogRecordsSameResource(count int) plog.Logs {
	ld := GenerateLogsOneEmptyLogRecord()
	logs := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	logs.EnsureCapacity(count)
	for i := 0; i < count; i++ {
		var l plog.LogRecord
		if i < logs.Len() {
			l = logs.At(i)
		} else {
			l = logs.AppendEmpty()
		}

		if i%2 == 0 {
			fillLogOne(l)
		} else {
			fillLogTwo(l)
		}
	}
	return ld
}
