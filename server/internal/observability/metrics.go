// Package observability holds the Prometheus metrics and OpenTelemetry
// tracing setup shared by the API and the worker (DD-205, DD-214).
//
// Metrics live on a private registry (no global default) and are served on a
// separate listener so /metrics never goes through the public proxy.
package observability

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics bundles every collector the process exports. A nil *Metrics is
// valid and turns every recording method into a no-op (tests, tools).
type Metrics struct {
	Registry *prometheus.Registry

	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
	httpInFlight prometheus.Gauge

	jobRuns     *prometheus.CounterVec
	jobDuration *prometheus.HistogramVec
	jobRows     *prometheus.CounterVec

	clientErrors *prometheus.CounterVec
}

// NewMetrics registers the process, Go runtime and Otter collectors.
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	m := &Metrics{
		Registry: reg,
		httpRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "otter_http_requests_total",
			Help: "HTTP requests by route pattern and status code.",
		}, []string{"method", "route", "code"}),
		httpDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "otter_endpoint_request_duration_seconds",
			Help:    "HTTP request latency by route pattern.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"method", "route"}),
		httpInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "otter_http_in_flight_requests",
			Help: "HTTP requests currently being served.",
		}),
		jobRuns: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "otter_worker_job_runs_total",
			Help: "Worker job runs by result (ok|error).",
		}, []string{"job", "result"}),
		jobDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "otter_worker_job_duration_seconds",
			Help:    "Worker job run latency.",
			Buckets: []float64{.01, .05, .1, .5, 1, 5, 15, 60},
		}, []string{"job"}),
		jobRows: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "otter_worker_job_rows_total",
			Help: "Rows affected by worker jobs.",
		}, []string{"job"}),
		clientErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "otter_client_errors_total",
			Help: "Uncaught client errors reported by the app (DD-235).",
		}, []string{"platform", "kind"}),
	}
	reg.MustRegister(m.httpRequests, m.httpDuration, m.httpInFlight, m.jobRuns, m.jobDuration, m.jobRows, m.clientErrors)
	return m
}

// ObserveJob records one worker job run.
func (m *Metrics) ObserveJob(job string, rows int64, dur time.Duration, err error) {
	if m == nil {
		return
	}
	result := "ok"
	if err != nil {
		result = "error"
	}
	m.jobRuns.WithLabelValues(job, result).Inc()
	m.jobDuration.WithLabelValues(job).Observe(dur.Seconds())
	if rows > 0 {
		m.jobRows.WithLabelValues(job).Add(float64(rows))
	}
}

// ObserveClientError counts one client error report (labels are allowlisted
// by the caller, so cardinality stays bounded).
func (m *Metrics) ObserveClientError(platform, kind string) {
	if m == nil {
		return
	}
	m.clientErrors.WithLabelValues(platform, kind).Inc()
}

// RegisterPool exports pgxpool stats, read lazily at scrape time.
func (m *Metrics) RegisterPool(pool *pgxpool.Pool) {
	if m == nil || pool == nil {
		return
	}
	gauge := func(name, help string, fn func(*pgxpool.Stat) float64) prometheus.Collector {
		return prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: name, Help: help},
			func() float64 { return fn(pool.Stat()) })
	}
	counter := func(name, help string, fn func(*pgxpool.Stat) float64) prometheus.Collector {
		return prometheus.NewCounterFunc(prometheus.CounterOpts{Name: name, Help: help},
			func() float64 { return fn(pool.Stat()) })
	}
	m.Registry.MustRegister(
		gauge("otter_db_pool_acquired_conns", "Connections currently checked out.",
			func(s *pgxpool.Stat) float64 { return float64(s.AcquiredConns()) }),
		gauge("otter_db_pool_idle_conns", "Idle connections.",
			func(s *pgxpool.Stat) float64 { return float64(s.IdleConns()) }),
		gauge("otter_db_pool_max_conns", "Pool size limit.",
			func(s *pgxpool.Stat) float64 { return float64(s.MaxConns()) }),
		counter("otter_db_pool_empty_acquire_total", "Acquires that had to wait for a connection.",
			func(s *pgxpool.Stat) float64 { return float64(s.EmptyAcquireCount()) }),
		counter("otter_db_pool_acquire_wait_seconds_total", "Total time spent waiting to acquire.",
			func(s *pgxpool.Stat) float64 { return s.AcquireDuration().Seconds() }),
	)
}

// GaugeFunc registers a scrape-time gauge (e.g. WebSocket client count).
func (m *Metrics) GaugeFunc(name, help string, fn func() float64) {
	if m == nil {
		return
	}
	m.Registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: name, Help: help}, fn))
}

// Handler serves the registry in the Prometheus exposition format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{})
}

// Serve runs the private metrics listener until ctx is done. An empty addr
// disables it. Listener errors are logged, never fatal — losing metrics must
// not take the game down.
func (m *Metrics) Serve(ctx context.Context, addr string, logger *slog.Logger) {
	if m == nil || addr == "" {
		return
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.Handler())
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()
	go func() {
		logger.Info("metrics listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("metrics listener failed", "err", err.Error())
		}
	}()
}
