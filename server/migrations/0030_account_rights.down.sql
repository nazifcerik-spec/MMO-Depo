DELETE FROM auth.account_token WHERE kind = 'data_export';
ALTER TABLE auth.account_token DROP CONSTRAINT account_token_kind_check;
ALTER TABLE auth.account_token ADD CONSTRAINT account_token_kind_check
    CHECK (kind IN ('password_reset', 'email_verify'));

DROP INDEX IF EXISTS auth.ix_account_deletion_due;
ALTER TABLE auth.account
    DROP COLUMN deletion_requested_at,
    DROP COLUMN terms_accepted_at,
    DROP COLUMN terms_version;
