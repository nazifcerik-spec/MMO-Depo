# Bug Fix Cycle

## 1. Amaç
Beta wave süresince + wave arası bug triage, fix, deploy, verify döngüsü.

## 2. Triage Matrix

| Severity | Kriter | SLA | Fix Path |
|---|---|---|---|
| **P0 Critical** | Server down, data loss, security exploit, cascading crash | 1 saat | Hotfix path + auto rollback |
| **P1 High** | Major feature broken, crash rate > %10, kullanıcı ilerleyemez | 24 saat | Fast track deploy (canary skip) |
| **P2 Medium** | Feature buggy workaround var, UX broken, minor exploit | Wave sonu | Normal deploy sonraki gate |
| **P3 Low** | Visual glitch, typo, edge case | Backlog | Prioritize next milestone |

## 3. Bug Kaynağı

- **Player report**: in-app + Discord + email.
- **Automated crash report**: Sentry / Firebase Crashlytics.
- **Server error alert**: high 5xx spike, exception.
- **Anomaly detection**: security rule fire.
- **QA smoke**: post-deploy verification fail.
- **Support triage**: repeating ticket pattern.

## 4. Bug Ticket Yapısı

```
Title: [P1] Vendor buy sometimes charges twice
Reporter: player_xyz + Sentry auto
Severity: P1
Steps to reproduce:
  1. Vendor sword_iron için 450c → buy
  2. 500ms sonra tekrar buy → wallet -900c yerine -450c
Expected: idempotent; second call E_IDEMPOTENT_REPLAY
Actual: charged twice
Attachments: screenshot, session log, trace_id d3f7…
Environment: staging w1-b
Correlation ID: d3f7…
```

## 5. Fix Cycle

```
Triage → Assign → Reproduce → Root Cause → Fix
  → Unit test (regression guard) → PR review
  → CI + integration test → Merge
  → Canary staging → Verify → Prod
  → Ticket close + player notify (if reporter)
```

## 6. Hotfix Path (P0)

Emergency:
- Skip full CI (unit + smoke only).
- Skip canary if urgent + evidence rollback capability.
- Direct deploy prod + immediate verify.
- Post-hoc full CI on hotfix branch.
- Post-mortem within 48h.

Kullanım: exploit patched, sum violation kill switch, kritik memory leak.

## 7. Rollback vs Hotfix

Karar tree:
- Deploy N sonrası regression → **rollback N** öncelikli (< 5dk).
- Regression eski dönemden (>N gün önce) → **hotfix** gerekli.
- Data corruption ise **freeze writes + backup restore** (bkz. [infra/rollback.md](../infra/rollback.md)).

## 8. Bug Backlog Management

- **Weekly grooming**: P3 → deferred to milestone; obsolete close.
- **Wave planning**: P2 wave sonu fix zorunlu; P3 next wave gate.
- **Regression prevention**: her fix regression test zorunlu (bkz. [testing/](../testing/)).

## 9. Player Communication

- P0 outage: **status page update** + Discord + Twitter/X in 15 dk.
- P1 broad impact: patch note "known issue" + ETA.
- P2/P3: silent fix; patch notes.

Patch notes format:
```
v1.4.1 (2026-11-05)
### Bug Fixes
- Vendor double-charge issue fixed (P1).
- Trade offer hash mismatch on rapid changes (P2).
- Chat scroll jump on incoming message (P3).
```

## 10. Metrikler

- `bug_open_gauge{severity}`
- `bug_closed_total{severity}`
- `bug_mttr_seconds{severity}` — median + p95
- `bug_regression_reintroduced_total`
- `hotfix_deployed_total`

## 11. Cross-Reference
- [closed-beta.md](closed-beta.md), [performance-tuning.md](performance-tuning.md)
- [infra/rollback.md](../infra/rollback.md), [infra/ci-cd-prod.md](../infra/ci-cd-prod.md)
- [security/audit.md](../security/audit.md) — post-mortem template

## 12. Bekleyen Kararlar
- Sentry vs Firebase Crashlytics client crash — Phase 01 client team karar.
- Bug tracking tool (Linear / GitHub Issues / JIRA) — biz karar.
- Post-mortem publish policy (public retrospective sonrası) Phase 02 karar.
