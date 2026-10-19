"""Raid sets (DD-277).

Every raid has a set per armor type (cloth, leather, plate): five armor
pieces and a necklace at the raid's level, 15% above the gear tier below it.
Worn together they add: 2 pieces +4% HP, 4 pieces +5% attack, 6 pieces a
type bonus (cloth crit damage, leather crit, plate damage reduction). The
raid boss drops one piece of the set the class wears. Writes
content/sets.yaml (pieces + sets) and names in en, tr, zh-CN and es.

Run from the repository root after gen_jewelry.py: python scripts/gen_sets.py
"""
import io
import os
import sys

import yaml

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)

# raid code -> (piece code prefix, epithet)
RAIDS = {
    'north_keep_siege': ('undying', L('Undying', 'Ölümsüz', '不死', 'inmortal')),
    'crownless_palace': ('tyrant', L("Tyrant's", 'Tiran', '暴君', 'del tirano')),
    'dragon_king_lair': ('dragonking', L("Dragon King's", 'Ejder Kral', '龙王', 'del rey dragón')),
    'time_lord_tower': ('timeless', L('Timeless', 'Zamansız', '永恒', 'atemporal')),
    'archangels_fall': ('fallen', L('Fallen', 'Düşmüş', '堕落', 'caído')),
    'fallen_god_throne': ('unbound', L('Unbound', 'Zincirsiz', '无缚', 'desatado')),
}
TYPES = {
    'cloth': (L('Cloth', 'Kumaş', '布甲', 'de tela'), {
        'helm': ('hood', L('Hood', 'Kukuleta', '兜帽', 'Capucha')),
        'chest': ('robe', L('Robe', 'Cübbe', '法袍', 'Túnica')),
        'legs': ('trousers', L('Trousers', 'Şalvar', '布裤', 'Calzas')),
        'boots': ('sandals', L('Sandals', 'Sandalet', '便鞋', 'Sandalias')),
        'gloves': ('wraps', L('Wraps', 'Sargı', '裹手', 'Vendas')),
        'necklace': ('pendant', L('Pendant', 'Kolye', '吊坠', 'Colgante')),
    }, {'kind': 'flat', 'target': 'crit_dmg_pct', 'amount': 10}),
    'leather': (L('Leather', 'Deri', '皮甲', 'de cuero'), {
        'helm': ('cowl', L('Cowl', 'Deri Başlık', '皮帽', 'Caperuza')),
        'chest': ('jerkin', L('Jerkin', 'Deri Yelek', '皮甲', 'Jubón')),
        'legs': ('leggings', L('Leggings', 'Deri Tozluk', '皮裤', 'Perneras')),
        'boots': ('treads', L('Treads', 'Deri Çizme', '皮靴', 'Botines')),
        'gloves': ('grips', L('Grips', 'Deri Eldiven', '皮手套', 'Mitones')),
        'necklace': ('torc', L('Torc', 'Gerdanlık', '项圈', 'Torques')),
    }, {'kind': 'flat', 'target': 'crit_pct', 'amount': 3}),
    'plate': (L('Plate', 'Plaka', '板甲', 'de placas'), {
        'helm': ('helm', L('Helm', 'Kask', '头盔', 'Yelmo')),
        'chest': ('chestplate', L('Chestplate', 'Göğüslük', '胸甲', 'Coraza')),
        'legs': ('legguards', L('Legguards', 'Bacaklık', '护腿', 'Grebas')),
        'boots': ('boots', L('Boots', 'Çizme', '靴', 'Botas')),
        'gloves': ('gauntlets', L('Gauntlets', 'Eldiven', '手套', 'Guanteletes')),
        'necklace': ('chain', L('Chain', 'Zincir Kolye', '项链', 'Cadena')),
    }, {'kind': 'flat', 'target': 'dmg_reduction_pct', 'amount': 3}),
}
SLOTS = ['helm', 'chest', 'legs', 'boots', 'gloves', 'necklace']
STAT_BOOST = 1.15  # above the tier below the raid


def load(name, key):
    with io.open(os.path.join(ROOT, name), encoding='utf-8') as f:
        return yaml.safe_load(f)[key]


gear = [g for g in load('gear.yaml', 'items') if not g.get('requirements')]
raids = {d['code']: d['min_level'] for d in load('dungeons.yaml', 'dungeons') if d.get('group_max', 1) > 1}
assert set(raids) == set(RAIDS), 'raids and set epithets differ: %s' % (set(raids) ^ set(RAIDS))


def source(kind, slot, level):
    """The ladder piece the set piece is built on: same type and slot, the
    highest tier at or below the raid's level."""
    want = 'necklace' if slot == 'necklace' else kind + '_armor'
    cands = [g for g in gear if g['equip_slot'] == slot and g['item_subclass'] == want and g['min_level'] <= level]
    return max(cands, key=lambda g: g['min_level'])


items, sets, names = [], [], {lang: {} for lang in ('en', 'tr', 'zh', 'es')}
for raid, (prefix, epithet) in RAIDS.items():
    level = raids[raid]
    for kind, (type_name, nouns, bonus6) in TYPES.items():
        pieces = []
        for slot in SLOTS:
            base, noun = nouns[slot]
            code = '%s_%s' % (prefix, base)
            src = source(kind, slot, level)
            jewel = slot == 'necklace'
            if jewel:
                code = '%s_%s' % (prefix, base)  # one necklace per type: pendant / torc / chain
            pieces.append(code)
            bonuses = [(b['target'], max(1, round(b['amount'] * STAT_BOOST))) for b in src['stat_bonuses']]
            items.append('\n'.join([
                '  - code: ' + code, '    display_name_key: item.' + code,
                '    item_class: ' + src['item_class'], '    item_subclass: ' + src['item_subclass'],
                '    min_level: %d' % level, '    base_value_copper: %d' % (2 * src['base_value_copper']),
                '    base_weight_g: %d' % src['base_weight_g'], '    bind_kind: boe',
                '    affixable: %s' % ('false' if jewel else 'true'), '    equip_slot: ' + slot,
                '    stat_bonuses:'] + ['      - { kind: flat, target: %s, amount: %d }' % b for b in bonuses]))
            for lang in names:
                names[lang]['item.' + code] = (
                    noun[lang] + ' ' + epithet[lang] if lang == 'es'
                    else epithet[lang] + noun[lang] if lang == 'zh' else epithet[lang] + ' ' + noun[lang])
        set_code = '%s_%s' % (raid, kind)
        for lang in names:
            names[lang]['itemset.' + set_code] = {
                'en': '%s %s' % (epithet['en'], type_name['en']),
                'tr': '%s %s Seti' % (epithet['tr'], type_name['tr']),
                'zh': '%s%s套装' % (epithet['zh'], type_name['zh']),
                'es': 'Conjunto %s %s' % (type_name['es'], epithet['es']),
            }[lang]
        sets.append('\n'.join([
            '  - code: ' + set_code, '    display_name_key: itemset.' + set_code, '    source: ' + raid,
            '    pieces: [%s]' % ', '.join(pieces), '    bonuses:',
            '      - { pieces: 2, modifiers: [{ kind: mult_add, target: hp_max, amount: 0.04 }] }',
            '      - { pieces: 4, modifiers: [{ kind: mult_add, target: attack, amount: 0.05 }] }',
            '      - { pieces: 6, modifiers: [{ kind: %s, target: %s, amount: %s }] }'
            % (bonus6['kind'], bonus6['target'], bonus6['amount'])]))

out = '\n'.join([
    '# Raid sets (DD-277). Generated by scripts/gen_sets.py — edit the script.',
    '# items: the pieces (same fields as items.yaml); sets: which pieces form a',
    '# set, the raid whose boss drops them and the bonuses at piece counts.',
    '', 'items:', '\n\n'.join(items), '', 'sets:', '\n\n'.join(sets), ''])
with io.open(os.path.join(ROOT, 'sets.yaml'), 'w', encoding='utf-8', newline='\n') as f:
    f.write(out)
merge(names)
print('pieces', len(items), 'sets', len(sets))
