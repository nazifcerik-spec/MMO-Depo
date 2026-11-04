DROP TABLE social.guild_withdrawal;
DROP TABLE social.guild_vault_log;
DELETE FROM inventory.item_instance WHERE guild_bank_id IS NOT NULL;
ALTER TABLE inventory.item_instance
  DROP CONSTRAINT item_guild_bank_alone,
  DROP CONSTRAINT item_guild_bank_both,
  DROP COLUMN guild_bank_slot,
  DROP COLUMN guild_bank_id;
DROP TABLE economy.guild_wallet;
