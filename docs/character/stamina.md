# Stamina

## 1. Amaç
Uzun süreli fiziksel efor sınırlayıcısı. **v1'de opsiyonel** (bkz. [database/activity.md §8](../database/activity.md)); ancak bu doküman resmi tasarımı sunar — Phase 02'de aktive edilecek.

**v1 aktif mi?** — HAYIR. Bu doküman "hazır iskelet". Aktivasyon kararı Phase 02 ekonomi/gameplay simülasyonu sonrası.

## 2. Neden Var
- AFK hunting'in "sonsuz farm" hissini yumuşatmak.
- Aktif oyun ve dinlenme arasında ritim yaratmak.
- Consumable (food/drink) ekonomisine anlam vermek.

## 3. Neden v1'de Kapalı
- Idle MMO'da stamina "wait wall" gibi hissedilebilir.
- Ana pillar "persistent progression" (DD-003) ile çelişme riski.
- Balance simülasyonu Phase 02'de yapılacak; aktivasyon o zamana.

## 4. Model (Aktif Edildiğinde)
- **Max Stamina**: 100 (base), gear/mastery ile artabilir.
- **Current Stamina**: runtime state.
- **Recovery Rate**: idle 1/dk; food/potion boost.

Schema (Phase 02 aktif):
```sql
ALTER TABLE character.character
  ADD COLUMN current_stamina SMALLINT NOT NULL DEFAULT 100,
  ADD COLUMN max_stamina SMALLINT NOT NULL DEFAULT 100,
  ADD COLUMN stamina_last_recovery_at TIMESTAMPTZ NOT NULL DEFAULT now();
```

## 5. Kullanım Kuralları

### 5.1 Aktivite Konsumsiyonu
| Aktivite | Stamina/dk |
|---|---|
| Mining/Woodcutting/Foraging | 1 |
| Hunting | 2 |
| Fishing | 0.5 |
| Adventure | 1.5 |
| Travel (walking) | 0.5 |
| Dungeon idle-run | 3 |
| Manuel combat | 4 |
| Rest (city) | 0 (+recovery) |
| Crafting | 0 |

### 5.2 Depletion
- Stamina 0'a düştüğünde: activity `paused`, `pause_reason='stamina'`.
- Notification "Stamina depleted".

### 5.3 Recovery
- Idle out-of-activity: 1/dk baseline.
- Rest in city inn: 5/dk.
- Food consumable: +30 stamina instant.
- Drink consumable: +10/dk 5 dk boost.

## 6. Endpoint

**GET `/v1/characters/{id}/stamina`**
```json
{
  "current": 62,
  "max": 100,
  "recovery_per_min": 1,
  "expected_recovery_to_max_min": 38
}
```

## 7. AFK Interaction (Phase 02 aktif hali)
Reveal sırasında:
```
elapsed_min = 60
consumption = activity_rate * elapsed_min
if consumption > current_stamina:
    active_min = current_stamina / activity_rate
    yield only for active_min; then paused
    current_stamina = 0
else:
    current_stamina -= consumption
```

## 8. Consumables
- **Bread / Meal**: `+30 stamina` instant; cooldown 60 sn; food category.
- **Ration** (crafted): `+50 stamina` + 5 dk regen buff.
- **Rare Feast** (guild event): `+100` capped.

## 9. Config (Phase 02)
```yaml
stamina:
  base_max: 100
  base_recovery_per_min: 1.0
  rest_recovery_per_min: 5.0
  consumption_per_activity:
    mining: 1.0
    hunting: 2.0
    ...
```

## 10. Balance Considerations
- Aktive edildiğinde: normal aktif oyuncu için 8 saatlik AFK "yeterli" ancak sonsuz değil.
- Casual player deneyimini kırmamalı; food temini kolay.
- Sink olarak ekonomiye katkı (food talep).

## 11. Kararlar Bekliyor
- v2'de aktive et vs asla — GDD §2.1 ("Persistent Progression" pillar) ile çelişki riski dikkatli test edilmeli.
- Alternatif model: soft cap only (24h AFK sonrası verim azal, stamina depleted değil sadece yield düşer). Öneri: **bu alternatif tercih**.
- Yani bu doküman muhtemelen "hazır ama etkinleştirilmeyecek" olarak kalır; soft cap sistemi zaten [AFK Engine §Soft Cap](../systems/afk-engine.md) tarafından çözer.

## 12. Cross-Reference
- [AFK Engine §Soft Cap](../systems/afk-engine.md) — soft cap alternative
- [Activity Engine](../systems/activity-engine.md) — pause_reason path
- [database/activity.md §8](../database/activity.md) — schema hazırlığı

## 13. İlişkiler
- [Health](health.md), [Energy](energy.md), [Activity Engine](../systems/activity-engine.md), [game-design/hunting.md](../game-design/hunting.md), [game-design/professions.md](../game-design/professions.md)

## 14. Bekleyen Kararlar (Phase 02)
- Aktive et mi kalıcı olarak devre dışı bırak mı.
- Aktif olursa consumption rate tune.
- Food economy denge testi.
- "Rest bonus" pasif fikri (Phase 02).
