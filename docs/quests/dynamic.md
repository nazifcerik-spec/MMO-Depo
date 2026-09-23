# Dynamic (Procedural / Random Quest)

## 1. Amaç
Procedural quest — content template'den runtime yaratılan quest instance. Bounty board, random NPC dialog, character-generated quest.

**v1'de basit**; Phase 02 tam procedural.

## 2. Karakteristik

- **Procedural**: template-based ancak parameters random.
- **Time-limited**: kısa süre (1-4 saat).
- **Bounty board**: NPC-generated (kill target X, gather Y).
- **Rare spawn quest**: rare mob spawn → auto-quest.

## 3. Content Template

```yaml
# content/quests/dynamic_bounty_template.yaml
code: dynamic_bounty_template
category: dynamic
subcategory: bounty
display_name_key: quest.dynamic.bounty.name

# Parameter placeholders (runtime substitute)
parameters:
  target_mob_family: {type: choice, choices: [beast, humanoid, undead]}
  kill_count: {type: range, min: 10, max: 30}
  region_id: {type: character_current_region}

objectives_template:
  - kind: kill
    target_npc_family: "{target_mob_family}"
    count: "{kill_count}"

reward_template:
  xp: "{kill_count * 20}"
  gold_copper: "{kill_count * 50}"
  currency_delta:
    dynamic_token: "{floor(kill_count/5)}"

time_limit_ms: 14400000          # 4 hour
tags: [dynamic, bounty]
```

## 4. Runtime Instance

Server (bounty board interaction):
- Karakter board'a gelir.
- Server generates dynamic quest instance:
  - Parameter substitute (random within range).
  - `progression.quest_progress` insert (custom quest instance with generated code).
- Karakter kabul → normal tracking.

## 5. Bounty Board

City bounty board:
- 3-5 bounty offered per refresh.
- Refresh: 4 hour.
- Character accept: 1 slot at a time (v1).

## 6. Rare Mob Spawn Quest

Rare mob spawn detected (bkz. [hunting-adventure/rare-enemy.md](../hunting-adventure/rare-enemy.md)):
- Auto-generate quest instance: "Defeat <rare_mob>".
- Reward: rare loot + achievement.
- Time-limited (rare mob despawn).

## 7. Character-Generated Quest (Phase 02+)

- Player creates own quest (kill X for Y gold reward).
- Vendor board post.
- Other player accepts.
- Escrow + fulfillment.

v1'de yok.

## 8. Dynamic Quest Storage

Runtime instance:
```sql
-- Uses progression.quest_progress with custom code
character_id, quest_id (references dynamic instance),
generated_from_template_id, generated_params_json
```

Ek Phase 22:
```sql
CREATE TABLE progression.dynamic_quest_instance (
  id UUID PRIMARY KEY,
  template_id UUID NOT NULL,
  character_id UUID NOT NULL,
  generated_params_json JSONB NOT NULL,
  generated_objectives_json JSONB NOT NULL,
  generated_rewards_json JSONB NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## 9. Endpoint

**GET `/v1/characters/{id}/quests/bounty-board`**
Server generates on-demand:
```json
{
  "bounties": [
    {
      "instance_id": "...",
      "template_code": "dynamic_bounty_template",
      "generated_display": "Kill 15 beasts",
      "reward_summary": "300 XP + 750 gold + 3 dynamic tokens",
      "expires_at": "..."
    }
  ]
}
```

**POST `/v1/characters/{id}/quests/{instance_id}/accept`** — accept dynamic quest.

## 10. AFK Tick

Dynamic quest objective aynı objective_progress mekaniği ile tick.

## 11. Metrics
- `quest_dynamic_generated_total{template,subcategory}`.
- `quest_dynamic_completed_total`.
- `quest_dynamic_expired_total`.

## 12. Cross-Reference
- [Main](main.md), [Side](side.md)
- [hunting-adventure/rare-enemy.md](../hunting-adventure/rare-enemy.md)
- [hunting-adventure/hunting.md](../hunting-adventure/hunting.md)

## 13. Bekleyen Kararlar
- Dynamic quest complexity depth Phase 02.
- Character-generated quest (player bounty) Phase 02+.
- Cross-region bounty (travel-based Phase 02).
- Dynamic reward scaling final tune.
