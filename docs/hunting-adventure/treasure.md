# Treasure

## 1. Amaç
Hidden chest, world drop, cache spawn. Explore ödülü + AFK reveal sürprizi.

## 2. Treasure Türleri

| Kind | Örnek | Trigger |
|---|---|---|
| `hidden_chest` | Forest cache | Random encounter |
| `world_drop_item` | Rare item on ground | Ultra rare (hunting) |
| `spawn_chest` | Location chest | Time-based location |
| `guarded_chest` | Elite guards drop | Elite kill co-loot |
| `treasure_hunt` (Phase 02) | Map-based quest | Explicit search |

## 3. Model

```sql
CREATE TABLE content.treasure_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  kind VARCHAR(24) NOT NULL,
  location_ids UUID[] NOT NULL DEFAULT '{}',        -- spawn location'lar
  region_ids UUID[] NOT NULL DEFAULT '{}',
  loot_table_id UUID NOT NULL,
  spawn_chance_per_activity_hour NUMERIC(5,3) NOT NULL,
  cooldown_ms BIGINT NOT NULL DEFAULT 21600000,     -- 6h per character
  requires_lockpick BOOLEAN NOT NULL DEFAULT false,   -- Phase 02
  min_character_level SMALLINT NULL,
  version INT NOT NULL DEFAULT 1
);
```

## 4. Content Pack

```yaml
# content/treasures/hidden_forest_chest.yaml
code: hidden_forest_chest
kind: hidden_chest
region_ids: [ashen_hills, emerald_coast_forest]
location_tags: [hunt_spot, gather_hub]

loot_table_code: hidden_forest_chest_loot
spawn_chance_per_activity_hour: 0.2       # ~%20 per hour

cooldown_ms: 21600000                      # 6h per character
requires_lockpick: false                    # v1 kapalı; Phase 02 opt
min_character_level: 10
```

## 5. Spawn Trigger

Reveal içinde (hunting/adventure/gathering aktivite):
```go
for _, treasure := range treasuresFor(activity.LocationID, snapshot.RegionID) {
    // Character cooldown
    if character.LastTreasureAt[treasure.Code].Within(cooldown_ms): continue
    // Level check
    if snapshot.Level < treasure.MinLevel: continue

    expected := treasure.SpawnChancePerActivityHour * elapsedHours
    count := samplePoisson(expected, 1.0, seed.Sub("treasure-"+treasure.Code))
    for i := 0; i < count; i++ {
        drops := lootSystem.RollTable(treasure.LootTableID, seed.Sub(fmt.Sprintf("treasure-%d", i)))
        addToReveal(drops)
        character.LastTreasureAt[treasure.Code] = now
    }
}
```

## 6. Loot Composition

Treasure loot:
- Gold pile (small).
- Random material.
- Rare item (%5 chance).
- Currency-item (event/dungeon token, opsiyonel).

Example loot table:
```yaml
# content/loot_tables/hidden_forest_chest_loot.yaml
code: hidden_forest_chest_loot
source_type: chest
entries:
  - item_code: gold_pile_small
    weight: 100
    min_qty: 100
    max_qty: 500
  - item_code: iron_ore
    weight: 50
    min_qty: 5
    max_qty: 15
  - item_code: forest_essence
    weight: 20
    min_qty: 1
    max_qty: 3
  - item_code: rare_gem
    weight: 5
    min_qty: 1
    max_qty: 1
```

## 7. World Drop (Ultra Rare)

`kind: world_drop_item`:
- Chance per activity hour ~0.005 (%0.5).
- Unique item — usually legendary+.
- No cooldown per character; global scarcity via chance.
- Reveal payload highlight.

Örnek: "Ancient Mystery Coin" — rare currency for special vendor.

## 8. Guarded Chest

Elite kill sonrası opsiyonel chest:
- Elite mob defeated → chest spawn (%50 chance).
- Chest reveal içinde auto-open.
- Rare drop guaranteed.

Content link: elite template'ye chest ref.

## 9. Lockpick (Phase 02)

Bazı chest'ler lockpick required:
- Consumable "Lockpick" (rogue-friendly).
- Craftable (Phase 02 rogue trainer).
- Success chance based on DEX + lockpick tier.

v1 kapalı; tüm chest auto-open.

## 10. Treasure Hunt Quest (Phase 02)

- Player buys treasure map (vendor rare, drop).
- Map hints location + hint clues.
- Player travel + gather → chest.

v1'de yok; Phase 02 quest chain.

## 11. Reveal Payload

```json
{
  "notable_encounters": [
    {
      "kind": "treasure_chest",
      "chest_code": "hidden_forest_chest",
      "items": [
        { "item_code": "gold_pile_small", "quantity": 320 },
        { "item_code": "iron_ore", "quantity": 12 },
        { "item_code": "forest_essence", "quantity": 2 }
      ]
    }
  ]
}
```

## 12. Metrics
- `treasure_spawned_total{kind,region}`.
- `treasure_looted_total{item}`.
- `world_drop_awarded_total{item}`.

## 13. Audit
- `treasure.spawned` (severity=debug).
- `treasure.world_drop` (severity=notable — ultra rare).

## 14. Cross-Reference
- [Hunting](hunting.md), [Encounter](encounter.md), [Adventure](adventure.md)
- [systems/loot-system.md](../systems/loot-system.md)
- [world/location.md](../world/location.md)

## 15. Bekleyen Kararlar
- Lockpick Phase 02.
- Treasure hunt quest Phase 02.
- Guarded chest per elite % chance final (v1: 50%).
- World drop rate tune (0.005/hour v1; playtest).
