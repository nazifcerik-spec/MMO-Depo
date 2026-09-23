# Combat Tests

## 1. Amaç
15-step damage pipeline, threat sistemi, CC diminishing, cooldown, ölüm/diriliş, party combat.

## 2. Kritik Invariant'lar

- **Damage pipeline sıra bağlı** — modifier stack apply order kesin (flat → mult_add → mult_prod → clamp).
- **Threat table Redis sorted set** — top threat = current target.
- **Cooldown deterministik** — server monotonic time; race olmaz.
- **CC diminishing 4-cascade** (DD-107) — 4. CC immune 30s.
- **Death penalty deterministik** — XP debt hesap.

## 3. Unit — Damage Pipeline (15 Step)

Table-driven test (bkz. [combat/damage.md](../combat/damage.md)):

- ✓ Step 1-15 her adımda beklenen değer.
- ✓ Base attack: 100 STR × formula.
- ✓ Weapon damage roll deterministik (seed).
- ✓ Skill multiplier apply.
- ✓ Critical roll (crit_chance %25 → seed 3/12).
- ✓ Armor mitigation (`dmg / (1 + armor/100)`).
- ✓ Resistance (fire, cold, magic per-type).
- ✓ Absorb shield.
- ✓ Reflect (opsiyonel).
- ✓ Overkill = damage capped at target HP.

Örnek:
```
Case: str=100, weapon 20-30, skill mult 1.5, crit=true, target armor=50, resist=%20
Expected: base=100, weapon_roll(seed=42)=25, base+w=125, skill=125*1.5=187.5,
          crit=×2=375, armor mitigation=375/(1+0.5)=250, resist=250*0.8=200
```

## 4. Unit — Threat System

- ✓ Damage → threat add to sorted set.
- ✓ Heal → threat halved (config).
- ✓ Taunt → self threat +50% of top.
- ✓ Top-of-set = current target.

## 5. Unit — CC Diminishing (DD-107)

- ✓ 1st CC: full duration.
- ✓ 2nd CC: %50 duration.
- ✓ 3rd CC: %25 duration.
- ✓ 4th CC: immune (0 duration + status effect "cc_immune" 30s).
- ✓ 30s window reset; sonra sırf sıfırdan.
- ✓ CC kategori bağımsız (stun vs sleep ayrı cascade).

## 6. Unit — Cooldown

- ✓ Skill cast → cooldown started (server time).
- ✓ Duplicate cast within CD → `E_ON_COOLDOWN`.
- ✓ CD reduction stat (haste) → CD * (1 - haste%).
- ✓ Global CD (1.5s) all abilities.

## 7. Unit — Death / Revive

- ✓ HP ≤ 0 → death event; state='dead'.
- ✓ XP debt = level_xp × 5%.
- ✓ Durability -10% all equipped.
- ✓ Revive at spot: 10dk timer opts, cost consumable or 30s wait.
- ✓ Revive at city: instant + no penalty above death penalty.

## 8. Integration — Combat Session

- Testcontainer.
- Start combat session, execute 10 skill casts, verify HP delta, threat sorted set, cooldown state.
- Enemy death → loot roll (bkz. loot.md).
- Session state 'completed' → audit entry.

## 9. Integration — Party Combat

- 3-member party, single boss.
- Verify: threat per-member independent; damage credit per-member; loot personal rule apply (DD-012).
- Combat log event stream — timeline correct.

## 10. E2E — Hunting → Death → Revive → Continue

- Login, start hunt encounter, take damage, die, revive at city, resume hunt.
- Verify: XP debt applied, durability drop, character state normal after revive.

## 11. Fuzz — Skill Parameters

- Random skill_code + target_id → server never panic; error kod uygun.

## 12. Regression

- `testdata/combat/fireball_lvl30.json` — deterministic damage sequence snapshot.

## 13. Coverage
- `domain/combat`: ≥ 90%.
- `domain/skill-engine`: ≥ 85%.

## 14. Cross-Reference
- [combat/damage.md](../combat/damage.md), [combat/threat.md](../combat/threat.md), [combat/death-revive.md](../combat/death-revive.md), [combat/status-effects.md](../combat/status-effects.md)
