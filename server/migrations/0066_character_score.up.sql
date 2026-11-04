-- Character score and ranking tables (DD-357). A worker recomputes the rows;
-- reads never aggregate the game tables.
CREATE TABLE progression.character_score (
    character_id       UUID PRIMARY KEY REFERENCES character.character(id) ON DELETE CASCADE,
    name               TEXT        NOT NULL,
    class_code         VARCHAR(32) NOT NULL,
    spec_code          VARCHAR(48),
    level              SMALLINT    NOT NULL,
    achievement_points INT         NOT NULL DEFAULT 0,
    mastery_tiers      INT         NOT NULL DEFAULT 0,
    profession_levels  INT         NOT NULL DEFAULT 0,
    dungeons_cleared   INT         NOT NULL DEFAULT 0,
    raids_cleared      INT         NOT NULL DEFAULT 0,
    score              BIGINT      NOT NULL DEFAULT 0,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ix_score_score        ON progression.character_score (score DESC, character_id);
CREATE INDEX ix_score_level        ON progression.character_score (level DESC, character_id);
CREATE INDEX ix_score_achievements ON progression.character_score (achievement_points DESC, character_id);
CREATE INDEX ix_score_mastery      ON progression.character_score (mastery_tiers DESC, character_id);
CREATE INDEX ix_score_professions  ON progression.character_score (profession_levels DESC, character_id);
CREATE INDEX ix_score_class        ON progression.character_score (class_code, score DESC, character_id);
