package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/ratelimit"
)

// Bucket re-exports for callers that want to construct policies with the
// ratelimit constructors without importing the store package explicitly.
type Bucket = ratelimit.Bucket

// Policy defines the per-endpoint rate limits. A nil dimension = no limit for
// that scope. Endpoint is the label used in cache keys and metrics — usually
// the route pattern with method prefix.
//
// Example:
//
//	Policy{
//	  Endpoint:   "POST /v1/vendor/{code}/buy",
//	  PerIP:      &Bucket{Capacity: 60, RefillPerSec: 1.0},   // 60/min
//	  PerAccount: &Bucket{Capacity: 60, RefillPerSec: 1.0},
//	}
type Policy struct {
	Endpoint     string
	PerIP        *Bucket
	PerAccount   *Bucket
	PerCharacter *Bucket
}

// RateLimit enforces DD-187 token buckets across the 3 configured scopes.
// Any denied dimension → 429 with Retry-After + X-RateLimit-* headers.
//
// Fail policy: on store error the middleware fails OPEN with a warn log
// (docs/security/rate-limit.md — Redis down should not cause total outage).
// Callers that need strict enforcement (payment, security-critical) can wrap
// with their own check.
func RateLimit(store ratelimit.Store, policy Policy, logger *slog.Logger) func(http.Handler) http.Handler {
	if policy.Endpoint == "" {
		policy.Endpoint = "unknown"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, dim := range collectDimensions(r, policy) {
				res, err := store.Consume(r.Context(), dim.key, dim.bucket, 1)
				if err != nil {
					// Fail-open: log + allow. Never take the site down over Redis.
					logger.WarnContext(r.Context(), "ratelimit store failure — fail-open",
						slog.String("correlation_id", chimw.GetReqID(r.Context())),
						slog.String("scope", dim.scope),
						slog.String("err", err.Error()))
					continue
				}
				setRateHeaders(w, dim.bucket, res)
				if !res.Allowed {
					logger.InfoContext(r.Context(), "rate limit hit",
						slog.String("correlation_id", chimw.GetReqID(r.Context())),
						slog.String("endpoint", policy.Endpoint),
						slog.String("scope", dim.scope),
						slog.Duration("retry_after", res.RetryAfter))
					apierr.Write(w, r, http.StatusTooManyRequests,
						"E_RATE_LIMIT",
						fmt.Sprintf("rate limit exceeded (%s scope)", dim.scope))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

type dimension struct {
	scope  string // "ip" | "account" | "character"
	key    string // redis key
	bucket ratelimit.Bucket
}

func collectDimensions(r *http.Request, p Policy) []dimension {
	dims := make([]dimension, 0, 3)

	if p.PerIP != nil {
		ip := httpctx.ClientIP(r)
		dims = append(dims, dimension{
			scope:  "ip",
			key:    fmt.Sprintf("rl:ip:%s:%s", p.Endpoint, ip),
			bucket: *p.PerIP,
		})
	}
	if p.PerAccount != nil {
		if acct := httpctx.AccountID(r.Context()); acct != uuid.Nil {
			dims = append(dims, dimension{
				scope:  "account",
				key:    fmt.Sprintf("rl:acct:%s:%s", p.Endpoint, acct),
				bucket: *p.PerAccount,
			})
		}
	}
	if p.PerCharacter != nil {
		if char := httpctx.CharacterID(r.Context()); char != uuid.Nil {
			dims = append(dims, dimension{
				scope:  "character",
				key:    fmt.Sprintf("rl:char:%s:%s", p.Endpoint, char),
				bucket: *p.PerCharacter,
			})
		}
	}
	return dims
}

func setRateHeaders(w http.ResponseWriter, b ratelimit.Bucket, res ratelimit.Result) {
	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(b.Capacity, 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(max0(res.Remaining), 10))
	if !res.Allowed && res.RetryAfter > 0 {
		// Retry-After per RFC 7231: integer seconds.
		secs := int64(res.RetryAfter.Seconds())
		if secs < 1 {
			secs = 1
		}
		w.Header().Set("Retry-After", strconv.FormatInt(secs, 10))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(res.RetryAfter).Unix(), 10))
	}
}

func max0(n int64) int64 {
	if n < 0 {
		return 0
	}
	return n
}
