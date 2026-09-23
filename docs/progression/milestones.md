# Milestones

## 1. Amaç
Karakter/account bazında zaman veya hacim eşiklerine (100 saat oyun, 1M gold earned, 10k enemy kill vs) ulaşan retro reward + hikaye markerları. Achievement'ın kuzeni; "aktivite akümülatörleri".

## 2. Karakteristik

- **Passive** — oyuncu explicit tamamlamaz; ekonomi/aktivite/oturum sayaçları biriktikçe unlock.
- **Sayaç bazlı** — her milestone bir counter kind + threshold.
- **Bir kez unlock** — geçmiş sayaç retroactive; ancak reward 1 kez.
- **Reward paketi** — gold, event token, cosmetic, title, achievement.
- **Chain milestone** — 100h → 500h → 1000h zinciri.

## 3. Milestone Counter Türleri (v1)

| Kind | Kaynak | Örnek Threshold |
|---|---|---|
| `total_playtime_hours` | session tick | 100, 500, 1000, 3000 |
| `total_gold_earned_copper` | ledger source events | 1M, 10M, 100M, 1B copper |
| `total_gold_spent_copper` | ledger sink events | Aynı |
| `total_enemies_killed` | event.combat.enemy_kill | 1k, 10k, 100k |
| `total_ore_mined` | event.gather.mining | 1k, 10k, 100k |
| `total_items_crafted` | event.craft.complete | 100, 1k, 10k |
| `total_distance_traveled_m` | travel complete | 100k, 1M, 10M |
| `total_dungeon_completed` | event.dungeon.complete | 10, 100, 500 |
| `total_quests_completed` | event.quest.complete | 100, 500, 2000 |
| `total_login_days` | daily login | 30, 100, 365 |
| `total_afk_ticks` | activity tick complete | 100k, 1M, 10M |

## 4. Model

```sql
CREATE TABLE content.milestone_template (
  id UUID PRIMARY KEY,
  code VARCHAR(64) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  counter_kind VARCHAR(48) NOT NULL,
  threshold BIGINT NOT NULL,
  scope VARCHAR(16) NOT NULL DEFAULT 'character',   -- 'character'|'account'
  reward JSONB NOT NULL,
  chain_next_code VARCHAR(64) NULL,
  version INT DEFAULT 1
);

CREATE TABLE progression.milestone_counter (
  scope_kind VARCHAR(16) NOT NULL,     -- 'character'|'account'
  scope_ref UUID NOT NULL,
  counter_kind VARCHAR(48) NOT NULL,
  value BIGINT NOT NULL DEFAULT 0,
  last_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (scope_kind, scope_ref, counter_kind)
);

CREATE TABLE progression.milestone_reached (
  scope_kind VARCHAR(16) NOT NULL,
  scope_ref UUID NOT NULL,
  milestone_id UUID NOT NULL,
  reached_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (scope_kind, scope_ref, milestone_id)
);
```

## 5. Counter Update

Event bus consumer (`milestone_counter_apply`):
- Event → counter map (`enemy_kill → total_enemies_killed +1`).
- `milestone_counter` upsert.
- Threshold check: `SELECT * FROM milestone_template WHERE counter_kind=? AND threshold <= new_value AND NOT EXISTS (SELECT 1 FROM milestone_reached WHERE ...)`.
- Reach unlock → reward apply + WS "milestone_reached".

Aggregate high-volume (kill count) → batched (500ms) upsert.

## 6. Reward Delivery

- Auto-claim v1 default (opt Phase 02 manual claim).
- Reward JSONB parse:
  - Gold → ledger.
  - Event token → currency.
  - Title → title unlock.
  - Cosmetic → collection unlock.
  - Achievement → trigger.

## 7. Chain

```yaml
- code: "playtime_100h"
  counter_kind: "total_playtime_hours"
  threshold: 100
  reward: { gold: 1000c, title: "the_dedicated" }
  chain_next_code: "playtime_500h"

- code: "playtime_500h"
  counter_kind: "total_playtime_hours"
  threshold: 500
  reward: { gold: 5000c, cosmetic: "aura_veteran" }
  chain_next_code: "playtime_1000h"
```

Chain UI'da progress bar next threshold gösterir.

## 8. Retroactive

Yeni milestone content release'te retroactive `SELECT value FROM milestone_counter WHERE counter_kind=?` scan; threshold aşıldıysa retroactive unlock.

## 9. Endpoint

**GET `/v1/character/{id}/milestones?category=`**
**GET `/v1/account/{id}/milestones`**
**GET `/v1/character/{id}/milestones/counters`** — tüm sayaçlar.

## 10. Metrikler
- `milestone_reached_total{code}`
- `milestone_counter_update_total{kind}`
- `milestone_reward_grant_total{reward_kind}`

## 11. Cross-Reference
- [achievements.md](achievements.md), [titles.md](titles.md), [collections.md](collections.md)
- [systems/activity-engine.md](../systems/activity-engine.md), [economy/transaction.md](../economy/transaction.md)

## 12. Bekleyen Kararlar
- Manual claim vs auto-claim (v1: auto).
- Counter retention (Phase 03'te 5+ yıl account için parquet arşiv).
- Achievement vs milestone overlap (bazı achievement aslında milestone) — v1 karışık; Phase 02 taxonomy pass.
- Cross-character counter (account scope) hangi kind'lar (v1: playtime + login days account; diğerleri character).
