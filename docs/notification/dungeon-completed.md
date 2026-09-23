# Dungeon Completed Notification

## 1. Amaç
Dungeon AFK run tamamlandığında push. Dungeon reveal daha uzun sürer (multi-floor) — ayrı kind.

## 2. Trigger

`event.dungeon.afk_run_completed`
```json
{
  "character_id": "...",
  "dungeon_run_id": "...",
  "dungeon_template_code": "iron_mine_lvl_20",
  "floors_cleared": 5,
  "boss_defeated": true,
  "yield_summary": { "gold_copper": 4500, "epic_drops": 1 }
}
```

## 3. Payload

```json
{
  "title_key": "notif.dungeon_completed.title",
  "body_key": "notif.dungeon_completed.body.success",
  "body_vars": {
    "dungeon_name": "Iron Mine",
    "floors": 5,
    "loot_hint": "1 epic + 4500c"
  },
  "category": "activity",
  "route": "character/dungeon/<dungeon_run_id>/reveal",
  "badge_delta": 1,
  "priority": "high",
  "collapse_key": "dungeon_<character_id>"
}
```

Priority `high` — normally 30-60 dk dungeon; oyuncu daha uzun offline kaldıysa "loot bekliyor" sinyali güçlü.

## 4. Fail Case

Dungeon fail (wipe):
```json
"body_key": "notif.dungeon_completed.body.failed",
"body_vars": { "dungeon_name": "Iron Mine", "floor": 3, "cause": "boss_wipe" }
```

Priority `normal` (kötü haber; push spam etme).

## 5. Rate Limit

Bkz. push-notification.md: activity kategorisi ile aynı bucket (6/saat). Dungeon ayrı category istenirse Phase 02.

## 6. Cross-Reference
- [dungeon/afk-run.md](../dungeon/afk-run.md)
- [push-notification.md](push-notification.md), [activity-completed.md](activity-completed.md)

## 7. Bekleyen Kararlar
- Fail durumunda ayrı category (opt-out müstakil) Phase 02.
- Best-run rekor push (yeni high score) Phase 02.
