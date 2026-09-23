-- 0004_character.up.sql
-- Character domain — auth'un ilk tüketicisi.
--
-- state:
--   active         — oyunda kullanılabilir.
--   delete_grace   — kullanıcı sildi; 48h içinde cancel edebilir.
--   frozen         — admin freeze (economy drift, forensic).
--
-- Hard delete: worker `character_delete_expire` grace bitince kayıt siler
-- (name UNIQUE index tekrar kullanılabilir olur).

CREATE TABLE character.character (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id        UUID NOT NULL REFERENCES auth.account(id) ON DELETE CASCADE,
  name              TEXT NOT NULL,
  name_normalized   TEXT NOT NULL,          -- lowercase + trimmed, unique lookup
  race_code         VARCHAR(48) NOT NULL,
  class_code        VARCHAR(48) NOT NULL,
  spec_code         VARCHAR(48) NULL,       -- level 10'da seçilir
  level             SMALLINT NOT NULL DEFAULT 1 CHECK (level >= 1 AND level <= 60),
  xp                BIGINT   NOT NULL DEFAULT 0 CHECK (xp >= 0),
  state             VARCHAR(16) NOT NULL DEFAULT 'active',
  location_id       UUID NULL,              -- world/location bind (Faz 01 sprint 12+)
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_played_at    TIMESTAMPTZ NULL,
  delete_grace_at   TIMESTAMPTZ NULL,       -- delete_grace expire zamanı
  version           INT NOT NULL DEFAULT 1
);

-- Name UNIQUE partial: silinen karakter adı serbest kalır (hard delete sonrası).
CREATE UNIQUE INDEX ux_character_name_normalized ON character.character (name_normalized);

-- Account başına aktif karakter listeleme.
CREATE INDEX idx_character_account_active
  ON character.character (account_id)
  WHERE state IN ('active', 'delete_grace');

-- Worker: expire olan delete_grace'leri hard delete için tarar.
CREATE INDEX idx_character_delete_grace_expire
  ON character.character (delete_grace_at)
  WHERE state = 'delete_grace';
