# Daily (Daily Quest)

## 1. Amaç
Günlük havuz — sabit sayıda quest, günlük reset. AFK-oyuncular için "günlük checklist"; ancak zorunlu değil.

## 2. Karakteristik

- **Rotasyonlu havuz**: günlük 3 quest sunulur (havuzdan random).
- **Daily reset**: her gün 00:00 UTC.
- **Repeatable**: evet (her gün fresh).
- **Cooldown**: 24 saat per quest.
- **Reward-focused**: gold + tokens + moderate XP.

## 3. Content Format

```yaml
# content/quests/daily_ashen_hills_wolf_hunt.yaml
code: daily_ashen_hills_wolf_hunt
category: daily
display_name_key: quest.daily.wolf_hunt.name

giver_npc_ref: ashfield_daily_board       # sanal npc; board interaction
turn_in_npc_ref: ashfield_daily_board

prerequisites:
  min_level: 10
  region_visited: ashen_hills

objectives:
  - kind: kill
    target_npc_family: beast
    count: 20

rewards:
  xp: 200
  gold_copper: 500
  currency_delta:
    daily_token: 5
  faction_rep:
    faction_code: emerald_wardens
    delta: 10

repeatable: true
cooldown_ms: 86400000            # 24 hour
time_limit_ms: 86400000          # complete within 24 hour or expire
tags: [daily, kill, ashen_hills]
```

## 4. Daily Selection

Server (worker/scheduler):
- Her karakter için günlük 3 daily quest slot.
- Havuzdan random selection (region-appropriate + level-appropriate).
- Slot dolduğunda karakter yeni daily kabul edemez (v1 3 slot cap).

Reset:
- 00:00 UTC → tüm karakter daily progress reset.
- Yeni daily selection pool refresh.

## 5. Daily Board Interaction

Karakter city board'a gelir:
- 3 available daily.
- Accept → progress tracking başlar.
- Complete → turn-in board.

Alternate v1: auto-accept opsiyonu (Phase 02).

## 6. Daily Token

`daily_token` currency-item:
- Vendor exclusive item için.
- Weekly accumulate → tier reward.

Bkz. [character/reputation.md](../character/reputation.md) faction rep interaction.

## 7. Faction Rep Farming

Daily quests faction rep grind için primary path:
- Emerald Wardens daily → rep grind.
- Grand milestone: Exalted → unique title + cosmetic.

## 8. Missed Daily Handling

Daily 24h içinde tamamlanmazsa:
- Expire.
- Progress reset.
- Karakter yeniden accept edebilir sonraki gün (yeni selection).

**Missed daily = missed reward**; makeup yok (design intent).

## 9. Endpoint

**GET `/v1/characters/{id}/quests/daily`**
```json
{
  "available_daily_slots": 3,
  "daily_quests": [
    {
      "code": "daily_ashen_hills_wolf_hunt",
      "expires_at": "2026-09-10T00:00:00Z",
      "objectives_progress": [
        { "current": 8, "goal": 20 }
      ]
    }
  ],
  "cooldown_until_next_daily_reset_sec": 43200
}
```

## 10. AFK Tick

Aktivite reveal içinde daily objective tick:
- Hunting kill → beast family counter.
- Auto-progression AFK-friendly.

## 11. Metrics
- `quest_daily_accepted_total{quest}`.
- `quest_daily_completed_total{quest}`.
- `quest_daily_expired_total`.
- `quest_daily_slots_used_avg` gauge.

## 12. Cross-Reference
- [Weekly](weekly.md), [Side](side.md), [Event](event.md)
- [game-design/quests.md](../game-design/quests.md)

## 13. Bekleyen Kararlar
- Auto-accept daily Phase 02.
- Slot expansion (Phase 02 monetization opt cosmetic).
- Daily quest tier (bronze/silver/gold reward tier) Phase 02.
- Cross-region daily (roam bonus) Phase 02.
