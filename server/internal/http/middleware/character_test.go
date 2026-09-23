package middleware_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
	appmw "github.com/otter/server/internal/http/middleware"
	"github.com/otter/server/internal/httpctx"
)

type fakeOwner struct {
	err error
}

func (f fakeOwner) Owns(_ context.Context, _, _ uuid.UUID) error { return f.err }

func newCharRouter(mw func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(mw).Get("/scoped", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Character", httpctx.CharacterID(req.Context()).String())
		w.WriteHeader(http.StatusOK)
	})
	return r
}

func withAccount(r *http.Request, id uuid.UUID) *http.Request {
	return r.WithContext(httpctx.WithAccountID(r.Context(), id))
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRequireCharacter_MissingHeader(t *testing.T) {
	router := newCharRouter(appmw.RequireCharacter(fakeOwner{}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "E_CHARACTER_REQUIRED")
}

func TestRequireCharacter_UnauthenticatedRejected(t *testing.T) {
	router := newCharRouter(appmw.RequireCharacter(fakeOwner{}, discardLogger()))
	req := httptest.NewRequest(http.MethodGet, "/scoped", nil)
	req.Header.Set(appmw.CharacterHeader, uuid.NewString())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "E_UNAUTHENTICATED")
}

func TestRequireCharacter_InvalidUUIDRejected(t *testing.T) {
	// Invalid header is treated as "not provided" → E_CHARACTER_REQUIRED.
	router := newCharRouter(appmw.RequireCharacter(fakeOwner{}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())
	req.Header.Set(appmw.CharacterHeader, "not-a-uuid")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRequireCharacter_OwnershipFailReturns404(t *testing.T) {
	router := newCharRouter(appmw.RequireCharacter(fakeOwner{err: character.ErrForbidden}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())
	req.Header.Set(appmw.CharacterHeader, uuid.NewString())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Existence not leaked — 404 not 403.
	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "E_NOT_FOUND")
}

func TestRequireCharacter_Success(t *testing.T) {
	charID := uuid.New()
	router := newCharRouter(appmw.RequireCharacter(fakeOwner{}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())
	req.Header.Set(appmw.CharacterHeader, charID.String())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, charID.String(), rr.Header().Get("X-Character"))
}

func TestOptionalCharacter_MissingHeaderPassesAnonymous(t *testing.T) {
	router := newCharRouter(appmw.OptionalCharacter(fakeOwner{}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, uuid.Nil.String(), rr.Header().Get("X-Character"))
}

func TestOptionalCharacter_ForeignOwnedPassesAnonymous(t *testing.T) {
	router := newCharRouter(appmw.OptionalCharacter(fakeOwner{err: character.ErrForbidden}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())
	req.Header.Set(appmw.CharacterHeader, uuid.NewString())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Silent pass-through — no error surfaced to caller.
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, uuid.Nil.String(), rr.Header().Get("X-Character"))
}

func TestOptionalCharacter_ValidHeaderSetsContext(t *testing.T) {
	charID := uuid.New()
	router := newCharRouter(appmw.OptionalCharacter(fakeOwner{}, discardLogger()))
	req := withAccount(httptest.NewRequest(http.MethodGet, "/scoped", nil), uuid.New())
	req.Header.Set(appmw.CharacterHeader, charID.String())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, charID.String(), rr.Header().Get("X-Character"))
}

// Path-scoped ownership guard: a foreign or unknown {id} must look exactly
// like a missing character (404) and never reach the handler (IDOR guard).
func TestRequireCharacterParam(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	acct, char := uuid.New(), uuid.New()
	router := func(owner appmw.CharacterOwner) (http.Handler, *int) {
		hits := 0
		r := chi.NewRouter()
		r.Route("/character/{id}", func(r chi.Router) {
			r.Use(appmw.RequireCharacterParam(owner, "id", logger))
			r.Delete("/inventory/{item}", func(w http.ResponseWriter, req *http.Request) {
				hits++
				w.Header().Set("X-Character", httpctx.CharacterID(req.Context()).String())
			})
		})
		return r, &hits
	}
	do := func(h http.Handler, path string, withAcct bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodDelete, path, nil)
		if withAcct {
			req = withAccount(req, acct)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}

	own, hits := router(fakeOwner{})
	rec := do(own, "/character/"+char.String()+"/inventory/x", true)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, char.String(), rec.Header().Get("X-Character"))
	require.Equal(t, 1, *hits)

	for name, err := range map[string]error{"foreign": character.ErrForbidden, "unknown": character.ErrNotFound} {
		h, hits := router(fakeOwner{err: err})
		rec := do(h, "/character/"+char.String()+"/inventory/x", true)
		require.Equal(t, http.StatusNotFound, rec.Code, name)
		require.Zero(t, *hits, name)
	}

	rec = do(own, "/character/not-a-uuid/inventory/x", true)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	rec = do(own, "/character/"+char.String()+"/inventory/x", false)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
