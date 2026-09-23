# Member

## 1. Amaç
Party member state — join, leave, ready, role hint.

## 2. Model

Kaynak: [database/social-guilds-parties.md §10](../database/social-guilds-parties.md).

```sql
social.party_member
  party_id, character_id (PK-part)
  role_hint (nullable)
  joined_at, left_at (nullable)
  ready (bool)
  UNIQUE (character_id) WHERE left_at IS NULL
```

## 3. Join Party

**POST `/v1/parties/{id}/join`**
Server (accepting invite):
1. Validate invite exists + not expired.
2. Validate party size < size_cap.
3. Validate character not in another party.
4. Insert social.party_member row.
5. Delete invite row.
6. Emit event `party.member_joined`.
7. WS push to all party members.

## 4. Leave Party

**POST `/v1/parties/{id}/leave`**
Server:
1. Load party_member.
2. Set `left_at = now()` (soft delete).
3. If character was leader → auto-promote next.
4. If last member → auto-disband.
5. Emit event `party.member_left`.
6. WS push.

## 5. Kick Member

**POST `/v1/parties/{id}/kick`**
Leader only.
```json
{ "target_character_id": "..." }
```

Server:
1. Validate leader.
2. Set target's `left_at = now()`.
3. WS push to target (kicked notification).
4. Emit event `party.member_kicked`.

## 6. Role Hint

Bkz. [roles.md](roles.md).

`role_hint`: tank/heal/dps/support (opsiyonel).
Character self-assign; leader override (Phase 02).

## 7. Ready State

`ready` boolean.
- Leader ready check → prompt.
- Member accepts → ready=true.
- 30 sec timeout → auto false.

## 8. Cross-Region Membership

Members farklı region OK.
Ancak aktivite (dungeon/adventure) location-bound → tüm üyeler oraya gelmeli.

Party leader dungeon enter → members auto-teleport prompt (Phase 02); v1 manuel travel.

## 9. Member Presence

Live state:
- `character.last_seen_at` party UI'da display.
- Offline member: still in party until leave.
- 10 min offline → auto-remove (opsiyonel; Phase 02).

## 10. Concurrency

Character.party membership unique (single active party).

## 11. Metrics
- `party_member_joined_total`.
- `party_member_left_total{cause}`.
- `party_ready_check_avg_response_time_sec`.

## 12. Cross-Reference
- [Party](party.md), [Leader](leader.md), [Invite](invite.md), [Roles](roles.md)
- [database/social-guilds-parties.md §10](../database/social-guilds-parties.md)

## 13. Bekleyen Kararlar
- Cross-region auto-teleport Phase 02.
- Offline auto-remove threshold.
- Member profile detail visibility.
