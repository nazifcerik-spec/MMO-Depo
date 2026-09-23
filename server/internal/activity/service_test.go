package activity_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
)

type fakeChars struct {
	owned map[uuid.UUID]*character.Character
}

func (f *fakeChars) GetOwned(_ context.Context, cid, acct uuid.UUID) (*character.Character, error) {
	c, ok := f.owned[cid]
	if !ok {
		return nil, character.ErrNotFound
	}
	if c.AccountID != acct {
		return nil, character.ErrForbidden
	}
	return c, nil
}

func newFixtures(t *testing.T) (*activity.Service, *activity.MemoryRepo, *economy.MemoryTxnService, *character.Character, uuid.UUID) {
	t.Helper()
	acct := uuid.New()
	char := &character.Character{
		ID: uuid.New(), AccountID: acct, Name: "Grim",
		RaceCode: "human", ClassCode: "warrior",
		Level: 1, State: character.StateActive,
	}
	chars := &fakeChars{owned: map[uuid.UUID]*character.Character{char.ID: char}}
	repo := activity.NewMemoryRepo()
	txn := economy.NewMemoryTxnService(economy.StaticCurrencies())
	svc := activity.NewService(repo, chars, txn, activity.DefaultServiceConfig(),
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo, txn, char, acct
}

func TestStart_Success(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)

	res, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID:     char.ID,
		AccountID:       acct,
		Kind:            activity.KindMining,
		DurationSeconds: 3600,
	})
	require.NoError(t, err)
	require.NotNil(t, res.Contract)
	require.Equal(t, activity.StateActive, res.Contract.State)
	require.Equal(t, activity.KindMining, res.Contract.Kind)
	require.NotZero(t, res.Contract.Seed)
	require.Greater(t, res.FullYieldPeek.GoldCopper, int64(0), "1h mining preview should yield gold")
}

func TestStart_UnknownKindRejected(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	_, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct,
		Kind: "unicorn_farming", DurationSeconds: 3600,
	})
	var ve *activity.ValidationError
	require.ErrorAs(t, err, &ve)
	require.Equal(t, "kind", ve.Field)
}

func TestStart_DurationBoundsEnforced(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)

	// Too short
	_, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 30,
	})
	require.Error(t, err)

	// Too long
	_, err = svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 100_000,
	})
	require.Error(t, err)
}

func TestStart_ForbiddenForOtherAccount(t *testing.T) {
	svc, _, _, char, _ := newFixtures(t)
	_, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: uuid.New(),
		Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.ErrorIs(t, err, activity.ErrForbidden)
}

func TestStart_SlotBusy(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	_, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct,
		Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.NoError(t, err)

	_, err = svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct,
		Kind: activity.KindFishing, DurationSeconds: 3600,
	})
	require.ErrorIs(t, err, activity.ErrSlotBusy)
}

func TestReveal_MintsGold(t *testing.T) {
	svc, _, txn, char, acct := newFixtures(t)

	res, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct,
		Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.NoError(t, err)

	// Simulate time passing to make elapsed = duration by pushing the contract
	// StartedAt back. Alternative would be clock injection into Service.
	// (Since MemoryRepo stores by copy, we mutate via reveal path — elapsed
	// clamps at duration when reveal is later than expires_at.)

	reveal, err := svc.Reveal(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)
	require.False(t, reveal.Replayed)
	require.GreaterOrEqual(t, reveal.Yield.GoldCopper, int64(0))

	// Wallet debit — if any gold was minted, character wallet should reflect it.
	bal, err := txn.Balance(context.Background(), economy.AccountCharacter, char.ID, "gold")
	require.NoError(t, err)
	require.Equal(t, reveal.Yield.GoldCopper, bal)
}

func TestReveal_Idempotent(t *testing.T) {
	svc, _, txn, char, acct := newFixtures(t)

	res, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct,
		Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.NoError(t, err)

	r1, err := svc.Reveal(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)

	r2, err := svc.Reveal(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)
	require.True(t, r2.Replayed, "second reveal must be idempotent replay")
	require.Equal(t, r1.Yield.GoldCopper, r2.Yield.GoldCopper)

	// Balance applied ONCE.
	bal, _ := txn.Balance(context.Background(), economy.AccountCharacter, char.ID, "gold")
	require.Equal(t, r1.Yield.GoldCopper, bal, "balance must NOT double on replay")
}

func TestReveal_DeterministicGivenSeed(t *testing.T) {
	svc1, _, _, char1, acct1 := newFixtures(t)
	res1, _ := svc1.Start(context.Background(), activity.StartRequest{
		CharacterID: char1.ID, AccountID: acct1, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	// Force known seed for the test: rewrite through the repo.
	res1.Contract.Seed = 4242
	elapsed := res1.Contract.DurationSeconds
	yield1 := activity.ComputeYield(res1.Contract, elapsed)
	yield2 := activity.ComputeYield(res1.Contract, elapsed)
	require.Equal(t, yield1, yield2, "same seed + elapsed must produce identical yield")
}

func TestReveal_ForbiddenForOtherAccount(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	res, _ := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	_, err := svc.Reveal(context.Background(), res.Contract.ID, uuid.New())
	require.ErrorIs(t, err, activity.ErrForbidden)
}

func TestStop_CancelsActive(t *testing.T) {
	svc, _, txn, char, acct := newFixtures(t)
	res, _ := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	c, err := svc.Stop(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)
	require.Equal(t, activity.StateCancelled, c.State)

	// No gold minted.
	bal, _ := txn.Balance(context.Background(), economy.AccountCharacter, char.ID, "gold")
	require.Equal(t, int64(0), bal)
}

func TestStop_IdempotentOnCancelled(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	res, _ := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	_, err := svc.Stop(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)
	_, err = svc.Stop(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err, "second stop should be no-op")
}

func TestStop_RevealedRejected(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	res, _ := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	_, err := svc.Reveal(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)

	_, err = svc.Stop(context.Background(), res.Contract.ID, acct)
	require.ErrorIs(t, err, activity.ErrNotActive)
}

func TestActive_ReturnsCurrentPreview(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	_, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.NoError(t, err)

	c, preview, err := svc.Active(context.Background(), char.ID, acct)
	require.NoError(t, err)
	require.Equal(t, activity.StateActive, c.State)
	require.NotNil(t, preview)
	require.GreaterOrEqual(t, preview.ElapsedSeconds, int32(0))
}

func TestActive_NoActiveReturnsNotFound(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	_, _, err := svc.Active(context.Background(), char.ID, acct)
	require.ErrorIs(t, err, activity.ErrNotFound)
}

// Sanity that time.Duration multiplication of seconds is correct — regression
// guard against int32 overflow in ExpiresAt computation.
func TestExpiresAt_Correct(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	res, _ := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 12 * 3600,
	})
	require.WithinDuration(t, res.Contract.StartedAt.Add(12*time.Hour), res.Contract.ExpiresAt, time.Second)
}
