-- Intentionally a no-op (DD-239): moving the tables back into `auth` would
-- fail on Supabase, and 0003/0001 down already drop them from `identity`.
SELECT 1;
