package contentpack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Every file loader: a missing file is "no content", malformed YAML and
// content that fails validation are load errors (a reload must not half-apply).
func TestLoaders_MissingMalformedInvalid(t *testing.T) {
	loaders := map[string]struct {
		load    func(string) error
		invalid string
	}{
		"professions": {func(p string) error { _, err := loadProfessionsFile(p); return err }, "professions:\n  - code: ''\n"},
		"recipes":     {func(p string) error { _, err := loadRecipesFile(p); return err }, "recipes:\n  - code: r\n"},
		"regions":     {func(p string) error { _, _, err := loadRegionsFile(p); return err }, "regions:\n  - code: a\n    kind: moon\n"},
		"dungeons":    {func(p string) error { _, err := loadDungeonsFile(p); return err }, "dungeons:\n  - code: d\n"},
		"mobs":        {func(p string) error { _, err := loadMobsFile(p); return err }, "mobs:\n  - code: m\n"},
		"vendors":     {func(p string) error { _, err := loadVendorsFile(p); return err }, "vendors:\n  - code: ''\n"},
		"items":       {func(p string) error { _, err := loadItemsFile(p); return err }, "items:\n  - code: ''\n"},
		"skills":      {func(p string) error { _, err := loadSkillsFile(p); return err }, ""},
		"quests":      {func(p string) error { _, err := loadQuestsFile(p); return err }, ""},
	}
	for name, l := range loaders {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, l.load(filepath.Join(dir, "missing.yaml")), "missing file is tolerated")

			bad := filepath.Join(dir, "bad.yaml")
			require.NoError(t, os.WriteFile(bad, []byte("{not: [yaml"), 0o600))
			require.Error(t, l.load(bad))

			if l.invalid != "" {
				inv := filepath.Join(dir, "invalid.yaml")
				require.NoError(t, os.WriteFile(inv, []byte(l.invalid), 0o600))
				require.Error(t, l.load(inv))
			}

			// A directory where a file is expected is a read error, not "missing".
			require.Error(t, l.load(dir))
		})
	}
}

func TestValidateProfessions(t *testing.T) {
	require.NoError(t, validateProfessions([]Profession{{Code: "a", MaxLevel: 60}}))
	for name, ps := range map[string][]Profession{
		"duplicate": {{Code: "a", MaxLevel: 1}, {Code: "a", MaxLevel: 1}},
		"max level": {{Code: "a", MaxLevel: 101}},
		"zero":      {{Code: "a"}},
	} {
		require.Error(t, validateProfessions(ps), name)
	}
}

func TestValidateRecipes(t *testing.T) {
	ok := func() Recipe {
		return Recipe{Code: "r", ProfessionCode: "p", MinSkill: 1, OutputTemplateCode: "o", OutputStack: 1, CraftSeconds: 30,
			Inputs: []RecipeInput{{TemplateCode: "a", Quantity: 1}}}
	}
	require.NoError(t, validateRecipes([]Recipe{ok()}))
	for name, mut := range map[string]func(*Recipe){
		"empty code":      func(r *Recipe) { r.Code = "" },
		"no profession":   func(r *Recipe) { r.ProfessionCode = "" },
		"min skill":       func(r *Recipe) { r.MinSkill = 0 },
		"no output":       func(r *Recipe) { r.OutputTemplateCode = "" },
		"output stack":    func(r *Recipe) { r.OutputStack = 0 },
		"negative gold":   func(r *Recipe) { r.GoldCostCopper = -1 },
		"negative xp":     func(r *Recipe) { r.ProfessionXPReward = -1 },
		"no craft time":   func(r *Recipe) { r.CraftSeconds = 0 },
		"craft too long":  func(r *Recipe) { r.CraftSeconds = 3601 },
		"no inputs":       func(r *Recipe) { r.Inputs = nil },
		"empty input":     func(r *Recipe) { r.Inputs[0].TemplateCode = "" },
		"zero quantity":   func(r *Recipe) { r.Inputs[0].Quantity = 0 },
		"duplicate input": func(r *Recipe) { r.Inputs = append(r.Inputs, RecipeInput{TemplateCode: "a", Quantity: 1}) },
	} {
		r := ok()
		mut(&r)
		require.Error(t, validateRecipes([]Recipe{r}), name)
	}
	require.Error(t, validateRecipes([]Recipe{ok(), ok()}), "duplicate code")
}

func TestValidateRegions(t *testing.T) {
	regions := func() []Region {
		return []Region{{Code: "a", Kind: "city", IsStarting: true}, {Code: "b", Kind: "zone", ParentCode: "a"}, {Code: "c", Kind: "wilderness"}}
	}
	require.NoError(t, validateRegions(regions(), []RegionLink{{RegionACode: "a", RegionBCode: "c", TravelSeconds: 60}}))
	for name, mut := range map[string]func([]Region) []Region{
		"empty code":  func(r []Region) []Region { r[0].Code = ""; return r },
		"duplicate":   func(r []Region) []Region { return append(r, r[0]) },
		"kind":        func(r []Region) []Region { r[0].Kind = "moon"; return r },
		"travel":      func(r []Region) []Region { r[1].TravelSecondsFromParent = 3601; return r },
		"two starts":  func(r []Region) []Region { r[1].IsStarting = true; return r },
		"self parent": func(r []Region) []Region { r[1].ParentCode = "b"; return r },
		"no parent":   func(r []Region) []Region { r[1].ParentCode = "zz"; return r },
	} {
		require.Error(t, validateRegions(mut(regions()), nil), name)
	}
	for name, l := range map[string][]RegionLink{
		"empty":     {{RegionACode: "", RegionBCode: "b"}},
		"self":      {{RegionACode: "a", RegionBCode: "a"}},
		"order":     {{RegionACode: "c", RegionBCode: "a"}},
		"unknown a": {{RegionACode: "0", RegionBCode: "a"}},
		"unknown b": {{RegionACode: "a", RegionBCode: "z"}},
		"duplicate": {{RegionACode: "a", RegionBCode: "c"}, {RegionACode: "a", RegionBCode: "c"}},
		"travel":    {{RegionACode: "a", RegionBCode: "c", TravelSeconds: -1}},
	} {
		require.Error(t, validateRegions(regions(), l), name)
	}
}

func TestValidateDungeons(t *testing.T) {
	ok := func() Dungeon {
		return Dungeon{Code: "d", DurationSeconds: 600, BossMobCode: "boss", Waves: []DungeonWave{{MobCode: "m", Count: 3}}}
	}
	require.NoError(t, validateDungeons([]Dungeon{ok()}))
	for name, mut := range map[string]func(*Dungeon){
		"empty code": func(d *Dungeon) { d.Code = "" },
		"duration":   func(d *Dungeon) { d.DurationSeconds = 0 },
		"no waves":   func(d *Dungeon) { d.Waves = nil },
		"no boss":    func(d *Dungeon) { d.BossMobCode = "" },
		"wave mob":   func(d *Dungeon) { d.Waves[0].MobCode = "" },
		"wave count": func(d *Dungeon) { d.Waves[0].Count = 33 },
		"bonus":      func(d *Dungeon) { d.BonusXP = -1 },
	} {
		d := ok()
		mut(&d)
		require.Error(t, validateDungeons([]Dungeon{d}), name)
	}
	require.Error(t, validateDungeons([]Dungeon{ok(), ok()}), "duplicate")
}

func TestValidateMobs(t *testing.T) {
	ok := func() Mob { return Mob{Code: "m", HPMax: 10, HitPct: 80, GoldCopperMin: 1, GoldCopperMax: 2} }
	require.NoError(t, validateMobs([]Mob{ok()}))
	for name, mut := range map[string]func(*Mob){
		"empty code": func(m *Mob) { m.Code = "" },
		"hp":         func(m *Mob) { m.HPMax = 0 },
		"attack":     func(m *Mob) { m.Attack = -1 },
		"crit":       func(m *Mob) { m.CritPct = 101 },
		"hit":        func(m *Mob) { m.HitPct = -1 },
		"evasion":    func(m *Mob) { m.EvasionPct = 101 },
		"gold":       func(m *Mob) { m.GoldCopperMax = 0 },
	} {
		m := ok()
		mut(&m)
		require.Error(t, validateMobs([]Mob{m}), name)
	}
	require.Error(t, validateMobs([]Mob{ok(), ok()}), "duplicate")
}
