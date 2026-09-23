// Package log builds the process-wide slog.Logger emitting structured JSON on stdout.
// Every record carries service + env + commit + timestamp; correlation_id is added
// at request scope by the http middleware.
package log

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

// New returns a JSON logger honoring the given level ("debug"|"info"|"warn"|"error").
func New(level, env, commit string) *slog.Logger {
	lv := parseLevel(level)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lv,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// stdlib emits "time"/"level"/"msg"; keep them.
			return a
		},
	})

	return slog.New(traceHandler{handler}).With(
		slog.String("service", "otter"),
		slog.String("env", env),
		slog.String("commit", commit),
	)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// traceHandler stamps trace_id/span_id on every *Context log call made inside
// a traced request, so Grafana can jump from a trace to its logs (DD-205).
type traceHandler struct{ slog.Handler }

func (h traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		r.AddAttrs(slog.String("trace_id", sc.TraceID().String()), slog.String("span_id", sc.SpanID().String()))
	}
	return h.Handler.Handle(ctx, r)
}

func (h traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return traceHandler{h.Handler.WithAttrs(attrs)}
}

func (h traceHandler) WithGroup(name string) slog.Handler {
	return traceHandler{h.Handler.WithGroup(name)}
}
