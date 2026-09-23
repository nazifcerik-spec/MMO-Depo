# Boss (Guild Boss)

> ⚠️ **Phase 02** — guild-summoned boss encounter.

## 1. Amaç
Guild-exclusive boss — guild members summon; guild boss reward.

## 2. Model (Phase 02)

```sql
CREATE TABLE content.guild_boss_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  min_guild_level SMALLINT NOT NULL,
  hp_base BIGINT NOT NULL,
  phase_definitions JSONB NOT NULL,
  enrage_timer_ms INT NOT NULL,
  loot_table_id UUID NOT NULL,
  guild_xp_reward BIGINT NOT NULL,
  guild_treasury_reward_copper BIGINT NOT NULL,
  summon_cost_treasury_copper BIGINT NOT NULL,
  weekly_lockout BOOLEAN NOT NULL DEFAULT true,
  version INT DEFAULT 1
);

CREATE TABLE progression.guild_boss_spawn (
  id UUID PRIMARY KEY,
  guild_id UUID NOT NULL,
  boss_template_id UUID NOT NULL,
  spawned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  defeated_at TIMESTAMPTZ NULL,
  participants JSONB NOT NULL DEFAULT '[]'
);
```

## 3. Summon Flow

**POST `/v1/guilds/{id}/boss/summon`**
Officer+ only.
```json
{ "boss_template_code": "guild_boss_ashen_titan" }
```

Server:
1. Validate guild level ≥ boss.min_guild_level.
2. Validate weekly lockout not active.
3. Debit treasury summon_cost.
4. Insert guild_boss_spawn.
5. WS broadcast guild-wide.
6. Emit event `guild.boss_summoned`.

Guild members within X minutes attend (bkz. world boss pattern).

## 4. Combat

Bkz. [combat/pve-encounters.md §7](../combat/pve-encounters.md):
- Phase-based (multi-phase).
- Enrage timer.
- Contribution tracking (damage + heal + participation).

## 5. Reward

Boss defeated:
- Guild XP + treasury deposit.
- Participants (>10% contribution): personal loot roll from boss table.
- First guild kill server-wide anons + achievement.

## 6. Weekly Lockout

Per guild:
- Boss weekly lockout counter.
- Subsequent summons: rare loot penalty.

## 7. Endpoint (Phase 02)

**GET `/v1/guilds/{id}/boss/active`** — currently spawned bosses.

## 8. Metrics (Phase 02)
- `guild_boss_summoned_total`.
- `guild_boss_defeated_total`.
- `guild_boss_wipe_total`.

## 9. Cross-Reference
- [Guild](guild.md), [XP](xp.md), [Bank](bank.md), [Raid](raid.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)

## 10. Bekleyen Kararlar
- Aktivasyon Phase 02.
- Summon cost balance.
- Boss content depth (v1: 3-5 guild boss).
