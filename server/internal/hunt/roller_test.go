package hunt

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
)

func newRNG(seed int64) combat.RNG { return activity.NewRNG(seed) }

func sturdy() combat.Combatant {
	return combat.Combatant{Name: "player", HPMax: 500, HPCurrent: 500,
		Attack: 40, Defense: 15, CritPct: 10, HitPct: 90, EvasionPct: 10}
}

func boar() combat.Mob {
	return combat.Mob{Code: "boar", HPMax: 40, Attack: 6, Defense: 1, CritPct: 3, HitPct: 80, EvasionPct: 3,
		XPReward: 12, GoldCopperMin: 5, GoldCopperMax: 15}
}

// Fights run back to back (DD-281): an hour holds as many kills as their
// turns allow, each paying the mob's own XP; gold and drops follow the time.
func TestRoll_BackToBack(t *testing.T) {
	res := Roll(1234, sturdy(), boar(), 3600, combat.DefaultMaxTurns, newRNG)
	require.Zero(t, res.Defeats)
	require.False(t, res.StoppedOnDeath)
	require.Greater(t, res.Victories, 100, "short fights, no cadence")
	require.Equal(t, res.Victories, res.Completed)
	require.Equal(t, int64(12*res.Victories), res.TotalXP, "every kill pays the mob's XP")
	require.Equal(t, int32(res.FightTurns*TurnSeconds)+res.RestSecs, res.ElapsedSecs)
	require.LessOrEqual(t, res.ElapsedSecs, int32(3600))
	require.Greater(t, res.ElapsedSecs, int32(3600-60), "the hour is used up")
	require.InDelta(t, int(res.ElapsedSecs)/LootSeconds, res.LootRolls, 1, "the wins' time in LootSeconds")
	// Gold per hour of hunting: 30 × 5..15 copper, however fast the kills.
	require.InDelta(t, 300, res.TotalGoldCopper, 120)
	require.Len(t, res.Encounters, MaxLoggedEncounters, "the log keeps the first fights")
	require.Greater(t, res.PlayerHPFinal, int32(0))
}

// Combat power sets the kill time: a stronger hunter kills more in an hour.
func TestRoll_StrongerHuntsFaster(t *testing.T) {
	mob := boar()
	mob.HPMax = 400
	weak, strong := sturdy(), sturdy()
	weak.Attack, strong.Attack = 30, 90
	a := Roll(5, weak, mob, 3600, combat.DefaultMaxTurns, newRNG)
	b := Roll(5, strong, mob, 3600, combat.DefaultMaxTurns, newRNG)
	require.Greater(t, b.Victories, 2*a.Victories)
	require.Less(t, float64(b.FightTurns)/float64(b.Victories), float64(a.FightTurns)/float64(a.Victories))
}

// A hunter worn under RestBelowPct camps to full; the camp counts as time.
func TestRoll_RestsWhenWounded(t *testing.T) {
	p := sturdy()
	mob := combat.Mob{Code: "wolf", HPMax: 120, Attack: 70, Defense: 2, HitPct: 90, XPReward: 20}
	res := Roll(9, p, mob, 3600, combat.DefaultMaxTurns, newRNG)
	require.Positive(t, res.RestSecs)
	require.Positive(t, res.Victories)
	p.RegenPerSec = 50 // regen that outpaces the damage: no camp needed
	healed := Roll(9, p, mob, 3600, combat.DefaultMaxTurns, newRNG)
	require.Less(t, healed.RestSecs, res.RestSecs)
	require.Greater(t, healed.Victories, res.Victories)
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
	res := Roll(42, player, mob, 3600, combat.DefaultMaxTurns, newRNG)

	require.True(t, res.StoppedOnDeath, "player MUST die vs vastly stronger mob")
	require.Equal(t, 1, res.Defeats)
	require.Zero(t, res.Victories)
	require.Zero(t, res.TotalXP, "no reward on defeat-only run")
	require.Equal(t, int32(0), res.PlayerHPFinal)
	require.Equal(t, 1, res.Completed)
}

// A mob nothing kills ends the hunt on the turn cap without reward.
func TestRoll_TimeoutStops(t *testing.T) {
	p := combat.Combatant{Name: "p", HPMax: 1000, HPCurrent: 1000, Attack: 1, Defense: 100, HitPct: 50}
	mob := combat.Mob{Code: "wall", HPMax: 100000, Attack: 1, Defense: 1000, HitPct: 50, XPReward: 5}
	res := Roll(3, p, mob, 3600, 10, newRNG)
	require.Equal(t, 1, res.Timeouts)
	require.Zero(t, res.TotalXP)
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
	a := Roll(777, player, mob, 1800, combat.DefaultMaxTurns, newRNG)
	b := Roll(777, player, mob, 1800, combat.DefaultMaxTurns, newRNG)
	require.Equal(t, a, b)
}

// No time, no fight: an empty log without touching the RNG.
func TestRoll_ZeroBudget(t *testing.T) {
	player := combat.Combatant{Name: "p", HPCurrent: 100, HPMax: 100, Attack: 10, HitPct: 90}
	mob := combat.Mob{Code: "any", HPMax: 20, Attack: 3}
	res := Roll(0, player, mob, 0, combat.DefaultMaxTurns, newRNG)
	require.Equal(t, 0, res.Completed)
	require.Empty(t, res.Encounters)
	require.Equal(t, int32(100), res.PlayerHPFinal)
	short := Roll(0, player, mob, 2, combat.DefaultMaxTurns, newRNG)
	require.Zero(t, short.Completed, "a fight longer than the budget is not fought")
}

// A kill pays the mob's own XP scaled by the configured falloff (DD-281).
func TestHuntXP(t *testing.T) {
	require.Equal(t, int32(200), HuntXP(200, 10, 10), "default config: no falloff")
}
