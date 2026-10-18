package combat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-261 group encounters.

func duelist() Combatant {
	return Combatant{Name: "p", HPMax: 300, HPCurrent: 300, Attack: 40, Defense: 10, HitPct: 85, CritPct: 10, MPMax: 50,
		Skills: []Skill{{Code: "cut", Effect: EffectStrike, Power: 150, MPCost: 10, CooldownTurns: 3},
			{Code: "mend", Effect: EffectHeal, Power: 30, MPCost: 20, CooldownTurns: 5}}}
}

func TestGroupOfOneMatchesDuel(t *testing.T) {
	for seed := int64(1); seed <= 40; seed++ {
		mob := Combatant{Name: "m", HPMax: 900, HPCurrent: 900, Attack: 45, Defense: 8, HitPct: 85, CritPct: 5}
		s := &Session{Player: duelist().Clone(), Mob: mob}
		want, err := s.Resolve(newRNG, seed)
		require.NoError(t, err)
		g := &Group{Players: []Combatant{duelist().Clone()}, Mob: mob}
		got, err := g.Resolve(newRNG, seed)
		require.NoError(t, err)
		require.Equal(t, want.Outcome, got.Outcome, "seed %d", seed)
		require.Equal(t, want.PlayerHPEnd, got.PlayersHPEnd[0], "seed %d", seed)
		require.Equal(t, want.MobHPEnd, got.MobHPEnd, "seed %d", seed)
		require.Equal(t, want.Duration, got.Duration, "seed %d", seed)
		require.Len(t, got.Turns, len(want.Turns))
	}
}

func TestGroupHealerHealsTheWoundedAndMobHitsTheTank(t *testing.T) {
	tankP := Combatant{Name: "tank", HPMax: 1000, HPCurrent: 400, Attack: 5, Defense: 60, HitPct: 95}
	healer := Combatant{Name: "healer", HPMax: 300, HPCurrent: 300, Attack: 5, Defense: 5, HitPct: 95, MPMax: 100,
		Skills: []Skill{{Code: "mend", Effect: EffectHeal, Power: 20, MPCost: 10, CooldownTurns: 1},
			{Code: "ward", Effect: EffectBarrier, Power: 10, MPCost: 10, CooldownTurns: 1, DurationTurns: 2}}}
	mob := Combatant{Name: "m", HPMax: 100000, HPCurrent: 100000, Attack: 80, HitPct: 95}
	g := &Group{Players: []Combatant{tankP, healer}, Mob: mob, MaxTurns: 6}
	res, err := g.Resolve(newRNG, 7)
	require.NoError(t, err)
	require.Equal(t, OutcomeTimeout, res.Outcome)
	var healed, wardedTank bool
	for _, rec := range res.Turns {
		if rec.Actor == "m" {
			require.Equal(t, "tank", rec.Target, "the mob holds on the highest defense")
		}
		if rec.Skill == "mend" {
			require.Equal(t, "tank", rec.Target, "the most wounded ally is healed")
			healed = true
		}
		if rec.Skill == "ward" && rec.Target == "tank" {
			wardedTank = true
		}
	}
	require.True(t, healed)
	require.True(t, wardedTank)
}

func TestGroupAreaAttackHitsEveryone(t *testing.T) {
	a := Combatant{Name: "a", HPMax: 500, HPCurrent: 500, Attack: 1, Defense: 30, HitPct: 95}
	b := Combatant{Name: "b", HPMax: 500, HPCurrent: 500, Attack: 1, Defense: 0, HitPct: 95}
	mob := Combatant{Name: "m", HPMax: 100000, HPCurrent: 100000, Attack: 60, HitPct: 95}
	g := &Group{Players: []Combatant{a, b}, Mob: mob, MaxTurns: 4, AoEEvery: 2, AoEPct: 50}
	res, err := g.Resolve(newRNG, 3)
	require.NoError(t, err)
	hit := map[int]map[string]bool{}
	for _, rec := range res.Turns {
		if rec.Skill == SkillAoE {
			if hit[rec.Turn] == nil {
				hit[rec.Turn] = map[string]bool{}
			}
			hit[rec.Turn][rec.Target] = true
		}
	}
	require.Equal(t, map[int]map[string]bool{2: {"a": true, "b": true}, 4: {"a": true, "b": true}}, hit)
}

func TestGroupStunImmuneAndWipe(t *testing.T) {
	stunner := Combatant{Name: "s", HPMax: 50, HPCurrent: 50, Attack: 10, HitPct: 95, MPMax: 100,
		Skills: []Skill{{Code: "bash", Effect: EffectStun, Power: 100, MPCost: 1, CooldownTurns: 1, DurationTurns: 3}}}
	frail := Combatant{Name: "f", HPMax: 40, HPCurrent: 40, Attack: 10, HitPct: 95}
	mob := Combatant{Name: "m", HPMax: 100000, HPCurrent: 100000, Attack: 200, HitPct: 95}
	g := &Group{Players: []Combatant{stunner, frail}, Mob: mob, MaxTurns: 20, StunImmune: true}
	res, err := g.Resolve(newRNG, 11)
	require.NoError(t, err)
	require.Equal(t, OutcomeDefeat, res.Outcome)
	require.Equal(t, []int32{0, 0}, res.PlayersHPEnd)
	for _, rec := range res.Turns {
		require.False(t, rec.Stunned, "a stun-immune mob never loses a swing")
	}

	// Without immunity the same stun holds the mob.
	g = &Group{Players: []Combatant{stunner, frail}, Mob: mob, MaxTurns: 3}
	res, err = g.Resolve(newRNG, 11)
	require.NoError(t, err)
	stunned := false
	for _, rec := range res.Turns {
		stunned = stunned || rec.Stunned
	}
	require.True(t, stunned)

	_, err = (&Group{Players: []Combatant{{HPCurrent: 0}}, Mob: mob}).Resolve(newRNG, 1)
	require.Error(t, err)
	_, err = (&Group{Players: []Combatant{frail}, Mob: Combatant{}}).Resolve(newRNG, 1)
	require.Error(t, err)
}
