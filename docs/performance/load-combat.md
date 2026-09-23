# Load Test: Combat

## 1. Amaç
Live combat skill cast throughput + WS event fanout + damage pipeline latency + concurrent session cap.

## 2. Senaryolar

### 2.1 Steady Combat Session
- Yük: T3'te 3000 active × 40 skill/saat = **33 skill/sn** system-wide.
- Süre: 60 dk.
- Beklenti: skill cast p95 < 100ms.

### 2.2 Combat Peak (Boss Fight)
- Yük: %20 active playerı aynı boss encounter'ında (rare, world boss).
- Skill cast rate: 600 active × 60 skill/dk = 600/sn boss peak.
- Süre: 3 dk (boss fight duration).
- Beklenti: p95 < 200ms; WS event fanout p95 < 400ms.

### 2.3 Concurrent Combat Session
- Yük: 3000 concurrent combat session (T3).
- Beklenti: session state Redis lookup < 10ms; damage pipeline p95 < 100ms.

### 2.4 Party Combat Coordination
- Yük: 1000 party (5 kişi) combat aynı anda; WS event party-scoped fanout.
- Beklenti: fanout latency p95 < 200ms per party.

## 3. Kritik Path

Skill cast flow:
1. HTTP receive.
2. Auth + rate limit.
3. Session state fetch (Redis).
4. Cooldown check (Redis).
5. Skill definition lookup (cache).
6. Damage pipeline 15-step.
7. Target HP update (Redis + DB async).
8. Threat table update (Redis sorted set).
9. Ledger (XP source, if kill) — opsiyonel.
10. WS event fanout to party (Redis pubsub).
11. Response.

**Bottleneck adayları:**
- **Redis lookup fan-out** (skill cast → 5-10 Redis op).
- **WS pubsub subscriber count** (party of 5 → 5 subs).
- **DB write eventual** (HP update → async → drift risk).

## 4. Metrics

- `combat_skill_cast_p95_ms`.
- `combat_pipeline_step_duration_ms{step}`.
- `combat_ws_event_fanout_latency_ms`.
- `redis_cluster_get_latency_ms`.
- `combat_session_active_gauge`.

## 5. Damage Pipeline Perf

15-step pipeline in-memory + Redis:
- Steps 1-10 (base/weapon/skill/crit/armor/resist) — pure CPU < 1ms.
- Steps 11-13 (threat, HP update, event emit) — Redis 3 op ~5ms.
- Step 14-15 (ledger if kill, loot roll if kill) — DB write async 20-50ms.

Kritik: HP update Redis primary; DB async writer batch (100ms window) — reveal-safe.

## 6. WS Fanout Model

Party of 5 → skill cast → 5 WS message.
Guild-scope events (5 kişi party × 30 kişi guild eğer olay guild) → daha büyük fanout.

Fanout mimarisi:
- Redis pubsub → application worker → active WS conn map lookup → serialize + send.
- Backpressure: WS conn buffer full → drop old events or disconnect (last resort).

## 7. Combat Session Cache

Redis: `combat:session:{sid}` HSET (participants, current_target, hp, threat_top).
- TTL 30 dk (inactive session cleanup).
- Serialize compact (msgpack veya JSON).

Session count 3000 × ~2 KB = 6 MB Redis (küçük).

## 8. Cooldown Check

Redis EXPIRE-based:
- `cd:{character_id}:{skill_code}` → PX (millisecond TTL).
- EXISTS → true = on cooldown.
- SET NX PX → set + check atomic.

Per skill cast 1 Redis op — sub-ms latency.

## 9. AFK Combat (Hunting)

AFK combat reveal analytical (bkz. [testing/afk.md](../testing/afk.md)):
- Reveal batch worker (5 dk pencerede tüm expiring contract).
- Yük: T3'te 3000 active × 0.1 hunt/dk = 300/dk = 5/sn reveal.
- Kolay throughput.

## 10. Beklenen Sonuç

T3 boss peak 600 skill/sn:
- p50: ~40ms.
- p95: ~120ms.
- p99: ~300ms.
- WS fanout p95: ~200ms.
- Redis cluster CPU ~40%.

## 11. Cross-Reference
- [combat/damage.md](../combat/damage.md), [combat/threat.md](../combat/threat.md)
- [backend/websocket.md](../backend/websocket.md), [backend/redis.md](../backend/redis.md)

## 12. Bekleyen Kararlar
- HP update batch window (100ms v1) Phase 02 tune.
- WS conn buffer size Phase 02.
- Redis pubsub → NATS migration threshold Phase 03.
