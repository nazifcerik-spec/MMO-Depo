package middleware_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
	appmw "github.com/otter/server/internal/http/middleware"
	"github.com/otter/server/internal/httpctx"
)

// fakeVerifier lets tests script the Verify outcome.
type fakeVerifier struct {
	acct uuid.UUID
	sid  uuid.UUID
	err  error
}

func (f fakeVerifier) Verify(_ context.Context, _ string) (uuid.UUID, uuid.UUID, error) {
	return f.acct, f.sid, f.err
}

func newAuthRouter(mw func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.With(mw).Get("/protected", func(w http.ResponseWriter, req *http.Request) {
		acct := httpctx.AccountID(req.Context())
		sid := httpctx.SessionID(req.Context())
		w.Header().Set("X-Account", acct.String())
		w.Header().Set("X-Session", sid.String())
		w.WriteHeader(http.StatusOK)
	})
	return r
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	v := fakeVerifier{}
	router := newAuthRouter(appmw.RequireAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "E_UNAUTHENTICATED")
}

func TestRequireAuth_ValidToken(t *testing.T) {
	acct := uuid.New()
	sid := uuid.New()
	v := fakeVerifier{acct: acct, sid: sid}
	router := newAuthRouter(appmw.RequireAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, acct.String(), rr.Header().Get("X-Account"))
	require.Equal(t, sid.String(), rr.Header().Get("X-Session"))
}

func TestRequireAuth_TokenExpired(t *testing.T) {
	v := fakeVerifier{err: auth.ErrTokenExpired}
	router := newAuthRouter(appmw.RequireAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "E_TOKEN_EXPIRED")
}

func TestRequireAuth_SessionRevoked(t *testing.T) {
	v := fakeVerifier{err: auth.ErrSessionRevoked}
	router := newAuthRouter(appmw.RequireAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-but-revoked")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "E_SESSION_REVOKED")
}

func TestRequireAuth_GenericError(t *testing.T) {
	v := fakeVerifier{err: errors.New("db down")}
	router := newAuthRouter(appmw.RequireAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer any")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "E_UNAUTHENTICATED")
}

func TestOptionalAuth_NoHeaderPassesAnonymous(t *testing.T) {
	v := fakeVerifier{}
	router := newAuthRouter(appmw.OptionalAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, uuid.Nil.String(), rr.Header().Get("X-Account"))
}

func TestOptionalAuth_InvalidTokenPassesAnonymous(t *testing.T) {
	v := fakeVerifier{err: auth.ErrTokenInvalid}
	router := newAuthRouter(appmw.OptionalAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Optional auth swallows the error → anonymous request.
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, uuid.Nil.String(), rr.Header().Get("X-Account"))
}

func TestOptionalAuth_ValidTokenSetsContext(t *testing.T) {
	acct := uuid.New()
	sid := uuid.New()
	v := fakeVerifier{acct: acct, sid: sid}
	router := newAuthRouter(appmw.OptionalAuth(v, testLogger()))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer good")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, acct.String(), rr.Header().Get("X-Account"))
}

func TestBearerToken_CaseInsensitiveScheme(t *testing.T) {
	acct := uuid.New()
	sid := uuid.New()
	v := fakeVerifier{acct: acct, sid: sid}
	router := newAuthRouter(appmw.RequireAuth(v, testLogger()))

	for _, scheme := range []string{"Bearer", "bearer", "BEARER", "BeArEr"} {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", scheme+" tok")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code, "scheme %q", scheme)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
