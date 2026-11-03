"""Bounty quest medallions (DD-343) from the mob's own medallion: a red
bounty ring inside the gold frame, like the Gemini bounty art (art_sets.py
quests). For the bounties whose mob joined the hunting ladder later.

    python scripts/bounty_art.py
"""
import os

from PIL import Image, ImageDraw, ImageFilter

ASSETS = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'client', 'assets')
MOBS = ['bog_witch', 'dwarf_ghost', 'war_bear', 'mummy_priest', 'siren', 'shadow_elf', 'volcano_shaman',
        'dark_angel', 'wind_elemental', 'light_hunter', 'ancient_golem']

for mob in MOBS:
    im = Image.open(os.path.join(ASSETS, 'mobs', mob + '.webp')).convert('RGB')
    w, h = im.size
    s = 4  # supersample
    ring = Image.new('L', (w * s, h * s), 0)
    d = ImageDraw.Draw(ring)
    cx, cy = w * s / 2, h * s / 2
    r = w * s * 0.405  # just inside the gold frame
    d.ellipse([cx - r, cy - r, cx + r, cy + r], outline=255, width=int(w * s * 0.035))
    ring = ring.resize((w, h), Image.LANCZOS)
    glow = ring.filter(ImageFilter.GaussianBlur(w * 0.04))
    red = Image.new('RGB', (w, h), (220, 30, 30))
    hot = Image.new('RGB', (w, h), (255, 120, 110))
    out = Image.composite(red, im, glow.point(lambda v: int(v * 0.85)))
    out = Image.composite(hot, out, ring.point(lambda v: int(v * 0.7)))
    out.save(os.path.join(ASSETS, 'quests', 'side_bounty_%s.webp' % mob), 'WEBP', quality=90)
    print(mob)
