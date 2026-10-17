// Package market implements the fixed-price marketplace (Sprint 22,
// docs/marketplace/*): post with listing fee, buy with 5% sale cut, tiered
// cancel refund, expiry, search and price history.
package market

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

const (
	ListingCap        = 50
	ListingFeeRate    = 0.05
	ListingFeeMin     = 10
	SaleCutRate       = 0.05
	PendingSaleTTL    = 5 * time.Minute // worker resolves older pending sales
	minPriceMarkup    = 1.10            // docs/marketplace/price.md
	maxPriceFactor    = 10_000
	vendorSellMult    = 0.25 // DD-139 default sell multiplier
	vendorBuyMult     = 1.0
	absoluteMaxPrice  = 1_000_000_000
	defaultSearchSize = 25
	maxSearchSize     = 50
)

// DurationMultiplier maps allowed listing durations (hours) to fee multipliers.
var DurationMultiplier = map[int16]float64{12: 1.0, 24: 1.5, 48: 2.5, 168: 5.0}

func round(f float64) int64 { return int64(math.Round(f)) }

// ListingFee = max(10c, round(unit × qty × 5% × duration multiplier)).
func ListingFee(unit int64, qty int32, hours int16) int64 {
	fee := round(float64(unit) * float64(qty) * ListingFeeRate * DurationMultiplier[hours])
	if fee < ListingFeeMin {
		return ListingFeeMin
	}
	return fee
}

// SaleCut is the 5% taken from the buyer's payment.
func SaleCut(total int64) int64 { return round(float64(total) * SaleCutRate) }

// CancelRefund returns the part of the listing fee refunded on cancel:
// 50% under 1h, 25% under 6h, nothing after.
func CancelRefund(fee int64, age time.Duration) int64 {
	switch {
	case age < time.Hour:
		return round(float64(fee) * 0.50)
	case age < 6*time.Hour:
		return round(float64(fee) * 0.25)
	}
	return 0
}

// PriceBounds is the allowed unit price range for an item of base value.
func PriceBounds(baseValue int64) (min, max int64) {
	min = round(float64(baseValue) * vendorSellMult * minPriceMarkup)
	if min < 1 {
		min = 1
	}
	max = round(float64(baseValue) * vendorBuyMult * maxPriceFactor)
	if max < 1_000_000 {
		max = 1_000_000
	}
	if max > absoluteMaxPrice {
		max = absoluteMaxPrice
	}
	return min, max
}

type Listing struct {
	ID               uuid.UUID  `json:"id"`
	SellerID         uuid.UUID  `json:"seller_character_id"`
	SellerName       string     `json:"seller_name"`
	ItemInstanceID   *uuid.UUID `json:"item_instance_id"` // nil once the item merged into another stack
	TemplateCode     string     `json:"template_code"`
	ItemClass        string     `json:"item_class"`
	Rarity           string     `json:"rarity"`
	PrefixCode       *string    `json:"prefix_code,omitempty"`
	SuffixCode       *string    `json:"suffix_code,omitempty"`
	StackQuantity    int32      `json:"stack_quantity"`
	UnitPriceCopper  int64      `json:"unit_price_copper"`
	TotalPriceCopper int64      `json:"total_price_copper"`
	ListingFeeCopper int64      `json:"listing_fee_copper"`
	DurationHours    int16      `json:"duration_hours"`
	State            string     `json:"state"`
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        time.Time  `json:"expires_at"`
	ClosedAt         *time.Time `json:"closed_at,omitempty"`
}

type PriceStats struct {
	TemplateCode string `json:"template_code"`
	Sales        int64  `json:"sales_30d"`
	Quantity     int64  `json:"quantity_30d"`
	AvgUnit      int64  `json:"avg_unit_price_copper"`
	MinUnit      int64  `json:"min_unit_price_copper"`
	MaxUnit      int64  `json:"max_unit_price_copper"`
	LowestActive *int64 `json:"lowest_active_unit_price_copper,omitempty"`
}

var (
	ErrNotFound           = errors.New("market: listing not found")
	ErrUnavailable        = errors.New("market: listing is no longer available")
	ErrItemNotFound       = errors.New("market: item not found or not yours")
	ErrNotTradable        = errors.New("market: item is bound, locked, equipped or no-trade")
	ErrBadQuantity        = errors.New("market: invalid stack quantity")
	ErrBadDuration        = errors.New("market: duration must be 12, 24, 48 or 168 hours")
	ErrPriceOutOfRange    = errors.New("market: unit price out of allowed range")
	ErrPriceMismatch      = errors.New("market: price changed")
	ErrListingCap         = errors.New("market: 50 active listings max")
	ErrSelfBuy            = errors.New("market: cannot buy your own listing")
	ErrBagFull            = errors.New("market: bag full")
	ErrNotSeller          = errors.New("market: not your listing")
	ErrIdempotencyMissing = errors.New("market: Idempotency-Key header required")
)
