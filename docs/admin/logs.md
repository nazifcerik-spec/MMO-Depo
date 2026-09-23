# Logs (Audit Log Viewer)

## 1. Amaç
Tüm sistem audit log'larını (admin action, ekonomi, güvenlik, moderation) filter + query + export.

## 2. Kaynaklar (Aggregate)

- `audit.admin_action` — admin panel'de yapılan aksiyonlar.
- `audit.ledger_audit` — ekonomi transaction mirror.
- `audit.chat_report` — kullanıcı chat report.
- `audit.trade_log` — P2P trade complete/cancel.
- `audit.mail_action` (Phase 02).
- `audit.security_event` — login fail, ban trigger, device revoke.
- `audit.guild_action` — leader change, disband, MOTD edit, bank withdraw.
- `audit.market_admin_action` — admin cancel bundle.

Retention: her audit tür 24 ay hot; cold arşiv sonra parquet (compliance).

## 3. Query Builder

```
Source: [multi-select]
Actor: (admin_user_id / character_id / system)
Target: (kind + id)
Action: (enum by source)
From / To: date range
Search text: (reason field)
```

Result grid + expand row (before/after JSON diff).

## 4. Export

CSV / JSONL export; row cap 100k / query.
`GET /internal/admin/logs/export?…&format=csv` — asenkron worker; email link when ready.

## 5. Correlation Trace

Bir action_id verildiğinde ilgili tüm audit satırları (aynı `correlation_id` veya `target_id`) tek view.

Ör: Ledger reversal → önceki original txn + admin_action + character notification event.

## 6. Aksiyon

Audit üzerinde yalnızca **read-only**. Immutable — silme veya edit yok.

## 7. Endpoint

- `GET /internal/admin/logs?source=&actor=&target_kind=&target_id=&action=&from=&to=&search=&limit=&cursor=`
- `GET /internal/admin/logs/correlation/{correlation_id}`
- `POST /internal/admin/logs/export` `{ query, format }` → job_id
- `GET /internal/admin/logs/export/{job_id}` — status + download link

## 8. Metrikler
- `admin_logs_query_total{source}`
- `admin_logs_export_total`
- `admin_logs_correlation_query_total`

## 9. Retention Politikası

- 24 ay hot Postgres.
- 24-60 ay parquet S3 cold arşiv (compliance query slow).
- 60+ ay silme (yasal min saklama süresi hariç).

## 10. Cross-Reference
- [database/audit.md](../database/audit.md)
- [economy/transaction.md](../economy/transaction.md), [auth-rbac.md](auth-rbac.md)

## 11. Bekleyen Kararlar
- Correlation ID propagation tüm event bus'a (Phase 01 audit schema).
- Full-text search Postgres vs Elastic (Phase 02).
- GDPR / KVKK right-to-erasure entegrasyonu Phase 02.
- Export approval workflow (PII içeriyorsa) Phase 02.
