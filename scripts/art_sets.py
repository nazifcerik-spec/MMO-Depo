"""Picture sets beyond gear (DD-332): class path nodes, map regions, crafting
professions, helper NPCs and a few materials, each painted by Gemini in 3x3
sheets in the game's style and cut into client/assets/<dir>/<code>.webp.

    python scripts/art_sets.py plan
    python scripts/art_sets.py prompt <set> <sheet>   # Gemini prompt
    python scripts/art_sets.py cut <set> [sheet ...]  # client/art_originals/<set>/sheetNN.jpg

Medallion sets (nodes, regions, professions) match the class emblems
(DD-325): one round medallion with an ornate gold rim per cell, cut to the
medallion. Portraits match the avatars (DD-324); items match the gear icons
(DD-329) and land among them in assets/items.
"""
import os
import sys

from PIL import Image

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from emblem_assets import medallion  # noqa: E402
from item_art import gutters  # noqa: E402

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
GRID = 3
POS = ['top-left', 'top-middle', 'top-right', 'middle-left', 'center', 'middle-right',
       'bottom-left', 'bottom-middle', 'bottom-right']

# Base class palettes, as the class emblems use them.
PAL = {
    'warrior': 'steel and crimson', 'rogue': 'black and toxic green', 'ranger': 'forest green and brown leather',
    'mage': 'blue and violet', 'monk': 'orange and gold', 'cleric': 'gold and white',
    'paladin': 'gold and royal blue', 'druid': 'green and amber', 'bard': 'magenta and gold',
    'shaman': 'teal and orange',
}

# Promotions and specializations (DD-252), by base class: what each emblem shows.
NODES = [
    ('warrior', [
        ('guardian', 'a massive tower shield with a lion crest'),
        ('executioner', 'a huge executioner greataxe with a red glow'),
        ('iron_fortress', 'a fortress gate tower behind a riveted iron shield'),
        ('warlord', 'a war banner on a spear over two crossed swords'),
        ('berserker', 'a roaring horned helm wreathed in red rage flames'),
        ('weapon_master', 'a fan of five different weapons: sword, axe, spear, hammer and glaive'),
    ]),
    ('rogue', [
        ('assassin', 'a hooded skull mask with a dagger behind it'),
        ('nightblade', 'a curved black blade under a crescent moon'),
        ('poison_master', 'a bubbling poison vial wrapped by a green snake'),
        ('duelist', 'two crossed rapiers with a red rose'),
        ('blade_dancer', 'twin spinning blades inside a swirling silk ribbon'),
        ('trickster', 'a grinning jester mask with fanned playing cards'),
    ]),
    ('ranger', [
        ('marksman', 'a bullseye target pierced by an arrow'),
        ('sharpshooter', 'a fierce eagle eye inside a crosshair ring'),
        ('crossbow_master', 'a heavy loaded crossbow'),
        ('scout', 'a brass spyglass over a compass rose'),
        ('beastmaster', 'a wolf head side by side with a hawk'),
        ('trapper', 'a steel bear trap with a rope snare'),
    ]),
    ('mage', [
        ('arcanist', 'an open spellbook with floating glowing runes'),
        ('rune_master', 'a carved standing runestone glowing with runes'),
        ('chronomancer', 'a clock face with swirling sands of time'),
        ('elementalist', 'four small orbs of fire, water, earth and air circling each other'),
        ('pyromancer', 'a blazing fireball held in a clawed hand of flame'),
        ('cryomancer', 'a jagged ice crystal snowflake with frost mist'),
    ]),
    ('monk', [
        ('disciple', 'a calm meditating figure seated on a lotus'),
        ('iron_body', 'an iron-banded fist pressed on a cracked stone'),
        ('spirit_walker', 'a glowing spirit wisp above a ring of prayer beads'),
        ('fighter', 'two cloth-wrapped fists bumping together'),
        ('combo_master', 'a fist trailing a chain of three impact rings'),
        ('storm_fist', 'a fist crackling with blue lightning'),
    ]),
    ('cleric', [
        ('saint', 'a radiant halo above two praying hands'),
        ('high_priest', 'an ornate golden mitre with a sun emblem'),
        ('aegis_priest', 'a glowing dome of holy light over a small sun sigil'),
        ('oracle', 'an all-seeing eye inside a starry crystal ball'),
        ('prophet', 'an unrolled glowing scroll with a bright star'),
        ('exorcist', 'a golden chain binding a dark shadow wraith'),
    ]),
    ('paladin', [
        ('protector', 'a kite shield with spread golden wings'),
        ('aegis_knight', 'a knight great helm in front of a radiant round shield'),
        ('lightwarden', 'a lantern of holy light hanging over a shield'),
        ('crusader', 'a raised longsword with a long blue banner'),
        ('templar', 'a white temple pillar with a sword leaning on it'),
        ('justicar', 'golden scales of justice balanced on a sword'),
    ]),
    ('druid', [
        ('grove_keeper', 'a great oak tree with glowing roots'),
        ('life_binder', 'a sprouting seedling cupped in vines around a heart of light'),
        ('thorn_sage', 'a twisted thorny vine with a purple bloom'),
        ('shapeshifter', 'a face split half wolf and half human'),
        ('bearwarden', 'a roaring bear head armoured in bark'),
        ('mooncaller', 'a howling wolf under a huge full moon'),
    ]),
    ('bard', [
        ('minstrel', 'a lute with floating musical notes'),
        ('virtuoso', 'a violin and bow with golden sparkles'),
        ('harmonist', 'a small harp with glowing green healing notes'),
        ('skald', 'a curved war horn with knotwork carvings'),
        ('war_chanter', 'a war drum with two crossed drumsticks'),
        ('dirge_singer', 'a black raven perched on a broken lyre'),
    ]),
    ('shaman', [
        ('totemist', 'a carved totem pole with eagle wings'),
        ('earthwarden', 'a huge stone fist rising out of the ground'),
        ('storm_totemist', 'a totem struck by lightning under storm clouds'),
        ('spirit_caller', 'a ghostly ancestral spirit mask'),
        ('ancestor_sage', 'an elder spirit face in smoke with feathers'),
        ('hexer', 'a cursed charm of bones and beads with purple glow'),
    ]),
]

# Map regions (DD-279): a painted landscape vignette inside each medallion.
REGIONS = [
    ('sundermere', 'a green island continent with rolling hills and a castle seen from above'),
    ('ironhold_city', 'a stone fortress city with iron gates and smoking chimneys'),
    ('greenwood', 'a lush green forest with sunbeams'),
    ('whispering_reach', 'a misty river valley with tall reeds'),
    ('crypt_gates', 'a crypt gate carved with skulls'),
    ('emberpeak', 'a volcanic mountain peak glowing with embers'),
    ('ashen_hollow', 'a cave of grey ash with a red glow inside'),
    ('broken_crown', 'a cracked golden crown above a kingdom of castles'),
    ('dragonspine_isles', 'a chain of islands shaped like a dragon spine in the sea'),
    ('starless_lands', 'dark plains under a black starless sky'),
    ('sky_rift', 'floating rocks around a glowing tear in the sky'),
    ('void_shore', 'a dark shore against a purple cosmic void'),
    ('misty_marshes', 'a foggy swamp with dead trees'),
    ('sunken_temple', 'a half-sunken stone temple in green water'),
    ('iron_mountains', 'grey jagged mountains with a mine cart track'),
    ('forgotten_mine', 'an abandoned mine entrance with old lanterns'),
    ('northern_keep_ruins', 'snowy castle ruins'),
    ('northern_keep', 'a frozen keep gate with tattered banners'),
    ('crown_plains', 'golden wheat plains with a royal road'),
    ('rebel_camp', 'tents and campfires with a torn banner'),
    ('burnt_forest', 'charred black trees with glowing embers'),
    ('fire_root', 'a cavern of burning tree roots'),
    ('salt_desert', 'white salt flats with cracked ground'),
    ('buried_palace', 'a palace dome half buried in sand'),
    ('crown_graveyard', 'a graveyard of tombstones topped with crowns'),
    ('royal_crypt', 'an ornate royal crypt door with gold trim'),
    ('storm_coast', 'stormy sea cliffs with lightning'),
    ('sunken_fleet', 'old shipwrecks under blue water'),
    ('dragon_slopes', 'rocky slopes scattered with dragon bones'),
    ('brood_cavern', 'a cavern nest with dragon eggs'),
    ('mist_forest', 'a misty blue forest of ancient trees'),
    ('lost_grove', 'a hidden glowing grove with a stone circle'),
    ('crystal_caves', 'a cave full of blue crystals'),
    ('crystal_heart', 'a giant glowing crystal heart in a cavern'),
    ('volcanic_isle', 'an island volcano erupting'),
    ('volcano_heart', 'a lava chamber inside a volcano'),
    ('throne_peaks', 'snowy mountain peaks with a throne carved in the rock'),
    ('dragon_throne', 'a huge throne made of a dragon skull'),
    ('dark_plains', 'dark plains with dead grass under a black moon'),
    ('darkened_observatory', 'a ruined observatory dome with a broken telescope'),
    ('lost_city', 'overgrown ruined city towers'),
    ('clock_tower', 'a crooked clock tower with glowing gears'),
    ('floating_isles', 'floating islands with waterfalls in the sky'),
    ('storm_fortress', 'a fortress inside thunderclouds'),
    ('rift_edge', 'a cliff edge facing a glowing purple sky rift'),
    ('rift_core', 'a swirling purple vortex core'),
    ('celestial_temple', 'a white and gold temple above the clouds'),
    ('celestial_dome', 'a starlit dome with angel statues'),
    ('ash_sea', 'a grey sea of ash with burning ships'),
    ('sunken_throne', 'a throne half drowned in ash waves'),
    ('last_bastion', 'towering bastion walls at the edge of a dark void'),
    ('bastion_keep', 'a fortified keep gate glowing with blue wards'),
    ('edge_of_eternity', 'a stairway rising into endless starlight'),
    ('eternal_throne', 'a cosmic throne made of stars'),
]

# Crafting professions (Workshop list).
PROFESSIONS = [
    ('blacksmithing', 'a hammer striking a glowing sword blade on an anvil'),
    ('leatherworking', 'a tanned leather hide with an awl and stitching thread'),
    ('tailoring', 'a needle, a spool of thread and a folded cloth roll'),
    ('woodworking', 'a carving chisel and a wooden bow stave with curled shavings'),
    ('jewelcrafting', 'a cut gemstone with a jeweller loupe and a gold ring'),
    ('alchemy', 'a round bubbling flask with colourful vapour'),
    ('cooking', 'a steaming cooking pot with a ladle'),
    ('smelting', 'a crucible pouring molten metal into an ingot mould'),
    ('enchanting', 'a glowing parchment scroll above a magic circle'),
]

# Hired helpers (DD-319): head-and-shoulders portraits of the workers.
HELPERS = [
    ('miner', 'a stout bearded miner with a lamp helmet and a pickaxe over his shoulder'),
    ('woodsman', 'a burly lumberjack with a red plaid shirt and an axe'),
    ('angler', 'an old fisherman with a straw hat and a fishing rod'),
    ('herbalist', 'a young woman herbalist with a flower crown and a basket of herbs'),
    ('trapper', 'a hooded trapper in a fur cloak with a coil of rope'),
    ('smith', 'a muscular blacksmith in a leather apron with soot on his face'),
    ('tanner', 'a woman tanner in a leather apron holding a hide'),
    ('weaver', 'a woman weaver with a colourful shawl and thread spools'),
    ('carver', 'an elf woodcarver with wood shavings in her hair and a chisel'),
    ('apothecary', 'a bespectacled gnome apothecary holding a potion flask'),
    ('cook', 'a jolly round cook with a chef hat and a ladle'),
    ('smelter', 'an orc smelter with goggles lit by furnace glow'),
    ('scribe', 'an elderly scribe with a quill and a glowing scroll'),
    ('gemcutter', 'a dwarf woman gemcutter with a loupe and a sparkling gem'),
]

# Materials still without art; None = a filler cell, not cut.
MATERIALS = [
    ('elm_log', 'a short bundle of pale elm logs with light grey bark'),
    ('yew_log', 'a stack of reddish-brown yew logs with dark flaky bark'),
    ('linseed', 'a small burlap pouch spilling golden-brown flax seeds'),
    ('mandrake', 'a mandrake root with a gnarled humanlike root body and green leaves'),
    ('king_bloom', 'a regal golden and purple flower with a crown-shaped bloom'),
    ('mercury_ore', 'a chunk of grey rock with beads of liquid silver mercury'),
    (None, 'a coil of hemp rope'),
    (None, 'a lump of black coal'),
    (None, 'a bundle of white cotton'),
]


# Redraws in the last sheet's free cells: a later sheet's cut wins.
NODE_FIXES = [
    ('monk', 'disciple', 'a lotus blossom holding a floating glowing pearl'),
]


def nodes():
    return [(code, '%s, in %s' % (d, PAL[root])) for root, ns in NODES for code, d in ns] + [
        (code, '%s, in %s' % (d, PAL[root])) for root, code, d in NODE_FIXES]


MEDALLION = ('Each cell is one round medallion that fills its square, with an ornate gold rim, on a near-black '
             'background. Stylized hand-painted high-fantasy MMORPG icon art, bold readable shapes that stay clear '
             'when tiny, vivid colours, dramatic glow and rim lighting.')
NO_TEXT = ('IMPORTANT: no captions, no labels, no names, no text or letters of any kind anywhere in the image. '
           'No logos, no watermark, no outer frame.')

SETS = {
    'nodes': dict(dir='class_nodes', mode='medallion', size=160, entries=nodes,
                  head='fantasy class emblem icons. ' + MEDALLION + ' Each emblem shows one symbol:'),
    'regions': dict(dir='regions', mode='medallion', size=160, entries=lambda: REGIONS,
                    head='fantasy map location icons. ' + MEDALLION + ' Inside each medallion is a small painted '
                         'landscape vignette of the place:'),
    'professions': dict(dir='professions', mode='medallion', size=160, entries=lambda: PROFESSIONS,
                        head='fantasy crafting profession emblem icons. ' + MEDALLION + ' Each emblem shows the '
                             'tools of one craft, warm workshop light:'),
    'helpers': dict(dir='helpers', mode='portrait', size=192, entries=lambda: HELPERS,
                    head='fantasy character portraits of friendly hired workers. Each portrait is head and '
                         'shoulders, face centered, looking at the viewer, on a dark simple vignette background. '
                         'Stylized hand-painted high-fantasy MMORPG art, bold readable shapes, warm lighting, '
                         'every person clearly different in age, race, face and clothes:'),
    'materials': dict(dir='items', mode='item', size=128, entries=lambda: MATERIALS,
                      head='fantasy RPG crafting material icons. Each cell shows exactly one item, whole and '
                           'centered, at a slight diagonal, on the same plain dark slate-blue backdrop with a soft '
                           'vignette. Stylized hand-painted game icon art like polished mobile MMORPG inventory '
                           'icons: clean bold silhouettes, soft cel shading, crisp rim light. Every item looks '
                           'unique:'),
}


def sheets(name):
    e = SETS[name]['entries']()
    return [e[i:i + GRID * GRID] for i in range(0, len(e), GRID * GRID)]


def prompt(name, n):
    cells = sheets(name)[n]
    parts = ['%s: %s' % (POS[k], d) for k, (_, d) in enumerate(cells)]
    empty = '' if len(cells) == GRID * GRID else ' The remaining cells show only the empty background.'
    return ('Generate a square image: a %dx%d grid of %d separate square %s %s.%s Divided by thin solid black '
            'lines, cells of equal size. %s' % (GRID, GRID, GRID * GRID, SETS[name]['head'], '; '.join(parts),
                                                empty, NO_TEXT))


def trim(im, dark=14):
    """Drops near-black rows/columns at the edges (grid lines)."""
    gray = im.convert('L')
    w, h = gray.size
    px = gray.load()
    col = lambda x: sum(px[x, y] for y in range(0, h, 4)) / len(range(0, h, 4))
    row = lambda y: sum(px[x, y] for x in range(0, w, 4)) / len(range(0, w, 4))
    l, r, t, b = 0, w - 1, 0, h - 1
    while l < w // 6 and col(l) < dark:
        l += 1
    while r > w * 5 // 6 and col(r) < dark:
        r -= 1
    while t < h // 6 and row(t) < dark:
        t += 1
    while b > h * 5 // 6 and row(b) < dark:
        b -= 1
    return im.crop((l, t, r + 1, b + 1))


def square(im):
    side = min(im.size)
    x0, y0 = (im.width - side) // 2, (im.height - side) // 2
    return im.crop((x0, y0, x0 + side, y0 + side))


def cut(name, numbers):
    s = SETS[name]
    src = os.path.join(ROOT, 'client', 'art_originals', name)
    dst = os.path.join(ROOT, 'client', 'assets', s['dir'])
    os.makedirs(dst, exist_ok=True)
    for n, cells in enumerate(sheets(name)):
        path = os.path.join(src, 'sheet%02d.jpg' % n)
        if (numbers and n not in numbers) or not os.path.exists(path):
            continue
        im = Image.open(path).convert('RGB')
        gray = im.convert('L')
        w, h = im.size
        xs = [0] + [p for a, b in gutters(gray, 0, GRID) for p in (a, b + 1)] + [w]
        ys = [0] + [p for a, b in gutters(gray, 1, GRID) for p in (a, b + 1)] + [h]
        inset = max(2, w // 200)
        for k, (code, _) in enumerate(cells):
            if code is None:
                continue
            r, c = divmod(k, GRID)
            cell = im.crop((xs[2 * c] + inset, ys[2 * r] + inset, xs[2 * c + 1] - inset, ys[2 * r + 1] - inset))
            cell = medallion(cell) if s['mode'] == 'medallion' else square(trim(cell))
            cell.resize((s['size'], s['size']), Image.LANCZOS).save(
                os.path.join(dst, code + '.webp'), 'WEBP', quality=84, method=6)
        print(name, 'sheet%02d' % n, sum(1 for c, _ in cells if c))


if __name__ == '__main__':
    cmd = sys.argv[1] if len(sys.argv) > 1 else 'plan'
    if cmd == 'plan':
        for name in SETS:
            print(name, [len(x) for x in sheets(name)])
    elif cmd == 'prompt':
        print(prompt(sys.argv[2], int(sys.argv[3])))
    elif cmd == 'cut':
        cut(sys.argv[2], [int(x) for x in sys.argv[3:]])
