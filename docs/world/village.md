# Village

## 1. Amaç
Küçük yerleşim; kısıtlı servis. City-lite. Bölgesel kaynak + basit quest hub.

## 2. Model

`world.settlement` `kind='village'`, `size='village'`.

Aynı tablo, farklı flag.

## 3. Content Pack

```yaml
# content/world/settlements/wheatcross.yaml
code: wheatcross
territory_code: ashen_flats
kind: village
size: village
display_name_key: village.wheatcross.name
lore_key: village.wheatcross.lore

spawn_point:
  x: 100
  y: 50
  heading: 0

services:
  bank: false                     # opsiyonel small bank Phase 02
  marketplace: false              # sadece city
  guild_hall: false
  inn: true
  waygate: false                  # bazı köylerde var Phase 02
  vendors:
    - food
    - materials                   # local goods only
  trainers:
    - profession                  # sadece basic trainer (tier 1-30)
  quest_hub: true
  home_city_available: false

specialty: fishing                # opsiyonel — village-özel
ambient_theme: rural_daytime
```

## 4. Village Servis Kapsamı

City'nin subset'i:
- ✓ Inn (rest — Phase 02).
- ✓ Basic vendor (food, materials).
- ✓ Basic profession trainer (tier 1-30).
- ✓ Quest hub (bölgesel quest'ler).
- ✗ Bank (opsiyonel Phase 02).
- ✗ Marketplace (sadece city).
- ✗ Guild hall.
- ✗ Waygate (opsiyonel Phase 02).
- ✗ Home city.
- ✗ Advanced trainer.

## 5. Specialty

Village-özel özel drop:
- Fishing village → fresh water fish exclusive.
- Farming village → grain/herb exclusive material.
- Border village → faction rep quest specific.

## 6. Population Density

Village visual olarak küçük (10-20 NPC vs city 100+). Layer/sharding gerekli değil (düşük yoğunluk).

## 7. Endpoint

**GET `/v1/content/world/settlements/{code}`**
City'nin aynı endpoint; `kind` alanı `village` döner.

## 8. Village Development (Phase 02+)

Village upgrade mekaniği (Phase 02 fikri):
- Guild territory war → village bank service unlock.
- Faction reputation → village exclusive quest chain açılır.
- Season event → temporary vendor.

v1'de static content.

## 9. Metrics
- `village_active_visitors{village}` gauge.
- `village_quest_completed_total{village}`.

## 10. Cross-Reference
- [City](city.md), [Outpost](outpost.md), [Territory](territory.md)
- [game-design/cities.md](../game-design/cities.md)
- [database/world.md §4](../database/world.md)

## 11. Bekleyen Kararlar
- Village bank enable Phase 02 (basic tier).
- Village waygate expansion Phase 02.
- Village-guild territory bonus Phase 02.
