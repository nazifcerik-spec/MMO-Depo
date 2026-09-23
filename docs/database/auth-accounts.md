# Database — `auth` Schema (Users & Accounts)

## Kapsam
Hesap sistemi, oturum yönetimi, kimlik doğrulama. **User** = kişi; **Account** = login credential + oyun-hesabı bağlam; **Character** ise `character` schema'sında.

`User (kişi) → Account (oyun-hesabı) → Character (avatar)` üç katmanlı ayrım.

---

## 1. `auth.user`
Kişi. Bir kişi zamanla birden fazla hesap işletebilir (v1: 1:1 varsayımı, ancak schema esnek).

| Kolon | Tip | Constraint | Not |
|---|---|---|---|
| `id` | `UUID` | PK, default `gen_random_uuid()` | |
| `email` | `CITEXT` | UNIQUE, NOT NULL | Case-insensitive |
| `email_verified_at` | `TIMESTAMPTZ` | NULL | Verification akışı |
| `password_hash` | `TEXT` | NOT NULL | argon2id |
| `password_updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() | Force rotation için |
| `display_name` | `VARCHAR(48)` | NOT NULL | Public görünmez; sadece account setting |
| `locale` | `VARCHAR(8)` | NOT NULL DEFAULT `'tr-TR'` | i18n |
| `country_code` | `CHAR(2)` | NULL | ISO 3166 |
| `mfa_enabled` | `BOOLEAN` | NOT NULL DEFAULT false | |
| `mfa_secret` | `TEXT` | NULL | Encrypted at rest (app-level) |
| `status` | `auth.user_status` | NOT NULL DEFAULT `'active'` | enum: active, disabled, banned, deleted |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() | |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() | |
| `last_login_at` | `TIMESTAMPTZ` | NULL | |
| `deleted_at` | `TIMESTAMPTZ` | NULL | Soft delete |

Enum:
```sql
CREATE TYPE auth.user_status AS ENUM ('active','disabled','banned','deleted');
```

Indexler:
- `uq_user__email` (UNIQUE, otomatik CITEXT constraint).
- `ix_user__status` (partial `WHERE status <> 'active'`).
- `ix_user__deleted_at` (partial `WHERE deleted_at IS NOT NULL`).

---

## 2. `auth.account`
Oyun-hesabı katmanı. v1'de user ile 1:1 varsayılır; ancak schema, bir kişinin (user) birden fazla oyun-hesabı olma ihtimalini destekler (aile hesabı, senaryo test).

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK, default `gen_random_uuid()` |
| `user_id` | `UUID` | FK → `auth.user(id)` ON DELETE RESTRICT |
| `region` | `VARCHAR(8)` | NOT NULL DEFAULT `'global'` — market/region routing için |
| `character_slots` | `SMALLINT` | NOT NULL DEFAULT 2 — bkz. GDD monetization |
| `bag_slot_bonus` | `SMALLINT` | NOT NULL DEFAULT 0 |
| `bank_tab_purchased` | `SMALLINT` | NOT NULL DEFAULT 0 |
| `subscription_status` | `auth.subscription_status` | NOT NULL DEFAULT `'none'` — Phase 02 |
| `notification_prefs` | `JSONB` | NOT NULL DEFAULT `'{}'` — kanal/kategori pref (bkz. Notification System) |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `deleted_at` | `TIMESTAMPTZ` | NULL |

Enum:
```sql
CREATE TYPE auth.subscription_status AS ENUM ('none','trial','active','expired');
```

Indexler:
- `ix_account__user_id`.
- `ix_account__region`.

---

## 3. `auth.session`
Aktif oturumlar (JWT dolaşımı). Refresh token'lar burada tutulur; access token stateless.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK — session id (jti) |
| `account_id` | `UUID` | FK → `auth.account(id)` ON DELETE CASCADE |
| `refresh_token_hash` | `TEXT` | NOT NULL — sha256 |
| `client_kind` | `auth.client_kind` | NOT NULL — enum: web, ios, android, desktop, cli |
| `client_version` | `VARCHAR(32)` | NOT NULL |
| `ip_hash` | `TEXT` | NULL — sha256, PII minimizasyon |
| `user_agent` | `TEXT` | NULL — truncated (255 char) |
| `issued_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `expires_at` | `TIMESTAMPTZ` | NOT NULL |
| `revoked_at` | `TIMESTAMPTZ` | NULL |
| `last_seen_at` | `TIMESTAMPTZ` | NULL |

Enum:
```sql
CREATE TYPE auth.client_kind AS ENUM ('web','ios','android','desktop','cli');
```

Indexler:
- `ix_session__account_id`.
- `uq_session__refresh_token_hash` (UNIQUE).
- `ix_session__expires_at` (cleanup job).

Notlar:
- Access token stateless (JWT); revocation Redis blocklist (`sess:blk:<jti>`) + bu tabloda `revoked_at` mark.
- Session rotation: refresh her kullanımda yeni token + eski revoke.

---

## 4. `auth.jwt_key`
JWT imza key'lerinin versiyonu; rotation destekli.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `VARCHAR(32)` | PK — key id (kid) |
| `algorithm` | `VARCHAR(16)` | NOT NULL — `EdDSA` veya `RS256` |
| `public_key_pem` | `TEXT` | NOT NULL |
| `private_key_ref` | `TEXT` | NOT NULL — secret manager reference (URI değil private key) |
| `status` | `auth.key_status` | NOT NULL DEFAULT `'active'` |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `activated_at` | `TIMESTAMPTZ` | NULL |
| `retired_at` | `TIMESTAMPTZ` | NULL |

Enum:
```sql
CREATE TYPE auth.key_status AS ENUM ('pending','active','retired');
```

Notlar: JWT `kid` header alanında bu id; validation kütüphanesi bu tabloyu (in-memory cache) okur.

---

## 5. `auth.password_reset`
Şifre sıfırlama tokenları.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `user_id` | `UUID` | FK → `auth.user(id)` |
| `token_hash` | `TEXT` | NOT NULL UNIQUE — sha256 |
| `expires_at` | `TIMESTAMPTZ` | NOT NULL — 1 saat |
| `used_at` | `TIMESTAMPTZ` | NULL — bir kez kullanılır |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

Indexler:
- `ix_password_reset__user_id`.
- `ix_password_reset__expires_at` (cleanup).

---

## 6. `auth.email_verification`
E-posta doğrulama tokenları — password_reset paterni.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `user_id` | `UUID` FK |
| `email` | `CITEXT` NOT NULL — verify edilen |
| `token_hash` | `TEXT` UNIQUE |
| `expires_at` | `TIMESTAMPTZ` NOT NULL |
| `verified_at` | `TIMESTAMPTZ` NULL |
| `created_at` | `TIMESTAMPTZ` DEFAULT now() |

---

## 7. `auth.rate_limit_bucket` (opsiyonel; alternatif Redis-only)
Not: Ana rate limit **Redis** üstünde. Bu tablo yalnız audit / historical trend için (Phase 02 aday). v1'de skip.

---

## 8. Kullanım Örneği

**Login akışı:**
```
1. POST /auth/login {email, password}
2. Load user by email (CITEXT)
3. Verify password_hash (argon2id)
4. If mfa_enabled → challenge
5. Create session row (refresh_token_hash, expires_at = now + AUTH_REFRESH_TTL)
6. Return {access_token (JWT), refresh_token}
7. Update user.last_login_at
```

**Refresh akışı:**
```
1. POST /auth/refresh {refresh_token}
2. Hash → lookup auth.session
3. If expired or revoked → 401
4. Rotate: revoke old, insert new
5. Return new {access_token, refresh_token}
```

**Logout:**
```
1. Revoke session (revoked_at = now)
2. Add jti to Redis blocklist until access_token TTL expires
```

---

## 9. İlişkiler
- `auth.user` (1) → `auth.account` (0..N — v1'de 1)
- `auth.account` (1) → `character.character` (0..N, cap = `character_slots`)
- `auth.account` (1) → `auth.session` (0..N)
- `auth.user` (1) → `auth.password_reset` / `auth.email_verification` (0..N)

---

## 10. Notlar
- **PII minimizasyon**: email hash'lenmez (login için gerekli), ancak log'a asla; IP hash'li tutulur.
- **Password**: argon2id parametreleri (m=64MB, t=3, p=1) config'te.
- **Account deletion**: soft delete (`deleted_at`); 30 gün sonra hard delete job (`audit.event_log`'a "user_hard_deleted").
- **Ban**: `user.status = 'banned'`; ilgili account/character login reddedilir; ancak veriler korunur.
- **GDPR**: right-to-be-forgotten → PII anonimize (email → hashed uuid), oyun içi karakter kayıtları sürekli anon olarak korunur (ekonomi bütünlüğü için).

---

## 11. Bekleyen Kararlar
- MFA türü: TOTP zorunlu mu, WebAuthn opsiyon mu.
- Sosyal login (Google/Discord) — Phase 02.
- Multi-user-per-account (aile hesabı) modeli — Phase 02.
- Region field'ı gerçek kullanıma girecek mi (v1'de sabit `global`).
