"""Item art for every equippable item (DD-329): weapons, armor, shields,
jewelry and raid sets, each its own painted icon.

The items are ordered (family, then level) and packed GRID×GRID to a sheet;
each sheet is one image generated from the prompt this script prints, saved
as client/art_originals/items/sheetNNN.jpg (git-ignored), then cut into
client/assets/items/<code>.webp.

    python scripts/item_art.py plan          # sheets and their items
    python scripts/item_art.py prompt N      # the prompt for sheet N
    python scripts/item_art.py cut [N ...]   # cut the saved sheets (all by default)
    python scripts/item_art.py eggs-prompt   # pet eggs (DD-331): one sheet, eggs000.jpg
    python scripts/item_art.py eggs-cut
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
    'fist_weapon': 'clawed gauntlet fist weapon worn over the hand (never a blade on a handle)', 'holy_tome': 'holy spell tome', 'lute': 'magical lute',
    'mace': 'one-handed mace', 'nature_focus': 'druidic nature focus', 'rapier': 'rapier',
    'ritual_axe': 'shamanic ritual axe', 'short_sword': 'short sword', 'sickle': 'sickle',
    'staff': 'mage staff', 'sword': 'sword', 'tonfa': 'pair of tonfas (straight batons with a short side handle)', 'totem_focus': 'shaman totem',
    'two_handed_mace': 'two-handed warhammer', 'wand': 'magic wand', 'shield': 'shield',
    'necklace': 'necklace', 'ring': 'ring', 'trinket': 'charm trinket',
}
SLOT = {'helm': 'helmet', 'chest': 'chest armor', 'legs': 'leg armor', 'gloves': 'gloves', 'boots': 'boots'}
ARMOR = {'cloth_armor': 'cloth', 'leather_armor': 'leather', 'plate_armor': 'plate'}
# What each armor type is made of, slot by slot: the material in a name
# only trims and colours cloth and leather.
ARMOR_KIND = {
    'cloth_armor': {
        'helm': 'soft cloth hood or cowl, draped fabric, no metal', 'chest': 'flowing cloth robe (fabric only)',
        'legs': 'pair of loose cloth trousers, soft fabric, no metal', 'gloves': 'pair of empty cloth hand wraps shaped like fingerless gloves, wound fabric bands, no hands inside, no metal',
        'boots': 'cloth sandals / soft slippers (fabric only)'},
    'leather_armor': {
        'helm': 'leather cowl (stitched leather, no metal plates)', 'chest': 'leather jerkin (stitched leather, no metal plates)',
        'legs': 'pair of leather leggings (stitched leather, no metal plates)',
        'gloves': 'pair of empty leather gloves (stitched leather, no metal plates)',
        'boots': 'pair of leather boots (stitched leather, no metal plates)'},
    'plate_armor': {
        'helm': 'metal plate helmet', 'chest': 'metal plate breastplate', 'legs': 'metal plate leg guards',
        'gloves': 'metal plate gauntlets', 'boots': 'metal plate sabatons'},
}


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
        return ARMOR_KIND[sub][i['equip_slot']]
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
    if i['item_subclass'] in ('cloth_armor', 'leather_armor'):
        look = 'colours and trims taken from ' + look  # soft pieces stay fabric or leather
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
        'sickle, boots are always boots); the shape note only varies that kind. For cloth and leather pieces the material only sets '
        'colours, trims, embroidery and gems; they stay soft fabric or leather. '
                'Cells, in reading order: %s. '
        'IMPORTANT: no captions, no labels, no item names, no text or letters of any kind anywhere in the '
        'image; each cell holds only the painted item on its backdrop, its lower edge empty backdrop. No logos, '
        'no watermark, no outer frame.'
        % (GRID, GRID, GRID * GRID, '; '.join(parts))
    )


# Pet eggs (DD-331), weakest to best: each its creature's look, the glow and
# ornament rising with the rank. The sheet's last row is filler, not cut.
EGGS = [
    ('ember_fox', 'smooth oval egg of warm russet orange with glowing ember-red cracks and a curling fox-tail flame wisp'),
    ('moss_turtle', 'round squat egg with a hexagonal turtle-shell pattern in mossy green and bark brown, tufts of moss'),
    ('frost_owl', 'pale icy blue egg with white feather-shaped frost patterns and a few small ice crystals'),
    ('thunder_hare', 'tall narrow egg of stormy grey-violet with crackling yellow lightning zigzags'),
    ('shadow_cat', 'sleek glossy black egg with faint purple smoke wisps and two glowing amber cat-eye slits'),
    ('crystal_beetle', 'faceted egg of translucent teal crystal with an iridescent green-purple beetle-wing sheen'),
    ('storm_hawk', 'slate-blue egg streaked with white wind swirls, a single silver hawk feather beside it'),
    ('iron_boar', 'heavy dark iron-grey egg bound with riveted metal bands and two small ivory tusks'),
    ('spirit_wolf', 'softly glowing ghostly cyan egg, half transparent, a faint wolf-head rune inside'),
    ('sand_wyrmling', 'sandy gold egg with dune-ripple bands and tiny ridged dragon scales'),
    ('coral_serpent', 'coral pink and red egg encrusted with branching coral and small shells, a sea-serpent coil around it'),
    ('lava_salamander', 'black basalt egg with bright molten orange lava veins and rising smoke'),
    ('moon_stag', 'pearly silver-white egg with a crescent moon emblem and silver antler-shaped filigree'),
    ('void_raven', 'deep indigo egg that swallows light, purple starry void specks, black raven feathers around its base'),
    ('griffin_cub', 'cream and golden-brown egg, eagle-feather pattern above and tawny lion-fur texture below, gold trim'),
    ('frost_drake', 'large white-blue egg with dragon scales and jagged ice spikes, cold mist, bright glow'),
    ('celestial_lynx', 'midnight blue egg with glowing golden constellation lines and starlight tufts like lynx ears'),
    ('phoenix_chick', 'radiant gold and crimson egg wreathed in living flames with a phoenix-feather crest, strong glow'),
    ('star_dragon', 'majestic egg of deep violet dragon scales studded with glowing star gems and a nebula swirl, '
                    'held in an ornate gold claw cradle'),
    ('eternal_kirin', 'legendary jade and gold egg with flowing cloud patterns, a kirin horn crest and a rainbow aura, '
                      'resting on an ornate golden lotus cradle'),
]
EGG_FILLER = ['plain brown speckled egg', 'small white egg in straw', 'grey stone egg', 'blue robin egg',
              'cracked empty eggshell']


def eggs_prompt():
    parts = ['%s: %s' % (POS(k), d) for k, d in enumerate([d for _, d in EGGS] + EGG_FILLER)]
    return (
        'Generate a square image: a %dx%d grid of %d separate square fantasy RPG creature egg icons of equal size, '
        'divided by thin solid black lines. Each cell shows exactly one whole egg, centered, upright, on the same '
        'plain dark slate-blue backdrop with a soft vignette. Stylized hand-painted game icon art like polished '
        'mobile MMORPG inventory icons: clean bold silhouettes, soft cel shading, crisp rim light, ornate details, '
        'magical glow growing stronger from the first cell to the twentieth. Every egg must look unique: a '
        'different shape, size, shell texture, ornament and colour palette from every other cell, never a recolour. '
        'Each cell is always an egg (no hatched creature, no creature outside the shell). '
        'Cells, in reading order: %s. '
        'IMPORTANT: no captions, no labels, no names, no text or letters of any kind anywhere in the image; each '
        'cell holds only the painted egg on its backdrop, its lower edge empty backdrop. No logos, no watermark, '
        'no outer frame.' % (GRID, GRID, GRID * GRID, '; '.join(parts)))


def gutters(gray, axis, grid=GRID):
    """The grid-1 dark lines across axis 0 (x) or 1 (y), each as (start, end)."""
    from PIL import ImageStat
    w, h = gray.size
    n = w if axis == 0 else h

    def mean(i):
        box = (i, 0, i + 1, h) if axis == 0 else (0, i, w, i + 1)
        return ImageStat.Stat(gray.crop(box)).mean[0]
    out = []
    for g in range(1, grid):
        c = n * g // grid
        lo, hi = c - n // (grid * 6), c + n // (grid * 6)
        best = min(range(lo, hi), key=mean)
        floor, a, b = mean(best), best, best
        while a > lo and mean(a - 1) < floor + 14:
            a -= 1
        while b < hi and mean(b + 1) < floor + 14:
            b += 1
        out.append((a, b))
    return out


def cut_sheet(path, codes):
    """Cuts one saved sheet into <code>.webp, cells in reading order."""
    from PIL import Image
    os.makedirs(DST, exist_ok=True)
    im = Image.open(path).convert('RGB')
    gray = im.convert('L')
    w, h = im.size
    xs = [0] + [p for a, b in gutters(gray, 0) for p in (a, b + 1)] + [w]
    ys = [0] + [p for a, b in gutters(gray, 1) for p in (a, b + 1)] + [h]
    inset = max(2, w // 200)
    for k, code in enumerate(codes):
        r, c = divmod(k, GRID)
        box = (xs[2 * c] + inset, ys[2 * r] + inset, xs[2 * c + 1] - inset, ys[2 * r + 1] - inset)
        cell = im.crop(box)
        side = min(cell.size)
        x0, y0 = (cell.width - side) // 2, (cell.height - side) // 2
        cell = cell.crop((x0, y0, x0 + side, y0 + side)).resize((SIZE, SIZE), Image.LANCZOS)
        cell.save(os.path.join(DST, code + '.webp'), 'WEBP', quality=82, method=6)


def cut(numbers):
    sheets = plan()
    for path in sorted(glob.glob(os.path.join(SRC, 'sheet[0-9][0-9][0-9].*'))):
        no = int(os.path.basename(path)[5:8])
        if numbers and no not in numbers:
            continue
        cut_sheet(path, [i['code'] for i in sheets[no]])
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
    elif cmd == 'eggs-prompt':
        print(eggs_prompt())
    elif cmd == 'eggs-cut':
        cut_sheet(os.path.join(SRC, 'eggs000.jpg'), ['pet_egg_' + c for c, _ in EGGS])
        print('eggs000', len(EGGS), 'eggs')
