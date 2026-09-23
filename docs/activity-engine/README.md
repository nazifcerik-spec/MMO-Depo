# Activity / AFK Engine (Faz 12) — Oyunun Kalbi

Karakterin tüm zaman-tabanlı işlerinin (mining, fishing, crafting, hunting, travel, combat, dungeon) tek soyutlama altında yürütüldüğü motor. **Online veya offline** — engine aynı.

Foundation dokümanlar: [systems/afk-engine.md](../systems/afk-engine.md), [systems/activity-engine.md](../systems/activity-engine.md). Bu klasör Faz 12 implementation-hazır derinlik.

---

## 🚫 Kritik İlke: NEVER SIMULATE EVERY OFFLINE TICK

**Bu motor asla per-tick simülasyon yapmaz.**

Yani:
- ❌ `for t := 0; t < elapsed_seconds; t++ { simulateOneTick() }` — **YAP-MA**.
- ✅ **Analitik hesap**: `expected_yield = rate × elapsed × probability` — O(1).
- ✅ **Sampled variance**: RNG **aggregate distribution**'dan sample; per-tick roll yok.
- ✅ **Threshold events** (envanter dolma, stamina bitme): **binary search** ile "ne zaman doldu" hesabı; sonra o ana kadar simulate.

### Neden?
- **Ölçek**: 100k karakter × 8 saatlik AFK = 100k × 28800 tick simülasyonu = 2.88 milyar iterasyon. Reveal yığılır.
- **Determinism**: analitik hesap deterministic + tekrar edilebilir (audit + replay).
- **Latency**: login p50 < 300ms hedefi (bkz. GDD §18) tick simülasyonuyla imkansız.
- **Performance budget**: reveal ≥ 5000 char/sn (bkz. system-architecture.md §4) sadece O(1) hesapla mümkün.

### Nasıl?
1. **Rate compute** — karakter derived stat + spot template + modifier stack.
2. **Elapsed compute** — `now - last_reveal_at`, soft cap uygula.
3. **Analytical yield** — `expected = rate × elapsed_effective`.
4. **Stochastic variance** — poisson/binomial toplam sample (tek RNG draw yerine).
5. **Threshold detection** — inventory full, stamina drain, mob dry → **çözüm formülü** (linear inverse).
6. **Ledger + reveal_log** — sonuç bir kez yazılır.

Detay: [offline-progress.md](offline-progress.md).

---

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| **Kritik ilke + genel bakış** (bu doküman) | [README.md](README.md) |
| [Activity](activity.md) | Ortak activity abstraction, tipler, tree |
| [Mining](mining.md) | Gathering: cevher, taş node |
| [Fishing](fishing.md) | Gathering: kıyı/su node |
| [Crafting](crafting.md) | Recipe-based üretim |
| [Hunting](hunting.md) | AFK combat: mob spawn + rotation |
| [Travel](travel.md) | Konum değiştirme aktivitesi |
| [Combat](combat.md) | Aktif combat aktivitesi (rotation) |
| [Dungeon](dungeon.md) | Idle-run dungeon |
| [Lifecycle](lifecycle.md) | Start / Pause / Resume / Cancel / Complete |
| [Offline Progress](offline-progress.md) | ★ Analitik hesap (per-tick simülasyon YOK) |
| [Rewards](rewards.md) | Reward payload compute (item + gold + rep) |
| [XP](xp.md) | XP compute (character + profession + mastery) |

## 2. Aktivite Ağacı

```
Activity
 ├── Gathering
 │   ├── Mining
 │   ├── Fishing
 │   ├── Woodcutting  (Phase 02 profession)
 │   ├── Foraging     (Phase 02)
 │   └── Skinning     (hunting + skinning yan aktivite)
 ├── Crafting
 │   ├── Blacksmithing / Leatherworking / ... (recipe-driven)
 │   └── Cooking / Alchemy / Enchanting / Jewelcrafting
 ├── Hunting          (AFK combat vs mob spawn)
 ├── Travel           (konum değişimi)
 ├── Combat           (manuel + rotation AFK)
 ├── Dungeon
 │   ├── Manuel dungeon run
 │   └── Idle-run dungeon
 └── Adventure        (story-based mixed check)
```

Her aktivite ortak `Activity` interface'ini implement eder (bkz. [activity.md](activity.md)).

## 3. Motor Interface Genel

```go
package activityengine

type Engine interface {
    // Lifecycle
    Start(ctx context.Context, req StartRequest) (Contract, error)
    Pause(ctx context.Context, contractID uuid.UUID, reason PauseReason) error
    Resume(ctx context.Context, contractID uuid.UUID) error
    Cancel(ctx context.Context, contractID uuid.UUID, mode CancelMode) (Settlement, error)
    Complete(ctx context.Context, contractID uuid.UUID) (Settlement, error)

    // Reveal (offline progress + rewards)
    Reveal(ctx context.Context, characterID uuid.UUID) (RevealReport, error)

    // Query
    CurrentContract(ctx context.Context, characterID uuid.UUID, channel Channel) (*Contract, error)
    ExpectedYield(ctx context.Context, contract Contract, projectionMs int64) (YieldProjection, error)
}
```

Detay implementation: [lifecycle.md](lifecycle.md), [offline-progress.md](offline-progress.md).

## 4. Contract Model

Kaynak: [database/activity.md §1](../database/activity.md).

```
activity.contract
  id, character_id, activity_type,
  params_json, status,
  started_at, last_reveal_at, paused_at, pause_reason,
  seed, rng_version_id, snapshot_ref,
  location_id, channel (primary/pet/mount_travel),
  soft_cap_override_ms,
  correlation_id_seed,
  version
```

Karakter aynı anda **max 1 primary + 1 pet + 1 travel** contract (partial unique constraint bkz. [database/activity.md §1](../database/activity.md)).

## 5. Zamanlama Katmanı

```
now — server monotonic clock (UTC ms; NTP-safe)
started_at — contract create time
last_reveal_at — son hesap noktası
elapsed_raw = now - last_reveal_at
elapsed_effective = softCap(elapsed_raw)   # bkz. offline-progress.md §Soft Cap
```

Tüm hesap `elapsed_effective` üstünde; **tick sayımı YOK**.

## 6. Snapshot Discipline

`activity.snapshot` — contract oluştuğunda karakterin dondurulmuş görüntüsü:
- Stat cache
- Equipment
- Learned skills
- Race + class + spec + ascension node'ları
- Content version tag

Snapshot analitik hesabın **giriş parametreleri**. Karakter arada değişse bile reveal snapshot'ı kullanır (contract süresince kararlı hesap).

Aktif oyun (Reveal → yeniden Start) → yeni snapshot alınır.

## 7. Kullanılan Diğer Sistemler

| Bileşen | Rolü |
|---|---|
| [Skill+Stat Engine](../skill-engine/README.md) | Rate compute, combat cast (hunting) |
| [Item Engine](../item-engine/README.md) | Loot instance yaratımı |
| [Inventory](../inventory/README.md) | Yield destination, capacity check |
| [Reward System](../systems/reward-system.md) | Payload dispatch |
| [Loot System](../systems/loot-system.md) | Loot table roll |
| [Random System](../systems/random-system.md) | Deterministic seed stream |
| [Notification System](../systems/notification-system.md) | Reveal ready push, activity paused alert |
| [Event Bus](../backend/event-bus.md) | activity.contract_started/completed/paused |

## 8. Reveal (Login Trigger) — Yüksek Seviye

```
POST /v1/characters/{id}/reveal   (auto on login veya explicit)

Server:
1. Load character contracts (primary + pet + travel channels).
2. For each active/paused contract:
   a. elapsed = now - last_reveal_at
   b. elapsed_effective = softCap(elapsed)
   c. yield = computeAnalytical(contract, snapshot, elapsed_effective)
   d. threshold detection (inventory full, stamina, mob dry)
      - if hit: yield_valid_until = t_hit; recompute
   e. rewards + xp payload build
   f. write reveal_log (idempotent)
   g. dispatch rewards (ledger, inventory, xp ledger, mastery)
   h. contract.last_reveal_at = t_effective_end
   i. contract.status = paused if threshold hit else active
3. Return RevealReport (aggregate yield summary).
```

Detay: [offline-progress.md](offline-progress.md).

## 9. Design Non-Goals

- **Per-tick simulation** — asla.
- **Wall clock bağımlılık** — server monotonic time; client asla.
- **Non-deterministic RNG** — asla; seed stream.
- **Cross-contract dependency** — her contract bağımsız hesaplanır.
- **Real-time position** — v1 aktivite lokasyon anchor'lı; positional yok.

## 10. Ortak İlkeler

1. **Server-authoritative** — reveal + settlement always server.
2. **Idempotent reveal** — aynı `idempotency_key` → no-op.
3. **Deterministic** — snapshot + seed + elapsed → aynı yield.
4. **Analytical > iterative** — per-tick loop asla.
5. **Threshold-aware** — envanter/stamina limit'i "ne zaman ulaştı" hesabı.
6. **Config-driven rates** — yield/xp coefficient config, hot-tunable.

## 11. Config Kataloğu (Örnek)

```yaml
# config.activity-engine.yaml
soft_cap:
  full_hours: 24
  medium_hours_range: [24, 48]
  low_hours_range: [48, 168]
  zero_after_hours: 168
  medium_multiplier: 0.6
  low_multiplier: 0.3

rates:
  mining:
    base_yield_per_hour_by_tier:
      tier_1: 60
      tier_2: 90
      tier_3: 120
    xp_per_yield: 5
  fishing:
    base_yield_per_hour_by_tier:
      tier_1: 40
      tier_2: 65
  hunting:
    kills_per_hour_base: 60
    dps_multiplier: 1.0
```

Rate simülasyon Phase 01 balance pass tune eder.

## 12. Cross-Reference

- [systems/afk-engine.md](../systems/afk-engine.md), [systems/activity-engine.md](../systems/activity-engine.md)
- [database/activity.md](../database/activity.md)
- [skill-engine/README.md](../skill-engine/README.md)
- [item-engine/README.md](../item-engine/README.md)
- [inventory/README.md](../inventory/README.md)
- [character/xp.md](../character/xp.md), [character/level.md](../character/level.md)
- [game-design/hunting.md](../game-design/hunting.md), [game-design/professions.md](../game-design/professions.md)

## 13. Ortak Bekleyen Kararlar
- Reveal max parallel contracts per character (v1: 3 kanal).
- Threshold detection precision (approx vs binary search precision).
- Content patch mid-contract: contract archive vs re-simulate policy.
- Progression fairness monitoring (spike detection).
