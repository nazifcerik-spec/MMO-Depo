# Level (Guild Level)

## 1. Amaç
Guild level cap + benefit progression.

## 2. Guild Level Cap

- Base: 1.
- Cap: **10** v1; Phase 02 cap 20 (raid/war unlock).

## 3. Level Benefits

Per level unlocks:
- **Level 1**: base 50 member cap, bank 1 tab, chat.
- **Level 3**: bank tab 2 unlock; guild logo customization.
- **Level 5**: member cap 55; bank tab 3.
- **Level 7**: bank tab 4; guild research unlock (Phase 02).
- **Level 10**: member cap 60; guild boss unlock (Phase 02); tabard customization.
- **Level 15** (Phase 02): member cap 75; raid unlock.
- **Level 20** (Phase 02): member cap 100; guild war eligibility.

## 4. XP Formula

Bkz. [xp.md](xp.md).

```
guild_xp_to_next(L) = base × pow(factor, L-1)
base = 5000
factor = 1.30
```

Level 10 cumulative ~50k XP.
Level 20 (Phase 02) ~500k XP.

## 5. Level Up Flow

XP threshold reached (via aggregate contributions):
```go
if guild.XP >= xpToNext(guild.Level) {
    guild.XP -= xpToNext(guild.Level)
    guild.Level++
    applyLevelUpBenefits(guild)
    emit event 'guild.leveled_up'
    guild-wide notification
}
```

## 6. Benefits Apply

Level up → member cap update; tab unlock; feature unlock.

Members receive notification:
- "Guild leveled up to X — new bank tab available!"

## 7. Reverse Progression

Guild XP loss (member kick, disband period):
- v1: no XP loss on member leave.
- Phase 02: contribution decay if member leaves.

Level down: **YOK**. Cap sonrası XP birikimi kaybolur (level cap reached).

## 8. Level Cap Behavior

Guild max level → XP overflow = 0.
Phase 02: prestige system (raid boss defeats → guild "veteran" tier).

## 9. Metrics
- `guild_leveled_up_total{guild,level}`.
- `guild_avg_level` gauge.
- `guild_reached_cap_total`.

## 10. Cross-Reference
- [Guild](guild.md), [XP](xp.md), [Research](research.md), [Bank](bank.md)
- [database/social-guilds-parties.md §1](../database/social-guilds-parties.md)

## 11. Bekleyen Kararlar
- Cap 10 vs 20 v1 (öneri: 10; Phase 02 expand).
- XP curve tune (base 5000).
- Prestige system Phase 02+.
- Cross-guild level (alliance) Phase 03.
