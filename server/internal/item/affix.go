package item

import "github.com/otter/server/internal/activity"

// Affix v1: cosmetic only. Stat modifier lands with Sprint 11 skill engine.
// Content-driven pool in Sprint 12.

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
