# Side (Side Quest)

## 1. Amaç
Bölgesel, opsiyonel quest'ler. Ana anlatı dışı; keşif + karakter zenginlik.

## 2. Karakteristik

- **Opsiyonel**: main flow'u etkilemez.
- **Region-bound**: belirli region NPC'lerde.
- **Non-sequential**: sıralı değil (some chains still exist).
- **Repeatable**: **hayır** (bazı side chain repeatable — bkz. daily).
- **Reward-moderate**: xp + gold + occasional rare item.

## 3. Content Format

```yaml
# content/quests/side_ashen_hills_missing_traveler.yaml
code: side_ashen_hills_missing_traveler
category: side
display_name_key: quest.side.missing_traveler.name

giver_npc_ref: ashfield_worried_wife
turn_in_npc_ref: ashfield_worried_wife
region_code: ashen_hills

prerequisites:
  min_level: 5

objectives:
  - kind: discover
    target_poi_code: ashen_hills_forest_camp
  - kind: dialog
    target_npc_ref: forest_camp_wanderer
  - kind: deliver
    target_item_code: wanderer_note
    target_npc_ref: ashfield_worried_wife

rewards:
  xp: 300
  gold_copper: 800
  faction_rep:
    faction_code: emerald_wardens
    delta: 50
  items:
    - item_code: leather_bracers
      qty: 1

time_limit_ms: null
repeatable: false
tags: [exploration, mystery]
```

## 4. Region Content Volume

Per region:
- 15-25 side quest (v1 baseline).
- Mix: kill/gather/deliver/discover/dialog.

## 5. Side Chain (Opsiyonel Mini-Story)

Bazı side quest'ler mini-chain:
- Quest A tamamlandığında NPC yeni quest B verir.
- 3-5 quest zincir.
- Chain end: unique item + achievement.

## 6. Endpoint

**GET `/v1/characters/{id}/quests/side`** — active + eligible side quests.

## 7. AFK Tick

Aktivite reveal → objective tick:
- Kill quest active → hunting kill counter.
- Gather quest → mining/fishing.
- Discover quest → travel arrival trigger.

Auto-progress kolaylığı.

## 8. Achievement

Side quest zinciri tamamlanınca achievement:
- "Explorer of Ashen Hills" (all side quests in region).
- "Master of Mysteries" (all mystery-tagged side quests).

## 9. Metrics
- `quest_side_accepted_total{quest}`.
- `quest_side_completed_total{quest}`.
- `quest_side_per_region_completed{region}` gauge.

## 10. Cross-Reference
- [Main](main.md), [Daily](daily.md), [Weekly](weekly.md)
- [game-design/quests.md](../game-design/quests.md)

## 11. Bekleyen Kararlar
- Side quest tracker UI density (Phase 02 quest tags filter).
- Faction-hostile side quest (opposing faction reward Phase 02).
- Side chain length final.
