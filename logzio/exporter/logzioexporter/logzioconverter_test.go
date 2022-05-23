package logzioexporter

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"testing"
)

func newTestTracesWithAttributes() ptrace.Traces {
	td := ptrace.NewTraces()
	for i := 0; i < 10; i++ {
		s := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
		s.SetName(fmt.Sprintf("%s-%d", testOperation, i))
		s.SetTraceID(pcommon.NewTraceID([16]byte{byte(i), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}))
		s.SetSpanID(pcommon.NewSpanID([8]byte{byte(i), 0, 0, 0, 0, 0, 0, 2}))
		for j := 0; j < 5; j++ {
			s.Attributes().Insert(fmt.Sprintf("k%d", j), pcommon.NewValueString(fmt.Sprintf("v%d", j)))
		}
		s.SetKind(ptrace.SpanKindServer)
	}
	return td
}

func TestLogzioConverter(t *testing.T) {
	cfg := Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		TracesToken:      "",
		Region:           "us",
	}
	td := newTestTracesWithAttributes()
	params := componenttest.NewNopExporterCreateSettings()
	exporter, err := createTracesExporter(context.Background(), params, &cfg)
	err = exporter.Start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		return
	}
	require.NoError(t, err)

	ctx := context.Background()
	err = exporter.ConsumeTraces(ctx, td)
	require.NoError(t, err)
	err = exporter.Shutdown(ctx)
	require.NoError(t, err)
}
