# Room

## 1. Amaç
Floor içindeki encounter birimi. Room kind = ne olduğu (enemy pack, elite, boss, trap, treasure, event).

## 2. Room Kinds

| Kind | İçerik |
|---|---|
| `enemy_pack` | 3-6 standart mob |
| `elite` | 1 elite mob veya 2-3 mob elite group |
| `mini_boss` | Elite boss-lite; mid-floor |
| `boss` | Floor final boss (phase system) |
| `treasure` | Chest room; combat opsiyonel |
| `trap` | Trap check; skill/luck bypass |
| `random_event` | Event pool inject |
| `puzzle` (Phase 02) | Puzzle check; skill/knowledge |
| `rest` (Phase 02) | Safe zone; heal + buff |

## 3. Room Content Pack (Layout)

```yaml
# content/dungeons/ashen_keep.yaml (excerpt)
floors:
  - floor: 1
    rooms:
      - kind: enemy_pack
        mob_pool: [skeleton_warrior, skeleton_archer]
        mob_count: 4
      - kind: enemy_pack
        mob_pool: [skeleton_warrior]
        mob_count: 3
      - kind: treasure
        treasure_code: ashen_keep_small_chest
      - kind: elite
        elite_code: skeleton_captain
      - kind: trap
        trap_code: floor_spikes
      - kind: enemy_pack
        mob_pool: [skeleton_archer]
        mob_count: 5
      - kind: boss
        boss_code: keep_guardian
```

## 4. Room Trigger

**Manuel path**:
1. Enter room (server checks current_floor/room_index).
2. Room kind dispatch:
   - `enemy_pack`: spawn mobs → combat session.
   - `elite`: spawn elite (harder combat).
   - `trap`: skill check (bkz. [trap.md](trap.md)); pass or damage.
   - `treasure`: chest interact (bkz. [treasure.md](treasure.md)).
   - `boss`: boss encounter (bkz. [boss.md](boss.md)).
   - `random_event`: event pool roll.
3. Room complete condition:
   - Combat: all mobs defeated.
   - Treasure: chest opened.
   - Trap: passed or damage absorbed.
   - Boss: boss defeated → floor complete.
4. Advance room_index.

**Idle path**:
- Analytical clear per room; total floor clear time / room count.
- Trap damage: analytical HP loss.
- Treasure: auto-open + loot.
- Random event: Poisson sample injection.

## 5. Room Time Estimate

Per room clear time:
```yaml
room_kind_time_ms:
  enemy_pack: 60000       # 1 min
  elite: 120000            # 2 min
  mini_boss: 240000        # 4 min
  boss: 300000             # 5 min
  treasure: 10000          # 10 sec
  trap: 5000               # 5 sec
  random_event: 30000      # 30 sec
```

Dungeon total time ≈ sum room times.

## 6. Room Reward

- `enemy_pack`: standart loot per mob.
- `elite`: elite loot table + rare drop chance.
- `mini_boss`: rare drop guaranteed + achievement.
- `boss`: boss loot table + first-clear bonus + set item drop.
- `treasure`: chest loot table.
- `trap`: no direct reward (avoid pattern) + rare "trap disarm bonus" (rogue skill Phase 02).
- `random_event`: event-specific.

Detay: [loot.md](loot.md).

## 7. Party Room Interaction

Party members enter room together (v1); split not allowed.
- Combat session shared.
- Loot: personal (DD-012).
- Treasure: personal chest instance per character.

## 8. Metrics
- `dungeon_room_completed_total{dungeon,floor,kind}`.
- `dungeon_room_wipe_total{kind}`.

## 9. Cross-Reference
- [Dungeon](dungeon.md), [Floor](floor.md), [Enemy](enemy.md), [Elite](elite.md), [Boss](boss.md), [Treasure](treasure.md), [Trap](trap.md), [Random Event](random-event.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)

## 10. Bekleyen Kararlar
- Puzzle room Phase 02 (skill-based check).
- Rest room Phase 02 (safe zone regen).
- Branching path (multiple room paths within floor) Phase 02.
- Optional room (skippable) Phase 02.
