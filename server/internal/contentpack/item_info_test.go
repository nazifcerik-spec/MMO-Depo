package contentpack

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/leveling"
)

// Item tags (DD-339) follow the drop and crafting rules of the shipped content.
func TestItemInfoShipped(t *testing.T) {
	p := shippedPack(t)
	cfg := p.Progression
	if cfg == nil {
		cfg = leveling.Current()
	}
	ii := BuildItemInfo(p, cfg)
	require.Len(t, ii, len(p.Items))

	sword := ii["sword_iron"]
	require.Equal(t, int16(5), sword.Level)
	require.Contains(t, sword.Classes, "warrior")
	require.NotEmpty(t, sword.Mobs, "random gear drops from monsters whose ladder band reaches its tier")
	for _, m := range sword.Mobs {
		require.Greater(t, m.Level+combat.LadderStep, sword.Level)
	}
	require.Contains(t, sword.Crafted, ItemCraft{Profession: "blacksmithing", Skill: 1})

	require.Contains(t, ii["potion_hp_small"].Vendors, "general_ironhold")
	ore := ii["copper_ore"]
	require.Equal(t, "mining", ore.Gather.Kind)
	require.Contains(t, ore.UsedBy, "smelting")

	for code, info := range ii {
		if strings.HasPrefix(code, "pet_egg_") {
			require.NotEmpty(t, info.Dungeons, "%s drops from a raid boss", code)
			for _, d := range info.Dungeons {
				require.True(t, d.Raid, code)
			}
		}
		if info.Set != "" {
			require.Len(t, info.Dungeons, 1, "%s comes from its raid", code)
			require.Empty(t, info.Mobs, "%s never drops on hunts", code)
		}
		if info.Slot != "" && info.Set == "" {
			require.True(t, info.AnyClass || len(info.Classes) > 0, "%s has users", code)
			// Every piece of equipment can be obtained somehow (DD-340).
			require.True(t, len(info.Mobs)+len(info.Dungeons)+len(info.Crafted)+len(info.Vendors) > 0,
				"%s has no source", code)
		}
	}
}
