# Dungeon

## 1. Amaç
Dungeon giriş, floor progress, AFK run izleme, reveal.

## 2. Layout — Dungeon List

```
┌────────────────────────────┐
│ ← Zindanlar                │
├────────────────────────────┤
│ Iron Mine Lvl 20           │
│  Normal | Heroic | Mythic  │
│  Weekly lockout: Reset 6d  │
│  [Gir Solo] [Party]        │
├────────────────────────────┤
│ Ashen Peak Lvl 35          │
│  Normal | Heroic           │
│  [Gir]                     │
└────────────────────────────┘
```

## 3. Dungeon Detail

```
┌────────────────────────────┐
│ ← Iron Mine (Heroic)       │
│ Boss: Rusted Warden         │
│ Kat: 5 · Süre ~45 dk        │
├────────────────────────────┤
│ Loot preview:               │
│  Uncommon: 4 slot           │
│  Rare: 2 slot               │
│  Epic: 1 (boss %35)         │
│  Legendary: %2              │
├────────────────────────────┤
│ Party gerekli değil (solo) │
│ Enerji: 30                  │
│ [Solo Gir] [Party ile Gir]  │
│ [AFK Modu (Auto Combat)]    │
└────────────────────────────┘
```

## 4. AFK Dungeon Progress

```
Aktif Zindan · Iron Mine
Kat 3/5 · Room 12/20
▓▓▓░░ 32%   ~30dk kaldı
Şimdiye kadar: 24 kill, 3 trap defuse
[Detay]  [İptal (loot kayıp)]
```

## 5. Live Dungeon Combat

Live combat panel (bkz. [combat.md](combat.md)) + üstte floor/room breadcrumb + minimap küçük.

## 6. Reveal (AFK)

Dungeon reveal:
- Boss defeat animation (video / spine).
- Loot listesi (personal roll — DD-115).
- XP + achievement + first-clear bonus.

## 7. State

- Dungeon templates provider — `GET /v1/dungeon/templates?level_max=`.
- AFK run provider — `GET /v1/character/{id}/dungeon/active` + WS updates.

## 8. API

- `GET /v1/dungeon/templates`
- `POST /v1/dungeon/{template_code}/enter` body: `{ difficulty, party_id?, afk_mode }`
- `POST /v1/dungeon/{run_id}/afk-progress` (server-side background)
- `POST /v1/dungeon/{run_id}/cancel`
- `POST /v1/dungeon/{run_id}/reveal`

## 9. Metrikler
- `ui_dungeon_list_view_total`
- `ui_dungeon_enter_total{template, difficulty, afk_mode}`
- `ui_dungeon_reveal_open_total`

## 10. Cross-Reference
- [dungeon](../dungeon/README.md), [dungeon/afk-run.md](../dungeon/afk-run.md)
- [combat.md](combat.md), [party.md](party.md)
