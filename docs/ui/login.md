# Login / Register / Logout

## 1. Amaç
Oyuncunun hesabına giriş, yeni hesap kayıt, şifre reset, OTP challenge, logout.

## 2. Ekranlar

- **Splash** → auth token restore attempt → success: character select; fail: login screen.
- **Login** — email + password, "Remember me" (default ON), "Forgot password" link.
- **Register** — email + username + password + password confirm + ToS/Privacy checkbox.
- **OTP Challenge** (2FA / new device) — 6-digit code, resend 60s cooldown.
- **Forgot Password** — email input → send reset link.
- **Reset Password** — deep link → new password + confirm.

## 3. Layout (Login — mobile 375×812)

```
[bg gradient / logo]
─────────────────────
        [Logo]
       "Otter"
─────────────────────
E-posta:       [__________]
Şifre:         [__________] 👁
[✓] Beni hatırla       Şifremi unuttum
─────────────────────
[   GİRİŞ YAP (accent) ]
─────────────────────
Yeni misin? Kayıt ol
```

- Padding 20, form field 48px height, primary CTA 52px.
- Inline error mesajları field altında; global error snackbar üstte.
- Loading state: CTA disable + spinner içeride.

## 4. State

- Field controllers + form validator (email regex, password min 8).
- Auth service (Riverpod provider) — `login()`, `register()`, `logout()`.
- Token cache secure storage (Keychain iOS / EncryptedSharedPreferences Android).

## 5. API

- `POST /v1/auth/register`
- `POST /v1/auth/login` → `{ access_token, refresh_token, expires_at, requires_otp?, otp_challenge_id? }`
- `POST /v1/auth/otp/verify`
- `POST /v1/auth/password/reset-request`
- `POST /v1/auth/password/reset-confirm`
- `POST /v1/auth/logout` (refresh token revoke)

Bkz. [auth](../auth/README.md).

## 6. Error Handling

| Kod | UI davranış |
|---|---|
| `E_INVALID_CREDENTIALS` | Inline "E-posta veya şifre hatalı" |
| `E_ACCOUNT_LOCKED` | Modal "Hesap 24 saat kilitli, güvenlik nedeniyle" |
| `E_OTP_REQUIRED` | OTP ekranına yönlendir |
| `E_RATE_LIMIT` | Snackbar "Çok fazla deneme, X sn sonra" |
| Network fail | Retry button + offline banner |

## 7. Onboarding After Register

Register success → OTP verify → character select (empty) → **new character wizard** açılır.

## 8. Deep Link
`game://route/auth/reset?token=...` → Reset Password screen.

## 9. Metrikler
- `ui_login_success_total`, `ui_login_fail_total{reason}`
- `ui_register_complete_total`
- `ui_password_reset_request_total`

## 10. Cross-Reference
- [auth/register.md](../auth/register.md), [auth/login.md](../auth/login.md), [auth/session.md](../auth/session.md)
- [design-system.md](design-system.md)
