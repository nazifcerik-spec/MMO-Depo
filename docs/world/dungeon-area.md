# Dungeon Area

## 1. Amaç
Dungeon giriş noktası. Location'a bağlı, template ile dungeon instance yaratma anchor'ı.

Kaynak: [database/content-templates.md §14](../database/content-templates.md), [activity-engine/dungeon.md](../activity-engine/dungeon.md).

## 2. Model

Dungeon template zaten `content.dungeon_template` içinde. Dungeon area = location'ın hook'u:
```yaml
# location activity_hooks
activity_hooks:
  - kind: dungeon_area
    template_code: ashen_keep
    entrance_flavor: "Ashen Keep Gates"
```

Ek olarak yeni tablo (opsiyonel):
```sql
CREATE TABLE content.dungeon_area (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  location_id UUID NOT NULL,
  dungeon_template_id UUID NOT NULL,
  entrance_flavor_key VARCHAR(64) NULL,
  entrance_art_ref VARCHAR(128) NULL,
  discovered_by_first_char UUID NULL,       -- fun trivia
  version INT NOT NULL DEFAULT 1
);
```

v1 tercih: dedicated tablo minimalist; location activity_hook JSON yeter.

## 3. Content Pack (Location Bağlamında)

```yaml
# content/world/locations/ashen_keep_gates.yaml
code: ashen_keep_gates
parent_code: ashen_hills           # region
parent_kind: region
display_name_key: location.ashen_keep_gates.name
tags: [dungeon_entry, safe_zone]

activity_hooks:
  - kind: dungeon_area
    template_code: ashen_keep
    entrance_flavor_key: dungeon.ashen_keep.entrance_flavor
```

## 4. Dungeon Entry Flow

**POST `/v1/dungeons/enter`**
Auth + CharacterScope.

Request:
```json
{
  "dungeon_template_code": "ashen_keep",
  "difficulty": "normal",
  "mode": "manual",           // manual | idle
  "party_id": null
}
```

Server:
1. Load character.
2. Validate:
   - Character at dungeon area location (v1 opsiyonel; Phase 02 zorunlu).
   - Character level >= dungeon.min_level.
   - Party size ok (mode=manual).
   - Weekly lockout not full.
3. Mode dispatch:
   - **manual**: create `activity.instance` (bkz. [activity-engine/dungeon.md §5](../activity-engine/dungeon.md)).
   - **idle**: create `activity.contract` (channel=primary).
4. Response 200 + instance/contract ref.

## 5. Location Entrance Requirements

v1 opsiyonel; Phase 02 zorunlu:
- Karakter dungeon area location'unda olmalı (physical presence).
- Alternative: waygate teleport ile giriş.

v1: `current_location_id` any → OK.
Phase 02: `current_location_id == dungeon_area.location_id` zorunlu.

## 6. Dungeon Instance Management

Manuel dungeon:
- `activity.instance` oluşur, party bağlı.
- Expire 2h (config).
- Multiple concurrent instances mümkün (different parties).

Idle dungeon:
- `activity.contract` bir karakter kanalı; party ise her üye kendi contract'ını başlatır.

Detay: [activity-engine/dungeon.md](../activity-engine/dungeon.md).

## 7. Dungeon Zorluk

Content:
- `content.dungeon_template.estimated_total_hp` alanı per difficulty.
- `boss_encounters` array.
- `weekly_lockout` flag.

Difficulty:
- **normal**: baseline.
- **heroic**: +HP +damage; better loot.
- **mythic** (Phase 02): affix system.

## 8. Weekly Lockout

Bkz. [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md).

- `progression.dungeon_lockout` per character + dungeon + kind (weekly/daily).
- Weekly reset Salı 06:00 UTC (DD-016).
- Lockout aktif iken loot azalır (rare drop penalty); XP normal.

## 9. First Clear Bonus

- Server-wide first clear → title reward, cosmetic drop.
- Character first weekly clear → full reward.
- Sonraki weekly clear → reduced rare loot.

## 10. Endpoint

**GET `/v1/content/dungeons/{code}`**
```json
{
  "code": "ashen_keep",
  "display_name_key": "dungeon.ashen_keep.name",
  "region_code": "ashen_hills",
  "min_level": 20,
  "max_level": 30,
  "party_size_range": [1, 5],
  "weekly_lockout": true,
  "boss_encounters": [
    { "code": "keep_guardian", "order": 1 },
    { "code": "fire_wraith", "order": 2 },
    { "code": "ashen_lich_king", "order": 3 }
  ],
  "difficulty_options": ["normal", "heroic"],
  "estimated_clear_time_min": { "normal": 20, "heroic": 35 }
}
```

**GET `/v1/content/world/locations/{code}/dungeon`**
Dungeon area info specifically for location.

## 11. Cross-Region Party Dungeon

Party farklı region karakterler:
- Dungeon area location'u party leader'ın erişimi.
- Party members auto-teleport to dungeon area (v1'de basit) veya karakter kendi travel.

v1 tercih: leader coordinates; party members explicit accept.

## 12. Metrics
- `dungeon_area_entry_total{dungeon,mode}`.
- `dungeon_completion_total{dungeon,difficulty}`.
- `dungeon_wipe_total{dungeon,encounter}`.

## 13. Audit
- `dungeon.entered` (character + party).
- `dungeon.completed`.
- `dungeon.first_server_clear`.

## 14. Cross-Reference
- [Location](location.md), [Region](region.md), [Monster Area](monster-area.md)
- [game-design/dungeons.md](../game-design/dungeons.md), [game-design/party.md](../game-design/party.md)
- [database/content-templates.md §14-15](../database/content-templates.md), [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md)

## 15. Bekleyen Kararlar
- Entrance location zorunluluğu v1 vs Phase 02.
- Party member auto-teleport policy.
- Mythic affix Phase 02+.
- Raid (10+) area (Phase 02).
- Dungeon reset item (weekly bypass consumable) Phase 02 monetization.
