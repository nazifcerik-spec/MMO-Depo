# Guilds (Loncalar)

## 1. Amaç
Kalıcı, orta-büyük ölçekli sosyal birim. Ortak progression, bank, chat, coop aktiviteler için ana yapı.

## 2. Oyuncuya Fayda
- Sürekli sosyal ev.
- Ortak hedefler (guild boss, guild research).
- Kaynak paylaşımı (guild bank).
- Kimlik (guild tag, tabard).

## 3. Gameplay Loop
- Session'da: chat, günlük guild quest, guild bank etkileşimi.
- Meta'da: guild level, guild hall, guild war (Phase 02).

## 4. Temel Kurallar
- Guild size: 50 başlangıç; guild level ile 100'e kadar (Phase 02).
- Guild ranks: leader, officer, veteran, member, initiate (customizable Phase 02).
- Guild bank: tab bazlı; rank permissions per tab.
- Guild treasury (gold): shared purse.
- Guild XP: aktivite katkılarından biriktirilir.
- Guild disband: leader tarafından; officer confirmation.

## 5. Kullanıcı Akışı
```
Guild Panel → tabs (roster / bank / research / quests / events)
 ├─ Chat
 ├─ MOTD
 ├─ Contribution ledger
```

## 6. İlişkiler
- [Social](social.md), [Party](party.md), [Economy](economy.md), [Progression](progression.md), [Cities](cities.md) (guild hall), [World Bosses](world-bosses.md), (future) Guild War, Territory War

## 7. Gereken Veriler
```
guild
  id, name (unique), tag (unique short)
  leader_char_id, created_at
  level, xp
  treasury_gold, treasury_currency (json)
  motd, description
  home_city_id (guild hall parent)
  disband_scheduled_at (nullable, soft-delete)

guild_member
  guild_id, character_id, rank, joined_at, contribution_points, last_active_at

guild_bank_tab
  guild_id, tab_index, name, permissions_json

guild_research (Phase 02)
  guild_id, research_id, level, xp
```

## 8. I/O
- Input: create, invite, promote, demote, kick, deposit, withdraw, MOTD update, disband.
- Output: guild event stream, chat, contribution changes.

## 9. Progression
- Guild XP: member aktivitelerinden.
- Guild level → tab unlocks, member cap, guild research unlocks (Phase 02).
- Guild hall upgrade (Phase 02).

## 10. Reward
- Guild bank share.
- Guild research passive: örneğin +2% guild-wide gathering yield.
- Guild boss loot (Phase 02).

## 11. Failure Durumları
- Leader inaktif: officer auto-promote (30 gün rule).
- Insufficient contribution: rank downgrade opsiyonel (leader tercihi).
- Disband: 48 saat grace, tüm mal iadesi/split policy'ye göre.

## 12. Edge Cases
- Name collision: reject.
- Mass exit (raid quit): guild disband olmaz; sadece member azalır.
- Bank exploit (rapid deposit-withdraw): audit + rate limit.

## 13. AFK Davranışı
Guild aktiviteleri member karakterlerinin aktivitelerinden pasif olarak beslenir. Guild research (Phase 02) offline devam eder.

## 14. Server Authoritative Gereksinimleri
- Guild state server; membership index'li.
- Bank tab permissions rank-based enforce edilir.
- Contribution log audit-safe.

## 15. Genişletilebilirlik
- Guild War (Phase 02+): saldırgan/savunmacı.
- Territory Control (Phase 02+).
- Guild Alliance (Phase 03).
- Guild housing / hall customization (Phase 02).

## 16. Örnek Senaryo
"Emerald Wardens" guild, 32 member, level 6; günlük guild quest complete → 800 guild XP; guild bank tab 3 unlock cost 25k gold.

## 17. Balance Considerations
- Guild bonusları küçük ancak sürekli hissettirilmeli (2–5% aralığı).
- Guild size ile power arasındaki ilişki dengeli; büyük guild dominance önlenmeli.
- Guild research + guild bank + guild hall ekosistem: pay-off süresi haftalar mertebesi.
