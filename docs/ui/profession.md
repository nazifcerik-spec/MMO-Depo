# Profession Dashboard

## 1. Amaç
Karakterin mesleklerini (mining, fishing, hunting, herbalism, woodcutting, farming, skinning, archaeology) tek panelde göster; level, XP, unlock recipe, mastery tier, aktif contract.

## 2. Layout

```
┌────────────────────────────┐
│ ← Meslekler                │
├────────────────────────────┤
│ ⛏ Madencilik  Lvl 42       │
│    ▓▓▓▓▓░░ XP 62%           │
│    Mastery Tier 3           │
│    [Detay]                  │
├────────────────────────────┤
│ 🌲 Oduncu   Lvl 28          │
│    ▓▓▓░░░ XP 41%            │
│    Mastery Tier 2           │
├────────────────────────────┤
│ 🎣 Balıkçılık Lvl 15        │
│ 🏹 Avcılık   Lvl 22          │
│ 🌿 Bitkicilik Lvl —  [Aç]   │
│ 🌾 Çiftçilik  Phase 02      │
└────────────────────────────┘
```

## 3. Profession Detail

Ör. Mining Detail:
- Current level + XP bar + next tier XP.
- Mastery tier + next tier XP + passive bonus list.
- Unlocked node types (Copper, Iron, Silver, ...).
- Recent gather sessions history.
- **Aktif contract** — buradan hızlı start (Mining L15 Wildwood, 4 saat, tahmini yield).
- Unlockable next recipe/node preview.

## 4. Unlock (Level-Gated Content)

Profession level UP → yeni node type unlock (ör. iron @ Lvl 15) → toast + inbox message.

## 5. State

- Profession provider — `GET /v1/character/{id}/professions`.
- WS `character:{id}/profession` events (xp_gain, level_up).

## 6. API

- `GET /v1/character/{id}/professions`
- `GET /v1/character/{id}/profession/{code}` — detail.
- `POST /v1/activity/start` (kind='mining'/'fishing'/…) — contract start.

## 7. Metrikler
- `ui_profession_view_total{profession}`
- `ui_profession_contract_start_total{profession}`

## 8. Cross-Reference
- [professions](../professions/README.md), [progression/mastery.md](../progression/mastery.md)
- [activity.md](activity.md)
