# Notification + Mobile Background — Katalog

> **Faz 30**: Push notification altyapısı + Android/iOS mobile lifecycle. Idle MMO'da **retention'ın omurgası**.

## 1. Kapsam
AFK-first oyun için push notification bir "hatırlatma" değil, **oyun döngüsünün parçasıdır**. Karakter aktivitesi bitti → oyuncu bilinçli offline olsa da bilir → geri döner → reveal → yeni aktivite başlar.

Bu klasör:
- Push channels (FCM, APNs, WebPush — Phase 02).
- In-game notification (WS + inbox).
- Notification türleri ve trigger'ları.
- Mobile lifecycle (Doze mode, background execution).
- Kullanıcı preferences (opt-out granular).

## 2. Notification Türleri

| Kod | Kaynak | Tip | Default |
|---|---|---|---|
| `activity_completed` | AFK activity reveal | Push + inbox | ON |
| `dungeon_completed` | Dungeon AFK reveal | Push + inbox | ON |
| `craft_completed` | Crafting queue slot done | Push + inbox | ON |
| `travel_completed` | Travel arrived | Push + inbox | ON |
| `energy_restored` | Energy 100% | Push (silent opt) | OFF (opt-in) |
| `world_event_starting` | World boss / invasion 15-30 dk before | Push | ON (opt-out) |
| `world_event_active` | World boss spawned | Push | ON |
| `guild_event` | Guild boss/raid/war | Push | ON |
| `mail_received` | Yeni mail (attach varsa daha yüksek priority) | Push + inbox | ON |
| `whisper_received` | DM (offline hedef) | Push + inbox | ON |
| `friend_online` | Favorite friend online | Push | OFF |
| `market_sold` | Marketplace listing sold | Push + inbox | ON |
| `market_expiring` | Listing 2h kaldı | Push | OFF |
| `party_invite` | Party invite | Push | ON |
| `guild_invite` | Guild invite | Push + inbox | ON |
| `achievement_unlocked` | Notable (rare+) | Push (silent) | OFF |
| `system_broadcast` | Server-wide ops mesajı | Push + inbox | ON (uncontrollable maintenance) |

## 3. Dokümanlar
- [push-notification.md](push-notification.md) — FCM/APNs entegrasyonu + token yönetim + delivery.
- [activity-completed.md](activity-completed.md) — AFK activity reveal notification.
- [dungeon-completed.md](dungeon-completed.md) — Dungeon AFK completion.
- [craft-completed.md](craft-completed.md) — Crafting queue done.
- [travel-completed.md](travel-completed.md) — Travel arrived.
- [energy-restored.md](energy-restored.md) — Enerji dolduğunda opt-in.
- [world-event.md](world-event.md) — World boss / invasion / seasonal event.
- [guild-event.md](guild-event.md) — Guild boss/raid/war/mail.
- [mobile-lifecycle.md](mobile-lifecycle.md) — Android Doze + iOS Background Modes.
- [preferences.md](preferences.md) — Kullanıcı granular opt-out + quiet hours.

## 4. Cross-Reference
- [systems/notification.md](../systems/notification.md) — v1 in-game notification temel
- [auth/device.md](../auth/device.md) — Device registration + push token
- [backend/websocket.md](../backend/websocket.md) — In-game live push
- [systems/activity-engine.md](../systems/activity-engine.md) — Trigger source

## 5. Prensipler

- **Server-authoritative trigger** — client hiçbir notification kendisi tetiklemez.
- **Idempotent delivery** — aynı notification 2 kez düşmesin (dedup key).
- **Granular opt-out** — her tür için kullanıcı ayarı (default kısmında).
- **Quiet hours** — kullanıcı belirlediği saat aralığında **critical hariç** push sessiz.
- **Category-based** (Android channel / iOS category) — sistem seviyesinde kullanıcı yönetsin.
- **Delivery fallback** — push başarısız ise inbox message queue.
- **Cost-aware** — FCM/APNs quota önemsiz ancak spam retention için zararlı; server-side rate limit.

## 6. Metrikler
- `notif_sent_total{kind, transport}`
- `notif_delivered_total{kind, transport}`
- `notif_failed_total{kind, transport, reason}`
- `notif_click_open_total{kind}`
- `notif_opt_out_gauge{kind}`

## 7. Bekleyen Kararlar
- Rich notification (image / action button) Phase 02.
- WebPush (browser) Phase 02.
- SMS/E-posta kritik uyarı Phase 03.
- Analytics: click-through open rate KPI Phase 02.
