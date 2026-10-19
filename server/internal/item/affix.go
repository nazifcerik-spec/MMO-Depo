package item

import (
	"math"

	"github.com/otter/server/internal/activity"
)

// Affixes: rolled by rarity at drop time; each adds one stat that grows with
// the item's level (DD-275).

// affix pools per item class. Empty class → all pools available.
var prefixPool = map[string][]string{
	"weapon": {"sharp", "cursed", "ancient", "swift", "burning", "frozen", "vicious"},
	"armor":  {"reinforced", "gilded", "ancient", "shadowed", "runic", "warded"},
	// materials + consumables don't get affix in v1 (Affixable=false)
}

var suffixPool = map[string][]string{
	"weapon": {"of_the_bear", "of_precision", "of_the_wolf", "of_ruin", "of_haste"},
	"armor":  {"of_the_turtle", "of_wisdom", "of_the_ox", "of_shielding", "of_endurance"},
}

// affixStat is what one affix adds to an item of level L: base + perLevel × L,
// rounded, at least 1. Percent targets grow slowly so the caps hold at 1000.
type affixStat struct {
	target         string
	base, perLevel float64
}

var affixStats = map[string]affixStat{
	// weapon prefixes
	"sharp":   {"attack", 2, 0.10},
	"burning": {"attack", 1, 0.08},
	"vicious": {"crit_dmg_pct", 3, 0.02},
	"swift":   {"hit_pct", 1, 0.005},
	"frozen":  {"crit_pct", 1, 0.004},
	"cursed":  {"lifesteal_pct", 1, 0.003},
	"ancient": {"vit", 1, 0.06}, // weapons and armor
	// armor prefixes
	"reinforced": {"defense", 1, 0.05},
	"gilded":     {"hp_max", 5, 0.4},
	"shadowed":   {"evasion_pct", 1, 0.004},
	"runic":      {"mp_max", 5, 0.3},
	"warded":     {"dmg_reduction_pct", 1, 0.003},
	// weapon suffixes
	"of_the_bear":  {"str", 1, 0.08},
	"of_precision": {"hit_pct", 1, 0.005},
	"of_the_wolf":  {"agi", 1, 0.08},
	"of_ruin":      {"crit_dmg_pct", 2, 0.02},
	"of_haste":     {"crit_pct", 1, 0.004},
	// armor suffixes
	"of_the_turtle": {"defense", 1, 0.04},
	"of_wisdom":     {"int", 1, 0.08},
	"of_the_ox":     {"vit", 1, 0.08},
	"of_shielding":  {"block_pct", 1, 0.004},
	"of_endurance":  {"hp_max", 5, 0.3},
}

// StatBonus is one flat stat an item adds (target, amount).
type StatBonus struct {
	Target string
	Amount float64
}

// AffixBonuses is what an item's prefix and suffix add at its level
// (DD-275). Unknown codes add nothing.
func AffixBonuses(prefix, suffix *string, level int16) []StatBonus {
	var out []StatBonus
	for _, code := range []*string{prefix, suffix} {
		if code == nil {
			continue
		}
		if a, ok := affixStats[*code]; ok {
			out = append(out, StatBonus{Target: a.target,
				Amount: max(1, math.Round(a.base+a.perLevel*float64(max(level, 1))))})
		}
	}
	return out
}

// rollAffix picks a prefix + suffix based on rarity rules using the provided
// deterministic RNG. Returns (prefix, suffix) — either may be nil.
func rollAffix(rng *activity.RNG, tmpl Template, rarity Rarity) (*string, *string) {
	if !tmpl.Affixable {
		return nil, nil
	}
	rule := AffixRuleFor(rarity)

	var prefix, suffix *string

	pool := prefixPool[tmpl.ItemClass]
	sufP := suffixPool[tmpl.ItemClass]

	if rule.HasPrefix || (rule.OptionalPrefix && rng.Float64() < 0.5) {
		if len(pool) > 0 {
			p := pool[rng.IntN(len(pool))]
			prefix = &p
		}
	}
	if rule.HasSuffix || (rule.OptionalSuffix && rng.Float64() < 0.5) {
		if len(sufP) > 0 {
			s := sufP[rng.IntN(len(sufP))]
			suffix = &s
		}
	}
	return prefix, suffix
}
