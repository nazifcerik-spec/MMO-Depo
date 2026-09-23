# Flutter UI/UX — Katalog

> **Faz 31**: Client tarafı. Şimdiye kadar backend/gameplay engine tanımlandı; bu klasör oyuncunun **gördüğü ekranları** tarif eder.

## 1. Kapsam
Flutter tabanlı, **mobile-first** (Android/iOS), responsive (tablet + web) uygulama. Her ekran doküman: amaç, wireframe metin, state kaynağı, API bağlantısı, kritik interaksiyon.

## 2. Prensipler (DD-176)

- **Mobile-first**: 375×812 hedef (iPhone 11 / Android 5.8"); tablet + web scale-up.
- **Data-driven UI**: server hâlâ hakikat; UI reactive rebuild WS event ile.
- **Offline-tolerant**: cached snapshot ile açılış; server sync geldiğinde reconcile.
- **Localization**: TR primary, EN teknik; `arb` files; body_key + body_vars kalıbı server ↔ client (DD-171 notification uyumlu).
- **Deep link routing**: `game://route/...` her ekran adreslenebilir.
- **State management**: Riverpod (öneri; final Phase 01 client team).
- **Design token**: color/typography/spacing/radius/shadow tek yerden ([design-system.md](design-system.md)).
- **Accessibility**: min tap target 44×44, contrast ratio AA (4.5:1), screen reader semantic label.
- **Motion**: 200ms default easing (Curves.easeInOut); reveal/loot 400ms.

## 3. Bilgi Mimarisi

```
Splash → Login/Register
  ↓
Character Select
  ↓
Main Shell (bottom tab / drawer)
  ├── Home (dashboard: aktif contract, energy, günlük ödül)
  ├── Character (sheet + inventory + equipment sub-tabs)
  ├── World (map + activity + travel + profession + crafting)
  ├── Combat/PvE (combat + dungeon + party)
  └── Social (guild + friends + mail + chat + market + quest)
```

## 4. Dokümanlar
- [design-system.md](design-system.md) — Token, tipografi, ikonografi, motion.
- [login.md](login.md) — Login/register/logout/OTP.
- [character.md](character.md) — Character select + create + delete.
- [character-sheet.md](character-sheet.md) — Stats/XP/HP/energy dashboard.
- [inventory.md](inventory.md) — Grid inventory + weight + sort/filter.
- [equipment.md](equipment.md) — Paper doll + slot equip/unequip.
- [map.md](map.md) — World/region/territory/location navigator.
- [activity.md](activity.md) — AFK contract start/stop + reveal screen.
- [profession.md](profession.md) — Profession dashboard (mining/fishing/hunting/…).
- [crafting.md](crafting.md) — Recipe list + queue + critical craft alert.
- [combat.md](combat.md) — Live combat panel + AFK combat state.
- [dungeon.md](dungeon.md) — Dungeon entry + floor progress + reveal.
- [pet.md](pet.md) — Pet management + skills + happiness (v1 subset).
- [mount.md](mount.md) — Mount management + travel speed + skin.
- [party.md](party.md) — Party UI + invite + roles + chat + hunt.
- [guild.md](guild.md) — Guild dashboard + members + bank + chat + MOTD.
- [market.md](market.md) — Marketplace browse + list + buy + history chart.
- [quest.md](quest.md) — Quest log + tracker + turn-in.
- [social.md](social.md) — Friends + block + mail inbox + chat.

## 5. Cross-Reference
- [architecture/flutter.md](../architecture/flutter.md) — Client stack detay
- [systems/notification.md](../systems/notification.md) — Push + deep link
- [notification/mobile-lifecycle.md](../notification/mobile-lifecycle.md) — Background lifecycle
- Backend/domain docs — her ekranın API kaynağı orada

## 6. Metrikler (Client-side reported)
- `ui_screen_view_total{screen}`
- `ui_screen_load_ms{screen, p50, p95}`
- `ui_action_total{screen, action}`
- `ui_error_shown_total{screen, code}`
- `ui_rebuild_count_per_frame{screen}`

## 7. Bekleyen Kararlar
- State management final (Riverpod öneri).
- Design system Figma dosyası kaynak (Phase 01 art team).
- Web layout branch (mobile scale vs desktop layout) Phase 02.
- Accessibility audit final (screen reader label pass) Phase 02.
