# Quest System (Faz 22)

Karakteri hedeflere bağlayan structured content. AFK-uyumlu (objective tick reveal), party-share (Phase 02), reward-rich.

Foundation:
- [game-design/quests.md](../game-design/quests.md) — vizyon.
- [database/content-templates.md §23](../database/content-templates.md), [database/progression-quests-achievements.md §1-2](../database/progression-quests-achievements.md) — schema.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Main](main.md) | Main story quest — sıralı ana anlatı |
| [Side](side.md) | Side quest — bölgesel opsiyonel |
| [Daily](daily.md) | Günlük havuz + rotation |
| [Weekly](weekly.md) | Haftalık büyük hedef |
| [Dynamic](dynamic.md) | Procedural / random-generated |
| [Event](event.md) | Event-tied quest |
| [Guild](guild.md) | Guild-shared quest (Phase 02) |

## 2. Ortak Quest Model

Kaynak: [database/content-templates.md §23](../database/content-templates.md).

`content.quest_template`:
- `code`, `category` (main/side/daily/weekly/repeatable/guild/world/achievement).
- `giver_npc_ref`, `region_id`.
- `prerequisites_json` (level, quest, item, rep).
- `objectives_json` (`[{type, target, count}]`).
- `rewards_json` (xp, gold, items, rep).
- `time_limit_ms`, `repeatable`, `cooldown_ms`.

Character state: [database/progression-quests-achievements.md §1](../database/progression-quests-achievements.md) `progression.quest_progress`.

## 3. Ortak Sözleşme

1. **Accept**: player explicit acceptance (main/side); auto-accept opsiyonel daily.
2. **Objective tick**: aktivite/kill/gather/deliver — AFK reveal içinde de tick eder.
3. **Complete**: all objectives met → turn-in eligible.
4. **Turn-in**: NPC visit (v1); auto-turn-in Phase 02.
5. **Reward**: XP + gold + item + rep + title.

## 4. Objective Türleri

Bkz. [game-design/quests.md §5](../game-design/quests.md):
- `kill` — mob N sayı.
- `gather` — resource N.
- `deliver` — item to NPC.
- `discover` — POI visit.
- `craft` — recipe iteration.
- `use_skill` — skill cast N.
- `dialog` — NPC talk.
- `escort` (Phase 02) — NPC follow.

## 5. Ortak İlkeler

1. **AFK-integrated**: reveal içinde objective tick.
2. **Multi-quest concurrent**: aynı anda ~25 quest max (v1).
3. **Data-driven** — YAML content pack.
4. **Deterministic** — cross-session state consistent.

## 6. Endpoint Genel

```
GET  /v1/characters/{id}/quests                 # aktif + eligible
POST /v1/characters/{id}/quests/{code}/accept
POST /v1/characters/{id}/quests/{code}/abandon
POST /v1/characters/{id}/quests/{code}/turn-in
GET  /v1/content/quests                          # public discovery
GET  /v1/content/quests/{code}
```

## 7. Cross-Reference

- [game-design/quests.md](../game-design/quests.md)
- [database/content-templates.md §23](../database/content-templates.md), [database/progression-quests-achievements.md §1](../database/progression-quests-achievements.md)
- [systems/reward-system.md](../systems/reward-system.md)
- [activity-engine/README.md](../activity-engine/README.md) — objective tick

## 8. Ortak Bekleyen Kararlar

- Auto-accept daily Phase 02.
- Auto-turn-in Phase 02.
- Quest chain graph Phase 02.
- Quest share (party) Phase 02.
- Quest tracker UI advanced Phase 02.
