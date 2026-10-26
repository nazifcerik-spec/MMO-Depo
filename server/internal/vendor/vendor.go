// Package vendor implements NPC merchant buy/sell + 2h buyback (DD-139).
//
// Layering:
//
//	vendor.go        — types, errors, price formula.
//	repository.go    — Repository interface + MemoryRepo (tests).
//	pg_repository.go — Postgres implementation.
//	service.go       — Buy / Sell / Buyback / Info + TxnService + item.Service integration.
//	http.go          — HTTP handlers.
//
// Flow:
//
//	Buy  → check bag room → mint item (source_ref = ledger txn id) → ledger:
//	         character.wallet -X   → system.vendor_sink +X    (DD-139 implicit sink via sell_multiplier gap)
//	Sell → snapshot to buyback (2h TTL) → delete item → ledger:
//	         system.reward_source -Y → character.wallet +Y
//	Buyback → within 2h → recreate item from snapshot → ledger:
//	         character.wallet -Y   → system.reward_source +Y  (refunds the sell credit)
package vendor

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// Template mirrors content.vendor_template.
type Template struct {
	Code           string
	DisplayNameKey string
	VendorType     string
	BuyMultiplier  float64
	SellMultiplier float64
	Sells          []SellsEntry
}

// SellsEntry is one item in vendor's inventory (unlimited stock v1).
type SellsEntry struct {
	TemplateCode        string `json:"template_code"`
	PriceOverrideCopper *int64 `json:"price_override_copper,omitempty"` // nil → base × buy_multiplier
}

// SellsIndex maps template code → entry for O(1) lookup.
func (t *Template) SellsIndex() map[string]SellsEntry {
	m := make(map[string]SellsEntry, len(t.Sells))
	for _, e := range t.Sells {
		m[e.TemplateCode] = e
	}
	return m
}

// BuybackEntry mirrors economy.vendor_buyback.
type BuybackEntry struct {
	ID              uuid.UUID
	CharacterID     uuid.UUID
	VendorCode      string
	ItemSnapshot    json.RawMessage
	SoldPriceCopper int64
	SoldAt          time.Time
	ExpiresAt       time.Time
	ClaimedAt       *time.Time
}

// ItemSnapshot is the persistent shape stored in vendor_buyback.item_snapshot.
type ItemSnapshot struct {
	TemplateCode  string  `json:"template_code"`
	Rarity        string  `json:"rarity"`
	Prefix        *string `json:"prefix,omitempty"`
	Suffix        *string `json:"suffix,omitempty"`
	StackQuantity int32   `json:"stack_quantity"`
	IsBound       bool    `json:"is_bound"`
	Enchant       *string `json:"enchant,omitempty"` // DD-293
}

// --- Price formula ----------------------------------------------------------

// BuyPrice returns what the player pays for one unit.
// `overridePrice` (from SellsEntry.PriceOverrideCopper) supersedes template
// base × vendor.buy_multiplier when set.
func BuyPrice(baseValue int64, buyMultiplier float64, overridePrice *int64) int64 {
	if overridePrice != nil {
		return *overridePrice
	}
	return int64(math.Round(float64(baseValue) * buyMultiplier))
}

// SellPrice returns what the player receives per unit sold. Rarity multiplier
// applies (DD-139 rarity-scaled sell value).
func SellPrice(baseValue int64, sellMultiplier float64, rarityValueMultiplier float64) int64 {
	return int64(math.Round(float64(baseValue) * sellMultiplier * rarityValueMultiplier))
}

// --- Errors ------------------------------------------------------------------

var (
	// ErrVendorNotFound — unknown vendor code.
	ErrVendorNotFound = errors.New("vendor: not found")

	// ErrItemNotSold — vendor doesn't sell this template.
	ErrItemNotSold = errors.New("vendor: does not sell this item")

	// ErrItemNoSell — the item's template is marked no_sell (quest / soulbound).
	ErrItemNoSell = errors.New("vendor: item cannot be sold to vendor")

	// ErrItemBound — no longer returned by Sell (DD-245: bound items sell
	// to vendors); kept for the HTTP code mapping of older clients.
	ErrItemBound = errors.New("vendor: bound items cannot be sold")

	// ErrItemLocked — item is in escrow (market/trade/mail/craft).
	ErrItemLocked = errors.New("vendor: item is locked")

	// ErrItemNotOwned — item not owned by the caller.
	ErrItemNotOwned = errors.New("vendor: item not owned by caller")

	// ErrBuybackExpired — buyback window elapsed.
	ErrBuybackExpired = errors.New("vendor: buyback expired")

	// ErrBuybackClaimed — buyback already used.
	ErrBuybackClaimed = errors.New("vendor: buyback already claimed")

	// ErrBuybackNotFound — id doesn't match.
	ErrBuybackNotFound = errors.New("vendor: buyback not found")

	// ErrValidation — request shape rejection.
	ErrValidation = errors.New("vendor: validation")
)

// ValidationError attaches a field name to ErrValidation.
type ValidationError struct {
	Field, Message string
}

func (e *ValidationError) Error() string { return "vendor: validation: " + e.Field + ": " + e.Message }
func (e *ValidationError) Unwrap() error { return ErrValidation }
