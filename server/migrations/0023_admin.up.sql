-- 0023_admin.up.sql
-- Admin backend (Sprint 26): admin accounts, roles, sessions, audit,
-- player sanctions. docs/admin/*, DD-213 (password + TOTP instead of SSO).

CREATE TABLE ops.admin_user (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email           VARCHAR(320) UNIQUE NOT NULL,
  display_name    VARCHAR(64) NOT NULL DEFAULT '',
  password_hash   TEXT NOT NULL,
  totp_secret     VARCHAR(64) NOT NULL,     -- base32; KMS encryption Phase 02 (DD-213)
  last_totp_step  BIGINT NOT NULL DEFAULT 0, -- replay guard: codes must be newer
  is_active       BOOLEAN NOT NULL DEFAULT true,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  disabled_at     TIMESTAMPTZ NULL,
  version         INT NOT NULL DEFAULT 1
);

CREATE TABLE ops.admin_user_role (
  admin_user_id       UUID NOT NULL REFERENCES ops.admin_user(id) ON DELETE CASCADE,
  role_code           VARCHAR(48) NOT NULL
                      CHECK (role_code IN ('viewer','support_l1','support_l2','game_designer','ops_engineer','admin')),
  granted_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  granted_by_admin_id UUID NULL REFERENCES ops.admin_user(id) ON DELETE SET NULL,
  PRIMARY KEY (admin_user_id, role_code)
);

CREATE TABLE ops.admin_session (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  admin_user_id  UUID NOT NULL REFERENCES ops.admin_user(id) ON DELETE CASCADE,
  token_hash     BYTEA UNIQUE NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at     TIMESTAMPTZ NOT NULL,
  revoked_at     TIMESTAMPTZ NULL,
  stepup_until   TIMESTAMPTZ NULL,
  ip_hash        BYTEA NULL,
  user_agent     VARCHAR(256) NULL
);
CREATE INDEX idx_admin_session_user ON ops.admin_session (admin_user_id) WHERE revoked_at IS NULL;

CREATE TABLE audit.admin_action (
  id                  BIGSERIAL PRIMARY KEY,
  admin_user_id       UUID NOT NULL,
  roles_at_time       TEXT NOT NULL,
  action              VARCHAR(64) NOT NULL,
  target_kind         VARCHAR(24) NULL,
  target_id           UUID NULL,
  before_state        JSONB NULL,
  after_state         JSONB NULL,
  reason              TEXT NULL,
  ip_hash             BYTEA NULL,
  user_agent          VARCHAR(256) NULL,
  requires_stepup     BOOLEAN NOT NULL DEFAULT false,
  stepup_verified_at  TIMESTAMPTZ NULL,
  request_id          VARCHAR(64) NULL,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_admin_action_ts     ON audit.admin_action (created_at DESC);
CREATE INDEX idx_admin_action_actor  ON audit.admin_action (admin_user_id, created_at DESC);
CREATE INDEX idx_admin_action_target ON audit.admin_action (target_kind, target_id, created_at DESC);

CREATE TABLE ops.account_sanction (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id           UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  kind                 VARCHAR(8) NOT NULL CHECK (kind IN ('mute','ban')),
  until                TIMESTAMPTZ NULL,    -- NULL = permanent
  reason               TEXT NOT NULL,
  created_by_admin_id  UUID NOT NULL REFERENCES ops.admin_user(id),
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  revoked_at           TIMESTAMPTZ NULL,
  revoked_by_admin_id  UUID NULL REFERENCES ops.admin_user(id)
);
CREATE INDEX idx_sanction_active ON ops.account_sanction (account_id, kind) WHERE revoked_at IS NULL;
CREATE INDEX idx_sanction_expiry ON ops.account_sanction (until) WHERE revoked_at IS NULL AND until IS NOT NULL;
