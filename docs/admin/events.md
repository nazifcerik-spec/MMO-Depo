# Events (World Event Scheduler)

## 1. Amaç
World event (world boss, invasion, regional, seasonal, dynamic) admin scheduler + manual trigger + monitor.

## 2. Scheduler View

Sekmeler: **Upcoming** (next 7d) · **Active** · **Recent** (last 7d).

### Upcoming
```
Wed 12:00 UTC  World Boss: Ashen Titan  Blackspire Peak
Fri 20:00 UTC  Invasion: Goblin Horde   Glenmoor
Sat 18:00 UTC  Regional: Harvest Festival  Greenfield
```

### Active
```
World Boss: Ashen Titan
  Started 14:00, ends 15:00
  Attendance: 342 characters
  HP: 45% remaining
  [Details]
```

### Recent
Sonuç özet: attendance count, tier grant dağılımı, first-kill.

## 3. Manual Trigger

Admin bir event'i schedule dışında tetikleyebilir (game_designer role):
- Template seçim.
- Location/region seçim.
- Duration override.
- Reason.
- `POST /internal/ops/world-event/trigger`.

**Not:** Bu ledger etkileyici (reward grant) olduğu için admin karar + audit + rare use case (event testing veya server-wide compensation).

## 4. Seasonal Event Rollout

Seasonal event start/end date genellikle content pack ile deploy edilir; admin manual override:
- Extend end date (extension case).
- Suspend event (bug case).
- Reveal all event tokens (rare compensation).

## 5. Event Detail (Active)

- Attendance stream (yeni katılım son 5 dk).
- HP progress (world boss).
- Wave progress (invasion).
- Contribution top 20 (rare cosmetic leaderboard v1 opt).

## 6. Aksiyonlar

| Action | Permission | Step-up |
|---|---|---|
| `event.trigger` | game_designer+ | yes |
| `event.cancel_active` | admin | yes | attendance ödül minimum verilir |
| `event.extend` | game_designer+ | no |
| `event.seasonal_suspend` | admin | yes |

## 7. Endpoint

- `GET /internal/admin/world-event/scheduled`
- `GET /internal/admin/world-event/active`
- `GET /internal/admin/world-event/recent?…`
- `GET /internal/admin/world-event/{spawn_id}` compound
- `POST /internal/ops/world-event/trigger` (step-up)
- `POST /internal/ops/world-event/{spawn_id}/cancel`
- `POST /internal/ops/world-event/{spawn_id}/extend` `{ additional_seconds }`

## 8. Cross-Reference
- [world-events](../world-events/README.md)

## 9. Bekleyen Kararlar
- Seasonal event content pack pipeline Phase 02.
- Server-wide compensation event (rare mass reward) runbook Phase 02.
- A/B test event flavor Phase 03.
