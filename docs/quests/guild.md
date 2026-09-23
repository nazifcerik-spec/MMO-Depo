# Guild (Guild Quest)

> ⚠️ **Phase 02** — schema hazır (`content.quest_template.category='guild'`); v1'de basit.

## 1. Amaç
Guild-shared quest — tüm guild üyeleri katkı verir. Guild progression + guild reward.

## 2. Karakteristik

- **Guild-shared**: contribution guild-wide.
- **Objectives multi-member**: kill 500 dungeon boss (guild-wide sum).
- **Reward**: guild XP + guild bank deposit + member-özel small reward.
- **Weekly / long-term**: 7-14 gün duration.

## 3. Content Format (Phase 02 design)

```yaml
# content/quests/guild_weekly_dungeon.yaml
code: guild_weekly_dungeon
category: guild
display_name_key: quest.guild.weekly_dungeon.name

giver_npc_ref: guild_officer          # opsiyonel; guild master accepts
turn_in_npc_ref: guild_officer

prerequisites:
  guild_min_level: 5

objectives:
  - kind: kill_boss
    boss_code: ashen_lich_king
    count: 20                         # guild total
  - kind: gather
    target_item_code: shadow_essence
    count: 100

rewards:
  guild_xp: 5000
  guild_bank_deposit_gold_copper: 100000
  per_member_reward:
    xp: 1000
    gold_copper: 500
    currency_delta:
      guild_seal: 5

repeatable: true
cooldown_ms: 604800000                # weekly
time_limit_ms: 604800000
tags: [guild, weekly, boss]
```

## 4. Guild Progression

Guild XP:
- Guild level (bkz. [game-design/guilds.md](../game-design/guilds.md)).
- Level 1..cap.
- Guild bank tab unlock.

## 5. Contribution Ledger

`progression.guild_quest_progress` (Phase 02):
```sql
CREATE TABLE progression.guild_quest_progress (
  guild_id UUID NOT NULL,
  quest_code VARCHAR(48) NOT NULL,
  objectives_progress_json JSONB NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  expires_at TIMESTAMPTZ NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'active',
  PRIMARY KEY (guild_id, quest_code)
);

CREATE TABLE progression.guild_quest_contribution (
  guild_id UUID NOT NULL,
  quest_code VARCHAR(48) NOT NULL,
  character_id UUID NOT NULL,
  contribution_amount INT NOT NULL DEFAULT 0,
  PRIMARY KEY (guild_id, quest_code, character_id)
);
```

## 6. Accept

Guild master/officer accepts quest for guild.
- Contribution tracking automatic for members.
- Complete when total objective met.

## 7. Reward Distribution

- **Guild-wide**: guild XP + bank deposit.
- **Per-member**: XP + gold + guild seal token (participation-based).
- Contributing members prioritized (contribution ≥ threshold → full reward; below → half).

## 8. Guild Seal

Currency `guild_seal`:
- Guild vendor exclusive.
- Guild mount/tabard/hall decoration.
- Long-term prestige currency.

## 9. Endpoint (Phase 02)

**GET `/v1/guilds/{id}/quests`** — active + eligible guild quests.

**POST `/v1/guilds/{id}/quests/{code}/accept`** — officer accept.

**GET `/v1/guilds/{id}/quests/{code}/contribution`** — member contribution ledger.

## 10. AFK Tick

Member aktivite yield → guild quest objective tick:
- Kill boss → guild boss kill count.
- Gather → guild pool total.

Auto-progression AFK-friendly.

## 11. Metrics (Phase 02)
- `quest_guild_completed_total{guild,quest}`.
- `guild_contribution_avg_per_member`.
- `guild_seal_gained_total`.

## 12. Cross-Reference
- [game-design/guilds.md](../game-design/guilds.md)
- [database/social-guilds-parties.md](../database/social-guilds-parties.md)
- [Weekly](weekly.md)

## 13. Bekleyen Kararlar
- Aktivasyon Phase 02.
- Contribution threshold formulü.
- Guild seal vendor content depth.
- Guild-only questline (main story for guilds) Phase 03.
