"""Quest chains through the world (DD-278).

The main chain (hand-written to level 10) continues zone by zone: for each
zone, hunt its first monster and clear its dungeon; the next link needs the
previous one. Each zone also has a side bounty on its elite. Rewards follow
the XP curve: a main link pays half a level at its level, a bounty 30%; gold
scales like mob gold. Writes a generated block into content/quests.yaml and
names in en, tr, zh-CN and es.

Run from the repository root after gen_world.py: python scripts/gen_quests.py
"""
import io
import math
import os
import sys

import yaml

sys.path.insert(0, os.path.dirname(__file__))
from i18n_util import merge  # noqa: E402

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.join(HERE, '..', 'content')
START, END = '  # >>> generated: quests (scripts/gen_quests.py)\n', '  # <<< generated: quests\n'
LANGS = {'en': 'en', 'tr': 'tr', 'zh': 'zh-CN', 'es': 'es'}


def world():
    """ZONES from gen_world.py, without running its writer."""
    with io.open(os.path.join(HERE, 'gen_world.py'), encoding='utf-8') as f:
        src = f.read()
    ns = {'__file__': os.path.join(HERE, 'gen_world.py')}
    exec(compile(src[:src.index('def mob_yaml')], 'gen_world', 'exec'), ns)
    return ns['ZONES']


def xp_to_next(level):  # leveling: ceil((a + b·L^c) × baseline)
    with io.open(os.path.join(ROOT, 'progression.yaml'), encoding='utf-8') as f:
        c = yaml.safe_load(f)['xp_curve']
    return math.ceil((c['a'] + c['b'] * level ** c['c']) * c['baseline_xp_per_minute'])


def gold(level):  # ~25 kills' worth of mob gold (gen_world: 503·(L/25)^1.5 a mob)
    return round(25 * 503 * (level / 25) ** 1.5)


i18n = {}
for lang, f in LANGS.items():
    with io.open(os.path.join(ROOT, 'i18n', f + '.yaml'), encoding='utf-8') as fh:
        i18n[lang] = yaml.safe_load(fh)


def name_of(key):
    """A content name in every language, English where one is missing."""
    return {lang: i18n[lang].get(key) or i18n['en'][key] for lang in LANGS}


# (hunt mob, dungeon, level) links before the generated zones: levels 10-25
# live in hand-written regions.
LINKS = [('cultist_acolyte', 'shadowfen_ruins', 17), ('frost_giant', 'obsidian_spire', 23)]
for z in sorted(world(), key=lambda z: z['dungeon']['level']):
    LINKS.append((z['mobs'][0][0], z['dungeon']['code'], z['dungeon']['level']))
BOUNTIES = [(z['mobs'][5][0], z['mobs'][5][1]) for z in world()]

out, names, prev, order = [], {lang: {} for lang in LANGS}, 'main_06_hollow', 70
for i, (mob, dungeon, level) in enumerate(LINKS, start=7):
    code = 'main_%02d_%s' % (i, dungeon)
    d = name_of('dungeon.' + dungeon)
    for lang, title in {'en': 'The Road to %s', 'tr': '%s Yolu', 'zh': '%s之路', 'es': 'Camino a %s'}.items():
        names[lang]['quest.' + code] = title % d[lang]
    out.append('\n'.join([
        '  - code: ' + code, '    display_name_key: quest.' + code, '    category: main',
        '    min_level: %d' % level, '    requires_quest: ' + prev, '    objectives:',
        '      - { event: "mob_killed:%s", count: 20 }' % mob,
        '      - { event: "dungeon_cleared:%s", count: 1 }' % dungeon,
        '    reward_xp: %d' % round(0.5 * xp_to_next(level)), '    reward_gold_copper: %d' % gold(level),
        '    sort_order: %d' % order]))
    prev, order = code, order + 10
for mob, level in sorted(BOUNTIES, key=lambda b: b[1]):
    code = 'side_bounty_' + mob
    m = name_of('mob.' + mob)
    for lang, title in {'en': 'Bounty: %s', 'tr': 'Ödül Avı: %s', 'zh': '悬赏：%s', 'es': 'Recompensa: %s'}.items():
        names[lang]['quest.' + code] = title % m[lang]
    out.append('\n'.join([
        '  - code: ' + code, '    display_name_key: quest.' + code, '    category: side',
        '    min_level: %d' % max(1, level - 5), '    objectives:',
        '      - { event: "mob_killed:%s", count: 10 }' % mob,
        '    reward_xp: %d' % round(0.3 * xp_to_next(level)), '    reward_gold_copper: %d' % (gold(level) // 2),
        '    sort_order: %d' % order]))
    order += 10

p = os.path.join(ROOT, 'quests.yaml')
with io.open(p, encoding='utf-8', newline='') as f:
    raw = f.read()
crlf = '\r\n' in raw
s = raw.replace('\r\n', '\n')
block = START + '\n\n'.join(out) + '\n' + END
if START in s:
    s = s[:s.index(START)] + block + s[s.index(END) + len(END):]
else:
    s = s.rstrip('\n') + '\n\n' + block
with io.open(p, 'w', encoding='utf-8', newline='') as f:
    f.write(s.replace('\n', '\r\n') if crlf else s)
merge(names)
print('quests', len(out))
