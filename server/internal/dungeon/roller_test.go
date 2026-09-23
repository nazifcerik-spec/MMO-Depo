package dungeon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
)

// --- fake mob catalog -------------------------------------------------------

type fakeMobs map[string]combat.Mob

func (f fakeMobs) Get(code string) (combat.Mob, bool) {
	m, ok := f[code]
	return m, ok
}

func newRNG(seed int64) combat.RNG { return activity.NewRNG(seed) }

func strong() combat.Combatant {
	return combat.Combatant{
		Name: "hero", HPMax: 1000, HPCurrent: 1000,
		Attack: 80, Defense: 30, CritPct: 15, HitPct: 92, EvasionPct: 15,
	}
}

func fragile() combat.Combatant {
	return combat.Combatant{
		Name: "novice", HPMax: 30, HPCurrent: 30,
		Attack: 3, Defense: 0, CritPct: 0, HitPct: 55, EvasionPct: 2,
	}
}

// --- Tests ------------------------------------------------------------------

func TestRoll_ClearedRunGrantsBossBonus(t *testing.T) {
	mobs := fakeMobs{
		"boar": {Code: "boar", HPMax: 40, Attack: 5, Defense: 1, CritPct: 3, HitPct: 75, EvasionPct: 3,
			XPReward: 10, GoldCopperMin: 5, GoldCopperMax: 10},
		"boss": {Code: "boss", HPMax: 200, Attack: 20, Defense: 5, CritPct: 5, HitPct: 85, EvasionPct: 5,
			XPReward: 200, GoldCopperMin: 100, GoldCopperMax: 100},
	}
	tmpl := Template{
		Code: "d1", DurationSecs: 900,
		Waves:           []Wave{{MobCode: "boar", Count: 3}},
		BossMobCode:     "boss",
		BonusGoldCopper: 500, BonusXP: 250,
	}
	res := Roll(1, strong(), tmpl, mobs, combat.DefaultMaxTurns, newRNG)
	require.True(t, res.Cleared)
	require.True(t, res.BossKilled)
	require.False(t, res.StoppedOnDeath)
	require.Equal(t, 4, res.Victories, "3 waves + 1 boss")
	require.GreaterOrEqual(t, res.TotalXP, int64(3*10+200+250), "wave xp + boss xp + bonus xp")
	require.GreaterOrEqual(t, res.TotalGoldCopper, int64(3*5+100+500))
	require.GreaterOrEqual(t, res.RareDrops, 1, "boss kill grants guaranteed rare")
}

func TestRoll_StopsOnDeath(t *testing.T) {
	mobs := fakeMobs{
		"drake": {Code: "drake", HPMax: 500, Attack: 50, Defense: 20, CritPct: 20, HitPct: 95, EvasionPct: 5,
			XPReward: 999, GoldCopperMin: 500, GoldCopperMax: 500},
		"boss": {Code: "boss", HPMax: 800, Attack: 100, Defense: 40, HitPct: 95, EvasionPct: 5},
	}
	tmpl := Template{
		Code: "d2", DurationSecs: 900,
		Waves:           []Wave{{MobCode: "drake", Count: 3}},
		BossMobCode:     "boss",
		BonusGoldCopper: 5000, BonusXP: 5000,
	}
	res := Roll(2, fragile(), tmpl, mobs, combat.DefaultMaxTurns, newRNG)
	require.True(t, res.StoppedOnDeath)
	require.False(t, res.Cleared)
	require.False(t, res.BossKilled)
	require.Zero(t, res.Victories)
	require.Zero(t, res.TotalXP, "defeat first wave → no bonus xp granted")
}

func TestRoll_DeterministicReplay(t *testing.T) {
	mobs := fakeMobs{
		"grunt": {Code: "grunt", HPMax: 80, Attack: 10, Defense: 3, HitPct: 82, EvasionPct: 6,
			XPReward: 20, GoldCopperMin: 8, GoldCopperMax: 20},
		"boss": {Code: "boss", HPMax: 300, Attack: 25, Defense: 10, HitPct: 85, EvasionPct: 5,
			XPReward: 300, GoldCopperMin: 200, GoldCopperMax: 200},
	}
	tmpl := Template{
		Code: "d3", DurationSecs: 1200,
		Waves:       []Wave{{MobCode: "grunt", Count: 4}},
		BossMobCode: "boss",
	}
	a := Roll(777, strong(), tmpl, mobs, combat.DefaultMaxTurns, newRNG)
	b := Roll(777, strong(), tmpl, mobs, combat.DefaultMaxTurns, newRNG)
	require.Equal(t, a.Cleared, b.Cleared)
	require.Equal(t, a.Victories, b.Victories)
	require.Equal(t, a.TotalGoldCopper, b.TotalGoldCopper)
	require.Equal(t, a.TotalXP, b.TotalXP)
	require.Equal(t, len(a.Encounters), len(b.Encounters))
	for i := range a.Encounters {
		require.Equal(t, a.Encounters[i], b.Encounters[i], "encounter %d", i)
	}
}

func TestRoll_UnknownMobHaltsGracefully(t *testing.T) {
	mobs := fakeMobs{
		"boss": {Code: "boss", HPMax: 50, Attack: 5},
	}
	tmpl := Template{
		Code: "d4", DurationSecs: 900,
		Waves:       []Wave{{MobCode: "ghost", Count: 1}}, // not in catalog
		BossMobCode: "boss",
	}
	res := Roll(3, strong(), tmpl, mobs, combat.DefaultMaxTurns, newRNG)
	require.False(t, res.Cleared, "ghost mob is missing → roll aborts before boss")
	require.Empty(t, res.Encounters)
}

func TestISOWeek_PadsAndSortsChronologically(t *testing.T) {
	early := time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC)
	late := time.Date(2026, 12, 28, 12, 0, 0, 0, time.UTC)
	require.Equal(t, "2026-02", ISOWeek(early), "first week of January often becomes ISO week 2")
	// The precise ISO year for late-December may be 2026 or 2027 depending on
	// dates — we assert the padding + sort-friendliness rather than the value.
	got := ISOWeek(late)
	require.Len(t, got, 7)
	require.Less(t, ISOWeek(early), got, "lexical sort must match chronological")
}

func TestTemplate_TotalEncounters(t *testing.T) {
	tmpl := Template{
		Waves:       []Wave{{MobCode: "a", Count: 3}, {MobCode: "b", Count: 2}},
		BossMobCode: "boss",
	}
	require.Equal(t, 6, tmpl.TotalEncounters(), "3 + 2 + boss")
}

func TestContract_IsReady(t *testing.T) {
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	c := &Contract{StartedAt: start, DurationSecs: 900}
	require.False(t, c.IsReady(start.Add(500*time.Second)))
	require.True(t, c.IsReady(start.Add(900*time.Second)))
	require.True(t, c.IsReady(start.Add(9999*time.Second)))
}
