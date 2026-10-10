# Password

## 1. Amaç
Şifre hash, verify, reset, change politikası. Kullanıcı dostu ancak güvenli.

## 2. Hash Politikası
- **Algoritma**: argon2id (Argon2 v1.3, tip=id).
- **Parametreler (v1 — DD-238 benchmark sonrası OWASP tabanı)**:
  - Memory: 19 MiB (önceki taslak 64 MB)
  - Iterations (time): 2 (önceki 3)
  - Parallelism: 1
  - Salt: 16 byte random
  - Key length: 32 byte
- **Ayarlanabilir**: config'de (upgrade path için).
- Hash format: PHC string (`$argon2id$v=19$m=19456,t=2,p=1$<salt>$<hash>`). Parametreler hash'in içinde; eski hash'ler kendi parametreleriyle doğrulanır.
- **Eşzamanlılık sınırı (DD-238)**: aynı anda en fazla `AUTH_HASH_CONCURRENCY` (varsayılan CPU sayısı) argon2id; 500 ms içinde yer açılmazsa 503 `E_BUSY` + `Retry-After: 1`.

## 3. Verify
- Argon2 lib `VerifyEncoded(hash, password)`.
- Timing-safe: user yoksa bile fake verify (constant time).
- **Rehash on login**: hash parametreleri eskiyse yeni param ile rehash + persist.

```go
ok, needsRehash, err := argon2id.Verify(user.PasswordHash, input, currentParams)
if ok && needsRehash {
    newHash, _ := argon2id.Hash(input, currentParams)
    _ = userRepo.UpdatePassword(ctx, user.ID, newHash, now())
}
```

## 4. Password Kuralları
| Kural | Değer |
|---|---|
| Min length | 10 |
| Max length | 128 |
| Unicode | Allowed |
| Common password check | Top-1M list (haveibeenpwned-like local, hashed) |
| Regex complexity | YOK (length + not-common tercih) |
| Password history | Son 5 password reuse yasak (Phase 02) |

## 5. Change Password (Authenticated)

**POST `/v1/auth/password/change`**
Auth required.

Request:
```json
{ "current_password": "...", "new_password": "..." }
```

Server:
1. Verify current password (constant-time).
2. Validate new password (rules).
3. Rate limit: 5/hour per account.
4. Hash new password.
5. Update `auth.user.password_hash`, `password_updated_at = now()`.
6. **Logout-all** — tüm session'lar revoke (bkz. [logout.md §8](logout.md)); current device için tekrar login gerekir veya "keep this device" opsiyonu (v2).
7. Publish event `account.password_changed`.
8. Audit event.
9. E-mail bildirim.
10. Return 204.

Response 401 `auth.credentials_invalid` if current password wrong.

## 6. Password Reset (Unauthenticated)

### 6.1 Request Reset

**POST `/v1/auth/password/reset-request`**
Rate limit: 3/hour per e-mail + 5/hour per IP.

Request:
```json
{ "email": "..." }
```

Response 202 (always, generic):
```json
{ "status": "If an account exists, a reset e-mail has been sent." }
```
(Enumeration koruma: e-mail yok da olsa 202 döner.)

Server:
1. Lookup user by email (silent if not found; still 202).
2. If found:
   - Invalidate active reset tokens for user.
   - Generate reset token: 32 byte random.
   - Insert `auth.password_reset` (`token_hash = sha256`, `expires_at = now + 1h`).
   - Enqueue reset e-mail (via `notification_dispatch` worker).
3. Delay response 200-500ms (timing consistency).

### 6.2 Confirm Reset

**POST `/v1/auth/password/reset-confirm`**

Request:
```json
{ "token": "...", "new_password": "..." }
```

Server:
1. Hash token → lookup.
2. `expires_at > now && used_at IS NULL`.
3. Validate new password.
4. Update `auth.user.password_hash`, `password_updated_at = now()`.
5. Mark reset used.
6. **Logout-all** for that user (invalidate all sessions).
7. Publish event `account.password_reset`.
8. Audit + e-mail confirmation.
9. Return 204.

Fail: 400 `auth.token_invalid`, 410 `auth.token_expired`.

## 7. E-mail Templates

### Reset Link
```
Subject: MMO — Password reset
Body:
Hi <display_name>,

A password reset was requested for your account. Click below to set a new password:
https://<web>/reset-password?token=<url_safe_b64>

Mobile: mmo://reset-password?token=<...>

This link expires in 1 hour. If you didn't request this, ignore this e-mail.
```

### Password Changed Notice
```
Subject: MMO — Your password was changed
Body:
Your account password was recently changed.

If this wasn't you, contact support immediately.

Time: <utc>
Device: <client_kind>, <country_hint>
```

## 8. Rate Limit Sıkı Katmanı
- Reset request: 3/hour per e-mail + 5/hour per IP.
- Reset confirm: 5/hour per IP (token brute force).
- Change: 5/hour per account.
- Verify (login): [login.md §10](login.md).

## 9. Common-Password Check
- Local hashed blacklist file (haveibeenpwned top-1M SHA-1 first-5 char).
- Runtime: input password sha1 → prefix lookup → suffix compare.
- Match → reject `auth.password_weak`.

Kaynak yenileme: quarterly.

## 10. Password Strength UX
- Client-side hint (green bar) — strength meter (zxcvbn-like veya basit rule).
- Server-side kesin karar (client bypass mümkün).

## 11. Storage
- Sadece hash `auth.user.password_hash`.
- Plaintext asla persist / log / metric.
- Backup DB → hash şifreli (Postgres TDE veya volume encryption).

## 12. Fail Modları
- argon2 hesap yoğun → CPU spike; concurrent limit (semaphore).
- Reset e-mail gönderilemez → 202 zaten döndürülüyor (kullanıcıya deneme öğütlenir).

## 13. Metrics
- `auth_password_change_total{result}`.
- `auth_password_reset_request_total`.
- `auth_password_reset_confirm_total{result}`.
- `auth_password_verify_duration_seconds`.

## 14. Audit
- `account.password_changed`.
- `account.password_reset`.
- `account.password_change_failed` (severity=warn) — brute force sinyal.

## 15. Test Stratejisi
- Unit: hash + verify + rehash trigger.
- Integration: reset happy path e2e.
- Security: timing test (user exists vs not).
- Load: 100 concurrent hash → CPU throttle.

## 16. İlişkiler
- [Login](login.md), [Logout](logout.md), [Register](register.md), [Session](session.md), [database/auth-accounts.md](../database/auth-accounts.md), [systems/notification-system.md](../systems/notification-system.md)

## 17. Bekleyen Kararlar
- argon2 params final (memory/time trade-off with concurrent user goal).
- Password history (last 5) — Phase 02.
- Passphrase mode (BIP39-like) — Phase 03 opsiyonel.
- Post-change "keep this device logged in" — Phase 02.
