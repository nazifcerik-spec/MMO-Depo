# World Bosses (Dünya Boss'ları)

## 1. Amaç
Sunucu genelinde community event boyutunda karşılaşmalar sunmak. Sosyal fokus, guild organizasyonu ve prestij drop.

## 2. Oyuncuya Fayda
- MMO'nun "büyük an" hissi.
- Nadir loot, kozmetik, sıralamaya girmek.
- Guild etkinliği için ideal hedef.

## 3. Gameplay Loop
- Meta ölçekli. Spawn timer + region.

## 4. Boss Türleri
- **Regional World Boss** — bölgesel; haftalık respawn; katılım açık.
- **Continental Champion** — nadir; 2 haftada bir; büyük ödül.
- **Server Titan** (Phase 02) — 1 ay/1 sezon; efsanevi drop, unique first-kill kozmetik.

## 5. Temel Kurallar
- Spawn: pencere içinde random (announced range).
- Boss HP: server-scale; katılım sayısına göre soft-scale (10–200 oyuncu).
- Contribution scoring: damage + heal + assist + participation time.
- Loot: kill credit + participation threshold (min damage or min time).
- Weekly cap: bir karakter belirli boss'tan bir kez rare loot roll edebilir.

## 6. Kullanıcı Akışı
```
Notification: "Boss spawn imminent in [region]"
 ↓
Travel → coordinate → fight
 ↓
Kill → loot roll → announced kill credit
```

## 7. İlişkiler
- [Combat](combat.md), [Guilds](guilds.md), [Notification System](../systems/notification-system.md), [Loot System](../systems/loot-system.md), [World](world.md), [Events](events.md)

## 8. Gereken Veriler
```
world_boss_template
  id, code, name, region_id, tier
  hp_base, hp_scale_formula, damage_profile
  ability_ids [], loot_table_id
  spawn_windows [] (recurrence rules)

world_boss_spawn
  id, template_id, spawned_at, defeated_at (nullable)
  participants [] { char_id, contribution, kill_credit }
```

## 9. I/O
- Input: engage boss, contribute.
- Output: contribution log, kill event, loot payout, world announcement.

## 10. Progression
- Boss achievements (kill counts, first strike, execute).
- Boss-only mount/pet drops.

## 11. Reward
- Gold, materyal, epic+ item drops.
- Cosmetic (title, dye, mount).
- Guild XP contribution.

## 12. Failure Durumları
- Enrage timer hit: boss despawn, participants küçük iade + notification.
- Under-participation: soft-cancel spawn (Phase 02 policy).
- Server-side disconnect during fight: contribution korunur if reconnected within 5 min.

## 13. Edge Cases
- Zerg outnumber: soft-cap effective participants (200+ diminishing returns).
- Under-attack denied: boss uzun süre unengaged → repositions, spawns adjacent quest.
- Griefing (aggro drag): server anti-pattern.

## 14. AFK Davranışı
World boss AFK oynanamaz — manuel katılım gerekir. Ancak spawn notification "5 dk önce gelir" ki oyuncu login olabilsin.

## 15. Server Authoritative Gereksinimleri
- Contribution ledger server.
- Loot roll deterministik + auditable.
- Announcement server broadcast.

## 16. Genişletilebilirlik
- Coop mechanics (Phase 02: multi-phase, positional).
- Cross-guild raid (Phase 02).
- Server-first system with permanent leaderboard.

## 17. Örnek Senaryo
"Ashen Wyrm" boss Ashfield civarında spawn; 5 dk önce anons; 48 oyuncu katıldı; kill 12 dk; contribution top 20 rare item drop; guild "First Kill" achievement.

## 18. Balance Considerations
- Casual friendly: min participation threshold makul; solo player da katılıp reward alabilmeli (küçük ölçekli).
- Guild advantage: coordination bonusu (+5% dmg for grouped party), zorunlu değil.
- Loot rarity: rare %8, epic %3, legendary %0.5.
