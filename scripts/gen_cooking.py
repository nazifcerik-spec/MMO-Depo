"""Cooking content (DD-268): dishes cooked from a fishing catch, a herb and a
log for the fire, one per fishing band (fish level 1..100). Every dish is a
consumable that restores HP — about a third of a same-level mob's HP — and a
recipe makes two. Names in en, tr, zh-CN and es.

Run from the repository root: python scripts/gen_cooking.py
"""
import io
import os
import sys

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

ROOT = os.path.join(os.path.dirname(__file__), '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)

# (dish, fish, fish level, fish tier, herb, log, name)
DISHES = [
    ('grilled_guppy', 'tiny_guppy', 1, 0, 'common_grass', 'young_sapling',
     L('Grilled Guppy', 'Izgara Lepistes', '烤孔雀鱼', 'Gupi a la brasa')),
    ('sardine_skewer', 'river_sardine', 4, 0, 'dandelion', 'poplar_log',
     L('Sardine Skewer', 'Sardalya Şiş', '沙丁鱼串', 'Brocheta de sardinas')),
    ('minnow_soup', 'minnow', 7, 0, 'wild_mint', 'willow_log',
     L('Minty Minnow Soup', 'Naneli Balık Çorbası', '薄荷小鱼汤', 'Sopa de pececillos a la menta')),
    ('pike_fillet', 'small_pike', 10, 1, 'chamomile', 'fir_log',
     L('Pike Fillet', 'Turna Fileto', '梭鱼鱼排', 'Filete de lucio')),
    ('lavender_carp', 'common_carp', 12, 1, 'lavender', 'pine_log',
     L('Lavender Carp', 'Lavantalı Sazan', '薰衣草鲤鱼', 'Carpa a la lavanda')),
    ('herb_trout', 'river_trout', 16, 1, 'bitterroot', 'cedar_log',
     L('Herb-Crusted Trout', 'Otlu Alabalık', '香草烤鳟鱼', 'Trucha en costra de hierbas')),
    ('smoked_salmon', 'salmon', 18, 1, 'healing_moss', 'oak_log',
     L('Smoked Salmon', 'Füme Somon', '烟熏三文鱼', 'Salmón ahumado')),
    ('bass_stew', 'bass', 20, 2, 'marsh_reed', 'beech_log',
     L('Bass Stew', 'Levrek Güveci', '炖鲈鱼', 'Guiso de lubina')),
    ('tuna_steak', 'bluefin_tuna', 28, 2, 'wild_ginseng', 'teak_log',
     L('Ginseng Tuna Steak', 'Ginsengli Ton Balığı', '人参金枪鱼排', 'Filete de atún al ginseng')),
    ('royal_carp_feast', 'royal_carp', 33, 3, 'dreamroot', 'ebony_log',
     L('Royal Carp Feast', 'Kraliyet Sazanı Ziyafeti', '皇家鲤鱼宴', 'Festín de carpa real')),
    ('sturgeon_pie', 'golden_sturgeon', 42, 3, 'ghost_mushroom', 'blackwood_log',
     L('Sturgeon Mushroom Pie', 'Mantarlı Mersin Turtası', '鲟鱼蘑菇派', 'Pastel de esturión y setas')),
    ('swordfish_grill', 'crimson_swordfish', 45, 4, 'stormflower', 'crystal_birch_log',
     L('Storm-Grilled Swordfish', 'Fırtına Izgarası Kılıçbalığı', '风暴烤剑鱼', 'Pez espada a la tormenta')),
    ('frostfin_chowder', 'frostfin_salmon', 55, 4, 'frostbloom', 'frostpine_log',
     L('Frostfin Chowder', 'Ayaz Somon Çorbası', '霜鳍鱼羹', 'Crema de salmón escarcha')),
    ('marlin_platter', 'ghost_marlin', 65, 5, 'spirit_lotus', 'spiritwood_log',
     L('Spirit Marlin Platter', 'Ruh Marlini Tabağı', '灵魂枪鱼拼盘', 'Bandeja de marlín espiritual')),
    ('dragonblood_serpentfish', 'crystal_serpentfish', 75, 6, 'dragonblood_herb', 'starbark_log',
     L('Dragonblood Serpentfish', 'Ejderkanlı Yılanbalığı', '龙血蛇鱼', 'Pez serpiente a la sangre de dragón')),
    ('moonfish_delicacy', 'astral_moonfish', 85, 6, 'void_orchid', 'dragonheart_log',
     L('Astral Moonfish Delicacy', 'Yıldız Aybalığı Lezzeti', '星界月鱼珍馐', 'Manjar de pez luna astral')),
    ('phoenix_koi_feast', 'phoenix_koi', 90, 7, 'celestial_lotus', 'voidwood_log',
     L('Phoenix Koi Feast', 'Anka Koi Ziyafeti', '凤凰锦鲤宴', 'Festín de koi fénix')),
    ('leviathan_banquet', 'ancient_leviathan_scale', 100, 7, 'worldroot', 'worldtree_log',
     L('Leviathan Royal Banquet', 'Leviathan Kraliyet Şöleni', '利维坦皇家盛宴', 'Banquete real de leviatán')),
]

# By fish tier, as smelting pays by ore tier.
XP = [40, 60, 90, 180, 260, 460, 760, 1100]
GOLD = [2, 5, 10, 20, 40, 80, 150, 250]


def round5(x):
    return int(5 * round(x / 5))


def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        s = f.read()
    return s.replace('\r\n', '\n'), '\r\n' in s


def write(path, s, crlf):
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s.replace('\n', '\r\n') if crlf else s)


def block(s, name, body):
    start, end = '  # >>> generated: %s (scripts/gen_cooking.py)\n' % name, '  # <<< generated: %s\n' % name
    b = start + body.rstrip('\n') + '\n' + end
    if start in s:
        return s[:s.index(start)] + b + s[s.index(end) + len(end):]
    return s.rstrip('\n') + '\n\n' + b


def values():
    v = {}
    for f in ('gathering.yaml',):
        import yaml
        for it in yaml.safe_load(io.open(os.path.join(ROOT, f), encoding='utf-8'))['items']:
            v[it['code']] = it['base_value_copper']
    return v


vals = values()
items, recipes, names = [], [], {k: {} for k in ('en', 'tr', 'zh', 'es')}
for dish, fish, level, tier, herb, log, name in DISHES:
    code = 'food_' + dish
    skill = max(1, round(level * 0.9))
    heal = round5(0.35 * (150 + 25 * level))
    value = max(1, round((2 * vals[fish] + vals[herb] + vals[log]) * 1.5 / 2))
    for lang in names:
        names[lang]['item.' + code] = name[lang]
    items.append('\n'.join([
        '  - code: ' + code, '    display_name_key: item.' + code, '    item_class: consumable',
        '    item_subclass: food', '    min_level: %d' % level, '    base_value_copper: %d' % value,
        '    base_weight_g: 200', '    stackable: true', '    max_stack: 50', '    bind_kind: boe',
        '    use_effect: { kind: heal_hp, amount: %d }' % heal]))
    recipes.append('\n'.join([
        '  - code: recipe_cook_' + dish, '    display_name_key: recipe.cook_' + dish, '    profession_code: cooking',
        '    min_skill: %d' % skill, '    inputs:',
        '      - { template_code: %s, quantity: 2 }' % fish,
        '      - { template_code: %s, quantity: 1 }' % herb,
        '      - { template_code: %s, quantity: 1 }' % log,
        '    gold_cost_copper: %d' % GOLD[tier], '    output_template_code: ' + code, '    output_stack: 2',
        '    profession_xp_reward: %d' % XP[tier], '    craft_seconds: %d' % round5(20 + 2 * skill)]))

for fname, body in (('items.yaml', items), ('recipes.yaml', recipes)):
    p = os.path.join(ROOT, fname)
    s, crlf = read(p)
    write(p, block(s, 'cooking', '\n\n'.join(body)), crlf)
merge(names)
print('dishes', len(items))
