# Long-Term Progression — Katalog

> **Faz 29**: Endgame ve horizontal progression sistemleri. Level cap ötesi motivasyon.

## 1. Kapsam
Level tamamlandıktan (ve tamamlanmadan da) sonra oyuncunun **binlerce saat** motive kalmasını sağlayan hedef ve reward katmanları. AFK-first uyumlu — pek çok progression saati doğal aktivite ile artar.

## 2. Sistemler ve Rolleri

| Sistem | Ekseni | Kaynak | Hedef |
|---|---|---|---|
| **Mastery** | Aktivite derinliği | Aktivite tick (mining, combat, crafting) | Passive bonuslar per specialty |
| **Ascension** | Rebirth cycle | Level cap + ek koşul | Reset karşılığı kalıcı skill puan |
| **Achievement** | Completionist | Trigger event (ilk kill, 100k ore) | Puan + title + kozmetik |
| **Title** | Prestij | Achievement, quest, event, ascension | Görsel isim yanı |
| **Collections** | Kolleksiyoncu | Pet, mount, appearance, mount skin | UI galeri + set bonus |
| **Milestones** | Zaman/hacim | Toplam saat, gold earned, distance | Periyodik reward + hikaye momenti |
| **Reputation** | Faction bağlılığı | Quest, kill, delivery | Vendor, quest, cosmetic access |
| **Faction** | Grupla ilişki | Quest tree + reputation | Story branching + hostility rules |

## 3. Ortak Prensipler (DD-162)

- **Karar-yumuşak**: bir sistem oyuncuya tam kilit vermez; her sistem alternatif yol sunar (bir achievement kaçtı → başka içerik).
- **AFK-first uyumluluk**: mastery ve reputation ana kazançları AFK aktivite üstünden; achievement ise online moment (attend + boss kill vs) ile karışık.
- **Server-authoritative event stream**: her progression trigger `event.progression.*` bus üzerinden yazılır → durum çift yazımı olmaz.
- **Personal ledger**: her karakter/account bazında append-only history.
- **Kalıcılık**: ascension hariç geri döndürülemez; ascension bilinçli reset.

## 4. Dokümanlar
- [mastery.md](mastery.md) — Cross-system mastery + tier + passive bonus.
- [ascension.md](ascension.md) — Rebirth + kalıcı puan + prestij tier.
- [achievements.md](achievements.md) — Trigger + point + reward + category.
- [titles.md](titles.md) — Title kaynak + display + retirement.
- [collections.md](collections.md) — Pet/mount/appearance/mount-skin galeri + set bonus.
- [milestones.md](milestones.md) — Zaman/hacim eşikleri + retro reward.
- [reputation.md](reputation.md) — Faction rep tier + gain/loss + decay.
- [faction.md](faction.md) — Faction taxonomy + hostility + story branch.

## 5. Cross-Reference
- [character/level.md](../character/level.md), [character/xp.md](../character/xp.md)
- [character/reputation.md](../character/reputation.md), [character/titles.md](../character/titles.md) — v1 temel (bu klasör extend eder)
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)
- [quests](../quests/README.md), [economy/currency.md](../economy/currency.md), [pets/README.md](../pets/README.md), [mounts/README.md](../mounts/README.md)

## 6. Metrikler (v1)
- `progression_event_total{system, event_kind}`
- `progression_reward_grant_total{system, tier}`
- `mastery_tier_up_total{track}`
- `ascension_completed_total{prestige_tier}`
- `achievement_unlocked_total{category}`
- `reputation_tier_up_total{faction}`

## 7. Bekleyen Kararlar
- Ascension frequency final tune (DD-OPEN-5).
- Achievement puan currency mi cosmetic mi (v1: puan hem cosmetic hem leaderboard input).
- Collection set bonus rotation (kalıcı vs seasonal).
- Faction hostility PvP tetikleyicisi (Phase 02+ PvP ile birleşecek).
