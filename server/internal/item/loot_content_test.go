package item_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/leveling"
)

// The shipped gathering tables (DD-244): every profession starts at a
// level-1 resource, and a combat hide exists for every mob level.
func TestShippedGatheringTables(t *testing.T) {
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	var tmpls []item.Template
	for _, it := range p.Items {
		if it.Gather != nil {
			require.Equal(t, int32(999), it.MaxStack, "%s: gathered stacks hold 999", it.Code)
			tmpls = append(tmpls, item.Template{Code: it.Code, GatherKind: it.Gather.Kind,
				GatherLevel: int16(it.Gather.Level), GatherTier: int16(it.Gather.Tier)})
		}
	}
	cat := item.NewCatalog(tmpls)
	for kind := range contentpack.GatherKinds {
		nodes := cat.Nodes(kind)
		require.NotEmpty(t, nodes, kind)
		require.Equal(t, int16(1), nodes[0].GatherLevel, "%s needs a level-1 resource", kind)
	}
	for lvl := 1; lvl <= 30; lvl++ {
		require.NotEmpty(t, item.CombatHides(cat, 1, lvl, 1, 1, 0), "mob level %d", lvl)
	}
}

// Shipped gear drops only what the class may equip (DD-271, DD-272): its
// armor type in all five body slots and its two weapon families.
func TestShippedGearPoolHonoursPermissions(t *testing.T) {
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	prev := leveling.Current()
	require.NoError(t, leveling.Install(p.Progression))
	t.Cleanup(func() { _ = leveling.Install(prev) })
	var tmpls []item.Template
	for _, it := range p.Items {
		if it.EquipSlot == "" {
			continue
		}
		slot, class := it.EquipSlot, it.ClassRestriction
		tm := item.Template{Code: it.Code, ItemClass: it.ItemClass, ItemSubclass: it.ItemSubclass,
			MinLevel: int16(it.MinLevel), EquipSlot: &slot}
		if class != "" {
			tm.ClassRestriction = &class
		}
		if it.Requirements != nil {
			tm.Requirements = &item.Requirements{Spec: it.Requirements.Spec}
		}
		tmpls = append(tmpls, tm)
	}
	cat := item.NewCatalog(tmpls)
	armor := map[string]string{"warrior": "plate_armor", "mage": "cloth_armor", "rogue": "leather_armor"}
	for class, want := range armor {
		for _, lvl := range []int{1, 12, 300, 1000} {
			pool := item.GearPool(cat, class, lvl)
			families := map[string]int{}
			for _, code := range pool {
				tm, err := cat.Get(code)
				require.NoError(t, err)
				families[tm.ItemSubclass]++
				require.True(t, p.Progression.CanUse(class, tm.ItemSubclass), "%s got %s", class, code)
			}
			require.Equal(t, 5, families[want], "%s Sv %d: five %s pieces", class, lvl, want)
			require.Len(t, families, 3, "%s Sv %d: armor + two weapon families", class, lvl)
		}
	}
}
