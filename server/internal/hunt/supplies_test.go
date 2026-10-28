package hunt

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/item"
)

func wolf() combat.Mob {
	return combat.Mob{Code: "wolf", HPMax: 120, Attack: 70, Defense: 2, HitPct: 90, XPReward: 20,
		GoldCopperMin: 5, GoldCopperMax: 15}
}

func used(res *HuntResult, code string) int32 {
	for _, u := range res.SuppliesUsed {
		if u.Code == code {
			return u.Qty
		}
	}
	return 0
}

// Potions are drunk before a camp would be needed (DD-288): a hunter with
// enough of them camps less and kills more in the hour.
func TestRollWith_PotionsReplaceCamping(t *testing.T) {
	bare := Roll(9, sturdy(), wolf(), 3600, combat.DefaultMaxTurns, newRNG)
	load := Loadout{Supplies: []Supply{
		{Code: "small", Kind: SupplyHP, Qty: 50, HealPct: 20},
		{Code: "big", Kind: SupplyHP, Qty: 50, HealPct: 40},
	}}
	res := RollWith(9, sturdy(), wolf(), 3600, combat.DefaultMaxTurns, newRNG, load)
	require.Less(t, res.RestSecs, bare.RestSecs)
	require.Greater(t, res.Victories, bare.Victories)
	require.Positive(t, used(res, "big"), "the strongest potion first")
	require.LessOrEqual(t, used(res, "big"), int32(50))

	// Run out and it camps again; never more than it brought.
	few := RollWith(9, sturdy(), wolf(), 3600, combat.DefaultMaxTurns, newRNG,
		Loadout{Supplies: []Supply{{Code: "big", Kind: SupplyHP, Qty: 2, HealPct: 40}}})
	require.Equal(t, int32(2), used(few, "big"))
	require.Positive(t, few.RestSecs)
}

// Before the first fight a wounded hunter tops up with potions that would
// not overheal.
func TestRollWith_TopsUpFirst(t *testing.T) {
	p := sturdy()
	p.HPCurrent = 380 // 120 missing: one 20% potion (100) fits, a second would overheal
	res := RollWith(1, p, boar(), 60, combat.DefaultMaxTurns, newRNG,
		Loadout{Supplies: []Supply{{Code: "small", Kind: SupplyHP, Qty: 5, HealPct: 20}}})
	require.Equal(t, int32(1), used(res, "small"))
}

// Food heals and buffs attack and defense while its time runs; a new dish
// is eaten when it runs out.
func TestRollWith_FoodBuff(t *testing.T) {
	mob := boar()
	mob.HPMax = 600
	bare := Roll(3, sturdy(), mob, 3600, combat.DefaultMaxTurns, newRNG)
	fed := RollWith(3, sturdy(), mob, 3600, combat.DefaultMaxTurns, newRNG,
		Loadout{Supplies: []Supply{{Code: "stew", Kind: SupplyFood, Qty: 3, HealPct: 10, BuffAtk: 50, BuffDef: 20, BuffSecs: 900}}})
	require.Greater(t, fed.Victories, bare.Victories, "the buff speeds kills up")
	require.Equal(t, int32(3), used(fed, "stew"), "eaten at 0, 15 and 30 minutes")

	// A dish eaten before the hunt counts for the time it has left.
	early := RollWith(3, sturdy(), mob, 3600, combat.DefaultMaxTurns, newRNG,
		Loadout{Buff: &Buff{AttackPct: 50, DefensePct: 20, Secs: 900}})
	require.Greater(t, early.Victories, bare.Victories)
	require.Less(t, early.Victories, fed.Victories)
}

// A carried pool (mana) does not refill between fights: once dry the
// hunter fights on with basic attacks until mana potions (or a camp) bring
// it back (DD-288).
func TestRollWith_ManaCarriesOver(t *testing.T) {
	mage := sturdy()
	mage.Attack = 20
	mage.MPMax, mage.Resource, mage.ResourceCarry = 100, "mana", true
	mage.ResourceStart, mage.ResourcePerTurn = 70, 8 // the mage's pool (progression.yaml)
	mage.Skills = []combat.Skill{{Code: "bolt", Effect: combat.EffectStrike, Power: 400, MPCost: 25, CooldownTurns: 1}}
	mob := boar()
	mob.HPMax = 200
	dry := Roll(4, mage, mob, 3600, combat.DefaultMaxTurns, newRNG)

	potions := RollWith(4, mage, mob, 3600, combat.DefaultMaxTurns, newRNG,
		Loadout{Supplies: []Supply{{Code: "mana", Kind: SupplyMP, Qty: 200, MPPct: 50}}})
	require.Positive(t, used(potions, "mana"))
	require.Greater(t, potions.Victories, dry.Victories)

	refill := mage
	refill.ResourceCarry = false // old snapshots: the pool starts every fight anew
	fresh := Roll(4, refill, mob, 3600, combat.DefaultMaxTurns, newRNG)
	require.Greater(t, fresh.Victories, dry.Victories, "a full pool every fight casts more")
}

// Under FleeBelowPct HP after a swing the hunter breaks off: the hunt ends,
// alive, without the fight's reward.
func TestRollWith_FleesNearlyDead(t *testing.T) {
	p := sturdy()
	p.HPMax, p.HPCurrent, p.Defense = 1000, 1000, 0
	// Hits of at most ~33 cannot jump the 0..50 window: it always flees.
	mob := combat.Mob{Code: "ogre", HPMax: 50000, Attack: 30, HitPct: 95, XPReward: 50}
	res := Roll(7, p, mob, 3600, combat.DefaultMaxTurns, newRNG)
	require.True(t, res.Fled)
	require.False(t, res.StoppedOnDeath)
	require.Equal(t, combat.EndFled, res.Report.End)
	require.Equal(t, "ogre", res.Report.EndMob)
	require.Less(t, res.Report.LowestHPPct, FleeBelowPct)
	require.Zero(t, res.Victories)
	require.Zero(t, res.TotalXP)
	require.Positive(t, res.PlayerHPFinal)
	require.Less(t, res.PlayerHPFinal*100, p.HPMax*FleeBelowPct)
	require.Equal(t, combat.OutcomeFled, res.Encounters[len(res.Encounters)-1].Outcome)
}

func TestSupplyFromTemplate(t *testing.T) {
	cases := map[string]struct {
		ue   *item.UseEffect
		want Supply
		ok   bool
	}{
		"potion pct": {&item.UseEffect{Kind: "heal_pct", Amount: 30}, Supply{Code: "x", Kind: SupplyHP, Qty: 2, HealPct: 30}, true},
		"potion hp":  {&item.UseEffect{Kind: "heal_hp", Amount: 60}, Supply{Code: "x", Kind: SupplyHP, Qty: 2, HealFlat: 60}, true},
		"mana":       {&item.UseEffect{Kind: "mana_pct", Amount: 40}, Supply{Code: "x", Kind: SupplyMP, Qty: 2, MPPct: 40}, true},
		"food": {&item.UseEffect{Kind: "heal_pct", Amount: 15, Buff: &item.UseBuff{AttackPct: 4, DefensePct: 3, Minutes: 60}},
			Supply{Code: "x", Kind: SupplyFood, Qty: 2, HealPct: 15, BuffAtk: 4, BuffDef: 3, BuffSecs: 3600}, true},
		"no effect": {nil, Supply{}, false},
		"unknown":   {&item.UseEffect{Kind: "teleport", Amount: 1}, Supply{}, false},
	}
	for name, c := range cases {
		got, ok := SupplyFromTemplate(item.Template{Code: "x", UseEffect: c.ue}, 2)
		require.Equal(t, c.ok, ok, name)
		require.Equal(t, c.want, got, name)
	}
}
