package combat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// fixedRNG returns the same value every draw.
type fixedRNG uint64

func (f fixedRNG) Uint64() uint64 { return uint64(f) }

// DD-250: Rage starts empty, fills on landed hits and skill gains, and pays
// for skills; HP-threshold triggers gate auto-casts.
func TestRageAndTriggers(t *testing.T) {
	self := &Combatant{Name: "w", HPMax: 100, HPCurrent: 100, Attack: 50, HitPct: 95, MPMax: 100,
		Resource: "rage", ResourcePerHit: 10, Skills: []Skill{
			{Code: "execute", Effect: EffectStrike, Power: 300, MPCost: 30, CooldownTurns: 1, TargetHPBelow: 35},
			{Code: "strike", Effect: EffectStrike, Power: 110, CooldownTurns: 1, Gain: 15},
		}}
	mob := &Combatant{Name: "m", HPMax: 1000, HPCurrent: 1000, Defense: 0}
	k := newCaster(self)
	require.Zero(t, k.mp, "Rage starts empty")

	rec := k.act(fixedRNG(0), 1, self, mob)
	require.Equal(t, "strike", rec.Skill, "execute waits for a low target")
	require.Equal(t, int32(25), k.mp, "10 per hit + 15 from Strike")

	mob.HPCurrent = 300 // 30%
	rec = k.act(fixedRNG(0), 2, self, mob)
	require.Equal(t, "strike", rec.Skill, "not enough Rage yet (25 < 30)")
	rec = k.act(fixedRNG(0), 3, self, mob)
	require.Equal(t, "execute", rec.Skill)
	require.Equal(t, int32(30), k.mp, "50 - 30 + 10 for the hit")
}

func TestGuardAndEmpowerBuffsExpire(t *testing.T) {
	self := &Combatant{Name: "w", HPMax: 100, HPCurrent: 50, Attack: 100, Defense: 40, HitPct: 95, MPMax: 100,
		Skills: []Skill{
			{Code: "wall", Effect: EffectGuard, Power: 50, MPCost: 10, CooldownTurns: 5, DurationTurns: 2, SelfHPBelow: 60},
			{Code: "rush", Effect: EffectEmpower, Power: 20, MPCost: 10, CooldownTurns: 5, DurationTurns: 2},
		}}
	mob := &Combatant{Name: "m", HPMax: 1000, HPCurrent: 1000}
	k := newCaster(self)
	require.Equal(t, "wall", k.act(fixedRNG(0), 1, self, mob).Skill)
	require.Equal(t, int32(60), self.Defense)
	require.Equal(t, "rush", k.act(fixedRNG(0), 2, self, mob).Skill)
	require.Equal(t, int32(120), self.Attack)
	k.act(fixedRNG(0), 3, self, mob)
	require.Equal(t, int32(40), self.Defense, "guard ended")
	k.act(fixedRNG(0), 4, self, mob)
	require.Equal(t, int32(100), self.Attack, "empower ended")

	self.HPCurrent = 90
	k2 := newCaster(self)
	self.Skills = self.Skills[:1]
	require.Empty(t, k2.act(fixedRNG(0), 1, self, mob).Skill, "guard waits for its HP threshold")
}

func TestProgressionDamage(t *testing.T) {
	a := &Combatant{HPMax: 100, HPCurrent: 50}
	d := &Combatant{HPMax: 100, HPCurrent: 100}
	require.Equal(t, int32(100), progressionDamage(fixedRNG(0), 100, true, a, d), "all zero: untouched")

	a.CritDmgPct, a.ExecutePct, a.MissingHPDmgPct, a.LifestealPct = 50, 20, 30, 10
	d.HPCurrent = 30 // executable
	got := progressionDamage(fixedRNG(99), 100, true, a, d)
	require.Equal(t, int32(100*(100+50+20+15)/100), got, "crit + execute + half the missing-HP bonus")
	require.Equal(t, int32(50+18), a.HPCurrent, "10% lifesteal")

	a = &Combatant{HPMax: 100, HPCurrent: 100}
	d = &Combatant{HPMax: 100, HPCurrent: 100, DmgReductionPct: 20, BlockPct: 50}
	require.Equal(t, int32(40), progressionDamage(fixedRNG(10), 100, false, a, d), "20% reduced, then blocked (roll 10 < 50)")
	require.Equal(t, int32(80), progressionDamage(fixedRNG(60), 100, false, a, d), "not blocked (roll 60)")
	require.Equal(t, int32(1), progressionDamage(fixedRNG(0), 1, false, a, d), "never below 1")
}
