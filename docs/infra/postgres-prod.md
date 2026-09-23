# PostgreSQL Production

## 1. Amaç
Prod Postgres deployment: HA, replication, connection pool, tuning, extension, migration.

## 2. Sürüm

- **PostgreSQL 16.x** (LTS-benzeri stability); Phase 02 minor upgrade cycle.
- Extension: `pg_stat_statements`, `pg_partman` (Phase 02), `pgcrypto`, `uuid-ossp`.

## 3. Topology

### T3 Launch
- 1 primary (write) + 1 hot standby (read replica, sync/async karar).
- WAL streaming replication.
- Async replica lag < 5s hedef.

### T4 Peak
- 1 primary + 2 read replicas.
- Load balancer (pgpool veya PgBouncer routing).

### T5 Phase 03
- Sharding karar (Citus / manuel partitioning / cross-account bölünme).

## 4. Deployment

- **Managed**: RDS / Cloud SQL — v1 tercih (ops simplification).
- **Self-managed K8s**: Patroni + PostgreSQL Operator (Phase 02 cost optimize).

## 5. Connection Pool

**PgBouncer** application ↔ Postgres arasında:
- Mode: `transaction` (session state SET LOCAL kullanılmaz).
- Max client conn: 1000; DB backend conn: 100 (per DB).
- App pool (pgx) → PgBouncer → Postgres.

Total backend conn cap:
- API app: 200.
- Worker: 50.
- Admin: 30.
- Migration: 10.
- Analytics: 20.
Total 300 backend conn @ Postgres.

## 6. Tuning (T3 baseline)

`postgresql.conf` kritikler (8vCPU/32GB):
```
shared_buffers = 8GB              # ~25% RAM
effective_cache_size = 24GB       # ~75% RAM
maintenance_work_mem = 1GB
work_mem = 32MB                   # sortu, hash
max_connections = 300
wal_buffers = 16MB
max_wal_size = 4GB
min_wal_size = 1GB
checkpoint_completion_target = 0.9
random_page_cost = 1.1            # SSD
default_statistics_target = 100
log_min_duration_statement = 500  # slow query log
```

pg_stat_statements + `auto_explain.log_min_duration = 500ms` slow query dashboard.

## 7. High Availability

- Failover automatic (RDS Multi-AZ veya Patroni).
- RTO 1-5 dk; RPO ~seconds (streaming replication).
- Replica promote on primary failure.
- Read traffic (search, admin, analytics) → replica.
- Write always primary.

## 8. Migration

- **Tool**: `sqlx migrate` veya `golang-migrate` — Go stack uyumlu.
- **Naming**: `<timestamp>_<slug>.up.sql` + `_.down.sql`.
- **Reversible**: her migration up + down çift; test `up → down → up` idempotent.
- **Deploy order**: migration önce (backward compat), sonra app deploy.
- **Zero-downtime rules**:
  - Add column: allowed (nullable veya default).
  - Drop column: 2-step (app önce stop kullanımı → sonra migration).
  - Rename column: 3-step (add new + dual write + drop old).
  - Add index CONCURRENTLY.
  - No table rewrite lock in one deploy.

Bkz. [rollback.md](rollback.md) migration reversibility.

## 9. Backup

Bkz. [backup.md](backup.md):
- Continuous WAL archive → S3.
- Nightly full base backup.
- PITR (Point-In-Time Recovery).
- RPO 15 dk (WAL archive interval); RTO ~1 saat full restore.

## 10. Monitoring

- pg_exporter → Prometheus.
- Key metrics:
  - `pg_stat_activity_count` (conn count).
  - `pg_stat_bgwriter_*` (checkpoint pressure).
  - `pg_replication_lag_bytes/seconds`.
  - `pg_stat_user_tables_*` (bloat, dead tuples).
  - `pg_stat_statements_top` (slow query).
- Grafana dashboard: connections, throughput, cache hit, replication lag.

Alerts:
- Replica lag > 60s: warn.
- Connection pool > 90%: warn.
- Slow query rate spike: warn.
- Failover: critical.

## 11. Vacuum / Autovacuum

- `autovacuum` ON (default) — tune per table:
  - High-write tables (ledger_entry, chat.message): `autovacuum_vacuum_scale_factor = 0.02`.
  - Cold tables: default.
- Reindex CONCURRENTLY monthly (bloat control) — cron.

## 12. Partitioning (Phase 02)

Büyük tablolar (`chat.message`, `analytics.combat_event`, `economy.ledger_entry`) → aylık partition:
- pg_partman extension.
- Cold partitions detach → S3 parquet arşiv.

## 13. Security

- TLS required (`ssl = on`).
- Auth: `scram-sha-256`.
- Role least-privilege: app_writer, app_reader, migration_admin.
- Password rotation quarterly (Phase 02 via Vault).
- No superuser in app.

## 14. Cross-Reference
- [backup.md](backup.md), [monitoring.md](monitoring.md)
- [architecture/postgresql.md](../architecture/postgresql.md) — v1 temel karar
- [performance/capacity-planning.md](../performance/capacity-planning.md)

## 15. Bekleyen Kararlar
- Managed RDS vs self-managed Patroni — ops cost/control karar.
- Sync vs async replication — sync yavaş, async lag risk (v1 async + short lag).
- Partition strategy content pack Phase 02.
- Sharding roadmap Phase 03.
