package objects

import (
	"encoding/json"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"io/ioutil"
	"testing"
)

func TestTransformToLogzioSpanBytes(tester *testing.T) {
	inStr, err := ioutil.ReadFile("./testresources/span.json")
	if err != nil {
		panic(fmt.Sprintf("error opening sample span file %s", err.Error()))
	}

	var span model.Span
	json.Unmarshal(inStr, &span)
	logzioSpan, err := TransformToLogzioSpanBytes(&span)
	m := make(map[string]interface{})
	err = json.Unmarshal(logzioSpan, &m)
	if _, ok := m["JaegerTag"]; !ok {
		tester.Error("error converting span to logzioSpan, JaegerTag is not found")
	}
}
