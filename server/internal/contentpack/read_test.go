package contentpack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/itemset"
)

// The shipped content directory parses, validates and cross-references.
func TestShippedPackReads(t *testing.T) {
	p, err := Read("../../../content")
	require.NoError(t, err)
	require.NotEmpty(t, p.Items)
	require.NotEmpty(t, p.Mobs)
}

func TestValidateRefs_Rejects(t *testing.T) {
	base := func() *Pack {
		return &Pack{
			Items:       []Item{{Code: "a"}, {Code: "b"}},
			Mobs:        []Mob{{Code: "m", RegionCode: "g"}},
			Regions:     []Region{{Code: "g"}},
			Professions: []Profession{{Code: "p"}},
			Vendors:     []Vendor{{Code: "v", Sells: []VendorSell{{TemplateCode: "a"}}}},
			Recipes:     []Recipe{{Code: "r", ProfessionCode: "p", OutputTemplateCode: "b", Inputs: []RecipeInput{{TemplateCode: "a"}}}},
			Dungeons:    []Dungeon{{Code: "d", BossMobCode: "m", Waves: []DungeonWave{{MobCode: "m"}}}},
		}
	}
	check := func(p *Pack) error {
		return validateRefs(p, refs{
			mobs:  codeSet(p.Mobs, func(m Mob) string { return m.Code }),
			items: codeSet(p.Items, func(i Item) string { return i.Code }),
		})
	}
	require.NoError(t, check(base()))
	for name, mut := range map[string]func(*Pack){
		"vendor item":       func(p *Pack) { p.Vendors[0].Sells[0].TemplateCode = "zz" },
		"recipe profession": func(p *Pack) { p.Recipes[0].ProfessionCode = "zz" },
		"recipe output":     func(p *Pack) { p.Recipes[0].OutputTemplateCode = "zz" },
		"recipe input":      func(p *Pack) { p.Recipes[0].Inputs[0].TemplateCode = "zz" },
		"dungeon boss":      func(p *Pack) { p.Dungeons[0].BossMobCode = "zz" },
		"dungeon wave":      func(p *Pack) { p.Dungeons[0].Waves[0].MobCode = "zz" },
		"mob region":        func(p *Pack) { p.Mobs[0].RegionCode = "zz" },
		"mob no region":     func(p *Pack) { p.Mobs[0].RegionCode = "" },
		"raid boss region":  func(p *Pack) { p.Mobs[0].RaidBoss = true },
	} {
		p := base()
		mut(p)
		require.Error(t, check(p), name)
	}
	_, err := Read("")
	require.Error(t, err)
	_, err = Read("contentpack.go")
	require.Error(t, err)
	_, err = Read("missing-dir")
	require.Error(t, err)
}

// Sets need a raid source, equippable pieces and known modifiers (DD-277).
func TestValidateSets_Rejects(t *testing.T) {
	base := func() *Pack {
		return &Pack{
			Items:    []Item{{Code: "helm", EquipSlot: "helm"}, {Code: "chest", EquipSlot: "chest"}, {Code: "ore"}},
			Dungeons: []Dungeon{{Code: "raid", GroupMax: 5}, {Code: "solo", GroupMax: 1}},
			ItemSets: []itemset.Set{{Code: "s", Source: "raid", Pieces: []string{"helm", "chest"},
				Bonuses: []itemset.Bonus{{Pieces: 2, Modifiers: []itemset.Modifier{{Kind: "flat", Target: "attack", Amount: 1}}}}}},
		}
	}
	require.NoError(t, validateSets(base()))
	for name, mut := range map[string]func(*Pack){
		"solo source":    func(p *Pack) { p.ItemSets[0].Source = "solo" },
		"unknown source": func(p *Pack) { p.ItemSets[0].Source = "zz" },
		"material piece": func(p *Pack) { p.ItemSets[0].Pieces[1] = "ore" },
		"unknown piece":  func(p *Pack) { p.ItemSets[0].Pieces[1] = "zz" },
		"unknown kind":   func(p *Pack) { p.ItemSets[0].Bonuses[0].Modifiers[0].Kind = "zz" },
		"unknown target": func(p *Pack) { p.ItemSets[0].Bonuses[0].Modifiers[0].Target = "zz" },
	} {
		p := base()
		mut(p)
		require.Error(t, validateSets(p), name)
	}
	_, err := loadSetsFile("missing.yaml")
	require.NoError(t, err, "no sets file is fine")
	dir := t.TempDir()
	bad := filepath.Join(dir, "sets.yaml")
	require.NoError(t, os.WriteFile(bad, []byte("sets: [{code: a, pieces: [x]}]"), 0o600))
	_, err = loadSetsFile(bad)
	require.Error(t, err, "a one-piece set is invalid")
	require.NoError(t, os.WriteFile(bad, []byte("sets: {"), 0o600))
	_, err = loadSetsFile(bad)
	require.Error(t, err)
}

// Every region has a place on the travel map and every position a region
// (DD-279).
func TestApplyMap(t *testing.T) {
	regions := func() []Region { return []Region{{Code: "a"}, {Code: "b"}} }
	dir := t.TempDir()
	path := filepath.Join(dir, "map.yaml")
	write := func(s string) { require.NoError(t, os.WriteFile(path, []byte(s), 0o600)) }

	rs := regions()
	require.NoError(t, applyMap(filepath.Join(dir, "none.yaml"), rs), "no map file: regions stay at 0,0")
	write("positions: {a: [10, 20], b: [30, 40]}")
	rs = regions()
	require.NoError(t, applyMap(path, rs))
	require.Equal(t, [2]int16{30, 40}, [2]int16{rs[1].MapX, rs[1].MapY})
	for name, body := range map[string]string{
		"missing region": "positions: {a: [1, 1]}",
		"unknown region": "positions: {a: [1, 1], b: [1, 1], c: [1, 1]}",
		"out of bounds":  "positions: {a: [1, 1], b: [1, 4001]}",
		"bad yaml":       "positions: {",
	} {
		write(body)
		require.Error(t, applyMap(path, regions()), name)
	}
	p, err := Read("../../../content")
	require.NoError(t, err)
	for _, r := range p.Regions {
		require.True(t, r.MapX > 0 && r.MapY > 0, "%s is on the map", r.Code)
	}
}
