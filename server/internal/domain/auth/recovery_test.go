package auth_test

import (
	"context"
	"io"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

var tokenRe = regexp.MustCompile(`token=([A-Za-z0-9_-]+)`)

// clockService is newTestService with a controllable clock.
func clockService(t *testing.T) (*auth.Service, *auth.MemoryRepo, *time.Time) {
	t.Helper()
	repo := auth.NewMemoryRepo()
	seed, _ := auth.GenerateSeed()
	signer, err := auth.NewSigner(seed, "otter-test", 15*time.Minute)
	require.NoError(t, err)
	now := time.Now().UTC()
	svc := auth.NewService(repo, signer, auth.ServiceConfig{
		Argon2Params: fastParams, AccessTTL: 15 * time.Minute, RefreshTTL: 24 * time.Hour,
		SessionTTL: 24 * time.Hour, MinPasswordLen: 8, PublicBaseURL: "https://play.otter.gg/",
		Now: func() time.Time { return now },
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo, &now
}

func lastToken(t *testing.T, repo *auth.MemoryRepo) string {
	t.Helper()
	require.NotEmpty(t, repo.Outbox)
	m := tokenRe.FindStringSubmatch(repo.Outbox[len(repo.Outbox)-1].Body)
	require.Len(t, m, 2)
	return m[1]
}

func register(t *testing.T, svc *auth.Service) *auth.LoginResult {
	t.Helper()
	res, err := svc.Register(context.Background(), auth.RegisterRequest{
		Email: "grim@example.com", Username: "grim", Password: "hunter22-strong",
	})
	require.NoError(t, err)
	return res
}

func TestRegister_QueuesVerificationAndVerifyIsSingleUse(t *testing.T) {
	svc, repo, _ := clockService(t)
	ctx := context.Background()
	res := register(t, svc)
	require.False(t, res.EmailVerified)
	require.Len(t, repo.Outbox, 1)
	require.Equal(t, "grim@example.com", repo.Outbox[0].To)
	require.Contains(t, repo.Outbox[0].Body, "https://play.otter.gg/verify-email?token=")

	token := lastToken(t, repo)
	require.NoError(t, svc.VerifyEmail(ctx, token))
	ok, err := svc.EmailVerified(ctx, res.AccountID)
	require.NoError(t, err)
	require.True(t, ok)
	require.ErrorIs(t, svc.VerifyEmail(ctx, token), auth.ErrRecoveryTokenInvalid, "single use")
	require.ErrorIs(t, svc.ResendVerification(ctx, res.AccountID), auth.ErrEmailAlreadyVerified)

	login, err := svc.Login(ctx, auth.LoginRequest{EmailOrUsername: "grim", Password: "hunter22-strong"})
	require.NoError(t, err)
	require.True(t, login.EmailVerified)
}

func TestPasswordReset_FullFlow(t *testing.T) {
	svc, repo, _ := clockService(t)
	ctx := context.Background()
	res := register(t, svc)
	mails := len(repo.Outbox)

	require.NoError(t, svc.RequestPasswordReset(ctx, "nobody@example.com"))
	require.Len(t, repo.Outbox, mails, "unknown e-mail: silent, nothing sent")

	require.NoError(t, svc.RequestPasswordReset(ctx, "  GRIM@example.com "))
	first := lastToken(t, repo)
	require.Contains(t, repo.Outbox[len(repo.Outbox)-1].Body, "/reset-password?token=")
	require.NoError(t, svc.RequestPasswordReset(ctx, "grim@example.com"))
	second := lastToken(t, repo)

	var ve *auth.ValidationError
	require.ErrorAs(t, svc.ConfirmPasswordReset(ctx, second, "short"), &ve)
	require.NoError(t, svc.ConfirmPasswordReset(ctx, second, "brand-new-pass"), "weak attempt did not burn the token")
	require.ErrorIs(t, svc.ConfirmPasswordReset(ctx, second, "another-pass-1"), auth.ErrRecoveryTokenInvalid)
	require.ErrorIs(t, svc.ConfirmPasswordReset(ctx, first, "another-pass-1"), auth.ErrRecoveryTokenInvalid,
		"older links die when one is used")

	_, err := svc.Login(ctx, auth.LoginRequest{EmailOrUsername: "grim", Password: "hunter22-strong"})
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	_, err = svc.Login(ctx, auth.LoginRequest{EmailOrUsername: "grim", Password: "brand-new-pass"})
	require.NoError(t, err)

	_, err = svc.Refresh(ctx, res.RefreshToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked, "reset signs out every existing session")
}

func TestRecoveryTokens_ExpireAndAreCapped(t *testing.T) {
	svc, repo, now := clockService(t)
	ctx := context.Background()
	register(t, svc)
	base := len(repo.Outbox)

	for range 5 {
		require.NoError(t, svc.RequestPasswordReset(ctx, "grim@example.com"))
	}
	require.Len(t, repo.Outbox, base+3, "3 reset mails per hour")

	token := lastToken(t, repo)
	*now = now.Add(61 * time.Minute)
	require.ErrorIs(t, svc.ConfirmPasswordReset(ctx, token, "brand-new-pass"), auth.ErrRecoveryTokenInvalid, "1h TTL")

	require.NoError(t, svc.RequestPasswordReset(ctx, "grim@example.com"))
	require.Len(t, repo.Outbox, base+4, "cap window slid")
}

func TestChangePassword_KeepsCurrentSessionOnly(t *testing.T) {
	svc, _, _ := clockService(t)
	ctx := context.Background()
	current := register(t, svc)
	other, err := svc.Login(ctx, auth.LoginRequest{EmailOrUsername: "grim", Password: "hunter22-strong"})
	require.NoError(t, err)

	require.ErrorIs(t, svc.ChangePassword(ctx, current.AccountID, current.SessionID, "wrong-password", "brand-new-pass"),
		auth.ErrInvalidCredentials)
	require.NoError(t, svc.ChangePassword(ctx, current.AccountID, current.SessionID, "hunter22-strong", "brand-new-pass"))

	_, err = svc.Refresh(ctx, other.RefreshToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked)
	_, err = svc.Refresh(ctx, current.RefreshToken)
	require.NoError(t, err, "the device that changed the password stays signed in")
}
