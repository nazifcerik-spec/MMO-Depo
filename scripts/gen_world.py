"""World content past level 25 (DD-254): zones, mobs and dungeons, with
names in en, tr, zh-CN and es. Mob stats follow the content's level curve
(hp 150+25L, attack 28+2.1L, defense 1+0.8L, xp 10·L^1.6) × archetype — the
same curve the level 1-25 mobs use, which the balance package confirms still
matches the class-average player with the gear ladder to level 100.

Run from the repository root: python scripts/gen_world.py
"""
import io
import os
import re

ROOT = os.path.join(os.path.dirname(__file__), '..', 'content')
L = lambda en, tr, zh, es: dict(en=en, tr=tr, zh=zh, es=es)

# archetype: hp, attack, defense multipliers; crit, evasion
ARCH = {
    'normal': (1.0, 1.0, 1.0, 5, 5),
    'tank': (1.75, 1.0, 1.9, 5, 2),
    'glass': (0.9, 1.5, 0.75, 11, 5),
    'evasive': (0.95, 1.15, 1.0, 9, 20),
    'crit': (1.05, 1.18, 1.0, 17, 10),
    'elite': (1.4, 1.15, 1.2, 9, 5),
    'chief': (1.2, 1.08, 1.1, 7, 5),  # dungeon bosses (the dungeon doubles HP, +20% attack)
}

ZONES = [
    dict(code='misty_marshes', level=25, travel=420, link='emberpeak',
         name=L('Misty Marshes', 'Sisli Bataklıklar', '迷雾沼泽', 'Marismas Brumosas'),
         mobs=[
             ('mist_leech', 30, 'evasive', L('Mist Leech', 'Sis Sülüğü', '雾水蛭', 'Sanguijuela de Niebla')),
             ('mud_golem', 35, 'tank', L('Mud Golem', 'Çamur Golemi', '泥魔像', 'Gólem de Barro')),
             ('bog_witch', 40, 'glass', L('Bog Witch', 'Bataklık Cadısı', '沼泽女巫', 'Bruja del Pantano')),
             ('marsh_queen', 42, 'chief', L('Marsh Queen', 'Bataklık Kraliçesi', '沼泽女王', 'Reina del Pantano')),
             ('venom_lizard', 45, 'crit', L('Venom Lizard', 'Zehirli Kertenkele', '毒蜥', 'Lagarto Venenoso')),
             ('drowned_knight', 50, 'elite', L('Drowned Knight', 'Batık Şövalye', '溺亡骑士', 'Caballero Ahogado')),
         ],
         dungeon=dict(code='sunken_temple', level=38, minutes=45, waves=['mist_leech', 'mud_golem', 'bog_witch'],
                      boss='marsh_queen', gold=8000, xp=12000,
                      name=L('Sunken Temple', 'Batık Tapınak', '沉没神殿', 'Templo Hundido'))),
    dict(code='iron_mountains', level=50, travel=480, link='misty_marshes',
         name=L('Iron Mountains', 'Demir Dağlar', '铁山', 'Montañas de Hierro'),
         mobs=[
             ('mountain_troll', 55, 'tank', L('Mountain Troll', 'Dağ Trolü', '山岭巨魔', 'Trol de Montaña')),
             ('dwarf_ghost', 60, 'evasive', L('Dwarf Ghost', 'Cüce Hayaleti', '矮人幽魂', 'Fantasma Enano')),
             ('rock_eagle', 65, 'glass', L('Rock Eagle', 'Kaya Kartalı', '岩鹰', 'Águila de Roca')),
             ('iron_king', 67, 'chief', L('Iron King', 'Demir Kral', '钢铁之王', 'Rey de Hierro')),
             ('iron_spider', 70, 'crit', L('Iron Spider', 'Demir Örümcek', '铁蛛', 'Araña de Hierro')),
             ('mine_warden', 75, 'elite', L('Mine Warden', 'Maden Bekçisi', '矿井看守', 'Guardián de la Mina')),
         ],
         dungeon=dict(code='forgotten_mine', level=63, minutes=50, waves=['mountain_troll', 'dwarf_ghost', 'rock_eagle'],
                      boss='iron_king', gold=15000, xp=20000,
                      name=L('Forgotten Mine', 'Unutulmuş Maden', '被遗忘的矿井', 'Mina Olvidada'))),
    dict(code='northern_keep_ruins', level=75, travel=540, link='iron_mountains',
         name=L('Northern Keep Ruins', 'Kuzey Kalesi Harabeleri', '北方要塞废墟', 'Ruinas de la Fortaleza del Norte'),
         mobs=[
             ('ghost_guard', 80, 'normal', L('Ghost Guard', 'Hayalet Muhafız', '幽灵卫兵', 'Guardia Fantasma')),
             ('cursed_archer', 85, 'glass', L('Cursed Archer', 'Lanetli Okçu', '诅咒弓手', 'Arquero Maldito')),
             ('bone_rider', 90, 'crit', L('Bone Rider', 'Kemik Atlı', '骸骨骑手', 'Jinete de Hueso')),
             ('fallen_lord_varn', 92, 'chief', L('Fallen Lord Varn', 'Düşmüş Lord Varn', '堕落领主瓦恩', 'Señor Caído Varn')),
             ('keep_gargoyle', 95, 'tank', L('Keep Gargoyle', 'Kale Gargoyli', '要塞石像鬼', 'Gárgola de la Fortaleza')),
             ('fallen_commander', 100, 'elite', L('Fallen Commander', 'Düşmüş Komutan', '堕落指挥官', 'Comandante Caído')),
         ],
         dungeon=dict(code='northern_keep', level=88, minutes=60, waves=['ghost_guard', 'cursed_archer', 'bone_rider'],
                      boss='fallen_lord_varn', gold=25000, xp=30000,
                      name=L('Northern Keep', 'Kuzey Kalesi', '北方要塞', 'Fortaleza del Norte'))),
]


def mob_yaml(code, level, arch):
    hm, am, dm, crit, eva = ARCH[arch]
    gold = round(503 * (level / 25) ** 1.5)
    return '\n'.join([
        '  - code: ' + code, '    display_name_key: mob.' + code, '    level: %d' % level,
        '    hp_max: %d' % round((150 + 25 * level) * hm), '    attack: %d' % round((28 + 2.1 * level) * am),
        '    defense: %d' % round((1 + 0.8 * level) * dm), '    crit_pct: %d' % crit,
        '    hit_pct: %d' % min(92, 82 + level // 5), '    evasion_pct: %d' % eva,
        '    xp_reward: %d' % round(10 * level ** 1.6), '    gold_copper_min: %d' % gold,
        '    gold_copper_max: %d' % round(gold * 2.5)])


mobs, dungeons, regions, links = [], [], [], []
names = {k: {} for k in ('en', 'tr', 'zh', 'es')}


def put(key, n):
    for lang in names:
        names[lang][key] = n[lang]


for z in ZONES:
    put('region.' + z['code'], z['name'])
    regions.append('\n'.join(['  - code: ' + z['code'], '    display_name_key: region.' + z['code'], '    kind: zone',
                              '    parent_code: sundermere', '    min_level: %d' % z['level'],
                              '    travel_seconds_from_parent: %d' % z['travel'], '    safe: false']))
    a, b = sorted([z['code'], z['link']])
    links.append('  - { region_a_code: %s, region_b_code: %s, travel_seconds: 300 }' % (a, b))
    for code, level, arch, n in z['mobs']:
        put('mob.' + code, n)
        mobs.append(mob_yaml(code, level, arch))
    d = z['dungeon']
    put('dungeon.' + d['code'], d['name'])
    put('region.' + d['code'], d['name'])
    regions.append('\n'.join(['  - code: ' + d['code'], '    display_name_key: region.' + d['code'], '    kind: dungeon',
                              '    parent_code: ' + z['code'], '    min_level: %d' % d['level'],
                              '    travel_seconds_from_parent: 300', '    safe: false']))
    dungeons.append('\n'.join(['  - code: ' + d['code'], '    display_name_key: dungeon.' + d['code'],
                               '    min_level: %d' % d['level'],
                               '    duration_seconds: %d         # %d minutes' % (d['minutes'] * 60, d['minutes']),
                               '    waves:',
                               '      - { mob_code: %s, count: 3 }' % d['waves'][0],
                               '      - { mob_code: %s, count: 3 }' % d['waves'][1],
                               '      - { mob_code: %s, count: 2 }' % d['waves'][2],
                               '    boss_mob_code: ' + d['boss'], '    bonus_gold_copper: %d' % d['gold'],
                               '    bonus_xp: %d' % d['xp']]))


def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        s = f.read()
    return s.replace('\r\n', '\n'), '\r\n' in s


def write(path, s, crlf):
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s.replace('\n', '\r\n') if crlf else s)


def block(s, name, body, before=None):
    start, end = '  # >>> generated: %s (scripts/gen_world.py)\n' % name, '  # <<< generated: %s\n' % name
    b = start + body.rstrip('\n') + '\n' + end
    if start in s:
        return s[:s.index(start)] + b + s[s.index(end) + len(end):]
    if before:
        assert s.count(before) == 1, before
        return s.replace(before, b + '\n' + before)
    return s.rstrip('\n') + '\n\n' + b


for fname, name, body, before in [
    ('mobs.yaml', 'mobs 25-100', '\n\n'.join(mobs), None),
    ('dungeons.yaml', 'dungeons 25-100', '\n\n'.join(dungeons), None),
]:
    p = os.path.join(ROOT, fname)
    s, crlf = read(p)
    write(p, block(s, name, body, before), crlf)
p = os.path.join(ROOT, 'regions.yaml')
s, crlf = read(p)
s = block(s, 'regions 25-100', '\n\n'.join(regions), before='# Non-parent shortcuts.')
s = block(s, 'links 25-100', '\n'.join(links))
write(p, s, crlf)

FILES = {'en': 'en', 'tr': 'tr', 'zh': 'zh-CN', 'es': 'es'}
for lang, fname in FILES.items():
    p = os.path.join(ROOT, 'i18n', fname + '.yaml')
    s, crlf = read(p)
    lines = s.split('\n')
    for key, value in sorted(names[lang].items()):
        line = '%s: "%s"' % (key, value)
        idx = {re.match(r'^([\w.\-]+):', l).group(1): i for i, l in enumerate(lines) if re.match(r'^[\w.\-]+:\s*"', l)}
        if key in idx:
            lines[idx[key]] = line
            continue
        kind = key.split('.', 1)[0]
        same = sorted(k for k in idx if k.split('.', 1)[0] == kind)
        if same:
            after = [k for k in same if k > key]
            lines.insert(idx[after[0]] if after else idx[same[-1]] + 1, line)
        else:  # partial languages: new section at the end
            lines += ['', '# --- %s ---' % kind, line] if lines[-1] != '' else ['# --- %s ---' % kind, line, '']
    write(p, '\n'.join(lines).rstrip('\n') + '\n', crlf)
print('mobs', len(mobs), 'dungeons', len(dungeons), 'regions', len(regions))
