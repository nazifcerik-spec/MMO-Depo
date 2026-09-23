# Session

## 1. Amaç
Sunucu tarafı kalıcı oturum kaydı. Client'a bakış: sadece token; server bakış: `auth.session` satırı. Revoke, listing, anomaly detection buradan.

## 2. Model
Kaynak: [database/auth-accounts.md §3](../database/auth-accounts.md).

Ana alanlar:
- `id` (session id, ULID).
- `account_id`.
- `refresh_token_hash` (rotating).
- `client_kind` (web/ios/android/desktop/cli).
- `client_version`.
- `ip_hash` (SHA-256).
- `user_agent` (truncated 255).
- `issued_at`, `expires_at`, `revoked_at`, `last_seen_at`.

## 3. Lifecycle
```
CREATE (login)
  ↓
ACTIVE (refresh her seferinde last_seen_at + expires_at reset)
  ↓
REVOKED (logout / password change / admin / anomaly / replay)
  ↓
EXPIRED (expires_at geçti, cleanup job)
```

## 4. Session Cap
- Aynı account maksimum aktif session sayısı: **25**.
- Per client_kind: **5**.
- Aşıldığında: en eski session revoke edilir (`reason='cap_exceeded'`).

## 5. Endpoint'ler

### 5.1 Aktif Session Listesi

**GET `/v1/auth/sessions`**
Auth required.

Response 200:
```json
{
  "sessions": [
    {
      "id": "s-...",
      "current": true,
      "client_kind": "ios",
      "client_version": "1.2.3",
      "country_hint": "TR",
      "device_label": "iPhone 15 Pro",
      "issued_at": "2026-09-01T10:00:00Z",
      "last_seen_at": "2026-09-08T13:22:00Z",
      "expires_at": "2026-10-01T10:00:00Z"
    },
    ...
  ]
}
```

Kaynak: `auth.session` + `auth.device` join (device_label için).

### 5.2 Session Revoke

**DELETE `/v1/auth/sessions/{id}`**
Auth required; session must belong to authenticated account.

Response 204.

Server:
1. Validate ownership.
2. UPDATE session revoked_at = now().
3. Redis blocklist için o session'ın jti prefix'i (jti session-scoped olduğundan session_id ile mapping).
4. WebSocket kick.
5. Push device unbind.
6. Audit event.

Error:
- 404 session bulunamadı veya başka account'un.
- 400 kendi mevcut session'ını revoke etmek istiyorsa: uyarı (opsiyonel confirm client tarafı) — server yine izin verir.

## 6. Session Details (Diagnostic)

**GET `/v1/auth/sessions/{id}`** (opsiyonel Phase 02)
Auth required.

Response 200:
```json
{
  "id": "s-...",
  "client_kind": "android",
  "client_version": "1.0.0",
  "device": { "model": "Pixel 8", "os": "Android 14" },
  "geo_history": [
    { "at": "2026-09-01T10:00:00Z", "country": "TR", "city_hint": "Istanbul" },
    ...
  ],
  "last_activity_endpoint": "GET /v1/characters/xxx",
  ...
}
```

## 7. Anomaly Detection
Session için:
- Yeni ülke değişimi (previous last_seen country ≠ current).
- Client_kind ani değişim (aynı session_id ile).
- IP hash farklı ancak kısa süre içinde birkaç değişim.

Aksiyon:
- Bilgi seviyesi: audit + e-mail.
- Yüksek risk: MFA challenge zorla veya session revoke.

## 8. Session Fingerprint (Phase 02 opsiyonel)
- IP hash + user agent + client_kind = "fingerprint hash".
- Mismatch (device change) → soft warn.
- Kesin binding (hard revoke on mismatch) opsiyonel; mobile network IP değişimi false-positive üretebilir.

## 9. Cleanup Job
Worker: `session_cleanup`:
- Her saat başı; `expires_at < now() OR revoked_at < now() - 30 days` → hard delete.
- İlgili `auth.login_history` korunur (retention 6 ay).

## 10. Cross-Session Data
- **Karakter session claim** ayrı (bkz. [Login §14](login.md)).
- **Chat WS connection** session'a bağlı; session revoke → conn kick.

## 11. Metrics
- `auth_active_sessions{client_kind}` gauge.
- `auth_session_revoked_total{reason}` counter.
- `auth_session_avg_lifetime_days`.

## 12. Audit
- `auth.login_history` login/logout/refresh event.
- `audit.event_log account.session_ended` reason ile.

## 13. Fail Modları
- Session cap aşımı: oldest revoke; e-mail'de "yeni cihaz eklendi, eski logout edildi" bildirim.
- Race condition (aynı anda iki logout): idempotent.

## 14. Test Stratejisi
- Unit: cap enforcement.
- Integration: list → revoke → verify blocklist + WS kick.
- Load: 25 concurrent login → 26. cap trigger.

## 15. İlişkiler
- [Login](login.md), [Logout](logout.md), [Refresh Token](refresh-token.md), [Device](device.md), [Account](account.md), [backend/websocket.md](../backend/websocket.md), [database/auth-accounts.md](../database/auth-accounts.md)

## 16. Bekleyen Kararlar
- Session cap final değer (25 önerisi).
- Geo/ASN lookup provider seçimi (MaxMind vs IP2Location vs cloud provider).
- Fingerprint binding sıkılığı (Phase 02).
- Session extend/renewal manuel opsiyonu ("stay logged in longer").
