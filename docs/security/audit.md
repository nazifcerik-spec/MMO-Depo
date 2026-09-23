# Security Audit + Anomaly Detection + Response

## 1. Amaç
Güvenlik olaylarının kayıt altına alınması, korelasyon, anomaly detection, incident response.

## 2. Audit Kaynakları (Security Ekseni)

Bkz. [database/audit.md](../database/audit.md) + [admin/logs.md](../admin/logs.md):

- `audit.security_event` — login fail, token reuse, geo anomaly, session revoke.
- `audit.admin_action` — admin panel her aksiyon.
- `audit.ledger_audit` — ekonomi transaction mirror.
- `audit.trade_log` — P2P trade lifecycle.
- `audit.chat_report` — user report.
- `analytics.combat_event` — combat forensic (24h hot).
- `analytics.market_search_event` — market query pattern.

Retention: 24 ay hot; parquet cold arşiv.

## 3. Security Event Model

```sql
CREATE TABLE audit.security_event (
  id BIGSERIAL PRIMARY KEY,
  ts TIMESTAMPTZ NOT NULL DEFAULT now(),
  event_kind VARCHAR(48) NOT NULL,      -- 'login_fail'|'refresh_reuse'|'geo_anomaly'|'idempotency_mismatch'|...
  severity VARCHAR(16) NOT NULL,        -- 'info'|'warn'|'high'|'critical'
  account_id UUID NULL,
  character_id UUID NULL,
  session_id UUID NULL,
  ip_hash BYTEA NULL,
  device_id UUID NULL,
  user_agent VARCHAR(256),
  detail JSONB,
  correlation_id UUID NULL
);
CREATE INDEX idx_sec_event_ts ON audit.security_event (ts DESC);
CREATE INDEX idx_sec_event_kind ON audit.security_event (event_kind, ts DESC);
CREATE INDEX idx_sec_event_account ON audit.security_event (account_id, ts DESC);
CREATE INDEX idx_sec_event_severity ON audit.security_event (severity, ts DESC) WHERE severity IN ('high','critical');
```

## 4. Anomaly Detection (v1 Rule-Based)

Worker `security_anomaly_watcher` her 1 dk:

### Kural Örnekleri
- **Credential stuffing**: aynı IP 5dk içinde 20+ login fail (farklı account) → auto-block IP 1h + alert.
- **Brute force**: aynı account 15dk içinde 5+ login fail → account lock 15dk.
- **Impossible geo**: aynı account 10dk içinde iki farklı country login → session invalidate + email.
- **Refresh reuse**: aynı refresh token 2. kullanım → session family revoke + P1 alert.
- **RMT bulk trade**: yeni account 24h içinde 5M+ gold gönderdi → flag + admin queue.
- **Ledger drift**: reconciler ≥ 1c → character freeze + P0 page.
- **Sum violation**: ledger sum≠0 write attempt → P0 alert (bug veya intrusion).
- **Chat spam**: aynı account 5dk 100+ chat → auto-mute 24h + report.
- **Item flow spike**: aynı character 1h 100+ item transfer → RMT flag.

### Rule Yönetimi
Config `security/rules.yml` hot reload. Yeni kural eklemek code deploy gerektirmez (data-driven).

## 5. Alert Sınıfları

| Severity | Örnek | Response |
|---|---|---|
| `info` | Login fail | Log only |
| `warn` | Geo anomaly | Inbox notify player |
| `high` | Credential stuffing | Slack #sec channel + admin dashboard |
| `critical` | Sum violation, mass ban trigger | PagerDuty + on-call oncall |

Alert delivery: metrics + Prometheus alertmanager → PagerDuty/Slack/email.

## 6. Incident Response Runbook (Özet)

**Level 1 (auto)**: rule tetiklendi → auto action (IP block, session revoke, character freeze).

**Level 2 (support)**: admin dashboard'da flag → support l2 review → action (ban, reversal).

**Level 3 (ops)**: P0/P1 alert → on-call → root cause + hotfix + post-mortem.

**Post-mortem**: her critical/high olay için doc (Notion / GitHub issue) — timeline, impact, fix, prevention.

## 7. Report Akışı (User)

Player `POST /v1/chat/report`, `POST /v1/user/report` (character):
```json
{ "target_character_id": "...", "reason": "harassment", "context_ref": "chat_msg_id or trade_id" }
```

Insert `audit.chat_report` / `audit.user_report`; aggregation:
- 3+ report / 1 saat aynı target → auto-mute 24h.
- Support review queue.

## 8. Bug Bounty (Phase 02)

- Dış program (HackerOne / Bugcrowd).
- Kapsam: production API, admin, auth.
- Ödül matrisi: critical $$$, high $$, medium $, low swag.
- SLA: initial response 48h; fix critical 7d.

Program açılışı ops + legal ready olduktan sonra Phase 02.

## 9. Compliance

- **GDPR / KVKK**: right-to-erasure endpoint (`POST /v1/account/delete-request`); 30 gün grace; audit'te retention 24-60 ay compliance minimum.
- **COPPA** (US 13 yaş altı): register 13+ zorunlu (self-declared v1); parental consent Phase 02.
- **PCI-DSS**: v1'de payment yok; Phase 02 store için gerekli — third-party (Stripe) handle.

## 10. Metrikler
- `sec_event_total{kind, severity}`
- `sec_alert_fired_total{severity}`
- `sec_auto_action_total{action}` (IP block, session revoke, character freeze)
- `sec_report_filed_total{reason}`
- `sec_incident_post_mortem_total`

## 11. Cross-Reference
- [database/audit.md](../database/audit.md), [admin/logs.md](../admin/logs.md)
- [replay-protection.md](replay-protection.md), [session-security.md](session-security.md)
- [economy/transaction.md](../economy/transaction.md) — reconciler

## 12. Bekleyen Kararlar
- ML anomaly detection Phase 03.
- Bug bounty program launch Phase 02.
- Compliance legal review Phase 02.
- Post-mortem tool (Notion vs GitHub) infra karar.
- Retention final policy compliance ile Phase 02.
