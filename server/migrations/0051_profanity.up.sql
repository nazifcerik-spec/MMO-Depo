-- Chat profanity (DD-312): the swear words an account used lately, system
-- mutes without an admin, and system lines in chat (the mute announcement).
CREATE TABLE social.profanity_strike (
  account_id UUID PRIMARY KEY,
  words      INT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
ALTER TABLE ops.account_sanction ALTER COLUMN created_by_admin_id DROP NOT NULL;
ALTER TABLE social.chat_message ADD COLUMN system_kind VARCHAR(32) NULL;
