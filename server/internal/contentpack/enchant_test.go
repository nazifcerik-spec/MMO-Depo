package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/item"
)

// Every crafting profession keeps a recipe that still teaches at every
// skill (DD-294): grey recipes pay no XP, so a gap in a ladder would stall
// the profession for good.
func TestRecipeLaddersNeverGoGrey(t *testing.T) {
	p, err := Read("../../../content")
	require.NoError(t, err)
	skills := map[string][]int16{}
	for _, r := range p.Recipes {
		skills[r.ProfessionCode] = append(skills[r.ProfessionCode], r.MinSkill)
	}
	for _, prof := range p.Professions {
		if prof.Kind == "gathering" {
			continue
		}
		require.NotEmpty(t, skills[prof.Code], prof.Code)
		for s := int16(1); s <= prof.MaxLevel; s++ {
			teaches := false
			for _, r := range skills[prof.Code] {
				if color, _ := crafting.Difficulty(r, s); r <= s && color != crafting.DifficultyGrey {
					teaches = true
					break
				}
			}
			require.True(t, teaches, "%s: nothing to learn from at skill %d", prof.Code, s)
		}
	}
}

// Enchanting (DD-293) ships ten tiers of scrolls for every class's
// attributes, the combat stats and each profession, each with a recipe.
func TestShippedEnchanting(t *testing.T) {
	p, err := Read("../../../content")
	require.NoError(t, err)
	targets := map[string]bool{}
	scrolls := map[string]bool{}
	for _, it := range p.Items {
		if it.Enchant == nil {
			continue
		}
		scrolls[it.Code] = true
		for _, m := range it.Enchant.Modifiers {
			targets[m.Target] = true
		}
	}
	require.Len(t, scrolls, 210)
	for _, want := range []string{"str", "agi", "int", "cha", "vit", "attack", "defense", "hp_max", "crit_pct",
		"gather_yield_mining", "gather_yield_herbalism", "gather_yield_woodcutting", "gather_yield_fishing",
		"gather_yield_hunting", "craft_speed"} {
		require.True(t, targets[want], want)
	}
	made := map[string]bool{}
	for _, r := range p.Recipes {
		if r.ProfessionCode == "enchanting" {
			made[r.OutputTemplateCode] = true
		}
	}
	require.Equal(t, scrolls, made, "every scroll has a recipe")
}

func TestValidateEnchant(t *testing.T) {
	ok := Item{Code: "s", ItemClass: "consumable", ItemSubclass: "enchant_scroll", Enchant: &item.Enchant{
		Slots: []string{"weapon"}, Modifiers: []item.EnchantModifier{{Kind: "mult_add", Target: "attack", Amount: 0.02}}}}
	require.NoError(t, validateItems([]Item{ok}))
	for name, mutate := range map[string]func(*Item){
		"not a scroll":   func(i *Item) { i.ItemSubclass = "potion" },
		"no slots":       func(i *Item) { i.Enchant.Slots = nil },
		"unknown group":  func(i *Item) { i.Enchant.Slots = []string{"boots"} },
		"unknown target": func(i *Item) { i.Enchant.Modifiers[0].Target = "charm" },
		"too big":        func(i *Item) { i.Enchant.Modifiers[0].Amount = 0.9 },
		"flat profession": func(i *Item) {
			i.Enchant.Modifiers[0] = item.EnchantModifier{Kind: "flat", Target: "craft_speed", Amount: 1}
		},
		"bad kind":        func(i *Item) { i.Enchant.Modifiers[0].Kind = "set" },
		"with use effect": func(i *Item) { i.UseEffect = &UseEffect{Kind: "heal_hp", Amount: 5} },
	} {
		bad := ok
		e := *ok.Enchant
		e.Slots = append([]string(nil), e.Slots...)
		e.Modifiers = append([]item.EnchantModifier(nil), e.Modifiers...)
		bad.Enchant = &e
		mutate(&bad)
		require.Error(t, validateItems([]Item{bad}), name)
	}
}

func TestEnchantFits(t *testing.T) {
	w := item.Enchant{Slots: []string{"weapon"}}
	require.True(t, w.Fits("mainhand"))
	require.True(t, w.Fits("offhand"))
	require.False(t, w.Fits("chest"))
	j := item.Enchant{Slots: []string{"jewelry"}}
	require.True(t, j.Fits("ring"))
	require.False(t, j.Fits(""))
	require.True(t, item.Enchant{Slots: []string{"any"}}.Fits("boots"))
}
