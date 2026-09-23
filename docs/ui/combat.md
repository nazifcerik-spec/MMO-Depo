# Combat

## 1. Amaç
Live combat panel (aktif dövüş) ve AFK combat state (hunting adventure) UI.

## 2. Live Combat Layout

```
┌────────────────────────────┐
│ ← Combat · Goblin Raider   │
├────────────────────────────┤
│ Hedef:                     │
│ Goblin Raider Lvl 15       │
│ HP ▓▓▓▓▓▓▓░░░ 720/1000   │
├────────────────────────────┤
│    [enemy sprite]          │
│    [combat log strip]      │
│  ⚔ 42 dmg  ★ crit          │
│  🩸 poison 12 tick         │
├────────────────────────────┤
│ Sen:                       │
│ HP ▓▓▓▓▓░░ 700/950         │
│ MP ▓▓▓░░░ 150/300          │
├────────────────────────────┤
│ [⚔ Auto] [Skill1][2][3][4] │
│ [Pot ↓]   [Kaç]  [Pause]   │
└────────────────────────────┘
```

## 3. Bileşenler

- **Target frame**: HP bar, level, class icon, buff/debuff.
- **Combat log strip**: son 5 aksiyon (damage/heal/crit/miss).
- **Self frame**: HP/MP/ST bar.
- **Skill bar**: 4-8 skill slot, cooldown radial, resource cost.
- **Consumable slot**: potion quick-use (2 slot).
- **Escape/pause**: escape roll or PvE de kaç.

## 4. AFK Combat (Hunting Adventure)

Live combat yerine "hunting AFK" state:
```
Aktif AFK Combat · Wildwood Boar Pack
▓▓▓▓░░ 2s 15dk / 4s
Şimdiye kadar: 42 kill · 18 gold · 3 uncommon
[Detay]  [Durdur]
```

Reveal (activity ekranına benzer): loot listesi + combat highlights (biggest crit, rare drop moment).

## 5. Auto-Combat Toggle

Live combat'te "Auto" toggle — skill AI rotasyon (mastery-tier göre) otomatik cast. Kullanıcı manuel skill kullanırsa auto pauze; timeout sonra resume.

## 6. State

- Combat session provider — WS `combat:{session_id}` events (damage_dealt, damage_taken, skill_cast, buff_add, death).
- Skill cooldown local timer (server ack ile senkron).
- Optimistic skill cast; server reject → toast + rollback.

## 7. API

- `POST /v1/combat/start` (hunt / dungeon start)
- `POST /v1/combat/{session_id}/skill` body: `{ skill_code, target_id }`
- `POST /v1/combat/{session_id}/consumable` body: `{ item_instance_id }`
- `POST /v1/combat/{session_id}/flee`

## 8. Death Screen

Karakter öldüğünde:
```
YENİLDİN
Death penalty: XP -5%, durability -10%
[Ölüm Yerinde Diril (10dk)] [Şehirde Diril]
```

Bkz. [combat/death-revive.md](../combat/death-revive.md).

## 9. Metrikler
- `ui_combat_view_total`
- `ui_combat_skill_cast_total{skill}`
- `ui_combat_auto_toggle_total`
- `ui_combat_death_total`

## 10. Cross-Reference
- [combat](../combat/README.md), [skill-engine](../skill-engine/README.md)
- [hunting-adventure](../hunting-adventure/README.md), [activity.md](activity.md)
