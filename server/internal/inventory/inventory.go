// Package inventory manages character bag: slot assignment, weight budget,
// stack merge/split, move, delete.
//
// Layering:
//
//	inventory.go     — types, constants, errors.
//	repository.go    — Repository interface + MemoryRepo.
//	pg_repository.go — Postgres implementation.
//	service.go       — Pickup (auto-slot+merge) / Move / Split / Delete / View.
//	http.go          — HTTP handlers.
//
// Locked items (market/trade/mail/craft escrow) still occupy their slot but
// mutation is rejected. Unlock happens when the escrow resolves.
package inventory

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// v1 caps. Phase 02 makes these config-driven per character (bag expansion,
// STR-based weight, etc.).
const (
	BagSlotCap     = 100    // DD-246
	WeightCapGrams = 25_000 // 25 kg baseline (soft cap in v1 — no reject yet)
)

// Slot is a zero-based index into the bag.
type Slot = int16

// BagItem is one instance placed in the character's bag OR equipment. Mirrors
// inventory.item_instance for the slot query path — richer than item.Instance
// because it also carries the template metadata clients want for rendering
// (display name, weight, value).
//
// Bag items have Slot ≥ 0 + EquipmentSlot nil; equipped items have Slot = -1
// (irrelevant) + EquipmentSlot set.
type BagItem struct {
	InstanceID       uuid.UUID       `json:"instance_id"`
	TemplateCode     string          `json:"template_code"`
	DisplayNameKey   string          `json:"display_name_key"`
	ItemClass        string          `json:"item_class"`
	ItemSubclass     string          `json:"item_subclass"`
	Slot             Slot            `json:"slot"`
	EquipmentSlot    *string         `json:"equipment_slot,omitempty"` // 'mainhand'|'helm'|... nil if bagged
	StackQuantity    int32           `json:"stack_quantity"`
	MaxStack         int32           `json:"max_stack"`
	Rarity           string          `json:"rarity"`
	Prefix           *string         `json:"prefix,omitempty"`
	Suffix           *string         `json:"suffix,omitempty"`
	IsBound          bool            `json:"is_bound"`
	IsLocked         bool            `json:"is_locked"`
	LockOwner        *string         `json:"lock_owner,omitempty"`
	WeightGramsPer   int32           `json:"weight_grams_per"`
	BaseValueCopper  int64           `json:"base_value_copper"`
	StatBonusesRaw   json.RawMessage `json:"stat_bonuses,omitempty"`      // JSONB — decoded by skill.DecodeModifiers on demand
	EquipSlotOnTmpl  *string         `json:"equip_slot,omitempty"`        // template.equip_slot (for equip validation, character sheet compute)
	ClassRestriction *string         `json:"class_restriction,omitempty"` // template.class_restriction
	MinLevel         int16           `json:"min_level"`                   // template.min_level
	UseEffect        json.RawMessage `json:"use_effect,omitempty"`        // template.use_effect (DD-218); nil = not usable
}

// TotalWeight returns the item's total contribution to bag weight.
func (b BagItem) TotalWeight() int64 {
	return int64(b.WeightGramsPer) * int64(b.StackQuantity)
}

// View is the compound response for GET /v1/character/{id}/inventory.
type View struct {
	CharacterID    uuid.UUID `json:"character_id"`
	Items          []BagItem `json:"items"`
	SlotCap        int       `json:"slot_cap"`
	SlotsUsed      int       `json:"slots_used"`
	WeightGrams    int64     `json:"weight_grams"`
	WeightCapGrams int64     `json:"weight_cap_grams"`
}

// --- Errors ------------------------------------------------------------------

var (
	// ErrNotFound is returned when an item id doesn't belong to the caller or
	// doesn't exist. Existence-leak-safe.
	ErrNotFound = errors.New("inventory: not found")

	// ErrLocked is returned when a mutation targets a locked item (market /
	// trade / mail / craft / admin escrow).
	ErrLocked = errors.New("inventory: item locked by " + "")

	// ErrSpecRequirement / ErrAttributeRequirement: the item needs another
	// class branch or more raw attributes (DD-250 §34).
	ErrSpecRequirement      = errors.New("inventory: class branch requirement not met")
	ErrAttributeRequirement = errors.New("inventory: attribute requirement not met")

	// ErrBagFull is returned when the bag has no free slot for a new item.
	ErrBagFull = errors.New("inventory: bag full")

	// ErrSlotOccupied is returned by Move when target slot has an item that
	// cannot be swapped (locked).
	ErrSlotOccupied = errors.New("inventory: target slot occupied by a locked item")

	// ErrSlotOutOfRange is returned when a caller passes a slot beyond BagSlotCap.
	ErrSlotOutOfRange = errors.New("inventory: slot out of range")

	// ErrNotStackable is returned when Split targets a non-stackable item.
	ErrNotStackable = errors.New("inventory: item is not stackable")

	// ErrInvalidSplit is returned when split amount is <= 0 or >= stack.
	ErrInvalidSplit = errors.New("inventory: invalid split amount")

	// ErrInsufficient is returned by ConsumeByTemplate when the character
	// does not have enough unlocked units of the requested template across
	// their bag. No partial consumption happens — the whole operation is
	// atomic-or-nothing.
	ErrInsufficient = errors.New("inventory: insufficient quantity")
)
