# Event (Event Quest)

## 1. Amaç
Server-wide/region event ile bağlı özel quest — event participation reward.

Foundation: [hunting-adventure/random-event.md](../hunting-adventure/random-event.md), [game-design/events.md](../game-design/events.md).

## 2. Karakteristik

- **Event window**: event start-end arasında aktif.
- **Participation-based**: event contribution ile progress.
- **Rewards event-özel**: event token → cosmetic vendor.
- **Time-limited**: event ends → quest expires.
- **Broadcast**: event notification'ta quest hint.

## 3. Content Format

```yaml
# content/quests/event_blood_moon_kill.yaml
code: event_blood_moon_kill
category: event
event_template_code: blood_moon

display_name_key: quest.event.blood_moon_kill.name

giver_npc_ref: event_capital_board       # sanal event board
turn_in_npc_ref: event_capital_board

prerequisites:
  event_active: blood_moon              # event running during quest

objectives:
  - kind: kill
    target_npc_family: undead            # blood moon → undead spawn boost
    count: 100

rewards:
  xp: 3000
  gold_copper: 5000
  currency_delta:
    blood_moon_token: 20

time_limit_ms: 86400000              # matches event duration (24 hour)
repeatable: false                     # once per event
tags: [event, blood_moon, seasonal]
```

## 4. Event Trigger

Random event start (bkz. [hunting-adventure/random-event.md](../hunting-adventure/random-event.md)):
- Event quest'ler auto-eligible.
- Karakter event board'da accept.
- Notification "New event quest: <name>".

## 5. Event Token

Event-özel currency-item:
- `blood_moon_token`, `harvest_token`, ... .
- Event vendor exclusive item.
- Cosmetic mount/pet/transmog.

Token expires? Karar:
- **v1**: token expire yok; permanent currency.
- **Phase 02**: opsiyonel expiring event token.

## 6. Multi-Event Overlap

Event overlap (blood moon + weekly world event):
- Multiple event quest concurrently.
- Contribution paralel.

## 7. Event Missed

Event ended, quest incomplete:
- Auto-expire.
- Progress lost.
- Season-end make-up rare (Phase 02).

## 8. Event Chain (Seasonal)

Season event chain:
- Harvest Festival: 10 daily-like quest tamamlanınca title unlock.
- Chain achievement track per season.

## 9. Endpoint

**GET `/v1/characters/{id}/quests/event`**
```json
{
  "active_events": [
    { "event_code": "blood_moon", "ends_at": "..." }
  ],
  "event_quests": [
    {
      "code": "event_blood_moon_kill",
      "objectives_progress": [
        { "current": 35, "goal": 100 }
      ]
    }
  ]
}
```

## 10. AFK Tick

Event active + aktivite yield:
- Kill events → objective tick.
- Reveal aggregate event participation.

## 11. Metrics
- `quest_event_accepted_total{event,quest}`.
- `quest_event_completed_total{event}`.
- `event_token_gained_total{token}`.

## 12. Cross-Reference
- [hunting-adventure/random-event.md](../hunting-adventure/random-event.md), [game-design/events.md](../game-design/events.md)
- [database/progression-quests-achievements.md §5](../database/progression-quests-achievements.md)

## 13. Bekleyen Kararlar
- Token expire Phase 02 (event-specific).
- Season chain achievement content pack.
- Cross-character event progress BoA Phase 02.
- Event leaderboard Phase 02.
