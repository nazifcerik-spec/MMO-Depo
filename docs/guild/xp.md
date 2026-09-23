# XP (Guild XP)

## 1. Amaç
Guild XP kaynakları + grant flow.

## 2. XP Kaynakları

| Kaynak | XP |
|---|---|
| Member level up | +50 per level |
| Member quest completion (main) | +100 |
| Member quest completion (side/daily) | +25 |
| Member weekly quest | +200 |
| Member dungeon completion | +100 |
| Member world boss kill credit | +250 |
| Guild bank deposit (gold) | +1 per 1000 copper (max 100/day per member) |
| Guild quest completion | +1000 |
| Guild boss defeat (Phase 02) | +2500 |
| Achievement (guild-scoped) | +500 |

## 3. Grant Flow

Member aktivite tetiklendiğinde:
```go
func onCharacterActivity(char, activity, result) {
    if char.GuildID != nil {
        xp := computeGuildXP(activity, result)
        guildService.GrantXP(char.GuildID, xp, source=activity, char=char)
        // Track contribution
        updateContributionPoints(char.GuildID, char.ID, xp)
    }
}
```

## 4. Contribution Tracking

Bkz. [members.md §5](members.md).

`social.guild_member.contribution_points` accumulate.

## 5. Per-Member XP Cap

Prevent bot abuse:
- Per-member per-day XP contribution cap: 5000 (v1).
- Exceeded → capped (silent); log warn.

## 6. Weekly Contribution Reset (Phase 02)

- Weekly counter reset Salı 06:00 UTC (DD-016).
- Guild rank promotion criteria (contribution ≥ X weekly).

v1'de no reset; accumulate lifetime.

## 7. XP Aggregate Reveal

Character reveal içinde:
- Activity result → guild XP grant.
- Reveal payload includes `guild_xp_contributed`.

## 8. Endpoint

**GET `/v1/guilds/{id}/xp-log`** — recent XP entries with source + contributor.

## 9. Metrics
- `guild_xp_granted_total{source}`.
- `guild_xp_per_member_avg` gauge.
- `guild_contribution_cap_hit_total`.

## 10. Cross-Reference
- [Level](level.md), [Members](members.md), [Quest](quest.md)
- [database/social-guilds-parties.md §1, 3](../database/social-guilds-parties.md)

## 11. Bekleyen Kararlar
- XP source balance final Playtest.
- Contribution reset Phase 02.
- Anti-bot cap 5000/day final.
- Cross-guild alliance XP share Phase 03.
