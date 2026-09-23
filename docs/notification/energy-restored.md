# Energy Restored Notification

## 1. Amaç
Karakter enerjisi 100%'e ulaştığında opt-in push. Enerji-consumer aktivite hemen başlatmak isteyenler için.

## 2. Karakteristik

- **Default OFF** — çünkü çok sık tetiklenir (birkaç saatte bir).
- Kullanıcı explicit opt-in yapmalı.
- Rate limit: 2/saat, cooldown 30 dk.
- Silent push (Android priority=normal, iOS content-available) — badge only, no sound.

## 3. Trigger

`event.character.energy_full`
```json
{ "character_id": "...", "energy_kind": "energy"|"stamina"|"mana", "capped_at": "<ts>" }
```

**Not:** enerji-benzeri kaynak birden fazla (energy/stamina/mana). Kullanıcı her biri için ayrı opt-in yapabilir.

## 4. Prediction ve Erken Push

Enerji artışı hesaplaması analitik (bkz. AFK engine DD-084). Server, karakterin regen rate'inden `energy_full_at` timestamp'ini hesaplar; worker `energy_watcher` her 5 dk tarayıp `energy_full_at <= now() + 5min` olan karakterlere push scheduler'a eklenmesi.

```sql
CREATE INDEX idx_char_energy_full_at ON character.character (energy_full_at)
  WHERE energy_notification_enabled = true;
```

## 5. Payload

```json
{
  "title_key": "notif.energy_restored.title",
  "body_key": "notif.energy_restored.body",
  "body_vars": { "energy_kind": "energy", "character_name": "Grim" },
  "category": "silent",
  "route": "character/dashboard",
  "priority": "silent",
  "collapse_key": "energy_<character_id>"
}
```

Android channel `silent`, iOS `apns-priority=5, content-available=1, no alert body` (badge only).

## 6. Opt-Out Kolaylık

Preferences UI'da "energy_restored" toggle en üstte kolay erişilebilir; çünkü aşırı bildirim gelirse opt-out müşteri satisfaction için önemli.

## 7. Cross-Reference
- [character/energy.md](../character/energy.md), [character/stamina.md](../character/stamina.md), [character/mana.md](../character/mana.md)
- [preferences.md](preferences.md), [push-notification.md](push-notification.md)

## 8. Bekleyen Kararlar
- Silent push default vs sound (Phase 02 A/B test).
- Energy tipi başına ayrı opt-in mi ortak mı (v1: ortak; Phase 02 ayrı).
- Erken push threshold (5dk before full? %90 full?) Phase 02 UX pass.
