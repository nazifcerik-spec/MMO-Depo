# Party (Entity + Lifecycle)

## 1. Amaç
Party entity — creation, size, activity binding, disband.

## 2. Model

Kaynak: [database/social-guilds-parties.md §9](../database/social-guilds-parties.md).

```sql
social.party
  id, leader_character_id, created_at
  disbanded_at NULL
  loot_rule (personal/round_robin/need_before_greed/master)
  activity_instance_id NULL         -- bound activity
```

Ek Phase 23:
```sql
ALTER TABLE social.party
  ADD COLUMN size_cap SMALLINT NOT NULL DEFAULT 5,
  ADD COLUMN persistent BOOLEAN NOT NULL DEFAULT false,        -- Phase 02
  ADD COLUMN name VARCHAR(48) NULL,                             -- Phase 02 named party
  ADD COLUMN version INT NOT NULL DEFAULT 1;
```

## 3. Lifecycle

```
1. Create (POST /v1/parties).
   - Leader = creator.
   - Empty member list (leader only).
2. Invite members (bkz. invite.md).
3. Members accept join.
4. Bind activity (dungeon enter, hunt start).
5. Members leave/kicked or activity end.
6. Disband:
   - Explicit (leader disband).
   - Auto (all members left).
   - Timeout (session-scope; 5 min no activity).
```

## 4. Create Endpoint

**POST `/v1/parties`**
Auth + CharacterScope.

Server:
1. Validate character not in party.
2. Insert social.party (leader_character_id = char.id).
3. Insert social.party_member (leader as first member).
4. Emit event `party.created`.
5. Response 201 with party detail.

## 5. Disband Endpoint

**POST `/v1/parties/{id}/disband`**
Auth + leader only.

Server:
1. Load party.
2. Set `disbanded_at = now()`.
3. Delete party_member rows (or soft-delete `left_at`).
4. Emit event `party.disbanded`.
5. WS push to all members.

## 6. Party Detail

**GET `/v1/characters/{id}/party`**
```json
{
  "id": "...",
  "leader_character_id": "...",
  "members": [
    { "character_id": "...", "role_hint": "tank", "ready": true, "level": 25, "class": "warrior" },
    { "character_id": "...", "role_hint": "dps", "ready": false, "level": 24, "class": "ranger" }
  ],
  "loot_rule": "personal",
  "activity_instance_id": null,
  "created_at": "..."
}
```

## 7. Size Constraint

`size_cap` = 5 (v1 hard cap).
Party full → invite reject.

Party 1 kişi = "solo party" (opsiyonel; v1'de allow — dungeon solo idle-run için).

## 8. Cross-Character Consistency

Character.party membership unique:
- `social.party_member` unique constraint per character (aktif membership).
- Character birleşemez ikinci party'ye (leave first).

## 9. Activity Binding

Party activity_instance_id set:
- Dungeon enter → instance ref.
- Adventure start → instance ref.
- Hunt aktivite: her karakter kendi contract'ı (party sync at aura level).

Party activity leave:
- Aktif üye instance'dan çıkar.
- Party disband olabilir (last member left).

## 10. Persistent Party (Phase 02)

- Party name + logo.
- Save composition (preferred members list).
- Auto-invite reconnect.

v1'de yok (ephemeral).

## 11. Metrics
- `party_created_total`.
- `party_disbanded_total{cause}`.
- `party_active_count` gauge.
- `party_avg_size` gauge.

## 12. Cross-Reference
- [Leader](leader.md), [Member](member.md), [Invite](invite.md), [Roles](roles.md)
- [database/social-guilds-parties.md §9-11](../database/social-guilds-parties.md)
- [game-design/party.md](../game-design/party.md)

## 13. Bekleyen Kararlar
- Persistent party Phase 02.
- Party name customization.
- Cross-account party (aynı user'ın 2 hesabı) Phase 02.
- Party timeout (5 min no activity) — v1 evet; Phase 02 config.
