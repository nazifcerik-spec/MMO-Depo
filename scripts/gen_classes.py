"""Generate the class progression content of the nine non-Warrior base
classes (DD-252) from scripts/class_data.py:

  content/progression.yaml   class nodes + talent trees (between markers)
  content/skills.yaml        class skills (between markers)
  content/gear.yaml          weapon lines of the new classes (between markers)
  content/vendors.yaml       tier-1 weapons on sale
  content/i18n/*.yaml        names + descriptions in en, tr, zh-CN, es

Run from the repository root: python scripts/gen_classes.py
Warrior (DD-250) stays hand-written; the generator only adds its mastery
achievement names. Idempotent: re-running rewrites the generated blocks.
"""
import io
import os
import re
import sys

sys.path.insert(0, os.path.dirname(__file__))
from class_data import CLASSES, ROLE_NAMES, RESOURCE_NAMES, TALENTS, WEAPON_STATS  # noqa: E402

ROOT = os.path.join(os.path.dirname(__file__), '..')
CONTENT = os.path.join(ROOT, 'content')
LANGS = {'en': 'en', 'tr': 'tr', 'zh': 'zh-CN', 'es': 'es'}

MOD = {
    'atk': ('mult_add', 'attack', 0.01), 'hp': ('mult_add', 'hp_max', 0.01), 'def': ('mult_add', 'defense', 0.01),
    'mp': ('mult_add', 'mp_max', 0.01), 'crit': ('flat', 'crit_pct', 1), 'cdmg': ('flat', 'crit_dmg_pct', 1),
    'hit': ('flat', 'hit_pct', 1), 'eva': ('flat', 'evasion_pct', 1), 'blk': ('flat', 'block_pct', 1),
    'ls': ('flat', 'lifesteal_pct', 1), 'dr': ('flat', 'dmg_reduction_pct', 1), 'exe': ('flat', 'execute_pct', 1),
    'miss': ('flat', 'missing_hp_dmg_pct', 1),
}
EFFECT = {'S': 'strike', 'X': 'stun', 'G': 'guard', 'E': 'empower', 'H': 'heal', 'O': 'hot', 'B': 'barrier',
          'D': 'dot', 'W': 'weaken'}
DEFAULT_TYPE = {'strike': 'DAMAGE', 'stun': 'STUN', 'guard': 'BUFF', 'empower': 'BUFF', 'heal': 'HEAL',
                'hot': 'HOT', 'barrier': 'BARRIER', 'dot': 'DOT', 'weaken': 'DEBUFF'}

# --- descriptions -----------------------------------------------------------
MOD_TEXT = {
    'attack': ('+{v}% attack', '+%{v} saldırı', '攻击 +{v}%', '+{v}% de ataque'),
    'hp_max': ('+{v}% max HP', '+%{v} en yüksek can', '最大生命值 +{v}%', '+{v}% de vida máxima'),
    'defense': ('+{v}% armor', '+%{v} zırh', '护甲 +{v}%', '+{v}% de armadura'),
    'mp_max': ('+{v}% mana', '+%{v} mana', '法力 +{v}%', '+{v}% de maná'),
    'crit_pct': ('+{v}% crit chance', '+%{v} kritik şansı', '暴击率 +{v}%', '+{v}% de probabilidad crítica'),
    'crit_dmg_pct': ('+{v}% crit damage', '+%{v} kritik hasarı', '暴击伤害 +{v}%', '+{v}% de daño crítico'),
    'hit_pct': ('+{v}% accuracy', '+%{v} isabet', '命中 +{v}%', '+{v}% de precisión'),
    'evasion_pct': ('+{v}% dodge', '+%{v} kaçınma', '闪避 +{v}%', '+{v}% de esquiva'),
    'block_pct': ('+{v}% block chance', '+%{v} blok şansı', '格挡几率 +{v}%', '+{v}% de probabilidad de bloqueo'),
    'lifesteal_pct': ('+{v}% lifesteal', '+%{v} can çalma', '吸血 +{v}%', '+{v}% de robo de vida'),
    'dmg_reduction_pct': ('+{v}% damage reduction', '+%{v} hasar azaltma', '伤害减免 +{v}%', '+{v}% de reducción de daño'),
    'execute_pct': ('+{v}% damage to targets below 35% HP', '%35 canın altındaki hedeflere +%{v} hasar',
                    '对生命低于35%的目标伤害 +{v}%', '+{v}% de daño a objetivos con menos del 35% de vida'),
    'missing_hp_dmg_pct': ('up to +{v}% damage as your HP falls', 'canın azaldıkça +%{v} hasara kadar',
                           '随生命降低伤害最多 +{v}%', 'hasta +{v}% de daño según pierdes vida'),
}
JOIN = (', ', ', ', '，', ', ')
PER_RANK = (lambda s: s + ' per rank', lambda s: 'Kademe başı ' + s, lambda s: '每级' + s, lambda s: s + ' por rango')
CAPSTONE = ('Capstone: ', 'Zirve: ', '终极天赋：', 'Cúspide: ')
ACTIVE_TEXT = {
    'strike': ('Deals {p}% weapon damage.', '%{p} silah hasarı verir.', '造成{p}%武器伤害。', 'Inflige un {p}% de daño de arma.'),
    'stun': ('Deals {p}% weapon damage and stuns for {d} turn(s).', '%{p} silah hasarı verir ve {d} tur sersemletir.',
             '造成{p}%武器伤害并眩晕{d}回合。', 'Inflige un {p}% de daño de arma y aturde {d} turno(s).'),
    'guard': ('+{p}% armor for {d} turns.', '{d} tur +%{p} zırh.', '{d}回合护甲 +{p}%。', '+{p}% de armadura durante {d} turnos.'),
    'empower': ('+{p}% attack for {d} turns.', '{d} tur +%{p} saldırı.', '{d}回合攻击 +{p}%。', '+{p}% de ataque durante {d} turnos.'),
    'heal': ('Restores {p}% max HP.', 'En yüksek canın %{p} kadarını yeniler.', '恢复{p}%最大生命值。', 'Restaura un {p}% de la vida máxima.'),
    'hot': ('Restores {p}% max HP each turn for {d} turns.', '{d} tur boyunca her tur en yüksek canın %{p} kadarını yeniler.',
            '{d}回合内每回合恢复{p}%最大生命值。', 'Restaura un {p}% de la vida máxima por turno durante {d} turnos.'),
    'barrier': ('Barrier absorbing {p}% max HP for {d} turns.', '{d} tur boyunca en yüksek canın %{p} kadarını emen bariyer.',
                '获得吸收{p}%最大生命值的屏障，持续{d}回合。', 'Barrera que absorbe un {p}% de la vida máxima durante {d} turnos.'),
    'dot': ('Deals {p}% damage each turn for {d} turns.', '{d} tur boyunca her tur %{p} hasar verir.',
            '{d}回合内每回合造成{p}%伤害。', 'Inflige un {p}% de daño por turno durante {d} turnos.'),
    'weaken': ('Lowers enemy armor by {p}% for {d} turns.', 'Düşman zırhını {d} tur %{p} düşürür.',
               '{d}回合内敌人护甲降低{p}%。', 'Reduce la armadura enemiga un {p}% durante {d} turnos.'),
}
GAIN_TEXT = (' Generates {g} {r}.', ' {g} {r} üretir.', '产生{g}点{r}。', ' Genera {g} de {r}.')
TARGET_TEXT = (' Used on targets below {t}% HP.', ' Canı %{t} altındaki hedeflere kullanılır.',
               '对生命低于{t}%的目标使用。', ' Se usa contra objetivos con menos del {t}% de vida.')
SELF_TEXT = (' Auto: at or below {h}% HP.', ' Otomatik: can %{h} veya altındayken.', '自动：生命不高于{h}%时。',
             ' Automático: con {h}% de vida o menos.')
ALLY_TEXT = (' Auto: an ally at or below {h}% HP.', ' Otomatik: bir müttefikin canı %{h} veya altındayken.',
             '自动：友方生命不高于{h}%时。', ' Automático: un aliado con {h}% de vida o menos.')
MASTERY_ACH = ('{n} Mastery', '{n} Ustalığı', '{n}精通', 'Maestría de {n}')
LANG_ORDER = ('en', 'tr', 'zh', 'es')
WARRIOR_SPECS = ('iron_fortress', 'warlord', 'berserker', 'weapon_master')


def fmt_num(v):
    return str(int(v)) if float(v).is_integer() else ('%g' % v)


def parse_mods(spec):
    out = []
    for part in spec.split(','):
        m = re.fullmatch(r'([a-z]+)(\d+)', part.strip())
        assert m and m.group(1) in MOD, spec
        kind, target, scale = MOD[m.group(1)]
        amount = int(m.group(2)) * scale
        out.append((kind, target, round(amount, 4)))
    return out


def mods_text(mods, i, per_rank=False):
    parts = []
    for kind, target, amount in mods:
        v = amount * 100 if kind == 'mult_add' else amount
        parts.append(MOD_TEXT[target][i].format(v=fmt_num(round(v, 2))))
    s = JOIN[i].join(parts)
    return PER_RANK[i](s) if per_rank else s


def parse_active(note):
    toks = note.split()
    m = re.fullmatch(r'([A-Z])(\d+)/(\d+)', toks[0])
    assert m, note
    a = dict(effect=EFFECT[m.group(1)], power=int(m.group(2)), cd=int(m.group(3)), d=0, h=0, a=0, t=0,
             cost=None, gain=0, type=None)
    for t in toks[1:]:
        if t.startswith('T:'):
            a['type'] = t[2:]
        elif t[0] in 'dhatcg' and t[1:].isdigit():
            a[{'d': 'd', 'h': 'h', 'a': 'a', 't': 't', 'c': 'cost', 'g': 'gain'}[t[0]]] = int(t[1:])
        else:
            raise ValueError(note)
    return a


def cost_of(a, ultimate):
    if a['cost'] is not None:
        c = a['cost']
    elif a['gain'] > 0:
        c = 0
    else:
        e, p = a['effect'], a['power']
        c = {
            'strike': lambda: min(60, max(5, (p - 100) // 5 + 5)),
            'stun': lambda: min(60, max(10, (p - 100) // 5 + 10)),
            'dot': lambda: min(50, 10 + p // 4),
            'empower': lambda: 10 + p // 2, 'guard': lambda: 10 + p // 2, 'weaken': lambda: 10 + p // 2,
            'heal': lambda: 10 + p // 2,
            'hot': lambda: 10 + p * max(1, a['d']) // 3,
            'barrier': lambda: 10 + p // 3,
        }[e]()
    if ultimate:
        c = max(c, 50)
    return min(c, 90)


def sort_of(a, category):
    if a['h'] and a['h'] <= 30 or a['a'] and a['a'] <= 30:
        return 1
    if category == 'ultimate':
        return 3
    if a['h'] or a['a']:
        return 6
    if a['t']:
        return 10
    if category == 'awakening':
        return 15
    e = a['effect']
    if e in ('dot', 'stun'):
        return 30
    if e in ('empower', 'weaken'):
        return 86  # set-up casts fill turns when no strike is ready
    if e in ('guard', 'barrier', 'hot', 'heal'):
        return 87
    if a['gain'] > 0 or a['power'] <= 130:
        return 90
    return max(40, 85 - a['power'] // 10)


def y(v):
    return '"' + str(v).replace('\\', '\\\\').replace('"', '\\"') + '"'


# --- build ------------------------------------------------------------------
strings = {lang: {} for lang in LANGS}


def put(key, names):
    for lang in LANGS:
        strings[lang][key] = names[lang]


def nodes(c):
    yield c, 0, None
    for p in c['promotions']:
        yield p, 1, c
        for s in p['specs']:
            yield s, 2, p


skills_out, classes_out, trees_out = [], [], []
for c in CLASSES:
    res_names = RESOURCE_NAMES[c['resource']['kind']]
    for node, tier, parent in nodes(c):
        code = node['code']
        put('class.' + code, node['name'])
        put('class.' + code + '.role', {lang: ' / '.join(ROLE_NAMES[r][lang] for r in node['roles'][:3]) for lang in LANGS})
        # class node yaml
        lines = ['  - code: ' + code]
        if parent:
            lines += ['    parent: ' + parent['code'], '    required_level: %d' % (100 if tier == 1 else 300)]
        lines += ['    name_key: class.' + code, '    role_tags: [%s]' % ', '.join(node['roles'])]
        if tier == 0:
            r = c['resource']
            lines.append('    resource: { kind: %s, max: %d, start_pct: %d, per_hit: %d, per_turn: %d }' % (
                r['kind'], r['max'], r['start_pct'], r['per_hit'], r['per_turn']))
            b = c['base']
            lines.append('    base_stats: { strength: %d, vitality: %d, dexterity: %d, intelligence: %d, wisdom: %d, luck: %d }' % (
                b['strength'], b['vitality'], b['dexterity'], b['intelligence'], b['wisdom'], b['luck']))
            lines.append('    derived:')
            for f in ('hp', 'attack', 'armor', 'crit', 'hit', 'dodge', 'block'):
                base, per_level, per = c['derived'][f]
                parts = []
                if base:
                    parts.append('base: %s' % fmt_num(base))
                if per_level:
                    parts.append('per_level: %s' % fmt_num(per_level))
                if per:
                    parts.append('per: { %s }' % ', '.join('%s: %s' % (k, fmt_num(v)) for k, v in per.items()))
                lines.append('      %-7s { %s }' % (f + ':', ', '.join(parts)))
            lines.append('    primary_attribute: ' + c['primary'])
            lines.append('    secondary_attributes: [%s]' % ', '.join(c['secondary']))
            lines.append('    equipment_permissions: [%s]' % ', '.join(c['perms']))
        support = node.get('support', 0)
        if support:
            lines.append('    support_power: %s' % fmt_num(support))
        if tier == 2:
            ks, sv, re_, gv, bv, av, af = node['afk']
            lines.append('    afk_profile: { kill_speed: %d, survivability: %d, resource_efficiency: %d, group_value: %d, '
                         'boss_value: %d, aoe_value: %d, afk_safety: %d }' % (ks, sv, re_, gv, bv, av, af))
            lines.append('    strong_in: [%s]' % ', '.join(node['strong']))
            lines.append('    weak_in: [%s]' % ', '.join(node['weak']))
            title, tnames, tmods = node['mastery']
            put('title.' + title, tnames)
            lines += ['    mastery:', '      title: ' + title, '      badge: badge_' + code, '      frame: frame_' + code,
                      '      achievement: mastery_' + code, '      prestige_eligible: true', '      modifiers:']
            for kind, target, amount in parse_mods(tmods):
                lines.append('        - { kind: %s, target: %s, amount: %s }' % (kind, target, fmt_num(amount)))
        classes_out.append('\n'.join(lines))

        # skills
        damage = node.get('damage') or (parent and parent.get('damage')) or c['damage']
        for (lvl, key, note, *names) in node['skills']:
            names = dict(zip(LANG_ORDER, names))
            scode = code + '_' + key
            category = {600: 'awakening', 700: 'ultimate'}.get(lvl)
            put('skill.' + scode, names)
            sl = ['  - code: ' + scode, '    display_name_key: skill.' + scode, '    class: ' + c['code']]
            if tier > 0:
                sl.append('    spec: ' + code)
            sl.append('    unlock_level: %d' % lvl)
            if note.startswith('P '):
                rest = note[2:].split(' @')
                mods = parse_mods(rest[0])
                trig, stacks = 'ON_COMBAT_START', 0
                if len(rest) > 1:
                    trig, _, st = rest[1].partition('/')
                    stacks = int(st) if st else 0
                sl.append('    kind: passive')
                if category:
                    sl.append('    category: ' + category)
                sl.append('    type: STAT_MODIFIER')
                sl.append('    passive: { trigger: %s%s, scaling: steady_state }' % (trig, ', stack_limit: %d' % stacks if stacks else ''))
                sl.append('    tags: [passive%s]' % (', ' + category if category else ''))
                sl.append('    modifiers:')
                for kind, target, amount in mods:
                    sl.append('      - { kind: %s, target: %s, amount: %s }' % (kind, target, fmt_num(amount)))
                for i, lang in enumerate(LANG_ORDER):
                    strings[lang]['skill.' + scode + '.desc'] = mods_text(mods, i) + ('.' if lang != 'zh' else '。')
            else:
                a = parse_active(note)
                e = a['effect']
                sl += ['    kind: active']
                if category:
                    sl.append('    category: ' + category)
                stance = e in ('guard',) and a['d'] >= 4
                sl += ['    activation_mode: ' + ('stance' if stance else 'active'),
                       '    sort_order: %d' % sort_of(a, category), '    effect: ' + e, '    power: %d' % a['power'],
                       '    mp_cost: %d' % cost_of(a, category == 'ultimate'), '    cooldown_turns: %d' % a['cd']]
                if a['gain']:
                    sl.append('    resource_gain: %d' % a['gain'])
                if a['d']:
                    sl.append('    duration_turns: %d' % a['d'])
                trig = []
                if a['h']:
                    trig.append('self_hp_below: %d' % a['h'])
                if a['a']:
                    trig.append('ally_hp_below: %d' % a['a'])
                if a['t']:
                    trig.append('target_hp_below: %d' % a['t'])
                if not trig:
                    trig.append('on_cooldown: true')
                sl.append('    auto_trigger: { %s }' % ', '.join(trig))
                typ = a['type'] or DEFAULT_TYPE[e]
                sl.append('    type: ' + typ)
                if e in ('strike', 'stun', 'dot'):
                    sl.append('    damage_type: ' + damage)
                target = 'enemy' if e in ('strike', 'stun', 'dot', 'weaken') else ('ally' if a['a'] else 'self')
                if e == 'empower' and ('buff_support' in node['roles'] or 'group_support' in node['roles']):
                    target = 'party'
                sl += ['    target_type: ' + target, '    range: ' + ('melee' if c['range'] == 'melee' and target == 'enemy' else 'ranged')]
                sl.append('    tags: [%s]' % ', '.join([typ.lower()] + ([category] if category else [])))
                for i, lang in enumerate(LANG_ORDER):
                    t = ACTIVE_TEXT[e][i].format(p=a['power'], d=a['d'])
                    if a['gain']:
                        t += GAIN_TEXT[i].format(g=a['gain'], r=res_names[lang])
                    if a['t']:
                        t += TARGET_TEXT[i].format(t=a['t'])
                    if a['h']:
                        t += SELF_TEXT[i].format(h=a['h'])
                    if a['a']:
                        t += ALLY_TEXT[i].format(h=a['a'])
                    strings[lang]['skill.' + scode + '.desc'] = t
            skills_out.append('\n'.join(sl))

    # talent trees
    for (tkey, tnames, archetypes, (ckey, cnames, cmods)) in c['trees']:
        tcode = c['code'] + '_' + tkey
        put('talent_tree.' + tcode, tnames)
        tl = ['  - code: ' + tcode, '    class: ' + c['code'], '    talents:']
        assert len(archetypes) == 5 and len(set(archetypes)) == 5, tcode
        for arch in archetypes:
            mods_spec, anames = TALENTS[arch]
            mods = parse_mods(mods_spec)
            acode = tcode + '_' + arch
            put('talent.' + acode, anames)
            for i, lang in enumerate(LANG_ORDER):
                strings[lang]['talent.' + acode + '.desc'] = mods_text(mods, i, per_rank=True)
            tl.append('      - { code: %s, max_rank: 5, modifiers: [%s] }' % (acode, ', '.join(
                '{ kind: %s, target: %s, amount: %s }' % (k, t, fmt_num(v)) for k, t, v in mods)))
        ccode = tcode + '_' + ckey
        put('talent.' + ccode, cnames)
        mods = parse_mods(cmods)
        for i, lang in enumerate(LANG_ORDER):
            strings[lang]['talent.' + ccode + '.desc'] = CAPSTONE[i] + mods_text(mods, i)
        tl += ['      - code: ' + ccode, '        max_rank: 1', '        capstone: true', '        modifiers:']
        for k, t, v in mods:
            tl.append('          - { kind: %s, target: %s, amount: %s }' % (k, t, fmt_num(v)))
        trees_out.append('\n'.join(tl))

for kind, names in RESOURCE_NAMES.items():
    put('resource.' + kind, names)
put('category.mastery', dict(en='Mastery', tr='Ustalık', zh='精通', es='Maestría'))  # mastery achievements

# Mastery achievement names for all 40 final specializations.
for c in CLASSES:
    for p in c['promotions']:
        for s in p['specs']:
            for i, lang in enumerate(LANG_ORDER):
                strings[lang]['achievement.mastery_' + s['code']] = MASTERY_ACH[i].format(n=s['name'][lang])


# --- write content files ----------------------------------------------------
def read(path):
    with io.open(path, encoding='utf-8', newline='') as f:
        s = f.read()
    return s.replace('\r\n', '\n'), '\r\n' in s


def write(path, s, crlf):
    if crlf:
        s = s.replace('\n', '\r\n')
    with io.open(path, 'w', encoding='utf-8', newline='') as f:
        f.write(s)


def replace_block(s, name, body, anchor_before=None, append=False):
    start, end = '  # >>> generated: %s (scripts/gen_classes.py)\n' % name, '  # <<< generated: %s\n' % name
    block = start + body.rstrip('\n') + '\n' + end
    if start in s:
        a = s.index(start)
        b = s.index(end) + len(end)
        return s[:a] + block + s[b:]
    if append:
        return s.rstrip('\n') + '\n\n' + block
    assert s.count(anchor_before) == 1, anchor_before
    return s.replace(anchor_before, block + '\n' + anchor_before)


p = os.path.join(CONTENT, 'progression.yaml')
s, crlf = read(p)
s = replace_block(s, 'classes', '\n\n'.join(classes_out), anchor_before='# Talent points start at 100.')
s = replace_block(s, 'talent trees', '\n\n'.join(trees_out), append=True)
write(p, s, crlf)

p = os.path.join(CONTENT, 'skills.yaml')
s, crlf = read(p)
old_marker = '  # --- Mage ------'
removed_skills = []
if old_marker in s:  # one-time: drop the v1 skills of the old classes
    i = s.index(old_marker)
    removed_skills = re.findall(r'^  - code: (\S+)', s[i:], re.M)
    s = s[:i].rstrip('\n') + '\n'
s = replace_block(s, 'skills', '\n\n'.join(skills_out), append=True)
write(p, s, crlf)

# gear: class renames, wisdom maces, new weapon lines
p = os.path.join(CONTENT, 'gear.yaml')
s, crlf = read(p)
s = s.replace('class_restriction: hunter', 'class_restriction: ranger').replace('class_restriction: priest', 'class_restriction: cleric')
s = s.replace('# --- hunter: bow ---', '# --- ranger: bow ---').replace('# --- priest: mace ---', '# --- cleric: mace ---')
s = re.sub(r'(item_subclass: mace\n(?:.*\n){8}    stat_bonuses:\n      - \{ kind: flat, target: attack, amount: \d+ \}\n      - \{ kind: flat, target: )int',
           r'\1cha', s)
gear = []
new_vendor = []
for c in CLASSES:
    if not c['weapon']:
        continue
    subclass, weight, prim, sec, items = c['weapon']
    gear.append('  # --- %s: %s ---' % (c['code'], subclass))
    for (lvl, value, atk, pa, sa), (code, en, tr) in zip(WEAPON_STATS, items):
        strings['en']['item.' + code] = en
        strings['tr']['item.' + code] = tr
        gear.append('\n'.join([
            '  - code: ' + code, '    display_name_key: item.' + code, '    item_class: weapon',
            '    item_subclass: ' + subclass, '    min_level: %d' % lvl, '    base_value_copper: %d' % value,
            '    base_weight_g: %d' % weight, '    bind_kind: boe', '    affixable: true', '    equip_slot: mainhand',
            '    class_restriction: ' + c['code'], '    stat_bonuses:',
            '      - { kind: flat, target: attack, amount: %d }' % atk,
            '      - { kind: flat, target: %s, amount: %d }' % (prim, pa),
            '      - { kind: flat, target: %s, amount: %d }' % (sec, sa)]))
    new_vendor.append(items[0][0])
s = replace_block(s, 'weapons', '\n\n'.join(gear), append=True)
write(p, s, crlf)

p = os.path.join(CONTENT, 'vendors.yaml')
s, crlf = read(p)
for code in new_vendor:
    line = '      - { template_code: %s }\n' % code
    if line not in s:
        s = s.replace('      - { template_code: dagger_rusty }\n', '      - { template_code: dagger_rusty }\n' + line, 1)
write(p, s, crlf)


# --- i18n -------------------------------------------------------------------
def load_strings(path):
    s, crlf = read(path)
    header, kv = [], {}
    for line in s.split('\n'):
        m = re.match(r'^([A-Za-z0-9_.\-]+):\s*"(.*)"\s*$', line)
        if m:
            kv[m.group(1)] = m.group(2).replace('\\"', '"').replace('\\\\', '\\')
        elif not kv and line.startswith('#') and not line.startswith('# ---'):
            header.append(line)
    return header, kv, crlf


for lang, fname in LANGS.items():
    path = os.path.join(CONTENT, 'i18n', fname + '.yaml')
    header, kv, crlf = load_strings(path)
    for code in removed_skills:
        kv.pop('skill.' + code, None)
        kv.pop('skill.' + code + '.desc', None)
    kv.update(strings[lang])
    for spec in WARRIOR_SPECS:  # hand-written Warrior (DD-250): mastery achievement names
        kv['achievement.mastery_' + spec] = MASTERY_ACH[LANG_ORDER.index(lang)].format(n=kv['class.' + spec])
    out = header[:] if header else []
    kinds = {}
    for k in sorted(kv):
        kinds.setdefault(k.split('.', 1)[0], []).append(k)
    for kind in sorted(kinds):
        out += ['', '# --- %s ---' % kind]
        out += ['%s: %s' % (k, y(kv[k])) for k in kinds[kind]]
    write(path, '\n'.join(out).lstrip('\n') + '\n', crlf)

total = sum(len(v) for v in strings.values())
print('classes %d, skills %d, trees %d, strings %d, removed %d' % (len(classes_out), len(skills_out), len(trees_out), total, len(removed_skills)))
