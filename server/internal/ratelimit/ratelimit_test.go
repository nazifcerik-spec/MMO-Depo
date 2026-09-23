package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/ratelimit"
)

func TestMemoryStore_AllowsInsideCapacity(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	bucket := ratelimit.PerMinute(5) // capacity 5, refill 5/60 per sec

	for i := 0; i < 5; i++ {
		r, err := store.Consume(context.Background(), "key1", bucket, 1)
		require.NoError(t, err)
		require.True(t, r.Allowed, "call %d must be allowed", i)
	}
}

func TestMemoryStore_RejectsWhenExhausted(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	bucket := ratelimit.PerMinute(3)

	// Drain
	for i := 0; i < 3; i++ {
		r, err := store.Consume(context.Background(), "k", bucket, 1)
		require.NoError(t, err)
		require.True(t, r.Allowed)
	}

	// Next → denied
	r, err := store.Consume(context.Background(), "k", bucket, 1)
	require.NoError(t, err)
	require.False(t, r.Allowed)
	require.Equal(t, int64(0), r.Remaining)
	require.Greater(t, r.RetryAfter, time.Duration(0))
}

func TestMemoryStore_KeysAreIndependent(t *testing.T) {
	store := ratelimit.NewMemoryStore()
	bucket := ratelimit.PerMinute(1)

	// Exhaust key A
	r1, _ := store.Consume(context.Background(), "A", bucket, 1)
	require.True(t, r1.Allowed)
	r2, _ := store.Consume(context.Background(), "A", bucket, 1)
	require.False(t, r2.Allowed)

	// Key B still fresh
	r3, _ := store.Consume(context.Background(), "B", bucket, 1)
	require.True(t, r3.Allowed)
}

func TestMemoryStore_RefillOverTime(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	clock := func() time.Time { return now }
	store := ratelimit.NewMemoryStoreWithClock(clock)

	bucket := ratelimit.Bucket{Capacity: 2, RefillPerSec: 1.0} // 1 token/sec

	// Drain
	r, _ := store.Consume(context.Background(), "k", bucket, 1)
	require.True(t, r.Allowed)
	r, _ = store.Consume(context.Background(), "k", bucket, 1)
	require.True(t, r.Allowed)

	// Empty
	r, _ = store.Consume(context.Background(), "k", bucket, 1)
	require.False(t, r.Allowed)

	// Advance 1.5s → 1 token refilled (floored)
	now = now.Add(1500 * time.Millisecond)
	r, _ = store.Consume(context.Background(), "k", bucket, 1)
	require.True(t, r.Allowed)

	// Empty again
	r, _ = store.Consume(context.Background(), "k", bucket, 1)
	require.False(t, r.Allowed)
}

func TestMemoryStore_RetryAfterEstimate(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	store := ratelimit.NewMemoryStoreWithClock(func() time.Time { return now })
	bucket := ratelimit.Bucket{Capacity: 1, RefillPerSec: 0.5} // 1 token every 2s

	r, _ := store.Consume(context.Background(), "k", bucket, 1)
	require.True(t, r.Allowed)

	r, _ = store.Consume(context.Background(), "k", bucket, 1)
	require.False(t, r.Allowed)
	// Need 1 token, refill 0.5/sec → 2s.
	require.InDelta(t, 2*time.Second, r.RetryAfter, float64(50*time.Millisecond))
}

func TestMemoryStore_CapacityClampsRefill(t *testing.T) {
	// Long idle should not overflow beyond capacity.
	now := time.Unix(1_700_000_000, 0)
	store := ratelimit.NewMemoryStoreWithClock(func() time.Time { return now })
	bucket := ratelimit.Bucket{Capacity: 3, RefillPerSec: 1.0}

	// Drain
	for i := 0; i < 3; i++ {
		_, _ = store.Consume(context.Background(), "k", bucket, 1)
	}

	// Advance 100 seconds — should refill only to capacity=3
	now = now.Add(100 * time.Second)

	// 3 allowed
	for i := 0; i < 3; i++ {
		r, _ := store.Consume(context.Background(), "k", bucket, 1)
		require.True(t, r.Allowed, "call %d", i)
	}
	// 4th denied
	r, _ := store.Consume(context.Background(), "k", bucket, 1)
	require.False(t, r.Allowed)
}

func TestPerMinuteConstructor(t *testing.T) {
	b := ratelimit.PerMinute(60)
	require.Equal(t, int64(60), b.Capacity)
	require.InDelta(t, 1.0, b.RefillPerSec, 0.0001)
}

func TestPerSecondConstructor(t *testing.T) {
	b := ratelimit.PerSecond(10)
	require.Equal(t, int64(10), b.Capacity)
	require.Equal(t, 10.0, b.RefillPerSec)
}
