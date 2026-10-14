-- Gathering tables (DD-244): an item can be a gathering node of one
-- gathering profession, unlocked at a level (1..100) with a tier (0..7);
-- professions are crafting (learned, recipes) or gathering (levelled by
-- gathering, never learned).
ALTER TABLE content.item_template
    ADD COLUMN gather_kind  VARCHAR(16),
    ADD COLUMN gather_level SMALLINT CHECK (gather_level BETWEEN 1 AND 100),
    ADD COLUMN gather_tier  SMALLINT CHECK (gather_tier BETWEEN 0 AND 7);

CREATE INDEX ix_item_template_gather ON content.item_template (gather_kind, gather_level)
    WHERE gather_kind IS NOT NULL;

ALTER TABLE content.profession
    ADD COLUMN kind VARCHAR(16) NOT NULL DEFAULT 'crafting' CHECK (kind IN ('crafting', 'gathering'));
