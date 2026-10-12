// Package middleware carries the request-scoped HTTP helpers: structured access log,
// panic recovery, and (Faz 01 sprint 2) idempotency + rate limit + auth.
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/otter/server/internal/httpctx"
)

// Logger emits one JSON record per response with method, path, status, latency, and
// the chi request id (correlation_id in the audit trail — DD-205).
func Logger(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			base.LogAttrs(r.Context(), slog.LevelInfo, "http",
				slog.String("correlation_id", chimw.GetReqID(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("remote", r.RemoteAddr),
				// The IP rate limits key on (RealIP): shows whether the proxy
				// chain is trusted correctly behind a host's load balancer.
				slog.String("client_ip", httpctx.ClientIP(r)),
			)
		})
	}
}
