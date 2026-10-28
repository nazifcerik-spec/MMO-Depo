package combat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReportFight(t *testing.T) {
	you := Combatant{Name: "player", HPMax: 200, HitPct: 90, EvasionPct: 10, Level: 5}
	mob := Combatant{Name: "wolf", HPMax: 100, HitPct: 80, EvasionPct: 20, Level: 6}
	r := NewReport()
	r.Fight(you, nil, mob, []TurnRecord{
		{Actor: "player", Target: "wolf", Hit: true, Damage: 30, Crit: true},
		{Actor: "wolf", Target: "player", Hit: true, Damage: 50, TargetHPAft: 150},
		{Actor: "player", Target: "wolf", Hit: false},
		{Actor: "wolf", Target: "player", Stunned: true},
		{Actor: "player", Target: "player", Hit: true, Heal: 40, Skill: "heal"}, // no swing
		{Actor: "player", Target: "wolf", Hit: true, Tick: 7, TargetHPAft: 0},   // a damage-over-time kill
	})
	require.Equal(t, 1, r.Fights)
	require.Equal(t, Side{Swings: 2, Hits: 1, Crits: 1, Damage: 37, MaxHit: 30, Healed: 40}, r.Yours)
	require.Equal(t, Side{Swings: 1, Hits: 1, Damage: 50, MaxHit: 50}, r.Theirs)
	require.Equal(t, 75, r.LowestHPPct)
	require.Equal(t, int32(70), r.HitChance)
	require.Equal(t, int32(70), r.EnemyHitChance)
	require.Equal(t, "wolf", r.EnemyID)
	require.Equal(t, int16(6), r.Enemy.Level)
}

func TestHitChanceClamps(t *testing.T) {
	require.Equal(t, int32(MinHitPct), HitChance(50, 90))
	require.Equal(t, int32(MaxHitPct), HitChance(200, 0))
}

func TestReportForMember(t *testing.T) {
	r := NewReport()
	r.Enemy = Stats{HitPct: 90, EvasionPct: 30}
	me := r.For(Combatant{HPMax: 500, HitPct: 100, EvasionPct: 40})
	require.Equal(t, int32(500), me.You.HPMax)
	require.Equal(t, int32(70), me.HitChance)
	require.Equal(t, int32(50), me.EnemyHitChance)
	require.Zero(t, r.You.HPMax, "the shared report is untouched")
}
