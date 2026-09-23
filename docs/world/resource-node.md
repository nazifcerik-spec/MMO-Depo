# Resource Node

## 1. Amaç
Gathering aktivitesi için node instance modelleme. Mining/fishing/woodcutting/foraging node'ları location'a bağlı, respawn timer'lı.

Kaynak: [database/world.md §7](../database/world.md), [activity-engine/mining.md](../activity-engine/mining.md), [activity-engine/fishing.md](../activity-engine/fishing.md).

## 2. Model

**Template** (`content.resource_node_template`, kalıcı):
```sql
content.resource_node_template
  id, code, kind (mining/fishing/foraging/woodcutting/skinning)
  tier SMALLINT, display_name_key, art_ref
  respawn_ms INT
  yield_curve_json                -- (ör. base yield per rank)
  loot_table_id (rare drop)
  primary_output_item_code
  secondary_output_item_code
  primary_share NUMERIC, secondary_share NUMERIC
```

**Instance** (`world.resource_node`, runtime):
```sql
world.resource_node
  id, template_id (FK)
  location_id (FK)
  next_respawn_at TIMESTAMPTZ NULL     -- dry ise ne zaman geri gelir
  spawn_group_id UUID NULL              -- cluster spawn kontrolü
  created_at
```

## 3. Personal Instance (DD-022)

Kural: **her karakter node'un kendi görüntüsünü görür** (personal instance). Node "dry" durumu karakter-özel.

Aslında server node instance'ları paylaşımlı satır; karakter perspektifi:
```
availability_for_char = f(node_last_yield_by_char, node.respawn_ms)
```

Yani DB'de tek satır; ancak "yielded_at" per-character track:
```sql
CREATE TABLE world.resource_node_character_state (
  character_id UUID NOT NULL,
  resource_node_id UUID NOT NULL,
  last_yield_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  next_available_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (character_id, resource_node_id)
);
```

Alternative simplification (v1 tercih):
- Node instance shared row (as-is schema).
- Character-özel state Redis'te (`char:node:<char>:<node>` TTL respawn_ms).
- Redis kayıp → node available.
- **Personal instance = Redis TTL-based** (basit + performant).

## 4. Content Pack — Template

```yaml
# content/resource_node_templates/iron_ore_node.yaml
code: iron_ore_node
kind: mining
tier: 2
display_name_key: node.iron_ore.name
art_ref: nodes/iron_ore.png

respawn_ms: 60000                # 60 saniye
yield_curve:
  base_yield: 3
  per_profession_level_bonus: 0.05
loot_table_code: iron_ore_rare_drop

primary_output_item_code: iron_ore
secondary_output_item_code: gray_stone
primary_share: 0.75
secondary_share: 0.25
```

## 5. Location Bağlama

Location'ın `activity_hooks` alanı:
```yaml
activity_hooks:
  - kind: resource_node
    template_code: iron_ore_node
    count: 6                # bu location'da 6 node instance
    respawn_ms: 60000
```

Location seed → runtime N instance oluşur (initial deploy):
```sql
INSERT INTO world.resource_node (id, template_id, location_id, next_respawn_at, spawn_group_id)
SELECT gen_random_uuid(), template_id, $loc_id, NULL, gen_random_uuid()
FROM generate_series(1, count);
```

Node group: aynı seed grup respawn cadence sync (opsiyonel).

## 6. Runtime Availability

Karakter gather aktivite başlatınca (bkz. [activity-engine/mining.md](../activity-engine/mining.md)):
- Node personal availability (Redis) check.
- Yield analytical hesap (per-tick loop yok).
- Post-yield: Redis TTL set (`char:node:<char>:<node>` EX respawn_ms).

Aktivite sürekli — karakter bir gather aktivite'de belirli spot'ta; server "her N saniye bir node yield" formuluyle çalışır:
```
nodes_per_hour = (total_nodes / respawn_sec) * 3600
yield_per_hour = nodes_per_hour * avg_node_yield
```

Bu yield rate compute'a katkı; per-node simulation yok.

## 7. Rare Node (Server-Wide)

Bazı node "rare" — nadir spawn, server-wide announcement:
- `content.resource_node_template.rare_flag = true`.
- Spawn interval çok uzun (saatler).
- Spawn olduğunda server chat anons (opsiyonel).
- İlk gather ile despawn (yaşam boyu tek karakter).

v1'de rare node opsiyonel; Phase 02'de tam mekanik.

## 8. Endpoint

**GET `/v1/content/world/locations/{code}/nodes`**
```json
{
  "location_code": "iron_valley",
  "nodes": [
    {
      "instance_id": "...",
      "template_code": "iron_ore_node",
      "tier": 2,
      "available_for_character": true,     // Redis check
      "next_respawn_at": null              // personal
    }
  ]
}
```

## 9. Depletion & Respawn

- **Personal** (v1 tercih): Redis TTL; karakter için node yield sonrası respawn süresi bekle.
- **Global** (v2 opsiyonel): world.resource_node.next_respawn_at set; tüm karakterlere aynı respawn.

Global depletion sadece rare node için (event context).

## 10. Yield Rate Formula (Aktivite Bağlamı)

Bkz. [activity-engine/mining.md §4](../activity-engine/mining.md).

Rate = f(total_nodes, respawn_ms, character_gather_yield_pct, mastery_bonus, race_affinity).

## 11. Node Migration (Content Patch)

Location'ın node profili değişince (patch):
- Migration: `world.resource_node` satır yeniden seed.
- Redis state invalidate.
- Karakter aktif contract snapshot content_version bağlı; yeni template devreye Phase 01 sonu.

## 12. Cross-Reference
- [Location](location.md), [Region](region.md)
- [database/world.md §7](../database/world.md)
- [database/content-templates.md §20](../database/content-templates.md)
- [activity-engine/mining.md](../activity-engine/mining.md), [activity-engine/fishing.md](../activity-engine/fishing.md)
- [game-design/professions.md](../game-design/professions.md)

## 13. Metrics
- `resource_node_yield_total{kind,tier}` counter.
- `resource_node_rare_spawn_total`.
- `resource_node_available_ratio{location}` gauge.

## 14. Bekleyen Kararlar
- Personal Redis vs global DB depletion final (v1 personal Redis).
- Rare node mechanic Phase 02 tam.
- Tool durability drain per gather (v1 yok).
- Weather effect (Phase 02 rain → fishing bonus).
