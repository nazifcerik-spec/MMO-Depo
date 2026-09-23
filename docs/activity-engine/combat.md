# Combat Activity (Manuel + Rotation AFK)

## 1. Amaç
Aktif oynanış aktivitesi. Manuel skill cast + rotation AFK combat + party context. `Hunting` bir alt kategori; ancak "combat" burada **generic** AFK combat davranışını tanımlar (dungeon boss encounter'da rotation, PvP Phase 02 vs).

## 2. Combat Modları

| Mod | Kategori | Notlar |
|---|---|---|
| **Manuel** | interactive | Client cast pipeline; server otoritedir |
| **Rotation AFK** | contract-based | Auto-cast rotation; hunting/dungeon idle içinde |
| **Party** | shared | Party activity; her karakter kendi contract'ı ancak coordinated |
| **World Boss** | event | Contribution scoring; timed |
| **PvP** (Phase 02) | direct | duel/arena/open-world |

## 3. Manuel Combat Flow

Client → cast intent:
```
POST /v1/characters/{id}/cast
  body: { skill_code, target, source: "manual" }
```

Handler → [skill-engine](../skill-engine/README.md) engine.Cast → sonuç WS push.

Manuel combat aktivite contract'ı **oluşturmaz** — bir "session-level" state. Ancak damage/heal/kill sayaçları için opsiyonel "session summary" tutulur:

```
character.combat_session (opsiyonel Phase 02):
  character_id, started_at, ended_at,
  total_damage_dealt, total_healing_done, total_deaths, kill_count
```

v1'de: `audit.event_log` üstünde combat events (per-cast); post-hoc aggregate.

## 4. Rotation AFK Combat

Contract-based; hunting/dungeon idle içinde eval. Detay [hunting.md §7 DPS estimate](hunting.md), [dungeon.md](dungeon.md).

Rotation eval loop (AFK reveal içinde):
- DPS estimate (analytical, per-tick loop yok).
- Threshold detection (death cascade, low HP potion out).
- Yield compute.

## 5. Party Combat

Party içinde her karakter kendi contract'ını çalıştırır (party synergy passive/aura ile).

Party AFK dungeon:
- `activity.instance` party_id ile.
- Her karakter için `activity.contract` party_dungeon type.
- Reveal her karakter için ayrı (party leader'ın reveal'ı çekmez; her karakter kendi login'inde).

Cross-character sync:
- Party aura → snapshot içinde `party_bonus_modifiers` field.
- Party leaves → aura kaldırılır (contract snapshot güncellemez ancak yield'e etkisi limitli).

## 6. World Boss Combat

Bkz. [game-design/world-bosses.md](../game-design/world-bosses.md).

- Manuel katılım gerekir (AFK yasak).
- Contract oluşturulmaz; ephemeral session.
- Contribution `progression.world_boss_participation` tablosunda.

## 7. Combat AFK Rotation Editor

`character.rotation.priority_list`:
```json
[
  { "cond": { "hp_below_pct": 40 }, "skill_code": "healing_potion" },
  { "cond": { "resource_at_least": 60 }, "skill_code": "berserker_rampage" },
  { "cond": { "any": true }, "skill_code": "warrior_bash" }
]
```

Condition types:
- `hp_below_pct: N`
- `hp_above_pct: N`
- `resource_at_least: N`
- `resource_below: N`
- `target_hp_below_pct: N`
- `buff_active: <buff_code>`
- `debuff_active_on_target: <debuff_code>`
- `combo_at_least: N` (Phase 02)
- `any: true` (fallback)

Client rotation editor UI; save to DB.

Rate compute AFK contract (hunting):
- `skillEngine.EstimateDPS(rotation, snap, targets)` → analytical.

## 8. Combat Events

Combat event bus fires (bkz. [skill-engine/passives.md §3](../skill-engine/passives.md)):
- `on_hit_dealt`, `on_hit_taken`, `on_crit_dealt`, `on_kill`, `on_death`, `on_hp_below`, etc.

Passive triggers evaluated per event (in manuel combat); AFK contract reveal aggregate.

## 9. Death Handling (Combat Context)

- Damage taken → HP down → 0 → death.
- Death event: `character.death_log`.
- Respawn timer 30s (self) veya 5min auto.
- XP debt applied.
- Contract paused if AFK; manuel combat session ended.

Detay: [character/health.md §5](../character/health.md).

## 10. Combat Metrics Aggregation

Per combat session (manuel):
- `combat.dps_dealt_avg` histogram.
- `combat.tps_taken_avg` histogram.
- `combat.kills_per_session`.

Per AFK contract:
- Already covered in hunting metrics.

## 11. Combat Balance Considerations

- Manuel oyun ~%15-25 daha verimli AFK'ya kıyasla (DD-open balance).
- Rotation editor karmaşıklık cap: max 20 entry (config).
- Priority list eval per cast time cap: 100 iterations (safety).

## 12. Cross-Reference
- [Activity](activity.md), [Hunting](hunting.md), [Dungeon](dungeon.md)
- [skill-engine/active-skills.md](../skill-engine/active-skills.md), [skill-engine/passives.md](../skill-engine/passives.md), [skill-engine/scaling.md](../skill-engine/scaling.md)
- [character/health.md](../character/health.md)
- [game-design/combat.md](../game-design/combat.md), [game-design/skills.md](../game-design/skills.md)

## 13. Bekleyen Kararlar
- Combat session tracking (persist vs event only) Phase 02.
- PvP damage split (Phase 02 config).
- Party sync aura across contracts (v1 snapshot-based; Phase 02 dynamic).
- Interrupt system (skill lock on damage) Phase 02.
- Combo point mekaniği (Rogue spec) Phase 02.
