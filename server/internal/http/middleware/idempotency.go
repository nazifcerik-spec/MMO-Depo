package middleware

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/idempotency"
)

// IdempotencyKeyHeader is the client-provided key per DD-188.
const IdempotencyKeyHeader = "Idempotency-Key"

// MaxBodyBytes caps request bodies read by the idempotency middleware. Anything
// larger is rejected before touching the store or the handler.
const MaxBodyBytes = 256 << 10 // 256 KB

// idempotencyKeyPattern accepts 16-128 characters of ASCII alphanumeric + _ - .
// This admits UUIDs (32 hex) and ULIDs (26 base32) as well as client-picked keys.
var idempotencyKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_.\-]{16,128}$`)

// Idempotency enforces DD-188 on mutating methods (POST/PUT/PATCH/DELETE).
// Non-mutating methods pass through unchanged.
//
// Flow:
//  1. Extract and validate Idempotency-Key header (required, format).
//  2. Read+hash request body (SHA-256), rebuffer for the handler.
//  3. Look up (key, route pattern, account_id) in the store.
//     - HIT + same hash → replay cached status/body; header Idempotency-Replayed: true.
//     - HIT + different hash → 409 E_IDEMPOTENT_MISMATCH (security event candidate).
//     - MISS → execute handler with a recording ResponseWriter; on 2xx, INSERT record.
//
// Anonymous requests use uuid.Nil as the account scope. Once auth middleware is
// wired ahead of this in the chain, httpctx.AccountID(ctx) returns the caller.
func Idempotency(store idempotency.Store, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isMutating(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get(IdempotencyKeyHeader)
			if key == "" {
				apierr.Write(w, r, http.StatusBadRequest,
					"E_IDEMPOTENCY_KEY_REQUIRED",
					"Idempotency-Key header is required for mutating requests")
				return
			}
			if !idempotencyKeyPattern.MatchString(key) {
				apierr.Write(w, r, http.StatusBadRequest,
					"E_IDEMPOTENCY_KEY_INVALID",
					"Idempotency-Key must be 16-128 characters, [A-Za-z0-9_.-]")
				return
			}

			// Read body, cap at MaxBodyBytes+1 so we can detect overflow.
			body, err := io.ReadAll(io.LimitReader(r.Body, MaxBodyBytes+1))
			if err != nil {
				apierr.Write(w, r, http.StatusBadRequest,
					"E_BODY_READ", "failed to read request body")
				return
			}
			if len(body) > MaxBodyBytes {
				apierr.Write(w, r, http.StatusRequestEntityTooLarge,
					"E_BODY_TOO_LARGE", "request body exceeds 256KB")
				return
			}
			// Rebuffer so downstream handlers can Read the body normally.
			r.Body = io.NopCloser(bytes.NewReader(body))

			sum := sha256.Sum256(body)
			hash := sum[:]

			// Route pattern (e.g. /v1/vendor/{code}/buy) — resolved by chi before
			// middleware runs. Falls back to raw path if no chi context.
			endpoint := r.Method + " " + routePattern(r)
			accountID := httpctx.AccountID(r.Context())

			// Lookup
			existing, err := store.Get(r.Context(), key, endpoint, accountID)
			if err != nil && !errors.Is(err, idempotency.ErrNotFound) {
				logger.ErrorContext(r.Context(), "idempotency lookup",
					slog.String("correlation_id", chimw.GetReqID(r.Context())),
					slog.String("err", err.Error()))
				apierr.Write(w, r, http.StatusInternalServerError,
					"E_INTERNAL", "idempotency store unavailable")
				return
			}

			if existing != nil {
				if !bytes.Equal(existing.RequestHash, hash) {
					logger.WarnContext(r.Context(), "idempotency mismatch",
						slog.String("correlation_id", chimw.GetReqID(r.Context())),
						slog.String("endpoint", endpoint),
						slog.String("key_prefix", safeKeyPrefix(key)))
					apierr.Write(w, r, http.StatusConflict,
						"E_IDEMPOTENT_MISMATCH",
						"same Idempotency-Key was used with a different request body")
					return
				}
				// Replay cached response.
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("Idempotency-Replayed", "true")
				w.WriteHeader(int(existing.ResponseStatus))
				_, _ = w.Write(existing.ResponseBody)
				return
			}

			// MISS — execute handler with a recorder.
			rec := newResponseRecorder(w)
			next.ServeHTTP(rec, r)

			// Cache only successful responses. Retries after a 4xx/5xx must be
			// allowed to hit the real handler again.
			if rec.status >= 200 && rec.status < 300 {
				if err := store.Set(r.Context(), idempotency.Record{
					Key:            key,
					Endpoint:       endpoint,
					AccountID:      accountID,
					RequestHash:    hash,
					ResponseStatus: int16(rec.status),
					ResponseBody:   rec.body.Bytes(),
				}); err != nil {
					// Cache write failure is not fatal — response already sent.
					logger.WarnContext(r.Context(), "idempotency store set",
						slog.String("correlation_id", chimw.GetReqID(r.Context())),
						slog.String("err", err.Error()))
				}
			}
		})
	}
}

func isMutating(m string) bool {
	switch m {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// routePattern returns the chi-resolved pattern; falls back to URL path if no
// chi routing context (unusual, e.g. tests calling middleware directly).
func routePattern(r *http.Request) string {
	if rc := chi.RouteContext(r.Context()); rc != nil {
		if p := rc.RoutePattern(); p != "" {
			return p
		}
	}
	return r.URL.Path
}

// safeKeyPrefix returns the first 8 chars for log context without exposing the
// full key (defense in depth; keys are often UUIDs but treated as sensitive).
func safeKeyPrefix(k string) string {
	if len(k) > 8 {
		return k[:8] + "…"
	}
	return k
}

// --- response recorder ------------------------------------------------------

// responseRecorder is a passthrough http.ResponseWriter that captures the
// status code and body for caching. Body is buffered and forwarded — the client
// sees the response as if the recorder were not there.
type responseRecorder struct {
	http.ResponseWriter
	status  int
	body    bytes.Buffer
	written bool
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (r *responseRecorder) WriteHeader(status int) {
	if r.written {
		return
	}
	r.status = status
	r.written = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if !r.written {
		// Implicit 200 on first Write without WriteHeader.
		r.written = true
	}
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
