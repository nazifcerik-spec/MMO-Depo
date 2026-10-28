ALTER TABLE social.chat_message DROP COLUMN IF EXISTS system_kind;
DELETE FROM ops.account_sanction WHERE created_by_admin_id IS NULL;
ALTER TABLE ops.account_sanction ALTER COLUMN created_by_admin_id SET NOT NULL;
DROP TABLE IF EXISTS social.profanity_strike;
