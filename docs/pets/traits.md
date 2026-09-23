# Traits (Pet Traits)

## 1. Amaç
Pet'e özel passive modifier'lar — species pool'undan random roll (deterministic seed).

## 2. Trait Sistem

Trait = passive modifier:
- Effect DSL benzer skill effect pattern.
- Kalıcı (pet lifetime).
- Rerollable (Phase 02 consumable).

## 3. Trait Kataloğu (Örnek)

| Trait Code | Effect | Species Bias |
|---|---|---|
| `keen_senses` | +5% crit chance | Beast, avian |
| `swift_paw` | +10% movement speed | Beast |
| `iron_hide` | +15% pet HP | Reptile, mystical |
| `pack_leader` | +5% party damage aura | Beast (rare) |
| `blazing_bite` | Auto-attacks fire dmg | Elemental |
| `mistwalk` | +10% pet evasion | Avian, mystical |
| `master_gatherer` | +5% profession yield | Utility species |
| `treasure_finder` | +2% treasure spawn chance | Utility |
| `hunter_instinct` | +5% loot rare drop | Beast |
| `unbreakable_bond` | +10% pet damage; requires happiness > 80 | All |

## 4. Trait Model

```sql
CREATE TABLE content.pet_trait (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  description_key VARCHAR(96) NULL,
  effect_json JSONB NOT NULL,       -- modifier list
  rarity VARCHAR(16) NOT NULL DEFAULT 'common',
  conditional_requirement_json JSONB NULL,   -- ör. happiness > 80
  species_family_bias VARCHAR(24)[] NOT NULL DEFAULT '{}',
  version INT NOT NULL DEFAULT 1
);
```

## 5. Trait Pool per Species

Species `base_traits_pool_json`:
```json
[
  { "trait_code": "keen_senses", "weight": 40 },
  { "trait_code": "swift_paw", "weight": 30 },
  { "trait_code": "iron_hide", "weight": 20 },
  { "trait_code": "pack_leader", "weight": 10 }
]
```

## 6. Trait Roll on Hatch

Pet instance yaratımında:
```go
func rollTraits(species PetSpecies, seed Seed, count int) []string {
    traits := []string{}
    remaining := species.BaseTraitsPool
    for i := 0; i < count; i++ {
        picked := seed.WeightedPick(remaining, weight fn)
        traits = append(traits, picked.TraitCode)
        remaining = remove(remaining, picked)  // no duplicates
    }
    return traits
}
```

Trait count:
- Common pet: 1 trait.
- Uncommon: 2.
- Rare: 2.
- Epic: 3.
- Legendary: 3-4.
- Mythic: 4.

## 7. Reroll Consumable (Phase 02)

Rare item "Trait Reroll Scroll":
- Consume → traits re-rolled from pool.
- Existing traits lost.

v1'de reroll yok; permanent trait.

## 8. Effect Application

Pet aktif iken:
- Trait effect stack'e ekle (bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md)).
- Conditional check (ör. happiness > 80).
- Scope: pet-specific (pet effect kontribüsyonu).

Owner character'a passive contribution:
- Pack Leader → party aura.
- Master Gatherer → character gather yield.

## 9. Trait Progression (Phase 02)

- Level 5, 15, 30 milestone: extra trait unlock.
- Bond level milestone: unique trait (unbreakable_bond).

v1'de fixed 2 trait.

## 10. Endpoint

**GET `/v1/content/pet-traits`**
Public.

**GET `/v1/characters/{id}/pets/{instance_id}`** — traits görünür instance detail'da.

## 11. Metrics
- `pet_trait_rolled_total{code,species}`.
- `pet_trait_distribution` histogram.

## 12. Cross-Reference
- [Species](species.md), [Instance](instance.md), [Skills](skills.md), [Combat](combat.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md) — modifier stack

## 13. Bekleyen Kararlar
- Reroll consumable Phase 02.
- Trait unlock progression Phase 02.
- Cross-species trait sharing (breeding) Phase 03.
- Signature legendary trait (unique) Phase 02.
