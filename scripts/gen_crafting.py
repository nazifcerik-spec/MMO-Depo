"""Gear and potion crafting over the whole ladder (DD-274).

Profession levels 1-100 stand for character levels 1-1000 (K1): a recipe for
level-L gear needs skill ceil(L/10) and the materials gathered at about that
level. Every piece of random-drop gear gets a recipe:

  blacksmithing   plate armor + metal weapons   bars (+ hide straps, log grips)
  leatherworking  leather armor                 hides and pelts (+ herb tanning)
  tailoring       cloth armor + holy tomes      herbs as fibre (+ hide lining)
  woodworking     staves, bows, lutes, totems   logs (+ bar fittings, feathers)

Healing potions heal a share of max HP (heal_pct) and mana potions refill a
share of a hunter's carried pool (mana_pct, DD-288); both spread over the
alchemy levels. Writes generated blocks into content/recipes.yaml and
content/items.yaml and names in en, tr, zh-CN and es.

Run from the repository root after gen_gear_types.py: python scripts/gen_crafting.py
"""
import io
import math
import os
import sys

import yaml

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)

PROFESSIONS = {
    'leatherworking': L('Leatherworking', 'Deri İşçiliği', '制皮', 'Peletería'),
    'tailoring': L('Tailoring', 'Terzilik', '裁缝', 'Sastrería'),
    'woodworking': L('Woodworking', 'Ağaç İşçiliği', '木工', 'Carpintería'),
    'jewelcrafting': L('Jewelcrafting', 'Kuyumculuk', '珠宝加工', 'Joyería'),
}
METAL = {'plate_armor', 'sword', 'axe', 'mace', 'two_handed_mace', 'dagger', 'short_sword', 'rapier',
         'fist_weapon', 'sickle', 'ritual_axe', 'shield'}
WOOD = {'staff', 'wand', 'bow', 'crossbow', 'lute', 'totem_focus', 'nature_focus', 'tonfa'}
FLETCHED = {'bow', 'crossbow'}
JEWELRY = {'necklace', 'ring', 'trinket'}  # DD-276
# Material count and XP/time weight per slot.
PIECE = {'chest': (5, 1.6), 'legs': (4, 1.4), 'helm': (3, 1.0), 'boots': (3, 1.0), 'gloves': (2, 0.8),
         'mainhand': (4, 1.6), 'offhand': (4, 1.4), 'necklace': (3, 1.2), 'ring': (2, 1.0), 'trinket': (2, 1.0)}

# Healing potions (heal_pct): code, level, percent of max HP, value, names.
POTIONS = [
    ('potion_hp_small', 1, 20, 20, None),
    ('potion_hp_medium', 100, 30, 60, None),
    ('potion_hp_large', 250, 40, 150, None),
    ('potion_hp_xlarge', 450, 55, 320, None),
    ('potion_hp_huge', 650, 70, 600, None),
    ('potion_hp_supreme', 850, 85, 1000,
     L('Supreme Healing Potion', 'Yüce Can İksiri', '至高治疗药水', 'Poción de curación suprema')),
]
# Mana potions (mana_pct): code, level, percent of the pool, value, names.
MANA_POTIONS = [
    ('potion_mp_small', 1, 30, 20, L('Small Mana Potion', 'Küçük Mana İksiri', '小型法力药水', 'Poción de maná pequeña')),
    ('potion_mp_medium', 100, 40, 60, L('Medium Mana Potion', 'Orta Boy Mana İksiri', '中型法力药水', 'Poción de maná mediana')),
    ('potion_mp_large', 250, 50, 150, L('Large Mana Potion', 'Büyük Boy Mana İksiri', '大型法力药水', 'Poción de maná grande')),
    ('potion_mp_xlarge', 450, 60, 320,
     L('Extra Large Mana Potion', 'Ekstra Büyük Mana İksiri', '特大法力药水', 'Poción de maná extragrande')),
    ('potion_mp_huge', 650, 70, 600, L('Huge Mana Potion', 'Dev Mana İksiri', '巨型法力药水', 'Poción de maná enorme')),
    ('potion_mp_supreme', 850, 85, 1000, L('Supreme Mana Potion', 'Yüce Mana İksiri', '至高法力药水', 'Poción de maná suprema')),
]


def skill_for(level):
    return max(1, math.ceil(level / 10))


def xp_to_next(level):  # crafting.ProfessionXPToNext
    return round(60 * 1.06 ** max(1, level))


def round5(x):
    return int(5 * round(x / 5))


def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        raw = f.read()
    return raw.replace('\r\n', '\n'), '\r\n' in raw


def write(path, s, crlf):
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s.replace('\n', '\r\n') if crlf else s)


def put_block(s, name, body, indent='  '):
    start = '%s# >>> generated: %s (scripts/gen_crafting.py)\n' % (indent, name)
    end = '%s# <<< generated: %s\n' % (indent, name)
    block = start + body.rstrip('\n') + '\n' + end
    if start in s:
        return s[:s.index(start)] + block + s[s.index(end) + len(end):]
    return s.rstrip('\n') + '\n\n' + block


def load(name):
    return yaml.safe_load(read(os.path.join(ROOT, name))[0])


nodes = load('gathering.yaml')['items']


def ladder(*subclasses):
    """(level, code) of the gathering nodes of these subclasses, by level."""
    return sorted((n['gather']['level'], n['code']) for n in nodes if n['item_subclass'] in subclasses)


def pick(lad, g):
    """The best material gathered at or below level g; the first one when none is."""
    best = [c for lvl, c in lad if lvl <= g]
    return best[-1] if best else lad[0][1]


def runner_up(lad, g):
    """The material just below pick(lad, g): metal pieces alloy the tier's bar
    with the one before it, so every bar the smelter makes has a use."""
    best = [c for lvl, c in lad if lvl <= g]
    return best[-2] if len(best) > 1 else pick(lad, g)


ores = {n['code'][:-len('_ore')]: n['gather']['level'] for n in nodes if n['item_subclass'] == 'ore'}
bars = sorted((ores[i['code'][:-len('_bar')]], i['code']) for i in load('items.yaml')['items']
              if i.get('item_subclass') == 'bar' and i['code'][:-len('_bar')] in ores)
LOGS, HERBS, FEATHERS = ladder('log'), ladder('herb'), ladder('feather')
HIDES = ladder('pelt', 'hide', 'scale')


def recipe(out, profession, skill, inputs, gold, xp, seconds, stack=1):
    lines = ['  - code: recipe_' + out, '    display_name_key: recipe.' + out, '    profession_code: ' + profession,
             '    min_skill: %d' % skill, '    inputs:']
    lines += ['      - { template_code: %s, quantity: %d }' % (c, q) for c, q in inputs]
    lines += ['    gold_cost_copper: %d' % gold, '    output_template_code: ' + out, '    output_stack: %d' % stack,
              '    profession_xp_reward: %d' % xp, '    craft_seconds: %d' % seconds]
    return '\n'.join(lines)


def merged(inputs):
    """Sum repeated materials (low tiers pick the same node twice)."""
    out = {}
    for c, q in inputs:
        out[c] = out.get(c, 0) + q
    return list(out.items())


gear_recipes = []
STONES = ladder('stone')
# Branch-gated gear (guardian shields) is crafted too; attribute-gated
# showpieces (the Titan Bulwark) are not.
gear = [g for g in load('gear.yaml')['items'] if not set(g.get('requirements') or {}) - {'spec'}]
for g in sorted(gear, key=lambda g: (g['min_level'], g['item_subclass'], g['code'])):
    sub, slot, lvl = g['item_subclass'], g['equip_slot'], g['min_level']
    skill = skill_for(lvl)
    q, w = PIECE[slot]
    # Materials one level ahead of the skill, so the top tier (skill 99) uses
    # the level-100 nodes.
    reach = min(100, skill + 1)
    bar, log, hide, herb = pick(bars, reach), pick(LOGS, reach), pick(HIDES, reach), pick(HERBS, reach)
    if sub in METAL:
        prof = 'blacksmithing'
        alloy = [(bar, q - 1), (runner_up(bars, reach), 1)]
        inputs = alloy + ([(hide, 1)] if slot != 'mainhand' else [(log, 2)])
    elif sub == 'leather_armor':
        prof, inputs = 'leatherworking', [(hide, q), (herb, 1)]
    elif sub in ('cloth_armor', 'holy_tome'):
        prof, inputs = 'tailoring', [(herb, 2 * q), (hide, 1)]
    elif sub in WOOD:
        prof, inputs = 'woodworking', [(log, q), (bar, 1)]
        if sub in FLETCHED:
            inputs.append((pick(FEATHERS, reach), 2))
    elif sub in JEWELRY:
        prof = 'jewelcrafting'
        inputs = [(bar, q), (runner_up(bars, reach), 1), (pick(STONES, reach), 2)]
    else:
        raise SystemExit('no profession crafts %s (%s)' % (g['code'], sub))
    gear_recipes.append(recipe(g['code'], prof, skill, merged(inputs), max(1, round(g['base_value_copper'] * 0.1)),
                               max(10, round(0.25 * xp_to_next(skill) * w)), round5((20 + 1.6 * skill) * w)))

potion_items, potion_recipes, names = [], [], {lang: {} for lang in ('en', 'tr', 'zh', 'es')}
for code, lvl, pct, value, title in POTIONS:
    skill = skill_for(lvl)
    potion_items.append('\n'.join([
        '  - code: ' + code, '    display_name_key: item.' + code, '    item_class: consumable',
        '    item_subclass: potion', '    min_level: %d' % lvl, '    base_value_copper: %d' % value,
        '    base_weight_g: 250', '    stackable: true', '    max_stack: 20', '    bind_kind: boe',
        '    use_effect: { kind: heal_pct, amount: %d }' % pct]))
    if title:
        for lang in names:
            names[lang]['item.' + code] = title[lang]
    inputs = merged([(pick(HERBS, skill), 2), (pick(HERBS, max(1, skill // 2)), 1)])
    potion_recipes.append(recipe(code, 'alchemy', skill, inputs, max(1, value // 4),
                                 max(15, round(0.25 * xp_to_next(skill))), round5(10 + 0.4 * skill)))
for code, lvl, pct, value, title in MANA_POTIONS:
    skill = skill_for(lvl)
    potion_items.append('\n'.join([
        '  - code: ' + code, '    display_name_key: item.' + code, '    item_class: consumable',
        '    item_subclass: potion', '    min_level: %d' % lvl, '    base_value_copper: %d' % value,
        '    base_weight_g: 250', '    stackable: true', '    max_stack: 20', '    bind_kind: boe',
        '    use_effect: { kind: mana_pct, amount: %d }' % pct]))
    for lang in names:
        names[lang]['item.' + code] = title[lang]
    inputs = merged([(pick(HERBS, skill), 1), (pick(HERBS, max(1, skill // 2)), 2)])
    potion_recipes.append(recipe(code, 'alchemy', skill, inputs, max(1, value // 4),
                                 max(15, round(0.25 * xp_to_next(skill))), round5(10 + 0.4 * skill)))
for code, title in PROFESSIONS.items():
    for lang in names:
        names[lang]['profession.' + code] = title[lang]

p = os.path.join(ROOT, 'recipes.yaml')
s, crlf = read(p)
s = put_block(s, 'alchemy', '\n\n'.join(potion_recipes))
s = put_block(s, 'gear', '\n\n'.join(gear_recipes))
write(p, s, crlf)
p = os.path.join(ROOT, 'items.yaml')
s, crlf = read(p)
s = put_block(s, 'potions', '\n\n'.join(potion_items))
write(p, s, crlf)
merge(names)
by_prof = {}
for r in gear_recipes:
    prof = r.split('profession_code: ')[1].split('\n')[0]
    by_prof[prof] = by_prof.get(prof, 0) + 1
print('gear recipes', by_prof, 'potions', len(potion_items))
