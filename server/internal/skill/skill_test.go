package skill_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/skill"
)

func TestComputeSheet_NoModifiers_ProducesBaseStats(t *testing.T) {
	base := skill.Primary{STR: 15, AGI: 10, INT: 8, VIT: 15, CHA: 10, Luck: 5}
	s := skill.ComputeSheet(base, nil)

	require.Equal(t, 15, s.Primary.STR)
	require.Equal(t, 10, s.Primary.AGI)
	require.Equal(t, 8, s.Primary.INT)
	require.Equal(t, 15, s.Primary.VIT)

	// Derived: attack = STR*2, defense = VIT, crit = 5+Luck/2, hit = 90+AGI/2 (DD-229)
	require.Equal(t, 30, s.Derived.Attack)
	require.Equal(t, 15, s.Derived.Defense)
	require.Equal(t, 7, s.Derived.CritPct) // 5 + 5/2 = 7 (int round)
	require.Equal(t, 95, s.Derived.HitPct) // 90 + 10/2 = 95

	// Resources: hp=100+VIT*10=250, mp=50+INT*5=90
	require.Equal(t, 250, s.Resources.HPMax)
	require.Equal(t, 90, s.Resources.MPMax)
	require.Equal(t, 100, s.Resources.EnMax)
}

func TestComputeSheet_FlatBonuses(t *testing.T) {
	base := skill.Primary{STR: 15, AGI: 10, INT: 8, VIT: 15, CHA: 10, Luck: 5}
	// sword_iron: +10 attack, +3 STR
	mods := []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetSTR, Amount: 3},
		{Kind: skill.KindFlat, Target: skill.TargetAttack, Amount: 10},
	}
	s := skill.ComputeSheet(base, mods)

	require.Equal(t, 18, s.Primary.STR)
	// Attack recomputed from new STR: 18*2 = 36 + flat 10 = 46
	require.Equal(t, 46, s.Derived.Attack)
}

func TestComputeSheet_MultAddCombines(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 10, INT: 10, VIT: 10, CHA: 10, Luck: 10}
	// Two mult_add: +10% and +5% — additive, so 15% total.
	mods := []skill.Modifier{
		{Kind: skill.KindMultAdd, Target: skill.TargetSTR, Amount: 0.10},
		{Kind: skill.KindMultAdd, Target: skill.TargetSTR, Amount: 0.05},
	}
	s := skill.ComputeSheet(base, mods)
	// 10 * (1 + 0.15) = 11.5 → round to 12
	require.Equal(t, 12, s.Primary.STR)
}

func TestComputeSheet_MultProdCompounds(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 10, INT: 10, VIT: 10, CHA: 10, Luck: 10}
	// Two mult_prod: ×1.10 and ×1.10 = ×1.21 (compound)
	mods := []skill.Modifier{
		{Kind: skill.KindMultProd, Target: skill.TargetSTR, Amount: 0.10},
		{Kind: skill.KindMultProd, Target: skill.TargetSTR, Amount: 0.10},
	}
	s := skill.ComputeSheet(base, mods)
	// 10 × 1.10 × 1.10 = 12.1 → 12
	require.Equal(t, 12, s.Primary.STR)
}

func TestComputeSheet_OrderMatters_FlatThenMult(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 10, INT: 10, VIT: 10, CHA: 10, Luck: 10}
	// (10 + 5) × 1.20 = 18
	mods := []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetSTR, Amount: 5},
		{Kind: skill.KindMultAdd, Target: skill.TargetSTR, Amount: 0.20},
	}
	s := skill.ComputeSheet(base, mods)
	require.Equal(t, 18, s.Primary.STR)
}

func TestComputeSheet_ClampMaxCapsValue(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 10, INT: 10, VIT: 10, CHA: 10, Luck: 10}
	// crit base 5 + Luck/2 = 10; +50 flat = 60; clamp_max 25 = 25
	mods := []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetCritPct, Amount: 50},
		{Kind: skill.KindClampMax, Target: skill.TargetCritPct, Amount: 25},
	}
	s := skill.ComputeSheet(base, mods)
	require.Equal(t, 25, s.Derived.CritPct)
}

func TestComputeSheet_SetOverrides(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 10, INT: 10, VIT: 10, CHA: 10, Luck: 10}
	mods := []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetSTR, Amount: 100}, // would take STR to 110
		{Kind: skill.KindSet, Target: skill.TargetSTR, Amount: 42},   // overrides
	}
	s := skill.ComputeSheet(base, mods)
	require.Equal(t, 42, s.Primary.STR)
}

func TestComputeSheet_MinimumClamp_PrimaryNeverBelow1(t *testing.T) {
	base := skill.Primary{STR: 10, AGI: 10, INT: 10, VIT: 10, CHA: 10, Luck: 10}
	// Massive negative flat.
	mods := []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetSTR, Amount: -100},
	}
	s := skill.ComputeSheet(base, mods)
	// Auto-clamped to 1 by ComputeSheet's intClamp on primary.
	require.Equal(t, 1, s.Primary.STR)
}

func TestComputeSheet_FullEquipmentLoadout(t *testing.T) {
	// Warrior baseline + iron sword + iron helm + leather chest + iron boots.
	base := skill.Primary{STR: 15, AGI: 10, INT: 8, VIT: 15, CHA: 10, Luck: 5}
	mods := []skill.Modifier{
		// sword_iron
		{Kind: skill.KindFlat, Target: skill.TargetAttack, Amount: 10},
		{Kind: skill.KindFlat, Target: skill.TargetSTR, Amount: 3},
		// helm_iron
		{Kind: skill.KindFlat, Target: skill.TargetDefense, Amount: 5},
		{Kind: skill.KindFlat, Target: skill.TargetVIT, Amount: 2},
		// chest_leather
		{Kind: skill.KindFlat, Target: skill.TargetDefense, Amount: 8},
		{Kind: skill.KindFlat, Target: skill.TargetVIT, Amount: 1},
		{Kind: skill.KindFlat, Target: skill.TargetHPMax, Amount: 25},
		// boots_iron
		{Kind: skill.KindFlat, Target: skill.TargetDefense, Amount: 3},
		{Kind: skill.KindFlat, Target: skill.TargetAGI, Amount: 2},
	}
	s := skill.ComputeSheet(base, mods)

	require.Equal(t, 18, s.Primary.STR) // 15 + 3
	require.Equal(t, 12, s.Primary.AGI) // 10 + 2
	require.Equal(t, 18, s.Primary.VIT) // 15 + 3

	// Attack: 18*2 + 10 = 46
	require.Equal(t, 46, s.Derived.Attack)
	// Defense: 18 + (5+8+3) = 34
	require.Equal(t, 34, s.Derived.Defense)
	// HP: 100 + 18*10 + 25 = 305
	require.Equal(t, 305, s.Resources.HPMax)
}
