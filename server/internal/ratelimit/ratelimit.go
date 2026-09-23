// Package ratelimit implements DD-187: token-bucket rate limiter with three
// scope dimensions (IP, account, character). Any dimension exhaustion → 429.
//
// The Store interface separates policy from persistence so unit tests use
// MemoryStore and production wires RedisStore (Lua atomic).
package ratelimit

import (
	"context"
	"math"
	"sync"
	"time"
)

// Bucket is the token-bucket configuration for one scope+endpoint pair.
type Bucket struct {
	// Capacity is the burst size — max tokens the bucket ever holds.
	Capacity int64

	// RefillPerSec is the steady-state rate at which tokens are added.
	// Example: 20 requests/minute → 20.0/60 = 0.333 tokens/sec.
	RefillPerSec float64
}

// PerMinute returns a Bucket where capacity == n and refill == n/60 sec.
// This is the ergonomic constructor for "X requests per minute" limits.
func PerMinute(n int64) Bucket {
	return Bucket{Capacity: n, RefillPerSec: float64(n) / 60.0}
}

// PerHour returns a Bucket for "X requests per hour" (account recovery, DD-219).
func PerHour(n int64) Bucket {
	return Bucket{Capacity: n, RefillPerSec: float64(n) / 3600.0}
}

// PerSecond returns a Bucket for burstable X requests/second.
func PerSecond(n int64) Bucket {
	return Bucket{Capacity: n, RefillPerSec: float64(n)}
}

// Result is the outcome of one Consume attempt.
type Result struct {
	Allowed    bool
	Remaining  int64         // tokens left in the bucket after consumption attempt
	RetryAfter time.Duration // time until enough tokens exist for the requested amount
}

// Store consumes tokens from a bucket identified by key. Implementations must
// be safe for concurrent use and MUST be atomic per key (Redis Lua, or Go mutex
// per key). The bucket is created on first access with full capacity.
type Store interface {
	Consume(ctx context.Context, key string, bucket Bucket, tokens int64) (Result, error)
}

// --- MemoryStore (tests) -----------------------------------------------------

type memBucket struct {
	tokens float64
	last   time.Time
}

// MemoryStore is a goroutine-safe in-memory Store. Test-only.
type MemoryStore struct {
	mu    sync.Mutex
	rows  map[string]*memBucket
	clock func() time.Time
}

// NewMemoryStore returns an empty store using real time.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rows:  make(map[string]*memBucket),
		clock: time.Now,
	}
}

// NewMemoryStoreWithClock lets tests advance time deterministically.
func NewMemoryStoreWithClock(clock func() time.Time) *MemoryStore {
	return &MemoryStore{
		rows:  make(map[string]*memBucket),
		clock: clock,
	}
}

func (m *MemoryStore) Consume(_ context.Context, key string, bucket Bucket, tokens int64) (Result, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.clock()
	b, ok := m.rows[key]
	if !ok {
		b = &memBucket{tokens: float64(bucket.Capacity), last: now}
		m.rows[key] = b
	}

	// Refill based on elapsed time (clamped at capacity).
	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens = math.Min(float64(bucket.Capacity), b.tokens+elapsed*bucket.RefillPerSec)
		b.last = now
	}

	req := float64(tokens)
	if b.tokens >= req {
		b.tokens -= req
		return Result{
			Allowed:   true,
			Remaining: int64(math.Floor(b.tokens)),
		}, nil
	}

	// Not enough tokens — compute delay until refill catches up.
	needed := req - b.tokens
	var retryAfter time.Duration
	if bucket.RefillPerSec > 0 {
		retryAfter = time.Duration(math.Ceil(needed/bucket.RefillPerSec*1000)) * time.Millisecond
	}
	return Result{
		Allowed:    false,
		Remaining:  int64(math.Floor(b.tokens)),
		RetryAfter: retryAfter,
	}, nil
}
