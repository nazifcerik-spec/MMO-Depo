// Package contentpack loads YAML content (items, vendors, ...) from a directory
// and upserts it into Postgres content.* tables. Idempotent — same content = no
// row mutation. Runs at startup after migrations; dev endpoint reloads it.
//
// Design (Sprint 12, DD-045 spirit):
//   - YAML is the source of truth. Migration seed exists only as a bootstrap
//     fallback so a fresh DB is playable without the loader having run yet.
//   - Loader is upsert-only: never DELETE. Template retirement goes through a
//     separate migration path (would need lifecycle handling in Phase 02).
//   - Upsert wrapped in a single tx per file so partial YAML rot cannot leave
//     the catalog half-updated.
//   - After DB upsert, in-process caches (item.Catalog, vendor cache) get told
//     to refresh via a Reloader callback.
package contentpack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"

	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/skill"
)

// --- YAML shapes ------------------------------------------------------------

// Item mirrors content.item_template + Sprint 11 columns.
type Item struct {
	Code             string      `yaml:"code"`
	DisplayNameKey   string      `yaml:"display_name_key"`
	ItemClass        string      `yaml:"item_class"`
	ItemSubclass     string      `yaml:"item_subclass"`
	MinLevel         int16       `yaml:"min_level"`
	BaseValueCopper  int64       `yaml:"base_value_copper"`
	BaseWeightG      int32       `yaml:"base_weight_g"`
	Stackable        bool        `yaml:"stackable"`
	MaxStack         int32       `yaml:"max_stack"`
	BindKind         string      `yaml:"bind_kind"`
	NoTrade          bool        `yaml:"no_trade"`
	NoSell           bool        `yaml:"no_sell"`
	Affixable        bool        `yaml:"affixable"`
	EquipSlot        string      `yaml:"equip_slot"`
	ClassRestriction string      `yaml:"class_restriction"`
	StatBonuses      []StatBonus `yaml:"stat_bonuses"`
	UseEffect        *UseEffect  `yaml:"use_effect"`
}

// UseEffect is what consuming one unit of a consumable does (DD-218).
type UseEffect struct {
	Kind   string `yaml:"kind"   json:"kind"`
	Amount int    `yaml:"amount" json:"amount"`
}

// validUseEffects: only effects with persisted state (HP, DD-218) exist yet.
var validUseEffects = map[string]bool{"heal_hp": true}

// StatBonus is one entry of an item's equipped modifier stack.
type StatBonus struct {
	Kind   string  `yaml:"kind"   json:"kind"`
	Target string  `yaml:"target" json:"target"`
	Amount float64 `yaml:"amount" json:"amount"`
}

// itemsFile is the top-level shape of content/items.yaml.
type itemsFile struct {
	Items []Item `yaml:"items"`
}

// Vendor mirrors content.vendor_template.
type Vendor struct {
	Code           string       `yaml:"code"`
	DisplayNameKey string       `yaml:"display_name_key"`
	VendorType     string       `yaml:"vendor_type"`
	BuyMultiplier  float64      `yaml:"buy_multiplier"`
	SellMultiplier float64      `yaml:"sell_multiplier"`
	Sells          []VendorSell `yaml:"sells"`
}

type VendorSell struct {
	TemplateCode        string `yaml:"template_code"         json:"template_code"`
	PriceOverrideCopper *int64 `yaml:"price_override_copper,omitempty" json:"price_override_copper,omitempty"`
}

type vendorsFile struct {
	Vendors []Vendor `yaml:"vendors"`
}

// Mob mirrors content.mob_template (Sprint 13).
type Mob struct {
	Code           string `yaml:"code"`
	DisplayNameKey string `yaml:"display_name_key"`
	Level          int16  `yaml:"level"`
	HPMax          int32  `yaml:"hp_max"`
	Attack         int32  `yaml:"attack"`
	Defense        int32  `yaml:"defense"`
	CritPct        int16  `yaml:"crit_pct"`
	HitPct         int16  `yaml:"hit_pct"`
	EvasionPct     int16  `yaml:"evasion_pct"`
	XPReward       int32  `yaml:"xp_reward"`
	GoldCopperMin  int64  `yaml:"gold_copper_min"`
	GoldCopperMax  int64  `yaml:"gold_copper_max"`
}

type mobsFile struct {
	Mobs []Mob `yaml:"mobs"`
}

// Dungeon mirrors content.dungeon_template (Sprint 15).
type Dungeon struct {
	Code            string        `yaml:"code"`
	DisplayNameKey  string        `yaml:"display_name_key"`
	MinLevel        int16         `yaml:"min_level"`
	DurationSeconds int32         `yaml:"duration_seconds"`
	Waves           []DungeonWave `yaml:"waves"`
	BossMobCode     string        `yaml:"boss_mob_code"`
	BonusGoldCopper int64         `yaml:"bonus_gold_copper"`
	BonusXP         int32         `yaml:"bonus_xp"`
}

type DungeonWave struct {
	MobCode string `yaml:"mob_code" json:"mob_code"`
	Count   int    `yaml:"count"    json:"count"`
}

type dungeonsFile struct {
	Dungeons []Dungeon `yaml:"dungeons"`
}

// Region mirrors world.region (Sprint 16).
type Region struct {
	Code                    string `yaml:"code"`
	DisplayNameKey          string `yaml:"display_name_key"`
	Kind                    string `yaml:"kind"`
	ParentCode              string `yaml:"parent_code"`
	MinLevel                int16  `yaml:"min_level"`
	TravelSecondsFromParent int32  `yaml:"travel_seconds_from_parent"`
	Safe                    bool   `yaml:"safe"`
	IsStarting              bool   `yaml:"is_starting"`
}

// RegionLink mirrors world.region_link.
type RegionLink struct {
	RegionACode   string `yaml:"region_a_code"`
	RegionBCode   string `yaml:"region_b_code"`
	TravelSeconds int32  `yaml:"travel_seconds"`
}

type regionsFile struct {
	Regions     []Region     `yaml:"regions"`
	RegionLinks []RegionLink `yaml:"region_links"`
}

// Profession mirrors content.profession (Sprint 17).
type Profession struct {
	Code           string `yaml:"code"`
	DisplayNameKey string `yaml:"display_name_key"`
	MaxLevel       int16  `yaml:"max_level"`
}

type professionsFile struct {
	Professions []Profession `yaml:"professions"`
}

// Recipe mirrors content.recipe.
type Recipe struct {
	Code               string        `yaml:"code"`
	DisplayNameKey     string        `yaml:"display_name_key"`
	ProfessionCode     string        `yaml:"profession_code"`
	MinSkill           int16         `yaml:"min_skill"`
	Inputs             []RecipeInput `yaml:"inputs"`
	GoldCostCopper     int64         `yaml:"gold_cost_copper"`
	OutputTemplateCode string        `yaml:"output_template_code"`
	OutputStack        int32         `yaml:"output_stack"`
	ProfessionXPReward int32         `yaml:"profession_xp_reward"`
}

type RecipeInput struct {
	TemplateCode string `yaml:"template_code" json:"template_code"`
	Quantity     int32  `yaml:"quantity"      json:"quantity"`
}

type recipesFile struct {
	Recipes []Recipe `yaml:"recipes"`
}

// --- Load result ------------------------------------------------------------

// LoadResult reports what the loader upserted this pass — used by callers to
// log stats and by tests to assert idempotency.
type LoadResult struct {
	ItemsSeen       int
	VendorsSeen     int
	MobsSeen        int
	DungeonsSeen    int
	RegionsSeen     int
	RegionLinksSeen int
	ProfessionsSeen int
	RecipesSeen     int
	QuestsSeen      int
	SkillsSeen      int
	SourceDir       string
}

// Reloader is called AFTER a successful DB upsert so in-process caches
// (item.Catalog, vendor list) can rebuild themselves from the fresh DB rows.
// Every registered callback runs in order; the first error aborts the reload.
type Reloader interface {
	Reload(ctx context.Context) error
}

// --- Loader -----------------------------------------------------------------

// Loader reads YAML files from Dir and upserts them into Postgres. Same
// instance can be called multiple times (startup + hot reload).
type Loader struct {
	Dir       string
	DB        *pgxpool.Pool
	Log       *slog.Logger
	reloaders []Reloader
}

func New(dir string, pool *pgxpool.Pool, logger *slog.Logger) *Loader {
	if logger == nil {
		logger = slog.Default()
	}
	return &Loader{Dir: dir, DB: pool, Log: logger}
}

// Register attaches a cache reloader. Called after every successful Load.
func (l *Loader) Register(r Reloader) { l.reloaders = append(l.reloaders, r) }

// Load reads every YAML file, upserts the DB, then fires reloaders. Returns
// stats or the first error encountered. Missing directory returns an error;
// missing individual files are tolerated (some servers may not need vendors,
// etc).
func (l *Loader) Load(ctx context.Context) (LoadResult, error) {
	res := LoadResult{SourceDir: l.Dir}
	p, err := Read(l.Dir)
	if err != nil {
		return res, err
	}
	items, vendors, mobs, dungeons := p.Items, p.Vendors, p.Mobs, p.Dungeons
	regions, regionLinks, professions, recipes := p.Regions, p.RegionLinks, p.Professions, p.Recipes
	quests, skills := p.Quests, p.Skills
	res.ItemsSeen, res.VendorsSeen, res.MobsSeen, res.DungeonsSeen = len(items), len(vendors), len(mobs), len(dungeons)
	res.RegionsSeen, res.RegionLinksSeen = len(regions), len(regionLinks)
	res.ProfessionsSeen, res.RecipesSeen = len(professions), len(recipes)
	res.QuestsSeen, res.SkillsSeen = len(quests), len(skills)

	if len(items) > 0 {
		if err := l.upsertItems(ctx, items); err != nil {
			return res, fmt.Errorf("contentpack: upsert items: %w", err)
		}
	}
	if len(vendors) > 0 {
		if err := l.upsertVendors(ctx, vendors); err != nil {
			return res, fmt.Errorf("contentpack: upsert vendors: %w", err)
		}
	}
	if len(mobs) > 0 {
		if err := l.upsertMobs(ctx, mobs); err != nil {
			return res, fmt.Errorf("contentpack: upsert mobs: %w", err)
		}
	}
	if len(dungeons) > 0 {
		if err := l.upsertDungeons(ctx, dungeons); err != nil {
			return res, fmt.Errorf("contentpack: upsert dungeons: %w", err)
		}
	}
	// Region UPSERT must happen BEFORE region_link UPSERT (FK dependency).
	// Both are guarded inside a single upsertRegions call so a mid-batch failure
	// rolls back cleanly.
	if len(regions) > 0 || len(regionLinks) > 0 {
		if err := l.upsertRegions(ctx, regions, regionLinks); err != nil {
			return res, fmt.Errorf("contentpack: upsert regions: %w", err)
		}
	}
	if len(professions) > 0 {
		if err := l.upsertProfessions(ctx, professions); err != nil {
			return res, fmt.Errorf("contentpack: upsert professions: %w", err)
		}
	}
	if len(recipes) > 0 {
		if err := l.upsertRecipes(ctx, recipes); err != nil {
			return res, fmt.Errorf("contentpack: upsert recipes: %w", err)
		}
	}

	if len(skills) > 0 {
		if err := l.upsertSkills(ctx, skills); err != nil {
			return res, fmt.Errorf("contentpack: upsert skills: %w", err)
		}
	}
	if len(quests) > 0 {
		if err := l.upsertQuests(ctx, quests); err != nil {
			return res, fmt.Errorf("contentpack: upsert quests: %w", err)
		}
	}

	for _, r := range l.reloaders {
		if err := r.Reload(ctx); err != nil {
			return res, fmt.Errorf("contentpack: reloader: %w", err)
		}
	}

	l.Log.InfoContext(ctx, "content pack loaded",
		slog.String("dir", l.Dir),
		slog.Int("items", res.ItemsSeen),
		slog.Int("vendors", res.VendorsSeen),
		slog.Int("mobs", res.MobsSeen),
		slog.Int("dungeons", res.DungeonsSeen),
		slog.Int("regions", res.RegionsSeen),
		slog.Int("region_links", res.RegionLinksSeen),
		slog.Int("professions", res.ProfessionsSeen),
		slog.Int("recipes", res.RecipesSeen),
		slog.Int("quests", res.QuestsSeen),
		slog.Int("skills", res.SkillsSeen),
	)
	return res, nil
}

// Pack is a parsed and validated content directory.
type Pack struct {
	Items       []Item
	Vendors     []Vendor
	Mobs        []Mob
	Dungeons    []Dungeon
	Regions     []Region
	RegionLinks []RegionLink
	Professions []Profession
	Recipes     []Recipe
	Quests      []Quest
	Skills      []ClassSkill
}

// Read parses and validates every file in dir, including the references
// between files, without touching the database. Load uses it; tests and the
// balance simulation use it directly.
func Read(dir string) (*Pack, error) {
	if dir == "" {
		return nil, errors.New("contentpack: empty directory")
	}
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("contentpack: stat %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("contentpack: %s is not a directory", dir)
	}
	p := &Pack{}
	f := func(name string) string { return filepath.Join(dir, name) }
	if p.Items, err = loadItemsFile(f("items.yaml")); err != nil {
		return nil, err
	}
	// Equipment lives in its own file (gear ladder, DD-229); codes are unique
	// across both.
	gear, err := loadItemsFile(f("gear.yaml"))
	if err != nil {
		return nil, err
	}
	p.Items = append(p.Items, gear...)
	if err := validateItems(p.Items); err != nil {
		return nil, fmt.Errorf("validate items + gear: %w", err)
	}
	if p.Vendors, err = loadVendorsFile(f("vendors.yaml")); err != nil {
		return nil, err
	}
	if p.Mobs, err = loadMobsFile(f("mobs.yaml")); err != nil {
		return nil, err
	}
	if p.Dungeons, err = loadDungeonsFile(f("dungeons.yaml")); err != nil {
		return nil, err
	}
	if p.Regions, p.RegionLinks, err = loadRegionsFile(f("regions.yaml")); err != nil {
		return nil, err
	}
	if p.Professions, err = loadProfessionsFile(f("professions.yaml")); err != nil {
		return nil, err
	}
	if p.Recipes, err = loadRecipesFile(f("recipes.yaml")); err != nil {
		return nil, err
	}
	if p.Quests, err = loadQuestsFile(f("quests.yaml")); err != nil {
		return nil, err
	}
	if p.Skills, err = loadSkillsFile(f("skills.yaml")); err != nil {
		return nil, err
	}
	r := refs{
		mobs:     codeSet(p.Mobs, func(m Mob) string { return m.Code }),
		dungeons: codeSet(p.Dungeons, func(d Dungeon) string { return d.Code }),
		items:    codeSet(p.Items, func(i Item) string { return i.Code }),
	}
	if err := validateQuests(p.Quests, r); err != nil {
		return nil, fmt.Errorf("validate quests.yaml: %w", err)
	}
	if err := validateRefs(p, r); err != nil {
		return nil, err
	}
	return p, nil
}

// validateRefs checks the references between files: a typo would otherwise
// load fine and fail only when a player buys, crafts or enters.
func validateRefs(p *Pack, r refs) error {
	for _, v := range p.Vendors {
		for _, s := range v.Sells {
			if !r.items[s.TemplateCode] {
				return fmt.Errorf("vendor %s: sells unknown item %q", v.Code, s.TemplateCode)
			}
		}
	}
	profs := codeSet(p.Professions, func(x Profession) string { return x.Code })
	for _, rc := range p.Recipes {
		if !profs[rc.ProfessionCode] {
			return fmt.Errorf("recipe %s: unknown profession %q", rc.Code, rc.ProfessionCode)
		}
		if !r.items[rc.OutputTemplateCode] {
			return fmt.Errorf("recipe %s: unknown output %q", rc.Code, rc.OutputTemplateCode)
		}
		for _, in := range rc.Inputs {
			if !r.items[in.TemplateCode] {
				return fmt.Errorf("recipe %s: unknown input %q", rc.Code, in.TemplateCode)
			}
		}
	}
	for _, d := range p.Dungeons {
		if !r.mobs[d.BossMobCode] {
			return fmt.Errorf("dungeon %s: unknown boss %q", d.Code, d.BossMobCode)
		}
		for _, w := range d.Waves {
			if !r.mobs[w.MobCode] {
				return fmt.Errorf("dungeon %s: unknown wave mob %q", d.Code, w.MobCode)
			}
		}
	}
	return nil
}

// --- File I/O ---------------------------------------------------------------

func loadItemsFile(path string) ([]Item, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f itemsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateItems(f.Items); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Items, nil
}

func loadProfessionsFile(path string) ([]Profession, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f professionsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateProfessions(f.Professions); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Professions, nil
}

func loadRecipesFile(path string) ([]Recipe, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f recipesFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateRecipes(f.Recipes); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Recipes, nil
}

func loadRegionsFile(path string) ([]Region, []RegionLink, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f regionsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateRegions(f.Regions, f.RegionLinks); err != nil {
		return nil, nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Regions, f.RegionLinks, nil
}

func loadDungeonsFile(path string) ([]Dungeon, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f dungeonsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateDungeons(f.Dungeons); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Dungeons, nil
}

func loadMobsFile(path string) ([]Mob, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f mobsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateMobs(f.Mobs); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Mobs, nil
}

func loadVendorsFile(path string) ([]Vendor, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var f vendorsFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validateVendors(f.Vendors); err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}
	return f.Vendors, nil
}

// --- Validation -------------------------------------------------------------

var (
	validItemClasses = map[string]bool{
		"weapon": true, "armor": true, "consumable": true, "material": true, "misc": true,
	}
	validBindKinds = map[string]bool{"boe": true, "bop": true, "boa": true}
	// Must match what the equip service accepts, or equippable content can
	// never be equipped (E_UNKNOWN_SLOT).
	validEquipSlots    = inventory.KnownEquipSlots
	validModifierKinds = func(k string) bool { return skill.KnownKinds[skill.ModifierKind(k)] }
	validStatTargets   = func(t string) bool { return skill.KnownTargets[skill.Target(t)] }
)

func validateItems(items []Item) error {
	seen := make(map[string]bool, len(items))
	for i, it := range items {
		if it.Code == "" {
			return fmt.Errorf("item[%d]: empty code", i)
		}
		if seen[it.Code] {
			return fmt.Errorf("item[%d] %s: duplicate code", i, it.Code)
		}
		seen[it.Code] = true
		if !validItemClasses[it.ItemClass] {
			return fmt.Errorf("item %s: unknown item_class %q", it.Code, it.ItemClass)
		}
		if it.BindKind != "" && !validBindKinds[it.BindKind] {
			return fmt.Errorf("item %s: unknown bind_kind %q", it.Code, it.BindKind)
		}
		if it.EquipSlot != "" && !validEquipSlots[it.EquipSlot] {
			return fmt.Errorf("item %s: unknown equip_slot %q", it.Code, it.EquipSlot)
		}
		if it.Stackable && it.MaxStack < 1 {
			return fmt.Errorf("item %s: stackable requires max_stack >= 1", it.Code)
		}
		if ue := it.UseEffect; ue != nil {
			if it.ItemClass != "consumable" {
				return fmt.Errorf("item %s: use_effect only allowed on consumables", it.Code)
			}
			if !validUseEffects[ue.Kind] {
				return fmt.Errorf("item %s: unknown use_effect kind %q", it.Code, ue.Kind)
			}
			if ue.Amount <= 0 || ue.Amount > 100_000 {
				return fmt.Errorf("item %s: use_effect amount out of (0, 100000]", it.Code)
			}
		}
		for j, sb := range it.StatBonuses {
			if !validModifierKinds(sb.Kind) {
				return fmt.Errorf("item %s: stat_bonuses[%d] unknown kind %q", it.Code, j, sb.Kind)
			}
			if !validStatTargets(sb.Target) {
				return fmt.Errorf("item %s: stat_bonuses[%d] unknown target %q", it.Code, j, sb.Target)
			}
		}
	}
	return nil
}

func validateProfessions(profs []Profession) error {
	seen := make(map[string]bool, len(profs))
	for i, p := range profs {
		if p.Code == "" {
			return fmt.Errorf("profession[%d]: empty code", i)
		}
		if seen[p.Code] {
			return fmt.Errorf("profession[%d] %s: duplicate code", i, p.Code)
		}
		seen[p.Code] = true
		if p.MaxLevel <= 0 || p.MaxLevel > 100 {
			return fmt.Errorf("profession %s: max_level out of (0, 100]", p.Code)
		}
	}
	return nil
}

func validateRecipes(recipes []Recipe) error {
	seen := make(map[string]bool, len(recipes))
	for i, r := range recipes {
		if r.Code == "" {
			return fmt.Errorf("recipe[%d]: empty code", i)
		}
		if seen[r.Code] {
			return fmt.Errorf("recipe[%d] %s: duplicate code", i, r.Code)
		}
		seen[r.Code] = true
		if r.ProfessionCode == "" {
			return fmt.Errorf("recipe %s: profession_code required", r.Code)
		}
		if r.MinSkill < 1 {
			return fmt.Errorf("recipe %s: min_skill must be >= 1", r.Code)
		}
		if r.OutputTemplateCode == "" {
			return fmt.Errorf("recipe %s: output_template_code required", r.Code)
		}
		if r.OutputStack < 1 {
			return fmt.Errorf("recipe %s: output_stack must be >= 1", r.Code)
		}
		if r.GoldCostCopper < 0 {
			return fmt.Errorf("recipe %s: gold_cost_copper must be >= 0", r.Code)
		}
		if r.ProfessionXPReward < 0 {
			return fmt.Errorf("recipe %s: profession_xp_reward must be >= 0", r.Code)
		}
		if len(r.Inputs) == 0 {
			return fmt.Errorf("recipe %s: at least one input required", r.Code)
		}
		inputSeen := make(map[string]bool, len(r.Inputs))
		for j, in := range r.Inputs {
			if in.TemplateCode == "" {
				return fmt.Errorf("recipe %s: inputs[%d].template_code empty", r.Code, j)
			}
			if in.Quantity <= 0 {
				return fmt.Errorf("recipe %s: inputs[%d].quantity must be > 0", r.Code, j)
			}
			if inputSeen[in.TemplateCode] {
				return fmt.Errorf("recipe %s: inputs[%d] duplicate template %q — merge into one row", r.Code, j, in.TemplateCode)
			}
			inputSeen[in.TemplateCode] = true
		}
	}
	return nil
}

var validRegionKinds = map[string]bool{
	"continent": true, "zone": true, "city": true, "wilderness": true, "dungeon": true,
}

func validateRegions(regions []Region, links []RegionLink) error {
	seen := make(map[string]bool, len(regions))
	startingCount := 0
	for i, r := range regions {
		if r.Code == "" {
			return fmt.Errorf("region[%d]: empty code", i)
		}
		if seen[r.Code] {
			return fmt.Errorf("region[%d] %s: duplicate code", i, r.Code)
		}
		seen[r.Code] = true
		if !validRegionKinds[r.Kind] {
			return fmt.Errorf("region %s: unknown kind %q", r.Code, r.Kind)
		}
		if r.TravelSecondsFromParent < 0 || r.TravelSecondsFromParent > 3600 {
			return fmt.Errorf("region %s: travel_seconds_from_parent out of [0, 3600]", r.Code)
		}
		if r.IsStarting {
			startingCount++
		}
	}
	if startingCount > 1 {
		return fmt.Errorf("regions: exactly one is_starting allowed, got %d", startingCount)
	}
	// Parent existence + non-self-reference.
	for _, r := range regions {
		if r.ParentCode != "" {
			if r.ParentCode == r.Code {
				return fmt.Errorf("region %s: parent_code points to self", r.Code)
			}
			if !seen[r.ParentCode] {
				return fmt.Errorf("region %s: parent_code %q not in regions list", r.Code, r.ParentCode)
			}
		}
	}
	// Links: both endpoints exist, canonical order (a < b), no self-link.
	linkSeen := make(map[[2]string]bool, len(links))
	for i, l := range links {
		if l.RegionACode == "" || l.RegionBCode == "" {
			return fmt.Errorf("region_link[%d]: empty code", i)
		}
		if l.RegionACode == l.RegionBCode {
			return fmt.Errorf("region_link[%d]: self-link on %q", i, l.RegionACode)
		}
		if l.RegionACode >= l.RegionBCode {
			return fmt.Errorf("region_link[%d]: (%q, %q) not in canonical order (a < b)",
				i, l.RegionACode, l.RegionBCode)
		}
		if !seen[l.RegionACode] {
			return fmt.Errorf("region_link[%d]: unknown region %q", i, l.RegionACode)
		}
		if !seen[l.RegionBCode] {
			return fmt.Errorf("region_link[%d]: unknown region %q", i, l.RegionBCode)
		}
		key := [2]string{l.RegionACode, l.RegionBCode}
		if linkSeen[key] {
			return fmt.Errorf("region_link[%d]: duplicate (%q, %q)", i, l.RegionACode, l.RegionBCode)
		}
		linkSeen[key] = true
		if l.TravelSeconds < 0 || l.TravelSeconds > 3600 {
			return fmt.Errorf("region_link[%d]: travel_seconds out of [0, 3600]", i)
		}
	}
	return nil
}

func validateDungeons(dungeons []Dungeon) error {
	seen := make(map[string]bool, len(dungeons))
	for i, d := range dungeons {
		if d.Code == "" {
			return fmt.Errorf("dungeon[%d]: empty code", i)
		}
		if seen[d.Code] {
			return fmt.Errorf("dungeon[%d] %s: duplicate code", i, d.Code)
		}
		seen[d.Code] = true
		if d.DurationSeconds <= 0 || d.DurationSeconds > 3600 {
			return fmt.Errorf("dungeon %s: duration_seconds out of (0,3600]", d.Code)
		}
		if len(d.Waves) == 0 {
			return fmt.Errorf("dungeon %s: at least one wave required", d.Code)
		}
		if d.BossMobCode == "" {
			return fmt.Errorf("dungeon %s: boss_mob_code required", d.Code)
		}
		for j, w := range d.Waves {
			if w.MobCode == "" {
				return fmt.Errorf("dungeon %s: waves[%d] mob_code empty", d.Code, j)
			}
			if w.Count <= 0 || w.Count > 32 {
				return fmt.Errorf("dungeon %s: waves[%d] count out of (0,32]", d.Code, j)
			}
		}
		if d.BonusGoldCopper < 0 || d.BonusXP < 0 {
			return fmt.Errorf("dungeon %s: bonus rewards must be >= 0", d.Code)
		}
	}
	return nil
}

func validateMobs(mobs []Mob) error {
	seen := make(map[string]bool, len(mobs))
	for i, m := range mobs {
		if m.Code == "" {
			return fmt.Errorf("mob[%d]: empty code", i)
		}
		if seen[m.Code] {
			return fmt.Errorf("mob[%d] %s: duplicate code", i, m.Code)
		}
		seen[m.Code] = true
		if m.HPMax <= 0 {
			return fmt.Errorf("mob %s: hp_max must be > 0", m.Code)
		}
		if m.Attack < 0 || m.Defense < 0 {
			return fmt.Errorf("mob %s: negative attack/defense", m.Code)
		}
		if m.CritPct < 0 || m.CritPct > 100 {
			return fmt.Errorf("mob %s: crit_pct out of range 0..100", m.Code)
		}
		if m.HitPct < 0 || m.HitPct > 100 {
			return fmt.Errorf("mob %s: hit_pct out of range 0..100", m.Code)
		}
		if m.EvasionPct < 0 || m.EvasionPct > 100 {
			return fmt.Errorf("mob %s: evasion_pct out of range 0..100", m.Code)
		}
		if m.GoldCopperMax < m.GoldCopperMin {
			return fmt.Errorf("mob %s: gold_copper_max < gold_copper_min", m.Code)
		}
	}
	return nil
}

func validateVendors(vendors []Vendor) error {
	seen := make(map[string]bool, len(vendors))
	for i, v := range vendors {
		if v.Code == "" {
			return fmt.Errorf("vendor[%d]: empty code", i)
		}
		if seen[v.Code] {
			return fmt.Errorf("vendor[%d] %s: duplicate code", i, v.Code)
		}
		seen[v.Code] = true
		if v.BuyMultiplier < 0 || v.SellMultiplier < 0 {
			return fmt.Errorf("vendor %s: negative multiplier", v.Code)
		}
	}
	return nil
}

// --- DB upsert --------------------------------------------------------------

func (l *Loader) upsertItems(ctx context.Context, items []Item) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO content.item_template (
  code, display_name_key, item_class, item_subclass, min_level,
  base_value_copper, base_weight_g, stackable, max_stack, bind_kind,
  no_trade, no_sell, affixable, equip_slot, class_restriction, stat_bonuses, use_effect
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
ON CONFLICT (code) DO UPDATE SET
  display_name_key   = EXCLUDED.display_name_key,
  item_class         = EXCLUDED.item_class,
  item_subclass      = EXCLUDED.item_subclass,
  min_level          = EXCLUDED.min_level,
  base_value_copper  = EXCLUDED.base_value_copper,
  base_weight_g      = EXCLUDED.base_weight_g,
  stackable          = EXCLUDED.stackable,
  max_stack          = EXCLUDED.max_stack,
  bind_kind          = EXCLUDED.bind_kind,
  no_trade           = EXCLUDED.no_trade,
  no_sell            = EXCLUDED.no_sell,
  affixable          = EXCLUDED.affixable,
  equip_slot         = EXCLUDED.equip_slot,
  class_restriction  = EXCLUDED.class_restriction,
  stat_bonuses       = EXCLUDED.stat_bonuses,
  use_effect         = EXCLUDED.use_effect,
  version            = content.item_template.version + 1`

	for _, it := range items {
		bindKind := it.BindKind
		if bindKind == "" {
			bindKind = "boe"
		}
		maxStack := it.MaxStack
		if maxStack < 1 {
			maxStack = 1
		}
		bonuses, err := json.Marshal(it.StatBonuses)
		if err != nil {
			return fmt.Errorf("item %s: marshal stat_bonuses: %w", it.Code, err)
		}
		var useEffect []byte // nil → SQL NULL (not usable)
		if it.UseEffect != nil {
			if useEffect, err = json.Marshal(it.UseEffect); err != nil {
				return fmt.Errorf("item %s: marshal use_effect: %w", it.Code, err)
			}
		}
		_, err = tx.Exec(ctx, q,
			it.Code, it.DisplayNameKey, it.ItemClass, it.ItemSubclass, it.MinLevel,
			it.BaseValueCopper, it.BaseWeightG, it.Stackable, maxStack, bindKind,
			it.NoTrade, it.NoSell, it.Affixable,
			nullableStr(it.EquipSlot), nullableStr(it.ClassRestriction), bonuses, useEffect,
		)
		if err != nil {
			return fmt.Errorf("upsert item %s: %w", it.Code, err)
		}
	}
	return tx.Commit(ctx)
}

func (l *Loader) upsertVendors(ctx context.Context, vendors []Vendor) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO content.vendor_template (
  code, display_name_key, vendor_type, buy_multiplier, sell_multiplier, sells
) VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (code) DO UPDATE SET
  display_name_key = EXCLUDED.display_name_key,
  vendor_type      = EXCLUDED.vendor_type,
  buy_multiplier   = EXCLUDED.buy_multiplier,
  sell_multiplier  = EXCLUDED.sell_multiplier,
  sells            = EXCLUDED.sells,
  version          = content.vendor_template.version + 1`

	for _, v := range vendors {
		buy := v.BuyMultiplier
		if buy == 0 {
			buy = 1.0
		}
		sell := v.SellMultiplier
		if sell == 0 {
			sell = 0.25
		}
		sells, err := json.Marshal(v.Sells)
		if err != nil {
			return fmt.Errorf("vendor %s: marshal sells: %w", v.Code, err)
		}
		_, err = tx.Exec(ctx, q,
			v.Code, v.DisplayNameKey, v.VendorType, buy, sell, sells,
		)
		if err != nil {
			return fmt.Errorf("upsert vendor %s: %w", v.Code, err)
		}
	}
	return tx.Commit(ctx)
}

func (l *Loader) upsertMobs(ctx context.Context, mobs []Mob) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO content.mob_template (
  code, display_name_key, level,
  hp_max, attack, defense, crit_pct, hit_pct, evasion_pct,
  xp_reward, gold_copper_min, gold_copper_max
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
ON CONFLICT (code) DO UPDATE SET
  display_name_key = EXCLUDED.display_name_key,
  level            = EXCLUDED.level,
  hp_max           = EXCLUDED.hp_max,
  attack           = EXCLUDED.attack,
  defense          = EXCLUDED.defense,
  crit_pct         = EXCLUDED.crit_pct,
  hit_pct          = EXCLUDED.hit_pct,
  evasion_pct      = EXCLUDED.evasion_pct,
  xp_reward        = EXCLUDED.xp_reward,
  gold_copper_min  = EXCLUDED.gold_copper_min,
  gold_copper_max  = EXCLUDED.gold_copper_max,
  version          = content.mob_template.version + 1`

	for _, m := range mobs {
		_, err = tx.Exec(ctx, q,
			m.Code, m.DisplayNameKey, m.Level,
			m.HPMax, m.Attack, m.Defense, m.CritPct, m.HitPct, m.EvasionPct,
			m.XPReward, m.GoldCopperMin, m.GoldCopperMax,
		)
		if err != nil {
			return fmt.Errorf("upsert mob %s: %w", m.Code, err)
		}
	}
	return tx.Commit(ctx)
}

func (l *Loader) upsertDungeons(ctx context.Context, dungeons []Dungeon) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
INSERT INTO content.dungeon_template (
  code, display_name_key, min_level, duration_seconds,
  waves, boss_mob_code, bonus_gold_copper, bonus_xp
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (code) DO UPDATE SET
  display_name_key   = EXCLUDED.display_name_key,
  min_level          = EXCLUDED.min_level,
  duration_seconds   = EXCLUDED.duration_seconds,
  waves              = EXCLUDED.waves,
  boss_mob_code      = EXCLUDED.boss_mob_code,
  bonus_gold_copper  = EXCLUDED.bonus_gold_copper,
  bonus_xp           = EXCLUDED.bonus_xp,
  version            = content.dungeon_template.version + 1`

	for _, d := range dungeons {
		waves, err := json.Marshal(d.Waves)
		if err != nil {
			return fmt.Errorf("dungeon %s: marshal waves: %w", d.Code, err)
		}
		_, err = tx.Exec(ctx, q,
			d.Code, d.DisplayNameKey, d.MinLevel, d.DurationSeconds,
			waves, d.BossMobCode, d.BonusGoldCopper, d.BonusXP,
		)
		if err != nil {
			return fmt.Errorf("upsert dungeon %s: %w", d.Code, err)
		}
	}
	return tx.Commit(ctx)
}

func (l *Loader) upsertProfessions(ctx context.Context, professions []Profession) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const q = `
INSERT INTO content.profession (code, display_name_key, max_level)
VALUES ($1, $2, $3)
ON CONFLICT (code) DO UPDATE SET
  display_name_key = EXCLUDED.display_name_key,
  max_level        = EXCLUDED.max_level,
  version          = content.profession.version + 1`
	for _, p := range professions {
		if _, err := tx.Exec(ctx, q, p.Code, p.DisplayNameKey, p.MaxLevel); err != nil {
			return fmt.Errorf("upsert profession %s: %w", p.Code, err)
		}
	}
	return tx.Commit(ctx)
}

func (l *Loader) upsertRecipes(ctx context.Context, recipes []Recipe) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const q = `
INSERT INTO content.recipe (
  code, display_name_key, profession_code, min_skill,
  inputs, gold_cost_copper,
  output_template_code, output_stack, profession_xp_reward
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
ON CONFLICT (code) DO UPDATE SET
  display_name_key     = EXCLUDED.display_name_key,
  profession_code      = EXCLUDED.profession_code,
  min_skill            = EXCLUDED.min_skill,
  inputs               = EXCLUDED.inputs,
  gold_cost_copper     = EXCLUDED.gold_cost_copper,
  output_template_code = EXCLUDED.output_template_code,
  output_stack         = EXCLUDED.output_stack,
  profession_xp_reward = EXCLUDED.profession_xp_reward,
  version              = content.recipe.version + 1`
	for _, r := range recipes {
		inputs, err := json.Marshal(r.Inputs)
		if err != nil {
			return fmt.Errorf("recipe %s: marshal inputs: %w", r.Code, err)
		}
		_, err = tx.Exec(ctx, q,
			r.Code, r.DisplayNameKey, r.ProfessionCode, r.MinSkill,
			inputs, r.GoldCostCopper,
			r.OutputTemplateCode, r.OutputStack, r.ProfessionXPReward,
		)
		if err != nil {
			return fmt.Errorf("upsert recipe %s: %w", r.Code, err)
		}
	}
	return tx.Commit(ctx)
}

func (l *Loader) upsertRegions(ctx context.Context, regions []Region, links []RegionLink) error {
	tx, err := l.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Regions must be inserted parent-first so the FK is satisfied. Simple
	// two-pass: parents (no parent_code) first, then everything else. Deep
	// hierarchies would need a proper topological sort — content is shallow
	// enough (~4 levels) that this is fine.
	orderedRegions := make([]Region, 0, len(regions))
	added := make(map[string]bool, len(regions))
	remaining := append([]Region(nil), regions...)
	// Multi-pass until stable.
	for {
		progress := false
		next := remaining[:0]
		for _, r := range remaining {
			if r.ParentCode == "" || added[r.ParentCode] {
				orderedRegions = append(orderedRegions, r)
				added[r.Code] = true
				progress = true
			} else {
				next = append(next, r)
			}
		}
		remaining = next
		if len(remaining) == 0 {
			break
		}
		if !progress {
			return fmt.Errorf("regions: unresolvable parent chain, remaining=%d", len(remaining))
		}
	}

	const q = `
INSERT INTO world.region (
  code, display_name_key, kind, parent_code, min_level,
  travel_seconds_from_parent, safe, is_starting
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (code) DO UPDATE SET
  display_name_key           = EXCLUDED.display_name_key,
  kind                       = EXCLUDED.kind,
  parent_code                = EXCLUDED.parent_code,
  min_level                  = EXCLUDED.min_level,
  travel_seconds_from_parent = EXCLUDED.travel_seconds_from_parent,
  safe                       = EXCLUDED.safe,
  is_starting                = EXCLUDED.is_starting,
  version                    = world.region.version + 1`

	for _, r := range orderedRegions {
		_, err = tx.Exec(ctx, q,
			r.Code, r.DisplayNameKey, r.Kind, nullableStr(r.ParentCode),
			r.MinLevel, r.TravelSecondsFromParent, r.Safe, r.IsStarting,
		)
		if err != nil {
			return fmt.Errorf("upsert region %s: %w", r.Code, err)
		}
	}

	// Links — full replace strategy: DELETE all + INSERT current set. Simpler
	// than per-row UPSERT and correct for the "yaml is authoritative" model.
	// Cheap: link count stays in the low hundreds.
	if _, err := tx.Exec(ctx, `DELETE FROM world.region_link`); err != nil {
		return fmt.Errorf("delete region_link: %w", err)
	}
	const linkQ = `INSERT INTO world.region_link (region_a_code, region_b_code, travel_seconds)
	               VALUES ($1, $2, $3)`
	for _, lk := range links {
		if _, err := tx.Exec(ctx, linkQ, lk.RegionACode, lk.RegionBCode, lk.TravelSeconds); err != nil {
			return fmt.Errorf("insert region_link %s↔%s: %w", lk.RegionACode, lk.RegionBCode, err)
		}
	}
	return tx.Commit(ctx)
}

func nullableStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
