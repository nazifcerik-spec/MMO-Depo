package auth_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
)

func TestRegister_RequiresTerms(t *testing.T) {
	svc, repo, _ := clockService(t)
	ctx := context.Background()
	_, err := svc.Register(ctx, auth.RegisterRequest{Email: "t@example.com", Username: "terms", Password: "long-enough-1"})
	var ve *auth.ValidationError
	require.ErrorAs(t, err, &ve)
	require.Equal(t, "accept_terms", ve.Field)

	res, err := svc.Register(ctx, auth.RegisterRequest{AcceptTerms: true, Email: "t@example.com", Username: "terms", Password: "long-enough-1"})
	require.NoError(t, err)
	a, err := repo.AccountByID(ctx, res.AccountID)
	require.NoError(t, err)
	require.Equal(t, auth.DefaultTermsVersion, a.TermsVersion)
	require.NotNil(t, a.TermsAcceptedAt)
}

func TestDeletion_RequestThenCancelBySigningIn(t *testing.T) {
	svc, repo, _ := clockService(t)
	ctx := context.Background()
	res, err := svc.Register(ctx, auth.RegisterRequest{AcceptTerms: true, Email: "d@example.com", Username: "deleter", Password: "long-enough-1"})
	require.NoError(t, err)

	_, err = svc.RequestDeletion(ctx, res.AccountID, "wrong")
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)

	before := len(repo.Outbox)
	at, err := svc.RequestDeletion(ctx, res.AccountID, "long-enough-1")
	require.NoError(t, err)
	require.WithinDuration(t, time.Now().Add(auth.DeletionGrace), at, time.Minute)
	require.Len(t, repo.Outbox, before+1)
	require.Contains(t, repo.Outbox[len(repo.Outbox)-1].Subject, "Hesap silme")
	_, _, err = svc.Verify(ctx, res.AccessToken)
	require.ErrorIs(t, err, auth.ErrSessionRevoked, "every device is signed out")

	again, err := svc.RequestDeletion(ctx, res.AccountID, "long-enough-1")
	require.NoError(t, err)
	require.Equal(t, at, again, "a repeated request keeps the original date")
	require.Len(t, repo.Outbox, before+1, "and sends no second mail")

	in, err := svc.Login(ctx, auth.LoginRequest{EmailOrUsername: "deleter", Password: "long-enough-1"})
	require.NoError(t, err)
	require.True(t, in.DeletionCancelled)
	a, _ := repo.AccountByID(ctx, res.AccountID)
	require.Nil(t, a.DeletionRequestedAt)
	in, err = svc.Login(ctx, auth.LoginRequest{EmailOrUsername: "deleter", Password: "long-enough-1"})
	require.NoError(t, err)
	require.False(t, in.DeletionCancelled)
}

func TestExportLink_OneTimeAndCapped(t *testing.T) {
	svc, _, _ := clockService(t)
	ctx := context.Background()
	res, err := svc.Register(ctx, auth.RegisterRequest{AcceptTerms: true, Email: "e@example.com", Username: "exporter", Password: "long-enough-1"})
	require.NoError(t, err)

	url, exp, err := svc.IssueExportLink(ctx, res.AccountID)
	require.NoError(t, err)
	require.WithinDuration(t, time.Now().Add(auth.ExportLinkTTL), exp, time.Minute)
	require.True(t, strings.HasPrefix(url, "https://play.otter.gg/v1/me/export?token="), url)
	token := url[strings.Index(url, "token=")+len("token="):]
	id, err := svc.ConsumeExportToken(ctx, token)
	require.NoError(t, err)
	require.Equal(t, res.AccountID, id)
	_, err = svc.ConsumeExportToken(ctx, token)
	require.ErrorIs(t, err, auth.ErrRecoveryTokenInvalid, "single use")
	_, err = svc.ConsumeExportToken(ctx, "")
	require.ErrorIs(t, err, auth.ErrRecoveryTokenInvalid)

	for i := 0; i < 4; i++ {
		_, _, err = svc.IssueExportLink(ctx, res.AccountID)
		require.NoError(t, err)
	}
	_, _, err = svc.IssueExportLink(ctx, res.AccountID)
	require.ErrorIs(t, err, auth.ErrExportLimit)
}
