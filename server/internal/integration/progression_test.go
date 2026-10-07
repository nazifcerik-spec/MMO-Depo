//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/progression"
)

func TestProgressionRecord(t *testing.T) {
	w := newWorld(t)
	cat, err := progression.LoadCatalog(w.ctx, w.p)
	require.NoError(t, err)
	svc := progression.NewService(w.p, cat, w.chars, w.txn, nil)
	p := w.player("warrior", 1)
	ev := func(kind string, n int64) progression.Event { return progression.Event{Kind: kind, Amount: n} }

	// First kill: mastery XP + first_blood (500c) paid once even if replayed.
	svc.Record(w.ctx, p.ID, ev("mob_killed", 1))
	require.Equal(t, int64(500), w.balance(p))
	svc.Record(w.ctx, p.ID, ev("mob_killed", 1))
	require.Equal(t, int64(500), w.balance(p), "an unlocked achievement pays once")

	// 600 ore: mining mastery tier 3 and miner_500 (title + 5000c).
	svc.Record(w.ctx, p.ID, ev("gather_mining", 600))
	sum, err := svc.Summary(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	var mining, combat progression.MasteryView
	for _, m := range sum.Mastery {
		switch m.Track {
		case "mining":
			mining = m
		case "combat":
			combat = m
		}
	}
	require.Equal(t, int64(6000), mining.XP)
	require.Equal(t, int16(3), mining.Tier)
	require.NotEmpty(t, mining.ActiveBonus, "tier-1 mining bonus is active")
	require.Equal(t, int64(30), combat.XP, "2 kills × 15 xp")
	require.Equal(t, int64(5500), w.balance(p))
	require.GreaterOrEqual(t, sum.AchievementPoints, 15)

	titles, err := svc.Titles(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	require.Equal(t, "prospector", titles[0].Code)
	require.ErrorIs(t, svc.SetActiveTitle(w.ctx, p.ID, p.Account, "bossbane"), progression.ErrTitleNotOwned)
	require.NoError(t, svc.SetActiveTitle(w.ctx, p.ID, p.Account, "prospector"))
	sum, err = svc.Summary(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	require.Equal(t, "prospector", *sum.ActiveTitle)
	require.NoError(t, svc.SetActiveTitle(w.ctx, p.ID, p.Account, ""))

	achievements, err := svc.Achievements(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	unlocked := map[string]bool{}
	for _, a := range achievements {
		unlocked[a.Code] = a.UnlockedAt != nil
	}
	require.True(t, unlocked["first_blood"])
	require.True(t, unlocked["miner_500"])
	require.False(t, unlocked["slayer_100"])

	// Reputation: a dungeon clear plus a kill flood is capped per day.
	svc.Record(w.ctx, p.ID, ev("dungeon_cleared:crypt_of_ashes", 1), ev("mob_killed", 5000))
	sum, err = svc.Summary(w.ctx, p.ID, p.Account)
	require.NoError(t, err)
	var guard progression.ReputationView
	for _, r := range sum.Reputation {
		if r.Faction == "ironhold_guard" {
			guard = r
		}
	}
	require.Equal(t, 1500, guard.Rep, "daily gain cap")
	require.Equal(t, 0, guard.DailyLeft)

	// Ownership is enforced on reads.
	other := w.player("mage", 1)
	_, err = svc.Summary(w.ctx, p.ID, other.Account)
	require.Error(t, err)
}
