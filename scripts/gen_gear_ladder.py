"""Extend the gear ladder past level 25 (DD-254): tiers at 35, 50, 65, 80
and 95 for every class weapon line and the five armor slots, continuing the
level 1-25 progression. Writes a generated block into content/gear.yaml and
item names into content/i18n/{en,tr}.yaml.

Run from the repository root: python scripts/gen_gear_ladder.py
"""
import io
import os
import re

ROOT = os.path.join(os.path.dirname(__file__), '..', 'content')

# level, material code, en, tr, value multiplier vs. tier 25
TIERS = [
    (35, 'mithril', 'Mithril', 'Mithril', 1.55),
    (50, 'adamant', 'Adamant', 'Adamant', 2.46),
    (65, 'wyrmforged', 'Wyrmforged', 'Ejder Dövme', 3.46),
    (80, 'starforged', 'Starforged', 'Yıldız Dövme', 4.5),
    (95, 'titan', 'Titan', 'Titan', 5.66),
]
WEAPON_ATK = [88, 120, 152, 184, 216]
PRIMARY = [12, 16, 21, 25, 30]
SECONDARY = [7, 9, 12, 14, 17]

# class, item subclass, base code, en, tr, weight, primary target, secondary (target, per-tier values)
WEAPONS = [
    ('warrior', 'sword', 'sword', 'Sword', 'Kılıç', 2500, 'str', ('vit', SECONDARY)),
    ('mage', 'staff', 'staff', 'Staff', 'Asa', 1800, 'int', ('mp_max', [150, 200, 250, 300, 350])),
    ('ranger', 'bow', 'bow', 'Bow', 'Yay', 1400, 'agi', ('hit_pct', [5, 5, 6, 6, 7])),
    ('cleric', 'mace', 'mace', 'Mace', 'Topuz', 2800, 'cha', ('vit', SECONDARY)),
    ('rogue', 'dagger', 'dagger', 'Dagger', 'Hançer', 800, 'agi', ('crit_pct', [4, 5, 5, 6, 6])),
    ('monk', 'fist_weapon', 'fist', 'Knuckles', 'Muşta', 1200, 'agi', ('vit', SECONDARY)),
    ('paladin', 'two_handed_mace', 'hammer', 'Warhammer', 'Savaş Çekici', 3500, 'vit', ('str', SECONDARY)),
    ('druid', 'sickle', 'sickle', 'Sickle', 'Orak', 1500, 'cha', ('vit', SECONDARY)),
    ('bard', 'lute', 'lute', 'Lute', 'Lavta', 1600, 'luck', ('agi', SECONDARY)),
    ('shaman', 'totem_focus', 'totem', 'Totem', 'Totem', 1700, 'cha', ('int', SECONDARY)),
]
# slot, en, tr, weight, tier-25 value, defense per tier, (bonus target, per-tier values)
ARMOR = [
    ('helm', 'Helm', 'Kask', 1800, 765, [8, 10, 13, 15, 18], ('vit', [5, 7, 9, 11, 13])),
    ('chest', 'Chestplate', 'Göğüslük', 4000, 1224, [14, 19, 24, 29, 34], ('hp_max', [150, 210, 270, 330, 390])),
    ('legs', 'Legguards', 'Bacaklık', 3000, 1071, [11, 15, 19, 23, 27], ('vit', [7, 9, 11, 13, 15])),
    ('boots', 'Boots', 'Çizme', 1500, 688, [7, 9, 11, 13, 15], ('agi', [5, 7, 9, 11, 13])),
    ('gloves', 'Gloves', 'Eldiven', 700, 612, [7, 9, 11, 13, 15], ('hit_pct', [4, 5, 6, 7, 8])),
]
WEAPON_VALUE_25 = 1530


def item(code, cls, subclass, level, value, weight, slot, restriction, bonuses):
    lines = ['  - code: ' + code, '    display_name_key: item.' + code, '    item_class: ' + cls,
             '    item_subclass: ' + subclass, '    min_level: %d' % level, '    base_value_copper: %d' % value,
             '    base_weight_g: %d' % weight, '    bind_kind: boe', '    affixable: true', '    equip_slot: ' + slot]
    if restriction:
        lines.append('    class_restriction: ' + restriction)
    lines.append('    stat_bonuses:')
    lines += ['      - { kind: flat, target: %s, amount: %d }' % (t, a) for t, a in bonuses]
    return '\n'.join(lines)


items, names = [], {'en': {}, 'tr': {}}
for i, (lvl, mat, en_mat, tr_mat, mult) in enumerate(TIERS):
    items.append('  # --- tier %d: %s ---' % (lvl, mat))
    for cls, subclass, base, en, tr, weight, prim, (sec, secs) in WEAPONS:
        code = '%s_%s' % (base, mat)
        items.append(item(code, 'weapon', subclass, lvl, round(WEAPON_VALUE_25 * mult), weight, 'mainhand', cls,
                          [('attack', WEAPON_ATK[i]), (prim, PRIMARY[i]), (sec, secs[i])]))
        names['en'][code], names['tr'][code] = en_mat + ' ' + en, tr_mat + ' ' + tr
    for slot, en, tr, weight, value25, defs, (bt, bvals) in ARMOR:
        code = '%s_%s' % (slot, mat)
        items.append(item(code, 'armor', slot, lvl, round(value25 * mult), weight, slot, None,
                          [('defense', defs[i]), (bt, bvals[i])]))
        names['en'][code], names['tr'][code] = en_mat + ' ' + en, tr_mat + ' ' + tr


def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        s = f.read()
    return s.replace('\r\n', '\n'), '\r\n' in s


def write(path, s, crlf):
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s.replace('\n', '\r\n') if crlf else s)


p = os.path.join(ROOT, 'gear.yaml')
s, crlf = read(p)
start, end = '  # >>> generated: gear ladder 35-95 (scripts/gen_gear_ladder.py)\n', '  # <<< generated: gear ladder 35-95\n'
block = start + '\n\n'.join(items) + '\n' + end
if start in s:
    s = s[:s.index(start)] + block + s[s.index(end) + len(end):]
else:
    s = s.rstrip('\n') + '\n\n' + block
write(p, s, crlf)

for lang in ('en', 'tr'):
    p = os.path.join(ROOT, 'i18n', lang + '.yaml')
    s, crlf = read(p)
    lines = s.split('\n')
    keys = {re.match(r'^([\w.\-]+):', l).group(1): i for i, l in enumerate(lines) if re.match(r'^[\w.\-]+:\s*"', l)}
    for code, name in sorted(names[lang].items()):
        key = 'item.' + code
        line = '%s: "%s"' % (key, name)
        if key in keys:
            lines[keys[key]] = line
        else:
            # keep the "# --- item ---" section sorted: insert before the first larger item key
            item_keys = sorted(k for k in keys if k.startswith('item.'))
            after = [k for k in item_keys if k > key]
            at = keys[after[0]] if after else keys[item_keys[-1]] + 1
            lines.insert(at, line)
            keys = {re.match(r'^([\w.\-]+):', l).group(1): i for i, l in enumerate(lines) if re.match(r'^[\w.\-]+:\s*"', l)}
    write(p, '\n'.join(lines), crlf)
print('items', len(items) - len(TIERS))
