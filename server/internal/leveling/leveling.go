// Package leveling is the data-driven character progression engine (DD-250):
// XP curve (formula or table), the level reward registry, class promotion
// tree, talent trees, skill/milestone unlocks and attribute scaling. It holds
// definitions and pure rules only; persistence lives in the character
// domain, which applies them inside one transaction. Definitions come from
// content/progression.yaml (+ skills.yaml unlocks) and are swapped
// atomically on content reload.
package leveling

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync/atomic"
)

// Curve is the XP needed for each level: a formula
// (minutes = a + b·L^c, XP = ceil(minutes × baseline)) or an explicit table
// (index 0 = level 1).
type Curve struct {
	Kind                string  `yaml:"kind" json:"kind"` // formula | table
	BaselineXPPerMinute float64 `yaml:"baseline_xp_per_minute" json:"baseline_xp_per_minute"`
	A                   float64 `yaml:"a" json:"a"`
	B                   float64 `yaml:"b" json:"b"`
	C                   float64 `yaml:"c" json:"c"`
	Table               []int64 `yaml:"table" json:"table,omitempty"`
}

// AttributeRule: Amount points every Every levels, plus milestone bonuses.
type AttributeRule struct {
	Every      int         `yaml:"every" json:"every"`
	Amount     int         `yaml:"amount" json:"amount"`
	Milestones map[int]int `yaml:"milestones" json:"milestones"`
}

// TalentRule: Amount points at From, From+Every, ... up to Until.
type TalentRule struct {
	From   int `yaml:"from" json:"from"`
	Every  int `yaml:"every" json:"every"`
	Until  int `yaml:"until" json:"until"`
	Amount int `yaml:"amount" json:"amount"`
}

// Bracket: raw attribute points up to UpTo count at Effectiveness (0 = no limit).
type Bracket struct {
	UpTo          int     `yaml:"up_to" json:"up_to"`
	Effectiveness float64 `yaml:"effectiveness" json:"effectiveness"`
}

// Modifier is a stat modifier (same shape as skill.Modifier: flat | mult_add | mult_prod).
type Modifier struct {
	Kind   string  `yaml:"kind" json:"kind"`
	Target string  `yaml:"target" json:"target"`
	Amount float64 `yaml:"amount" json:"amount"`
}

// Stats are the six raw attributes.
type Stats struct {
	Strength     int `yaml:"strength" json:"strength"`
	Vitality     int `yaml:"vitality" json:"vitality"`
	Dexterity    int `yaml:"dexterity" json:"dexterity"`
	Intelligence int `yaml:"intelligence" json:"intelligence"`
	Wisdom       int `yaml:"wisdom" json:"wisdom"`
	Luck         int `yaml:"luck" json:"luck"`
}

// Total is the sum of all six.
func (s Stats) Total() int {
	return s.Strength + s.Vitality + s.Dexterity + s.Intelligence + s.Wisdom + s.Luck
}

// Attributes are the allocatable attribute names (API / localization keys).
var Attributes = []string{"strength", "vitality", "dexterity", "intelligence", "wisdom", "luck"}

// Get / Add by attribute name.
func (s Stats) Get(name string) int {
	switch name {
	case "strength":
		return s.Strength
	case "vitality":
		return s.Vitality
	case "dexterity":
		return s.Dexterity
	case "intelligence":
		return s.Intelligence
	case "wisdom":
		return s.Wisdom
	case "luck":
		return s.Luck
	}
	return 0
}

func (s *Stats) Add(name string, n int) bool {
	switch name {
	case "strength":
		s.Strength += n
	case "vitality":
		s.Vitality += n
	case "dexterity":
		s.Dexterity += n
	case "intelligence":
		s.Intelligence += n
	case "wisdom":
		s.Wisdom += n
	case "luck":
		s.Luck += n
	default:
		return false
	}
	return true
}

// Derived is a class's configuration-driven derived-stat formula. Percent
// values are percentage points; caps live in Config.Caps.
type Derived struct {
	HPBase         float64 `yaml:"hp_base" json:"hp_base"`
	HPPerLevel     float64 `yaml:"hp_per_level" json:"hp_per_level"`
	HPPerVit       float64 `yaml:"hp_per_vit" json:"hp_per_vit"`
	AttackBase     float64 `yaml:"attack_base" json:"attack_base"`
	ArmorBase      float64 `yaml:"armor_base" json:"armor_base"`
	AttackPerLevel float64 `yaml:"attack_per_level" json:"attack_per_level"`
	AttackPerStr   float64 `yaml:"attack_per_str" json:"attack_per_str"`
	AttackPerDex   float64 `yaml:"attack_per_dex" json:"attack_per_dex"`
	ArmorPerLevel  float64 `yaml:"armor_per_level" json:"armor_per_level"`
	ArmorPerVit    float64 `yaml:"armor_per_vit" json:"armor_per_vit"`
	CritBase       float64 `yaml:"crit_base" json:"crit_base"`
	CritPerDex     float64 `yaml:"crit_per_dex" json:"crit_per_dex"`
	HitBase        float64 `yaml:"hit_base" json:"hit_base"`
	HitPerDex      float64 `yaml:"hit_per_dex" json:"hit_per_dex"`
	DodgeBase      float64 `yaml:"dodge_base" json:"dodge_base"`
	DodgePerDex    float64 `yaml:"dodge_per_dex" json:"dodge_per_dex"`
	BlockPerStr    float64 `yaml:"block_per_str" json:"block_per_str"`
}

// Resource is a class's combat resource (Warrior: Rage, starts empty).
type Resource struct {
	Kind   string `yaml:"kind" json:"kind"` // rage | mp
	Max    int    `yaml:"max" json:"max"`
	PerHit int    `yaml:"per_hit" json:"per_hit"` // gained on each landed hit
}

// Mastery is the Level-cap class mastery of a final specialization.
type Mastery struct {
	Title     string     `yaml:"title" json:"title"`
	Badge     string     `yaml:"badge" json:"badge"`
	Frame     string     `yaml:"frame" json:"frame"`
	Modifiers []Modifier `yaml:"modifiers" json:"modifiers"`
}

// Class is a node of the class tree: a base class (no parent) or a promotion.
type Class struct {
	Code                 string     `yaml:"code" json:"code"`
	Parent               string     `yaml:"parent" json:"parent,omitempty"`
	RequiredLevel        int        `yaml:"required_level" json:"required_level,omitempty"`
	NameKey              string     `yaml:"name_key" json:"name_key"`
	RoleTags             []string   `yaml:"role_tags" json:"role_tags"`
	Resource             *Resource  `yaml:"resource" json:"resource,omitempty"`
	BaseStats            *Stats     `yaml:"base_stats" json:"base_stats,omitempty"`
	Derived              *Derived   `yaml:"derived" json:"derived,omitempty"`
	StatModifiers        []Modifier `yaml:"stat_modifiers" json:"stat_modifiers,omitempty"`
	EquipmentPermissions []string   `yaml:"equipment_permissions" json:"equipment_permissions,omitempty"`
	Mastery              *Mastery   `yaml:"mastery" json:"mastery,omitempty"`

	tier int
	root string
}

// Tier is the depth in the tree (base class 0, first promotion 1, ...).
func (c Class) Tier() int { return c.tier }

// Root is the base class code.
func (c Class) Root() string { return c.root }

// Talent is one talent node; Modifiers apply once per rank.
type Talent struct {
	Code      string     `yaml:"code" json:"code"`
	MaxRank   int        `yaml:"max_rank" json:"max_rank"`
	Capstone  bool       `yaml:"capstone" json:"capstone,omitempty"`
	Modifiers []Modifier `yaml:"modifiers" json:"modifiers"`
}

// Tree is a class talent tree.
type Tree struct {
	Code    string   `yaml:"code" json:"code"`
	Class   string   `yaml:"class" json:"class"`
	Talents []Talent `yaml:"talents" json:"talents"`
}

// Unlock is a level-gated unlock of a class node (skills, passives,
// awakenings, ultimates). Filled from the class skill content.
type Unlock struct {
	Code     string `json:"code"`
	Spec     string `json:"spec"`
	Level    int    `json:"level"`
	Category string `json:"category"` // skill | passive | awakening | ultimate
}

// Milestones are the post-specialization major levels.
type Milestones struct {
	Awakening int `yaml:"awakening" json:"awakening"`
	Ultimate  int `yaml:"ultimate" json:"ultimate"`
	Capstone  int `yaml:"capstone" json:"capstone"`
	Mastery   int `yaml:"mastery" json:"mastery"`
}

// LevelReward is one explicit registry entry beyond the generated ones.
type LevelReward struct {
	Level  int    `yaml:"level" json:"level"`
	Type   string `yaml:"type" json:"type"`
	Amount int    `yaml:"amount" json:"amount,omitempty"`
	Ref    string `yaml:"ref" json:"ref,omitempty"`
}

// Config is the whole progression definition.
type Config struct {
	LevelCap             int                `yaml:"level_cap" json:"level_cap"`
	Curve                Curve              `yaml:"xp_curve" json:"xp_curve"`
	XPMultiplier         float64            `yaml:"xp_multiplier" json:"xp_multiplier"`
	AFKCombatXPPerMinute float64            `yaml:"afk_combat_xp_per_minute" json:"afk_combat_xp_per_minute"`
	AttributePoints      AttributeRule      `yaml:"attribute_points" json:"attribute_points"`
	TalentPoints         TalentRule         `yaml:"talent_points" json:"talent_points"`
	TalentTreePoints     int                `yaml:"talent_tree_points" json:"talent_tree_points"`
	Titles               map[int]string     `yaml:"titles" json:"titles"`
	Milestones           Milestones         `yaml:"milestones" json:"milestones"`
	Scaling              []Bracket          `yaml:"attribute_scaling" json:"attribute_scaling"`
	Caps                 map[string]float64 `yaml:"caps" json:"caps"`
	LevelRewards         []LevelReward      `yaml:"level_rewards" json:"level_rewards,omitempty"`
	Classes              []Class            `yaml:"classes" json:"classes"`
	Trees                []Tree             `yaml:"talent_trees" json:"talent_trees"`
	Unlocks              []Unlock           `yaml:"-" json:"unlocks"`

	table    []int64
	classes  map[string]*Class
	children map[string][]*Class
	talents  map[string]*Talent
	treeOf   map[string]*Tree
	rewards  map[int][]Reward
}

// Reward kinds (the registry, DD-250 §32) and progression event kinds.
const (
	RewardAttributePoints = "ATTRIBUTE_POINTS"
	RewardTalentPoint     = "TALENT_POINT"
	RewardTitle           = "TITLE"
	RewardPromotion       = "PROMOTION_UNLOCK"
	RewardAwakening       = "AWAKENING"
	RewardUltimate        = "ULTIMATE"
	RewardCapstone        = "CAPSTONE_ELIGIBLE"
	RewardMastery         = "CLASS_MASTERY"
)

// Reward is one registry entry for a level.
type Reward struct {
	Type   string `json:"type"`
	Amount int    `json:"amount,omitempty"`
	Ref    string `json:"ref,omitempty"`
}

var current atomic.Pointer[Config]

// Current is the installed configuration (the built-in default until content loads).
func Current() *Config {
	if c := current.Load(); c != nil {
		return c
	}
	return defaultConfig
}

// Install validates and swaps in a configuration.
func Install(c *Config) error {
	if err := c.Prepare(); err != nil {
		return err
	}
	current.Store(c)
	return nil
}

// Prepare validates the configuration and builds its caches.
func (c *Config) Prepare() error {
	if c.LevelCap < 2 {
		return errors.New("leveling: level_cap must be >= 2")
	}
	if c.XPMultiplier == 0 {
		c.XPMultiplier = 1
	}
	if err := c.buildTable(); err != nil {
		return err
	}
	if c.AttributePoints.Every < 1 || c.TalentPoints.Every < 1 {
		return errors.New("leveling: attribute/talent point 'every' must be >= 1")
	}
	c.classes = map[string]*Class{}
	c.children = map[string][]*Class{}
	for i := range c.Classes {
		k := &c.Classes[i]
		if k.Code == "" {
			return errors.New("leveling: class without code")
		}
		if _, dup := c.classes[k.Code]; dup {
			return fmt.Errorf("leveling: duplicate class %s", k.Code)
		}
		c.classes[k.Code] = k
	}
	for _, k := range c.classes {
		if k.Parent == "" {
			continue
		}
		p, ok := c.classes[k.Parent]
		if !ok {
			return fmt.Errorf("leveling: class %s: unknown parent %s", k.Code, k.Parent)
		}
		if k.RequiredLevel < 2 || k.RequiredLevel > c.LevelCap {
			return fmt.Errorf("leveling: class %s: required_level out of range", k.Code)
		}
		c.children[p.Code] = append(c.children[p.Code], k)
	}
	for _, list := range c.children {
		sort.Slice(list, func(i, j int) bool { return list[i].Code < list[j].Code })
		for _, k := range list {
			if k.RequiredLevel != list[0].RequiredLevel {
				return fmt.Errorf("leveling: promotions of %s must share one required_level", k.Parent)
			}
		}
	}
	for _, k := range c.classes {
		depth, root, seen := 0, k, map[string]bool{}
		for root.Parent != "" {
			if seen[root.Code] {
				return fmt.Errorf("leveling: class cycle at %s", k.Code)
			}
			seen[root.Code] = true
			depth++
			root = c.classes[root.Parent]
		}
		k.tier, k.root = depth, root.Code
		if k.Parent != "" && c.classes[k.Parent].RequiredLevel >= k.RequiredLevel {
			return fmt.Errorf("leveling: class %s must require a higher level than its parent", k.Code)
		}
	}
	c.talents = map[string]*Talent{}
	c.treeOf = map[string]*Tree{}
	for i := range c.Trees {
		t := &c.Trees[i]
		if _, ok := c.classes[t.Class]; !ok {
			return fmt.Errorf("leveling: tree %s: unknown class %s", t.Code, t.Class)
		}
		spendable, capstones := 0, 0
		for j := range t.Talents {
			tl := &t.Talents[j]
			if tl.MaxRank < 1 {
				return fmt.Errorf("leveling: talent %s: max_rank must be >= 1", tl.Code)
			}
			if _, dup := c.talents[tl.Code]; dup {
				return fmt.Errorf("leveling: duplicate talent %s", tl.Code)
			}
			c.talents[tl.Code], c.treeOf[tl.Code] = tl, t
			if tl.Capstone {
				capstones++
			} else {
				spendable += tl.MaxRank
			}
		}
		if spendable != c.TalentTreePoints || capstones != 1 {
			return fmt.Errorf("leveling: tree %s must hold %d spendable points and one capstone", t.Code, c.TalentTreePoints)
		}
	}
	for _, u := range c.Unlocks {
		if u.Level < 1 || u.Level > c.LevelCap {
			return fmt.Errorf("leveling: unlock %s: level out of range", u.Code)
		}
	}
	c.buildRewards()
	return nil
}

func (c *Config) buildTable() error {
	c.table = make([]int64, c.LevelCap-1)
	switch c.Curve.Kind {
	case "table":
		if len(c.Curve.Table) != c.LevelCap-1 {
			return fmt.Errorf("leveling: xp table needs %d entries (levels 1..%d)", c.LevelCap-1, c.LevelCap-1)
		}
		copy(c.table, c.Curve.Table)
	case "formula", "":
		if c.Curve.BaselineXPPerMinute <= 0 {
			return errors.New("leveling: baseline_xp_per_minute must be > 0")
		}
		for l := 1; l < c.LevelCap; l++ {
			minutes := c.Curve.A + c.Curve.B*math.Pow(float64(l), c.Curve.C)
			c.table[l-1] = int64(math.Ceil(minutes * c.Curve.BaselineXPPerMinute))
		}
	default:
		return fmt.Errorf("leveling: unknown xp_curve kind %q", c.Curve.Kind)
	}
	for i, v := range c.table {
		if v < 1 {
			return fmt.Errorf("leveling: xp for level %d must be >= 1", i+1)
		}
	}
	return nil
}

// XPToNext is the XP from level to level+1 (0 at the cap).
func (c *Config) XPToNext(level int) int64 {
	if level < 1 {
		level = 1
	}
	if level >= c.LevelCap {
		return 0
	}
	return c.table[level-1]
}

// Table returns a copy of the cached XP table (index 0 = level 1).
func (c *Config) Table() []int64 { return append([]int64(nil), c.table...) }

// Class looks up a class node.
func (c *Config) Class(code string) (*Class, bool) {
	k, ok := c.classes[code]
	return k, ok
}

// Path is the chain from the base class to spec (spec's own code last).
// Unknown specs are their own one-element path.
func (c *Config) Path(spec string) []string {
	var out []string
	for k, ok := c.classes[spec]; ok; k, ok = c.classes[k.Parent] {
		out = append([]string{k.Code}, out...)
		if k.Parent == "" {
			break
		}
	}
	if len(out) == 0 {
		return []string{spec}
	}
	return out
}

// Choices are the promotions open to spec at level (empty below the gate).
func (c *Config) Choices(spec string, level int) []*Class {
	kids := c.children[spec]
	if len(kids) == 0 || level < kids[0].RequiredLevel {
		return nil
	}
	return kids
}

// gate is the promotion level ahead of spec (0 = none).
func (c *Config) gate(spec string) int {
	if kids := c.children[spec]; len(kids) > 0 {
		return kids[0].RequiredLevel
	}
	return 0
}

// State is a character's level state as the engine sees it.
type State struct {
	Spec     string
	Level    int
	XP       int64
	Overflow int64 // held at a promotion gate
	PostCap  int64 // earned at the level cap
}

// Advance adds amount XP: levels up as far as it reaches, stops at a pending
// promotion gate (the rest is kept as Overflow, never discarded) and at the
// cap (the rest goes to PostCap). Returns the new state and levels gained.
func (c *Config) Advance(st State, amount int64) (State, []int) {
	if amount < 0 {
		amount = 0
	}
	if st.Level >= c.LevelCap {
		st.Level, st.XP = c.LevelCap, 0
		st.PostCap += amount
		return st, nil
	}
	if g := c.gate(st.Spec); g > 0 && st.Level >= g {
		st.Overflow += amount
		return st, nil
	}
	st.XP += amount
	var gained []int
	for {
		if st.Level >= c.LevelCap {
			st.PostCap += st.XP
			st.XP = 0
			break
		}
		need := c.XPToNext(st.Level)
		if st.XP < need {
			break
		}
		st.XP -= need
		st.Level++
		gained = append(gained, st.Level)
		if g := c.gate(st.Spec); g > 0 && st.Level >= g {
			st.Overflow += st.XP
			st.XP = 0
			break
		}
	}
	return st, gained
}

// AttributePointsEarned is the total attribute points granted through level.
func (c *Config) AttributePointsEarned(level int) int {
	r := c.AttributePoints
	n := (level / r.Every) * r.Amount
	for l, amount := range r.Milestones {
		if l <= level {
			n += amount
		}
	}
	return n
}

// TalentPointsEarned is the total talent points granted through level.
func (c *Config) TalentPointsEarned(level int) int {
	r := c.TalentPoints
	if level < r.From {
		return 0
	}
	last := min(level, r.Until)
	return ((last-r.From)/r.Every + 1) * r.Amount
}

// EffectiveAttribute applies diminishing returns to a raw attribute value.
func (c *Config) EffectiveAttribute(raw int) float64 {
	if raw <= 0 {
		return 0
	}
	eff, prev := 0.0, 0
	for _, b := range c.Scaling {
		top := b.UpTo
		if top == 0 || raw <= top {
			return eff + float64(raw-prev)*b.Effectiveness
		}
		eff += float64(top-prev) * b.Effectiveness
		prev = top
	}
	return eff + float64(raw-prev) // past the last bracket: full value
}

// Cap clamps a percent stat by its configured cap (no cap = unchanged).
func (c *Config) Cap(stat string, v float64) float64 {
	if limit, ok := c.Caps[stat]; ok && v > limit {
		return limit
	}
	return v
}

// Talent looks up a talent and its tree.
func (c *Config) Talent(code string) (*Talent, *Tree, bool) {
	t, ok := c.talents[code]
	if !ok {
		return nil, nil, false
	}
	return t, c.treeOf[code], true
}

// TreesFor lists the talent trees of a class path (the base class's).
func (c *Config) TreesFor(spec string) []Tree {
	path := c.Path(spec)
	var out []Tree
	for _, t := range c.Trees {
		for _, p := range path {
			if t.Class == p {
				out = append(out, t)
			}
		}
	}
	return out
}

// buildRewards fills the registry from the rules plus explicit entries.
func (c *Config) buildRewards() {
	c.rewards = map[int][]Reward{}
	add := func(level int, r Reward) {
		if level >= 1 && level <= c.LevelCap {
			c.rewards[level] = append(c.rewards[level], r)
		}
	}
	for l := 2; l <= c.LevelCap; l++ {
		if n := c.AttributePointsEarned(l) - c.AttributePointsEarned(l-1); n > 0 {
			add(l, Reward{Type: RewardAttributePoints, Amount: n})
		}
		if n := c.TalentPointsEarned(l) - c.TalentPointsEarned(l-1); n > 0 {
			add(l, Reward{Type: RewardTalentPoint, Amount: n})
		}
	}
	for l, title := range c.Titles {
		add(l, Reward{Type: RewardTitle, Ref: title})
	}
	for parent, kids := range c.children {
		add(kids[0].RequiredLevel, Reward{Type: RewardPromotion, Ref: parent})
	}
	m := c.Milestones
	add(m.Awakening, Reward{Type: RewardAwakening})
	add(m.Ultimate, Reward{Type: RewardUltimate})
	add(m.Capstone, Reward{Type: RewardCapstone})
	add(m.Mastery, Reward{Type: RewardMastery})
	for _, r := range c.LevelRewards {
		add(r.Level, Reward{Type: r.Type, Amount: r.Amount, Ref: r.Ref})
	}
	for _, list := range c.rewards {
		sort.SliceStable(list, func(i, j int) bool { return list[i].Type < list[j].Type })
	}
}

// RewardsAt is the registry entry for one level.
func (c *Config) RewardsAt(level int) []Reward { return c.rewards[level] }

// UnlocksAt lists the unlocks of a class path at exactly level.
func (c *Config) UnlocksAt(path []string, level int) []Unlock {
	var out []Unlock
	for _, u := range c.Unlocks {
		if u.Level == level && contains(path, u.Spec) {
			out = append(out, u)
		}
	}
	return out
}

// UnlockedBy lists every unlock a path has reached by level.
func (c *Config) UnlockedBy(path []string, level int) []Unlock {
	var out []Unlock
	for _, u := range c.Unlocks {
		if u.Level <= level && contains(path, u.Spec) {
			out = append(out, u)
		}
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// Milestone is the next major objective (DD-250 §44).
type Milestone struct {
	Level   int      `json:"level"`
	Kind    string   `json:"kind"` // promotion | awakening | ultimate | capstone | mastery | skill | passive | title
	Ref     string   `json:"ref,omitempty"`
	Choices []string `json:"choices,omitempty"`
}

var milestoneRank = map[string]int{
	"promotion": 0, "awakening": 1, "ultimate": 2, "capstone": 3, "mastery": 4, "skill": 5, "passive": 6, "title": 7,
}

// NextMilestone is the nearest major unlock above level for spec's path
// (promotion gates, awakening, ultimate, capstone, mastery, skills, titles).
// A pending promotion at the current level is itself the objective.
func (c *Config) NextMilestone(spec string, level int) (Milestone, bool) {
	if kids := c.Choices(spec, level); len(kids) > 0 {
		return Milestone{Level: kids[0].RequiredLevel, Kind: "promotion", Choices: codes(kids)}, true
	}
	path := c.Path(spec)
	var best *Milestone
	consider := func(m Milestone) {
		if m.Level <= level || m.Level > c.LevelCap {
			return
		}
		if best == nil || m.Level < best.Level ||
			(m.Level == best.Level && milestoneRank[m.Kind] < milestoneRank[best.Kind]) {
			best = &m
		}
	}
	if kids := c.children[spec]; len(kids) > 0 {
		consider(Milestone{Level: kids[0].RequiredLevel, Kind: "promotion", Choices: codes(kids)})
	}
	for _, u := range c.Unlocks {
		if contains(path, u.Spec) {
			kind := u.Category
			if kind == "" {
				kind = "skill"
			}
			consider(Milestone{Level: u.Level, Kind: kind, Ref: u.Code})
		}
	}
	m := c.Milestones
	if c.isFinal(spec) || len(c.children[spec]) == 0 {
		consider(Milestone{Level: m.Capstone, Kind: "capstone"})
		consider(Milestone{Level: m.Mastery, Kind: "mastery"})
	}
	for l, t := range c.Titles {
		consider(Milestone{Level: l, Kind: "title", Ref: t})
	}
	if best == nil {
		return Milestone{}, false
	}
	return *best, true
}

// isFinal: a promoted class with no further promotions.
func (c *Config) isFinal(spec string) bool {
	k, ok := c.classes[spec]
	return ok && k.Parent != "" && len(c.children[spec]) == 0
}

func codes(list []*Class) []string {
	out := make([]string, len(list))
	for i, k := range list {
		out[i] = k.Code
	}
	return out
}

// EstimateMinutes is the baseline time from (level, xp) to target level.
func (c *Config) EstimateMinutes(level int, xp int64, target int) float64 {
	rate := c.Curve.BaselineXPPerMinute
	if rate <= 0 || target <= level {
		return 0
	}
	total := -float64(xp)
	for l := level; l < target && l < c.LevelCap; l++ {
		total += float64(c.XPToNext(l))
	}
	return math.Max(0, total/rate)
}

// PromotionPending reports a reached, unchosen promotion gate.
func (c *Config) PromotionPending(spec string, level int) bool {
	return len(c.Choices(spec, level)) > 0
}
