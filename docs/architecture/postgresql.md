# PostgreSQL (Durable State)

## 1. Amaç
Karakter, envanter, ekonomi, sosyal grafiği, dünya içeriği ve audit için kaynak-of-truth. ACID zorunluluğu bu tercihi yönlendirir.

## 2. Sürüm
- **PostgreSQL 16+** (Phase 01 başlangıcı; LTS-benzeri sürüm politikası: major upgrade planlı yılda 1).
- Extension'lar (varsayılan aktif):
  - `uuid-ossp` (uuid_generate_v4) veya `pgcrypto` (gen_random_uuid).
  - `pg_stat_statements` (perf analiz).
  - `pg_trgm` (chat/isim arama).
  - `citext` (case-insensitive index — karakter/guild adı).

## 3. Şema Grupları (PostgreSQL "schema" bazlı)
```
character   -- character profile, stats, inventory, equipment, activity_contract, mastery, ascension
content     -- item_template, race, class, spec, skill, recipe, dungeon, race, npc_template, loot_table
economy     -- currency_ledger, market_listing, market_transaction, mail
social      -- guild, guild_member, guild_bank, party, friendship, chat_message, report
world       -- continent, region, territory, city, location, poi, resource_node
audit       -- event_log, admin_action, security_signal
system      -- migrations, feature_flags, schema_version
```

Neden schema-bazlı? Ownership + permission granularity + query prefix okunabilirliği.

## 4. Kritik Design Kararları
- **UUID PK** (v4/v7) çoğu tablo için; `bigserial` yalnızca yüksek-throughput audit tablo (ledger event id) için değerlendirilir.
- **Timestamp**: her mutable satırda `created_at`, `updated_at`; `TIMESTAMPTZ` (UTC).
- **Soft delete**: `deleted_at TIMESTAMPTZ NULL`, indexli — hard delete rarely.
- **Optimistic locking**: `version INT` sütunu concurrency kritik tablolarda (character, inventory).
- **JSONB**: content template (item.affix pool, race.stat modifier) ve esnek payload'lar için — sorgulanan alanlar generated columns ile indexlenir.
- **Enum**: küçük, kararlı setler için PG `ENUM`; büyük/değişken setler için lookup table.

## 5. Transaction Politikaları
- **Serializable** isolation kritik ekonomik işlemlerde: currency_ledger writes, market_transaction.
- **Read Committed** (default) çoğu okuma.
- Retry policy: serializable serialization_failure → 3 kez retry (exponential jitter).
- Uzun transaction yasak (> 500ms warning; > 5s auto-abort).

## 6. Indexler
- Her FK sütununa index (Postgres otomatik değil).
- Marketplace hot-path: `(item_template_id, status, unit_price)` composite.
- Currency ledger: `(character_id, created_at DESC)`.
- Activity contract: `(character_id, status)` partial index.
- Full-text: chat_message ve item search için `pg_trgm` gin index (Phase 02 tune).

## 7. Partitioning (Phase 02 aday)
- Audit / ledger tabloları: month-range partition (retention + archive kolaylığı).
- Chat message: month-range.
- Karakter tabloları partition yok (v1 tek shard; ölçek başlayınca değerlendir).

## 8. Ledger Discipline
- `economy.currency_ledger` append-only.
- Idempotency key unique index.
- Trigger yok — application layer'da audit; transparency.
- Daily reconcile job: `SUM(delta) per character == wallet balance`; fark alarm.

## 9. Content Data Loading
- `content.*` tabloları content patch'te versionlu.
- Hot reload: workers/CLI ile stage → validate → swap versioning; running server API'sinde şüpheli ise version parametresi ile eski/yeni serve.
- Migration: forward-only tercih; rollback path test edilmiş olmalı.

## 10. Backup & Recovery
- **PITR**: WAL archive (S3-benzeri object store); RPO ≤ 1 dk.
- **Nightly logical dump** (`pg_dump`): kolay restore + smoke test.
- **RTO ≤ 15 dk**: playbook + rehearsal (çeyrek yılda 1).
- Test restore: staging'e ayda 1 kez.

## 11. Migration Kuralları
- `goose` veya `atlas` — tek dil (SQL forward).
- Migration adları: `NNNN_short_description.sql`.
- Naming: schema.tablo.column snake_case.
- Breaking migration (drop column) 2-adımlı: (a) dual-write / dual-read, (b) drop.

## 12. Query Style
- **Explicit column list**; `SELECT *` yasak (audit/analytics dışı).
- Prepared statement (pgx default).
- N+1 avoid: `IN (...)` batch veya `LATERAL JOIN` gerektiğinde.
- Query timeout: default 3s (SET statement_timeout).

## 13. Security
- Ayrık DB user'lar: `app_rw`, `app_ro`, `migrate`, `admin`. Her app node kendi role'ünü kullanır.
- TLS zorunlu (in-transit).
- Row-level security: v1 gerekli değil; v2'de admin panel scenarios için.
- PII sütunları (email) yalnız `auth` schema'sında, kısıtlı erişim.

## 14. Monitoring
- `pg_stat_statements` top-N query.
- Connection count, autovacuum lag, replication lag.
- Slow query log > 500 ms.
- Table bloat & vacuum health.

## 15. Anti-Patterns
- Application-side JOIN yerine 5 ayrı query.
- ORM'in "lazy load" mode (N+1).
- `SELECT FOR UPDATE` uzun transaction içinde.
- Trigger-based business logic (opaque).
- Datetime string'i olarak saklamak (mutlaka `TIMESTAMPTZ`).

## 16. İlişkiler
- [System Architecture](system-architecture.md), [Go Server](go.md), [Redis](redis.md), [Logging](logging.md), [Environment](environment.md), [Economy](../game-design/economy.md), [Marketplace](../game-design/marketplace.md)

## 17. Bekleyen Kararlar
- Migration tool final: `goose` vs `atlas`.
- Query codegen: `sqlc` vs raw.
- v1 replica sayısı (1 primary + 1 hot standby minimum).
- Multi-region strategy (Phase 02).
