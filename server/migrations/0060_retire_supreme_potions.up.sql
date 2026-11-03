-- The supreme healing and mana potions are retired (DD-341). The content
-- loader never deletes, so their recipes go here; potions players already
-- hold become the next tier down (huge). The templates stay for history
-- rows (buybacks, price history) that name them.
DELETE FROM content.recipe WHERE code IN ('recipe_potion_hp_supreme', 'recipe_potion_mp_supreme');
UPDATE inventory.item_instance SET template_code = 'potion_hp_huge' WHERE template_code = 'potion_hp_supreme';
UPDATE inventory.item_instance SET template_code = 'potion_mp_huge' WHERE template_code = 'potion_mp_supreme';
