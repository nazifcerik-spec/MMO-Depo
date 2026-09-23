# Party

## 1. Amaç
Party create/invite, member yönetimi, party chat, party hunt/dungeon başlatma.

## 2. Layout — Party HUD (In-World)

```
┌────────────────────────┐
│ Party (3/5)           │
│ • Grim (leader) HP▓▓▓ │
│ • Aria         HP▓▓░  │
│ • Kael         HP▓░░  │
│ [Ayrıl] [Chat] [Hunt] │
└────────────────────────┘
```

## 3. Party Screen (Full)

```
┌────────────────────────────┐
│ ← Party (3/5)              │
├────────────────────────────┤
│ Grim ★ Leader   Warrior 42 │
│  HP 850/950  MP 100/300    │
│  Role: Tank                │
├────────────────────────────┤
│ Aria           Mage 40     │
│  HP 500/600  MP 400/500    │
│  Role: DPS                 │
├────────────────────────────┤
│ Kael           Priest 41   │
│  HP 600/700  MP 350/400    │
│  Role: Healer              │
├────────────────────────────┤
│ Loot Rule: Personal        │
│ [Rol Ata] [Kick] [Promote] │
│ [Davet Et]  [Party Ayrıl]  │
└────────────────────────────┘
```

## 4. Bileşenler

- **Member card**: name + class + level + HP/MP + role + leader crown.
- **Loot rule indicator**: personal / round-robin / master loot.
- **Actions**: invite (name search), promote leader, kick (leader only), assign role.
- **Party chat inline**: bottom drawer.
- **Party hunt / dungeon CTA**: aktif party ile activity.

## 5. Invite Flow

- **Send invite**: name input → `POST /v1/party/invite` (target target character) → target push notification.
- **Receive invite**: modal "Grim seni Party'ye davet etti — [Kabul Et] [Reddet] [Bloke Et]".

## 6. State

- Party provider — `GET /v1/character/{id}/party` + WS `party:{id}` events (member_joined, left, hp_change, role_change).

## 7. API

- `GET /v1/character/{id}/party`
- `POST /v1/party/create`
- `POST /v1/party/invite` body: `{ target_character_name }`
- `POST /v1/party/invite/{id}/accept`
- `POST /v1/party/{id}/leave`
- `POST /v1/party/{id}/member/{character_id}/kick`
- `POST /v1/party/{id}/member/{character_id}/promote`
- `PUT /v1/party/{id}/member/{character_id}/role`

## 8. Cross-Reference
- [party](../party/README.md), [social/party-chat.md](../social/party-chat.md)
- [combat.md](combat.md), [dungeon.md](dungeon.md)
