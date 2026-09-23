# Database — `activity` Schema

## Kapsam
Karakterin AFK / aktif aktivitelerinin kalıcı state'i. Contract (uzun ömürlü AFK), instance (session-scope, dungeon run), reveal log, event log. Server AFK Engine ve Activity Engine buradan besleniyor (bkz. [afk-engine.md](../systems/afk-engine.md), [activity-engine.md](../systems/activity-engine.md)).

---

## 1. `activity.contract`

Aktivitenin uzun ömürlü tanımı; karakter online/offline farkı gözetmez.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `character_id` | `UUID` | FK → character.character NOT NULL |
| `activity_type` | `activity.activity_type` | NOT NULL |
| `params_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `status` | `activity.contract_status` | NOT NULL DEFAULT `'active'` |
| `started_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `last_reveal_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `paused_at` | `TIMESTAMPTZ` | NULL |
| `pause_reason` | `VARCHAR(32)` | NULL — inventory_full, stamina, death, mob_dry |
| `seed` | `BYTEA` | NOT NULL — 32 byte |
| `rng_version_id` | `SMALLINT` | NOT NULL |
| `snapshot_ref` | `UUID` | NULL — activity.snapshot(id) |
| `soft_cap_override_ms` | `BIGINT` | NULL |
| `correlation_id_seed` | `UUID` | NOT NULL — reveal idempotency key üretimi için |
| `location_id` | `UUID` | FK NOT NULL |
| `channel` | `activity.contract_channel` | NOT NULL DEFAULT `'primary'` |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enums:
```sql
CREATE TYPE activity.activity_type AS ENUM (
  'mining','woodcutting','fishing','foraging','skinning',
  'blacksmithing','leatherworking','tailoring','alchemy','cooking','enchanting','jewelcrafting',
  'hunting','adventure','travel','dungeon_idle','crafting_generic',
  'rest_city'
);
CREATE TYPE activity.contract_status AS ENUM (
  'active','paused','completed','interrupted','archived'
);
CREATE TYPE activity.contract_channel AS ENUM ('primary','pet','mount_travel');
```

Constraint:
- Bir karakterin aynı anda **max 1 primary** contract'ı olabilir (partial unique):
  ```sql
  CREATE UNIQUE INDEX uq_contract__character_primary_active
  ON activity.contract (character_id)
  WHERE channel = 'primary' AND status IN ('active','paused');
  ```
- `pet` ve `mount_travel` kanalları için benzer partial unique.

Indexler:
- `ix_contract__character_id__status`.
- `ix_contract__status__last_reveal_at` (worker sweep).

---

## 2. `activity.snapshot`

Contract oluştuğunda karakter state'in dondurulmuş versiyonu. Reveal deterministic simülasyonu için gerekli girdiler.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `contract_id` | `UUID` FK NOT NULL |
| `taken_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `character_snapshot_json` | `JSONB` NOT NULL — stat, equip, skill rank, professions, mastery, race/class/spec ids |
| `content_version` | `VARCHAR(32)` NOT NULL — content.content_version.version_tag |

Notlar:
- Contract süresince stat değişmiş olsa da AFK simülasyonu snapshot'ı kullanır. Aktif oynama başlarsa contract yeniden şekillenir.
- Snapshot storage optimizasyonu: JSONB compress (Postgres TOAST).

---

## 3. `activity.instance`

Session-scope aktivite (dungeon run, party hunt); contract'tan farklı — kısa ömürlü, birden fazla üye içerebilir.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `template_id` | `UUID` NOT NULL — dungeon_template veya adventure_template |
| `template_kind` | `activity.instance_template_kind` NOT NULL |
| `party_id` | `UUID` FK → social.party NULL |
| `owner_character_id` | `UUID` FK — solo veya party leader |
| `difficulty` | `VARCHAR(16)` NULL — normal/heroic/mythic |
| `seed` | `BYTEA` NOT NULL |
| `rng_version_id` | `SMALLINT` NOT NULL |
| `current_encounter_index` | `SMALLINT` NOT NULL DEFAULT 0 |
| `state` | `activity.instance_state` NOT NULL DEFAULT `'active'` |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `expires_at` | `TIMESTAMPTZ` NOT NULL |
| `completed_at` | `TIMESTAMPTZ` NULL |

Enums:
```sql
CREATE TYPE activity.instance_template_kind AS ENUM ('dungeon','adventure','event_scenario');
CREATE TYPE activity.instance_state AS ENUM ('active','wiped','completed','abandoned','expired');
```

---

## 4. `activity.instance_participant`

Instance içindeki üyeler.

| Kolon | Tip |
|---|---|
| `instance_id` | `UUID` PK-part |
| `character_id` | `UUID` PK-part |
| `role_hint` | `VARCHAR(16)` NULL — tank/heal/dps |
| `joined_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `left_at` | `TIMESTAMPTZ` NULL |
| `is_leader` | `BOOLEAN` NOT NULL DEFAULT false |

---

## 5. `activity.reveal_log`

Her reveal işleminin özet kaydı. Idempotency ve audit için.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `character_id` | `UUID` NOT NULL |
| `contract_id` | `UUID` NOT NULL |
| `reveal_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `elapsed_ms` | `BIGINT` NOT NULL |
| `capped_ms` | `BIGINT` NOT NULL — soft cap uygulanmış efektif |
| `yield_summary_json` | `JSONB` NOT NULL — xp, gold, items[], mastery, ... |
| `events_ref` | `UUID` NULL — activity.event_batch(id), varsa detaylı log |
| `idempotency_key` | `TEXT` NOT NULL UNIQUE |
| `worker_id` | `VARCHAR(64)` NULL |

Indexler:
- `uq_reveal_log__idempotency_key`.
- `ix_reveal_log__character_id__reveal_at`.

Partitioning (Phase 02 aday): month-range on `reveal_at`.

---

## 6. `activity.event_batch`

Reveal veya instance içi event kümesi. Compressed detay.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `contract_id` | `UUID` NULL |
| `instance_id` | `UUID` NULL |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `events_compressed` | `BYTEA` NOT NULL — msgpack + zstd |
| `event_count` | `INT` NOT NULL |
| `hash` | `TEXT` NOT NULL — audit için |

Notlar:
- Detaylı per-tick log çoğu zaman gerekli değil; sadece disputed reveal veya rare drop için tam replay.
- Storage optimization: 30 gün sonra event_batch arşive taşınır.

---

## 7. `activity.death_event`

Aktivite sırasında ölüm.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `character_id` | `UUID` NOT NULL |
| `contract_id` | `UUID` NULL |
| `instance_id` | `UUID` NULL |
| `died_at` | `TIMESTAMPTZ` NOT NULL |
| `cause_kind` | `VARCHAR(24)` NOT NULL — mob/boss/env/other |
| `cause_ref_id` | `UUID` NULL |
| `xp_debt_amount` | `BIGINT` NOT NULL |
| `resurrected_at` | `TIMESTAMPTZ` NULL |

---

## 8. `activity.stamina` (Phase 02, opsiyonel)

Stamina sistemi eklenecekse.

```
character_id UUID PK,
current_stamina INT NOT NULL,
max_stamina INT NOT NULL,
last_recovery_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL
```

v1'de stamina yok; contract "inventory full" veya "location resource dry" ile pause olur.

---

## 9. `activity.rng_version`

RNG algoritma versiyonu tracking (Random System).

```
id SMALLSERIAL PK,
algorithm_code VARCHAR(32) NOT NULL — xoshiro256pp | chacha8 | ...
active BOOLEAN NOT NULL DEFAULT true,
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

---

## 10. Kullanım Örnekleri

**Login reveal:**
```sql
BEGIN;

SELECT c.id, c.character_id, c.activity_type, c.params_json, c.seed,
       c.last_reveal_at, c.status, c.rng_version_id
FROM activity.contract c
WHERE c.character_id = $1
  AND c.channel = 'primary'
  AND c.status IN ('active','paused')
FOR UPDATE;

-- Simulate reveal in application; then:
INSERT INTO activity.reveal_log (
  id, character_id, contract_id, reveal_at, elapsed_ms, capped_ms,
  yield_summary_json, idempotency_key
) VALUES (...);

UPDATE activity.contract
SET last_reveal_at = now(), status = $status, paused_at = $pausedAt, version = version + 1
WHERE id = $contract_id AND version = $expected_version;

-- Grant rewards via reward system (economy.ledger_entry, character.wallet, character.inventory_slot)
COMMIT;
```

**Aktif contract kontrolü (aktivite değiştirmeden önce):**
```sql
SELECT COUNT(*) FROM activity.contract
WHERE character_id = $1 AND channel = 'primary' AND status IN ('active','paused');
```

---

## 11. İlişkiler
- Contract → character, location, snapshot, activity_type.
- Instance → template (dungeon/adventure), party, participants.
- Reveal log → contract, event batch.
- Ödüller: `economy.ledger_entry`, `character.wallet`, `character.inventory_slot` üzerinden dağıtılır (Reward System soyutlamasıyla).

---

## 12. Notlar
- **Hard cap 30 gün**: contract archived; snapshot arşive taşınır; reveal safe fallback (bkz. AFK Engine §Fallback).
- **Contract rotation ile veri şişer**: eski contract'lar `archived` state'e alınır; ayrı archive tablo (partition benzeri) v2.
- **Yield summary şeması** versiyonlu (`yield_summary_json._v` alanı).
- **Test golden replay**: aynı contract + seed + snapshot + rng_version → aynı yield_summary. CI'da doğrulanır.

---

## 13. Bekleyen Kararlar
- Activity event batch storage: DB'de mi, object store'da mı (v1 DB, v2 hot→cold tier).
- Partition stratejisi (reveal_log, death_event) — Phase 02.
- Stamina modeli açılacak mı — Phase 02 opsiyonel.
