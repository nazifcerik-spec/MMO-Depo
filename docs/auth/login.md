# Login (Giriş)

## 1. Amaç
Kimlik doğrulama + access/refresh token üretimi + session/device kaydı.

## 2. Kullanıcı Akışı
```
Client
  1. E-mail + password gir → POST /v1/auth/login
     ↓
  Server:
     a. Rate limit (per-IP + per-email)
     b. Load user by email (case-insensitive; timing-safe fail path)
     c. Verify password (argon2id)
     d. Check user.status ('active' değilse özel error)
     e. Check e-mail verified (unverified'a limited access)
     f. MFA gerekiyorsa → mfa_required + challenge_id
     g. Aksi halde:
         - Session insert (auth.session)
         - Access token (JWT) + refresh token üret
         - Push token bind (device)
         - Audit login history
         - Return tokens
  Client:
     h. Access token memory'e; refresh token secure storage
     i. WebSocket bağla (Authorization: Bearer)
```

## 3. Endpoint

**POST `/v1/auth/login`**

Request:
```json
{
  "email": "user@example.com",
  "password": "...",
  "device": {
    "kind": "web|ios|android|desktop",
    "install_id": "uuid-stable-per-install",
    "client_version": "1.0.0",
    "os": "iOS 17.4",
    "model": "iPhone15,3",
    "push_token": "..."  // opsiyonel
  },
  "mfa_code": "123456"  // sadece MFA gerekiyorsa
}
```

Response 200 (başarılı, MFA yok veya doğrulandı):
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 900,
  "refresh_token": "opaque-32-byte-b64",
  "refresh_expires_in": 2592000,
  "session_id": "s-...",
  "account": {
    "id": "a-...",
    "email_verified": true,
    "character_slots": 2
  }
}
```

Response 200 (MFA gerekli):
```json
{
  "mfa_required": true,
  "mfa_challenge_id": "mfa-...",
  "mfa_methods": ["totp"]
}
```

Response 401 (invalid credentials):
```json
{ "error": { "code": "auth.credentials_invalid", "message": "..." } }
```

## 4. Password Verify
- argon2id `VerifyEncoded(password_hash, input_password)`.
- Timing-safe: user yoksa bile fake verify çalıştır (constant time).
- Rehash tetikleyici: hash param eski ise yeni param ile rehash + persist.

## 5. MFA (Phase 01 baseline: TOTP)

### 5.1 Challenge
- Password OK + `user.mfa_enabled=true` → `auth.mfa_required` response + challenge_id.
- Server Redis'te `mfa:chal:<id>` → 5 dk TTL, `{user_id, remaining_attempts=5}`.

### 5.2 Verify

**POST `/v1/auth/mfa/verify`**
```json
{ "challenge_id": "mfa-...", "code": "123456" }
```
- TOTP validate (secret'i encrypt at rest, app-level).
- OK → complete login flow (session, tokens).
- 3 fail → challenge invalidate; login baştan.

### 5.3 Recovery Code
- Phase 01'de opsiyonel; Phase 02'de zorunlu (10 tek-kullanımlık code).

## 6. Session Insert
```
INSERT INTO auth.session (
  id, account_id, refresh_token_hash, client_kind, client_version,
  ip_hash, user_agent, issued_at, expires_at
) VALUES (
  ulid(), $account, sha256($rt), 'ios', '1.0.0',
  sha256($ip), truncate($ua, 255), now(), now() + interval '30 days'
);
```

## 7. Access Token (JWT)
Claim:
```json
{
  "iss": "mmo",
  "sub": "<account_id>",
  "aud": ["mmo-api"],
  "exp": 1725720000,
  "iat": 1725719100,
  "jti": "<session_id-scoped-jti>",
  "sid": "<session_id>",
  "uid": "<user_id>",
  "kind": "access",
  "roles": ["player"]
}
```
- Alg: `EdDSA` (kid from `auth.jwt_key`).
- TTL: 15 dk.
- Signing key rotation: kid ile versionlu; validation cache 60s.

## 8. Refresh Token
- Opaque 32 byte random (crypto/rand).
- Base64url encode → client'a.
- Server sadece **hash** tutar (`refresh_token_hash = sha256(rt)`).
- TTL 30 gün.
- **Rotation zorunlu** (bkz. [refresh-token.md](refresh-token.md)).

## 9. Cross-Platform Storage Kuralı
| Platform | Access Token | Refresh Token |
|---|---|---|
| Web | Memory (JS var; sekme kapanınca kaybolur) | HttpOnly Secure SameSite=Strict cookie **veya** Memory (SPA state) |
| iOS | Memory | Keychain (SecItem) |
| Android | Memory | EncryptedSharedPreferences (Keystore) |
| Desktop | Memory | OS credential store |

Öneri (v1): **her platformda bearer body**; access token memory; refresh token secure storage. Web'de cookie route Phase 02 opsiyonu.

## 10. Rate Limit
- Per-IP: 20/dk.
- Per-e-mail: 10/dk (brute-force koruma).
- Password fail: 5. hatalı denemede geçici lock 15 dk (login lock, `user_status=disabled` DEĞİL).
- CAPTCHA (Phase 02) after N fails.

## 11. Login Failure Codes
| Kod | Nedeni |
|---|---|
| `auth.credentials_invalid` | E-mail veya password hatalı (generic) |
| `auth.rate_limited` | Rate limit |
| `auth.account_disabled` | user.status='disabled' |
| `auth.account_banned` | user.status='banned' |
| `auth.email_not_verified` | Restricted access uyarısı — token DÖNER ama restricted flag |
| `auth.mfa_required` | MFA challenge |
| `auth.mfa_invalid` | MFA kod hatalı |

## 12. Session Cap
- Aynı account'ta aktif session cap: **25**.
- Aynı client_kind için cap: **5** (5 web + 5 iOS + 5 android + 5 desktop + 5 CLI).
- Cap aşımı: en eski session revoke edilir (auto-logout that device).

## 13. Anomaly Detection
- Yeni ülke (country hint IP lookup) + değişik client_kind → e-mail bildirim.
- Aynı account 5 dk içinde 3+ farklı ülke → challenge (MFA force) veya block.
- `audit.security_signal` insert.

## 14. Cross-Character Session
- Login karakterlerinden birini otomatik seçmez. Session hesap seviyesi.
- Karakter seçimi ayrı endpoint (`POST /v1/characters/{id}/enter`) — o karakter için `current_session_claim` alır; başka session aynı karaktere `enter` ederse birinci karakter kick.

## 15. Metrics
- `auth_login_total{result,reason}` (result: success/fail; reason: bad_password, mfa_pending, ...).
- `auth_login_duration_seconds`.
- `auth_active_sessions` gauge (per client_kind).

## 16. Audit
- `auth.login_history` her attempt (success + fail + fail_reason).
- Push event `account.session_created`.

## 17. Test Stratejisi
- Unit: password verify + timing constant.
- Integration: full happy path + MFA path.
- Security: brute force → lock; enumeration timing test.

## 18. İlişkiler
- [Register](register.md), [Refresh Token](refresh-token.md), [Session](session.md), [Device](device.md), [Password](password.md), [Logout](logout.md), [Account](account.md), [database/auth-accounts.md](../database/auth-accounts.md)

## 19. Bekleyen Kararlar
- Recovery code Phase 01 zorunlu mu?
- Login lock süresi ve fail eşiği final.
- Cookie-based refresh (web) opt-in mı zorunlu mu — v1 bearer, v2 hybrid?
- Passkey / WebAuthn — Phase 02.
