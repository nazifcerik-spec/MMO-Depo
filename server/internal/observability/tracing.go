package observability

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/otter/server"

// InitTracing installs the W3C traceparent propagator and, when endpoint is
// set, an OTLP/HTTP exporter with a parent-based ratio sampler. Without an
// endpoint the global provider stays a no-op: inbound trace context is still
// propagated to logs, but no spans are recorded or exported.
//
// DD-214: head sampling defaults to 100%; the 10% baseline / 100% error / 100%
// slow policy of DD-205 is applied by tail sampling in the collector.
func InitTracing(ctx context.Context, service, env, version, endpoint string, ratio float64) (func(context.Context) error, error) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint+"/v1/traces"))
	if err != nil {
		return nil, fmt.Errorf("otlp exporter: %w", err)
	}
	res := resource.NewSchemaless(
		attribute.String("service.name", service),
		attribute.String("service.version", version),
		attribute.String("deployment.environment", env),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

// StartSpan starts a span on the global tracer (worker jobs, domain code).
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(tracerName).Start(ctx, name, trace.WithAttributes(attrs...))
}

// HTTP is the outermost request middleware: it continues (or starts) the W3C
// trace, makes the trace id the correlation id when the client sent none
// (so chi's RequestID and every log line agree), and records RED metrics
// labelled by chi route pattern — never the raw path, to bound cardinality.
func HTTP(m *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := otel.Tracer(tracerName).Start(ctx, r.Method, trace.WithSpanKind(trace.SpanKindServer))
			defer span.End()

			if sc := span.SpanContext(); sc.IsValid() {
				if r.Header.Get(chimw.RequestIDHeader) == "" {
					r.Header.Set(chimw.RequestIDHeader, sc.TraceID().String())
				}
				w.Header().Set("traceparent", "00-"+sc.TraceID().String()+"-"+sc.SpanID().String()+"-"+sc.TraceFlags().String())
			}

			if m != nil {
				m.httpInFlight.Inc()
				defer m.httpInFlight.Dec()
			}
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			r = r.WithContext(ctx)
			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK // handler wrote nothing (e.g. hijacked WS)
			}
			route := "unmatched"
			if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
				route = rc.RoutePattern()
			}
			span.SetName(r.Method + " " + route)
			span.SetAttributes(
				attribute.String("http.request.method", r.Method),
				attribute.String("http.route", route),
				attribute.Int("http.response.status_code", status),
			)
			if id := r.Header.Get(chimw.RequestIDHeader); id != "" {
				span.SetAttributes(attribute.String("correlation_id", id))
			}
			if status >= 500 {
				span.SetStatus(codes.Error, http.StatusText(status))
			}
			if m != nil {
				m.httpRequests.WithLabelValues(r.Method, route, strconv.Itoa(status)).Inc()
				m.httpDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
			}
		})
	}
}
