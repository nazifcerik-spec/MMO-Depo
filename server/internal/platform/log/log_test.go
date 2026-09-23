package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestTraceHandler_AddsTraceIDs(t *testing.T) {
	var buf bytes.Buffer
	l := slog.New(traceHandler{slog.NewJSONHandler(&buf, nil)}).With("service", "otter")

	tid, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	sid, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	ctx := trace.ContextWithSpanContext(context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid, SpanID: sid}))

	l.InfoContext(ctx, "traced")
	var rec map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
	require.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", rec["trace_id"])
	require.Equal(t, "00f067aa0ba902b7", rec["span_id"])
	require.Equal(t, "otter", rec["service"])

	buf.Reset()
	l.Info("untraced")
	require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
	require.NotContains(t, buf.String(), "trace_id")
}
