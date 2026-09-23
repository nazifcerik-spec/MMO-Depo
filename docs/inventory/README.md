# Inventory + Equipment (Faz 11)

Karakterin item ownership'ini yönetim katmanı. Item Engine item'ları yaratır (bkz. [item-engine/](../item-engine/README.md)); bu klasör onların **nerede** tutulduğunu, **nasıl** yer değiştirdiğini, **kim erişebildiğini** modeller.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Inventory](inventory.md) | Karakter üstü bag — slot bazlı |
| [Bank](bank.md) | Şehir bank storage + tab sistemi |
| [Storage](storage.md) | Kavramsal genel çerçeve (mailbox, guild bank, loadout) |
| [Equipment](equipment.md) | Ekipman slot + equip/unequip flow |
| [Slots](slots.md) | Slot mekaniği + capacity + expansion |
| [Stacking](stacking.md) | Stack rules + split + merge |
| [Weight](weight.md) | Weight sistemi (v1 kapalı; design hazır) |
| [Item Transfer](item-transfer.md) | Cross-owner transfer flow + atomicity |
| [Item Lock](item-lock.md) | Accidental sell/destroy/trade koruma |

## 2. Storage Türleri Genel Bakış

| Kaynak | Sahip | Persistence | Cap (v1) | Erişim |
|---|---|---|---|---|
| **Character Inventory** | Char | `character.inventory_slot` | 40..120 (bag upgrade ile) | Her yerde |
| **Character Bank** | Char | `character.bank_slot` + `bank_tab` | Tab başına 40, N tab | Sadece şehir bankasında |
| **Character Loadout** (Phase 02) | Char | `character.loadout` | Sınırsız preset | Şehirde swap |
| **Guild Bank** | Guild | `social.guild_bank_tab` + `guild_bank_slot` | Tab başına 100, N tab | Guild hall'da; rank permission |
| **Mailbox** | Char | `economy.mail` + attachments | 100 mail (server cap) | Herhangi bir yer |
| **Marketplace Listing** | Char (soft) | `economy.market_listing` | Aktif listing 25 max | Her yerde |
| **Void (Deleted)** | Sistem | `content.item_instance` owner_kind='void' | – | – (arşiv) |

## 3. Ownership Katmanları

Tüm storage'lar **Item Engine** üstünde çalışır (bkz. [item-engine/item-instance.md §3](../item-engine/item-instance.md)):
- `content.item_instance.owner_kind` (character_inventory / character_bank / character_equipment / guild_bank / mailbox / market_listing / void).
- `owner_id` polymorphic (character.id, guild.id, mail.id, listing.id).

Application layer sorumlu: owner_kind ile owner_id tutarlılık.

## 4. Cross-Storage Kural Seti

1. **Item bir anda tek storage'da**: aynı instance başka storage'a taşındığında transfer atomic.
2. **Bound rules geçerli**: BoP item guild_bank / mail'e geçemez.
3. **Storage transfer atomic**: source remove + destination add tek transaction.
4. **Audit log her transfer için**: `economy.item_transfer_log` (bkz. [database/economy-market.md §2](../database/economy-market.md)).
5. **Slot dolu = activity paused**: AFK aktivite envanter dolduğunda durur (bkz. [systems/afk-engine.md](../systems/afk-engine.md)).

## 5. Ortak Endpoint Listesi

```
GET    /v1/characters/{id}/inventory
POST   /v1/characters/{id}/inventory/move          # slot içi hareket
POST   /v1/characters/{id}/inventory/split         # stack split
POST   /v1/characters/{id}/inventory/merge         # stack merge
POST   /v1/characters/{id}/inventory/salvage       # multi item salvage
POST   /v1/characters/{id}/inventory/destroy       # confirm delete
POST   /v1/characters/{id}/inventory/lock          # lock/unlock
POST   /v1/characters/{id}/inventory/vendor-sell   # bulk sell

GET    /v1/characters/{id}/bank
POST   /v1/characters/{id}/bank/deposit
POST   /v1/characters/{id}/bank/withdraw
POST   /v1/characters/{id}/bank/purchase-tab
POST   /v1/characters/{id}/bank/rename-tab

GET    /v1/characters/{id}/equipment
POST   /v1/characters/{id}/equipment                # equip
DELETE /v1/characters/{id}/equipment/{slot}         # unequip

GET    /v1/characters/{id}/mail
POST   /v1/characters/{id}/mail                     # send
POST   /v1/characters/{id}/mail/{id}/claim          # attachments
DELETE /v1/characters/{id}/mail/{id}
```

Guild bank: [social/guilds.md](../game-design/guilds.md) tarafı.

## 6. Ortak İlkeler

1. **Server authoritative** — client sadece intent gönderir.
2. **Optimistic locking** — `item_instance.version`, `character.character.version`.
3. **Transaction safety** — SERIALIZABLE isolation cross-storage transfer için.
4. **Idempotent operations** — `Idempotency-Key` header + application-level dedup.
5. **Rate limit** — bulk operations (salvage, sell) 60/dk cap.
6. **Undo penceresi yok** — destroy/salvage geri alınamaz; kullanıcı çift onay verir.

## 7. Cross-Reference
- [game-design/inventory.md](../game-design/inventory.md), [game-design/equipment.md](../game-design/equipment.md)
- [database/character.md](../database/character.md) (slot tabloları)
- [database/economy-market.md](../database/economy-market.md) (transfer log, mailbox)
- [database/social-guilds-parties.md](../database/social-guilds-parties.md) (guild bank)
- [item-engine/item-instance.md](../item-engine/item-instance.md) (ownership transfer)
- [systems/afk-engine.md](../systems/afk-engine.md) (envanter dolu → pause)

## 8. Bekleyen / Yönlendirilecek Kararlar

- Weight sistemi v1 kapalı (bkz. [weight.md](weight.md)); Phase 02 opsiyonu.
- Loadout preset (equipment + skill bar snapshot) — Phase 02.
- Auto-salvage rule engine — Phase 02.
- Wishlist / smart search — Phase 02.
