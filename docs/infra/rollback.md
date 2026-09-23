# Rollback Procedures

## 1. Amaç
Kötü deploy veya data corruption durumunda hızlı, güvenli, veri-koruyucu geri dönüş.

## 2. Rollback Türleri

| Tür | Trigger | MTTR |
|---|---|---|
| **App version** | 5xx spike, latency regression | < 5 dk |
| **Feature flag** | Behavior bug | < 1 dk |
| **Database migration** | Schema issue, data corruption | 15-60 dk |
| **Content pack** | Balance bug, exploit content | 5-10 dk |
| **Full data restore** | Widespread data loss / ransomware | 1-2 saat (RTO) |

## 3. App Version Rollback

### Automatic (Canary Abort)
Argo Rollouts detected regression → abort → revert to previous ReplicaSet automatic.

### Manual
```bash
# Option 1: Argo Rollouts abort
kubectl argo rollouts abort api-rollout
kubectl argo rollouts undo api-rollout

# Option 2: kubectl rollout undo
kubectl rollout undo deployment/api

# Option 3: Deploy previous image tag (safest — explicit version)
kubectl set image deployment/api api=registry/api:v1.3.9
kubectl rollout status deployment/api
```

Verify:
- `/live` + `/ready` OK.
- Error rate baseline.
- SLO recovery.

## 4. Feature Flag Rollback

Fastest — no redeploy:
```
GrowthBook / LaunchDarkly dashboard → flag toggle OFF
```

Immediate effect (SDK poll interval 30s).

## 5. Database Migration Rollback

### Reversible Migration
Every migration has `.down.sql`:
```bash
sqlx migrate down 1
```

### Irreversible (Data Loss Risk)
Migration flagged `IRREVERSIBLE` → down step not provided → **restore from backup**.

### Prevention
- Reversibility test in CI (`up → down → up` idempotent).
- Destructive migrations 2-step deploy (backward-compat first).

## 6. Content Pack Rollback

Data-driven content:
```bash
# Git revert commit
git revert <content-commit>
git push
# CI redeploys content pack → CDN
# App polls content_version → auto reload
```

Time < 10 dk.

## 7. Full Data Restore (Disaster Recovery)

Trigger: widespread data corruption, ransomware, primary + replica loss.

Procedure:
1. **Declare incident** — PagerDuty critical + war room.
2. **Freeze writes** — put app in read-only mode (config flag).
3. **Determine RPO target** — hangi zamana restore (WAL replay).
4. **Provision fresh DB**.
5. **Restore base backup** → `pg_basebackup` from S3.
6. **WAL replay** to target timestamp.
7. **Verify data integrity** — sample queries + reconciler.
8. **Redirect app** to new DB (config swap).
9. **Enable writes**.
10. **Monitor** 24h.

RTO target: 1-2 saat. Full runbook wiki'de.

## 8. Redis Rollback

- Redis less critical (rebuildable from Postgres).
- Snapshot restore:
  ```bash
  aws s3 cp s3://backup/redis/dump.rdb.gz .
  gunzip dump.rdb.gz
  redis-cli DEBUG RELOAD
  ```
- User session invalidate → forced re-login. Acceptable.

## 9. Rollback Decision Tree

```
Something wrong?
  ├── Config / flag issue? → Feature flag toggle OFF (1 dk)
  ├── App code bug (recent deploy)?
  │     ├── Auto-rollback triggered? → wait; monitor
  │     └── Not triggered? → Manual `kubectl rollout undo` (5 dk)
  ├── Migration bad?
  │     ├── Reversible? → `sqlx migrate down` (5 dk)
  │     └── Irreversible? → Restore from backup (60 dk)
  ├── Content bug?
  │     └── Git revert + redeploy content pack (10 dk)
  └── Widespread data loss?
        └── DR runbook (1-2 saat)
```

## 10. Communication

Rollback in progress:
- Slack #incidents post start.
- Status page update (customers.io / Statuspage).
- ETA + impact summary.
- Post-mortem within 48h.

## 11. Post-Rollback

- Verify SLO recovery.
- Post-mortem doc (blameless):
  - Timeline.
  - Impact (users affected, revenue impact).
  - Root cause.
  - Fix.
  - Prevention (test coverage add, monitoring improve).
  - Action items + owners + deadlines.

## 12. Rollback Drill

Quarterly:
- Practice scenario: app rollback, DB restore, content revert.
- New team members participate.
- Update runbook based on findings.

## 13. Metrikler
- `rollback_total{kind, reason}`
- `rollback_duration_seconds{kind}`
- `rollback_success_total`
- `mttr_seconds{severity}`

## 14. Cross-Reference
- [ci-cd-prod.md](ci-cd-prod.md), [backup.md](backup.md)
- [postgres-prod.md](postgres-prod.md), [redis-prod.md](redis-prod.md)
- [alerting.md](alerting.md), [monitoring.md](monitoring.md)

## 15. Bekleyen Kararlar
- Status page tool (Statuspage / Betteruptime) — Phase 01.
- War room tooling (Zoom / Slack Huddle) — ops karar.
- Rollback drill sıklık (quarterly v1) Phase 02 monthly.
- Automated DR failover script Phase 02.
