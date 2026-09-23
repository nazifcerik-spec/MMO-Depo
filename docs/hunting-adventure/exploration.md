# Exploration

## 1. Amaç
Discovery, fog of war, region mastery, POI hunt. AFK ile ilerlemez ancak session'da odaklanan explore aktivitesi.

Foundation: [world/location.md](../world/location.md), [world/map.md](../world/map.md), [character/mastery.md](../character/mastery.md).

## 2. Discovery

Discovery = karakter yeni POI/location keşfeder → discovery XP + achievement + region mastery XP.

### 2.1 v1 Behavior

- Karakter travel/hunt/gather ile location'a ilk kez varır.
- `world.character_discovery` row insert.
- Discovery XP grant (100-500 based on region tier).
- Notification "New location discovered: <name>".

### 2.2 Fog of War (Phase 02)

- POI `hidden_by_default = true`: karakter yakınlaşana kadar map'te görünmez.
- Discovery event → POI reveal.
- Client map animation.

v1'de tüm POI görünür ancak "not visited" işaretli.

## 3. Model

Kaynak: [database/world.md §11](../database/world.md).

```sql
world.character_discovery
  character_id, poi_id
  discovered_at
  PRIMARY KEY (character_id, poi_id)
```

Ek Phase 18:
```sql
CREATE TABLE character.exploration_progress (
  character_id UUID NOT NULL,
  region_id UUID NOT NULL,
  poi_discovered_count INT NOT NULL DEFAULT 0,
  poi_total_count INT NOT NULL,
  completion_pct SMALLINT NOT NULL DEFAULT 0,
  first_completed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (character_id, region_id)
);
```

## 4. Discovery Endpoint

**POST `/v1/characters/{id}/discover`** — implicit; travel arrival ile tetiklenir.

Alternate: **POST `/v1/characters/{id}/interact/{poi_id}`** — POI-based interaction discovery.

## 5. Region Completion

- `poi_discovered_count / poi_total_count = completion_pct`.
- %100 → region completion achievement + reward.

Reward:
- Region-specific title ("Explorer of Ashen Hills").
- Cosmetic (dye, tabard).
- Currency (event token).

## 6. Region Mastery Link

Bkz. [character/mastery.md](../character/mastery.md), [skill-engine/mastery.md](../skill-engine/mastery.md):
- Region mastery track: time spent in region + kills + discovery.
- Discovery events → mastery XP boost.

## 7. Hidden POI (Phase 02)

Special hidden POI'ler:
- Discovery quest.
- Random encounter via exploration.
- Rare treasure.

v1'de skip.

## 8. Cartography (Phase 02)

Cartographer profession (Phase 02):
- Map notes create (personal + guild shared).
- Sell to other players.
- Detailed sub-region maps unlock.

v1'de yok.

## 9. Waypoint Personal

Bkz. [database/world.md §10](../database/world.md):
- Character up to 20 personal waypoints.
- Custom name + location ref.
- Fast lookup from map UI.

## 10. Discovery XP Config

```yaml
exploration:
  discovery_xp_per_poi:
    tier_1: 100
    tier_2: 200
    tier_3: 300
    tier_4: 400
    tier_5: 500
  region_completion_bonus_xp: 5000
  region_completion_reward_token_code: exploration_token
  fog_of_war_enabled: false     # v1 kapalı
```

## 11. Cross-Reference
- [world/location.md](../world/location.md), [world/map.md](../world/map.md), [world/region.md](../world/region.md)
- [character/mastery.md](../character/mastery.md), [game-design/mastery.md](../game-design/mastery.md)
- [database/world.md §10-11](../database/world.md)
- [Adventure](adventure.md), [Encounter](encounter.md)

## 12. Metrics
- `poi_discovered_total{region}`.
- `region_completion_total{region}` counter.
- `waypoint_created_total`.

## 13. Audit
- `character.poi_discovered`.
- `character.region_completed`.

## 14. Bekleyen Kararlar
- Fog of war Phase 02 aktivasyon.
- Cartographer profession Phase 02.
- Region completion reward tier'lar (5000 XP + token — Phase 02 balance).
- Cross-region "grand exploration" (all continents completed → prestige) Phase 03.
