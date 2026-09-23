# Floor

## 1. Amaç
Dungeon'un ana bölünüm birimi. Her floor N room + 1 boss (genelde) içerir. Floor complete → sonraki unlock.

## 2. Floor Model

Content pack içinde `dungeon.layout.floors[]`:
```yaml
floors:
  - floor: 1
    rooms: [...]
    boss_code: keep_guardian
    unlock_condition: none
  - floor: 2
    rooms: [...]
    boss_code: fire_wraith
    unlock_condition: floor_1_boss_defeated
```

Runtime state `activity.instance.current_floor`.

## 3. Floor Progression

```
1. Enter floor N.
2. Progress rooms sequentially (manuel) veya analytical clear (idle).
3. Boss encounter at floor end.
4. Boss defeated → floor complete.
5. Advance to floor N+1 (unlock condition satisfied).
```

## 4. Floor Difficulty Scaling

Deeper floor → harder:
- Mob HP + damage +%20 per floor.
- Loot tier +1 per floor.
- Room count more (v1: sabit; Phase 02 scaling).

## 5. Skip Ahead (Party Wipe Recovery)

- Party wipe → floor reset (mobs full HP); floor restart.
- Boss defeated already → boss stays defeated; sadece mobs restart.

## 6. Checkpoint (Save State)

- Boss defeated = checkpoint (persistent).
- Client leave & rejoin → resume at latest floor.
- Instance expire (2h) → progress lost.

## 7. Manuel Floor Progress

Client → `POST /v1/dungeons/instances/{id}/room/{room_id}/enter`:
- Server validate room in current floor.
- Room encounter trigger (combat/trap/treasure).
- Complete → next room OR floor boss.

## 8. Idle Floor Progress

Reveal içinde:
```
elapsed_ms
current_floor = character.instance.current_floor
per_floor_clear_time_ms = dungeon.estimated_clear_time / floor_count
floors_progressed = elapsed_ms / per_floor_clear_time_ms
new_floor = min(dungeon.floor_count, current_floor + floors_progressed)
```

Yield: floors_progressed × per_floor_reward.

## 9. Boss per Floor

Her floor genelde 1 boss (final). Alt-boss (mini-boss) opsiyonel:
```yaml
floor: 2
rooms:
  - kind: enemy_pack
  - kind: mini_boss     # mini-boss encounter mid-floor
    boss_code: skeletal_champion
  - kind: enemy_pack
  - kind: boss
    boss_code: fire_wraith
```

## 10. Metrics
- `dungeon_floor_progressed_total{dungeon,floor}`.
- `dungeon_floor_clear_time_avg{dungeon,floor}` histogram.

## 11. Cross-Reference
- [Dungeon](dungeon.md), [Room](room.md), [Boss](boss.md), [AFK Run](afk-run.md)

## 12. Bekleyen Kararlar
- Mid-floor mini-boss Phase 02 detay.
- Floor difficulty scaling per-floor factor.
- Checkpoint persistence duration (v1 2h expire).
- Descending vs ascending floor (dungeon layout).
