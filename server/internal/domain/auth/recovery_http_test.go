package auth_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

func TestRecoveryPages(t *testing.T) {
	svc, repo, _ := clockService(t)
	h := auth.NewHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	register(t, svc)
	token := lastToken(t, repo)

	rec := httptest.NewRecorder()
	h.ResetPasswordPage(rec, httptest.NewRequest("GET", "/reset-password?token=x", nil))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "no-referrer", rec.Header().Get("Referrer-Policy"), "token in URL must not leak")
	csp := rec.Header().Get("Content-Security-Policy")
	require.Contains(t, csp, "default-src 'none'")
	nonce := strings.TrimSuffix(strings.SplitN(csp, "'nonce-", 2)[1], "'")
	require.Contains(t, rec.Body.String(), `<script nonce="`+nonce+`">`)

	rec = httptest.NewRecorder()
	h.VerifyEmailPage(rec, httptest.NewRequest("GET", "/verify-email?token=bogus", nil))
	require.Equal(t, http.StatusBadRequest, rec.Code)

	rec = httptest.NewRecorder()
	h.VerifyEmailPage(rec, httptest.NewRequest("GET", "/verify-email?token="+token, nil))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "verified")
	require.NotContains(t, rec.Header().Get("Content-Security-Policy"), "script-src", "result page runs no script")
}
