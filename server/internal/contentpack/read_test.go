package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"
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
