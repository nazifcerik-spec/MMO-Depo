# Notification Preferences

## 1. Amaç
Kullanıcı push ve in-app notification'larını granular kontrol edebilsin. Retention için sağlıklı denge — kullanıcı doğru bilgiyi alır, spam yaşamaz.

## 2. Yapı

Her notification kind için ayrı ayar; ek olarak:
- **Global mute** (all off — critical hariç).
- **Quiet hours** — belirli saat aralığında sessiz.
- **Per-device override** — sadece bir cihazda push (v1'de tüm cihazlar aynı).

## 3. Model

```sql
CREATE TABLE notification.preference (
  account_id UUID PRIMARY KEY,
  quiet_hours_start VARCHAR(5) NULL,     -- "22:30"
  quiet_hours_end VARCHAR(5) NULL,       -- "08:00"
  quiet_hours_tz VARCHAR(32) NULL,       -- "Europe/Istanbul"
  global_mute_until TIMESTAMPTZ NULL,    -- global mute expire
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE notification.preference_kind (
  account_id UUID NOT NULL,
  kind VARCHAR(48) NOT NULL,
  mode VARCHAR(16) NOT NULL,             -- 'push_and_inbox'|'inbox_only'|'silent_push'|'off'
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (account_id, kind)
);
```

Default değerler `content.notification_kind_defaults`:
```yaml
activity_completed:     push_and_inbox
dungeon_completed:      push_and_inbox
craft_completed:        push_and_inbox
travel_completed:       push_and_inbox
energy_restored:        off
world_event_starting:   push_and_inbox
world_event_active:     push_and_inbox
guild_event:            push_and_inbox
mail_received:          push_and_inbox
whisper_received:       push_and_inbox
friend_online:          off
market_sold:            push_and_inbox
market_expiring:        off
party_invite:           push_and_inbox
guild_invite:           push_and_inbox
achievement_unlocked:   silent_push
system_broadcast:       push_and_inbox     # override edilemez (bkz. §5)
```

## 4. Quiet Hours

- Kullanıcı `22:30-08:00 Europe/Istanbul` gibi tanımlar.
- Bu pencerede tüm push priority `silent` (badge only, no sound); inbox delivery normal.
- **Kritik istisna**: `system_broadcast` (server maintenance, security alert) her koşulda ses açık.

## 5. Uncontrollable Kinds

Aşağıdaki notification'lar opt-out edilemez (yasal / güvenlik):
- `system_broadcast` (server maintenance, security)
- `account_security` (login from new device, password change)
- `payment_receipt` (Phase 02 store)

UI'da bu kind'lar disabled toggle gösterilir + açıklama.

## 6. Global Mute (Vacation Mode)

Kullanıcı X gün push kapatmak isterse:
`POST /v1/notification/preference/global-mute body: { hours: 24 }` → 24 saat push tamamen sessiz (inbox delivery devam). Expire otomatik.

## 7. Endpoint

**GET `/v1/notification/preference`** — full snapshot.
**PUT `/v1/notification/preference`** body:
```json
{
  "quiet_hours_start": "22:30",
  "quiet_hours_end": "08:00",
  "quiet_hours_tz": "Europe/Istanbul",
  "kinds": {
    "energy_restored": "silent_push",
    "friend_online": "off"
  }
}
```
**POST `/v1/notification/preference/global-mute`** body: `{ hours: 24 }`.
**DELETE `/v1/notification/preference/global-mute`** — mute iptal.

## 8. UI Guidance

Ayarlar → Bildirimler ekranı önerisi:
- En üstte "Global Mute" toggle + quiet hours.
- Kategori grupları (Aktivite, Etkinlik, Sosyal, Guild, Sistem).
- Her kind: 4 mode (push+inbox / inbox / silent / off) + short description.
- Cihaz OS bildirim ayarları için deep link (Android/iOS system settings).

## 9. Metrikler
- `notif_pref_updated_total{kind, mode}`
- `notif_pref_global_mute_active_gauge`
- `notif_pref_opt_out_ratio{kind}`
- `notif_quiet_hours_suppressed_total`

## 10. Cross-Reference
- [push-notification.md](push-notification.md), [mobile-lifecycle.md](mobile-lifecycle.md)
- [auth/account.md](../auth/account.md) — security kinds

## 11. Bekleyen Kararlar
- Per-device override Phase 02 (aynı account, farklı cihazda farklı pref).
- Category-level toggle (Aktivite ALL off) Phase 02.
- Notification snooze (bir tık kapat 1 saat) Phase 02.
- Localization tz auto-detect Phase 02.
