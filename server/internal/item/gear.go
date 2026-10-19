package item

import (
	"sort"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/itemset"
	"github.com/otter/server/internal/leveling"
)

// GearPool lists, sorted, the equipment a class can use from the highest
// gear tier a source level reaches (DD-229, DD-254): a level-12 mob drops the
// level-10 tier, a level-45 mob the level-35 tier. Only the weapon families
// and armor type the class may equip drop (DD-271). Items with extra
// requirements (branch or attribute gated) are never random drops.
func GearPool(cat *Catalog, class string, level int) []string {
	perms, sets := leveling.Current(), itemset.Current()
	usable := func(t Template) bool {
		if _, inSet := sets.Of(t.Code); inSet {
			return false // set pieces come from their raid only (DD-277)
		}
		return t.IsEquippable() && t.Requirements == nil && int(t.MinLevel) <= level &&
			(t.ClassRestriction == nil || *t.ClassRestriction == "" || *t.ClassRestriction == class) &&
			perms.CanUse(class, t.ItemSubclass)
	}
	top := 0
	for _, t := range cat.All() {
		if usable(t) && int(t.MinLevel) > top {
			top = int(t.MinLevel)
		}
	}
	var out []string
	for _, t := range cat.All() {
		if usable(t) && int(t.MinLevel) == top {
			out = append(out, t.Code)
		}
	}
	sort.Strings(out)
	return out
}

// rarityOdds is one row of a drop's rarity table (percent).
type rarityOdds struct {
	r   Rarity
	pct int
}

// Gear rarity odds (percent, cumulative order common → legendary).
var (
	gearOdds     = []rarityOdds{{RarityCommon, 60}, {RarityUncommon, 25}, {RarityRare, 11}, {RarityEpic, 4}}
	bossGearOdds = []rarityOdds{{RarityUncommon, 50}, {RarityRare, 35}, {RarityEpic, 13}, {RarityLegendary, 2}}
)

// RollGear resolves a GearDrop into drops: the guaranteed pieces plus one
// Pct% chance per roll. Pure and deterministic for (seed, g); its RNG stream
// is separate from the material roll's.
func RollGear(cat *Catalog, seed int64, g activity.GearDrop) []DropSpec {
	rng := activity.NewRNG(seed + 4001)
	out := rollPool(cat, rng, g)
	if piece := setPiece(cat, rng, g); piece != "" {
		out = append(out, rollDrop(cat, rng, piece, bossGearOdds))
	}
	return out
}

// setPiece is one piece of the raid's set the class can wear (DD-277): the
// set whose armor the class may equip. Empty when the drop has no set.
func setPiece(cat *Catalog, rng *activity.RNG, g activity.GearDrop) string {
	if !g.Boss || g.SetSource == "" {
		return ""
	}
	perms := leveling.Current()
	for _, s := range itemset.Current().FromSource(g.SetSource) {
		wearable := true
		for _, code := range s.Pieces {
			if t, err := cat.Get(code); err != nil || !perms.CanUse(g.Class, t.ItemSubclass) {
				wearable = false
				break
			}
		}
		if wearable {
			return s.Pieces[rng.IntN(len(s.Pieces))]
		}
	}
	return ""
}

func rollPool(cat *Catalog, rng *activity.RNG, g activity.GearDrop) []DropSpec {
	pool := GearPool(cat, g.Class, g.Level)
	if len(pool) == 0 {
		return nil
	}
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
		out = append(out, rollDrop(cat, rng, pool[rng.IntN(len(pool))], odds))
	}
	return out
}

// rollDrop rolls a template's rarity on the odds table, then its affixes.
func rollDrop(cat *Catalog, rng *activity.RNG, code string, odds []rarityOdds) DropSpec {
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
	return DropSpec{TemplateCode: code, Rarity: rarity, Prefix: prefix, Suffix: suffix}
}
