package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Verifier is the auth capability the middleware needs — one method, implemented
// by *auth.Service. Kept as an interface so tests can swap in a fake.
type Verifier interface {
	Verify(ctx context.Context, accessToken string) (accountID, sessionID uuid.UUID, err error)
}

// RequireAuth rejects any request without a valid Bearer token. On success it
// stores account_id + session_id in the request context and forwards.
func RequireAuth(verifier Verifier, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r)
			if !ok {
				apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "missing bearer token")
				return
			}
			acct, sid, err := verifier.Verify(r.Context(), token)
			if err != nil {
				respondAuthError(w, r, logger, err)
				return
			}
			ctx := httpctx.WithAccountID(r.Context(), acct)
			ctx = httpctx.WithSessionID(ctx, sid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth tries to resolve a Bearer token if present, populating the
// context on success. Missing or invalid tokens do not fail the request —
// downstream handlers see uuid.Nil in that case. Useful for routes that behave
// differently for anonymous vs authenticated clients (e.g. rate-limit account
// scope enrichment on public endpoints).
func OptionalAuth(verifier Verifier, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			acct, sid, err := verifier.Verify(r.Context(), token)
			if err != nil {
				// Silent pass-through — the request continues as anonymous.
				logger.DebugContext(r.Context(), "optional auth failed — treating as anonymous",
					slog.String("correlation_id", chimw.GetReqID(r.Context())),
					slog.String("err", err.Error()))
				next.ServeHTTP(w, r)
				return
			}
			ctx := httpctx.WithAccountID(r.Context(), acct)
			ctx = httpctx.WithSessionID(ctx, sid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// bearerToken extracts the token from `Authorization: Bearer <token>`.
// Case-insensitive on the scheme; strict on the space separator.
func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	const prefix = "bearer "
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return "", false
	}
	tok := strings.TrimSpace(h[len(prefix):])
	if tok == "" {
		return "", false
	}
	return tok, true
}

func respondAuthError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	switch {
	case errors.Is(err, auth.ErrTokenExpired):
		apierr.Write(w, r, http.StatusUnauthorized, "E_TOKEN_EXPIRED", "token expired")
	case errors.Is(err, auth.ErrTokenInvalid):
		apierr.Write(w, r, http.StatusUnauthorized, "E_TOKEN_INVALID", "token invalid")
	case errors.Is(err, auth.ErrSessionRevoked):
		apierr.Write(w, r, http.StatusUnauthorized, "E_SESSION_REVOKED", "session revoked")
	case errors.Is(err, auth.ErrSessionExpired):
		apierr.Write(w, r, http.StatusUnauthorized, "E_SESSION_EXPIRED", "session expired")
	default:
		logger.WarnContext(r.Context(), "auth verify failed",
			slog.String("correlation_id", chimw.GetReqID(r.Context())),
			slog.String("err", err.Error()))
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "authentication failed")
	}
}
