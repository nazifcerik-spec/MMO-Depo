package auth_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

// Unknown-account logins must cost what real ones do: the dummy carries the
// service's own parameters, and real hashes use the defaults.
func TestDummyHashMatchesParams(t *testing.T) {
	svc, _, _ := newTestService(t)
	require.Contains(t, svc.DummyHash(), "$m=8192,t=1,p=1$") // fastParams
	real, err := auth.HashPassword("pw", auth.DefaultArgon2Params)
	require.NoError(t, err)
	require.Contains(t, real, "$m=19456,t=2,p=1$")
	require.ErrorIs(t, auth.VerifyPassword("anything", svc.DummyHash()), auth.ErrInvalidCredentials)
}

// When every argon2id slot is taken, hashing and verifying fail fast with
// ErrBusy, login answers 503 + Retry-After, and an unknown account answers
// exactly like a known one (no enumeration under load).
func TestKDFBusy(t *testing.T) {
	svc, _, _ := newTestService(t)
	mustRegister(t, svc)
	hash, err := auth.HashPassword("pw", fastParams)
	require.NoError(t, err)

	auth.SaturateKDF(t)
	_, err = auth.HashPassword("pw", fastParams)
	require.ErrorIs(t, err, auth.ErrBusy)
	require.ErrorIs(t, auth.VerifyPassword("pw", hash), auth.ErrBusy)

	for _, who := range []string{"grim", "nobody"} {
		_, err := svc.Login(context.Background(), auth.LoginRequest{EmailOrUsername: who, Password: "hunter22-strong"})
		require.ErrorIs(t, err, auth.ErrBusy, who)
	}

	h := auth.NewHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	rec := httptest.NewRecorder()
	h.Login(rec, httptest.NewRequest(http.MethodPost, "/v1/auth/login",
		strings.NewReader(`{"email_or_username":"grim","password":"hunter22-strong"}`)))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Equal(t, "1", rec.Header().Get("Retry-After"))
	require.Contains(t, rec.Body.String(), "E_BUSY")
}
