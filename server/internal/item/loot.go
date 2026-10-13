package item

import "github.com/otter/server/internal/activity"

// LootPool defines what an activity kind can drop, grouped by rarity tier.
// A code listed twice in a tier is twice as likely (weights by repetition).
type LootPool struct {
	Common    []string
	Uncommon  []string
	Rare      []string
	Epic      []string
	Legendary []string
	Mythic    []string
}

// PoolPerKind is the loot table per activity kind. Mining follows the ore
// ladder (DD-243): very common/common → Common, medium/precious → Uncommon,
// rare → Rare, very rare → Epic, legendary → Legendary, mythic → Mythic.
var PoolPerKind = map[activity.Kind]LootPool{
	activity.KindMining: {
		Common:    []string{"coal", "coal", "stone_ore", "stone_ore", "copper_ore", "iron_ore", "tin_ore", "lead_ore", "zinc_ore"},
		Uncommon:  []string{"cobalt_ore", "cobalt_ore", "nickel_ore", "nickel_ore", "silver_ore", "gold_ore"},
		Rare:      []string{"palladium_ore", "rhodium_ore", "titanium_ore"},
		Epic:      []string{"platinum_ore", "iridium_ore", "osmium_ore"},
		Legendary: []string{"orichalcum_ore", "mythril_ore"},
		Mythic:    []string{"adamantite_ore"},
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

// Share of each yield tier that steps up one rarity (DD-243). The yield
// model counts common / rare / legendary finds; without this split the
// uncommon, epic and mythic tiers never dropped (copper, tin, mandrake and
// wolf pelt were unobtainable, so the recipes using them could not be made).
const (
	upgradeCommonPct    = 25 // common → uncommon
	upgradeRarePct      = 25 // rare → epic
	upgradeLegendaryPct = 15 // legendary → mythic
)

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
	out = append(out, rollTier(cat, rng, tier{pool.Common, RarityCommon}, tier{pool.Uncommon, RarityUncommon}, upgradeCommonPct, commonCount)...)
	out = append(out, rollTier(cat, rng, tier{pool.Rare, RarityRare}, tier{pool.Epic, RarityEpic}, upgradeRarePct, rareCount)...)
	out = append(out, rollTier(cat, rng, tier{pool.Legendary, RarityLegendary}, tier{pool.Mythic, RarityMythic}, upgradeLegendaryPct, legendaryCount)...)
	return out
}

type tier struct {
	codes  []string
	rarity Rarity
}

// rollTier draws n finds of a yield tier; each steps up to the next rarity
// with upPct% when that tier has items. An empty base tier uses the upper.
func rollTier(cat *Catalog, rng *activity.RNG, base, up tier, upPct, n int) []DropSpec {
	if n <= 0 || (len(base.codes) == 0 && len(up.codes) == 0) {
		return nil
	}
	out := make([]DropSpec, 0, n)
	for i := 0; i < n; i++ {
		t := base
		if len(up.codes) > 0 && (len(base.codes) == 0 || rng.IntN(100) < upPct) {
			t = up
		}
		pool, rarity := t.codes, t.rarity
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
