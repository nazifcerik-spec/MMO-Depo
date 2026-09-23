# Roles (Rank Hierarchy)

## 1. Amaç
Guild rank sistemi — 5 sabit rank v1; Phase 02 customizable.

## 2. Rank Model

Kaynak: [database/social-guilds-parties.md §2](../database/social-guilds-parties.md).

```sql
social.guild_rank
  guild_id, rank_index (0..9)
  name, permissions_json
```

## 3. v1 Default Rank Set

| rank_index | Name | Permissions Summary |
|---|---|---|
| 0 | Leader | All permissions |
| 1 | Officer | Invite, kick (below), promote (up to 2), bank all tabs, disband schedule |
| 2 | Veteran | Invite, chat, bank tab 0-2 read/write |
| 3 | Member | Chat, bank tab 0 read/write, contribution |
| 4 | Initiate | Chat read, no bank write |

Content-driven per guild (Phase 02 customization).

## 4. Promote / Demote

**POST `/v1/guilds/{id}/promote`**
```json
{ "target_character_id": "...", "new_rank_index": 1 }
```

Server (leader/officer only):
1. Validate promoter rank < new_rank.
2. Validate target rank > new_rank (only promote up).
3. Validate promoter can promote to new_rank (officer max promote to rank 2).
4. Update guild_member.rank_index.
5. Emit event `guild.member_promoted`.

**POST `/v1/guilds/{id}/demote`**
Similar; opposite direction.

## 5. Leader Transfer

Leader promotes another member to rank 0:
- Old leader automatically demoted to rank 1 (officer).
- New leader takes 0.
- Transaction.

**Alternate**: Auto-migration on leader leave (5min timeout):
- Longest-active member auto-promoted.
- Notification.

## 6. Custom Rank (Phase 02)

- Guild leader creates custom ranks (max 10 total).
- Custom names + permission grid.
- Rename existing ranks.

v1'de sabit 5 rank; Phase 02 esnek.

## 7. Rank Prestige

- Rank rank up milestone: unique title ("Officer of Emerald Wardens").
- Officer badge on character nameplate.

## 8. Contribution Requirement (Phase 02)

Automatic promotion:
- Contribution X points + Y hafta member → auto-promote to rank Z.
- Optional; opt-in.

v1'de manuel promotion only.

## 9. Metrics
- `guild_member_rank_distribution{guild,rank}` gauge.
- `guild_promotion_total`.

## 10. Cross-Reference
- [Members](members.md), [Permissions](permissions.md)
- [database/social-guilds-parties.md §2](../database/social-guilds-parties.md)

## 11. Bekleyen Kararlar
- Custom rank Phase 02.
- Auto-promotion Phase 02.
- Rank badge cosmetic (nameplate icon).
- Officer count limit (v1 unlimited; Phase 02 max 5).
