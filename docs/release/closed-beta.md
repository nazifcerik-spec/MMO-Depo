# Closed Beta

## 1. Amaç
Kontrollü ölçekte gerçek oyuncu geri bildirimi + prod stack canlı stress test + kritik bug + balance data toplama.

## 2. Wave Modeli

| Wave | Oyuncu Sayısı | Süre | Odak |
|---|---|---|---|
| **W1** | 200-500 | 2 hafta | Kritik bug + crash rate + core loop feedback |
| **W2** | 1.000-2.000 | 2 hafta | Performance + balance + retention D1/D7 |
| **W3** | 3.000-5.000 | 2 hafta | Ekonomi denge + guild + market + PvE endgame |

Waves arası **bug fix + balance pass** 1-2 hafta.

## 3. Invite Modeli

- **Kaynak listeleri**:
  - Discord community waitlist.
  - Newsletter subscribers.
  - Influencer allocation (%10 slot).
  - Random signup pool (%40).
  - Targeted (curator seçim) (%50).
- **NDA** opsiyonel; player sharing izinli (public marketing benefit).
- **Access**: email code + client build download link (Play internal test track + TestFlight iOS + web).

## 4. Telemetry Odağı (Wave 1)

- **Crash-free session rate** ≥ %98 hedef.
- **CCU peak** vs planned capacity.
- **Login success rate**.
- **Reveal open rate** (push → app).
- **Session length dist.**
- **D1 retention** ≥ %35 hedef.
- **Bug report volume** + severity dağılımı.

Data platform (bkz. Phase 02 analytics) — Grafana + custom event stream.

## 5. Feedback Kanalları

- **Discord**: dedicated #beta-feedback channel.
- **In-app**: `POST /v1/support/feedback` → ticket system.
- **Bug report form**: screenshot + description + auto-attach client log.
- **Weekly survey**: NPS + top 3 pain + top 3 love.

## 6. Kritik Bug Response

- **P0**: server down, data loss, exploit → 1 saat içinde hotfix veya rollback.
- **P1**: major feature broken, crash rate >%10 → 24 saat içinde fix.
- **P2**: annoying bug, workaround yok → next wave öncesi fix.
- **P3**: minor visual → backlog.

Bkz. [bug-fix.md](bug-fix.md).

## 7. Wipe Politikası

Waves arası **karakter wipe olası** (progression reset):
- W1 → W2: wipe zorunlu (balance değişikliği geçici veri geçersiz).
- W2 → W3: wipe olası; player'a önceden duyuru.
- W3 → Open Beta: wipe zorunlu (fresh start).
- Open Beta → Release: **wipe yok** (progression korunur; ancak testers rewards Founder title vs).

## 8. Wave Kapanış Kriterleri

Wave X sonu için:
- Kritik bug ≥ %90 close.
- CCU peak ≥ planlanan capacity handled.
- Retention D1 ≥ hedef.
- Support ticket per 1k user < 50.
- SLO uptime > %99.
- Team review meeting → decision W(X+1) GO.

## 9. Reward for Testers

Post-launch:
- **Founder** title (permanent, retired after launch — DD-181 uyumlu).
- Beta-exclusive cosmetic (mount skin, tabard, banner).
- Small currency bonus (10k gold or gem package).
- Special role Discord.

## 10. Marketing Handling

- Closed Beta = public knowledge; streamers + review OK.
- NDA yalnızca **build stability** dışında yok (dev iterasyon geri bildirim izinli).
- **Do not share**: server-side implementation detail, admin credentials.

## 11. Cross-Reference
- [production-readiness.md](production-readiness.md), [bug-fix.md](bug-fix.md), [balance.md](balance.md)
- [open-beta.md](open-beta.md)

## 12. Metrikler
- `beta_wave_active_users_gauge{wave}`
- `beta_bug_reported_total{severity}`
- `beta_ticket_resolved_total{severity}`
- `beta_survey_nps_score`
- `beta_wave_retention_d1_pct`, `_d7_pct`

## 13. Bekleyen Kararlar
- Wave sayısı 3 default; ihtiyaç halinde W4 opt.
- Wipe politikası final duyuru — biz karar.
- Founder reward tier — biz karar.
- Beta build web + mobile parallel mi mobile-only başlangıç — Phase 01 client karar.
