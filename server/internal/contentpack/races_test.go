package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-349: the shipped race table, exactly as designed.
func TestRaceTableShipped(t *testing.T) {
	p := shippedPack(t)
	type want struct {
		classes    []string
		passives   map[string]float64
		profession string
	}
	table := map[string]want{
		"human":    {[]string{"warrior", "rogue", "ranger", "mage", "monk", "cleric", "paladin", "bard"}, map[string]float64{"str": .01, "int": .01}, "cooking"},
		"elf":      {[]string{"warrior", "rogue", "ranger", "mage", "druid", "bard"}, map[string]float64{"str": .01, "agi": .02}, "enchanting"},
		"dwarf":    {[]string{"warrior", "cleric", "paladin", "bard"}, map[string]float64{"str": .01, "vit": .01}, "blacksmithing"},
		"orc":      {[]string{"warrior", "mage", "shaman"}, map[string]float64{"str": .02, "vit": .01}, "smelting"},
		"troll":    {[]string{"rogue", "ranger", "bard", "shaman"}, map[string]float64{"str": .01, "agi": .01}, "leatherworking"},
		"undead":   {[]string{"warrior", "rogue", "ranger", "mage"}, map[string]float64{"str": .01, "agi": .01, "luck": .01}, "alchemy"},
		"dark_elf": {[]string{"warrior", "rogue", "ranger", "mage"}, map[string]float64{"str": .01, "agi": .02}, "alchemy"},
		"drakari":  {[]string{"mage", "monk", "cleric", "bard", "shaman"}, map[string]float64{"int": .01, "cha": .02}, "enchanting"},
		"tauren":   {[]string{"warrior", "monk", "druid", "shaman"}, map[string]float64{"str": .01, "vit": .02}, "woodworking"},
		"gnome":    {[]string{"rogue", "mage", "bard"}, map[string]float64{"int": .01, "agi": .01}, "jewelcrafting"},
	}
	require.Len(t, p.Progression.Races, len(table))
	for code, w := range table {
		r, ok := p.Progression.Races[code]
		require.True(t, ok, code)
		require.ElementsMatch(t, w.classes, r.Classes, code)
		got := map[string]float64{}
		for _, m := range r.Passives {
			require.Equal(t, "mult_add", m.Kind, code)
			got[m.Target] = m.Amount
		}
		require.Equal(t, w.passives, got, code)
		require.Equal(t, map[string]float64{w.profession: .02}, r.ProfessionXP, code)
	}
}

// Every base class stays reachable, and an unlisted combination is refused.
func TestRaceAllows(t *testing.T) {
	p := shippedPack(t)
	for _, class := range baseClasses {
		reachable := false
		for race := range p.Progression.Races {
			reachable = reachable || p.Progression.RaceAllows(race, class)
		}
		require.True(t, reachable, class)
	}
	require.False(t, p.Progression.RaceAllows("human", "shaman"))
	require.True(t, p.Progression.RaceAllows("tauren", "shaman"))
	require.True(t, p.Progression.RaceAllows("unknown", "shaman"), "a race not listed has no limits")
}
