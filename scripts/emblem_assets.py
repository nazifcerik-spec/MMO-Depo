"""Emblem art: client/art_originals/<dir>/sheetNN.jpg (2x2 grids of round
medallions, git-ignored) -> client/assets/<dir>/<code>.webp. Class icons
(DD-325) and XP boost boxes (DD-326).

Each quarter's medallion is found as the bounding box of its bright rows and
columns (the background is near-black), cut square around it and saved at
192 px (the game shows emblems at most ~88 dp, clipped to a circle).

    python scripts/emblem_assets.py
"""
import os

from PIL import Image

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
ART = os.path.join(ROOT, 'client', 'art_originals')
ASSETS = os.path.join(ROOT, 'client', 'assets')
SIZE = 192
BRIGHT = 70  # a medallion pixel; the backdrop stays well below
SHARE = 0.06  # a row/column belongs to the medallion past this bright share

# <dir>/<sheet> -> codes in reading order; None leaves a quarter unused.
SHEETS = {
    'class_icons/sheet00': ['warrior', 'mage', 'cleric', 'rogue'],
    'class_icons/sheet01': ['ranger', 'paladin', 'shaman', 'druid'],
    'class_icons/sheet02': ['monk', 'bard', None, None],
    'boosts/sheet00': ['gather', 'adventure', 'craft', None],
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


total = 0
for name, codes in SHEETS.items():
    folder = name.split('/')[0]
    os.makedirs(os.path.join(ASSETS, folder), exist_ok=True)
    im = Image.open(os.path.join(ART, name + '.jpg')).convert('RGB')
    w, h = im.size
    for k, code in enumerate(codes):
        if code is None:
            continue
        x, y = (k % 2) * w // 2, (k // 2) * h // 2
        icon = medallion(im.crop((x, y, x + w // 2, y + h // 2))).resize((SIZE, SIZE), Image.LANCZOS)
        out = os.path.join(ASSETS, folder, code + '.webp')
        icon.save(out, 'WEBP', quality=85, method=6)
        total += os.path.getsize(out)
        print(folder, code)
print('total', total // 1024, 'KB')
