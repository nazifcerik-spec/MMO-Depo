# Storage (Genel Kavramsal Çerçeve)

## 1. Amaç
Tüm depolama alanlarının ortak soyutlaması. Inventory, Bank, Mailbox, Guild Bank, Marketplace Listing — hepsi Storage. Owner + capacity + access rules ile modellenir.

## 2. Storage Abstraksiyonu

```go
package storage

type Kind string
const (
    KindCharacterInventory Kind = "character_inventory"
    KindCharacterBank      Kind = "character_bank"
    KindCharacterEquipment Kind = "character_equipment"
    KindGuildBank          Kind = "guild_bank"
    KindMailbox            Kind = "mailbox"
    KindMarketListing      Kind = "market_listing"
    KindVoid               Kind = "void"
)

type Ref struct {
    Kind Kind
    ID   uuid.UUID
    TabIndex *int16   // opsiyonel (bank ve guild bank için)
    SlotIndex *int16  // opsiyonel (slotted stores için)
}

type Storage interface {
    Kind() Kind
    Capacity(ctx context.Context, ownerID uuid.UUID) (int, error)
    Used(ctx context.Context, ownerID uuid.UUID) (int, error)
    List(ctx context.Context, ownerID uuid.UUID, filter Filter) ([]StoredItem, error)
    Put(ctx context.Context, tx db.TX, ownerID uuid.UUID, slot int16, itemID uuid.UUID, count int) error
    Take(ctx context.Context, tx db.TX, ownerID uuid.UUID, slot int16, count int) (uuid.UUID, error)
    CanAccess(ctx context.Context, char Character, ownerID uuid.UUID) (bool, string, error)
}
```

Her storage kind ayrı implementation:
- `inventoryStorage`
- `bankStorage`
- `equipmentStorage`
- `guildBankStorage`
- `mailboxStorage`
- `marketListingStorage`

## 3. Storage Locator

Belirli item nerede?
```go
func (s *ItemEngine) LocationOf(ctx context.Context, itemID uuid.UUID) (storage.Ref, error) {
    inst := load(itemID)
    return storage.Ref{
        Kind:      storage.Kind(inst.OwnerKind),
        ID:        inst.OwnerID,
        SlotIndex: inst.SlotIndex,      // varsa
        TabIndex:  inst.TabIndex,        // varsa
    }, nil
}
```

Not: SlotIndex + TabIndex `item_instance` tablosunda direkt tutulmayabilir; storage-specific tablolarda tutulur (bkz. [database/character.md §3-4](../database/character.md), [database/social-guilds-parties.md §5](../database/social-guilds-parties.md)).

## 4. Cross-Storage Transfer

Ortak akış (bkz. [item-transfer.md](item-transfer.md)):
```
TX BEGIN;
1. Source storage: Take(item)
2. Destination storage: Put(item, slot)
3. item_instance.owner_kind + owner_id update
4. item_transfer_log insert
TX COMMIT;
```

Isolation: SERIALIZABLE (envanter dolma race önlenir).

## 5. Access Control

Her storage `CanAccess` implementation:
- **Inventory**: owner character == authenticated character.
- **Bank**: owner character + at bank location.
- **Equipment**: owner character.
- **Guild Bank**: guild member + rank permission for tab.
- **Mailbox**: owner character.
- **Market Listing**: seller character (yalnız cancel), buyer character (satın al).

## 6. Storage-Specific Kısıtlar

| Kind | BoP allowed? | Bind on move? | Max stack | Trade allowed after? |
|---|---|---|---|---|
| Inventory | evet | – | template.stack_max | evet |
| Bank | evet | – | template.stack_max | withdraw sonrası |
| Equipment | evet | BoE → BoP on equip | 1 | ownership değişmez |
| Guild Bank | **hayır** (BoP block) | – | template.stack_max | withdraw eden character |
| Mailbox | evet | – | attachment cap | claim sonrası |
| Market Listing | hayır (BoP değil listable) | – | template.stack_max | satın alanın |

## 7. Overflow Zinciri

Item yerleşirken:
```
Primary: hedef storage slot
  ↓ full
Fallback 1: mailbox (system mail 7 gün)
  ↓ mail cap full
Fallback 2: void (with audit warn)
```

v1 tercih: mailbox overflow (nadir ölçek). Void'a düşme rare + audit.

## 8. Loadout Storage (Phase 02)

Loadout = equipment + skill_bar preset:
```
character.loadout
  id, character_id, name, equipment_json, skill_bar_json, rotation_id
```

Loadout apply: transaction — mevcut equipment/inventory swap ile preset uygular.

v1'de yok; schema hazır.

## 9. Guild Bank Notes

Bkz. [database/social-guilds-parties.md §4-6](../database/social-guilds-parties.md).

- 6 tab max (v1); guild level ile artabilir (Phase 02).
- Tab başına permission (rank-based) — read/write ayrı.
- Deposit log audit-safe (kimin ne koyduğu takip edilir).

## 10. Marketplace Listing Storage

Marketplace listing sırasında item:
- Owner → `market_listing`; slot yok, listing id.
- Satış olursa: mail'e (buyer'a) transfer.
- Iptal: envanter'e geri.

Bkz. [database/economy-market.md §3](../database/economy-market.md).

## 11. Mailbox Storage

- Item attachment cap: 12 per mail (config).
- Mail cap 100 (server).
- Expire 30 gün → sender'a iade.
- Claim = attach'ler envanter'e (kısmi claim mümkün Phase 02).

## 12. Storage Migration (Rare)

Content patch bir storage kind kaldırırsa (nadir):
- Migration job → içerik primary storage (inventory) veya mailbox'a taşı.
- Audit + notification.

## 13. Fail Modları

- Access denied → 403.
- Capacity exceed → 400 `storage.full`.
- Cross-storage inconsistency (ör. item_instance.owner_kind = 'character_inventory' but no inventory_slot row) → data integrity job (worker) tespit + fix.

## 14. Test Stratejisi

- Storage interface implementation per kind.
- Cross-storage transfer round-trip.
- Access control matrix.
- Overflow zinciri.

## 15. Metrics
- `storage_capacity_used{kind}` gauge.
- `storage_transfer_total{from,to,reason}` counter.
- `storage_overflow_total{primary,fallback}` counter.

## 16. Cross-References
- [Inventory](inventory.md), [Bank](bank.md), [Equipment](equipment.md), [Item Transfer](item-transfer.md), [Slots](slots.md)
- [database/character.md](../database/character.md), [database/social-guilds-parties.md](../database/social-guilds-parties.md), [database/economy-market.md](../database/economy-market.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)

## 17. Bekleyen Kararlar
- Loadout storage Phase 02.
- Cross-character bank (BoA) Phase 02.
- Cold storage (arşiv item'lar for inactive chars) Phase 03.
- Storage-integrity worker frequency.
