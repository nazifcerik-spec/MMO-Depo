-- Mana refills at the start of every fight, so mana potions have no effect
-- yet (DD-241). Retire the recipe until mana persists; the item template
-- stays (existing stacks remain sellable). The content loader never
-- deletes, so retirement is a migration (contentpack design note).
DELETE FROM content.recipe WHERE code = 'recipe_potion_mp_small';
