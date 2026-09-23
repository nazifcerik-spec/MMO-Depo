# Map (World / Region / Location Navigator)

## 1. Amaç
Dünya haritası üzerinden yolculuk hedef seçimi, POI (resource node, monster area, dungeon, city), aktif event marker, travel initiate.

## 2. Layout (Mobile)

```
┌────────────────────────────┐
│ ← Harita        🔍  📌     │
├────────────────────────────┤
│                            │
│    [interactive map]       │
│    [region → territory     │
│      → location hierarchy] │
│                            │
│  📍 sen: Ironhold          │
│                            │
├────────────────────────────┤
│ Aktif:                     │
│ 🔥 World Boss (Blackspire) │
│    2 saat 14 dk kaldı      │
├────────────────────────────┤
│ [Hedef Seç]  [Aktif Et]   │
└────────────────────────────┘
```

## 3. Zoom Levels

- **World**: continent hierarchy, kıtalar.
- **Continent**: region.
- **Region**: territory → settlement, dungeon, POI.
- **Territory**: location detail; resource node ve monster area cluster'ları.

Pinch zoom / +/- buttons.

## 4. POI Türleri

| İkon | POI |
|---|---|
| 🏰 | City |
| 🏘 | Village / Outpost |
| ⛏ | Resource node cluster (mining) |
| 🌲 | Woodcutting |
| 🎣 | Fishing spot |
| 🌿 | Herbalism |
| 🏹 | Hunting ground |
| 🕳 | Dungeon entrance |
| ⚔ | Monster area |
| 🔥 | Active world event |
| ⭐ | Quest marker |

Tap POI → info sheet (name + level range + activity type + travel time estimate).

## 5. Travel Initiate

Location tap → "Buraya Yolculuk" CTA → travel confirm modal:

```
Ironhold → Wildwood
Mesafe: 4200 m
Süre (mount+equipment+weather): 28 dk
Enerji maliyeti: 12 EN
Yol: paved (bonus -10%)
──────────────
[İptal]        [Başlat]
```

Confirm → `POST /v1/travel/start` → travel state ekranına geçiş.

## 6. Travel Progress Bar

Global banner (üstte):
```
🐎 Yolculukta → Wildwood ─▓▓▓▓░░ 62% · 10dk kaldı [Detay]
```

Tıklama → travel detail (bkz. [travel/README.md](../travel/README.md)).

## 7. State

- Map provider — cached world hierarchy (content pack; hot reload).
- Character location subscribe → own marker update.
- Active events subscribe → marker update.

## 8. API

- `GET /v1/world` — hierarchy.
- `GET /v1/world/region/{id}` — detail.
- `GET /v1/world/location/{id}/travel-preview?from_location_id=...` — mesafe/süre estimate.
- `POST /v1/travel/start`

## 9. Metrikler
- `ui_map_view_total{zoom_level}`
- `ui_map_poi_tap_total{poi_kind}`
- `ui_travel_initiate_total`

## 10. Cross-Reference
- [world](../world/README.md), [travel](../travel/README.md), [world-events](../world-events/README.md)
