# Admin Auth + RBAC

## 1. Amaç
Admin panel giriş, oturum, rol-yetki matrisi, 2FA, step-up authentication (kritik aksiyon).

## 2. Auth Yolu

- **SSO** (OIDC): Google Workspace / Okta / Azure AD entegrasyonu.
- Player account sisteminden **ayrı**; admin.game.internal domain ayrı OIDC client.
- **WebAuthn 2FA** zorunlu — hardware key veya platform authenticator (TouchID/FaceID/Windows Hello).
- Session cookie HttpOnly + SameSite=Strict + Secure; TTL 8 saat.

## 3. Model

```sql
CREATE TABLE ops.admin_user (
  id UUID PRIMARY KEY,
  external_sso_sub VARCHAR(128) UNIQUE NOT NULL,   -- SSO subject
  email VARCHAR(320) UNIQUE NOT NULL,
  display_name VARCHAR(64),
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  disabled_at TIMESTAMPTZ NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE ops.admin_role (
  code VARCHAR(48) PRIMARY KEY,
  display_name VARCHAR(64) NOT NULL,
  permissions JSONB NOT NULL,             -- ['player.read', 'player.mute', ...]
  version INT DEFAULT 1
);

CREATE TABLE ops.admin_user_role (
  admin_user_id UUID NOT NULL,
  role_code VARCHAR(48) NOT NULL,
  granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  granted_by_admin_id UUID NOT NULL,
  expires_at TIMESTAMPTZ NULL,
  PRIMARY KEY (admin_user_id, role_code)
);

CREATE TABLE ops.admin_webauthn_credential (
  id UUID PRIMARY KEY,
  admin_user_id UUID NOT NULL,
  credential_id BYTEA UNIQUE NOT NULL,
  public_key BYTEA NOT NULL,
  aaguid UUID NULL,
  registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_used_at TIMESTAMPTZ NULL
);

CREATE TABLE audit.admin_action (
  id BIGSERIAL PRIMARY KEY,
  admin_user_id UUID NOT NULL,
  role_at_time VARCHAR(48) NOT NULL,
  action VARCHAR(64) NOT NULL,             -- 'player.mute', 'ledger.reversal', ...
  target_kind VARCHAR(24) NULL,
  target_id UUID NULL,
  before_state JSONB NULL,
  after_state JSONB NULL,
  reason TEXT NULL,
  ip_hash BYTEA NULL,
  user_agent VARCHAR(256) NULL,
  requires_stepup BOOLEAN NOT NULL DEFAULT false,
  stepup_verified_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_admin_action_ts ON audit.admin_action (created_at DESC);
CREATE INDEX idx_admin_action_actor ON audit.admin_action (admin_user_id, created_at DESC);
CREATE INDEX idx_admin_action_target ON audit.admin_action (target_kind, target_id, created_at DESC);
```

## 4. Role Matrix (v1)

| Permission | viewer | support_l1 | support_l2 | game_designer | ops_engineer | admin |
|---|---|---|---|---|---|---|
| `player.read` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `player.pii_reveal` | ✗ | ✓ (audit) | ✓ | ✗ | ✓ | ✓ |
| `player.mute<=24h` | ✗ | ✓ | ✓ | ✗ | ✗ | ✓ |
| `player.ban<=7d` | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| `player.ban_permanent` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| `character.read` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `character.freeze` | ✗ | ✗ | ✓ | ✗ | ✓ | ✓ |
| `character.item_restore` | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| `ledger.read` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `ledger.reversal<=100k_gold` | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| `ledger.reversal_any` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| `market.listing_cancel` | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| `guild.disband_override` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| `content.event_trigger` | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ |
| `content.pack_preview` | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ |
| `ops.config_read` | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ |
| `ops.config_write` | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ |
| `ops.feature_flag` | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ |
| `ops.worker_restart` | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ |
| `admin.role_grant` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |

## 5. Step-Up Authentication (Kritik Aksiyon)

Belirli aksiyonlar (reversal, ban_permanent, disband_override, config_write) SSO oturumundan **bağımsız** WebAuthn re-verification ister:
- Aksiyon başlangıcında `POST /internal/admin/stepup/challenge` → challenge.
- Client WebAuthn assertion → `POST /internal/admin/stepup/verify`.
- Verify success → `stepup_token` (TTL 5 dk); aynı session içindeki aksiyon token ile.

`audit.admin_action.requires_stepup=true, stepup_verified_at` zorunlu.

## 6. Session

- SSO callback → server session cookie (HttpOnly).
- Session data Redis (`admin:session:{sid}` TTL 8 saat).
- Logout → session invalidate + audit.

## 7. IP Allowlist / VPN

Production admin.game.internal:
- VPN gate zorunlu (WireGuard veya corp VPN).
- Alternatif: IP allowlist (statik IP).
- Development environment daha esnek (dev VPN opt).

## 8. Rate Limit

- Admin API: 300 req/dk per admin_user_id.
- Aşım → 429 + audit + alert to Slack.

## 9. Endpoint

**`GET /internal/admin/session/me`** — current admin + roles + permissions.
**`POST /internal/admin/logout`**
**`GET /internal/admin/users`** — admin list (admin role only).
**`POST /internal/admin/users/{id}/role`** — role grant.
**`POST /internal/admin/stepup/challenge`** / **`verify`**.

## 10. Metrikler
- `admin_login_total{sso_provider}`
- `admin_login_fail_total{reason}`
- `admin_stepup_challenge_total`
- `admin_stepup_fail_total`
- `admin_permission_denied_total{permission}`

## 11. Cross-Reference
- [database/audit.md](../database/audit.md)
- [auth](../auth/README.md) — player auth (ayrı)

## 12. Bekleyen Kararlar
- SSO provider (Google/Okta/Azure) ops karar.
- Role approval flow (dual approval kritik role grant için) Phase 02.
- Session inactivity timeout (v1: 8h fixed; Phase 02 rolling).
- IP allowlist yaklaşımı (VPN vs statik) infra team karar.
