package item_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/itemset"
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
	require.NoError(t, itemset.Install(p.ItemSets)) // set pieces are raid drops only (DD-277)
	t.Cleanup(func() {
		_ = leveling.Install(prev)
		_ = itemset.Install(nil)
	})
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
			// Jewelry is for everyone from the level-25 tier (DD-276).
			jewelry := 0
			for _, sub := range []string{"necklace", "ring", "trinket"} {
				jewelry += families[sub]
				delete(families, sub)
			}
			require.Equal(t, map[bool]int{true: 3, false: 0}[lvl >= 25], jewelry, "%s Sv %d: jewelry", class, lvl)
			require.Len(t, families, 3, "%s Sv %d: armor + two weapon families", class, lvl)
		}
	}
}

// A raid boss drops one piece of the raid's set for the class's armor
// (DD-277); set pieces never come from the ordinary gear pool.
func TestShippedRaidSetPieces(t *testing.T) {
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	prev := leveling.Current()
	require.NoError(t, leveling.Install(p.Progression))
	require.NoError(t, itemset.Install(p.ItemSets))
	t.Cleanup(func() {
		_ = leveling.Install(prev)
		_ = itemset.Install(nil)
	})
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
		tmpls = append(tmpls, tm)
	}
	cat := item.NewCatalog(tmpls)
	for class, set := range map[string]string{"warrior": "north_keep_siege_plate", "mage": "north_keep_siege_cloth", "rogue": "north_keep_siege_leather"} {
		drop := activity.GearDrop{Class: class, Level: 100, Guaranteed: 2, Boss: true, SetSource: "north_keep_siege"}
		got := item.RollGear(cat, 7, drop)
		require.Len(t, got, 3, "two gear pieces and a set piece")
		s, ok := itemset.Current().Of(got[2].TemplateCode)
		require.True(t, ok, "%s: %s is a set piece", class, got[2].TemplateCode)
		require.Equal(t, set, s.Code)
		for _, d := range got[:2] {
			_, inSet := itemset.Current().Of(d.TemplateCode)
			require.False(t, inSet, "ordinary gear is not a set piece")
		}
		drop.SetSource = ""
		require.Len(t, item.RollGear(cat, 7, drop), 2)
	}
}
