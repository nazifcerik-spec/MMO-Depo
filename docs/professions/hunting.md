# Hunting (Profession)

## 1. Amaç
Mob kill'den gelen tracking/skinning/beast lore. **Hunting aktivite ≠ hunting profession**:
- **Hunting aktivite**: [activity-engine/hunting.md](../activity-engine/hunting.md) — AFK combat.
- **Hunting profession**: aktivite yan-effect; kill'lerden extra bounty, beast track.

## 2. Profession Data

```yaml
# content/professions/hunting.yaml
code: hunting
category: gathering
display_name_key: profession.hunting.name
max_level: 100
xp_curve: { base: 50, factor: 1.10 }
tool_slot: main_hand              # ranged weapon (bow/crossbow)
tool_alternate_slots: [main_hand] # melee also OK; hunting profession pasif
tool_template_codes: []            # any weapon works
race_affinity_bonus_pct:
  orc: 8
```

## 3. Profession Yield vs Combat Yield

Hunting aktivitesi 2 ayrı yield üretir:
- **Combat yield** (loot table): item drops, character XP — [activity-engine/hunting.md §5](../activity-engine/hunting.md).
- **Hunting profession yield**: bounty tracking, bestiary progression, rare mob tag.

Profession yield hesabı:
```
hunting_profession_xp_per_kill = 5 (baseline)
+ bonus per mob family unlocked
```

## 4. Bestiary Integration

Kill'ler mob family bestiary'e katkı (bkz. [skill-engine/mastery.md](../skill-engine/mastery.md)):
- Her mob family unique bestiary track.
- Milestone: +damage vs family, drop bonus.

Hunting profession level bestiary track'lerini boost eder.

## 5. Bounty Contracts (Phase 02)

Hunting profession trainer → bounty contract:
- "Kill 50 shadow wolves for 100 gold + rare drop token".
- Daily contracts.
- Reputation with hunter faction (Phase 02).

v1'de basic; Phase 02 tam bounty system.

## 6. Race Affinity

**Orc**: `orcish_bloodlust` +5% crit dmg combat + affinity +8% hunting profession XP.

## 7. Rank Tier'lar

Aynı pattern. Grandmaster: **"Master Hunter"** — bestiary progress ×1.5, unique pet unlock (Phase 02).

## 8. Tool

Ranged weapon (bow/crossbow) tercih; ancak herhangi weapon accept.

Ranged weapon +5% hunting XP passive bonus.

## 9. Skinning Cross-Profession

Hunting kill'i skinning için materyal üretir (bkz. [skinning.md](skinning.md)):
- Skinning profession aktif ise hunt kill'inden leather drop.
- Cross-profession sinerji.

## 10. AFK Aktivite

Combat: [activity-engine/hunting.md](../activity-engine/hunting.md).

Profession XP hunting aktivite yield'inde per kill hesap:
```
yield.XPProfession[hunting_id] += killsCount * 5
```

## 11. Cross-Reference

- [Skinning](skinning.md), [Mining](mining.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md)
- [game-design/hunting.md](../game-design/hunting.md)
- [race/racial-bonuses.md](../race/racial-bonuses.md) — orc
- [skill-engine/mastery.md](../skill-engine/mastery.md) — bestiary

## 12. Metrics
- `profession_hunting_kills_total`.
- `profession_hunting_bounty_completed_total` (Phase 02).

## 13. Bekleyen Kararlar
- Bounty contract system Phase 02.
- Rare mob tracking passive Phase 02.
- Hunter faction rep Phase 02.
- Beast taming (pet capture) Phase 03.
