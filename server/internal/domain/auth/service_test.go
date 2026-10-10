package auth_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

func newTestService(t *testing.T) (*auth.Service, *auth.MemoryRepo, *auth.Signer) {
	t.Helper()
	repo := auth.NewMemoryRepo()
	seed, _ := auth.GenerateSeed()
	signer, err := auth.NewSigner(seed, "otter-test", 15*time.Minute)
	require.NoError(t, err)

	cfg := auth.ServiceConfig{
		Argon2Params:   fastParams, // reuse fast params from password_test.go
		AccessTTL:      15 * time.Minute,
		RefreshTTL:     24 * time.Hour,
		SessionTTL:     24 * time.Hour,
		Now:            time.Now,
		MinPasswordLen: 8,
	}
	svc := auth.NewService(repo, signer, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo, signer
}

func mustRegister(t *testing.T, svc *auth.Service) *auth.LoginResult {
	t.Helper()
	res, err := svc.Register(context.Background(), auth.RegisterRequest{AcceptTerms: true,
		Email:    "grim@example.com",
		Username: "grim",
		Password: "hunter22-strong",
	})
	require.NoError(t, err)
	return res
}

func TestRegister_Success(t *testing.T) {
	svc, _, _ := newTestService(t)
	res := mustRegister(t, svc)

	require.NotEqual(t, "", res.AccessToken)
	require.NotEqual(t, "", res.RefreshToken)
	require.True(t, res.AccessExpiresAt.After(time.Now()))
	require.True(t, res.RefreshExpires.After(res.AccessExpiresAt))
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _, _ := newTestService(t)
	mustRegister(t, svc)

	_, err := svc.Register(context.Background(), auth.RegisterRequest{AcceptTerms: true,
		Email:    "GRIM@example.com", // case variant of taken email
		Username: "grim2",
		Password: "another-password",
	})
	require.ErrorIs(t, err, auth.ErrAccountExists)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	svc, _, _ := newTestService(t)
	mustRegister(t, svc)

	_, err := svc.Register(context.Background(), auth.RegisterRequest{AcceptTerms: true,
		Email:    "another@example.com",
		Username: "GRIM", // case variant
		Password: "another-password",
	})
	require.ErrorIs(t, err, auth.ErrAccountExists)
}

func TestRegister_ValidationErrors(t *testing.T) {
	svc, _, _ := newTestService(t)

	cases := []struct {
		name  string
		req   auth.RegisterRequest
		field string
	}{
		{"empty email", auth.RegisterRequest{AcceptTerms: true, Email: "", Username: "grim", Password: "hunter22-strong"}, "email"},
		{"bad email", auth.RegisterRequest{AcceptTerms: true, Email: "not-an-email", Username: "grim", Password: "hunter22-strong"}, "email"},
		{"empty username", auth.RegisterRequest{AcceptTerms: true, Email: "a@b.co", Username: "", Password: "hunter22-strong"}, "username"},
		{"bad username chars", auth.RegisterRequest{AcceptTerms: true, Email: "a@b.co", Username: "grim!!", Password: "hunter22-strong"}, "username"},
		{"short password", auth.RegisterRequest{AcceptTerms: true, Email: "a@b.co", Username: "grim", Password: "short"}, "password"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), tc.req)
			require.Error(t, err)
			var ve *auth.ValidationError
			require.ErrorAs(t, err, &ve)
			require.Equal(t, tc.field, ve.Field)
		})
	}
}

func TestLogin_SuccessByEmail(t *testing.T) {
	svc, _, _ := newTestService(t)
	mustRegister(t, svc)

	res, err := svc.Login(context.Background(), auth.LoginRequest{
		EmailOrUsername: "grim@example.com",
		Password:        "hunter22-strong",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
}

func TestLogin_SuccessByUsername(t *testing.T) {
	svc, _, _ := newTestService(t)
	mustRegister(t, svc)

	res, err := svc.Login(context.Background(), auth.LoginRequest{
		EmailOrUsername: "GRIM",
		Password:        "hunter22-strong",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
}

func TestLogin_UnknownAccount(t *testing.T) {
	svc, _, _ := newTestService(t)
	_, err := svc.Login(context.Background(), auth.LoginRequest{
		EmailOrUsername: "ghost@example.com",
		Password:        "any-password",
	})
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _ := newTestService(t)
	mustRegister(t, svc)

	_, err := svc.Login(context.Background(), auth.LoginRequest{
		EmailOrUsername: "grim@example.com",
		Password:        "wrong-password",
	})
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)
}

func TestVerify_ReturnsAccountAndSession(t *testing.T) {
	svc, _, _ := newTestService(t)
	res := mustRegister(t, svc)

	acct, sess, err := svc.Verify(context.Background(), res.AccessToken)
	require.NoError(t, err)
	require.Equal(t, res.AccountID, acct)
	require.Equal(t, res.SessionID, sess)
}

func TestRefresh_RotatesTokens(t *testing.T) {
	svc, _, _ := newTestService(t)
	first := mustRegister(t, svc)

	second, err := svc.Refresh(context.Background(), first.RefreshToken)
	require.NoError(t, err)

	require.NotEqual(t, first.RefreshToken, second.RefreshToken, "refresh must rotate")
	require.NotEqual(t, first.AccessToken, second.AccessToken, "access must rotate")
	require.Equal(t, first.SessionID, second.SessionID, "session preserved")
	require.Equal(t, first.AccountID, second.AccountID)
}

func TestRefresh_ReuseDetectedRevokesSession(t *testing.T) {
	// DD-189 core guarantee: presenting an already-rotated refresh token is
	// treated as compromise. Session and all chained tokens revoked.
	svc, _, _ := newTestService(t)
	first := mustRegister(t, svc)

	second, err := svc.Refresh(context.Background(), first.RefreshToken)
	require.NoError(t, err)

	// Attacker replays the first (already used) refresh token.
	_, err = svc.Refresh(context.Background(), first.RefreshToken)
	require.ErrorIs(t, err, auth.ErrRefreshTokenReuse)

	// The legit second token must also fail now — session tree revoked.
	_, err = svc.Refresh(context.Background(), second.RefreshToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked)

	// Access token verification also fails after session revocation.
	_, _, err = svc.Verify(context.Background(), second.AccessToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked)
}

func TestRefresh_UnknownToken(t *testing.T) {
	svc, _, _ := newTestService(t)
	_, err := svc.Refresh(context.Background(), "unknown-token-nonexistent")
	require.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)
}

func TestLogout_RevokesSession(t *testing.T) {
	svc, _, _ := newTestService(t)
	res := mustRegister(t, svc)

	// Verify works pre-logout
	_, _, err := svc.Verify(context.Background(), res.AccessToken)
	require.NoError(t, err)

	require.NoError(t, svc.Logout(context.Background(), res.SessionID))

	// Verify fails post-logout
	_, _, err = svc.Verify(context.Background(), res.AccessToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked)

	// Refresh also fails
	_, err = svc.Refresh(context.Background(), res.RefreshToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked)
}

func TestLogout_Idempotent(t *testing.T) {
	svc, _, _ := newTestService(t)
	res := mustRegister(t, svc)

	require.NoError(t, svc.Logout(context.Background(), res.SessionID))
	require.NoError(t, svc.Logout(context.Background(), res.SessionID), "second logout must be no-op")
}

// Note: locked/banned account rejection is verified in the pg integration
// suite once the admin state-transition endpoint lands (Faz 01 sprint 5+).
// The service branches themselves (ErrAccountLocked/Banned) are covered by
// the Login handler mapping test.
