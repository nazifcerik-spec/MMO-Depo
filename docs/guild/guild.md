# Guild (Entity + Lifecycle)

## 1. Amaç
Guild entity — create, lifecycle, disband, name/tag.

## 2. Model

Kaynak: [database/social-guilds-parties.md §1](../database/social-guilds-parties.md).

```sql
social.guild
  id, name (CITEXT UNIQUE), tag (CITEXT UNIQUE)
  leader_character_id
  motd, description
  level, xp, treasury_gold_copper
  home_city_id
  member_cap
  created_at
  disband_scheduled_at, disbanded_at
```

Ek Phase 24:
```sql
ALTER TABLE social.guild
  ADD COLUMN logo_ref VARCHAR(128) NULL,
  ADD COLUMN tabard_json JSONB NOT NULL DEFAULT '{}',
  ADD COLUMN faction_id UUID NULL,                    -- guild faction alignment Phase 02
  ADD COLUMN alliance_id UUID NULL,                    -- Phase 03
  ADD COLUMN founded_by_user_id UUID NULL,             -- audit
  ADD COLUMN version INT NOT NULL DEFAULT 1;
```

## 3. Create Endpoint

**POST `/v1/guilds`**
Auth + CharacterScope.

Request:
```json
{
  "name": "Emerald Wardens",
  "tag": "EWD",
  "description": "Casual guild for adventurers"
}
```

Server:
1. Validate character not in guild.
2. Validate name unique (CITEXT) 3-24 char.
3. Validate tag unique 2-5 char.
4. Cost check: 10000 gold (100000 copper) creation fee.
5. Transaction:
   - Debit gold via ledger.
   - Insert social.guild (leader = char).
   - Insert social.guild_rank (5 default ranks; bkz. [roles.md](roles.md)).
   - Insert social.guild_member (leader rank 0).
   - Emit event `guild.created`.
   - Audit.
6. Response 201.

## 4. Disband Endpoint

**POST `/v1/guilds/{id}/disband`**
Leader only.

Server:
1. Validate leader.
2. Set `disband_scheduled_at = now() + interval '48 hour'`.
3. Notification all members.
4. Grace period 48h — leader cancel possible.

Worker (48h later):
- If not cancelled: `disbanded_at = now()`.
- Distribute treasury (rank-based split; v1: leader takes all).
- Return bank items to members (mail per member; v1: skipped, leader claim).
- Delete members.
- Emit event `guild.disbanded`.

**POST `/v1/guilds/{id}/disband/cancel`** — leader within 48h.

## 5. Guild Detail

**GET `/v1/guilds/{id}`**
```json
{
  "id": "...",
  "name": "Emerald Wardens",
  "tag": "EWD",
  "leader_character_id": "...",
  "motd": "Welcome!",
  "level": 6,
  "xp": 12500,
  "treasury_gold_copper": 5000000,
  "member_count": 32,
  "member_cap": 50,
  "created_at": "...",
  "home_city_id": "...",
  "logo_ref": "guilds/emerald.png"
}
```

## 6. MOTD Update

**POST `/v1/guilds/{id}/motd`**
Leader/officer.
```json
{ "motd": "Weekly guild raid Friday 8pm UTC" }
```

## 7. Guild Rename (Phase 02)

- Cost: 50000 gold.
- Cooldown: 30 gün.
- Name unique enforce.
- Old name reservation 90 gün.

v1'de yok.

## 8. Home City

Guild home city:
- Guild hall (Phase 02) located there.
- Guild members auto-shown "you are home".
- Change: leader only; 30 gün cooldown.

## 9. Metrics
- `guild_created_total`.
- `guild_disbanded_total`.
- `guild_active_count` gauge.
- `guild_avg_member_count` gauge.

## 10. Audit
- `guild.created`, `guild.disband_scheduled`, `guild.disbanded`, `guild.motd_updated`.

## 11. Cross-Reference
- [Members](members.md), [Roles](roles.md), [Bank](bank.md), [Level](level.md)
- [database/social-guilds-parties.md](../database/social-guilds-parties.md)
- [game-design/guilds.md](../game-design/guilds.md)

## 12. Bekleyen Kararlar
- Create fee (10000 gold v1; playtest).
- Disband grace 48h vs 24h.
- Treasury distribution formula on disband (v1 leader; Phase 02 rank-based split).
- Alliance Phase 03.
