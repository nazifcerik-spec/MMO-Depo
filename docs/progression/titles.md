# Titles (Extended)

> [character/titles.md](../character/titles.md) v1 temel; bu doküman genişletir.

## 1. Amaç
Karakterin ismi yanında görsel prestij marker'ı. Achievement, quest, event, mastery, ascension, reputation sources.

## 2. Karakteristik

- **Prefix veya suffix** yapı (v1'de sadece **suffix**: `Grim, the Ironmaster`).
- **Scope**: character-bound (v1); account-bound titles opt Phase 02.
- **Concurrent unlock cap**: sınırsız (koleksiyon amaçlı) ama aktif tek title.
- **Rarity**: `common | uncommon | rare | epic | legendary | mythic` (görsel renk).
- **Retire**: legendary+ titles bazıları "retired" (event limited); yeni oyuncu unlock edemez.

## 3. Model

```sql
CREATE TABLE content.title_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_key VARCHAR(64) NOT NULL,          -- "the Ironmaster"
  display_key_female VARCHAR(64) NULL,       -- lokalizasyon için
  placement VARCHAR(16) NOT NULL DEFAULT 'suffix',
  rarity VARCHAR(16) NOT NULL,
  source_kind VARCHAR(24) NOT NULL,          -- achievement|quest|event|mastery|ascension|reputation|donation
  is_retired BOOLEAN NOT NULL DEFAULT false,
  is_hidden_until_unlock BOOLEAN NOT NULL DEFAULT false,
  version INT DEFAULT 1
);

CREATE TABLE progression.character_title (
  character_id UUID NOT NULL,
  title_id UUID NOT NULL,
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  source_ref_kind VARCHAR(24) NOT NULL,
  source_ref_id UUID NULL,
  is_favorite BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (character_id, title_id)
);

ALTER TABLE character.character
  ADD COLUMN active_title_id UUID NULL;
```

## 4. Kaynaklar (v1 örnek katalog)

| Source | Örnek | Rarity |
|---|---|---|
| Achievement | `the_ironmaster` (mining tier 10) | rare |
| Quest | `path_of_ascension_complete` → `the Ascended` | epic |
| Event (seasonal) | `winter_champion` | epic |
| Mastery | `master_gatherer` (all gathering tier 10) | epic |
| Ascension | `ascended_i`, `ascended_v`, `eternal` (10+) | legendary/mythic |
| Reputation | `friend_of_faction_X` (rep tier max) | uncommon |
| Donation/support | `founder`, `alpha_tester` | retired legendary |
| Guild | `guildmaster_of_X` | rare (dinamik display) |

## 5. Aktif Title Seçim

`PUT /v1/character/{id}/title/active`
```json
{ "title_code": "the_ironmaster" }
```
- Karakter title'ı unlock etmiş olmalı (`character_title` join).
- `character.active_title_id` update.
- Aktif title kaldırma: body `{ title_code: null }`.

Chat/UI'da display: `{character_name} {, title_display_key}` (title varsa).

## 6. Retirement Politikası

- Alpha/beta test titles: launch sonrası retire (yeni unlock kapalı).
- Seasonal event unique: event bitince retire (nostalji koruma).
- Retired title `is_retired=true`; unlock edenler yine kullanabilir.

## 7. Fav / Filter

Karakter title galeri UI:
- Favorite flag ile pin.
- Category filter (source_kind).
- Rarity filter.
- Search by name.

## 8. Endpoint

**GET `/v1/character/{id}/titles?rarity=&source=&favorite=`**
**PUT `/v1/character/{id}/title/active`**
**PUT `/v1/character/{id}/title/{code}/favorite`** body: `{ is_favorite: true }`

## 9. Metrikler
- `title_unlock_total{source, rarity}`
- `title_active_change_total`
- `title_retired_gauge`

## 10. Cross-Reference
- [character/titles.md](../character/titles.md) — v1 temel şema (bu doküman extend eder)
- [achievements.md](achievements.md), [mastery.md](mastery.md), [ascension.md](ascension.md), [reputation.md](reputation.md)
- [collections.md](collections.md)

## 11. Bekleyen Kararlar
- Prefix title (Phase 02'de "Champion Grim").
- Account-wide title (Phase 02).
- Title dinamik display (guildmaster_of_{guild_name}) rendering Phase 02.
- Retirement policy (yıllık review) ops runbook Phase 02.
