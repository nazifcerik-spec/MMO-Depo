# Hunting

## 1. Amaç
Hunt aktivite deep-dive: params, filter, party context, edge cases.

Foundation: [activity-engine/hunting.md](../activity-engine/hunting.md) — activity contract + analitik yield.

## 2. Hunt Aktivite Params

```json
{
  "spot_id": "loc-whispering-forest",
  "mob_type_codes": ["ashen_wolf", "shadow_wolf"],
  "difficulty_filter": "normal",       // easy/normal/hard/elite
  "autopotion": true,
  "stop_on_hp_below_pct": 30,
  "stop_after_n_deaths": 3,
  "use_pet": true,
  "party_id": null                      // party ile hunt
}
```

## 3. Filter Behavior

- **mob_type_codes**: pool subset; sadece bunlar target.
- **difficulty_filter**: spot danger tier'ıyla eşleşen mob'lar (normal/hard elite ayrı).
- **stop_on_hp_below_pct**: safety threshold — reveal death cascade prevention.

Boş `mob_type_codes`: spot'un tüm pool'u aktif.

## 4. Yield Rate Compute

Bkz. [activity-engine/hunting.md §4-5](../activity-engine/hunting.md).

Özet formula:
```
avg_mob_hp = averageMobHP(spot, params.mob_type_codes)
dps = skillEngine.EstimateDPS(snap, rotation, filter)
time_per_kill_sec = avg_mob_hp / dps + spot.spawn_interval_sec
kills_per_hour = 3600 / time_per_kill_sec
```

Party bonus: bkz. [party.md] (Phase 02).

## 5. Party Hunt

Party context:
- Party leader hunt spot'ta; members karakter bazlı yield (v1 solo compute her karakter için).
- Party synergy: +%5-10 party bonus (rotation + pet aura).
- Loot rule: personal loot (DD-012).

Party contract:
- Her karakter kendi contract'ı (v1).
- Phase 02: shared party contract (leader-anchored).

## 6. Pet Assist

Snapshot `ActivePetInstanceID` var:
- Combat pet: DPS + threat.
- Utility pet: yield bonus.

Bkz. [game-design/pets.md](../game-design/pets.md), [activity-engine/hunting.md §8](../activity-engine/hunting.md).

## 7. Consumables Auto-Use

- Autopotion (health/mana) — envanterde varsa combat kill'inde consume.
- Estimated potion consumption: `kills × consume_per_kill` (bkz. threshold detection).
- Threshold: potion out → paused veya continue without.

## 8. Death Handling

Bkz. [activity-engine/hunting.md §10](../activity-engine/hunting.md), [combat/death-revive.md](../combat/death-revive.md).

Cascade prevention:
- `stop_after_n_deaths: 3` — 3 death sonrası contract pause.
- Safety önemli AFK reveal içinde.

## 9. Endpoint

**POST `/v1/characters/{id}/hunt`** — start hunt contract (bkz. activity-engine).

**GET `/v1/characters/{id}/hunt/current`** — aktif hunt state + progress.

## 10. Reveal Örnek

Bkz. [activity-engine/hunting.md §11](../activity-engine/hunting.md).

## 11. Rare Mob Interaction

Hunt sırasında rare mob spawn (bkz. [rare-enemy.md](rare-enemy.md)):
- Reveal içinde rare mob spawn detection.
- Rare mob kill guaranteed (character DPS yeter ise).
- Rare drop tetiklenir.
- Notification "Rare mob defeated!"

## 12. Random Event Interaction

Hunt aktivitesi sırasında random event tetiklenebilir (bkz. [random-event.md](random-event.md)):
- Bandit ambush, treasure spawn, event mob.
- Reveal aggregate; per-event narration özet.

## 13. Metrics

Bkz. [activity-engine/hunting.md §12](../activity-engine/hunting.md).

## 14. Cross-Reference
- [activity-engine/hunting.md](../activity-engine/hunting.md), [combat/README.md](../combat/README.md)
- [Rare Enemy](rare-enemy.md), [Elite Enemy](elite-enemy.md), [Random Event](random-event.md), [Treasure](treasure.md)
- [world/monster-area.md](../world/monster-area.md)

## 15. Bekleyen Kararlar
- Party shared contract Phase 02.
- Autopotion tier priority (cheap first vs strongest first).
- Consume-per-kill cap Phase 02 tune.
