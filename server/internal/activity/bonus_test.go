package activity_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
)

type fakeBonus struct {
	flat, mult float64
	err        error
	kinds      []string
}

func (f *fakeBonus) GatherBonus(_ context.Context, _ uuid.UUID, kind string) (float64, float64, error) {
	f.kinds = append(f.kinds, kind)
	return f.flat, f.mult, f.err
}

func TestStart_FreezesGatherBonusInSnapshot(t *testing.T) {
	svc, _, _, char, acct := newFixtures(t)
	b := &fakeBonus{flat: 1, mult: 0.05}
	svc.WithGatherBonus(b)

	res, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.NoError(t, err)
	require.Equal(t, []string{"mining"}, b.kinds)
	var snap activity.Snapshot
	require.NoError(t, json.Unmarshal(res.Contract.Snapshot, &snap))
	require.Equal(t, 1.0, snap.YieldFlatPerHour)
	require.Equal(t, 0.05, snap.YieldMult)
}

func TestStart_GatherBonusErrorFailsStart(t *testing.T) {
	svc, repo, _, char, acct := newFixtures(t)
	svc.WithGatherBonus(&fakeBonus{err: errors.New("db down")})
	_, err := svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char.ID, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600,
	})
	require.Error(t, err)
	_, err = repo.ActiveByCharacter(context.Background(), char.ID)
	require.Error(t, err, "no contract is created without its snapshot")
}

func TestComputeYield_MasteryBonusRaisesRate(t *testing.T) {
	c := activity.Contract{Kind: activity.KindMining, DurationSeconds: 8 * 3600, Seed: 42}
	c.Snapshot = activity.MarshalSnapshot(activity.Snapshot{})
	plain := activity.ComputeYield(&c, c.DurationSeconds)

	c.Snapshot = activity.MarshalSnapshot(activity.Snapshot{YieldFlatPerHour: 1, YieldMult: 0.10})
	boosted := activity.ComputeYield(&c, c.DurationSeconds)
	again := activity.ComputeYield(&c, c.DurationSeconds)

	require.Equal(t, boosted, again, "deterministic for the same seed + snapshot")
	require.Greater(t, boosted.CommonCount, plain.CommonCount)

	c.Snapshot = json.RawMessage(`{"character_level":3}`) // pre-DD-220 contract
	require.Equal(t, plain, activity.ComputeYield(&c, c.DurationSeconds))
}
