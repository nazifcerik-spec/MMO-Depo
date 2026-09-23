# Push Notification (FCM + APNs)

## 1. Amaç
Mobil (Android/iOS) push transport'unu abstract eder. Uygulama backend push queue → provider gateway → cihaz.

## 2. Provider

| Platform | Provider | Notlar |
|---|---|---|
| Android | **FCM (Firebase Cloud Messaging)** HTTP v1 API | Data + notification message |
| iOS | **APNs (Apple Push Notification service)** HTTP/2 API + JWT | Alert + background wake |
| Web | WebPush (VAPID) | Phase 02 |

## 3. Device Token Yönetimi

Kayıt:
```
POST /v1/device/push-token
body: {
  device_id: <uuid>,
  platform: "android"|"ios"|"web",
  push_token: "<provider_token>",
  locale: "tr-TR",
  app_version: "1.4.0"
}
```

Model (bkz. [auth/device.md](../auth/device.md) ile birleşecek):
```sql
CREATE TABLE notification.device_push_token (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL,
  device_id UUID NOT NULL,
  platform VARCHAR(16) NOT NULL,
  push_token TEXT NOT NULL,
  locale VARCHAR(16),
  app_version VARCHAR(32),
  is_active BOOLEAN NOT NULL DEFAULT true,
  registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_used_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_failure_at TIMESTAMPTZ NULL,
  failure_count SMALLINT NOT NULL DEFAULT 0,
  UNIQUE (device_id, platform)
);
CREATE INDEX idx_dpt_account ON notification.device_push_token (account_id, is_active);
```

Token refresh: cihaz FCM/APNs token değiştiğinde re-register; eski token `is_active=false`.

Token invalidation (provider "NotRegistered" / "InvalidToken" cevap): `failure_count++`; 3 ardışık fail → `is_active=false`.

## 4. Notification Model (Send-Side)

```sql
CREATE TABLE notification.push_message (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL,
  character_id UUID NULL,
  kind VARCHAR(48) NOT NULL,        -- 'activity_completed', 'world_event', ...
  dedup_key VARCHAR(96) NOT NULL,   -- unique per (kind + ref) idempotent
  priority VARCHAR(16) NOT NULL DEFAULT 'normal',  -- 'high'|'normal'|'silent'
  payload JSONB NOT NULL,           -- provider-agnostic body
  scheduled_for TIMESTAMPTZ NOT NULL DEFAULT now(),
  sent_at TIMESTAMPTZ NULL,
  delivered_at TIMESTAMPTZ NULL,    -- provider ack
  fail_reason VARCHAR(64) NULL,
  UNIQUE (account_id, dedup_key)
);
CREATE INDEX idx_push_pending ON notification.push_message (scheduled_for) WHERE sent_at IS NULL;
```

Dedup key örneği: `activity_completed:<activity_contract_id>` — retry aynı key'i duplicate düşürmez.

## 5. Payload Format (Provider-Agnostic)

```json
{
  "title_key": "notif.activity_completed.title",
  "body_key": "notif.activity_completed.body",
  "body_vars": { "yield": "42 iron ore" },
  "category": "activity",           // Android channel / iOS category
  "route": "character/activity/{contract_id}",  // client deep link
  "badge_delta": 1,                 // iOS badge count
  "sound": "default",
  "collapse_key": "activity_<char_id>"  // FCM: aynı char'a birden fazla düşünce merge
}
```

Client localization: `title_key` + `body_key` + `body_vars` client'ta resolve edilir (server locale-agnostic).

## 6. Send Pipeline

```
Event bus: event.notification.dispatch
  { kind, account_id, character_id?, payload, dedup_key, priority }
  ↓
push_dispatcher worker:
  1. Preferences check (user opted-out for kind? quiet hours?)
  2. Rate limit (per account per kind, cooldown table)
  3. Insert notification.push_message (ON CONFLICT dedup_key DO NOTHING)
  4. Foreach active device_push_token:
     - Platform-specific format
     - Provider API call (FCM/APNs)
     - Update sent_at, failure_count
  5. Metrics
```

## 7. FCM Format (Android)

HTTP v1 body örnek:
```json
{
  "message": {
    "token": "<device_token>",
    "notification": { "title":"Aktivite Tamamlandı", "body":"42 iron ore" },
    "data": { "route":"character/activity/abc", "kind":"activity_completed" },
    "android": {
      "notification": { "channel_id": "activity", "click_action":"OPEN_ROUTE" },
      "priority": "high"
    },
    "fcm_options": { "analytics_label": "activity_completed" }
  }
}
```

Android channel'lar (client tarafında register):
- `activity` — activity/craft/travel/dungeon completed
- `event` — world/guild event
- `social` — mail/whisper/friend/party/guild invite
- `system` — server broadcast (uncontrollable)
- `silent` — background wake (energy_restored, achievement_unlocked)

## 8. APNs Format (iOS)

```json
{
  "aps": {
    "alert": { "title":"Aktivite Tamamlandı", "body":"42 iron ore" },
    "badge": 3,
    "sound": "default",
    "thread-id": "activity",
    "category": "ACTIVITY_ACTIONS",
    "content-available": 1
  },
  "route": "character/activity/abc",
  "kind": "activity_completed"
}
```

- `apns-priority`: 10 (high) veya 5 (normal).
- `apns-push-type`: `alert` veya `background`.
- Silent push (`content-available: 1`, no alert) — enerji restored gibi.

## 9. Rate Limit

Per account per kind:
```yaml
push_rate_limit:
  activity_completed:  { max_per_hour: 6, cooldown_seconds: 300 }
  craft_completed:     { max_per_hour: 12, cooldown_seconds: 120 }
  travel_completed:    { max_per_hour: 8, cooldown_seconds: 120 }
  energy_restored:     { max_per_hour: 2, cooldown_seconds: 1800 }
  world_event_starting:{ max_per_hour: 2, cooldown_seconds: 600 }
  world_event_active:  { max_per_hour: 2, cooldown_seconds: 600 }
  guild_event:         { max_per_hour: 4, cooldown_seconds: 300 }
  mail_received:       { max_per_hour: 10, cooldown_seconds: 60 }
  whisper_received:    { max_per_hour: 20, cooldown_seconds: 30 }
  friend_online:       { max_per_hour: 5, cooldown_seconds: 60 }
  market_sold:         { max_per_hour: 20, cooldown_seconds: 60 }
```

Cap aşımı → notification `push_message`'a yazılmaz; inbox delivery yine yapılır (fallback).

## 10. Delivery Fallback

Push fail → inbox message queue (in-app "notifications" tab):
```sql
CREATE TABLE notification.inbox_message (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL,
  character_id UUID NULL,
  kind VARCHAR(48) NOT NULL,
  body JSONB NOT NULL,
  route VARCHAR(256) NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  read_at TIMESTAMPTZ NULL,
  expires_at TIMESTAMPTZ NULL       -- default +30d
);
CREATE INDEX idx_inbox_unread ON notification.inbox_message (account_id, read_at) WHERE read_at IS NULL;
```

## 11. Endpoint

**POST `/v1/device/push-token`** — register/refresh.
**DELETE `/v1/device/push-token/{device_id}`** — unregister (logout).
**GET `/v1/inbox?limit=50&cursor=`** — in-app notification list.
**POST `/v1/inbox/{id}/read`**.
**POST `/v1/inbox/read-all`**.

## 12. Metrikler
- `push_sent_total{kind, platform}`
- `push_delivered_total{kind, platform}`
- `push_failed_total{kind, platform, reason}`
- `push_dedup_dropped_total`
- `push_rate_limit_dropped_total{kind}`
- `push_click_open_total{kind}`

## 13. Cross-Reference
- [preferences.md](preferences.md), [mobile-lifecycle.md](mobile-lifecycle.md)
- [auth/device.md](../auth/device.md)
- [systems/notification.md](../systems/notification.md)

## 14. Bekleyen Kararlar
- Rich notification (image + action button) Phase 02.
- WebPush Phase 02 (browser).
- Notification analytics (open rate, dismiss rate) KPI dashboard Phase 02.
- Priority tune (yüksek/normal/silent) balance Phase 02 tune.
