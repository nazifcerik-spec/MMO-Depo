"""Art audit (DD-334): which game content ships without a picture.

Reads the content YAML and the client asset folders and writes
docs/art-missing.md — per category: how many exist, how many are missing,
the missing codes, and whether the client already shows that art (a "UI
hook") or would need one first.

Run from the repository root: python scripts/art_audit.py
"""
import glob
import io
import os
import re

ROOT = os.path.join(os.path.dirname(__file__), '..')
A = os.path.join(ROOT, 'client', 'assets')


def codes(file, pattern=r'^  - code: (\w+)'):
    with io.open(os.path.join(ROOT, 'content', file), encoding='utf-8') as f:
        return re.findall(pattern, f.read(), re.M)


def have(folder, exts=('.png', '.webp')):
    return {os.path.splitext(os.path.basename(p))[0] for p in glob.glob(os.path.join(A, folder, '*'))
            if p.endswith(exts)}


def section(items_file, cls):
    """Item codes of the given item_class in a content file."""
    with io.open(os.path.join(ROOT, 'content', items_file), encoding='utf-8') as f:
        text = f.read()
    out = []
    for block in re.split(r'\n(?=  - code: )', text):
        m = re.match(r'\s*- code: (\w+)', block)
        c = re.search(r'item_class: (\w+)', block)
        if m and c and (cls is None or c.group(1) in cls):
            out.append(m.group(1))
    return out


items = have('items')
rows = []  # (title, total, missing list, hooked)


def add(title, all_codes, present, hooked=True):
    all_codes = list(dict.fromkeys(all_codes))
    rows.append((title, len(all_codes), [c for c in all_codes if c not in present], hooked))


item_files = ['items.yaml', 'gear.yaml', 'gathering.yaml', 'sets.yaml', 'pets.yaml']
add('Ekipman (silah, zırh, takı, set)', [c for f in item_files for c in section(f, {'weapon', 'armor', 'jewelry'})], items)
add('İksirler + yemekler (tüketilebilir)', [c for f in item_files for c in section(f, {'consumable'})], items)
add('Malzemeler (cevher, kütük, ot, deri, külçe)', [c for f in item_files for c in section(f, {'material', 'reagent'})], items)
add('Efsun parşömenleri', [c for c in section('items.yaml', None) if c.startswith('enchant_')], items)
add('Pet yumurtaları', section('pets.yaml', None), items)
other = [c for f in item_files for c in section(f, None)
         if c not in {x for r in rows for x in []}]
pets = codes('pets.yaml', r'^  - code: (\w+)\n    display_name_key: pet\.')
add('Pet evreleri (tür × 3 evre)', ['%s_%d' % (p, s) for p in pets for s in (1, 2, 3)], have('pets'))
add('Zindanlar + baskınlar', codes('dungeons.yaml'), have('dungeons'))
add('Görevler', codes('quests.yaml'), have('quests'))
add('Bölgeler (harita madalyonu)', codes('regions.yaml'), have('regions'))
add('Sınıf amblemleri', ['warrior', 'paladin', 'ranger', 'rogue', 'mage', 'cleric', 'druid', 'shaman', 'monk', 'bard'],
    have('class_icons'))
add('Sınıf portreleri (oluşturma ekranı, 3:4)',
    ['warrior', 'paladin', 'ranger', 'rogue', 'mage', 'cleric', 'druid', 'shaman', 'monk', 'bard'], have('classes'))
nodes = codes('progression.yaml', r'^\s+- code: (\w+)\s+parent:')
add('Terfi + uzmanlık amblemleri', nodes, have('class_nodes'))
prof = codes('professions.yaml')
add('Üretim meslekleri', prof[:9], have('professions'))
add('Toplama meslekleri (Toplama sekmesi)', prof[9:], have('professions'))
add('Yardımcı NPC portreleri', ['miner', 'woodsman', 'angler', 'herbalist', 'trapper', 'smith', 'tanner', 'weaver',
                                 'carver', 'apothecary', 'cook', 'smelter', 'scribe', 'gemcutter'], have('helpers'))
add('Canavarlar (av listesi)', codes('mobs.yaml'), have('mobs'), hooked=False)
add('Sınıf yetenekleri', codes('skills.yaml'), have('skills'), hooked=False)

out = ['# Eksik görsel listesi', '',
       '`python scripts/art_audit.py` üretir (DD-334). "Arayüzde" sütunu: istemci bu görseli zaten gösteriyor mu '
       '(hayırsa önce bir gösterim yeri eklenmeli). Eksik görselde istemci ikon gösterir; hiçbir ekran bozulmaz.', '',
       '| Alan | Toplam | Var | Eksik | Arayüzde |', '|---|---:|---:|---:|---|']
for title, total, miss, hooked in sorted(rows, key=lambda r: len(r[2])):
    out.append('| %s | %d | %d | %d | %s |' % (title, total, total - len(miss), len(miss), 'evet' if hooked else 'hayır'))
out += ['', '## Eksik kodlar', '']
for title, total, miss, hooked in sorted(rows, key=lambda r: len(r[2])):
    if miss:
        out += ['### %s (%d)' % (title, len(miss)), '', ', '.join('`%s`' % c for c in miss), '']
with io.open(os.path.join(ROOT, 'docs', 'art-missing.md'), 'w', encoding='utf-8', newline='\n') as f:
    f.write('\n'.join(out))
for title, total, miss, hooked in sorted(rows, key=lambda r: len(r[2])):
    print('%-45s %5d %5d %s' % (title, total, len(miss), '' if hooked else '(UI yok)'))
