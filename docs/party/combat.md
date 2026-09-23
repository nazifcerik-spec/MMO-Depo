# Party Combat

## 1. Amaç
Party context'inde combat — aura sync, party bonus, kill credit.

Foundation: [combat/combat-loop.md](../combat/combat-loop.md).

## 2. Party Combat Session

- Party members combat session'a birlikte katılır (aynı dungeon instance veya hunt spot).
- Combat sistem her karakter için ayrı state ancak session shared.

## 3. Party Aura

Party synergy passives:
- Cleric aura → +5% all resist party-wide.
- Warden signature → damage share.
- Pack Leader (racial) → +5% party stat 20s active.

Aura apply:
- Passive: modifier stack per member (bkz. [skill-engine/passives.md §9](../skill-engine/passives.md)).
- Party leaves → aura removed.

## 4. Party Bonus (Flat v1)

Party ile combat (2+ members active):
- +5% XP per kill.
- +5% loot rare chance.

Config:
```yaml
party_bonus:
  xp_multiplier: 1.05
  loot_rare_multiplier: 1.05
```

Optimal composition bonus Phase 02.

## 5. Kill Credit

Bkz. [combat/death-revive.md §4](../combat/death-revive.md).

- Last-hit attributor kill credit (v1).
- Party assist bonus: any member >10% damage → credit share.
- XP grant to kill-credit character; party bonus applied.

## 6. Loot

Bkz. [loot.md](loot.md).

Personal loot (v1 default):
- Each character kill loot roll independent.
- No competition.

Alternative rules Phase 02.

## 7. Party DPS Aggregate

AFK reveal:
- Party members combined DPS estimate.
- Faster kill; higher yield per hour per member.

Formula:
```
party_dps = Σ member.DPS
kills_per_hour = 3600 / (avg_mob_hp / party_dps + spawn_interval)
```

## 8. Death & Wipe

- Death: individual XP debt, durability loss (bkz. [combat/death-revive.md](../combat/death-revive.md)).
- Party wipe: all members dead → dungeon encounter reset OR hunt spot restart.
- Cleric resurrect skill can revive mid-encounter.

## 9. Party Chat

Combat event chat (opsiyonel; v1 basic):
- "Boss defeated!"
- "Player X died."
- Chat channel `party`.

## 10. Cross-Region Combat

Members different region OK ancak combat instance location-bound. Party dungeon activity anchor.

## 11. Metrics
- `party_combat_kills_total`.
- `party_xp_bonus_applied_total`.
- `party_wipe_total`.

## 12. Cross-Reference
- [Party](party.md), [Roles](roles.md), [Hunt](hunt.md), [Dungeon](dungeon.md), [Loot](loot.md)
- [combat/combat-loop.md](../combat/combat-loop.md), [combat/death-revive.md](../combat/death-revive.md)
- [skill-engine/passives.md](../skill-engine/passives.md)

## 13. Bekleyen Kararlar
- Party composition bonus formula Phase 02.
- Voice chat integration Phase 03.
- Combat log per-member Phase 02.
