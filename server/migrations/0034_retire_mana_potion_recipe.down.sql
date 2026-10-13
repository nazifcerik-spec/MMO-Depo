INSERT INTO content.recipe (code, display_name_key, profession_code, min_skill,
                            inputs, gold_cost_copper,
                            output_template_code, output_stack, profession_xp_reward)
VALUES ('recipe_potion_mp_small', 'recipe.potion_mp_small', 'alchemy', 3,
        '[{"template_code":"mandrake","quantity":2}]'::jsonb, 5,
        'potion_mp_small', 1, 20)
ON CONFLICT (code) DO NOTHING;
