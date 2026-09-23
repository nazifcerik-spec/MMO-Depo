# Skills (Pet Abilities)

## 1. Amaç
Pet'in aktif/passive abilities — species base skill + rank up.

## 2. Skill Structure

Pet skill = `content.skill_template` benzeri; ancak pet-özel:
- `type = 'pet_active' | 'pet_passive'`.
- `usable_by_pet = true` flag.
- No character resource cost (pet auto-cast).

## 3. Skill Slot

Pet 2 skill slot:
- **Active skill**: cooldown-based; combat cast.
- **Passive skill**: sürekli aura.

Species `base_skill_codes` array — pet hatched'de öğrenilmiş.

## 4. Örnek Pet Skills

### Combat Pet Skills

| Code | Type | Effect | Cooldown |
|---|---|---|---|
| `pet_bite` | active | melee damage + bleed | 6s |
| `pet_pounce` | active | leap + stun 1s | 15s |
| `pet_growl` | passive | +5% owner threat | – |
| `pet_ferocity` | passive | +10% pet damage vs low HP target | – |
| `pet_pack_hunter` | passive | Owner + pet +5% dmg when both engaged | – |

### Gathering Pet Skills

| Code | Type | Effect |
|---|---|---|
| `pet_forage_boost` | passive | +5% herbalism yield |
| `pet_mining_boost` | passive | +5% mining yield |
| `pet_carry_capacity` | passive | +10 inventory slots |
| `pet_lucky_paw` | passive | +3% treasure spawn chance |

### Utility Pet Skills

| Code | Type | Effect |
|---|---|---|
| `pet_swift_wings` | passive | +5% movement speed |
| `pet_night_vision` | passive | +1 dungeon vision tile (Phase 02) |
| `pet_bag_bonus` | passive | +8 inventory slot |

## 5. Rank Up

Skill ranks (1..5):
- Rank up: pet leveling + happiness > 60 gate.
- Each rank +20% effect magnitude.
- Level 10/20/30/40 unlocks rank 2/3/4/5.

`companion.pet_owned.skill_ranks_json`:
```json
{
  "pet_bite": 3,
  "pet_growl": 2
}
```

## 6. Combat Skill Cast

Pet aktif ise combat tick loop pet cast:
```go
if pet.IsActive && pet.CurrentTarget && pet.CooldownReady(activeSkill) {
    castPetSkill(pet, activeSkill, pet.CurrentTarget)
    pet.SetCooldown(activeSkill)
}
```

Pet skill damage:
```
damage = base × (1 + rank × 0.2) × pet.EffectivePower × happiness_multiplier
```

## 7. Passive Skill Apply

Pet passive skill stack modifier'a ekle:
- Owner character modifier stack'e katkı (aura scope).
- Pet stat contribution.

Cache invalidation: pet activate/deactivate.

## 8. Skill Content

Pet skill'ler `content.skill_template` ile aynı tablo; `type='pet_active'` veya `type='pet_passive'` filter.

## 9. AFK Skill Contribution

Pet aktif AFK hunting sırasında:
- Passive skill: yield modifier stack.
- Active skill: DPS estimate'e ekle (rate compute).

Analitik; per-cast simulate yok.

## 10. Metrics
- `pet_skill_cast_total{skill}`.
- `pet_skill_rank_up_total`.

## 11. Cross-Reference
- [Species](species.md), [Instance](instance.md), [Level](level.md), [Traits](traits.md), [Combat](combat.md), [Hunting](hunting.md)
- [skill-engine/README.md](../skill-engine/README.md)

## 12. Bekleyen Kararlar
- Signature legendary pet skill Phase 02.
- Skill point spend system (Phase 02 talent).
- Skill loadout swap (multiple skill sets) Phase 02.
