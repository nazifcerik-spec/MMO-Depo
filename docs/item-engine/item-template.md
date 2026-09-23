# Item Template

## 1. Amaç
Item'ın kalıcı tanımı — art, base stats, affix pool referansları, requirements, binding type. Instance yaratımının şablonu.

Kaynak: [database/content-templates.md §8](../database/content-templates.md).

## 2. Model
Ana alanlar:
- `id`, `code` (unique).
- `display_name_key`, `description_key`, `art_ref`.
- `category` (equipment/consumable/material/quest/currency/cosmetic/utility/recipe).
- `subcategory` (weapon/armor/accessory + weapon_kind).
- `rarity` (poor/common/uncommon/rare/epic/legendary/mythic).
- `base_stats` JSONB (min-max ranges için).
- `affix_pool_ids` (prefix + suffix pool refs).
- `socket_slots` (rarity-driven; template'te opsiyonel override).
- `set_id` (set üyesi ise).
- `binding_type` (no_bind/boe/bop/boa).
- `requirements` (min_level, class_ids, spec_ids, race_ids, stat_mins).
- `base_durability` (nullable; consumables için null).
- `stack_max`, `weight`.
- `vendor_buy_price_copper`, `vendor_sell_price_copper`.
- `unique_effect_id` (legendary+ için).
- `equipment_slot`, `weapon_kind`, `armor_kind`.
- `version`, `deprecated_at`.

## 3. Content Pack Format (YAML)

```yaml
# content/items/steel_longsword.yaml
code: steel_longsword
version: 1

display_name_key: item.steel_longsword.name
description_key: item.steel_longsword.desc
art_ref: items/weapons/steel_longsword.png

category: equipment
subcategory: weapon
equipment_slot: main_hand
weapon_kind: sword
armor_kind: null

rarity: uncommon

base_stats:
  weapon_damage:
    min: 45
    max: 60
  attack_speed_ms: 2400

base_durability: 100
stack_max: 1
weight: 10

vendor_buy_price_copper: 12000
vendor_sell_price_copper: 3000

binding_type: boe

affix_pool_refs:
  prefix: uncommon_weapon_prefix_pool_v1
  suffix: uncommon_weapon_suffix_pool_v1

socket_slots: 0     # rarity default kullanılabilir; explicit override

set_id: null
unique_effect_id: null

requirements:
  min_level: 12
  class_codes: [warrior, ranger, rogue]
  spec_codes: null
  race_codes: null
  stat_mins:
    str: 15

tags: [melee, one_hand]
```

## 4. Rarity-Driven Defaults

Bazı alanlar template'te belirtilmezse rarity'den defaults gelir (bkz. [rarity.md](rarity.md)):
```
uncommon: affix_count=1-2, socket_slots=0
rare:     affix_count=2-3, socket_slots=0
epic:     affix_count=3-4, socket_slots=1
legendary: affix_count=4, socket_slots=1-2, unique_effect required
mythic:   affix_count=4, socket_slots=2, unique_effect + progression gated
```

Template explicit override rarity default'unu ezer.

## 5. Requirement Check

Equip endpoint'inde ve loot hot-path'te:
```go
func canEquip(char Character, tmpl TemplateDef) error {
    r := tmpl.Requirements
    if char.Level < r.MinLevel { return errs.ErrForbidden }
    if len(r.ClassCodes) > 0 && !slices.Contains(r.ClassCodes, char.ClassCode) { return errs.ErrForbidden }
    if len(r.SpecCodes) > 0 && !specIn(char, r.SpecCodes) { return errs.ErrForbidden }
    if len(r.RaceCodes) > 0 && !slices.Contains(r.RaceCodes, char.RaceCode) { return errs.ErrForbidden }
    for stat, min := range r.StatMins {
        if char.Derived.Primary[stat] < min { return errs.ErrForbidden }
    }
    return nil
}
```

## 6. Template Registry

```go
type TemplateRegistry interface {
    ByID(ctx context.Context, id uuid.UUID) (TemplateDef, error)
    ByCode(ctx context.Context, code string) (TemplateDef, error)
    ForCategory(ctx context.Context, cat Category) ([]TemplateDef, error)
    ForEquipmentSlot(ctx context.Context, slot EquipmentSlot) ([]TemplateDef, error)
    ReloadOnVersionBump(ctx context.Context) error
}
```

In-memory cache; content_version bump → reload.

## 7. Template Snapshot Hash

Template değişse bile instance state'i korur:
- Instance yaratımında `template_snapshot_hash = sha256(template_content_serialized)`.
- Instance okumada template'in hash'i mevcut hash'le uyumsuzsa "template drift" flag.
- Drift durumu: instance rolled_stats + affixes ile çalışır; yeni template'in mekaniğine göre re-render (yalnız cosmetic/UX).

Ekonomi kritik alanlar (rarity, binding_type) instance'ta snapshot değil — güncel template kullanır. Değişiklik oyuncu için görünür (item değeri artabilir/azalabilir) — patch note zorunluluğu.

## 8. Template Versioning

Content patch:
1. YAML draft.
2. Validation:
   - Zorunlu alanlar.
   - Rarity ↔ affix pool tutarlılığı.
   - Affix pool refs var.
   - Set ref valid (set upsert edildiyse).
   - Unique effect handler kayıtlı.
   - Vendor price budget check.
3. Content build hash generate.
4. Upsert `content.item_template`.
5. `version` bump.
6. Publish event `content.item_template_updated`.

Deprecation:
- `deprecated_at = now()`.
- Yeni loot roll'da drop yok.
- Mevcut instance kullanılmaya devam.
- Marketplace listing OK.

## 9. Content Loader Validation

CI check:
- `code` unique.
- `rarity` enum'da.
- `category` valid + subcategory pattern.
- `equipment_slot` (equipment ise required).
- `weapon_kind` (weapon ise required).
- `affix_pool_refs` var (rarity gerektiriyorsa).
- `binding_type` enum.
- `requirements.min_level` 1..cap.
- Vendor price non-negative.
- Set ref valid.
- Unique effect var (rarity >= legendary ise required).

## 10. Endpoint (Public Content)

**GET `/v1/content/items/{code}`**
Public — item tooltip için.
```json
{
  "id": "...",
  "code": "steel_longsword",
  "display_name_key": "item.steel_longsword.name",
  "art_ref": "...",
  "category": "equipment",
  "rarity": "uncommon",
  "base_stats_summary": { "weapon_damage": "45-60", "attack_speed_ms": 2400 },
  "requirements": { "min_level": 12, "class_codes": [...] },
  "binding_type": "boe",
  "affix_pool_summary": { "possible_prefixes_count": 12, "possible_suffixes_count": 15 }
}
```

## 11. Metrics
- `item_template_count` gauge.
- `item_template_reload_total`.
- `item_template_deprecated_total`.

## 12. Audit
- `content.item_template_created`, `updated`, `deprecated`.

## 13. Test Stratejisi
- Loader validation: valid → row insert; invalid → build fail.
- Requirement matrix.
- Snapshot hash generation deterministic.
- Registry lookup performance (10k templates).

## 14. Cross-References
- [Item Instance](item-instance.md), [Rarity](rarity.md), [Affix](affix.md), [Set](set.md), [Unique Effect](unique-effect.md)
- [database/content-templates.md](../database/content-templates.md)
- [game-design/items.md](../game-design/items.md)

## 15. Bekleyen Kararlar
- Snapshot alanları listesi final (rolled_stats snapshot; rarity live).
- Template i18n resolver server-side mi client-side (öneri client).
- Craft output template ile loot template ayrı tutulsun mu (öneri: aynı sistem).
