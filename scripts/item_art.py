"""Item art for every equippable item (DD-329): weapons, armor, shields,
jewelry and raid sets, each its own painted icon.

The items are ordered (family, then level) and packed GRID×GRID to a sheet;
each sheet is one image generated from the prompt this script prints, saved
as client/art_originals/items/sheetNNN.jpg (git-ignored), then cut into
client/assets/items/<code>.webp.

    python scripts/item_art.py plan          # sheets and their items
    python scripts/item_art.py prompt N      # the prompt for sheet N
    python scripts/item_art.py cut [N ...]   # cut the saved sheets (all by default)
"""
import glob
import io
import os
import re
import sys

import yaml

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from item_motifs import MOTIF  # noqa: E402

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CONTENT = os.path.join(ROOT, 'content')
SRC = os.path.join(ROOT, 'client', 'art_originals', 'items')
DST = os.path.join(ROOT, 'client', 'assets', 'items')
GRID = 5
SIZE = 128  # the bag shows icons at most ~64 dp

# What each kind is, for the prompt (the item name carries the material).
KIND = {
    'axe': 'one-handed battle axe', 'bow': 'longbow', 'crossbow': 'crossbow', 'dagger': 'dagger',
    'fist_weapon': 'clawed fist weapon', 'holy_tome': 'holy spell tome', 'lute': 'magical lute',
    'mace': 'one-handed mace', 'nature_focus': 'druidic nature focus', 'rapier': 'rapier',
    'ritual_axe': 'shamanic ritual axe', 'short_sword': 'short sword', 'sickle': 'sickle',
    'staff': 'mage staff', 'sword': 'sword', 'tonfa': 'pair of tonfas (straight batons with a short side handle)', 'totem_focus': 'shaman totem',
    'two_handed_mace': 'two-handed warhammer', 'wand': 'magic wand', 'shield': 'shield',
    'necklace': 'necklace', 'ring': 'ring', 'trinket': 'charm trinket',
}
SLOT = {'helm': 'helmet', 'chest': 'chest armor', 'legs': 'leg armor', 'gloves': 'gloves', 'boots': 'boots'}
ARMOR = {'cloth_armor': 'cloth', 'leather_armor': 'leather', 'plate_armor': 'plate'}


def items():
    out = []
    for f in ('items.yaml', 'gear.yaml', 'sets.yaml'):
        with io.open(os.path.join(CONTENT, f), encoding='utf-8') as fh:
            out += [i for i in yaml.safe_load(fh).get('items', []) if i.get('equip_slot')]
    return out


def names():
    out = {}
    with io.open(os.path.join(CONTENT, 'i18n', 'en.yaml'), encoding='utf-8') as fh:
        for line in fh:
            m = re.match(r'^item\.([\w-]+):\s*"(.*)"\s*$', line)
            if m:
                out[m.group(1)] = m.group(2)
    return out


def kind(i):
    sub = i['item_subclass']
    if sub in ARMOR:
        return '%s %s' % (ARMOR[sub], SLOT[i['equip_slot']])
    return KIND.get(sub, sub)


def order(i):
    """Weapons, shields, armor, jewelry; a family's items by level."""
    sub = i['item_subclass']
    group = 0 if i['item_class'] == 'weapon' else 1 if sub == 'shield' else 2 if sub in ARMOR else 3
    family = i['code'].split('_')[0]
    raid = 1 if 'raid' in (i.get('_set') or '') else 0
    return (group, sub, i['equip_slot'], raid, family if i['min_level'] else '', i['min_level'], i['code'])


def plan():
    its = sorted(items(), key=order)
    n = GRID * GRID
    return [its[k:k + n] for k in range(0, len(its), n)]


POS = lambda k: 'row %d col %d' % (k // GRID + 1, k % GRID + 1)


# Silhouette rules, one per item in family order, so a family's items never
# share an outline even where their materials look alike.
FORMS = {
    'weapon': ['long and slender', 'short and broad', 'curved', 'wavy flame-like', 'asymmetric', 'hooked',
               'forked tip', 'leaf-shaped', 'angular faceted', 'spiked', 'crescent-shaped', 'twisted spiral',
               'segmented', 'ornate open-work', 'heavy and blocky', 'needle-thin', 'flared', 'stepped tiers',
               'double-headed', 'barbed'],
    'armor': ['rounded', 'angular', 'layered plates', 'spiked', 'flared', 'slim fitted', 'bulky', 'asymmetric',
              'winged', 'horned', 'fur-trimmed', 'segmented', 'open-work filigree', 'banded straps', 'scalloped',
              'crested', 'quilted', 'studded', 'draped', 'ridged'],
    'jewelry': ['teardrop', 'round', 'square-cut', 'twisted band', 'claw setting', 'cluster of gems', 'hoop',
                'crescent', 'star-shaped', 'knotwork', 'chain-linked', 'coiled serpent', 'triangular', 'oval cameo',
                'spiral', 'winged', 'skull', 'flower', 'eye', 'key-shaped'],
}


def form(i, rank):
    group = 'armor' if i['item_subclass'] in ARMOR or i['item_subclass'] == 'shield' else         'jewelry' if i['item_class'] == 'jewelry' else 'weapon'
    forms = FORMS[group]
    # Offset by family too: rank N of two families must not share a shape.
    family = sum(map(ord, i['code'].split('_')[0]))
    return forms[(rank * 7 + family) % len(forms)]


def ranks():
    """Each item's position within its family (subclass and slot)."""
    out, seen = {}, {}
    for i in sorted(items(), key=order):
        key = (i['item_subclass'], i['equip_slot'])
        out[i['code']] = seen.get(key, 0)
        seen[key] = out[i['code']] + 1
    return out


def motif(i, word):
    """The material's colours, and its figure on a third of the families
    (which third depends on the material) so an emblem does not repeat on
    every weapon."""
    look, figure = MOTIF.get(word, (word.lower(), ''))
    family = i['code'].split('_')[0]
    if figure and (sum(map(ord, family)) + sum(map(ord, word))) % 3 == 0:
        return '%s, featuring %s' % (look, figure)
    return look


def prompt(sheet):
    nm = names()
    # Descriptions, not names: names get painted on as captions. The
    # material's motif gives each its own shape and colours.
    rk = ranks()

    def describe(i):
        return '%s, %s shape, %s' % (kind(i), form(i, rk[i['code']]), motif(i, nm.get(i['code'], i['code']).split()[0]))
    parts = ['%s: %s' % (POS(k), describe(i)) for k, i in enumerate(sheet)]
    return (
        'Generate a square image: a %dx%d grid of %d separate square fantasy RPG item icons of equal size, '
        'divided by thin solid black lines. Each cell shows exactly one item, whole and centered, at a slight '
        'diagonal, on the same plain dark slate-blue backdrop with a soft vignette. Stylized hand-painted game '
        'icon art like polished mobile MMORPG inventory icons: clean bold silhouettes, soft cel shading, crisp '
        'rim light, ornate trims, glowing gems or runes on the stronger ones. Every item must look unique: a '
        'different silhouette, shape, ornament and colour palette from every other cell, never a recolour; even the plain early ones differ in head or blade shape, size, handle and wraps. '
        'Each item must stay clearly recognisable as the kind named first in its cell (a sickle is always a '
        'sickle, boots are always boots); the shape note only varies that kind. '
                'Cells, in reading order: %s. '
        'IMPORTANT: no captions, no labels, no item names, no text or letters of any kind anywhere in the '
        'image; each cell holds only the painted item on its backdrop, its lower edge empty backdrop. No logos, '
        'no watermark, no outer frame.'
        % (GRID, GRID, GRID * GRID, '; '.join(parts))
    )


def gutters(gray, axis):
    """The GRID-1 dark lines across axis 0 (x) or 1 (y), each as (start, end)."""
    from PIL import ImageStat
    w, h = gray.size
    n = w if axis == 0 else h

    def mean(i):
        box = (i, 0, i + 1, h) if axis == 0 else (0, i, w, i + 1)
        return ImageStat.Stat(gray.crop(box)).mean[0]
    out = []
    for g in range(1, GRID):
        c = n * g // GRID
        lo, hi = c - n // (GRID * 6), c + n // (GRID * 6)
        best = min(range(lo, hi), key=mean)
        floor, a, b = mean(best), best, best
        while a > lo and mean(a - 1) < floor + 14:
            a -= 1
        while b < hi and mean(b + 1) < floor + 14:
            b += 1
        out.append((a, b))
    return out


def cut(numbers):
    from PIL import Image
    sheets = plan()
    os.makedirs(DST, exist_ok=True)
    for path in sorted(glob.glob(os.path.join(SRC, 'sheet[0-9][0-9][0-9].*'))):
        no = int(os.path.basename(path)[5:8])
        if numbers and no not in numbers:
            continue
        im = Image.open(path).convert('RGB')
        gray = im.convert('L')
        w, h = im.size
        xs = [0] + [p for a, b in gutters(gray, 0) for p in (a, b + 1)] + [w]
        ys = [0] + [p for a, b in gutters(gray, 1) for p in (a, b + 1)] + [h]
        inset = max(2, w // 200)
        for k, item in enumerate(sheets[no]):
            r, c = divmod(k, GRID)
            box = (xs[2 * c] + inset, ys[2 * r] + inset, xs[2 * c + 1] - inset, ys[2 * r + 1] - inset)
            cell = im.crop(box)
            side = min(cell.size)
            x0, y0 = (cell.width - side) // 2, (cell.height - side) // 2
            cell = cell.crop((x0, y0, x0 + side, y0 + side)).resize((SIZE, SIZE), Image.LANCZOS)
            cell.save(os.path.join(DST, item['code'] + '.webp'), 'WEBP', quality=82, method=6)
        print('sheet%03d' % no, len(sheets[no]), 'items')


if __name__ == '__main__':
    cmd = sys.argv[1] if len(sys.argv) > 1 else 'plan'
    if cmd == 'plan':
        for no, s in enumerate(plan()):
            print('%03d' % no, s[0]['code'], '..', s[-1]['code'], len(s))
    elif cmd == 'prompt':
        print(prompt(plan()[int(sys.argv[2])]))
    elif cmd == 'cut':
        cut({int(a) for a in sys.argv[2:]})
