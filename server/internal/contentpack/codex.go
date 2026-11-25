package contentpack

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/pet"
)

// The codex (DD-372): one page per piece of content with everything that is
// connected to it, so every tag in the client can open what it names.

// CodexRef points at a page. Level is the target's level, Skill the
// profession level a craft or a node needs, Count a quantity; Sub narrows
// the kind for the client's icon (raid, item class, region kind).
type CodexRef struct {
	Kind  string `json:"kind"`
	Code  string `json:"code"`
	Sub   string `json:"sub,omitempty"`
	Level int    `json:"level,omitempty"`
	Skill int    `json:"skill,omitempty"`
	Count int    `json:"count,omitempty"`
}

// CodexGroup is one relation, e.g. "drops" → items.
type CodexGroup struct {
	Key  string     `json:"key"`
	Refs []CodexRef `json:"refs"`
}

// CodexGoal is something to do: a progression event and how often, or
// "level" / "class_mastery:<spec>" with the character level as the count.
type CodexGoal struct {
	Event string `json:"event"`
	Count int64  `json:"count"`
}

// CodexModifier is a bonus; At is the piece count, tier or level it needs.
type CodexModifier struct {
	Kind   string  `json:"kind"`
	Target string  `json:"target"`
	Amount float64 `json:"amount"`
	At     int     `json:"at,omitempty"`
}

// CodexPage is everything about one piece of content.
type CodexPage struct {
	Kind  string `json:"kind"`
	Code  string `json:"code"`
	Sub   string `json:"sub,omitempty"`
	Level int    `json:"level,omitempty"`
	// Facts are numbers the client labels (hp, attack, gold, seconds...);
	// Traits are codes it names (slot, subclass, category, element).
	Facts     map[string]int64  `json:"facts,omitempty"`
	Traits    map[string]string `json:"traits,omitempty"`
	Goals     []CodexGoal       `json:"goals,omitempty"`
	Modifiers []CodexModifier   `json:"modifiers,omitempty"`
	Groups    []CodexGroup      `json:"groups,omitempty"`
}

// CodexSeeds is the progression content that lives in the database
// (migrations), not in the pack.
type CodexSeeds struct {
	Achievements []SeedAchievement
	Factions     []SeedSource
	Tracks       []SeedTrack
	Events       []SeedEvent
}

type SeedAchievement struct {
	Code, Category, Trigger, Title string
	Points                         int16
	Target                         int32
	Gold                           int64
}

type SeedSource struct {
	Code   string
	Events []string
}

type SeedTrack struct {
	SeedSource
	MaxTier int16
	Bonuses []CodexModifier
}

type SeedEvent struct {
	Code, Kind, Region string
	MinLevel           int16
}

// sourceEvents reads the event kinds of a [{kind, ...}] column.
func sourceEvents(raw []byte) ([]string, error) {
	var src []struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(raw, &src); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(src))
	for _, x := range src {
		out = append(out, x.Kind)
	}
	return out, nil
}

func (l *Loader) codexSeeds(ctx context.Context) (CodexSeeds, error) {
	var s CodexSeeds
	rows, err := l.DB.Query(ctx, `
		SELECT code, category, trigger_kind, COALESCE(reward_title_code, ''), point_value, progress_target, reward_gold_copper
		  FROM content.achievement_template WHERE NOT is_hidden ORDER BY code`)
	if err != nil {
		return s, err
	}
	s.Achievements, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (a SeedAchievement, err error) {
		err = r.Scan(&a.Code, &a.Category, &a.Trigger, &a.Title, &a.Points, &a.Target, &a.Gold)
		return a, err
	})
	if err != nil {
		return s, err
	}
	if rows, err = l.DB.Query(ctx, `SELECT code, sources FROM content.faction ORDER BY code`); err != nil {
		return s, err
	}
	s.Factions, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (f SeedSource, err error) {
		var raw []byte
		if err = r.Scan(&f.Code, &raw); err == nil {
			f.Events, err = sourceEvents(raw)
		}
		return f, err
	})
	if err != nil {
		return s, fmt.Errorf("factions: %w", err)
	}
	if rows, err = l.DB.Query(ctx, `SELECT code, max_tier, sources, tier_bonuses FROM content.mastery_track ORDER BY code`); err != nil {
		return s, err
	}
	s.Tracks, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (t SeedTrack, err error) {
		var raw, bonuses []byte
		if err = r.Scan(&t.Code, &t.MaxTier, &raw, &bonuses); err != nil {
			return t, err
		}
		if t.Events, err = sourceEvents(raw); err != nil {
			return t, err
		}
		var tiers []struct {
			Tier    int             `json:"tier"`
			Effects []CodexModifier `json:"bonus_effects"`
		}
		if err = json.Unmarshal(bonuses, &tiers); err != nil {
			return t, err
		}
		for _, b := range tiers {
			for _, e := range b.Effects {
				e.At = b.Tier
				t.Bonuses = append(t.Bonuses, e)
			}
		}
		return t, nil
	})
	if err != nil {
		return s, fmt.Errorf("mastery tracks: %w", err)
	}
	if rows, err = l.DB.Query(ctx, `
		SELECT code, kind, COALESCE(region_code, ''), min_level FROM content.world_event_template ORDER BY code`); err != nil {
		return s, err
	}
	s.Events, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (e SeedEvent, err error) {
		err = r.Scan(&e.Code, &e.Kind, &e.Region, &e.MinLevel)
		return e, err
	})
	return s, err
}

// Codex is the page index by kind and code.
type Codex map[string]map[string]*CodexPage

// Page returns one page, nil when unknown.
func (c Codex) Page(kind, code string) *CodexPage { return c[kind][code] }

func (c Codex) page(kind, code string) *CodexPage {
	if c[kind] == nil {
		c[kind] = map[string]*CodexPage{}
	}
	p := c[kind][code]
	if p == nil {
		p = &CodexPage{Kind: kind, Code: code}
		c[kind][code] = p
	}
	return p
}

func (p *CodexPage) fact(key string, v int64) {
	if v == 0 {
		return
	}
	if p.Facts == nil {
		p.Facts = map[string]int64{}
	}
	p.Facts[key] = v
}

func (p *CodexPage) trait(key, v string) {
	if v == "" {
		return
	}
	if p.Traits == nil {
		p.Traits = map[string]string{}
	}
	p.Traits[key] = v
}

// add puts ref under a relation once; a page never lists itself.
func (p *CodexPage) add(key string, ref CodexRef) {
	if key == "" || (ref.Kind == p.Kind && ref.Code == p.Code) {
		return
	}
	for i := range p.Groups {
		g := &p.Groups[i]
		if g.Key != key {
			continue
		}
		for _, r := range g.Refs {
			if r.Kind == ref.Kind && r.Code == ref.Code {
				return
			}
		}
		g.Refs = append(g.Refs, ref)
		return
	}
	p.Groups = append(p.Groups, CodexGroup{Key: key, Refs: []CodexRef{ref}})
}

// link connects two pages both ways: a lists b under key, b lists a under
// inverse (an empty side is left out). A missing page links nothing.
func (c Codex) link(a *CodexPage, key string, b *CodexPage, inverse string, skill, count int) {
	if a == nil || b == nil {
		return
	}
	a.add(key, CodexRef{Kind: b.Kind, Code: b.Code, Skill: skill, Count: count})
	b.add(inverse, CodexRef{Kind: a.Kind, Code: a.Code, Skill: skill, Count: count})
}

// eventTargets: what a scoped event ("mob_killed:grey_wolf") names, and the
// list page of everything that counts for the bare one.
var eventTargets = map[string]struct{ kind, list string }{
	"mob_killed":      {"mob", "mobs"},
	"dungeon_cleared": {"dungeon", "dungeons"},
	"raid_cleared":    {"dungeon", "dungeons"},
	"boss_killed":     {"", "dungeons"},
	"item_crafted":    {"item", "crafts"},
	"class_mastery":   {"class", ""},
}

// eventPage is the page an event names: the monster, dungeon, item or class
// it is scoped to, the profession it counts, or the list of everything that
// counts ("any monster").
func (c Codex) eventPage(event string) *CodexPage {
	kind, code, scoped := strings.Cut(event, ":")
	if g, ok := strings.CutPrefix(kind, "gather_"); ok {
		return c.Page("profession", g)
	}
	t := eventTargets[kind]
	if scoped {
		return c.Page(t.kind, code)
	}
	return c.Page("list", t.list)
}

// goal adds something to do to p and links what it names.
func (c Codex) goal(p *CodexPage, event string, count int64, inverse string) {
	p.Goals = append(p.Goals, CodexGoal{Event: event, Count: count})
	if t := c.eventPage(event); t != nil {
		t.add(inverse, CodexRef{Kind: p.Kind, Code: p.Code})
	}
}

// BuildCodex indexes the whole pack. Item sources mirror BuildItemInfo, the
// region lists BuildRegionContent. Pet eggs name their raids but a raid does
// not name its eggs: the pet index keeps undiscovered species hidden (DD-342).
func BuildCodex(p *Pack, cfg *leveling.Config, seeds CodexSeeds) Codex {
	c := Codex{}
	c.addContent(p)
	c.addProgression(cfg, seeds)
	c.linkItems(p, cfg)
	c.linkWorld(p)
	c.linkProgression(p, cfg, seeds)
	c.finish()
	return c
}

// addContent creates the pack's pages with their own numbers, so links find
// both ends.
func (c Codex) addContent(p *Pack) {
	for _, it := range p.Items {
		pg := c.page("item", it.Code)
		pg.Sub, pg.Level = it.ItemClass, int(it.MinLevel)
		pg.trait("slot", it.EquipSlot)
		pg.trait("subclass", it.ItemSubclass)
		pg.fact("weight_g", int64(it.BaseWeightG))
		pg.fact("value", it.BaseValueCopper)
		for _, b := range it.StatBonuses {
			pg.Modifiers = append(pg.Modifiers, CodexModifier{Kind: b.Kind, Target: b.Target, Amount: b.Amount})
		}
		if it.Enchant != nil {
			for _, b := range it.Enchant.Modifiers {
				pg.Modifiers = append(pg.Modifiers, CodexModifier{Kind: b.Kind, Target: b.Target, Amount: b.Amount})
			}
		}
		if u := it.UseEffect; u != nil {
			pg.fact(u.Kind, int64(u.Amount))
			if u.Buff != nil {
				pg.fact("buff_attack_pct", int64(u.Buff.AttackPct))
				pg.fact("buff_defense_pct", int64(u.Buff.DefensePct))
				pg.fact("buff_minutes", int64(u.Buff.Minutes))
			}
		}
	}
	for _, m := range p.Mobs {
		pg := c.page("mob", m.Code)
		pg.Level = int(m.Level)
		if m.RaidBoss {
			pg.Sub = "boss"
		}
		e := combat.Mob{Level: m.Level, HPMax: m.HPMax, Attack: m.Attack, Defense: m.Defense}.Effective()
		pg.fact("hp", int64(e.HPMax))
		pg.fact("attack", int64(e.Attack))
		pg.fact("defense", int64(e.Defense))
		if m.Hunted() {
			pg.fact("xp", int64(m.XPReward))
		}
	}
	for _, d := range p.Dungeons {
		pg := c.page("dungeon", d.Code)
		pg.Level = int(d.MinLevel)
		pg.fact("seconds", int64(d.DurationSeconds))
		pg.fact("gold", d.BonusGoldCopper)
		if d.GroupMax > 1 {
			pg.Sub = "raid"
			pg.fact("group_min", int64(d.GroupMin))
			pg.fact("group_max", int64(d.GroupMax))
		}
	}
	for _, r := range p.Regions {
		pg := c.page("region", r.Code)
		pg.Sub, pg.Level = r.Kind, int(r.MinLevel)
	}
	for _, x := range p.Professions {
		pg := c.page("profession", x.Code)
		pg.Sub = x.Kind
		pg.fact("max_level", int64(x.MaxLevel))
	}
	for _, v := range p.Vendors {
		c.page("vendor", v.Code)
	}
	for _, q := range p.Quests {
		pg := c.page("quest", q.Code)
		pg.Sub, pg.Level = q.Category, int(q.MinLevel)
		if q.Category != "daily" && q.Category != "side" { // theirs scale with the character (DD-318)
			pg.fact("xp", int64(q.RewardXP))
			pg.fact("gold", q.RewardGoldCopper)
		}
	}
	for _, s := range p.ItemSets {
		pg := c.page("itemset", s.Code)
		for _, b := range s.Bonuses {
			for _, m := range b.Modifiers {
				pg.Modifiers = append(pg.Modifiers, CodexModifier{Kind: m.Kind, Target: m.Target, Amount: m.Amount, At: b.Pieces})
			}
		}
	}
	for _, s := range p.Pets {
		pg := c.page("pet", s.Code)
		pg.trait("element", s.Element)
		pg.fact("rank", int64(s.Rank))
		pg.fact("seconds", int64(s.HatchSeconds))
		pg.Modifiers = []CodexModifier{{Kind: s.Bonus.Kind, Target: s.Bonus.Target, Amount: s.Bonus.Amount, At: pet.MaxLevel}}
	}
	for _, s := range p.Skills {
		pg := c.page("skill", s.Code)
		pg.Sub, pg.Level = s.category(), int(s.UnlockLevel)
		for _, m := range s.Modifiers {
			pg.Modifiers = append(pg.Modifiers, CodexModifier{Kind: m.Kind, Target: m.Target, Amount: m.Amount})
		}
	}
	for _, code := range []string{"mobs", "dungeons", "crafts"} {
		c.page("list", code)
	}
}

// addProgression creates the pages of classes, races and the seeds.
func (c Codex) addProgression(cfg *leveling.Config, seeds CodexSeeds) {
	for _, a := range seeds.Achievements {
		pg := c.page("achievement", a.Code)
		pg.trait("category", a.Category)
		pg.fact("points", int64(a.Points))
		pg.fact("gold", a.Gold)
	}
	for _, f := range seeds.Factions {
		c.page("faction", f.Code)
	}
	for _, t := range seeds.Tracks {
		pg := c.page("track", t.Code)
		pg.Modifiers = t.Bonuses
		pg.fact("max_tier", int64(t.MaxTier))
	}
	for _, e := range seeds.Events {
		pg := c.page("worldevent", e.Code)
		pg.Sub, pg.Level = e.Kind, int(e.MinLevel)
	}
	if cfg == nil {
		return
	}
	for _, k := range cfg.Classes {
		pg := c.page("class", k.Code)
		pg.Sub, pg.Level = k.Root(), k.RequiredLevel
		for _, f := range k.EquipmentPermissions {
			c.page("equipment", f)
		}
	}
	for code, r := range cfg.Races {
		pg := c.page("race", code)
		for _, m := range r.Passives {
			pg.Modifiers = append(pg.Modifiers, CodexModifier{Kind: m.Kind, Target: m.Target, Amount: m.Amount})
		}
	}
}

func (c Codex) linkItems(p *Pack, cfg *leveling.Config) {
	info := BuildItemInfo(p, cfg)
	for _, it := range p.Items {
		pg, in := c.Page("item", it.Code), info[it.Code]
		for _, k := range in.Classes {
			c.link(pg, "classes", c.Page("class", k), "", 0, 0)
		}
		c.link(pg, "family", c.Page("equipment", it.ItemSubclass), "items", 0, 0)
		c.link(pg, "set", c.Page("itemset", in.Set), "pieces", 0, 0)
		if g := in.Gather; g != nil {
			c.link(pg, "gathered", c.Page("profession", g.Kind), "nodes", int(g.Level), 0)
			c.link(pg, "grows", c.Page("region", g.Region), "nodes", 0, 0)
		}
		for _, m := range in.Mobs {
			c.link(pg, "dropped_by", c.Page("mob", m.Code), "drops", 0, 0)
		}
		inverse := "drops"
		if strings.HasPrefix(it.Code, pet.EggPrefix) {
			inverse = "" // a raid does not name its eggs
		}
		for _, d := range in.Dungeons {
			c.link(pg, "dropped_in", c.Page("dungeon", d.Code), inverse, 0, 0)
		}
		for _, v := range in.Vendors {
			c.link(pg, "sold_by", c.Page("vendor", v), "sells", 0, 0)
		}
	}
	for _, r := range p.Recipes {
		out, prof := c.Page("item", r.OutputTemplateCode), c.Page("profession", r.ProfessionCode)
		c.link(out, "crafted_by", prof, "makes", int(r.MinSkill), 0)
		for _, in := range r.Inputs {
			mat := c.Page("item", in.TemplateCode)
			c.link(out, "materials", mat, "used_in", 0, int(in.Quantity))
			c.link(mat, "used_by", prof, "", 0, 0)
		}
	}
	for _, s := range p.ItemSets {
		c.link(c.Page("itemset", s.Code), "dropped_in", c.Page("dungeon", s.Source), "sets", 0, 0)
	}
	for _, s := range p.Pets {
		pg := c.Page("pet", s.Code)
		c.link(pg, "egg", c.Page("item", s.EggCode()), "hatches", 0, 0)
		if in := info[s.EggCode()]; in != nil {
			for _, d := range in.Dungeons {
				c.link(pg, "dropped_in", c.Page("dungeon", d.Code), "", 0, 0)
			}
		}
	}
}

func (c Codex) linkWorld(p *Pack) {
	mobs, dungeons := c.Page("list", "mobs"), c.Page("list", "dungeons")
	for _, d := range p.Dungeons {
		pg, key := c.Page("dungeon", d.Code), "dungeons"
		if d.GroupMax > 1 {
			key = "raids"
		}
		c.link(pg, "boss", c.Page("mob", d.BossMobCode), "boss_of", 0, 0)
		for _, w := range d.Waves {
			c.link(pg, "enemies", c.Page("mob", w.MobCode), "fought_in", 0, w.Count)
		}
		c.link(pg, "region", c.Page("region", d.RegionCode), key, 0, 0)
		c.link(dungeons, key, pg, "", 0, 0)
	}
	for _, m := range p.Mobs {
		if m.Hunted() {
			c.link(mobs, "mobs", c.Page("mob", m.Code), "", 0, 0)
		}
	}
	for code, rc := range BuildRegionContent(p) {
		for _, m := range rc.Mobs {
			c.link(c.Page("region", code), "mobs", c.Page("mob", m.Code), "hunted_in", 0, 0)
		}
	}
	for _, r := range p.Regions {
		c.link(c.Page("region", r.Code), "within", c.Page("region", r.ParentCode), "contains", 0, 0)
	}
	for _, l := range p.RegionLinks {
		c.link(c.Page("region", l.RegionACode), "roads", c.Page("region", l.RegionBCode), "roads", 0, 0)
	}
	for _, x := range p.Professions {
		if x.Kind != "gathering" {
			c.link(c.Page("list", "crafts"), "professions", c.Page("profession", x.Code), "", 0, 0)
		}
	}
	for _, q := range p.Quests {
		pg := c.Page("quest", q.Code)
		for _, o := range q.Objectives {
			c.goal(pg, o.Event, o.Count, "quests")
		}
		c.link(pg, "requires", c.Page("quest", q.RequiresQuest), "leads_to", 0, 0)
	}
}

func (c Codex) linkProgression(p *Pack, cfg *leveling.Config, seeds CodexSeeds) {
	for _, a := range seeds.Achievements {
		pg, count := c.Page("achievement", a.Code), int64(a.Target)
		if strings.HasPrefix(a.Trigger, "class_mastery:") && cfg != nil {
			count = int64(cfg.Milestones.Mastery) // reached at a level, not "once"
		}
		c.goal(pg, a.Trigger, count, "achievements")
		if a.Title == "" {
			continue
		}
		title := c.page("title", a.Title)
		if len(title.Goals) == 0 {
			title.Goals = slices.Clone(pg.Goals)
			if t := c.eventPage(a.Trigger); t != nil {
				t.add("titles", CodexRef{Kind: "title", Code: a.Title})
			}
		}
		c.link(pg, "titles", title, "achievements", 0, 0)
	}
	for _, f := range seeds.Factions {
		for _, e := range f.Events {
			c.link(c.Page("faction", f.Code), "sources", c.eventPage(e), "reputation", 0, 0)
		}
	}
	for _, t := range seeds.Tracks {
		for _, e := range t.Events {
			c.link(c.Page("track", t.Code), "sources", c.eventPage(e), "mastery", 0, 0)
		}
	}
	for _, e := range seeds.Events {
		c.link(c.Page("worldevent", e.Code), "region", c.Page("region", e.Region), "events", 0, 0)
	}
	if cfg == nil {
		return
	}
	for level, t := range cfg.Titles {
		c.page("title", t).Goals = []CodexGoal{{Event: "level", Count: int64(level)}}
	}
	for _, k := range cfg.Classes {
		pg := c.Page("class", k.Code)
		c.link(pg, "promotes_from", c.Page("class", k.Parent), "branches", 0, 0)
		for _, f := range cfg.Permissions(k.Code) {
			c.link(pg, "equipment", c.Page("equipment", f), "", 0, 0)
		}
		for _, f := range k.EquipmentPermissions {
			// A family names the node that gains it, not every branch below.
			if k.Parent == "" || !cfg.CanUse(k.Parent, f) {
				c.Page("equipment", f).add("classes", CodexRef{Kind: "class", Code: k.Code})
			}
		}
	}
	for _, s := range p.Skills {
		c.link(c.Page("class", s.specCode()), "skills", c.Page("skill", s.Code), "classes", 0, 0)
	}
	for code, r := range cfg.Races {
		pg := c.Page("race", code)
		for _, k := range r.Classes {
			c.link(pg, "classes", c.Page("class", k), "races", 0, 0)
		}
		for prof := range r.ProfessionXP {
			c.link(pg, "professions", c.Page("profession", prof), "races", 0, 0)
		}
	}
}

// finish fills each ref's level and sub from its page and orders everything:
// groups by key, refs by skill, level, then code.
func (c Codex) finish() {
	for _, pages := range c {
		for _, p := range pages {
			for gi := range p.Groups {
				refs := p.Groups[gi].Refs
				for i := range refs {
					t := c.Page(refs[i].Kind, refs[i].Code)
					refs[i].Level, refs[i].Sub = t.Level, t.Sub
				}
				sort.SliceStable(refs, func(i, j int) bool {
					a, b := refs[i], refs[j]
					if a.Skill != b.Skill {
						return a.Skill < b.Skill
					}
					if a.Level != b.Level {
						return a.Level < b.Level
					}
					return a.Code < b.Code
				})
			}
			sort.SliceStable(p.Groups, func(i, j int) bool { return p.Groups[i].Key < p.Groups[j].Key })
		}
	}
}
