package hunt

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
)

func newRNG(seed int64) combat.RNG { return activity.NewRNG(seed) }

// Sturdy player who reliably beats a weak mob — victory chain with no defeats.
func TestRoll_VictoryChain(t *testing.T) {
	player := combat.Combatant{
		Name: "player", HPMax: 500, HPCurrent: 500,
		Attack: 40, Defense: 15, CritPct: 10, HitPct: 90, EvasionPct: 10,
	}
	mob := combat.Mob{
		Code: "boar", HPMax: 40, Attack: 6, Defense: 1,
		CritPct: 3, HitPct: 80, EvasionPct: 3,
		XPReward: 12, GoldCopperMin: 5, GoldCopperMax: 15,
	}
	res := Roll(1234, player, mob, 10, combat.DefaultMaxTurns, newRNG)

	require.Equal(t, 10, res.PlannedEncounters)
	require.Equal(t, 10, res.Completed)
	require.Equal(t, 10, res.Victories)
	require.Equal(t, 0, res.Defeats)
	require.False(t, res.StoppedOnDeath)
	require.Equal(t, int64(120), res.TotalXP, "10 kills × 12 xp")
	require.GreaterOrEqual(t, res.TotalGoldCopper, int64(50), "10 kills × min 5c")
	require.LessOrEqual(t, res.TotalGoldCopper, int64(150), "10 kills × max 15c")
	require.Greater(t, res.PlayerHPFinal, int32(0))
}

// Fragile player fighting a heavy mob — defeat stops the chain early.
func TestRoll_StopsOnDeath(t *testing.T) {
	player := combat.Combatant{
		Name: "player", HPMax: 30, HPCurrent: 30,
		Attack: 3, Defense: 0, CritPct: 0, HitPct: 60, EvasionPct: 2,
	}
	mob := combat.Mob{
		Code: "lich", HPMax: 500, Attack: 50, Defense: 20,
		CritPct: 25, HitPct: 95, EvasionPct: 5,
		XPReward: 999, GoldCopperMin: 500, GoldCopperMax: 500,
	}
	res := Roll(42, player, mob, 5, combat.DefaultMaxTurns, newRNG)

	require.True(t, res.StoppedOnDeath, "player MUST die vs vastly stronger mob")
	require.Equal(t, 1, res.Defeats)
	require.Zero(t, res.Victories)
	require.Zero(t, res.TotalXP, "no reward on defeat-only run")
	require.Equal(t, int32(0), res.PlayerHPFinal)
	require.LessOrEqual(t, res.Completed, 5)
}

// Same inputs + seed → identical result (deterministic replay).
func TestRoll_DeterministicReplay(t *testing.T) {
	player := combat.Combatant{
		Name: "p", HPMax: 200, HPCurrent: 200,
		Attack: 20, Defense: 5, CritPct: 10, HitPct: 85, EvasionPct: 10,
	}
	mob := combat.Mob{
		Code: "wolf", HPMax: 60, Attack: 9, Defense: 2,
		CritPct: 5, HitPct: 82, EvasionPct: 6,
		XPReward: 22, GoldCopperMin: 10, GoldCopperMax: 25,
	}
	a := Roll(777, player, mob, 6, combat.DefaultMaxTurns, newRNG)
	b := Roll(777, player, mob, 6, combat.DefaultMaxTurns, newRNG)
	require.Equal(t, a.Victories, b.Victories)
	require.Equal(t, a.Defeats, b.Defeats)
	require.Equal(t, a.TotalXP, b.TotalXP)
	require.Equal(t, a.TotalGoldCopper, b.TotalGoldCopper)
	require.Equal(t, a.CommonDrops, b.CommonDrops)
	require.Equal(t, a.RareDrops, b.RareDrops)
	require.Equal(t, a.LegendaryDrops, b.LegendaryDrops)
	require.Equal(t, len(a.Encounters), len(b.Encounters))
	for i := range a.Encounters {
		require.Equal(t, a.Encounters[i], b.Encounters[i], "encounter %d must match", i)
	}
}

// Zero planned encounters returns an empty log without touching RNG.
func TestRoll_ZeroPlanned(t *testing.T) {
	player := combat.Combatant{Name: "p", HPCurrent: 100, HPMax: 100, Attack: 10, HitPct: 90}
	mob := combat.Mob{Code: "any", HPMax: 20, Attack: 3}
	res := Roll(0, player, mob, 0, combat.DefaultMaxTurns, newRNG)
	require.Equal(t, 0, res.PlannedEncounters)
	require.Equal(t, 0, res.Completed)
	require.Empty(t, res.Encounters)
	require.Equal(t, int32(100), res.PlayerHPFinal)
}

// Contract helper: EncountersPlanned rounds down.
func TestContract_EncountersPlannedRoundsDown(t *testing.T) {
	c := &Contract{EncounterSecs: 120, DurationSecs: 3600}
	require.Equal(t, 30, c.EncountersPlanned(3600))
	require.Equal(t, 30, c.EncountersPlanned(3719)) // 30.99 → 30
	require.Equal(t, 0, c.EncountersPlanned(119))   // less than one full cadence
}

func TestContract_EncountersPlannedZeroCadenceSafe(t *testing.T) {
	c := &Contract{EncounterSecs: 0, DurationSecs: 3600}
	require.Equal(t, 0, c.EncountersPlanned(3600), "zero cadence must not divide by zero")
}
