# Stamina (Pet Stamina)

## 1. Amaç
Pet efort sınırı. Combat/hunting sırasında tüketilir; rest ile geri gelir. Balance mechanic.

**v1'de aktif**; pet için stamina meaningful.

## 2. Model

`companion.pet_owned.stamina_current SMALLINT 0..100`, `stamina_max SMALLINT DEFAULT 100`.

Baseline: 100/100.

## 3. Stamina Kaynağı & Tüketim

Tüketim:
- **Combat cast**: -1 per skill cast (aktif skill).
- **Combat auto-attack**: -0.2 per hit.
- **Hunting AFK**: -2 per hour (aggregate).

Recovery:
- **Rest** (city, stable): +5/min.
- **Idle** (aktif değil, character offline): +2/min.
- **Deactivated in slot**: +5/min.
- **Feed**: +10 stamina bonus (food side-effect).

## 4. Effect

Stamina 0 → pet **exhausted**:
- Skill cast blocked.
- Auto-attack -%50 damage.
- Cannot activate (must rest).

Notification: "Your pet is exhausted; rest required."

## 5. Rest Flow

**POST `/v1/characters/{id}/pets/{instance_id}/rest`**
Server:
1. Load pet.
2. Requires:
   - Character at city (safe zone).
   - Pet activated OK (rest sırasında pet deactivated).
3. Start rest state: `stamina_last_recovery_at = now`.
4. Passive recovery per tick (worker or on-read compute).

Rest duration: full recovery ~20 min from 0 → 100.

## 6. AFK Hunting Stamina

Reveal aggregate:
```
hunting_duration_hours := elapsed / 3600000
stamina_drain := hunting_duration_hours × 2
if pet.stamina_current - stamina_drain < 0:
    hunting_active_duration = pet.stamina_current / 2   // hours before exhaust
    threshold: pet exhausts at t_exhaust
    yield capped at t_exhaust; pet contribution stops
```

Threshold detection linear inverse (DD-086).

## 7. Stamina Cap Upgrade (Phase 02)

Base 100 cap; upgrade via:
- Species evolution → +20 cap.
- Bond level milestones → +10 cap per milestone.
- Rare consumable "Endurance Elixir" → +5 permanent cap.

Max cap: 200 (Phase 02).

## 8. Alternative: Stamina v1 Disabled Option

**Karar tercih**: v1'de aktif — pet care meaningful. Character stamina v1'de kapalı (DD-054) ancak pet için mantıklı (idle-frequent activation).

## 9. Metrics
- `pet_stamina_avg` gauge.
- `pet_exhausted_total`.
- `pet_rest_started_total`.

## 10. Cross-Reference
- [Instance](instance.md), [Happiness](happiness.md), [Combat](combat.md), [Hunting](hunting.md)
- [character/stamina.md](../character/stamina.md) — character stamina (v1 kapalı; contrast)
- [game-design/pets.md](../game-design/pets.md)

## 11. Bekleyen Kararlar
- Stamina v1 aktive final (öneri: evet).
- Rest location gating (v1 city sadece; Phase 02 anywhere out-of-combat).
- Auto-rest passive tick worker interval.
- Endurance Elixir Phase 02.
