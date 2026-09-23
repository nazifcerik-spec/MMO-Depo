//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/mail"
)

func newAccount(t *testing.T, ctx context.Context, repo *auth.PgRepo) auth.Account {
	t.Helper()
	suffix := strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	now := time.Now().UTC()
	a := auth.Account{
		ID: uuid.New(), Email: "it-" + suffix + "@example.com", EmailNormalized: "it-" + suffix + "@example.com",
		Username: "it_" + suffix, UsernameNormalized: "it_" + suffix, PasswordHash: "x",
		State: auth.AccountActive, CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	require.NoError(t, repo.CreateAccount(ctx, a))
	return a
}

type recordingSender struct{ sent []string }

func (r *recordingSender) Send(_ context.Context, to, _, _ string) error {
	r.sent = append(r.sent, to)
	return nil
}

// DD-219 SQL: token issue + outbox in one tx, consume burns siblings, the
// outbox body is cleared after delivery.
func TestAccountRecoverySQL(t *testing.T) {
	p, ctx := pool(t)
	repo := auth.NewPgRepo(p)
	a := newAccount(t, ctx, repo)
	now := time.Now().UTC()

	issue := func(plain string) {
		tok := auth.AccountToken{
			ID: uuid.New(), AccountID: a.ID, Kind: auth.TokenPasswordReset,
			TokenHash: []byte(plain + strings.Repeat("0", 32-len(plain))), ExpiresAt: now.Add(time.Hour), CreatedAt: now,
		}
		require.NoError(t, repo.IssueAccountToken(ctx, tok, auth.OutboundEmail{To: a.Email, Subject: "s", Body: "b"}))
	}
	issue("first")
	issue("second")
	n, err := repo.CountAccountTokens(ctx, a.ID, auth.TokenPasswordReset, now.Add(-time.Minute))
	require.NoError(t, err)
	require.Equal(t, 2, n)

	got, err := repo.ConsumeAccountToken(ctx, auth.TokenPasswordReset, []byte("second"+strings.Repeat("0", 26)), now)
	require.NoError(t, err)
	require.Equal(t, a.ID, got)
	_, err = repo.ConsumeAccountToken(ctx, auth.TokenPasswordReset, []byte("first"+strings.Repeat("0", 27)), now)
	require.ErrorIs(t, err, auth.ErrRecoveryTokenInvalid, "sibling token burned")

	require.NoError(t, repo.UpdatePassword(ctx, a.ID, "new-hash", now))
	require.NoError(t, repo.MarkEmailVerified(ctx, a.ID, now))
	require.NoError(t, repo.RevokeAccountSessions(ctx, a.ID, uuid.Nil, "it", now))
	back, err := repo.AccountByID(ctx, a.ID)
	require.NoError(t, err)
	require.Equal(t, "new-hash", back.PasswordHash)
	require.NotNil(t, back.EmailVerifiedAt)

	sender := &recordingSender{}
	_, err = mail.Dispatch(ctx, p, sender, 500)
	require.NoError(t, err)
	require.Contains(t, sender.sent, a.Email)
	var pendingBodies int
	require.NoError(t, p.QueryRow(ctx,
		`SELECT count(*) FROM ops.email_outbox WHERE to_address = $1 AND body IS NOT NULL`, a.Email).Scan(&pendingBodies))
	require.Zero(t, pendingBodies, "delivered mail keeps no body at rest")
}

// DD-218 SQL: HP columns round-trip through the character repository.
func TestCharacterHPSQL(t *testing.T) {
	p, ctx := pool(t)
	a := newAccount(t, ctx, auth.NewPgRepo(p))
	repo := character.NewPgRepo(p)
	now := time.Now().UTC()
	c := character.Character{
		ID: uuid.New(), AccountID: a.ID, Name: "Hp" + a.Username[3:9], NameNormalized: strings.ToLower("hp" + a.Username[3:9]),
		RaceCode: "human", ClassCode: "warrior", Level: 1, State: character.StateActive,
		CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	require.NoError(t, repo.Create(ctx, c))

	fresh, err := repo.Get(ctx, c.ID)
	require.NoError(t, err)
	require.Nil(t, fresh.HPStored, "never damaged = NULL = full")

	require.NoError(t, repo.SetHP(ctx, c.ID, 42, now))
	hurt, err := repo.Get(ctx, c.ID)
	require.NoError(t, err)
	require.Equal(t, int32(42), *hurt.HPStored)
	require.WithinDuration(t, now, *hurt.HPUpdatedAt, time.Second)
}
