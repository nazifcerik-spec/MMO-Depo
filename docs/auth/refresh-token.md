# Refresh Token

## 1. Amaç
Uzun süreli oturum sürdürme; short-lived access token'ı yenileme; replay detection ile hesap güvenliği.

## 2. Prensipler
- **Rotation zorunlu**: her refresh yeni token üretir; eski invalidate.
- **Opaque**: JWT değil; random bytes; server-side hashed.
- **One-time use**: kullanılmış refresh yeniden kullanılamaz.
- **Replay detection**: kullanılmış refresh ikinci kez gelirse → tüm family revoke + alarm.

## 3. Token Formatı
- 32 byte crypto/rand.
- Base64url encode → client'a `refresh_token`.
- Server:
  - `refresh_token_hash = sha256(raw_bytes)`.
  - Bind: `auth.session.id`, `account_id`, `client_kind`, `expires_at`, `refresh_generation` (integer, artan).

## 4. Endpoint

**POST `/v1/auth/refresh`**

Auth: **hiçbir access token gerekmez**; sadece refresh token.

Request:
```json
{
  "refresh_token": "opaque-b64",
  "device": {
    "install_id": "...",
    "client_version": "..."
  }
}
```

Response 200:
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 900,
  "refresh_token": "new-opaque",
  "refresh_expires_in": 2592000,
  "session_id": "s-..."
}
```

Response 401:
```json
{ "error": { "code": "auth.token_invalid" } }
```

## 5. Server Akışı
```
1. Parse body.refresh_token → sha256 → lookup auth.session by refresh_token_hash.
2. If not found → 401 (auth.token_invalid).
   AND publish security signal (potential replay).
3. If session.revoked_at IS NOT NULL:
     → REPLAY DETECTED (used-before, then reused).
     → Revoke ALL sessions for this account_id (family compromise).
     → Insert audit.security_signal severity='critical'.
     → E-mail user with "your account was accessed from unknown source" template.
     → Return 401.
4. If session.expires_at < now() → 401 (auth.token_expired).
5. If client_version mismatch beklenen check — soft warn only, not reject.
6. Rotate: generate new refresh token; sha256 → update auth.session:
     - refresh_token_hash = new_hash
     - refresh_generation += 1
     - expires_at = now() + 30 days  (sliding window)
     - last_seen_at = now()
7. Also insert a "revoked" old row? — HAYIR, session id sabit tutulur; hash rotation ile yeni token.
   Alternatif: family model — her refresh yeni session row + parent_session_id.
   Karar: **hash rotation on same session row** (v1 basit). Phase 02'de family model.
8. Issue new access token (JWT, 15dk).
9. Return.
```

## 6. Replay Detection Detay
Eski (revoked) refresh token gelirse:
- **Anlamı**: birisi refresh token'ı çaldı ve kullandı; sonra gerçek kullanıcı da kullandı (veya tersi).
- **Aksiyon (v1)**:
  - Tüm session'ları account_id için revoke.
  - Blocklist tüm aktif jti (Redis pub/sub → hub kick).
  - E-mail + push notif "security alert".
  - `audit.security_signal severity=critical`.
- **Aksiyon (v2)**: family model — sadece o family compromise (tek device).

## 7. Rate Limit
- Per session_id: 60 refresh/hour (normal ~4/hour; 15dk access TTL sayesinde).
- Per account: 300/hour (multi-device).
- Aşım → 429.

## 8. Access Token TTL Politikası
- **15 dk** — kısa; revoke etkisi hızlı.
- Trade-off: refresh sıklığı ↑; ancak blocklist yükü ↓ (kısa TTL nedeniyle blocklist entry hızlı expire).

## 9. Refresh Token TTL Politikası
- **30 gün sliding window** (her refresh'te reset).
- 30 gün inaktif → expire; kullanıcı re-login.
- Absolute max (v2): 90 gün — inactivity'ye rağmen force re-login güvenlik için.

## 10. Cross-Platform Storage
Bkz. [login.md §9](login.md).

## 11. Fail Modları
- Refresh token invalid → 401 + timing consistent.
- Session revoked → 401 + potential security signal.
- Session expired → 401.
- Race condition (aynı anda iki refresh): transaction + optimistic lock (`refresh_generation`).

## 12. Race Condition Handling
```sql
UPDATE auth.session
SET refresh_token_hash = $new_hash,
    refresh_generation = refresh_generation + 1,
    expires_at = now() + interval '30 days',
    last_seen_at = now()
WHERE id = $sid
  AND refresh_token_hash = $expected_old_hash
  AND revoked_at IS NULL
RETURNING refresh_generation;
```
- `RETURNING` 0 satır → hash mismatch (concurrent refresh olmuş; birinci kazandı).
- İkinci refresh: eski hash artık DB'de yok → 401 (as if invalid).

Bu yaklaşım ile "aynı anda iki cihaz refresh atarsa" birinci kazanır, ikinci geri gitmelidir. Nadir; UX etkilenmez çünkü mobile client refresh'i tek yerde tutar.

## 13. Client Behavior
```
Access token expiring (< 60 sn) → proactively refresh (background)
401 error on API → refresh flow → retry original request
Refresh 401 → logout local; redirect to login
```

Concurrent request handling (SPA):
- Global "refresh in-flight" promise; ikinci request beklerken kuyruğa alınır.
- Refresh success → all queued retry.
- Refresh fail → all queued fail with 401.

## 14. Metrics
- `auth_refresh_total{result}` (success/expired/invalid/replay_detected).
- `auth_refresh_duration_seconds`.
- `auth_replay_detected_total` — ALERT threshold.

## 15. Audit
- Başarılı refresh normalde audit'e girmez (yüksek volume). Sadece:
  - `first_refresh_of_session` (session'ın ilk refresh'i).
  - `replay_detected`.
  - `refresh_after_long_inactive` (> 7 gün) — bilgi amaçlı.

## 16. Security
- Refresh token asla query string.
- HTTPS zorunlu.
- Sha256 hash olduğundan DB leak → token direkt kullanılamaz (rainbow tablo pratik değil random 32 byte için).
- Server-side hash için pepper (env secret) — Phase 02.

## 17. Test Stratejisi
- Happy path: login → refresh → refresh → access.
- Replay: refresh sonra eski token yeniden → replay alarm.
- Race: goroutine 2 concurrent refresh → biri success biri 401.
- Expired: TTL geçmiş → 401.

## 18. İlişkiler
- [Login](login.md), [Logout](logout.md), [Session](session.md), [Password](password.md) (change → all sessions revoke), [backend/middleware.md](../backend/middleware.md)

## 19. Bekleyen Kararlar
- Family model migration (v2) — session lineage.
- Pepper for hash — Phase 02.
- Absolute refresh cap (90 gün force re-login) — Phase 02.
- Refresh token binding to device fingerprint (Phase 02 anti-theft).
