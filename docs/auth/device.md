# Device

## 1. Amaç
Cihaz kaydı: install id, model/os, push token, trusted device flag. Session ile 1:1 değil; bir device birden fazla session olabilir (aynı device'a farklı zamanlarda login).

## 2. Neden Ayrı Bir Kavram?
- **Session** = auth token life-cycle.
- **Device** = fiziksel/logical istemci — install id sabit; refresh sonrası bile aynı.
- **Push notification** hedefi device seviyesindedir.
- **Trusted device** flag → yeni login'de MFA skip (Phase 02).

## 3. Data Model

Yeni tablo `auth.device` (Phase 5 eklentisi):

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `account_id` | `UUID` | FK → auth.account NOT NULL |
| `install_id` | `TEXT` | NOT NULL — client-generated stable |
| `kind` | `auth.client_kind` | NOT NULL |
| `os_name` | `VARCHAR(32)` | NULL — "iOS 17.4" |
| `model_hint` | `VARCHAR(64)` | NULL — "iPhone15,3" |
| `client_version_last` | `VARCHAR(32)` | NULL |
| `label` | `VARCHAR(48)` | NULL — user-set "Ali'nin telefonu" |
| `push_provider` | `auth.push_provider` | NULL — fcm/apns/webpush |
| `push_token` | `TEXT` | NULL — provider token |
| `push_token_updated_at` | `TIMESTAMPTZ` | NULL |
| `trusted` | `BOOLEAN` | NOT NULL DEFAULT false |
| `first_seen_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `last_seen_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `revoked_at` | `TIMESTAMPTZ` | NULL |

Enum:
```sql
CREATE TYPE auth.push_provider AS ENUM ('fcm','apns','webpush');
```

Unique: `(account_id, install_id)`.
Indexler:
- `ix_device__account_id`.
- `ix_device__push_provider__push_token` (send notification lookup).

## 4. Install ID
- Client tarafından oluşturulur (UUID v4).
- Persist:
  - **iOS**: Keychain (uninstall'da opsiyonel kalır).
  - **Android**: Encrypted SharedPreferences veya `Settings.Secure.ANDROID_ID` (uygunsuz — reset edilebilir; UUID persist tercih).
  - **Web**: localStorage + fallback cookie. `first_use_at` cookie ile beraber.
  - **Desktop**: OS credential store.
- Uninstall & yeniden install → yeni install_id → yeni device row.

## 5. Login Sırasında Device Sync

Login endpoint body içindeki `device` alanı işlenir:
```
1. Lookup device by (account_id, install_id).
2. Yoksa create.
3. Update last_seen_at, client_version_last, os_name, model_hint.
4. Push_token varsa update push_token + provider.
5. Return device_id (session içinde saklanır — auth.session'a device_id FK opsiyonel eklenebilir).
```

## 6. Endpoint'ler

### 6.1 Device Listesi

**GET `/v1/auth/devices`**
Auth required.

Response 200:
```json
{
  "devices": [
    {
      "id": "d-...",
      "kind": "ios",
      "model_hint": "iPhone15,3",
      "os_name": "iOS 17.4",
      "label": "Ali'nin telefonu",
      "trusted": true,
      "push_enabled": true,
      "first_seen_at": "...",
      "last_seen_at": "...",
      "active_session_count": 1
    }
  ]
}
```

### 6.2 Device Revoke

**DELETE `/v1/auth/devices/{id}`**
Auth required.

Server:
1. Ownership check.
2. `auth.device.revoked_at = now()`.
3. İlgili session'lar (device_id ile eşleşen) revoke.
4. Push token unbind (provider unsubscribe opsiyonel).
5. Audit.

Response 204.

### 6.3 Device Label Update

**PATCH `/v1/auth/devices/{id}`**
```json
{ "label": "Kişisel iPhone" }
```

### 6.4 Trust Device (Phase 02 MFA-skip)

**POST `/v1/auth/devices/{id}/trust`**
Auth required + MFA verified recently (< 5 dk).
- `trusted = true`.
- Süre: 30 gün (auto revert).

### 6.5 Push Registration

**POST `/v1/auth/devices/register-push`**
Auth required.

Request:
```json
{
  "install_id": "...",
  "push_provider": "fcm|apns|webpush",
  "push_token": "..."
}
```

Server:
1. Lookup device by install_id + account_id.
2. Update push_provider, push_token, push_token_updated_at.
3. Return 204.

## 7. Push Notification Delivery
Notification service (worker `notification_dispatch`):
```
For target character:
  → account
  → devices (revoked_at IS NULL AND push_token IS NOT NULL)
  → per push_provider: FCM API / APNS / Web Push
```

Provider-specific:
- **FCM (Android + Web)**: server key veya OAuth token; batch send.
- **APNS (iOS)**: p8 key + team id + bundle id; HTTP/2 stream.
- **Web Push**: VAPID key; endpoint per subscription.

## 8. Push Token Rotation
- iOS: APNS token gerçek zamanlı değişebilir (iOS event); client server'a `register-push` gönderir.
- FCM: benzer.
- Server: eski token deactivate (yeni gelenle üzerine yazılır).

## 9. Delivery Fail Handling
- 4xx invalid token → device push_token = NULL + audit.
- 5xx transient → retry (worker job backoff).

## 10. Rate Limit
- Push register: 10/hour per device (rare update).
- Label update: 30/hour per account.
- Device revoke: 30/hour.

## 11. Fail Modları
- Push provider outage → notification degrade (in-game inbox yine dolu).
- Install_id collision (nadir; UUID v4) — retry with new UUID client side.

## 12. Cross-Platform Notlar
| Platform | Push Provider | Notes |
|---|---|---|
| iOS | APNS | Real device zorunlu (simulator sınırlı) |
| Android | FCM | Google Play Services (China: alternative Phase 03) |
| Web | Web Push (VAPID) | Service Worker gerekli |
| Desktop | Native OS notif (Phase 02) | Toast API |

## 13. Metrics
- `auth_devices_active` gauge per kind.
- `push_send_total{provider,result}`.
- `push_delivery_duration_seconds`.

## 14. Audit
- `account.device_registered`.
- `account.device_revoked`.
- `account.push_token_updated`.

## 15. Privacy
- IP address device'a bind YAPILMAZ (session'da hash tutulur).
- Model_hint fine-grained fingerprint değil (jenerik: "iPhone15,3" OK; hardware serial YASAK).

## 16. Test Stratejisi
- Unit: install_id upsert.
- Integration: register-push → send notification → provider mock.
- Cross-platform: iOS+Android+Web aynı account'ta 3 device → 3 push.

## 17. İlişkiler
- [Login](login.md), [Session](session.md), [Logout](logout.md), [Account](account.md), [systems/notification-system.md](../systems/notification-system.md), [database/auth-accounts.md](../database/auth-accounts.md)

## 18. Bekleyen Kararlar
- `auth.device.id` `auth.session`'a FK eklenmesi (device_id) — Phase 01'de karar.
- Trusted device süresi (30 gün önerisi).
- Web Push subscription persistence detayları.
- Multi-account per device (kullanıcı aynı device'ta iki account) — desteklenir, ancak per-account row.
