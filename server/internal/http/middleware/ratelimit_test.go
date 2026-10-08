package middleware_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	appmw "github.com/otter/server/internal/http/middleware"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/ratelimit"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func buildRateRouter(store ratelimit.Store, policy appmw.Policy) (http.Handler, *atomic.Int32) {
	var counter atomic.Int32
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(appmw.RateLimit(store, policy, newTestLogger())).Get("/ping", func(w http.ResponseWriter, req *http.Request) {
		counter.Add(1)
		w.WriteHeader(http.StatusOK)
	})
	return r, &counter
}

func doPing(t *testing.T, h http.Handler, ip string, acct uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = ip + ":40000"
	if acct != uuid.Nil {
		req = req.WithContext(httpctx.WithAccountID(req.Context(), acct))
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestRateLimit_AllowsWithinIPCapacity(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	b := ratelimit.PerMinute(3)
	router, counter := buildRateRouter(store, appmw.Policy{
		Endpoint: "GET /ping",
		PerIP:    &b,
	})

	for i := 0; i < 3; i++ {
		rr := doPing(t, router, "10.0.0.1", uuid.Nil)
		require.Equal(t, http.StatusOK, rr.Code, "call %d", i)
	}
	require.Equal(t, int32(3), counter.Load())
}

func TestRateLimit_RejectsWhenIPExhausted(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	b := ratelimit.PerMinute(2)
	router, counter := buildRateRouter(store, appmw.Policy{
		Endpoint: "GET /ping",
		PerIP:    &b,
	})

	// Drain
	for i := 0; i < 2; i++ {
		rr := doPing(t, router, "10.0.0.1", uuid.Nil)
		require.Equal(t, http.StatusOK, rr.Code)
	}

	// Next → 429
	rr := doPing(t, router, "10.0.0.1", uuid.Nil)
	require.Equal(t, http.StatusTooManyRequests, rr.Code)
	require.Equal(t, int32(2), counter.Load(), "handler must NOT execute on 429")
	require.Contains(t, rr.Body.String(), "E_RATE_LIMIT")

	// Response headers
	require.Equal(t, "2", rr.Header().Get("X-RateLimit-Limit"))
	require.NotEmpty(t, rr.Header().Get("Retry-After"))
	require.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
}

func TestRateLimit_DifferentIPsIndependent(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	b := ratelimit.PerMinute(1)
	router, _ := buildRateRouter(store, appmw.Policy{
		Endpoint: "GET /ping",
		PerIP:    &b,
	})

	// Exhaust IP A
	require.Equal(t, http.StatusOK, doPing(t, router, "10.0.0.1", uuid.Nil).Code)
	require.Equal(t, http.StatusTooManyRequests, doPing(t, router, "10.0.0.1", uuid.Nil).Code)

	// IP B fresh
	require.Equal(t, http.StatusOK, doPing(t, router, "10.0.0.2", uuid.Nil).Code)
}

func TestRateLimit_AccountScopeCatchesShared(t *testing.T) {
	// Same account across different IPs (NAT scenario) → account bucket enforces.
	store := ratelimit.NewMemoryStore()
	ipB := ratelimit.PerMinute(100) // irrelevant — high
	acctB := ratelimit.PerMinute(2) // strict
	router, _ := buildRateRouter(store, appmw.Policy{
		Endpoint:   "GET /ping",
		PerIP:      &ipB,
		PerAccount: &acctB,
	})

	acct := uuid.New()

	// Both allowed from different IPs but same account
	require.Equal(t, http.StatusOK, doPing(t, router, "10.0.0.1", acct).Code)
	require.Equal(t, http.StatusOK, doPing(t, router, "10.0.0.2", acct).Code)
	// Third from a third IP → 429 (account scope hit)
	rr := doPing(t, router, "10.0.0.3", acct)
	require.Equal(t, http.StatusTooManyRequests, rr.Code)
}

func TestRateLimit_AnonymousSkipsAccountScope(t *testing.T) {
	// PerAccount configured but no account in context → dimension skipped.
	store := ratelimit.NewMemoryStore()
	strictAcct := ratelimit.PerMinute(1)
	looseIP := ratelimit.PerMinute(10)
	router, counter := buildRateRouter(store, appmw.Policy{
		Endpoint:   "GET /ping",
		PerIP:      &looseIP,
		PerAccount: &strictAcct,
	})

	// 5 anonymous requests — account limit ignored
	for i := 0; i < 5; i++ {
		rr := doPing(t, router, "10.0.0.1", uuid.Nil)
		require.Equal(t, http.StatusOK, rr.Code, "call %d", i)
	}
	require.Equal(t, int32(5), counter.Load())
}

// failStore always errors on Consume; used to test fail-open behavior.
type failStore struct{}

func (failStore) Consume(_ context.Context, _ string, _ ratelimit.Bucket, _ int64) (ratelimit.Result, error) {
	return ratelimit.Result{}, errors.New("simulated redis outage")
}

func TestRateLimit_FailOpenOnStoreError(t *testing.T) {
	// DD-187 fail-open: Redis outage must not cause total service outage.
	b := ratelimit.PerMinute(1)
	router, counter := buildRateRouter(failStore{}, appmw.Policy{
		Endpoint: "GET /ping",
		PerIP:    &b,
	})

	for i := 0; i < 5; i++ {
		rr := doPing(t, router, "10.0.0.1", uuid.Nil)
		require.Equal(t, http.StatusOK, rr.Code)
	}
	require.Equal(t, int32(5), counter.Load())
}

func TestRateLimit_HeadersOnSuccess(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	b := ratelimit.PerMinute(10)
	router, _ := buildRateRouter(store, appmw.Policy{
		Endpoint: "GET /ping",
		PerIP:    &b,
	})

	rr := doPing(t, router, "10.0.0.1", uuid.Nil)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "10", rr.Header().Get("X-RateLimit-Limit"))
	require.Equal(t, "9", rr.Header().Get("X-RateLimit-Remaining"))
	require.Empty(t, rr.Header().Get("Retry-After"), "no retry-after when allowed")
}

// Guard: verify the assertion helpers compile against strings we care about.
func TestPolicyStringFields(t *testing.T) {
	p := appmw.Policy{Endpoint: "POST /v1/echo"}
	require.True(t, strings.HasPrefix(p.Endpoint, "POST"))
}
