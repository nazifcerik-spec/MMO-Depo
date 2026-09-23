# Players (Account View)

## 1. Amaç
Support ekibinin oyuncu hesabı üzerinde teşhis ve moderasyon yapması.

## 2. Sayfalar

### List View
- Search: email, username, account_id, character name → account.
- Filter: state (active / locked / banned), created_at range, last_login range.
- Sort: created, last_login, character count.
- Columns: email (masked default) · username · state · character count · last_login · account_id.

### Detail View
```
Account: acct_e3d…4a
Email: g***@example.com     [🔓 Reveal (audit)]
Username: Grim
Created: 2026-04-12
State: active
2FA: enabled (TOTP)
Devices: 3 active
Characters: 2 (see below)

[Mute 24h] [Ban 7d] [Lock] [Reset Password Link] [Delete Account (grace)]

Recent Sessions (10)
Recent Login Failures
Payment History (Phase 02)
```

## 3. PII Reveal

- Email default masked (`g***@example.com`).
- `player.pii_reveal` permission gerektir; her reveal ayrı audit entry.
- Full email display 60 sn sonra otomatik re-mask.

## 4. Aksiyonlar

| Action | Permission | Step-up | Notes |
|---|---|---|---|
| `player.mute<=24h` | support_l1+ | no | reason zorunlu |
| `player.ban<=7d` | support_l2+ | yes | reason + expires_at |
| `player.ban_permanent` | admin | yes | dual approval Phase 02 |
| `player.unlock` | support_l2+ | no | login attempt sayaç reset |
| `player.password_reset_link` | support_l1+ | no | email delivery |
| `player.device_revoke` | support_l1+ | no | tek cihaz push token invalidate |
| `player.delete_account` | admin | yes | 30 gün grace; compliance |

Her aksiyon audit `audit.admin_action` + player'a in-game notification (ban dışında).

## 5. Model Refs

`auth.account`, `auth.session`, `auth.mute_state`, `auth.ban_state` — bkz. [auth](../auth/README.md).

Admin panel `GET /internal/admin/player/{account_id}` compound endpoint.

## 6. Karakter Listesi Embed

Detail sayfasında account'un karakterlerinin özet listesi + karaktere navigate.

## 7. Bulk Aksiyonlar (Phase 02)

RMT bot dalgası vs — CSV import → bulk ban (her satır ayrı audit).

## 8. Metrikler
- `admin_player_view_total`
- `admin_player_pii_reveal_total`
- `admin_player_mute_total`, `admin_player_ban_total`

## 9. Endpoint

- `GET /internal/admin/player?q=&state=&cursor=`
- `GET /internal/admin/player/{account_id}` (compound)
- `POST /internal/admin/player/{account_id}/mute` `{ hours, reason }`
- `POST /internal/admin/player/{account_id}/ban` `{ hours or 'permanent', reason }`
- `POST /internal/admin/player/{account_id}/unlock`
- `POST /internal/admin/player/{account_id}/pii-reveal-audit`
- `POST /internal/admin/player/{account_id}/device/{device_id}/revoke`

## 10. Cross-Reference
- [auth](../auth/README.md), [auth-rbac.md](auth-rbac.md)
- [characters.md](characters.md)

## 11. Bekleyen Kararlar
- Payment history entegrasyon Phase 02.
- Bulk aksiyon UI + CSV Phase 02.
- Dual approval permanent ban Phase 02.
