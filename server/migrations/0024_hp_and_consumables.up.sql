-- Sprint 36 (DD-218): persistent HP with analytical regen + consumable effects.
--
-- hp_current NULL means "full" (fresh character, never damaged). Regen is not
-- ticked: readers compute min(max, hp_current + floor(elapsed/2s) × rate)
-- from hp_updated_at, so nothing has to run while the player is away.
ALTER TABLE character.character
  ADD COLUMN hp_current    INT         NULL CHECK (hp_current >= 0),
  ADD COLUMN hp_updated_at TIMESTAMPTZ NULL;

-- Use effect of a consumable template, e.g. {"kind":"heal_hp","amount":60}.
-- NULL = item cannot be used.
ALTER TABLE content.item_template
  ADD COLUMN use_effect JSONB NULL;
