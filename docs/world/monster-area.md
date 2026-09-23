# Monster Area (Hunt Spot)

## 1. Amaç
Karakterin AFK hunting yapabildiği mob spawn bölgesi. Location'a bağlı; mob pool + spawn interval + danger tier tanımlar.

Kaynak: [database/content-templates.md §17](../database/content-templates.md), [activity-engine/hunting.md](../activity-engine/hunting.md).

## 2. Model

**Spawn Template** (`content.monster_area_spawn`, yeni tablo Phase 13):
```sql
CREATE TABLE content.monster_area_spawn (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  location_id UUID NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  mob_pool_json JSONB NOT NULL,        -- [{npc_template_code, weight}]
  max_active_mobs SMALLINT NOT NULL DEFAULT 10,
  respawn_interval_ms INT NOT NULL DEFAULT 5000,
  recommended_level_min SMALLINT,
  recommended_level_max SMALLINT,
  danger_tier VARCHAR(16) NOT NULL DEFAULT 'normal',   -- easy/normal/hard/elite
  rare_mob_template_code VARCHAR(48) NULL,
  rare_mob_spawn_windows_json JSONB NULL,               -- opsiyonel rare mob
  version INT NOT NULL DEFAULT 1
);
```

**Runtime** (Redis-based, no DB persist):
- `spawn:<location_id>` sorted set: aktif mob sayısı.
- Personal instance mantığı: karakter perspektifinden aktivite bağımsız yield.

## 3. Content Pack

```yaml
# content/monster_area_spawns/iron_valley_wolves.yaml
code: iron_valley_wolves
location_code: iron_valley
display_name_key: monster_area.iron_valley_wolves.name

mob_pool:
  - npc_template_code: ashen_wolf
    weight: 60
  - npc_template_code: shadow_wolf
    weight: 25
  - npc_template_code: alpha_wolf
    weight: 15

max_active_mobs: 15
respawn_interval_ms: 5000
recommended_level_min: 15
recommended_level_max: 25
danger_tier: normal

rare_mob_template_code: shadowfang_alpha
rare_mob_spawn_windows:
  - cron: "0 */6 * * *"       # her 6 saatte bir potansiyel spawn
    chance_pct: 30
```

## 4. Mob Pool Weight

Weight-based selection (aktivite hunting'de mob type filter):
- Character `params.mob_type_codes` filter → pool subset.
- Kill count analytical hesap (bkz. [activity-engine/hunting.md §5](../activity-engine/hunting.md)).

## 5. Danger Tier

| Tier | Notlar |
|---|---|
| `easy` | Baş harci level -2..+2 |
| `normal` | Standart level match |
| `hard` | +5 level; group hunt recommended |
| `elite` | +10 level; party required tercih |

Danger tier → mob HP/damage scale + rare drop bonus.

## 6. Personal Instance vs Shared

- **Personal** (v1 tercih): karakter perspektifinden mob spawn interval sabit; her kille yeni mob spawn (Redis TTL).
- **Shared** (Phase 02 opsiyonel): tüm karakterler aynı spawn pool; competition + party bonuslar.

Karar: **personal** (idle MMO psikolojisi + AFK yield tutarlılığı).

## 7. Rare Mob

- Spawn window cron-based.
- Kill by first character; despawn.
- Rare drop chance yüksek (rare weapon, unique pet).
- Server-wide anons (opsiyonel Phase 02).

## 8. Location Bağlama

Location activity_hooks:
```yaml
activity_hooks:
  - kind: monster_area
    spawn_code: iron_valley_wolves
```

Bir location birden fazla monster area barındırabilir (Phase 02):
- Farklı mob family, farklı zone içinde.

## 9. Endpoint

**GET `/v1/content/world/locations/{code}/monsters`**
```json
{
  "location_code": "iron_valley",
  "monster_areas": [
    {
      "spawn_code": "iron_valley_wolves",
      "display_name_key": "...",
      "mob_pool": [ ... ],
      "danger_tier": "normal",
      "recommended_level_min": 15,
      "recommended_level_max": 25,
      "current_rare_active": false
    }
  ]
}
```

## 10. Hunt Aktivite Integration

Bkz. [activity-engine/hunting.md](../activity-engine/hunting.md).

Karakter hunt spot'a gelir → aktivite başlatır → analytical yield.

## 11. World Boss Spawn (Farklı Sistem)

World boss ayrı — bkz. [game-design/world-bosses.md](../game-design/world-bosses.md) + [database/progression-quests-achievements.md §7](../database/progression-quests-achievements.md).

Monster area world boss'la karışmaz; world boss dedicated location + spawn window.

## 12. Mob Family Bestiary

Kill'ler bestiary mastery'e katkı verir (bkz. [character/mastery.md](../character/mastery.md) via game-design; skill-engine/mastery.md).

## 13. Metrics
- `monster_area_kills_total{spawn,mob}` counter.
- `monster_area_rare_mob_spawn_total{spawn}`.
- `monster_area_active_character_hunts{spawn}` gauge.

## 14. Audit
- `content.monster_area_created/updated`.
- `combat.rare_mob_killed`.

## 15. Cross-Reference
- [Location](location.md), [Region](region.md), [Dungeon Area](dungeon-area.md)
- [game-design/hunting.md](../game-design/hunting.md)
- [database/content-templates.md §17](../database/content-templates.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)

## 16. Bekleyen Kararlar
- Personal vs shared spawn v1'de personal; Phase 02'de opsiyon.
- Rare mob server-wide anons Phase 02.
- Elite spawn party gating enforce (v1: soft recommendation; Phase 02 harder).
- Weather effect (nocturnal mob'lar) Phase 02.
