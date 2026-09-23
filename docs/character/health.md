# Health (HP)

## 1. Amaç
Karakterin can havuzu; hasar, iyileşme, regen, ölüm ve dirilme akışı.

## 2. Model
- **Max HP**: derived stat (bkz. [stats.md §3.2](stats.md)).
- **Current HP**: runtime state; server-authoritative.
- Persistence:
  - **Ephemeral (Redis)**: aktif combat/aktivite sırasında güncel HP.
  - **Snapshot (Postgres)**: karakter logout/aktivite değişimi/save point.
- Death durumunda: HP = 0 → death event.

## 3. HP Kaynakları

### 3.1 Max HP Formülü
```
max_hp = 50 + 10 * effective_STA + 5 * level + equipment_hp_bonus + buff_hp_bonus
```
Örnek Warrior lv 25, STA 40, gear +150 HP:
```
max_hp = 50 + 10*40 + 5*25 + 150 = 725
```

### 3.2 Regen Formülü
```
hp_regen_per_2s = 1 + 0.1 * effective_STA
combat_regen_multiplier = 0.2   -- combat sırasında %20
```

Regen event her 2 sn:
- Out-of-combat: `current_hp = min(max_hp, current_hp + hp_regen_per_2s)`.
- In-combat: `... + hp_regen_per_2s * 0.2`.

Sitting/rest bonus (Phase 02): +50% regen.

## 4. Damage & Heal Uygulaması

### 4.1 Damage Path
```
raw_damage = ability.base_damage + scaling
after_armor = raw_damage * (1 - armor_dr(raw_damage, target_armor, attacker_level))
after_resist = after_armor * (1 - magic_resist_dr(...))  // only if magic
after_evasion = target.evaded ? 0 : after_resist
after_block = after_evasion * (1 - block_reduction if blocked)
after_lifesteal_credit = attacker gains X% of after_block
final_damage = clamp(after_block, 1, MAX_INT)  // min 1 always
current_hp -= final_damage
if current_hp <= 0: trigger_death()
```

### 4.2 Heal Path
```
raw_heal = ability.base_heal + heal_scaling
if target.dead: reject or resurrect skill only
new_hp = min(max_hp, current_hp + raw_heal)
overheal = raw_heal - (new_hp - current_hp)  // metrics
current_hp = new_hp
```

### 4.3 Percent-Based Damage
Effects: `dmg_pct_max_hp` (bosses' execute):
```
final = target.max_hp * pct
```
Cap % damage per second (e.g. max 30% max_hp/sec) — nuke koruma.

### 4.4 Shield Absorb
Absorb buff:
- Insert `stats_buff` with `shield_amount`.
- Damage first hits shield; residual → HP.
- Shield expire when consumed veya duration.

## 5. Death

### 5.1 Trigger
`current_hp <= 0` (server-side).

### 5.2 State Transition
```
character.death_state:
  is_dead BOOLEAN
  died_at TIMESTAMPTZ
  respawn_available_at TIMESTAMPTZ
  cause_kind, cause_ref_id
```

Schema önerisi:
```sql
ALTER TABLE character.character
  ADD COLUMN is_dead BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN died_at TIMESTAMPTZ NULL,
  ADD COLUMN respawn_available_at TIMESTAMPTZ NULL;
```

### 5.3 Death Sequence
```
BEGIN;
1. UPDATE character SET is_dead=true, died_at=now(),
      respawn_available_at=now() + interval '30 seconds',
      version = version + 1
   WHERE id=$1;
2. Insert character.death_log (bkz. database/character.md §16).
3. Compute XP debt (bkz. xp.md §7):
   - level >= 20: xp_debt += level_xp_required * 0.05
4. Durability loss: her equipped item durability -= max(1, durability_max * 0.05).
5. Interrupt aktif activity contract (status='paused', pause_reason='death').
6. Party/dungeon encounter: kill credit rules; wipe check.
7. Publish 'combat.death_recorded' event.
8. WS push 'death' frame to client (respawn UI).
COMMIT;
```

### 5.4 Respawn / Revive

#### Self Respawn
**POST `/v1/characters/{id}/revive`**
```json
{ "mode": "graveyard" }  // "graveyard" | "resurrect_scroll" | "party_revive_accept"
```

Server (graveyard mode):
1. `now >= respawn_available_at` check.
2. Teleport to `home_city_id.spawn_point`.
3. `current_hp = max_hp * 0.5`.
4. `current_mana/resource = 0`.
5. `is_dead=false`, timestamps clear.
6. Apply "Weakness" debuff 60 sn (-10% dmg).
7. Response 200: character state + location.

#### Party Revive Skill
- Cleric `Resurrect` skill: target dead within 60 sn.
- Accept flow: dead player alır UI prompt; onayla → respawn at death location.

#### Resurrect Scroll (Consumable)
- Item consume → self revive at death location.

#### Auto-Revive (Phase 02)
- 5 dk sonra otomatik graveyard revive (idle-friendly).

### 5.5 Weakness Debuff
- 60 sn `-10% damage output`, `-10% healing received`.
- Stacking: her ölümde refresh (max 3 stack).

## 6. Death Penalty (Özet)
Bkz. [xp.md §7](xp.md), [combat.md §Failure](../game-design/combat.md).

- XP debt (level 20+, %5 of level).
- Durability -5% her equipped item.
- Respawn timer 30 sn (self); 5 dk auto.
- Weakness debuff 60 sn.
- Item loss YOK (DD-019).

## 7. AFK Combat Death
AFK hunting sırasında ölüm:
- Activity contract → paused, pause_reason='death'.
- Death log insert.
- Login sırasında reveal → "died N times" özet + XP debt applied.
- Otomatik revive: 5 dk sonra self at graveyard (activity resume opsiyonu client'a).

## 8. Ephemeral vs Persistent State

**Ephemeral (Redis)**:
- Key: `char:hp:<char_id>` → `{current_hp, current_mana, current_resource, last_updated_at}`.
- TTL: 24h (keep-alive on read).
- Aktif combat/session sırasında read/write.

**Snapshot (Postgres)**:
- `character.character` alanları: `current_hp_snapshot`, `current_mana_snapshot`.
- Save triggers:
  - Session end (logout).
  - Aktivite değişimi.
  - Death.
  - Heartbeat (5 dk).

Redis kayıp → snapshot'tan restore (biraz yaşlı ama safe).

## 9. Endpoint

### 9.1 HP Status

**GET `/v1/characters/{id}/vitals`**
```json
{
  "hp": { "current": 425, "max": 725, "regen_per_2s": 5.0 },
  "mana": { "current": 78, "max": 155, "regen_per_2s": 4.0 },
  "resource": { "kind": "rage", "current": 60, "max": 100 },
  "is_dead": false
}
```

### 9.2 Revive (bkz. §5.4)

## 10. WebSocket Push

Frame:
```json
{ "type": "vitals_update", "payload": { "hp": 425, "hp_max": 725, "delta": -50, "cause": "mob:wolf_alpha" } }
```

Push sıklığı:
- Combat sırasında: her damage/heal.
- Combat dışı: her regen tick (opsiyonel — client hesabı yapabilir).
- Client throttle: max 10Hz.

## 11. Formül Config
Config'te (`config.combat.hp.yaml`):
```yaml
base_max_hp: 50
sta_per_hp: 10
level_per_hp: 5
regen_base: 1.0
regen_per_sta: 0.1
combat_regen_mult: 0.2
weakness_pct: 10
weakness_duration_sec: 60
respawn_grace_sec: 30
auto_revive_sec: 300
death_xp_debt_pct: 0.05
death_durability_pct: 0.05
death_penalty_min_level: 20
```

## 12. Fail Modları
- Concurrent damage → optimistic lock on HP; last-writer-wins with delta merge (relative delta safe).
- Redis unreachable → fallback to snapshot; write-back queue.
- Death race condition (multi-hit last blow): kill credit rules by highest damage attributed.

## 13. Metrics
- `character_death_total{cause}` counter.
- `character_dps_dealt_per_sec` histogram.
- `character_hp_ratio` gauge (histogram of current/max).
- `character_overheal_ratio` gauge.

## 14. Audit
- `character.died` (`audit.event_log` + `character.death_log`).
- `character.revived`.
- Suspicious: rapid death cycle (potential exploit) → security_signal.

## 15. Cross-Platform
Server tek gerçek; her platform aynı HP. WS push senkron.

## 16. Test Stratejisi
- Formul unit test (damage/heal/regen).
- Death → respawn happy path.
- XP debt integration.
- Concurrent damage (goroutine).
- AFK death (reveal senaryosu).

## 17. İlişkiler
- [XP](xp.md), [Stats](stats.md), [Mana](mana.md), [game-design/combat.md](../game-design/combat.md), [systems/afk-engine.md](../systems/afk-engine.md), [database/character.md](../database/character.md), [backend/websocket.md](../backend/websocket.md)

## 18. Bekleyen Kararlar
- Auto-revive Phase 01'de aktif mi (5 dk).
- Weakness debuff süresi (60 sn önerisi).
- Elemental damage sistemi (Phase 02).
- Hardcore mode permadeath (Phase 03).
- Death-cam UI (opsiyonel).
