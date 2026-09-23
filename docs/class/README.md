# Class + Specialization System (Faz 8)

Karakterin oynanış kimliğinin merkezi. Base class → Tier I spec (2 seçenek) → Tier II spec (Tier I'e bağlı 3 seçenek) ağacı. **Data-driven**: yeni class/spec eklemek kod değişikliği gerektirmez.

Yüksek seviye vizyon: [game-design/classes.md](../game-design/classes.md), [game-design/specializations.md](../game-design/specializations.md). Bu klasör Phase 8'de implementation-hazır detay.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış + tree yapısı + data-driven (bu doküman) | [README.md](README.md) |
| [Class](class.md) | Base class entity + creation + starter kit |
| [Specialization](specialization.md) | Tier I & II unlock + tree + data-driven |
| [Respec](respec.md) | Spec değişim akışı + cost + cooldown |
| [Placeholder Content](placeholder-content.md) | **Geçici v1 class + spec listesi (revize edilecek)** |

## 2. Tree Yapısı

Her base class:
```
BASE CLASS
 ├── Spec I (option A)
 │   ├── Spec II option A1
 │   ├── Spec II option A2
 │   └── Spec II option A3
 └── Spec I (option B)
     ├── Spec II option B1
     ├── Spec II option B2
     └── Spec II option B3
```

**Toplam per base**: 2 Spec I × 3 Spec II = **6 leaf identity**.

**6 base class × 6 leaf = 36 unique identity** (v1 hedef).

Example (kullanıcının verdiği):
```
WARRIOR
 ├── Berserker
 │   ├── Ravager
 │   ├── Warlord
 │   └── Blood Reaver
 └── Guardian
     ├── ??? (3 slot)
     ├── ???
     └── ???
```

**Not:** Placeholder isimler [placeholder-content.md](placeholder-content.md) altında; user revize edecek.

## 3. Unlock Timing

- **Base class**: character creation'da seçilir; kalıcı (bkz. DD-004).
- **Spec I**: Level 5'te unlock; oyuncu seçer.
- **Spec II**: Level 10'da unlock; oyuncu Spec I'e bağlı 3 seçenekten seçer.

## 4. Data-Driven İlkeleri (Class + Spec)

1. **Kod class code'u bilmez** — domain code'da `if class == "warrior"` yasak (DD-057 racial ile aynı prensip).
2. **Effect DSL** — Class stat bias, spec passive/skill unlock hepsi content data.
3. **Content pack YAML** — `content/classes/*.yaml`, `content/specializations/*.yaml`.
4. **Hot reload** — content_version bump ile canlı sunucuya yansır.
5. **Backward compat** — mevcut karakter için spec deprecation compensating (free respec).

## 5. Content Pack Format

### 5.1 Base Class YAML
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

specializations_tier1:
  - berserker
  - guardian

# spec_tier2 seçimleri specialization YAML'larından parent link ile bağlanır
```

### 5.2 Specialization YAML (Tier 1)
```yaml
# content/specializations/berserker.yaml
code: berserker
version: 1
class_code: warrior
tier: 1
display_name_key: spec.berserker.name
lore_key: spec.berserker.lore

passive_skill_codes:
  - berserker_frenzy_passive
active_skill_unlock_codes:
  - berserker_rampage

stat_bias:
  str: 3

playstyle_tags: [offense, dual_wield, sustain_dmg]

tier2_options:
  - ravager
  - warlord
  - blood_reaver
```

### 5.3 Specialization YAML (Tier 2)
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

playstyle_tags: [execute, burst, low_hp_bonus]
```

## 6. Data Loader Validation

Content build (CI):
- `class_code` valid content.class kaydı.
- `parent_spec_code` tier 1 spec olmalı (tier 2 için).
- `tier2_options` listesi 3 element (v1 kesin; Phase 02'de 2-4 esnek).
- Skill code'ları var (skill_template).
- Stat bias budget check (bkz. §11).
- İsim/kod unique.

## 7. Runtime API

Character seçtiğinde:
- Level 5 → UI unlock notif → oyuncu Spec I seçer → `POST /v1/characters/{id}/specializations/tier1`.
- Level 10 → UI unlock notif → Spec II seçer → `POST /v1/characters/{id}/specializations/tier2`.

Detay: [specialization.md](specialization.md).

## 8. Passive/Active Skill Grant

Spec seçildiğinde:
- `spec.passive_skill_codes` → her biri `character.skill_learn` (rank=1).
- `spec.active_skill_unlock_codes` → aynı; ayrıca skill_bar slot suggest (opsiyonel).
- `spec.signature_skill_code` (Tier 2 only) → skill_bar slot 10 (reserved signature).

## 9. Respec

Detay: [respec.md](respec.md).

- İlk respec **ücretsiz** (öğrenme).
- Sonrası: `500 * pow(2, respec_count)` gold, cap 20000; cooldown 24h.
- Free respec token: content patch (rework), quest reward, event.

## 10. Balance Budget

Spec bias toplam bütçesi:
- **Tier 1**: +3 stat toplam.
- **Tier 2**: +2 stat toplam.
- Ayrıca: 1-2 passive + 1-2 active unlock + Tier 2 için 1 signature skill.

Class base stat bias + spec bias toplamı **+7-8 stat maks** (net); race ile beraber +10-12 (bkz. [race/base-stats.md](../race/base-stats.md) §5).

## 11. Class Rework Compensation

Spec/class rework durumu:
- Etkilenen karakterlere **free respec token** (30 gün TTL).
- Version bump + audit event.
- Notification.

Spec kaldırma:
- Deprecation: mevcut karakter için legacy olarak korunur veya en yakın Spec II'ye migrate (oyuncu onayı) — bkz. DD-004.

## 12. Cross-Reference

- [game-design/classes.md](../game-design/classes.md), [game-design/specializations.md](../game-design/specializations.md)
- [character/character.md](../character/character.md), [character/attributes.md](../character/attributes.md)
- [race/base-stats.md](../race/base-stats.md) — combined budget
- [database/content-templates.md §2-3](../database/content-templates.md)
- [game-design/skills.md](../game-design/skills.md), [game-design/combat.md](../game-design/combat.md)

## 13. Bekleyen / Revize Konuları

- **Placeholder isimler final değil** — [placeholder-content.md](placeholder-content.md) content oyuncu deneyimi + IP çakışma testi sonrası revize.
- Multi-class (dual class Phase 03) desteklenmiyor.
- Tier 3 spec (Phase 02+) planlı; schema `tier INT` alanı hazır.
- Talent point (Phase 02) spec ağacında düğüm bazlı puan.
