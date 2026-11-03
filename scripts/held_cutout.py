"""Cut held items (weapons, shields) out of their icon's painted background
for the full-body portrait (DD-344): region-grow the smooth blue-grey
background from the border (the icons have a dark outline), then feather.

    python scripts/held_cutout.py            # every mainhand/offhand icon
    python scripts/held_cutout.py CODE ...   # only these
Writes client/assets/held/<code>.webp (RGBA).
"""
import collections
import glob
import io
import os
import sys

import yaml
from PIL import Image, ImageFilter

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ITEMS = os.path.join(ROOT, 'client', 'assets', 'items')
OUT = os.path.join(ROOT, 'client', 'assets', 'held')
STEP = 14  # largest colour step between neighbouring background pixels
HOLE = 90  # brightness band (r+g+b) around the background's mean
HOLE_TINT = 25  # blue over red that marks the background
SIZE = 96  # shown at most this big in the portrait
HOLE_MIN = 25  # pixels in a hole (smaller patches are shading on the item)


def held_codes():
    out = []
    for f in ('items.yaml', 'gear.yaml', 'sets.yaml'):
        with io.open(os.path.join(ROOT, 'content', f), encoding='utf-8') as fh:
            for it in (yaml.safe_load(fh) or {}).get('items') or []:
                if it.get('equip_slot') in ('mainhand', 'offhand'):
                    out.append(it['code'])
    return out


def cutout(src):
    im = Image.open(src).convert('RGB')
    w, h = im.size
    px = im.load()
    bg = [[False] * w for _ in range(h)]
    q = collections.deque()
    for x in range(w):
        q.extend([(x, 0), (x, h - 1)])
    for y in range(h):
        q.extend([(0, y), (w - 1, y)])
    for x, y in q:
        bg[y][x] = True
    while q:
        x, y = q.popleft()
        r, g, b = px[x, y]
        for nx, ny in ((x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)):
            if 0 <= nx < w and 0 <= ny < h and not bg[ny][nx]:
                r2, g2, b2 = px[nx, ny]
                if abs(r - r2) + abs(g - g2) + abs(b - b2) <= STEP and b2 >= r2:  # the background is blue-grey
                    bg[ny][nx] = True
                    q.append((nx, ny))
    # Holes inside the item (an axe's cut-outs, behind a bowstring) show the
    # same background: pixels close to its mean colour go too.
    n = sum(row.count(True) for row in bg)
    mr = sum(px[x, y][0] for y in range(h) for x in range(w) if bg[y][x]) / n
    mg = sum(px[x, y][1] for y in range(h) for x in range(w) if bg[y][x]) / n
    mb = sum(px[x, y][2] for y in range(h) for x in range(w) if bg[y][x]) / n
    def bluish(c):  # the background's blue-grey: blue well above red, green between
        r, g, b = c
        return b - r >= HOLE_TINT and r <= g <= b and abs(r + g + b - (mr + mg + mb)) <= HOLE

    near = [[not bg[y][x] and bluish(px[x, y]) for x in range(w)] for y in range(h)]
    seen = [[False] * w for _ in range(h)]
    for y0 in range(h):
        for x0 in range(w):
            if not near[y0][x0] or seen[y0][x0]:
                continue
            comp, q = [], collections.deque([(x0, y0)])
            seen[y0][x0] = True
            while q:  # one hole: a connected patch of background-coloured pixels
                x, y = q.popleft()
                comp.append((x, y))
                for nx, ny in ((x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)):
                    if 0 <= nx < w and 0 <= ny < h and near[ny][nx] and not seen[ny][nx]:
                        seen[ny][nx] = True
                        q.append((nx, ny))
            if len(comp) >= HOLE_MIN:
                for x, y in comp:
                    bg[y][x] = True
    mask = Image.new('L', (w, h), 0)
    m = mask.load()
    for y in range(h):
        for x in range(w):
            if not bg[y][x]:
                m[x, y] = 255
    mask = mask.filter(ImageFilter.MinFilter(3)).filter(ImageFilter.GaussianBlur(0.8))
    out = im.convert('RGBA')
    out.putalpha(mask)
    return out


def main(codes):
    os.makedirs(OUT, exist_ok=True)
    done = 0
    for code in codes:
        src = os.path.join(ITEMS, code + '.webp')
        if not os.path.exists(src):
            continue
        cutout(src).resize((SIZE, SIZE), Image.LANCZOS).save(os.path.join(OUT, code + '.webp'), 'WEBP', quality=82)
        done += 1
    print('held', done)


if __name__ == '__main__':
    main(sys.argv[1:] or held_codes())
