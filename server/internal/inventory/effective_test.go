package inventory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/skill"
)

// A worn item adds its bonuses scaled by rarity plus its affixes (DD-275).
func TestEffectiveModifiers(t *testing.T) {
	raw := json.RawMessage(`[{"kind":"flat","target":"attack","amount":100},{"kind":"mult_add","target":"hp_max","amount":0.1}]`)
	sharp := "sharp"
	common := EffectiveModifiers(BagItem{StatBonusesRaw: raw, Rarity: "common", MinLevel: 100})
	require.Equal(t, []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetAttack, Amount: 100},
		{Kind: skill.KindMultAdd, Target: skill.TargetHPMax, Amount: 0.1},
	}, common)
	epic := EffectiveModifiers(BagItem{StatBonusesRaw: raw, Rarity: "epic", Prefix: &sharp, MinLevel: 100})
	require.Equal(t, []skill.Modifier{
		{Kind: skill.KindFlat, Target: skill.TargetAttack, Amount: 120},
		{Kind: skill.KindMultAdd, Target: skill.TargetHPMax, Amount: 0.1}, // percent kinds are not scaled
		{Kind: skill.KindFlat, Target: skill.TargetAttack, Amount: 12},
	}, epic)
	require.Nil(t, EffectiveModifiers(BagItem{StatBonusesRaw: json.RawMessage(`{bad`)}))
}
