package combat

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
)

// activity.NewRNG returns *activity.RNG which has Uint64() — satisfies our
// combat.RNG. Wrap it as the injectable factory Session.Resolve expects.
func newRNG(seed int64) RNG { return activity.NewRNG(seed) }

func TestTryHit_MissAlwaysWhenHitFloorMinus(t *testing.T) {
	// hit_pct 0 vs evasion 0 → effective becomes MinHitPct (15), so hits are
	// possible but rare. Sample a bunch and ensure at least some misses.
	r := activity.NewRNG(1)
	misses := 0
	for i := 0; i < 400; i++ {
		hit, _ := TryHit(r, 0, 0, 0)
		if !hit {
			misses++
		}
	}
	require.Greater(t, misses, 100, "with effective 15%% hit, expect majority misses")
}

func TestTryHit_MaxCapAlwaysAllowsSomeMiss(t *testing.T) {
	// hit_pct 100 vs evasion 0 → capped to MaxHitPct (95). Sample enough that
	// we should statistically see at least one miss.
	r := activity.NewRNG(2)
	misses := 0
	for i := 0; i < 2000; i++ {
		hit, _ := TryHit(r, 100, 0, 0)
		if !hit {
			misses++
		}
	}
	require.Greater(t, misses, 0, "cap should keep miss possible")
}

func TestRollDamage_NeverBelowFloor(t *testing.T) {
	r := activity.NewRNG(3)
	for i := 0; i < 100; i++ {
		// attack < defense * 0.5 — raw would go negative; floor keeps it at 1.
		dmg := RollDamage(r, 1, 999, false)
		require.GreaterOrEqual(t, dmg, int32(MinDamageFloor))
	}
}

func TestRollDamage_CritDoublesRoughly(t *testing.T) {
	// With variance ±10%, crit=2x should still be clearly separable at scale.
	base := activity.NewRNG(4)
	crit := activity.NewRNG(4) // same seed on purpose — same variance sequence

	var sumBase, sumCrit int64
	for i := 0; i < 500; i++ {
		sumBase += int64(RollDamage(base, 100, 10, false))
		sumCrit += int64(RollDamage(crit, 100, 10, true))
	}
	// Crit sum should be ~2× base sum, within 10%.
	ratio := float64(sumCrit) / float64(sumBase)
	require.InDelta(t, 2.0, ratio, 0.1)
}

func TestSession_ResolveVictoryDeterministic(t *testing.T) {
	newSession := func() *Session {
		return &Session{
			Player: Combatant{Name: "hero", HPMax: 100, HPCurrent: 100, Attack: 20, Defense: 5, CritPct: 10, HitPct: 90, EvasionPct: 8},
			Mob:    Combatant{Name: "boar", HPMax: 40, HPCurrent: 40, Attack: 6, Defense: 1, CritPct: 3, HitPct: 80, EvasionPct: 3},
		}
	}
	a, err := newSession().Resolve(newRNG, 42)
	require.NoError(t, err)
	b, err := newSession().Resolve(newRNG, 42)
	require.NoError(t, err)
	require.Equal(t, a.Outcome, b.Outcome, "same seed → same outcome")
	require.Equal(t, a.Duration, b.Duration, "same seed → same duration")
	require.Equal(t, len(a.Turns), len(b.Turns), "same seed → identical log length")
	for i := range a.Turns {
		require.Equal(t, a.Turns[i], b.Turns[i], "turn %d must match", i)
	}
	require.Equal(t, OutcomeVictory, a.Outcome)
}

func TestSession_ResolveDefeatWhenOutclassed(t *testing.T) {
	s := &Session{
		Player: Combatant{Name: "novice", HPMax: 30, HPCurrent: 30, Attack: 3, Defense: 0, HitPct: 60, EvasionPct: 2},
		Mob:    Combatant{Name: "lich", HPMax: 500, HPCurrent: 500, Attack: 40, Defense: 20, CritPct: 20, HitPct: 95, EvasionPct: 5},
	}
	res, err := s.Resolve(newRNG, 7)
	require.NoError(t, err)
	require.Equal(t, OutcomeDefeat, res.Outcome)
	require.Zero(t, res.PlayerHPEnd)
	require.Greater(t, res.MobHPEnd, int32(0))
}

func TestSession_ResolveTimeoutOnStalemate(t *testing.T) {
	// Both sides so tanky they can't dent each other in the given MaxTurns.
	s := &Session{
		Player:   Combatant{Name: "wall_a", HPMax: 100000, HPCurrent: 100000, Attack: 5, Defense: 999, HitPct: 90},
		Mob:      Combatant{Name: "wall_b", HPMax: 100000, HPCurrent: 100000, Attack: 5, Defense: 999, HitPct: 90},
		MaxTurns: 5,
	}
	res, err := s.Resolve(newRNG, 99)
	require.NoError(t, err)
	require.Equal(t, OutcomeTimeout, res.Outcome)
	require.Equal(t, 5, res.Duration)
}

func TestSession_ValidateRejectsBrokenInput(t *testing.T) {
	s := &Session{
		Player: Combatant{Name: "ghost", HPCurrent: 0, Attack: 10},
		Mob:    Combatant{Name: "boar", HPCurrent: 40, Attack: 5},
	}
	_, err := s.Resolve(newRNG, 1)
	require.ErrorContains(t, err, "player HPCurrent")

	s = &Session{
		Player: Combatant{Name: "hero", HPCurrent: 10, Attack: 0},
		Mob:    Combatant{Name: "wisp", HPCurrent: 10, Attack: 0},
	}
	_, err = s.Resolve(newRNG, 1)
	require.ErrorContains(t, err, "zero attack")
}

func TestSession_PlayerKillsBeforeCounter(t *testing.T) {
	// Glass-cannon hero one-shots weak mob — should NEVER die on turn 1.
	s := &Session{
		Player: Combatant{Name: "hero", HPMax: 10, HPCurrent: 10, Attack: 999, HitPct: 100},
		Mob:    Combatant{Name: "slime", HPMax: 5, HPCurrent: 5, Attack: 999, HitPct: 100},
	}
	res, err := s.Resolve(newRNG, 1)
	require.NoError(t, err)
	require.Equal(t, OutcomeVictory, res.Outcome, "player attacks first — mob never gets its swing")
	require.Equal(t, int32(10), res.PlayerHPEnd)
	require.Equal(t, 1, res.Duration)
}

func TestMobAsCombatant_FullHP(t *testing.T) {
	m := Mob{Code: "test", HPMax: 100, Attack: 10, Defense: 3, CritPct: 5, HitPct: 80, EvasionPct: 10}
	c := m.AsCombatant()
	require.Equal(t, int32(100), c.HPCurrent)
	require.Equal(t, int32(100), c.HPMax)
	require.Equal(t, "test", c.Name)
	require.Equal(t, int32(10), c.Attack)
}
