# Random System (RNG & Seed)

## 1. Amaç
Determinism, adalet, replayability ve audit için tek bir RNG stratejisi tanımlamak.

## 2. Oyuncuya Fayda
- Şeffaf, tekrar edilebilir sistemler.
- "Server hile yapıyor" hissini engelleyen açık şansların formalize edilmiş dokümantasyonu.

## 3. İlkeler
- **Deterministic**: `seed + counter` → değişmez sonuç.
- **Separable Streams**: her aktivite/roll için ayrı sub-stream (contamination yok).
- **Auditable**: seed + tick log ile replay mümkün.
- **Server-only**: hiçbir zaman client tarafında ödül belirleyen roll olmaz.

## 4. Seed Kaynağı
- **Character seed**: karakter oluşturulduğunda üretilir; uzun ömürlü.
- **Contract seed**: aktivite başlangıcında üretilir (`hash(character_seed, activity_type, started_at, salt)`).
- **Tick seed**: `hash(contract_seed, tick_index)`.
- **Event seed**: server-wide event için independent stream.

## 5. Algoritmalar
- Primary: **xoshiro256++** veya **ChaCha8-based CSPRNG**, versiyonlu (Phase 01 karar).
- Version bump: mevcut kontratlar için sabitlenmiş algoritma korunur; yeni contract yeni versionu kullanır.
- Distribution utilities: uniform, weighted, gaussian (rare use), logarithmic.

## 6. RNG API
```
rng(seed).uniform(a,b)
rng(seed).weighted_pick(list_of_(item, weight))
rng(seed).chance(pct)
rng(seed).roll_range(min,max)
```

## 7. Data
```
rng_version_registry
  version_id, algorithm_id, created_at
```
Her contract `rng_version_id` alanı taşır.

## 8. Failure Durumları
- Version mismatch: server refuses replay; fresh reveal.
- Seed re-use across contracts: cross-check + reject.

## 9. Edge Cases
- Aynı tick multiple roll: sub-stream `hash(tick_seed, purpose)`.
- Bad-luck protection: streak counter (bkz. Loot System).

## 10. Server Authoritative Gereksinimleri
- RNG all server; asla client.
- Seed leak riski: seed public expose edilmez, sadece hash.

## 11. Genişletilebilirlik
- Server-wide event RNG (world boss rare drop).
- Cross-region RNG independent.

## 12. Örnek
Contract seed = `sha256(char_seed + "mining" + "1725720000" + "salt-v1")`.
Tick 42 loot roll → sub-seed `sha256(contract_seed + 42 + "loot")` → uniform 0..1 → 0.0032 < 0.005 (rare drop). Item picked.

## 13. Balance Considerations
- Drop rates dokümante edilmeli (transparency).
- Streak protection açıkça belirtilmeli.
- Loot pity systems Phase 02'de opt-in olabilir.
