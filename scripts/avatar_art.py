"""Avatar portraits (DD-324, redone in DD-330): 80 head-and-shoulders
portraits, each with its own age, expression, angle, hair, mark, headgear,
light and backdrop colour, so no two read as copies.

Avatar k (0-based) is class (k // 4) % 10; the first 40 are the free ones,
the last 40 their legendary versions (other races and genders). Sheets hold
GRID×GRID avatars in order, saved as client/art_originals/avatars/gridNN.jpg
(git-ignored) and cut into client/assets/avatars/avNN.webp.

    python scripts/avatar_art.py prompt N    # the prompt for sheet N
    python scripts/avatar_art.py cut [N ...] # cut the saved sheets
"""
import glob
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from item_art import gutters  # noqa: E402

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, 'client', 'art_originals', 'avatars')
DST = os.path.join(ROOT, 'client', 'assets', 'avatars')
GRID = 3
COUNT = 80
SIZE = 256

RACES = [
    ('human man', 'human woman'),
    ('green-skinned orc man with tusks', 'green-skinned orc woman with small tusks'),
    ('high elf man with long pointed ears', 'high elf woman with long pointed ears'),
    ('stout dwarf man', 'stout dwarf woman'),
    ('dark elf man with violet skin', 'dark elf woman with violet skin'),
    ('dragon-kin man with scales and horns', 'dragon-kin woman with scales and small horns'),
    ('small gnome man with big ears', 'small gnome woman with big ears'),
    ('bull-headed minotaur man', 'bull-headed minotaur woman'),
    ('blue-skinned jungle troll man with tusks', 'blue-skinned jungle troll woman with tusks'),
    ('undead man with pale grey skin', 'undead woman with pale grey skin'),
]
CLASSES = [
    'warrior in battle-worn plate', 'mage in arcane robes', 'cleric in holy vestments',
    'rogue in dark leathers', 'ranger in leather and fur', 'paladin in gold and blue plate',
    'shaman in tribal mail', 'druid in leafy robes', 'monk in martial wraps', 'bard in an embroidered doublet',
]
AGE = ['young', 'in their prime', 'weathered veteran', 'elderly', 'youthful', 'battle-hardened']
MOOD = ['stern', 'sly smirk', 'fierce battle shout', 'calm and serene', 'warm smile', 'brooding',
        'determined glare', 'amused grin', 'wary sidelong look', 'proud and noble']
VIEW = ['facing forward', 'three-quarter view turned left', 'three-quarter view turned right',
        'slight low angle, looking down at the viewer', 'head tilted, looking over the shoulder']
HAIR = ['long braided hair', 'shaved head', 'wild mane', 'short cropped hair', 'topknot', 'flowing silver hair',
        'mohawk', 'twin braids with beads', 'curly red hair', 'slicked-back black hair', 'hair hidden by headgear']
MARK = ['a scar across the eye', 'glowing tattoos', 'war paint stripes', 'an eyepatch', 'freckles', 'nose ring and earrings',
        'a braided beard or jewelled chin', 'face runes', 'a burn mark on the cheek', 'clean features']
GEAR = ['hood', 'horned helmet', 'circlet with a gem', 'feathered headdress', 'bare head', 'crown of thorns',
        'bandana', 'antler crown', 'wide-brimmed hat', 'winged helm', 'skull mask pushed up', 'veil']
LIGHT = ['warm firelight', 'cold moonlight', 'arcane blue glow', 'sunset orange', 'toxic green glow', 'holy golden light',
         'crimson ember light', 'teal sea light', 'violet starlight']
BACK = ['deep crimson', 'midnight blue', 'forest green', 'royal purple', 'burnt orange', 'teal', 'charcoal',
        'golden ochre', 'icy blue', 'wine red', 'olive', 'indigo', 'rust brown', 'jade', 'slate', 'plum']


def avatar(k):
    """The description of avatar k (0-based)."""
    sheet, slot = divmod(k, 4)  # the old 2x2 pairing of class, race and gender
    epic = k >= COUNT // 2
    shift = 1 if epic else 0
    race = RACES[(4 * sheet + slot + 2 * shift) % 10][(sheet + slot + shift) % 2]
    cls = CLASSES[sheet % 10]
    pick = lambda xs, step: xs[(k * step + sheet) % len(xs)]
    parts = [
        '%s %s' % (race, cls), pick(AGE, 5), pick(MOOD, 3), pick(VIEW, 2), pick(HAIR, 7), pick(MARK, 3),
        'wearing a ' + pick(GEAR, 5), pick(LIGHT, 4), '%s backdrop' % pick(BACK, 7),
    ]
    if epic:
        parts.append('legendary: ornate glowing runes, radiant aura, jewelled trims')
    return ', '.join(parts)


def prompt(n):
    first = n * GRID * GRID
    cells = ['row %d col %d: %s' % (j // GRID + 1, j % GRID + 1, avatar(first + j))
             for j in range(GRID * GRID) if first + j < COUNT]
    return (
        'Generate a square image: a %dx%d grid of %d separate square fantasy character portrait avatars of equal '
        'size, divided by thin solid black lines. Each is a head-and-shoulders portrait, face fully inside its cell. '
        'Stylized hand-painted high-fantasy MMORPG portrait art: clean bold shapes, soft cel shading, crisp rim light, '
        'expressive faces. Every portrait must be a clearly different individual: different face shape, age, '
        'expression, pose, hair, colours, headgear, lighting and plain backdrop colour as described, never the '
        'same face twice. Cells, in reading order: %s. IMPORTANT: no captions, no labels, no text or letters of '
        'any kind; no logos, no watermark, no outer frame.' % (GRID, GRID, len(cells), '; '.join(cells))
    )


def cut(numbers):
    from PIL import Image
    os.makedirs(DST, exist_ok=True)
    for path in sorted(glob.glob(os.path.join(SRC, 'grid[0-9][0-9].*'))):
        no = int(os.path.basename(path)[4:6])
        if numbers and no not in numbers:
            continue
        im = Image.open(path).convert('RGB')
        gray = im.convert('L')
        w, h = im.size
        xs = [0] + [p for a, b in gutters(gray, 0, GRID) for p in (a, b + 1)] + [w]
        ys = [0] + [p for a, b in gutters(gray, 1, GRID) for p in (a, b + 1)] + [h]
        inset = max(3, w // 160)
        for j in range(GRID * GRID):
            k = no * GRID * GRID + j
            if k >= COUNT:
                break
            r, c = divmod(j, GRID)
            cell = im.crop((xs[2 * c] + inset, ys[2 * r] + inset, xs[2 * c + 1] - inset, ys[2 * r + 1] - inset))
            side = min(cell.size)
            x0, y0 = (cell.width - side) // 2, (cell.height - side) // 2
            cell = cell.crop((x0, y0, x0 + side, y0 + side)).resize((SIZE, SIZE), Image.LANCZOS)
            cell.save(os.path.join(DST, 'av%02d.webp' % (k + 1)), 'WEBP', quality=82, method=6)
        print('grid%02d' % no)


if __name__ == '__main__':
    if sys.argv[1] == 'prompt':
        print(prompt(int(sys.argv[2])))
    elif sys.argv[1] == 'cut':
        cut({int(a) for a in sys.argv[2:]})
