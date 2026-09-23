# Quest Log

## 1. Amaç
Aktif quest tracker, quest log (all), quest turn-in, category filter (main/side/daily/weekly/dynamic/event/guild).

## 2. Layout — Quest Tracker (HUD)

```
┌──────────────────────────┐
│ Aktif Görev              │
│ [Main] Ashen King        │
│  Objective: 3/5 rune bul │
│ ▓▓▓░░ 60%                │
│ [Detay]                  │
└──────────────────────────┘
```

Ekranın kenarında minimal widget (Home dashboard).

## 3. Quest Log (Full)

```
┌────────────────────────────┐
│ ← Görevler                 │
│ [Main][Side][Daily][Weekly]│
│ [Dynamic][Event][Guild]    │
├────────────────────────────┤
│ [Main] Ashen King's Return │
│  Lvl 40 · 3/5 rune bulundu │
│  Ödül: 2,500c + epic item │
│  [Track] [Detay]           │
├────────────────────────────┤
│ [Side] Missing Miner       │
│  Lvl 15 · Miner ile konuş │
│  [Track] [Detay]           │
├────────────────────────────┤
│ [Daily] Fresh Meat         │
│  0/10 boar killed          │
│  Reset: 12h                │
│  [Track] [Vazgeç]          │
└────────────────────────────┘
```

## 4. Quest Detail

- Story text (scrollable).
- Objectives (checklist).
- Rewards preview (gold, xp, item, reputation, event_token).
- NPC turn-in target.
- Deep link to relevant map location.

## 5. Turn-In Flow

Quest complete flag on → notification "Görev tamamlandı" → turn-in NPC (map marker highlight).
NPC etkileşim → reward modal → apply.

## 6. Track / Untrack

Kullanıcı 3'e kadar track edebilir (HUD widget). Fazla → replace prompt.

## 7. State

- Quest log provider — `GET /v1/character/{id}/quests`.
- WS `quest:{id}` events (objective_progress, complete, fail).

## 8. API

- `GET /v1/character/{id}/quests?category=`
- `GET /v1/quest/{code}` — detail.
- `POST /v1/character/{id}/quest/{id}/track`
- `POST /v1/character/{id}/quest/{id}/abandon` (daily/side only).
- `POST /v1/character/{id}/quest/{id}/turn-in`

## 9. Cross-Reference
- [quests](../quests/README.md)
