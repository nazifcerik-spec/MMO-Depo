-- 0011_combat.up.sql
-- Combat Engine — mob template catalog for Sprint 13.
--
-- Mob templates yaşam boyu değişmez (content pack ile deploy). Sprint 14
-- hunting bu tablodan seçilen mob'u sahaya rulo eder; Sprint 15 dungeon
-- boss encounter'ları da bu tabloyu ölçekleyerek kullanır.
--
-- v1 kapsamı:
--   - level + stat bloğu
--   - reward: xp + gold aralığı (min, max copper)
-- v1 dışı (Phase 02):
--   - special abilities (dot, stun, aoe)
--   - loot table (item drop) — v1'de reveal path activity üzerinden
--   - resistance table (fire/frost/holy)
--   - faction

CREATE TABLE content.mob_template (
  code             VARCHAR(48) PRIMARY KEY,
  display_name_key VARCHAR(64) NOT NULL,
  level            SMALLINT NOT NULL DEFAULT 1,

  -- Stat block (mirrors combat.Combatant columns; INT to leave headroom for
  -- boss tuning).
  hp_max           INT NOT NULL,
  attack           INT NOT NULL,
  defense          INT NOT NULL DEFAULT 0,
  crit_pct         SMALLINT NOT NULL DEFAULT 5,
  hit_pct          SMALLINT NOT NULL DEFAULT 85,
  evasion_pct      SMALLINT NOT NULL DEFAULT 5,

  -- Reward on kill. Range keeps loot from being fully predictable but stays
  -- deterministic given the encounter's RNG seed.
  xp_reward        INT NOT NULL DEFAULT 0,
  gold_copper_min  BIGINT NOT NULL DEFAULT 0,
  gold_copper_max  BIGINT NOT NULL DEFAULT 0,

  version          INT NOT NULL DEFAULT 1,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_mob_template_level ON content.mob_template (level);

-- Baseline seed (12 mob spanning level 1..25) — content pack YAML in Sprint 12
-- is the authoritative source; this seed only guarantees a fresh DB is
-- playable if the loader has not run yet.
INSERT INTO content.mob_template (
  code, display_name_key, level, hp_max, attack, defense, crit_pct, hit_pct, evasion_pct,
  xp_reward, gold_copper_min, gold_copper_max
) VALUES
  ('young_boar',    'mob.young_boar',    1,   40,  6, 1,  3, 80, 3,  12,   5,  15),
  ('grey_wolf',     'mob.grey_wolf',     3,   60,  9, 2,  5, 82, 6,  22,  10,  25),
  ('bandit_scout',  'mob.bandit_scout',  5,   90, 14, 3,  8, 85, 10, 45,  20,  60),
  ('cave_bat',      'mob.cave_bat',      6,   50, 12, 1, 12, 82, 20, 40,  15,  40),
  ('goblin_raider', 'mob.goblin_raider', 8,  120, 18, 5,  7, 84, 8,  70,  40, 100),
  ('forest_troll',  'mob.forest_troll',  12, 260, 24, 10, 5, 82, 4, 160,  90, 220),
  ('dire_wolf',     'mob.dire_wolf',     14, 220, 30, 6, 15, 86, 12, 200, 120, 260),
  ('cultist_acolyte','mob.cultist_acolyte',15, 200, 34, 5,  8, 85, 8, 240, 150, 320),
  ('rock_ogre',     'mob.rock_ogre',     18, 480, 42, 18, 5, 80, 3, 380, 250, 500),
  ('shadow_stalker','mob.shadow_stalker',20, 350, 55, 8, 22, 88, 25, 500, 400, 700),
  ('elder_wyrm',    'mob.elder_wyrm',    24, 800, 65, 22, 12, 86, 8, 900,1000,2000),
  ('ashen_lich',    'mob.ashen_lich',    25, 700, 75, 12, 20, 90, 15, 1200,1500, 3000);
