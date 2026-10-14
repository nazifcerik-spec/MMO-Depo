ALTER TABLE content.profession DROP COLUMN IF EXISTS kind;
DROP INDEX IF EXISTS content.ix_item_template_gather;
ALTER TABLE content.item_template
    DROP COLUMN IF EXISTS gather_tier,
    DROP COLUMN IF EXISTS gather_level,
    DROP COLUMN IF EXISTS gather_kind;
