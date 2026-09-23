# Weather

> ⚠️ **v1'de DEVRE DIŞI**. Design hazır; aktivasyon Phase 02.

## 1. Amaç
Ortam koşullarının travel süresini etkilemesi. Yağmur/kar/fırtına gibi durumlar hız azaltıcı.

## 2. Neden Phase 02?

v1'de dünya statik; weather sistem karmaşıklık ekler:
- Server tarafında weather state (region-wide, timed).
- Content pack: weather transitions + severity.
- UI: weather display + travel warning.

v1'de skip; Phase 02'de rich open-world hissi için.

## 3. Model (Phase 02)

```sql
CREATE TABLE world.region_weather_state (
  region_id UUID PRIMARY KEY,
  current_weather VARCHAR(16) NOT NULL DEFAULT 'clear',
  severity SMALLINT NOT NULL DEFAULT 0,      -- 0..3
  since TIMESTAMPTZ NOT NULL DEFAULT now(),
  next_change_at TIMESTAMPTZ NOT NULL
);
```

Weather types:
- `clear` — no effect.
- `rain` — light penalty.
- `storm` — moderate penalty.
- `snow` — moderate penalty (winter regions).
- `blizzard` — heavy penalty.
- `fog` — light penalty; ancak stealth bonus (Phase 02).
- `sandstorm` — heavy (desert).

## 4. Weather Multiplier

```yaml
weather_multipliers:
  clear: 1.0
  rain: 0.95            # -5% speed
  fog: 0.95
  snow: 0.90            # -10%
  storm: 0.85           # -15%
  blizzard: 0.75        # -25%
  sandstorm: 0.75
```

Multiplier `<1.0` → speed düşer (travel süresi artar).

## 5. Rotation

Weather periodic:
- Worker scheduler her N dakikada region_weather_state güncelle (random state transition).
- Season'a göre bias (winter → snow/blizzard tercih).

## 6. Weather Interaction with Terrain

Kombine etki:
- Rain + forest = 0.95 × 1.30 = terrain penalty artar.
- Snow + mountain = ekstra severe.

Compute segment-based; her edge terrain + region weather look-up.

## 7. Race + Weather Interaction

Bazı race weather'a dirençli:
- Dwarf snow/blizzard'a immun (+%0 penalty).
- Merfolk rain +5% bonus (su seven).
- Config-based race weather resist table.

## 8. UI Notifications

- Travel estimate weather-aware.
- Weather değişimi mid-travel → notification "Storm hit; travel slower".

## 9. Contract Snapshot vs Dynamic

- **Snapshot** (v1 fallback): weather travel start anında dondurulur.
- **Dynamic** (Phase 02): reveal her seferinde current weather'a bakılır.

Phase 02 tercih: dynamic; ancak snapshot fallback.

## 10. Cross-Reference
- [Travel](travel.md), [Base Time](base-time.md), [Terrain](terrain.md), [Race](race.md)
- [world/region.md](../world/region.md)

## 11. Bekleyen Kararlar
- Weather Phase 02 aktivasyon plan.
- Weather content pack (region default weather profile).
- Weather-race resist table.
- Weather season correlation (spring/summer/fall/winter).
- Weather event opportunities (e.g. rare mob spawn during storm).
