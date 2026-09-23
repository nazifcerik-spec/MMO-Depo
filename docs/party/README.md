# Party System (Faz 23)

Küçük ölçekli (2-5 kişi) geçici sosyal birim. Coop dungeon, party hunt, adventure, event katılım.

Foundation:
- [game-design/party.md](../game-design/party.md) — vizyon.
- [database/social-guilds-parties.md §9-11](../database/social-guilds-parties.md) — schema.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Party](party.md) | Party entity + lifecycle |
| [Leader](leader.md) | Leader yetkileri + migration |
| [Member](member.md) | Member state + join/leave |
| [Invite](invite.md) | Invite flow + LFG |
| [Roles](roles.md) | Tank/heal/dps role hint |
| [Combat](combat.md) | Party combat sync + aura |
| [Hunt](hunt.md) | Party hunt aktivite |
| [Dungeon](dungeon.md) | Party dungeon integration |
| [Loot](loot.md) | Loot rule + distribution |

## 2. Party Karakteristik

- Size: **2-5** (v1); Phase 02 raid 10+ ayrı system.
- Ephemeral (session-scope) veya persistent (Phase 02 saved party).
- Cross-region members OK; ancak aktivite location-bound.
- Chat channel: `party`.
- Loot rule: personal loot default (DD-012).

## 3. Endpoint Genel

```
GET  /v1/characters/{id}/party                     # current party
POST /v1/parties                                    # create
POST /v1/parties/{id}/invite                        # invite member
POST /v1/parties/{id}/join                          # accept invite
POST /v1/parties/{id}/leave                         # leave party
POST /v1/parties/{id}/kick                          # leader kick
POST /v1/parties/{id}/promote                       # leader promote
POST /v1/parties/{id}/disband                       # leader disband
POST /v1/parties/{id}/loot-rule                     # change loot rule
POST /v1/lfg/list                                   # LFG listing create
GET  /v1/lfg/browse                                 # browse LFG listings
```

## 4. Cross-Reference
- [game-design/party.md](../game-design/party.md)
- [database/social-guilds-parties.md §9-11](../database/social-guilds-parties.md)
- [activity-engine/dungeon.md](../activity-engine/dungeon.md), [activity-engine/hunting.md](../activity-engine/hunting.md)
- [combat/combat-loop.md](../combat/combat-loop.md)
- [systems/loot-system.md](../systems/loot-system.md)

## 5. Ortak Bekleyen Kararlar

- Persistent party Phase 02 (saved composition).
- Raid (10+ kişi) Phase 02+.
- Cross-guild party (default OK).
- Party finder / LFG matchmaking algorithm Phase 02.
- Voice chat integration Phase 03.
