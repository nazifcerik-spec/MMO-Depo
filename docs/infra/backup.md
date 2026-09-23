# Backup + Restore

## 1. Amaç
Data loss risk minimize. Postgres, Redis, content pack, admin secret backup + restore drill.

## 2. RTO / RPO Hedefleri (v1)

| System | RTO (recovery time) | RPO (data loss) |
|---|---|---|
| Postgres primary fail | 5 dk (replica promote) | ~seconds (streaming) |
| Postgres full restore | 1-2 saat | 15 dk (WAL archive) |
| Redis fail | 10 dk | 1 saat (snapshot) |
| Content pack | 30 dk | git commit (versioned) |
| Admin secret | 15 dk | Vault snapshot 1 saat |
| Static assets | 15 dk | Immutable (git+CDN) |

## 3. Postgres Backup

### Continuous WAL Archive
- `archive_command`: WAL segment → S3 (wal-g, wal-e, veya pgBackRest).
- Interval: her WAL segment (~16MB) veya 15 dk force switch.

### Base Backup
- Nightly `pg_basebackup` → S3.
- Retention: 7 gün daily + 4 hafta weekly + 12 ay monthly.

### PITR (Point-In-Time Recovery)
- Herhangi bir zamana restore edilebilir (WAL replay).
- Kullanım: bug/exploit data corruption → belirli anı geri al.
- Test: nightly canary restore (5 GB test DB, çalışıyor mu doğrula).

### Managed Provider
- RDS: automated backup 7-35 gün retention; snapshot on-demand.
- Cloud SQL: similar.

## 4. Redis Backup

### AOF + RDB
- AOF everysec (durability).
- RDB snapshot hourly + boot.

### Offsite Copy
- Snapshot → S3 hourly (self-managed).
- Managed: snapshot daily default; increase to hourly Phase 02.

### Restore
- Snapshot download → `redis-server --loadmodule ...` → verify.

## 5. Content Pack

- Git repository (versioned YAML files).
- Deploy: signed tarball → CDN.
- Restore: git checkout + redeploy.

## 6. Admin Secret

- Vault (HashiCorp) veya AWS Secrets Manager veya Doppler.
- Backup: Vault snapshot hourly → offsite bucket (encrypted).
- Restore: Vault init + snapshot load.

## 7. Static Asset

- S3 versioning ON.
- Immutable hash path → historical asset accessible.
- Restore: hiç değişmez; ancak S3 accidental delete → object versioning ile restore.

## 8. Offsite / Cross-Region

- Backup dedicated **cross-region S3 bucket** (production region ≠ backup region).
- Prevents regional AWS outage total loss.
- Retention same 12 ay.

## 9. Encryption

- **At rest**: S3 SSE-KMS zorunlu; RDS storage encrypted.
- **In transit**: TLS her transfer.
- **Key rotation**: KMS auto-rotate yearly.

## 10. Restore Drill

Aylık zorunlu:
1. Fresh env spin up.
2. Restore Postgres from S3 backup (random gün seçim).
3. Restore Redis snapshot.
4. Data integrity check (row counts, sample queries).
5. Report: RTO/RPO actual vs target.

Runbook: [rollback.md](rollback.md) referans.

## 11. Immutable / WORM Backup

- S3 Object Lock (compliance mode) retention 30 gün — ransomware defense.
- Even admin can't delete during lock period.

## 12. Metrikler

- `backup_success_total{kind, target}`
- `backup_last_success_timestamp{kind}`
- `backup_size_bytes{kind}`
- `backup_restore_test_success_total`
- `backup_restore_duration_seconds`

Alerts:
- Backup age > 25 saat (daily): warn.
- Backup fail: high.
- WAL archive fail: critical (RPO risk).

## 13. Compliance

- GDPR/KVKK: kullanıcı silme request → backup içinde de silinmeli (Phase 02 tooling).
- Retention max: 60 ay (yasal minimum sonrası temizle).

## 14. Cross-Reference
- [postgres-prod.md](postgres-prod.md), [redis-prod.md](redis-prod.md)
- [rollback.md](rollback.md), [monitoring.md](monitoring.md)
- [security/audit.md](../security/audit.md) — retention compliance

## 15. Bekleyen Kararlar
- Backup tool: wal-g vs pgBackRest vs Barman — ops karar.
- Secret manager final (Vault / AWS SM / Doppler) — Phase 01.
- Restore drill sıklık (aylık v1; Phase 02 haftalık).
- GDPR erasure tooling backup dahil Phase 02.
