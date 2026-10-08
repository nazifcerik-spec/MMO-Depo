package item

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
)

func gearCatalog() *Catalog {
	str := func(s string) *string { return &s }
	return NewCatalog([]Template{
		{Code: "sword_10", MinLevel: 10, EquipSlot: str("mainhand"), ClassRestriction: str("warrior"), ItemClass: "weapon", Affixable: true},
		{Code: "staff_10", MinLevel: 10, EquipSlot: str("mainhand"), ClassRestriction: str("mage"), ItemClass: "weapon"},
		{Code: "helm_10", MinLevel: 10, EquipSlot: str("helm"), ItemClass: "armor"},
		{Code: "helm_5", MinLevel: 5, EquipSlot: str("helm"), ItemClass: "armor"},
		{Code: "helm_15", MinLevel: 15, EquipSlot: str("helm"), ItemClass: "armor"},
		{Code: "ore", MinLevel: 1, ItemClass: "material", Stackable: true},
	})
}

func TestGearPool_TierAndClass(t *testing.T) {
	cat := gearCatalog()
	require.Equal(t, []string{"helm_10", "sword_10"}, GearPool(cat, "warrior", 12))
	require.Equal(t, []string{"helm_10", "staff_10"}, GearPool(cat, "mage", 14))
	require.Equal(t, []string{"helm_5"}, GearPool(cat, "rogue", 9))
	require.Empty(t, GearPool(cat, "warrior", 3))
}

func TestRollGear(t *testing.T) {
	cat := gearCatalog()
	require.Empty(t, RollGear(cat, 1, activity.GearDrop{}))
	require.Empty(t, RollGear(cat, 1, activity.GearDrop{Class: "warrior", Level: 3, Guaranteed: 5}), "nothing in the tier")

	g := activity.GearDrop{Class: "warrior", Level: 12, Rolls: 2000, Pct: 4}
	a, b := RollGear(cat, 7, g), RollGear(cat, 7, g)
	require.Equal(t, a, b, "deterministic")
	require.InDelta(t, 80, len(a), 30, "about Pct%% of the rolls")
	seen := map[Rarity]int{}
	for _, d := range a {
		require.Contains(t, []string{"helm_10", "sword_10"}, d.TemplateCode)
		seen[d.Rarity]++
	}
	require.Greater(t, seen[RarityCommon], seen[RarityRare])

	boss := RollGear(cat, 9, activity.GearDrop{Class: "warrior", Level: 12, Guaranteed: 200, Boss: true})
	require.Len(t, boss, 200)
	for _, d := range boss {
		require.NotEqual(t, RarityCommon, d.Rarity, "boss gear is never common")
	}
}
