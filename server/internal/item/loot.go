package item

import "github.com/otter/server/internal/activity"

// LootPool defines what an activity kind can drop, grouped by rarity tier.
// Each roll (common count → rare count → legendary count) uniformly picks
// from the tier's template codes.
type LootPool struct {
	Common    []string
	Uncommon  []string
	Rare      []string
	Epic      []string
	Legendary []string
}

// PoolPerKind is the v1 baseline loot table per activity kind. Content-driven
// version lands in Sprint 12.
var PoolPerKind = map[activity.Kind]LootPool{
	activity.KindMining: {
		Common:    []string{"iron_ore"},
		Uncommon:  []string{"copper_ore", "tin_ore"},
		Rare:      []string{"silver_ore", "gold_ore"},
		Legendary: []string{"mythril_ore"},
	},
	activity.KindFishing: {
		Common:    []string{"salmon"},
		Uncommon:  []string{"bass"},
		Rare:      []string{"golden_carp"},
		Legendary: []string{"ancient_leviathan_scale"},
	},
	activity.KindWoodcutting: {
		Common:    []string{"oak_log"},
		Uncommon:  []string{"elm_log"},
		Rare:      []string{"yew_log"},
		Legendary: []string{"ironwood_log"},
	},
	activity.KindHerbalism: {
		Common:    []string{"linseed"},
		Uncommon:  []string{"mandrake"},
		Rare:      []string{"king_bloom"},
		Legendary: []string{"phoenix_petal"},
	},
	activity.KindHunting: {
		Common:    []string{"boar_hide"},
		Uncommon:  []string{"wolf_pelt"},
		Rare:      []string{"dire_bear_hide"},
		Legendary: []string{"dragon_scale"},
	},
}

// DropSpec is a single item to be created — resolved template + rarity + affix
// choices. Instance persistence is a separate concern (service.MaterializeDrops).
type DropSpec struct {
	TemplateCode string
	Rarity       Rarity
	Prefix       *string
	Suffix       *string
}

// RollDrops is a pure function: given (seed, kind, rarity counts) it produces
// the DropSpec list deterministically. No DB, no side effects — tests and
// preview code both call this.
//
// Seed derivation: uses seed+1 as the RNG start so item rolls don't collide
// with yield-count rolls in the caller (activity/yield.go uses seed).
func RollDrops(cat *Catalog, seed int64, kind activity.Kind, commonCount, rareCount, legendaryCount int) []DropSpec {
	pool, ok := PoolPerKind[kind]
	if !ok {
		return nil
	}
	rng := activity.NewRNG(seed + 1) // separate stream from yield count roll

	var out []DropSpec
	out = append(out, rollTier(cat, rng, pool.Common, RarityCommon, commonCount)...)
	out = append(out, rollTier(cat, rng, pool.Rare, RarityRare, rareCount)...)
	out = append(out, rollTier(cat, rng, pool.Legendary, RarityLegendary, legendaryCount)...)
	return out
}

func rollTier(cat *Catalog, rng *activity.RNG, pool []string, rarity Rarity, n int) []DropSpec {
	if n <= 0 || len(pool) == 0 {
		return nil
	}
	out := make([]DropSpec, 0, n)
	for i := 0; i < n; i++ {
		code := pool[rng.IntN(len(pool))]
		tmpl, err := cat.Get(code)
		if err != nil {
			continue
		}
		prefix, suffix := rollAffix(rng, tmpl, rarity)
		out = append(out, DropSpec{
			TemplateCode: code,
			Rarity:       rarity,
			Prefix:       prefix,
			Suffix:       suffix,
		})
	}
	return out
}
