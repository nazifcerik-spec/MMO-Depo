-- Raid enemies' strength per raid (DD-307): HP and attack in percent, on top
-- of the field mob scaling and the group-size table. Solo dungeons stay 100.
ALTER TABLE content.dungeon_template
  ADD COLUMN power_pct SMALLINT NOT NULL DEFAULT 100 CHECK (power_pct BETWEEN 25 AND 400);
