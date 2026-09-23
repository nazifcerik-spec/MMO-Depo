# Weekly (Weekly Quest)

## 1. Amaç
Haftalık büyük hedef — 5-10 quest per hafta; büyük ödül. AFK-oyuncular için "orta-vadeli plan".

## 2. Karakteristik

- **Haftalık reset**: Salı 06:00 UTC (DD-016).
- **Slot cap**: 5-10 weekly quest.
- **Repeatable**: evet (haftalık).
- **Cooldown**: 7 gün.
- **Reward-heavy**: item + gold + faction rep + kozmetik.

## 3. Content Format

```yaml
# content/quests/weekly_boss_hunt.yaml
code: weekly_boss_hunt
category: weekly
display_name_key: quest.weekly.boss_hunt.name

giver_npc_ref: capital_weekly_board
turn_in_npc_ref: capital_weekly_board

prerequisites:
  min_level: 20

objectives:
  - kind: kill_boss
    boss_code: ashen_lich_king
    count: 1
  - kind: kill_boss
    boss_code: forest_guardian
    count: 1
  - kind: kill_boss
    boss_code: fire_wraith
    count: 1

rewards:
  xp: 5000
  gold_copper: 15000
  items:
    - item_code: weekly_reward_chest
      qty: 1
  currency_delta:
    dungeon_token: 100
    weekly_seal: 10

repeatable: true
cooldown_ms: 604800000           # 7 day
time_limit_ms: 604800000
tags: [weekly, boss, capstone]
```

## 4. Weekly Board

Capital city weekly board:
- 5-10 weekly quest available.
- Accept per week.
- Turn-in end of week.

## 5. Reset Sync

Weekly reset Salı 06:00 UTC:
- Dungeon lockout reset paralel (bkz. DD-016).
- Weekly quest progress reset.
- Yeni pool selection.

## 6. Weekly Seal

`weekly_seal` currency-item:
- Weekly quest reward.
- Season-end vendor exclusive mount/pet.
- Seal accumulate 12 hafta → legendary reward.

## 7. Multi-Character Weekly

Aynı account multiple character weekly:
- Her karakter kendi weekly.
- Account-wide "weekly master" achievement (Phase 02).

## 8. Missed Weekly

Weekly 7 gün içinde tamamlanmazsa:
- Expire; no makeup.
- Slot cleared for next week.

## 9. Endpoint

**GET `/v1/characters/{id}/quests/weekly`**
```json
{
  "available_weekly_slots": 10,
  "weekly_quests": [
    {
      "code": "weekly_boss_hunt",
      "expires_at": "2026-09-15T06:00:00Z",
      "objectives_progress": [
        { "current": 2, "goal": 3 }
      ]
    }
  ],
  "next_reset_at": "2026-09-15T06:00:00Z"
}
```

## 10. AFK Tick

Weekly boss kill: dungeon aktivite tetiklemez unless aktif dungeon run manuel.
- Boss kill objective tick sırasında dungeon manuel run.
- Alternate: idle dungeon boss kill counted (v1 evet).

## 11. Metrics
- `quest_weekly_accepted_total{quest}`.
- `quest_weekly_completed_total{quest}`.
- `weekly_seal_gained_total`.

## 12. Cross-Reference
- [Daily](daily.md), [Guild](guild.md), [Event](event.md)
- [game-design/dungeons.md](../game-design/dungeons.md), [game-design/world-bosses.md](../game-design/world-bosses.md)

## 13. Bekleyen Kararlar
- Weekly seal season vendor content Phase 02.
- Cross-character progress (BoA quest) Phase 02.
- Weekly leaderboard (fastest completion) Phase 02.
