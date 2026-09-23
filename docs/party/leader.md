# Leader

## 1. Amaç
Party leader — yönetim yetkileri + migration.

## 2. Leader Yetkileri

- Invite member.
- Kick member.
- Promote member to leader.
- Disband party.
- Set loot rule.
- Start party activity (dungeon enter, adventure).
- Ready check trigger.
- Party name change (Phase 02 persistent).

## 3. Leader Transfer (Promote)

**POST `/v1/parties/{id}/promote`**
```json
{ "target_character_id": "..." }
```

Server (leader only):
1. Validate target is party member.
2. Update `social.party.leader_character_id = target`.
3. Emit event `party.leader_promoted`.
4. WS push to all members.

Optional: 2-tier permission (Phase 02 co-leader Assistant).

## 4. Auto Migration

Leader disconnected 5 min:
- Auto-promote longest-active member.
- Notification "New leader: <name>".

Leader leave manual:
- Explicit promote first, OR auto-promote next.

## 5. Ready Check

Leader can trigger ready check before activity:
```
POST /v1/parties/{id}/ready-check
```

Members receive prompt; 30 sec window.
Response updates `party_member.ready`.

Leader can see "3/5 ready".

## 6. Vote Kick

Non-leader can initiate vote kick:
```
POST /v1/parties/{id}/vote-kick
  body: { target_character_id: "..." }
```

Requires 2/3 approval (voting other members).
Leader can override (kick directly).

## 7. Metrics
- `party_leader_transferred_total`.
- `party_leader_auto_migrated_total`.
- `party_vote_kick_total{result}`.

## 8. Cross-Reference
- [Party](party.md), [Member](member.md)
- [game-design/party.md](../game-design/party.md)

## 9. Bekleyen Kararlar
- Co-leader Phase 02.
- Auto-migration timeout (v1: 5 min).
- Vote kick threshold (2/3 v1).
