package middleware_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/otter/server/internal/idempotency"
)

const validKey = "a1b2c3d4e5f6a7b8" // 16 chars

// buildRouter wires the middleware onto a POST /echo route that returns the
// body wrapped in {"body": <parsed>, "seq": <counter>}. Seq lets tests assert
// that a replayed response is the cached one, not a re-execution.
func buildRouter(store idempotency.Store) (http.Handler, *atomic.Int32) {
	var counter atomic.Int32
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(appmw.Idempotency(store, logger)).Post("/echo", func(w http.ResponseWriter, req *http.Request) {
		n := counter.Add(1)
		body, _ := io.ReadAll(req.Body)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"seq":%d,"echo":%s}`, n, string(body))
	})
	return r, &counter
}

func doPost(t *testing.T, h http.Handler, key, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	if key != "" {
		req.Header.Set(appmw.IdempotencyKeyHeader, key)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestIdempotency_MissingKeyRejects(t *testing.T) {
	store := idempotency.NewMemoryStore()
	router, _ := buildRouter(store)

	rr := doPost(t, router, "", `{"x":1}`)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "E_IDEMPOTENCY_KEY_REQUIRED")
}

func TestIdempotency_InvalidKeyRejects(t *testing.T) {
	store := idempotency.NewMemoryStore()
	router, _ := buildRouter(store)

	// Too short
	rr := doPost(t, router, "short", `{"x":1}`)
	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "E_IDEMPOTENCY_KEY_INVALID")

	// Illegal char
	rr = doPost(t, router, "bad key with spaces!", `{"x":1}`)
	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "E_IDEMPOTENCY_KEY_INVALID")
}

func TestIdempotency_MissExecutesHandler(t *testing.T) {
	store := idempotency.NewMemoryStore()
	router, counter := buildRouter(store)

	rr := doPost(t, router, validKey, `{"x":1}`)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, int32(1), counter.Load(), "handler must run once on MISS")
	require.NotEqual(t, "true", rr.Header().Get("Idempotency-Replayed"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, float64(1), body["seq"])
}

func TestIdempotency_HitReplaysCachedResponse(t *testing.T) {
	store := idempotency.NewMemoryStore()
	router, counter := buildRouter(store)

	first := doPost(t, router, validKey, `{"x":1}`)
	require.Equal(t, http.StatusOK, first.Code)

	second := doPost(t, router, validKey, `{"x":1}`)

	require.Equal(t, http.StatusOK, second.Code)
	require.Equal(t, int32(1), counter.Load(), "handler must NOT re-run on cached HIT")
	require.Equal(t, "true", second.Header().Get("Idempotency-Replayed"))
	require.Equal(t, first.Body.String(), second.Body.String(), "cached body must match original")
}

func TestIdempotency_MismatchRejects(t *testing.T) {
	store := idempotency.NewMemoryStore()
	router, counter := buildRouter(store)

	// First call establishes cache with body {"x":1}
	first := doPost(t, router, validKey, `{"x":1}`)
	require.Equal(t, http.StatusOK, first.Code)

	// Same key, different body → 409
	second := doPost(t, router, validKey, `{"x":2}`)

	require.Equal(t, http.StatusConflict, second.Code)
	require.Contains(t, second.Body.String(), "E_IDEMPOTENT_MISMATCH")
	require.Equal(t, int32(1), counter.Load(), "handler must NOT run on mismatch")
}

func TestIdempotency_ErrorResponsesNotCached(t *testing.T) {
	store := idempotency.NewMemoryStore()
	var counter atomic.Int32
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(appmw.Idempotency(store, logger)).Post("/fail", func(w http.ResponseWriter, req *http.Request) {
		counter.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"code":"E_INTERNAL"}}`))
	})

	// First request: 500
	req := httptest.NewRequest(http.MethodPost, "/fail", strings.NewReader(`{}`))
	req.Header.Set(appmw.IdempotencyKeyHeader, validKey)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	// Second request with same key should re-execute (5xx is not cacheable).
	req2 := httptest.NewRequest(http.MethodPost, "/fail", strings.NewReader(`{}`))
	req2.Header.Set(appmw.IdempotencyKeyHeader, validKey)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusInternalServerError, rr2.Code)
	require.Equal(t, int32(2), counter.Load(), "5xx responses must not be cached — retry executes")
}

func TestIdempotency_ScopeIsolatedByAccount(t *testing.T) {
	store := idempotency.NewMemoryStore()
	var counter atomic.Int32
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	acctA := uuid.New()
	acctB := uuid.New()
	handler := chi.NewRouter()
	handler.Use(chimw.RequestID)
	handler.With(appmw.Idempotency(store, logger)).Post("/echo", func(w http.ResponseWriter, req *http.Request) {
		counter.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	do := func(acct uuid.UUID) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"a":1}`))
		req.Header.Set(appmw.IdempotencyKeyHeader, validKey)
		req = req.WithContext(httpctx.WithAccountID(req.Context(), acct))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		return rr
	}

	do(acctA)
	do(acctB)

	require.Equal(t, int32(2), counter.Load(), "same key with different accounts must execute separately")
}

func TestIdempotency_NonMutatingMethodsPassthrough(t *testing.T) {
	store := idempotency.NewMemoryStore()
	var counter atomic.Int32
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(appmw.Idempotency(store, logger)).Get("/read", func(w http.ResponseWriter, req *http.Request) {
		counter.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/read", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "GET should pass through without requiring key")
	require.Equal(t, int32(1), counter.Load())
}

func TestIdempotency_BodyIsRebufferedForHandler(t *testing.T) {
	// Middleware reads the body to hash it. Handler must still see the same bytes.
	store := idempotency.NewMemoryStore()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(appmw.Idempotency(store, logger)).Post("/echo", func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})

	payload := `{"complex":{"nested":[1,2,3]}}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(payload))
	req.Header.Set(appmw.IdempotencyKeyHeader, validKey)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, payload, rr.Body.String())
}

func TestIdempotency_BodyTooLargeRejects(t *testing.T) {
	store := idempotency.NewMemoryStore()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(appmw.Idempotency(store, logger)).Post("/echo", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	oversized := bytes.Repeat([]byte("x"), appmw.MaxBodyBytes+10)
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(oversized))
	req.Header.Set(appmw.IdempotencyKeyHeader, validKey)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rr.Code)
	require.Contains(t, rr.Body.String(), "E_BODY_TOO_LARGE")
}
