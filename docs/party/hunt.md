# Party Hunt

## 1. Amaç
Party ile hunt aktivite — coop mob farming; her karakter ayrı contract, party synergy.

Foundation: [activity-engine/hunting.md](../activity-engine/hunting.md), [hunting-adventure/hunting.md](../hunting-adventure/hunting.md).

## 2. Party Hunt Model

- Party leader hunt spot'ta start.
- Members follow (kendi contract'ları primary channel).
- Aura + party bonus applied.

## 3. Contract Model

Her karakter kendi `activity.contract` (activity_type='hunting'):
- Same spot_id.
- Same mob_type_codes (opsiyonel; leader chooses default).
- Party bonus modifier applied per character snapshot.

Alternate Phase 02: shared party contract (leader-anchored, members share yield).

v1 tercih: **her karakter kendi contract**. Simple + fair.

## 4. Rate Compute (Party Aggregate)

Party DPS aggregate:
```
party_dps = Σ member.DPS
```

Per character kill rate:
```
kills_per_hour_per_char = party_dps / avg_mob_hp / party_size × 3600
```

Actually simpler: party DPS combined → kills happen faster; loot roll per character:
```
combined_kills_per_hour = party_dps / avg_mob_hp × 3600
per_member_loot_rolls_per_hour = combined_kills_per_hour   (personal loot each)
```

## 5. Party Bonus in AFK Reveal

Each character reveal includes party bonus modifier snapshot:
- `+5% XP`, `+5% rare loot chance`.

Snapshot taken at contract start; party disband mid-hunt → bonus removed at next snapshot refresh (v1: contract lifetime bonus locked; Phase 02 dynamic).

## 6. Pet & Mount Interaction

- Party members' active pets aggregate DPS.
- Party member mounts don't stack.

## 7. Rare Mob / Random Event

Party shared rare mob spawn (v1 personal per bkz. DD-110; ancak party context Phase 02 shared spawn opsiyonu):
- v1: rare mob personal per character.
- Phase 02: party context shared rare (first-kill share reward).

Random event: apply party-wide (blessing/curse).

## 8. Party Hunt Endpoint

**POST `/v1/characters/{id}/hunt`** — start with `party_id` param.
Server:
- Validate party membership.
- Create hunt contract with party_bonus flag.
- Emit event `party.hunt_started`.

## 9. Metrics
- `party_hunt_active_total`.
- `party_hunt_kills_per_hour_avg`.

## 10. Cross-Reference
- [Party](party.md), [Combat](combat.md), [Loot](loot.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)
- [hunting-adventure/hunting.md](../hunting-adventure/hunting.md)

## 11. Bekleyen Kararlar
- Shared party contract Phase 02.
- Party rare mob spawn shared Phase 02.
- Party hunt XP share formula (v1: each char full XP + party bonus; alt fair share Phase 02).
