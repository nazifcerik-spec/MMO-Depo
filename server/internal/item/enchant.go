package item

import "slices"

// Enchant is what an enchanting scroll puts on a piece of gear (DD-293):
// its modifiers count as part of the item while it is worn. A piece holds
// one enchant; a new one replaces it.
type Enchant struct {
	// Slots: the gear groups the scroll fits (EnchantGroups keys).
	Slots     []string          `json:"slots" yaml:"slots"`
	Modifiers []EnchantModifier `json:"modifiers" yaml:"modifiers"`
}

// EnchantModifier mirrors skill.Modifier (item does not import skill).
type EnchantModifier struct {
	Kind   string  `json:"kind" yaml:"kind"`
	Target string  `json:"target" yaml:"target"`
	Amount float64 `json:"amount" yaml:"amount"`
}

// EnchantGroups maps a scroll's slot group to the template equip slots it
// covers; "any" fits every piece of gear.
var EnchantGroups = map[string][]string{
	"weapon":  {"mainhand", "offhand"},
	"armor":   {"helm", "chest", "legs", "boots", "gloves"},
	"jewelry": {"necklace", "ring", "trinket"},
	"any":     nil,
}

// ProfessionTargets are the enchant targets outside the character sheet:
// gathering yield per profession and crafting speed (DD-293).
var ProfessionTargets = map[string]bool{
	"gather_yield_mining": true, "gather_yield_herbalism": true, "gather_yield_woodcutting": true,
	"gather_yield_fishing": true, "gather_yield_hunting": true, "craft_speed": true,
}

// Fits reports whether the enchant can go on gear of a template equip slot.
func (e Enchant) Fits(equipSlot string) bool {
	if equipSlot == "" {
		return false
	}
	for _, g := range e.Slots {
		if g == "any" || slices.Contains(EnchantGroups[g], equipSlot) {
			return true
		}
	}
	return false
}
