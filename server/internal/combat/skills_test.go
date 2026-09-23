package combat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func fighter(skills ...Skill) Combatant {
	return Combatant{Name: "player", HPMax: 200, HPCurrent: 200, Attack: 30, Defense: 5, CritPct: 0, HitPct: 95,
		MPMax: 30, Skills: skills}
}

func mob() Combatant {
	return Combatant{Name: "mob", HPMax: 400, HPCurrent: 400, Attack: 20, Defense: 5, HitPct: 90}
}

func TestSkills_StrikeCooldownAndMP(t *testing.T) {
	cleave := Skill{Code: "cleave", Effect: EffectStrike, Power: 200, MPCost: 10, CooldownTurns: 2}
	s := Session{Player: fighter(cleave), Mob: mob(), MaxTurns: 8}
	res, err := s.Resolve(newRNG, 7)
	require.NoError(t, err)

	var casts []int
	for _, tr := range res.Turns {
		if tr.Actor == "player" && tr.Skill == "cleave" {
			casts = append(casts, tr.Turn)
		}
	}
	require.Equal(t, []int{1, 3, 5}, casts, "every 2 turns until 30 MP runs out")
}

func TestSkills_StrikeHitsHarder(t *testing.T) {
	plain := Session{Player: fighter(), Mob: mob(), MaxTurns: 1}
	a, err := plain.Resolve(newRNG, 11)
	require.NoError(t, err)
	strong := Session{Player: fighter(Skill{Code: "x", Effect: EffectStrike, Power: 300, CooldownTurns: 1}), Mob: mob(), MaxTurns: 1}
	b, err := strong.Resolve(newRNG, 11)
	require.NoError(t, err)
	require.True(t, a.Turns[0].Hit && b.Turns[0].Hit, "same seed → same hit roll")
	require.Equal(t, max(1, a.Turns[0].Damage*3), b.Turns[0].Damage)
}

func TestSkills_HealOnlyWhenLow(t *testing.T) {
	mend := Skill{Code: "mend", Effect: EffectHeal, Power: 40, MPCost: 10, CooldownTurns: 3}
	p := fighter(mend)
	s := Session{Player: p, Mob: mob(), MaxTurns: 1}
	res, err := s.Resolve(newRNG, 3)
	require.NoError(t, err)
	require.Empty(t, res.Turns[0].Skill, "full HP → attack instead")

	p.HPCurrent = 60
	s = Session{Player: p, Mob: mob(), MaxTurns: 1}
	res, err = s.Resolve(newRNG, 3)
	require.NoError(t, err)
	require.Equal(t, "mend", res.Turns[0].Skill)
	require.Equal(t, int32(80), res.Turns[0].Heal)
	require.Equal(t, int32(140), res.Turns[0].TargetHPAft)

	p.HPCurrent = 90
	p.HPMax = 100
	p.Skills = []Skill{{Code: "big", Effect: EffectHeal, Power: 90, CooldownTurns: 1}}
	s = Session{Player: p, Mob: mob(), MaxTurns: 1}
	res, err = s.Resolve(newRNG, 3)
	require.NoError(t, err)
	require.Empty(t, res.Turns[0].Skill, "90/100 is above half")

	p.HPCurrent = 40
	s = Session{Player: p, Mob: mob(), MaxTurns: 1}
	res, err = s.Resolve(newRNG, 3)
	require.NoError(t, err)
	require.Equal(t, int32(60), res.Turns[0].Heal, "capped at missing HP")
}

func TestSkills_DeterministicAndCloneIndependent(t *testing.T) {
	p := fighter(Skill{Code: "cleave", Effect: EffectStrike, Power: 150, MPCost: 5, CooldownTurns: 2})
	c := p.Clone()
	c.Skills[0].Power = 999
	require.Equal(t, int32(150), p.Skills[0].Power)

	s1 := Session{Player: p.Clone(), Mob: mob()}
	s2 := Session{Player: p.Clone(), Mob: mob()}
	r1, err := s1.Resolve(newRNG, 99)
	require.NoError(t, err)
	r2, err := s2.Resolve(newRNG, 99)
	require.NoError(t, err)
	require.Equal(t, r1, r2)
}
