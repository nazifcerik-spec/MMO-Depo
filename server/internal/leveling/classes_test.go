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

// Hunting XP falls off with the level gap either way (DD-281).
func TestHuntXPFalloff(t *testing.T) {
	f := leveling.HuntXPFalloff{FreePct: 0.05, Power: 4, Floor: 0.1}
	for _, tc := range []struct {
		player, mob int
		want        float64
	}{
		{40, 40, 1}, {40, 42, 1}, {40, 38, 1}, // 5% of 40: two levels free
		{40, 46, 0.6561}, {40, 34, 0.6561}, // 4 past free: (1 - 0.1)^4
		{10, 9, 1}, {10, 5, 0.1296}, {10, 3, 0.1}, // one level free at 10; the floor
		{500, 475, 1}, {500, 460, 0.88529281}, // 25 free at 500; 15 past: (1 - 0.03)^4
		{500, 1, 0.1}, {0, 5, 1}, // the floor; unknown level: full
	} {
		require.InDelta(t, tc.want, f.Factor(tc.player, tc.mob), 1e-9, "player %d mob %d", tc.player, tc.mob)
	}
	require.Equal(t, 1.0, leveling.HuntXPFalloff{}.Factor(500, 1), "no falloff configured: full")

	for name, bad := range map[string]leveling.HuntXPFalloff{
		"negative free":  {FreePct: -1},
		"negative power": {Power: -1},
		"floor > 1":      {Power: 2, Floor: 1.5},
	} {
		c := validConfig()
		c.HuntLevelXP = bad
		require.Error(t, c.Prepare(), name)
	}
}

// Dungeon XP follows the member's level (DD-303): mobs below pay as in
// hunting, mobs above up to UpCap more.
func TestDungeonXPScale(t *testing.T) {
	d := leveling.DungeonXPScale{Down: leveling.HuntXPFalloff{FreePct: 0.05, Power: 4, Floor: 0.1}, UpCap: 1.5}
	for _, tc := range []struct {
		player, mob int
		want        float64
	}{
		{40, 40, 1}, {40, 42, 1}, {40, 38, 1}, // the free band both ways
		{40, 34, 0.6561},         // below: the hunting falloff
		{40, 46, 1.1},            // above: 4 past free over 40
		{40, 80, 1.5}, {0, 5, 1}, // capped; unknown level: full
	} {
		require.InDelta(t, tc.want, d.Factor(tc.player, tc.mob), 1e-9, "player %d mob %d", tc.player, tc.mob)
	}
	require.Equal(t, 1.0, leveling.DungeonXPScale{}.Factor(40, 80), "not configured: full")
	c := validConfig()
	c.DungeonLevelXP = d
	require.NoError(t, c.Prepare())
	require.EqualValues(t, 110, c.ScaleDungeonXP(100, 40, 46))
	for name, bad := range map[string]leveling.DungeonXPScale{
		"negative cap": {UpCap: -1},
		"huge cap":     {UpCap: 9},
		"bad down":     {Down: leveling.HuntXPFalloff{Power: -1}},
	} {
		c := validConfig()
		c.DungeonLevelXP = bad
		require.Error(t, c.Prepare(), name)
	}
}
