# Travel (Mount Travel Integration)

## 1. Amaç
Mount'un travel süresine katkısı — speed multiplier, terrain compatibility, snapshot lock.

Foundation: [travel/mount.md](../travel/mount.md), [travel/base-time.md](../travel/base-time.md).

## 2. Travel Method: Mount

Travel method options:
- `walking` — mount ignore.
- `mount` — active mount speed multiplier used.
- `road` — mount speed × road bonus.
- Others (wagon/ship/waygate) — mount ignore, fixed method.

## 3. Speed Contribution

Bkz. [travel/mount.md](../travel/mount.md):
```
mount_multiplier = mount.EffectiveSpeed × (1 + barding_bonus)
```

Level-scaled:
```
mount.EffectiveSpeed = species.base_speed × (1 + level × 0.02)
```

## 4. Snapshot Lock

Travel contract snapshot includes:
- `mount_species_code`
- `mount_speed_multiplier`
- `mount_type`
- `mount_terrain_compatibility_json`

Mid-travel mount change → duration korunur (DD-095).

## 5. Terrain Compatibility

Bkz. [terrain.md](terrain.md).

Segment-based compute — her edge terrain'e göre mount penalty veya bonus.

## 6. AFK Travel

Reveal:
```go
if travel.MountID != nil {
    mount := loadFromSnapshot(travel.SnapshotID)
    speed_mult := mount.EffectiveSpeed
    // Segment-based; terrain apply
}
```

Analytical duration compute; per-segment iteration OK (küçük N).

## 7. Sprint Mode (Phase 02)

Special active skill "Sprint":
- +50% speed for 15 sec.
- Cost: 20 stamina.
- Cooldown 5 min.

Travel içinde sprint aktive:
- Rare optimization (short burst).

v1'de yok.

## 8. Combat Interaction

Combat state → auto-dismount:
- Redis `char:mount:<char>` clear.
- Cooldown 30 sec re-summon.
- Travel activity: if mid-travel + combat interrupt (Phase 02 ambush) → travel pause + dismount.

## 9. Multi-Player Party Travel

- Each character kendi mount.
- Party leader travel speed → slowest member (Phase 02 party travel).

v1'de her karakter bağımsız travel.

## 10. Waygate Bypass

Waygate teleport mount ignore (instant teleport).

## 11. Metrics
- `mount_travel_used_total{species}`.
- `mount_travel_duration_avg{species}`.

## 12. Cross-Reference
- [Mount](mount.md), [Stats](stats.md), [Terrain](terrain.md), [Skills](skills.md)
- [travel/mount.md](../travel/mount.md), [travel/base-time.md](../travel/base-time.md), [travel/afk-integration.md](../travel/afk-integration.md)

## 13. Bekleyen Kararlar
- Sprint mode Phase 02.
- Multi-passenger mount Phase 03.
- Mount-only region gate (Phase 02 winged/aquatic).
