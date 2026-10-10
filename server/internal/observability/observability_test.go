package observability

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const parent = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"

func router(m *Metrics) http.Handler {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	r := chi.NewRouter()
	r.Use(HTTP(m))
	r.Use(chimw.RequestID)
	r.Get("/v1/character/{id}/sheet", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(chimw.GetReqID(r.Context())))
	})
	r.Get("/boom", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(500) })
	return r
}

func TestHTTP_RouteLabelsAreBounded(t *testing.T) {
	m := NewMetrics()
	h := router(m)
	for _, p := range []string{"/v1/character/a/sheet", "/v1/character/b/sheet", "/nope", "/boom"} {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", p, nil))
	}
	require.Equal(t, 2.0, testutil.ToFloat64(m.httpRequests.WithLabelValues("GET", "/v1/character/{id}/sheet", "200")))
	require.Equal(t, 1.0, testutil.ToFloat64(m.httpRequests.WithLabelValues("GET", "unmatched", "404")))
	require.Equal(t, 1.0, testutil.ToFloat64(m.httpRequests.WithLabelValues("GET", "/boom", "500")))
	require.Equal(t, 0.0, testutil.ToFloat64(m.httpInFlight))
	require.Equal(t, 3, testutil.CollectAndCount(m.httpRequests))
}

func TestHTTP_TraceIDBecomesCorrelationID(t *testing.T) {
	h := router(nil)
	req := httptest.NewRequest("GET", "/v1/character/a/sheet", nil)
	req.Header.Set("traceparent", parent)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", rec.Body.String())
	require.True(t, strings.HasPrefix(rec.Header().Get("traceparent"), "00-4bf92f3577b34da6a3ce929d0e0e4736-"))

	// Client-supplied correlation id wins.
	req = httptest.NewRequest("GET", "/v1/character/a/sheet", nil)
	req.Header.Set("traceparent", parent)
	req.Header.Set("X-Request-Id", "client-123")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, "client-123", rec.Body.String())
}

func TestObserveJob(t *testing.T) {
	m := NewMetrics()
	m.ObserveJob("chat_gc", 5, 10*time.Millisecond, nil)
	m.ObserveJob("chat_gc", 0, time.Millisecond, errors.New("x"))
	require.Equal(t, 1.0, testutil.ToFloat64(m.jobRuns.WithLabelValues("chat_gc", "ok")))
	require.Equal(t, 1.0, testutil.ToFloat64(m.jobRuns.WithLabelValues("chat_gc", "error")))
	require.Equal(t, 5.0, testutil.ToFloat64(m.jobRows.WithLabelValues("chat_gc")))

	require.Positive(t, testutil.ToFloat64(m.jobLastOK.WithLabelValues("chat_gc")), "success stamps the time")
	m.JobInterval("chat_gc", time.Hour)
	require.Equal(t, 3600.0, testutil.ToFloat64(m.jobInterval.WithLabelValues("chat_gc")))
	m.ObserveClientError("android", "flutter")
	require.Equal(t, 1.0, testutil.ToFloat64(m.clientErrors.WithLabelValues("android", "flutter")))

	var nilM *Metrics
	nilM.ObserveJob("x", 1, time.Second, nil) // must not panic
	nilM.JobInterval("x", time.Second)
	nilM.ObserveClientError("x", "y")
}

func TestMetricsHandlerExposition(t *testing.T) {
	m := NewMetrics()
	m.GaugeFunc("otter_ws_clients", "test", func() float64 { return 7 })
	rec := httptest.NewRecorder()
	m.Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/metrics", nil))
	require.Contains(t, rec.Body.String(), "otter_ws_clients 7")
	require.Contains(t, rec.Body.String(), "go_goroutines")
}
