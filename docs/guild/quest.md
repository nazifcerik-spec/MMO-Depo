# Quest (Guild Quest)

> ⚠️ **Phase 02**. Bkz. [quests/guild.md](../quests/guild.md) tam detay.

## 1. Amaç
Guild-shared quest — üyeler ortak katkı; guild XP + guild bank deposit.

## 2. Karakteristik

- Weekly cycle (Salı 06:00 UTC reset).
- Officer accepts on behalf.
- All members contribute.
- Reward: guild XP + guild bank + per-member reward.

## 3. Endpoint (Phase 02)

**POST `/v1/guilds/{id}/quests/{code}/accept`** — officer.
**GET `/v1/guilds/{id}/quests`** — active + eligible.
**GET `/v1/guilds/{id}/quests/{code}/contribution`** — per-member ledger.

## 4. Model

Bkz. [quests/guild.md §5-6](../quests/guild.md):
```sql
progression.guild_quest_progress
  guild_id, quest_code, objectives_progress_json, started_at, expires_at, status

progression.guild_quest_contribution
  guild_id, quest_code, character_id, contribution_amount
```

## 5. Reward Distribution

- Guild XP + guild bank deposit (aggregate).
- Per-member reward (contribution threshold ≥ N → full; below → half).

## 6. Cross-Reference
- [quests/guild.md](../quests/guild.md) — tam detay
- [Guild](guild.md), [XP](xp.md), [Bank](bank.md)

## 7. Bekleyen Kararlar
- Aktivasyon Phase 02.
- Contribution threshold formulü.
- Cross-guild quest (alliance) Phase 03.
