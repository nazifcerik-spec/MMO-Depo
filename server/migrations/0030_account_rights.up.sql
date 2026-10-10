-- Account rights (DD-231): terms acceptance at sign-up, self-service account
-- deletion with a 30-day grace, one-time data export links.
ALTER TABLE auth.account
    ADD COLUMN terms_version         VARCHAR(32),
    ADD COLUMN terms_accepted_at     TIMESTAMPTZ,
    ADD COLUMN deletion_requested_at TIMESTAMPTZ;

CREATE INDEX ix_account_deletion_due ON auth.account (deletion_requested_at)
    WHERE deletion_requested_at IS NOT NULL AND state <> 'deleted';

ALTER TABLE auth.account_token DROP CONSTRAINT account_token_kind_check;
ALTER TABLE auth.account_token ADD CONSTRAINT account_token_kind_check
    CHECK (kind IN ('password_reset', 'email_verify', 'data_export'));
