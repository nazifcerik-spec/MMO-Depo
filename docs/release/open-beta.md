# Open Beta

## 1. Amaç
Public erişim + marketing kick-off + T3/T4 capacity doğrulama + Release öncesi son büyük stress test.

## 2. Kapsam

- **Public**: kayıt açık; herkes indirebilir.
- **Süre**: 4-8 hafta.
- **Ölçek hedef**: 10-25k CCU peak; 50-100k DAU.
- **Region**: launch region (TR + EU) primary; global opt-in.
- **Wipe**: Open Beta sonu **wipe zorunlu**; Release fresh start (Founder tester rewards korunur).

## 3. Marketing Hazırlık

- **Landing page live** (waitlist → open signup).
- **Discord community** genişletme.
- **Content creator kit** (banner, video, key art).
- **Press outreach** (mmorpg news sites).
- **Store listing** (Play Store + App Store + Web).
- **Social launch trailer** (30-60 sn).
- **Website + FAQ + support docs public**.

## 4. Infra Hazırlık

Bkz. [performance/capacity-planning.md](../performance/capacity-planning.md):
- **T3 baseline** infra active + **T4 hazır** (auto-scaling threshold).
- Multi-node WS (3-4 node target).
- Postgres 2 read replica.
- Redis cluster 3-shard.
- CDN cache pre-warmed content.
- On-call rotation reinforced (2 SRE + eng lead backup).
- Status page monitoring live.

## 5. Kritik Metrikler

- **Peak CCU**.
- **CCU growth curve** (week over week).
- **D1/D7/D30 retention**.
- **Crash-free session rate** ≥ %99.
- **SLO uptime** ≥ %99.5.
- **Support ticket/1k**.
- **Bug open/close rate**.
- **Community sentiment** (Discord NPS, social mention).

## 6. Content Pipeline

Open Beta sırasında:
- **Weekly patch** (bug fix + minor balance).
- **Bi-weekly content update** (yeni quest, event opt).
- **Seasonal event** (opsiyonel launch prep) — beta test.

## 7. Community Management

- **Community Manager** aktif Discord + social.
- **Weekly dev log** (blog post + video özet).
- **Bug bounty small** (Discord role reward reporting).
- **AMA sessions** (Reddit / Twitch weekly).

## 8. Payment Test (Phase 02 Prep)

Open Beta'da monetization **YOK** (free) — Release'te v1'de free-to-play + Phase 02'de cosmetic store.
Ancak **Founder Pack pre-order** opt (Phase 02 monetization prep) — biz karar.

## 9. Wipe Duyuru

Open Beta start:
- Blog + Discord + in-app banner: "Open Beta karakteri Release'te wipe olacak."
- Founder reward (bkz. [closed-beta.md](closed-beta.md)) korunur.

Open Beta sonu 1 hafta önce:
- Countdown banner.
- Screenshot / achievement backup export (Phase 02 opt).

## 10. Kapanış Kriterleri (Release GO/NO-GO)

Release'e geçiş için:
- Kritik bug ≥ %95 close.
- Crash-free session ≥ %99.
- SLO uptime ≥ %99.5.
- T4 (50k CCU) load stress passed live.
- Retention D7 ≥ %20.
- Community NPS ≥ 40.
- Support scale (ticket/1k) < 30.
- Legal + biz + ops sign-off.

## 11. Cross-Reference
- [closed-beta.md](closed-beta.md), [release.md](release.md)
- [performance/capacity-planning.md](../performance/capacity-planning.md)
- [balance.md](balance.md), [bug-fix.md](bug-fix.md)

## 12. Metrikler
- `open_beta_ccu_gauge`
- `open_beta_dau_gauge`
- `open_beta_retention_d1/d7/d30_pct`
- `open_beta_new_signup_total`
- `open_beta_crash_free_pct`

## 13. Bekleyen Kararlar
- Region launch order — biz karar.
- Founder Pack pre-order Phase 02 fiyat + içerik.
- Content update cadence (weekly vs bi-weekly) — content team.
- Community bug bounty scope — Phase 02.
