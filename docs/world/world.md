# World

## 1. Amaç
Tüm dünyayı temsil eden **singleton** — meta root. v1'de tek dünya (single-shard, bkz. DD-020); Phase 03+ multi-shard/multi-world olasılığı için hazır.

## 2. Model

Şema (basit):
```sql
CREATE TABLE world.world (
  id UUID PRIMARY KEY,
  code VARCHAR(32) UNIQUE NOT NULL,       -- "aetherion" (default v1)
  display_name_key VARCHAR(64) NOT NULL,
  version INT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  metadata_json JSONB NOT NULL DEFAULT '{}'
);
```

v1: tek satır (code='aetherion' — placeholder).

## 3. Config

```yaml
# content/world/world.yaml
code: aetherion
display_name_key: world.aetherion.name
metadata:
  server_time_zone_hint: UTC
  weekly_reset_day: tuesday
  weekly_reset_hour_utc: 6
```

## 4. Cross-Continent Kural

- **Tek dünya, tek shard** (v1).
- Continents `world_id = <single_world_id>`.
- Karakter herhangi bir continent'ta olabilir.
- Market kıta bazlı global; ancak dünya seviyesinde tek pool.

## 5. Endpoint

**GET `/v1/content/world`**
Public.
```json
{
  "code": "aetherion",
  "display_name_key": "world.aetherion.name",
  "continent_count": 2,
  "metadata": { "weekly_reset": "tuesday_06:00_utc" }
}
```

## 6. Multi-World (Phase 03)

Not planned v1. Ancak schema `world_id` FK'sini continents/regions'da tutar — multi-world eklenirse migration kolay.

## 7. Server Time

`weekly_reset_day + hour_utc` → scheduler `system.scheduled_task` içine kayıt (bkz. [database/audit-system.md §12](../database/audit-system.md)).

Weekly reset workflow (bkz. [backend/worker.md §8.8](../backend/worker.md)):
- Dungeon lockout reset.
- Weekly quest rotation.
- Event rotation.

## 8. Metrics
- `world_singleton_healthy` gauge (should be 1).
- `world_content_version` gauge (numeric version).

## 9. Cross-Reference
- [Continent](continent.md), [Region](region.md)
- [database/world.md §1](../database/world.md)
- DD-016 (weekly reset time), DD-020 (single shard v1)

## 10. Bekleyen Kararlar
- Multi-world (Phase 03) — sharding stratejisi.
- Cross-world event (Phase 03) — special event framework.
