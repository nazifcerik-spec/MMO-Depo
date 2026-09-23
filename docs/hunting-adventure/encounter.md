# Encounter (Random Encounter Framework)

## 1. Amaç
Hunt/adventure/travel sırasında **inject edilen** olaylar. Rare mob spawn, treasure chest, ambush, npc dialog, mini-quest. AFK reveal içinde ortaya çıkar; player'a "sürprizle" hissi.

## 2. Encounter Türleri

| Kind | Örnek | Trigger |
|---|---|---|
| `rare_mob_spawn` | Alpha wolf, shadowfang | Time-based random |
| `elite_group` | 3 elite bandit | Region danger scale |
| `treasure_chest` | Hidden chest | Random discovery |
| `world_drop` | Rare item on ground | Ultra-rare |
| `npc_dialog` | Wandering merchant | Story lore |
| `mini_quest` | "Help this NPC" popup | Random branch |
| `ambush` | Bandit attack | Travel/hunt injection |
| `event_flag` | Server-wide event tag | Region-wide |

## 3. Model

```sql
CREATE TABLE content.encounter_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  kind VARCHAR(24) NOT NULL,
  region_ids UUID[] NOT NULL DEFAULT '{}',      -- hangi region'larda aktif
  location_tags VARCHAR(24)[] NOT NULL DEFAULT '{}',  -- location filter
  activity_kinds VARCHAR(24)[] NOT NULL DEFAULT '{}',  -- hangi aktivite tipinde inject
  chance_per_hour NUMERIC(5,2) NOT NULL,        -- baz spawn oranı
  cooldown_ms BIGINT NOT NULL DEFAULT 3600000,   -- character-per-encounter cooldown
  payload_json JSONB NOT NULL,                    -- encounter özel data
  reward_pool_id UUID NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE character.encounter_history (
  character_id UUID NOT NULL,
  encounter_template_id UUID NOT NULL,
  last_triggered_at TIMESTAMPTZ NOT NULL,
  total_triggered INT NOT NULL DEFAULT 0,
  PRIMARY KEY (character_id, encounter_template_id)
);
```

## 4. Content Pack

```yaml
# content/encounters/wandering_merchant.yaml
code: wandering_merchant
kind: npc_dialog
region_ids: [ashen_hills, emerald_coast]
location_tags: [hunt_spot, travel_route]
activity_kinds: [hunting, travel]
chance_per_hour: 0.5           # ~%50 per hour of activity
cooldown_ms: 21600000           # 6 hours per character

payload:
  npc_code: merchant_orin
  dialog_key: encounter.wandering_merchant.dialog
  offer_type: rare_potion_discount
  discount_pct: 50
  duration_ms: 1800000

reward_pool_id: null
```

## 5. Injection Kural Seti

Reveal içinde:
```go
for _, activity := range contract.Activities {
    elapsedHours := activity.ElapsedMs / 3600000.0
    for _, encTemplate := range encountersFor(activity, snapshot.Location, snapshot.RegionID) {
        // Character cooldown check
        history := charEncounterHistory[encTemplate.ID]
        if history.LastTriggeredAt.After(now.Add(-encTemplate.CooldownMs)) { continue }
        // Poisson sample: expected = chance_per_hour × elapsedHours
        expected := encTemplate.ChancePerHour * elapsedHours
        count := samplePoisson(expected, 1.0, seed.Sub("enc-"+encTemplate.Code))
        if count > 0 {
            triggerEncounter(encTemplate, count, ...)
        }
    }
}
```

Cooldown per character-per-encounter template — spam koruma.

## 6. Trigger Effects

Trigger'e göre:
- **rare_mob_spawn**: kill loot roll (bkz. [rare-enemy.md](rare-enemy.md)).
- **treasure_chest**: loot table roll (bkz. [treasure.md](treasure.md)).
- **npc_dialog**: reveal narration; opsiyonel offer (potion discount, quest hint).
- **mini_quest**: quest offer (Phase 02); auto-decline v1.
- **ambush**: combat inject (bkz. [Death & Revive](../combat/death-revive.md)); potential death.
- **event_flag**: server-wide event participation credit.

## 7. Reveal Aggregation

Encounter events reveal summary'e eklenir:
```json
{
  "notable_encounters": [
    {
      "kind": "rare_mob",
      "mob_code": "shadowfang_alpha",
      "defeated": true,
      "rare_drop": { "item_code": "shadowfang_tooth" }
    },
    {
      "kind": "treasure",
      "chest_code": "hidden_forest_chest",
      "items": [ ... ]
    }
  ]
}
```

Client "reveal panel" UI'da encounters highlight.

## 8. Encounter Cooldown vs Chance

Balance:
- **Common encounters** (merchant): chance_per_hour ~0.5-1, cooldown 6h.
- **Rare encounters** (rare mob spawn): chance_per_hour ~0.1, cooldown 12h.
- **Ultra rare** (world drop): chance_per_hour ~0.01, cooldown 48h.

Config playtest tune.

## 9. Region-Specific Encounter Pools

Region content:
- Ashen Hills: bandits, wolves, wandering merchant, hidden mine chest.
- Emerald Coast: pirates, sea creatures, treasure chest, sunken artifact.
- Frostspire: yeti, ice ambush, frozen relic.
- ...

Content team her region için encounter pool tanımlar.

## 10. Cross-Activity Injection

Encounter injection aktivitesine göre:
- **Hunting**: rare_mob_spawn, elite_group, ambush.
- **Adventure**: stage-specific narrative encounter.
- **Travel**: ambush, wandering_merchant, treasure_chest.
- **Gathering** (mining/fishing): treasure_chest, npc_dialog, rare_node.

Content `activity_kinds` filter.

## 11. Endpoint

**GET `/v1/characters/{id}/encounters/recent`** — son 30 gün encounter history.

## 12. Metrics
- `encounter_triggered_total{kind}` counter.
- `encounter_cooldown_blocked_total{kind}`.
- `encounter_rare_mob_spawn_total`.

## 13. Cross-Reference
- [Hunting](hunting.md), [Rare Enemy](rare-enemy.md), [Elite Enemy](elite-enemy.md), [Treasure](treasure.md), [Random Event](random-event.md), [Adventure](adventure.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/adventure.md](../activity-engine/adventure.md), [activity-engine/travel.md](../activity-engine/travel.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 14. Bekleyen Kararlar
- Encounter cooldown per template vs global.
- Content team encounter pool ihtiyaç ölçek (per region ~10 encounter v1).
- Ambush combat outcome (auto-loss vs win-based on DPS).
- Encounter chain (encounter triggers new encounter) Phase 02.
