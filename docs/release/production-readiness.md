# Production Readiness (Go-Live Checklist)

## 1. Amaç
Closed Beta'nın önündeki tüm zorunlu gate'ler. Her madde `done` olmadan Wave 1 başlamaz.

## 2. Technical Checklist

### Kod Kalitesi
- [ ] Core domain coverage ≥ 80% (bkz. [testing/README.md](../testing/README.md)).
- [ ] Ledger + RNG coverage ≥ 90-95%.
- [ ] Property-based tests: nightly 10k iteration passing.
- [ ] Golden fixture snapshot regression suite green.
- [ ] Lint (staticcheck + golangci-lint) zero critical.
- [ ] Contract test schema stable (breaking change guard).

### Performance
- [ ] Load test T3 (10k CCU) SLO passing (bkz. [performance/slo-slis.md](../performance/slo-slis.md)).
- [ ] Peak burst tested (login storm + world boss).
- [ ] 4h nightly soak drift < %5.
- [ ] Auto-scaling policies configured.

### Security
- [ ] All security kural aktif (rate limit, idempotency, validation, session).
- [ ] Threat model reviewed.
- [ ] Anomaly detection rule config live.
- [ ] Bug bounty policy legal reviewed (Phase 02 launch).
- [ ] PII audit + redaction verified.
- [ ] Argon2id parameters benchmarked prod-scale.

### Infra
- [ ] Terraform apply prod clean (dev + staging + perf + prod).
- [ ] HAProxy + TLS + WS proxy operational.
- [ ] CDN cache rules validated.
- [ ] Postgres HA + replica lag < 5s.
- [ ] Redis Sentinel/Cluster healthy.
- [ ] Backup + PITR + restore drill success (aylık).
- [ ] Monitoring dashboards deployed.
- [ ] Alertmanager routing PagerDuty + Slack test.
- [ ] Runbook wiki per alert.
- [ ] Deploy pipeline CI/CD green + canary tested.
- [ ] Feature flag provider (GrowthBook) live.
- [ ] Secret manager (Vault/AWS SM) populated + rotation policy.

### Data
- [ ] Content pack v1.0 signed + deployed CDN.
- [ ] Item template ~500 (weapon/armor/consumable/material).
- [ ] Recipe ~200 (crafting).
- [ ] Enemy template ~150.
- [ ] Dungeon template ~15.
- [ ] Region + territory + location content complete.
- [ ] Race (~8) + Class (~5) + Spec (~15) content complete.
- [ ] Quest content main line ~30 + side ~100.
- [ ] Vendor template ~50 (per major settlement).
- [ ] Drop table balance pass done.

### Client
- [ ] Flutter build (Android/iOS/Web) production release.
- [ ] Play Store + App Store account + app registered.
- [ ] App signing certificate + upload key secured.
- [ ] Localization (TR + EN) ≥ %95 coverage.
- [x] Client crash reporting — Dart hataları birinci taraf toplayıcıya (DD-235). Native (NDK/iOS) çöküşler için Crashlytics: Firebase hesabı bekliyor.
- [ ] Deep link routing test all screens.
- [ ] Notification (FCM + APNs) live.
- [ ] Cold start < 3sn hedef doğrula.

### Legal / Business
- [ ] Terms of Service published.
- [ ] Privacy Policy published (GDPR/KVKK compliant).
- [ ] EULA final.
- [ ] Age gate (13+) enforcement.
- [ ] Region compliance (KVKK for TR, GDPR for EU).
- [ ] Community Code of Conduct.
- [ ] Support email + ticket system live.
- [ ] Beta NDA (opsiyonel closed beta için).

### Ops Support
- [ ] Admin panel operational + SSO + WebAuthn.
- [ ] RBAC 6 rol tanımlı + admin kullanıcı listesi.
- [ ] Support runbook (mute, ban, reversal, item restore) hazır.
- [ ] Incident response runbook (P0-P3).
- [ ] Post-mortem template.
- [ ] On-call rotation weekly schedule.
- [ ] Escalation policy PagerDuty.
- [ ] Status page (Statuspage / Betteruptime) live.

## 3. Go/No-Go Meeting

Production Readiness sonu — cross-functional review:
- **Attendees**: eng lead, SRE lead, QA lead, product owner, legal, support lead, CEO.
- **Decision**: GO / NO-GO Closed Beta W1 açılışı.
- **Criteria**: yukarıdaki checklist ≥ %95 done + kritik %100 done.
- **Sign-off**: her rol yazılı onay (Slack thread veya doc).

## 4. Kritik Checklist (Zorunlu %100)

Bu maddeler eksik → NO-GO:
- [ ] Backup + restore drill success.
- [ ] Rollback procedure tested.
- [ ] Ledger reconciler operational.
- [ ] Security anomaly rules active.
- [ ] Legal (ToS + Privacy) published.
- [ ] Age gate enforcement.
- [ ] Support ticket flow.
- [ ] On-call rotation.
- [ ] PagerDuty alerts.
- [ ] Monitoring dashboards.

## 5. Nice-to-Have (Beta İçin Ertelenebilir)

- Auction (Phase 02 primary).
- Full transmog feature.
- Housing.
- Cross-region.
- WebAuthn player-side (v1 TOTP yeter).
- Rich notification (image + action button).
- Achievement leaderboard.

## 7. Durum (2026-10-10)

Kodla karşılanabilen kritik maddeler:
- [x] Core domain coverage ≥ 80% — birleşik birim + entegrasyon kapısı CI'da (DD-223).
- [x] Backup + restore drill success — `scripts/restore-drill.sh` yerelde geçti (DD-232); üretimde aylık tekrarlanmalı.
- [x] Rollback procedure — `scripts/rollback.sh`, ileri uyumlu migration kuralı (DD-232); üretimde bir kez denenmeli.
- [x] Ledger reconciler operational — worker `economy_reconciler`.
- [x] Age gate (13+) — kayıtta zorunlu onay (DD-231).
- [~] Legal (ToS + Privacy) published — `/terms`, `/privacy` taslak; hukuki inceleme + `LEGAL_OPERATOR`/`SUPPORT_EMAIL` gerekli.
- [x] Hesap silme + veri dışa aktarma (KVKK/GDPR, Play Store şartı) — DD-231.

Hesap / karar bekleyenler (geliştirme dışı): sunucu + alan adı, Play Store / App Store hesapları ve uygulama imzalama anahtarı, Firebase (push), SMTP sağlayıcısı, nesne depolama (yedek kopyası), destek e-postası, hukuki inceleme, uyarı kanalı (PagerDuty/Slack), durum sayfası.

## 6. Cross-Reference
- [closed-beta.md](closed-beta.md), [release.md](release.md)
- [infra/](../infra/), [security/](../security/), [testing/](../testing/), [performance/](../performance/)

## 7. Metrikler (Readiness Progress)
- `readiness_checklist_completion_pct{category}`
- `readiness_critical_blocker_count`
- `readiness_go_meeting_scheduled_at`

## 8. Bekleyen Kararlar
- Kritik %100 vs total %95 dengesi Go/No-Go — team final karar.
- Beta NDA gerekli mi (closed beta) — legal karar.
- Founders program (paid access early) — biz karar Phase 02.
