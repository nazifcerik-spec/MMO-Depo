# Load Test: Dungeon

## 1. Amaç
Dungeon entry + AFK run reveal + party dungeon + weekly lockout enforcement.

## 2. Senaryolar

### 2.1 Steady Dungeon Entry
- Yük: T3'te 3000 active × 0.3 dungeon/saat = **~0.25/sn**.
- Solo (%60) + party (%40).
- Beklenti: enter p95 < 600ms.

### 2.2 Dungeon AFK Run Reveal
- Yük: 3000 × 0.3 × 0.7 AFK mode = **0.18/sn reveal**.
- Analitik floor progression + loot rolls.
- Beklenti: reveal p95 < 500ms.

### 2.3 Weekly Reset Burst
- Salı 06:00 UTC — reset job → tüm lockout marker clear + oyuncular rush enter.
- 10 dk içinde ~30% active player heroic dungeon enter.
- Peak: ~50 enter/sn (T3).
- Beklenti: p95 < 1s tolerated; DB write throughput OK.

### 2.4 Party Dungeon Coordination
- 5 üye party dungeon enter → 5 concurrent DB write.
- Beklenti: party members sync join; no drift.

## 3. Kritik Path — Enter

1. Auth + rate limit.
2. Template fetch (content cache).
3. Level range check.
4. Lockout check (heroic/mythic — Redis: `dungeon:lockout:{cid}:{tmpl}:{diff}`).
5. Party validation (if party mode).
6. Energy cost debit (ledger).
7. Dungeon run insert (DB).
8. WS notify party members.
9. Response.

## 4. Kritik Path — AFK Reveal

1. Auth + idempotency lookup.
2. Run fetch + snapshot.
3. **Analitik floor/room progression** — character DPS proxy × time / floor difficulty.
4. Boss defeat probability (final floor + DPS threshold).
5. Loot rolls per member (personal loot; deterministic seed).
6. Ledger + inventory insert batch.
7. Achievement/first-clear trigger.
8. Idempotency cache.
9. COMMIT + outbox.

## 5. Weekly Lockout Storm

Salı 06:00 UTC reset:
- Worker `dungeon_lockout_reset` at 06:00 → Redis MASSDELETE lockout keys.
- 10.000 CCU aktif olabilir; 3000 active × 0.3 rush = 900 enter attempt / dk.
- Peak burst: ~50 enter/sn.
- Preemptive scale + read replica prep önerilir.

## 6. Metrics

- `dungeon_enter_p95_ms`.
- `dungeon_reveal_p95_ms`.
- `dungeon_lockout_reject_total`.
- `dungeon_first_clear_total`.
- `dungeon_run_active_gauge`.
- `dungeon_run_analytical_calc_ms`.

## 7. Party Coordination

Party dungeon enter:
- Leader inits `POST /v1/dungeon/{tmpl}/enter { party_id }`.
- Server: her party üye için ayrı lockout check (personal lockout DD-115'e uygun).
- Ready check WS (v1'de opsiyonel; leader auto-confirm).
- Run instance member array frozen at start.

## 8. Beklenen Sonuç

T3 steady:
- Enter p95: ~400ms.
- Reveal p95: ~300ms.
- Lockout check < 5ms Redis.

T3 weekly reset burst 50/sn:
- p95: ~800ms.
- Peak recovery 10 dk.

## 9. Cross-Reference
- [dungeon](../dungeon/README.md), [dungeon/afk-run.md](../dungeon/afk-run.md)
- [testing/dungeon.md](../testing/dungeon.md)

## 10. Bekleyen Kararlar
- Ready check UI (v1 auto; Phase 02 explicit).
- Cross-region dungeon queue Phase 03.
- Peak burst throttle (rush-hour queue) Phase 02.
