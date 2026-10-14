package activity_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/domain/character"
)

var now = time.Date(2026, 10, 14, 12, 0, 0, 0, time.UTC)

type fakeGathering struct {
	level  int
	xp     []int
	nodes  []activity.GatherNode
	levels map[activity.Kind]int
}

func (f *fakeGathering) Nodes(activity.Kind) []activity.GatherNode { return f.nodes }
func (f *fakeGathering) Level(context.Context, uuid.UUID, activity.Kind) (int, error) {
	return f.level, nil
}
func (f *fakeGathering) AddXP(_ context.Context, _ uuid.UUID, _ activity.Kind, xp int) (activity.GatherProgress, error) {
	f.xp = append(f.xp, xp)
	return activity.GatherProgress{Level: f.level + 1, Gained: xp, LeveledUp: true}, nil
}

type fakeGranter struct{ got []activity.Drops }

func (f *fakeGranter) Grant(_ context.Context, _, _ uuid.UUID, _ int64, d activity.Drops) ([]activity.GrantedItem, error) {
	f.got = append(f.got, d)
	return nil, nil
}

// newGatherFixtures runs on a clock the test moves with *now.
func newGatherFixtures(t *testing.T, level int) (*activity.Service, *fakeGathering, *fakeGranter, uuid.UUID, uuid.UUID) {
	_, repo, txn, char, acct := newFixtures(t)
	chars := &fakeChars{owned: map[uuid.UUID]*character.Character{char.ID: char}}
	svc := activity.NewService(repo, chars, txn, activity.ServiceConfig{Now: func() time.Time { return now }},
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	g := &fakeGathering{level: level, nodes: []activity.GatherNode{
		{Code: "raw_stone", Level: 1, Tier: 0}, {Code: "iron_ore", Level: 16, Tier: 1}, {Code: "mythril_ore", Level: 90, Tier: 7},
	}}
	gr := &fakeGranter{}
	svc.WithGathering(g).WithItemGranter(gr)
	return svc, g, gr, char.ID, acct
}

func startTarget(svc *activity.Service, char, acct uuid.UUID, target string) (*activity.StartResult, error) {
	return svc.Start(context.Background(), activity.StartRequest{
		CharacterID: char, AccountID: acct, Kind: activity.KindMining, DurationSeconds: 3600, Target: target,
	})
}

func TestStart_TargetGatedByGatheringLevel(t *testing.T) {
	svc, _, _, char, acct := newGatherFixtures(t, 20)
	_, err := startTarget(svc, char, acct, "mythril_ore")
	require.ErrorIs(t, err, activity.ErrGatherLevelLow)
	_, err = startTarget(svc, char, acct, "gold_dust")
	require.ErrorIs(t, err, activity.ErrUnknownTarget)

	res, err := startTarget(svc, char, acct, "raw_stone")
	require.NoError(t, err)
	var snap activity.Snapshot
	require.NoError(t, json.Unmarshal(res.Contract.Snapshot, &snap))
	require.Equal(t, "raw_stone", snap.Target, "a lower resource can still be chosen")
}

func TestStart_NoTargetPicksBestUnlocked(t *testing.T) {
	svc, _, _, char, acct := newGatherFixtures(t, 20)
	res, err := startTarget(svc, char, acct, "")
	require.NoError(t, err)
	var snap activity.Snapshot
	require.NoError(t, json.Unmarshal(res.Contract.Snapshot, &snap))
	require.Equal(t, "iron_ore", snap.Target)
	require.Equal(t, 1, snap.TargetTier)
}

func TestReveal_GrantsTargetAndGatheringXP(t *testing.T) {
	svc, g, gr, char, acct := newGatherFixtures(t, 20)
	res, err := startTarget(svc, char, acct, "iron_ore")
	require.NoError(t, err)
	now = now.Add(2 * time.Hour)
	out, err := svc.Reveal(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)
	units := out.Yield.CommonCount + out.Yield.RareCount + out.Yield.LegendaryCount
	require.Positive(t, units)
	require.Len(t, gr.got, 1)
	require.Equal(t, "iron_ore", gr.got[0].Target)
	require.Equal(t, units, gr.got[0].Common+gr.got[0].Rare+gr.got[0].Legendary)
	require.Equal(t, []int{units * 12}, g.xp, "tier 1 gives 12 XP a unit")
	require.NotNil(t, out.Gather)
	require.True(t, out.Gather.LeveledUp)

	again, err := svc.Reveal(context.Background(), res.Contract.ID, acct)
	require.NoError(t, err)
	require.True(t, again.Replayed)
	require.Len(t, g.xp, 1, "a replay grants no XP")
}

func TestComputeYield_HigherTierIsSlower(t *testing.T) {
	c := activity.Contract{Kind: activity.KindMining, DurationSeconds: 8 * 3600, Seed: 7}
	c.Snapshot = activity.MarshalSnapshot(activity.Snapshot{Target: "raw_stone", TargetTier: 0})
	low := activity.ComputeYield(&c, c.DurationSeconds)
	c.Snapshot = activity.MarshalSnapshot(activity.Snapshot{Target: "mythril_ore", TargetTier: 7})
	high := activity.ComputeYield(&c, c.DurationSeconds)
	require.Equal(t, "mythril_ore", high.Target)
	require.Less(t, high.CommonCount+high.RareCount+high.LegendaryCount, low.CommonCount+low.RareCount+low.LegendaryCount)
}
