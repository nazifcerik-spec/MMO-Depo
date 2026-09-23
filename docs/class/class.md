# Class (Base Class Entity)

## 1. Amaç
Karakterin oynanış kimliğinin kök seçimi: silah/zırh yatkınlığı, kaynak türü, rol yatkınlığı, ana stat bias, starter skill seti.

## 2. Model
Kaynak: [database/content-templates.md §2](../database/content-templates.md).

Ana alanlar (özet):
- `id`, `code` (unique).
- `display_name_key`, `lore_key`, `art_ref`.
- `resource_type` (rage/mana/energy/focus/nature).
- `armor_affinity` (cloth/leather/mail/plate).
- `weapon_affinities` (array).
- `role_tags` (tank/dps/support/hybrid).
- `base_stat_bias` JSONB.
- `starter_skill_ids` (array).
- `spec_tier1_ids` (array of Spec I).
- `version`, `deprecated_at`.

Ek alanlar (Phase 8):
```sql
ALTER TABLE content.class
  ADD COLUMN starter_home_city_id UUID NULL REFERENCES world.settlement(id),
  ADD COLUMN starter_equipment_pack_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL,
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}';
```

## 3. Class Life-Cycle

```
DRAFT (yaml, not deployed)
  ↓ content_build validate
STAGED
  ↓ promote
LIVE
  ↓ (rare)
DEPRECATED (deprecated_at set; creation'da gizli; mevcut karakter korunur)
```

Character creation'da sadece `deprecated_at IS NULL` class listelenir.

## 4. Character Creation Integration

Bkz. [character/character.md §4](../character/character.md).

Class seçimi validation:
- `class_id` exists + not deprecated.
- Race → class kilit yok (DD-005).
- Cosmetic uyum kontrolü opsiyonel (Phase 02).

Insert path:
```
Character.class_id = <selected>
Character.stats.base = 10 + race.stat_modifiers + class.base_stat_bias
character.skill_learn ← class.starter_skill_ids (rank=1)
character.skill_bar ← starter default (slot 0..N)
character.equipment ← class.starter_equipment_pack_json → item_instance(s)
character.rotation ← default rotation for class (spec sonrası revize edilir)
```

## 5. Starter Equipment Pack

`starter_equipment_pack_json` format:
```json
{
  "_v": 1,
  "items": [
    { "template_code": "training_sword", "slot": "main_hand", "count": 1 },
    { "template_code": "training_shield", "slot": "off_hand", "count": 1 },
    { "template_code": "cloth_shirt", "slot": "chest", "count": 1 },
    { "template_code": "leather_pants", "slot": "legs", "count": 1 },
    { "template_code": "leather_boots", "slot": "feet", "count": 1 },
    { "template_code": "minor_health_potion", "slot": "inventory", "count": 5 },
    { "template_code": "bread_loaf", "slot": "inventory", "count": 3 }
  ]
}
```

Content loader item_template varlığını doğrular; yoksa build fail.

## 6. Class Change

**v1'de KAPALI.** DD-004: class kalıcı. Spec değişir (bkz. [respec.md](respec.md)).

Phase 03+ tartışılabilir (feature flag `class_change_enabled`); şu an kararı yok.

## 7. Endpoint

### 7.1 Class Discovery

**GET `/v1/content/classes`**
Public.

Response:
```json
{
  "content_version": "2026.09.08-01",
  "classes": [
    {
      "id": "cls-warrior",
      "code": "warrior",
      "display_name_key": "class.warrior.name",
      "lore_key": "class.warrior.lore",
      "art_ref": "classes/warrior/portrait.png",
      "resource_type": "rage",
      "armor_affinity": "plate",
      "weapon_affinities": ["sword", "axe", "mace", "shield"],
      "role_tags": ["tank", "dps_melee"],
      "base_stat_bias": { "str": 3, "sta": 2 },
      "specializations_tier1": [
        { "code": "berserker", "display_name_key": "spec.berserker.name" },
        { "code": "guardian",  "display_name_key": "spec.guardian.name" }
      ],
      "tags": ["melee", "frontline"]
    }
  ]
}
```

### 7.2 Class Detail

**GET `/v1/content/classes/{code}`**
Public. Genişletilmiş detay: starter skills, spec tree preview.

## 8. Class Rework Compensation

Balance rework:
1. Content patch draft (yaml).
2. Simülasyon: mevcut karakterlerin etkilenen build'leri.
3. Impact > threshold ise: **free respec token** (spec) — class kalıcı olduğundan respec spec seviyesinde.
4. Rollout: content_version bump + notification.
5. Extreme rework (class fundamentally changed): 30 gün free respec pencere + duyuru.

## 9. Deprecation

Class kaldırma (nadir, tercih edilmez):
- `deprecated_at = now()`.
- Character creation UI'da gizli.
- Mevcut karakter etkilenmez (legacy).
- 30 gün free respec + duyuru; oyuncu isterse yeni karakter seçebilir (compensation slot?).

## 10. Fail Modları
- Class not found / deprecated → 400 `class.unavailable`.
- Insufficient starter items (content missing) → 500 (build should have prevented; audit critical).
- Skill code invalid → build fail.

## 11. Metrics
- `character_class_distribution{class}` gauge.
- `content_class_reload_total`.
- Character creation source distribution per class.

## 12. Audit
- `character.class_selected` (creation).
- `content.class_updated`, `content.class_deprecated`.

## 13. Test Stratejisi
- Content loader: valid yaml → row insert; invalid → fail.
- Character creation: race + class combinations (6×6 = 36 case).
- Starter equipment integrity: her item_template var, equip requirements karşılanır.

## 14. Cross-Platform
Class content her platformda aynı (server single source).

## 15. Balance Constraint

- **Base stat bias toplam +5 stat bütçesi** (yumuşak sınır; ±1 tolerate).
- **Weapon affinity minimum 3** (versatility).
- **Starter equipment tier**: her class için eşdeğer basit set (level 1 baseline).
- **Role tag**: max 2 (single-focus veya hybrid; 3+ tag hedef kaybı).

## 16. İlişkiler
- [Specialization](specialization.md), [Respec](respec.md), [Placeholder Content](placeholder-content.md)
- [character/character.md](../character/character.md), [character/attributes.md](../character/attributes.md)
- [race/base-stats.md](../race/base-stats.md) — combined budget
- [database/content-templates.md](../database/content-templates.md)
- [game-design/classes.md](../game-design/classes.md)

## 17. Bekleyen Kararlar
- Class change Phase 03+ karar (v1 kesin hayır).
- Dual class hybrid Phase 03+ (currently unsupported).
- Class-locked cosmetic (Phase 02).
- Class-race preferred pairing UI hint (öneri, kilitleme yok).
