# Continent

## 1. Amaç
Dünyanın en üst coğrafi bölünümü. Farklı biyom setleri, kültürler, ilerleme yolları. v1'de 2-3 continent (kontrol edilmiş kapsam).

## 2. Model
Kaynak: [database/world.md §1](../database/world.md).

```sql
world.continent
  id UUID PK
  code VARCHAR(32) UNIQUE
  display_name_key VARCHAR(64) NOT NULL
  art_ref VARCHAR(128) NULL
  order_index SMALLINT NOT NULL DEFAULT 0
  unlock_condition_json JSONB NOT NULL DEFAULT '{}'
  version INT NOT NULL DEFAULT 1
```

Ek alanlar (Phase 13):
```sql
ALTER TABLE world.continent
  ADD COLUMN lore_key VARCHAR(64) NULL,
  ADD COLUMN recommended_level_min SMALLINT NULL,
  ADD COLUMN recommended_level_max SMALLINT NULL,
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}',
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL;
```

## 3. Content Pack (YAML)

```yaml
# content/world/continents/emberreach.yaml
code: emberreach
display_name_key: continent.emberreach.name
lore_key: continent.emberreach.lore
art_ref: continents/emberreach/map.png
order_index: 1

recommended_level_min: 1
recommended_level_max: 50

unlock_condition:
  # v1 varsayılan: hepsi açık
  # Phase 02'de bazı continent'lar quest/level gated:
  # min_level: 40
  # quest_completed: journey_to_farshore
  always_unlocked: true

tags: [starter, temperate, human_dominant]
```

v1 için önerilen 2 continent (placeholder):
- `emberreach` — starter continent, tempered biyomlar, tüm race spawn.
- `frostspire` — mid-level, buzul + dağ, dwarf/orc-friendly.
- (Phase 02'de bir 3.: `verdant_expanse` — high-level tropikal.)

## 4. Unlock

- Continent-level gate:
  - `always_unlocked: true` → default aç.
  - `min_level: N` → karakter level shart.
  - `quest_completed: quest_code` → prerequisite.
  - `faction_reputation: {faction, tier}` — Phase 02.

Character creation'da starter continent (race'e göre) — `content.race.starter_location_city` içinden inferred.

## 5. Endpoint

**GET `/v1/content/world/continents`**
Public.
```json
{
  "content_version": "...",
  "continents": [
    {
      "id": "...",
      "code": "emberreach",
      "display_name_key": "continent.emberreach.name",
      "art_ref": "continents/emberreach/map.png",
      "order_index": 1,
      "unlocked_for_character": true,
      "regions_count": 5
    }
  ]
}
```

**GET `/v1/content/world/continents/{code}`**
```json
{
  "code": "emberreach",
  "regions": [
    { "code": "ashen_hills", "biome": "hills", "recommended_level_min": 1, "recommended_level_max": 10 },
    ...
  ]
}
```

`unlocked_for_character` boolean: character-context'inde bakılır (Auth ile session character); public generic endpoint her zaman `true` döner + client karakter-ile lookup ayrı.

Alternative:
- **GET `/v1/characters/{id}/world/continents`** — character-scoped unlock status.

## 6. Character Location Mapping

`character.character.current_location_id` → location → territory → region → continent.

Discovery zinciri:
- Yeni continent açıldığında → notification "New continent unlocked: Frostspire".
- Achievement track.

## 7. Deprecation

Content rework nadir (continent kaldırma imkansıza yakın).
Deprecation:
- `deprecated_at` set.
- Karakter safe city'ye teleport (v1 starter continent'a).
- Notification + free respec token (opsiyonel; extreme senaryo).

## 8. Metrics
- `continent_content_active_total` gauge.
- `continent_character_distribution{continent}` gauge (karakter counts).
- `continent_unlocked_avg_level{continent}` gauge.

## 9. Audit
- `content.continent_created/updated/deprecated`.
- `character.continent_unlocked` (character milestone).

## 10. Test Stratejisi
- Content loader: continent + regions ref valid.
- Unlock condition check matrix.
- Deprecation flow (teleport + notification).

## 11. Cross-Reference
- [World](world.md), [Region](region.md)
- [game-design/world.md](../game-design/world.md)
- [database/world.md §1](../database/world.md)

## 12. Bekleyen Kararlar
- Continent count v1 (2 önerisi; playtest tune).
- Cross-continent travel süresi (ship — Phase 02).
- Continent-özel event (Phase 02 seasonal).
- Discovery XP for continent-first visit (Phase 02).
