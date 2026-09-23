# Achievements (Başarımlar)

## 1. Amaç
Uzun soluklu meta-hedefler; oyuncunun ilerlemesini görselleştiren koleksiyon katmanı.

## 2. Oyuncuya Fayda
- "Ne yaptım?" görünürlüğü.
- Kozmetik unlock (title, badge).
- Uzun vadeli motivasyon.

## 3. Gameplay Loop
- Sessions'da pop-up tetiklenir; meta'da liste tamamlama hedefi.

## 4. Kategoriler
- Combat (kill counts, boss firsts)
- Professions (materyal biriktirme, master craft)
- Exploration (region discovery, POI)
- Dungeons (clear counts, no-death runs)
- Events (season participation)
- Social (guild growth, party runs)
- Meta (level ups, ascension milestones)
- Collections (mount, pet, transmog Phase 02)

## 5. Temel Kurallar
- Achievement point: kategori bazlı bütçe; toplam gösterge.
- Bazı achievement'lar hidden (unlock condition sadece tamamlandığında görünür).
- Server-first: sunucuda ilk tamamlayan özel kozmetik.

## 6. Kullanıcı Akışı
```
Achievement panel → kategori → liste (locked/progress/complete)
 ↓
Detail: kriter, ödül
```

## 7. İlişkiler
- Tüm sistemler tetikleyici; özellikle [Events](events.md), [World Bosses](world-bosses.md), [Professions](professions.md), [Dungeons](dungeons.md), [Progression](progression.md)

## 8. Gereken Veriler
```
achievement
  id, code, category, points, hidden (bool)
  criteria_json, reward_json

character_achievement
  character_id, achievement_id, progress_json, completed_at
```

## 9. I/O
- Input: internal event trigger.
- Output: pop-up, reward payout.

## 10. Progression
Kategorilerde ilerleme çubuğu; server-first için leaderboard.

## 11. Reward
Title, kozmetik, achievement point, occasional gold/token.

## 12. Failure Durumları
- Hatalı trigger: retro-fix job.
- Server-first çekişme: timestamp resolution.

## 13. Edge Cases
- Achievement patch değişikliği: progression migrate.
- Rolled-back exploit: achievement geri alınabilir (nadir).

## 14. AFK Davranışı
AFK aktivite tetikleyiciler achievement'ı da tetikler.

## 15. Server Authoritative Gereksinimleri
- Progression ledger server; hile korumasız değil.

## 16. Genişletilebilirlik
- Meta-achievement (Phase 02): achievement kombine ödül.
- Cross-character account achievement.

## 17. Örnek Senaryo
"Ashen Wyrm Slayer" — first kill within server; unique title + tabard.

## 18. Balance Considerations
- Achievement ödülleri kozmetik; power değil.
- Hidden achievement oranı %10; keşif tatmini.
- Uzun süre biriken achievement point'lerin bir "point vendor" karşılığı (Phase 02).
