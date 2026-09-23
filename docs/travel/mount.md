# Mount (Speed Contribution)

## 1. Amaç
Aktif mount'un travel hızına katkısı.

Yüksek seviye: [game-design/mounts.md](../game-design/mounts.md).

## 2. Model

Mount template (`content.mount_template`):
- `speed_multiplier NUMERIC(4,2)` — 1.5 / 2.0 / 2.5 tier'ları.
- `cargo_bonus_slots SMALLINT` — inventory cap bonus.
- `type` — ground / cargo / aquatic / winged / cosmetic.

Character `active_mount_id` → `companion.mount_owned` → template lookup.

## 3. Speed Multiplier Etkisi

```
if character.active_mount_id != null:
    mount = load(character.active_mount_id)
    mount_multiplier = mount.template.speed_multiplier
else:
    mount_multiplier = 1.0  (walking)
```

Base tiers (v1 önerisi):
```
tier 1 (basic): 1.5×
tier 2 (mid):   2.0×
tier 3 (fast):  2.5×
```

Cap: `mount_multiplier <= 3.0` (Phase 02 rare mount).

## 4. Mount Type Kısıtları

| Type | Terrain compatible | Notlar |
|---|---|---|
| **Ground** | grass, road, forest, hills | Water/mountain penalty |
| **Cargo** | grass, road (limited) | Slow ancak inventory bonus |
| **Aquatic** | water, coast | Su üstünde full speed; kara ×0.5 |
| **Winged** (Phase 02) | all + air | Cross-terrain best |
| **Cosmetic** | all | Sadece görsel; speed base |

Terrain-mount mismatch → speed penalty (bkz. [terrain.md](terrain.md)).

## 5. Cargo Mount Bonus

`cargo_bonus_slots`:
- Envanter kapasitesi + N slot travel süresince.
- Ancak speed_multiplier düşük (1.2×).
- Trade-off: yavaş ama daha çok taşı.

Aktive: mount summon anında; dismount ile normale döner.

## 6. Mount Constraint (Travel Method)

Travel `method='mount'` seçildiyse:
- `character.active_mount_id NOT NULL` gerekir.
- Aksi 400: `travel.no_active_mount`.

`method='walking'`:
- Mount ignore edilir; walking baseline.

## 7. Combat Dismount

Ambush event (Phase 02):
- Combat başlarsa mount dismount + cooldown.
- Travel pause; combat resolve sonrası resume.

## 8. Multi-Mount Kanal (Phase 02+)

`activity.contract.channel = 'mount_travel'` alternate channel:
- Primary aktivite gathering; travel paralel.
- v1'de kapalı (primary tek aktivite).

## 9. Endpoint

Mount aktivasyonu:
**POST `/v1/characters/{id}/mount/activate`**
```json
{ "mount_instance_id": "..." }
```
Bkz. [game-design/mounts.md](../game-design/mounts.md).

Aktivasyon travel süresini etkiler; ancak travel mid-flight mount değişimi (Phase 02):
- v1'de yasak; travel bitene kadar mount fixed (snapshot).

## 10. Snapshot Dahil

Travel contract snapshot'ında `mount_speed_multiplier` alanı dondurulur. Travel süresince mount değişse bile duration korunur.

## 11. Fail Modları
- `travel.no_active_mount` (method=mount ancak mount yok).
- `mount.incompatible_terrain` (aquatic mount karada uzun süre).

## 12. Metrics
- `travel_mount_type_distribution` gauge.
- `travel_mount_speed_avg` gauge.

## 13. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Terrain](terrain.md)
- [game-design/mounts.md](../game-design/mounts.md)
- [companion/mount_owned](../database/companions.md)

## 14. Bekleyen Kararlar
- Winged mount v1 kapsam dışı; Phase 02.
- Mount ability (dismount jump, sprint burst) Phase 02.
- Mount barding equip (mount armor) Phase 02.
- Rare mount speed cap final (3.0× v1).
