# CI/CD (Continuous Integration & Delivery)

## 1. Amaç
Her PR'ın kalite kapılarından geçmesi, her merge'in deploy-edilebilir bir artifact üretmesi, prod deploy'un tek tuşla veya otomatik olması.

## 2. Prensipler
- **Trunk always green**: `main` her zaman deploy-edilebilir.
- **Fast feedback**: PR CI < 10 dk hedef; ideal < 6 dk.
- **Reproducible builds**: sha-pinned image, lockfiles committed.
- **Same artifact promoted**: dev → staging → prod aynı container image.
- **Rollback = redeploy previous tag**; asla "hotfix on prod host".

## 3. Pipeline Türleri

### 3.1 PR CI (her PR açıldığında + push)
Stages:
1. **Setup** — checkout, cache restore.
2. **Lint**
   - Go: `golangci-lint run`, `go vet`, `govulncheck`.
   - Flutter: `dart format --set-exit-if-changed`, `flutter analyze`.
   - Yaml/JSON schema check (content dosyaları).
   - Commit message check (Conventional Commits).
3. **Test**
   - Go unit: `go test -race -cover ./...`.
   - Go integration: testcontainers Postgres + Redis.
   - Flutter test: `flutter test --coverage`.
   - Deterministic replay: AFK Engine golden test suite.
4. **Build**
   - Go binaries (server, worker, migrate) — `-trimpath -ldflags`.
   - Flutter web build (verify only; upload artifact opsiyonel).
5. **Security scan**
   - `trivy` (dependency + image), `gitleaks` (secret scan).
6. **Report**
   - Coverage upload, size diff, SBOM artifact.

Fail-fast: lint fail → test skip. Test fail → build skip.

### 3.2 Main CI (merge sonrası)
- Full PR pipeline yeniden.
- **Docker image build & push** (tag: `git-sha`, `main`).
- Dev environment auto-deploy (`git-sha` image).
- Smoke test dev'de: health check + 3 kritik endpoint.

### 3.3 Release Pipeline (tag `v*.*.*`)
- Full pipeline.
- Signed image tag (`cosign` — Phase 02).
- SBOM + release notes generation.
- Staging deploy (manual approve → prod).
- Rollout monitoring (5–15 dk soak; auto-rollback thresholds).

### 3.4 Nightly / Scheduled
- Full integration suite (uzun).
- Load test (küçük seed'li).
- Dependency update PR (`renovate` veya `dependabot`).
- Backup restore rehearsal (haftalık).

## 4. Environment Promotion
```
PR         → ephemeral preview (opsiyonel; Phase 02)
merge main → dev (auto)
tag v*     → staging (auto) → prod (manual approve)
hotfix     → prod (approve, fast path)
```

## 5. Deployment Stratejisi
- **v1**: rolling deployment (1 pod at a time; readiness gate).
- **v2 hedef**: blue/green veya canary (ör. %5 → %25 → %100 progressive rollout).
- Post-deploy check:
  - Health probe green N sn.
  - Error rate < baseline + %5.
  - Latency p99 < baseline + %10.
- Auto-rollback: metrics trip → previous tag.

## 6. Migration Deploy
- **Migration önce, uygulama sonra**:
  1. `migrate` job runs → apply forward migration.
  2. Server pods rolling update.
- Breaking migration 2 aşamalı deploy (expand-contract):
  - Deploy A: dual-write (yeni sütun eklenir, kod hem eskiye hem yeniye yazar).
  - Data backfill job.
  - Deploy B: eskiyi drop.
- Migration idempotent olmalı; retry safe.

## 7. Secrets in CI
- Provider secret store (GitHub Actions Secrets / equivalent).
- OIDC-based cloud auth (federated identity) — long-lived key yok.
- Env variable injection at runtime only.
- Secret rotation quarterly; alarm on last-used-90-days.

## 8. Caching
- Go module cache.
- `go build` cache (persist across job).
- Flutter pub cache.
- Docker layer cache (BuildKit).
- Test container images.

## 9. Observability
- CI metrics: pipeline duration, success rate, flakiness count.
- Deploy events → Prometheus/Grafana annotation (correlation with prod metrics).
- Slack / e-mail notification: deploy start/success/failure/rollback.

## 10. Feature Flags
- Runtime feature flag service (Phase 02 hedef; v1'de env-based toggle yeterli).
- Flag TTL: 30 gün; süresi dolan flag PR ile cleanup.

## 11. Manual Gates
- Prod deploy manual approve (öneri: v1'de zorunlu; v2'de canary green ise otomatik).
- Migration destructive (drop column) manual approve.

## 12. Rollback
- **Definition**: previous image tag re-deploy.
- Migration rollback: rare; expand-contract sayesinde forward-fix tercih.
- Post-incident: retrospective + doc (`docs/incidents/YYYY-MM-DD.md`).

## 13. Anti-Patterns
- Pipeline'a arbitrary `sleep`.
- Manual step arasına unversioned artifact.
- Prod'da hotfix'i CI'ya sokmadan uygulamak.
- Aynı tag'i overwrite etmek (mutable tag).
- Test'i "flaky diye" retry ile örtmek (root cause fix).

## 14. İlişkiler
- [Git](git.md), [Docker](docker.md), [Environment](environment.md), [Config](config.md), [Logging](logging.md), [System Architecture](system-architecture.md)

## 15. Bekleyen Kararlar
- CI provider: GitHub Actions (öneri) vs GitLab CI vs CircleCI vs Buildkite.
- Ephemeral preview env: v1'de değer/maliyet karşılaştırması.
- Load test framework: `k6` (öneri) vs `Locust`.
- Rollout stratejisi (rolling → canary geçiş milestone).
