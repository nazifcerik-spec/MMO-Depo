# Friends

## 1. Amaç
Karakterlerin (opsiyonel: account) birbirini takip etmesi. Presence gösterimi + quick invite + note.

## 2. Karakteristik

- **Scope (v1)**: karakter bazlı friend list (account-friend Phase 02 opt).
- **Cap**: 100 friend karakter başına (v1).
- **Bidirectional**: her iki taraf ekli olmadan da "one-way follow" yok — friend request accept ile karşılıklı.
- **Note field**: kendiye ait 128 char özel not (karşı taraf görmez).
- **Presence**: online / away / offline / dnd + son görülme (privacy setting ile gizlenebilir).
- **Cross-region visibility**: karakter aynı shard/region'da olmasa da görülür (v1 tek shard).

## 3. Model

```sql
CREATE TABLE social.friend_request (
  id UUID PRIMARY KEY,
  requester_character_id UUID NOT NULL,
  target_character_id UUID NOT NULL,
  message VARCHAR(128) NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'pending',   -- pending|accepted|declined|expired|cancelled
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  resolved_at TIMESTAMPTZ NULL,
  expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '14 days',
  UNIQUE (requester_character_id, target_character_id)
);

CREATE TABLE social.friend (
  character_id UUID NOT NULL,
  friend_character_id UUID NOT NULL,
  friended_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  note VARCHAR(128) NULL,
  favorite BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (character_id, friend_character_id)
);
```

Accept anında iki satır insert (bidirectional).

## 4. Presence

`character.presence` Redis'te tutulur (ephemeral):
```
key: presence:character:{id}
value: { state: 'online', last_seen: <ts>, location_region_id?, activity_kind? }
TTL: 5 dk (heartbeat ile yenilenir)
```

Presence subscription (WS): friend ekranı açıkken client subscribes to friend id list; presence değişiminde WS event.

Privacy:
- `character.presence_visibility`: `friends` (default) / `guild` / `everyone` / `no_one`.
- `no_one` → presence her zaman "offline" görünür (kendi hariç).

## 5. Akış

### Add
`POST /v1/friends/request`
```json
{ "target_character_name": "Grim", "message": "hey, dun grubu?" }
```
Server:
1. Target var mı, aktif mi (soft-delete değil).
2. Block ilişkisi var mı (her iki yön) → var ise silent success (drop; requester "sent" görür, target hiç görmez — targeted user pattern).
3. Zaten friend mi → `E_ALREADY_FRIENDS`.
4. Pending request var mı → replace (aynı request expiry uzatılır).
5. Insert `friend_request`, WS push target'a (privacy blocking değilse).

### Accept
`POST /v1/friends/request/{id}/accept`
- Bidirectional insert.
- Requester notify.
- Metric `friend_request_accepted_total`.

### Decline
Silent (requester notify edilmez — anti-shame).

### Remove
`DELETE /v1/friends/{friend_character_id}` — iki taraflı sil.
Karşı tarafa notify **yok** (silent removal).

## 6. Endpoint Özet

**GET `/v1/friends`** — kendi list + presence snapshot.
**GET `/v1/friends/requests`** — pending in/out.
**POST `/v1/friends/request`** — add.
**POST `/v1/friends/request/{id}/accept`** / **decline** / **cancel**.
**DELETE `/v1/friends/{friend_character_id}`**.
**PUT `/v1/friends/{friend_character_id}`** — note / favorite update.

## 7. Rate Limit

- Friend request send: 10/dk karakter, 30/gün karakter.
- Friend list read: 60/dk (client polling ok).

## 8. Metrikler
- `friend_request_sent_total`, `friend_request_accepted_total`, `friend_request_declined_total`
- `friend_count_avg`
- `friend_presence_update_rate`

## 9. Cross-Reference
- [block.md](block.md) — Block precedence
- [chat.md](chat.md) — Friend online notify → WS event
- [systems/notification.md](../systems/notification.md)
- [character/character.md](../character/character.md)

## 10. Bekleyen Kararlar
- Account-level friend (aynı player'ın tüm karakterleri) Phase 02.
- Friend cap 100 v1 (Phase 02: guild perk ile artabilir).
- Cross-server friend Phase 03.
- Presence "custom status text" Phase 02.
