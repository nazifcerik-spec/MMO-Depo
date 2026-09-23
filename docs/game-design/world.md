# World (Dünya)

## 1. Amaç
Oyunun mekan hiyerarşisini tanımlamak; her lokasyonun bir kimliği ve fonksiyonu olmasını sağlamak.

## 2. Oyuncuya Fayda
- Keşfin ödüllendirici olması.
- Bölgeler arasında ekonomik ve içerik farkları (gitmeye değmesi).
- Her lokasyonun "burada ne yapılır" cevabı olması.

## 3. Gameplay Loop
- Session'da lokasyon seçimi.
- Meta'da yeni bölge unlock; kıta ilerleme; faction rep.

## 4. Hiyerarşi
```
World
  Continent
    Region
      Territory
        City / Village / Outpost
          Location (POI, resource node, dungeon entrance, quest hub)
```

- **World**: tek dünya, tek shard (Phase 01).
- **Continent**: farklı biyom setleri; ilerleme yolları.
- **Region**: seviye bandı DEĞİL kimlik (klimat, kültür, kaynak).
- **Territory**: political bölünüm (guild territory Phase 02 hedefi).
- **City / Village / Outpost**: yerleşim tipi ölçek farkı.
- **Location / POI**: bir aktivitenin "adresi".

## 5. Temel Kurallar
- Karakter herhangi bir zamanda bir `current_location_id` taşır.
- Bir aktivite başlatmak için lokasyon uygun olmalı (adventure X location'a bağlı).
- Bölge unlock: bazı bölgeler quest + level + fiziksel travel gate ile açılır.
- Sub-region hava/durum farkları (Phase 02 opsiyonel).

## 6. Kullanıcı Akışı
```
World Map → zoom in continent → region → territory → city
  ↓
Location seç → travel başlat
```

## 7. İlişkiler
- [Map](map.md), [Cities](cities.md), [Travel](travel.md), [Hunting](hunting.md), [Adventure](adventure.md), [Professions](professions.md), [Quests](quests.md)

## 8. Gereken Veriler
```
continent { id, code, name }
region { id, continent_id, name, biome, unlock_condition_json }
territory { id, region_id, name, faction_id (nullable) }
city { id, territory_id, name, services_json }
village { id, territory_id, name, services_json }
outpost { id, territory_id, name, purpose }
location { id, parent_id, parent_type, name, tags[], activity_hooks[] }
```

## 9. I/O
- Input: navigate map, request travel.
- Output: available activities at location, ambient descriptors.

## 10. Progression
- Continent unlock: milestone (level + quest + faction).
- Region "discovery" achievement.
- Location tag'leri açılınca ilgili aktiviteler görünür.

## 11. Reward
Dünya, ödül vermez; ödüllerin sahnesidir. Discovery bonus xp Phase 02.

## 12. Failure Durumları
- Locked region'a travel: reddedilir, unlock listesi gösterilir.
- Kaotik zone (Phase 02 PvP açılırsa): opt-in warning.

## 13. Edge Cases
- Location kapatılması (patch): oyuncu safe city'ye teleport, mail notification.
- Time-limited seasonal region (Phase 02): expire olunca arşive.

## 14. AFK Davranışı
Dünya travel aktivite ile birleşir; AFK travel süresi tick'lenir.

## 15. Server Authoritative Gereksinimleri
- Location transitions server.
- Region unlock kontrol server.

## 16. Genişletilebilirlik
- Yeni kıta = data set.
- Instance-only zones (raid).
- Housing plots (Phase 02+): territory içinde slot.

## 17. Örnek Senaryo
"Emerald Coast" region: kıyı, fishing/foraging bonus, orta-tier hunting. Kıyı outpost'unda dock var → aquatic travel start.

## 18. Balance Considerations
- Content dağılımı: her region'da min 1 dungeon + 2 hunt spot + 3 gather node type + 2 adventure lane.
- Bölge dengesi: her tier band için en az 2 bölge; oyuncunun "tek yol" hissetmemesi.
