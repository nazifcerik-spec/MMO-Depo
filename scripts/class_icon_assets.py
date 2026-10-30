"""Class icon art (DD-325): client/art_originals/class_icons/sheetNN.jpg (2x2
grids of round medallions, git-ignored) -> client/assets/class_icons/<class>.webp.

Each quarter's medallion is found as the bounding box of its bright rows and
columns (the background is near-black), cut square around it and saved at
192 px (the game shows class icons at most ~56 dp, clipped to a circle).

    python scripts/class_icon_assets.py
"""
import os

from PIL import Image

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, 'client', 'art_originals', 'class_icons')
DST = os.path.join(ROOT, 'client', 'assets', 'class_icons')
SIZE = 192
BRIGHT = 70  # a medallion pixel; the backdrop stays well below
SHARE = 0.06  # a row/column belongs to the medallion past this bright share

# Sheet file -> classes in reading order; None leaves a quarter unused.
SHEETS = {
    'sheet00': ['warrior', 'mage', 'cleric', 'rogue'],
    'sheet01': ['ranger', 'paladin', 'shaman', 'druid'],
    'sheet02': ['monk', 'bard', None, None],
}


def medallion(part):
    """The square around part's bright medallion."""
    gray = part.convert('L').point(lambda v: 255 if v > BRIGHT else 0)
    w, h = gray.size
    px = gray.load()
    cols = [x for x in range(w) if sum(px[x, y] > 0 for y in range(0, h, 2)) > SHARE * h / 2]
    rows = [y for y in range(h) if sum(px[x, y] > 0 for x in range(0, w, 2)) > SHARE * w / 2]
    l, r, t, b = cols[0], cols[-1], rows[0], rows[-1]
    side = max(r - l, b - t) + 1
    cx, cy = (l + r) // 2, (t + b) // 2
    x0 = max(0, min(w - side, cx - side // 2))
    y0 = max(0, min(h - side, cy - side // 2))
    return part.crop((x0, y0, x0 + side, y0 + side))


os.makedirs(DST, exist_ok=True)
total = 0
for name, classes in SHEETS.items():
    im = Image.open(os.path.join(SRC, name + '.jpg')).convert('RGB')
    w, h = im.size
    for k, code in enumerate(classes):
        if code is None:
            continue
        x, y = (k % 2) * w // 2, (k // 2) * h // 2
        icon = medallion(im.crop((x, y, x + w // 2, y + h // 2))).resize((SIZE, SIZE), Image.LANCZOS)
        out = os.path.join(DST, code + '.webp')
        icon.save(out, 'WEBP', quality=85, method=6)
        total += os.path.getsize(out)
        print(code, icon.size)
print('total', total // 1024, 'KB')
