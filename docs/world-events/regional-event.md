# Regional Event

## 1. Amaç
Bölge-bazlı, periyodik olarak tekrarlayan orta ölçekli event. Örnek: **Wildlife Migration** (bir region'da rare hayvan yoğunluğu), **Harvest Festival** (herb yenilenme boost), **Storm Surge** (fishing bonus).

Boss/invasion kadar dramatik değil; günlük/haftalık rutin çeşnisi.

## 2. Karakteristik

- **Kapsam**: 1 region veya 1 territory.
- **Süre**: 1-3 saat.
- **Frekans**: haftada 3-6 kez region başına (rotation ile çakışmaz).
- **Buff / spawn density değişimi**: gather node hızı, monster spawn oranı, XP bonus, loot bonus.
- **AFK dostu**: mevcut aktiviteyi (mining/fishing/hunting) etkiler; ekstra karar gerekmez.
- **No fail state**: sadece window; kazan-kazan.

## 3. Model

```sql
CREATE TABLE content.regional_event_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  region_id UUID NULL,
  territory_id UUID NULL,
  duration_seconds INT NOT NULL DEFAULT 5400,     -- 1.5h
  activity_effects JSONB NOT NULL,                -- [{ kind:'mining', multiplier:1.5 }, ...]
  buff_code VARCHAR(48) NULL,                     -- opsiyonel karakter buff
  buff_duration_seconds INT NULL,
  version INT DEFAULT 1
);

CREATE TABLE progression.regional_event_spawn (
  id UUID PRIMARY KEY,
  template_id UUID NOT NULL,
  region_id UUID NULL,
  territory_id UUID NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'scheduled', -- scheduled|active|ended
  version INT DEFAULT 1
);
```

## 4. Activity Effects

Activity Engine event boyunca active event'i `context.regional_events[]` içine ekler. Reveal / hesap sırasında modifier stack (bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md)) tarafından uygulanır:

```yaml
effects:
  - kind: gathering
    profession: mining
    multiplier: 1.5           # +50% yield
  - kind: monster_spawn
    rate_multiplier: 1.3
  - kind: xp_gain
    activity: hunting
    multiplier: 1.25
  - kind: loot_bonus
    rarity_boost_permil: 20   # +2% rare chance
```

## 5. Buff (Opsiyonel)

Bazı event'ler karakter buff verir (region'a girer girmez):
- `festival_cheer`: +5% all XP, 30dk.
- `storm_focus`: +10% fishing speed, 45dk.

Bkz. [combat/status-effects.md](../combat/status-effects.md); buff region-conditional (region çıkışında iptal).

## 6. Örnek Katalog (Content — v1)

- `wildlife_migration_glenmoor` — Hunting +30%, rare mob rate +15%, 90dk.
- `harvest_festival_greenfield` — Farming/Herbalism +50%, 2h.
- `mining_boom_ironhold` — Mining +40%, ore quality +10%, 90dk.
- `storm_surge_seaport` — Fishing +50% rare fish, 60dk.
- `xp_carnival_capital` — XP +20% all sources, 3h (haftada bir).

## 7. Schedule

- Cron content-driven: `[template_code, region_id, day_of_week, hour_utc]`.
- Worker `regional_event_scheduler` 6h ahead insert.
- WS "regional_event_active" event notification (opt-in).

## 8. Endpoint

**GET `/v1/regional-event/active?region_id=`**
**GET `/v1/regional-event/upcoming?region_id=`** — sıradaki 48h.

Explicit "attend" yok — event window'da region'da aktivite yapan otomatik faydalanır.

## 9. Metrikler
- `regional_event_active_gauge{template}`
- `regional_event_activity_gain_ratio` — effect ölçüm

## 10. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md)
- [world/region.md](../world/region.md), [world/territory.md](../world/territory.md)
- [combat/status-effects.md](../combat/status-effects.md)

## 11. Bekleyen Kararlar
- Overlap policy — aynı region'da concurrent 2 event (v1: allowed, effects stackable via mult_add).
- WS opt-in default v1'de kapalı (spam koruma); Phase 02'de user prefs.
- Reveal-side effect apply: mid-activity başlarsa geriye dönük mü ileriden mi (v1: sadece event window içi tick'lere).
