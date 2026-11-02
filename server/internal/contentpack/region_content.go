package contentpack

import (
	"sort"
	"strings"
)

// RegionContent is everything found in one region (DD-338), for the map's
// region card: the monsters hunted there, the dungeons and raids entered
// there, the resources that grow there and the quests sent there.
type RegionContent struct {
	Mobs     []RegionMob     `json:"mobs"`
	Dungeons []RegionDungeon `json:"dungeons"`
	Gather   []RegionNode    `json:"gather"`
	Quests   []RegionQuest   `json:"quests"`
}

type RegionMob struct {
	Code  string `json:"code"`
	Level int16  `json:"level"`
}

type RegionDungeon struct {
	Code     string `json:"code"`
	MinLevel int16  `json:"min_level"`
	Raid     bool   `json:"raid"`
}

type RegionNode struct {
	Code  string `json:"code"`
	Kind  string `json:"kind"`
	Level int16  `json:"level"`
}

type RegionQuest struct {
	Code     string `json:"code"`
	Category string `json:"category"`
	MinLevel int16  `json:"min_level"`
}

// BuildRegionContent indexes p by region. A quest belongs to the region of
// the first monster or dungeon its objectives name; quests about gathering
// or crafting alone have no region and are left out. A dungeon-kind region
// is an entrance on the map: it shows the dungeons and raids of the zone it
// opens from, where they are entered.
func BuildRegionContent(p *Pack) map[string]*RegionContent {
	out := make(map[string]*RegionContent, len(p.Regions))
	at := func(code string) *RegionContent {
		if out[code] == nil {
			out[code] = &RegionContent{Mobs: []RegionMob{}, Dungeons: []RegionDungeon{}, Gather: []RegionNode{}, Quests: []RegionQuest{}}
		}
		return out[code]
	}
	for _, r := range p.Regions {
		at(r.Code)
	}
	mobRegion := make(map[string]string, len(p.Mobs))
	for _, m := range p.Mobs {
		mobRegion[m.Code] = m.RegionCode
		if m.RegionCode != "" && !m.RaidBoss { // raid bosses show as their raid
			c := at(m.RegionCode)
			c.Mobs = append(c.Mobs, RegionMob{Code: m.Code, Level: m.Level})
		}
	}
	dungeonRegion := make(map[string]string, len(p.Dungeons))
	for _, d := range p.Dungeons {
		dungeonRegion[d.Code] = d.RegionCode
		if d.RegionCode != "" {
			c := at(d.RegionCode)
			c.Dungeons = append(c.Dungeons, RegionDungeon{Code: d.Code, MinLevel: d.MinLevel, Raid: d.GroupMax > 1})
		}
	}
	for _, it := range p.Items {
		if g := it.Gather; g != nil && g.Region != "" {
			c := at(g.Region)
			c.Gather = append(c.Gather, RegionNode{Code: it.Code, Kind: g.Kind, Level: g.Level})
		}
	}
	for _, q := range p.Quests {
		region := ""
		for _, o := range q.Objectives {
			kind, target, _ := strings.Cut(o.Event, ":")
			switch kind {
			case "mob_killed":
				region = mobRegion[target]
			case "dungeon_cleared":
				region = dungeonRegion[target]
			}
			if region != "" {
				break
			}
		}
		if region != "" {
			c := at(region)
			c.Quests = append(c.Quests, RegionQuest{Code: q.Code, Category: q.Category, MinLevel: q.MinLevel})
		}
	}
	for _, reg := range p.Regions {
		if reg.Kind == "dungeon" && reg.ParentCode != "" {
			at(reg.Code).Dungeons = append(at(reg.Code).Dungeons, at(reg.ParentCode).Dungeons...)
		}
	}
	category := map[string]int{"main": 0, "side": 1, "daily": 2}
	for _, c := range out {
		sort.SliceStable(c.Mobs, func(i, j int) bool { return c.Mobs[i].Level < c.Mobs[j].Level })
		sort.SliceStable(c.Dungeons, func(i, j int) bool { return c.Dungeons[i].MinLevel < c.Dungeons[j].MinLevel })
		sort.SliceStable(c.Gather, func(i, j int) bool {
			if c.Gather[i].Kind != c.Gather[j].Kind {
				return c.Gather[i].Kind < c.Gather[j].Kind
			}
			return c.Gather[i].Level < c.Gather[j].Level
		})
		sort.SliceStable(c.Quests, func(i, j int) bool {
			if a, b := category[c.Quests[i].Category], category[c.Quests[j].Category]; a != b {
				return a < b
			}
			return c.Quests[i].MinLevel < c.Quests[j].MinLevel
		})
	}
	return out
}
