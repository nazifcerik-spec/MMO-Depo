# Race (Entity & Lifecycle)

## 1. Amaç
Race entity'sinin kalıcı tanımı; oyuncu karakterine bağlanma flow'u; content pack üzerinden ekleme/güncelleme/eskitme akışı.

## 2. Model
Kaynak: [database/content-templates.md §1](../database/content-templates.md).

Ana alanlar:
- `id`, `code` (unique).
- `display_name_key`, `lore_key`, `art_ref`.
- `stat_modifiers` (JSONB).
- `passive_ability_id`, `active_ability_id` (FK → `content.skill_template`).
- `profession_affinity_id`, `profession_affinity_bonus_pct`.
- `unique_mechanic_code`.
- `version`, `deprecated_at` (soft-deprecate).

Ek alanlar (Phase 7 eklentileri):
```sql
ALTER TABLE content.race
  ADD COLUMN allowed_cosmetics_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN starter_location_city_id UUID NULL REFERENCES world.settlement(id),
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}',
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL,
  ADD COLUMN unique_mechanic_params_json JSONB NOT NULL DEFAULT '{}';
```

## 3. Race Life-Cycle

```
DRAFT (yaml only, not deployed)
  ↓ content_build validate
STAGED (staging content_version)
  ↓ promote
LIVE (production content_version bump)
  ↓ (opsiyonel, gelecekte)
DEPRECATED (deprecated_at set; character creation dışı, mevcut kararaktere etkisi yok)
  ↓ (nadir)
ARCHIVED (unreachable; mevcut karakterler zaten snapshot ile korunur)
```

Character creation'da sadece `deprecated_at IS NULL` race listelenir.

## 4. Character Creation Integration

Bkz. [character/character.md §4](../character/character.md).

Race seçimi çakışması:
- Client `GET /v1/content/races` → aktif race listesi.
- Client body → `race_id`.
- Server validate: race var + deprecated değil + allowed_cosmetics uygun.
- Insert character → race_id set.
- Compute base stats:
  ```
  primary[i] = 10 + race.stat_modifiers[i] + class.base_stat_bias[i]
  ```
- Grant racial abilities:
  - Insert `character.skill_learn (character_id, skill_id=race.passive_ability_id, rank=1)`.
  - Insert `character.skill_learn (character_id, skill_id=race.active_ability_id, rank=1)`.
  - Insert `character.skill_bar (character_id, bar_index=11, skill_id=race.active_ability_id)` — son slot reserve racial active.
- Set `character.home_city_id = race.starter_location_city_id`.

## 5. Content Endpoint (Race Discovery)

**GET `/v1/content/races`**
Public (no auth) — karakter creation ekranı için.

Response:
```json
{
  "content_version": "2026.09.01-01",
  "races": [
    {
      "id": "r-...",
      "code": "human",
      "display_name_key": "race.human.name",
      "lore_key": "race.human.lore",
      "art_ref": "races/human/portrait.png",
      "stat_modifiers": { "str": 1, "dex": 1, ... },
      "passive_summary_key": "race.human.passive.summary",
      "active_summary_key": "race.human.active.summary",
      "profession_affinity": {
        "profession_code": "trading",
        "bonus_pct": 5
      },
      "unique_mechanic_key": "race.human.unique.summary",
      "tags": ["balanced", "versatile"]
    },
    ...
  ]
}
```

Cache: content_version key; ETag. UI localizes with i18n keys.

## 6. Race Change (Phase 02+)

**POST `/v1/characters/{id}/race-change`**
Auth + CharacterScope.

Request:
```json
{ "new_race_id": "r-elf" }
```

Server:
1. Feature flag check: `race_change_enabled`.
2. Cost: 5000 gold **veya** rename token (Phase 02 monetization).
3. Cooldown: 90 gün per character.
4. Transaction:
   - Debit currency.
   - Old racial ability'leri character.skill_learn'dan kaldır (rank sıfırlama uyarısı).
   - New racial'ları öğret.
   - Update `character.race_id`.
   - Recompute base stats (level-up puanları etkilenmez; sadece racial baz).
   - `character.stats.derived_cache_invalidated_at = now()`.
   - Publish event `character.race_changed`.
   - Audit event.
5. Response 200.

v1'de **kapalı**; flag `race_change_enabled=false`.

## 7. Race Nerf / Rework Compensation

Rework durumu:
- Server-side detection: `content.race.version` bump + mevcut karakter sayısı N > 0.
- Worker `race_rework_notify`:
  - Etkilenen karakterlere in-game mail + kompanzasyon token.
  - Token: 1x free race change (30 gün TTL).
- Audit: `content.race_reworked`.

## 8. Deprecation

**POST `/v1/admin/content/races/{id}/deprecate`** (admin route)

Server:
1. `deprecated_at = now()`.
2. Publish event `content.race_deprecated`.
3. Character creation UI'da gizle.
4. Mevcut karakter: hiç değişmez.

## 9. Fail Modları
- Race not found / deprecated → 404 / 400 `race.unavailable`.
- Cosmetic not allowed for race → 400 `race.cosmetic_invalid`.
- Race change cooldown → 429.
- Insufficient funds → 402.

## 10. Metrics
- `character_race_distribution{race}` gauge (character sayısı per race).
- `race_change_total{from,to}` counter.
- `race_content_reload_total`.

## 11. Audit
- `character.race_selected` (creation'da).
- `character.race_changed` (Phase 02).
- `content.race_created`, `content.race_updated`, `content.race_deprecated`.

## 12. Event
```json
{
  "_v": 1,
  "character_id": "c-...",
  "old_race_id": "r-human",
  "new_race_id": "r-elf",
  "reason": "player_change"
}
```

## 13. Test Stratejisi
- Content loader: valid yaml → row insert.
- Invalid yaml (missing ability code) → build fail.
- Character creation: race constraint enforcement.
- Race change transaction: partial fail rollback.
- Deprecation: creation reject; existing char unaffected.

## 14. Cross-Platform
Race data client'a cache-friendly (i18n resolver client-side).

## 15. İlişkiler
- [Base Stats](base-stats.md), [Passive](passive.md), [Racial Skills](racial-skills.md), [Racial Bonuses](racial-bonuses.md)
- [character/character.md](../character/character.md), [character/attributes.md](../character/attributes.md)
- [database/content-templates.md](../database/content-templates.md)
- [game-design/races.md](../game-design/races.md)

## 16. Bekleyen Kararlar
- Race change Phase 02 aktivasyonu (feature flag).
- Race compendium (koleksiyon achievement) Phase 03.
- Sub-race varyantları modellemesi (`race_variant` table Phase 02).
- Home city override — race default'u vs oyuncu tercih.
