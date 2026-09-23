# Guild

## 1. Amaç
Guild dashboard: overview, members, bank, chat, MOTD, permissions, events (boss/raid/war Phase 02).

## 2. Layout — Guild Home

```
┌────────────────────────────┐
│ ← Iron Fist [tabard]      │
│ Lvl 7 · 34/50 üye          │
│ Treasury: 145,320 gold     │
├────────────────────────────┤
│ MOTD:                      │
│ "Salı 20:00 boss run!"     │
├────────────────────────────┤
│ [Members][Bank][Chat][Info]│
├────────────────────────────┤
│ Online: 8                  │
│  • Grim (Leader)           │
│  • Aria (Officer)          │
│  • Kael                    │
│  ...                       │
├────────────────────────────┤
│ Weekly Contribution: ▓▓▓░ │
│ Sen: 4,500 / 5,000 cap     │
└────────────────────────────┘
```

## 3. Sub-tabs

### Members
- Full member list with rank/last-online/contribution.
- Filter by rank.
- Actions (rank permission göre): invite, kick, promote, demote.

### Bank
- Tab selector (5 tab, level-gated).
- Item grid per tab.
- Deposit / withdraw with daily cap indicator.
- Audit log link (leader only).

### Chat
- Ana kanal + officer kanal switcher.
- Scroll + input.

### Info
- Guild name, tabard, faction (Phase 24), founded date, achievement, rank progression.
- Research / boss / raid / war (Phase 02 sub-screens).

## 4. Invite Flow

Officer/leader → `POST /v1/guild/{id}/invite` (target character name) → target push.

## 5. State

- Guild provider — `GET /v1/character/{id}/guild` + WS `guild:{id}` events.

## 6. API

- `GET /v1/guild/{id}`
- `GET /v1/guild/{id}/members`
- `GET /v1/guild/{id}/bank?tab=`
- `POST /v1/guild/{id}/bank/deposit` / `withdraw`
- `POST /v1/guild/{id}/invite`
- `POST /v1/guild/{id}/leave`
- `PUT /v1/guild/{id}/motd`

## 7. Cross-Reference
- [guild](../guild/README.md), [social/guild-chat.md](../social/guild-chat.md)
