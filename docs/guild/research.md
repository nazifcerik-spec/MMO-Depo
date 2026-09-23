# Research (Guild Perks)

> ⚠️ **Phase 02** — schema hazır; v1'de kapalı.

## 1. Amaç
Guild-wide passive perks — long-term investment. Yield bonuses, party synergy, unique unlock.

## 2. Model (Phase 02)

```sql
CREATE TABLE content.guild_research_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  category VARCHAR(24) NOT NULL,       -- combat/gathering/social/economic
  max_level SMALLINT NOT NULL DEFAULT 5,
  xp_curve_json JSONB NOT NULL,
  effect_per_level_json JSONB NOT NULL,   -- modifier list per level
  prerequisites_json JSONB NULL,
  cost_gold_per_level_json JSONB NOT NULL,
  version INT DEFAULT 1
);

CREATE TABLE social.guild_research (
  guild_id UUID NOT NULL,
  research_id UUID NOT NULL,
  level SMALLINT NOT NULL DEFAULT 0,
  xp BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (guild_id, research_id)
);
```

## 3. Örnek Research

```yaml
# content/guild_research/master_gatherer.yaml
code: master_gatherer
category: gathering
max_level: 5
xp_curve:
  base: 10000
  factor: 1.5
effect_per_level:
  - level: 1
    modifiers:
      - { kind: mult_add, target: gather_yield_pct, amount: 1 }
  - level: 2
    modifiers:
      - { kind: mult_add, target: gather_yield_pct, amount: 2 }
  - level: 3
    modifiers:
      - { kind: mult_add, target: gather_yield_pct, amount: 3 }
  # ... total +5% at level 5
cost_gold_per_level:
  1: 50000
  2: 150000
  3: 400000
  4: 1000000
  5: 2500000
```

## 4. Research Categories

- **Combat**: +damage pct, +party HP.
- **Gathering**: +yield pct per profession.
- **Social**: +XP gain, +quest reward.
- **Economic**: +vendor discount, +market fee reduction.

## 5. Activation Flow (Phase 02)

**POST `/v1/guilds/{id}/research/{code}/activate`**
Officer+ only.

Server:
1. Validate guild level meets prereq.
2. Debit treasury cost for level 1.
3. Insert social.guild_research (level=1, xp=0).
4. Apply modifier guild-wide.

**POST `/v1/guilds/{id}/research/{code}/contribute`**
Any member:
```json
{ "xp_amount": 100 }
```
- Member contributes XP (from personal contribution).
- Alternate: automatic contribution from member activity aggregate.

Level up when XP threshold met + additional cost.

## 6. Modifier Application

Guild research active → modifier applied to all guild member characters:
- Modifier stack contribution (`source_kind='guild_research'`).
- Cache invalidation guild-wide.

## 7. Research Cap

- v1 (Phase 02 aktif): guild has cap `N active research`; default 3.
- Research tree Phase 03 (dependency tree).

## 8. Endpoint (Phase 02)

**GET `/v1/guilds/{id}/research`**
```json
{
  "active": [
    { "code": "master_gatherer", "level": 3, "xp": 22000, "next_level_cost": 400000 }
  ],
  "available_to_activate": [ ... ]
}
```

## 9. Metrics (Phase 02)
- `guild_research_activated_total`.
- `guild_research_leveled_up_total`.
- `guild_research_avg_active_per_guild` gauge.

## 10. Cross-Reference
- [Guild](guild.md), [XP](xp.md), [Bank](bank.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md)
- [database/social-guilds-parties.md §8](../database/social-guilds-parties.md)

## 11. Bekleyen Kararlar
- Aktivasyon Phase 02.
- Research tree dependency Phase 03.
- Cost tune balance.
- Alliance-shared research Phase 03.
