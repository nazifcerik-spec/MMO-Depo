# Logging

## 1. Amaç
Structured JSON log — collect, aggregate, query, retention, correlation.

## 2. Prensipler

- **Structured JSON** — key/value; no free-form text.
- **Correlation ID** her log satırında.
- **Levels**: `debug, info, warn, error, fatal`.
- **PII minimize** — email hash, IP hash, no password.
- **Sampling** high-volume (chat send info logs) — 10% sample.
- **No secret log EVER** — password, token, cvv.

## 3. Format

```json
{
  "ts": "2026-09-09T14:32:15.234Z",
  "level": "info",
  "service": "api",
  "version": "v1.4.0",
  "commit": "abc123",
  "trace_id": "d3f7...",
  "span_id": "f42b...",
  "correlation_id": "d3f7...",
  "account_id": "acct_...",
  "character_id": "char_...",
  "endpoint": "POST /v1/vendor/buy",
  "duration_ms": 82,
  "status_code": 200,
  "message": "vendor buy success",
  "extra": { "vendor_code": "general_ironhold", "item_code": "iron_sword" }
}
```

## 4. Stack

- **Collector**: Vector / Fluent Bit — sidecar or DaemonSet.
- **Storage**: **Loki** (Grafana) — cost-effective; log level filter + label index.
- Alternate: ELK (Elasticsearch + Logstash + Kibana) — more expensive.

## 5. Sink

App → stdout JSON → k8s collect → Loki.
Distinct streams per service label.

## 6. Retention

| Log Level | Retention |
|---|---|
| debug | 3 gün (staging only; prod OFF) |
| info | 14 gün |
| warn | 30 gün |
| error / fatal | 90 gün |
| audit (bkz. audit source) | 24 ay hot + cold |

Loki config tier retention.

## 7. Sampling

High-volume log categories:
- `chat_send` info: 10% sample.
- `ws_ping` info: 1% sample.
- Others: 100%.

Errors + warns: always 100%.

## 8. Correlation

- Trace ID = correlation_id (OpenTelemetry).
- Log include trace_id → Grafana Tempo "log to trace" navigate.
- Audit action correlate: `audit.admin_action.correlation_id` == log `correlation_id`.

## 9. Sensitive Field Handling

**Never log**:
- Password, hash, salt.
- Refresh token, access token (only prefix `<sha256[:8]>`).
- Payment card details.
- Full email (hashed or masked).

**Redact on collect** (Vector filter):
```
replace `Bearer\s+[A-Za-z0-9\.\-\_]+` → `Bearer <redacted>`
replace `\b[\w\.\-\+]+@[\w\.\-]+\.[a-z]{2,}\b` → `<email_masked>`
```

## 10. Query Examples

Loki LogQL:
```
{service="api", level="error"} |= "E_LEDGER"
{service="api"} | json | duration_ms > 1000
{service="worker", correlation_id="d3f7..."}  # trace all logs of one request
```

## 11. Alerts from Logs

- Error rate spike → Grafana Loki alert → Slack.
- Fatal log → PagerDuty (critical).
- Kritik keyword (`sec_ledger_sum_violation`) → PagerDuty.

## 12. Cost

Loki cheaper than ELK (label index only).
Prod T3 log volume: ~50-100 GB/gün.
Retention 14 gün info = 700 GB-1.4 TB.
Loki S3 backend cheap; annual ~$200-500.

## 13. Metrikler (Meta)
- `log_lines_total{service, level}`
- `log_lines_dropped_total{reason}`
- `log_shipping_lag_seconds`

## 14. Cross-Reference
- [monitoring.md](monitoring.md), [alerting.md](alerting.md)
- [architecture/logging.md](../architecture/logging.md) — v1 karar
- [security/audit.md](../security/audit.md)

## 15. Bekleyen Kararlar
- Loki vs ELK vs cloud provider — ops cost karar.
- Structured logger Go lib final (`slog` v1.22+; alternatif `zap`, `zerolog`) — Phase 01.
- Debug log prod'da OFF politikası enforce.
- PII field auto-redact policy final Phase 02.
