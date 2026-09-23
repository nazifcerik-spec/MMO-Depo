# Encounter (Dungeon Encounter Framework)

## 1. Amaç
Room içindeki fiili "olay" tanımı. Encounter = fight/trap/interact unit; room bir veya birden fazla encounter içerebilir (Phase 02).

v1'de: 1 room = 1 encounter. Basit.

## 2. Encounter Türleri

Room kind ile aynı katalog (bkz. [room.md §2](room.md)):
- `combat` (enemy_pack, elite, boss).
- `trap`.
- `treasure_interact`.
- `puzzle` (Phase 02).
- `narrative` (npc dialog; Phase 02).

## 3. Encounter Life-Cycle

```
1. Encounter start (room enter).
2. Encounter tick (combat 250ms; trap instant; treasure instant).
3. Encounter resolve:
   - Success → reward + advance.
   - Fail (combat wipe): restart room (mobs reset).
   - Trap fail: damage absorbed + advance.
4. Emit event 'dungeon.encounter_completed'.
```

## 4. Encounter Instance

`activity.instance.state` içinde encounter progress:
```
{
  "current_floor": 2,
  "current_room_index": 3,
  "current_encounter_state": {
    "kind": "combat",
    "mobs_alive": [id1, id2],
    "mobs_defeated": [id3],
    "started_at": "..."
  }
}
```

## 5. Combat Encounter

- Combat session start (session sticky to instance).
- Mob spawn per template.
- Party members join session.
- Combat loop tick (250ms; bkz. [combat/combat-loop.md](../combat/combat-loop.md)).
- All mobs defeated → encounter complete.

## 6. Trap Encounter

Bkz. [trap.md](trap.md).

- Trap trigger on room enter.
- Skill/luck check per character.
- Damage or CC on fail.
- Auto-advance to next room.

## 7. Treasure Interact

- Chest visible.
- Character interact → open.
- Loot roll.
- Advance.

Party: each member personal chest instance (v1); shared chest Phase 02 option.

## 8. Puzzle Encounter (Phase 02)

- Riddle/mechanic solve.
- Party cooperation.
- Skill/knowledge check.

v1'de yok.

## 9. Narrative Encounter (Phase 02)

- NPC dialog.
- Choice affects branch (Phase 02+).
- Reward: quest hint, faction rep.

## 10. AFK Encounter Compute

Reveal analytical:
```
for each floor progressed:
    for each room in floor:
        encounter_result = deterministic_resolve(room, snapshot, seed.sub(...))
        aggregate yields
```

Bounded loop (~50 room per dungeon); per-tick simulation YOK.

## 11. Wipe Handling

Combat wipe:
- Party dead → wipe.
- Instance state: current_room progress reset.
- Boss checkpoint (floor cleared) korunur.
- Party respawn → re-enter dungeon.

## 12. Cross-Reference
- [Room](room.md), [Enemy](enemy.md), [Elite](elite.md), [Boss](boss.md), [Trap](trap.md), [Treasure](treasure.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md), [combat/death-revive.md](../combat/death-revive.md)

## 13. Bekleyen Kararlar
- Multi-encounter room Phase 02.
- Puzzle content pack Phase 02.
- Narrative branching Phase 02.
- Timed encounter (must complete within N sec) Phase 02.
