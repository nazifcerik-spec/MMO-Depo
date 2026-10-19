"""Jewelry and guardian shields (DD-276).

Jewelry fills the necklace, ring and trinket slots for every class from the
level-25 tier up (below it the balance targets keep stronger mobs
dangerous): a necklace adds HP, a ring attack and hit, a charm defense and
evasion. Rings and charms are worn in pairs. Shields are the offhand of the
Guardian line (warrior specialization, level 100): a ladder from the
level-110 tier, crafted, never dropped (branch-gated gear is no random drop).
Writes a generated block into content/gear.yaml and names in en, tr, zh-CN
and es.

Run from the repository root after gen_gear_types.py: python scripts/gen_jewelry.py
"""
import io
import os
import sys

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


TIERS = {25: ('runed', L('Runed', 'Rünlü', '符文', 'runas'))}  # as gen_gear_types.LOW
TIERS.update(ladder_tiers())

# base code, slot, names, weight, tier-25 value, stats(level) -> [(target, amount)]
JEWELRY = [
    ('necklace', 'necklace', L('Necklace', 'Kolye', '项链', 'Collar'), 100, 700,
     lambda lv: [('hp_max', round(1.2 * lv + 10))]),
    ('ring', 'ring', L('Ring', 'Yüzük', '戒指', 'Anillo'), 30, 600,
     lambda lv: [('attack', round(0.12 * lv + 2)), ('hit_pct', round(1 + lv / 250))]),
    ('charm', 'trinket', L('Charm', 'Tılsım', '护符', 'Amuleto'), 80, 600,
     lambda lv: [('defense', round(0.08 * lv + 1)), ('evasion_pct', round(1 + lv / 500))]),
]
SHIELD = ('shield', L('Shield', 'Kalkan', '盾', 'Escudo'), 4500, 1100,
          lambda lv: [('defense', round(0.2 * lv + 2)), ('block_pct', round(2 + lv / 100))])
JEWELRY_FROM, SHIELD_FROM = 25, 110
START, END = '  # >>> generated: jewelry (scripts/gen_jewelry.py)\n', '  # <<< generated: jewelry\n'


def name(mat, base, lang):  # as gen_gear_ladder.name
    if lang == 'zh':
        return mat[lang] + base[lang]
    if lang == 'es':
        return '%s de %s' % (base[lang], mat[lang])
    return '%s %s' % (mat[lang], base[lang])


def value(v25, level):
    return round(v25 * (level / 25) ** 1.3)


def item(code, cls, subclass, level, val, weight, slot, bonuses, extra=()):
    lines = ['  - code: ' + code, '    display_name_key: item.' + code, '    item_class: ' + cls,
             '    item_subclass: ' + subclass, '    min_level: %d' % level, '    base_value_copper: %d' % val,
             '    base_weight_g: %d' % weight, '    bind_kind: boe', '    affixable: %s' % ('true' if cls == 'armor' else 'false'),
             '    equip_slot: ' + slot] + list(extra)
    lines.append('    stat_bonuses:')
    lines += ['      - { kind: flat, target: %s, amount: %d }' % (t, a) for t, a in bonuses]
    return '\n'.join(lines)


out, names = [], {lang: {} for lang in ('en', 'tr', 'zh', 'es')}
for lvl in sorted(TIERS):
    mat, mat_names = TIERS[lvl]
    if lvl < JEWELRY_FROM:
        continue
    out.append('  # --- tier %d ---' % lvl)
    for base, slot, base_names, weight, v25, stats in JEWELRY:
        code = '%s_%s' % (base, mat)
        out.append(item(code, 'jewelry', base if base != 'charm' else 'trinket', lvl, value(v25, lvl), weight, slot,
                        stats(lvl)))
        for lang in names:
            names[lang]['item.' + code] = name(mat_names, base_names, lang)
    if lvl >= SHIELD_FROM:
        base, base_names, weight, v25, stats = SHIELD
        code = '%s_%s' % (base, mat)
        out.append(item(code, 'armor', 'shield', lvl, value(v25, lvl), weight, 'offhand', stats(lvl),
                        ['    class_restriction: warrior', '    requirements: { spec: guardian }']))
        for lang in names:
            names[lang]['item.' + code] = name(mat_names, base_names, lang)

p = os.path.join(ROOT, 'gear.yaml')
with io.open(p, encoding='utf-8', newline='') as f:
    raw = f.read()
crlf = '\r\n' in raw
s = raw.replace('\r\n', '\n')
block = START + '\n\n'.join(out) + '\n' + END
if START in s:
    s = s[:s.index(START)] + block + s[s.index(END) + len(END):]
else:
    s = s.rstrip('\n') + '\n\n' + block
with io.open(p, 'w', encoding='utf-8', newline='') as f:
    f.write(s.replace('\n', '\r\n') if crlf else s)
merge(names)
print('items', sum(1 for o in out if o.startswith('  - code')))
