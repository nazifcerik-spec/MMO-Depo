# Redis (Ephemeral State + Pub/Sub)

## 1. Amaç
Sunucu içi geçici state, canlı push, rate-limit, kısa süreli lock ve leaderboard. **Redis kaybı canlılığı etkiler; kalıcı veriyi etkilemez.**

## 2. Sürüm ve Mod
- **Redis 7.x** (stable).
- Mod:
  - Phase 01: single primary + 1 replica (yüksek erişilebilirlik).
  - Phase 02: Redis Sentinel veya managed alternatif; ölçek ihtiyacına göre Cluster.
- Persistence: RDB + AOF (everysec). Kayıp toleransı: < 1 sn.

## 3. Kullanım Alanları
| Alan | TTL | Key Pattern |
|---|---|---|
| Session token blocklist | refresh TTL | `sess:blk:<jti>` |
| Rate limit buckets | 60 sn window | `rl:<endpoint>:<subject>` |
| Active activity ticker | 2 dk | `act:tick:<contract_id>` |
| WebSocket connection registry | süre-boyu | `ws:conn:<node>:<char_id>` |
| Pub/Sub channels | – | `pub:chat:<channel>`, `pub:notif:<char_id>` |
| Live leaderboards (boss contribution) | span | `lb:boss:<spawn_id>` (sorted set) |
| Distributed lock (kısa) | 1–30 sn | `lock:<resource_id>` |
| Marketplace search cache | 30 sn | `mkt:search:<hash>` |
| Content data hash cache | patch süresi | `content:hash:<name>` |

## 4. Kurallar
- **Her key TTL almalı** (uzun-yaşayan istisnalar açıkça listelenir).
- Key naming: `<domain>:<type>:<id>`; snake_case, `:` separator.
- Value: küçük tut (< 10 KB); büyük payload → Postgres.
- **No source-of-truth data in Redis** — Redis crash olsa oyun temel state'i kaybolmaz.
- Concurrent kritik section için `SET NX PX` veya Redlock (bounded scope, short hold).

## 5. Pub/Sub
- Chat fan-out: her channel için Redis pub/sub topic; app nodes subscribe.
- Notification fan-out: per-character subscription; delivery WebSocket hub üzerinden.
- World event announcement: server-wide topic.
- Alternatif (v2 aday): Redis Streams — replay + consumer group ihtiyacı ortaya çıkarsa.

## 6. Rate Limiting
- Token bucket veya sliding window (Lua script atomic).
- Per endpoint + per account + per IP.
- Response header: `X-RateLimit-Remaining`, `X-RateLimit-Reset`.

## 7. Locking (Kısa Süreli)
- `SET key value NX PX 5000` pattern.
- Value = requester unique id (release safe).
- Uzun süreli lock yasak (> 30 sn); business logic çevir.
- Redlock yalnızca gerçek multi-node kritik section için.

## 8. Leaderboards (Live)
- Sorted set (`ZADD`, `ZINCRBY`).
- World boss contribution: `lb:boss:<spawn_id>` — key expiry span.
- Season leaderboards (Phase 02) sorted set + snapshot to Postgres.

## 9. Cache Invalidation
- Kural: **write-through**, not lazy. Postgres yazımı sonrası Redis update/delete.
- Content patch: version bump → Redis key prefix'i yeni versiyon → eski otomatik unused.

## 10. Client Konvansiyonları
- `go-redis` client; pool sizing config'te.
- Command timeout 100 ms; retry disabled by default (fail fast).
- Lua script'ler versiyonlu (`SCRIPT LOAD` startup'ta).

## 11. Observability
- Prometheus exporter (Redis).
- Slow log > 10 ms.
- Memory usage, evicted keys, cluster health.

## 12. Failure Modu
- Redis unreachable → app node fall back:
  - Rate limit: fail-open? fail-closed? — **kritik endpoint fail-closed** (login, market), diğerleri fail-open.
  - WebSocket: mesajlar in-memory buffer'a düşer; kısa outage tolere edilir.
  - Live activity tick cache: Postgres query'e düşer (biraz yavaş; kabul edilebilir).
- Restart: RDB + AOF replay.

## 13. Anti-Patterns
- Redis'i primary DB gibi kullanmak.
- Uzun TTL'siz key birikimi.
- Blocking komut (KEYS *, unbounded SCAN) production'da.
- Serialize/deserialize edilmiş büyük struct saklamak (> 10 KB).

## 14. İlişkiler
- [System Architecture](system-architecture.md), [Go Server](go.md), [Notification System](../systems/notification-system.md), [Marketplace](../game-design/marketplace.md), [World Bosses](../game-design/world-bosses.md)

## 15. Bekleyen Kararlar
- Managed (ElastiCache/Cloud Memory Store) mi self-hosted mı — v1 için ilk cloud provider seçimi bağlı.
- Cluster mode gerçek ihtiyaç eşiği (memory > 8 GB veya throughput > 50k ops/sec).
- Streams vs classic Pub/Sub — durability ihtiyacı ortaya çıkarsa migration path.
