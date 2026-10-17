// Package item is the Item Engine — templates (immutable content) + instances
// (per-character mutable) + rarity + affix.
//
// Layering:
//
//	template.go     — v1 hardcoded catalog (Sprint 12 YAML loader replaces).
//	rarity.go       — rarity enum + value multiplier + affix rules.
//	affix.go        — prefix/suffix pools + deterministic roll.
//	loot.go         — per-kind loot pool + RollDrops.
//	repository.go   — Repository interface + MemoryRepo.
//	pg_repository.go — Postgres implementation.
//	service.go      — CreateInstance / RollDrops / GetInstance.
//
// v1 scope:
//
//	✓ Template + instance
//	✓ Rarity (common → mythic)
//	✓ Affix (prefix + suffix) — cosmetic only, stats Sprint 11
//	✓ Loot pool per activity kind
//	✗ Socket + gem (Phase 02)
//	✗ Enchantment (Phase 02)
//	✗ Set bonus (Phase 02)
//	✗ Durability (Phase 02)
package item

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// BindKind controls trade/mail/market permission.
type BindKind string

const (
	BindOnEquip   BindKind = "boe" // Bind on first equip
	BindOnPickup  BindKind = "bop" // Bound immediately at drop/craft
	BindOnAccount BindKind = "boa" // Account-wide, but transferable across own chars
)

// Template is the immutable definition. One row per catalog entry.
type Template struct {
	Code             string
	DisplayNameKey   string
	ItemClass        string // weapon|armor|consumable|material|misc
	ItemSubclass     string
	MinLevel         int16
	BaseValueCopper  int64
	BaseWeightG      int32
	Stackable        bool
	MaxStack         int32
	BindKind         BindKind
	NoTrade          bool
	NoSell           bool
	Affixable        bool
	EquipSlot        *string // 'mainhand'|'helm'|'chest'|'boots'|... nil for non-equippable
	ClassRestriction *string // 'warrior'|... nil for universal
	StatBonuses      []byte  // JSONB — parsed via skill.DecodeModifiers on demand
	UseEffect        *UseEffect
	Requirements     *Requirements // DD-250: class branch + raw attribute minimums

	// Gathering node (DD-244): profession, unlock level 1..100, tier 0..7.
	// Empty kind = not gathered.
	GatherKind  string
	GatherLevel int16
	GatherTier  int16
}

// UseEffect is what consuming one unit does (content use_effect, DD-218).
// Templates without one cannot be used.
type UseEffect struct {
	Kind   string `json:"kind"` // heal_hp
	Amount int    `json:"amount"`
}

// IsEquippable returns true when the template has an equip_slot.
func (t Template) IsEquippable() bool { return t.EquipSlot != nil && *t.EquipSlot != "" }

// Instance is a concrete owned item. Rolled from a Template + Rarity + optional
// affix at creation time. Fields are frozen; upgrades / re-rolls happen only
// via explicit services.
type Instance struct {
	ID               uuid.UUID
	TemplateCode     string
	OwnerCharacterID *uuid.UUID
	StackQuantity    int32
	Rarity           Rarity
	PrefixCode       *string
	SuffixCode       *string
	IsBound          bool
	IsLocked         bool
	LockOwner        *string
	LockedRefID      *uuid.UUID
	SourceKind       *string
	SourceRefID      *uuid.UUID
	EquipmentSlot    *string // Sprint 11: nil = in bag or detached; set = equipped
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Version          int32
}

// IsEquipped returns true when the item is currently in an equipment slot.
func (i Instance) IsEquipped() bool { return i.EquipmentSlot != nil }

// --- Errors ------------------------------------------------------------------

var (
	// ErrNotFound is returned when an instance or template cannot be resolved.
	ErrNotFound = errors.New("item: not found")

	// ErrUnknownTemplate is returned when a template code is not in the catalog.
	ErrUnknownTemplate = errors.New("item: unknown template code")

	// ErrValidation wraps input-shape rejection.
	ErrValidation = errors.New("item: validation")
)

// Requirements are an item's equip requirements beyond MinLevel and
// ClassRestriction (DD-250 §34). Attribute minimums use RAW attributes.
type Requirements struct {
	Spec      string `json:"spec,omitempty"` // a class node the wearer's path must include
	Strength  int    `json:"strength,omitempty"`
	Vitality  int    `json:"vitality,omitempty"`
	Dexterity int    `json:"dexterity,omitempty"`
}
