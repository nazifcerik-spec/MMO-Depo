package combat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-252 shared effects: barrier, damage/heal over time, armor debuff,
// stun; generic resources (start share, per-turn regeneration).

func TestBarrierAbsorbsAndFades(t *testing.T) {
	self := &Combatant{Name: "p", HPMax: 200, HPCurrent: 200, Attack: 10, HitPct: 95, MPMax: 100,
		Skills: []Skill{{Code: "ward", Effect: EffectBarrier, Power: 25, MPCost: 10, CooldownTurns: 9, DurationTurns: 2}}}
	mob := &Combatant{Name: "m", HPMax: 100, HPCurrent: 100, Attack: 30, HitPct: 95}
	k := newCaster(self)
	require.Equal(t, "ward", k.act(fixedRNG(0), 1, self, mob).Skill)
	require.Equal(t, int32(50), self.Barrier)
	rec := swing(fixedRNG(0), 1, mob, self, 100)
	require.Positive(t, rec.Absorbed)
	require.Equal(t, 200-rec.Damage, self.HPCurrent, "only the unabsorbed part reaches HP")
	require.Equal(t, 50-rec.Absorbed, self.Barrier)
	k.act(fixedRNG(0), 3, self, mob)
	require.Zero(t, self.Barrier, "faded after its turns")
}

func TestDamageAndHealOverTime(t *testing.T) {
	self := &Combatant{Name: "p", HPMax: 100, HPCurrent: 60, Attack: 100, HitPct: 95, MPMax: 100,
		Skills: []Skill{
			{Code: "mend", Effect: EffectHot, Power: 10, MPCost: 5, CooldownTurns: 6, DurationTurns: 3},
			{Code: "rot", Effect: EffectDot, Power: 50, MPCost: 5, CooldownTurns: 6, DurationTurns: 3},
		}}
	mob := &Combatant{Name: "m", HPMax: 1000, HPCurrent: 1000, Defense: 20}
	k := newCaster(self)
	rec := k.act(fixedRNG(0), 1, self, mob)
	require.Equal(t, "mend", rec.Skill)
	require.Equal(t, int32(10), rec.Heal, "first tick on cast")
	rec = k.act(fixedRNG(0), 2, self, mob)
	require.Equal(t, "rot", rec.Skill)
	require.Equal(t, int32(10), rec.TickHeal)
	require.Equal(t, int32((100-10)*50/100), rec.Damage, "mitigated attack × power")
	rec = k.act(fixedRNG(0), 3, self, mob)
	require.Equal(t, int32(45), rec.Tick)
	require.Equal(t, int32(10), rec.TickHeal)
	rec = k.act(fixedRNG(0), 4, self, mob)
	require.Equal(t, int32(45), rec.Tick)
	require.Zero(t, rec.TickHeal, "heal over time ran its 3 turns")
	rec = k.act(fixedRNG(0), 5, self, mob)
	require.Zero(t, rec.Tick, "damage over time ran its 3 turns")

	// A tick can finish the fight before the action.
	mob.HPCurrent = 10
	k.dot = tick{left: 1, amount: 45}
	rec = k.act(fixedRNG(0), 6, self, mob)
	require.Equal(t, int32(10), rec.Tick)
	require.Zero(t, mob.HPCurrent)
}

func TestWeakenAndSetupFloor(t *testing.T) {
	self := &Combatant{Name: "p", HPMax: 100, HPCurrent: 100, Attack: 50, HitPct: 95, MPMax: 100,
		Skills: []Skill{{Code: "sunder", Effect: EffectWeaken, Power: 50, MPCost: 5, CooldownTurns: 3, DurationTurns: 2}}}
	mob := &Combatant{Name: "m", HPMax: 1000, HPCurrent: 1000, Defense: 40}
	k := newCaster(self)
	require.Equal(t, "sunder", k.act(fixedRNG(0), 1, self, mob).Skill)
	require.Equal(t, int32(20), mob.Defense)
	k.act(fixedRNG(0), 3, self, mob)
	require.Equal(t, int32(40), mob.Defense, "debuff expired")
	mob.HPCurrent = 300
	k.ready[0] = 0
	require.Empty(t, k.act(fixedRNG(0), 9, self, mob).Skill, "no set-up against a nearly dead target")
}

func TestStunSkipsMobSwings(t *testing.T) {
	p := Combatant{Name: "player", HPMax: 500, HPCurrent: 500, Attack: 20, HitPct: 95, MPMax: 100,
		Skills: []Skill{{Code: "bash", Effect: EffectStun, Power: 100, MPCost: 0, CooldownTurns: 4, DurationTurns: 2}}}
	s := Session{Player: p, Mob: Combatant{Name: "mob", HPMax: 5000, HPCurrent: 5000, Attack: 30, HitPct: 95}, MaxTurns: 4}
	res, err := s.Resolve(newRNG, 3)
	require.NoError(t, err)
	var stunned []int
	for _, tr := range res.Turns {
		if tr.Stunned {
			stunned = append(stunned, tr.Turn)
		}
	}
	if res.Turns[0].Hit {
		require.Equal(t, []int{1, 2}, stunned, "two lost swings, no chaining")
	}
}

func TestGenericResourceStartAndRegen(t *testing.T) {
	self := &Combatant{Name: "p", HPMax: 100, HPCurrent: 100, Attack: 10, HitPct: 95, MPMax: 100,
		Resource: "energy", ResourceStart: 40, ResourcePerTurn: 15,
		Skills: []Skill{{Code: "big", Effect: EffectStrike, Power: 200, MPCost: 60, CooldownTurns: 1}}}
	mob := &Combatant{Name: "m", HPMax: 1000, HPCurrent: 1000}
	k := newCaster(self)
	require.Equal(t, int32(40), k.mp)
	require.Empty(t, k.act(fixedRNG(0), 1, self, mob).Skill, "55 < 60")
	require.Equal(t, "big", k.act(fixedRNG(0), 2, self, mob).Skill, "70 after another turn")
	require.Equal(t, int32(10), k.mp)
	require.True(t, Effects[EffectStun] && Timed(EffectDot) && !Timed(EffectStrike))
}

// DD-254: past five levels of gap the player also loses attack.
func TestLevelGapAttackPenalty(t *testing.T) {
	base := Combatant{Level: 50, HitPct: 90, Defense: 100, Attack: 1000}
	p := base
	applyLevelGap(&p, 54)
	require.Equal(t, int32(1000), p.Attack, "within the grace: attack kept")
	require.Equal(t, int32(90-4*4), p.HitPct)
	p = base
	applyLevelGap(&p, 60)
	require.Equal(t, int32(1000*(100-10*5)/100), p.Attack, "10 levels up: 5 beyond the grace")
	require.Equal(t, int32(100*(100-6*10)/100), p.Defense)
	p = base
	applyLevelGap(&p, 90)
	require.Equal(t, int32(100), p.Attack, "never below 10%")
	p = Combatant{Attack: 1000}
	applyLevelGap(&p, 90)
	require.Equal(t, int32(1000), p.Attack, "unknown level: no gap rule")
}
