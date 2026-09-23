-- 0020_notification.up.sql
-- Notifications (Sprint 23): device push tokens, push queue, in-game inbox,
-- preferences. docs/notification/push-notification.md, preferences.md.

CREATE SCHEMA IF NOT EXISTS notification;

CREATE TABLE notification.device_push_token (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id      UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  device_id       UUID NOT NULL,
  platform        VARCHAR(16) NOT NULL CHECK (platform IN ('android','ios')),
  push_token      TEXT NOT NULL,
  locale          VARCHAR(16) NULL,
  app_version     VARCHAR(32) NULL,
  is_active       BOOLEAN NOT NULL DEFAULT true,
  registered_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_used_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_failure_at TIMESTAMPTZ NULL,
  failure_count   SMALLINT NOT NULL DEFAULT 0,
  UNIQUE (device_id, platform)
);
CREATE INDEX idx_dpt_account ON notification.device_push_token (account_id) WHERE is_active;

CREATE TABLE notification.push_message (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id     UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  character_id   UUID NULL,
  kind           VARCHAR(48) NOT NULL,
  dedup_key      VARCHAR(128) NOT NULL,
  priority       VARCHAR(16) NOT NULL DEFAULT 'normal' CHECK (priority IN ('high','normal','silent')),
  payload        JSONB NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  scheduled_for  TIMESTAMPTZ NOT NULL DEFAULT now(),
  sent_at        TIMESTAMPTZ NULL,
  devices_ok     SMALLINT NOT NULL DEFAULT 0,
  fail_reason    VARCHAR(64) NULL,
  UNIQUE (account_id, dedup_key)
);
CREATE INDEX idx_push_pending ON notification.push_message (scheduled_for) WHERE sent_at IS NULL;
CREATE INDEX idx_push_rate    ON notification.push_message (account_id, kind, created_at DESC);

CREATE TABLE notification.inbox_message (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id    UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  character_id  UUID NULL,
  kind          VARCHAR(48) NOT NULL,
  dedup_key     VARCHAR(128) NOT NULL,
  body          JSONB NOT NULL,
  route         VARCHAR(256) NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  read_at       TIMESTAMPTZ NULL,
  expires_at    TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30 days',
  UNIQUE (account_id, dedup_key)
);
CREATE INDEX idx_inbox_account ON notification.inbox_message (account_id, created_at DESC);
CREATE INDEX idx_inbox_unread  ON notification.inbox_message (account_id) WHERE read_at IS NULL;

CREATE TABLE notification.preference (
  account_id        UUID PRIMARY KEY REFERENCES auth.account(id) ON DELETE CASCADE,
  quiet_hours_start VARCHAR(5) NULL,
  quiet_hours_end   VARCHAR(5) NULL,
  quiet_hours_tz    VARCHAR(48) NULL,
  global_mute_until TIMESTAMPTZ NULL,
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  version           INT NOT NULL DEFAULT 1
);

CREATE TABLE notification.preference_kind (
  account_id  UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  kind        VARCHAR(48) NOT NULL,
  mode        VARCHAR(16) NOT NULL CHECK (mode IN ('push_and_inbox','inbox_only','silent_push','off')),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (account_id, kind)
);
