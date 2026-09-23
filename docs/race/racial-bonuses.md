# Racial Bonuses (Profession Affinity + Unique Mechanics)

## 1. Amaç
Race'in stat/skill dışındaki **konfor / niş** bonusları: profession affinity, keşif kolaylıkları, benzersiz mekanikler (dwarf dark_vision, merfolk water_travel).

## 2. Kategori
1. **Profession Affinity** — bir meslek dalında %5-8 verim bonusu.
2. **Exploration Bonus** — travel/discovery konforları.
3. **Unique Mechanic** — race'e özgü mekanik hokka (kimlik).

## 3. Data Format
Profession affinity ve unique mechanic race satırında:
```yaml
# content/races/dwarf.yaml (özet)
profession_affinity:
  profession_code: mining
  bonus_pct: 8

unique_mechanic:
  code: dwarven_dark_vision
  params:
    dungeon_tile_bonus: 1
    description_key: race.dwarf.unique.desc
```

Server schema:
- `content.race.profession_affinity_id` (FK).
- `content.race.profession_affinity_bonus_pct`.
- `content.race.unique_mechanic_code` (VARCHAR).
- `content.race.unique_mechanic_params_json`.

## 4. Profession Affinity

### 4.1 Etki
- İlgili meslek AFK yield hesabına eklenir (bkz. [systems/afk-engine.md](../systems/afk-engine.md)).
- Skill point / XP gain (opsiyonel Phase 02) — v1'de sadece yield.
- Effect DSL üzerinden pattern:
  ```
  { kind: flat_bonus, target: gather_yield_pct, amount: 8, scope: profession:mining }
  ```
  (Aslında bu passive ile aynı sistemde uygulanır; buranın avantajı: yalın form + explicit content field.)

### 4.2 v1 Race Affinity Tablosu

| Race | Profession | Bonus |
|---|---|---|
| `human` | Trading (Phase 02 profession) | +5% vendor sell price |
| `dwarf` | Mining | +8% yield |
| `elf` | Woodcutting | +8% yield |
| `orc` | Hunting (skinning yield) | +8% yield |
| `merfolk` | Fishing | +8% yield |
| `arcanborn` | Enchanting | +8% success chance |

Not: `trading` v1'de resmi meslek değil; human'ın affinity'si `vendor_sell_price_pct` olarak modellenir:
```yaml
profession_affinity: null
unique_mechanic:
  code: human_market_savvy
  params:
    vendor_sell_pct: 5
```

Bu esneklik: affinity nullable, unique mechanic ile ifade et.

## 5. Exploration Bonuses

Race'e göre keşif konforu:

### 5.1 Water Traversal (Merfolk)
- Su travel süresi -%50.
- Diving stamina consumption -%100.
- Underwater breathing süresiz.

### 5.2 Mountain / Cave (Dwarf)
- Mining node yield ekstra rare drop şansı +%2.
- Dark cave "vision" bonus (dungeon tile radius +1).
- Fall damage reduction -%25 (Phase 02).

### 5.3 Forest (Elf)
- Forest biome hunt spot mob spawn interval -%10.
- Herb node discover chance +%10 (foraging).
- Silent movement (Phase 02 stealth mekaniği).

### 5.4 Plains / Open (Orc)
- Kill streak XP bonus (+%2 per consecutive kill, max 10).
- Mount speed +5% (Phase 02).

### 5.5 Urban (Human)
- City vendor price -%5.
- Guild reputation gain +%10 (Phase 02).
- Faction reputation gain +%5.

### 5.6 Magical (Arcanborn)
- Enchant fail rate -%20.
- Rare enchant material drop +%5.
- Magical mob damage taken -%5.

## 6. Unique Mechanic Implementation

Kod tarafı: `unique_mechanic_code` switch (kod'un tanıdığı sabit set):

```go
package racemechanic

// Kayıtlı handler'lar. Yeni race unique için handler eklenir.
var handlers = map[string]Handler{
    "human_versatility_flavor":   handleHumanVersatility,
    "dwarven_dark_vision":        handleDarkVision,
    "elven_forest_stride":        handleForestStride,
    "orcish_battle_hardened":     handleBattleHardened,
    "merfolk_tide_born":          handleTideBorn,
    "arcanborn_arcane_intuition": handleArcaneIntuition,
    "human_market_savvy":         handleMarketSavvy,
}

type Handler interface {
    // Aktivite yield hesaplanırken çağrılır
    ModifyYield(ctx context.Context, activityKind string, params map[string]any, yield *YieldReport)
    // Combat context (damage/resist)
    ModifyCombat(ctx context.Context, evt *CombatEvent)
    // Exploration (travel time)
    ModifyTravel(ctx context.Context, route *TravelRoute)
    // Vendor / market
    ModifyVendor(ctx context.Context, tx *VendorTx)
}
```

Handler yoksa: `noop` — logu warn, deploy engellemez.

### 6.1 Data-Driven Görünüm, Kod-Tarafı Kayıt
- **Race YAML** yeni race ekleyebilir (stat + passive + active + affinity).
- **Unique mechanic** için: code yaml'da referans; ancak handler kod'da olmalı.
- Yeni "custom mechanic" tipi eklemek isteniyorsa: kod PR gerekir. **Handler ekleme** trade-off; alternatif ise generic effect DSL'in genişlemesi.

## 7. Cross-Cutting Points

Unique mechanic çağrı noktaları:
- AFK Engine reveal (yield modifier).
- Combat tick (damage/heal modifier).
- Travel activity start (route duration modifier).
- Vendor transaction (price modifier).
- Discovery event (drop chance modifier).

Her nokta:
```go
if raceMechCode := char.Race.UniqueMechanic; raceMechCode != "" {
    if h, ok := racemechanic.handlers[raceMechCode]; ok {
        h.ModifyYield(ctx, activity, char.Race.UniqueMechanicParams, &yield)
    }
}
```

## 8. Config

```yaml
race_bonuses:
  affinity_default_bonus_pct: 5
  unique_mechanic_registry_check_on_load: true
```

## 9. Test Stratejisi
- Handler unit test her mekaniği için (mock yield/damage/vendor).
- Loader validation: unknown mechanic code → build fail (staging).
- Cross-race noop: race'in mechanic'i başka race'i etkilemez.

## 10. Metrics
- `race_affinity_yield_bonus_applied_total{race,profession}`.
- `race_unique_mechanic_active_count{race}`.
- `race_exploration_bonus_used_total{race,kind}`.

## 11. Balance
- Toplam race avantajı (passive + active + affinity + unique) toplu meta-impact %5-10 aralığı.
- Herhangi bir race "en iyi seçim" olmamalı; farklı build hedefleri için farklı en iyi race.

## 12. Cross-Platform
Tüm hesap server'da; her platform aynı bonus.

## 13. Fail Modları
- Unknown mechanic code → noop + warn (soft fail).
- Handler panic → recover + audit signal.

## 14. Audit
- `race.affinity_applied` (severity=debug, high volume).
- `race.unique_mechanic_invoked` (debug).

## 15. Content Update Flow
1. YAML: yeni race, affinity, unique mechanic code.
2. Loader validate.
3. Handler var mı check → yoksa CI fail.
4. Test suite genişletme.
5. Staging deploy → simülasyon.
6. Production.

## 16. İlişkiler
- [Race](race.md), [Passive](passive.md), [Racial Skills](racial-skills.md), [Base Stats](base-stats.md)
- [systems/afk-engine.md](../systems/afk-engine.md), [systems/reward-system.md](../systems/reward-system.md)
- [game-design/professions.md](../game-design/professions.md), [game-design/economy.md](../game-design/economy.md)
- [game-design/travel.md](../game-design/travel.md), [game-design/world.md](../game-design/world.md)

## 17. Bekleyen Kararlar
- Unique mechanic tam data-driven (effect DSL'in genişlemesi) vs handler kod modeli — v1 handler kayıt.
- Trading profession (Human affinity'nin ana ev'i) v1'de mevcut mesleklerin bir alt kategorisi mi Phase 02'de yeni meslek mi.
- Race exploration bonus'unun UI görünürlüğü (character sheet detay panel).
- Cross-race hybrid (Phase 03) mechanic sharing.
