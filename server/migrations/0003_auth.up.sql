-- 0003_auth.up.sql
-- DD-189 uygulaması: JWT EdDSA + refresh rotation + reuse detection.
--
-- Şema:
--   auth.account         — kullanıcı hesabı (email, username, argon2id hash, state)
--   auth.session         — bir cihaz + browser başına bir aktif oturum
--   auth.refresh_token   — rotation chain (parent_token_id), reuse detection için

CREATE TABLE auth.account (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email                 TEXT UNIQUE NOT NULL,
  email_normalized      TEXT UNIQUE NOT NULL,   -- lowercased + trimmed, lookup için
  username              TEXT UNIQUE NOT NULL,
  username_normalized   TEXT UNIQUE NOT NULL,
  password_hash         TEXT NOT NULL,          -- $argon2id$… encoded string
  state                 VARCHAR(16) NOT NULL DEFAULT 'active',  -- active|locked|banned
  email_verified_at     TIMESTAMPTZ NULL,
  last_login_at         TIMESTAMPTZ NULL,
  created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  version               INT NOT NULL DEFAULT 1
);

CREATE INDEX idx_account_state ON auth.account (state) WHERE state <> 'active';

CREATE TABLE auth.session (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id        UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  device_id         UUID NULL,
  ip_hash           BYTEA NULL,             -- SHA-256(client_ip), audit forensic
  user_agent        VARCHAR(256) NULL,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_active_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at        TIMESTAMPTZ NOT NULL,
  revoked_at        TIMESTAMPTZ NULL,
  revoke_reason     VARCHAR(48) NULL
);

CREATE INDEX idx_session_account_active
  ON auth.session (account_id)
  WHERE revoked_at IS NULL;

CREATE INDEX idx_session_expire
  ON auth.session (expires_at)
  WHERE revoked_at IS NULL;

CREATE TABLE auth.refresh_token (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id        UUID NOT NULL REFERENCES auth.session(id) ON DELETE CASCADE,
  token_hash        BYTEA UNIQUE NOT NULL,  -- SHA-256(token) — plaintext DB'de yok
  parent_token_id   UUID NULL REFERENCES auth.refresh_token(id) ON DELETE SET NULL,
  issued_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at        TIMESTAMPTZ NOT NULL,
  used_at           TIMESTAMPTZ NULL,       -- rotation'da işaretlenir; ikinci kullanım = compromise
  revoked_at        TIMESTAMPTZ NULL,
  revoke_reason     VARCHAR(48) NULL
);

CREATE INDEX idx_refresh_session_active
  ON auth.refresh_token (session_id)
  WHERE revoked_at IS NULL AND used_at IS NULL;
