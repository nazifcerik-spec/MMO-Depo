# Database — `audit` + `system` Schemas

## Kapsam
Denetim (audit) izi, admin aksiyonları, güvenlik sinyalleri, sistem/schema metadatası, feature flag, migration versionu, background job kuyruğu.

---

# audit.*

## 1. `audit.event_log`

Genel amaçlı event log. Ekonomi ledger özel amaçlı (economy.ledger_entry); burada iş olayları + admin + güvenlik.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `BIGSERIAL` | PK |
| `event_kind` | `VARCHAR(48)` | NOT NULL — `character.created`, `character.ascended`, `guild.created`, `market.listing_created`, `security.suspicious_reveal`, ... |
| `severity` | `audit.severity` | NOT NULL DEFAULT `'info'` |
| `actor_kind` | `audit.actor_kind` | NOT NULL |
| `actor_id` | `UUID` | NULL |
| `subject_kind` | `VARCHAR(24)` | NULL — character/guild/item/... |
| `subject_id` | `UUID` | NULL |
| `payload_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `correlation_id` | `UUID` | NULL |
| `request_id` | `TEXT` | NULL |
| `at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Enums:
```sql
CREATE TYPE audit.severity AS ENUM ('debug','info','warn','error','critical');
CREATE TYPE audit.actor_kind AS ENUM ('character','account','user','admin','system','worker');
```

Indexler:
- `ix_event_log__event_kind__at DESC`.
- `ix_event_log__subject_kind__subject_id__at DESC`.
- `ix_event_log__severity__at DESC` (partial `WHERE severity IN ('warn','error','critical')`).
- BRIN on `at`.

Partitioning (Phase 02): month-range.
Retention: 12 ay+; downgrade of hot → cold storage.

---

## 2. `audit.admin_action`

Admin/GM işlemleri (grant item, ban user, edit character).

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `admin_user_id` | `UUID` FK → auth.user NOT NULL |
| `action_kind` | `VARCHAR(48)` NOT NULL — grant_item, revoke_item, adjust_currency, mute, kick_session, force_disband_guild, ... |
| `target_kind` | `VARCHAR(24)` NOT NULL |
| `target_id` | `UUID` NULL |
| `params_json` | `JSONB` NOT NULL DEFAULT `'{}'` |
| `reason` | `TEXT` NOT NULL |
| `ticket_ref` | `VARCHAR(64)` NULL — CRM ticket id |
| `at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |

Indexler: `ix_admin_action__admin_user_id__at`, `ix_admin_action__target_kind__target_id`.

Not: admin_action bir subset audit.event_log ile duplike; ancak dedicated table denetim / raporlama için ergonomik.

---

## 3. `audit.security_signal`

Otomatik güvenlik / anti-cheat sinyal.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `signal_kind` | `VARCHAR(48)` NOT NULL — reveal_yield_anomaly, market_flood, rapid_listing_cancel, botting_pattern, geo_login_anomaly |
| `severity` | `audit.severity` NOT NULL |
| `subject_kind` | `VARCHAR(24)` NOT NULL — account/character/session/ip |
| `subject_id` | `UUID` NULL |
| `evidence_json` | `JSONB` NOT NULL DEFAULT `'{}'` |
| `status` | `VARCHAR(16)` NOT NULL DEFAULT 'open' — open/reviewing/dismissed/actioned |
| `assigned_to_user_id` | `UUID` NULL |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `updated_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |

Indexler: `ix_security_signal__status__severity__created_at`, `ix_security_signal__subject_kind__subject_id`.

---

## 4. `audit.login_history`

Login/logout kayıtları — GDPR audit + security.

| Kolon | Tip |
|---|---|
| `id` | `BIGSERIAL` PK |
| `account_id` | `UUID` NOT NULL |
| `event_kind` | `VARCHAR(16)` NOT NULL — login/logout/login_fail/refresh/mfa_challenge |
| `client_kind` | `auth.client_kind` NOT NULL |
| `client_version` | `VARCHAR(32)` NULL |
| `ip_hash` | `TEXT` NULL |
| `country_hint` | `CHAR(2)` NULL — geo ip lookup snapshot |
| `success` | `BOOLEAN` NOT NULL |
| `fail_reason` | `VARCHAR(48)` NULL |
| `session_id` | `UUID` NULL |
| `at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |

Indexler: `ix_login_history__account_id__at DESC`.

---

# system.*

## 5. `system.schema_migrations`

Migration tool tarafından maintain (goose/atlas).

```
version BIGINT PK,
name VARCHAR(128) NOT NULL,
applied_at TIMESTAMPTZ NOT NULL DEFAULT now(),
checksum TEXT NOT NULL,
by_user VARCHAR(64) NULL
```

## 6. `system.content_version`

Content patch versiyonu — `content.content_version` referans; duplike/alias olarak. Tercih: sadece `content.content_version`.

## 7. `system.feature_flag`

v1'de flag mode `env` (config.md). Ancak service mode gerekirse:

| Kolon | Tip |
|---|---|
| `key` | `VARCHAR(64)` PK |
| `description` | `TEXT` NULL |
| `enabled_default` | `BOOLEAN` NOT NULL DEFAULT false |
| `variants_json` | `JSONB` NOT NULL DEFAULT `'{}'` — {"on":100,"off":0} percentage |
| `updated_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `updated_by` | `VARCHAR(64)` NULL |
| `notes` | `TEXT` NULL |
| `expires_at` | `TIMESTAMPTZ` NULL — 30 gün cleanup |

## 8. `system.feature_flag_targeting` (v2)

```
flag_key VARCHAR(64) PK-part,
target_kind VARCHAR(16) PK-part — account/character/guild/percent
target_id UUID NULL,
percent SMALLINT NULL,
variant VARCHAR(24) NOT NULL,
updated_at TIMESTAMPTZ NOT NULL
```

## 9. `system.job_queue`

Postgres SKIP LOCKED-based job kuyruğu (DD-OPEN-10 karar alınana kadar baseline).

| Kolon | Tip |
|---|---|
| `id` | `BIGSERIAL` PK |
| `job_kind` | `VARCHAR(48)` NOT NULL |
| `payload_json` | `JSONB` NOT NULL DEFAULT `'{}'` |
| `state` | `system.job_state` NOT NULL DEFAULT `'queued'` |
| `run_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `attempts` | `SMALLINT` NOT NULL DEFAULT 0 |
| `max_attempts` | `SMALLINT` NOT NULL DEFAULT 5 |
| `last_error` | `TEXT` NULL |
| `locked_by` | `TEXT` NULL |
| `locked_at` | `TIMESTAMPTZ` NULL |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `completed_at` | `TIMESTAMPTZ` NULL |
| `idempotency_key` | `TEXT` UNIQUE NULL |

Enum:
```sql
CREATE TYPE system.job_state AS ENUM ('queued','running','completed','failed','dead');
```

Indexler:
- `ix_job_queue__state__run_at` (worker pick).
- `uq_job_queue__idempotency_key` (unique partial).

Worker pattern:
```sql
UPDATE system.job_queue
SET state='running', locked_by=$worker_id, locked_at=now(), attempts=attempts+1
WHERE id IN (
  SELECT id FROM system.job_queue
  WHERE state='queued' AND run_at <= now()
  ORDER BY run_at ASC
  LIMIT 10 FOR UPDATE SKIP LOCKED
)
RETURNING *;
```

## 10. `system.job_history` (opsiyonel)

Uzun süreli iş kaydı için archive; live queue clean tutulur.

## 11. `system.server_health`

Server node heartbeat.

```
node_id VARCHAR(64) PK,
role VARCHAR(24) NOT NULL — api/worker/scheduler,
version VARCHAR(48) NOT NULL,
started_at TIMESTAMPTZ NOT NULL,
last_heartbeat_at TIMESTAMPTZ NOT NULL,
metrics_json JSONB NOT NULL DEFAULT '{}'
```

## 12. `system.scheduled_task`

Cron-benzeri scheduled tasks (event start, world boss spawn window, weekly reset).

```
id UUID PK,
task_kind VARCHAR(48) NOT NULL,
schedule_expr VARCHAR(64) NOT NULL — cron expr
next_run_at TIMESTAMPTZ NOT NULL,
last_run_at TIMESTAMPTZ NULL,
enabled BOOLEAN NOT NULL DEFAULT true,
notes TEXT NULL
```

## 13. `system.setting`

Runtime tweakable non-secret setting (v1'de config env; ancak baz için tablo tutulabilir).

```
key VARCHAR(64) PK,
value TEXT NOT NULL,
description TEXT NULL,
updated_at TIMESTAMPTZ NOT NULL,
updated_by VARCHAR(64) NULL
```

---

## 14. Kullanım Örnekleri

**Suspicious reveal alert:**
```sql
INSERT INTO audit.security_signal (
  signal_kind, severity, subject_kind, subject_id, evidence_json
) VALUES (
  'reveal_yield_anomaly', 'warn', 'character', $1,
  '{"expected_yield_max":5000,"actual":12400,"contract_id":"..."}'
);
```

**Job dispatch (reveal batch):**
```sql
INSERT INTO system.job_queue (job_kind, payload_json, run_at, idempotency_key)
VALUES ('reveal_batch', '{"character_id":"..."}', now(), $key)
ON CONFLICT (idempotency_key) DO NOTHING;
```

**Weekly reset trigger:**
```sql
UPDATE system.scheduled_task
SET last_run_at = now(),
    next_run_at = now() + interval '7 days'
WHERE task_kind = 'weekly_reset' AND enabled;
```

---

## 15. İlişkiler
- audit.event_log: cross-schema referans (subject_kind + subject_id).
- audit.admin_action → auth.user (admin), + target.
- audit.security_signal → account/character/session referansları.
- system.job_queue: any domain queues jobs; payload_json domain-specific.
- system.scheduled_task: worker consumer.

---

## 16. Notlar
- **Append-only audit**: UPDATE/DELETE yasak (role permission).
- **Retention**: audit.event_log 12 ay (compliance ile revise); login_history 6 ay.
- **PII**: ip hash'li, email hiç değil.
- **Job idempotency**: idempotency_key üretici tarafında correlation-based.
- **DR**: audit tabloları PITR kapsamında; ancak backup priority downstream.
- **Job kuyruğu alternatifleri (DD-OPEN-10)**: Redis Streams, NATS. Şeması ayrı — v1 Postgres SKIP LOCKED yeterli.

---

## 17. Bekleyen Kararlar
- Log aggregator gateway'i: audit → cold storage (S3-benzeri) rotation.
- Feature flag service (v2) → schema tam olgunlaştıktan sonra.
- Real-time anomaly detection pipeline (Phase 02).
- Retention policy compliance (region-specific): Phase 02 legal review.
