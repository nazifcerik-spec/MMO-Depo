# Elite Enemy

## 1. Amaç
Standart mob'dan güçlü ancak boss'dan hafif — party recommendation'lı zorluk katmanı. Group hunt teşvik.

## 2. Elite Model

`content.npc_template.elite_flag = true`:
- HP × 3 (baseline vs regular).
- Damage × 1.5.
- Special abilities (2-3 unique per elite).
- Loot table richer.
- XP reward × 3.

Ek Phase 18:
```sql
ALTER TABLE content.npc_template
  ADD COLUMN elite_tier SMALLINT NOT NULL DEFAULT 0,     -- 0=normal, 1=elite, 2=champion, 3=mini-boss
  ADD COLUMN party_recommended_min SMALLINT NOT NULL DEFAULT 1,
  ADD COLUMN party_recommended_max SMALLINT NOT NULL DEFAULT 5,
  ADD COLUMN elite_group_ids UUID[] NOT NULL DEFAULT '{}';    -- co-spawn group (3 elite bandit)
```

## 3. Elite Tier'lar

| Tier | Ad | HP Mult | Damage Mult | Party Rec | Notlar |
|---|---|---|---|---|---|
| 0 | Normal | 1.0 | 1.0 | Solo | Standard |
| 1 | Elite | 3.0 | 1.5 | 1-2 | Special ability |
| 2 | Champion | 5.0 | 2.0 | 2-3 | Rare drop guaranteed |
| 3 | Mini-boss | 10.0 | 2.5 | 3-5 | Region unique |

## 4. Elite Group

Elite mob'lar bazen "co-spawn" grubu:
- 3 elite bandit birlikte spawn.
- `elite_group_ids` array — grup üyeleri.
- Kill all → group loot bonus.

## 5. Spawn

Monster area (bkz. [world/monster-area.md](../world/monster-area.md)):
- Elite mob spawn rate düşük (~%5 baseline pool).
- Content pool weight: normal 90, elite 10.

Location'a bağlı:
- Hunt spot `danger_tier` = elite → elite oranı yüksek.
- Normal spot elite %5-10 pool weight.

## 6. AFK Hunting Interaction

Reveal hesabı:
- Elite mob kill time ↑ (HP × 3); rate düşer.
- Solo hunter elite spot'ta yield %60 (undertier); party hunt daha iyi.
- Elite kill loot table extra rare drop.

Character DPS check:
```
if avg_time_per_kill > safety_threshold (60 sec elite):
    yield_penalty = 0.6    // solo underpowered
```

## 7. Party Bonus

Party ile elite hunt:
- DPS toplamı ~= party size × individual DPS.
- Party kill time daha kısa.
- Loot: personal loot each; party members eşit rare drop şans.

## 8. Champion / Mini-Boss

Champion: elite'nin güçlü versiyonu. Region'un "boss-lite" mob'u.
- Rare drop guaranteed.
- Achievement kill.
- Weekly cooldown per character (loot penalty).

Mini-boss: region unique; boss encounter-lite.
- Dungeon-outside boss.
- Uses [combat/pve-encounters.md §7 phase system](../combat/pve-encounters.md).
- World boss'tan hafif (no server-wide anons).

## 9. Content Örneği

```yaml
# content/npc_templates/elite_bandit_captain.yaml
code: elite_bandit_captain
family: humanoid
level: 25
hp_base: 3500
damage_base: 220
elite_flag: true
elite_tier: 1
party_recommended_min: 2
party_recommended_max: 3
ai_profile_code: elite_bruiser
abilities_json: ["bandit_slash", "bandit_stun_bash"]
loot_table_code: elite_bandit_loot
xp_reward: 450
elite_group_ids: []
```

Elite group example:
```yaml
# content/npc_templates/elite_bandit_trio.yaml
code: elite_bandit_group
kind: elite_group
group_members:
  - elite_bandit_captain
  - elite_bandit_healer
  - elite_bandit_scout
co_spawn_at_area: forest_bandit_camp
```

## 10. Cross-Reference
- [Hunting](hunting.md), [Rare Enemy](rare-enemy.md), [Encounter](encounter.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [world/monster-area.md](../world/monster-area.md)
- [database/content-templates.md §17](../database/content-templates.md)

## 11. Metrics
- `elite_mob_killed_total{tier}`.
- `elite_group_defeated_total`.
- `champion_kill_weekly_avg`.

## 12. Bekleyen Kararlar
- Elite group co-spawn coordination Phase 02.
- Champion weekly lockout — evet (v1); mini-boss weekly (v1).
- Elite tier balance simulation.
- Elite scaling by party size (Phase 02 dynamic).
