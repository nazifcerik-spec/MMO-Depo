# Racial Skills (Active Abilities)

## 1. Amaç
Race'in **aktif** kullanılabilir yeteneği. Class ability bütçesinden küçük; sinerji katmanı; kimlik hissi güçlendirir.

## 2. Model
- `content.race.active_ability_id` FK → `content.skill_template`.
- `skill_template.type = 'racial'`.
- Character oluşturulduğunda `character.skill_learn` (rank=1) + `character.skill_bar` slot 11 (son slot reserve).

## 3. Ortak Kural Seti

### 3.1 Cooldown
- Baz cooldown: **180..300 sn** (3-5 dk).
- CDR ile azalabilir (min 60 sn cap).
- Uzun cooldown = "büyük moment" hissi.

### 3.2 Resource Cost
- Racial ability **kaynak tüketmez** (rage/mana/energy).
- İstisna: özel race (arcanborn `arcane_surge` mana yer).

### 3.3 Cast Time
- Genellikle **instant** (0 cast time).
- Sadece `arcanborn` gibi caster-race'lerde 0.5-1.0 sn cast time.

### 3.4 Global Cooldown Bypass
- Bazı utility racial'lar (dwarf `stone_skin` defensive burst) GCD bypass.
- Combat DPS racial'lar GCD respect.

### 3.5 Range
- Self target (buff) çoğunlukla.
- AoE (rally) party wide.
- Line/cone (attack) nadir.

## 4. Effect Bütçesi
- **Combat racial**:
  - Damage burst = weapon damage × 1.5 (short burst) veya
  - Buff = +%5-10 stat / +%10 damage output 15-30 sn.
- **Utility racial**:
  - Defense burst: +50% armor 8 sn (dwarf stone_skin) veya
  - Movement: +40% speed 6 sn (elf wind_step) veya
  - Support: party +5% stat 20 sn (human rally).
- Racial ability class'ın ultimate'ından **hep daha küçük**.

## 5. Data Format

Skill template yaml (örnek human):
```yaml
# content/skills/human_rally.yaml
code: human_rally
type: racial
race_code: human
display_name_key: race.human.active.name
description_key: race.human.active.desc

resource_type: null
resource_cost: 0
cooldown_ms: 180000
gcd_bypass: false
cast_time_ms: 0
range_kind: mid
aoe_shape: circle
aoe_radius: 15

effects:
  _v: 1
  effects:
    - kind: buff
      target: party_members
      buff:
        target: primary_all
        amount_pct: 5
        duration_ms: 20000
        stacks: 1
    - kind: buff
      target: self
      buff:
        target: xp_gain_pct
        amount: 5
        duration_ms: 60000
        stacks: 1

scaling: {}

unlock_condition:
  min_level: 1
  race_required: human

max_rank: 1
```

## 6. v1 Race Active Skills (Placeholder)

| Race | Active Skill | Effect | Cooldown |
|---|---|---|---|
| `human` | `human_rally` | Party +5% primary stat + self +5% XP for 20 sn | 3 min |
| `dwarf` | `dwarven_stone_skin` | Self +50% armor for 8 sn | 3 min |
| `elf` | `elven_wind_step` | Self +40% movement + evasion +20% for 6 sn | 4 min |
| `orc` | `orcish_blood_roar` | Self +15% damage + fear enemies 3 sn (aoe 8m) | 5 min |
| `merfolk` | `tide_grace` | Self +25% evasion + cleanse root/slow 8 sn | 4 min |
| `arcanborn` | `arcane_surge` | Self next 3 spells +25% dmg, no cast time | 5 min |

## 7. Character Skill Bar

Character oluşturulduğunda:
- Slot 0..9: class starter skills.
- Slot 10: opsiyonel spec skill slot.
- **Slot 11: racial active** (fixed default; oyuncu kaldırabilir).

Rotation editor: racial ability priority listesine eklenebilir; ancak default örnek rotation'da yer almaz.

## 8. Cast & Effect Flow

Skill cast (bkz. [game-design/skills.md](../game-design/skills.md)):
```
Client → POST /v1/characters/{id}/skills/{skill_id}/cast
         (or WS 'skill_cast' frame)
  ↓
Server:
1. Load character, skill, current context.
2. Cooldown check (Redis `char:cd:<char>:<skill_id>`).
3. Resource check (0 for racial).
4. Range/target validate.
5. Apply effects (buff insert, damage compute).
6. Set cooldown in Redis.
7. Publish 'combat.skill_used' event.
8. WS push 'skill_effect' frame.
```

Cooldown Redis: `char:cd:<char>:<skill_id>` → `1 EX <cooldown_sec>`.

## 9. AFK Racial Skills

AFK aktivite sırasında racial'lar rotation'a dahil ise cast edilir:
- Combat racial (orc blood_roar): normal rotation eval.
- Utility racial (elf wind_step): kullanışsız aktif rotation'da; ancak kesin kaçınma yok.
- Human rally: solo AFK'de kendine buff kısmı çalışır.

## 10. Passive vs Active Ayrımı
- Passive (bkz. [passive.md](passive.md)): her zaman aktif, cost yok, cooldown yok.
- Active: manuel/rotation tetiklenir, cooldown var.

Bir race'in ikisi de vardır — her ikisi de character.skill_learn'a insert.

## 11. Endpoint

**GET `/v1/content/skills/{code}`** (public)
Skill detay (i18n key + effect summary).

**POST `/v1/characters/{id}/skills/{skill_id}/cast`**
Auth + CharacterScope.
Body:
```json
{ "target_kind": "self|party|enemy", "target_id": "..." }
```

## 12. Fail Modları
- Cooldown active → 400 `skill.cooldown_active` + remaining ms.
- Resource insufficient → 400 (racial'lar için nadir).
- Race mismatch (skill başka race'e ait) → 403 `skill.race_mismatch`.
- Out of range → 400.

## 13. Metrics
- `racial_skill_cast_total{race,skill}`.
- `racial_skill_cooldown_avg`.
- `racial_skill_effectiveness` — weekly meta report.

## 14. Audit
- `combat.skill_used` (high volume; debug).
- Suspicious rapid-cast → security signal.

## 15. Test Stratejisi
- Cooldown enforce (2 cast in <cooldown → 2. reject).
- Effect apply/expire timing.
- AoE target selection (party members).
- Race mismatch reject.

## 16. Balance Considerations
- Racial ability class rotation'ın kalbi olmamalı; opsiyonel katman.
- Cooldown 3+ dk → dungeon boss başına 1 kez tercih.
- PvP (Phase 02+) balance ayrı; racial abilities öldürücü combo olmamalı.

## 17. İlişkiler
- [Race](race.md), [Passive](passive.md), [Racial Bonuses](racial-bonuses.md), [character/character.md](../character/character.md), [game-design/skills.md](../game-design/skills.md), [game-design/combat.md](../game-design/combat.md), [database/content-templates.md](../database/content-templates.md)

## 18. Bekleyen Kararlar
- Racial ability rank up mümkün mü (Phase 02 opsiyonel — evet öneri).
- Racial "spec" ekleme (Phase 03): race içinde alt seçim.
- Combined "racial + class" signature ability (Phase 03+).
