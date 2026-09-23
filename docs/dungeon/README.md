# 🔥 Dungeon Engine (Faz 19)

Dungeon = multi-floor, multi-room structure. Enemy, elite, boss, treasure, trap, random event içerir. **AFK çalışır** — DD-084 uyumlu analitik hesap.

Foundation:
- [combat/README.md](../combat/README.md) — combat engine.
- [activity-engine/dungeon.md](../activity-engine/dungeon.md) — activity contract (manuel + idle mode).
- [game-design/dungeons.md](../game-design/dungeons.md) — vizyon.
- [world/dungeon-area.md](../world/dungeon-area.md) — entrance location.

## ⚠️ AFK Prensibi: Analytical Per-Floor / Per-Room

**Per-room simulation YASAK** — bounded loop OK (floor sayısı < 20, room sayısı < 100 per dungeon).

Yaklaşım:
- Room = "encounter unit" — küçük N iteration; deterministic seed + snapshot.
- Boss = phase-based analytical (HP × DPS = clear time).
- Reveal: analytical clear progression O(floor × avg_room_count).

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Dungeon](dungeon.md) | Dungeon entity + template + instance |
| [Floor](floor.md) | Floor structure + progression |
| [Room](room.md) | Room mekaniği + kind çeşitliliği |
| [Encounter](encounter.md) | Encounter within room |
| [Enemy](enemy.md) | Dungeon mob standart |
| [Elite](elite.md) | Elite mob dungeon-özel |
| [Boss](boss.md) | Boss encounter + phase + enrage |
| [Treasure](treasure.md) | Dungeon chest + reward |
| [Trap](trap.md) | Trap mekaniği + damage/CC |
| [Random Event](random-event.md) | Dungeon-inject random event |
| [Loot](loot.md) | Loot distribution + party rule + lockout |
| [AFK Run](afk-run.md) | AFK-özel deep-dive + analitik yield |

## 2. Dungeon Yapısı

```
Dungeon
 └── Floor (1..N; default 3)
     └── Room (1..M per floor; default 5-8)
         ├── Enemy (standart mob pack)
         ├── Elite (opsiyonel per room)
         ├── Trap (opsiyonel)
         ├── Treasure (opsiyonel)
         └── Boss (floor end veya dungeon end)
```

Dungeon size örnekleri:
- **Small dungeon**: 1 floor, 5 room, 1 boss.
- **Standard**: 3 floor, 5-8 room per floor, 3 boss (1 per floor).
- **Mega dungeon** (Phase 02+): 10+ floor, 50+ room, 5+ boss.

## 3. İki Mod

- **Manuel Run**: `activity.instance`; party ile aktif oyun.
- **AFK Run** (idle): `activity.contract`; primary channel; ~%50 yield.

Bkz. [activity-engine/dungeon.md](../activity-engine/dungeon.md).

## 4. Ortak İlkeler

1. **AFK-first friendly** — idle mode content-parallel.
2. **Analytical** — DD-084 uyumlu; bounded loop OK.
3. **Deterministic RNG** — seed-based; reveal audit.
4. **Content-driven** — YAML dungeon layout.
5. **Party or solo** — 1-5 party size; solo scaled.
6. **Weekly lockout** — DD-016 Salı 06:00 UTC.

## 5. Ana Interface

```go
package dungeonengine

type Engine interface {
    // Instance yaratma (manuel)
    CreateInstance(ctx context.Context, req CreateInstanceRequest) (Instance, error)

    // Aktif progress (manuel)
    ProgressRoom(ctx context.Context, instance Instance, roomID uuid.UUID) (RoomResult, error)
    EncounterBoss(ctx context.Context, instance Instance, bossID uuid.UUID) (BossResult, error)

    // AFK contract yield compute
    ComputeAFKYield(ctx context.Context, contract Contract, snapshot Snapshot, elapsedMs int64) (AFKDungeonYield, error)

    // Lockout
    RecordCompletion(ctx context.Context, char Character, dungeonID uuid.UUID, difficulty Difficulty) error
}
```

## 6. Cross-Reference
- [game-design/dungeons.md](../game-design/dungeons.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md)
- [combat/README.md](../combat/README.md), [combat/pve-encounters.md](../combat/pve-encounters.md)
- [world/dungeon-area.md](../world/dungeon-area.md)
- [database/content-templates.md §14-15](../database/content-templates.md)
- [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 7. Ortak Bekleyen Kararlar

- Mega dungeon Phase 02+.
- Mythic+ affix Phase 02+.
- Cross-realm dungeon queue Phase 03.
- Dungeon reset scroll (weekly bypass) Phase 02 monetization.
- Positional room mechanics Phase 02.
