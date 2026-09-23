# Database — `world` Schema

## Kapsam
Coğrafi hiyerarşi: kıta → bölge → territory → şehir/köy/outpost → location → POI. Ayrıca resource node **instance**'ları (template `content.resource_node_template`), city services.

---

## 1. `world.continent`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(32)` | UNIQUE |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `art_ref` | `VARCHAR(128)` | NULL |
| `order_index` | `SMALLINT` | NOT NULL DEFAULT 0 |
| `unlock_condition_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `version` | `INT` | NOT NULL DEFAULT 1 |

---

## 2. `world.region`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(32)` | UNIQUE |
| `continent_id` | `UUID` | FK NOT NULL |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `biome` | `VARCHAR(32)` | NOT NULL — forest/desert/tundra/coast/... |
| `recommended_level_min` | `SMALLINT` | NULL |
| `recommended_level_max` | `SMALLINT` | NULL |
| `unlock_condition_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |
| `art_ref` | `VARCHAR(128)` | NULL |

Index: `ix_region__continent_id`.

---

## 3. `world.territory`

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `code` | `VARCHAR(32)` UNIQUE |
| `region_id` | `UUID` FK NOT NULL |
| `display_name_key` | `VARCHAR(64)` NOT NULL |
| `faction_id` | `UUID` FK → content.faction NULL |
| `controlling_guild_id` | `UUID` FK → social.guild NULL — Phase 02 territory war |
| `art_ref` | `VARCHAR(128)` NULL |

---

## 4. `world.city` / `world.village` / `world.outpost`

Ortak `world.settlement` tablosu + kind ayrımı.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(48)` | UNIQUE |
| `territory_id` | `UUID` | FK NOT NULL |
| `kind` | `world.settlement_kind` | NOT NULL |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `size` | `world.settlement_size` | NOT NULL |
| `spawn_point_json` | `JSONB` | NOT NULL — {x,y,heading} |
| `services_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — bank, market, trainer, vendor, guild_hall bool flags |
| `ambient_theme` | `VARCHAR(32)` | NULL |
| `home_faction_id` | `UUID` | FK NULL |
| `art_ref` | `VARCHAR(128)` | NULL |

Enums:
```sql
CREATE TYPE world.settlement_kind AS ENUM ('city','village','outpost');
CREATE TYPE world.settlement_size AS ENUM ('capital','regional','village','outpost');
```

Indexler: `ix_settlement__territory_id`, `ix_settlement__kind`.

---

## 5. `world.location`

Genel POI + activity anchor.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `code` | `VARCHAR(48)` | UNIQUE |
| `parent_id` | `UUID` | NOT NULL — territory/settlement/region id (polymorphic) |
| `parent_kind` | `world.location_parent_kind` | NOT NULL |
| `display_name_key` | `VARCHAR(64)` | NOT NULL |
| `coords` | `POINT` | NULL — Postgres geometric type; opsiyonel v1 |
| `tags` | `VARCHAR(24)[]` | NOT NULL DEFAULT `'{}'` — hunt_spot, gather_hub, dungeon_entry, quest_hub |
| `activity_hooks` | `JSONB` | NOT NULL DEFAULT `'[]'` — [{kind, template_id}] |
| `unlock_condition_json` | `JSONB` | NOT NULL DEFAULT `'{}'` |

Enum:
```sql
CREATE TYPE world.location_parent_kind AS ENUM ('region','territory','settlement');
```

Indexler: `ix_location__parent_id__parent_kind`.

---

## 6. `world.poi`

Map üzerinde gösterilen ilgi noktası.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `location_id` | `UUID` FK NOT NULL |
| `category` | `VARCHAR(32)` NOT NULL — resource, vendor, dungeon, quest_hub, guild_hall, event |
| `display_name_key` | `VARCHAR(64)` NOT NULL |
| `icon_ref` | `VARCHAR(64)` NULL |
| `hidden_by_default` | `BOOLEAN` NOT NULL DEFAULT false — discovery gate |

---

## 7. `world.resource_node`

**Instance** — dünyada spawn'lanmış gerçek node. Template `content.resource_node_template`.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `template_id` | `UUID` FK → content.resource_node_template NOT NULL |
| `location_id` | `UUID` FK NOT NULL |
| `next_respawn_at` | `TIMESTAMPTZ` NULL — depleted ise ne zaman geri gelir |
| `spawn_group_id` | `UUID` NULL — grup halinde spawn kontrolü |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |

Notlar:
- Personal instance (DD-022): row bir "referans"; gerçek "sana özgü" node state Redis cache + activity_state içindedir.
- Rare node: server-wide announcement.

Indexler: `ix_resource_node__location_id`, `ix_resource_node__next_respawn_at`.

---

## 8. `world.waygate` (Teleport düğümü)

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `code` | `VARCHAR(48)` UNIQUE |
| `location_id` | `UUID` FK NOT NULL |
| `display_name_key` | `VARCHAR(64)` NOT NULL |
| `cost_gold_copper` | `BIGINT` NOT NULL |
| `unlock_condition_json` | `JSONB` NOT NULL DEFAULT `'{}'` |

---

## 9. `world.travel_route` (Wagon/Caravan/Ship)

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `code` | `VARCHAR(48)` UNIQUE |
| `origin_location_id` | `UUID` FK NOT NULL |
| `destination_location_id` | `UUID` FK NOT NULL |
| `method` | `world.travel_method` NOT NULL |
| `duration_ms` | `INT` NOT NULL |
| `cost_gold_copper` | `BIGINT` NOT NULL |
| `schedule_json` | `JSONB` NULL — recurring schedule for ships |

Enum:
```sql
CREATE TYPE world.travel_method AS ENUM ('walking','mount','road','wagon','ship','waygate');
```

---

## 10. `world.character_waypoint`

Kişisel waypoint'ler.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `character_id` | `UUID` FK NOT NULL |
| `name` | `VARCHAR(48)` NOT NULL |
| `location_id` | `UUID` FK NOT NULL |
| `is_shared_with_guild` | `BOOLEAN` NOT NULL DEFAULT false |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |

Constraint: max 20 per character (application; v2 expand).

---

## 11. `world.character_discovery`

Fog of war — keşfedilmiş POI (Phase 02 fog aktifse).

| Kolon | Tip |
|---|---|
| `character_id` | `UUID` PK-part |
| `poi_id` | `UUID` PK-part |
| `discovered_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |

---

## 12. `world.city_vendor_slot`

Şehirde satış yapan NPC vendor tanımı (basit).

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `settlement_id` | `UUID` FK NOT NULL |
| `display_name_key` | `VARCHAR(64)` NOT NULL |
| `inventory_json` | `JSONB` NOT NULL — [{item_template_id, price_copper, stock: infinite/int}] |
| `buys_categories` | `content.item_category[]` NOT NULL DEFAULT '{}' |

---

## 13. `world.city_trainer_slot`

Class/spec/profession trainer NPC.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `settlement_id` | `UUID` FK NOT NULL |
| `kind` | `VARCHAR(24)` NOT NULL — class, spec, profession, recipe |
| `service_json` | `JSONB` NOT NULL |

---

## 14. `world.quest_hub` (opsiyonel; quest_giver referansı)

```
id, settlement_id, display_name_key, offered_quest_template_ids UUID[]
```

---

## 15. Kullanım Örnekleri

**Şehir servis listesi:**
```sql
SELECT services_json FROM world.settlement WHERE id = $1;
```

**Karakter için müsait aktivite (location bazlı):**
```sql
SELECT l.activity_hooks
FROM world.location l
WHERE l.id = $1;
```

**Rare node spawn arama:**
```sql
SELECT rn.*, rnt.tier, rnt.kind
FROM world.resource_node rn
JOIN content.resource_node_template rnt ON rnt.id = rn.template_id
WHERE rnt.tier >= 4
  AND (rn.next_respawn_at IS NULL OR rn.next_respawn_at < now());
```

---

## 16. Notlar
- **Coord system**: `POINT` opsiyonel; MMO'da tam positional gameplay yok (idle-first). v1'de UI için mock coords yeter.
- **Location vs POI**: location oyun mekaniği (activity anchor); POI harita göstergesi (birden fazla POI aynı location'ı işaret edebilir).
- **Waygate teleport maliyeti** economy sink (DD-014).
- **Territory ownership** (Phase 02) `controlling_guild_id` alanı hazır.

---

## 17. Bekleyen Kararlar
- Personal instance node modelleme: pure Redis mi, satır bazlı mı — v1 basit personal computed at query time.
- Weather/time-of-day (Phase 02).
- Hidden POI unlock modu (fog).
