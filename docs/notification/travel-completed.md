# Travel Completed Notification

## 1. Amaç
Travel bitip karakter hedefe vardığında push. Uzun travel (30dk+) için özellikle önemli.

## 2. Trigger

`event.travel.arrived`
```json
{
  "character_id": "...",
  "travel_id": "...",
  "from_location_id": "...",
  "to_location_id": "...",
  "duration_actual_seconds": 1830,
  "interrupted": false
}
```

## 3. Payload

```json
{
  "title_key": "notif.travel_completed.title",
  "body_key": "notif.travel_completed.body.arrived",
  "body_vars": {
    "location_name": "Ironhold",
    "duration_hm": "30 dk 30 sn"
  },
  "category": "activity",
  "route": "character/location/<to_location_id>",
  "badge_delta": 1,
  "collapse_key": "travel_<character_id>"
}
```

## 4. Interrupted Case

Travel interrupted (encounter, mount die, storm):
```json
"body_key": "notif.travel_completed.body.interrupted",
"body_vars": { "cause": "wolf_encounter", "at_location": "wildwood_border" }
```

Priority `high` — kullanıcının karar vermesi gerekir (resume vs cancel).

## 5. Kısa Travel Threshold

< 60 sn travel için push atma (spam koruma). `travel.duration_actual_seconds < 60` → suppress; sadece in-app WS event.

## 6. Cross-Reference
- [travel/README.md](../travel/README.md)
- [push-notification.md](push-notification.md)

## 7. Bekleyen Kararlar
- Interrupted travel'de "quick resume" action button Phase 02.
- Multi-leg travel (waypoint) intermediate push opt Phase 02.
- Suppress threshold (60 sn) v1; Phase 02 tune (user prefs).
