# Auto-Attack

## 1. Amaç
Karakter combat'a girdiğinde otomatik weapon-based attack; her `attack_speed_ms` interval'ında bir. Skill cast'inden bağımsız continuous DPS baseline.

## 2. Trigger Kural Seti

Auto-attack başlar:
- Character combat state'de (target hostile lock).
- Weapon equipped.
- Not silenced.
- Not casting (channeled skill sırasında pause; instant skill sırasında paralel devam).

Auto-attack durur:
- Combat state exit.
- Target dead/despawn.
- Weapon unequip.
- Sleep/stun/crowd-control.

## 3. Interval Formula

```
attack_speed_ms = weapon.template.attack_speed_ms   (default 2400 ms sword; 3200 ms 2H; 1600 ms dagger)
effective_interval = attack_speed_ms × (1 - haste_pct/100)
min_interval = 500 ms (haste cap %70)
```

Örnek Warrior sword + 20% haste:
```
2400 × (1 - 0.20) = 1920 ms → attack every 1.92 sec
```

## 4. Damage Compute

Auto-attack damage:
```
weapon_damage = uniform(weapon.min_dmg, weapon.max_dmg)   // seed'li roll
scaling = weapon_damage × (1 + primary_stat × 0.005)      // ör. STR based warrior
crit = random < crit_chance_pct
if crit:
    damage = scaling × crit_dmg_mult
else:
    damage = scaling

// Target defense mitigation (bkz. damage.md)
final = applyMitigation(damage, target)
```

Formul detay: [damage.md](damage.md) + [skill-engine/scaling.md](../skill-engine/scaling.md).

## 5. Weapon Kind Behavior

| Weapon Kind | attack_speed_ms | Damage Range | Behavior |
|---|---|---|---|
| **Dagger** | 1600 | Fast + low | High proc rate |
| **Sword (1H)** | 2400 | Balanced | Standard |
| **Axe (1H)** | 2400 | Balanced + crit bias | +2% crit |
| **Mace (1H)** | 2600 | Slow + armor pen | +5% armor pen |
| **Sword (2H)** | 3200 | Slow + high | Big burst |
| **Bow** | 2800 | Ranged | Range: mid |
| **Crossbow** | 3600 | Slow ranged | Big damage per shot |
| **Staff** | 3000 | Ranged + spell scaling | Caster auto |
| **Wand** | 2000 | Fast ranged spell | Caster low burst |
| **Polearm** | 3400 | AoE potential | Cleave (Phase 02) |
| **Fist (unarmed)** | 1800 | Fastest + lowest | Monk stub Phase 03 |
| **Shield** | – | – | No auto (off_hand) |

## 6. Dual Wield

Character dual wield weapon:
- Main hand + off hand alternate attacks (staggered).
- Off hand damage × 0.5 (default).
- Effective DPS ~ 1.5x single weapon.

Class-specific:
- Rogue Duelist spec: off-hand damage × 0.8.
- Warrior Berserker Ravager: off-hand damage × 0.7.

## 7. Ranged Auto-Attack

- Range max: 30 units (arbitrary v1).
- Ammo (Phase 02): arrows/bolts consumable.
- v1'de infinite ammo.

## 8. AFK Auto-Attack (Aggregate)

AFK reveal'de auto-attack **aggregate**:
```
auto_attacks_per_hour = 3600 / (attack_speed_sec)
dps_auto_attack_contribution = auto_attacks_per_hour × avg_damage_per_hit
```

Analytical; per-attack loop yok (DD-084).

## 9. Threat Contribution

Auto-attack threat:
```
threat_generated = damage × threat_modifier   (default 1.0)
```

Tank auto-attack threat modifier higher (bkz. [threat.md](threat.md)).

## 10. Interruption

Cast time skill sırasında auto-attack:
- Instant cast: paralel devam.
- Cast time skill: auto-attack paused; resume post-cast.
- Channeled skill: auto-attack disabled.

## 11. Config

```yaml
auto_attack:
  min_interval_ms: 500
  haste_cap_pct: 70
  off_hand_damage_multiplier: 0.5
  ranged_max_units: 30
  ammo_required: false     # v1 kapalı
```

## 12. Event

Per auto-attack:
```json
{
  "kind": "auto_attack",
  "actor_id": "...",
  "target_id": "...",
  "weapon_kind": "sword",
  "damage": 45,
  "crit": false
}
```

Emit: `combat.auto_attack` (high volume — sample or aggregate).

## 13. Metrics
- `auto_attack_total{weapon}` counter.
- `auto_attack_dps_avg` gauge.
- `auto_attack_crit_ratio` gauge.

## 14. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [Threat](threat.md)
- [item-engine/item-template.md](../item-engine/item-template.md) — weapon
- [skill-engine/scaling.md](../skill-engine/scaling.md) — damage formula
- [character/stats.md](../character/stats.md) — haste

## 15. Bekleyen Kararlar
- Ammo Phase 02 (ranged consumable).
- Polearm cleave AoE Phase 02.
- Dual wield off-hand damage class-based tune final.
- Auto-attack proc gear (on-hit effect Phase 02).
