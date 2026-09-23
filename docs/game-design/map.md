# Map (Harita)

## 1. Amaç
Dünyayı oyuncuya görsel olarak sunmak; keşif, planlama ve travel arayüzü olarak çalışmak.

## 2. Oyuncuya Fayda
- Bir yerin nerede olduğunu görme.
- Aktivite planlama (nereye gideyim → orada ne var).
- Discovery çabası ve ödülü (fog of war Phase 02).

## 3. Gameplay Loop
- Sık açılır: hangi hunt spot / gather node / dungeon nerede?
- Meta'da keşif tamamlanır, achievement.

## 4. Katmanlar
- **World map** (kıtalar)
- **Continent map**
- **Region map** — POI icon overlay
- **Territory / City map**

## 5. Temel Kurallar
- Filtreler: profession node, dungeon, quest hub, vendor, guild hall.
- Waypoints: oyuncu kişisel işaretler; guild shared waypoints (Phase 02).
- Fog of war: keşfedilmemiş POI görünmez (Phase 02); Phase 00'da hepsi görünür ancak "not visited" işaretli.

## 6. Kullanıcı Akışı
```
Map open → filter set → POI click → info panel → "Travel here" button
```

## 7. İlişkiler
- [World](world.md), [Travel](travel.md), [Cities](cities.md), [Quests](quests.md)

## 8. Gereken Veriler
```
poi { id, location_id, category, coords }
waypoint { id, character_id, name, location_id, is_shared }
map_tile { region_id, tile_ref, art_ref }
```

## 9. I/O
- Input: filter, waypoint add, travel request.
- Output: rendered map layer, POI info popup.

## 10. Progression
- Region discovery: %100 açınca achievement.
- Special hidden POI'ler unlock ile ortaya çıkar.

## 11. Reward
Discovery mini xp + achievement.

## 12. Failure Durumları
- Locked POI: click → unlock koşulu görünür.
- Waypoint limit doldu: yeni eklemek için biri silinmeli.

## 13. Edge Cases
- POI relocation (patch): waypoint'lerin auto-fallback.
- Instance-only POI (dungeon): map'te görünür ama travel doğrudan LFG'ye açılır.

## 14. AFK Davranışı
Harita statik; AFK ile ilişki yok. Travel aktivitesi tetiklendiğinde AFK olur.

## 15. Server Authoritative Gereksinimleri
- POI ownership değil; kişisel waypoint kaydı.
- Fog of war state karakter başına persist.

## 16. Genişletilebilirlik
- Cartography profession (Phase 02) map notes create.
- Guild map (Phase 02).

## 17. Örnek Senaryo
Karakter map açar, "Mining Node" filter, en yakın iron node'a travel start.

## 18. Balance Considerations
- Map'ten "travel here" fonksiyonu ekonomik olarak ücretlendirilir (opsiyonel gold) veya süre bazlıdır.
- Waypoint sayısı sınırı: 20 kişisel (Phase 02'de expand).
