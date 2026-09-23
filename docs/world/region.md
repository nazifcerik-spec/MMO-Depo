# Region

## 1. Amaç
Continent'ın alt bölümü. Biyom, kimlik, kaynak profili + seviye bandı. Content'in ana yerleştirme birimi.

## 2. Model
Kaynak: [database/world.md §2](../database/world.md).

```sql
world.region
  id, code, continent_id (FK)
  display_name_key, biome, art_ref
  recommended_level_min, recommended_level_max
  unlock_condition_json
```

Ek (Phase 13):
```sql
ALTER TABLE world.region
  ADD COLUMN lore_key VARCHAR(64) NULL,
  ADD COLUMN climate VARCHAR(24) NULL,          -- temperate/arid/frigid/tropical
  ADD COLUMN dominant_resource_kinds VARCHAR(24)[] NOT NULL DEFAULT '{}',
  ADD COLUMN dominant_mob_families VARCHAR(24)[] NOT NULL DEFAULT '{}',
  ADD COLUMN faction_control_id UUID NULL,      -- default faction; territory war Phase 02
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL,
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}';
```

## 3. Content Pack

```yaml
# content/world/regions/ashen_hills.yaml
code: ashen_hills
continent_code: emberreach
display_name_key: region.ashen_hills.name
lore_key: region.ashen_hills.lore
biome: hills
climate: temperate
art_ref: regions/ashen_hills/map.png

recommended_level_min: 5
recommended_level_max: 15

dominant_resource_kinds: [mining, foraging]
dominant_mob_families: [beast, humanoid]

unlock_condition:
  always_unlocked: true

tags: [starter, iron_ore, wolves]
```

## 4. Region Content İhtiyaç Listesi (v1)

Her region için minimum content:
- 1 city + 1-2 village + 0-2 outpost.
- 2+ monster area (hunt spot).
- 3+ resource node (varying kinds).
- 1-2 dungeon area.
- 2-5 adventure lane (Phase 02'de detay).
- 15-25 side quest.
- 1 world boss spawn point (opsiyonel).

## 5. Unlock

- v1 tüm region unlocked continent unlocked'sa.
- Phase 02'de bazı region gated:
  - `min_level`
  - `quest_completed`
  - `faction_reputation`
  - `previous_region_discovered`

## 6. Endpoint

**GET `/v1/content/world/regions/{code}`**
```json
{
  "code": "ashen_hills",
  "continent_code": "emberreach",
  "display_name_key": "region.ashen_hills.name",
  "biome": "hills",
  "recommended_level_min": 5,
  "recommended_level_max": 15,
  "territories": [
    { "code": "ashen_flats", "faction_code": null, "settlements_count": 3 }
  ],
  "monster_areas_count": 2,
  "resource_nodes_count": 4,
  "dungeon_areas_count": 1
}
```

## 7. Region Mastery

Bkz. [game-design/mastery.md](../game-design/mastery.md).

Karakter belirli region'da geçirilen zaman + kill sayısı ile "Region Mastery" ilerletir:
- Milestone: +%2 gather yield, +%1 XP gain in region, discovery bonus.
- Track: `content.mastery_track.category='region'`.

## 8. Deprecation

- `deprecated_at` set.
- Karakter safe location'a teleport.
- Content re-route (varsa migration).

## 9. Territory War (Phase 02+)

`faction_control_id` alanı runtime dynamic:
- Territory war sonucu değişir.
- Faction bonuses region-wide.
- Detay Phase 02 aayrı sistem.

## 10. Metrics
- `region_character_time_spent_avg{region}` gauge.
- `region_content_asset_count{region}` gauge (health check).
- `region_unlock_total{region}`.

## 11. Audit
- `content.region_created/updated/deprecated`.
- `character.region_discovered`.

## 12. Cross-Reference
- [Continent](continent.md), [Territory](territory.md), [City](city.md), [Location](location.md)
- [game-design/world.md](../game-design/world.md)
- [database/world.md §2](../database/world.md)
- [game-design/mastery.md](../game-design/mastery.md)

## 13. Bekleyen Kararlar
- Region weather (Phase 02 seasonal cycle).
- Region-wide event (Phase 02).
- Region reputation (faction control impact) Phase 02+.
