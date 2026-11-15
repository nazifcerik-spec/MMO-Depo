"""Story quests (DD-367): writes the `quests:` block at the end of
content/story.yaml from the table below. Rewards follow the XP curve: a
dungeon or raid pays half a level at the quest's level, any other quest 30%;
gold is ~25 kills' worth of mob gold (80% / 50%).

Objectives may only name hunted mobs (the 20-level ladder, DD-343); mobs
that live in dungeon waves only are counted as plain kills.

Run from the repository root: python scripts/gen_story_quests.py
"""
import io
import math
import os

import yaml

HERE = os.path.dirname(os.path.abspath(__file__))
STORY = os.path.join(HERE, '..', 'content', 'story.yaml')

# code, level, objectives
Q = [
    ('story_p_bread', 1, [('gather_mining', 10), ('gather_herbalism', 5)]),
    ('story_p_kettle', 1, [('mob_killed:young_boar', 6), ('item_crafted:food_grilled_guppy', 1)]),
    ('story_c1_wolves', 2, [('mob_killed', 8), ('gather_hunting', 3)]),
    ('story_c1_bandits', 5, [('mob_killed', 12)]),
    ('story_c1_crypt', 6, [('dungeon_cleared:crypt_of_ashes', 1)]),
    ('story_c1_prototype', 8, [('gather_woodcutting', 10), ('gather_mining', 10), ('item_crafted:aluminium_bar', 1)]),
    ('story_c1_ledger', 12, [('mob_killed', 15)]),
    ('story_c1_peak', 18, [('mob_killed:shadow_stalker', 6), ('gather_mining', 20)]),
    ('story_c1_spire', 23, [('dungeon_cleared:obsidian_spire', 1)]),
    ('story_c1_marsh', 30, [('mob_killed', 25), ('gather_herbalism', 10), ('item_crafted:potion_hp_small', 1)]),
    ('story_c1_temple', 38, [('dungeon_cleared:sunken_temple', 1)]),
    ('story_c1_buried', 55, [('mob_killed:dwarf_ghost', 10), ('gather_mining', 10)]),
    ('story_c1_mine', 63, [('dungeon_cleared:forgotten_mine', 1)]),
    ('story_c1_anvil', 63, [('item_crafted:iron_bar', 10), ('item_crafted:copper_bar', 5)]),
    ('story_c1_garrison', 80, [('mob_killed:ghost_guard', 10), ('mob_killed', 30)]),
    ('story_c1_varn', 88, [('dungeon_cleared:northern_keep', 1)]),
    ('story_c1_siege', 100, [('raid_cleared:north_keep_siege', 1)]),
    # Chapter II
    ('story_c2_plains', 100, [('mob_killed:plains_lion', 10), ('gather_herbalism', 15)]),
    ('story_c2_orrin', 128, [('dungeon_cleared:rebel_camp', 1)]),
    ('story_c2_burnt', 160, [('mob_killed:ash_wolf', 10), ('mob_killed', 20), ('gather_woodcutting', 15)]),
    ('story_c2_fireroot', 178, [('dungeon_cleared:fire_root', 1)]),
    ('story_c2_rootseal', 178, [('item_crafted', 5)]),
    ('story_c2_desert', 220, [('mob_killed:desert_viper', 10), ('mob_killed', 15), ('gather_herbalism', 15)]),
    ('story_c2_palace', 228, [('dungeon_cleared:buried_palace', 1)]),
    ('story_c2_inner_gate', 250, [('mob_killed:mummy_priest', 10), ('item_crafted', 3)]),
    ('story_c2_crypt', 278, [('dungeon_cleared:royal_crypt', 1)]),
    ('story_c2_tribute', 280, [('mob_killed:royal_skeleton', 15), ('mob_killed:soul_drinker', 10)]),
    ('story_c2_aldric', 295, [('mob_killed:crownless_king', 1)]),
    ('story_c2_tyrant', 300, [('raid_cleared:crownless_palace', 1)]),
    # Chapter III
    ('story_c3_coast', 320, [('mob_killed:storm_bird', 10), ('mob_killed:siren', 6), ('gather_fishing', 20)]),
    ('story_c3_fleet', 328, [('dungeon_cleared:sunken_fleet', 1)]),
    ('story_c3_slopes', 360, [('mob_killed:young_dragon', 10), ('mob_killed', 25)]),
    ('story_c3_brood', 378, [('dungeon_cleared:brood_cavern', 1)]),
    ('story_c3_mist', 420, [('mob_killed:mist_stalker', 10), ('mob_killed:shadow_elf', 8), ('gather_herbalism', 20)]),
    ('story_c3_grove', 428, [('dungeon_cleared:lost_grove', 1)]),
    ('story_c3_caves', 460, [('mob_killed:crystal_spider', 12), ('mob_killed:light_elemental', 10),
                             ('gather_mining', 20), ('item_crafted', 3)]),
    ('story_c3_heart', 478, [('dungeon_cleared:crystal_heart', 1)]),
    ('story_c3_volcano', 520, [('mob_killed:fire_giant', 10), ('mob_killed:volcano_shaman', 10), ('gather_mining', 20)]),
    ('story_c3_ignar', 528, [('dungeon_cleared:volcano_heart', 1)]),
    ('story_c3_peaks', 560, [('mob_killed:dragon_knight', 12), ('mob_killed:dragon_sorcerer', 10),
                             ('gather_woodcutting', 20), ('item_crafted', 1)]),
    ('story_c3_throne', 578, [('dungeon_cleared:dragon_throne', 1)]),
    ('story_c3_lair', 600, [('raid_cleared:dragon_king_lair', 1)]),
    # Chapter IV
    ('story_c4_plains', 610, [('mob_killed:void_wolf', 12), ('mob_killed', 30), ('gather_herbalism', 20)]),
    ('story_c4_observatory', 628, [('dungeon_cleared:darkened_observatory', 1)]),
    ('story_c4_lost_map', 650, [('mob_killed:ghost_noble', 10), ('mob_killed', 10)]),
    ('story_c4_city', 660, [('mob_killed:ghost_noble', 12), ('mob_killed:cursed_sorcerer', 10),
                            ('gather_mining', 20), ('item_crafted', 3)]),
    ('story_c4_tower', 678, [('dungeon_cleared:clock_tower', 1)]),
    ('story_c4_timeless', 700, [('raid_cleared:time_lord_tower', 1)]),
    # Chapter V
    ('story_c5_isles', 710, [('mob_killed:storm_dragon', 10), ('mob_killed:wind_elemental', 10),
                             ('gather_woodcutting', 25), ('gather_mining', 25), ('item_crafted', 2)]),
    ('story_c5_fortress', 728, [('dungeon_cleared:storm_fortress', 1)]),
    ('story_c5_edge', 760, [('mob_killed:void_creature', 15), ('mob_killed:chaos_elemental', 12),
                            ('gather_herbalism', 25), ('item_crafted', 1)]),
    ('story_c5_lullaby', 770, [('mob_killed:chaos_elemental', 4)]),
    ('story_c5_core', 778, [('dungeon_cleared:rift_core', 1)]),
    ('story_c5_temple', 820, [('mob_killed:temple_guardian', 12), ('mob_killed:light_hunter', 10), ('gather_mining', 25)]),
    ('story_c5_dome', 828, [('dungeon_cleared:celestial_dome', 1)]),
    ('story_c5_fall', 850, [('raid_cleared:archangels_fall', 1)]),
    # Chapter VI
    ('story_c6_sea', 860, [('mob_killed:ash_leviathan', 12), ('mob_killed:undying_soldier', 15), ('gather_fishing', 25)]),
    ('story_c6_throne', 878, [('dungeon_cleared:sunken_throne', 1)]),
    ('story_c6_walls', 900, [('mob_killed:eternal_sentinel', 12), ('mob_killed:void_titan', 10),
                             ('gather_mining', 30), ('item_crafted', 3)]),
    ('story_c6_bastion', 928, [('dungeon_cleared:bastion_keep', 1)]),
    ('story_c6_gate', 950, [('mob_killed:void_dragon', 10), ('mob_killed', 40)]),
    ('story_c6_edge', 978, [('dungeon_cleared:eternal_throne', 1)]),
    # Endings
    ('story_e_unbound', 1000, [('raid_cleared:fallen_god_throne', 1)]),
]


def xp_to_next(level):  # leveling: ceil((a + b·L^c) × baseline × (1 + growth·(L−1)))
    with io.open(os.path.join(HERE, '..', 'content', 'progression.yaml'), encoding='utf-8') as f:
        c = yaml.safe_load(f)['xp_curve']
    rate = c['baseline_xp_per_minute'] * (1 + c.get('rate_growth', 0) * (level - 1))
    return math.ceil((c['a'] + c['b'] * level ** c['c']) * rate)


def gold(level):  # ~25 kills' worth of mob gold (gen_world: 503·(L/25)^1.5 a mob)
    return round(25 * 503 * (level / 25) ** 1.5)


def block():
    out = []
    for i, (code, level, obj) in enumerate(Q):
        big = any(e.startswith(('dungeon_cleared', 'raid_cleared')) for e, _ in obj)
        out.append('\n'.join(
            ['  - code: ' + code, '    display_name_key: quest.' + code, '    category: story',
             '    min_level: %d' % level, '    objectives:']
            + ['      - { event: "%s", count: %d }' % o for o in obj]
            + ['    reward_xp: %d' % round(xp_to_next(level) * (0.5 if big else 0.3)),
               '    reward_gold_copper: %d' % round(gold(level) * (0.8 if big else 0.5)),
               '    sort_order: %d' % (2000 + i * 10)]))
    return '\n\n'.join(out) + '\n'


def main():
    with io.open(STORY, encoding='utf-8') as f:
        src = f.read()
    head = src[:src.index('\nquests:\n') + len('\nquests:\n')]
    with io.open(STORY, 'w', encoding='utf-8', newline='\n') as f:
        f.write(head + block())


if __name__ == '__main__':
    main()
