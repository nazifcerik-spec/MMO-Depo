# Permissions (Rank Permission Matrix)

## 1. Amaç
Rank başına yetki matriksi — content-driven; v1 sabit default.

## 2. Permission Kataloğu

```yaml
permissions:
  # Roster
  invite: bool
  kick_below_own_rank: bool
  promote_up_to_rank: int    # e.g. officer can promote to rank 2
  demote_below_own_rank: bool

  # Guild admin
  set_motd: bool
  set_description: bool
  disband_schedule: bool    # only leader
  transfer_leadership: bool  # only leader

  # Bank
  bank_tab_read: array of tab_index
  bank_tab_write: array of tab_index
  bank_withdraw_gold_per_day_max: int   # 0 = no limit; opsiyonel
  bank_deposit_gold: bool
  purchase_bank_tab: bool

  # Chat
  chat_read: bool
  chat_write: bool
  broadcast_moderator: bool   # send guild-wide announcements

  # Research (Phase 02)
  research_contribute: bool
  research_activate: bool

  # Quest (Phase 02)
  guild_quest_accept: bool
  guild_quest_abandon: bool

  # Boss / Raid (Phase 02)
  guild_boss_summon: bool
  raid_lead: bool

  # War (Phase 02+)
  war_declare: bool
  war_surrender: bool
```

## 3. v1 Default Matrix

| Permission | Leader | Officer | Veteran | Member | Initiate |
|---|---|---|---|---|---|
| invite | ✓ | ✓ | ✓ | ✗ | ✗ |
| kick_below_own_rank | ✓ | ✓ | ✗ | ✗ | ✗ |
| promote_up_to_rank | 0 (all) | 2 (up to veteran) | – | – | – |
| set_motd | ✓ | ✓ | ✗ | ✗ | ✗ |
| disband_schedule | ✓ | ✗ | ✗ | ✗ | ✗ |
| transfer_leadership | ✓ | ✗ | ✗ | ✗ | ✗ |
| bank_tab_read | all | all | 0-2 | 0 | ✗ |
| bank_tab_write | all | all | 0 | 0 | ✗ |
| bank_withdraw_gold_per_day_max | 0 (∞) | 100k | 20k | 5k | 0 |
| bank_deposit_gold | ✓ | ✓ | ✓ | ✓ | ✓ |
| purchase_bank_tab | ✓ | ✓ | ✗ | ✗ | ✗ |
| chat_read | ✓ | ✓ | ✓ | ✓ | ✓ |
| chat_write | ✓ | ✓ | ✓ | ✓ | ✓ |
| broadcast_moderator | ✓ | ✓ | ✗ | ✗ | ✗ |
| research_contribute | ✓ | ✓ | ✓ | ✓ | ✓ |
| research_activate | ✓ | ✓ | ✗ | ✗ | ✗ |
| guild_quest_accept | ✓ | ✓ | ✗ | ✗ | ✗ |
| guild_boss_summon | ✓ | ✓ | ✗ | ✗ | ✗ |
| raid_lead | ✓ | ✓ | ✗ | ✗ | ✗ |
| war_declare | ✓ | ✗ | ✗ | ✗ | ✗ |
| war_surrender | ✓ | ✗ | ✗ | ✗ | ✗ |

## 4. Permission Check

Server middleware:
```go
func checkGuildPermission(char Character, guildID UUID, permission string) error {
    member, _ := repo.LoadGuildMember(guildID, char.ID)
    if member == nil {
        return errs.ErrForbidden
    }
    rank, _ := repo.LoadGuildRank(guildID, member.RankIndex)
    if !rank.Permissions[permission] {
        return errs.ErrForbidden
    }
    return nil
}
```

## 5. Bank Withdraw Daily Cap

- Per rank daily withdraw gold cap.
- Track via ledger correlation:
  ```sql
  SELECT SUM(delta_negative) FROM economy.ledger_entry
  WHERE character_id=X AND source_kind='guild_withdraw'
    AND applied_at > now() - interval '24 hour'
  ```
- Cap exceeded → 429 `guild.withdraw_cap_exceeded`.

## 6. Custom Permission (Phase 02)

Leader edit rank permissions:
```
POST /v1/guilds/{id}/ranks/{rank_index}/permissions
  body: { permissions: {...} }
```

v1'de sabit; Phase 02 esnek.

## 7. Audit

All permission-based actions log to audit.event_log with actor + permission used + timestamp.

## 8. Metrics
- `guild_permission_denied_total{permission}`.
- `guild_bank_withdraw_cap_hit_total`.

## 9. Cross-Reference
- [Roles](roles.md), [Bank](bank.md), [Members](members.md)
- [database/social-guilds-parties.md §2, 6](../database/social-guilds-parties.md)

## 10. Bekleyen Kararlar
- Custom permission Phase 02.
- Withdraw cap tune (100k officer v1; playtest).
- Permission audit compliance UI (Phase 02).
- Cross-rank permission override (leader temporary grant).
