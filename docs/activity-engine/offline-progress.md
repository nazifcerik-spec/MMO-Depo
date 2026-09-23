# Offline Progress — Analitik Hesap

## ⚠️ ALTIN KURAL: NEVER SIMULATE EVERY OFFLINE TICK.

Bu doküman aktivite yield'inin **analitik** (analytical, closed-form) hesabını tanımlar. Per-tick loop kesin **YASAK**.

## 1. Neden?

- **Ölçek**: 100k karakter × 8h AFK = 100k × 28800 tick = 2.88 milyar iterasyon.
- **Latency**: login p50 < 300ms hedefi tick simülasyonuyla imkansız.
- **Determinism**: analitik + seed → tekrarlanabilir.
- **Test edilebilirlik**: pure function; golden replay.

## 2. Ana Pipeline

```
Reveal(contract, char):
  1. Load contract + snapshot.
  2. elapsed_raw = now - contract.last_reveal_at
  3. elapsed_effective = softCap(elapsed_raw)
  4. rate = Activity.ComputeRate(snapshot, params)
  5. threshold = Activity.DetectThreshold(contract, snapshot, elapsed_effective)
  6. if threshold != nil:
       elapsed_effective = threshold.HitAtMs
  7. yield = Activity.ComputeYield(contract, snapshot, elapsed_effective)
  8. rewards + xp dispatch
  9. contract.last_reveal_at = now (adjusted by elapsed_effective)
  10. contract.status = updated (paused/completed/active)
```

**Tüm hesap O(1)** — elapsed'e bağlı iteration yok.

## 3. Soft Cap

`softCap(elapsedMs)`:
```go
func softCap(elapsedMs int64) int64 {
    const HOUR = 3600 * 1000
    fullCap := 24 * HOUR
    mediumCap := 48 * HOUR
    lowCap := 168 * HOUR

    if elapsedMs <= fullCap {
        return elapsedMs  // 24h tam
    }
    if elapsedMs <= mediumCap {
        return fullCap + int64(float64(elapsedMs-fullCap) * 0.6)  // 24-48h %60
    }
    if elapsedMs <= lowCap {
        return fullCap + int64(float64(mediumCap-fullCap)*0.6) + int64(float64(elapsedMs-mediumCap)*0.3)
    }
    // 168h+ (7 gün) → 0 ek
    return fullCap + int64(float64(mediumCap-fullCap)*0.6) + int64(float64(lowCap-mediumCap)*0.3)
}
```

Bu formul curve DD-014'te tanımlı.

## 4. Rate Compute (Örnek Mining)

```go
rate = base_yield_per_hour(spot, tier)
     × (1 + Σ gather_yield_pct_modifiers / 100)
     × (1 + mastery_bonus_fraction)
     × racial_multiplier
```

Rate scalar; elapsed'den bağımsız.

## 5. Yield Compute — Analytical Expected Value

Basit case:
```
expected_yield = rate × elapsed_hours
```

Deterministic path (raw expected, no variance):
```
yield = expected_yield rounded down
```

## 6. Variance Sampling — Aggregate Distribution

RNG variance için tek sample, distribution'dan:

### 6.1 Poisson (Gathering, Kill Count)
Kaynak: gathering yield, kill count — Poisson dağılımı λ = rate × elapsed.

```go
func samplePoisson(ratePerHour float64, elapsedHours float64, seed Seed) int32 {
    lambda := ratePerHour * elapsedHours
    if lambda < 20 {
        // Knuth's algorithm — O(lambda) küçük lambda için OK
        return knuthPoisson(lambda, seed)
    }
    // Büyük lambda → normal approximation
    return normalApproxPoisson(lambda, seed)
}
```

**Not**: `knuthPoisson` küçük lambda için iterasyon yapar ancak lambda < 20 zaten küçük. Büyük lambda → normal approximation O(1).

Cap: lambda > 1000 → sadece deterministic mean (variance ihmal).

### 6.2 Binomial (Craft Success, Drop Chance)
```go
func sampleBinomial(n int32, p float64, seed Seed) int32 {
    if n < 20 {
        // BinomTable veya BTPE
        return smallBinomial(n, p, seed)
    }
    // Normal approximation
    mean := float64(n) * p
    variance := mean * (1 - p)
    return int32(normalSample(mean, sqrt(variance), seed))
}
```

Cap: n > 100000 → deterministic mean.

### 6.3 Rare Drop (Multi-Roll Poisson)
```go
rareDropCount := samplePoisson(rareDropRatePerHour, elapsedHours, seed)
// Sonra rareDropCount kadar loot table roll — bounded küçük N.
```

## 7. Threshold Detection

### 7.1 Inventory Full
Linear inverse:
```
inv_free_slots = char.InventoryCap - char.InventoryUsed
slots_per_hour = rate.YieldPerHour × avg_slots_per_yield
ms_until_full = inv_free_slots / slots_per_hour × 3600000
```

Threshold set: `{Kind: inventory_full, HitAtMs: ms_until_full}`.

### 7.2 Stamina Zero (v1 kapalı; Phase 02 aktivasyon)
Linear:
```
stamina_drain_per_hour = rate.ResourceDrainPerHour
ms_until_stamina_zero = char.CurrentStamina / stamina_drain_per_hour × 3600000
```

### 7.3 Bait/Consumable Exhaustion
Linear:
```
consume_per_hour = rate.ConsumablePerHour
ms_until_out = char.ConsumableCount / consume_per_hour × 3600000
```

### 7.4 Death Cascade (Hunting)
Poisson-based:
```
expected_deaths = death_rate_per_hour × elapsed_hours
if expected_deaths > max_deaths_allowed:
    // Determine when first "cascade" happens
    // Simplified: threshold at expected_time_of_first_death
    ms_until_first_death = 1 / death_rate_per_hour × 3600000
    if params.stop_on_death_cascade and death_rate > safety_threshold:
        threshold = {Kind: death_cascade, HitAtMs: ms_until_first_death}
```

### 7.5 Weekly Lockout (Dungeon)
Time-based; server reset schedule known.

## 8. Multi-Threshold Detection

Multiple threshold aday → smallest `HitAtMs` kazanır:
```go
thresholds := []Threshold{
    detectInventoryFull(),
    detectStaminaZero(),
    detectConsumableOut(),
    detectDeathCascade(),
    detectLockout(),
}
// Filter nil
// Min HitAtMs
```

Return smallest; `elapsed_effective = threshold.HitAtMs`.

## 9. Content Version Consistency

Reveal snapshot'ın `content_version`ı ile mevcut content_version farklıysa:
- Contract snapshot content'i kullanır (kararlı hesap).
- Post-reveal: uyarı log; opsiyonel content re-snapshot new contract için.

## 10. Idempotency

Reveal idempotent:
```
idempotency_key = hash(contract_id, contract.last_reveal_at, now_bucketed)
```

Aynı reveal 2 kez çağrılırsa (client retry, network glitch):
- `activity.reveal_log` unique key → duplicate no-op.
- Yield dispatch atomic within reveal transaction.

## 11. Reveal Log

Kaynak: [database/activity.md §5](../database/activity.md).

```
activity.reveal_log:
  id, character_id, contract_id, reveal_at,
  elapsed_ms, capped_ms, yield_summary_json,
  events_ref (optional detail), idempotency_key
```

Full audit trail.

## 12. RNG Seed Derivation

Contract seed → sub-streams:
```
contract_seed = hash(character_seed, activity_type, started_at, salt)
yield_seed    = hash(contract_seed, "yield")
loot_seed     = hash(contract_seed, "loot")
rare_seed     = hash(contract_seed, "rare")
death_seed    = hash(contract_seed, "death")
```

Her sub-stream deterministic + izole.

Detay: [systems/random-system.md](../systems/random-system.md).

## 13. Edge Cases

### 13.1 Elapsed = 0
No-op reveal; audit info.

### 13.2 Very Long Elapsed (weeks)
Soft cap uygular; hard cap 30 gün → contract archived (worker); reveal safe fallback.

### 13.3 Snapshot Corrupted
Fresh snapshot al; reveal partial + log warn.

### 13.4 Content Deprecated Mid-Contract
Snapshot content kullanılır; UI'da "content changed" hint.

### 13.5 Character Died Multiple Times
Death count from Poisson; cascade prevention flag pause.

### 13.6 Inventory Full at Start
elapsed_effective = 0; contract paused immediately.

## 14. Performance Targeting

- Reveal per contract: **< 50ms** (p50).
- Reveal batch (100 characters): < 5s.
- Reveal throughput per node: ≥ 5000 char/sn (bkz. GDD §18).

Bottleneck: RNG (Poisson small lambda) + DB write (reveal_log + ledger).

Optimization opsiyonları (Phase 02):
- Snapshot cache Redis short-lived.
- Reveal batch worker.
- Ledger write coalescing.

## 15. Test Stratejisi

- **Golden replay**: aynı input → aynı yield deterministic (hash compare).
- **Poisson distribution**: 10k trials → mean within ±%2 of lambda.
- **Threshold accuracy**: computed hit_at vs simulated → match ±100ms.
- **Soft cap curve**: elapsed 12h/24h/36h/72h/168h → expected multiplier.
- **Idempotent retry**: 2 reveal calls → 1 dispatch.
- **Overflow**: expected_yield >> inventory_cap → threshold detect.
- **Zero rate**: rate=0 → yield=0, no crash.

## 16. Anti-Patterns (Never Do)

- ❌ `for t := 0; t < elapsed_seconds; t++ { simulate() }`
- ❌ Per-second RNG roll toplama.
- ❌ Per-mob simulate + damage exchange.
- ❌ Sequential loot table roll her drop için (bounded aggregate OK).
- ❌ Wall clock elapsed compute (server monotonic only).

**Do**:
- ✅ Rate × elapsed = expected.
- ✅ Poisson/Binomial/Normal sample distribution'dan.
- ✅ Threshold: linear inverse solve.
- ✅ Bounded loop (küçük N drops iteration OK, rare count).

## 17. Metrics

- `reveal_compute_duration_seconds` histogram.
- `reveal_batch_throughput` gauge.
- `reveal_threshold_hit_total{kind}`.
- `reveal_soft_cap_applied_total`.
- `reveal_idempotent_no_op_total`.
- `reveal_variance_sample_calls_total`.

## 18. Config

```yaml
# config.offline-progress.yaml
soft_cap:
  full_hours: 24
  medium_multiplier: 0.6
  low_multiplier: 0.3
  medium_range_hours: [24, 48]
  low_range_hours: [48, 168]
  archive_after_hours: 720   # 30 gün

variance:
  poisson_normal_threshold_lambda: 20
  binomial_normal_threshold_n: 20
  deterministic_only_lambda_cap: 1000
  deterministic_only_binomial_cap: 100000

reveal_performance:
  target_p50_ms: 50
  target_p99_ms: 500
  batch_worker_size: 100
```

## 19. Cross-Reference

- [Activity](activity.md), [Lifecycle](lifecycle.md), [Rewards](rewards.md), [XP](xp.md)
- Aktivite dokümanları (yield formulası her aktivitede)
- [systems/afk-engine.md](../systems/afk-engine.md), [systems/random-system.md](../systems/random-system.md)
- [database/activity.md §5](../database/activity.md)
- DD-014 (soft cap), DD-067 (deterministic RNG)

## 20. Bekleyen Kararlar

- Poisson large lambda cap threshold final (v1: 1000).
- Snapshot cache Redis short-lived (Phase 02).
- Reveal batch worker enable Phase 02.
- Content-version drift compensation policy.
