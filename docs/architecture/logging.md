# Logging (Loglama)

## 1. Amaç
Sistem davranışını, hataları ve iş olaylarını yapısal, aranabilir, korelasyona uygun şekilde kaydetmek. Debug + audit + oyuncu destek + güvenlik ihtiyaçları için tek disiplin.

## 2. Prensipler
- **Structured JSON logs** — her satır makine-okunur.
- **Correlation everywhere**: her log satırı `request_id`, mümkünse `character_id`.
- **Log ≠ audit ≠ metric**: her biri farklı kanal (bkz. §11).
- **Level bilinçli**: `debug` gürültü değildir; sadece detaylı; `warn` aksiyona değer; `error` sistem sorunu.
- **PII minimize**: kişisel bilgiler asla loga (email, IP → hash / partial).

## 3. Level Kılavuzu

| Level | Ne zaman |
|---|---|
| `debug` | Geliştirici detayı; prod'da default kapalı |
| `info` | Domain akış: request tamamlandı, activity başladı, market işlemi |
| `warn` | Beklenen değil ama tekrar edilebilir: retry yaptı, rate limit hit, degraded |
| `error` | Uygulama hatası; incident/alert değerlendirilir |
| `fatal` | Boot fail; process exit — çok nadir |

## 4. Log Formatı (JSON)
```json
{
  "ts": "2026-09-08T12:34:56.789Z",
  "level": "info",
  "msg": "activity started",
  "service": "server",
  "env": "prod",
  "version": "v1.0.3-abc1234",
  "request_id": "01J8T2P...",
  "trace_id": "9a8b7c...",
  "span_id": "0f1e2d...",
  "character_id": "c-777",
  "activity_type": "mining",
  "spot_id": "iron_valley",
  "duration_ms": 42
}
```

Alanlar:
- `ts` — UTC ISO 8601, ms hassasiyet.
- `level` — `debug|info|warn|error|fatal`.
- `msg` — kısa, sabit ("activity started"), interpolate ETMEZ (parametreler alan olarak).
- `service` — `server|worker|migrate|client-web` vb.
- `env`, `version` — deploy izleme.
- `request_id`, `trace_id`, `span_id` — korelasyon.
- Domain alanları — modüle göre standard set (character_id, activity_id, order_id, ...).

**Kural**: msg değişken içermez. `"user 123 not found"` yerine `msg="user not found", user_id=123`.

## 5. Correlation
- HTTP: `X-Request-ID` header (client verirse honor, yoksa gen).
- Context propagation: her handler `ctx` alır; her log context'ten alanları çeker.
- Cross-service (worker ↔ server): job payload `request_id`, `trace_id` taşır.
- OpenTelemetry ile trace/span id'ler otomatik enjekte.

## 6. Domain Log Standartları (örnekler)
- Activity: `activity_started`, `activity_paused`, `activity_completed`, `activity_reveal`.
- Combat: `damage_calc` (debug only, high volume), `kill`, `death`.
- Economy: `ledger_write`, `market_listing_created`, `market_transaction`, `mail_delivered`.
- Auth: `login_ok`, `login_fail`, `token_refresh`, `session_revoked`.
- Guild: `guild_created`, `guild_member_joined`, `guild_bank_write`.
- Worker: `job_started`, `job_completed`, `job_failed`, `job_retry`.

## 7. Ne LOGLANMAZ
- **Şifre, token, secret, PII (email, IP)**.
- Tam item stack listeleri (sadece adedi).
- Chat mesaj gövdesi (`chat_message` metadata evet, body sadece report akışında).
- Kredi kartı bilgisi — hiçbir yerde tutulmaz (payment provider tokenize eder).
- Büyük payload (> 4 KB) — özet alan.
- Loop içinde her tick log (bulk aggregate).

## 8. Client (Flutter) Logging
- Development: verbose console.
- Production:
  - Sadece error/warn kanal (Sentry-benzeri).
  - Correlation id server ile eşleşir.
  - PII flag: user opt-in olmadan diagnostic gönderme.
- Local persistent log: sadece bug report akışı için, opt-in.

## 9. Log Pipeline
```
App stdout (JSON) → collector (Fluent Bit / Vector) → aggregator (Loki / OpenSearch / cloud) → dashboards (Grafana)
```
- Retention:
  - `debug`: 3 gün.
  - `info`: 14 gün.
  - `warn`/`error`: 30 gün.
  - `audit`: 12 ay+ (bkz. §11).
- Retention'ları cost budget'a göre revize (Phase 01).

## 10. Alerting
- Metric-based alert öncelikli; log-based ikincil.
- `error` rate spike (baseline + %50 in 5 min) → PagerDuty.
- `fatal` her zaman ping.
- `warn` sinyal panosu; alarm değil.
- False-positive review haftalık.

## 11. Log vs Audit vs Metric

| | Amaç | Kanal | Retention | Sample |
|---|---|---|---|---|
| **Log** | Debug + trace | Structured JSON stdout | 3–30 gün | full |
| **Audit** | Compliance + forensic | Postgres `audit.event_log` | 12 ay+ | full |
| **Metric** | Trend + alert | Prometheus | 90 gün+ | aggregated |
| **Trace** | Request-level path | OTEL / Tempo | 7 gün | sampled (1–10%) |

Audit tabloları: login events, currency ledger, admin actions, permission changes.

## 12. Rate Limiting Logs
- Yüksek-frekans event: bulk metric + periyodik summary log (per-minute).
- Sample debug logs (< %1) uzun-süreli inceleme için.

## 13. Zaman
- **Her zaman UTC** (`ts` alanında `Z` suffix zorunlu).
- Client timestamp'i ayrı alan olarak (varsa), otorite değildir.

## 14. Anti-Patterns
- `fmt.Println` production'da.
- `msg`'de string interpolation (grep zorlaşır).
- Stack trace her info log'a.
- Secret'ı bir kez de olsa log'a bırakmak (immutable — replace bile silinmiş sayılmaz).
- Log seviyesini kod içinde runtime değiştirmek.

## 15. Test
- Log satırlarının şeması test edilir (kritik event'ler için).
- Missing field regression test.

## 16. İlişkiler
- [Go Server](go.md), [CI/CD](cicd.md), [Environment](environment.md), [Config](config.md), [System Architecture](system-architecture.md)

## 17. Bekleyen Kararlar
- Aggregator: Loki + Grafana (öneri, k8s-native, ucuz) vs OpenSearch (arama zenginliği) vs cloud native.
- Client crash reporting: Sentry vs equivalent.
- Log encryption at-rest & in-transit (compliance ihtiyacı ortaya çıkarsa).
- Sample ratesi (trace + debug).
