# AFK Engine Tests

## 1. Amaç
AFK Engine'in analitik hesap (DD-084), contract snapshot (DD-088/095), threshold detection ve reveal tutarlılığını garantiye alma.

## 2. Kritik Invariant'lar

- **Never simulate every tick** (DD-084) — reveal O(1) analitik hesap (Poisson/Binomial sampling).
- **Snapshot immutable** — contract snapshot at start; runtime değişimler snapshot'ı etkilemez.
- **Idempotent reveal** — aynı reveal 2x → aynı response; ledger 2x yazmaz.
- **Threshold detection linear inverse** — threshold aşıldıysa exact tick hesap (iteration yok).
- **Time monotonic** — reveal sonucu server time'a bağlı; client time yok sayılır.

## 3. Unit Test Senaryoları

### 3.1 Contract Snapshot
- ✓ `NewContract` snapshot at `started_at`; character.gear_score değişse bile snapshot sabit.
- ✓ Snapshot içeriği zorunlu alan: gear_score, mastery_tier, race_bonus, class_bonus, active_buffs, energy_regen_rate.
- ✓ Snapshot JSON serialize/deserialize round-trip stable.

### 3.2 Analytical Yield Calc
- ✓ 1 saat mining, expected 40 ore ± 8 (Poisson mean=40): sampler mean 39.5-40.5.
- ✓ Rare drop chance 1% × 3600 tick → binomial N=3600, p=0.01 → expected 36 ± 6.
- ✓ Multi-resource: iron 80% + gem 15% + rare_gem 5% → multinomial samplers.
- ✓ Yield formula per-tick × ticks; bonus multiplier apply order (flat → mult_add → mult_prod → clamp).

### 3.3 Threshold Detection (Linear Inverse)
- ✓ Energy cap at 100; regen 0.5/s; contract 8s duration, start=50 EN → **cap reached at t=100s**; sonrası 0 EN cost activity.
- ✓ Inverse solve: `t_cap = (100-50)/0.5 = 100s`; assertion: worker doesn't iterate 8*3600 ticks.
- ✓ Multiple threshold cascade (energy 100 → aktivite dur; sonra regen 100 → resume).

### 3.4 Reveal Idempotency
- ✓ Reveal(contract_id) → response A.
- ✓ Reveal(contract_id) — 2. kez → aynı response A (cached); ledger entry sayısı sabit.
- ✓ Concurrent reveal 5 goroutine → hepsi aynı response; ledger race yok.

## 4. Property-Based Tests

### 4.1 Yield Distribution
- Property: 10.000 seed × 1 saat mining → mean(yield) within ±3σ of expected.
- Property: yield always ≥ 0.
- Property: sum(rare + common) ≤ total_tick.

### 4.2 Snapshot Determinism
- Property: aynı snapshot + aynı seed → aynı reveal (bit-exact).

## 5. Integration Tests

### 5.1 Contract Lifecycle
- Testcontainer Postgres + Redis.
- Create character + start contract → DB insert + Redis key.
- Simulate time skip (clock inject 8h) → reveal → check ledger + inventory + XP + mastery.
- Assert audit entries.

### 5.2 Concurrent AFK Slot
- Max concurrent AFK (DD-OPEN-3: 1 ana + 1 pet + 1 travel).
- Aynı karakter için 2. mining kontratı reddedilmeli → `E_CONCURRENT_SLOT_EXCEEDED`.

### 5.3 Cross-Region Regional Event Multiplier
- Regional event active period başında mining contract → yield +30%.
- Contract event window'unu aştıysa yield hesap event pencere kesim → doğru bölünme.

## 6. Fuzz Tests

- Contract params fuzzer: kind + duration + snapshot content → hiçbir input server panic yaratmaz.
- Time skip fuzzer: rastgele skip 1s-100y → threshold hesap doğru.

## 7. E2E Test (Smoke)

- Login → start mining 15dk → clock advance → reveal → inventory 30-50 ore + XP > 0.
- Login → start travel 30dk → clock advance → arrived at destination.

## 8. Regression Fixtures

- `testdata/afk/mining_1h_lvl15_snapshot.json` → deterministic reveal output snapshot.
- Snapshot değişmesi guard: yield hesap değişirse assertion fail + review.

## 9. Coverage Hedefi

- `systems/activity-engine`: ≥ 90%.
- `systems/reward-engine` (AFK path): ≥ 85%.
- Property tests: 10.000 iteration nightly; 1.000 CI.

## 10. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md), [systems/rng.md](../systems/rng.md)
