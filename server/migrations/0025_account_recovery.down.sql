DROP TABLE IF EXISTS ops.email_outbox;
ALTER TABLE identity.account DROP COLUMN IF EXISTS password_updated_at;
DROP TABLE IF EXISTS identity.account_token;
