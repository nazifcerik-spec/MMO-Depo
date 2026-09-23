# Social

## 1. Amaç
Friends + Block + Mail + Chat konsol tabları.

## 2. Layout — Social Tab

```
┌────────────────────────────┐
│ ← Sosyal                   │
│ [Friends][Mail][Chat][Block]│
├────────────────────────────┤
│ Online (5):                │
│  🟢 Aria    Warrior 42     │
│  🟢 Kael    Priest 41 (dnd)│
│  🟢 Zara    Hunter 38      │
│  ...                       │
│ Offline (12):              │
│  ⚫ Bran    Rogue 35        │
│  ...                       │
├────────────────────────────┤
│ Pending Requests: 2        │
│ [Arkadaş Ekle]             │
└────────────────────────────┘
```

## 3. Friends Sub-tab

- Presence indicator (green/yellow/red/grey).
- Note field (kendi görür).
- Favorite pin.
- Actions: whisper, invite party, invite guild, view profile, remove, block.

## 4. Mail Sub-tab

```
┌────────────────────────────┐
│ Inbox (12 unread)          │
├────────────────────────────┤
│ 🔴 [attach] Grim: Loot     │
│    Sun · 2h                │
│ ⚪ System: Server maint     │
│    Sun · 1d                │
│  ...                       │
├────────────────────────────┤
│ [Yeni Mail]                │
└────────────────────────────┘
```

Mail detail:
- Subject + body.
- Attachment slots (12 max) + gold.
- **Claim** → `POST /v1/mail/{id}/claim`.
- **Return** / **Reply** (Phase 02) / **Delete**.

## 5. Chat Sub-tab (Full-Screen Chat)

Kanal switcher (top tabs): Global / Region / Trade / LFG / Party / Guild / Whisper (list).

Scroll + input; slash commands (`/whisper name`, `/party`, `/guild`).

Whisper list ayrı — DM conversation aç.

## 6. Block Sub-tab

- Block list (account-level).
- Search + add (name → account resolve).
- Remove (unblock).

## 7. State

- Friends provider — `GET /v1/friends` + presence WS.
- Mail provider — `GET /v1/mail/inbox` + WS `new_mail`.
- Chat state — per-channel Redis rolling + WS.
- Block provider — `GET /v1/block`.

## 8. API

- Friends / block: [social/friends.md](../social/friends.md), [social/block.md](../social/block.md).
- Chat: [social/chat.md](../social/chat.md).
- Mail: [social/mail.md](../social/mail.md).

## 9. Cross-Reference
- [social](../social/README.md), [notification](../notification/README.md)
