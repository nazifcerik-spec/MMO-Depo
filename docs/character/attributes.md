# Attributes (Primary Stats)

## 1. Amaç
Karakterin altı ana stat'ının (STR/DEX/INT/WIS/STA/CHA) kaynak, allocation ve respec akışı. Derived stat'lar bu değerlere dayanır (bkz. [stats.md](stats.md)).

## 2. Model
Kaynak: [database/character.md §2](../database/character.md), [game-design/stats.md](../game-design/stats.md).

Tablo: `character.stats`.

Alanlar:
- Primary: `str`, `dex`, `int_`, `wis`, `sta`, `cha`.
- `unspent_points` — level up'lardan biriken atanmamış puan.
- `derived_cache_json` — hesaplanmış derived stat'lar (cache).
- `derived_cache_invalidated_at`.
- `version` — optimistic lock.

## 3. Primary Stats Rolleri
| Kod | Ad | Etki (özet) |
|---|---|---|
| `str` | Strength | Melee AP, taşıma kapasitesi |
| `dex` | Dexterity | Ranged AP, evasion, crit chance |
| `int_` | Intelligence | Spell Power, mana pool |
| `wis` | Wisdom | Mana regen, healing power, meslek verim |
| `sta` | Stamina | HP pool, physical resist |
| `cha` | Charisma | Vendor fiyat, quest reward %, party buff scale |

## 4. Kaynak Akışı

### 4.1 Base Değerler
- Character oluşturulduğunda:
  ```
  primary[i] = 10 (universal base)
        + race.stat_modifiers[i]
        + class.base_stat_bias[i]
  ```
- Bu değerler `character.stats` satırına insert edilir.

### 4.2 Level Up Puanları
- Her level up → `unspent_points += 5`.
- Oyuncu istediği stat'a dağıtır.

### 4.3 Ascension Puanları
- Ascension'lar `character.ascension_choice` tarafından; belirli node'lar `+1 STR` vs verir.
- Bu passive bonus derived_cache hesabında toplanır (**base değere yazılmaz**).

### 4.4 Item & Buff & Set Bonus
- Ekipman `content.item_instance.rolled_stats` + affix + set bonus.
- Buff/debuff geçici modifier.
- Bunlar da **base'e eklenmez**; derived hesabında toplanır.

## 5. Effective Stat Hesabı

```
effective[stat] = base_from_row
                + racial_passive_bonus
                + class_passive_bonus
                + spec_passive_bonus
                + ascension_bonus
                + equipment_sum (rolled + affix + enchant)
                + set_bonus_sum
                + active_buffs_sum
                - active_debuffs_sum
```

Sonuç `>= 0` (floor 0).

Cache alanı: `character.stats.derived_cache_json.effective_primary`.

## 6. Allocation Endpoint

**POST `/v1/characters/{id}/allocate-stats`**
Auth + CharacterScope.

Request:
```json
{
  "spend": { "str": 3, "sta": 2 }
}
```

Server:
1. Load `character.stats FOR UPDATE`.
2. Validate:
   - Her value >= 0.
   - Toplam <= `unspent_points`.
   - Sınırsız cap yok; ama tek stat'a >1000 spend warn (UX koruma).
3. Update:
   - Base primary += spend.
   - `unspent_points -= sum(spend)`.
   - `derived_cache_invalidated_at = now()`.
   - `version += 1`.
4. Recompute derived cache (async ok; ilk okumada gerçekleşir).
5. Publish event `character.stat_allocated`.
6. Audit event.

Response 200:
```json
{
  "stats": { "primary": {...}, "unspent_points": 0 },
  "version": 12
}
```

## 7. Respec (Stat)

**POST `/v1/characters/{id}/respec-stats`**
Auth + CharacterScope.

Server:
1. Cost check:
   - **İlk respec ücretsiz** (öğrenme fırsatı).
   - Sonrası: `cost = 500 * pow(2, respec_count)` gold, cap 20000 gold.
   - Cooldown: 24 saat.
2. Transaction:
   - Debit gold via ledger.
   - Reset primary'yi base (level bazlı):
     ```
     new_base[i] = 10 + race.stat_modifiers[i] + class.base_stat_bias[i]
     unspent_points = 5 * (character.level - 1)  // all level up points back
     ```
   - Publish event `character.stat_respec`.
   - Invalidate derived cache.
3. Response 200: new stats + unspent_points.

## 8. Derived Cache Invalidation

Trigger noktalar:
- Level up (yeni pool sayıları).
- Allocation.
- Respec.
- Equip / unequip.
- Buff apply / expire.
- Ascension node spend.
- Content patch (base formul değişimi).

Invalidate = `derived_cache_invalidated_at = now()`.
Next stat read → recompute + store.

## 9. Cap ve Diminishing Returns
- Primary stat cap: **yumuşak sonsuz** (sadece kaynak sınırlı: 5/level × 50 = 250 + ascension).
- Derived stat'lardaki cap (crit %, CDR %) — bkz. [stats.md](stats.md).

## 10. Race + Class Total Bütçe
- Design constraint: race + class base bias toplamı ±5 stat aralığında (bkz. [race](../game-design/races.md) DD-open-4).
- Character creation validate.

## 11. Concurrency
- Optimistic lock (`version`).
- Concurrent allocate → biri success, biri 409 `character.conflict_retry`.
- Client retry pattern.

## 12. Fail Modları
- Insufficient points → 400 `stats.insufficient_points`.
- Negative spend → 400.
- Version mismatch → 409, retry.
- Insufficient funds (respec) → 402 `economy.insufficient_funds`.

## 13. Rate Limit
- Allocate: 30/hour.
- Respec: 5/day (cooldown ayrıca zorlar).

## 14. Metrics
- `character_stat_allocated_total{stat}`.
- `character_stat_respec_total`.
- `character_stat_avg_distribution{class,stat}` — meta-analysis.

## 15. Audit
- `character.stat_allocated`.
- `character.stat_respec`.

## 16. Event Payload
```json
{
  "_v": 1,
  "character_id": "c-...",
  "delta": { "str": 3, "sta": 2 },
  "new_totals": { "str": 45, "dex": 20, "int_": 15, "wis": 15, "sta": 32, "cha": 10 }
}
```

Subscribers:
- `character.stat_cache_invalidate`.
- `combat.rotation_reeval` (rotation cache invalidate if any).

## 17. Cross-Platform
Aynı character sheet her platformda; stat sheet WS push (versiyon değiştiğinde diff).

## 18. Test Stratejisi
- Unit: base + race + class summation.
- Effective stat computation matrix.
- Respec cost curve.
- Race: concurrent allocate → version conflict.

## 19. İlişkiler
- [Level](level.md), [Stats](stats.md), [Health](health.md), [Mana](mana.md), [game-design/stats.md](../game-design/stats.md), [database/character.md](../database/character.md), [game-design/races.md](../game-design/races.md), [game-design/classes.md](../game-design/classes.md)

## 20. Bekleyen Kararlar
- Respec cost formul final.
- Free respec token drop (Phase 02).
- Rebalance token (content patch sonrası).
- Auto-allocate assistant (class-suggested build) — Phase 02 opsiyonel.
