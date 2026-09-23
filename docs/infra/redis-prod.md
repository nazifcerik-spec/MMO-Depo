# Redis Production

## 1. Amaç
Prod Redis deployment: HA (Sentinel veya Cluster), persistence, memory eviction, pubsub scale.

## 2. Sürüm

- **Redis 7.2+** (or Valkey 7.2 — open-source fork; Phase 02 evaluate).

## 3. Topology

### T3 Launch
- Sentinel HA: 1 primary + 2 replica + 3 Sentinel node.
- OR: 1 primary + 1 replica (yeterli tek AZ) + manual failover.
- Total memory: ~8 GB.

### T4 Peak
- **Cluster** mode 3-shard (primary + replica each) = 6 node.
- Hash slot dağılımı.
- Total memory: ~24 GB.

### v1 Managed
- AWS ElastiCache (Redis / Valkey).
- Multi-AZ replica; automatic failover.

## 4. Persistence

Karar tabanlı:
- **AOF everysec** (append-only file) — durability high, throughput OK.
- Snapshot (RDB) hourly + boot.
- Persistence YES çünkü: session, refresh token, idempotency, wallet cache — kritik.

Config:
```
appendonly yes
appendfsync everysec
save 3600 1 300 100 60 10000
```

## 5. Memory + Eviction

- `maxmemory 8gb` (or cluster shard-per size).
- `maxmemory-policy allkeys-lru` — kritik olmayan (cache) evict; kritik veri Postgres backup.

**Ancak kritik key namespace'lere eviction istenmez**:
- `session:*`, `idempotency:*`, `wallet:cache:*` — bunlar Postgres-authoritative; evict OK.
- `presence:*`, `chat:live:*` — TTL-managed; eviction risk düşük.

## 6. Pubsub Scale

v1 Redis pubsub tek node routing OK ~2k msg/sn.
T4'de pubsub cluster (Cluster mode pubsub cross-shard) — Redis 7 destekler.
T5 Phase 03 → NATS opsiyonel (DD-200).

## 7. Key Design

Naming convention: `<namespace>:<entity>:<id>[:<sub>]`
```
session:{sid}
presence:character:{cid}
idempotency:{endpoint}:{key}
rate:{bucket_key}:{actor_id}
combat:session:{sid}
chat:live:global:region:{region_id}
wallet:cache:character:{cid}
```

TTL zorunlu (memory güvenliği):
- session: 30d
- presence: 5m
- idempotency: 24h
- rate: 60s (bucket refill)
- combat.session: 30m
- chat.live: 24h
- wallet.cache: 5m

## 8. Connection

- Client: `go-redis/v9` — connection pool.
- Pool size: 50 per app node.
- Cluster mode: client-aware routing (hash slot).

## 9. Monitoring

redis_exporter → Prometheus:
- `redis_memory_used_bytes`, `redis_maxmemory_bytes`
- `redis_commands_processed_total`
- `redis_pubsub_channels`
- `redis_pubsub_subscribers`
- `redis_connected_clients`
- `redis_slowlog_length`
- Replication lag: `master_repl_offset - slave_repl_offset`.

Alerts:
- Memory > %90: warn.
- Replica lag > 5s: warn.
- Persistence disable / AOF rewrite fail: critical.

## 10. Backup

- AOF + RDB otomatik.
- Managed RDS/ElastiCache: automatic snapshot daily.
- Self-managed: cron snapshot → S3 offsite.
- RPO 1h (snapshot); RTO ~15 dk (restore).

## 11. Sentinel Setup (Self-managed)

Sentinel monitors primary:
- Quorum 2/3 for failover decision.
- Promote replica automatic.
- Client (go-redis Sentinel driver) discovers primary.

## 12. Cluster Mode Considerations

- Cross-shard commands limited (`MULTI/EXEC` same slot only).
- Lua scripts single-key (hash tag `{...}` grouping).
- Pubsub cluster-wide (Redis 7).

## 13. Security

- TLS zorunlu (`tls-port 6379`, no plain 6379).
- Auth: `requirepass <strong>`, `masterauth` for replica.
- ACL users (Redis 6+): app_write, app_read, monitoring.
- Network policy: only app subnet allowed.

## 14. Cost

Managed ElastiCache T3 (4vCPU / 8GB): ~$200/ay.
T4 cluster 3-shard: ~$800-1200/ay.

## 15. Cross-Reference
- [postgres-prod.md](postgres-prod.md), [backup.md](backup.md), [monitoring.md](monitoring.md)
- [architecture/redis.md](../architecture/redis.md) — v1 temel karar
- [performance/capacity-planning.md](../performance/capacity-planning.md)

## 16. Bekleyen Kararlar
- Managed (ElastiCache) vs self-managed (Sentinel) — ops karar.
- Valkey fork vs Redis (license concern Redis 7.4 SSPL) — Phase 01 karar.
- Cluster migration threshold T3 vs T4 — data.
- NATS pubsub migration Phase 03.
