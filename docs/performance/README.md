# Performance + Load Test — Katalog

> **Faz 35**: Server kapasite, throughput ve latency hedeflerinin ölçülmesi + doğrulanması. CCU tier bazlı senaryolar.

## 1. Kapsam
AFK-first MMO'nun **online + AFK karışık** yükünü, gerçekçi kullanıcı davranışı ile modelleyip her CCU seviyesinde SLO'ları doğrulama planı. Test kodu değil; **senaryo + hedef + tooling matrisi**.

## 2. Prensipler (DD-197)

- **SLO-driven**: her endpoint için latency + error budget tanımlı; test hedef bu SLO'ları doğrulamak.
- **Realistic mix**: bot script'i değil, gerçek oyuncu davranış oranı (online play + AFK reveal + market + chat).
- **Steady state + burst**: ısıtma → sabit yük → burst spike → recovery.
- **Long-running**: 1-2 saat testler drift yakalar (memory leak, connection pool eksilme, GC pauses).
- **Reproducible**: seed'li fake data + fixture; her koşu comparable.
- **Metric-first**: k6 client metrics + Prometheus server metrics + application traces; tek dashboard.
- **Baseline nightly**: main branch nightly test → regression < %5.

## 3. CCU Tier Hedefleri

| Tier | CCU | Concurrent Active | AFK Only | Notlar |
|---|---|---|---|---|
| **T1** | 1.000 | 300 | 700 | MVP dev env doğrulama |
| **T2** | 5.000 | 1.500 | 3.500 | Soft launch ready |
| **T3** | 10.000 | 3.000 | 7.000 | Launch capacity |
| **T4** | 50.000 | 15.000 | 35.000 | Peak season / event day |
| **T5** | 100.000 | 30.000 | 70.000 | Phase 03 stretch goal |

**Concurrent Active**: WS bağlı + actively playing (chat, combat, market).
**AFK only**: WS bağlı ama sadece contract ilerleyen (network overhead düşük).

Kullanıcı `active:afk = 3:7` oranı MMO idle genre benchmark; Phase 02 balance verisi ile tune.

## 4. SLO'lar (v1 taslak — [slo-slis.md](slo-slis.md))

| Endpoint / Sistem | Latency p95 | Latency p99 | Error rate |
|---|---|---|---|
| `POST /v1/auth/login` | 300ms | 800ms | < 0.5% |
| `GET /v1/character/{id}/sheet` | 150ms | 500ms | < 0.1% |
| `POST /v1/activity/{id}/reveal` | 250ms | 700ms | < 0.1% |
| `POST /v1/vendor/{code}/buy` | 200ms | 600ms | < 0.1% |
| `POST /v1/market/listing/{id}/buy` | 300ms | 900ms | < 0.5% |
| `GET /v1/market/search` | 200ms | 500ms | < 0.5% |
| `POST /v1/combat/{id}/skill` | 100ms | 300ms | < 0.5% |
| WebSocket ping p95 | 50ms | 200ms | — |
| WebSocket fanout p95 | 100ms | 400ms | — |

## 5. Dokümanlar
- [slo-slis.md](slo-slis.md) — Full SLI/SLO tablosu + error budget.
- [ccu-scenarios.md](ccu-scenarios.md) — CCU tier senaryo dağılımı.
- [load-login.md](load-login.md) — Login burst + credential fan-out.
- [load-combat.md](load-combat.md) — Live combat + AFK combat mix.
- [load-market.md](load-market.md) — Market search + list + buy karışım.
- [load-guild.md](load-guild.md) — Guild activity + chat + bank.
- [load-afk-calc.md](load-afk-calc.md) — Reveal throughput + batch worker.
- [load-dungeon.md](load-dungeon.md) — Dungeon start + AFK run + party.
- [load-websocket.md](load-websocket.md) — WS bağlantı + fanout + reconnect.
- [capacity-planning.md](capacity-planning.md) — Infrastructure sizing per tier.

## 6. Test Tooling

| Kullanım | Araç |
|---|---|
| Load generation | **k6** (JavaScript scripting) |
| WebSocket | k6 WS module + custom driver |
| Metrics collection | Prometheus + Grafana |
| Trace | OpenTelemetry + Jaeger/Tempo |
| Chaos | toxiproxy (network latency inject) |
| Profile | pprof (Go), async-profiler (Phase 02) |
| Report | Grafana dashboard export + Markdown summary |

k6 script'leri repo `perf/k6/` altında (`mix.js`, `login.js`); runner Docker — `scripts/loadtest.sh [smoke|t1|login]` (DD-238).

## 7. Test Ortamı

- **Staging**: production-like scale, tam veri seti, real Postgres/Redis.
- **Perf lab**: izole env (staging'i etkilemez), sadece perf test.
- **Prod canary**: yalnızca soak test (Phase 02); real traffic'in %5 shadow.

## 8. Cross-Reference
- [architecture/system-architecture.md](../architecture/system-architecture.md)
- [architecture/ci-cd.md](../architecture/ci-cd.md)
- [systems/activity-engine.md](../systems/activity-engine.md) — AFK reveal bottleneck kaynağı

## 9. Metrikler (Perf-Specific)
- `perf_test_run_total{scenario, tier, result}`
- `perf_test_p95_latency_ms{endpoint, scenario}`
- `perf_test_error_rate{endpoint, scenario}`
- `perf_test_ccu_gauge{tier}`
- `perf_test_regression_flag{scenario}` — nightly compare to baseline

## 10. Bekleyen Kararlar
- Perf lab bulut sağlayıcı (AWS/GCP) — infra karar.
- Prod canary %5 shadow onayı — Phase 02.
- Peak season capacity T4/T5 karar — launch data sonrası.
- Auto-scaling threshold Phase 02.
