-- 0062_mob_dungeon_only.up.sql
-- Hunting ladder (DD-343): hunted mobs stand 20 levels apart; the others
-- are fought only in dungeons and raids (waves), so hunts, pet hunts and the
-- hunt list leave them out.
ALTER TABLE content.mob_template ADD COLUMN dungeon_only BOOLEAN NOT NULL DEFAULT false;

-- Bounties on mobs that left the ladder are retired (their kills can no
-- longer be counted); each zone's bounty now names its last hunted mob.
DELETE FROM progression.character_quest WHERE quest_code IN (
  'side_bounty_drowned_knight', 'side_bounty_mine_warden', 'side_bounty_black_lancer',
  'side_bounty_sand_worm', 'side_bounty_kraken_spawn', 'side_bounty_forest_spirit',
  'side_bounty_volcano_warden', 'side_bounty_night_giant', 'side_bounty_thunderbird',
  'side_bounty_seraph', 'side_bounty_last_commander');
DELETE FROM content.quest_template WHERE code IN (
  'side_bounty_drowned_knight', 'side_bounty_mine_warden', 'side_bounty_black_lancer',
  'side_bounty_sand_worm', 'side_bounty_kraken_spawn', 'side_bounty_forest_spirit',
  'side_bounty_volcano_warden', 'side_bounty_night_giant', 'side_bounty_thunderbird',
  'side_bounty_seraph', 'side_bounty_last_commander');
