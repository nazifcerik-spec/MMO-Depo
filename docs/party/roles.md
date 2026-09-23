# Roles

## 1. Amaç
Party içindeki rol hint — tank/heal/dps/support. UI + matchmaking hint.

## 2. Role Types

| Role | Class Örnekleri |
|---|---|
| `tank` | Warrior Guardian, Cleric Templar |
| `heal` | Cleric Priest, Druid Grovekeeper |
| `dps` | Warrior Berserker, Mage, Rogue, Ranger |
| `support` | Druid Grovekeeper, Cleric Templar (buff) |
| `dps_ranged` | Ranger, Mage |
| `dps_melee` | Warrior Berserker, Rogue |

## 3. Model

`social.party_member.role_hint VARCHAR(16) NULL`.

Character self-assign; leader override (Phase 02).

## 4. Role Compatibility Suggestion

Class default role hint (character sheet):
- Warrior Guardian → tank.
- Warrior Berserker → dps_melee.
- Mage → dps_ranged.
- Cleric Priest → heal.

Client shows default; user can override.

## 5. Party Composition (Ideal)

Dungeon standard party (5-man):
- 1 tank + 1 heal + 3 dps.

Alternatif (small dungeon):
- 3-man: 1 tank + 1 heal + 1 dps.
- Solo idle: everything.

Server doesn't enforce composition; suggestion only.

## 6. Party Bonus by Composition (Phase 02)

Optimal composition:
- Tank + Heal + 3 DPS → party bonus +5% dmg (Phase 02).

v1'de flat party bonus (bkz. [combat.md](combat.md)).

## 7. Role Change

**POST `/v1/parties/{id}/role`** (self change)
```json
{ "role_hint": "dps" }
```

Server:
- Update party_member.role_hint.
- Emit event.

## 8. LFG Role Filter

LFG listing role_hint:
- Party leader searches "tank" — see LFG listings tagged as tank.

## 9. Class Restrictions

Bazı role class-locked:
- `heal` — Cleric/Druid only (v1).
- `tank` — Warrior/Cleric Templar (v1).
- `dps` — any.

Player rol seçerken class validation (soft warning, not blocked).

## 10. Metrics
- `party_role_distribution{role}` gauge.
- `party_composition_type_avg` histogram.

## 11. Cross-Reference
- [Party](party.md), [Member](member.md), [Combat](combat.md), [Dungeon](dungeon.md)
- [game-design/classes.md](../game-design/classes.md), [game-design/specializations.md](../game-design/specializations.md)

## 12. Bekleyen Kararlar
- Optimal composition bonus Phase 02.
- Role change combat lock (Phase 02: sadece out-of-combat).
- Auto-assign role by class default v1 (soft opt).
