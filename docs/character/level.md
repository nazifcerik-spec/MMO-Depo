# Level

## 1. Amaç
Karakter seviyesinin nasıl yükseldiği, hangi ödülleri getirdiği, cap'in nerede olduğu ve ascension'la nasıl birleştiği.

## 2. Model
- `character.character.level` (SMALLINT), başlangıç `1`, cap `50` (v1) — Phase 02'de 100 hedefi.
- `character.character.xp` (BIGINT).
- Threshold formulü ve ödül tablosu content data'da (`content.level_reward`).

## 3. XP Threshold Formülü (v1 önerisi)

Geometrik curve, kolay tune edilebilir:
```
xp_to_next(L) = round(BASE * pow(FACTOR, L-1))
BASE   = 100
FACTOR = 1.14
```

Örnek ilk 10:
| Level | XP to next | Cumulative |
|---|---|---|
| 1 | 100 | 100 |
| 2 | 114 | 214 |
| 3 | 130 | 344 |
| 4 | 148 | 492 |
| 5 | 169 | 661 |
| 6 | 193 | 854 |
| 7 | 220 | 1074 |
| 8 | 251 | 1325 |
| 9 | 286 | 1611 |
| 10 | 326 | 1937 |

Level 50 cumulative ~ 190k (ascension #1 hedefi ile örtüşür).

Ascension multiplier:
```
effective_xp_to_next(L, A) = xp_to_next(L) * pow(ASCENSION_MULT, A)
ASCENSION_MULT = 1.05  (DD-011 balance)
```

Ascension #5'te aynı level için ~%28 fazla XP gerekir.

## 4. Level Up Trigger

`character.xp >= xp_to_next(level)` olduğunda:

```
while character.xp >= xp_to_next(character.level):
    character.xp -= xp_to_next(character.level)
    character.level += 1
    grant_level_rewards(character.level)
    emit event 'character.leveled_up'
```

Multi-level jump destekli (bir defada 2+ level).

## 5. Level Rewards

Her level up'ta:
- **Stat point**: 5 (character.stats.unspent_points += 5).
- **Skill point** (Phase 02 talent sistemi): 1 (level 20+).
- **HP/Mana pool** re-derive (stat cache invalidate).
- **Spec unlock**:
  - Level 5 → Spec I unlock notification.
  - Level 10 → Spec II unlock notification.
- **New skill unlock**: content.skill_template `unlock_condition_json.min_level`.
- **Notification**: in-game "Level up!" pop-up.
- **Milestone rewards** (content driven, `content.level_reward`):
  ```
  level | reward_json
    5   | {"item_instance_template": "starter_bag", "gold_copper": 5000, "title_id": "..."}
   10   | {"gold_copper": 20000, "unlock": "spec_tier2"}
   20   | {"item": "adventurer_gear_set", "gold_copper": 50000}
   ...
   50   | {"unlock": "ascension_gate", "title": "Ascendant Aspirant"}
  ```

## 6. Cap Davranışı

Level cap'e ulaşınca:
- Fazla XP **çöpe atılmaz**; `character.xp` üstünde birikir.
- UI "XP → Ascension Progress" gösterir.
- Ascension quest chain unlock (bkz. [game-design/ascension.md](../game-design/ascension.md)).
- Ascension gerçekleştiğinde:
  - `character.level = 1`.
  - `character.xp = 0` (kalan XP kaybolur; ascension ödülü bunu telafi eder).
  - `character.ascension_level += 1`.
  - `character.ascension_points_total += 3` (DD-011).

## 7. Endpoint

Level up ayrı endpoint DEĞİLDİR — otomatik trigger. Ancak kullanıcı UI feedback için:

**GET `/v1/characters/{id}/progression`**
```json
{
  "level": 25,
  "xp": 15423,
  "xp_to_next": 32000,
  "ascension_level": 0,
  "next_milestone": {
    "at_level": 30,
    "reward_preview": "gold 100k + item epic_belt"
  }
}
```

## 8. Server Implementation

```go
func (s *CharacterService) ApplyXP(ctx context.Context, charID uuid.UUID, delta int64, source XPSource) error {
    return s.pool.WithTx(ctx, db.TxOptions{Iso: pgx.ReadCommitted}, func(tx db.TX) error {
        char, err := s.repo.WithTx(tx).ByIDForUpdate(ctx, charID)
        if err != nil { return err }

        newXP := char.XP + delta
        levelUps := 0
        for {
            need := xpToNext(char.Level, char.AscensionLevel)
            if newXP < need { break }
            newXP -= need
            char.Level++
            levelUps++
            if char.Level >= s.cfg.LevelCap {
                // cap'te — kalan XP birikimde tut
                break
            }
        }

        if err := s.repo.WithTx(tx).UpdateLevelXP(ctx, char.ID, char.Level, newXP, char.Version); err != nil {
            return err
        }

        if levelUps > 0 {
            if err := s.grantLevelRewards(ctx, tx, char, levelUps); err != nil { return err }
            s.bus.Publish(ctx, bus.Event{
                Kind: "character.leveled_up",
                SubjectID: char.ID,
                PayloadJSON: mustJSON(LevelUpPayload{
                    OldLevel: char.Level - levelUps,
                    NewLevel: char.Level,
                    StatPointsAwarded: 5 * levelUps,
                }),
            })
        }
        return nil
    })
}
```

Optimistic lock: `UpdateLevelXP` version check; conflict → retry (max 3).

## 9. XP Debt (Death) Interaction
Bkz. [xp.md §Debt](xp.md). Level down **yok**: XP debt sadece `character.xp` from current level başlangıcı üzerinden düşer; level azalmaz. Debt bir sonraki level up'a kadar birikir; sonraki XP grant önce debt öder.

## 10. Event Yayını
Her level up:
```json
{
  "_v": 1,
  "character_id": "c-777",
  "old_level": 24,
  "new_level": 25,
  "stat_points_awarded": 5,
  "unlocks": ["skill:fireball_rank3"]
}
```

Subscriber örnekleri:
- `achievement.leveled_up_milestones` (25/50).
- `notification.leveled_up` — in-game + push.
- `character.stat_cache_invalidate` (HP pool re-derive).

## 11. Fail Modları
- XP grant → char version mismatch → retry.
- Cap seviyesinde xp overflow → sadece birikir; hata değil.
- Ascension trigger için "gate quest not completed" → cap seviyesinde birikir; oyuncu quest'i tamamladığında ascend.

## 12. Cross-Platform Notlar
- Level up bildirimi WS push ile o an bağlı device'lara.
- Achievement/title unlock sonucu global (tüm device'larda aynı).

## 13. Test Stratejisi
- Unit: xpToNext(L, A) formul table.
- Integration: multi-level jump (deltaXP büyük).
- Ascension boundary: cap seviyesinde XP overflow davranışı.
- Race: concurrent ApplyXP → version retry.

## 14. Metrics
- `character_level_up_total{level}` counter (per achieved level).
- `character_at_level` gauge (histogram of active char levels).
- `character_at_ascension` gauge.

## 15. Audit
- `character.leveled_up` (event_log).
- `character.ascended` (ayrı akış).

## 16. İlişkiler
- [XP](xp.md), [Attributes](attributes.md), [Stats](stats.md), [Health](health.md), [game-design/ascension.md](../game-design/ascension.md), [game-design/progression.md](../game-design/progression.md), [database/character.md](../database/character.md), [backend/event-bus.md](../backend/event-bus.md)

## 17. Bekleyen Kararlar
- Level cap final (50 v1 vs 60/80).
- XP curve final BASE + FACTOR (simülasyon sonrası).
- Level milestone reward listesi (content pack).
- Multi-level jump per level up notification: her level ayrı toast mu tek özet mi (öneri: özet).
