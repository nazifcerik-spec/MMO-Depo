# Main (Story Quest)

## 1. Amaç
Ana anlatı zinciri. Sıralı, prerequisite-based, karakter dünyaya tanıtan story flow.

## 2. Karakteristik

- **Sıralı**: quest N tamamlandığında quest N+1 unlock.
- **Prerequisite**: level + quest completion.
- **Repeatable**: **hayır**.
- **Reward-heavy**: XP + gold + set item + title + cosmetic.
- **Location-hopping**: region tour.
- **Cutscene / dialog**: rich narrative (Phase 02 UI).

## 3. Content Format

```yaml
# content/quests/main_ashen_hills_intro.yaml
code: main_ashen_hills_intro
category: main
display_name_key: quest.main.ashen_hills_intro.name
description_key: quest.main.ashen_hills_intro.description
lore_key: quest.main.ashen_hills_intro.lore

giver_npc_ref: ashfield_captain_elara
turn_in_npc_ref: ashfield_captain_elara

prerequisites:
  min_level: 3
  # quest_completed: main_prologue   (varsa)

objectives:
  - kind: kill
    target_npc_code: shadow_wolf
    count: 5
  - kind: gather
    target_item_code: shadow_pelt
    count: 3
  - kind: dialog
    target_npc_ref: ashfield_captain_elara

rewards:
  xp: 500
  gold_copper: 2500
  items:
    - item_code: leather_boots_starter
      qty: 1
  next_quest_code: main_ashen_hills_investigation

time_limit_ms: null              # main quest genelde time limit yok
repeatable: false
```

## 4. Quest Chain

Main quest chain graph:
- Prologue → Ashen Hills Intro → Investigation → Confrontation → Resolution.
- Her tamamlandığında `next_quest_code` auto-unlock.

Server:
- Quest turn-in → next quest_code unlock (character.quest_progress).
- Notification "New main quest available: <name>".

## 5. Reward Composition

Main quest rewards zengin:
- XP: level level_xp_required(min_level) × 0.3-0.5 baseline.
- Gold: 500-5000 copper.
- Item: starter gear, tier-appropriate weapon/armor.
- Milestone chain end: title + cosmetic mount/pet unlock.

## 6. Prerequisites

- Level.
- Previous quest completion.
- Item (rare — quest item collection).
- Faction reputation tier (Phase 02).

## 7. Chain Break Recovery

Content patch veya bug ile chain broken:
- Admin support ticket.
- Auto-repair worker: verify chain sequence intact.

## 8. Storyline Skip (Alt Path)

Bazı chain'lerde alternate branch:
- v1'de single path.
- Phase 02: choice-based branch.

## 9. Endpoint

**GET `/v1/characters/{id}/quests/main`** — main quest tracker.

```json
{
  "current_main_quest": {
    "code": "main_ashen_hills_intro",
    "objectives_progress": [
      { "kind": "kill", "target": "shadow_wolf", "current": 3, "goal": 5 }
    ]
  },
  "completed_main_count": 12,
  "next_available": null
}
```

## 10. AFK Objective Tick

Aktivite reveal içinde main quest objective auto-tick:
- kill objective: hunting kill count → objective_progress.
- gather objective: mining/fishing yield.

Detay: [systems/reward-system.md](../systems/reward-system.md) benzer flow.

## 11. Metrics
- `quest_main_accepted_total{quest}`.
- `quest_main_completed_total{quest}`.
- `quest_main_chain_stuck_total`.

## 12. Cross-Reference
- [Side](side.md), [Daily](daily.md)
- [database/progression-quests-achievements.md §1](../database/progression-quests-achievements.md)
- [game-design/quests.md](../game-design/quests.md)

## 13. Bekleyen Kararlar
- Branching narrative Phase 02.
- Voice-over Phase 03.
- Story cutscene UI Phase 02.
- Chain length per region (v1: ~10-15 main quest per region).
