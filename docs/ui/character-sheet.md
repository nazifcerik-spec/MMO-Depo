# Character Sheet

## 1. Amaç
Karakterin sayısal ve görsel durumu: level, XP, HP/mana/stamina/energy bar, primary stats, secondary stats, aktif buff/debuff, aktif title.

## 2. Layout (Mobile 375×812)

```
┌────────────────────────────┐
│  ← Grim, the Ironmaster  ≡ │
│  Lvl 42 · Warrior · Human │
├────────────────────────────┤
│  ┌── 3D pose ──┐  HP 850/950│
│  │             │  MP 200/300│
│  │             │  ST 100/100│
│  └─────────────┘  EN  87/100│
│  XP ▓▓▓▓▓░░░░ 62% (12h) │
├────────────────────────────┤
│ [Stats] [Inventory] [Gear] │
├────────────────────────────┤
│  STR 45   INT 12   CHA 8   │
│  AGI 28   VIT 32   LUCK 5  │
│  ──────                    │
│  ⚔  Atk 128   🛡  Def 92   │
│  ★  Crit 12%  🎯 Hit 95%   │
│  ⏱  Cd -5%    🍀 Find +8%  │
├────────────────────────────┤
│ Aktif Buff/Debuff:         │
│ 🟢 Iron Vigor (28dk)       │
│ 🟢 Invasion Survivor (12h) │
│ 🔴 Poison III (24s)        │
└────────────────────────────┘
```

## 3. Bileşenler

- **Header**: name + title (active), level, class, race, location.
- **Resource bars**: HP/MP/ST/EN — radial or linear; label + fill %.
- **XP bar**: current/next + tooltip "AFK regen tahmini X saat".
- **Sub-tabs**: Stats / Inventory / Equipment (embedded).
- **Primary stats**: STR/AGI/INT/VIT/CHA/LUCK — grid 3×2.
- **Secondary stats**: derived (Attack, Def, Crit, Hit, CDR, MF).
- **Buff/debuff strip**: horizontal scroll icons + remaining duration.
- **Active Title** dropdown → aktif title değiştir modal.

## 4. State

- Character sheet provider — `GET /v1/character/{id}/sheet` (compound endpoint).
- WS subscribe `character:{id}` for live update (XP tick, HP change).
- Buff/debuff subscribe → strip real-time update.

## 5. Interaction

- Stat pill tap → tooltip breakdown (base + gear + buff + mastery bonus).
- Buff icon tap → modal detail (source, duration, effect).
- Title tap → title selector modal (bkz. [progression/titles.md](../progression/titles.md)).
- Long-press resource bar → detail chart (last 24h).

## 6. API

- `GET /v1/character/{id}/sheet` — snapshot compound (stats + resources + buffs + title).
- `PUT /v1/character/{id}/title/active`
- `GET /v1/character/{id}/buffs` (opsiyonel refresh).
- WS `character:{id}` events: xp_gain, hp_change, buff_added, buff_removed, level_up.

## 7. Level Up Animation

Level_up event → full-screen modal:
- "SEVİYE ATLAD!N" + rank badge.
- +N stat gain summary.
- +M skill point (Phase 02 class skill tree UI).
- "Devam" → return to sheet.

## 8. Metrikler
- `ui_character_sheet_view_total`
- `ui_title_change_total`
- `ui_stat_tooltip_open_total`

## 9. Cross-Reference
- [character](../character/README.md), [progression/titles.md](../progression/titles.md)
- [combat/status-effects.md](../combat/status-effects.md) — buff/debuff kaynağı
