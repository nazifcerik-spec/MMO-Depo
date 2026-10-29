"""Avatar art (DD-324): client/art_originals/avatars/sheetNN.jpg (2x2 grids,
git-ignored) -> client/assets/avatars/avXX.webp.

Sheet N holds avatars 4N+1..4N+4 in reading order (top-left, top-right,
bottom-left, bottom-right). The black gutters are found per sheet, each
portrait is trimmed of dark edges, cut square around its centre and saved at
256 px (the game shows avatars at most ~64 dp).

    python scripts/avatar_assets.py
"""
import glob
import os

from PIL import Image, ImageStat

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, 'client', 'art_originals', 'avatars')
DST = os.path.join(ROOT, 'client', 'assets', 'avatars')
SIZE = 256
DARK = 14  # mean brightness of a gutter / frame line
INSET = 6  # pixels shaved past a gutter so no line remains


def gutter(gray, axis):
    """(start, end) of the dark line crossing the middle along axis 0 (x) or 1 (y)."""
    w, h = gray.size
    n = w if axis == 0 else h
    def mean(i):
        box = (i, 0, i + 1, h) if axis == 0 else (0, i, w, i + 1)
        return ImageStat.Stat(gray.crop(box)).mean[0]
    lo, hi = int(n * 0.4), int(n * 0.6)
    best = min(range(lo, hi), key=mean)
    floor = mean(best)
    a = b = best
    while a > lo and mean(a - 1) < floor + DARK:
        a -= 1
    while b < hi and mean(b + 1) < floor + DARK:
        b += 1
    return a, b


def trim(im):
    """Drops near-black rows/columns at the edges (outer frame lines)."""
    gray = im.convert('L')
    w, h = gray.size
    col = lambda x: ImageStat.Stat(gray.crop((x, 0, x + 1, h))).mean[0]
    row = lambda y: ImageStat.Stat(gray.crop((0, y, w, y + 1))).mean[0]
    l, r, t, b = 0, w - 1, 0, h - 1
    while l < w // 4 and col(l) < DARK:
        l += 1
    while r > w * 3 // 4 and col(r) < DARK:
        r -= 1
    while t < h // 4 and row(t) < DARK:
        t += 1
    while b > h * 3 // 4 and row(b) < DARK:
        b -= 1
    return im.crop((l, t, r + 1, b + 1))


os.makedirs(DST, exist_ok=True)
total = 0
for path in sorted(glob.glob(os.path.join(SRC, 'sheet[0-9][0-9].*'))):
    sheet = int(os.path.basename(path)[5:7])
    im = Image.open(path).convert('RGB')
    gray = im.convert('L')
    gx, gy = gutter(gray, 0), gutter(gray, 1)
    w, h = im.size
    boxes = [
        (0, 0, gx[0] - INSET, gy[0] - INSET),
        (gx[1] + 1 + INSET, 0, w, gy[0] - INSET),
        (0, gy[1] + 1 + INSET, gx[0] - INSET, h),
        (gx[1] + 1 + INSET, gy[1] + 1 + INSET, w, h),
    ]
    for k, box in enumerate(boxes):
        part = trim(im.crop(box))
        side = min(part.size)
        x0, y0 = (part.width - side) // 2, (part.height - side) // 2
        part = part.crop((x0, y0, x0 + side, y0 + side)).resize((SIZE, SIZE), Image.LANCZOS)
        code = f'av{4 * sheet + k + 1:02d}'
        out = os.path.join(DST, code + '.webp')
        part.save(out, 'WEBP', quality=82, method=6)
        total += os.path.getsize(out)
    print(f'sheet{sheet:02d} gutters x={gx} y={gy}')
print('total', total // 1024, 'KB')
