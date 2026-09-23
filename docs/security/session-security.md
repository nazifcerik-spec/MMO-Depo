# Session Security

## 1. Amaç
Auth token lifecycle, rotation, revocation, device fingerprint, anomaly-based session invalidate.

## 2. Token Türleri

| Token | TTL | Refresh | Kullanım |
|---|---|---|---|
| **Access token** (JWT) | 15 dk | rotating refresh | API auth header `Bearer` |
| **Refresh token** | 30 gün | tek kullanımlık | new access + new refresh (rotasyon) |
| **Session token** (server-side) | 30 gün | refresh ile bağlı | Redis session state |
| **WebSocket auth token** | 5 dk | connect anında | WS handshake |
| **OTP challenge token** | 5 dk | tek kullanımlık | 2FA verify |
| **Password reset token** | 1 saat | tek kullanımlık | Reset endpoint |
| **Admin step-up token** | 5 dk | tek kullanımlık | Kritik admin aksiyon |

## 3. Access Token (JWT)

```json
{
  "sub": "<account_id>",
  "sid": "<session_id>",
  "cid": "<active_character_id>",
  "scope": ["character.read","character.write","chat.send"],
  "iat": 1720000000,
  "exp": 1720000900
}
```

- **Algorithm**: EdDSA (Ed25519) — RS256 alternatif.
- **Key rotation**: 90 gün cycle; overlap 30 gün (dual verify).
- **Kid header**: rotation için key id.
- Server public key sadece verify (private HSM veya AWS KMS Phase 02).

## 4. Refresh Token Rotation

- Refresh token 1 kez kullanılabilir.
- Refresh çağrısı → yeni access + yeni refresh döner; eski refresh iptal.
- **Reuse detection**: aynı refresh token 2. kez → **compromise assumption**; tüm session ailesi (aynı session tree) invalidate + inbox uyarı.

Model:
```sql
CREATE TABLE auth.refresh_token (
  id UUID PRIMARY KEY,
  session_id UUID NOT NULL,
  token_hash BYTEA UNIQUE NOT NULL,
  parent_token_id UUID NULL,        -- rotation chain
  issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  revoke_reason VARCHAR(32) NULL
);
```

## 5. Session Model

```sql
CREATE TABLE auth.session (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL,
  device_id UUID NOT NULL,
  device_fingerprint VARCHAR(128) NULL,   -- ua_hash + screen + tz opt
  ip_hash BYTEA NULL,
  country_code CHAR(2) NULL,              -- geo IP
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_active_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  revoke_reason VARCHAR(32) NULL
);
CREATE INDEX idx_session_account_active ON auth.session (account_id, revoked_at) WHERE revoked_at IS NULL;
```

Ek session state Redis:
- `session:{sid}` → { account_id, active_character_id, scope, last_heartbeat, ws_conn_count }.
- TTL 30 gün; refresh sırasında yenilenir.

## 6. Device Fingerprint

- **device_id**: client persistent (Keychain / EncryptedSharedPreferences). Kullanıcı silmedikçe stabil.
- **fingerprint**: user_agent + screen + tz + OS + app version (SHA-256 hash).
- **New device**: bilinmeyen device_id → OTP challenge zorunlu (login).

Fingerprint değişir (ör. tarayıcı update) → uyarı ama block değil.

## 7. Geo Anomaly

- Login country code değişimi → uyarı (inbox + email).
- Impossibly rapid geo change (10dk içinde İstanbul → Tokyo) → session invalidate + OTP challenge next login.

## 8. Session Revocation

Nedenler:
- Kullanıcı logout.
- Password change → tüm session revoke.
- Admin ban / lock.
- Refresh token reuse detection.
- Session inactivity > 30 gün (expire).
- Device revoke (kullanıcı UI'dan).

Revocation broadcast:
- Redis `session:revoked:{sid}` set 1h TTL — verify middleware çeker.
- WS aktif connection → force disconnect + client error state.

## 9. Password Change Flow

`POST /v1/auth/password/change`:
1. Current password verify.
2. New password strength check (min 8 char, 1 upper, 1 digit önerilir; NIST guideline).
3. Update `password_hash` (Argon2id — parametreler config).
4. **Revoke all sessions** except current.
5. Notification inbox + email "password changed".

## 10. 2FA (TOTP + WebAuthn)

Player:
- TOTP (Google Authenticator vb) — v1.
- WebAuthn (hardware key / platform auth) — Phase 02.
- Recovery codes — 10 tek kullanımlık; 2FA setup sırasında gösterilir.

New device login → 2FA challenge zorunlu (opsiyonel remember-device 30d).

## 11. Password Storage

- **Argon2id** — memory 64MB, iterations 3, parallelism 1 (v1 tune Phase 01).
- Salt per user (16 byte).
- Pepper server-side (Phase 02 opt).
- No plaintext log EVER.

## 12. Fraud Signal

- **Credential stuffing**: 5+ fail farklı account aynı IP 5dk → IP block 1h + captcha (Phase 02).
- **Brute force**: aynı account 5+ fail → account lock 15dk.
- **Session hijack attempt**: refresh reuse veya token fingerprint mismatch.
- **Geo impossibility**: session invalidate.

## 13. Metrikler
- `sec_login_fail_total{reason}`
- `sec_refresh_reuse_detected_total` (P1 alert)
- `sec_session_revoke_total{reason}`
- `sec_geo_anomaly_total`
- `sec_2fa_challenge_total`
- `sec_2fa_fail_total`

## 14. Cross-Reference
- [auth](../auth/README.md), [auth/session.md](../auth/session.md), [auth/device.md](../auth/device.md)
- [replay-protection.md](replay-protection.md), [audit.md](audit.md)

## 15. Bekleyen Kararlar
- Argon2id parametreleri (memory/iterations) Phase 01 tune.
- Recovery code sayısı (10 v1) Phase 02 review.
- WebAuthn player-side rollout Phase 02.
- Captcha sağlayıcı (hCaptcha / Turnstile) infra karar.
- Client integrity attestation Phase 02.
