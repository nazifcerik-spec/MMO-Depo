# Specialization (Tier I & II Tree)

## 1. Amaç
Base class'ı iki katmanlı seçim ile derinleştirmek. Aynı class'ın farklı oynanış yönlerini kanıtlamak. Data-driven; yeni spec eklemek kod değişikliği gerektirmez.

## 2. Model
Kaynak: [database/content-templates.md §3](../database/content-templates.md).

`content.specialization` alanları (özet):
- `id`, `code` (unique), `class_id` FK, `tier` (1 or 2), `parent_spec_id` (tier 2 için tier 1 ref).
- `display_name_key`, `lore_key`, `art_ref`.
- `passive_ids`, `active_skill_unlock_ids` (arrays).
- `stat_bias` JSONB.
- `playstyle_tags` array.

Ek alanlar (Phase 8):
```sql
ALTER TABLE content.specialization
  ADD COLUMN signature_skill_id UUID NULL REFERENCES content.skill_template(id),
  ADD COLUMN tier2_option_ids UUID[] NOT NULL DEFAULT '{}',    -- tier 1 için tier 2 seçenek listesi
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL,
  ADD COLUMN description_short_key VARCHAR(96) NULL,
  ADD COLUMN badge_ref VARCHAR(128) NULL;
```

## 3. Tree Constraint
```
Tier 1: parent_spec_id IS NULL
Tier 2: parent_spec_id IS NOT NULL AND parent.tier=1
```

Constraint (mevcut database/content-templates.md §3'te):
```sql
CHECK (tier = 1 OR parent_spec_id IS NOT NULL)
```

v1 kesin sayı:
- Her Tier 1: `array_length(tier2_option_ids) = 3`.
- Her class: `array_length(spec_tier1_ids) = 2`.

Loader validation zorlar.

## 4. Character Model Bağlantısı
`character.character.spec_tier1_id`, `character.character.spec_tier2_id` (nullable; level ile açılır).

Constraint (mevcut):
```sql
CHECK (spec_tier2_id IS NULL OR spec_tier1_id IS NOT NULL)
```

## 5. Unlock Timing

### 5.1 Tier I (Level 5)
Character level 5'te:
- `character.leveled_up` event → subscriber `specialization_unlock_check`.
- UI'da "Choose your path" pop-up.
- Oyuncu 2 seçenekten seçer.

### 5.2 Tier II (Level 10)
Character level 10'da:
- Aynı pattern; Tier I'e bağlı 3 seçenek gösterilir.

## 6. Endpoint

### 6.1 Spec Options

**GET `/v1/characters/{id}/specializations/options`**
Auth + CharacterScope.

Response:
```json
{
  "tier1": {
    "unlocked": true,
    "current_id": null,
    "options": [
      {
        "id": "sp-berserker",
        "code": "berserker",
        "display_name_key": "spec.berserker.name",
        "description_short_key": "spec.berserker.short",
        "passive_summary_key": "spec.berserker.passive.sum",
        "active_summary_key": "spec.berserker.active.sum",
        "stat_bias": { "str": 3 },
        "playstyle_tags": ["offense","sustain_dmg"]
      },
      {
        "id": "sp-guardian",
        "code": "guardian",
        ...
      }
    ]
  },
  "tier2": {
    "unlocked": false,      // level 10'a ulaşınca true
    "current_id": null,
    "options": []
  }
}
```

Character level >= 10 && spec_tier1_id set ise `tier2.options` doldurulur (parent'ın 3 çocuğu).

### 6.2 Select Tier I

**POST `/v1/characters/{id}/specializations/tier1`**
Auth + CharacterScope.

Request:
```json
{ "spec_id": "sp-berserker" }
```

Server:
1. Load character FOR UPDATE.
2. Validate:
   - `character.level >= 5`.
   - `spec_tier1_id IS NULL` (ilk seçim; respec ayrı endpoint).
   - `spec.class_id = character.class_id`.
   - `spec.tier = 1`.
   - `spec.deprecated_at IS NULL`.
3. Transaction:
   - `character.spec_tier1_id = spec_id`, `version += 1`.
   - Grant passive skills (`character.skill_learn`).
   - Grant active skill unlocks.
   - Apply stat bias (recompute derived cache).
   - Emit event `character.spec_selected` (tier=1).
   - Audit event.
4. Response 200: updated character.

### 6.3 Select Tier II

**POST `/v1/characters/{id}/specializations/tier2`**
Auth + CharacterScope.

Request:
```json
{ "spec_id": "sp-ravager" }
```

Server:
1. Load character.
2. Validate:
   - `character.level >= 10`.
   - `character.spec_tier1_id IS NOT NULL`.
   - `character.spec_tier2_id IS NULL` (first pick).
   - `spec.parent_spec_id = character.spec_tier1_id`.
   - `spec.tier = 2`.
3. Transaction:
   - `character.spec_tier2_id = spec_id`.
   - Grant passive, active, signature skill.
   - Signature skill_bar slot 10 (reserved; oyuncu değiştirebilir).
   - Stat bias apply.
   - Event + audit.
4. Response 200.

## 7. Skill Grant Detay

Spec seçildiğinde:
```
for each passive_id in spec.passive_skill_ids:
    upsert character.skill_learn (character_id, skill_id, rank=1)
for each active_id in spec.active_skill_unlock_ids:
    upsert character.skill_learn (rank=1)
    // opsiyonel: character.skill_bar boş slot varsa suggest
if spec.tier = 2 and spec.signature_skill_id:
    upsert character.skill_learn (signature, rank=1)
    upsert character.skill_bar (slot=10, skill=signature)
```

Zaten öğrenilmiş skill'ler (starter class veya önceki spec) rank korur.

## 8. Spec Change (Respec)

Detay: [respec.md](respec.md).

Kısaca:
- `POST /v1/characters/{id}/respec-spec` — Tier I ve/veya II değişimi.
- İlk free; sonra gold cost + 24h cooldown.
- Eski spec skill'leri **kaldırılır** (`character.skill_learn` delete).
- Skill bar temizlenir + spec varsayılan bar rebuild.
- Yeni spec skill'leri öğretilir.

## 9. Playstyle Tags
Standart tag seti (whitelist):
- `offense`, `defense`, `sustain_dmg`, `burst`, `execute`, `dot`, `aoe`, `single_target`
- `tank_main`, `tank_off`, `support`, `heal`
- `mobility`, `crowd_control`, `positional`, `stationary`
- `caster`, `melee`, `ranged`, `hybrid_melee_ranged`
- `dual_wield`, `two_handed`, `shield`
- `low_hp_bonus`, `high_hp_bonus`, `pet_focus`

## 10. Content Loader Validation

- `class_code` valid.
- Tier=1: `parent_spec_code=null`.
- Tier=2: `parent_spec_code` valid tier=1 spec.
- Tier=1 için `tier2_options` tam 3 element.
- Tier=2 için signature_skill_code opsiyonel ama önerilir.
- Skill code'ları var.
- Stat bias budget:
  - Tier 1: `sum(abs(stat_bias)) <= 3`.
  - Tier 2: `sum(abs(stat_bias)) <= 2`.
- Unique code.

## 11. Spec Rework / Deprecation

Spec rework:
1. YAML patch.
2. Version bump.
3. Etkilenen karakterlere free respec token (Tier I için Tier I respec token, Tier II için Tier II respec token).
4. Notification + audit.

Spec deprecation:
- `deprecated_at = now()`.
- Character creation UI'da yeni spec seçilemez.
- Mevcut karakterler için: 30 gün free respec pencere; sonra otomatik "closest surviving spec"e migrate (oyuncu onayı ile).

## 12. Metrics
- `character_spec_distribution{class,spec_tier1,spec_tier2}` gauge.
- `spec_selected_total{spec}` counter.
- `spec_respec_total{tier,from,to}` counter.
- `spec_meta_impact_pct{spec}` — weekly meta health check.

## 13. Audit
- `character.spec_selected` (tier=1|2).
- `character.spec_respec` (bkz. respec.md).
- `content.spec_updated`, `content.spec_deprecated`.

## 14. Event Payload
```json
{
  "_v": 1,
  "character_id": "c-...",
  "spec_id": "sp-...",
  "spec_code": "berserker",
  "tier": 1,
  "class_code": "warrior"
}
```

## 15. Fail Modları
- Level insufficient → 400 `spec.level_required`.
- Already picked → 409 `spec.already_selected` (respec kullan).
- Wrong class → 400 `spec.wrong_class`.
- Parent mismatch (tier 2) → 400 `spec.parent_mismatch`.
- Deprecated → 400 `spec.deprecated`.

## 16. Rate Limit
- Tier select: 5/hour per character (spam koruma; genelde 1 kez).
- Respec ayrı endpoint kendi limit'iyle.

## 17. Test Stratejisi
- Tree validation: her class'ın 2 Tier I × 3 Tier II var.
- Select Tier I → grant skills + bias.
- Select Tier II → grant signature.
- Wrong parent Tier II → reject.
- Level gate enforcement.

## 18. Cross-Platform
Standart character sheet; her platform aynı.

## 19. İlişkiler
- [Class](class.md), [Respec](respec.md), [Placeholder Content](placeholder-content.md)
- [character/character.md](../character/character.md), [character/level.md](../character/level.md), [character/attributes.md](../character/attributes.md)
- [game-design/specializations.md](../game-design/specializations.md), [game-design/skills.md](../game-design/skills.md)
- [database/content-templates.md](../database/content-templates.md)

## 20. Bekleyen Kararlar
- Tier III (Phase 02+) — hazır schema, henüz yok.
- Talent point per-tier düğüm (Phase 02).
- Cross-class hybrid (dual-spec) — Phase 03+ değerlendirme.
- Spec preview (deneme modu, geri dönüşsüz olmadan) — Phase 02 opsiyonel.
