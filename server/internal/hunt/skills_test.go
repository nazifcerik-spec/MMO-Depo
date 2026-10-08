package hunt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
)

// DD-222: actives + MP survive the JSON start snapshot and are cast in
// every encounter of the chain.
func TestRoll_UsesSnapshotSkills(t *testing.T) {
	sheet := &character.Sheet{Resources: character.Resources{HPCurrent: 500, HPMax: 500, MPMax: 40},
		Derived: character.DerivedStats{Attack: 40, Defense: 15, CritPct: 10, HitPct: 90}}
	p := combat.FromSheet(sheet)
	require.Equal(t, int32(40), p.MPMax)
	p.Skills = []combat.Skill{{Code: "cleave", Effect: combat.EffectStrike, Power: 180, MPCost: 10, CooldownTurns: 2}}

	raw, err := json.Marshal(p)
	require.NoError(t, err)
	var frozen combat.Combatant
	require.NoError(t, json.Unmarshal(raw, &frozen))
	require.Equal(t, p, frozen)

	mob := combat.Mob{Code: "boar", HPMax: 200, Attack: 6, Defense: 1, HitPct: 80, XPReward: 12}
	res := Roll(77, frozen, mob, 3, combat.DefaultMaxTurns, newRNG)
	require.Equal(t, 3, res.Victories)

	plain := frozen.Clone()
	plain.Skills = nil
	base := Roll(77, plain, mob, 3, combat.DefaultMaxTurns, newRNG)
	shorter := 0
	for i := range res.Encounters {
		if res.Encounters[i].Duration < base.Encounters[i].Duration {
			shorter++
		}
	}
	require.Positive(t, shorter, "skills end some fights sooner")
}
