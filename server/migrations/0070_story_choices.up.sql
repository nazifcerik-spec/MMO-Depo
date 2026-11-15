-- Story choices (DD-367): one row per passed scene with the chosen option
-- ('' = continued). Flags and counters are replayed from these rows.
CREATE TABLE progression.story_scene (
    character_id UUID        NOT NULL REFERENCES character.character(id) ON DELETE CASCADE,
    scene_code   VARCHAR(64) NOT NULL,
    option_code  VARCHAR(64) NOT NULL DEFAULT '',
    done_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (character_id, scene_code)
);

-- Story quests are offered by the current scene.
ALTER TABLE content.quest_template DROP CONSTRAINT quest_template_category_check;
ALTER TABLE content.quest_template ADD CONSTRAINT quest_template_category_check
    CHECK (category IN ('main','side','daily','story'));
