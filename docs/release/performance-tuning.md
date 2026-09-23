# Performance Tuning

## 1. Amaç
Beta wave süresince real user data ile perf regression tespit + fix. Load test'ler baseline; canlı player davranışı gerçek doğrulama.

## 2. Ana Ölçüm

| Kategori | Metrik |
|---|---|
| Latency | p50/p95/p99 per endpoint |
| Throughput | req/sec + WS msg/sec |
| Error rate | 5xx + auth fail + biz error |
| Resource | CPU/memory/disk/network per node |
| DB | query latency, cache hit, replication lag |
| Cache | Redis hit ratio, memory usage |
| GC | Go GC pause p99 |
| Cold start | Client cold start ms |

Grafana dashboard (bkz. [infra/monitoring.md](../infra/monitoring.md)).

## 3. Regression Detection

- **Baseline**: main branch nightly load test.
- **Comparison**: her deploy sonrası 24h ma vs baseline.
- **Threshold**: p95 latency +%10 → warn; +%25 → high; +%50 → critical (auto rollback trigger).
- **Endpoint-level**: her endpoint ayrı regression sinyal.

## 4. Common Bottleneck Kalıpları

### DB
- **N+1 query**: ORM lazy load → tek query'e batch (`IN` clause).
- **Missing index**: EXPLAIN plan seq scan → CREATE INDEX CONCURRENTLY.
- **Connection pool exhaustion**: pool size artır veya slow query fix.
- **Lock contention**: SERIALIZABLE hotspot → optimistic + retry.

### Redis
- **Hot key**: single key aşırı trafik → shard or split.
- **Big key**: >1MB value → paginate.
- **Pubsub subscriber sınırı**: partition subscribers.

### Application
- **GC pause**: allocation heavy → object pool.
- **Goroutine leak**: worker sızıntı → context cancel.
- **Sync.Mutex hot**: sharded map.
- **JSON parse cost**: pooled decoder + tight struct.

### WebSocket
- **Backpressure drop**: slow client → drop older or disconnect.
- **Fanout duplication**: subscriber map inefficient → batching.

### Client
- **Cold start**: bundle size → tree-shake + lazy load.
- **Re-render storm**: Riverpod overuse → memoization.
- **Image cost**: sprite sheet + WebP + CDN.

## 5. Profiling

Go:
- `pprof` HTTP endpoint (internal only) — CPU + heap + goroutine snapshot.
- `go tool pprof` analysis.
- **Continuous profiling** (Grafana Pyroscope Phase 02) — always-on low overhead.

Flutter:
- DevTools timeline + memory tab.
- Impeller renderer perf test.

## 6. Load Test Regression

Perf lab nightly:
- Run T3 baseline scenario.
- Compare to previous week baseline.
- Report Slack #perf.
- Regression > %10 → issue + investigate.

Bkz. [performance/](../performance/).

## 7. Real User Monitoring (RUM)

Client-side metric report to server:
- Cold start time.
- Screen transition time.
- API call latency (client-side view).
- Crash rate.

Aggregate Grafana; geo split (per region latency insight).

## 8. Optimization Prioritization

Karar:
- Top 5 slow endpoint (p95 × request count) → highest impact.
- Cost/benefit: 10% speedup önemsiz endpoint vs 3% önemli endpoint → önemli.
- User-facing perceived > backend efficient.

## 9. Post-Change Verification

Deploy sonrası:
- 24h observation window.
- Compare pre/post metric.
- Confirm target improvement + no regression elsewhere.

## 10. Documentation

Her perf work:
- Issue tracking (before/after metric).
- Techniques used (index add, query rewrite, cache add).
- Learnings post-mortem (best practices document).

## 11. Metrikler
- `perf_baseline_p95_ms{endpoint, week}`
- `perf_regression_detected_total{severity}`
- `perf_optimization_deployed_total`
- `perf_pprof_snapshot_total`

## 12. Cross-Reference
- [performance/](../performance/), [infra/monitoring.md](../infra/monitoring.md)
- [testing/README.md](../testing/README.md) — property-based perf

## 13. Bekleyen Kararlar
- Continuous profiling (Pyroscope) Phase 02.
- RUM SDK client-side (v1: manual metric) Phase 02.
- Load test regression auto issue create — Phase 01 CI integration.
- Client bundle size budget Phase 02.
