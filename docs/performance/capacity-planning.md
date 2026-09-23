# Capacity Planning

## 1. Amaç
Her CCU tier için infrastructure sizing rehberi: uygulama node, DB, Redis, worker, LB, storage, network.

## 2. Kaynaklar (v1 tek-region)

| Kaynak | T1 (1k CCU) | T2 (5k) | T3 (10k) | T4 (50k) |
|---|---|---|---|---|
| **App node (Go)** | 1× 2vCPU/4GB | 2× 4vCPU/8GB | 3× 8vCPU/16GB | 8× 8vCPU/16GB |
| **Postgres primary** | 1× 2vCPU/8GB | 1× 4vCPU/16GB | 1× 8vCPU/32GB | 1× 16vCPU/64GB + replica |
| **Postgres replica (read)** | — | 1× 4vCPU/16GB | 1× 8vCPU/32GB | 2× 16vCPU/64GB |
| **Redis** | 1× 1vCPU/2GB | 1× 2vCPU/4GB | 1× 4vCPU/8GB (or cluster) | 3-node cluster 4vCPU/16GB |
| **Worker (Go)** | 1× 2vCPU/4GB | 1× 4vCPU/8GB | 2× 4vCPU/8GB | 4× 8vCPU/16GB |
| **LB (nginx/HAProxy)** | 1× 2vCPU/2GB | 1× 2vCPU/4GB | 2× 4vCPU/8GB | 2× 8vCPU/16GB |
| **Object storage (S3)** | 10 GB | 100 GB | 500 GB | 2 TB |
| **Bandwidth egress** | 100 GB/mo | 500 GB/mo | 1 TB/mo | 5 TB/mo |

## 3. Postgres Sizing

### T3 Örnek (10k CCU)
- Tables active: characters (~50k), item_instance (~2M), ledger_entry (~20M/year), audit rows.
- Storage: ~100 GB (12 ay retention).
- Write throughput: 200-500 write/sn peak.
- Read throughput: 800-1500 read/sn peak (read replica).
- Connection pool: 200 (app) + 50 (worker) + 50 (admin) = 300 total.
- pg_stat_statements + auto_explain nightly review.

### T4 Ölçek
- Storage: 500 GB (multi-year retention).
- Read replicas: 2× to distribute search + admin queries.
- Vertical scale primary; sharding karar Phase 03.

## 4. Redis Sizing

### T3
- Data types:
  - Sessions: 10.000 × 1 KB = 10 MB.
  - Presence: 10.000 × 200 B = 2 MB.
  - Rate limit buckets: 10.000 × 200 B = 2 MB.
  - Cache (market search): 10.000 key × 5 KB avg = 50 MB.
  - Combat sessions: 3.000 × 2 KB = 6 MB.
  - Chat rolling: kanal × 200 msg × 500 B = ~200 MB.
- Total: ~500 MB - 1 GB working set.
- Pubsub: 300-500 msg/sn.
- Memory 8 GB güvenli buffer.

### T4
- Cluster 3-node; sharding by key hash.
- Total ~4-6 GB working set.
- Pubsub ~2000 msg/sn.

## 5. Worker Sizing

Worker'lar:
- `activity_reveal_dispatcher` (event-driven, sync).
- `activity_expire_batch` (30s cron).
- `notif_push_dispatcher` (FCM/APNs external I/O).
- `economy_reconciler` (nightly).
- `security_anomaly_watcher` (1dk).
- `mail_expire`, `market_listing_expire`, `dungeon_lockout_reset`.
- `outbox_publisher` (event bus).

Concurrent worker goroutine cap: app node başına 200-500 (I/O bound).

## 6. Network

Bandwidth per CCU:
- Active player: ~1-5 KB/sn (WS ping + occasional message).
- AFK player: ~0.2-0.5 KB/sn.
- Peak per player: ~10 KB/sn (combat burst).

T3 aggregate: 10.000 × ~1.5 KB/sn = **~15 MB/sn** peak = ~130 Mbps.
T4: ~750 Mbps peak; 1 Gbps NIC yeterli.

## 7. Storage Growth

- ledger_entry: 500 entry/sn peak × 200 B = 100 KB/sn = 8.6 GB/gün T3.
- audit tables: aynı ölçek.
- Analytics combat_event: 200 event/sn × 300 B = 60 KB/sn = 5 GB/gün.
- Total DB growth ~20-30 GB/gün T3.
- Cold arşiv Parquet S3: aylık ~1 TB (12 ay compressed).

## 8. Multi-Region (Phase 03)

- Region başına full stack; DB primary tek region + async replica (write forward).
- WS shard by region; cross-region chat delay tolerated.
- Karar Phase 03: peak season data ile.

## 9. Cost Projection (Approximate, AWS us-east-1)

| Tier | Aylık Cost Approx |
|---|---|
| T1 | ~$200 |
| T2 | ~$800 |
| T3 | ~$2.500 |
| T4 | ~$12.000 |
| T5 (Phase 03) | ~$40.000+ |

Exclude egress bandwidth (~$0.09/GB varies).

## 10. Auto-Scaling

- App node: CPU > 70% 5dk → scale up 1; < 30% 15dk → scale down 1.
- Worker: queue depth > 5000 → scale up 1.
- DB: vertical scale (manuel); horizontal Phase 03.
- Redis: manual cluster resize.

## 11. Failure Domain

- 1 app node down: LB reroute; capacity -1/N; alert warn.
- Primary DB down: promote replica (RTO 1-5 dk); RPO ~seconds (streaming replication).
- Redis down: sessions in-memory; users forced re-login; alert critical.

## 12. Cross-Reference
- [README.md](README.md), [ccu-scenarios.md](ccu-scenarios.md), [slo-slis.md](slo-slis.md)
- [architecture/system-architecture.md](../architecture/system-architecture.md), [architecture/docker.md](../architecture/docker.md)

## 13. Bekleyen Kararlar
- Cloud provider seçim (AWS/GCP/Hetzner) — ops karar.
- Kubernetes vs Nomad orchestrator (DD-OPEN-9) — v1 scale karar.
- Multi-region timeline Phase 03.
- Cold archive S3 vs Cloudflare R2 cost karar Phase 02.
