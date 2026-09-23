# Combat (Pet Combat Behavior)

## 1. Amaç
Pet aktif iken combat sırasında ne yapar — attack, cast, threat, downed state.

## 2. Combat Participation

Aktif combat pet:
- Slot: `combat`.
- Owner combat state ile senkron.
- Auto-attack + skill cast.
- Threat generation.

Utility pet: no combat participation (sadece passive stat).

## 3. Pet Combat Loop

Per combat tick:
```go
if pet.IsActive && pet.Slot == "combat":
    // Target selection
    if !pet.CurrentTarget || pet.CurrentTarget.IsDead:
        pet.CurrentTarget = pickTarget(owner.CurrentTarget)   // mirror owner
    // Auto-attack (weapon-based)
    if auto_attack_due(pet):
        petAutoAttack(pet, pet.CurrentTarget)
    // Skill cast (rotation)
    if pet.HasActiveSkill && pet.SkillCooldownReady:
        castPetSkill(pet, pet.ActiveSkill, pet.CurrentTarget)
```

## 4. Pet Auto-Attack

Species base_stats.attack_power ile:
```
attack_speed_ms = 2500 (baseline)   // pet-özel
damage = pet.base_ap × level_scale × happiness_mult × rarity_mult
```

Pet damage owner damage'a katkı (~%5-15).

## 5. Pet Skill Cast

Bkz. [skills.md](skills.md):
- Active skill per cooldown.
- Damage/CC/heal based.
- Stamina consume -1 per cast.

## 6. Pet Threat

Pet damage/heal threat generate:
- Threat modifier: 0.5 (halved vs character).
- Owner has +threat control (Guardian spec).
- Pet accidental aggro pull ihtimal; owner rescue.

## 7. Downed State

HP 0:
- Pet downed (not dead).
- Combat contribution 0.
- Redis flag `pet:downed:<id>` TTL 30 sec.
- Auto-revive to full HP after 30s.
- No permanent damage.

## 8. AFK Combat Contribution

Reveal (hunting AFK):
- Pet aktif → DPS contribution:
  ```
  pet_dps_contribution = pet.base_ap × level_scale × happiness_mult × rarity_mult × 0.10
  // (10% baseline; up to 25% for legendary pets)
  effective_party_dps = owner_dps × (1 + pet_dps_contribution/100)
  ```
- Faster kills; higher yield.

Analytical; per-cast simulate yok.

## 9. Pet Passive Aura

Pet passive skill aura owner + party (bkz. [skills.md §4-8](skills.md)):
- Modifier stack contribution (bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md)).
- Cache invalidation on activate/deactivate.

## 10. Pet Death (Not Real)

Pet ölmez (DD-024 pet lifecycle). Downed = temporary.

## 11. PvP Pet (Phase 02)

- Pet PvP damage (opsiyonel).
- Some pets PvP-only bonus.
- Balance: pet cap %10 damage in PvP (Phase 02 config).

## 12. Metrics
- `pet_damage_dealt_total`.
- `pet_downed_total`.
- `pet_dps_contribution_avg` gauge.

## 13. Cross-Reference
- [Instance](instance.md), [Skills](skills.md), [Stamina](stamina.md), [Happiness](happiness.md)
- [combat/README.md](../combat/README.md), [combat/combat-loop.md](../combat/combat-loop.md)

## 14. Bekleyen Kararlar
- Pet DPS cap final (%15 baseline, %25 legendary).
- Threat modifier 0.5 vs 0.7.
- Pet PvP Phase 02 balance.
- Multi-pet combat (2+ aktif) Phase 03.
