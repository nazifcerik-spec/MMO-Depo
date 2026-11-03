"""Character customization layers (DD-346): per body one Gemini 5x2 sheet, top
row five bald faces, bottom row the same plain head under five hair styles.
Hair, brows and beards are painted chroma green and irises chroma magenta, so
this script splits them into layers the client tints, and aligns every head on
its eyes, so any face fits any hair:

    <body>_face<n>_base.webp    face and clothes, background removed
    <body>_face<n>_brows.webp   brows and beard, grey with alpha
    <body>_face<n>_eyes.webp    irises, grey with alpha
    <body>_hair<n>.webp         hair, grey with alpha

    python scripts/head_art.py prompt male
    python scripts/head_art.py cut male    # client/art_originals/heads/male.png
Writes client/assets/heads/. The art is original: reference pictures are only
described in words, never uploaded or shipped.
"""
import os
import sys

from PIL import Image, ImageFilter

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
SRC = os.path.join(ROOT, 'client', 'art_originals', 'heads')
OUT = os.path.join(ROOT, 'client', 'assets', 'heads')
GRIDS = {'male': (2, 1), 'female': (2, 1)}  # rows on the sheet, the row of the hair styles; faces are row 0
SIZE = 256  # the square frame of a head
SKIN_FIX = {('male', 4): (1.02, 0.90, 1.04)}  # per-channel: Gemini paints this olive skin too yellow
EYE_Y, EYE_GAP = 0.50, 0.215  # eye line and distance between the eyes in the frame

LOOKS = {
    'male': dict(
        who='adult man',
        faces=['light skin, clean-shaven, sharp jaw, narrow eyes',
               'tan skin, strong square jaw, a thin scar across one cheek',
               'dark brown skin, broad friendly face, full lips',
               'warm olive skin, short full beard, heavy brow',
               'pale weathered skin, mature face with crow\'s feet, trimmed goatee'],
        hair=['short hair swept to one side', 'undercut with a longer top',
              'wavy chin-length hair', 'long loose hair past the shoulders', 'shoulder-length braided locs']),
    'female': dict(
        who='adult woman',
        faces=['light skin, soft oval face, large eyes',
               'tan skin, high cheekbones, narrow confident eyes',
               'dark brown skin, rounded face, full lips',
               'olive skin, strong jaw, straight brows',
               'pale skin with freckles, mature composed face'],
        hair=['straight chin-length bob with bangs', 'high ponytail', 'very long straight hair',
              'wavy shoulder-length hair', 'one thick side braid']),
}


def prompt(body):
    k = LOOKS[body]
    faces = '; '.join('%d) %s' % (i + 1, f) for i, f in enumerate(k['faces']))
    hair = '; '.join('%d) %s' % (i + 6, h) for i, h in enumerate(k['hair']))
    return ('Create an image. Original fantasy RPG character creator sheet in a modern anime concept-art style: clean '
            'confident ink lineart, soft cel shading. A grid of 5 columns and 2 rows of equal panels separated by '
            'thin black lines. Every panel shows the head and shoulders of an %s, facing the viewer perfectly '
            'straight on, calm neutral expression, in exactly the same framing and scale, head centered with a clear '
            'empty margin above the hair, wearing the same plain grey round-neck tunic, on a flat plain medium-grey background. Top row, five different '
            'people, each completely bald with no hair on the head: %s. Bottom row, always the same plain person with '
            'light skin, only the hair changes: %s. COLOUR RULE, very important: all hair, eyebrows and facial hair '
            'are painted only in shades of pure saturated chroma green, and the irises in every panel are large, '
            'round and painted only in solid bright pure magenta with white around them, as if for a green screen; '
            'nothing else in the picture is green or magenta. No captions, no labels, no numbers, no text or letters anywhere in the image.'
            % (k['who'], faces, hair))


def is_green(c):
    return c[1] - max(c[0], c[2]) > 12


def is_pink(c):
    return min(c[0], c[2]) - c[1] > 12


def eyes_of(cell):
    """Centres of the two irises: the magenta left and right of the middle."""
    w, h = cell.size
    px = cell.load()
    sides = [[0, 0, 0], [0, 0, 0]]
    for y in range(h):
        for x in range(w):
            if is_pink(px[x, y]):
                s = sides[0 if x < w / 2 else 1]
                s[0] += x
                s[1] += y
                s[2] += 1
    assert sides[0][2] and sides[1][2], 'no irises found'
    return [(s[0] / s[2], s[1] / s[2]) for s in sides]


def align(cell):
    """The head moved and scaled so its eyes sit on the frame's eye line."""
    (lx, ly), (rx, ry) = eyes_of(cell)
    scale = EYE_GAP * SIZE / max(1, rx - lx)
    cx, cy = (lx + rx) / 2, (ly + ry) / 2
    # Output pixel (x, y) comes from source pixel (a*x + c, e*y + f).
    a = 1 / scale
    c = cx - a * SIZE / 2
    f = cy - a * EYE_Y * SIZE
    bg = cell.getpixel((3, 3))
    return cell.transform((SIZE, SIZE), Image.AFFINE, (a, 0, c, 0, a, f), Image.BICUBIC, fillcolor=bg)


def key(cell, test, amount):
    """The keyed pixels as a grey layer with alpha, highlights at white."""
    w, h = cell.size
    px = cell.load()
    grey, alpha = Image.new('L', (w, h)), Image.new('L', (w, h), 0)
    gp, ap = grey.load(), alpha.load()
    vals = []
    for y in range(h):
        for x in range(w):
            c = px[x, y]
            if test(c):
                gp[x, y] = max(c)
                ap[x, y] = min(255, amount(c) * 10)
                vals.append(max(c))
    if vals:
        top = sorted(vals)[int(len(vals) * 0.98)] or 255
        grey = grey.point(lambda v: min(255, v * 255 // top))
    # Close pinholes: a shaded strand must not let the face show through.
    closed = alpha.filter(ImageFilter.MaxFilter(3)).filter(ImageFilter.MinFilter(3))
    return Image.merge('LA', (grey, closed)).convert('RGBA'), alpha


def outlined(layer, alpha, cell):
    """Hair keeps the dark ink line around it."""
    w, h = cell.size
    grown = alpha.filter(ImageFilter.MaxFilter(3)).load()
    px, lp, ap = cell.load(), layer.load(), alpha.load()
    for y in range(h):
        for x in range(w):
            if grown[x, y] > 60 and ap[x, y] < 60 and max(px[x, y]) < 70:
                lp[x, y] = (0, 0, 0, 255)
    return layer


def base_of(cell, masks):
    """The face without its keyed parts and without the flat background."""
    w, h = cell.size
    px = cell.load()
    base = cell.convert('RGBA')
    bp = base.load()
    for mask in masks:
        mp = mask.load()
        for y in range(h):
            for x in range(w):
                if mp[x, y]:
                    r, g, b = px[x, y]
                    d = min(r, g, b)  # an edge pixel keeps its ink, not its key colour
                    bp[x, y] = (d, d, d, 255 - mp[x, y] // 2)
    bgv = px[3, 3]
    seen = [[False] * w for _ in range(h)]
    q = [(x, y) for x in range(w) for y in (0, h - 1)] + [(x, y) for y in range(h) for x in (0, w - 1)]
    for x, y in q:
        seen[y][x] = True
    while q:
        x, y = q.pop()
        for nx, ny in ((x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)):
            if 0 <= nx < w and 0 <= ny < h and not seen[ny][nx] and sum(
                    abs(px[nx, ny][i] - bgv[i]) for i in range(3)) < 36:
                seen[ny][nx] = True
                q.append((nx, ny))
    keep = Image.new('L', (w, h), 255)
    kp = keep.load()
    for y in range(h):
        for x in range(w):
            if seen[y][x]:
                kp[x, y] = 0
    keep = keep.filter(ImageFilter.GaussianBlur(0.8))
    base.putalpha(Image.composite(base.getchannel('A'), Image.new('L', (w, h), 0), keep))
    return base


def cell_at(im, rows, row, col):
    """One cell of the sheet's even grid, inside its black lines."""
    cw, ch = im.width / 5, im.height / rows
    return im.crop((round(col * cw) + 5, round(row * ch) + 5, round((col + 1) * cw) - 5, round((row + 1) * ch) - 5))


def save(img, name):
    img.save(os.path.join(OUT, name + '.webp'), 'WEBP', quality=88)


def cut(body):
    im = Image.open(os.path.join(SRC, body + '.png')).convert('RGB')
    os.makedirs(OUT, exist_ok=True)
    rows, hair_row = GRIDS[body]
    green = (is_green, lambda c: c[1] - max(c[0], c[2]))
    pink = (is_pink, lambda c: min(c[0], c[2]) - c[1])
    for n in range(1, 6):
        cell = align(cell_at(im, rows, 0, n - 1))
        brows, ba = key(cell, *green)
        eyes, ea = key(cell, *pink)
        base = base_of(cell, [ba, ea])
        if (body, n) in SKIN_FIX:
            r, g, b, a = base.split()
            fix = SKIN_FIX[body, n]
            base = Image.merge('RGBA', [c.point(lambda v, f=f: min(255, int(v * f))) for c, f in zip((r, g, b), fix)] + [a])
        save(base, '%s_face%d_base' % (body, n))
        save(brows, '%s_face%d_brows' % (body, n))
        save(eyes, '%s_face%d_eyes' % (body, n))
    for n in range(1, 6):
        cell = align(cell_at(im, rows, hair_row, n - 1))
        hair, ha = key(cell, *green)
        save(outlined(hair, ha, cell), '%s_hair%d' % (body, n))
    print('cut', body)


if __name__ == '__main__':
    print(prompt(sys.argv[2])) if sys.argv[1] == 'prompt' else cut(sys.argv[2])
