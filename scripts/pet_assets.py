"""Pet art (DD-334): the pet design sheet's images become client assets.

The source folder holds one folder per pet, "NN_Name" in the order of
scripts/pet_data.py, each with 1_baslangic.png / 2_orta.png / 3_ust.png (the
three evolution stages, cut from a Gemini poster), yumurta.png (the egg) and
evrim_poster.png. It defaults to ~/Desktop/Pet_Evrim (copied on every run
to client/art_originals/pets, git-ignored) and falls back to that copy when
the desktop folder is gone. Writes
  client/assets/pets/<species>_<stage>.webp   stage art, 256 x 256
  client/assets/items/pet_egg_<species>.webp  egg icon, 128 x 128
A missing image is skipped (its art falls back to an icon in the game).

Run from the repository root: python scripts/pet_assets.py [source_dir]
"""
import glob
import os
import shutil
import sys

from PIL import Image, ImageStat

sys.path.insert(0, os.path.dirname(__file__))
from pet_data import SPECIES  # noqa: E402

CLIENT = os.path.join(os.path.dirname(__file__), '..', 'client')
ROOT = os.path.join(CLIENT, 'assets')
KEEP = os.path.join(CLIENT, 'art_originals', 'pets')
DESKTOP = os.path.join(os.path.expanduser('~'), 'Desktop', 'Pet_Evrim')
SRC = sys.argv[1] if len(sys.argv) > 1 else (DESKTOP if os.path.isdir(DESKTOP) else KEEP)
STAGES = ['1_baslangic.png', '2_orta.png', '3_ust.png']
PAD = 0.08  # margin added around art that touches its edges


def touches_edge(im):
    """True when the art runs into the top or a side (a busy strip, not
    backdrop); the bottom is often a painted floor and stays as it is."""
    g = im.convert('L')
    w, h = g.size
    return any(ImageStat.Stat(g.crop(b)).stddev[0] > 15 for b in ((0, 0, 3, h), (w - 3, 0, w, h), (0, 0, w, 3)))


def save(src, dst, size):
    im = Image.open(src).convert('RGB')
    side = max(im.size)
    if touches_edge(im):
        side = int(side * (1 + 2 * PAD))
    bg = tuple(int(v) for v in ImageStat.Stat(im.crop((0, 0, im.width, 4))).mean)
    sq = Image.new('RGB', (side, side), bg)
    sq.paste(im, ((side - im.width) // 2, (side - im.height) // 2))
    sq.resize((size, size), Image.LANCZOS).save(dst, 'WEBP', quality=88, method=6)


if SRC != KEEP and os.path.isdir(SRC):
    shutil.copytree(SRC, KEEP, dirs_exist_ok=True)
os.makedirs(os.path.join(ROOT, 'pets'), exist_ok=True)
stages = eggs = 0
missing = []
for n, (code, *_names) in enumerate(SPECIES, 1):
    folder = glob.glob(os.path.join(SRC, '%02d_*' % n))
    if not folder:
        missing.append(code)
        continue
    for s, f in enumerate(STAGES, 1):
        p = os.path.join(folder[0], f)
        if os.path.exists(p):
            save(p, os.path.join(ROOT, 'pets', '%s_%d.webp' % (code, s)), 256)
            stages += 1
        else:
            print('no stage', s, folder[0])
    egg = os.path.join(folder[0], 'yumurta.png')
    if os.path.exists(egg):
        save(egg, os.path.join(ROOT, 'items', 'pet_egg_%s.webp' % code), 128)
        eggs += 1
    else:
        print('no egg', folder[0])
print('stages:', stages, '/', 3 * len(SPECIES), ' eggs:', eggs, '/', len(SPECIES))
print('no art yet:', ', '.join(missing))
