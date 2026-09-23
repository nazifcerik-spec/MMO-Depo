# Register (Kayıt)

## 1. Amaç
Yeni hesap oluşturma. Cross-platform (Android/iOS/Web) tek endpoint.

## 2. Kullanıcı Akışı
```
Client (any platform)
  1. E-mail + password + display_name + locale gir
  2. POST /v1/auth/register
     ↓
  Server:
     a. Validate (format, password strength, blacklist)
     b. Rate limit (per-IP + per-email)
     c. Hash password (argon2id)
     d. Insert auth.user (status='active', email_verified_at=NULL)
     e. Insert auth.account (default character_slots=2)
     f. Generate email verification token → auth.email_verification
     g. Enqueue verification e-mail (worker: notification_dispatch)
     h. Optional: auto-login (issue access+refresh) — öneri: verification istenmeden login mümkün ama e-mail unverified iken restricted (karakter yaratamaz).
     i. Log audit event 'account.registered'
  Client:
     3. "Check your inbox" mesajı
     4. Kullanıcı e-mail linkine tıklar → verify-email
```

## 3. Endpoint

**POST `/v1/auth/register`**

Request:
```json
{
  "email": "user@example.com",
  "password": "correct-horse-battery-staple",
  "display_name": "Rhaegar",
  "locale": "tr-TR",
  "country_code": "TR",
  "device": {
    "kind": "web|ios|android|desktop",
    "install_id": "uuid-stable-per-install",
    "client_version": "1.0.0"
  },
  "accepts_terms": true,
  "accepts_privacy": true
}
```

Response 201:
```json
{
  "user_id": "u-...",
  "account_id": "a-...",
  "email_verification_required": true,
  "access_token": "eyJ...",  // opsiyonel — aşağıda §7
  "refresh_token": "opaque",
  "expires_in": 900,
  "session_id": "s-..."
}
```

Response 409 (e-mail taken):
```json
{ "error": { "code": "auth.email_taken", "message": "..." } }
```

## 4. Validation

| Alan | Kural |
|---|---|
| `email` | RFC 5321 lite check + MX opsiyonel (Phase 02) |
| `password` | length 10..128; NOT in top-1M common; entropy soft-check |
| `display_name` | 2..48 char; unicode-friendly; trimmed |
| `locale` | BCP-47; whitelist supported (tr-TR, en-US, ... — Phase 01 config) |
| `country_code` | ISO 3166-1 alpha-2 opsiyonel |
| `accepts_terms` | true zorunlu |
| `accepts_privacy` | true zorunlu |

Password strength (soft):
- Boş, sadece rakam, ardışık char → reject.
- Common-password blacklist (haveibeenpwned-like local list, hashed).
- Regex katı yerine "length + not-common" tercih (kullanıcı dostu).

## 5. Rate Limit
- Per-IP: 5 register/hour.
- Per-e-mail domain (Phase 02): domain-based anomaly.
- 429 → `Retry-After` header.

## 6. Enumeration Koruma
- E-mail taken durumunda **yavaş** ret veya **generic mesaj** — trade-off:
  - **Trade-off kararı**: klasik "account created; check your inbox" mesajı verilip gerçek durum e-mail üstünden bildirilebilir. Ancak UX kötüleşir.
  - **Öneri (v1)**: gerçek `409 auth.email_taken` döndürülür; **rate limit sıkı** + audit sinyali.

## 7. Auto-Login Karar
- **Karar (v1)**: register sonrası **otomatik login** yapılır (access+refresh döndürülür), ancak `email_verified_at IS NULL` olduğu için:
  - Karakter yaratma reddedilir (`auth.email_not_verified`).
  - Marketplace, guild, chat kapalı.
  - Sadece profile + verify + resend endpoint'leri açık.
- Verification tamamlandıktan sonra tam erişim.

## 8. E-mail Verification

### 8.1 Trigger
- Register veya `POST /v1/auth/resend-verification`.

### 8.2 Data
- `auth.email_verification` satırı:
  - `token_hash = sha256(random_32_byte)`.
  - `expires_at = now + 24h`.
  - `used_at` NULL.

### 8.3 E-mail Template
- Link: `https://<web>/verify-email?token=<url_safe_base64>`.
- Mobile deep link: `mmo://verify-email?token=...`.

### 8.4 Verify Endpoint

**POST `/v1/auth/verify-email`**
```json
{ "token": "..." }
```

Server:
1. hash token → lookup.
2. `expires_at > now` && `used_at IS NULL`.
3. Update `auth.user.email_verified_at = now()`.
4. Mark verification used.
5. Publish event `account.email_verified`.
6. Return 200.

Fail: 400 `auth.token_invalid` or 410 `auth.token_expired`.

### 8.5 Resend

**POST `/v1/auth/resend-verification`**
Auth required (session ile bağlı).
- Rate limit: 3/hour/account.
- Latest active verification revoke; yeni token.

## 9. Post-Register Session
- Cihaz için `auth.session` satırı oluşturulur.
- Push token varsa `auth.device.push_token` bind (bkz. [device.md](device.md)).
- Login history: `event_kind='login'` (register auto-login).

## 10. Cross-Platform Notes
- **Web**: refresh token opsiyonel olarak HttpOnly Secure SameSite=Lax cookie'de tutulabilir. Öneri: her platform bearer body (tek pattern; XSS koruma frontend'de).
- **iOS/Android**: refresh token secure storage (Keychain / Keystore).
- **Desktop**: OS credential store (Windows Credential Manager, macOS Keychain, GNOME Keyring).

## 11. Fail Modları
- E-mail service down → verification e-mail gönderilemez; kayıt yine olur; resend flow.
- DB down → 503; retry.
- Race condition (aynı anda iki register aynı email): unique index + `auth.email_taken`.

## 12. Metrics
- `auth_register_total{result}`
- `auth_register_duration_seconds`
- `auth_email_verification_sent_total`
- `auth_email_verification_success_total`

## 13. Audit Events
- `account.registered` — subject: user_id.
- `account.email_verified`.
- `account.verification_failed` (severity=warn).

## 14. Test Stratejisi
- Unit: password strength, e-mail parser.
- Integration: full register → verify happy path.
- Fuzz: exotic e-mail formats.
- Load: 100 concurrent register on same IP → rate-limited.

## 15. İlişkiler
- [Login](login.md), [Password](password.md), [Session](session.md), [Device](device.md), [Account](account.md), [database/auth-accounts.md](../database/auth-accounts.md), [backend/middleware.md](../backend/middleware.md)

## 16. Bekleyen Kararlar
- Auto-login post-register: v1 evet; v2'de enumeration riski bazlı revize.
- E-mail domain blacklist (disposable e-mail servisleri) — Phase 02.
- CAPTCHA gerekli mi (spam register) — Phase 02+ metric-driven.
