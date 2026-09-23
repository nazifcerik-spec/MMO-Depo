# Threat (Aggro System)

## 1. Amaç
NPC mob'ların hangi karakteri hedefleyeceğini belirleyen threat/aggro tablosu. Tank tanımlar; heal/dps balance.

## 2. Threat Table

Her hostile NPC bir threat table taşır:
```
mob.threat_table:
  { character_id: threat_amount }
```

Redis-based (`mob:threat:<mob_id>`) — live combat state.

Ephemeral; combat session bitince cleanup.

## 3. Threat Generation

Karakter aksiyonu threat üretir:
- **Damage dealt**: `threat = damage × threat_modifier`.
- **Heal dealt**: `threat = heal × heal_threat_ratio` (default 0.5).
- **Buff apply on ally**: `threat = 5 flat` (v1 sabit).
- **Aggro pull (approach)**: `threat = 1` (initial).
- **Taunt**: `threat = target's max threat + 1000` (guarantee top).

## 4. Threat Modifier

Character `threat_modifier`:
- Baz: 1.0.
- Tank spec (Guardian): +50% (`threat_modifier = 1.5`).
- Berserker spec: +10% (`1.1`).
- Rogue Assassin: -20% (`0.8` — reduce aggro).
- Cleric healing: default 1.0 (heal threat separate ratio).

Content-driven per spec.

## 5. Target Selection (Mob AI)

Her tick mob decides target:
```go
func selectTarget(mob):
    if mob.CurrentTarget == nil || mob.CurrentTarget.IsDead:
        // Pick highest threat
        top := mob.ThreatTable.HighestByAmount()
        return top
    // Threat swap: if higher threat > current × 1.1 (10% buffer)
    top := mob.ThreatTable.HighestByAmount()
    if top.Amount > mob.CurrentTarget.ThreatAmount × 1.1:
        return top
    return mob.CurrentTarget
```

## 6. Threat Decay

Combat inaktif periyodda threat decay:
- Combat sessions'da mob out-of-combat 5 sn → threat reset.
- Combat ongoing: threat kalıcı; decay yok (v1).

Phase 02: slow decay in combat (5% per 10 sec) tank rotation baskısı için.

## 7. Taunt Mechanic

Taunt skill:
```
apply:
  target_mob.ThreatTable[actor] = max(current) + 1000
  target_mob.CurrentTarget = actor
  target_mob.TauntedUntilMs = now + duration_ms
```

Taunt duration (3-5 sn):
- Mob forced to target; başka target değişemez.
- Post-taunt: normal threat table.

CC-like; DR: taunt category (Phase 02).

## 8. Aggro Pull (Initial)

Mob patrol area'da karakter yaklaşırsa (v1 no positional; alternatif: mob spawn ile combat contract'a katılım):
- Hunt aktivite başlatıldığında mob threat table karakter listelenir.
- Damage before hit → aggro pull chance (v1 explicit no; aktivite başlangıcı otomatik combat).

## 9. Multi-Target Aggro

AoE damage:
- Her hostile mob damage'la threat generate.
- Multiple mob combat state'ye girer.

## 10. Threat Redirect (Phase 02)

Cleric "Salvation" tarzı: threat generation halved for target 15 sn.

Guild Warden "Guardian Bond": teammate hasarının %30'u ally'e transfer + threat divert.

## 11. Config

```yaml
threat:
  heal_threat_ratio: 0.5
  buff_apply_flat_threat: 5
  aggro_pull_initial_threat: 1
  swap_buffer_pct: 10
  taunt_duration_default_ms: 4000
  taunt_forced_target: true
  decay_enabled: false            # v1 kapalı
  decay_pct_per_10s: 5             # Phase 02
```

## 12. AFK Aggregate

AFK reveal içinde threat tablosu **hesap yok** — kills aggregate; character solo default (party AFK dungeon party dynamics).

Party AFK dungeon: tank rolü hint snapshot'ta; loot distribution etkilenir.

## 13. Redis Schema

```
Key: mob:threat:<mob_instance_id>
Value: sorted set — member=character_id, score=threat_amount
TTL: 300 sec (5 min post combat)
```

Command usage:
```
ZADD mob:threat:xxx 152 char_id
ZINCRBY mob:threat:xxx 45 char_id
ZREVRANGE mob:threat:xxx 0 0 WITHSCORES  # highest
```

## 14. Event

`combat.threat_updated`:
```json
{
  "mob_id": "...",
  "character_id": "...",
  "delta": 152,
  "new_amount": 830
}
```

`combat.target_changed`:
```json
{
  "mob_id": "...",
  "old_target": "...",
  "new_target": "...",
  "reason": "threat_swap|taunt|target_dead"
}
```

## 15. Metrics
- `threat_generated_total`.
- `taunt_forced_total`.
- `threat_swap_total`.

## 16. Cross-Reference
- [Combat Loop](combat-loop.md), [Damage](damage.md), [Healing](healing.md), [PvE Encounters](pve-encounters.md)
- [character/attributes.md](../character/attributes.md)
- [game-design/combat.md](../game-design/combat.md)

## 17. Bekleyen Kararlar
- Threat decay Phase 02 aktivasyon.
- Taunt DR Phase 02.
- Ranged pull mechanic (Phase 02 positional).
- Boss-specific threat mechanic (fixation, resets) Phase 02+.
