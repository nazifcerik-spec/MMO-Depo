# Database — `content` Schema (Templates)

## Kapsam
Oyunun **data-driven** kısmı: race, class, specialization, skill, item template, recipe, dungeon, adventure, npc, loot table, pet template, mount template, ascension tree, mastery track, title, faction. Ayrıca item **instance** (oyuncu sahipli) burada — çünkü sahiplenilmiş item template'e referans verir ve global item id space'i buradadır.

Bu schema **hot-reloadable**; content patch'te version bump ile yenilenir.

---

## 1. `content.race`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(32)` | UNIQUE, NOT NULL |
| `display_name_key` | `VARCHAR(64)` | NOT NULL — i18n key |
| `lore_key` | `VARCHAR(64)` | NOT NULL |
| `stat_modifiers` | `JSONB` | NOT NULL DEFAULT `'{}'` — `{str:1, dex:-1, ...}` |
| `passive_ability_id` | `UUID` | FK → content.skill_template NULL |
| `active_ability_id` | `UUID` | FK → content.skill_template NULL |
| `profession_affinity_id` | `UUID` | FK → content.profession NULL |
| `profession_affinity_bonus_pct` | `SMALLINT` | NOT NULL DEFAULT 0 |
| `unique_mechanic_code` | `VARCHAR(64)` | NULL — kod tarafında case |
| `art_ref` | `VARCHAR(128)` | NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Indexler: `uq_race__code`.

---

## 2. `content.class`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(32)` | UNIQUE, NOT NULL |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `resource_type` | `content.resource_type` | NOT NULL |
| `armor_affinity` | `content.armor_kind` | NOT NULL |
| `weapon_affinities` | `content.weapon_kind[]` | NOT NULL DEFAULT `'{}'` |
| `role_tags` | `VARCHAR(16)[]` | NOT NULL DEFAULT `'{}'` — tank/dps/support/hybrid |
| `base_stat_bias` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `starter_skill_ids` | `UUID[]` | NOT NULL DEFAULT `'{}'` |
| `art_ref` | `VARCHAR(128)` | NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enums:
```sql
CREATE TYPE content.resource_type AS ENUM ('rage','mana','energy','focus','nature');
CREATE TYPE content.armor_kind AS ENUM ('cloth','leather','mail','plate');
CREATE TYPE content.weapon_kind AS ENUM (
  'sword','axe','mace','dagger','staff','wand','bow','crossbow',
  'shield','offhand_focus','polearm','fist'
);
```

---

## 3. `content.specialization`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(32)` | UNIQUE, NOT NULL |
| `class_id` | `UUID` | FK → content.class NOT NULL |
| `tier` | `SMALLINT` | NOT NULL `CHECK tier IN (1,2)` |
| `parent_spec_id` | `UUID` | FK → content.specialization NULL — tier 1 için null |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `passive_ids` | `UUID[]` | NOT NULL DEFAULT `'{}'` |
| `active_skill_unlock_ids` | `UUID[]` | NOT NULL DEFAULT `'{}'` |
| `stat_bias` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `playstyle_tags` | `VARCHAR(24)[]` | NOT NULL DEFAULT `'{}'` |

Constraint: `CHECK (tier = 1 OR parent_spec_id IS NOT NULL)`.

---

## 4. `content.skill_template`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(48)` | UNIQUE, NOT NULL |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `type` | `content.skill_type` | NOT NULL |
| `class_id` | `UUID` | FK NULL |
| `spec_id` | `UUID` | FK NULL |
| `race_id` | `UUID` | FK NULL |
| `resource_type` | `content.resource_type` | NULL |
| `resource_cost` | `INT` | NOT NULL DEFAULT 0 |
| `cooldown_ms` | `INT` | NOT NULL DEFAULT 0 |
| `gcd_bypass` | `BOOLEAN` | NOT NULL DEFAULT false |
| `cast_time_ms` | `INT` | NOT NULL DEFAULT 0 |
| `range_kind` | `content.skill_range` | NOT NULL |
| `aoe_shape` | `content.aoe_shape` | NULL |
| `effects_json` | `JSONB` | NOT NULL DEFAULT `'[]'` — damage/heal/buff/debuff/dot/hot/shield/teleport |
| `scaling_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — {str_coef: 1.0, int_coef: 0} |
| `unlock_condition_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `max_rank` | `SMALLINT` | NOT NULL DEFAULT 1 |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enums:
```sql
CREATE TYPE content.skill_type AS ENUM ('active','passive','utility','racial','ultimate');
CREATE TYPE content.skill_range AS ENUM ('melee','short','mid','long','global');
CREATE TYPE content.aoe_shape AS ENUM ('cone','circle','line','chain');
```

---

## 5. `content.skill_rank` (opsiyonel; ranklar per-skill farklı ise ayrı tablo)

| Kolon | Tip |
|---|---|
| `skill_id` | `UUID` PK-part |
| `rank` | `SMALLINT` PK-part |
| `effects_json_override` | `JSONB` NOT NULL |
| `scaling_json_override` | `JSONB` NOT NULL DEFAULT `'{}'` |

---

## 6. `content.profession`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(32)` | UNIQUE |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `category` | `content.profession_category` | NOT NULL |
| `max_level` | `SMALLINT` | NOT NULL DEFAULT 100 |
| `xp_curve` | `JSONB` | NOT NULL |

Enum:
```sql
CREATE TYPE content.profession_category AS ENUM ('gathering','crafting','service');
```

---

## 7. `content.recipe`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(48)` | UNIQUE |
| `profession_id` | `UUID` | FK NOT NULL |
| `tier_required` | `SMALLINT` | NOT NULL |
| `inputs_json` | `JSONB` | NOT NULL — `[{item_template_id, qty}]` |
| `outputs_json` | `JSONB` | NOT NULL — `[{item_template_id, qty, chance}]` |
| `time_ms` | `INT` | NOT NULL |
| `xp_reward` | `INT` | NOT NULL |
| `station_required_code` | `VARCHAR(32)` | NULL |
| `learn_source` | `VARCHAR(32)` | NOT NULL DEFAULT `'trainer'` — trainer, drop, quest, event |
| `version` | `INT` | NOT NULL DEFAULT 1 |

---

## 8. `content.item_template`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(64)` | UNIQUE |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `description_key` | `VARCHAR(96)` | NULL |
| `category` | `content.item_category` | NOT NULL |
| `subcategory` | `VARCHAR(32)` | NULL |
| `rarity` | `content.rarity` | NOT NULL |
| `base_stats_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `affix_pool_ids` | `UUID[]` | NOT NULL DEFAULT `'{}'` — prefix/suffix pool ref |
| `socket_slots` | `SMALLINT` | NOT NULL DEFAULT 0 |
| `set_id` | `UUID` | FK → content.item_set NULL |
| `binding_type` | `content.binding_type` | NOT NULL DEFAULT `'no_bind'` |
| `requirements_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — {min_level, class_ids, spec_ids, race_ids, stat_mins} |
| `base_durability` | `INT` | NULL |
| `stack_max` | `INT` | NOT NULL DEFAULT 1 |
| `weight` | `INT` | NOT NULL DEFAULT 1 |
| `vendor_buy_price_copper` | `BIGINT` | NULL |
| `vendor_sell_price_copper` | `BIGINT` | NULL |
| `art_ref` | `VARCHAR(128)` | NULL |
| `unique_effect_id` | `UUID` | FK NULL |
| `equipment_slot` | `character.equipment_slot` | NULL — sadece equipment ise |
| `weapon_kind` | `content.weapon_kind` | NULL |
| `armor_kind` | `content.armor_kind` | NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Enums:
```sql
CREATE TYPE content.item_category AS ENUM (
  'equipment','consumable','material','quest','currency','cosmetic','utility','recipe'
);
CREATE TYPE content.rarity AS ENUM (
  'poor','common','uncommon','rare','epic','legendary','mythic'
);
CREATE TYPE content.binding_type AS ENUM ('no_bind','boe','bop','boa');
```

Indexler:
- `uq_item_template__code`.
- `ix_item_template__category__rarity` (search).

---

## 9. `content.item_instance`

Sahiplenilmiş item; her drop bir instance.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `template_id` | `UUID` | FK → content.item_template NOT NULL |
| `owner_kind` | `content.item_owner_kind` | NOT NULL |
| `owner_id` | `UUID` | NOT NULL — character/guild/mailbox/market/listing |
| `rolled_stats_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `affixes_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — `{prefix_id, suffix_id}` |
| `sockets_json` | `JSONB` | NOT NULL DEFAULT `'[]'` — `[{slot, gem_instance_id}]` |
| `enchant_id` | `UUID` | FK → content.enchant NULL |
| `durability_current` | `INT` | NULL |
| `bound_to_character_id` | `UUID` | FK → character.character NULL |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `seed` | `BYTEA` | NOT NULL — RNG replay için |
| `template_snapshot_hash` | `TEXT` | NOT NULL — template değişiminde migration |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enum:
```sql
CREATE TYPE content.item_owner_kind AS ENUM (
  'character_inventory','character_bank','character_equipment',
  'guild_bank','mailbox','market_listing','void'
);
```

Indexler:
- `ix_item_instance__owner_kind__owner_id`.
- `ix_item_instance__template_id`.
- `ix_item_instance__bound_to_character_id`.

Notlar:
- `owner_id` polimorf FK (Postgres FK ile enforce edilemez); application layer sorumlu.
- Void owner: silinmiş ancak referans tarih için tutulan; garbage collect kural.

---

## 10. `content.item_set` + `content.item_set_bonus`

```
item_set: id, code, display_name_key
item_set_bonus: set_id, pieces_required (2/4/6), effect_json
```

---

## 11. `content.affix_pool` + `content.affix`

Affix havuzu tanımları.

```
affix_pool: id, code, kind (prefix/suffix), entries JSONB [{affix_id, weight, magnitude_min, magnitude_max}]
affix: id, code, kind, display_name_key, effect_json
```

---

## 12. `content.enchant`

```
enchant: id, code, display_name_key, applicable_slots [], effect_json, tier
```

---

## 13. `content.gem_template` + `content.gem_instance`

Socket sistemi (Phase 02 detay). Şimdilik iskelet:
```
gem_template: id, code, rarity, effect_json
gem_instance: id, template_id, owner_id, ...
```

---

## 14. `content.dungeon_template`

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `code` | `VARCHAR(48)` UNIQUE |
| `display_name_key` | `VARCHAR(64)` NOT NULL |
| `region_id` | `UUID` FK → world.region |
| `min_level` | `SMALLINT` NOT NULL |
| `max_level` | `SMALLINT` NOT NULL |
| `party_size_min` | `SMALLINT` NOT NULL DEFAULT 1 |
| `party_size_max` | `SMALLINT` NOT NULL DEFAULT 5 |
| `weekly_lockout` | `BOOLEAN` NOT NULL DEFAULT false |
| `idle_run_supported` | `BOOLEAN` NOT NULL DEFAULT true |
| `reward_pool_id` | `UUID` FK → content.loot_table NULL |
| `version` | `INT` NOT NULL DEFAULT 1 |

---

## 15. `content.boss_encounter`

```
id, dungeon_template_id, order_index,
phase_definitions JSONB,
loot_table_id UUID FK,
ai_profile_id UUID
```

---

## 16. `content.adventure_template` + `content.adventure_stage`

```
adventure_template:
  id, code, region_id, tier, duration_bucket (short/mid/long),
  reward_pool_id, faction_id (nullable), narration_key

adventure_stage:
  id, adventure_id, order_index,
  check_type (combat/skill/resource/luck), requirement_json,
  success_effect_json, failure_effect_json, narration_key
```

---

## 17. `content.npc_template`

Mob ve boss (world boss ayrı tabloda).

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `code` | `VARCHAR(48)` UNIQUE |
| `display_name_key` | `VARCHAR(64)` NOT NULL |
| `family` | `VARCHAR(32)` NOT NULL — beast/humanoid/undead/elemental/... |
| `level` | `SMALLINT` NOT NULL |
| `hp_base` | `BIGINT` NOT NULL |
| `damage_base` | `INT` NOT NULL |
| `armor` | `INT` NOT NULL DEFAULT 0 |
| `magic_resist` | `INT` NOT NULL DEFAULT 0 |
| `abilities_json` | `JSONB` NOT NULL DEFAULT `'[]'` |
| `ai_profile_id` | `UUID` NULL |
| `loot_table_id` | `UUID` FK → content.loot_table NULL |
| `xp_reward` | `INT` NOT NULL |
| `elite_flag` | `BOOLEAN` NOT NULL DEFAULT false |

---

## 18. `content.world_boss_template`

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `code` | `VARCHAR(48)` UNIQUE |
| `display_name_key` | `VARCHAR(64)` |
| `region_id` | `UUID` FK |
| `tier` | `SMALLINT` NOT NULL |
| `hp_base` | `BIGINT` NOT NULL |
| `hp_scale_formula` | `TEXT` NOT NULL — DSL/expr |
| `damage_profile_json` | `JSONB` NOT NULL |
| `ability_ids` | `UUID[]` NOT NULL DEFAULT '{}' |
| `loot_table_id` | `UUID` FK |
| `spawn_windows_json` | `JSONB` NOT NULL — recurrence rules |

---

## 19. `content.loot_table` + `content.loot_entry`

```
loot_table: id, code, source_type (mob/boss/node/chest/quest), rarity_dist_json, meta_json
loot_entry: id, loot_table_id, item_template_id, weight, min_qty, max_qty, guaranteed (bool)
```

Index: `ix_loot_entry__loot_table_id`.

---

## 20. `content.resource_node_template`

```
id, code, kind (mining/woodcutting/foraging/fishing),
tier, respawn_ms, yield_curve_json,
loot_table_id (drop çeşitliliği için opsiyonel)
```

---

## 21. `content.pet_template`

```
id, code, display_name_key, species, type (combat/gather/utility/cosmetic),
base_stats_json, rarity,
active_ability_id (nullable), passive_ability_id (nullable),
art_ref
```

## 22. `content.mount_template`

```
id, code, display_name_key, type (ground/cargo/aquatic/winged/cosmetic),
speed_multiplier NUMERIC(4,2), cargo_bonus_slots SMALLINT DEFAULT 0,
rarity, unlock_condition_json, art_ref
```

---

## 23. `content.quest_template`

```
id, code, category (main/side/daily/weekly/repeatable/guild/world/achievement),
giver_npc_ref VARCHAR(48) NULL,
region_id UUID FK NULL,
prerequisites_json JSONB,
objectives_json JSONB,  -- [{type,target,count}]
rewards_json JSONB,
time_limit_ms BIGINT NULL,
repeatable BOOLEAN NOT NULL DEFAULT false,
cooldown_ms BIGINT NULL,
version INT
```

---

## 24. `content.achievement`

```
id, code, category, points SMALLINT, hidden BOOLEAN DEFAULT false,
criteria_json, reward_json
```

---

## 25. `content.title`

```
id, code, display_name_key, description_key, source (achievement/event/purchase)
```

---

## 26. `content.faction` (Phase 02)

```
id, code, display_name_key, tier_thresholds_json
```

---

## 27. `content.mastery_track`

```
id, code, category (weapon/combat/profession/region/bestiary/collection),
max_level SMALLINT, xp_curve JSONB, milestones_json
```

---

## 28. `content.ascension_node`

```
id, tier SMALLINT, category (combat/gather/social/utility),
cost SMALLINT, effect_json, prerequisite_ids UUID[]
```

---

## 29. `content.event_template`

```
id, code, name_key, category (daily/weekly/seasonal/emergent/guild),
start_at TIMESTAMPTZ NULL, end_at TIMESTAMPTZ NULL,
activity_hook_id UUID NULL,
reward_pool_id UUID FK NULL,
token_currency_code VARCHAR(24) NULL,
narration_key VARCHAR(64) NULL
```

---

## 30. `content.currency_definition`

```
code VARCHAR(24) PK,
display_name_key VARCHAR(64) NOT NULL,
kind (primary/token/mark/honor/season_specific),
tradable BOOLEAN NOT NULL DEFAULT false,
cap BIGINT NULL,
description_key VARCHAR(96) NULL
```

---

## 31. `content.content_version`

Content patch tracking.

```
id BIGSERIAL PK,
version_tag VARCHAR(32) NOT NULL UNIQUE,
applied_at TIMESTAMPTZ NOT NULL DEFAULT now(),
manifest_hash TEXT NOT NULL,
notes TEXT
```

---

## 32. Kullanım Örneği

**Item drop simülasyonu:**
```sql
SELECT le.item_template_id, le.weight, le.min_qty, le.max_qty, le.guaranteed
FROM content.loot_entry le
WHERE le.loot_table_id = $1;
-- Weighted pick server-side (rng seed).
```

**Karakter class + spec + skill listeleme:**
```sql
SELECT s.*
FROM content.skill_template s
WHERE s.class_id = $1
  AND (s.spec_id IS NULL OR s.spec_id IN ($2, $3))
  AND (s.unlock_condition_json->>'min_level')::int <= $4;
```

---

## 33. Notlar
- **content.* tabloları hot reload**: patch sırasında yeni row insert, eski row `deprecated_at` alanı ile korunur (backward compat instance için).
- **Deprecation**: instance oluşturulmuş eski template'ler snapshot hash ile korunur.
- **Search performansı**: item_template için trigram GIN index (Phase 02).
- **RNG determinism**: item_instance.seed replay için tutulur.
