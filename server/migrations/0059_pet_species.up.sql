-- Pets (DD-334): the first 20 species give way to the 50 of the pet design
-- sheet. Owned pets and eggs move to the new species of the same standing
-- (old rank r -> new rank 1 + (r - 1) x 49/19). The new egg templates are
-- copied from the old rows here so instances can move; the content load then
-- writes their real values.
CREATE TEMP TABLE pet_species_map (old_code VARCHAR(32) PRIMARY KEY, new_code VARCHAR(32) NOT NULL);
INSERT INTO pet_species_map VALUES
  ('ember_fox', 'flamewing'),
  ('moss_turtle', 'stormbird'),
  ('frost_owl', 'sand_drake'),
  ('thunder_hare', 'venom_moth'),
  ('shadow_cat', 'gale_owl'),
  ('crystal_beetle', 'coral_snake'),
  ('storm_hawk', 'cosmic_whale'),
  ('iron_boar', 'silvermoon_hare'),
  ('spirit_wolf', 'ember_hyena'),
  ('sand_wyrmling', 'mushroom_squirrel'),
  ('coral_serpent', 'rainbow_parrot'),
  ('lava_salamander', 'emerald_peacock'),
  ('moon_stag', 'rune_golem'),
  ('void_raven', 'wheat_boar'),
  ('griffin_cub', 'star_octopus'),
  ('frost_drake', 'obsidian_frog'),
  ('celestial_lynx', 'seal_dragon'),
  ('phoenix_chick', 'autumn_stag'),
  ('star_dragon', 'graveyard_owl'),
  ('eternal_kirin', 'chaos_hydra');

INSERT INTO content.item_template (code, display_name_key, item_class, item_subclass, min_level, base_value_copper, base_weight_g, stackable, max_stack, bind_kind, no_trade, no_sell, affixable, stat_bonuses)
SELECT 'pet_egg_' || m.new_code, 'item.pet_egg_' || m.new_code, item_class, item_subclass, min_level,
       base_value_copper, base_weight_g, stackable, max_stack, bind_kind, no_trade, no_sell, affixable, stat_bonuses
  FROM content.item_template t JOIN pet_species_map m ON t.code = 'pet_egg_' || m.old_code
ON CONFLICT (code) DO NOTHING;

UPDATE inventory.item_instance i SET template_code = 'pet_egg_' || m.new_code
  FROM pet_species_map m WHERE i.template_code = 'pet_egg_' || m.old_code;
UPDATE market.listing l SET template_code = 'pet_egg_' || m.new_code
  FROM pet_species_map m WHERE l.template_code = 'pet_egg_' || m.old_code;
UPDATE market.price_history h SET template_code = 'pet_egg_' || m.new_code
  FROM pet_species_map m WHERE h.template_code = 'pet_egg_' || m.old_code;
UPDATE pet.pet p SET species_code = m.new_code
  FROM pet_species_map m WHERE p.species_code = m.old_code;

DELETE FROM content.item_template t USING pet_species_map m WHERE t.code = 'pet_egg_' || m.old_code;
DROP TABLE pet_species_map;
