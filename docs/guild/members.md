# Members (Guild Roster)

## 1. Amaç
Guild üye yönetimi — invite, join, leave, kick, contribution tracking.

## 2. Model

Kaynak: [database/social-guilds-parties.md §3](../database/social-guilds-parties.md).

```sql
social.guild_member
  guild_id, character_id (PK)
  rank_index (0=leader, ..., 4=initiate)
  joined_at
  contribution_points
  last_active_at
  note_private, note_officer
  UNIQUE (character_id) — bir karakter tek guild
```

## 3. Invite Flow

`social.guild_invite`:
```sql
guild_invite
  id, guild_id, target_character_id
  invited_by_character_id
  status (pending/accepted/declined/expired)
  expires_at (7 gün)
  created_at
```

## 4. Endpoints

### 4.1 Invite

**POST `/v1/guilds/{id}/invite`**
Officer+ (rank ≤ 1).
```json
{ "target_character_name": "Fenrir" }
```

Server:
1. Validate inviter permission.
2. Validate guild not full.
3. Validate target not in another guild.
4. Insert guild_invite (expires_at = now + 7 gün).
5. WS push + notification target.

### 4.2 Accept

**POST `/v1/guilds/{id}/join`**
Target character.

Server:
1. Validate invite pending + not expired.
2. Validate guild size < cap.
3. Insert social.guild_member (rank 4 initiate default).
4. Delete invite.
5. Emit event `guild.member_joined`.
6. Guild chat broadcast.

### 4.3 Leave

**POST `/v1/guilds/{id}/leave`**
Server:
1. Delete guild_member row.
2. If leader → auto-promote (bkz. [roles.md](roles.md)).
3. Emit event `guild.member_left`.

### 4.4 Kick

**POST `/v1/guilds/{id}/kick`**
Officer+ only.
```json
{ "target_character_id": "..." }
```

Server:
1. Validate permission.
2. Validate target rank < kicker rank.
3. Delete member.
4. Notification target.

## 5. Contribution Points

`contribution_points`:
- Guild XP verilen aktivite → contribution artar.
- Guild bank deposit gold/item → contribution.
- Guild quest completion → contribution.
- Weekly reset (opsiyonel Phase 02).

Contribution → guild rank promotion criteria (Phase 02).

## 6. Guild Cap

Base: 50 (v1).
Guild level artışı ile:
- Level 5: +5 (55).
- Level 10: +10 (60).
- Level 25: +25 (75).
- Level 50: +50 (100 cap).

## 7. Member Note

- `note_private` — üye'nin kendi notu (128 char).
- `note_officer` — officer notu (256 char, sadece officer görür).

## 8. Auto-Inactive

Member 60 gün inaktif:
- `last_active_at` check.
- Officer notification.
- Auto-kick (Phase 02 opsiyonel; v1 manuel).

## 9. Cross-Character Consistency

Character.guild membership unique (single guild).
Alt karakter aynı guild'e katılabilir mi? Karar:
- **v1**: farklı karakter farklı guild OK.
- **Phase 02**: aynı account cross-character guild membership hint.

## 10. Metrics
- `guild_member_joined_total`.
- `guild_member_left_total{cause}`.
- `guild_member_contribution_avg` gauge.

## 11. Audit
- `guild.member_joined`, `guild.member_left`, `guild.member_kicked`, `guild.member_promoted`.

## 12. Cross-Reference
- [Guild](guild.md), [Roles](roles.md), [Permissions](permissions.md), [XP](xp.md)
- [database/social-guilds-parties.md §3](../database/social-guilds-parties.md)

## 13. Bekleyen Kararlar
- Auto-kick inactivity 60d Phase 02.
- Contribution weekly reset Phase 02.
- Cross-character guild account link Phase 02.
- Recruit officer level Phase 02 (mid-rank permission).
