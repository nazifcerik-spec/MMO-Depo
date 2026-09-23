# Crafting Mastery

## 1. Amaç
Belirli crafting profession'da uzun soluklu mastery track. Milestone'lar quality/critical/craft time bonuses + unique unlock.

Foundation: [skill-engine/mastery.md](../skill-engine/mastery.md), [game-design/mastery.md](../game-design/mastery.md).

## 2. Mastery Track Kataloğu (Crafting)

Her crafting profession için 1 mastery track:
- `blacksmithing_mastery`
- `leatherworking_mastery`
- `tailoring_mastery`
- `alchemy_mastery`
- `cooking_mastery`
- `enchanting_mastery`
- `jewelcrafting_mastery`

`content.mastery_track.category='profession'` filter.

## 3. XP Kaynağı

Mastery XP crafting yield içinde:
```
mastery_xp_per_iteration = recipe.xp_reward × 0.3   (v1 önerisi)
```

Success iteration = full; fail = %50.

Anti-cheat: per-hour cap (bot koruma; bkz. [skill-engine/mastery.md §11](../skill-engine/mastery.md)).

## 4. Rank Tier'lar

| Level | Tier | Bonus |
|---|---|---|
| 1-19 | Novice | +0% |
| 20-39 | Apprentice | quality exceptional +2%, craft time -5% |
| 40-59 | Journeyman | +4%, -10% |
| 60-79 | Expert | +6%, -15%, critical +3% |
| 80-99 | Master | +8%, -20%, critical +6% |
| 100 | Grandmaster | +10%, -25%, critical +10%, unlock unique recipes |

## 5. Milestone Rewards

Milestone (level 25, 50, 75, 100):
```yaml
milestones:
  - level: 25
    reward:
      modifiers:
        - { kind: flat, target: craft_exceptional_pct, amount: 3 }
      title_code: apprentice_smith
  - level: 50
    reward:
      modifiers:
        - { kind: mult_add, target: craft_time_reduction_pct, amount: 10 }
      unlock_recipe: masterwork_steel_longsword
      title_code: journeyman_smith
  - level: 75
    reward:
      modifiers:
        - { kind: flat, target: craft_critical_pct, amount: 5 }
      unlock_recipe: enchanted_steel_longsword
  - level: 100
    reward:
      modifiers:
        - { kind: mult_add, target: craft_time_reduction_pct, amount: 25 }
        - { kind: flat, target: craft_critical_pct, amount: 10 }
      title_code: grandmaster_smith
      unlock_recipe_group: grandmaster_smith_signature
```

## 6. Grandmaster Signature

Level 100 unlock: signature recipes.
- Blacksmithing GM: "Sunforged Blade" — unique weapon recipe.
- Alchemy GM: "Phoenix Elixir" — legendary consumable.
- Enchanting GM: "Runic Awakening" — unique enchant.
- Jewelcrafting GM: "Star Gem Cutting" — legendary gem tier.
- Cooking GM: "Feast of the Ages" — party-wide long-duration food buff.
- Leatherworking GM: "Dragonhide Cloak" — legendary armor.
- Tailoring GM: "Weaver of Fates" — legendary robe.

Signature recipes BoP; personal prestij.

## 7. Modifier Stack Integration

Mastery milestone modifier'ları modifier stack'e katkı verir:
- Character stat compute'da apply (bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md)).
- Craft success, quality, critical calculation'da mastery bonusları kullanılır.

## 8. Cross-Profession Mastery

- Karakter tek anda 3 aktif profession (v1) — 3 mastery track paralel.
- Track başına ilerleme bağımsız.
- Cross-track bonus (Phase 02): "renaissance crafter" achievement — 3 GM = title.

## 9. Endpoint

Mastery endpoint zaten var (bkz. [skill-engine/mastery.md §10](../skill-engine/mastery.md)):
- `GET /v1/characters/{id}/mastery` → tüm track.
- `POST /v1/characters/{id}/mastery/{track_code}/claim-milestone` — milestone claim.

## 10. Milestone Claim

Milestone eligible ise auto-notification. Claim:
- Reward apply.
- Cache invalidate.
- Audit event.

Auto-claim vs opt-in tercih: **opt-in** (Phase 02 auto opsiyonu).

## 11. Metrics
- `crafting_mastery_leveled_up_total{profession,level}`.
- `crafting_mastery_milestone_claimed_total{profession,level}`.
- `crafting_grandmaster_count{profession}` gauge.

## 12. Cross-Reference
- [Recipes](recipes.md), [Quality](quality.md), [Critical Craft](critical-craft.md), [Crafting XP](crafting-xp.md)
- [skill-engine/mastery.md](../skill-engine/mastery.md)
- [game-design/mastery.md](../game-design/mastery.md)
- [character/attributes.md](../character/attributes.md) — WIS

## 13. Bekleyen Kararlar
- Mastery XP oranı final (0.3 × recipe_xp v1).
- Recipe rank up (Phase 02) — mastery bağlı.
- Cross-profession GM bonus (renaissance).
- Guild research profession bonus stack Phase 02.
