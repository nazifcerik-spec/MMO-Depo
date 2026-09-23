# Logout (Çıkış)

## 1. Amaç
Aktif session'ı sonlandırmak; access token blocklist'e; refresh token invalidate.

## 2. Türler
- **Logout (single)** — sadece bu cihaz.
- **Logout all** — tüm session'lar.
- **Revoke session by id** — hesap ayarlarından belirli bir session (bkz. [session.md](session.md)).

## 3. Endpoint'ler

### 3.1 Single Device

**POST `/v1/auth/logout`**
Auth: Bearer access token.
Body opsiyonel:
```json
{ "refresh_token": "..." }
```

Response 204 (no content).

Server:
1. Decode access token → session_id (jti/sid).
2. `UPDATE auth.session SET revoked_at = now() WHERE id = $sid AND revoked_at IS NULL`.
3. Redis blocklist: `SET sess:blk:<jti> 1 EX <access_token_remaining_ttl>`.
4. WebSocket hub → aktif conn kill (varsa).
5. Push token unbind (`auth.device.push_token = NULL` for this session's device).
6. Publish event `account.session_ended`.
7. Log `audit.login_history event='logout'`.

Idempotent: zaten revoked ise 204 döner (no-op).

### 3.2 Tüm Cihazlar

**POST `/v1/auth/logout-all`**
Auth: Bearer access token.
Response 204.

Server:
1. Load `account_id` from ctx.
2. `UPDATE auth.session SET revoked_at = now() WHERE account_id = $a AND revoked_at IS NULL RETURNING id, jti_prefix?`.
3. Redis blocklist tüm affected jti'lar (Phase 02 optimization: user-scope blocklist).
4. WebSocket hub → tüm conn kill (fanout via Redis).
5. Publish `account.all_sessions_ended`.
6. Audit event.

### 3.3 Belirli Session'ı Revoke

**DELETE `/v1/auth/sessions/{id}`**
Auth required.
Only owner (account_id match).

Response 204 or 404.

## 4. Access Token Blocklist Politikası
- JWT stateless olduğundan revoke için Redis blocklist.
- Key: `sess:blk:<jti>` veya `sess:blk:<sid>` (session_id).
- TTL: access token'ın kalan süresi (max 15 dk).
- Middleware her access token verify'da blocklist check.

Redis unreachable → auth failure (fail-closed on protected routes).

## 5. WebSocket Cleanup
- Hub tarafında `SessionRevoked` event listener:
  - Redis pub/sub `pub:session:revoked:<sid>`.
  - Aktif conn(lar) close frame 4001 (auth revoked) ile kapanır.

## 6. Push Token Unbind
- `auth.device.push_token = NULL` for that session's device.
- Provider (FCM/APNS) token unsubscribe (opsiyonel; gelen mesaj karakter unread stays).

## 7. Cross-Platform Client Davranışı
| Platform | Post-Logout Aksiyonu |
|---|---|
| Web | Memory'deki access + local refresh sil; login sayfasına yönlendir |
| iOS | Keychain'den refresh token sil; app root'a dön |
| Android | EncryptedSharedPreferences'ten sil; ana ekran |
| Desktop | OS credential store'dan sil; login screen |

## 8. Server-Side Logout (Otomatik)
Tetikleyiciler:
- **Password change / reset** → tüm session'lar logout-all.
- **Account banned/disabled** → aynı.
- **User request in support ticket** (admin action, audit event kaydı).
- **Anomaly detection critical severity** → tüm session'lar (Phase 02).

## 9. Rate Limit
- Logout: 30/dk/user (aşırı değil; auto-logout'lar için tolerant).
- Logout-all: 5/dk/user.

## 10. Fail Modları
- Token zaten expired → 204 (idempotent).
- Session zaten revoked → 204.
- Access token invalid signature → 401.
- Redis down → blocklist yazılamaz; fail-closed vs fail-open trade-off. Öneri: DB session revoke başarılıysa 204 döner; Redis retry background (worker outbox).

## 11. Race Condition
- Concurrent logout + refresh: refresh sırasında session revoked check yapılır (bkz. [refresh-token.md §replay detection](refresh-token.md)).

## 12. Metrics
- `auth_logout_total{scope}` (scope: single/all/admin/auto).
- `auth_active_sessions` (gauge, logout ile decrement).

## 13. Audit
- `auth.login_history event_kind='logout'`.
- `audit.event_log kind='account.session_ended'`.
- Admin-forced logout → `audit.admin_action`.

## 14. Test Stratejisi
- Unit: blocklist writer + TTL.
- Integration: login → logout → refresh (fail expected) → access token (fail).
- Redis outage: logout DB revoke works; blocklist eventual consistency.

## 15. İlişkiler
- [Login](login.md), [Refresh Token](refresh-token.md), [Session](session.md), [Device](device.md), [Password](password.md), [Account](account.md), [database/auth-accounts.md](../database/auth-accounts.md), [backend/websocket.md](../backend/websocket.md)

## 16. Bekleyen Kararlar
- User-scope blocklist (single Redis key with session set) mi per-jti mi — Phase 02.
- Provider token unsubscribe policy (FCM/APNS).
- Logout confirmation UI (özellikle mobile — accidental logout korumak) client tarafı.
