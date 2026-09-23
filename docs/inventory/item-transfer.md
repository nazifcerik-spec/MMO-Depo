# Item Transfer

## 1. Amaç
Item'ın owner değişimi (inventory↔bank, inventory↔equipment, character↔mail, character↔market, character↔guild bank) tek ortak sözleşme ile atomik + auditable.

## 2. Ana Sözleşme

Bkz. [item-engine/item-instance.md §3](../item-engine/item-instance.md):
```go
Transfer(ctx, itemInstanceID, to OwnerRef, reason string) error
```

Transaction akışı:
```
TX BEGIN;
1. Load item_instance FOR UPDATE.
2. Validate source owner matches.
3. Remove from source storage (slot delete or move to "in-transit").
4. Update item_instance.owner_kind + owner_id.
5. Insert into destination storage (allocate slot / mailbox / listing).
6. Log economy.item_transfer_log.
7. Trigger side effects (BoE → BoP bind on equip; cache invalidate; event emit).
TX COMMIT;
```

Isolation: **SERIALIZABLE** (envanter dolma race + double-transfer önlenir).

## 3. Transfer Kaynakları + Nedenleri

`reason` sabitleri:
| Reason | Detay |
|---|---|
| `loot` | mob/boss/node/chest drop |
| `craft_output` | crafting recipe output |
| `quest_reward` | quest turn-in |
| `admin_grant` | GM manual grant |
| `event_reward` | event participation |
| `mailbox_send` | player-to-player mail |
| `mailbox_claim` | mail attachments alındı |
| `market_list` | marketplace'e list |
| `market_cancel` | listing cancel |
| `market_sale` | satın alındı |
| `guild_deposit` | guild bank deposit |
| `guild_withdraw` | guild bank withdraw |
| `bank_deposit` | character bank |
| `bank_withdraw` | character bank |
| `equip` | inventory → equipment |
| `unequip` | equipment → inventory |
| `salvage_consume` | salvage input |
| `vendor_sale` | NPC vendor'a sat (item destroyed) |
| `destroyed` | destroy endpoint |

## 4. Transfer Türleri

### 4.1 Internal (Same Character)
- inventory ↔ bank.
- inventory ↔ equipment.
- inventory slot ↔ inventory slot (move).

Atomic — SERIALIZABLE.

### 4.2 Character → Mailbox → Character
- Send: source inv → mail (owner_kind=mailbox).
- Claim: mail → target inv (owner_kind=character_inventory).

Sender + Receiver ayrı tx (mail create ve claim iki adım).

### 4.3 Character → Market → Buyer Mail
- Listing create: inv → market_listing.
- Buy: market_listing → buyer mailbox (item), seller mailbox (gold).

Atomicity kritik — bkz. [database/economy-market.md §9](../database/economy-market.md).

### 4.4 Character → Guild Bank → Character
- Deposit: inv → guild_bank.
- Withdraw: guild_bank → inv.

Guild bank access rules (rank permission).

### 4.5 Character → Void
- Vendor sell, destroy, expiration.
- Item instance owner_kind='void', owner_id=NULL.
- Instance satırı korunur (audit); Phase 02 garbage collect worker.

## 5. Endpoint (Genel)

Transfer endpoint'leri storage-özel (bkz. inventory.md, bank.md, mail, market docs). Ortak internal wrapper:
```go
func (e *ItemEngine) Transfer(ctx, req TransferRequest) error
```

## 6. Idempotency

Transfer idempotent:
- `Idempotency-Key` header (client autogen UUID).
- Application layer'da:
  ```
  transfer_ledger_check(idempotency_key) → varsa no-op success.
  yoksa: transfer + log ile insert.
  ```

## 7. Bound Item Kısıtları

- **BoP** (bind on pickup): kimseye trade edilemez; guild bank ret; mail ret; market ret.
  - İstisna: **character owned within same account** (BoA sonrası Phase 02 aynı account karakterler arası).
- **BoE** (bind on equip): equip anına kadar tradeable; sonrasında BoP kuralı.
- **BoA** (bind on account): aynı account karakterleri arasında tradeable; farklı account ret.
- **no_bind**: her yere.

Transfer validation:
```go
if inst.BoundToCharacterID != nil {
    if to.Kind == KindGuildBank { return errs.ErrForbidden }
    if to.Kind == KindMailbox && to.OwnerCharID != inst.BoundToCharacterID { return errs.ErrForbidden }
    if to.Kind == KindMarketListing { return errs.ErrForbidden }
    // aynı character için içeride move OK
}
```

## 8. Locked Item

Bkz. [item-lock.md](item-lock.md).

- Locked item transfer edilemez (equip, sell, destroy, deposit, list).
- Lock override yok; kullanıcı önce unlock.

## 9. Cross-Storage Atomicity Guarantee

Postgres tek DB — SERIALIZABLE isolation yeterli. Redis (cooldown, live cache) etkilenmez.

Failure mode (partial commit yok — tx tam ya rollback):
- App crash mid-tx → rollback.
- Postgres crash → WAL recovery.
- Redis unreachable → transfer OK (Redis cache invalidate best-effort).

## 10. Bulk Transfer

Bulk operations (deposit-all, salvage multiple):
- Application layer transaction'da N transfer.
- Any failure → full rollback.
- Optimization: batch item_transfer_log insert (single INSERT with multiple rows).

## 11. Cross-Character Direct Trade

**Yok v1** (DD-023). Ticaret sadece marketplace + mail.

Phase 02 opsiyonu: escrow direct trade (dual confirmation).

## 12. Anti-Duplication Guards

- Item instance version optimistic lock.
- Owner check FOR UPDATE.
- Insert/delete atomic pair.
- Ledger append-only audit.

Reconcile worker (daily):
```
For each item_instance:
  - Verify owner_kind + owner_id has corresponding slot/mail/listing row.
  - Mismatch → audit critical, admin review.
```

## 13. Metrics
- `item_transfer_total{reason,from,to}`.
- `item_transfer_duration_seconds`.
- `item_transfer_rejected_total{reason}`.
- `item_transfer_atomicity_violation_total` (should be 0 always).

## 14. Audit

`economy.item_transfer_log` her transfer için satır:
- item_instance_id, from_owner_kind/id, to_owner_kind/id, reason, correlation_id, meta_json.

## 15. Test Stratejisi

- Happy path all transfer types.
- Bound item reject.
- Locked item reject.
- Concurrent transfer (goroutine 100) → biri success.
- Bulk deposit rollback.
- Idempotent retry.

## 16. Cross-References
- [Inventory](inventory.md), [Bank](bank.md), [Equipment](equipment.md), [Storage](storage.md), [Item Lock](item-lock.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)
- [database/economy-market.md §2](../database/economy-market.md) — item_transfer_log
- [systems/reward-system.md](../systems/reward-system.md)

## 17. Bekleyen Kararlar
- Direct trade (escrow) Phase 02.
- Cross-account transfer (BoA + linked accounts) Phase 03.
- Reconcile worker frequency (daily vs hourly).
- Transfer batch API (Phase 02 optimization).
