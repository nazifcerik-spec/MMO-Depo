# Stats (Derived + Tertiary)

## 1. Amaç
Primary stat'lardan türetilen ve doğrudan gameplay formüllerine giren tüm sayısal değişkenler. Bu doküman **formüller** ve **cap** disiplinine odaklanır. Primary stat yönetimi: [attributes.md](attributes.md).

## 2. Kategori Özet
- **Primary** (6): STR, DEX, INT, WIS, STA, CHA — bkz. [attributes.md](attributes.md).
- **Derived** — primary'den türetilen combat/vital.
- **Tertiary / Utility** — meslek, ekonomi, konfor.

## 3. Derived Stats (v1 baseline)

### 3.1 Combat Base
| Stat | Formül | Not |
|---|---|---|
| **Attack Power (AP)** | `2 * STR + 1 * DEX + weapon_AP` | Warrior/Rogue/Ranger melee |
| **Ranged AP** | `2 * DEX + 1 * STR + weapon_RangedAP` | Ranger |
| **Spell Power (SP)** | `2 * INT + 1 * WIS + weapon_SP` | Mage/Cleric/Druid |
| **Heal Power (HP+)** | `1.5 * WIS + 0.5 * INT + weapon_HealPow` | Cleric |

### 3.2 Vital
| Stat | Formül | Not |
|---|---|---|
| **Max HP** | `50 + 10 * STA + 5 * level + gear_bonus` | Universal |
| **Max Mana** | `20 + 8 * INT + 3 * level + gear_bonus` | Caster |
| **Max Rage/Focus/Energy/Nature** | Class-specific pool | Bkz. [energy.md](energy.md), [stamina.md](stamina.md) |
| **HP Regen** | `1 + 0.1 * STA` (per 2 sn, out-of-combat) | Combat'ta ×0.2 |
| **Mana Regen** | `1 + 0.2 * WIS` (per 2 sn) | Combat'ta ×0.5 |

### 3.3 Precision & Speed
| Stat | Formül | Cap |
|---|---|---|
| **Crit Chance** | `5 + 0.15 * DEX + gear_crit%` | Hard cap 75% (diminishing after 50%) |
| **Crit Damage Multiplier** | `1.5 + gear_critdmg%` | Hard cap 2.5× |
| **Haste %** | `gear_haste%` (no primary source) | Hard cap 40% |
| **Cooldown Reduction (CDR) %** | `gear_cdr%` | Hard cap 30% |
| **Accuracy %** | `95 + level_diff_bonus + gear_hit%` | Cap 100% |
| **Evasion %** | `2 + 0.1 * DEX + gear_evasion%` | Soft cap 30% (diminishing) |

### 3.4 Defense
| Stat | Formül | Not |
|---|---|---|
| **Armor** | `equipment_armor_sum + STA * 0.5` | Physical damage reduction |
| **Magic Resist** | `equipment_mr_sum + WIS * 0.3` | Magic damage reduction |
| **Block Chance %** | `shield_base + gear_block%` (shield varsa) | Cap 40% |
| **Parry Chance %** | `weapon_parry_base + gear_parry%` | Cap 30% |

### 3.5 Special
| Stat | Formül |
|---|---|
| **Life Steal %** | `gear_lifesteal%` |
| **Spell Vamp %** | `gear_spellvamp%` |
| **Threat Modifier %** | Tank spec passive |

### 3.6 Damage Reduction Formula
```
armor_dr(dmg, armor, attacker_level) = armor / (armor + K * attacker_level)
  where K = 60 (v1 tune)

Örnek: armor=300, atk_lv=20 → 300 / (300 + 1200) = 0.20 → %20 DR

Similar for magic_resist with mr.
```

Cap DR: **75%** hard cap.

## 4. Tertiary / Utility Stats

| Stat | Formül | Cap |
|---|---|---|
| **Movement Speed %** | `base + gear + mount + buff` | Soft cap +50% (out-of-combat), +25% (combat) |
| **Carry Capacity** | `40 + 2 * STR + bag_bonus` (slot count) | Bkz. [inventory.md](../game-design/inventory.md) |
| **Gathering Yield %** | `gear + race_affinity + mastery` | No hard cap |
| **Crafting Success %** | `85 + 0.5 * WIS + tool_bonus + rank_bonus` | Cap 100% |
| **XP Gain %** | Sum of bonuses (rested, race, buff, event) | Total cap ×3.0 |
| **Gold Find %** | `gear_gf% + luck_potion` | No cap |
| **Rarity Bonus %** | `gear_mf% + luck` | Effective cap %50 |
| **AFK Efficiency %** | Kozmetik ödül; **no combat power impact** | Cap +%20 |

## 5. Derived Cache

Storage: `character.stats.derived_cache_json`.

Yapı:
```json
{
  "_v": 2,
  "computed_at": "2026-09-08T12:00:00Z",
  "effective_primary": { "str": 45, "dex": 20, "int_": 15, "wis": 15, "sta": 32, "cha": 10 },
  "combat": {
    "attack_power": 122,
    "spell_power": 45,
    "heal_power": 30,
    "crit_chance_pct": 8.0,
    "crit_dmg_mult": 1.55,
    "haste_pct": 5.0,
    "cdr_pct": 8.0,
    "accuracy_pct": 100.0,
    "evasion_pct": 4.0
  },
  "vital": {
    "max_hp": 640,
    "max_mana": 155,
    "max_resource": 100,
    "resource_kind": "rage",
    "hp_regen_per_2s": 4.2,
    "mana_regen_per_2s": 4.0
  },
  "defense": {
    "armor": 316,
    "magic_resist": 68,
    "block_pct": 15.0,
    "parry_pct": 8.0
  },
  "tertiary": {
    "move_speed_pct": 100,
    "carry_capacity_slots": 50,
    "gather_yield_pct": 5,
    "craft_success_pct": 92,
    "xp_gain_pct": 108,
    "gold_find_pct": 3,
    "rarity_bonus_pct": 2
  }
}
```

Version alanı (`_v`) formül değişikliği ile artar; migration için işaretçi.

## 6. Cache Invalidation

Tetikleyiciler (bkz. [attributes.md §8](attributes.md)):
- Primary allocation / respec.
- Equip / unequip.
- Buff apply/expire.
- Ascension spend.
- Mastery milestone.
- Level up.

Invalidate: `derived_cache_invalidated_at = now()`.
Read path: eğer `derived_cache_invalidated_at IS NOT NULL` → recompute + `_at = NULL`.

Concurrency: recompute optimistic (version check); race harmless (aynı sonuç).

## 7. Recompute Cost
- Tam recompute: <1ms domain hesap; DB'ye 1 update.
- Batch: pipeline sırasında birden fazla mutation → tek recompute (defer flush).

## 8. Formül Config'te
- Katsayılar (2, 1, 60 (armor K), regen tick sürekleri) `config.stats.yaml` (data-driven).
- Content patch ile revize; version bump.

## 9. Buff/Debuff Modeli

```
character.stats_buff table:
  character_id, buff_id, magnitude, source (item/skill/food),
  applied_at, expires_at, stacks
```

Buff apply:
- Insert row.
- Cache invalidate.

Buff expire (worker or check on read):
- Row delete + cache invalidate.

Stacking:
- Same buff refresh (duration reset).
- Different buff same category: max 2 stack (config).

## 10. Fail Modları
- Cache corruption (schema drift) → recompute fresh.
- Formül overflow (int64 taşma) → clamp + warn.
- Missing gear reference → 0 contribution + warn log.

## 11. Endpoint

**GET `/v1/characters/{id}/stats`**
Auth + CharacterScope.

Response 200: `derived_cache_json` (fresh).

Cache invalidated iken: recompute + return (blocking, <1ms).

## 12. Event Yayını
- `character.stat_recomputed` (info; not hot event).

## 13. Cross-Platform
Server tek gerçek; client cache invalidate → refetch.

## 14. Test Stratejisi
- Table-driven: (primary, gear, buff) → expected derived.
- Cap enforce: crit > 75% → clamp.
- Race: concurrent recompute → idempotent.
- Migration: `_v` bump → old cache invalidated on read.

## 15. Metrics
- `stat_recompute_total` counter.
- `stat_recompute_duration_seconds` histogram.
- `stat_cache_hit_ratio` gauge.

## 16. İlişkiler
- [Attributes](attributes.md), [Health](health.md), [Mana](mana.md), [Energy](energy.md), [Stamina](stamina.md), [game-design/stats.md](../game-design/stats.md), [game-design/combat.md](../game-design/combat.md), [database/character.md](../database/character.md)

## 17. Bekleyen Kararlar
- Elemental resistance sistemi (fire/ice/nature/holy/shadow) — Phase 02.
- PvP damage split formülü (Phase 02).
- Combat rating (level scaling) sistemi — Phase 02.
- Diminishing return eğrileri final (crit, evasion) simülasyon sonrası.
