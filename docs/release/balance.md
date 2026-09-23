# Balance Pass

## 1. Amaç
Wave data ile ekonomi + progression + PvE difficulty + class win rate ayarları. Data-driven; content-driven; kod değişikliği minimum.

## 2. Ana Balance Ekseni

| Alan | Metrik | Target |
|---|---|---|
| Gold ekonomi | Sink/source ratio | 0.85-0.95 (haftalık) |
| Item drop rate | Legendary % | 0.5-1% aim (per-hour play) |
| XP curve | Level-up frequency | Player playtime dağılımı ile eşleş |
| Class win rate (PvE) | Damage / clear time | Class arası ±%5 farksızlık |
| Craft profit | Recipe profitability | Positive expected value (mat cost < sell price) |
| AFK yield | Rate per hour | Design assumption ± %10 |
| Quest completion rate | % players complete | > %70 main; > %40 side |

## 3. Data Kaynağı

- **Server telemetry**: ledger, item drops, XP grants, class win record.
- **Grafana dashboard**: balance-specific view.
- **Weekly report**: automated CSV + narrative.

## 4. Değişim Yöntemi

Balance değişikliği **content-driven** (kod değil):
- Item stat: `content.item_template` update → content pack redeploy.
- Drop rate: `content.drop_table` update.
- XP curve: `content.xp_curve` update.
- Recipe cost: `content.recipe` update.
- Class skill damage: `content.skill_definition` update.
- Fee: `config/economy/fee.yml` hot reload.

Rare cases kod değişikliği: formula değişimi (Phase 02 review).

## 5. Balance Cycle (Wave Arası)

Timeline:
- Wave X kapanış (Cuma 20:00).
- Cumartesi: data toplama + rapor.
- Pazar: game designer analiz + değişiklik önerisi.
- Pazartesi: team review + karar.
- Salı-Çarşamba: content pack update + staging test.
- Perşembe: content deploy prod.
- Cuma: Wave X+1 açılışı.

Toplam 1 hafta balance cycle.

## 6. A/B Testing (Phase 02)

- Feature flag ile 2 balance variant.
- Karşılaştırmalı retention + completion data.
- Winner variant permanent.

V1'de A/B **yok**; global change.

## 7. Class Balance

Warrior/Mage/Priest/Hunter/Rogue win rate:
- Solo dungeon clear time per class.
- PvE damage per second.
- Survivability (death rate).

Ayarlama:
- Underperform → +stat multiplier veya skill damage.
- Overperform → indirim yerine buff others (community-friendly).

## 8. Ekonomi Balance

- Sink/source ratio dashboard izle.
- Sink < 0.85: inflation risk → yeni sink (fee up, vendor sell multiplier düşür).
- Sink > 1.0: deflation risk → drop rate up veya vendor buy price down.

Aşırı ani değişim → player rage; küçük iterative adjustment (±5-10% per pass).

## 9. Recipe Profitability

Formula:
```
profit = expected_sell_price - Σ(material_cost × material_market_median)
```

Weekly compute; negative → recipe rework veya vendor buy fiyat artışı.

## 10. AFK Yield

Reveal telemetry — actual vs expected yield.
Deviation > %10:
- Formula bug (rare).
- Rare drop rate mis-set.
- Bonus stacking issue.

Ayarlama content pack.

## 11. Quest Completion Rate

Main quest chain step-by-step drop-off analiz:
- Step 3 → %20 drop-off → too hard / unclear objective.
- Fix: hint UI, difficulty tune, reward increase.

## 12. Documentation

Her balance pass:
- Rapor doc (data + rationale + değişiklik listesi).
- Patch notes (player-facing).
- Content pack version bump + git commit.

## 13. Metrikler
- `balance_pass_total{wave, category}`
- `content_pack_version_gauge`
- `class_win_rate_pct{class}`
- `economy_sink_source_ratio`
- `recipe_profit_gauge{recipe}`
- `quest_completion_rate_pct{quest}`

## 14. Cross-Reference
- [closed-beta.md](closed-beta.md), [bug-fix.md](bug-fix.md)
- [progression/](../progression/), [economy/](../economy/)
- [systems/reward-engine.md](../systems/reward-engine.md)

## 15. Bekleyen Kararlar
- A/B test framework Phase 02.
- Public balance changelog Phase 02.
- Player feedback loop (patch note comment section) Phase 02.
- Sim testing (AI player) Phase 03.
