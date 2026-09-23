# Guild (Admin View)

## 1. Amaç
Guild teşhis + moderasyon: liste, member, treasury, bank audit, disband override, MOTD moderation.

## 2. List View
- Search: name, tag, guild_id, member name.
- Filter: level range, member count range, state (active/disband_grace).
- Columns: name · level · members · treasury · leader · created_at.

## 3. Detail View
```
Guild: Iron Fist (guild_ab…)
Tag: [IRON] · Level 7
Members: 34/50
Treasury: 145,320 gold
Founded: 2026-04-15 by Grim

Tabs:
[Members] [Bank] [Chat Log] [Audit] [Events]
```

### Members
Rank + join date + last online + contribution + banka withdraw geçmişi.
Aksiyon: **remove member** (support), **promote/demote** (rare; guild leader normalde yapar).

### Bank
Tab bazında item + gold; withdraw audit trail full (kim aldı, ne zaman, ne).

### Chat Log
Guild + officer chat 30 gün retention (bkz. [social/guild-chat.md](../social/guild-chat.md)).
Read-only forensic; toxic content moderation input.

### Audit
Guild üzerine yapılmış admin aksiyonların history'si.

### Events
Guild-scope events (boss summon, raid create, war declare — Phase 02).

## 4. Aksiyonlar

| Action | Permission | Step-up |
|---|---|---|
| `guild.member_remove` | support_l2+ | no |
| `guild.motd_clear` | support_l2+ | no (toxic MOTD) |
| `guild.rename_force` | admin | yes (name violation) |
| `guild.disband_override` | admin | yes (rule violation; bypass 48h grace) |
| `guild.treasury_reversal` | admin | yes (bkz. economy) |
| `guild.leadership_transfer` | admin | yes (inaktif leader case) |

Auto-leader migration zaten var (DD-130 5min inaktif); admin manual override farklı case.

## 5. Endpoint

- `GET /internal/admin/guild?…`
- `GET /internal/admin/guild/{id}` compound
- `GET /internal/admin/guild/{id}/bank/audit?…`
- `POST /internal/admin/guild/{id}/member/{character_id}/remove`
- `POST /internal/admin/guild/{id}/motd/clear`
- `POST /internal/admin/guild/{id}/rename-force` `{ new_name, reason }`
- `POST /internal/admin/guild/{id}/disband-override` `{ reason }`
- `POST /internal/admin/guild/{id}/leader-transfer` `{ new_leader_character_id, reason }`

## 6. Cross-Reference
- [guild](../guild/README.md), [social/guild-chat.md](../social/guild-chat.md)

## 7. Bekleyen Kararlar
- Guild name violation policy (auto vs manual review) Phase 02.
- Toxic MOTD auto-detection Phase 02.
- Alliance admin (Phase 03) yeni yetki matrisi.
