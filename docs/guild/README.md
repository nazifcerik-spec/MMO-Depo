# 🔥 Guild System (Faz 24)

Kalıcı büyük sosyal birim. 50 üye başlangıç, level ile 100'e kadar. Bank, research, quest, boss, raid (Phase 02), war (Phase 02+).

Foundation:
- [game-design/guilds.md](../game-design/guilds.md) — vizyon.
- [database/social-guilds-parties.md §1-8](../database/social-guilds-parties.md) — schema.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Guild](guild.md) | Entity + lifecycle + disband grace |
| [Members](members.md) | Roster + invite + kick + contribution |
| [Roles](roles.md) | Rank hierarchy + customization Phase 02 |
| [Permissions](permissions.md) | Rank permission matrix |
| [Level](level.md) | Guild level + cap + benefits |
| [XP](xp.md) | Guild XP kaynakları + sink |
| [Bank](bank.md) | Guild bank tabs + audit + gold treasury |
| [Research](research.md) | **Phase 02** — passive guild perks |
| [Quest](quest.md) | Guild quest özet + refers quests/guild.md |
| [Boss](boss.md) | Guild boss encounter — Phase 02 |
| [Raid](raid.md) | 10+ kişi raid — Phase 02 |
| [War](war.md) | Guild war + territory control — Phase 02+ |

## 2. Karakteristik

- **50 üye başlangıç**, guild level ile 100'e (DD-018).
- **Persistent** — kalıcı sosyal ev.
- **Guild XP** aggregate member aktivite.
- **Guild treasury** — shared gold.
- **Guild bank** — tab-based item storage.
- **Chat channel**: `guild`.

## 3. Ortak Sözleşme

1. **Character 1 guild max** (aktif membership unique).
2. **Guild leader disband** — 48h grace.
3. **Rank permissions** — content-driven; v1 sabit 5 rank.
4. **Bank permissions per tab** — rank-based read/write.
5. **Audit log** — bank + admin action forensic.

## 4. Endpoint Genel

```
GET  /v1/guilds/{id}                                # detail
POST /v1/guilds                                     # create
POST /v1/guilds/{id}/disband
POST /v1/guilds/{id}/invite
POST /v1/guilds/{id}/join
POST /v1/guilds/{id}/leave
POST /v1/guilds/{id}/kick
POST /v1/guilds/{id}/promote
POST /v1/guilds/{id}/demote
POST /v1/guilds/{id}/motd
GET  /v1/guilds/{id}/roster
GET  /v1/guilds/{id}/bank
POST /v1/guilds/{id}/bank/deposit
POST /v1/guilds/{id}/bank/withdraw
POST /v1/guilds/{id}/bank/tab/purchase
GET  /v1/guilds/{id}/research                       # Phase 02
POST /v1/guilds/{id}/research/{code}/contribute
GET  /v1/guilds/{id}/quests
POST /v1/guilds/{id}/quests/{code}/accept
```

## 5. Cross-Reference

- [game-design/guilds.md](../game-design/guilds.md)
- [database/social-guilds-parties.md](../database/social-guilds-parties.md)
- [quests/guild.md](../quests/guild.md)

## 6. Ortak Bekleyen Kararlar

- Guild alliance Phase 03.
- Cross-realm guild Phase 03.
- Guild housing Phase 02+.
- Voice chat Phase 03.
- Custom guild UI (tabard, logo customizer) Phase 02.
