package contentpack

import (
	"slices"
	"sort"

	"github.com/otter/server/internal/leveling"
)

// ItemInfo is what the item tags show (DD-339): who can use an item and
// from what level, and every way to get it or use it.
type ItemInfo struct {
	Level int16  `json:"level"`
	Slot  string `json:"slot,omitempty"`
	// Classes that may equip it (base classes, or the branches a piece is
	// limited to); empty with AnyClass false = not equipment.
	Classes  []string `json:"classes,omitempty"`
	AnyClass bool     `json:"any_class,omitempty"`
	// Mobs drop it on hunts (gear of their tier, hides they skin).
	Mobs []ItemMob `json:"mobs,omitempty"`
	// Dungeons and raids whose boss drops it.
	Dungeons []ItemDungeon  `json:"dungeons,omitempty"`
	Crafted  []ItemCraft    `json:"crafted,omitempty"`
	Vendors  []string       `json:"vendors,omitempty"`
	Gather   *ItemGatherSrc `json:"gather,omitempty"`
	// UsedBy: the professions whose recipes take it as a material.
	UsedBy []string `json:"used_by,omitempty"`
	Set    string   `json:"set,omitempty"`
}

type ItemMob struct {
	Code  string `json:"code"`
	Level int16  `json:"level"`
}

type ItemDungeon struct {
	Code string `json:"code"`
	Raid bool   `json:"raid"`
}

type ItemCraft struct {
	Profession string `json:"profession"`
	Skill      int16  `json:"skill"`
}

type ItemGatherSrc struct {
	Kind   string `json:"kind"`
	Level  int16  `json:"level"`
	Region string `json:"region,omitempty"`
}

// maxItemMobs caps the monster list of one item (a gear tier spans a few).
const maxItemMobs = 8

// BuildItemInfo indexes every item's tags from the content pack, mirroring
// the drop rules: random gear comes from the highest gear tier at or below
// the source's level (item.GearPool), hides from the monsters whose skinning
// reach covers them (item.CombatHides), raid set pieces and pet eggs from
// raid bosses. cfg supplies the class equipment permissions.
func BuildItemInfo(p *Pack, cfg *leveling.Config) map[string]*ItemInfo {
	out := make(map[string]*ItemInfo, len(p.Items))
	setOf := map[string]string{}
	setSource := map[string]string{}
	for _, s := range p.ItemSets {
		for _, piece := range s.Pieces {
			setOf[piece], setSource[piece] = s.Code, s.Source
		}
	}
	mobLevel := make(map[string]int16, len(p.Mobs))
	for _, m := range p.Mobs {
		mobLevel[m.Code] = m.Level
	}

	// Random gear tiers: equipment with no extra requirement and no set.
	droppable := func(it Item) bool { return it.EquipSlot != "" && it.Requirements == nil && setOf[it.Code] == "" }
	tierSet := map[int16]bool{}
	for _, it := range p.Items {
		if droppable(it) {
			tierSet[it.MinLevel] = true
		}
	}
	tiers := make([]int16, 0, len(tierSet))
	for t := range tierSet {
		tiers = append(tiers, t)
	}
	sort.Slice(tiers, func(i, j int) bool { return tiers[i] < tiers[j] })
	tierAt := func(level int16) int16 { // highest tier ≤ level; -1 = none
		t := int16(-1)
		for _, x := range tiers {
			if x <= level {
				t = x
			}
		}
		return t
	}

	// Hunting nodes by level, for the skinning reach of each mob.
	var hides []Item
	for _, it := range p.Items {
		if it.Gather != nil && it.Gather.Kind == "hunting" {
			hides = append(hides, it)
		}
	}
	sort.SliceStable(hides, func(i, j int) bool { return hides[i].Gather.Level < hides[j].Gather.Level })
	skinned := func(level int16) []string { // item.CombatHides' span
		reach := min(int16(100), max(int16(1), (level+9)/10))
		var unlocked []string
		for _, h := range hides {
			if h.Gather.Level <= reach {
				unlocked = append(unlocked, h.Code)
			}
		}
		return unlocked[max(0, len(unlocked)-4):]
	}

	for _, it := range p.Items {
		info := &ItemInfo{Level: it.MinLevel, Slot: it.EquipSlot, Set: setOf[it.Code]}
		if it.EquipSlot != "" {
			info.Classes, info.AnyClass = usableBy(it, cfg)
		}
		if g := it.Gather; g != nil {
			info.Gather = &ItemGatherSrc{Kind: g.Kind, Level: g.Level, Region: g.Region}
		}
		out[it.Code] = info
	}

	addMob := func(code string, m Mob) {
		if info := out[code]; info != nil && len(info.Mobs) < maxItemMobs {
			info.Mobs = append(info.Mobs, ItemMob{Code: m.Code, Level: m.Level})
		}
	}
	gearOfTier := map[int16][]string{}
	for _, it := range p.Items {
		if droppable(it) {
			gearOfTier[it.MinLevel] = append(gearOfTier[it.MinLevel], it.Code)
		}
	}
	mobs := append([]Mob(nil), p.Mobs...)
	sort.SliceStable(mobs, func(i, j int) bool { return mobs[i].Level < mobs[j].Level })
	for _, m := range mobs {
		if m.RaidBoss {
			continue
		}
		for _, code := range gearOfTier[tierAt(m.Level)] {
			addMob(code, m)
		}
		for _, code := range skinned(m.Level) {
			addMob(code, m)
		}
	}

	addDungeon := func(code string, d Dungeon) {
		info := out[code]
		if info == nil {
			return
		}
		for _, x := range info.Dungeons {
			if x.Code == d.Code {
				return
			}
		}
		info.Dungeons = append(info.Dungeons, ItemDungeon{Code: d.Code, Raid: d.GroupMax > 1})
	}
	for _, d := range p.Dungeons {
		boss := mobLevel[d.BossMobCode]
		for _, code := range gearOfTier[tierAt(boss)] {
			addDungeon(code, d)
		}
		if d.GroupMax <= 1 {
			continue
		}
		for piece, src := range setSource {
			if src == d.Code {
				addDungeon(piece, d)
			}
		}
		for _, s := range p.Pets {
			if int(boss) >= s.MinMobLevel {
				addDungeon("pet_egg_"+s.Code, d)
			}
		}
	}

	for _, r := range p.Recipes {
		if info := out[r.OutputTemplateCode]; info != nil {
			info.Crafted = append(info.Crafted, ItemCraft{Profession: r.ProfessionCode, Skill: r.MinSkill})
		}
		for _, in := range r.Inputs {
			if info := out[in.TemplateCode]; info != nil && !slices.Contains(info.UsedBy, r.ProfessionCode) {
				info.UsedBy = append(info.UsedBy, r.ProfessionCode)
			}
		}
	}
	for _, v := range p.Vendors {
		for _, s := range v.Sells {
			if info := out[s.TemplateCode]; info != nil && !slices.Contains(info.Vendors, v.Code) {
				info.Vendors = append(info.Vendors, v.Code)
			}
		}
	}
	for _, info := range out {
		sort.Strings(info.UsedBy)
	}
	return out
}

// usableBy lists who may equip it: the branch it requires, the class it is
// restricted to, or the classes whose equipment permissions allow its family,
// each named at the highest node of its tree that gains it (a base class, or
// a promotion for a family only that branch uses). Jewelry: anyone.
func usableBy(it Item, cfg *leveling.Config) ([]string, bool) {
	if it.Requirements != nil && it.Requirements.Spec != "" {
		return []string{it.Requirements.Spec}, false
	}
	if it.ClassRestriction != "" {
		return []string{it.ClassRestriction}, false
	}
	if cfg == nil || !cfg.IsFamily(it.ItemSubclass) {
		return nil, true
	}
	var out []string
	for _, k := range cfg.Classes {
		if cfg.CanUse(k.Code, it.ItemSubclass) && (k.Parent == "" || !cfg.CanUse(k.Parent, it.ItemSubclass)) {
			out = append(out, k.Code)
		}
	}
	return out, false
}
