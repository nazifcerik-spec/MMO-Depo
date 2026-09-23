# Pet Management

## 1. Amaç
Pet listesi, aktif pet seçim, level/XP, trait, skill, evolution (Phase 02 detay), happiness/stamina (Phase 02 detay).

## 2. Layout — Pet List

```
┌────────────────────────────┐
│ ← Petler       ⚙          │
├────────────────────────────┤
│ Aktif: 🐺 Fenrir Lvl 24    │
│ [Aktif Değiştir]           │
├────────────────────────────┤
│ Sahip Olunan (7/30):        │
│  🐺 Fenrir     Lvl 24 ★★★  │
│  🐈 Whiskers   Lvl 12 ★★   │
│  🦉 Athena     Lvl 8 ★     │
│  ...                       │
│ [Katalog Gör]              │
└────────────────────────────┘
```

## 3. Pet Detail

- Species + rarity + level + XP bar.
- Trait list (species-specific + individual).
- Skill list (unlocked at level).
- Evolution track (Phase 02).
- Happiness / Stamina (Phase 02).
- Actions: **Aktif Yap**, **Yeniden Adlandır**, **Salma** (release, permanent).

## 4. State

- Pet provider — `GET /v1/character/{id}/pets` + WS `pet:{id}` events.
- Active pet — `character.active_pet_instance_id`.

## 5. API

- `GET /v1/character/{id}/pets`
- `GET /v1/character/{id}/pet/{instance_id}`
- `POST /v1/character/{id}/pet/{instance_id}/activate`
- `PUT /v1/character/{id}/pet/{instance_id}/rename`
- `DELETE /v1/character/{id}/pet/{instance_id}` (release, confirm modal)

## 6. Cross-Reference
- [pets](../pets/README.md), [pets/species.md](../pets/species.md)

## 7. Bekleyen
- Evolution UI Phase 02.
- Happiness/stamina UI Phase 02 (v1 dokümantasyon eksik — Pet Phase 20 gap).
- Pet combat inspector Phase 02.
