# Mastery (Skill + Weapon + Track)

## 1. Amaç
"Kullandıkça artar" tarzı ilerleme. Skill rank up (use-based), weapon mastery, mastery track ilerleme. Yatay derinlik ve uzun soluklu progression.

Yüksek seviye: [game-design/mastery.md](../game-design/mastery.md). Bu doküman engine implementation.

## 2. İki Katman

### 2.1 Skill Rank (Use-Based)
- Her skill'in `max_rank` (v1'de default 1; Phase 02'de 5).
- Skill kullanıldıkça `character.skill_learn.xp` artar.
- XP threshold'ta rank up.

### 2.2 Mastery Track
- `content.mastery_track` (weapon/combat/profession/region/bestiary/collection).
- `character.mastery.level`, `.xp`.
- Milestone bazlı ödüller (passive modifier).

## 3. Skill Rank Progression

XP formulu:
```
rank_xp_required(rank) = 100 * pow(2, rank-1)
```
| Rank | XP need | Cumulative |
|---|---|---|
| 1 → 2 | 100 | 100 |
| 2 → 3 | 200 | 300 |
| 3 → 4 | 400 | 700 |
| 4 → 5 | 800 | 1500 |

XP grant kaynakları:
- Skill cast success: `+1 XP` (baseline).
- Skill kill credit contribution: `+2 XP`.
- Passive proc: `+0.5 XP` (small, rounded up).

Cap: per skill max 100 XP/dk (bot koruma).

## 4. Skill Rank Effect

Rank up etkisi:
- **Effect scale**: base amount × rank_multiplier.
  ```
  rank_multiplier(r) = 1 + 0.15 * (r-1)   → rank 5 = 1.6× baz
  ```
- Alternatif: rank başına spesifik override (`content.skill_rank.effects_json_override`).

## 5. Skill Rank Engine

```go
func (m *MasteryEngine) GrantSkillXP(ctx context.Context, charID uuid.UUID, skillCode string, xp int32) error {
    return m.pool.WithTx(ctx, ..., func(tx db.TX) error {
        row, err := m.repo.WithTx(tx).SkillLearnFor(ctx, charID, skillCode)
        if err != nil { return err }
        row.XP += xp
        rankedUp := 0
        for {
            need := skillRankXPRequired(row.Rank)
            if row.XP < need { break }
            row.XP -= need
            row.Rank++
            rankedUp++
            if row.Rank >= skillMaxRank(skillCode) { break }
        }
        if err := m.repo.WithTx(tx).UpdateSkillLearn(ctx, row); err != nil { return err }
        if rankedUp > 0 {
            m.bus.Publish(ctx, bus.Event{
                Kind: "skill.ranked_up",
                Payload: {char_id, skill_code, old_rank, new_rank},
            })
        }
        return nil
    })
}
```

## 6. Mastery Track Yapısı

Kaynak: [database/content-templates.md §27](../database/content-templates.md), [database/character.md §12](../database/character.md).

Track'ler:
- **Weapon**: sword, axe, staff, bow, dagger, ...
- **Combat Discipline**: tanking, healing, dps.
- **Profession**: mining, blacksmithing, ...
- **Region**: mastery per region visited/played.
- **Bestiary**: mob family kills.
- **Collection**: mount/pet/transmog completion (Phase 02).

Cap: 100 level (v1); Phase 02'de 200.

## 7. Mastery Track XP Kaynakları

| Track | XP kaynağı |
|---|---|
| Weapon | Skill cast with that weapon equipped |
| Combat Discipline (tank) | Damage taken while spec_role=tank |
| Combat Discipline (dps) | Damage dealt to hostile |
| Combat Discipline (heal) | Effective healing |
| Profession | Craft/gather success |
| Region | Time spent in region (per minute tick) + kills |
| Bestiary | Kill count per family |

Formul örneği (weapon):
```
weapon_xp_per_cast = 1 + skill.cast_time_ms / 1000
```

## 8. Mastery Milestones

Level başına belirli milestone'lar (config-driven):
```yaml
# content/mastery_tracks/weapon_sword.yaml
milestones:
  - level: 10
    reward:
      modifiers:
        - { kind: mult_add, target: "damage_pct.weapon.sword", amount: 2 }
      title_code: sword_apprentice
  - level: 25
    reward:
      modifiers:
        - { kind: mult_add, target: "damage_pct.weapon.sword", amount: 4 }
        - { kind: flat, target: "crit_chance_pct", amount: 1 }
  - level: 50
    reward:
      title_code: sword_veteran
      unlock_skill: sword_perfect_strike
  - level: 100
    reward:
      title_code: sword_grandmaster
      unlock_signature: true
```

Milestone claim:
- Otomatik (level up ile) veya opt-in (UI'da claim button — Phase 02).
- `character.mastery.claimed_milestones` array kaydı.

## 9. Mastery Milestone Modifier Apply

Milestone'daki modifier'lar `SourceKind='mastery'`, `SourceID=track_id` olarak modifier stack'e katılır (bkz. [modifiers.md](modifiers.md)).

Cache invalidate: milestone claim → stat cache invalidate.

## 10. Endpoint

**GET `/v1/characters/{id}/mastery`**
Response:
```json
{
  "skills": [
    { "code": "berserker_rampage", "rank": 3, "xp": 350, "next_rank_xp": 800 }
  ],
  "tracks": [
    {
      "code": "weapon_sword", "category": "weapon",
      "level": 42, "xp": 12000, "next_level_xp": 15000,
      "claimed_milestones": [10, 25],
      "next_milestone_at": 50
    }
  ]
}
```

**POST `/v1/characters/{id}/mastery/{track_code}/claim-milestone`**
Body: `{"level": 25}`.
- Milestone eligible ise apply; ledger + audit.

## 11. Anti-Cheat

- Per skill hourly XP cap.
- Per track daily XP cap.
- Rapid cast + no cost skill (macro) → cap normalizes; combat contribution filter.

## 12. AFK Simülasyon Uyumu

- AFK combat/gathering yield sırasında skill cast + mastery XP granted.
- Deterministic seed; reveal içinde `mastery_xp_grants` özeti.

## 13. Fail Modları

- Unknown skill code → warn; skip.
- Milestone content updated post-claim → apply new (Phase 02 rework compensation).
- Concurrent claim race → optimistic (transaction).

## 14. Test Stratejisi

- Table-driven rank XP curve.
- Multi-rank jump.
- Milestone claim idempotency.
- Cap enforcement.

## 15. Metrics

- `skill_ranked_up_total{skill_code}`.
- `mastery_leveled_up_total{track_code}`.
- `mastery_milestone_claimed_total{track,level}`.
- `mastery_xp_per_hour{track_code}`.

## 16. Cross-References
- [Skills](skills.md), [Active Skills](active-skills.md), [Modifiers](modifiers.md), [Stats](stats.md)
- [game-design/mastery.md](../game-design/mastery.md), [character/xp.md](../character/xp.md)
- [database/character.md §12](../database/character.md), [database/content-templates.md §27](../database/content-templates.md)

## 17. Bekleyen Kararlar
- Skill rank v1'de aktif mi (`max_rank=1` sabit) — öneri: **evet aktif ama max_rank=1 (v1)**; content Phase 02'de rank 5'e çıkarır.
- Mastery track cap final (100 vs 200).
- Milestone auto-claim vs opt-in.
- Cross-character account-wide mastery (Phase 02).
- Weapon mastery calculation: weapon-specific damage bonus vs general — CI cap.
