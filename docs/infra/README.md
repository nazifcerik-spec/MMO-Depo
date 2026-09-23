# Production Infrastructure — Katalog

> **Faz 36**: Prod ortam altyapısı. Şimdiye kadar tanımlanan tüm sistemlerin **çalıştığı zemin**.

## 1. Kapsam
Bu klasör; container image, reverse proxy, CDN, TLS, Postgres/Redis prod deployment, backup, monitoring, logging, alerting, CI/CD ve rollback prosedürlerini tarif eder. `architecture/` klasörü mimari kararları içerir; **`infra/` operasyonel çalıştırma detayı**.

## 2. Prensipler (DD-202)

- **Immutable infrastructure** — sunucular tekrar tekrar üretilebilir; SSH ile canlı değişiklik yasak.
- **Everything as code** — Terraform (infra), Helm/Kustomize (deploy), GitOps (ArgoCD Phase 02).
- **12-factor** uygulama; config env var + secret manager.
- **Zero-downtime deploy** — blue/green veya rolling; DB migration reversible.
- **Observability by default** — metric + log + trace her serviste; correlation_id propagation.
- **Backup + restore drill** — aylık restore test; RTO/RPO explicit.
- **Least privilege** — service account, network policy, IAM role.
- **Cost-aware** — resource limit + auto-scale; sizing per tier ([performance/capacity-planning.md](../performance/capacity-planning.md)).

## 3. Dokümanlar
- [docker.md](docker.md) — Prod image build + distroless + registry + vulnerability scan.
- [reverse-proxy.md](reverse-proxy.md) — nginx/HAProxy + TLS termination + WS proxy + LB.
- [cdn.md](cdn.md) — Cloudflare + static asset + edge cache.
- [ssl.md](ssl.md) — TLS cert + ACME + rotation.
- [postgres-prod.md](postgres-prod.md) — HA + replication + connection pool + tuning.
- [redis-prod.md](redis-prod.md) — Cluster + persistence + eviction + Sentinel.
- [backup.md](backup.md) — Postgres WAL + Redis snapshot + S3 offsite + retention.
- [monitoring.md](monitoring.md) — Prometheus + Grafana + Jaeger + dashboards.
- [logging.md](logging.md) — Structured JSON + Loki + retention + query.
- [alerting.md](alerting.md) — Alertmanager + PagerDuty + Slack routing + runbook.
- [ci-cd-prod.md](ci-cd-prod.md) — Deploy pipeline + canary + feature flag.
- [rollback.md](rollback.md) — Rollback procedures + DB migration reverse + incident response.

## 4. Ortamlar

| Env | Amaç | Data | Access |
|---|---|---|---|
| **dev** | Local development | seed | herkes |
| **staging** | Prod-like test | anonymized snapshot | dev + qa |
| **perf** | Load test | fixture 250k | perf team |
| **prod** | Production | live | ops (SRE gate) |

Her env ayrı kubernetes namespace + ayrı secret vault.

## 5. Deploy Cadence

- **dev**: her commit auto-deploy.
- **staging**: her PR merge auto-deploy + smoke.
- **prod**: gated (SRE approval + rollout window; günde 1-3 deploy hedef).
- **hotfix**: emergency path (bkz. [rollback.md](rollback.md)).

## 6. Cross-Reference
- [architecture/](../architecture/) — Mimari kararlar (docker/ci-cd/env/config/logging temel)
- [performance/capacity-planning.md](../performance/capacity-planning.md)
- [security/](../security/) — Zero-trust internal
- [admin/](../admin/) — Ops panel

## 7. Metrikler (Infra KPI)
- `deploy_success_total{env, service}`
- `deploy_rollback_total{env, service, reason}`
- `deploy_duration_seconds{env, service}`
- `backup_success_total{kind, target}`
- `backup_restore_test_success_total`
- `sla_uptime_pct{service}`
- `incident_open_gauge{severity}`

## 8. Bekleyen Kararlar
- Cloud provider final (AWS / GCP / Hetzner) — ops karar.
- Orchestrator (Kubernetes / Nomad / ECS) — DD-OPEN-9.
- GitOps tool (ArgoCD / Flux) — Phase 02.
- Secret manager (Vault / AWS SM / Doppler) — Phase 01.
- Cost dashboard (Kubecost / Vantage) — Phase 02.
