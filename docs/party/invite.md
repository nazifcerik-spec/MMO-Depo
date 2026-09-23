# Invite (Invitation Flow)

## 1. Amaç
Party invite lifecycle + LFG (Looking For Group).

## 2. Invite Model (yeni Phase 23)

```sql
CREATE TABLE social.party_invite (
  id UUID PRIMARY KEY,
  party_id UUID NOT NULL,
  target_character_id UUID NOT NULL,
  invited_by_character_id UUID NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',      -- pending/accepted/declined/expired
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (party_id, target_character_id) WHERE status='pending'
);
```

## 3. Invite Endpoint

**POST `/v1/parties/{id}/invite`**
Leader (v1 leader only) or member (Phase 02 config).
```json
{ "target_character_name": "Fenrir" }
```

Server:
1. Validate party size < size_cap.
2. Validate target character exists + online (opsiyonel — offline OK Phase 02).
3. Validate target not in another party.
4. Insert social.party_invite (expires_at = now + 60 sec).
5. WS push to target character.
6. Emit event `party.invite_sent`.
7. Response 200.

## 4. Accept Invite

**POST `/v1/parties/{id}/join`**
Target character:
- Server validates invite pending + not expired.
- Add member (bkz. [member.md §3](member.md)).
- Mark invite `accepted`.

## 5. Decline Invite

**POST `/v1/parties/{id}/invite/decline`**
Target character:
- Mark invite `declined`.
- WS push to inviter.

## 6. Expire

Invite 60 sec no response → `expired`.
Cleanup worker daily archive.

## 7. LFG (Looking For Group)

Bkz. [database/social-guilds-parties.md §11](../database/social-guilds-parties.md).

```sql
social.lfg_listing
  id, character_id
  activity_kind (dungeon/adventure/hunt/world_boss)
  template_id
  role_hint
  message (160 char)
  created_at, expires_at
  status (open/matched/expired/cancelled)
```

## 8. LFG Create

**POST `/v1/lfg/list`**
```json
{
  "activity_kind": "dungeon",
  "template_id": "dungeon-ashen-keep",
  "role_hint": "tank",
  "message": "LFG Ashen Keep, chill run"
}
```

Server:
1. Validate character not in party (or leader looking for members).
2. Insert lfg_listing (expires_at = now + 1 hour).
3. Emit event `lfg.listed`.

## 9. LFG Browse

**GET `/v1/lfg/browse`**
```json
{
  "listings": [
    {
      "id": "...",
      "character_summary": { "name": "Fenrir", "level": 25, "class": "warrior" },
      "activity_kind": "dungeon",
      "template_code": "ashen_keep",
      "role_hint": "tank",
      "message": "...",
      "created_at": "..."
    }
  ]
}
```

Filter: activity_kind + template + role + level range.

## 10. LFG → Invite Flow

Browser sees listing → clicks "Invite":
- If browser has party → invite target to party.
- If browser solo → create party + invite target.

## 11. LFG Matchmaking (Phase 02)

- Auto-match by activity + role complementarity.
- Queue-based (like WoW LFR).
- Group pop notification.

v1'de manuel browse + invite.

## 12. Metrics
- `party_invite_sent_total`.
- `party_invite_accepted_ratio` gauge.
- `lfg_listed_total{activity}`.
- `lfg_matched_total`.

## 13. Cross-Reference
- [Party](party.md), [Member](member.md), [Roles](roles.md)
- [database/social-guilds-parties.md §11](../database/social-guilds-parties.md)

## 14. Bekleyen Kararlar
- Matchmaking Phase 02.
- Cross-guild recommendation.
- Invite spam protection (5 pending max Phase 02).
- Auto-invite based on friend list Phase 02.
