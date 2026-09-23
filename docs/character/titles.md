# Titles

## 1. Amaç
Karakterin isminin önünde/arkasında görünen kozmetik ünvan. Kimlik + prestij + koleksiyon.

## 2. Model
Kaynak: [database/character.md §17](../database/character.md), [database/content-templates.md §25](../database/content-templates.md), [database/progression-quests-achievements.md §4](../database/progression-quests-achievements.md).

Tablolar:
- `content.title` — title tanımı (kod, isim key, kaynak).
- `character.title_owned` — karakter × title × earned_at × is_active.
- `progression.title_earned` — earn history (log).

## 3. Title Türleri

### 3.1 Prefix (İsim önü)
- `"Grandmaster <Name>"` (profession mastery).
- `"The Ascendant <Name>"` (ascension milestone).

### 3.2 Suffix (İsim arkası)
- `"<Name>, Ashen Wyrm Slayer"` (world boss achievement).
- `"<Name> of the Emerald Wardens"` (faction exalted).

### 3.3 Server-First
- Sadece 1 karaktere verilir (kilitli slot).
- `content.title.source = 'server_first'` + `unique_holder_character_id`.

### 3.4 Seasonal
- Season'da earn edilir; season sonrası archive olur ancak sahip olan gösterebilir.

## 4. Kaynak (Earn Sources)

| Kaynak | Örnek Title |
|---|---|
| Achievement | "Adventurer" (level 25 first) |
| Quest turn-in | "Hero of Ashfield" |
| Event participation | "Harvest Champion 2026" |
| World boss first kill | "Ashen Wyrm Slayer" |
| Reputation exalted | "of the Emerald Wardens" |
| Ascension level | "The Ascendant" (asc 1), "The Twice-Born" (asc 2), ... |
| Profession grandmaster | "Grandmaster Miner" |
| Purchase (Phase 02) | Cosmetic shop |
| Admin grant | "Beta Tester" |

## 5. Content Şema Ekleme (v1)

```sql
CREATE TABLE content.title (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,      -- "title.grandmaster_miner"
  format VARCHAR(96) NOT NULL,                -- "Grandmaster {name}" veya "{name}, Ashen Wyrm Slayer"
  placement VARCHAR(16) NOT NULL,             -- 'prefix' | 'suffix'
  source VARCHAR(32) NOT NULL,                -- 'achievement' | 'quest' | 'event' | 'faction' | 'ascension' | 'profession' | 'purchase' | 'admin'
  source_ref_id UUID NULL,
  is_server_first BOOLEAN NOT NULL DEFAULT false,
  unique_holder_character_id UUID NULL,        -- server_first ise sahiplenildikten sonra kilit
  season_key VARCHAR(32) NULL,
  hidden_until_earned BOOLEAN NOT NULL DEFAULT false,
  version INT NOT NULL DEFAULT 1
);
```

`character.title_owned` zaten mevcut:
```sql
CREATE TABLE character.title_owned (
  character_id UUID NOT NULL,
  title_id UUID NOT NULL,
  earned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  is_active BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (character_id, title_id)
);
```

Constraint: bir karakterin tek `is_active = true` title'ı (partial unique):
```sql
CREATE UNIQUE INDEX uq_title_owned__active
ON character.title_owned (character_id)
WHERE is_active = true;
```

## 6. Earn Flow
Domain call:
```go
func (s *Service) GrantTitle(ctx context.Context, charID uuid.UUID, titleCode string, source string) error
```

Server:
1. Load title by code.
2. Server-first check: `unique_holder_character_id IS NULL` → set (transaction).
3. Insert `character.title_owned (character_id, title_id, is_active=false)`.
4. Insert `progression.title_earned` log.
5. Publish event `character.title_earned`.
6. Notification "New title earned!"

Idempotent:
- Duplicate insert → ON CONFLICT DO NOTHING.
- Server-first race: transaction; ilk kazanır.

## 7. Set Active Title

**POST `/v1/characters/{id}/set-title`**
Auth + CharacterScope.

Request:
```json
{ "title_id": "t-..." }
```
(veya `title_id = null` — hiçbir title aktif değil, sadece isim gösterilir.)

Server:
1. Verify ownership: `character.title_owned` row exists.
2. Transaction:
   - `UPDATE character.title_owned SET is_active = false WHERE character_id = $1 AND is_active = true`.
   - `UPDATE character.title_owned SET is_active = true WHERE character_id = $1 AND title_id = $2` (or skip if null).
3. Publish event `character.title_activated` (guild/party name cache invalidate).
4. Response 200.

Rate limit: 60/hour (spam koruma).

## 8. Display Composition

Client renders:
```
title.format.replace("{name}", character.name)
```

Örnek:
- Prefix: `"Grandmaster Rhaegar"`.
- Suffix: `"Rhaegar, Ashen Wyrm Slayer"`.

Both varsa (aktif tek); ancak Phase 02'de "prefix + suffix" combine (2 slot) opsiyonu.

## 9. Public Visibility
- Character sheet (kendi görür): tüm owned + active.
- Diğer oyuncu (public profile): aktif title.
- Nameplate over character (WebSocket render hint): aktif title formatted.

## 10. Hidden Titles
- `hidden_until_earned = true`: title'ın kendisi (achievement/lore spoiler) earn olmadan görünmez.
- Achievement panel'de gizli slot placeholder.

## 11. Server-First Race Condition
- Achievement engine: same trigger multiple characters aynı tick'te tetiklerse:
  - Transaction serializable.
  - İlk `unique_holder_character_id` set kazanır.
  - Diğerleri farklı title alır (örn. "Ashen Wyrm Hunter" — non-first version) veya hiç almaz (config).

## 12. Endpoint

### 12.1 Owned Titles List

**GET `/v1/characters/{id}/titles`**
Response:
```json
{
  "active_title_id": "t-...",
  "titles": [
    {
      "id": "t-...",
      "code": "grandmaster_miner",
      "display_name_key": "title.grandmaster_miner",
      "format": "Grandmaster {name}",
      "placement": "prefix",
      "source": "profession",
      "earned_at": "..."
    }
  ]
}
```

### 12.2 Available But Not Owned (Discovery)

**GET `/v1/characters/{id}/titles/available`** (opsiyonel)
- Owned + hidden=false + season=current.
- "How to earn" hint (achievement id or condition summary).

## 13. Rename Interaction
- Rename değişse bile title format `{name}` yeni ismi kullanır (dinamik).
- Cached UI invalidate.

## 14. Deletion Interaction
- Character silinirse `character.title_owned` cleanup.
- Server-first title: `unique_holder_character_id = NULL` (yeni sahiplenilebilir mi? — HAYIR; `unique_holder_recycled` flag ile önlenir; audit).

## 15. Event
```json
{
  "_v": 1,
  "character_id": "c-777",
  "title_id": "t-...",
  "title_code": "grandmaster_miner",
  "server_first": false
}
```

## 16. Cross-Platform
Aynı character sheet her platformda; title display standart.

## 17. Metrics
- `character_title_earned_total{code}`.
- `character_title_active_distribution` — hangi title'lar aktif.

## 18. Audit
- `character.title_earned`.
- `character.title_activated`.

## 19. Test Stratejisi
- Grant → owned → set active happy path.
- Server-first concurrent trigger.
- Rename with active title → format re-render.
- Delete → cleanup + server-first lock.

## 20. İlişkiler
- [Character](character.md), [Reputation](reputation.md), [game-design/achievements.md](../game-design/achievements.md), [game-design/mastery.md](../game-design/mastery.md), [game-design/ascension.md](../game-design/ascension.md), [database/character.md](../database/character.md), [database/content-templates.md](../database/content-templates.md)

## 21. Bekleyen Kararlar
- Prefix + suffix aynı anda kullanım (Phase 02).
- Kozmetik shop title'ları (monetization Phase 02).
- Guild tabard renk/simge kombine (Phase 02).
- Custom title (Phase 03 opsiyonel; moderation heavy).
