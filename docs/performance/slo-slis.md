# SLI + SLO Tablosu

## 1. Amaç
Service Level Indicator (SLI) ve Objective (SLO) tanımları — ölçülen + hedef. Perf test bunları doğrular; ops bunları izler.

## 2. Prensipler

- **Kullanıcı-facing SLI** öncelikli (login, sheet, reveal, buy, chat send).
- **Percentile-based** (p50, p95, p99) — ortalama yanıltıcı.
- **Error budget** aylık: SLO'nun tersinden hesap; aşımda değişiklik dondurma.
- **Alert tiering**: budget %50 tüketildiğinde warn; %80 high; %100 critical.

## 3. Full SLI/SLO Tablosu

### Auth / Session
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `POST /v1/auth/login` | 120ms | 300ms | 800ms | < 0.5% |
| `POST /v1/auth/refresh` | 80ms | 200ms | 500ms | < 0.1% |
| `POST /v1/auth/otp/verify` | 150ms | 400ms | 900ms | < 0.5% |

### Character
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `GET /v1/character/{id}/sheet` | 60ms | 150ms | 500ms | < 0.1% |
| `GET /v1/character/{id}/inventory` | 80ms | 250ms | 700ms | < 0.1% |
| `POST /v1/character/{id}/equip` | 100ms | 300ms | 800ms | < 0.5% |
| `POST /v1/character/{id}/ascend` | 500ms | 1500ms | 3000ms | < 0.5% |

### Activity (AFK)
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `POST /v1/activity/start` | 80ms | 200ms | 500ms | < 0.1% |
| `POST /v1/activity/{id}/reveal` | 100ms | 250ms | 700ms | < 0.1% |
| `POST /v1/activity/{id}/stop` | 80ms | 200ms | 500ms | < 0.1% |

### Combat
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `POST /v1/combat/{id}/skill` | 40ms | 100ms | 300ms | < 0.5% |
| `POST /v1/combat/{id}/consumable` | 50ms | 150ms | 400ms | < 0.5% |
| WS combat_event latency | 30ms | 100ms | 300ms | — |

### Market
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `GET /v1/market/search` | 60ms | 200ms | 500ms | < 0.5% |
| `POST /v1/market/listing` | 100ms | 300ms | 800ms | < 0.5% |
| `POST /v1/market/listing/{id}/buy` | 120ms | 300ms | 900ms | < 0.5% |
| `DELETE /v1/market/listing/{id}` | 80ms | 250ms | 700ms | < 0.5% |

### Vendor / Fee
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `POST /v1/vendor/{code}/buy` | 80ms | 200ms | 600ms | < 0.1% |
| `POST /v1/vendor/{code}/sell` | 80ms | 200ms | 600ms | < 0.1% |

### Guild
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `GET /v1/guild/{id}` | 80ms | 250ms | 700ms | < 0.1% |
| `POST /v1/guild/{id}/bank/withdraw` | 120ms | 300ms | 800ms | < 0.5% |
| `PUT /v1/guild/{id}/motd` | 80ms | 250ms | 700ms | < 0.5% |

### Chat / Social
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `POST /v1/chat/send` | 50ms | 150ms | 400ms | < 0.5% |
| WS chat delivery | 30ms | 100ms | 300ms | — |
| `POST /v1/mail` | 100ms | 300ms | 800ms | < 0.5% |
| `GET /v1/mail/inbox` | 60ms | 200ms | 500ms | < 0.1% |

### Dungeon
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| `POST /v1/dungeon/{template}/enter` | 200ms | 600ms | 1200ms | < 0.5% |
| `POST /v1/dungeon/{run}/reveal` (AFK) | 200ms | 500ms | 1200ms | < 0.5% |

### WebSocket
| SLI | p50 | p95 | p99 | Error |
|---|---|---|---|---|
| WS connect handshake | 50ms | 150ms | 400ms | < 0.5% |
| WS ping (RTT) | 20ms | 50ms | 200ms | — |
| WS reconnect success | — | — | — | > 99% (5 dk grace) |

## 4. Availability SLO

| Servis | Uptime Hedef | Notlar |
|---|---|---|
| Public API | 99.9% (~43 dk/ay downtime budget) | Region single v1 |
| WebSocket | 99.5% | Reconnect grace tolerant |
| Admin panel | 99.0% | İç kullanım |
| Worker pipeline | 99.9% | Reveal/reconciler critical |

## 5. Error Budget & Alerting

- **Budget** = (1 - SLO_target) × time_window.
- Örnek: 99.9% uptime → 43 dk/ay budget.
- Consumption tracking:
  - < 50%: normal
  - 50-80%: warn (Slack)
  - 80-100%: high (Slack + inbox)
  - > 100%: critical (PagerDuty + release freeze)

Release freeze: yeni feature deploy dondurulur; sadece hotfix + reliability work.

## 6. Metrikler
- `slo_target_pct{endpoint, quantile}`
- `slo_actual_pct{endpoint, quantile}`
- `slo_error_budget_remaining_pct{endpoint}`
- `slo_burn_rate{endpoint}` (window bazlı)

## 7. Cross-Reference
- [README.md](README.md), [ccu-scenarios.md](ccu-scenarios.md)
- [security/audit.md](../security/audit.md) — alert tiering
- [architecture/system-architecture.md](../architecture/system-architecture.md)

## 8. Bekleyen Kararlar
- Uptime SLO %99.95 (launch data sonrası tune).
- Multi-region SLO (Phase 03).
- SLO'ların user-facing status page'e yansıtımı Phase 02.
