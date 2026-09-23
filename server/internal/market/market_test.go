package market

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListingFee_DocExampleAndMinimum(t *testing.T) {
	// docs/marketplace/tax.md §2: 10.000c × 1 × 5% × 1.5 = 750c.
	require.Equal(t, int64(750), ListingFee(10_000, 1, 24))
	require.Equal(t, int64(500), ListingFee(10_000, 1, 12))
	require.Equal(t, int64(2500), ListingFee(10_000, 1, 168))
	require.Equal(t, int64(ListingFeeMin), ListingFee(10, 1, 12), "floored at 10c")
	require.Equal(t, int64(1500), ListingFee(1_000, 20, 24), "scales with quantity")
}

func TestSaleCut(t *testing.T) {
	require.Equal(t, int64(500), SaleCut(10_000))
	require.Equal(t, int64(0), SaleCut(9))
}

func TestCancelRefund_Tiers(t *testing.T) {
	require.Equal(t, int64(375), CancelRefund(750, 30*time.Minute))
	require.Equal(t, int64(188), CancelRefund(750, 2*time.Hour))
	require.Equal(t, int64(0), CancelRefund(750, 6*time.Hour))
}

func TestPriceBounds(t *testing.T) {
	lo, hi := PriceBounds(100)
	require.Equal(t, int64(28), lo) // 100 × 0.25 × 1.10 = 27.5 → 28
	require.Equal(t, int64(1_000_000), hi)
	lo, _ = PriceBounds(0)
	require.Equal(t, int64(1), lo)
	_, hi = PriceBounds(1_000_000)
	require.Equal(t, int64(absoluteMaxPrice), hi)
}

func TestDurationMultiplier_OnlyDocumentedDurations(t *testing.T) {
	require.Len(t, DurationMultiplier, 4)
	_, ok := DurationMultiplier[6]
	require.False(t, ok)
}
