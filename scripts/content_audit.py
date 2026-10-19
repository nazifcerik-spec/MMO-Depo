"""Content audit: every item, recipe, drop rule, gear tier, set and monster
in one Turkish report with automatic consistency checks. Read-only — it
changes no content. The drop rules mirror server/internal/item (GearPool,
CombatHides, TierRarity) and activity.DropLevel; update both together.

Run from the repository root: python scripts/content_audit.py
Writes docs/content-audit.md.
"""
import collections
import io
import os
import re

import yaml

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.join(HERE, '..')
CONTENT = os.path.join(ROOT, 'content')
OUT = os.path.join(ROOT, 'docs', 'content-audit.md')

SLOTS = ['mainhand', 'offhand', 'helm', 'chest', 'legs', 'boots', 'gloves',
         'necklace', 'ring1', 'ring2', 'trinket1', 'trinket2']  # inventory.KnownEquipSlots
RARITY_TR = ['sıradan', 'sıradan', 'nadir-altı', 'nadir-altı', 'nadir', 'destansı', 'efsanevi', 'mitik']
GEAR_ODDS = 'sıradan %60 · seçkin %25 · nadir %11 · destansı %4'
BOSS_ODDS = 'seçkin %50 · nadir %35 · destansı %13 · efsanevi %2'
HIDE_SPAN = 4  # item.combatHideSpan


def load(name, key):
    with io.open(os.path.join(CONTENT, name), encoding='utf-8') as f:
        return yaml.safe_load(f)[key]


def world_zones():
    """ZONES/RAIDS data from gen_world.py, without running its writer."""
    with io.open(os.path.join(HERE, 'gen_world.py'), encoding='utf-8') as f:
        src = f.read()
    ns = {'__file__': os.path.join(HERE, 'gen_world.py')}
    exec(compile(src[:src.index('def mob_yaml')], 'gen_world', 'exec'), ns)
    return ns['ZONES'], ns['RAIDS'], ns['ARCH']


items = load('items.yaml', 'items')
gear = load('gear.yaml', 'items')
nodes = load('gathering.yaml', 'items')
recipes = load('recipes.yaml', 'recipes')
profs = {p['code']: p for p in load('professions.yaml', 'professions')}
mobs = load('mobs.yaml', 'mobs')
dungeons = load('dungeons.yaml', 'dungeons')
vendors = load('vendors.yaml', 'vendors')
quests = load('quests.yaml', 'quests')
with io.open(os.path.join(CONTENT, 'i18n', 'tr.yaml'), encoding='utf-8') as f:
    TR = yaml.safe_load(f)
ZONES, RAIDS, ARCH = world_zones()

every = {i['code']: i for i in items + gear + nodes}
mob_by = {m['code']: m for m in mobs}


def tr(kind, code):
    return TR.get('%s.%s' % (kind, code), code)


def money(c):
    g, c = divmod(int(c), 10000)
    s, c = divmod(c, 100)
    return ' '.join(p for p in ('%dg' % g if g else '', '%ds' % s if s else '', '%dc' % c if c or not (g or s) else '') if p)


# --- Sources and sinks --------------------------------------------------------
vendor_sells = {s['template_code'] for v in vendors for s in v['sells']}
crafted_by = collections.defaultdict(list)  # output -> recipes
used_by = collections.defaultdict(list)  # input -> recipes
for r in recipes:
    crafted_by[r['output_template_code']].append(r)
    for i in r['inputs']:
        used_by[i['template_code']].append(r)
quest_mobs = {o['event'].split(':', 1)[1] for q in quests for o in q['objectives'] if o['event'].startswith('mob_killed:')}

dungeon_boss = {d['boss_mob_code']: d for d in dungeons}
wave_of = collections.defaultdict(list)
for d in dungeons:
    for w in d['waves']:
        wave_of[w['mob_code']].append(d['code'])

mob_zone, mob_arch = {}, {}
for z in ZONES:
    for code, _, arch, _ in z['mobs']:
        mob_zone[code], mob_arch[code] = z['code'], arch
for r in RAIDS:
    mob_zone[r['boss'][0]] = r['code']
    mob_arch[r['boss'][0]] = 'titan' if r['level'] <= 300 else 'elder_titan'


def gear_tier(level, cls=None):
    """item.GearPool's top tier for a drop level (class-agnostic armor counts)."""
    top = 0
    for g in gear:
        if g.get('requirements') or g['min_level'] > level:
            continue
        if cls and g.get('class_restriction') not in (None, '', cls):
            continue
        top = max(top, g['min_level'])
    return top


hunting = sorted((n for n in nodes if n['gather']['kind'] == 'hunting'), key=lambda n: (n['gather']['level'], n['code']))


def hides(mob_level):
    """item.CombatHides: common span and best animal for a mob level."""
    reach = min(100, max(1, mob_level * 4))
    unlocked = [n for n in hunting if n['gather']['level'] <= reach]
    return unlocked[-HIDE_SPAN:], unlocked[-1] if unlocked else None


gear_levels = sorted({g['min_level'] for g in gear if not g.get('requirements')})  # drop tiers
tier_name = {}
for g in gear:
    if g['item_subclass'] == 'chest' or (g['equip_slot'] == 'chest'):
        tier_name.setdefault(g['min_level'], tr('item', g['code']))


def drop_sources(code):
    """Where a template comes from, in words."""
    t = every[code]
    out = []
    if t.get('gather'):
        out.append('toplama (%s Sv %d)' % (t['gather']['kind'], t['gather']['level']))
        if t['gather']['kind'] == 'hunting':
            out.append('av/zindan derisi')
    if code in vendor_sells:
        out.append('dükkân')
    if code in crafted_by:
        out.append('üretim (%s)' % ', '.join(sorted({r['profession_code'] for r in crafted_by[code]})))
    if t['item_class'] in ('weapon', 'armor') and not t.get('requirements'):
        out.append('ganimet')
    return out


# --- Report -------------------------------------------------------------------
L = []
w = L.append
findings = []  # (severity, area, text)


def find(sev, area, text):
    findings.append((sev, area, text))


w('# İçerik Denetimi — eşya, üretim, ganimet, kademe, set, canavar\n')
w('> Otomatik üretilir: `python scripts/content_audit.py`. Elle düzenleme — betiği değiştir. '
  'Kurallar `server/internal/item` (GearPool, CombatHides, TierRarity) ve `activity.DropLevel` aynasıdır. '
  'Yol haritası: [content-roadmap.md](content-roadmap.md).\n')

# 0. Summary
non_raid = [m for m in mobs if not m.get('raid_boss')]
solo = [d for d in dungeons if d.get('group_max', 1) <= 1]
raids = [d for d in dungeons if d.get('group_max', 1) > 1]
w('## 0. Özet\n')
w('| Alan | Adet | Aralık |')
w('|---|---:|---|')
w('| Ekipman şablonu | %d | Sv %d–%d, %d kademe |' % (len(gear), gear_levels[0], gear_levels[-1], len(gear_levels)))
w('| Toplama kaynağı | %d | meslek Sv 1–100, kademe T0–T7 |' % len(nodes))
w('| Diğer eşya (külçe, iksir, yemek, eski malzeme) | %d | |' % len(items))
w('| Tarif | %d | %s |' % (len(recipes), ', '.join('%s %d' % (k, v) for k, v in collections.Counter(r['profession_code'] for r in recipes).most_common())))
w('| Canavar | %d | Sv %d–%d (+%d baskın boss’u) |' % (len(non_raid), min(m['level'] for m in non_raid), max(m['level'] for m in non_raid), len(mobs) - len(non_raid)))
w('| Zindan / baskın | %d / %d | Sv %d–%d |' % (len(solo), len(raids), min(d['min_level'] for d in dungeons), max(d['min_level'] for d in dungeons)))
w('| Görev | %d | Sv %d–%d |' % (len(quests), min(q['min_level'] for q in quests), max(q['min_level'] for q in quests)))
w('| Set | 0 | tasarım var (`docs/item-engine/set.md`), uygulanmadı |')
w('| Dükkân | %d | %d eşya |' % (len(vendors), len(vendor_sells)))
w('')

# 1. Gear tiers
w('## 1. Ekipman kademeleri (tier)\n')
w('Her kademe: 10 sınıf silahı (sınıfa kilitli) + 5 zırh parçası (sınıfsız). Kaynak sütunu: D = dükkân, Ü = üretim, G = ganimet.\n')
w('| Sv | Kademe (göğüs adı) | Silah | Zırh | Saldırı (kılıç) | Savunma (göğüs) | Kaynak | Bu kademeyi düşüren canavar Sv |')
w('|---:|---|---:|---:|---:|---:|---|---|')
by_level = collections.defaultdict(list)
for g in gear:
    by_level[g['min_level']].append(g)
mob_levels = sorted({m['level'] for m in non_raid})
for i, lv in enumerate(gear_levels):
    gs = by_level[lv]
    nxt = gear_levels[i + 1] if i + 1 < len(gear_levels) else 10 ** 9
    weapons = [g for g in gs if g['item_class'] == 'weapon']
    armor = [g for g in gs if g['item_class'] == 'armor']
    sword = next((g for g in weapons if g.get('class_restriction') == 'warrior'), None)
    chest = next((g for g in armor if g['equip_slot'] == 'chest'), None)
    atk = next((b['amount'] for b in sword['stat_bonuses'] if b['target'] == 'attack'), '') if sword else ''
    dfn = next((b['amount'] for b in chest['stat_bonuses'] if b['target'] == 'defense'), '') if chest else ''
    src = set()
    for g in gs:
        if g['code'] in vendor_sells:
            src.add('D')
        if g['code'] in crafted_by:
            src.add('Ü')
        if not g.get('requirements'):
            src.add('G')
    droppers = [m for m in mob_levels if lv <= m < nxt]
    rng = ('%d–%d' % (droppers[0], droppers[-1]) if len(droppers) > 1 else str(droppers[0])) if droppers else '—'
    w('| %d | %s | %d | %d | %s | %s | %s | %s |' % (lv, tier_name.get(lv, gs[0]['code']), len(weapons), len(armor), atk, dfn,
                                                    ''.join(sorted(src)), rng))
    if not droppers and not (src - {'G'}):
        find('Yüksek', 'Ekipman', 'Sv %d kademesini (%s) düşüren canavar yok ve başka kaynağı yok.' % (lv, tier_name.get(lv, '')))
w('')
gaps = [(a, b) for a, b in zip(gear_levels, gear_levels[1:]) if b - a > 15]
w('Gereksinimli (rastgele düşmeyen) eşyalar: %s.\n' % ', '.join('%s (Sv %d)' % (tr('item', g['code']), g['min_level']) for g in gear if g.get('requirements')))
w('Kademe aralığı: Sv 1–25 arası 5, 25–95 arası 10–15, 95 sonrası 25 seviye. '
  'En geniş boşluk: %s.\n' % ', '.join('Sv %d→%d' % g for g in gaps[:3]))

slot_count = collections.Counter(g['equip_slot'] for g in gear)
w('### 1.1 Slot doluluğu\n')
w('| Slot | Şablon | Not |')
w('|---|---:|---|')
for s in SLOTS:
    n = slot_count.get(s, 0)
    note = '' if n >= len(gear_levels) else ('**boş**' if n == 0 else 'yalnız %d şablon' % n)
    w('| %s | %d | %s |' % (s, n, note))
w('')
empty = [s for s in SLOTS if slot_count.get(s, 0) == 0]
if empty:
    find('Yüksek', 'Ekipman', '12 ekipman slotundan %d tanesi boş: %s. Karakter ekranında hep boş görünüyor, güç bütçesinin ~%%30’u kullanılmıyor.'
         % (len(empty), ', '.join(empty)))
if slot_count.get('offhand', 0) < 3:
    find('Orta', 'Ekipman', 'Yan el (offhand) yalnız %d şablon (Titan Siper Kalkanı, Sv 600, koruyucu uzmanlığı). Kalkan izni olan savaşçı/paladin/rahip 600’e kadar kalkan bulamıyor.'
         % slot_count.get('offhand', 0))

# 1.2 families
with io.open(os.path.join(CONTENT, 'progression.yaml'), encoding='utf-8') as f:
    prog = yaml.safe_load(f)
families = prog.get('equipment_families', [])
present = collections.Counter(g['item_subclass'] for g in gear)
w('### 1.2 Silah aileleri ve zırh türleri\n')
w('| Aile | Şablon |')
w('|---|---:|')
for fam in families:
    w('| %s | %d |' % (fam, present.get(fam, 0)))
for a in prog.get('armor_types', []):
    w('| %s_armor | %d |' % (a, present.get(a + '_armor', 0)))
w('')
missing_fam = [f for f in families if not present.get(f)]
find('Orta', 'Ekipman', '%d/%d silah ailesinin hiç eşyası yok (%s …). Sınıf izinleri (DD-252) her sınıfa 4–6 aile tanıyor, içerikte sınıf başı tek aile var.'
     % (len(missing_fam), len(families), ', '.join(missing_fam[:8])))
if not any(present.get(a + '_armor') for a in prog.get('armor_types', [])):
    find('Orta', 'Ekipman', 'Zırh türleri (kumaş/hafif/orta/ağır) tanımlı ama hiçbir zırhın alt sınıfı bu türlerden değil; `CanUse` aile dışı alt sınıfa izin verdiği için büyücü ile savaşçı aynı zırhı giyiyor.')

# 1.3 rarity
w('### 1.3 Nadirlik ve ek (affix)\n')
w('| Nadirlik | Değer çarpanı | Ek kuralı | Ekipman ganimetinde |')
w('|---|---:|---|---|')
w('| sıradan | ×1 | yok | %60 |')
w('| seçkin | ×1.5 | yok | %25 · boss %50 |')
w('| nadir | ×3 | %50 tek ek | %11 · boss %35 |')
w('| destansı | ×8 | ön ek + %50 son ek | %4 · boss %13 |')
w('| efsanevi | ×25 | ön + son ek | boss %2 |')
w('| mitik | ×100 | ön + son ek (+ benzersiz etki ayrılmış) | **hiç düşmez** (yalnız T7 malzeme) |')
w('')
find('Yüksek', 'Nadirlik', 'Ekler (sharp, of_the_bear …) yalnız isim: `item/affix.go` “cosmetic only”; stat vermiyor. Nadirlik ekipmanın gücünü değiştirmiyor, yalnız satış değerini — destansı kılıç ile sıradan kılıç aynı vuruyor.')
find('Orta', 'Nadirlik', 'Mitik ekipman hiçbir kaynaktan gelmiyor (boss tablosu efsanevide bitiyor); `unique-effect.md` tasarımı uygulanmadı.')

# 2. Item list
w('## 2. Eşya listesi\n')
w('### 2.1 Toplama kaynakları\n')
w('| Kod | Ad | Tür | Meslek Sv | Kademe | Nadirlik | Değer | Kullanan tarif |')
w('|---|---|---|---:|---:|---|---:|---|')
for n in sorted(nodes, key=lambda n: (n['gather']['kind'], n['gather']['level'])):
    g = n['gather']
    users = ', '.join(sorted({tr('item', r['output_template_code']) for r in used_by.get(n['code'], [])})) or '—'
    w('| %s | %s | %s/%s | %d | T%d | %s | %s | %s |' % (n['code'], tr('item', n['code']), g['kind'], n['item_subclass'],
                                                       g['level'], g['tier'], RARITY_TR[g['tier']], money(n['base_value_copper']), users))
w('')
unused = collections.defaultdict(list)
for n in nodes:
    if not used_by.get(n['code']):
        unused[n['item_subclass']].append(n['code'])
w('Kullanılmayan toplama kaynakları (hiçbir tarifte girdi değil — yalnız satılır):\n')
w('| Alt tür | Kullanılmayan / toplam |')
w('|---|---|')
for sub, cnt in sorted(collections.Counter(n['item_subclass'] for n in nodes).items()):
    w('| %s | %d / %d |' % (sub, len(unused[sub]), cnt))
w('')
total_unused = sum(len(v) for v in unused.values())
find('Yüksek', 'Ekonomi', '%d/%d toplama kaynağı hiçbir tarifte kullanılmıyor (deri/post/tüy/pul %d, taş %d, balık %d, ot %d, kütük %d). Avcılık mesleği ürününün tek değeri dükkân fiyatı.'
     % (total_unused, len(nodes), len(unused['pelt']) + len(unused['hide']) + len(unused['feather']) + len(unused['scale']),
        len(unused['stone']), len(unused['fish']), len(unused['herb']), len(unused['log'])))

w('### 2.2 Diğer eşyalar\n')
w('| Kod | Ad | Sınıf | Sv | Değer | Etki | Kaynak | Kullanan tarif |')
w('|---|---|---|---:|---:|---|---|---|')
retired = []
for i in sorted(items, key=lambda i: (i['item_class'], i.get('item_subclass', ''), i['min_level'])):
    eff = i.get('use_effect') or {}
    e = ''
    if eff:
        e = '+%d can' % eff.get('amount', 0)
        if eff.get('buff'):
            e += ', +%%%d sal/sav %d dk' % (eff['buff']['attack_pct'], eff['buff']['minutes'])
    src = ', '.join(drop_sources(i['code'])) or '**yok**'
    users = ', '.join(sorted({tr('item', r['output_template_code']) for r in used_by.get(i['code'], [])})) or '—'
    w('| %s | %s | %s/%s | %d | %s | %s | %s | %s |' % (i['code'], tr('item', i['code']), i['item_class'], i.get('item_subclass', ''),
                                                      i['min_level'], money(i['base_value_copper']), e, src, users))
    if src == '**yok**':
        retired.append(i['code'])
w('')
if retired:
    find('Düşük', 'Eşya', '%d eşya bilinçli olarak emekli (DD-241/244; eski yığınlar satılabilsin diye duruyor), kaynağı yok: %s. Temizlik isteğe bağlı.'
         % (len(retired), ', '.join(retired)))
bars = [i for i in items if i.get('item_subclass') == 'bar']
idle_bars = [b['code'] for b in bars if not used_by.get(b['code'])]
if idle_bars:
    find('Yüksek', 'Üretim', '%d/%d külçe hiçbir tarifte kullanılmıyor (%s …). Demircilik tarifleri külçe değil ham cevher istiyor; ergitme (24 tarif, beceri 100’e kadar) çıkmaz sokak.'
         % (len(idle_bars), len(bars), ', '.join(idle_bars[:6])))
potions = [i for i in items if i.get('item_subclass') == 'potion' and i.get('use_effect')]
top_heal = max(i['use_effect']['amount'] for i in potions)
top_food = max(i['use_effect']['amount'] for i in items if i.get('item_subclass') == 'food')
hp1000 = 150 + 25 * 1000
find('Orta', 'Tüketilebilir', 'En büyük iksir +%d, en büyük yemek +%d can; Sv 1000 canavar canı ~%d. İksir Sv 28, yemek Sv 100’de bitiyor — üst seviyede iyileştirme anlamsız.'
     % (top_heal, top_food, hp1000))

# 3. Crafting
w('## 3. Üretim listesi\n')
for pc in sorted({r['profession_code'] for r in recipes}):
    rs = sorted((r for r in recipes if r['profession_code'] == pc), key=lambda r: r['min_skill'])
    out_lv = [every[r['output_template_code']]['min_level'] for r in rs if r['output_template_code'] in every]
    w('### %s — %d tarif, beceri %d–%d (tavan %d), çıktı Sv %d–%d\n' % (
        tr('profession', pc), len(rs), rs[0]['min_skill'], rs[-1]['min_skill'], profs[pc]['max_level'], min(out_lv), max(out_lv)))
    w('| Tarif | Beceri | Girdi | Çıktı | Çıktı Sv | Süre | XP | Ücret |')
    w('|---|---:|---|---|---:|---:|---:|---:|')
    for r in rs:
        ins = ', '.join('%d× %s' % (i['quantity'], tr('item', i['template_code'])) for i in r['inputs'])
        o = every.get(r['output_template_code'], {})
        w('| %s | %d | %s | %d× %s | %s | %ds | %d | %s |' % (r['code'], r['min_skill'], ins, r.get('output_stack', 1),
                                                          tr('item', r['output_template_code']), o.get('min_level', '?'),
                                                          r.get('craft_seconds', 0), r['profession_xp_reward'], money(r['gold_cost_copper'])))
    w('')
    top = rs[-1]['min_skill']
    if top < profs[pc]['max_level'] * 0.8:
        find('Yüksek', 'Üretim', '%s: son tarif beceri %d, meslek tavanı %d — son %d beceri seviyesinde açılacak tarif yok.'
             % (tr('profession', pc), top, profs[pc]['max_level'], profs[pc]['max_level'] - top))
smith = [every[r['output_template_code']]['min_level'] for r in recipes if r['profession_code'] == 'blacksmithing']
find('Yüksek', 'Üretim', 'En yüksek üretilebilir ekipman Sv %d; ekipman merdiveni Sv %d. Sv %d üstü ekipmanın tek kaynağı %%4 av / boss ganimeti.'
     % (max(smith), gear_levels[-1], max(smith)))
crafted_gear_classes = {every[r['output_template_code']].get('class_restriction') for r in recipes if r['profession_code'] == 'blacksmithing'} - {None}
all_classes = {g.get('class_restriction') for g in gear} - {None}
if all_classes - crafted_gear_classes:
    find('Orta', 'Üretim', 'Üretilebilir silahı olmayan sınıflar: %s (yalnız %s silahı demircilikte).'
         % (', '.join(sorted(all_classes - crafted_gear_classes)), ', '.join(sorted(crafted_gear_classes))))
no_gather_prof = [p for p in profs.values() if p['code'] not in {r['profession_code'] for r in recipes} and p['max_level'] and p['code'] not in
                  ('mining', 'fishing', 'herbalism', 'woodcutting', 'hunting')]
for p in no_gather_prof:
    find('Orta', 'Üretim', '%s mesleği var ama tarifi yok.' % tr('profession', p['code']))
find('Orta', 'Üretim', 'Kuyumculuk / terzilik / deri işçiliği / oymacılık yok: deri, post, tüy, kütük ve taş için üretim yolu yok; takı slotları boş.')

# 4. Drops
w('## 4. Ganimet (drop) kuralları\n')
w('- **Ganimet seviyesi** = min(canavar Sv, karakter Sv) (`activity.DropLevel`).')
w('- **Ekipman**: ganimet seviyesine eşit ya da altındaki **en yüksek tek kademe**, karakterin sınıfına uygun (`item.GearPool`; gereksinimli eşya düşmez).')
w('  - Av: her zafer %%4 → %s.' % GEAR_ODDS)
w('  - Zindan: dalga öldürmesi başına %%4 + boss’ta 1 garanti → boss tablosu %s.' % BOSS_ODDS)
w('  - Baskın: dalga %4 + boss’ta 2 garanti, kişisel ganimet.')
w('- **Deri/post** (`item.CombatHides`): avcılık seviyesi = canavar Sv × 4 (tavan 100). Sıradan buluntu açılan son 4 hayvandan, nadir/efsanevi en iyisinden.')
w('- **Toplama**: seçilen kaynağın kendisi, kademe rengiyle (`TierRarity`): T0–1 sıradan, T2–3 seçkin, T4 nadir, T5 destansı, T6 efsanevi, T7 mitik.')
w('- **Altın / XP**: canavar tablosu; XP seviye farkı bantları (10 aşağı tam, 25 aşağı %50, ötesi %10).\n')
cap_mob = next(m for m in sorted(non_raid, key=lambda m: m['level']) if m['level'] * 4 >= 100)
find('Yüksek', 'Ganimet', 'Deri ganimeti Sv %d canavardan itibaren tavanda (avcılık 100): Sv %d–1000 arası tüm canavarlar aynı 4 hayvanın derisini düşürüyor, en iyisi %s.'
     % (cap_mob['level'], cap_mob['level'], tr('item', hides(1000)[1]['code'])))
find('Orta', 'Ganimet', 'Canavara özgü ganimet tablosu yok: aynı seviyedeki ejderha ile haydut aynı şeyi düşürüyor (yalnız seviye belirliyor). Görev/zanaat için özel malzeme (ejder pulu, boss mührü) yok.')
find('Orta', 'Ganimet', 'Ganimet tek kademeden: Sv 34 canavar yalnız Sv 25 kademesi düşürür; bir alt kademe/üst kademe şansı yok. Kademe sınırını geçince düşen eşya bir anda değişiyor.')

# 5. Monsters
w('## 5. Canavar listesi\n')
w('Bölge, `scripts/gen_world.py` ZONES’undan; Sv 1–25 canavarları elle yazılmış (`mobs.yaml` başı), bölgesi yok. '
  'Av ve zindan başlatmak bulunulan bölgeye bağlı değil (yolculuk yalnız konum değiştiriyor).\n')
w('| Sv | Kod | Ad | Arketip | Bölge | Can | Sal | Sav | XP | Altın | Ekipman kademesi | Deri (en iyi) | Rol |')
w('|---:|---|---|---|---|---:|---:|---:|---:|---|---:|---|---|')
for m in sorted(mobs, key=lambda m: (m['level'], m['code'])):
    roles = []
    if m['code'] in dungeon_boss:
        roles.append(('baskın boss’u ' if m.get('raid_boss') else 'boss ') + tr('dungeon', dungeon_boss[m['code']]['code']))
    if wave_of.get(m['code']):
        roles.append('dalga ×%d' % len(wave_of[m['code']]))
    if m['code'] in quest_mobs:
        roles.append('görev')
    best = hides(m['level'])[1]
    w('| %d | %s | %s | %s | %s | %d | %d | %d | %d | %s–%s | %d | %s | %s |' % (
        m['level'], m['code'], tr('mob', m['code']), mob_arch.get(m['code'], '—'), (tr('dungeon', mob_zone[m['code']]) if m.get('raid_boss') else tr('region', mob_zone[m['code']])) if m['code'] in mob_zone else '—',
        m['hp_max'], m['attack'], m['defense'], m['xp_reward'], money(m['gold_copper_min']), money(m['gold_copper_max']),
        gear_tier(m['level']), tr('item', best['code']) if best else '—', ', '.join(roles)))
w('')
lv = sorted(m['level'] for m in non_raid)
big = sorted(((b - a, a, b) for a, b in zip(lv, lv[1:])), reverse=True)[:3]
w('En geniş canavar seviye boşlukları: %s.\n' % ', '.join('Sv %d→%d' % (a, b) for _, a, b in big))
arch_count = collections.Counter(mob_arch.get(m['code'], 'el yapımı') for m in mobs)
w('Arketip dağılımı: %s.\n' % ', '.join('%s %d' % kv for kv in arch_count.most_common()))
find('Orta', 'Canavar', 'Av bölgeye bağlı değil: bölgeler yalnız yolculuk ekranında var; oyuncu herhangi bir yerden herhangi bir canavarı avlayabiliyor (DD-254/255 dünyaları anlamını yitiriyor).')
find('Düşük', 'Canavar', 'Her bölgede 6 canavar → Sv 100 sonrası her ~10 seviyede bir canavar; yan yana seçenek az (aynı seviyede ikinci canavar yok).')

# 6. Dungeons
w('## 6. Zindanlar ve baskınlar\n')
w('| Sv | Kod | Ad | Kişi | Süre | Dalgalar | Boss | Boss Sv | Bonus altın | Bonus XP | Boss ekipman kademesi |')
w('|---:|---|---|---|---:|---|---|---:|---|---:|---:|')
for d in sorted(dungeons, key=lambda d: (d['min_level'], d['code'])):
    boss = mob_by[d['boss_mob_code']]
    waves = ', '.join('%d× %s' % (x['count'], tr('mob', x['mob_code'])) for x in d['waves'])
    grp = '%d–%d' % (d.get('group_min', 1), d['group_max']) if d.get('group_max', 1) > 1 else '1'
    w('| %d | %s | %s | %s | %d dk | %s | %s | %d | %s | %d | %d |' % (d['min_level'], d['code'], tr('dungeon', d['code']), grp,
                                                                       d['duration_seconds'] // 60, waves, tr('mob', boss['code']), boss['level'],
                                                                       money(d['bonus_gold_copper']), d['bonus_xp'], gear_tier(boss['level'])))
w('')
find('Orta', 'Zindan', 'Zindana özgü ganimet yok: boss ekipmanı açık dünyadaki aynı kademe; zindan/baskın seti, mührü ya da benzersiz eşyası yok.')

# 7. Sets
w('## 7. Setler\n')
w('İçerikte set yok. Tasarım `docs/item-engine/set.md` (parça sayımı, kademeli bonus, imza tetikleyici, `content/item_sets/*.yaml`) hazır, motor ve içerik yok. '
  'Aynı malzeme adını taşıyan 5 zırh parçası (ör. demir başlık/göğüs/bacak/çizme/eldiven) set adayı.\n')
find('Yüksek', 'Set', 'Set sistemi yok (tasarım var). Zindan/baskın ödülü ve uzun vadeli hedef için en büyük boşluk.')
find('Orta', 'Eşya motoru', 'Tasarımı olup uygulanmayan: soket + mücevher (`gems.md`, `sockets.md`), büyü (`enchantments.md`), üretim kalitesi / kritik üretim (`quality.md`, `critical-craft.md`), benzersiz etki (`unique-effect.md`).')

# 8. Findings
w('## 8. Uyumsuzluklar ve eksikler\n')
order = {'Yüksek': 0, 'Orta': 1, 'Düşük': 2}
w('| # | Önem | Alan | Bulgu |')
w('|---:|---|---|---|')
for n, (sev, area, text) in enumerate(sorted(findings, key=lambda f: (order[f[0]], f[1])), 1):
    w('| %d | %s | %s | %s |' % (n, sev, area, text))
w('')

with io.open(OUT, 'w', encoding='utf-8', newline='\n') as f:
    f.write('\n'.join(L))
print('findings', len(findings), '->', os.path.relpath(OUT, ROOT))
for sev, area, text in sorted(findings, key=lambda f: (order[f[0]], f[1])):
    print('-', sev, area, re.sub(r'\s+', ' ', text)[:160])
