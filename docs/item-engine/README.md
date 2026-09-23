# Item Engine (Faz 10)

Loot doyumunun mekanik motoru: template tanımı, instance yaratımı, rarity roll, affix pool, socket/gem/enchant, set bonus, unique effect. Combat + economy + inventory bu motoru tüketir.

Yüksek seviye vizyon: [game-design/items.md](../game-design/items.md), [game-design/equipment.md](../game-design/equipment.md). Bu klasör Phase 10 implementation-hazır detay.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış + engine flow (bu doküman) | [README.md](README.md) |
| [Item Template](item-template.md) | Data-driven template + versioning + snapshot |
| [Item Instance](item-instance.md) | Instance lifecycle + ownership polymorphism |
| [Rarity](rarity.md) | Rarity bands + drop distribution + budget |
| [Affix](affix.md) | Affix system architecture + pool resolution |
| [Prefix](prefix.md) | Prefix pool + naming + roll |
| [Suffix](suffix.md) | Suffix pool + naming + roll |
| [Stats](stats.md) | Item base + rolled stats + stat modifier bağlantısı |
| [Sockets](sockets.md) | Socket slot types + install/remove flow |
| [Gems](gems.md) | Gem templates + gem instances + cutting |
| [Enchantments](enchantments.md) | Enchant apply + upgrade path + destruction risk |
| [Set](set.md) | Set bonus tiers (2/4/6) + activation |
| [Unique Effect](unique-effect.md) | Handler pattern (legendary+) |

## 2. Engine Katmanları

```
┌────────────────────────────────────────────────────────┐
│               Domain Entry Points                       │
│  (loot.roll, craft.output, quest.reward, admin.grant)   │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────────────────┐
│           ItemEngine (internal/item)                    │
│                                                          │
│  Template registry (in-memory + content_version)        │
│  Instance factory (seed'li roll)                        │
│  Affix resolver (prefix + suffix + roll)                │
│  Socket & gem manager                                    │
│  Enchant manager                                         │
│  Set bonus evaluator                                     │
│  Unique effect handler dispatcher                        │
│  Durability manager                                      │
│  Stat contribution (Modifier stack'e katkı)             │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────────────────┐
│  Data: content.item_template, item_instance,           │
│        affix_pool, affix, item_set, enchant, gem       │
└────────────────────────────────────────────────────────┘
```

## 3. Ortak İlkeler

1. **Template ↔ Instance** ayrımı (DD-035): kalıcı katalog + oyuncu-sahip snapshot.
2. **Data-driven** — item, affix, set, enchant, gem hepsi content YAML.
3. **Deterministic roll** — seed'li; replay/audit için tekrarlanabilir.
4. **Snapshot resistance** — template değişimi mevcut instance'ı korur (`template_snapshot_hash`).
5. **Modifier bağlantısı** — item stat'ları [skill-engine/modifiers.md](../skill-engine/modifiers.md) modifier stack'ine katkı olarak eklenir.
6. **Server authoritative** — item state (owner, stats, durability, sockets, enchant) hep server.

## 4. Ana Interface

```go
package item

type Engine interface {
    // Registry
    Template(ctx context.Context, id uuid.UUID) (TemplateDef, error)
    TemplateByCode(ctx context.Context, code string) (TemplateDef, error)

    // Instance factory (loot roll, craft output, quest reward)
    Create(ctx context.Context, req CreateRequest) (Instance, error)

    // Ownership transfer (transaction-safe)
    Transfer(ctx context.Context, instanceID uuid.UUID, to OwnerRef, reason string) error

    // Affix reroll (rare consumable)
    RerollAffixes(ctx context.Context, instanceID uuid.UUID, seed Seed) (Instance, error)

    // Sockets & Gems
    InstallGem(ctx context.Context, instanceID uuid.UUID, gemInstanceID uuid.UUID, slot int) error
    RemoveGem(ctx context.Context, instanceID uuid.UUID, slot int) error

    // Enchant
    ApplyEnchant(ctx context.Context, instanceID uuid.UUID, enchantCode string) (EnchantResult, error)
    RemoveEnchant(ctx context.Context, instanceID uuid.UUID) error

    // Durability
    LoseDurability(ctx context.Context, instanceID uuid.UUID, amount int32) error
    Repair(ctx context.Context, instanceID uuid.UUID, mode RepairMode) (RepairCost, error)

    // Stat contribution
    ModifiersFor(ctx context.Context, instance Instance) []engine.Modifier
    SetBonusesFor(ctx context.Context, equipped []Instance) []engine.Modifier
}
```

## 5. Data Flow — Örnek: Mob Kill Loot Roll

```
1. Combat → mob killed → loot.RollForKill(seed, mob, killer)
2. Loot System → picks entries from mob.loot_table (weighted)
3. For each entry:
     item.Create(template, seed_for_this_drop)
       - template.load
       - roll rarity band (from meta)
       - roll base stats (min..max range)
       - roll affix count (rarity budget)
       - roll prefix + suffix (from affix pool)
       - roll socket count (rarity-driven)
       - assign owner = character_inventory
       - insert item_instance row
     add to reward payload
4. Character receives → notification "New Item"
```

## 6. Data Flow — Örnek: Equip

```
1. Client POST /v1/characters/{id}/equipment {slot, item_instance_id}
2. Server:
   a. Verify ownership.
   b. Verify slot compatibility (template.equipment_slot).
   c. Verify requirements (level, class, spec, stat_min).
   d. If off-hand + 2H weapon → unequip conflict handling.
   e. Transaction:
      - Move current equipment (if any) → inventory slot.
      - Move new item → equipment.
      - Recompute set bonus evaluation.
      - Invalidate derived stat cache.
   f. Publish 'item.equipped'.
3. Response 200 + updated sheet.
```

## 7. Item Model Overview

Template (kalıcı katalog):
- `code`, `category`, `subcategory`, `rarity`, `base_stats`, `affix_pool_ids`, `socket_slots`, `set_id`, `binding_type`, `requirements`, `unique_effect_id`.

Instance (sahiplenilmiş):
- `template_id`, `owner_kind`, `owner_id`, `rolled_stats`, `affixes`, `sockets` (with gems), `enchant_id`, `durability_current`, `bound_to_character_id`, `seed`, `template_snapshot_hash`.

Detay: [item-template.md](item-template.md), [item-instance.md](item-instance.md).

## 8. Cache Stratejisi

- **Template cache** in-memory (content_version bump ile reload).
- **Affix pool cache** in-memory.
- **Set bonus cache** in-memory.
- **Instance** her seferinde DB'den (loot hot-path'inde template cache hit yeter).

## 9. Deterministic Roll

`Create` çağrısı deterministic:
- Input: template_id, seed, context (drop source).
- Output: aynı seed → aynı instance.
- Replay support: audit'ten seed alıp yeniden roll → aynı sonuç.

Random System (bkz. [systems/random-system.md](../systems/random-system.md)) sub-stream:
```
seed_rarity   = hash(seed, "rarity")
seed_base     = hash(seed, "base_stats")
seed_prefix   = hash(seed, "prefix")
seed_suffix   = hash(seed, "suffix")
seed_socket   = hash(seed, "socket_count")
```

## 10. Cross-Reference

- [game-design/items.md](../game-design/items.md), [game-design/equipment.md](../game-design/equipment.md), [game-design/inventory.md](../game-design/inventory.md)
- [database/content-templates.md §8-9](../database/content-templates.md)
- [systems/loot-system.md](../systems/loot-system.md)
- [systems/random-system.md](../systems/random-system.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md) — item stat modifier
- [skill-engine/stats.md](../skill-engine/stats.md) — cache invalidation
- [character/attributes.md](../character/attributes.md), [character/stats.md](../character/stats.md)

## 11. Bekleyen Kararlar
- Craft output vs loot roll aynı `Create` API'sini kullanmalı — evet (öneri).
- Reroll consumable / vendor Phase 02.
- Item transmog (görsel override) Phase 02 alan planı.
- Cross-server item exchange (Phase 03) — pooling stratejisi.
