# Location (POI Anchor)

## 1. Amaç
Aktivite anchor'ı ve map üstündeki genel POI. Resource node, monster area, dungeon area, waygate — hepsi location'a bağlıdır.

## 2. Model
Kaynak: [database/world.md §5](../database/world.md).

```sql
world.location
  id, code
  parent_id, parent_kind (region/territory/settlement)
  display_name_key
  coords POINT NULL
  tags VARCHAR(24)[]
  activity_hooks JSONB
  unlock_condition_json
```

`activity_hooks JSONB`:
```json
[
  { "kind": "resource_node", "template_code": "iron_valley_mining_node" },
  { "kind": "monster_area", "spawn_code": "whispering_forest_spawn" },
  { "kind": "dungeon_area", "template_code": "ashen_keep" },
  { "kind": "waygate", "code": "ashen_waygate" }
]
```

## 3. Content Pack

```yaml
# content/world/locations/iron_valley.yaml
code: iron_valley
parent_code: ashen_flats           # territory
parent_kind: territory
display_name_key: location.iron_valley.name
coords:
  x: 500
  y: 300
tags: [mining_spot, ore_rich, industrial]

activity_hooks:
  - kind: resource_node
    template_code: iron_ore_node
    tier: 2
    count: 6                        # kaç node instance spawn
    respawn_ms: 60000
  - kind: monster_area
    spawn_code: iron_valley_wolves
    max_active_mobs: 15
    respawn_ms: 5000

unlock_condition:
  always_unlocked: true
  # Alternatif:
  # min_level: 15
  # quest_completed: exploration_ashen_flats
```

## 4. Location Kategorileri (Tags)

- `hunt_spot` — monster area
- `gather_hub` — resource nodes cluster
- `dungeon_entry` — dungeon area
- `waygate` — waypoint
- `quest_hub` — quest giver POI
- `crafting_station` — outpost içerir
- `discovery_bonus` — first discovery XP
- `hidden` — fog of war (Phase 02)
- `event_spawn` — event trigger POI
- `boss_lair` — world boss spawn

Multiple tag mümkün (bir spot hem hunt hem gather olabilir).

## 5. Parent Polymorphism

Location'ın parent'ı:
- `region`: geniş wilderness.
- `territory`: territory içi.
- `settlement`: city/village/outpost içi.

`parent_kind` enum ile ayırt. FK enforce edilemez (polymorphic).

## 6. Activity Hooks

Bir location birden fazla activity hook taşıyabilir. Server aktivite başlatırken:
```go
func canStartActivity(char Character, activityType ActivityType, locationID uuid.UUID) error {
    loc := loadLocation(locationID)
    for _, hook := range loc.ActivityHooks {
        if hookMatches(hook, activityType) {
            return nil
        }
    }
    return errs.ErrForbidden  // activity.location_mismatch
}
```

## 7. Endpoint

**GET `/v1/content/world/locations/{code}`**
```json
{
  "code": "iron_valley",
  "parent_code": "ashen_flats",
  "parent_kind": "territory",
  "display_name_key": "location.iron_valley.name",
  "coords": { "x": 500, "y": 300 },
  "tags": ["mining_spot", "ore_rich"],
  "activity_hooks": [ ... ],
  "unlocked_for_character": true
}
```

**GET `/v1/content/world/locations/{code}/state`**
Runtime state (resource node respawn timers, monster count):
```json
{
  "resource_nodes": [
    { "instance_id": "...", "template_code": "iron_ore_node", "next_respawn_at": null, "available": true }
  ],
  "monster_area": {
    "spawn_code": "iron_valley_wolves",
    "current_active": 12,
    "max": 15
  }
}
```

## 8. Character Location

Character `current_location_id` → location id.

Location transition:
- Travel activity complete → location değişir.
- Death respawn → home_city location.
- Enter dungeon instance → temp location (dungeon_area).

## 9. POI vs Location Fark

- **Location**: mekanik anchor (aktivite başlatabilir).
- **POI** (`world.poi`): map üstünde görsel ikon (bazı POI location'sız olabilir — visual only).

Genelde: her POI bir location'a bağlıdır, ancak location'a bağlı POI olmaya da bilir (sadece anchor).

## 10. Coord System (v1 vs Phase 02)

- **v1**: abstract; `coords` opsiyonel map UI için. Aktivite hook mekanik.
- **Phase 02**: positional gameplay (mob spawn radius, node cluster) — koordinat kritik.

## 11. Unlock

- `unlock_condition_json`:
  - `always_unlocked: true` (default).
  - `min_level: N`.
  - `quest_completed: quest_code`.
  - `previous_location_discovered: code`.

Karakter'in unlock state'i `world.character_discovery` üstünden takip edilir (Phase 02 fog of war).

## 12. Hidden Location (Phase 02)

`tags` içinde `hidden`:
- Fog of war ile map'te görünmez.
- Discovery event ile açılır.
- Discovery XP grant + achievement.

## 13. Cross-Reference
- [Region](region.md), [Territory](territory.md), [Resource Node](resource-node.md), [Monster Area](monster-area.md), [Dungeon Area](dungeon-area.md)
- [game-design/map.md](../game-design/map.md), [game-design/world.md](../game-design/world.md)
- [database/world.md §5-6](../database/world.md)
- [activity-engine/README.md](../activity-engine/README.md) — location-gated activity check

## 14. Metrics
- `location_active_characters{location}` gauge.
- `location_activity_started_total{location,activity}`.
- `location_first_discovery_total{location}`.

## 15. Audit
- `content.location_created/updated/deprecated`.
- `character.location_discovered` (Phase 02).

## 16. Bekleyen Kararlar
- Positional coord Phase 02 aktivasyon.
- Fog of war Phase 02.
- Multi-parent location (nested) Phase 03.
- Dynamic locations (event-spawned, time-limited) Phase 02.
