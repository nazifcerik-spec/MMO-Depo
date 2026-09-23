# Closed Beta → Release — Katalog

> **Faz 37**: Launch pipeline. Şimdiye kadar tanımlanan tüm sistemler kod olarak hazır olduktan sonra oyuncuya götürmenin **kademeli gate'leri**.

## 1. Kapsam
Bu klasör; Production Readiness → Closed Beta → Bug Fix → Balance Pass → Performance Tuning → Open Beta → Release (v1.0) → Post-Launch geçiş kapılarını tarif eder.

## 2. Yolculuk Haritası

```
                 [Faz 01 kod uygulama]
                          ↓
                 Production Readiness
                          ↓
                 ┌────Closed Beta (Wave 1)──┐
                 │  (invite-only, 200-500)  │
                 └──────────────┬───────────┘
                                ↓
                 ┌────Bug Fix + Balance + Perf─┐
                 │   (2-6 hafta iterasyon)     │
                 └──────────────┬──────────────┘
                                ↓
                 ┌────Closed Beta (Wave 2/3)──┐
                 │  (invite genişleme, 2-5k) │
                 └──────────────┬─────────────┘
                                ↓
                 ┌────Open Beta──────────────┐
                 │  (public 10-25k, 4-8 hf) │
                 └──────────────┬────────────┘
                                ↓
                 ┌────Release v1.0───────────┐
                 │  (global launch)           │
                 └──────────────┬────────────┘
                                ↓
                       Post-Launch Ops
                       (Phase 02 kick-off)
```

## 3. Prensipler (DD-208)

- **Kademeli**: her gate önceki başarı olmadan geçilmez.
- **Data-driven decisions**: her gate metrik + geri bildirim ile geç/kal.
- **Player-first**: beta oyuncuları gerçek etki sahibi; onların feedback'i priority.
- **Reversible**: her gate rollback path'i açık.
- **Content veto**: content pack + balance her wave ile yenilenebilir.
- **Freeze windows**: bir wave sırasında breaking change **YOK**; sadece hotfix.
- **Public communication**: her aşama duyuru + shipping notes.

## 4. Dokümanlar
- [production-readiness.md](production-readiness.md) — Go-live checklist (tech + biz + legal).
- [closed-beta.md](closed-beta.md) — Invite wave planı + katılımcı + telemetry.
- [bug-fix.md](bug-fix.md) — Bug triage matrix + fix cycle + hotfix workflow.
- [balance.md](balance.md) — Balance pass methodology + telemetry input.
- [performance-tuning.md](performance-tuning.md) — Perf regression + baseline compare.
- [open-beta.md](open-beta.md) — Public beta capacity + marketing hazırlık.
- [release.md](release.md) — v1.0 Launch day playbook.
- [post-launch.md](post-launch.md) — Post-launch ops + Phase 02 roadmap kick-off.

## 5. Cross-Reference
- [performance/](../performance/) — Load test + capacity plan
- [infra/](../infra/) — Prod infra + backup + monitoring
- [security/](../security/) — Anti-cheat + audit
- [testing/](../testing/) — Test coverage
- [admin/](../admin/) — Support tooling

## 6. Ana Metrikler (Launch KPI)

- **Concurrent users (CCU)** peak.
- **DAU / MAU** ratio.
- **D1 / D7 / D30 retention**.
- **Session length + count per user**.
- **AFK contract completion rate**.
- **Reveal open rate** (push notification).
- **Crash-free session rate**.
- **Support ticket per 1k user**.
- **Bug open + close rate**.
- **Balance metric**: gold sink/source ratio, quest completion, class win rate.

## 7. Timeline (Realistic)

| Faz | Süre | Katılımcı |
|---|---|---|
| Production Readiness | 2-4 hafta | Full team |
| Closed Beta W1 | 2 hafta | 200-500 invite |
| Bug fix + balance | 2-4 hafta | Dev + QA |
| Closed Beta W2 | 2 hafta | 1-2k invite |
| Closed Beta W3 | 2 hafta | 3-5k invite |
| Open Beta | 4-8 hafta | 10-25k public |
| Release v1.0 | Launch day | Global |

Toplam: 3-5 ay Production Readiness → v1.0.

## 8. Bekleyen Kararlar
- Region priority launch (TR + EU önce; Phase 02 global) — biz karar.
- Store distribution (Play Store + App Store + Web) Phase 01 kayıt hazırlık.
- Monetization aktif mi Release'te (v1'de free-to-play, ancak cosmetic store Phase 02) — biz karar.
- Founder pack pre-order Phase 02.
