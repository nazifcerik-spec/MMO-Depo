package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// CharacterHeader is the request header carrying the active character id.
// Set by the client after /character/{id}/enter (or on any character-scoped
// request). Missing header → uuid.Nil in context (dimension skipped by rate
// limit and other character-scope middleware).
const CharacterHeader = "X-Character-Id"

// CharacterOwner is the ownership-check capability the middleware needs.
// Implemented by *character.Service.
type CharacterOwner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

// RequireCharacter parses X-Character-Id, verifies the caller's account owns
// it, and puts the id in ctx. Missing header or foreign-owned id → 400/404.
// Meant for routes that MUST have a character (activity, combat, inventory).
func RequireCharacter(owner CharacterOwner, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := parseCharacterHeader(r)
			if !ok {
				apierr.Write(w, r, http.StatusBadRequest,
					"E_CHARACTER_REQUIRED",
					"X-Character-Id header required")
				return
			}
			acct := httpctx.AccountID(r.Context())
			if acct == uuid.Nil {
				apierr.Write(w, r, http.StatusUnauthorized,
					"E_UNAUTHENTICATED",
					"auth required for character-scoped route")
				return
			}
			if err := owner.Owns(r.Context(), id, acct); err != nil {
				handleOwnershipErr(w, r, logger, err)
				return
			}
			ctx := httpctx.WithCharacterID(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireCharacterParam guards every /character/{param}/... route: the path id
// must belong to the caller's account, otherwise 404 (never 403 — existence is
// not leaked). Handlers below it may trust the id. Mounted once on the
// character route block so no handler can forget its own check.
func RequireCharacterParam(owner CharacterOwner, param string, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.Parse(chi.URLParam(r, param))
			if err != nil {
				apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", param+": invalid uuid")
				return
			}
			acct := httpctx.AccountID(r.Context())
			if acct == uuid.Nil {
				apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
				return
			}
			if err := owner.Owns(r.Context(), id, acct); err != nil {
				handleOwnershipErr(w, r, logger, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(httpctx.WithCharacterID(r.Context(), id)))
		})
	}
}

// OptionalCharacter parses X-Character-Id if present, verifies ownership, and
// populates ctx. Missing header or ownership failure → silently anonymous
// character scope (uuid.Nil). Useful for routes that behave differently for
// character-bound clients (e.g. rate-limit character scope enrichment).
func OptionalCharacter(owner CharacterOwner, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := parseCharacterHeader(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			acct := httpctx.AccountID(r.Context())
			if acct == uuid.Nil {
				next.ServeHTTP(w, r)
				return
			}
			if err := owner.Owns(r.Context(), id, acct); err != nil {
				logger.DebugContext(r.Context(), "optional character context — ownership check failed",
					slog.String("correlation_id", chimw.GetReqID(r.Context())),
					slog.String("err", err.Error()))
				next.ServeHTTP(w, r)
				return
			}
			ctx := httpctx.WithCharacterID(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseCharacterHeader(r *http.Request) (uuid.UUID, bool) {
	raw := r.Header.Get(CharacterHeader)
	if raw == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func handleOwnershipErr(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	switch {
	case errors.Is(err, character.ErrNotFound), errors.Is(err, character.ErrForbidden):
		// Do not distinguish — leaking existence enables enumeration attacks.
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
	default:
		logger.WarnContext(r.Context(), "character ownership check failed",
			slog.String("correlation_id", chimw.GetReqID(r.Context())),
			slog.String("err", err.Error()))
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", "internal error")
	}
}
