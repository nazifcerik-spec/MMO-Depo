DROP TABLE IF EXISTS ops.email_outbox;
ALTER TABLE auth.account DROP COLUMN IF EXISTS password_updated_at;
DROP TABLE IF EXISTS auth.account_token;
