# Energy

## 1. Amaç
Rogue class'ının birincil kaynağı. Hızlı regen, düşük cap, ritmik skill kullanımı.

## 2. Kim Kullanır
- **Rogue** — birincil.
- (Diğer class'ların energy'si yok; Warrior=Rage, Ranger=Focus, Mage/Cleric=Mana, Druid=Nature.)

## 3. Model
- **Max Energy**: sabit 100 (v1); Phase 02 spec/passive ile 120'ye çıkabilir.
- **Current Energy**: runtime state (Redis + snapshot).
- **Regen**: hızlı, sürekli — 10 energy/sn baseline.

## 4. Formüller

### 4.1 Max
```
max_energy = 100 + spec_bonus + gear_bonus  // v1'de 100
```

### 4.2 Regen
```
energy_regen_per_sec = 10 + gear_haste * 0.1 + spec_bonus
```
Hızlı; her saniye ~10 → cap dolumu 10 sn.

Combat/idle regen aynı (energy her zaman regen — burst kaynak).

## 5. Skill Kullanımı
- Rogue skill maliyetleri düşük ama sık: 25-40 energy.
- Combo point interaction (Phase 02 spec detay).

## 6. Rotation Interaction
Rogue AFK rotation örneği:
```json
[
  { "cond": { "energy_at_least": 40, "combo_less_than": 5 }, "skill_id": "sinister_strike" },
  { "cond": { "combo_at_least": 5 }, "skill_id": "eviscerate" }
]
```

## 7. Endpoint
`GET /v1/characters/{id}/vitals` → resource block:
```json
{
  "resource": {
    "kind": "energy",
    "current": 65,
    "max": 100,
    "regen_per_sec": 10.0
  }
}
```

## 8. Config
```yaml
energy:
  base_max: 100
  base_regen_per_sec: 10.0
  haste_regen_bonus: 0.1  # per haste %
```

## 9. Buff/Item
- Adrenaline Rush (spec skill, Phase 02): +100% regen 15 sn.
- Energy Elixir (rare consumable): instant +50.

## 10. Related: Rage / Focus / Nature (Overview)

Bu üç kaynak da resource_type; ayrı doküman şu an yok (Phase 6 listesinde yok). Kısa özet:

| Kaynak | Class | Pool | Regen Modeli |
|---|---|---|---|
| **Rage** | Warrior | 0..100 | Alınan/verilen hasar üretir (~5 rage per damage/hit); out-of-combat hızlı boşalır (2/sn drain) |
| **Focus** | Ranger | 0..100 | Zaman ile artar (5/sn); hızlı ranged shots tüketir; steady kaynak |
| **Nature** | Druid | 0..100 | Combat kill / heal +N; time regen 2/sn; form-specific tüketim |

Detay Phase 02 spec dokümanları eklenirken açılacak.

## 11. Fail Modları
- Insufficient energy → skill reject.
- Overflow: clamp to max.

## 12. WS Push
`vitals_update` frame içinde resource delta.

## 13. Metrics
- `character_resource_full_time_ratio{kind}` — resource'un cap'te olduğu zaman oranı (rotation etkinlik göstergesi).

## 14. Test Stratejisi
- Regen tick over time.
- Skill cost deduction.
- Combo integration (Phase 02).

## 15. İlişkiler
- [Mana](mana.md), [Stats](stats.md), [game-design/classes.md](../game-design/classes.md), [game-design/skills.md](../game-design/skills.md), [game-design/combat.md](../game-design/combat.md)

## 16. Bekleyen Kararlar
- Rage/Focus/Nature ayrı doküman açılması (Phase 02).
- Combo point sistem detay (Rogue spec Phase 02).
- Cross-class resource hybrid (Druid form-switching Phase 02).
