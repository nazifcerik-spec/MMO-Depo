# War (Guild War)

> ⚠️ **Phase 02+** — territory control + guild vs guild combat.

## 1. Amaç
Guild vs guild competition — territory ownership, resource control, PvP.

## 2. Karakteristik

- **Declared war** — attacker declares against defender.
- **Duration**: 7-14 gün.
- **Territory prize** — winner claims territory.
- **Faction/alliance tie** (Phase 03).

## 3. Model (Phase 02+)

```sql
CREATE TABLE social.guild_war (
  id UUID PRIMARY KEY,
  attacker_guild_id UUID NOT NULL,
  defender_guild_id UUID NOT NULL,
  territory_id UUID NULL,           -- contested territory
  declared_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ends_at TIMESTAMPTZ NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'active',  -- active/attacker_won/defender_won/draw/surrendered
  attacker_score INT NOT NULL DEFAULT 0,
  defender_score INT NOT NULL DEFAULT 0
);

CREATE TABLE social.guild_war_event (
  id BIGSERIAL PK,
  war_id UUID NOT NULL,
  event_kind VARCHAR(24) NOT NULL,  -- kill/capture/defend
  actor_guild_id UUID NOT NULL,
  actor_character_id UUID NULL,
  score_delta INT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## 4. War Declaration Flow

**POST `/v1/guilds/{id}/war/declare`**
Leader only.
```json
{
  "target_guild_id": "...",
  "contested_territory_code": "wildwood_border"
}
```

Server:
1. Validate leader.
2. Validate guild level ≥ 15.
3. Validate cost (500k gold treasury).
4. Insert guild_war.
5. Notify defender + broadcast server-wide.

## 5. War Mechanics

- **PvP zone** — territory tagged war-active.
- **Kill credit**: enemy guild member kill → war_event.
- **Capture point** — objective control.
- **Defense** — repel attackers.

Score accumulate; end of duration → higher score wins.

## 6. Reward / Penalty

Winner:
- Territory ownership (`world.territory.controlling_guild_id` update).
- Guild XP bonus.
- Season leaderboard rank.
- Cosmetic (war banner, tabard border).

Loser:
- Territory loss.
- Guild XP penalty (minor).
- No permanent damage.

## 7. Surrender

Leader can surrender:
- Set state='surrendered'.
- Half-penalty; ends war.

## 8. Alliance (Phase 03)

Multiple guilds ally:
- Alliance-wide war.
- Shared territory pool.

## 9. Balance (Phase 02+)

- Guild size similarity check (opsiyonel — anti-bullying).
- Weekly war cap per guild.
- Cross-region war logistics.

## 10. Endpoint (Phase 02+)

**POST `/v1/guilds/{id}/war/declare`**
**POST `/v1/guilds/{id}/war/surrender`**
**GET `/v1/guilds/{id}/war/active`**
**GET `/v1/wars/{id}/leaderboard`**

## 11. Cross-Reference
- [Guild](guild.md), [Raid](raid.md), [Boss](boss.md)
- [world/territory.md](../world/territory.md) — territory ownership
- PvP System (Phase 02+ ayrı doküman)

## 12. Bekleyen Kararlar
- Aktivasyon Phase 02+.
- Score formula (kill/capture/defend weight).
- Alliance war Phase 03.
- Cross-realm war Phase 03.
- Balance (bullying protection) Phase 02.
