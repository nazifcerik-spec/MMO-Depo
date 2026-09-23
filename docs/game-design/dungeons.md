# Dungeons (Zindanlar)

## 1. Amaç
Kısa, yoğun, party-orta ölçekli PvE encounter'ları modellemek. Ayrıca solo idle-run varyantı ile "farm dungeon" da destekler.

## 2. Oyuncuya Fayda
- Class kimliğinin party içindeki tatmini.
- Set item drop, tier ilerleme kaynağı.
- Sosyal etkileşim (guild party, LFG).
- Solo idle-run: aynı dungeon'un otomasyon ile daha uzun sürede farm edilebilir hali.

## 3. Gameplay Loop
- Session'da: LFG kur → dungeon manuel run (10–20 dk).
- Meta'da: haftalık lockout, dungeon mastery track, hardmode unlock.

## 4. Dungeon Türleri
- **Story Dungeon** — 1 kez per week weekly reward.
- **Farm Dungeon** — daily lockout, farm yield.
- **Heroic Dungeon** — hardmode, better reward, weekly.
- **Idle Dungeon Run** — solo, AFK, düşük tier reward ama sürekli.

## 5. Temel Kurallar
- Party size: 2–5 (varsayılan 3–5 için tune).
- Boss encounter'lar: bir "phase" fikri (Phase 1/2/enrage).
- Weekly reset: haftalık gün + saat sabit (öneri: Salı 06:00 UTC).
- Idle-run: tier azaltılmış difficulty; yield/hour manuel run'un ~%50'si.
- Ölüm/wipe: encounter reset; loot fırsatı tükenmez, encounter tekrar denenir; tam wipe → dungeon çıkışı opsiyonu.

## 6. Kullanıcı Akışı
```
Dungeon Panel → tier + difficulty seç
  ↓
Party (opsiyonel) → LFG queue veya guild party
  ↓
Enter → server instance oluşturur
  ↓ (manuel)
Boss encounters → loot roll → complete
  ↓ (idle-run)
Activity Engine tick → periodic reveal
```

## 7. İlişkiler
- [Combat](combat.md), [Party](party.md), [Loot System](../systems/loot-system.md), [Items](items.md), [Activity Engine](../systems/activity-engine.md), [Guilds](guilds.md)

## 8. Gereken Veriler
```
dungeon_template
  id, code, name, region_id
  min_level, max_level (recommended)
  party_size_min, party_size_max
  boss_encounters [] (ordered)
  reward_pool_id
  weekly_lockout (bool)
  idle_run_supported (bool)

dungeon_instance
  id, template_id, party_id, difficulty, seed
  created_at, expires_at
  current_encounter

boss_encounter
  id, dungeon_id, order_index
  phase_definitions (json)
  loot_table_id
  ai_profile_id
```

## 9. I/O
- Input: create instance, join, ready, attempt encounter, leave.
- Output: encounter events, loot rolls, completion event.

## 10. Progression
- Dungeon mastery: kill/completion sayısına göre level.
- Weekly lockout → shared loot restriction.
- Heroic unlock: normal completion X kez.

## 11. Reward
- Boss loot table: rare/epic/set item drop.
- Currency drop (dungeon token) → vendor unique item.
- Ilk kill "first clear" achievement + kozmetik.

## 12. Failure Durumları
- Wipe: encounter reset; N wipe sonrası "give up" opsiyonu.
- Party disband: instance timer devam; kalan üye solo devam edebilir (idle-run flag).
- Server-side crash: seed'li instance resume; encounter state persist.

## 13. Edge Cases
- Undergeared party: dungeon zorluk uyarısı.
- Overgeared party: idle-run yield capped.
- Wipe count leaderboard (kozmetik): "clean run" achievement.

## 14. AFK Davranışı
- Idle-run: AFK aktivitedir.
- Manuel dungeon: aktif oynanır; ancak internet kopunca 60 sn grace + rejoin.

## 15. Server Authoritative Gereksinimleri
- Instance oluşturma transaction.
- Loot roll seed'li ve audit-loglu.
- Weekly lockout server timer.

## 16. Genişletilebilirlik
- Mythic+ affix sistemi (Phase 02+).
- Raid (10+ oyuncu) alt kategori.
- Cross-guild raids.

## 17. Örnek Senaryo
3 kişilik party (Warrior/Cleric/Mage) "Ashen Keep" farm dungeon'a girer; 12 dk'da tamamlar; 2 rare item + 45 dungeon token + 8k XP her kişi.

## 18. Balance Considerations
- Manuel run yield ile idle-run yield oranı ~2:1 (aynı süre için).
- Set item drop rate: normal %5, heroic %8, mythic %12 (haftalık cap ile birleşince makul).
- Boss enrage timer: undergeared party için soft mesaj; DPS eşiği önerisi ekranda.
