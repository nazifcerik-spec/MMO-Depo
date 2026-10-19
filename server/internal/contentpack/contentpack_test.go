package contentpack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/inventory"
)

func writeTemp(t *testing.T, dir, name, body string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644))
}

func TestLoadItemsFile_ParsesCanonicalYAML(t *testing.T) {
	dir := t.TempDir()
	writeTemp(t, dir, "items.yaml", `
items:
  - code: iron_ore
    display_name_key: item.iron_ore
    item_class: material
    item_subclass: ore
    min_level: 1
    base_value_copper: 5
    base_weight_g: 200
    stackable: true
    max_stack: 99
    bind_kind: boe
  - code: sword_iron
    display_name_key: item.sword_iron
    item_class: weapon
    item_subclass: sword
    min_level: 5
    base_value_copper: 100
    base_weight_g: 2500
    bind_kind: boe
    affixable: true
    equip_slot: mainhand
    stat_bonuses:
      - { kind: flat, target: attack, amount: 10 }
      - { kind: flat, target: str,    amount: 3 }
`)
	items, err := loadItemsFile(filepath.Join(dir, "items.yaml"))
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, "iron_ore", items[0].Code)
	require.True(t, items[0].Stackable)

	sword := items[1]
	require.Equal(t, "sword_iron", sword.Code)
	require.Equal(t, "mainhand", sword.EquipSlot)
	require.Len(t, sword.StatBonuses, 2)
	require.Equal(t, "flat", sword.StatBonuses[0].Kind)
	require.Equal(t, "attack", sword.StatBonuses[0].Target)
	require.InDelta(t, 10.0, sword.StatBonuses[0].Amount, 0.001)
}

func TestLoadItemsFile_MissingFileIsTolerated(t *testing.T) {
	items, err := loadItemsFile(filepath.Join(t.TempDir(), "nope.yaml"))
	require.NoError(t, err)
	require.Nil(t, items)
}

func TestValidateItems_RejectsDuplicateCode(t *testing.T) {
	err := validateItems([]Item{
		{Code: "iron_ore", ItemClass: "material"},
		{Code: "iron_ore", ItemClass: "material"},
	})
	require.ErrorContains(t, err, "duplicate code")
}

func TestValidateItems_RejectsUnknownItemClass(t *testing.T) {
	err := validateItems([]Item{{Code: "x", ItemClass: "unicorn"}})
	require.ErrorContains(t, err, "unknown item_class")
}

func TestValidateItems_RejectsUnknownEquipSlot(t *testing.T) {
	err := validateItems([]Item{{Code: "x", ItemClass: "weapon", EquipSlot: "forehead"}})
	require.ErrorContains(t, err, "unknown equip_slot")
}

func TestValidateItems_RejectsUnknownModifierKind(t *testing.T) {
	err := validateItems([]Item{{
		Code: "x", ItemClass: "weapon",
		StatBonuses: []StatBonus{{Kind: "curse", Target: "str", Amount: 1}},
	}})
	require.ErrorContains(t, err, "unknown kind")
}

func TestValidateItems_RejectsUnknownStatTarget(t *testing.T) {
	err := validateItems([]Item{{
		Code: "x", ItemClass: "weapon",
		StatBonuses: []StatBonus{{Kind: "flat", Target: "swagger", Amount: 1}},
	}})
	require.ErrorContains(t, err, "unknown target")
}

func TestValidateItems_StackableRequiresMaxStack(t *testing.T) {
	err := validateItems([]Item{{
		Code: "x", ItemClass: "material", Stackable: true, MaxStack: 0,
	}})
	require.ErrorContains(t, err, "max_stack")
}

func TestValidateItems_AcceptsCanonicalCases(t *testing.T) {
	// Every valid combination should pass — regression guard for over-strict
	// validation.
	err := validateItems([]Item{
		{Code: "mat", ItemClass: "material", Stackable: true, MaxStack: 99, BindKind: "boe"},
		{Code: "wpn", ItemClass: "weapon", Affixable: true, EquipSlot: "mainhand",
			StatBonuses: []StatBonus{
				{Kind: "flat", Target: "attack", Amount: 10},
				{Kind: "mult_add", Target: "str", Amount: 0.10},
				{Kind: "mult_prod", Target: "attack", Amount: 0.05},
				{Kind: "clamp_max", Target: "crit_pct", Amount: 100},
			}},
		{Code: "arm", ItemClass: "armor", EquipSlot: "chest"},
		{Code: "pot", ItemClass: "consumable", Stackable: true, MaxStack: 20},
		{Code: "trk", ItemClass: "misc"},
	})
	require.NoError(t, err)
}

func TestLoadVendorsFile_ParsesCanonicalYAML(t *testing.T) {
	dir := t.TempDir()
	writeTemp(t, dir, "vendors.yaml", `
vendors:
  - code: general_ironhold
    display_name_key: vendor.general_ironhold
    vendor_type: general
    buy_multiplier: 1.000
    sell_multiplier: 0.250
    sells:
      - { template_code: potion_hp_small }
      - { template_code: sword_iron, price_override_copper: 150 }
`)
	vendors, err := loadVendorsFile(filepath.Join(dir, "vendors.yaml"))
	require.NoError(t, err)
	require.Len(t, vendors, 1)
	v := vendors[0]
	require.Equal(t, "general_ironhold", v.Code)
	require.InDelta(t, 0.25, v.SellMultiplier, 0.001)
	require.Len(t, v.Sells, 2)
	require.Nil(t, v.Sells[0].PriceOverrideCopper)
	require.NotNil(t, v.Sells[1].PriceOverrideCopper)
	require.Equal(t, int64(150), *v.Sells[1].PriceOverrideCopper)
}

func TestValidateVendors_RejectsDuplicateAndNegative(t *testing.T) {
	err := validateVendors([]Vendor{{Code: "a"}, {Code: "a"}})
	require.ErrorContains(t, err, "duplicate code")

	err = validateVendors([]Vendor{{Code: "a", BuyMultiplier: -1}})
	require.ErrorContains(t, err, "negative multiplier")
}

// The shipped content must validate and every equippable item must use a
// slot the equip service accepts (regression: head/feet vs helm/boots).
func TestShippedItemsValidate(t *testing.T) {
	items, err := loadItemsFile(filepath.Join("..", "..", "..", "content", "items.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, items)
	require.NoError(t, validateItems(items))
	for _, it := range items {
		if it.EquipSlot != "" {
			require.True(t, inventory.KnownEquipSlots[it.EquipSlot], "%s: equip_slot %q", it.Code, it.EquipSlot)
		}
	}
}

func TestValidateItems_UseEffect(t *testing.T) {
	ok := []Item{{Code: "pot", ItemClass: "consumable", UseEffect: &UseEffect{Kind: "heal_hp", Amount: 60}}}
	require.NoError(t, validateItems(ok))

	for name, it := range map[string]Item{
		"non-consumable": {Code: "x", ItemClass: "weapon", UseEffect: &UseEffect{Kind: "heal_hp", Amount: 5}},
		"unknown kind":   {Code: "x", ItemClass: "consumable", UseEffect: &UseEffect{Kind: "restore_mp", Amount: 5}},
		"zero amount":    {Code: "x", ItemClass: "consumable", UseEffect: &UseEffect{Kind: "heal_hp"}},
		"buff too big":   {Code: "x", ItemClass: "consumable", UseEffect: &UseEffect{Kind: "heal_hp", Amount: 5, Buff: &UseBuff{AttackPct: 51, Minutes: 60}}},
		"empty buff":     {Code: "x", ItemClass: "consumable", UseEffect: &UseEffect{Kind: "heal_hp", Amount: 5, Buff: &UseBuff{Minutes: 60}}},
		"buff no time":   {Code: "x", ItemClass: "consumable", UseEffect: &UseEffect{Kind: "heal_hp", Amount: 5, Buff: &UseBuff{AttackPct: 5}}},
	} {
		require.Error(t, validateItems([]Item{it}), name)
	}
}
