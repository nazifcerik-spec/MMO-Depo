# Account

## 1. Amaç
Kullanıcının kalıcı kimliği. User (kişi) → Account (oyun hesabı) → Character (avatar) üç katmanlı ayrım. Account seviyesindeki ayarlar, silme, ban, GDPR akışları.

## 2. User vs Account
- **User**: kimlik doğrulanabilir kişi (email + password).
- **Account**: bir user'ın oyun-hesabı bağlamı. v1'de 1:1; v2'de multi-account destekli (aile hesabı, senaryo hesabı).

## 3. Data Model
Kaynak: [database/auth-accounts.md](../database/auth-accounts.md).

Ana alanlar:
- `auth.user`: email, password_hash, display_name, locale, country_code, mfa_enabled, status, timestamps.
- `auth.account`: user_id, region, character_slots, bag_slot_bonus, bank_tab_purchased, notification_prefs, subscription_status.

## 4. Endpoint'ler

### 4.1 Kendi Profilim

**GET `/v1/me`**
Auth required.

Response 200:
```json
{
  "user": {
    "id": "u-...",
    "email": "user@example.com",
    "email_verified": true,
    "display_name": "Rhaegar",
    "locale": "tr-TR",
    "country_code": "TR",
    "mfa_enabled": false,
    "created_at": "..."
  },
  "account": {
    "id": "a-...",
    "region": "global",
    "character_slots": 2,
    "bag_slot_bonus": 0,
    "subscription_status": "none"
  },
  "characters": [
    { "id": "c-...", "name": "Rhaegar", "level": 25, "class": "warrior" }
  ],
  "notification_prefs": { ... }
}
```

### 4.2 Profil Update

**PATCH `/v1/me`**
Auth required.

Body (any subset):
```json
{
  "display_name": "New Name",
  "locale": "en-US",
  "country_code": "GB"
}
```

- E-mail değişimi ayrı endpoint (verification gerektirir).
- Password ayrı endpoint (bkz. [password.md](password.md)).

### 4.3 E-mail Change

**POST `/v1/me/email/change-request`**
```json
{ "new_email": "new@example.com", "current_password": "..." }
```
- Current password verify.
- New e-mail unique check.
- Verification token → new e-mail.
- Confirmation endpoint aktif olunca değişim uygulanır.
- Notify old e-mail: "Your e-mail was changed to <masked>".

**POST `/v1/me/email/change-confirm`**
```json
{ "token": "..." }
```

### 4.4 Account Deletion

**DELETE `/v1/me`**
Auth required + current password.

Body:
```json
{ "current_password": "...", "reason": "..." }
```

Server:
1. Verify password.
2. Rate limit 1/day per user.
3. `auth.user.status = 'active'` iken:
   - Set `auth.user.deleted_at = now()` (soft delete).
   - Set `auth.user.status = 'deleted'` — restricted, but not hard delete yet.
   - Revoke all sessions.
   - Enqueue 30-gün countdown job.
4. E-mail: "Your account will be permanently deleted in 30 days. Cancel by logging in."
5. Login yeniden başarılıysa `deleted_at = NULL, status='active'` — cancel.
6. 30 gün sonra worker `account_hard_delete`:
   - PII anonimize (`email = <hashed_uuid>@deleted.local`).
   - Characters:
     - Soft-delete (7 gün başlangıç grace) → hard delete (karakter isim havuza).
     - Marketplace listings iade / mail (30 gün).
     - Guild leadership: transfer prompt zaten yaşanmış olmalı (30 gün önce).
   - Log audit `account.hard_deleted`.

### 4.5 Cancel Deletion

**POST `/v1/me/cancel-deletion`**
Auth required (login mümkün olmalı — status='deleted' iken login limited).

Server:
1. `auth.user.status = 'active'`, `deleted_at = NULL`.
2. Cancel worker countdown.
3. E-mail confirmation.
4. Audit.

### 4.6 Data Export (GDPR — Phase 02)

**POST `/v1/me/data-export`**
Auth required.

Server:
1. Enqueue worker `account_data_export`.
2. Worker: user + account + characters + activity history + market history + chat history → JSON bundle → sign URL.
3. E-mail link to user (24h TTL).

## 5. Account Suspend / Ban (Admin)

**POST `/v1/admin/accounts/{id}/suspend`** (admin route)
```json
{ "duration_hours": 24, "reason": "chat_abuse" }
```
- `auth.user.status = 'disabled'` (temporary) or `'banned'` (permanent).
- Logout all sessions.
- E-mail user.
- Audit `admin_action`.

Auto-unlock: worker job when duration expires.

## 6. Character Slots
- Default: 2.
- Extra slot satın alma (Phase 02 monetization): `auth.account.character_slots += 1`.
- Cap (soft): 10 (Phase 02).

## 7. Region
- v1'de `region = 'global'` sabit.
- Phase 02+: region-locked account (regulatory compliance).
- Marketplace region-aware routing (Phase 02 opsiyonu).

## 8. Notification Preferences
- `auth.account.notification_prefs` JSONB.
- Kategori × kanal matrix (bkz. [systems/notification-system.md](../systems/notification-system.md)).
- Update: `PATCH /v1/me/notification-prefs`.

## 9. Subscription (Phase 02)
- `subscription_status`: none/trial/active/expired.
- Payment provider integration Phase 02.

## 10. Locale & Country
- Locale kullanım: e-mail templates, i18n key resolution.
- Country: GDPR jurisdictional; anomaly detection baseline; region routing.
- User değiştirebilir (`PATCH /v1/me`).

## 11. Multi-Account (Phase 02+)
- Bir user birden fazla account işletebilir (aile).
- Ana kısıtlar:
  - Aynı e-mail tek user (unchanged).
  - Account switcher UI.
  - Marketplace: cross-account trade legal audit gerektirir (Phase 02 karar).

## 12. Fail Modları
- E-mail change: yeni email collision → `auth.email_taken`.
- Account delete: guild leader ise "leadership transfer required" 409.
- Marketplace listing count > threshold → soft warn "listelerin iade olacak, devam?"

## 13. Ownership Contracts
Karakter/guild account'a bağlı; account silinirse:
- **Guild leader**: 30 gün grace → transfer or disband.
- **Active market listings**: 30 gün grace → refund to mail.
- **Guild bank contributions**: account-based değil; kalır.
- **Mailbox**: hard delete anında purge; unclaimed gold → void (economy sink).

## 14. Audit Events
- `account.created`.
- `account.email_changed`.
- `account.display_name_changed`.
- `account.password_changed`.
- `account.deletion_requested`.
- `account.deletion_cancelled`.
- `account.hard_deleted`.
- `account.suspended`.
- `account.banned`.
- `account.reactivated`.

## 15. Cross-Platform Notes
- Web + Mobile + Desktop aynı account.
- Character progression cloud-authoritative; cihaz-özel data (UI settings) local + opsiyonel sync (Phase 02).
- Payment / IAP (Phase 02): platform provider (App Store, Play Store, Steam) integration ayrı per platform; ancak entitlement server-side (Account level).

## 16. Rate Limits
- Profile update: 30/hour.
- E-mail change request: 3/day.
- Deletion request: 1/day.

## 17. Metrics
- `auth_accounts_active` gauge.
- `auth_account_created_total`.
- `auth_account_deleted_total`.
- `auth_account_suspended_total{reason}`.

## 18. Test Stratejisi
- Unit: soft-delete/cancel state transitions.
- Integration: full lifecycle (create → verify → play → delete → cancel → hard delete).
- Compliance: PII anonymization completeness.
- Cross-platform: aynı account 3 platformda paralel action.

## 19. İlişkiler
- [Register](register.md), [Login](login.md), [Session](session.md), [Password](password.md), [Device](device.md), [database/auth-accounts.md](../database/auth-accounts.md), [database/character.md](../database/character.md), [game-design/character.md](../game-design/character.md)

## 20. Bekleyen Kararlar
- Multi-account model Phase 02 detay.
- Region-lock policy.
- Data export format ve delivery mekanizması.
- Subscription tiers ve entitlement modeli.
- IAP idempotency + refund flow (Phase 02).
