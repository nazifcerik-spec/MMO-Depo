package item

import (
	"sort"

	"github.com/otter/server/internal/activity"
)

// GearTierSpan is how many levels a gear tier spans below its top: a drop
// comes from the highest tier the source level reaches, min levels in
// (top-5, top] — a level-12 mob drops the level-10 tier, a level-45 mob the
// level-35 tier once tiers grow sparser (DD-254).
const GearTierSpan = 5

// GearPool lists, sorted, the equipment a class can use from the tier at
// level (DD-229).
func GearPool(cat *Catalog, class string, level int) []string {
	usable := func(t Template) bool {
		return t.IsEquippable() && int(t.MinLevel) <= level &&
			(t.ClassRestriction == nil || *t.ClassRestriction == "" || *t.ClassRestriction == class)
	}
	top := 0
	for _, t := range cat.All() {
		if usable(t) && int(t.MinLevel) > top {
			top = int(t.MinLevel)
		}
	}
	var out []string
	for _, t := range cat.All() {
		if usable(t) && int(t.MinLevel) > top-GearTierSpan {
			out = append(out, t.Code)
		}
	}
	sort.Strings(out)
	return out
}

// Gear rarity odds (percent, cumulative order common → legendary).
var (
	gearOdds = []struct {
		r   Rarity
		pct int
	}{{RarityCommon, 60}, {RarityUncommon, 25}, {RarityRare, 11}, {RarityEpic, 4}}
	bossGearOdds = []struct {
		r   Rarity
		pct int
	}{{RarityUncommon, 50}, {RarityRare, 35}, {RarityEpic, 13}, {RarityLegendary, 2}}
)

// RollGear resolves a GearDrop into drops: the guaranteed pieces plus one
// Pct% chance per roll. Pure and deterministic for (seed, g); its RNG stream
// is separate from the material roll's.
func RollGear(cat *Catalog, seed int64, g activity.GearDrop) []DropSpec {
	pool := GearPool(cat, g.Class, g.Level)
	if len(pool) == 0 {
		return nil
	}
	rng := activity.NewRNG(seed + 4001)
	n := g.Guaranteed
	for i := 0; i < g.Rolls; i++ {
		if rng.IntN(100) < g.Pct {
			n++
		}
	}
	odds := gearOdds
	if g.Boss {
		odds = bossGearOdds
	}
	out := make([]DropSpec, 0, n)
	for i := 0; i < n; i++ {
		code := pool[rng.IntN(len(pool))]
		roll, rarity := rng.IntN(100), odds[len(odds)-1].r
		for _, o := range odds {
			if roll < o.pct {
				rarity = o.r
				break
			}
			roll -= o.pct
		}
		tmpl, _ := cat.Get(code)
		prefix, suffix := rollAffix(rng, tmpl, rarity)
		out = append(out, DropSpec{TemplateCode: code, Rarity: rarity, Prefix: prefix, Suffix: suffix})
	}
	return out
}
