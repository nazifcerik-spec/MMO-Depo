# Currency Validation

## 1. Amaç
Wallet + currency balance mutasyonlarının negatif olmaması, cap aşmaması, doğru currency kullanılması.

## 2. Zorunlu Check (Her Ledger Entry Öncesi)

| Check | Amaç |
|---|---|
| **Currency exists** | `economy.currency` içinde tanımlı |
| **Scope uyumu** | `character/account/guild` scope endpoint'e uyar |
| **Tradable** | Trade akışı için `tradable=true` (v1 sadece gold) |
| **Balance sufficiency (out)** | Delta < 0 için `balance + delta >= 0` (sistem hesabı hariç) |
| **Cap not exceeded (in)** | Delta > 0 için `balance + delta <= cap` |
| **Currency code enum** | Whitelist enum, unknown reject |
| **Amount positive** | `abs(delta) > 0` (0-delta entry yazılmaz) |
| **Amount integer** | BIGINT copper (float yok, DD-136) |

## 3. Sistem Hesabı Özel Kuralı

`account_kind='system'` hesaplar (`system.tax_sink, system.fee_sink, system.vendor_sink, system.reward_source, system.loot_source, system.mint_source, system.auction_escrow`):
- **Negatif balance izinli** (source/sink unbounded).
- Kısıtsız (rewrite hesap değil, akış hesap).
- Balance metric olarak tracked (source/sink dashboard).

## 4. Wallet Cap (v1)

- `character_wallet.gold_copper` cap: **10^12 copper** (10 milyar gold; DD-136).
- Currency cap: bkz. [economy/currency.md](../economy/currency.md) tablo.
- Aşım deneyen credit `E_CURRENCY_CAP_EXCEEDED` → grant reddi + inbox notification "cap reached".
- Overflow strategy (Phase 02):
  - Gold: `mail overflow` (extra credit auto-mail item vendor sale).
  - Event token: silent cap (extra kaybolur).
  - Faction rep: daily cap zaten var (DD-167).

## 5. Currency Trade Kısıt (DD-137)

`tradable=false` currency P2P trade/market listing yasak. Endpoint validation:
- Trade offer `gold_copper` accept; başka currency alan yok.
- Market listing `currency_code='gold'` zorunlu.
- Mail attachment gold only (v1); Phase 02'de opsiyonel event_token attach opt.

## 6. Currency Conversion Yasağı (v1)

Gem → gold otomatik dönüşüm **yok** (DD-137). Endpoint yok; ihtiyaç support runbook üzerinden manuel reversal + ayrı ledger entry (rare).

## 7. Wallet Freeze

Karakter freeze durumundayken (admin action, drift alert) tüm currency mutation reddedilir:
- `character.state='frozen'` → wallet servisleri `E_CHARACTER_FROZEN`.
- Aktif contract reveal grantleri de bekletilir (pending queue) — freeze kalkınca resume.

## 8. Fraud Signal

- **Cap-hit rush**: yeni account 24h içinde wallet cap'e yaklaştı → RMT sink candidate flag.
- **Currency scope violation attempt**: character endpoint'te account currency amount → security alert.
- **Sum violation**: currency conservation ihlali → P0 alert.

## 9. Metrikler
- `sec_currency_cap_exceeded_total{currency}`
- `sec_currency_tradable_reject_total`
- `sec_currency_scope_mismatch_total`
- `sec_wallet_freeze_reject_total`

## 10. Cross-Reference
- [economy/gold.md](../economy/gold.md), [economy/currency.md](../economy/currency.md), [economy/transaction.md](../economy/transaction.md)

## 11. Bekleyen Kararlar
- Overflow gold auto-mail Phase 02 opt.
- Currency conversion vendor Phase 02 (event token → gold gibi limited).
- Wallet cap dinamik (level ile artan) Phase 03.
