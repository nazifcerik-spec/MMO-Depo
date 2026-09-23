# Activity (AFK Contract)

## 1. Amaç
Aktif AFK contract başlatma/durdurma, reveal ekranı, tahmini ödül önizleme, geçmiş.

## 2. Layout — Home Dashboard (aktif contract mini widget)

```
┌────────────────────────────┐
│ Aktif Aktivite             │
│ ⛏ Iron Mining · Wildwood   │
│ ▓▓▓▓▓▓░░░ 5s 32dk / 8s     │
│ Tahmini: ~180 iron, ~12 rare│
│ [Detay]  [Değiştir]  [Dur] │
└────────────────────────────┘
```

## 3. Layout — Activity Detail

```
┌────────────────────────────┐
│ ← Aktivite: Iron Mining    │
├────────────────────────────┤
│  [context illustration]     │
│  Wildwood Copper Vein L15  │
│  Başlangıç: 12:34 UTC       │
│  Süre: 8 saat (28dk)         │
│  Enerji: 15/saat            │
├────────────────────────────┤
│ Tahmini Ödül (analitik):    │
│  ⛏ 180 iron ore (±10)      │
│  💎 12 rare gem (±3)        │
│  ⭐ 2 legendary şansı %8    │
│  XP: 1250 mining mastery   │
├────────────────────────────┤
│ [Şimdi Reveal] [Değiştir]  │
│ [Durdur]                   │
└────────────────────────────┘
```

## 4. Reveal Ekranı (Notification tap veya manuel)

Loot animation sırayla:
1. Common item cards yağar (0.5s stagger).
2. Rare item cards için parıltı + sound.
3. Epic/legendary için full-screen shine.
4. Summary: XP gain, mastery tier up, achievement unlock.
5. "Devam" → yeni contract seçim.

## 5. Contract Seçim Modal

Aktivite listesi (unlocked at level/region):
- Mining, Woodcutting, Fishing, Hunting, Herbalism, Farming, Combat area, Hunting adventure, Dungeon AFK, Travel.
- Her aktivite için: unlock level, energy/hr, tahmini yield.
- Süre slider: 15 dk / 30 dk / 1s / 2s / 4s / 8s / 12s (energy'ye göre cap).

## 6. Tahmini Ödül (Server Estimate)

Server hesaplar (analitik — DD-084) → `GET /v1/activity/estimate?kind=&duration=`. Sonuç varyans dahil (mean ± sd).

## 7. State

- Active contract provider — WS `character:{id}/activity` events.
- Reveal state — one-shot; cache 5 dk (accidental re-reveal idempotent).

## 8. API

- `GET /v1/character/{id}/activity/active`
- `POST /v1/activity/start` body: `{ kind, params, duration_seconds, idempotency_key }`
- `POST /v1/activity/{id}/reveal`
- `POST /v1/activity/{id}/stop`
- `GET /v1/character/{id}/activity/history?limit=50`

## 9. Concurrent Cap

Max concurrent AFK slot (DD-OPEN-3: 1 ana + 1 pet + 1 travel öneri). UI'da slot indicator; slot dolu → conflict modal.

## 10. Metrikler
- `ui_activity_start_total{kind}`
- `ui_activity_reveal_open_total`
- `ui_activity_reveal_latency_ms`
- `ui_activity_change_total`

## 11. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)
- [notification/activity-completed.md](../notification/activity-completed.md)
