# Treasure (Dungeon Chest)

## 1. Amaç
Dungeon room içinde chest — combat-free reward. Boss loot ile paralel; standart material + rare item.

Foundation: [hunting-adventure/treasure.md](../hunting-adventure/treasure.md).

## 2. Chest Kind

- `small_chest` — floor 1; common material + gold.
- `mid_chest` — floor 2; uncommon material + gold + rare drop chance.
- `large_chest` — floor 3 / hidden room; rare material + rare item.
- `boss_chest` — post-boss guaranteed chest.
- `locked_chest` (Phase 02) — lockpick required.
- `secret_chest` (Phase 02) — hidden room + puzzle.

## 3. Content Pack

```yaml
# content/treasures/ashen_keep_chest_small.yaml
code: ashen_keep_chest_small
kind: dungeon_chest
loot_table_code: ashen_keep_chest_small_loot
open_time_ms: 3000              # 3 sec to open
lockpick_required: false
guarantees:
  - item_code: dungeon_token
    qty: 5
    guaranteed: true
```

## 4. Chest Interaction

Manuel:
- Character interact → open animation (3 sec).
- Cannot be interrupted (v1); Phase 02: interrupt allowed.
- Loot popup.

Party:
- **Personal chest instance** per character (v1 tercih).
- Each character opens independently.
- No competition; fairness.

Alternative shared chest (Phase 02):
- One physical chest; first-come first-served.
- Round-robin.

## 5. AFK Chest

Reveal:
```
for each treasure room in cleared floors:
    treasure = load(room.treasure_code)
    for each character:
        loot = lootSystem.RollTable(treasure.LootTableID, seed.sub("treasure-<room_id>-<char>"))
        add to yield
```

Auto-open (no player action needed).

## 6. Weekly Lockout Interaction

Dungeon weekly lockout dolarsa treasure loot penalty:
- Full loot for first clear.
- Subsequent: gold only (rare item skip).

Weekly clear tracker character-özel.

## 7. Chest vs Boss Loot

- **Boss loot**: unique/set items; primary progression.
- **Treasure**: fillers (material, gold, tokens); progression secondary.

Together = complete dungeon reward pool.

## 8. Loot Config

```yaml
# content/loot_tables/ashen_keep_chest_small_loot.yaml
code: ashen_keep_chest_small_loot
source_type: chest
entries:
  - item_code: gold_pile_medium
    weight: 100
    min_qty: 200
    max_qty: 500
  - item_code: iron_ore
    weight: 50
    min_qty: 5
    max_qty: 12
  - item_code: dungeon_token
    weight: 100
    min_qty: 3
    max_qty: 5
    guaranteed: true
  - item_code: rough_ruby
    weight: 15
    min_qty: 1
    max_qty: 1
```

## 9. Metrics
- `dungeon_chest_opened_total{dungeon,chest_kind}`.
- `dungeon_treasure_gold_awarded_total`.
- `dungeon_treasure_rare_drop_total`.

## 10. Cross-Reference
- [Dungeon](dungeon.md), [Room](room.md), [Loot](loot.md)
- [hunting-adventure/treasure.md](../hunting-adventure/treasure.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 11. Bekleyen Kararlar
- Locked chest (lockpick) Phase 02.
- Secret chest hidden room Phase 02.
- Personal vs shared chest v1 personal.
- Chest interrupt (mid-open) Phase 02.
