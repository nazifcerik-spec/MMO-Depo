# Territory

## 1. Amaç
Region'ın alt bölümü — **politik/faction** anchor. Settlement'ları gruplandırır. Phase 02+ territory war için hazır alan.

## 2. Model
Kaynak: [database/world.md §3](../database/world.md).

```sql
world.territory
  id, code, region_id (FK)
  display_name_key
  faction_id (FK → content.faction NULL)
  controlling_guild_id (FK → social.guild NULL)   -- Phase 02
  art_ref
```

Ek:
```sql
ALTER TABLE world.territory
  ADD COLUMN lore_key VARCHAR(64) NULL,
  ADD COLUMN tags VARCHAR(24)[] NOT NULL DEFAULT '{}',
  ADD COLUMN pvp_enabled BOOLEAN NOT NULL DEFAULT false,   -- Phase 02 openworld PvP
  ADD COLUMN deprecated_at TIMESTAMPTZ NULL;
```

## 3. Content Pack

```yaml
# content/world/territories/ashen_flats.yaml
code: ashen_flats
region_code: ashen_hills
display_name_key: territory.ashen_flats.name

faction_code: emerald_wardens   # default faction control (v1); guild control Phase 02
tags: [starter, farmland]
```

## 4. Territory Composition

Territory içerir:
- 0-1 city (bazı territory sadece rural).
- 0-N village.
- 0-N outpost.
- 0-N location (POI, resource nodes, hunt spots, dungeon).

Constraint: territory içi settlement toplam max ~10 (v1); balance simülasyon.

## 5. Endpoint

**GET `/v1/content/world/territories/{code}`**
```json
{
  "code": "ashen_flats",
  "region_code": "ashen_hills",
  "display_name_key": "territory.ashen_flats.name",
  "faction": { "code": "emerald_wardens", "tier": "friendly_for_character" },
  "controlling_guild": null,
  "settlements": [
    { "code": "ashfield", "kind": "city", "size": "capital" },
    { "code": "wheatcross", "kind": "village", "size": "village" }
  ],
  "locations_count": 12
}
```

## 6. Faction Control

- **Static faction**: content'te `faction_code` — v1 default.
- **Guild control**: `controlling_guild_id` — Phase 02+ territory war sonucu.

Faction control etkileri:
- Vendor fiyat modifier (rep tier'a göre).
- Quest availability (faction quest chain).
- Bazı POI faction locked (rep required).

Guild control (Phase 02+):
- Guild tax on marketplace transactions in territory.
- Guild vendor exclusive item access.
- Guild banner visible.

## 7. PvP Territory (Phase 02+)

- `pvp_enabled=true` → open-world PvP zone (opt-in flag Phase 02).
- v1: tüm territory safe.
- Phase 02: bazı territory contested zone.

## 8. Deprecation

- Faction/guild control automatic revert to neutral.
- Karakter safe city'ye teleport.

## 9. Territory War Hazırlık (Phase 02+)

Alanlar hazır:
- `controlling_guild_id`.
- `is_at_war_since TIMESTAMPTZ NULL`.
- `war_target_guild_id UUID NULL`.
- `defense_meter INT NOT NULL DEFAULT 0`.

Schema önerisi:
```sql
ALTER TABLE world.territory
  ADD COLUMN is_at_war_since TIMESTAMPTZ NULL,
  ADD COLUMN war_attacker_guild_id UUID NULL,
  ADD COLUMN defense_meter INT NOT NULL DEFAULT 0;
```

## 10. Metrics
- `territory_faction_control_distribution{faction}` gauge.
- `territory_guild_control_active_total`.

## 11. Audit
- `content.territory_created/updated/deprecated`.
- `territory.faction_control_changed`.
- `territory.guild_control_changed` (Phase 02).

## 12. Cross-Reference
- [Region](region.md), [City](city.md), [Village](village.md), [Outpost](outpost.md)
- [character/reputation.md](../character/reputation.md)
- [game-design/guilds.md](../game-design/guilds.md) — future guild war
- [database/world.md §3](../database/world.md)

## 13. Bekleyen Kararlar
- Territory war Phase 02 tam mekanik.
- Guild tax %.
- PvP zone opt-in vs mandatory.
- Territory boss (unique boss per territory) Phase 02+.
