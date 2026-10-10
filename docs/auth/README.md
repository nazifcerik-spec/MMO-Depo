# Auth & Account (Faz 5)

Bu klasör hesap sisteminin uçtan uca tasarımını içerir. Android, iOS ve Web istemcileri **aynı** account sistemini kullanır; token formatı, endpoint sözleşmesi ve güvenlik kuralları platformlar arasında değişmez.

## 1. Kapsam
- **Register** (kayıt) — e-mail + password ile hesap oluşturma; e-mail verification.
- **Login** (giriş) — kimlik doğrulama, token verme, cihaz ilişkilendirme.
- **Logout** (çıkış) — session sonlandırma, token blocklist.
- **Refresh Token** — access token yenileme, rotation.
- **Password** — hash + rotation + reset + change.
- **Session** — sunucu tarafı oturum kaydı, revoke, listing.
- **Device** — cihaz kaydı, trusted device, push token binding.
- **Account** — user ↔ account model, silme, ban, GDPR.

## 2. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | README.md |
| [Register](register.md) | Kayıt akışı + verification |
| [Login](login.md) | Giriş, MFA, tokens |
| [Logout](logout.md) | Session sonlandırma, all-devices |
| [Refresh Token](refresh-token.md) | Rotation, replay detection |
| [Password](password.md) | Hash policy, reset, change |
| [Session](session.md) | Sunucu tarafı session model |
| [Device](device.md) | Device fingerprint + push binding |
| [Account](account.md) | User/Account model, silme, ban |

## 3. Cross-Platform Prensipleri

**Aynı account, farklı cihaz:**
- Access token JWT — platform-agnostic; alg + kid public.
- Refresh token opaque — platform-agnostic.
- Cihaza özel: `client_kind` (web/ios/android/desktop) `auth.session` satırında.

**Aynı hesabın birden çok cihazda aktif oturumu olabilir:**
- Web + Telefon + Tablet aynı anda.
- Session listesi kullanıcıya gösterilir; her birini ayrı ayrı revoke edebilir.

**Karakter senkronu:**
- Karakter server-side'da; cihaz fark etmez. Web'de başladığın aktivite iOS'te aynı anlık state ile devam eder.
- Aynı karakter iki cihazda paralel login mümkün mü? **Hayır (v1)**: karakter tekil "current session claim" tutar; ikinci login birinciyi kicker (opsiyonel: read-only spectator mode Phase 02).

## 4. Ana Sözleşme (Özet)

### 4.1 Endpoint Listesi
```
POST   /v1/auth/register
POST   /v1/auth/verify-email
POST   /v1/auth/resend-verification
POST   /v1/auth/login
POST   /v1/auth/mfa/challenge
POST   /v1/auth/mfa/verify
POST   /v1/auth/refresh
POST   /v1/auth/logout
POST   /v1/auth/logout-all
POST   /v1/auth/password/reset-request
POST   /v1/auth/password/reset-confirm
POST   /v1/auth/password/change
GET    /v1/auth/sessions
DELETE /v1/auth/sessions/{id}
GET    /v1/auth/devices
DELETE /v1/auth/devices/{id}
POST   /v1/auth/devices/register-push
GET    /v1/me
DELETE /v1/me      (account deletion request)
```

### 4.2 Token Model
- **Access token**: JWT, 15 dk TTL, `EdDSA` (öneri) veya `RS256` imzalı.
- **Refresh token**: 32 byte random, base64url; server-side hashed (sha256); 30 gün TTL; **rotation zorunlu**.
- **Blocklist**: Redis `sess:blk:<jti>` — logout / password change / revoke sonrası.

### 4.3 Password Model
- **argon2id** hash (m=19MiB, t=2, p=1 — DD-238).
- Length min 10, max 128; unicode allowed.
- Common-password list check (top-1M blacklist).
- Reset flow: token-hash + 1h TTL + single-use.

### 4.4 Session Model
- Server-side satır (`auth.session`) her aktif refresh için.
- Bir user 20+ eş zamanlı session'a sahip olabilir (cap: 25 client-kind başına 5).
- Anomaly detection (yeni geo, yeni device) → e-mail bildirim.

## 5. Güvenlik İlkeleri

1. **Password never in log**. E-mail hash'li log'da.
2. **Access token URL query'de yasak**. Header veya body.
3. **HTTPS zorunlu**. HTTP dev only.
4. **Rate limit her endpoint**. Detay her doküman.
5. **Enumeration attack koruma**: register + login + password reset **timing** benzer + generic error mesajları.
6. **CSRF**: token-based auth (JWT header) → CSRF risk düşük; ancak refresh cookie kullanılırsa SameSite=Strict + CSRF token.
7. **Audit**: her login/logout/failure `audit.login_history`'e; her password change/reset audit event.
8. **MFA**: opt-in (Phase 01 basic TOTP; Phase 02 WebAuthn).

## 6. GDPR / Compliance
- Data export request (Phase 02): kullanıcı tüm verisini indirir.
- Right-to-be-forgotten: PII anonimize; oyun içi karakter kayıtları ekonomi bütünlüğü için anon olarak korunur.
- Login history retention: 6 ay.

## 7. Ortak Error Response
Bkz. [backend/error-system.md](../backend/error-system.md). Auth-özel kod'lar:
```
auth.credentials_invalid
auth.email_taken
auth.password_weak
auth.rate_limited
auth.token_invalid
auth.token_expired
auth.session_revoked
auth.mfa_required
auth.mfa_invalid
auth.email_not_verified
auth.account_disabled
auth.account_banned
```

## 8. Bekleyen Kararlar
- **Sosyal login** (Google/Discord/Apple): Phase 02 (iOS App Store zorunluluğu: Apple Sign-In gerekebilir Sign-in with other providers varsa).
- **MFA türü v1**: TOTP zorunlu vs opt-in.
- **Session cap** başlangıç değeri.
- **Region-locked account** (Phase 02+ regulatory).
- **Refresh token cookie mi bearer mı** (web için) — öneri: bearer (mobile+web tek pattern).
