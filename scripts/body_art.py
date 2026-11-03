"""Full-body portrait art (DD-344): per race and body, the same figure in an
outfit for every chest armor type (cloth, leather, plate) and level band,
plus one without armor. Gemini draws two 5x2 sheets per body in one chat;
this script prints their prompts and cuts the saved sheets.

    python scripts/body_art.py prompt human male A
    python scripts/body_art.py cut human male A   # client/art_originals/bodies/human_male_A.png
Writes client/assets/bodies/<race>_<body>_<outfit>.webp, outfit = none or
<type><band> (cloth1..cloth6, leather1..6, plate1..6).
"""
import os
import sys

from PIL import Image

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
SRC = os.path.join(ROOT, 'client', 'art_originals', 'bodies')
OUT = os.path.join(ROOT, 'client', 'assets', 'bodies')

# Level bands of the chest piece (client bodyOutfit mirrors these).
BANDS = [1, 100, 300, 500, 700, 850]
W, H = 200, 320  # portrait frame
OUTFITS = {
    'none': 'a simple linen shirt, plain trousers and boots, no armor',
    'cloth1': 'an apprentice grey wool robe with a rope belt',
    'cloth2': 'a blue mage robe with silver trim',
    'cloth3': 'crimson sorcerer robes with gold embroidery',
    'cloth4': 'deep violet arcane robes with glowing rune patterns',
    'cloth5': 'white-and-gold celestial robes with soft floating light motes',
    'cloth6': 'legendary cosmic robes woven of starlight with a glowing aura',
    'leather1': 'a simple brown leather jerkin with straps',
    'leather2': 'dark green ranger leather armor with a hood down',
    'leather3': 'black studded leather armor with a fur collar',
    'leather4': 'crimson dragonhide leather armor',
    'leather5': 'midnight-blue shadow leather armor with silver filigree',
    'leather6': 'legendary void-touched leather armor with glowing purple seams',
    'plate1': 'simple iron plate armor',
    'plate2': 'polished steel full plate armor',
    'plate3': 'ornate blue-and-silver plate armor with a cape',
    'plate4': 'black-and-red dragon plate armor with horned pauldrons',
    'plate5': 'white-and-gold holy plate armor with wing motifs',
    'plate6': 'legendary golden plate armor with glowing runes and a flowing cape',
}
SHEETS = {
    'A': ['none', 'cloth1', 'cloth2', 'cloth3', 'cloth4', 'cloth5', 'cloth6', 'leather1', 'leather2', 'leather3'],
    'B': ['leather4', 'leather5', 'leather6', 'plate1', 'plate2', 'plate3', 'plate4', 'plate5', 'plate6', 'none'],
}
LOOKS = {
    ('human', 'male'): 'human male adventurer (short brown hair, light beard, athletic build)',
    ('human', 'female'): 'human female adventurer (long auburn hair tied back, athletic build)',
}


def prompt(race, body, sheet):
    looks = '; '.join('%d) %s' % (i + 1, OUTFITS[o]) for i, o in enumerate(SHEETS[sheet]))
    again = ('Now a new image of the SAME character as in the previous image (same face, hair, body and '
             'proportions), same layout and pose. ') if sheet != 'A' else 'Create an image. '
    return (again + 'Fantasy MMORPG character equipment sheet, painted digital illustration in the style of '
            'high-end mobile RPG character art: a grid of 5 columns and 2 rows of tall panels, evenly spaced, '
            'each panel shows the SAME %s (same face, same hair, same body, same height and scale) standing in '
            'the exact same relaxed front-facing full-body pose, arms slightly away from the body, hands empty '
            '(no weapons, no shields), whole figure visible from head to feet with a little space above and '
            'below, centered in the panel, on a plain dark charcoal background with a soft ground shadow. Only '
            'the outfit changes between panels, left to right, top row then bottom row: %s. IMPORTANT: no '
            'captions, no labels, no numbers, no text or letters anywhere in the image.' % (LOOKS[race, body], looks))


def runs(flags):
    """Spans of consecutive True flags longer than 20."""
    out, start = [], None
    for i, f in enumerate(flags + [False]):
        if f and start is None:
            start = i
        elif not f and start is not None:
            if i - start > 20:
                out.append((start, i))
            start = None
    return out


def panel_line(vals):
    """A line through panels holds some dark-charcoal background; a gutter
    line is all black or all pale paper."""
    return min(vals) < 100 and max(vals) > 22


def cut(race, body, sheet):
    """Panels are found by their black gutters (Gemini spaces them unevenly),
    the figure in each by its difference from the panel's background; every
    figure is scaled to the same height in a W x H frame, feet on the same
    line."""
    im = Image.open(os.path.join(SRC, '%s_%s_%s.png' % (race, body, sheet))).convert('RGB')
    g = im.convert('L')
    w, h = g.size
    px = g.load()
    os.makedirs(OUT, exist_ok=True)
    i = 0
    for half in ((0, h // 2), (h // 2, h)):  # each row of panels: columns first, then each one's height
        mid = (half[0] + half[1]) // 2
        cols = runs([panel_line([px[x, y] for y in range(half[0], half[1], 2)]) for x in range(w)])
        assert len(cols) == 5, cols
        for left, right in cols:
            xs = range(left, right, 3)
            span = runs([panel_line([px[x, y] for x in xs]) for y in range(h)])
            top, bottom = max((r for r in span if r[0] < mid < r[1]), key=lambda r: r[1] - r[0])
            outfit = SHEETS[sheet][i]
            i += 1
            if sheet == 'B' and outfit == 'none':
                continue  # the identity anchor; sheet A's is kept
            panel = im.crop((left + 3, top + 3, right - 3, bottom - 3))
            pg = panel.convert('L')
            hist, acc, bgv = pg.histogram(), 0, 0
            while acc < pg.width * pg.height / 3:  # background brightness: the darker third
                acc, bgv = acc + hist[bgv], bgv + 1
            mask = pg.point(lambda v: 255 if abs(v - bgv) > 22 else 0)
            x0, y0, x1, y1 = mask.getbbox()
            fh = y1 - y0
            # Frame: the figure fills 90% of the height, feet near the bottom.
            fw_out = round(fh / 0.9 * W / H)
            cx = (x0 + x1) / 2
            box = (round(cx - fw_out / 2), round(y1 - fh / 0.9 * 0.96), round(cx + fw_out / 2), round(y1 + fh / 0.9 * 0.04))
            frame = Image.new('RGB', (box[2] - box[0], box[3] - box[1]), panel.getpixel((2, 2)))
            frame.paste(panel.crop((max(0, box[0]), max(0, box[1]), min(panel.width, box[2]), min(panel.height, box[3]))),
                        (max(0, -box[0]), max(0, -box[1])))
            frame.resize((W, H), Image.LANCZOS).save(os.path.join(OUT, '%s_%s_%s.webp' % (race, body, outfit)), 'WEBP', quality=88)
    print('cut', race, body, sheet)


if __name__ == '__main__':
    cmd, race, body, sheet = sys.argv[1:5]
    print(prompt(race, body, sheet)) if cmd == 'prompt' else cut(race, body, sheet)
