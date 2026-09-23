# AFK Engine (Offline İlerleme Motoru)

## 1. Amaç
Oyuncu offline iken karakterin aktivite kontratını sürdürmek ve login anında **deterministik**, **auditable** bir şekilde ödülleri hesaplamak. Oyunun kalbi.

## 2. Oyuncuya Fayda
- Offline geçen zamanın "kayıp" olmaması.
- Login anında somut ödül (reveal).
- Zaman biriktikçe adil ödül eğrisi.

## 3. Gameplay Loop
- Session Loop'ta reveal ilk aksiyondur.
- Meta Loop'ta 8 saat +/- offline pencerelerine dayanır (uyku, iş).

## 4. Temel Kavramlar
- **Activity Contract**: karakterin en son başlattığı aktivite tanımı; server persist.
- **Reveal**: login sırasında contract'ın son bilinen zamanı ile şimdiki zaman arasında ödüllerin hesaplanması.
- **Deterministic Simulation**: seed + timestamp + activity state → same input → same output.
- **Soft Cap**: belirli sürenin ötesinde verim eğrisi yumuşar; hiçbir zaman sıfır olmaz.
- **Snapshot**: karakter state'in versionlanmış anlık görüntüsü (envanter kapasitesi, ekipman, stat).

## 5. Kurallar
1. Kesin zaman kaynağı server epoch (UTC ms).
2. Contract oluşturulduğunda: `start_at`, `activity_type`, `params`, `seed`, `character_snapshot_ref`.
3. Reveal her istekte idempotent (aynı reveal iki kez çalışmaz; correlation_id ile).
4. **Soft cap**: default 24 saat verim tam; 24–48 saat %60 verim; 48+ saat %30 verim; 168 saat (7 gün) sonra 0.
   - Not: soft cap parametreleri profession/aktivite bazlı override edilebilir.
5. **Hard cap**: 30 gün ötesinde snapshot arşive taşınır; reveal safe fallback döner.
6. **Interruption**: envanter dolu, stamina bitti, insufficient tool durability, mob out of range → activity `paused`.
7. Envanter dolduğunda overflow mailbox'a düşer (Phase 01 karar), veya activity yield discard edilmez, sadece durur.
8. Aynı contract farklı istemciden re-open olduğunda: server otoritedir; en son login geçerli.

## 6. Reveal Hesaplama Akışı
```
1. Load contract, snapshot, character state
2. Compute elapsed = now - last_reveal_at
3. Apply soft cap curve
4. Simulate ticks:
   for each tick (or bulk chunk):
     - Roll RNG (seed + tick_index)
     - Compute yield (base × character mods × race/class/spec passives × luck)
     - Check stop conditions (inv full, stamina, mob dry, death)
     - Emit "activity_event" records (compressed)
5. Persist reveal_result (idempotent)
6. Emit character deltas (xp, gold, items, mastery, progress)
7. Return reveal payload to client
```

## 7. Data Model
```
activity_contract
  id, character_id, activity_type, params_json, seed
  started_at, last_reveal_at, status (active/paused/completed)
  snapshot_ref, correlation_id_seed

reveal_log
  id, character_id, contract_id, reveal_at
  elapsed_ms, capped_ms, yield_summary_json
  idempotency_key (unique)
```

## 8. I/O
- Input: `reveal(character_id, now)` (login trigger).
- Output: reveal payload (rewards, events, deltas).

## 9. Failure Durumları
- Contract not found (patch): safe fallback → last known activity type default yield; user notified.
- Snapshot missing: recompute from current state; note in reveal.
- RNG mismatch (versioned): deterministik replay only within same RNG version.

## 10. Edge Cases
- Time desync (client-server): server otoritedir.
- Server saati düzeltmesi (leap seconds/DST N/A because UTC): monotonic time referansı korunur.
- Player interrupt mid-reveal: transaction; partial reveal reddedilir, retry safe.
- Multi-device login race: first-in-wins reveal; second connection sees consistent post-reveal state.
- Manipulation attempt (clock modification): irrelevant; server clock only.

## 11. Server Authoritative Gereksinimleri
- Time source: single monotonic; leap-safe.
- Deterministic RNG (bkz. [Random System](random-system.md)).
- Idempotency: aynı reveal 2. kez çalışmaz.
- Full audit trail: reveal → ledger entries.
- Snapshot versioning: karakter state'i migration'a dayanıklı.

## 12. Anti-Manipulation
- Client'a hiçbir XP/gold değer güvenilmez.
- Contract creation/modification transaction'lı.
- Rate limit: reveal max 1 per second per character (spam koruması).
- Anomaly detection: reveal yield beklenen banttan sapıyorsa flag.

## 13. Genişletilebilirlik
- Yeni activity_type = params şeması + simulate() implementation.
- Soft-cap curve modifier per event (Phase 02 event bonus).
- Multi-track parallel activity (Phase 02): main + pet + travel.

## 14. Örnek Senaryo — Mining AFK
```
Contract:
  activity_type = mining
  params = { spot_id = iron_valley, autoloot=true }
  started_at = t0
  seed = abc123

Reveal at t0+8h:
  elapsed = 28800000 ms
  soft_cap not triggered
  simulate ticks (every 30s → 960 ticks)
  yields: iron_ore ×420, stone ×80, rare_gem ×2, mining_xp 3800
  inventory: 46/40 → stopped at 6h30m, remaining 1h30m no yield
  status = paused (inventory_full)
  reveal_result persisted, mail overflow if configured
```

## 15. Örnek Senaryo — Hunting AFK
```
Contract:
  activity_type = hunting
  params = { spot = whispering_forest, mob_types = [wolf], autopotion=true }
  seed = def456

Reveal at t0+6h:
  simulate combat ticks (250 ms) but bulk chunks (60s aggregate)
  compute kills = ~620 wolves
  loot roll seeded per kill_index
  xp = 12400, deaths = 0
```

## 16. Balance Considerations
- Soft cap curve tuning: 24h full ideal; oyuncu her gün en az bir kez login olmalı motive olsun; ama 2 gün AFK de %60 verimle "kabul edilebilir".
- Yield hesap fonksiyonu her aktivite için deterministic ve saf; side effect yok.
- Death mid-AFK sıkıntısı: hunting sırasında death olursa 30 dk yield loss + respawn; tekrar aktif ol.

## 17. Interface
```
POST /reveal
  auth: character token
  body: { character_id, correlation_id }
  response: { reveal_summary, delta_events, next_state }
```

## 18. Test Stratejisi (Phase 01'e devir)
- Golden replay tests: aynı seed + params → aynı yield.
- Determinism test: 100k random contracts, tick simulation, hash match.
- Time drift injection: skew ±5min, sonuç değişmez.
- Overflow: envanter capacity edge test.
