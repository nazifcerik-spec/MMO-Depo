# Evolution (Pet Evolution)

## 1. Amaç
Pet species evolves — daha güçlü form. Baby → adult → elder progression.

**v1'de opsiyonel** (schema hazır); tam Phase 02 rollout.

## 2. Evolution Model

Species template `evolution_stage` + `next_evolution_template_id`:
- Wolf Cub (stage 1) → Grey Wolf (stage 2) → Alpha Wolf (stage 3).
- Chain her stage ayrı species template.

## 3. Evolution Conditions

`evolution_condition_json`:
```json
{
  "level_min": 30,
  "bond_level_min": 50,
  "quest_completed": "pack_leader_quest",
  "items_consumed": [
    { "item_code": "evolution_stone", "qty": 1 }
  ],
  "happiness_min": 70
}
```

Örnek:
- Wolf Cub → Grey Wolf: level 20 + happiness 60 + item consumed.
- Grey Wolf → Alpha Wolf: level 40 + bond 80 + special quest.

## 4. Evolution Flow

**POST `/v1/characters/{id}/pets/{instance_id}/evolve`**

Server:
1. Load pet.
2. Validate:
   - Species has `next_evolution_template_id`.
   - Level >= condition.
   - Happiness >= min.
   - Bond >= min (Phase 02).
   - Quest completed (opsiyonel).
   - Items consumed available.
3. Transaction:
   - Consume items.
   - Update pet_owned.template_id = next_species.id.
   - Level reset to 1 (new species level 1); XP = 0.
   - Traits: **preserved** (existing traits carry over).
   - Skills: reroll base skills from new species; existing custom skills preserved (opsiyonel v1: reset all).
   - `evolved_from_instance_id` = old template_id (tarihçe için opsiyonel).
   - Version bump.
4. Emit event `pet.evolved`.
5. Notification.

Response 200: updated pet detail.

## 5. Reset vs Preserve

Karar (v1):
- **Level reset**: 1'e döner; yeni species learning curve.
- **XP reset**: 0.
- **Traits preserved**: existing trait roll korunur; ancak trait pool yeni species'ten.
- **Skills reset**: yeni species base skills.
- **Bond preserved** (Phase 02): bond level korunur; hatta bonus artabilir.
- **Happiness preserved**: current happiness korunur.

## 6. Evolution Bonus Stats

Evolved species base stats fazla:
- Wolf Cub base_HP 200 → Alpha Wolf base_HP 500.
- Level 1 Alpha stronger than level 30 Cub.

Trade-off: reset level ancak higher ceiling.

## 7. Evolution Cost

`items_consumed`:
- Common evolution: 1 evolution stone (nadir drop).
- Rare/legendary: multiple items + quest.

Vendor için evolution stone rare item.

## 8. Reverse Evolution

**v1 yok** — evolution kalıcı; geri dönüş yok.
Phase 03 opsiyonu: "de-evolution" consumable (nostalgia).

## 9. Evolution Chain Length

- Standart 2-3 stage per species.
- Legendary: 4 stage.
- Mythic (Phase 02+): 5+ stage.

## 10. Multi-Branch Evolution (Phase 02)

Species evolve to multiple forms:
- Wolf Cub → **Alpha Wolf** (physical DPS) OR **Shadow Wolf** (magic DPS).
- Choice locked.

v1'de tek path.

## 11. Metrics
- `pet_evolved_total{from,to}`.
- `pet_evolution_condition_fail_total`.

## 12. Audit
- `pet.evolved` (severity=info).

## 13. Cross-Reference
- [Species](species.md), [Instance](instance.md), [Level](level.md), [Traits](traits.md), [Skills](skills.md), [Happiness](happiness.md)

## 14. Bekleyen Kararlar
- Skills reset vs preserve final.
- Multi-branch evolution Phase 02.
- Evolution stone drop rate + vendor availability.
- Reverse evolution Phase 03.
