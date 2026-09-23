# Skills (Mount Skills)

## 1. Amaç
Mount'un aktif/passive yetenekleri — travel bonus, sprint, terrain, combat (Phase 03).

## 2. Skill Türleri

- **Passive**: sürekli aktif (aura effect).
- **Active** (mount ability): character trigger (cooldown-based).

## 3. Örnek Mount Skills

### Universal Skills

| Code | Type | Effect | Cooldown |
|---|---|---|---|
| `mount_gallop` | active | +50% speed for 15 sec | 5 min |
| `mount_endurance` | passive | +20% stamina max | – |
| `mount_dash` | active | Instant burst forward (Phase 02 positional) | 30 sec |

### Species-Specific

| Species | Skill | Effect |
|---|---|---|
| Horse | `war_charge` (Phase 03) | Dismount + damage in cone |
| Wolf | `pack_instinct` | Party movement speed +5% aura |
| Bear | `carry_extra` | +5 cargo slots (temporary 30 min buff) |
| Dolphin | `wave_ride` | +100% aquatic speed |
| Turtle | `shell_shield` | Immune to ambush (Phase 02) |
| Griffon (Phase 02) | `airborne` | Cross-terrain (mountain bypass) |

## 4. Skill Structure

`content.skill_template` benzer:
- `type = 'mount_active' | 'mount_passive'`.
- `owner_kind = 'mount'`.

## 5. Skill Ranks

Level 10/20/30 → rank 2/3/4 unlock (matching pet skill pattern).

Each rank: +20% effect magnitude.

## 6. Active Skill Cast

Character triggers mount skill:
```
POST /v1/characters/{id}/mounts/{instance_id}/skill/{skill_code}/cast
```

Server:
- Validate mount summoned.
- Cooldown check (Redis).
- Stamina consume (opsiyonel).
- Apply effect.
- Emit event.

## 7. Passive Skill Apply

Mount aktif iken:
- Passive skill effect modifier stack'e ekle.
- Character speed_pct, cargo_pct, evasion_pct etkilenir.
- Cache invalidation on activate/deactivate.

## 8. Rare Signature Skill

Legendary mount unique ability (species-özel):
- Phoenix mount: "Rebirth" — mount HP 0 → auto-revive (Phase 03).
- Unicorn: "Purifying Aura" — cleanse debuff on party aura.
- Wyvern: "Fire Breath" — cone damage (Phase 03 mounted combat).

## 9. AFK Travel Skill

- Passive skills auto-active.
- Active skill (gallop) çoğunlukla manuel; AFK travel'da default gallop opsiyonu (Phase 02).

## 10. Metrics
- `mount_skill_cast_total{skill}`.
- `mount_skill_rank_up_total`.

## 11. Cross-Reference
- [Mount](mount.md), [Level](level.md), [Stats](stats.md), [Travel](travel.md)
- [skill-engine/README.md](../skill-engine/README.md)

## 12. Bekleyen Kararlar
- Mounted combat Phase 03 (attack while mounted).
- Sprint stamina drain final.
- Skill unlock tree (Phase 02).
- Legendary signature ability content depth.
