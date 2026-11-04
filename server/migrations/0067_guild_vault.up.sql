-- Guild treasury and bank (DD-358). Gold rests in a ledger-backed wallet of the
-- guild, items in the guild's slots; every movement leaves a log line.
CREATE TABLE economy.guild_wallet (
    guild_id      UUID        NOT NULL REFERENCES social.guild(id) ON DELETE CASCADE,
    currency_code VARCHAR(32) NOT NULL REFERENCES economy.currency(code),
    amount        BIGINT      NOT NULL DEFAULT 0 CHECK (amount >= 0),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (guild_id, currency_code)
);

ALTER TABLE inventory.item_instance
  ADD COLUMN guild_bank_id UUID NULL REFERENCES social.guild(id) ON DELETE CASCADE,
  ADD COLUMN guild_bank_slot SMALLINT NULL CHECK (guild_bank_slot BETWEEN 0 AND 99),
  ADD CONSTRAINT item_guild_bank_both CHECK ((guild_bank_id IS NULL) = (guild_bank_slot IS NULL)),
  ADD CONSTRAINT item_guild_bank_alone CHECK (guild_bank_id IS NULL OR (owner_character_id IS NULL AND bank_account_id IS NULL));
CREATE UNIQUE INDEX uq_item_guild_bank_slot ON inventory.item_instance (guild_bank_id, guild_bank_slot)
  WHERE guild_bank_id IS NOT NULL;

CREATE TABLE social.guild_vault_log (
    id             BIGSERIAL PRIMARY KEY,
    guild_id       UUID         NOT NULL REFERENCES social.guild(id) ON DELETE CASCADE,
    character_id   UUID         REFERENCES character.character(id) ON DELETE SET NULL,
    character_name TEXT         NOT NULL,
    action         VARCHAR(16)  NOT NULL CHECK (action IN ('gold_in', 'gold_out', 'item_in', 'item_out')),
    template_code  VARCHAR(64),
    rarity         VARCHAR(16),
    quantity       INT,
    gold_copper    BIGINT,
    at             TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX ix_guild_vault_log ON social.guild_vault_log (guild_id, at DESC);

-- What a member took out today, against their rank's limit.
CREATE TABLE social.guild_withdrawal (
    guild_id     UUID   NOT NULL REFERENCES social.guild(id) ON DELETE CASCADE,
    character_id UUID   NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
    day          DATE   NOT NULL,
    copper       BIGINT NOT NULL DEFAULT 0,
    items        INT    NOT NULL DEFAULT 0,
    PRIMARY KEY (guild_id, character_id, day)
);
