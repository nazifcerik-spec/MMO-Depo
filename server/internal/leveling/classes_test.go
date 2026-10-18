package leveling_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/leveling"
)

func TestTermEval(t *testing.T) {
	term := leveling.Term{Base: 10, PerLevel: 2, Per: map[string]float64{"strength": 1.5, "luck": 0.5}}
	attr := map[string]float64{"strength": 10, "luck": 4}
	require.InDelta(t, 10+2*5+15+2, term.Eval(5, func(n string) float64 { return attr[n] }), 1e-9)
}

func validConfig() *leveling.Config {
	return &leveling.Config{
		LevelCap: 1000, Curve: leveling.Curve{Kind: "formula", BaselineXPPerMinute: 1000, A: 3, B: 0.001, C: 1.7},
		AttributePoints: leveling.AttributeRule{Every: 2, Amount: 1}, TalentPoints: leveling.TalentRule{From: 100, Every: 20, Until: 980, Amount: 1},
		EquipmentFamilies: []string{"sword", "shield", "staff"}, ArmorTypes: []string{"heavy", "cloth"},
		SoloAdaptation: leveling.SoloAdaptation{Ratio: 0.35},
		Classes: []leveling.Class{
			{Code: "knight", RoleTags: []string{"tank"}, Resource: &leveling.Resource{Kind: "rage", Max: 100},
				EquipmentPermissions: []string{"heavy_armor", "sword", "shield"}, PrimaryAttribute: "strength"},
			{Code: "shieldbearer", Parent: "knight", RequiredLevel: 100, RoleTags: []string{"tank"}, SupportPower: 0.4},
			{Code: "blademaster", Parent: "knight", RequiredLevel: 100, RoleTags: []string{"melee_dps"},
				EquipmentPermissions: []string{"heavy_armor", "sword"}},
		},
	}
}

func TestClassValidationAndQueries(t *testing.T) {
	c := validConfig()
	require.NoError(t, c.Prepare())
	require.Equal(t, []string{"heavy_armor", "sword", "shield"}, c.Permissions("shieldbearer"), "inherited")
	require.True(t, c.CanUse("shieldbearer", "shield"))
	require.False(t, c.CanUse("blademaster", "shield"), "own narrower list")
	require.False(t, c.CanUse("knight", "staff"))
	require.True(t, c.CanUse("knight", "helm"), "not a governed family")
	require.True(t, c.IsFamily("staff"))
	require.Equal(t, "rage", c.ResourceOf("blademaster").Kind, "from the base class")
	require.Nil(t, c.ResourceOf("nobody"))
	require.Equal(t, 0.4, c.SupportPower("shieldbearer"))
	require.Len(t, c.SoloAdaptationModifiers("shieldbearer", false), 1)
	require.InDelta(t, 0.14, c.SoloAdaptationModifiers("shieldbearer", false)[0].Amount, 1e-9)
	require.Nil(t, c.SoloAdaptationModifiers("shieldbearer", true))
	require.Nil(t, c.SoloAdaptationModifiers("blademaster", false))

	for name, mut := range map[string]func(c *leveling.Config){
		"role tag":       func(c *leveling.Config) { c.Classes[0].RoleTags = []string{"wizard"} },
		"permission":     func(c *leveling.Config) { c.Classes[0].EquipmentPermissions = []string{"bazooka"} },
		"resource kind":  func(c *leveling.Config) { c.Classes[0].Resource.Kind = "ammo" },
		"resource start": func(c *leveling.Config) { c.Classes[0].Resource.StartPct = 120 },
		"attribute":      func(c *leveling.Config) { c.Classes[0].PrimaryAttribute = "charm" },
		"formula": func(c *leveling.Config) {
			c.Classes[0].Derived = &leveling.Derived{HP: leveling.Term{Per: map[string]float64{"charm": 1}}}
		},
		"afk profile": func(c *leveling.Config) { c.Classes[2].AFKProfile = &leveling.AFKProfile{KillSpeed: 9} },
		"support":     func(c *leveling.Config) { c.Classes[1].SupportPower = 2 },
		"content tag": func(c *leveling.Config) { c.Classes[2].StrongIn = []string{"raids"} },
		"solo ratio":  func(c *leveling.Config) { c.SoloAdaptation.Ratio = 2 },
	} {
		c := validConfig()
		mut(c)
		require.Error(t, c.Prepare(), name)
	}
}

// DD-254: combat XP falls off for mobs far below the character.
func TestMobXPFactor(t *testing.T) {
	c := validConfig()
	c.MobLevelXP = []leveling.MobXPBand{{MaxBelow: 10, Factor: 1}, {MaxBelow: 25, Factor: 0.5}, {Factor: 0.1}}
	require.NoError(t, c.Prepare())
	for _, tc := range []struct {
		player, mob int
		want        float64
	}{
		{40, 50, 1}, {40, 40, 1}, {40, 30, 1}, {40, 29, 0.5}, {40, 15, 0.5}, {40, 14, 0.1}, {500, 1, 0.1},
		{0, 1, 1}, {40, 0, 1}, // unknown level (older snapshots): full
	} {
		require.Equal(t, tc.want, c.MobXPFactor(tc.player, tc.mob), "player %d mob %d", tc.player, tc.mob)
	}
	require.Equal(t, int32(1000), c.ScaleMobXP(2000, 40, 20))
	require.Equal(t, int32(200), c.ScaleMobXP(2000, 40, 1))

	for name, bands := range map[string][]leveling.MobXPBand{
		"open band not last": {{MaxBelow: 0, Factor: 1}, {MaxBelow: 10, Factor: 0.5}},
		"no open band":       {{MaxBelow: 10, Factor: 1}},
		"descending":         {{MaxBelow: 20, Factor: 1}, {MaxBelow: 10, Factor: 0.5}, {Factor: 0.1}},
		"factor":             {{MaxBelow: 10, Factor: 2}, {Factor: 0.1}},
	} {
		c := validConfig()
		c.MobLevelXP = bands
		require.Error(t, c.Prepare(), name)
	}
	empty := validConfig()
	require.NoError(t, empty.Prepare())
	require.Equal(t, 1.0, empty.MobXPFactor(500, 1), "no bands: always full")
}
