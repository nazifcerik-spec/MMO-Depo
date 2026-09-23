# World Event Notification

## 1. Amaç
World boss, invasion, seasonal event, dynamic event başlangıcı → push. Oyuncunun "kaçırmama" opsiyonu.

## 2. Trigger'lar

| Event | Kind | Timing |
|---|---|---|
| World boss 15-30dk before | `world_event_starting` | Scheduler pre-alert |
| World boss spawn | `world_event_active` | Anında |
| Invasion start | `world_event_active` | Anında |
| Seasonal event open | `world_event_starting` | Anında |
| Dynamic event spawn | `world_event_active` | Sadece opt-in high-value |
| Regional event | (v1'de push YOK) | UI'da list, no push |

## 3. Eligibility Filter

Push tüm oyunculara **atılmaz**; filter:

- **Level range** (`event.min_level ≤ character.level`).
- **Region** (event region'da bulunmuş / recent 7d).
- **Preference** opt-in (kind göre).
- **Registered attendance** (world boss attend → daha yüksek priority).

## 4. Payload

```json
{
  "title_key": "notif.world_event.starting.title",
  "body_key": "notif.world_event.starting.body",
  "body_vars": {
    "event_name": "Ashen Titan Awakens",
    "starts_in_minutes": 15,
    "location": "Blackspire Peak"
  },
  "category": "event",
  "route": "world-event/<spawn_id>",
  "badge_delta": 1,
  "priority": "high",
  "collapse_key": "world_event_<spawn_id>"
}
```

## 5. Rate Limit

Bkz. push-notification.md: `world_event_*` 2/saat, cooldown 10dk. Cap dolarsa öncelik: attendance registered > opt-in filter > default.

## 6. Preferences Segmenti

Preferences UI'da:
- "World Boss başlangıcı" ON (default)
- "İstila (invasion) başlangıcı" ON
- "Seasonal event başlangıcı" ON
- "Dynamic event (küçük)" OFF (default) — spam koruma

## 7. Cross-Reference
- [world-events/README.md](../world-events/README.md)
- [world-events/world-boss.md](../world-events/world-boss.md), [invasion.md](../world-events/invasion.md)
- [preferences.md](preferences.md), [push-notification.md](push-notification.md)

## 8. Bekleyen Kararlar
- Regional event push v1 yok; Phase 02 opt-in.
- Player-triggered dynamic event chain push (rich context) Phase 02.
- Push i18n `location_name` proper noun handling Phase 02.
