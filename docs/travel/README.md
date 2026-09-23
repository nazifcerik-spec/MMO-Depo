# Travel (Faz 14)

Konum A → Konum B geçişi. AFK aktivite tipidir; süre karakterin **efektif hız katmanları** üzerinden hesaplanır.

```
Location A
    ↓
  Travel (base × mount × race × weather × terrain × equipment)
    ↓
Location B
```

Foundation: [activity-engine/travel.md](../activity-engine/travel.md) (activity contract entegrasyonu). Bu klasör Phase 14 implementation-hazır deep dive.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | [README.md](README.md) |
| [Travel](travel.md) | Travel entity + lifecycle + endpoint |
| [Distance](distance.md) | Mesafe hesabı + route model |
| [Base Time](base-time.md) | Base walking time formula |
| [Mount](mount.md) | Mount speed katkısı |
| [Race](race.md) | Racial speed bonus |
| [Weather](weather.md) | Weather modifier (Phase 02 design ready) |
| [Terrain](terrain.md) | Terrain type modifier |
| [Equipment](equipment.md) | Boots + barding + accessory katkısı |
| [AFK Integration](afk-integration.md) | Offline continuation + reveal integration |

## 2. Kritik İlkeler

1. **AFK-first**: travel doğuştan AFK aktivite; oyuncu logout olsa da devam eder.
2. **Analitik hesap**: süre formülle çözülür, per-tick simulate YOK (DD-084).
3. **Server-authoritative**: hız katmanları server tarafında; client asla.
4. **Deterministic**: aynı input + koşullar → aynı süre.
5. **Interruptable**: cancel edilebilir (partial refund opsiyonu Phase 02).

## 3. Efektif Süre Formülü

```
duration_ms = (distance_units / base_walking_speed) × 1000
              × (1 / mount_multiplier)
              × (1 / race_multiplier)
              × terrain_multiplier
              × weather_multiplier                (Phase 02; v1 = 1.0)
              × (1 / equipment_multiplier)
              × (1 / road_bonus_multiplier)
```

Detay her katman kendi dokümanında.

## 4. Travel Methods (Bkz. Base Time)

Method:
- `walking` — base.
- `mount` — mount speed.
- `road` — road bonus (hız çarpımı).
- `wagon` — fixed schedule + cost.
- `ship` — fixed schedule (cross-continent).
- `waygate` — instant/cast time; gold cost sink.

`walking + mount + road + equipment` çarpımsal; `wagon/ship/waygate` fixed durations (yukarıdaki formul uygulanmaz).

## 5. Contract Model

`activity.contract` `activity_type='travel'`:
- `channel='primary'` — travel primary channel işgal eder (v1'de basit; Phase 02 mount_travel channel opsiyonu).
- `params_json`:
  ```json
  {
    "destination_location_id": "...",
    "method": "mount",
    "route_json": [...],
    "cost_paid_copper": 0,
    "expected_duration_ms": 720000,
    "expected_arrival_at": "..."
  }
  ```

## 6. Endpoint Genel

```
POST /v1/characters/{id}/travel/estimate     # süre + maliyet preview
POST /v1/characters/{id}/travel              # travel start (activity contract create)
POST /v1/characters/{id}/travel/{id}/cancel  # cancel + partial handling
GET  /v1/characters/{id}/travel/current      # aktif travel state
```

Detay: [travel.md](travel.md).

## 7. Reveal / Progress

- Aktif travel contract reveal her login'de:
  - elapsed >= duration → arrive.
  - elapsed < duration → still in progress; yield yok.
- Travel yield: sadece arrival'da discovery XP (küçük).

Detay: [afk-integration.md](afk-integration.md).

## 8. Cross-Reference

- [game-design/travel.md](../game-design/travel.md), [game-design/mounts.md](../game-design/mounts.md), [game-design/world.md](../game-design/world.md), [game-design/map.md](../game-design/map.md)
- [activity-engine/travel.md](../activity-engine/travel.md), [activity-engine/lifecycle.md](../activity-engine/lifecycle.md)
- [world/location.md](../world/location.md), [world/city.md](../world/city.md)
- [character/character.md](../character/character.md) — current_location_id
- [race/racial-bonuses.md](../race/racial-bonuses.md) — race speed bonus
- [inventory/equipment.md](../inventory/equipment.md) — boots slot

## 9. Config

```yaml
# config.travel.yaml
base_walking_speed_units_per_sec: 1.0
road_bonus_multiplier: 1.15
waygate_cast_time_ms: 5000
travel_route_max_distance_units: 10000

cancel_policy:
  partial_yield_pct: 0                # v1: cancel = origin geri; hiç ilerleme yok
  cost_refund_pct: 0                  # waygate/ship cost non-refundable

weather_enabled: false                # v1 kapalı; Phase 02
ambush_enabled: false                 # Phase 02

race_speed_bonuses:
  elf: 1.05
  merfolk: 1.10  # scope=aquatic
  # ...
```

## 10. Bekleyen Kararlar (Ortak)

- Cancel refund politikası final (v1: no refund).
- Multi-leg pathfinding Phase 02.
- Weather aktive Phase 02.
- Ambush event system Phase 02.
- Fast-travel gold sink balance final.
