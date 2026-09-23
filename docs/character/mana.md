# Mana

## 1. Amaç
Caster class'ların (Mage, Cleric, Druid partial) birincil kaynağı. Doldurma-tabanlı büyük havuz.

## 2. Kim Kullanır
- **Mage** — birincil.
- **Cleric** — birincil.
- **Druid** — Nature kaynağı kullanır ancak bazı heal skill'leri mana tüketir (dual).

Warrior/Ranger/Rogue → mana kullanmaz.

## 3. Model
- **Max Mana**: derived stat.
- **Current Mana**: runtime state (Redis + snapshot).

## 4. Formüller

### 4.1 Max Mana
```
max_mana = 20 + 8 * effective_INT + 3 * level + equipment_mana_bonus
```
Örnek Mage lv 25, INT 60, gear +60:
```
max_mana = 20 + 8*60 + 3*25 + 60 = 635
```

### 4.2 Regen
```
mana_regen_per_2s = 1 + 0.2 * effective_WIS
combat_regen_multiplier = 0.5  -- mana combat'ta %50
```

### 4.3 Rest / Meditation (Phase 02)
- Karakter idle → +100% regen bonus 15 sn sonra.

## 5. Kullanım
- Skill cast: `resource_cost` (bkz. [content.skill_template](../database/content-templates.md)).
- Insufficient mana → skill reject (`skill.insufficient_resource`).
- Cast time süresince mana ayrılmaz; cast complete anında düşer.
- Interrupt cast (hasar alma → cast_break) mana geri alınır (partial reject).

## 6. Regen Path
- Server tick 2 sn (birleştirilebilir with vitals tick).
- Cache Redis: `char:mana:<char_id>`.
- Snapshot Postgres: `current_mana_snapshot`.

## 7. Rotation Interaction
AFK rotation:
- Priority list mana check condition: `{"cond": {"mana_below_pct": 30}, "skill_id": "mana_drink"}`.
- Consumable (mana potion) rotation'a eklenebilir.

## 8. Endpoint
`GET /v1/characters/{id}/vitals` (bkz. [health.md §9.1](health.md)) → mana blok içerir.

## 9. Config
```yaml
mana:
  base_max: 20
  int_per_mana: 8
  level_per_mana: 3
  regen_base: 1.0
  regen_per_wis: 0.2
  combat_regen_mult: 0.5
  rest_bonus_after_sec: 15
  rest_regen_mult: 2.0
```

## 10. Buff / Item
- Mana potion consumable: instant +N mana; cooldown 1.5 sn (short — combat use).
- Elixir (long buff): +30% mana regen 30 dk.
- Gear `mana_regen_per_2s` bonus additive.

## 11. Fail Modları
- Insufficient mana on cast → skill reject; UI feedback.
- Overflow: `current_mana` clamp to `max_mana`.

## 12. WS Push
`vitals_update` frame içinde mana delta.

## 13. Metrics
- `character_mana_deficit_pct` gauge (histogram of 1 - current/max).
- `skill_reject_mana_total`.

## 14. Cross-Class Skills Using Mana
- Druid heal spells (`Rejuvenation`).
- Cleric holy damage (rare).

## 15. Test Stratejisi
- Cast → mana deduct → regen tick.
- OOM (out-of-mana) reject.
- Rotation mana condition.

## 16. İlişkiler
- [Stats](stats.md), [Health](health.md), [Energy](energy.md), [game-design/skills.md](../game-design/skills.md), [game-design/combat.md](../game-design/combat.md)

## 17. Bekleyen Kararlar
- Meditation stance mekaniği (Phase 02).
- Mana burn debuff (Phase 02).
- Cross-class mana borrowing (druid dual mode) detay.
