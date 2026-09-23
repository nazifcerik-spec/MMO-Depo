# Automated Testing — Katalog

> **Faz 34**: Kritik sistemler için test stratejisi. Server-side Go + domain unit + integration + end-to-end.

## 1. Kapsam
Bu klasör, hangi sistem için hangi test seviyesinde hangi senaryoların **zorunlu** olduğunu tarif eder. Testler kod değil; **kural ve senaryo katalogudur**. Uygulama Phase 01 backend team.

## 2. Test Piramidi

```
       ┌─────────────┐
       │  E2E (5%)   │   flow-level, deployment smoke
       ├─────────────┤
       │  Contract   │   API şema + gRPC compat
       ├─────────────┤
       │Integration  │   DB + Redis + full domain wire
       │   (20%)     │
       ├─────────────┤
       │    Unit     │   fonksiyon + domain pure logic
       │   (70%)     │
       └─────────────┘
```

## 3. Katmanlar

| Katman | Amaç | Araç |
|---|---|---|
| **Unit** | Saf iş mantığı (formül, invariant) | Go `testing` + `testify` |
| **Property-based** | Invariant her random input için (RNG, sum=0) | `gopter` / `testing/quick` |
| **Integration** | Domain servis + Postgres + Redis (real) | testcontainers-go |
| **Contract** | API şema stable (breaking change guard) | JSON Schema snapshot + `oapi-codegen` |
| **E2E** | Kullanıcı senaryosu (login → play → logout) | k6 script + smoke playbook |
| **Load** | Concurrency + throughput | k6 + Grafana |
| **Chaos** | Failure injection (DB drop, Redis kill) | toxiproxy / littlekube (Phase 02) |
| **Security** | OWASP top 10 + fuzz | zap-baseline + go-fuzz |
| **Client widget** | Flutter UI logic + widget | `flutter_test` |

## 4. Prensipler (DD-192)

- **Fast unit** — < 100ms; TDD comfort.
- **Deterministic** — seed RNG (DD ile), zaman clock injection.
- **Isolate** — testcontainer fresh DB per test class.
- **No shared state** — parallel run safe.
- **Golden fixtures** — content pack test fixture ayrı (`testdata/fixtures/`).
- **Table-driven** — Go idiom; her senaryo satır.
- **Coverage ≥ 80%** core domain modülleri; utility ≥ 60%.
- **Mutasyon skoru ≥ 50%** kritik ledger (Phase 02 target).
- **CI fail on coverage drop > 2%**.

## 5. Zorunlu Test Kategorileri (Bu Klasör)

- [afk.md](afk.md) — Analytical hesap + contract snapshot.
- [combat.md](combat.md) — 15-step damage pipeline + threat.
- [loot.md](loot.md) — Deterministic RNG + personal roll + rarity.
- [crafting.md](crafting.md) — Queue + critical craft + mastery.
- [inventory.md](inventory.md) — Ownership + lock + transfer atomicity.
- [economy.md](economy.md) — Ledger integrity + idempotency + reconcile.
- [trade.md](trade.md) — P2P bait-and-switch + concurrent + block.
- [dungeon.md](dungeon.md) — Floor progression + AFK run reveal.
- [character.md](character.md) — Lifecycle + freeze + prestige.
- [class.md](class.md) — Spec tree + skill unlock + respec.

## 6. CI/CD Test Aşamaları

```
git push
  ↓
Lint (staticcheck, golangci-lint)
  ↓
Unit test (parallel; timeout 60s)
  ↓
Integration test (testcontainer; timeout 300s)
  ↓
Contract test (schema diff)
  ↓
[merge to main]
  ↓
E2E smoke (staging)
  ↓
Load baseline (nightly)
```

## 7. Test Data Yönetimi

- **Fixture pack**: `testdata/fixtures/*.yml` — race, class, item template, drop table, region.
- **Factory pattern**: `character.Factory().WithLevel(42).WithClass('warrior').Build()`.
- **Content pack seed**: her test suite başında `content.LoadTestPack('v1_baseline')`.
- **Time freeze**: `clock.Freeze(time.Date(2026,9,9,12,0,0))` deterministic contract expire.

## 8. Cross-Reference
- [architecture/ci-cd.md](../architecture/ci-cd.md)
- [backend/database.md](../backend/database.md) — testcontainer
- [systems/rng.md](../systems/rng.md) — seed injection

## 9. Metrikler
- `test_run_total{suite, result}`
- `test_duration_seconds{suite}`
- `test_coverage_percent{package}`
- `test_flake_total{test_name}` — Phase 02

## 10. Bekleyen Kararlar
- Mutation testing framework Go seçimi (Phase 02).
- Load test hedef (peak concurrent users) Phase 01 capacity planning.
- Chaos engineering tool (toxiproxy vs LitmusChaos) Phase 02.
- Client E2E aracı (integration_test vs Maestro) client team karar.
