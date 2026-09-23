# Admin Panel — Katalog

> **Faz 32**: Sunucu operatörlerinin oyun state'ini gözlemleme, moderasyon aksiyonu, support ticket çözme ve içerik doğrulama aracı. **Ayrı web tool** (Flutter uygulama değil).

## 1. Kapsam
Admin panel, geliştirici + ops + support ekibinin kullandığı **read-first, write-audited** araçtır. Oyunun canlı state'ini görür, sorunları teşhis eder, sınırlı düzeltmeler yapar.

**Prensipler (DD-181):**
- **Read-first**: default view read-only; write aksiyonlar açık ikinci onay ile.
- **Server API kullanır**: DB'ye direct query yasak; her aksiyon `/internal/ops/...` endpoint üzerinden ledger/audit yazar.
- **RBAC**: rol-tabanlı yetki; en az yetki prensibi.
- **2FA zorunlu**: her admin login + write aksiyonlar step-up authentication.
- **Full audit**: her sayfa view, her aksiyon `audit.admin_action`'a yazılır (operator_id, target, before/after diff, reason).
- **Rate limit**: bot admin tool koruma (dakikada N istek).
- **Ayrı domain**: `admin.game.internal` (VPN veya SSO gated); public'e açık değil.

## 2. Teknoloji

- **Frontend**: React + TypeScript + shadcn/ui (sunucu team'i familiarity); Phase 02'de opsiyonel Retool alternatifi.
- **Backend**: Go, aynı monolith içinde `/internal/ops/*` endpoint namespace.
- **Auth**: SSO (OIDC / Google Workspace / Okta) + WebAuthn 2FA; oyun account sisteminden ayrı.

## 3. Dokümanlar
- [auth-rbac.md](auth-rbac.md) — Admin auth + role matrix + 2FA.
- [players.md](players.md) — Account view + ban/mute/lock.
- [characters.md](characters.md) — Character inspect + inventory + freeze.
- [items.md](items.md) — Item template + instance browser + audit trail.
- [economy.md](economy.md) — Ledger viewer + wallet + sink/source dashboard + reversal.
- [guild.md](guild.md) — Guild list + members + treasury + disband override.
- [market.md](market.md) — Listing viewer + admin cancel + auction inspect.
- [combat.md](combat.md) — Active combat session inspect + replay (Phase 02).
- [activities.md](activities.md) — Active AFK contract inspect + force-reveal.
- [dungeons.md](dungeons.md) — Active runs + template rollout.
- [events.md](events.md) — World event scheduler + manual trigger.
- [logs.md](logs.md) — Audit log viewer + filter + export.

## 4. Cross-Reference
- [database/audit.md](../database/audit.md) — Audit schema
- [auth](../auth/README.md) — Player account (admin ayrı sistem)
- [economy/transaction.md](../economy/transaction.md) — Ledger
- [systems/*](../systems/) — İçerik hakkı olan sistemler

## 5. Erişim Kontrolü Özeti (Detay: auth-rbac.md)

| Role | Sample Access |
|---|---|
| `viewer` | Read-only tüm liste; PII gizli |
| `support_l1` | Player detail, item history, mail read, mute 24h |
| `support_l2` | Reversal < 100k gold, item restore, ban 7d |
| `game_designer` | Event trigger, content pack preview, no player write |
| `ops_engineer` | Server config, feature flag, worker restart |
| `admin` | Tüm aksiyon + role management |

## 6. Metrikler
- `admin_page_view_total{page, role}`
- `admin_action_total{action, role, target_kind}`
- `admin_action_reject_total{reason}` (RBAC / rate limit / audit fail)
- `admin_2fa_challenge_total`
- `admin_2fa_fail_total`

## 7. Bekleyen Kararlar
- SSO provider final seçim (Google vs Okta) — ops team.
- Retool vs custom React tool — Phase 02 karar.
- Admin activity anomaly detection (Phase 02).
- Screen recording / session replay (Phase 03 forensic).
