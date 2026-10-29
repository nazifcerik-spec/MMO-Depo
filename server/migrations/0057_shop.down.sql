DROP SCHEMA IF EXISTS cosmetic CASCADE;
DELETE FROM economy.currency WHERE code = 'credit'
  AND NOT EXISTS (SELECT 1 FROM economy.ledger_entry WHERE currency_code = 'credit');
