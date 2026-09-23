# Placeholder Content — v1 Class + Specialization Listesi

> ⚠️ **GEÇİCİ / PLACEHOLDER**
> Bu doküman v1 içerik hazırlığı için placeholder isim + kısa tanımlar içerir. **Kullanıcı revize edecek**. Final lore, isim, dengeleyici sayılar Phase 01 balance pass sonrası.

## 1. 6 Base Class

| Kod | Placeholder Ad | Rol | Kaynak | Kısa Kimlik |
|---|---|---|---|---|
| `warrior` | Warrior | Tank / Melee DPS | Rage | Klasik ön saf, kılıç + kalkan |
| `ranger` | Ranger | Ranged DPS / Skirmish | Focus | Yay + hafif zırh, hunting affinity |
| `mage` | Mage | Ranged Caster DPS | Mana | Staff + kumaş, elemental burst |
| `cleric` | Cleric | Support / Off-heal | Mana | Kalkan + topuz, party sustain |
| `rogue` | Rogue | Melee DPS / Burst | Energy | Hançer + hafif zırh, execute |
| `druid` | Druid | Hybrid (dps/support) | Nature | Doğa büyüsü, adventure affinity |

Her class için 2 Tier I × 3 Tier II = 6 leaf identity → **36 unique build**.

---

## 2. WARRIOR

```
WARRIOR
 ├── Berserker           (offense)
 │   ├── Ravager         (dual wield sustain)
 │   ├── Warlord         (aoe + party rally)
 │   └── Blood Reaver    (execute + lifesteal)
 └── Guardian            (defense)
     ├── Bulwark         (shield master, single-target tank)
     ├── Sentinel        (crowd control tank)
     └── Warden          (party protector, damage share)
```

### 2.1 Berserker (Tier I)
- Passive: `berserker_frenzy` — düşük HP'de +%15 attack speed.
- Active unlock: `berserker_rampage` — 20 sn +%30 damage.
- Stat bias: STR+3.

**Tier II:**
| Kod | İsim | Signature | Kısa |
|---|---|---|---|
| `ravager` | Ravager | `ravager_final_slaughter` | Dual wield burst; execute mob under %20 HP. |
| `warlord` | Warlord | `warlord_battle_cry` | AoE damage + party +%10 dmg for 15sn. |
| `blood_reaver` | Blood Reaver | `blood_reaver_soul_drain` | Lifesteal focus; +%15 lifesteal 20 sn. |

### 2.2 Guardian (Tier I)
- Passive: `guardian_iron_will` — +%20 threat, +%10 armor.
- Active: `guardian_taunt` — target enemy'yi 3 sn taunt.
- Stat bias: STA+3.

**Tier II:**
| Kod | İsim | Signature | Kısa |
|---|---|---|---|
| `bulwark` | Bulwark | `bulwark_shield_wall` | 8 sn %50 damage reduction. |
| `sentinel` | Sentinel | `sentinel_shockwave` | AoE stun 2 sn. |
| `warden` | Warden | `warden_guardian_bond` | Party içinde 1 hedefin aldığı hasarın %30'u sana. |

---

## 3. RANGER

```
RANGER
 ├── Marksman            (single target precision)
 │   ├── Deadeye         (crit burst, snipe)
 │   ├── Hunter          (pet + track)
 │   └── Ranger Captain  (mobility + kite)
 └── Skirmisher          (mobility + traps)
     ├── Trapmaster      (crowd control ile trap)
     ├── Duskstalker     (stealth openers)
     └── Wildkin         (nature summon)
```

### Placeholder pass — Tier II detay TBD by user; her Tier II için signature skill + 1 passive + 1 active gerekli.

---

## 4. MAGE

```
MAGE
 ├── Elementalist        (raw damage)
 │   ├── Firecaller      (DoT + AoE burn)
 │   ├── Stormbringer    (lightning + haste)
 │   └── Frostweaver     (slow + control)
 └── Mystic              (utility + support)
     ├── Chronomancer    (haste manipulation + rewind)
     ├── Illusionist     (fear + confuse)
     └── Arcanist        (mana burn + shield)
```

---

## 5. CLERIC

```
CLERIC
 ├── Priest              (heal focus)
 │   ├── Radiant         (single target heal + burst)
 │   ├── Sanctifier      (aoe heal + shield)
 │   └── Confessor       (heal-through-damage)
 └── Templar             (heal + hybrid melee)
     ├── Crusader        (holy damage + heal on hit)
     ├── Warpriest       (party buff + off-heal)
     └── Inquisitor      (undead specialist + damage)
```

---

## 6. ROGUE

```
ROGUE
 ├── Assassin            (single burst)
 │   ├── Shadowblade     (stealth openers)
 │   ├── Envenom         (poison DoT)
 │   └── Slayer          (execute focus)
 └── Duelist             (dual wield sustain)
     ├── Fencer          (parry-riposte)
     ├── Corsair         (mobility + pistol Phase 02)
     └── Trickster       (dodge + counter)
```

---

## 7. DRUID

```
DRUID
 ├── Wildshaper          (form-based combat)
 │   ├── Beastform       (melee bear form)
 │   ├── Sunform         (caster owl form — sustained dps)
 │   └── Tideform        (aquatic support form)
 └── Grovekeeper         (support + nature)
     ├── Bloombearer     (heal-over-time + HoT stack)
     ├── Thornwarden     (dot + slow)
     └── Skywalker       (winged form travel + kite)
```

---

## 8. Content Pack Örneği (Yaml)

Tam yaml pack örneği (Warrior/Berserker/Ravager):

```yaml
# content/classes/warrior.yaml
code: warrior
version: 1
display_name_key: class.warrior.name
lore_key: class.warrior.lore
art_ref: classes/warrior/portrait.png

resource_type: rage
armor_affinity: plate
weapon_affinities: [sword, axe, mace, shield, polearm, fist]
role_tags: [tank, dps_melee]

base_stat_bias:
  str: 3
  sta: 2

starter_skill_codes:
  - warrior_bash
  - warrior_defend

starter_home_city_code: ashfield
tags: [melee, frontline]

specializations_tier1:
  - berserker
  - guardian
```

```yaml
# content/specializations/berserker.yaml
code: berserker
version: 1
class_code: warrior
tier: 1
display_name_key: spec.berserker.name
lore_key: spec.berserker.lore

passive_skill_codes:
  - berserker_frenzy
active_skill_unlock_codes:
  - berserker_rampage

stat_bias:
  str: 3

playstyle_tags: [offense, sustain_dmg]

default_rotation:
  - cond: { hp_below_pct: 30 }
    skill_code: berserker_rampage
  - cond: { any: true }
    skill_code: warrior_bash

tier2_options:
  - ravager
  - warlord
  - blood_reaver
```

```yaml
# content/specializations/ravager.yaml
code: ravager
version: 1
class_code: warrior
tier: 2
parent_spec_code: berserker
display_name_key: spec.ravager.name

passive_skill_codes:
  - ravager_bloodthirst
active_skill_unlock_codes:
  - ravager_execute
signature_skill_code: ravager_final_slaughter

stat_bias:
  str: 2

playstyle_tags: [execute, dual_wield, low_hp_bonus]

default_rotation:
  - cond: { target_hp_below_pct: 20 }
    skill_code: ravager_execute
  - cond: { self_hp_below_pct: 40 }
    skill_code: berserker_rampage
  - cond: { any: true }
    skill_code: warrior_bash
```

## 9. Skill Placeholder Set

Her class başına ~15-20 skill (starter + passive + spec unlock + signature). Toplam ~120 skill v1 hedefi.

Skill isimleri placeholder — content team revizyon.

## 10. Balance Kontrol Listesi

Content pack finalize öncesi:
- [ ] Her class 2 Tier I + 6 Tier II.
- [ ] Stat bias budget uyumu (class +5, spec +3/+2).
- [ ] Her Tier II signature skill var.
- [ ] Default rotation makul (AFK-friendly).
- [ ] i18n key'ler present (tr-TR + en-US en az).
- [ ] Art reference (portrait, badge, icon) planı.
- [ ] Class-race sinerji simülasyonu (dominant "must-play" kombo yok).

## 11. Rework Discipline

Kullanıcı revizyonu (bu placeholder → final):
1. Yeni yaml pack yaz (draft branch).
2. Loader validate (CI).
3. Impact simulation (mevcut karakter distribution).
4. Compensation plan (free respec tokens).
5. Staging deploy.
6. Community feedback (Phase 02+ playtest).
7. Production rollout.

## 12. İlişkiler
- [Class](class.md), [Specialization](specialization.md), [Respec](respec.md), [README](README.md)
- [game-design/classes.md](../game-design/classes.md), [game-design/specializations.md](../game-design/specializations.md), [game-design/skills.md](../game-design/skills.md)
- [database/content-templates.md](../database/content-templates.md)

## 13. Notlar (User Revizyon İçin)

- **İsim tutarlılığı**: fantasy MMO klişesi vs özgün naming — kullanıcı tercihi belirtildiğinde revize.
- **IP çakışması**: bazı isimler (Warlord, Assassin) yaygın kullanılıyor; trademark risk yok ama distinctive olmak isterse rename.
- **Tema**: bazı specler "kanlı" (Blood Reaver, Bloodthirst) — ton yönetimi karar.
- **Alt-mekanikler**: signature skill fikirleri Phase 8+ derinleştirilebilir.
- **Boş yerler**: Ranger/Mage/Cleric/Rogue/Druid Tier II detay boş bırakıldı — placeholder çerçeve ile başla, kullanıcı yönlendirmesi bekleniyor.

Tüm bu içerik `content/` klasöründe yaml'lara dönüşecek; yaml patch'i git commit ile version'lu, content_version bump ile canlıya alınır (bkz. [architecture/postgresql.md §Content Data Loading](../architecture/postgresql.md)).
