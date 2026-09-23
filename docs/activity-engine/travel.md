# Travel Activity

## 1. Amaç
Konum değişimi aktivitesi. Character location A → location B; süre = mesafe / hız. Aktivite tamamlanınca character location güncellenir.

Yüksek seviye: [game-design/travel.md](../game-design/travel.md).

## 2. Params

```json
{
  "destination_location_id": "loc-emerald-coast",
  "method": "mount",                          // walking, mount, road, wagon, ship, waygate
  "route_json": [...],                         // opsiyonel; server compute
  "cost_paid_copper": 500                      // waygate/ship cost paid at start
}
```

## 3. Channel

Travel `Channel = "mount_travel"` — primary aktivite ile paralel çalışabilir (bir mining + bir travel aynı anda mümkün mü? → HAYIR; mining bir spot'a bağlı, travel spot'tan çıkar).

Aslında travel karakter'in `current_location_id`'sini değiştirdiği için primary aktivite pause olur (or restart at destination). Karar: **travel primary channel — diğer aktivite paralel değil**. Mount_travel channel adı yanıltıcı; v1'de basitleştirilir.

Not: `mount_travel` channel iki senaryo için hazır (Phase 02):
1. Karakter travel ediyor + pet aktivite ile gathering yapıyor.
2. Karakter travel'da + AFK dungeon queue bekliyor.

v1: **primary channel'da travel**.

## 4. Prerequisites

- Character at origin location.
- Destination location var + unlocked.
- Method available:
  - `walking`: her zaman.
  - `mount`: active mount var.
  - `wagon`/`ship`: schedule check + cost pay.
  - `waygate`: waygate at origin unlocked + cost pay.
- Sufficient funds (waygate/ship/wagon).

Fail:
- `travel.locked_destination`
- `travel.no_mount`
- `travel.insufficient_funds`
- `travel.route_unavailable`

## 5. Rate & Duration Compute

Duration formulü:
```go
distance := computeDistance(origin, destination)  // logical units (v1: fixed table)
baseSpeed := methodBaseSpeed(method)  // walking=1.0, mount tier scale, wagon fixed, ship fixed, waygate=∞

movementBuff := char.Derived.MovementSpeedPct / 100.0
effectiveSpeed := baseSpeed * (1.0 + movementBuff)

durationMs := int64(distance / effectiveSpeed * 1000)
```

Waygate: durationMs = 0 (instant) veya fixed cast time (5s).

## 6. Yield Compute (Analytical)

Travel yield genelde küçük (discovery XP, mini bonus):
```go
if elapsedMs >= contract.DurationMs {
    // Complete → arrival
    return Yield{
        XPCharacter: content.discoveryXPForRoute(origin, destination),
    }
}
// Still in progress → 0 yield
return Yield{}
```

Travel **incremental yield vermez**; sadece tamamlanınca discovery XP.

Alternatif: partial travel completed durum:
- Cancel edilirse → char position rollback to origin (v1 simple) veya midway waypoint (Phase 02).

## 7. Threshold Detection

Tek "threshold": duration complete.
```
ms_to_complete := contract.DurationMs - elapsedMsInContract
if ms_to_complete <= elapsedMs:
    return Threshold{Kind:"travel_complete", HitAtMs: contract.DurationMs}
```

Diğer threshold:
- Ambush event (Phase 02) — random combat interrupt.

## 8. Complete Behavior

Reveal veya `OnComplete`:
- `character.current_location_id = destination`.
- Auto-next activity (Phase 02 config): oyuncu "auto-fishing at destination" tanımlamışsa auto-start.
- Notification "You have arrived".

## 9. Interrupt / Cancel

- Manual cancel: char position rollback to origin (v1); partial yield 0.
- Ambush interrupt (Phase 02): combat başlar → travel pause; player resume post-combat.

## 10. Content

Route table:
```yaml
# content/travel_routes/ashfield_to_emerald_coast.yaml
code: ashfield_to_emerald_coast
origin_location_code: ashfield
destination_location_code: emerald_coast
method: wagon
duration_ms: 720000    # 12 min
cost_gold_copper: 500
schedule_json:
  interval_min: 15
```

Waygate:
```yaml
# content/waygates/ashfield_waygate.yaml
code: ashfield_waygate
location_code: ashfield
display_name_key: waygate.ashfield.name
cost_gold_copper: 250
unlock_condition:
  quest_completed: waygate_ashfield_unlock
```

## 11. Multi-Leg Route

Complex travel (ashfield → capital → emerald) multi-leg:
- v1: single-leg only.
- Phase 02: route pathfinding + intermediate stops.

## 12. Reveal Örnek

```
Travel: ashfield → emerald_coast, wagon, 12 min duration
Started 8 min ago; elapsed = 8 min → still in progress.
Yield: 0 (not arrived).
Contract remains active.

Alt scenario:
Started 15 min ago; elapsed > 12 min → threshold hit.
Yield: 50 discovery XP.
Contract completed.
char.current_location = emerald_coast.
Notification: "Arrived at Emerald Coast".
```

## 13. Metrics
- `activity_travel_started_total{method}`.
- `activity_travel_completed_total{method}`.
- `activity_travel_cancelled_total`.
- `activity_travel_duration_avg{method}` histogram.

## 14. Cross-Reference
- [Activity](activity.md), [Offline Progress](offline-progress.md)
- [game-design/travel.md](../game-design/travel.md), [game-design/mounts.md](../game-design/mounts.md)
- [game-design/world.md](../game-design/world.md), [game-design/map.md](../game-design/map.md)
- [database/world.md](../database/world.md)

## 15. Bekleyen Kararlar
- Multi-leg route Phase 02.
- Ambush event system Phase 02.
- Auto-next activity queue Phase 02.
- Partial cancel with midway position (Phase 02).
