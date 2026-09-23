# Item Instance

## 1. Amaç
Sahiplenilmiş bir item'ın kalıcı satırı: rolled stats, affixes, sockets, enchant, durability, owner. Template'in "gerçek dünya kopyası".

Kaynak: [database/content-templates.md §9](../database/content-templates.md).

## 2. Model
```
item_instance
  id (uuid)
  template_id (fk)
  owner_kind (character_inventory/character_bank/character_equipment/guild_bank/mailbox/market_listing/void)
  owner_id (polymorphic — application enforce)
  rolled_stats (json)
  affixes (json: {prefix_id, suffix_id})
  sockets (json: [{slot, gem_instance_id}])
  enchant_id (nullable)
  durability_current (nullable)
  bound_to_character_id (nullable)
  created_at
  seed (bytea — RNG replay)
  template_snapshot_hash (text)
  version (int — optimistic lock)
```

## 3. Ownership Polimorfizmi

`owner_kind` + `owner_id` tandem — Postgres FK enforce edemez (polymorphic).

Application layer sorumluluğu:
- **character_inventory**: bind to `character.character.id`; ilgili `character.inventory_slot` satırı var.
- **character_bank**: `character_id` + tab/slot bilgisi bank_slot'ta.
- **character_equipment**: `character_id`; equipment slot çakışması yok.
- **guild_bank**: `guild.id`.
- **mailbox**: `mail.id`.
- **market_listing**: `market_listing.id`.
- **void**: silinmiş / arşivlenmiş referans.

Transfer akışı bir source → destination atomic:
```go
func (e *ItemEngine) Transfer(ctx context.Context, instanceID uuid.UUID, to OwnerRef, reason string) error {
    return e.pool.WithSerializableTx(ctx, func(tx db.TX) error {
        inst, err := e.repo.WithTx(tx).LoadForUpdate(ctx, instanceID)
        if err != nil { return err }
        // 1. Remove from current owner slot
        if err := e.removeFromOwner(ctx, tx, inst); err != nil { return err }
        // 2. Update ownership fields
        inst.OwnerKind = to.Kind
        inst.OwnerID = to.ID
        inst.Version++
        if err := e.repo.WithTx(tx).Update(ctx, inst); err != nil { return err }
        // 3. Insert into new owner slot
        if err := e.addToOwner(ctx, tx, inst, to); err != nil { return err }
        // 4. Log transfer
        _ = e.repo.WithTx(tx).LogTransfer(ctx, inst, reason)
        return nil
    })
}
```

## 4. Instance Yaratımı (Create)

Girdi:
```go
type CreateRequest struct {
    TemplateID  uuid.UUID
    Owner       OwnerRef
    Seed        Seed          // deterministik
    RarityOverride *Rarity    // opsiyonel; loot table meta ise
    ForceBind   bool           // BoP quest reward
    Source      CreateSource   // loot/craft/quest/admin/event
}
```

Adımlar:
1. Load template.
2. Rarity determine (override ya da template.rarity).
3. Roll base stats (min-max range → seed_base).
4. Determine affix count (rarity budget).
5. Roll prefix (if applicable) — pool + seed_prefix.
6. Roll suffix — pool + seed_suffix.
7. Determine socket count.
8. Compute template_snapshot_hash.
9. Set durability_current = template.base_durability.
10. Insert row → `content.item_instance`.
11. Insert `economy.item_transfer_log` (source_kind='loot' vs).
12. Emit event `item.instance_created`.

Idempotent: `idempotency_key = hash(seed + template_id + owner_ref)`; retry aynı sonuç.

## 5. Instance Lifecycle

```
CREATED (owner set at creation)
  ↔ TRANSFERRED (owner change; log)
  ↓
BOUND (equipped or BoP triggered)
  ↓
VOIDED (destroyed, hard-deleted owner cascade, sink)
```

**Void** = owner_kind='void' + owner_id=NULL. Instance satırı arşiv için tutulur; garbage collect worker (Phase 02).

## 6. Binding Trigger

- **BoP (bind on pickup)**: template.binding_type='bop' → CreateRequest sırasında `bound_to_character_id = new_owner_character_id`.
- **BoE (bind on equip)**: equip event → `bound_to_character_id = character_id`.
- **BoA (bind on account)** (Phase 02): benzer, account_id.
- **no_bind**: hiç bind yok; free trade.

Bound olduktan sonra market'e/mail'e transfer engellenir:
```go
if inst.BoundToCharacterID != nil && to.Kind != CharacterOwned(inst.BoundToCharacterID) {
    return errs.ErrForbidden  // item.bound_to_other
}
```

## 7. Stack Rules

`stack_count > 1` template için (materials, consumables):
- Aynı template + aynı önemli flag (rarity=common yes; rolled_stats yok) → stack merge mümkün.
- Instance-only item'lar (weapon, armor with rolls) stack yok — her instance ayrı row.

Not: schema'da `stack_count` field'ı `character.inventory_slot` üstünde (bkz. [database/character.md §3](../database/character.md)).

Merge:
- Yeni loot geldiğinde, aynı template ile inventory'de mevcut stack varsa merge; instance yeni yaratılmaz (materials için).
- Alternatif: her materyal ayrı instance; `stack_count=1`. Karar: material'lar için **shared instance + stack_count** (basit).

## 8. Durability

- `durability_current` her hit / cast'te azalır (combat trigger).
- 0 olunca: item passive/proc devre dışı; base stat kalır (v1 tercih).
- Repair NPC (city) veya self-repair kit consumable.

Detay: item durability lifecycle [database/character.md §17](../database/character.md) death loss + this repair.

## 9. Endpoint

### 9.1 Item Instance Detail

**GET `/v1/items/{instance_id}`**
Auth required; ownership check.

Response:
```json
{
  "id": "ii-...",
  "template": { "code": "steel_longsword", "name_key": "..." },
  "rarity": "rare",
  "rolled_stats": { "weapon_damage": 53, "attack_speed_ms": 2400 },
  "affixes": {
    "prefix": { "code": "steel", "name_key": "affix.steel", "stats": { "str": 5 } },
    "suffix": { "code": "of_the_bear", "name_key": "affix.of_the_bear", "stats": { "sta": 10 } }
  },
  "sockets": [
    { "slot": 0, "gem": { "instance_id": "gi-...", "template_code": "ruby_cut", "stats": { "str": 3 } } }
  ],
  "enchant": { "code": "sharpness_ii", "stats": { "attack_power": 12 } },
  "durability": { "current": 87, "max": 100 },
  "binding": "boe",
  "bound_to_character_id": null,
  "owner_kind": "character_inventory",
  "created_at": "..."
}
```

### 9.2 Batch (Inventory)

**GET `/v1/characters/{id}/inventory`** — batch item summaries.

## 10. Snapshot Discipline

Rolled stats + affixes + sockets + enchant + durability = **instance'ın "hafıza"sıdır**; template değişse bile korunur.

Template değişimi (patch) → instance'ın effective stats bazlarını etkileyebilir (weapon_damage min-max değişimi mevcut instance'a yansır mı?):
- **Karar (v1)**: rolled_stats snapshot; template değişimi mevcut instance'a yansımaz.
- **Alternatif (Phase 02)**: content patch aksiyonu (rescale existing instances) worker ile toplu.

## 11. Fail Modları
- Ownership violation → 403.
- Template deprecated + instance oluşturma → reject (rare; audit).
- Snapshot hash mismatch → log warn; instance yine çalışır.
- Concurrent transfer → optimistic lock retry.

## 12. Test Stratejisi
- Create → detail happy path.
- Transfer between owners → transfer_log + owner update.
- BoP binding on create.
- BoE binding on equip.
- Stack merge (material).
- Deterministic seed → same rolled_stats + affixes.

## 13. Metrics
- `item_instance_created_total{template_code,rarity}`.
- `item_instance_transferred_total{reason}`.
- `item_instance_bound_total{binding_type}`.
- `item_instance_voided_total{cause}`.

## 14. Audit
- `item.instance_created` (source=loot/craft/quest/admin).
- `item.transferred` (source→target).
- `item.destroyed` (void with reason).

## 15. Cross-References
- [Item Template](item-template.md), [Rarity](rarity.md), [Affix](affix.md), [Sockets](sockets.md), [Gems](gems.md), [Enchantments](enchantments.md), [Set](set.md), [Unique Effect](unique-effect.md)
- [database/content-templates.md §9](../database/content-templates.md)
- [database/economy-market.md §2](../database/economy-market.md) — item_transfer_log
- [game-design/inventory.md](../game-design/inventory.md), [game-design/equipment.md](../game-design/equipment.md)

## 16. Bekleyen Kararlar
- Void garbage collect policy (retention süresi).
- Cross-account BoA implementation Phase 02.
- Rescale worker Phase 02 opsiyonu.
- Instance rename (kişisel isim koyma) Phase 03.
