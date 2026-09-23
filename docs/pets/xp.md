# XP (Pet XP Sources)

## 1. Amaç
Pet XP kaynakları + multiplier + happiness effect. Character XP'den ayrı ledger.

## 2. XP Kaynakları

| Kaynak | Rate |
|---|---|
| Character combat kill (aktif pet) | 20% of character XP per kill |
| Character AFK hunting (aktif pet) | 20% of yield character XP |
| Feeding (special treat consumable) | flat +N XP (Phase 02) |
| Quest reward (pet quest) | 500-5000 XP (Phase 02) |
| Event (pet-focused) | Variable |

Pet passive learning — karakter aktivite yaparken pet aktif ise yield paylaşır.

## 3. XP Grant Formula

```go
func petXPFromCharacterActivity(charXPGained int64, pet PetInstance) int64 {
    if !pet.IsActive:
        return 0
    baseShare := charXPGained × 0.20
    happinessMultiplier := 0.5 + (pet.Happiness / 100) × 0.5   // 0.5..1.0
    return int64(baseShare × happinessMultiplier)
}
```

Happiness → XP multiplier (motivate care):
- Happiness 100 → full XP share.
- Happiness 50 (baseline) → 75%.
- Happiness 0 → 50% (min).

## 4. XP Cap

Per hour cap: 20000 pet XP (bot koruma).

## 5. XP Ledger

Simple: `companion.pet_owned.xp` alanı; no dedicated ledger v1.

Phase 02 pet_xp_ledger:
```sql
CREATE TABLE companion.pet_xp_ledger (
  id BIGSERIAL PK,
  pet_instance_id UUID NOT NULL,
  delta BIGINT NOT NULL,
  source_kind VARCHAR(24) NOT NULL,
  source_ref_id UUID NULL,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## 6. Reveal AFK Integration

Character AFK activity reveal içinde:
```go
if activePet {
    petXP := computePetXPShare(activityYield.XPCharacter, activePet)
    petService.GrantXP(activePet.ID, petXP)
    revealPayload.pet_xp_gained = petXP
}
```

## 7. Multi-Pet Priority

- Combat pet: 20% XP share (primary).
- Utility pet: 10% XP share (secondary; Phase 02).

v1'de combat pet only (utility pet XP earn Phase 02).

## 8. XP Debt

Pet XP debt YOK (v1). Pet die-and-revive ceza yok.

## 9. Metrics
- `pet_xp_granted_total{source}`.
- `pet_xp_per_hour_avg{species}`.

## 10. Cross-Reference
- [Level](level.md), [Happiness](happiness.md), [Instance](instance.md)
- [character/xp.md](../character/xp.md) — pattern reference
- [activity-engine/xp.md](../activity-engine/xp.md) — reveal integration

## 11. Bekleyen Kararlar
- XP share ratio 20% final Playtest.
- Utility pet XP Phase 02.
- Pet-specific quest chain Phase 02.
- XP boost consumable (pet treat) Phase 02.
