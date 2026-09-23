# CI/CD Production Pipeline

## 1. Amaç
Kod commit → prod deploy sürecinin otomatik, test-gated, gözlenebilir, rollback-safe olması.

## 2. Pipeline Aşamaları

```
git push branch
  ↓
CI: lint + unit + build + trivy scan + cosign sign + push registry
  ↓
[PR merge to main]
  ↓
Auto deploy staging
  ↓
Smoke test staging (k6 lite + auth flow)
  ↓
[Manual gate: SRE approval] → prod deploy
  ↓
Canary deploy: 5% traffic → 30% → 100% (rolling)
  ↓
Post-deploy verification: health + SLO check 15 dk
  ↓
Auto-rollback if regression
```

## 3. Environment Progression

- **dev**: her commit auto-deploy.
- **staging**: her PR merge auto-deploy.
- **prod**: SRE approval + rollout window (business hours; Cuma 15:00 sonrası deploy freeze).
- **hotfix**: emergency fast-path (bkz. [rollback.md](rollback.md)).

## 4. Feature Flag

- **GrowthBook** (open-source) veya **LaunchDarkly** (managed).
- Deploy ≠ release; new feature behind flag OFF; SRE gradual enable.
- Rollback via flag toggle < 1 dk (kod redeploy gerektirmez).

## 5. Canary Deploy

- K8s `Deployment` + `Argo Rollouts` (canary strategy).
- Traffic split: 5% → 30% → 100%; her adım 10-30 dk observation.
- Analysis:
  - Error rate < baseline + 0.5%.
  - p95 latency < baseline × 1.2.
- Fail → auto rollback (bkz. [rollback.md](rollback.md)).

## 6. Deploy Trigger

- Git tag `v1.4.0` → prod pipeline.
- Manual `gh workflow run prod-deploy` (SRE approval workflow).
- Argo CD GitOps (Phase 02) — sync automated after approval.

## 7. Zero-Downtime Deploy

- Rolling update Kubernetes.
- Old pod graceful drain (SIGTERM → drain WS 30s → stop).
- New pod ready check `/ready` → 200 → traffic route.
- Database migration önce (backward compat) → app deploy sonra.
- **No breaking API change** in single deploy (bkz. contract versioning).

## 8. Database Migration Deploy

Bkz. [postgres-prod.md §Migration](postgres-prod.md#8-migration):
- Migration önce apply (backward compat safe).
- App new version sonra.
- 2 aşamalı destructive: mesaj kodu deploy 1'de mark deprecated; deploy 2'de remove.

## 9. Rollback

Auto-rollback triggers:
- Health check `/ready` fail 3 dk.
- 5xx rate > 10% for 5 dk.
- SLO error budget consume > 2× normal rate.

Manual rollback: Argo Rollouts abort + revert to previous ReplicaSet.

Full rollback: `kubectl rollout undo deploy/api` (< 1 dk).

## 10. Approval Gate

Prod deploy PR:
- CODEOWNERS review + SRE approve.
- CI green.
- Test coverage > threshold.
- Changelog entry.
- Deploy window check (no Cuma pm, no maintenance window overlap).

## 11. Deploy Metrics

- Frequency: hedef 3-5 deploy/hafta.
- Lead time: PR merge → prod < 4 saat.
- MTTR (Mean Time To Recovery): < 30 dk.
- Change failure rate: < 15% (auto-rollback triggered).

DORA metrics tracked.

## 12. Notification

- Slack #deploys: her deploy start + success/fail.
- Post-deploy report: changelog + SLO impact.

## 13. Content Pack Deploy

Data-driven content:
- YAML file git commit → CI validate → sign → CDN upload.
- App runtime `content_version` header check → hot reload.
- Rollback: revert git commit + redeploy content pack.

Content version stored per feature; deploy Boolean check.

## 14. Secret Deploy

- Vault (or AWS SM / Doppler) → env inject at pod start.
- Secret rotation: quarterly (or on incident).
- Rotation rolling deploy triggered.

## 15. Metrikler
- `deploy_success_total{env, service}`
- `deploy_rollback_total{env, service, reason}`
- `deploy_duration_seconds{env, service}`
- `canary_promoted_total`
- `canary_aborted_total{reason}`

## 16. Cross-Reference
- [rollback.md](rollback.md), [monitoring.md](monitoring.md)
- [docker.md](docker.md), [postgres-prod.md](postgres-prod.md)
- [architecture/ci-cd.md](../architecture/ci-cd.md) — v1 temel

## 17. Bekleyen Kararlar
- GitOps (Argo CD / Flux) Phase 02.
- Feature flag provider (GrowthBook v1 öneri) Phase 01.
- Canary tool (Argo Rollouts / Flagger) Phase 01.
- Deploy freeze policy weekend (Cuma pm freeze) — ops karar.
