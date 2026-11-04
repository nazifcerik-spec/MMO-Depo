-- Guardian shields below level 285 are retired (DD-363): the first promotion
-- moved to level 300 (DD-351), so nobody could ever wear them before outgrowing
-- them. The content loader never deletes, so their recipes go here. Shields
-- players hold, or have listed, become the level 285 shield. The templates stay
-- for history rows that name them; running craft jobs carry their own recipe.
DELETE FROM content.recipe WHERE code IN ('recipe_shield_royal', 'recipe_shield_sunforged', 'recipe_shield_obsidian', 'recipe_shield_sandsteel', 'recipe_shield_runebound', 'recipe_shield_phoenix', 'recipe_shield_dreadforged');
UPDATE inventory.item_instance SET template_code = 'shield_sovereign'
 WHERE template_code IN ('shield_royal', 'shield_sunforged', 'shield_obsidian', 'shield_sandsteel', 'shield_runebound', 'shield_phoenix', 'shield_dreadforged');
UPDATE market.listing SET template_code = 'shield_sovereign'
 WHERE template_code IN ('shield_royal', 'shield_sunforged', 'shield_obsidian', 'shield_sandsteel', 'shield_runebound', 'shield_phoenix', 'shield_dreadforged') AND state IN ('active', 'pending_sale');
