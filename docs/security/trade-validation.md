# Trade Validation

## 1. Amaç
P2P trade session güvenliği — bait-and-switch, çift transfer, race, block bypass, bind rule bypass önleme.

## 2. Session Life-Cycle Guard'lar

State machine (bkz. [marketplace/trade.md](../marketplace/trade.md)):

- **Update → auto-unlock** (DD-138): iki taraf lock/confirm reset.
- **Confirm hash**: `seen_partner_offer_hash` mismatch → `E_TRADE_OFFER_CHANGED`.
- **Version increment**: her offer update `version++`; concurrent update stale reject.
- **Cancel cooldown**: 30 sn same-pair spam guard.
- **Rate limit**: karakter başına 3 trade session açılışı/dk.

## 3. Finalize TX Check (SERIALIZABLE)

Finalize öncesi rechecks (state'in son 100ms'de değişme ihtimali):
1. Item ownership her attachment için hâlâ self.
2. Item lock — sadece bu trade session lock'u (başka lock varsa cancel).
3. Item bind rule (BoP/soulbound reject; double-check).
4. Wallet gold sufficient her iki taraf için.
5. Envanter kapasite karşı tarafta yeterli (yeni item alacak).
6. Session state == both_confirmed.
7. Both sides' character not frozen.

Fail → session state='cancelled', reason kaydı; UI'da her iki taraf bilgi.

## 4. Anti-Fraud Kuralları

### Bait-and-Switch
Kullanıcı offer'ı değiştirir → **iki taraf** confirm/lock sıfırlanır. Karşı taraf yeniden bakmadan onay veremez. Confirm request'in `seen_partner_offer_hash` payload'unda; server offer state hash'i ile eşleşmelidir.

### Item Swap Attempt
Trade lock sonrası satıcı item unequip / envanterde başka yere taşımaya çalışır → item lock (`lock_owner='trade', locked_ref_id=session_id`) engel; başka mutation `E_ITEM_LOCKED`.

### Same-Account
V1'de aynı account'un iki karakteri arasında trade **serbest** (self-fund, alt karakter transfer için). Phase 02'de anti-farm review (DD-138 alt karar).

### Block Bypass
`social.block_list` kontrol iki yönlü. Herhangi biri diğerini block etti ise `E_TARGET_BLOCKED`.

## 5. Concurrent Trade Session Limit

- Karakter başına aynı anda max 1 aktif trade session.
- Yeni request request pending: mevcut aktif session finalize/cancel edilene kadar reject.

## 6. Marketplace Trade Ayrı (DD-146)

Marketplace listing için de aynı prensipler:
- Same-account buy block (buy zamanında).
- Price race guard.
- Item escrow lock (post → sold arası).

Bkz. [marketplace](../marketplace/README.md).

## 7. Fraud Signal

- **High-value trade rush**: yeni account 24h içinde >1M gold trade → RMT flag.
- **Alt-farm pattern**: aynı account'un karakterleri arası gold+item trade yoğunluk (v1 legal ama flag) Phase 02 review.
- **Bait-and-switch attempt**: offer değişiminden sonra karşı taraf otomatik unlock; peş peşe 3+ attempt → session force cancel + `sec_trade_bait_attempt_total`.

## 8. Metrikler
- `sec_trade_bait_attempt_total`
- `sec_trade_offer_hash_mismatch_total`
- `sec_trade_finalize_recheck_fail_total{reason}`
- `sec_trade_block_reject_total`
- `sec_trade_high_value_flag_total`

## 9. Cross-Reference
- [marketplace/trade.md](../marketplace/trade.md), [marketplace](../marketplace/README.md)
- [social/block.md](../social/block.md), [inventory/item-lock.md](../inventory/item-lock.md)

## 10. Bekleyen Kararlar
- Same-account trade anti-farm review Phase 02.
- RMT-value threshold tune Phase 02.
- Trade replay (audit forensic) Phase 02.
- P2P COD (cash-on-delivery mail) trade alternatif Phase 02.
