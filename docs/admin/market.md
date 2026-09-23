# Market (Admin View)

## 1. Amaç
Aktif listing/auction gözlem, RMT şüphe listeleme cancel, price anomaly tespit, sale history forensic.

## 2. List View
- Filter: item_template, rarity, price range, seller, state (active/sold/expired/cancelled), listed_at range.
- Columns: listing_id · item · rarity · unit_price · qty · seller_char · state · listed_at · expires_at.

## 3. Suspicious Pattern Detection (Phase 02)

- **RMT flag**: same-account listing → sold → buyer new account 24h old.
- **Price outlier**: 10× median üstü/altı listing.
- **Bot pattern**: aynı seller 24h içinde 50+ listing.
- **Wash trade**: aynı kişi buy + sell (marketplace same-account block var v1 — DD-146; buyer alt char pattern hâlâ tespit edilebilir).

## 4. Detail View
```
Listing: listing_ab…
Item: sword_iron uncommon L15
Seller: Grim (acct_e3d…)
Price: 4500c
Duration: 24h
State: active (expires 2026-09-10 12:00)

Listing snapshot: {…full JSON}
Escrow item instance: item_ab…

Related audit / actions
```

## 5. Aksiyonlar

| Action | Permission | Step-up |
|---|---|---|
| `market.listing_cancel` | support_l2+ | no | reason zorunlu; fee iade default false, admin override |
| `market.listing_price_report` | viewer+ | no | anomaly review flag |
| `market.rmt_seller_ban_bundle` | admin | yes | seller ban + all listings cancel |

Cancel akış: bkz. [marketplace/cancel.md](../marketplace/cancel.md) §5 admin cancel.

## 6. Sale History Forensic

Query builder: item_template, seller, buyer, price range, sold_at range → CSV export.

## 7. Endpoint

- `GET /internal/admin/market/listing?…`
- `GET /internal/admin/market/listing/{id}`
- `POST /internal/admin/market/listing/{id}/cancel` `{ reason, refund_fee_override? }`
- `GET /internal/admin/market/history?…`
- `POST /internal/admin/market/seller/{account_id}/ban-and-cancel` (bundle; admin)

## 8. Cross-Reference
- [marketplace](../marketplace/README.md), [marketplace/cancel.md](../marketplace/cancel.md)
- [economy/transaction.md](../economy/transaction.md)

## 9. Bekleyen Kararlar
- RMT detection ML model Phase 03.
- Price anomaly alert threshold Phase 02.
- Bulk cancel bundle UI + audit Phase 02.
