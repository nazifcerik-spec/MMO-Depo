# World Events — Katalog

> **Faz 27**: Dünya çapında veya bölge çapında schedule/trigger'lı, birden çok oyuncunun aynı takvimde katıldığı olaylar.

## 1. Kapsam
Klasik MMO event türleri; **AFK-first** yorumlama ile: her event türü **online katılım + AFK katılım** kanalları sunar. World Event **shared timer + personal contribution instance** paterni kullanır (DD-091 ile uyumlu — competition frustration önlenir; timeline şareddir).

## 2. Türler

| Tür | Süre | Kapsam | Trigger | AFK Katılım |
|---|---|---|---|---|
| **World Boss** | 30-90 dk aktif window | Global (v1 tek shard) | Cron / rotation | Pre-registration + AFK "attend" contract |
| **Invasion** | 2-6 saat | Region | Cron / dynamic | Region-active char'lar auto-eligible; AFK ticks scale |
| **Regional Event** | 1-3 saat | Region/Territory | Periodic (haftada 1-3 kez) | AFK contract "participate" seçilebilir |
| **Seasonal Event** | 2-6 hafta | Global | Manual (event pack) | Login bonus + AFK contract "seasonal quest" |
| **Dynamic Event** | 15-60 dk | Location/Region | Emergent (chain trigger) | Sınırlı AFK; online player weight yüksek |

## 3. AFK Katılım Prensibi (DD-150)

- **Shared timer** — tüm oyuncular aynı server zamanında görür (`starts_at, ends_at UTC`).
- **Personal contribution instance** — her karakter için ayrı contribution ledger; başka oyuncuyla rekabet yok.
- **Reveal at open + at close** — event boyunca AFK ise `activity-engine` analytical hesap uygular: kaç tick attı, ne kadar contribution biriktirdi.
- **Reward tier** contribution eşiklerine göre (tier 1/2/3/4); tier 1 sadece "attended" olmayı gerektirir.
- **Rekabet yok** — leaderboard opsiyonel (cosmetic only), reward bağımsız.

## 4. Dokümanlar
- [world-boss.md](world-boss.md) — Global world boss + attendance contract + phase encounter.
- [invasion.md](invasion.md) — Region invasion + multi-wave + territory defense scoreboard.
- [regional-event.md](regional-event.md) — Periodic region event (harvest festival, wildlife migration).
- [seasonal-event.md](seasonal-event.md) — Uzun süreli seasonal (holiday, anniversary) + event token currency.
- [dynamic-event.md](dynamic-event.md) — Emergent chain event (crime → guard call → rescue).

## 5. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md) — Contract & reveal
- [systems/reward-engine.md](../systems/reward-engine.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [economy/currency.md](../economy/currency.md) — `event_<code>` token
- [systems/notification.md](../systems/notification.md) — pre-event heads-up
- [guild/boss.md](../guild/boss.md), [guild/raid.md](../guild/raid.md) — sınırdaş konseptler

## 6. Metrikler (v1)
- `world_event_scheduled_total{event_kind}`
- `world_event_active_gauge{event_kind}`
- `world_event_participant_total{event_kind}`
- `world_event_reward_grant_total{event_kind, tier}`
- `world_event_afk_share_ratio` — AFK vs online contribution ratio

## 7. Bekleyen Kararlar
- Cross-shard world boss (Phase 03).
- Dynamic event chain karmaşıklığı v1 kapsam sınırı.
- Seasonal event pack yayın modeli (client hot bundle vs full patch).
- Leaderboard cosmetic ödül ne (title, tabard, mount skin).
