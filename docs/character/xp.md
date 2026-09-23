# XP (Experience)

## 1. Amaç
Karakterin ana ilerleme sayacı; kaynakları, çarpanları, debt (borç) ve grant akışının server-authoritative modeli.

## 2. XP Türleri
Ayrı ledger — birbirine dönüşmez (bkz. [game-design/progression.md](../game-design/progression.md)):
- **Character XP** — level ilerlemesi (bu doküman).
- **Profession XP** — meslek başına ayrı.
- **Mastery XP** — track başına ayrı.

Bu doküman **character XP**'e odaklanır; diğerleri paralel motorla çalışır.

## 3. XP Kaynakları

| Kaynak | Formül (baz) | Not |
|---|---|---|
| Mob kill | `mob.xp_reward * class_mod * rested_mod * event_mod` | AFK hunting'de tick başına bulk |
| Quest turn-in | `quest.rewards.xp` | Story quest yüksek; daily orta |
| Adventure completion | `adventure_template.reward_pool_xp` | Stage bazlı, kısmi başarıda partial |
| Dungeon completion | `boss_encounter.xp * clear_multiplier` | Weekly first clear bonusu |
| World boss contribution | `contribution_share * boss.xp_pool` | Katılım orantısında |
| Event participation | event_template'e göre | Kozmetik odaklı; XP az |
| Discovery (POI first visit) | 100..500 (region tier) | Phase 02 opsiyonel |
| Crafting (grandmaster recipe) | `recipe.xp_reward` küçük character XP + büyük profession XP | XP paylaşımı |
| Gathering (rare node) | Küçük character XP | Rare drop ile beraber |

## 4. Çarpanlar (Multipliers)

Uygulama sırası (soldan sağa çarpımsal):
```
base * (1 + rested_bonus) * class_mod * spec_mod * race_mod * item_mod * event_mod * party_mod * ascension_dampen
```

- **Rested bonus** (Phase 02): şehir inn'de bekleme → +%50 XP (limited buffer).
- **Class/Spec mod**: bazı specler `+5% xp gain` passive.
- **Race mod**: örn. Human `+5% XP` (bkz. [races.md](../game-design/races.md)).
- **Item mod**: XP jewelry / consumable tonic.
- **Event mod**: server-wide event window (+%20 XP haftası).
- **Party mod**: party içindeyken +%5..%10.
- **Ascension dampen**: `ASCENSION_MULT` ile threshold artışı XP curve'de.

Cap: toplam çarpan `< 3.0` (aşırı stack koruma).

## 5. Rested / Fatigue Modeli (Phase 02)
- Şehir inn'de logout → biriktirdikçe "rested XP buffer" büyür (max: current level required XP × 1).
- Hunting/adventure yaparken buffer aynı miktar tüketilerek 2× XP gain.
- Buffer boş → normal rate.

v1'de yok; buffer sıfır, ama design'da alan var (`character.character.rested_xp_buffer` sonradan eklenir).

## 6. XP Grant Sözleşmesi

Domain servis:
```go
package character

type XPSource struct {
    Kind         string    // "mob", "quest", "adventure", "boss", "event", "discovery", "craft"
    RefID        uuid.UUID // source event id
    BaseAmount   int64     // pre-multiplier
    Multipliers  Multipliers
    IdempotencyKey string
}

func (s *Service) GrantXP(ctx context.Context, charID uuid.UUID, src XPSource) (Grant, error)

type Grant struct {
    Applied   int64 // final XP after multipliers - debt payoff
    DebtPaid  int64
    LevelUps  []int   // new levels achieved
    NewXP     int64
    NewLevel  int32
}
```

Idempotency:
- Reward System (bkz. [systems/reward-system.md](../systems/reward-system.md)) `idempotency_key` ile grant tekrarını engeller.
- Aynı `(source_kind, source_ref_id, character_id)` iki kez grant edemez.

## 7. XP Debt (Death Penalty)

Bkz. [combat.md §Death](../game-design/combat.md) ve [death §XP debt](health.md).

Kurallar:
- **Level 20 altı**: debt yok (öğrenme koruma).
- **Level 20+**: ölüm başına `debt = current_level_xp_required * 0.05` (%5 of level).
- Debt `character.character.xp_debt` alanında biriker (schema ekleme; migration).
- XP grant önce debt öder:
  ```
  effective_grant = max(0, delta - remaining_debt)
  remaining_debt = max(0, remaining_debt - delta)
  character.xp += effective_grant
  ```
- Debt cap: `current_level_xp_required` (bir levellik borçtan fazla birikmez).
- Debt hiçbir zaman level down yaptırmaz.
- Debt reset koşulu: sadece ödendikçe.

Schema ekleme:
```sql
ALTER TABLE character.character ADD COLUMN xp_debt BIGINT NOT NULL DEFAULT 0;
```

## 8. XP Ledger

Kalıcı log:
```sql
CREATE TABLE character.xp_ledger (
  id BIGSERIAL PRIMARY KEY,
  character_id UUID NOT NULL,
  delta BIGINT NOT NULL,
  applied_delta BIGINT NOT NULL,  -- debt sonrası fiili
  debt_paid BIGINT NOT NULL DEFAULT 0,
  source_kind VARCHAR(24) NOT NULL,
  source_ref_id UUID NULL,
  multipliers_json JSONB NOT NULL DEFAULT '{}',
  idempotency_key TEXT NOT NULL UNIQUE,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Not: Currency ledger'a benzer append-only pattern.

## 9. Grant Akışı (Reveal Örneği)

AFK hunting reveal:
```
Reveal yield summary:
  mob kills = 620
  xp per mob = 40 (base)
  total base XP = 24800
  multipliers:
    race 1.05, event 1.20, party 1.0 → net 1.26
  final XP = 31248

GrantXP:
  debt = 500
  applied = 31248 - 500 = 30748
  new level jumps: 2
  new xp = ...
  ledger insert (idempotency_key = "reveal:<contract_id>:<reveal_id>")
  publish 'character.leveled_up'
```

## 10. XP Boost Item (Phase 02)
- `content.item_template.category='consumable', subcategory='xp_boost'`.
- Effect: 30 dk +%20 XP.
- Item consume → apply buff (character.status_effect).
- Stack rules: aynı buff refresh; farklı buff aynı anda ancak toplam cap.

## 11. Guest / Trial XP (Phase 03)
Yeni oyuncular için ilk 7 gün +%25 XP boost. Retention aracı.

## 12. XP Anti-Cheat
- Reveal yield delta anomaly detection (bkz. [systems/afk-engine.md](../systems/afk-engine.md)).
- Grant per hour cap (per source_kind).
- Ledger reconcile job: `SUM(applied_delta) per character` — expected level XP progression karşılaştırma.

## 13. Cross-Platform
XP grant server-side; client sadece display. Multi-device active: hangi client bağlı olursa olsun grant aynı character state'ini etkiler.

## 14. Fail Modları
- Idempotency conflict → no-op (safe replay).
- Debt overflow → cap uygulanır; log warn.
- Multiplier calculation NaN/Inf → sanitize → base value.

## 15. Metrics
- `character_xp_granted_total{source}` counter.
- `character_xp_debt_incurred_total`.
- `character_xp_debt_paid_total`.
- `character_xp_per_hour` gauge (aggregated).

## 16. Test Stratejisi
- Unit: multiplier stacking table.
- Debt: grant amount vs debt scenarios.
- Idempotency: aynı source çift grant → tek insert.
- Level jump: büyük XP grant → doğru level + kalan XP.

## 17. Audit
- `audit.event_log character.xp_granted` (severity=debug; high volume).
- `character.xp_debt_incurred` (info).

## 18. İlişkiler
- [Level](level.md), [Attributes](attributes.md), [Health](health.md), [systems/reward-system.md](../systems/reward-system.md), [systems/afk-engine.md](../systems/afk-engine.md), [game-design/progression.md](../game-design/progression.md), [database/character.md](../database/character.md)

## 19. Bekleyen Kararlar
- Rested XP buffer resmi Phase 02 startı.
- Ascension XP dampener değeri (%5 per ascension) simülasyon sonrası.
- XP debt formülü % oranı final (5% önerisi).
- XP curve BASE/FACTOR final.
