# Release v1.0 (Launch Day)

## 1. Amaç
Global launch day playbook. Herşey hazır; oyun canlı; büyük kitle geliyor.

## 2. Timeline (Launch Day)

```
T-24h  Freeze all non-critical changes; final content deploy
T-12h  Final infra check + on-call all-hands standby
T-6h   Cache pre-warm + CDN edge purge
T-3h   Wipe execute (Open Beta data purge; Founder rewards migrate)
T-1h   Release build final deploy; staging + prod smoke
T-0    LAUNCH — landing page → live, store live, social announce
T+1h   Peak CCU monitor; scale response if needed
T+6h   First patch window (hotfix ready)
T+24h  Launch retrospective start
T+72h  Post-launch review meeting
```

## 3. War Room Setup

- **Discord voice channel** + **Zoom** aktif 24h.
- **Attendees**: SRE lead, eng lead, product owner, community manager, support lead, DBA.
- **Rotasyon**: 8h shift; 3 shift günde.
- **Dashboards**: on wall — CCU, SLO, error rate, alerts.
- **Slack channels**: #launch-war-room, #launch-alerts, #launch-support.

## 4. Deploy Frozen

Launch günü:
- **Kod deploy**: sadece hotfix (P0/P1).
- **Content deploy**: FROZEN (T-24h sonrası).
- **Config change**: sadece feature flag toggle (SRE + product joint approval).

Freeze period T-24h → T+72h.

## 5. Scaling Response

- Auto-scale kural aktif.
- Manuel intervention gerekli case:
  - **CCU spike > planlanan**: manual scale up (K8s replica up).
  - **DB slow**: read replica traffic switch.
  - **Redis memory**: cluster shard rebalance.
  - **CDN cache miss spike**: pre-warm additional edges.

## 6. Communication Playbook

### Success Message
Blog + Twitter/X + Discord:
> **"Otter v1.0 launched! Sunucularımız açıldı. Hoş geldin."**
> Landing page: [game.io/play](https://game.io/play)

### Delay / Issue
- < 30 dk delay: silent fix; no announcement.
- > 30 dk delay: Twitter/X + status page + Discord notice.
- Server down: PagerDuty + war room + immediate public status page.

Template pre-written; sadece SRE lead + community manager onayı ile publish.

## 7. First 24 Saat Ana Metrikler

- **Peak CCU**: hedef vs actual.
- **Signup rate**: yeni account/dk.
- **Crash-free session**: gerçek zamanlı dashboard.
- **Error rate**: 5xx ratio.
- **Payment errors** (if store live): 0 tolerance.
- **Support ticket volume**: dedicated queue.

## 8. Founder Tester Reward Migration

Closed Beta + Open Beta testerları:
- **Founder title** attach.
- **Beta-exclusive cosmetic** grant (mount skin, tabard, banner).
- **10k gold + 1 gem pack** starter reward.
- Special Discord role.

Migration script pre-tested T-48h; execute T-3h.

## 9. Hotfix Protocol

Launch day hotfix:
- SRE + eng lead joint approval.
- Skip full CI (unit + smoke only).
- Canary skip if urgent + rollback ready.
- Immediate rollback if regression.

Pre-approved hotfix list:
- Ledger sum violation kill switch (feature flag).
- Rate limit tune (config).
- Content pack minor balance revert (git revert).

## 10. Payment (if live)

- Payment provider (Stripe / Play Billing / App Store) live.
- Test transactions T-48h.
- Support ticket queue prep.
- Chargeback / refund policy documented.

**V1 default: free-to-play, no store**. Store Phase 02 opt.

## 11. Post-Launch Retrospective

T+72h meeting:
- Timeline recap.
- What went well.
- What broke.
- Metric summary.
- Action items → post-launch phase.

## 12. Sürprizin Küçültme

- **Load test T4** (50k CCU) doğrulanmış.
- **Chaos drill** done (Redis kill, DB primary fail).
- **Rollback drill** quarterly.
- **Support runbook** hazır.
- **Legal + biz sign-off** T-72h.

## 13. Version Naming

- Release tag: `v1.0.0`.
- Hotfix: `v1.0.1`, `v1.0.2`.
- Next feature: `v1.1.0` (Phase 02 first minor).

Semver strict; breaking change never in patch or minor.

## 14. Cross-Reference
- [open-beta.md](open-beta.md), [post-launch.md](post-launch.md)
- [infra/ci-cd-prod.md](../infra/ci-cd-prod.md), [infra/rollback.md](../infra/rollback.md)
- [performance/capacity-planning.md](../performance/capacity-planning.md)

## 15. Metrikler
- `launch_ccu_peak_gauge`
- `launch_signup_total`
- `launch_hotfix_deployed_total`
- `launch_incident_open_gauge{severity}`

## 16. Bekleyen Kararlar
- Payment live at launch mi Phase 02 mi — biz karar.
- Region rollout timing (global sim vs staggered) — biz karar.
- Founder Pack availability (Phase 02) — biz karar.
- Launch trailer video partner — marketing.
