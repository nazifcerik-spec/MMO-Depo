package item

import "github.com/otter/server/internal/activity"

// DropSpec is a single item to be created — resolved template + rarity + affix
// choices. Instance persistence is a separate concern (service.MaterializeDrops).
type DropSpec struct {
	TemplateCode string
	Rarity       Rarity
	Prefix       *string
	Suffix       *string
}

// TierRarity colours a gathering tier (DD-244): T0–T1 common, T2–T3
// uncommon, T4 rare, T5 epic, T6 legendary, T7 mythic.
func TierRarity(tier int16) Rarity {
	switch {
	case tier >= 7:
		return RarityMythic
	case tier == 6:
		return RarityLegendary
	case tier == 5:
		return RarityEpic
	case tier == 4:
		return RarityRare
	case tier >= 2:
		return RarityUncommon
	default:
		return RarityCommon
	}
}

// GatherDrops is a gathering contract's haul (DD-244): every find is one
// unit of the node the player chose, coloured by its tier. Unknown or
// non-node targets yield nothing.
func GatherDrops(cat *Catalog, target string, n int) []DropSpec {
	tmpl, err := cat.Get(target)
	if err != nil || tmpl.GatherKind == "" || n <= 0 {
		return nil
	}
	out := make([]DropSpec, n)
	for i := range out {
		out[i] = DropSpec{TemplateCode: target, Rarity: TierRarity(tmpl.GatherTier)}
	}
	return out
}

// combatHideSpan is how many of the best unlocked animals common hides
// come from after a fight.
const combatHideSpan = 4

// CombatHides is skinning after hunts and dungeons (DD-244): the hunting
// level the hides come from is the mob's level / 10, rounded up (profession
// levels 1-100 stand for character levels 1-1000, DD-274).
// Common finds pick among the best few animals it unlocks; rare and
// legendary finds are the best one. Deterministic per seed.
func CombatHides(cat *Catalog, seed int64, mobLevel, common, rare, legendary int) []DropSpec {
	reach := int16(min(100, max(1, (mobLevel+9)/10)))
	var unlocked []Template
	for _, t := range cat.Nodes(string(activity.KindHunting)) {
		if t.GatherLevel <= reach {
			unlocked = append(unlocked, t)
		}
	}
	if len(unlocked) == 0 {
		return nil
	}
	rng := activity.NewRNG(seed + 1) // separate stream from the yield roll
	best := unlocked[len(unlocked)-1]
	span := unlocked[max(0, len(unlocked)-combatHideSpan):]
	out := make([]DropSpec, 0, common+rare+legendary)
	for i := 0; i < common; i++ {
		t := span[rng.IntN(len(span))]
		out = append(out, DropSpec{TemplateCode: t.Code, Rarity: TierRarity(t.GatherTier)})
	}
	for i := 0; i < rare+legendary; i++ {
		out = append(out, DropSpec{TemplateCode: best.Code, Rarity: TierRarity(best.GatherTier)})
	}
	return out
}
