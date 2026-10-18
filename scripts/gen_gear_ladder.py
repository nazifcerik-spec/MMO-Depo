"""Extend the gear ladder past level 25 (DD-254, DD-255): a tier for every
class weapon line and the five armor slots at 35, 50, 65, 80, 95 and then
every 25 levels (110 ... 985), continuing the level 1-25 progression.
Writes a generated block into content/gear.yaml and item names in en, tr,
zh-CN and es.

Run from the repository root: python scripts/gen_gear_ladder.py
"""
import io
import os
import sys

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

ROOT = os.path.join(os.path.dirname(__file__), '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)

# Tier level, material code, material names (es as "de <material>").
TIERS = [
    (35, 'mithril', L('Mithril', 'Mithril', '秘银', 'mithril')),
    (50, 'adamant', L('Adamant', 'Adamant', '精金', 'adamantio')),
    (65, 'wyrmforged', L('Wyrmforged', 'Ejder Dövme', '龙铸', 'forja dracónica')),
    (80, 'starforged', L('Starforged', 'Yıldız Dövme', '星铸', 'forja estelar')),
    (95, 'titan', L('Titan', 'Titan', '泰坦', 'titán')),
    (110, 'royal', L('Royal', 'Kraliyet', '王家', 'la realeza')),
    (135, 'sunforged', L('Sunforged', 'Güneş Dövme', '日铸', 'forja solar')),
    (160, 'obsidian', L('Obsidian', 'Obsidyen', '黑曜石', 'obsidiana')),
    (185, 'sandsteel', L('Sandsteel', 'Kum Çeliği', '砂钢', 'acero de arena')),
    (210, 'runebound', L('Runebound', 'Rün Bağlı', '符文', 'runas')),
    (235, 'phoenix', L('Phoenix', 'Anka', '凤凰', 'fénix')),
    (260, 'dreadforged', L('Dreadforged', 'Dehşet Dövme', '恐铸', 'forja del terror')),
    (285, 'sovereign', L('Sovereign', 'Hükümdar', '君王', 'soberano')),
    (310, 'coral', L('Coral', 'Mercan', '珊瑚', 'coral')),
    (335, 'drakebone', L('Drakebone', 'Ejderha Kemiği', '龙骨', 'hueso de dragón')),
    (360, 'mistwoven', L('Mistwoven', 'Sis Örgü', '雾织', 'niebla')),
    (385, 'prismatic', L('Prismatic', 'Prizmatik', '棱彩', 'prisma')),
    (410, 'magmaforged', L('Magmaforged', 'Magma Dövme', '熔铸', 'magma')),
    (435, 'tempest', L('Tempest', 'Bora', '风暴', 'tempestad')),
    (460, 'voidtouched', L('Voidtouched', 'Hiçlik', '虚空', 'vacío')),
    (485, 'seraphic', L('Seraphic', 'Melek', '炽天使', 'serafín')),
    (510, 'eclipse', L('Eclipse', 'Tutulma', '日蚀', 'eclipse')),
    (535, 'wyrmlord', L('Wyrmlord', 'Ejder Lordu', '龙王', 'señor dragón')),
    (560, 'celestine', L('Celestine', 'Göksel', '天界', 'celestial')),
    (585, 'eternal', L('Eternal', 'Ebedi', '永恒', 'eternidad')),
    (610, 'starless', L('Starless', 'Yıldızsız', '无星', 'sin estrellas')),
    (635, 'nightsteel', L('Nightsteel', 'Gece Çeliği', '夜钢', 'acero nocturno')),
    (660, 'chronoforged', L('Chronoforged', 'Zaman Dövme', '时铸', 'forja temporal')),
    (685, 'ascendant', L('Ascendant', 'Yükselmiş', '升华', 'ascensión')),
    (710, 'stormborn', L('Stormborn', 'Fırtına Doğumlu', '风暴之子', 'tormenta')),
    (735, 'riftforged', L('Riftforged', 'Yarık Dövme', '裂隙铸', 'grieta')),
    (760, 'chaos', L('Chaos', 'Kaos', '混沌', 'caos')),
    (785, 'skyforged', L('Skyforged', 'Gök Dövme', '天铸', 'forja celeste')),
    (810, 'archangel', L('Archangel', 'Başmelek', '大天使', 'arcángel')),
    (835, 'aurora', L('Aurora', 'Şafak', '极光', 'aurora')),
    (860, 'ashen', L('Ashen', 'Kül', '灰烬', 'ceniza')),
    (885, 'leviathan', L('Leviathan', 'Leviathan', '利维坦', 'leviatán')),
    (910, 'doomforged', L('Doomforged', 'Kıyamet Dövme', '末日铸', 'forja del fin')),
    (935, 'voidforged', L('Voidforged', 'Hiçlik Dövme', '虚铸', 'forja del vacío')),
    (960, 'godslayer', L('Godslayer', 'Tanrı Katili', '弑神', 'matadioses')),
    (985, 'primordial', L('Primordial', 'Kadim Çağ', '太初', 'primordial')),
]

# Tiers 35-95 were tuned by hand (DD-254); later tiers follow the same lines.
FIXED_WEAPON = {35: (88, 12, 7), 50: (120, 16, 9), 65: (152, 21, 12), 80: (184, 25, 14), 95: (216, 30, 17)}
FIXED_SECONDARY = {
    'mp_max': {35: 150, 50: 200, 65: 250, 80: 300, 95: 350},
    'hit_pct': {35: 5, 50: 5, 65: 6, 80: 6, 95: 7},
    'crit_pct': {35: 4, 50: 5, 65: 5, 80: 6, 95: 6},
}


def weapon_stats(level, sec):
    if level in FIXED_WEAPON:
        atk, prim, s = FIXED_WEAPON[level]
    else:
        atk, prim, s = round(2.2 * level + 7), round(0.315 * level), round(0.18 * level)
    if sec in FIXED_SECONDARY:
        s = FIXED_SECONDARY[sec].get(level) or {
            'mp_max': round(3.7 * level), 'hit_pct': round(5 + level / 45), 'crit_pct': round(4 + level / 60)}[sec]
    return atk, prim, s


# class, item subclass, base code, names, weight, primary target, secondary target
WEAPONS = [
    ('warrior', 'sword', 'sword', L('Sword', 'Kılıç', '剑', 'Espada'), 2500, 'str', 'vit'),
    ('mage', 'staff', 'staff', L('Staff', 'Asa', '法杖', 'Bastón'), 1800, 'int', 'mp_max'),
    ('ranger', 'bow', 'bow', L('Bow', 'Yay', '弓', 'Arco'), 1400, 'agi', 'hit_pct'),
    ('cleric', 'mace', 'mace', L('Mace', 'Topuz', '锤', 'Maza'), 2800, 'cha', 'vit'),
    ('rogue', 'dagger', 'dagger', L('Dagger', 'Hançer', '匕首', 'Daga'), 800, 'agi', 'crit_pct'),
    ('monk', 'fist_weapon', 'fist', L('Knuckles', 'Muşta', '拳刃', 'Nudilleras'), 1200, 'agi', 'vit'),
    ('paladin', 'two_handed_mace', 'hammer', L('Warhammer', 'Savaş Çekici', '战锤', 'Martillo de guerra'), 3500, 'vit', 'str'),
    ('druid', 'sickle', 'sickle', L('Sickle', 'Orak', '镰刀', 'Hoz'), 1500, 'cha', 'vit'),
    ('bard', 'lute', 'lute', L('Lute', 'Lavta', '鲁特琴', 'Laúd'), 1600, 'luck', 'agi'),
    ('shaman', 'totem_focus', 'totem', L('Totem', 'Totem', '图腾', 'Tótem'), 1700, 'cha', 'int'),
]
# slot, names, weight, tier-25 value, (defense, bonus target, bonus) per level-95 tier
ARMOR = [
    ('helm', L('Helm', 'Kask', '头盔', 'Yelmo'), 1800, 765, [8, 10, 13, 15, 18], 'vit', [5, 7, 9, 11, 13]),
    ('chest', L('Chestplate', 'Göğüslük', '胸甲', 'Coraza'), 4000, 1224, [14, 19, 24, 29, 34], 'hp_max', [150, 210, 270, 330, 390]),
    ('legs', L('Legguards', 'Bacaklık', '护腿', 'Grebas'), 3000, 1071, [11, 15, 19, 23, 27], 'vit', [7, 9, 11, 13, 15]),
    ('boots', L('Boots', 'Çizme', '靴', 'Botas'), 1500, 688, [7, 9, 11, 13, 15], 'agi', [5, 7, 9, 11, 13]),
    ('gloves', L('Gloves', 'Eldiven', '手套', 'Guantes'), 700, 612, [7, 9, 11, 13, 15], 'hit_pct', [4, 5, 6, 7, 8]),
]
FIXED_LEVELS = [35, 50, 65, 80, 95]


def armor_stats(level, defs, bonus_target, bonuses):
    if level in FIXED_LEVELS:
        i = FIXED_LEVELS.index(level)
        return defs[i], bonuses[i]
    d = round(defs[-1] / 95 * level)
    b = round(5 + level / 40) if bonus_target == 'hit_pct' else round(bonuses[-1] / 95 * level)
    return d, b


FIXED_VALUE_MULT = {35: 1.55, 50: 2.46, 65: 3.46, 80: 4.5, 95: 5.66}


def value(v25, level):
    return round(v25 * FIXED_VALUE_MULT.get(level, (level / 25) ** 1.3))


def name(mat, base, lang):
    if lang == 'zh':
        return mat[lang] + base[lang]
    if lang == 'es':
        return '%s de %s' % (base[lang], mat[lang])
    return '%s %s' % (mat[lang], base[lang])


def item(code, cls, subclass, level, val, weight, slot, restriction, bonuses):
    lines = ['  - code: ' + code, '    display_name_key: item.' + code, '    item_class: ' + cls,
             '    item_subclass: ' + subclass, '    min_level: %d' % level, '    base_value_copper: %d' % val,
             '    base_weight_g: %d' % weight, '    bind_kind: boe', '    affixable: true', '    equip_slot: ' + slot]
    if restriction:
        lines.append('    class_restriction: ' + restriction)
    lines.append('    stat_bonuses:')
    lines += ['      - { kind: flat, target: %s, amount: %d }' % (t, a) for t, a in bonuses]
    return '\n'.join(lines)


items, names = [], {lang: {} for lang in ('en', 'tr', 'zh', 'es')}
for lvl, mat, mat_names in TIERS:
    items.append('  # --- tier %d: %s ---' % (lvl, mat))
    for cls, subclass, base, base_names, weight, prim, sec in WEAPONS:
        code = '%s_%s' % (base, mat)
        atk, p, s = weapon_stats(lvl, sec)
        items.append(item(code, 'weapon', subclass, lvl, value(1530, lvl), weight, 'mainhand', cls,
                          [('attack', atk), (prim, p), (sec, s)]))
        for lang in names:
            names[lang]['item.' + code] = name(mat_names, base_names, lang)
    for slot, slot_names, weight, value25, defs, bt, bvals in ARMOR:
        code = '%s_%s' % (slot, mat)
        d, b = armor_stats(lvl, defs, bt, bvals)
        items.append(item(code, 'armor', slot, lvl, value(value25, lvl), weight, slot, None, [('defense', d), (bt, b)]))
        for lang in names:
            names[lang]['item.' + code] = name(mat_names, slot_names, lang)

p = os.path.join(ROOT, 'gear.yaml')
with io.open(p, encoding='utf-8', newline='') as f:
    raw = f.read()
crlf = '\r\n' in raw
s = raw.replace('\r\n', '\n')
start, end = '  # >>> generated: gear ladder 35-95 (scripts/gen_gear_ladder.py)\n', '  # <<< generated: gear ladder 35-95\n'
new_start, new_end = '  # >>> generated: gear ladder (scripts/gen_gear_ladder.py)\n', '  # <<< generated: gear ladder\n'
if start in s:  # older marker
    s = s[:s.index(start)] + new_start + 'x\n' + new_end + s[s.index(end) + len(end):]
block = new_start + '\n\n'.join(items) + '\n' + new_end
if new_start in s:
    s = s[:s.index(new_start)] + block + s[s.index(new_end) + len(new_end):]
else:
    s = s.rstrip('\n') + '\n\n' + block
with io.open(p, 'w', encoding='utf-8', newline='') as f:
    f.write(s.replace('\n', '\r\n') if crlf else s)
merge(names)
print('items', sum(1 for i in items if i.startswith('  - code')))
