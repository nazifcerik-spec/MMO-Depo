# Post-Launch (Ops + Phase 02 Roadmap)

## 1. Amaç
Release sonrası ilk 30-90 gün operational stability + player retention + Phase 02 development kick-off.

## 2. First 30 Days Focus

### Stability
- **Uptime SLO** ≥ %99.9 target.
- **Crash-free session** ≥ %99.
- **P0/P1 bug close** < 24h.
- **Support ticket resolve** < 48h P2, < 24h P1.

### Community
- **Daily Discord activity**.
- **Weekly patch notes** (bug fix summary).
- **Bi-weekly dev log** (upcoming features).
- **Monthly AMA** (community Q&A).

### Retention
- **D7 retention** > %25 hedef.
- **D30 retention** > %15 hedef.
- **Session length** stabil / artıyor.
- **AFK contract completion** > %80.

## 3. First 30 Days Playbook

### Week 1
- War room continues; 24h coverage.
- Daily standup + metric review.
- Hotfix ready path active.
- Community listening (Discord + Reddit + Twitter).

### Week 2-4
- War room → normal on-call.
- Weekly patch cadence.
- First **balance pass** based on real data (bkz. [balance.md](balance.md)).
- First **content update** (new event / seasonal opt).

## 4. Phase 02 Roadmap Kick-Off

Post-launch stable → Phase 02 feature development start:

### High-Priority Phase 02 Features (release doc'larda `Phase 02` flag'li)
- **Auction house** (bkz. [marketplace/auction.md](../marketplace/auction.md)).
- **Guild boss + raid + war** (bkz. [guild/](../guild/)).
- **Cosmetic store** + payment integration.
- **Achievement leaderboard**.
- **Trade tax opt** (%2 >10k g — DD-141).
- **Cross-server chat + region migration**.
- **Transmog full feature**.
- **Pet evolution + happiness/stamina UI complete**.
- **Farming + archaeology profession**.
- **Housing** (Phase 02 or 03 karar).
- **PvP foundation**.

### Phase 02 Prep
- Requirements gathering + design review each feature.
- Feature flag rollout strategy per feature.
- A/B test framework.
- Retention analytics deeper (BI dashboard).

## 5. Monetization Aktif (Phase 02)

V1 free-to-play; Phase 02:
- **Cosmetic store** (mount skin, tabard, aura, effect).
- **Gem currency** activate.
- **Founder Pack** pre-order.
- **Battle Pass** seasonal (Phase 02 mid-late).

Compliance:
- Store TR PayTR / Iyzico.
- Global Stripe.
- Play Store / App Store IAP.
- Ethics: no pay-to-win, no gambling, transparent odds cosmetic bundle.

## 6. Ops Runbook Refinement

Real incident data ile runbook update:
- Each post-mortem → action items → runbook update.
- Runbook wiki review monthly.
- Runbook drill quarterly.

## 7. Team Growth

Post-launch success:
- Additional SRE hire (2 → 3-4).
- Content team growth (designer + writer).
- Support team scale to ticket volume.
- Community manager (1 → 2 opt).

## 8. Product Metrics Cadence

Weekly:
- CCU + DAU/MAU.
- Retention D1/D7/D30.
- Ekonomi sink/source.
- Support ticket volume + resolution.

Monthly:
- Retention D60/D90.
- Content completion rate.
- Community NPS.
- Financial (if store live).

Quarterly:
- Strategic review.
- Roadmap adjust.

## 9. Long-Term (6+ ay)

- **Phase 02 launch** (auction + guild advanced).
- **Region expansion** (TR + EU → US + Asia).
- **Seasonal content** cadence (quarterly major event).
- **First expansion** (Phase 03) — Ascension prestige tier extension, new zone, new class opt.

## 10. Version Cadence

- **Patch** (v1.0.x): bug fix + minor balance — weekly.
- **Minor** (v1.x.0): new feature + content — monthly.
- **Major** (v2.0.0): expansion + big content — 6-12 ay.

## 11. Player Communication

- **In-game news feed**: launch update, patch notes, events.
- **Email digest**: weekly retention email (opt-in).
- **Push notification**: event + Phase 02 feature launch.
- **Discord announcement channel**: major update highlights.

## 12. Post-Mortem Culture

- Blameless retrospective every incident.
- Learnings shared team-wide.
- Runbook + code + monitoring improvement action items.
- Quarterly team retrospective (strategic view).

## 13. Post-Launch Metrics
- `postlaunch_incident_total{severity}`
- `postlaunch_deploy_total`
- `postlaunch_content_update_total`
- `postlaunch_dau_gauge`
- `postlaunch_retention_d30_pct`
- `postlaunch_nps_score`
- `postlaunch_revenue_daily` (Phase 02)

## 14. Cross-Reference
- [release.md](release.md), [balance.md](balance.md)
- [infra/](../infra/), [performance/](../performance/)
- Phase 02 doc katmanları (henüz yazılmadı — Phase 02 kick-off ile başlar)

## 15. Bekleyen Kararlar
- Phase 02 feature priority final — biz karar.
- Monetization launch timing — biz karar.
- Region expansion order — biz karar.
- Long-term expansion (v2.0) scope — Phase 02 sonu review.
