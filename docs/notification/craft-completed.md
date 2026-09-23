# Craft Completed Notification

## 1. Amaç
Crafting queue slot tamamlandığında push. Çoklu slot varsa slot bazlı (queue tamamının değil, adım adım) veya "batch complete" (kullanıcı ayarı).

## 2. Trigger

`event.craft.slot_completed`
```json
{
  "character_id": "...",
  "queue_slot_id": "...",
  "recipe_code": "sword_iron",
  "quality": "normal",
  "critical_craft": false,
  "output_item_instance_id": "..."
}
```

Ayrıca queue tamamı bittiğinde: `event.craft.queue_drained`.

## 3. Payload

```json
{
  "title_key": "notif.craft_completed.title",
  "body_key": "notif.craft_completed.body.single",
  "body_vars": {
    "recipe_name": "Iron Sword",
    "quality_badge": "★ Critical!" 
  },
  "category": "activity",
  "route": "character/crafting/queue",
  "badge_delta": 1,
  "collapse_key": "craft_<character_id>"
}
```

Critical craft → badge ve body vurgusu; rare progression an.

## 4. Batch vs Per-Slot Mode

Preference `notif.craft_completed.mode`:
- `per_slot` — her slot bittiğinde push (kısa recipe için gürültülü).
- `queue_drained` (default) — queue tam boşaldığında tek push.
- `critical_only` — sadece critical craft push.
- `off`

`queue_drained` default → gürültü az; ancak uzun recipe (1h+) tek başına ise per_slot davranışı işe yarar → uygulama içi "always push if slot > 30dk" opt.

## 5. Payload Örnek (Queue Drained)

```
Başlık: Crafting kuyruğu tamamlandı — Grim
İçerik: 5 recipe hazır: 3 iron sword, 2 iron plate.
```

## 6. Cross-Reference
- [crafting/queue.md](../crafting/queue.md), [crafting/critical-craft.md](../crafting/critical-craft.md)
- [push-notification.md](push-notification.md)

## 7. Bekleyen Kararlar
- Critical craft için ayrı sound Phase 02.
- Recipe seçimi ile per-recipe push tercih Phase 03.
