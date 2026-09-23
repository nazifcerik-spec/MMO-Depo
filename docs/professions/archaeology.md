# Archaeology

> ⚠️ **v1'de kapalı**; **Phase 02** aktivasyon. Schema + design hazır.

## 1. Amaç
Kazı yaparak ancient artifact toplama. Rare item + quest chain + kozmetik unlock. **Unique** profession — koleksiyon odaklı, ekonomi'ye "nadir aristokrat" hediyeler.

## 2. Neden Phase 02?

- Content-heavy (dig site + artifact + lore).
- Uzun cycle mekaniği.
- v1 core loop odak; Phase 02 zenginlik.

## 3. Model (Phase 02 Schema)

```sql
CREATE TABLE content.dig_site_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  location_id UUID NOT NULL,
  display_name_key VARCHAR(64),
  culture VARCHAR(24) NOT NULL,      -- ancient_human, elven_ruins, dwarven_forge, etc.
  tier SMALLINT NOT NULL,
  min_profession_level SMALLINT NOT NULL,
  dig_time_ms BIGINT NOT NULL,        -- örn 30 min
  artifact_loot_table_id UUID NOT NULL,
  respawn_ms BIGINT NOT NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE character.dig_progress (
  character_id UUID NOT NULL,
  dig_site_id UUID NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  next_available_at TIMESTAMPTZ NOT NULL,
  fragments_collected INT NOT NULL DEFAULT 0,
  PRIMARY KEY (character_id, dig_site_id)
);
```

## 4. Profession Data

```yaml
# content/professions/archaeology.yaml
code: archaeology
category: gathering
display_name_key: profession.archaeology.name
max_level: 100
xp_curve: { base: 50, factor: 1.10 }
tool_slot: null                    # brush + shovel (Phase 02)
race_affinity_bonus_pct:
  arcanborn: 5                     # lore affinity
  human: 5                          # scholar bias
```

## 5. Culture Kataloğu

| Culture | Dig Site Örnek | Unique Reward |
|---|---|---|
| `ancient_human` | Ashen Ruins | Human kozmetik + old_coin currency |
| `elven_ruins` | Silverbough Vault | Elf transmog + magic essence |
| `dwarven_forge` | Deep Mine Ruins | Dwarven mount blueprint + iron_relic |
| `orcish_warcamp` | Warcamp Bones | Orc weapon skin + bone_relic |
| `merfolk_sunken` | Sunken City | Aquatic mount unlock + coral_shard |
| `arcanborn_spire` | Old Spire | Arcane recipe + arcane_shard |
| `unknown_precursor` | Precursor Anomaly | Legendary trinket (Phase 02+ endgame) |

Culture unlock: quest chain veya level gated.

## 6. Dig Flow

```
1. Character at dig_site location.
2. Start dig aktivite (30 min - 4 hour depending tier).
3. AFK devam (single dig activity contract; primary channel).
4. Complete → artifact roll from loot table.
5. Fragment based:
   - Bazı artifact "fragment"ler collected → tam artifact assemble için N fragment gerekir.
6. Respawn timer (site).
```

## 7. Artifact Türleri

- **Fragment**: 3-10 combine → complete artifact.
- **Complete artifact**: tradeable (BoE), museum donation (Phase 03), quest turn-in.
- **Unique relic**: legendary lore item.

## 8. AFK Aktivite

Aktivite tipi `archaeology_dig`; contract channel=primary; dig_time long duration.

## 9. Endpoint (Phase 02)

```
GET  /v1/characters/{id}/archaeology/dig-sites   # near-by
POST /v1/characters/{id}/archaeology/dig/start
POST /v1/characters/{id}/archaeology/dig/{contract_id}/harvest
GET  /v1/characters/{id}/archaeology/fragments
POST /v1/characters/{id}/archaeology/assemble    # fragment combine
```

## 10. Race Affinity

- **Arcanborn** + **Human**: secondary affinity (+5%; scholar bias).

## 11. Cross-Reference

- [Herbalism](herbalism.md), [Mining](mining.md)
- [activity-engine/README.md](../activity-engine/README.md) — pattern
- [game-design/professions.md](../game-design/professions.md)

## 12. Bekleyen Kararlar

- Aktivasyon Phase 02 milestone.
- Museum system (Phase 03) — donate artifact, achievement + reward.
- Cross-culture bonus (all 6 cultures grandmaster = "Historian" title).
- Trade rules (unique relic BoP vs BoE).
- Archaeology event (server-wide "expedition") Phase 02.
