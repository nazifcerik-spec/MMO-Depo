# Database — `character` Schema

## Kapsam
Karaktere sahip olduğu tüm veri: profil, statlar, envanter, ekipman, öğrenilmiş skill'ler, meslek seviyeleri, mastery, ascension, ayarlar.

Template/data-driven veriler `content` schema'sındadır (bkz. [content-templates.md](content-templates.md)); burada yalnızca **oyuncunun kişisel state'i**.

---

## 1. `character.character`

Ana karakter tablosu.

| Kolon | Tip | Constraint | Not |
|---|---|---|---|
| `id` | `UUID` | PK | |
| `account_id` | `UUID` | FK → `auth.account(id)` ON DELETE RESTRICT | Silmeden önce karakter arşive |
| `name` | `CITEXT` | UNIQUE, NOT NULL, `CHECK length between 3 and 16` | Case-insensitive unique |
| `race_id` | `UUID` | FK → `content.race(id)` ON DELETE RESTRICT | |
| `class_id` | `UUID` | FK → `content.class(id)` ON DELETE RESTRICT | |
| `spec_tier1_id` | `UUID` | FK → `content.specialization(id)` NULL | Level 5'te set |
| `spec_tier2_id` | `UUID` | FK → `content.specialization(id)` NULL | Level 10'da set |
| `level` | `SMALLINT` | NOT NULL DEFAULT 1, `CHECK level between 1 and 200` | |
| `xp` | `BIGINT` | NOT NULL DEFAULT 0 | |
| `ascension_level` | `SMALLINT` | NOT NULL DEFAULT 0 | |
| `ascension_points_total` | `INT` | NOT NULL DEFAULT 0 | |
| `ascension_points_spent` | `INT` | NOT NULL DEFAULT 0 | |
| `home_city_id` | `UUID` | FK → `world.city(id)` NULL | Respawn |
| `current_location_id` | `UUID` | FK → `world.location(id)` NULL | |
| `current_activity_id` | `UUID` | FK → `activity.contract(id)` NULL | Aktif contract |
| `active_pet_id` | `UUID` | FK → `companion.pet_owned(id)` NULL | |
| `active_mount_id` | `UUID` | FK → `companion.mount_owned(id)` NULL | |
| `cosmetic_json` | `JSONB` | NOT NULL DEFAULT `'{}'` | Görünüm (transmog Phase 02 placeholder) |
| `settings_json` | `JSONB` | NOT NULL DEFAULT `'{}'` | UI tercihleri |
| `flags_json` | `JSONB` | NOT NULL DEFAULT `'{}'` | Tutorial, unlock, bool flag |
| `version` | `INT` | NOT NULL DEFAULT 1 | Optimistic locking |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() | |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() | |
| `last_seen_at` | `TIMESTAMPTZ` | NULL | Aktif oyun anı |
| `last_reveal_at` | `TIMESTAMPTZ` | NULL | AFK Engine reveal noktası |
| `soft_deleted_at` | `TIMESTAMPTZ` | NULL | 7 gün grace |
| `seed` | `BYTEA` | NOT NULL | 32 byte; RNG için karakter kökü |

Indexler:
- `uq_character__name`.
- `ix_character__account_id`.
- `ix_character__soft_deleted_at` (partial).
- `ix_character__last_seen_at`.
- `ix_character__level` (leaderboard, filter).

Constraint:
- `CHECK (spec_tier2_id IS NULL OR spec_tier1_id IS NOT NULL)` — tier 2 bağımsız olamaz.
- `CHECK (ascension_points_spent <= ascension_points_total)`.

---

## 2. `character.stats`
Level up allocation ve derived stat cache.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK, FK → character.character |
| `str` | `INT` | NOT NULL DEFAULT 0 |
| `dex` | `INT` | NOT NULL DEFAULT 0 |
| `int_` | `INT` | NOT NULL DEFAULT 0 — kolon adı çakışması için trailing underscore |
| `wis` | `INT` | NOT NULL DEFAULT 0 |
| `sta` | `INT` | NOT NULL DEFAULT 0 |
| `cha` | `INT` | NOT NULL DEFAULT 0 |
| `unspent_points` | `INT` | NOT NULL DEFAULT 0 |
| `derived_cache_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — AP/SP/crit/haste/... |
| `derived_cache_invalidated_at` | `TIMESTAMPTZ` | NULL |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Notlar:
- Primary stat'lar 0'dan başlar; level up + item + race/class + ascension bonusu **application** tarafında toplanır; ancak baz allocation burada.
- Derived stat cache invalidation: item takıldığında `derived_cache_invalidated_at = now()`; sonraki okumada recompute.

---

## 3. `character.inventory_slot`
Karakter üstünde taşınan çantayla.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part, FK → character |
| `slot_index` | `SMALLINT` | PK-part, `CHECK slot_index between 0 and 119` |
| `item_instance_id` | `UUID` | FK → `content.item_instance(id)` ON DELETE RESTRICT |
| `stack_count` | `INT` | NOT NULL DEFAULT 1, `CHECK stack_count > 0` |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Composite PK: `(character_id, slot_index)`.
Indexler:
- `ix_inventory_slot__item_instance_id` (FK).

Notlar:
- Slot cap Phase 00'da 40 default; bag item ile 120'ye kadar.
- Boş slot = row yok; sparse representation.
- Ağırlık yok (DD-021).

---

## 4. `character.bank_slot`
Şehir bankası; tab bazlı.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `tab_index` | `SMALLINT` | PK-part, `CHECK tab_index between 0 and 20` |
| `slot_index` | `SMALLINT` | PK-part |
| `item_instance_id` | `UUID` | FK → item_instance |
| `stack_count` | `INT` | NOT NULL |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL |

Composite PK: `(character_id, tab_index, slot_index)`.

---

## 5. `character.bank_tab`
Bank tab metadata.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `tab_index` | `SMALLINT` | PK-part |
| `name` | `VARCHAR(32)` | NOT NULL DEFAULT '' |
| `purchased_at` | `TIMESTAMPTZ` | NOT NULL |
| `purchase_cost_gold` | `BIGINT` | NOT NULL |
| `capacity` | `SMALLINT` | NOT NULL DEFAULT 40 |

---

## 6. `character.equipment`
Ekipman slot ↔ item instance mapping.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `slot` | `character.equipment_slot` | PK-part |
| `item_instance_id` | `UUID` | FK → item_instance NOT NULL |
| `equipped_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Enum:
```sql
CREATE TYPE character.equipment_slot AS ENUM (
  'head','neck','shoulder','back','chest','wrist','hands','waist','legs','feet',
  'ring1','ring2','trinket1','trinket2',
  'main_hand','off_hand','ranged'
);
```

Notlar:
- `ring1/ring2`, `trinket1/trinket2` iki slot; UI hangisine takılacağı kararı.
- 2-handed weapon `main_hand` slot; equipment logic `off_hand` boşaltır (transaction).

---

## 7. `character.skill_learn`
Öğrenilmiş skill'ler ve rank'ları.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `skill_id` | `UUID` | PK-part, FK → content.skill_template |
| `rank` | `SMALLINT` | NOT NULL DEFAULT 1, `CHECK rank between 1 and 10` |
| `learned_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `xp` | `INT` | NOT NULL DEFAULT 0 — use-based rank up için |

Indexler:
- `ix_skill_learn__character_id` (skill listeleme).

---

## 8. `character.skill_bar`
Aktif skill slot bar (UI için).

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `bar_index` | `SMALLINT` | PK-part `CHECK bar_index between 0 and 11` |
| `skill_id` | `UUID` | FK → content.skill_template NOT NULL |

Max 12 slot; boş slot yok row.

---

## 9. `character.rotation`
AFK / auto rotation editor.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK |
| `priority_list` | `JSONB` | NOT NULL DEFAULT `'[]'` |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Priority list şeması (JSON):
```json
[
  { "cond": {"hp_below_pct": 50}, "skill_id": "..." },
  { "cond": {"any": true}, "skill_id": "..." }
]
```

---

## 10. `character.profession`
Karakterin öğrendiği meslekler.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `profession_id` | `UUID` | PK-part, FK → content.profession |
| `level` | `SMALLINT` | NOT NULL DEFAULT 1 |
| `xp` | `BIGINT` | NOT NULL DEFAULT 0 |
| `is_active` | `BOOLEAN` | NOT NULL DEFAULT true |
| `learned_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Constraint: max 3 active profession per character (application check).
Index: `ix_profession__character_id`.

---

## 11. `character.recipe_known`
Öğrenilmiş recipe'ler.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `recipe_id` | `UUID` | PK-part, FK → content.recipe |
| `learned_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

---

## 12. `character.mastery`
Mastery track ilerlemesi.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `track_id` | `UUID` | PK-part, FK → content.mastery_track |
| `level` | `SMALLINT` | NOT NULL DEFAULT 0 |
| `xp` | `BIGINT` | NOT NULL DEFAULT 0 |
| `claimed_milestones` | `SMALLINT[]` | NOT NULL DEFAULT `'{}'` |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

---

## 13. `character.ascension_choice`
Ascension tree'de harcanmış puanlar.

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `node_id` | `UUID` | PK-part, FK → content.ascension_node |
| `points_spent` | `SMALLINT` | NOT NULL |

---

## 14. `character.wallet`
Currency bakiyeleri (denormalize snapshot; source of truth = economy.ledger_entry).

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `currency_code` | `VARCHAR(24)` | PK-part |
| `balance` | `BIGINT` | NOT NULL DEFAULT 0, `CHECK balance >= 0` |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Notlar:
- Ledger'da toplam ≠ wallet ise reconcile job alarm.
- `currency_code`: `gold`, `dungeon_token`, `event_token_harvest2026`, ...

---

## 15. `character.reputation` (Phase 02)
Faction reputation.

| Kolon | Tip |
|---|---|
| `character_id` | `UUID` PK-part |
| `faction_id` | `UUID` PK-part |
| `value` | `INT` NOT NULL DEFAULT 0 |
| `tier` | `SMALLINT` NOT NULL DEFAULT 0 |

---

## 16. `character.death_log`
Ölüm audit.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `character_id` | `UUID` FK |
| `died_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `cause_type` | `VARCHAR(32)` NOT NULL — `mob`, `boss`, `pvp` (Phase 02), `env` |
| `cause_ref_id` | `UUID` NULL |
| `xp_debt_incurred` | `BIGINT` NOT NULL |
| `durability_lost` | `INT` NOT NULL DEFAULT 0 |

Index: `ix_death_log__character_id__died_at`.

---

## 17. `character.title_owned`
Kazanılmış title'lar (achievement/kozmetik).

| Kolon | Tip |
|---|---|
| `character_id` | `UUID` PK-part |
| `title_id` | `UUID` PK-part, FK → content.title |
| `earned_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `is_active` | `BOOLEAN` NOT NULL DEFAULT false |

---

## 18. Kullanım Örnekleri

**Karakter sheet load:**
```sql
-- Ana profil + stats + equipment tek query (JOIN):
SELECT c.*, s.*
FROM character.character c
LEFT JOIN character.stats s ON s.character_id = c.id
WHERE c.id = $1 AND c.soft_deleted_at IS NULL;
-- Equipment ayrı:
SELECT slot, item_instance_id FROM character.equipment WHERE character_id = $1;
-- Inventory paged:
SELECT slot_index, item_instance_id, stack_count
FROM character.inventory_slot
WHERE character_id = $1
ORDER BY slot_index;
```

**AFK reveal öncesi state load:**
```sql
SELECT c.id, c.level, c.xp, c.current_activity_id, c.last_reveal_at, c.seed,
       s.derived_cache_json
FROM character.character c
JOIN character.stats s ON s.character_id = c.id
WHERE c.id = $1
FOR UPDATE;
```

---

## 19. İlişkiler
- Owner: `auth.account`.
- Aktif contract: `activity.contract`.
- Envanter/ekipman item'ları: `content.item_instance`.
- Skill/profession/recipe/mastery: `content.*_template`.
- Ölüm nedeni: `content.npc_template` (mob) veya `world_boss_spawn`.

---

## 20. Notlar
- Karakter silme: soft delete (`soft_deleted_at`); 7 gün sonra hard delete job (audit).
- Karakter rename: unique name enforce; ledger event.
- Ascension anında level 1'e döner ancak envanter/equipment korunur; `character.character` update transaction'ı `ascension_level` +1, `xp` = 0, `level` = 1.
