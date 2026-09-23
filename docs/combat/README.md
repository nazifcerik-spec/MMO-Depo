# 🔥 Combat Engine (Faz 17)

Karakterin ve düşmanların dövüşme motoru. Character + Stats + Class + Skill + Item hazır olduğuna göre bu motor hepsini birleştirir.

Combat, oyunun **loot doyumu + kimlik + risk** üçgeninin merkezi. Aktif oyun ve AFK reveal — **aynı motor**.

Foundation:
- [game-design/combat.md](../game-design/combat.md) — vizyon.
- [skill-engine/README.md](../skill-engine/README.md) — cast + damage compute.
- [item-engine/README.md](../item-engine/README.md) — gear stat contribution.
- [character/health.md](../character/health.md) — HP/regen/death.
- [activity-engine/hunting.md](../activity-engine/hunting.md) + [activity-engine/combat.md](../activity-engine/combat.md) — AFK entegrasyonu.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Combat Loop](combat-loop.md) | Server tick + scheduler + rotation eval |
| [Auto-Attack](auto-attack.md) | Weapon-based auto-attack |
| [Damage](damage.md) | Damage pipeline detay (armor/resist/DR/crit) |
| [Healing](healing.md) | Heal + shield + regen flow |
| [Status Effects](status-effects.md) | Buff/debuff/DoT/HoT/CC |
| [Threat](threat.md) | Aggro + threat modifier + taunt |
| [Death & Revive](death-revive.md) | Death sequence + respawn + XP debt |
| [PvE Encounters](pve-encounters.md) | Mob/boss AI + phases + world boss |
| [Combat Events](combat-events.md) | Event bus emissions + subscriber pattern |

## 2. Combat Modları

| Mod | Aktivite | Katılım |
|---|---|---|
| **Manuel PvE** | Manual cast + client input | Aktif oyun |
| **AFK Rotation** | Activity contract eval | AFK reveal |
| **Party Dungeon** | Instance + party synergy | Aktif or AFK-idle |
| **World Boss** | Server-wide event; contribution | Manuel katılım |
| **PvP** (Phase 02) | Duel/arena/open-world | Aktif |

Tümü **aynı motor**; farklı entry point.

## 3. Ana Interface

```go
package combatengine

type Engine interface {
    // Live tick — session-based combat (manuel + party live)
    Tick(ctx context.Context, deltaMs int64) (TickReport, error)

    // AFK aggregate compute — activity reveal içinden
    ComputeAFKYield(ctx context.Context, contract Contract, snapshot Snapshot, elapsedMs int64) (AFKYield, error)

    // Skill cast (manuel path)
    CastSkill(ctx context.Context, req CastRequest) (CastResult, error)

    // Damage/heal apply (skill effect resolver from)
    ApplyDamage(ctx context.Context, req DamageRequest) (DamageResult, error)
    ApplyHeal(ctx context.Context, req HealRequest) (HealResult, error)

    // Status effects
    ApplyStatusEffect(ctx context.Context, req ApplyStatusRequest) (StatusInstance, error)
    ExpireStatusEffect(ctx context.Context, statusID uuid.UUID) error

    // Threat
    AddThreat(ctx context.Context, actor Actor, target Target, amount int32) error
    Taunt(ctx context.Context, actor Actor, target Target, durationMs int32) error

    // Death
    RegisterDeath(ctx context.Context, req DeathRequest) error
    Revive(ctx context.Context, req ReviveRequest) error
}
```

## 4. Ortak İlkeler

1. **Server-authoritative** — damage, heal, RNG, cooldown hep server (bkz. DD-002, DD-066).
2. **Deterministic RNG** — seed-based; replay + audit (bkz. DD-007).
3. **Aynı motor manuel + AFK** — DD-066: engine combat'ın parçası değil, combat engine'i tüketir.
4. **Analitik AFK** — DD-084: per-tick simulate yok; combat aggregate.
5. **Event-emitting** — her önemli olay event bus'ta (bkz. [combat-events.md](combat-events.md)).

## 5. Combat Entry / Exit

Karakter **combat state**'e girer:
- Damage taken from hostile.
- Damage dealt to hostile.
- CC (crowd control) applied to.

Combat state:
- Redis flag `char:combat:<char>` TTL 5 sn.
- Her damage/attack event refresh.
- Combat regen ×0.2 (bkz. [character/stats.md §3.2](../character/stats.md)).
- Bazı aksiyon (mount, respec, city entry) combat'ta blokludur.

Combat exit:
- 5 sn hostile action yoksa → out-of-combat.
- Regen normal.

## 6. Tick Rate

Server tick: **250 ms** (v1; DD-target GDD §7).
- Manuel combat session: tick her 250ms.
- AFK reveal: aggregate over elapsed_ms.
- WebSocket push: her tick veya event triggered (throttled 10 Hz max).

## 7. Cross-Reference

- [game-design/combat.md](../game-design/combat.md) — vizyon
- [skill-engine/scaling.md](../skill-engine/scaling.md) — damage formula
- [skill-engine/active-skills.md](../skill-engine/active-skills.md) — cast pipeline
- [skill-engine/passives.md](../skill-engine/passives.md) — trigger events
- [skill-engine/cooldown.md](../skill-engine/cooldown.md) — CD manager
- [character/health.md](../character/health.md), [character/mana.md](../character/mana.md), [character/energy.md](../character/energy.md)
- [item-engine/stats.md](../item-engine/stats.md), [item-engine/unique-effect.md](../item-engine/unique-effect.md)
- [activity-engine/hunting.md](../activity-engine/hunting.md), [activity-engine/dungeon.md](../activity-engine/dungeon.md)
- [game-design/world-bosses.md](../game-design/world-bosses.md)

## 8. Bekleyen (Ortak) Kararlar

- Combat tick rate final (v1: 250ms).
- Positional gameplay v1 vs Phase 02 (v1: abstract; range check target validity level).
- PvP damage split (Phase 02).
- Combat replay UI (Phase 02+).
- AI advanced pattern (state machine) Phase 02.
