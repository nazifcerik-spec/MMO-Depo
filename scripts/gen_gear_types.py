"""Armor types and second weapon families (DD-271, DD-272).

Every gear tier gets cloth, leather and plate armor, and every class a second
weapon family. The hand-written and ladder gear is the source: plate armor
(tier 1: the leather set) and each class's first weapon line. New pieces copy
their counterpart's level, stats and value, so the balance stays as it is —
the type only decides who may wear it. Writes a generated block into
content/gear.yaml, the tier-1 pieces into the general vendor, and names in en,
tr, zh-CN and es.

Run from the repository root after gen_gear_ladder.py: python scripts/gen_gear_types.py
"""
import io
import os
import sys

import yaml

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.join(HERE, '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)


def ladder_tiers():
    """Tier materials 35..985 from gen_gear_ladder.py, without running it."""
    with io.open(os.path.join(HERE, 'gen_gear_ladder.py'), encoding='utf-8') as f:
        src = f.read()
    ns = {'__file__': os.path.join(HERE, 'gen_gear_ladder.py')}
    exec(compile(src[:src.index('# Tiers 35-95 were tuned')], 'gen_gear_ladder', 'exec'), ns)
    return {lvl: (mat, names) for lvl, mat, names in ns['TIERS']}


# Tiers 1-25 were named per line by hand; new lines use these materials.
LOW = {
    1: ('novice', L('Novice', 'Acemi', '新手', 'novato')),
    5: ('journeyman', L('Journeyman', 'Kalfa', '熟手', 'oficial')),
    10: ('fine', L('Fine', 'Kaliteli', '精制', 'calidad')),
    15: ('veteran', L('Veteran', 'Kıdemli', '老兵', 'veterano')),
    20: ('champion', L('Champion', 'Şampiyon', '勇士', 'campeón')),
    25: ('runed', L('Runed', 'Rünlü', '符文', 'runas')),
}

# Armor type -> slot -> (code base, names); weight relative to plate.
ARMOR = {
    'cloth': (0.4, {
        'helm': ('hood', L('Hood', 'Kukuleta', '兜帽', 'Capucha')),
        'chest': ('robe', L('Robe', 'Cübbe', '法袍', 'Túnica')),
        'legs': ('trousers', L('Trousers', 'Şalvar', '布裤', 'Calzas')),
        'boots': ('sandals', L('Sandals', 'Sandalet', '便鞋', 'Sandalias')),
        'gloves': ('wraps', L('Wraps', 'Sargı', '裹手', 'Vendas')),
    }),
    'leather': (0.7, {
        'helm': ('cowl', L('Cowl', 'Deri Başlık', '皮帽', 'Caperuza')),
        'chest': ('jerkin', L('Jerkin', 'Deri Yelek', '皮甲', 'Jubón')),
        'legs': ('leggings', L('Leggings', 'Deri Tozluk', '皮裤', 'Perneras')),
        'boots': ('treads', L('Treads', 'Deri Çizme', '皮靴', 'Botines')),
        'gloves': ('grips', L('Grips', 'Deri Eldiven', '皮手套', 'Mitones')),
    }),
    'plate': (1.0, {  # tier 1 only: the rest of the ladder is plate already
        'helm': ('helm', L('Helm', 'Kask', '头盔', 'Yelmo')),
        'chest': ('chest', L('Chestplate', 'Göğüslük', '胸甲', 'Coraza')),
        'legs': ('legs', L('Legguards', 'Bacaklık', '护腿', 'Grebas')),
        'boots': ('boots', L('Boots', 'Çizme', '靴', 'Botas')),
        'gloves': ('gloves', L('Gloves', 'Eldiven', '手套', 'Guantes')),
    }),
}

# class -> (first family, second family, code base, names, weight)
WEAPONS = {
    'warrior': ('sword', 'axe', 'axe', L('Axe', 'Balta', '斧', 'Hacha'), 2700),
    'paladin': ('two_handed_mace', 'sword', 'longsword', L('Longsword', 'Uzun Kılıç', '长剑', 'Espada larga'), 2600),
    'shaman': ('totem_focus', 'ritual_axe', 'ritualaxe', L('Ritual Axe', 'Ayin Baltası', '仪式斧', 'Hacha ritual'), 2400),
    'rogue': ('dagger', 'short_sword', 'shortsword', L('Short Sword', 'Kısa Kılıç', '短剑', 'Espada corta'), 1500),
    'ranger': ('bow', 'crossbow', 'crossbow', L('Crossbow', 'Arbalet', '弩', 'Ballesta'), 2600),
    'monk': ('fist_weapon', 'tonfa', 'tonfa', L('Tonfa', 'Tonfa', '拐', 'Tonfa'), 1300),
    'bard': ('lute', 'rapier', 'rapier', L('Rapier', 'Meç', '细剑', 'Estoque'), 1200),
    'mage': ('staff', 'wand', 'wand', L('Wand', 'Değnek', '魔杖', 'Varita'), 600),
    'cleric': ('mace', 'holy_tome', 'tome', L('Holy Tome', 'Kutsal Kitap', '圣典', 'Tomo sagrado'), 1000),
    'druid': ('sickle', 'nature_focus', 'focus', L('Nature Focus', 'Doğa Odağı', '自然法器', 'Foco natural'), 900),
}
VENDOR = 'general_ironhold'
START, END = '  # >>> generated: gear types (scripts/gen_gear_types.py)\n', '  # <<< generated: gear types\n'


def name(mat, base, lang):
    if lang == 'zh':
        return mat[lang] + base[lang]
    if lang == 'es':
        return '%s de %s' % (base[lang], mat[lang])
    return '%s %s' % (mat[lang], base[lang])


def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        raw = f.read()
    return raw.replace('\r\n', '\n'), '\r\n' in raw


def write(path, s, crlf):
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s.replace('\n', '\r\n') if crlf else s)


def emit(code, src, subclass, weight, restriction=None):
    lines = ['  - code: ' + code, '    display_name_key: item.' + code, '    item_class: ' + src['item_class'],
             '    item_subclass: ' + subclass, '    min_level: %d' % src['min_level'],
             '    base_value_copper: %d' % src['base_value_copper'], '    base_weight_g: %d' % weight,
             '    bind_kind: boe', '    affixable: true', '    equip_slot: ' + src['equip_slot']]
    if restriction:
        lines.append('    class_restriction: ' + restriction)
    lines.append('    stat_bonuses:')
    lines += ['      - { kind: %s, target: %s, amount: %s }' % (b['kind'], b['target'], b['amount'])
              for b in src['stat_bonuses']]
    return '\n'.join(lines)


gear_path = os.path.join(ROOT, 'gear.yaml')
gear_src, gear_crlf = read(gear_path)
if START in gear_src:
    gear_src = gear_src[:gear_src.index(START)].rstrip('\n') + '\n' + gear_src[gear_src.index(END) + len(END):]
source = [g for g in yaml.safe_load(gear_src)['items'] if not g.get('requirements')]
known = {g['code'] for g in source}
for path in ('items.yaml', 'gathering.yaml'):
    known |= {i['code'] for i in yaml.safe_load(read(os.path.join(ROOT, path))[0])['items']}

tiers = dict(LOW)
tiers.update(ladder_tiers())
levels = sorted({g['min_level'] for g in source})
assert levels == sorted(tiers), 'gear tiers and tier materials differ: %s' % sorted(set(levels) ^ set(tiers))

out, names, vendor = [], {lang: {} for lang in ('en', 'tr', 'zh', 'es')}, []


def add(code, src, subclass, weight, title, restriction=None):
    assert code not in known, 'duplicate code ' + code
    known.add(code)
    out.append(emit(code, src, subclass, weight, restriction))
    for lang in names:
        names[lang]['item.' + code] = title[lang]
    if src['min_level'] == 1:
        vendor.append(code)


for lvl in levels:
    mat, mat_names = tiers[lvl]
    out.append('  # --- tier %d ---' % lvl)
    at = [g for g in source if g['min_level'] == lvl]
    # The tier's body armor: plate, or the leather set at tier 1.
    body = {g['equip_slot']: g for g in at if g['item_class'] == 'armor' and g['item_subclass'].endswith('_armor')}
    have = body[next(iter(body))]['item_subclass'][:-len('_armor')]
    for kind, (wmult, slots) in ARMOR.items():
        if kind == have:
            continue
        for slot, (base, base_names) in slots.items():
            src = body[slot]
            add('%s_%s' % (base, mat), src, kind + '_armor', round(src['base_weight_g'] * wmult),
                {lang: name(mat_names, base_names, lang) for lang in names})
    for cls, (first, second, base, base_names, weight) in WEAPONS.items():
        src = next(g for g in at if g.get('class_restriction') == cls and g['item_subclass'] == first)
        add('%s_%s' % (base, mat), src, second, weight, {lang: name(mat_names, base_names, lang) for lang in names}, cls)

gear_src = gear_src.rstrip('\n') + '\n\n' + START + '\n\n'.join(out) + '\n' + END
write(gear_path, gear_src, gear_crlf)

vendor_path = os.path.join(ROOT, 'vendors.yaml')
vs, vcrlf = read(vendor_path)
vstart, vend = '      # >>> generated: tier-1 gear types (scripts/gen_gear_types.py)\n', '      # <<< generated: tier-1 gear types\n'
block = vstart + ''.join('      - { template_code: %s }\n' % c for c in vendor) + vend
if vstart in vs:
    vs = vs[:vs.index(vstart)] + block + vs[vs.index(vend) + len(vend):]
else:
    anchor = '      - { template_code: gloves_leather }\n'
    assert anchor in vs, 'vendor anchor missing'
    vs = vs.replace(anchor, anchor + block)
write(vendor_path, vs, vcrlf)
merge(names)
print('items', sum(1 for o in out if o.startswith('  - code')), 'vendor', len(vendor))
