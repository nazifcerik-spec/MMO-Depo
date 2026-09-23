# Level (Mount Level Progression)

## 1. Amaç
Mount leveling — speed/cargo/stamina cap artışı, skill unlock.

## 2. Level Cap

- v1: **30** (default; species-özel override).
- Legendary mount: 40.
- Mythic: 50.

## 3. XP Curve

Formula (species field):
```
xp_to_next(L) = base × pow(factor, L-1)
base = 100
factor = 1.10
```

Level 30 cumulative ~15k XP.

## 4. XP Kaynakları

- **Travel active with mount**: 1 XP per unit distance (100 unit walk = 100 XP).
- **Mount aktif combat kill** (mounted combat Phase 03): pending.
- **Feed treat** (Phase 02): +N XP consumable.
- **Quest reward** (mount-specific quest): 500-2000 XP.

Travel primary source; long travel → mount levels naturally.

## 5. Level Up Rewards

Per level:
- Base stat scale (+2% speed, +1 stamina_max).
- Every 10 level: skill rank up eligible.
- Level 30 (cap): unlock unique mount ability (species-özel).

Stat scale example:
```
level_1_speed_multiplier = species.base_speed_multiplier  (e.g. 2.0)
level_30_speed_multiplier = base × 1.5   // +50% at cap
```

## 6. Level Up Flow

```go
func (s *MountService) GrantXP(ctx, mountID, xp) error {
    mount := load(mountID)
    mount.XP += xp
    for {
        need := xpToNext(mount.Level, species.XPCurve)
        if mount.XP < need: break
        mount.XP -= need
        mount.Level++
        if mount.Level >= species.MaxLevel: break
        applyLevelUpRewards(mount)
    }
    save(mount)
    emit event 'mount.leveled_up'
}
```

## 7. AFK Travel Integration

Travel activity reveal içinde:
```go
if travel.ActiveMountID != nil {
    distance_units := travel.RouteDistance
    mount_xp := int64(distance_units)  // 1 XP per unit
    mountService.GrantXP(travel.ActiveMountID, mount_xp)
    revealPayload.mount_xp_gained = mount_xp
}
```

Uzun travel → mount level up otomatik.

## 8. Level Cap Bonus (Rare Ability Unlock)

Level 30 → unique skill unlock:
- Steel Warhorse: "War Charge" — combat pull dismount + damage aura (Phase 03).
- Sea Turtle: "Deep Dive" — %50 aquatic bonus speed.
- Dire Wolf: "Pack Instinct" — party movement speed +5% aura.

Species-özel unique ability content pack.

## 9. Metrics
- `mount_leveled_up_total{species,level}`.
- `mount_reached_cap_total{species}`.
- `mount_xp_per_hour_avg`.

## 10. Cross-Reference
- [Mount](mount.md), [Stats](stats.md), [Skills](skills.md), [Travel](travel.md)
- [character/level.md](../character/level.md) — pattern

## 11. Bekleyen Kararlar
- XP per unit tune (v1: 1; playtest).
- Level cap 30 vs 40 (Phase 02 balance).
- Cap unique ability content depth.
- Multi-mount level share (rare Phase 03).
