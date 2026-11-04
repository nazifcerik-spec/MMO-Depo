package classskill

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func activeCodes(c *Catalog, path []string, level int) []string {
	var out []string
	for _, s := range c.Actives(path, level, "shield") {
		out = append(out, s.Code)
	}
	return out
}

// DD-352: on the way up the path a stronger skill replaces a weaker one of the
// same job; a weaker new one never replaces, other jobs are left alone.
func TestActivesKeepTheStrongestOfAJob(t *testing.T) {
	strike := func(code, spec string, level int16, power int32, cd int) Skill {
		return Skill{Code: code, Class: "mage", Spec: spec, UnlockLevel: level, Kind: "active", Effect: "strike",
			Power: power, CooldownTurns: cd, Config: Config{OnCooldown: true, Type: "DAMAGE"}}
	}
	c := NewCatalog([]Skill{
		strike("bolt", "", 1, 130, 1),
		strike("burst", "", 25, 230, 3),
		strike("element_bolt", "elementalist", 300, 250, 3),
		strike("weak_bolt", "elementalist", 400, 200, 3),
		strike("fire_ball", "pyromancer", 700, 320, 3),
		strike("other_branch", "cryomancer", 700, 900, 3),
		strike("converge", "elementalist", 600, 340, 5),
		{Code: "meteor", Class: "mage", Spec: "pyromancer", UnlockLevel: 950, Kind: "active", Effect: "strike", Power: 640,
			CooldownTurns: 16, Config: Config{Category: "ultimate", OnCooldown: true, Type: "DAMAGE"}},
		{Code: "finisher", Class: "mage", UnlockLevel: 75, Kind: "active", Effect: "strike", Power: 300, CooldownTurns: 2,
			Config: Config{TargetHPBelow: 35, Type: "DAMAGE"}},
		{Code: "calm", Class: "mage", UnlockLevel: 50, Kind: "passive"},
	})
	base := []string{"mage"}
	promo := []string{"mage", "elementalist"}
	spec := []string{"mage", "elementalist", "pyromancer"}

	require.ElementsMatch(t, []string{"bolt"}, activeCodes(c, base, 24))
	require.ElementsMatch(t, []string{"bolt", "burst", "finisher"}, activeCodes(c, base, 299))
	require.ElementsMatch(t, []string{"bolt", "element_bolt", "finisher"}, activeCodes(c, promo, 300),
		"the promotion's stronger skill replaces the base one")
	require.ElementsMatch(t, []string{"bolt", "element_bolt", "finisher"}, activeCodes(c, promo, 500),
		"a weaker new skill replaces nothing")
	require.ElementsMatch(t, []string{"bolt", "fire_ball", "finisher", "converge", "meteor"}, activeCodes(c, spec, 950),
		"another branch's skill does not count; bands and the ultimate stay")
	require.ElementsMatch(t, []string{"bolt", "burst", "finisher"}, activeCodes(c, promo, 299), "only at its level")
	require.Len(t, c.ForClass("mage"), 10, "the catalog still lists every skill")
}

// Same strength: the newest stays; a longer effect counts for more.
func TestStrongestComparesWholeEffect(t *testing.T) {
	guard := func(code string, level int16, power int32, turns int) Skill {
		return Skill{Code: code, Class: "monk", UnlockLevel: level, Kind: "active", Effect: "guard", Power: power,
			CooldownTurns: 8, Config: Config{OnCooldown: true, DurationTurns: turns}}
	}
	c := NewCatalog([]Skill{guard("short", 1, 40, 2), guard("long", 2, 30, 4), guard("same", 3, 60, 2)})
	require.Equal(t, []string{"long"}, activeCodes(c, []string{"monk"}, 2))
	require.Equal(t, []string{"same"}, activeCodes(c, []string{"monk"}, 3))
}
