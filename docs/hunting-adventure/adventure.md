# Adventure

## 1. Amaç
Story + stage-based aktivite. Hunt'un yerine "lore + narrative" alternatif. AFK-uyumlu; login'de her stage için narration reveal.

Foundation: [game-design/adventure.md](../game-design/adventure.md), [activity-engine/adventure.md](../activity-engine/adventure.md — sadece hunting varsa, ekleme).

Not: `activity-engine/adventure.md` mevcut değil ise (v1 dokümantasyon başlangıcında), bu doküman adventure activity'sinin core detayını üstlenir.

## 2. Params

```json
{
  "adventure_template_id": "...",
  "party_id": null,
  "difficulty": "normal"
}
```

## 3. Duration Buckets

`content.adventure_template.duration_bucket`:
- `short` — 15 min.
- `mid` — 1 hour.
- `long` — 4 hours.
- `epic` (Phase 02) — 8 hours.

Duration total = duration_bucket.

## 4. Stage Model

`content.adventure_stage`:
```sql
adventure_stage
  id, adventure_id, order_index
  check_type (combat/skill/resource/luck)
  requirement_json
  success_effect_json
  failure_effect_json
  narration_key
```

Adventure N stage; each stage bir "check" (combat/skill/resource/luck).

## 5. Check Türleri

| Check Type | Requirement | Örnek |
|---|---|---|
| `combat` | DPS >= threshold | "Defeat forest guardian" |
| `skill` | Skill level >= X | "Decipher ancient runes (INT >= 40)" |
| `resource` | Item consume | "Pay 500 gold to bribe guard" |
| `luck` | RNG check | "Traverse dangerous path (%70 chance)" |

Requirement JSONB per check type.

## 6. Yield Compute

Reveal (analytical):
```go
elapsed = now - contract.started_at
total_stages = len(adventure.Stages)
stages_completed = min(total_stages, int(elapsed / (adventure.DurationMs / total_stages)))

results := []StageResult{}
for i := 0; i < stages_completed; i++ {
    stage := adventure.Stages[i]
    // Deterministic check (snapshot + seed)
    passed := evaluateCheck(stage, snapshot, seed.Sub(fmt.Sprintf("stage-%d", i)))
    results = append(results, StageResult{Stage: stage, Passed: passed})
}

// Reward: reward_pool × completion_ratio × success_ratio
completion_ratio = stages_completed / total_stages
success_ratio = successCount / stages_completed
finalReward = adventure.RewardPool × completion_ratio × (0.5 + success_ratio × 0.5)
```

## 7. Check Deterministic Evaluation

Combat check:
```
if snapshot.DPS >= stage.MinDPS: passed = true
```

Skill check:
```
if snapshot.StatValue(stage.RequiredStat) >= stage.MinStatValue: passed = true
```

Resource check:
```
if character has stage.RequiredItem >= stage.RequiredQty:
    consume; passed = true
else: passed = false
```

Luck check:
```
if seed.NextFloat() < stage.SuccessChance: passed = true
```

## 8. Failure Handling

Fail stage:
- Alternate branch (Phase 02 opsiyonel).
- Partial yield (fail effect_json).
- Adventure continues to next stage (soft fail v1).
- Cumulative fails: 3/5 fail → adventure "poor completion" (reward %50).

## 9. Narration Reveal

Login sonrası reveal panel:
```json
{
  "adventure_completed": true,
  "adventure_code": "ancient_grove_whispers",
  "stages": [
    { "narration_key": "adv.grove.stage1.enter", "passed": true },
    { "narration_key": "adv.grove.stage2.idol", "passed": true },
    { "narration_key": "adv.grove.stage3.guardian_fight", "passed": true, "note": "close call!" },
    { "narration_key": "adv.grove.stage4.runes", "passed": false, "note": "failed to decipher" },
    { "narration_key": "adv.grove.stage5.relic", "passed": true }
  ],
  "success_ratio": 0.8,
  "reward_summary": { ... }
}
```

Client narrative story flow.

## 10. Content Örneği

```yaml
# content/adventures/ancient_grove_whispers.yaml
code: ancient_grove_whispers
region_code: ashen_hills
tier: 20
duration_bucket: mid           # 1 hour

stages:
  - order: 1
    check_type: skill
    requirement:
      stat: dex
      min: 25
    narration_key: adv.grove.stage1
  - order: 2
    check_type: luck
    requirement:
      chance_pct: 70
    narration_key: adv.grove.stage2
  - order: 3
    check_type: combat
    requirement:
      min_dps: 150
    narration_key: adv.grove.stage3
  - order: 4
    check_type: skill
    requirement:
      stat: int
      min: 30
    narration_key: adv.grove.stage4
  - order: 5
    check_type: resource
    requirement:
      item_code: ancient_key
      qty: 1
    narration_key: adv.grove.stage5

reward_pool_code: ancient_grove_reward
faction_code: emerald_wardens
```

## 11. Party Adventure

Party adventure:
- Party leader'ın snapshot check'e girer (v1 basit).
- Phase 02: party avg check.
- Loot personal.

## 12. Endpoint

**POST `/v1/characters/{id}/adventure`** — start.
**GET `/v1/characters/{id}/adventure/current`** — progress.
**POST `/v1/characters/{id}/adventure/{cid}/cancel`** — partial yield.

## 13. Adventure Chain (Phase 02)

- Adventure A completion → Adventure B unlocks.
- Story arc; long-form questline.

v1'de bağımsız adventure'lar.

## 14. Metrics
- `adventure_completed_total{code}`.
- `adventure_stage_success_ratio_avg{code}`.
- `adventure_partial_completion_ratio` gauge.

## 15. Cross-Reference
- [game-design/adventure.md](../game-design/adventure.md), [game-design/quests.md](../game-design/quests.md)
- [Hunting](hunting.md), [Encounter](encounter.md), [Exploration](exploration.md)
- [database/content-templates.md §16](../database/content-templates.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 16. Bekleyen Kararlar
- Branching narrative Phase 02.
- Adventure chain (story arc) Phase 02.
- Coop-only adventures Phase 02.
- Seasonal adventure line Phase 02.
- Adventure repeatable cooldown (v1: rest window 24h).
