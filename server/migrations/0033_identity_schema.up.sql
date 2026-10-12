-- Account tables live in `identity`, not `auth` (DD-239): Supabase reserves
-- the `auth` schema for its own auth service and does not let us create
-- objects there. New databases already create them in `identity` (0001,
-- 0003 and later migrations); this moves the tables of a database migrated
-- before the rename. Indexes, constraints and foreign keys move with them.
-- The `auth` schema itself is left alone: on Supabase it is not ours.
CREATE SCHEMA IF NOT EXISTS identity;

DO $$
DECLARE
    t text;
BEGIN
    FOREACH t IN ARRAY ARRAY['account', 'session', 'refresh_token', 'account_token'] LOOP
        IF to_regclass('auth.' || t) IS NOT NULL AND to_regclass('identity.' || t) IS NULL THEN
            EXECUTE format('ALTER TABLE auth.%I SET SCHEMA identity', t);
        END IF;
    END LOOP;
END $$;
