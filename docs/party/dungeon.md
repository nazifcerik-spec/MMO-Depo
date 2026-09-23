# Party Dungeon

## 1. Amaç
Party ile dungeon — manuel run instance + idle-run per-character contract.

Foundation: [activity-engine/dungeon.md](../activity-engine/dungeon.md), [dungeon/README.md](../dungeon/README.md).

## 2. Manuel Dungeon

Instance-based (activity.instance):
- Party leader enter → instance create.
- Members join instance (party_id ref).
- All members combat session shared.
- Encounter progression together.

Bkz. [dungeon/dungeon.md](../dungeon/dungeon.md).

## 3. Idle-Run Party Dungeon

- Party members each start dungeon_idle contract.
- Same dungeon template, same difficulty.
- Party bonus applied.
- Loot per character.

v1'de her karakter bağımsız contract; Phase 02 shared party contract opsiyonu.

## 4. Party Composition Requirement

Bazı dungeon party size min:
- `content.dungeon_template.party_size_min` (v1 default 1; some dungeons 3).
- Underparty → 400 `dungeon.underparty` on enter.

## 5. Ready Check

Leader triggers ready check before dungeon enter:
- Members prompt.
- All ready → enter.

## 6. Party Loot Rules

Bkz. [loot.md](loot.md).

Default: personal loot (DD-012).
Boss loot: each character rolls own loot table.

Master loot / round-robin Phase 02 opt.

## 7. Wipe Recovery

Party wipe (all dead):
- Encounter reset (mobs full HP).
- Party respawn (graveyard).
- Manual re-enter dungeon.
- Boss checkpoint preserved (defeated boss stays down).

## 8. Weekly Lockout

Per character lockout counter (not per party):
- Character joins party with lockout ⇒ still gets reward reduced.
- Fresh character in party ⇒ full reward.

Bkz. [dungeon/loot.md §5](../dungeon/loot.md).

## 9. Member Leave Mid-Dungeon

- Party members can leave mid-dungeon.
- Remaining party continues.
- Left member: instance ejected + travel back.

Auto-refill (LFG replacement Phase 02).

## 10. Cross-Region Party

Members different region:
- Dungeon area location-bound.
- All members must travel to dungeon entrance.
- v1 manuel travel; Phase 02 auto-teleport prompt.

## 11. Metrics
- `party_dungeon_completed_total{dungeon,size}`.
- `party_dungeon_wipe_total`.
- `party_dungeon_time_avg_min` histogram.

## 12. Cross-Reference
- [Party](party.md), [Combat](combat.md), [Loot](loot.md)
- [dungeon/README.md](../dungeon/README.md), [dungeon/loot.md](../dungeon/loot.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md)

## 13. Bekleyen Kararlar
- Party auto-teleport Phase 02.
- Auto-refill LFG replacement Phase 02.
- Shared party dungeon contract idle mode Phase 02.
- Cross-realm party (Phase 03).
