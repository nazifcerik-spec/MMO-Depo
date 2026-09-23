# Redis (Client Wrapper)

## 1. Amaç
Redis bağlantı havuzu; convention-driven key namespace; Lua script yönetimi; health check. Ephemeral state kanalları için tek giriş.

Genel Redis kullanım politikası: [architecture/redis.md](../architecture/redis.md). Bu doküman Go client tarafını.

## 2. Sözleşme
```go
package cache

type Client struct {
    rdb *redis.Client
    log *slog.Logger
    scripts map[string]*redis.Script
}

type Config struct {
    URL             string
    PoolSize        int
    MinIdleConns    int
    DialTimeout     time.Duration
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    IdleTimeout     time.Duration
    CommandTimeout  time.Duration // default 100ms
}

func New(ctx context.Context, cfg Config, log *slog.Logger) (*Client, error)
func (c *Client) Ping(ctx context.Context) error
func (c *Client) Close() error

// Key builders (convention enforce)
func (c *Client) SessionBlockKey(jti string) string           // "sess:blk:<jti>"
func (c *Client) RateLimitKey(bucket, subject string) string  // "rl:<bucket>:<subject>"
func (c *Client) ActivityTickKey(contractID string) string    // "act:tick:<contract_id>"
func (c *Client) WSConnKey(node, charID string) string        // "ws:conn:<node>:<char>"
func (c *Client) PubChatChannel(kind, ref string) string      // "pub:chat:<kind>:<ref>"
func (c *Client) PubNotifChannel(charID string) string        // "pub:notif:<char_id>"
func (c *Client) LeaderboardKey(kind, id string) string       // "lb:<kind>:<id>"
func (c *Client) LockKey(resource string) string              // "lock:<resource>"
func (c *Client) MarketSearchCacheKey(hash string) string     // "mkt:search:<hash>"

// Rate limit (token bucket, Lua)
func (c *Client) AllowRate(ctx context.Context, bucket, subject string, limit, windowSec, cost int) (allowed bool, remaining int, resetAt time.Time, err error)

// Distributed lock (short-hold)
func (c *Client) Lock(ctx context.Context, resource, holder string, ttl time.Duration) (bool, error)
func (c *Client) Unlock(ctx context.Context, resource, holder string) error

// Pub/Sub
func (c *Client) Publish(ctx context.Context, channel string, payload []byte) error
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub
```

## 3. Key Convention (özet)
Bkz. [architecture/redis.md §3](../architecture/redis.md). Her key şu prensiplere uyar:
- `<domain>:<type>:<id>` snake_case.
- TTL zorunlu (opt-in exception listelenmiş).
- Value < 10KB.

## 4. Boot
```go
opt, err := redis.ParseURL(cfg.URL)
if err != nil { return nil, err }
opt.PoolSize = cfg.PoolSize
opt.MinIdleConns = cfg.MinIdleConns
opt.DialTimeout = cfg.DialTimeout
opt.ReadTimeout = cfg.ReadTimeout
opt.WriteTimeout = cfg.WriteTimeout
opt.IdleTimeout = cfg.IdleTimeout

rdb := redis.NewClient(opt)
if err := rdb.Ping(ctx).Err(); err != nil { return nil, err }

c := &Client{rdb: rdb, log: log, scripts: map[string]*redis.Script{}}
c.loadScripts()  // rate_limit.lua, ...
return c, nil
```

## 5. Lua Scripts

### 5.1 Token Bucket Rate Limit (`rate_limit.lua`)
```lua
-- KEYS[1] bucket key
-- ARGV[1] limit, ARGV[2] window (sec), ARGV[3] cost, ARGV[4] now_ms
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local cost = tonumber(ARGV[3])
local now_ms = tonumber(ARGV[4])

local bucket = redis.call("HMGET", key, "tokens", "ts")
local tokens = tonumber(bucket[1])
local ts = tonumber(bucket[2])

if tokens == nil then
  tokens = limit
  ts = now_ms
end

-- refill
local elapsed_ms = math.max(0, now_ms - ts)
local refill = elapsed_ms * limit / (window * 1000)
tokens = math.min(limit, tokens + refill)

local allowed = 0
if tokens >= cost then
  tokens = tokens - cost
  allowed = 1
end

redis.call("HMSET", key, "tokens", tokens, "ts", now_ms)
redis.call("PEXPIRE", key, window * 1000 * 2)

local reset_ms = math.ceil((limit - tokens) * (window * 1000) / limit)
return {allowed, math.floor(tokens), reset_ms}
```

Go call:
```go
res, err := c.scripts["rate_limit"].Run(ctx, c.rdb, []string{key}, limit, windowSec, cost, time.Now().UnixMilli()).Result()
```

### 5.2 Lock Release Guard
```lua
-- KEYS[1] lock key, ARGV[1] holder id
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end
```

## 6. Timeouts & Fail Policy
- `CommandTimeout` default 100 ms.
- Retry disabled by default (fail fast).
- Fail-open vs fail-closed matrix (bkz. architecture/redis.md §12).
- Redis unreachable → `errs.ErrUnavailable`; middleware karar verir.

## 7. Pub/Sub
- Subscribe pattern:
```go
sub := c.Subscribe(ctx, "pub:chat:guild:*", "pub:notif:*")
defer sub.Close()
for msg := range sub.Channel() {
    handle(msg.Channel, []byte(msg.Payload))
}
```
- WebSocket hub kullanır: chat message + notification fanout.

## 8. Leaderboard
```go
// increment contribution
c.rdb.ZIncrBy(ctx, key, delta, memberID)
// top N
c.rdb.ZRevRangeWithScores(ctx, key, 0, 19)
```

## 9. Health Check
- `Ping` her 30 sn.
- Metrics: `redis_pool_hits`, `redis_pool_misses`, `redis_cmd_duration_seconds{cmd}`.
- Slow command > 10 ms log.

## 10. Anti-Patterns
- `KEYS *` production'da.
- Blocking `BLPOP` uzun timeout ile.
- Uzun TTL'siz key birikimi.
- Serialize edilmiş büyük struct (> 10 KB) — Postgres'e gitmeli.
- Business logic Redis'te (source of truth ihlali).

## 11. Test Stratejisi
- Unit: `miniredis` in-memory fake.
- Integration: testcontainers-go Redis 7-alpine.
- Lua script: table-driven test (allow/deny/refill).

## 12. İlişkiler
- [Middleware](middleware.md), [WebSocket](websocket.md), [Event Bus](event-bus.md), [Worker](worker.md), [architecture/redis.md](../architecture/redis.md)

## 13. Bekleyen Kararlar
- Cluster mode geçiş eşiği.
- Streams-based event bus (Phase 02).
- Managed vs self-hosted final.
