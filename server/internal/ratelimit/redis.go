package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// tokenBucketScript is the atomic Redis Lua for one Consume operation.
// Fetches current (tokens, last_ms), refills based on elapsed time, and either
// decrements + returns allowed=1 or computes retry_after_ms.
//
// KEYS[1] = bucket key
// ARGV[1] = capacity          (int)
// ARGV[2] = refill_per_sec    (float)
// ARGV[3] = now_ms            (int)
// ARGV[4] = requested_tokens  (int)
// ARGV[5] = ttl_seconds       (int)
//
// Returns: {allowed_int, remaining_int, retry_after_ms_int}
const tokenBucketScript = `
local capacity = tonumber(ARGV[1])
local refill   = tonumber(ARGV[2])
local now      = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local ttl      = tonumber(ARGV[5])

local data = redis.call('HMGET', KEYS[1], 'tokens', 'last')
local tokens = tonumber(data[1])
local last   = tonumber(data[2])

if tokens == nil then
  tokens = capacity
  last   = now
end

local elapsed_sec = (now - last) / 1000.0
if elapsed_sec > 0 then
  tokens = math.min(capacity, tokens + elapsed_sec * refill)
end

local allowed = 0
local retry_after_ms = 0
if tokens >= requested then
  tokens = tokens - requested
  allowed = 1
elseif refill > 0 then
  local needed = requested - tokens
  retry_after_ms = math.ceil((needed / refill) * 1000)
end

redis.call('HMSET', KEYS[1], 'tokens', tokens, 'last', now)
redis.call('EXPIRE', KEYS[1], ttl)

return {allowed, math.floor(tokens), retry_after_ms}
`

// RedisStore runs the token-bucket algorithm as an atomic Redis Lua script.
// The script is loaded via EVALSHA on hot path with EVAL fallback (NOSCRIPT).
type RedisStore struct {
	client *redis.Client
	script *redis.Script
	clock  func() time.Time
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
		script: redis.NewScript(tokenBucketScript),
		clock:  time.Now,
	}
}

// keyTTL returns a safe TTL for the bucket key: long enough to survive silent
// idle periods but short enough that cold buckets don't accumulate. Rule of
// thumb: 10× the time it takes to refill from 0 to capacity, floored at 60s.
func keyTTL(b Bucket) time.Duration {
	if b.RefillPerSec <= 0 {
		return time.Hour
	}
	fillSec := float64(b.Capacity) / b.RefillPerSec
	ttl := time.Duration(fillSec*10) * time.Second
	if ttl < 60*time.Second {
		return 60 * time.Second
	}
	return ttl
}

func (r *RedisStore) Consume(ctx context.Context, key string, bucket Bucket, tokens int64) (Result, error) {
	nowMs := r.clock().UnixMilli()
	ttl := int64(keyTTL(bucket).Seconds())

	raw, err := r.script.Run(ctx, r.client, []string{key},
		bucket.Capacity,
		bucket.RefillPerSec,
		nowMs,
		tokens,
		ttl,
	).Result()
	if err != nil {
		return Result{}, fmt.Errorf("ratelimit script: %w", err)
	}

	arr, ok := raw.([]interface{})
	if !ok || len(arr) != 3 {
		return Result{}, fmt.Errorf("ratelimit unexpected reply: %T", raw)
	}
	allowed := toInt(arr[0]) == 1
	remaining := toInt(arr[1])
	retryMs := toInt(arr[2])

	return Result{
		Allowed:    allowed,
		Remaining:  remaining,
		RetryAfter: time.Duration(retryMs) * time.Millisecond,
	}, nil
}

// toInt coerces the interface{} Redis returns (usually int64) to int64.
func toInt(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case string:
		var i int64
		_, _ = fmt.Sscanf(n, "%d", &i)
		return i
	}
	return 0
}
