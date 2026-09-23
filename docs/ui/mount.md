# Mount Management

## 1. Amaç
Mount listesi, aktif mount, travel speed etkisi, skin/kozmetik, carry capacity ek.

## 2. Layout — Mount List

```
┌────────────────────────────┐
│ ← Bineği̇ler                │
├────────────────────────────┤
│ Aktif: 🐎 Shadowmane Lvl 18│
│  Hız +25%  Kapasite +5 kg  │
│  Terrain bonus: paved/plain│
│ [Aktif Değiştir]           │
├────────────────────────────┤
│ Sahip (5/20):              │
│  🐎 Shadowmane Lvl 18       │
│  🦌 Elk        Lvl 8        │
│  🐫 Camel      Lvl 12       │
│  🦎 Riding Lizard Lvl 5     │
│  ...                       │
└────────────────────────────┘
```

## 3. Mount Detail

- Species + rarity + level + XP.
- Stats: base speed, terrain multiplier (per terrain), carry capacity.
- Skill (mount ability — dash, jump).
- Skin (skin change modal).
- Actions: **Aktif Yap**, **Rename**, **Skin Değiştir**, **Sat/Trade** (Phase 02).

## 4. State

- Mount provider — `GET /v1/character/{id}/mounts`.
- Active mount — `character.active_mount_instance_id`.

## 5. API

- `GET /v1/character/{id}/mounts`
- `POST /v1/character/{id}/mount/{instance_id}/activate`
- `PUT /v1/character/{id}/mount/{instance_id}/skin` body: `{ skin_code }`
- `GET /v1/mount/skins?mount_species=` — unlocked skins list.

## 6. Terrain Bonus Visualization

Detail'de terrain grid: her terrain (plain / forest / mountain / desert / swamp / snow) için multiplier badge (kırmızı slow, yeşil fast).

## 7. Cross-Reference
- [mounts](../mounts/README.md), [travel](../travel/README.md), [progression/collections.md](../progression/collections.md)
