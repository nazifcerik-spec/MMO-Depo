"""Page background art (DD-323): client/art_originals/backgrounds/bgNN.* ->
client/assets/backgrounds/bgNN.webp.

The originals (git-ignored) are portrait ~9:16. They are kept at their own
size (at most 768 px wide; upscaling adds bytes, not detail) and stored as
WebP. The app covers the screen with them, cropping the edges a phone's
aspect ratio does not show.

    python scripts/bg_assets.py
"""
import glob
import os

from PIL import Image

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, 'client', 'art_originals', 'backgrounds')
DST = os.path.join(ROOT, 'client', 'assets', 'backgrounds')
MAX_WIDTH = 768

os.makedirs(DST, exist_ok=True)
total = 0
for path in sorted(glob.glob(os.path.join(SRC, 'bg[0-9][0-9].*'))):
    code = os.path.splitext(os.path.basename(path))[0]
    im = Image.open(path).convert('RGB')
    if im.width > MAX_WIDTH:
        im = im.resize((MAX_WIDTH, round(im.height * MAX_WIDTH / im.width)), Image.LANCZOS)
    out = os.path.join(DST, code + '.webp')
    im.save(out, 'WEBP', quality=80, method=6)
    total += os.path.getsize(out)
    print(code, im.size, os.path.getsize(out) // 1024, 'KB')
print('total', total // 1024, 'KB')
