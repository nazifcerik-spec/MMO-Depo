# Modifiers (Stack + Apply Order)

## 1. Amaç
Stat, damage, cost, cooldown gibi değerleri değiştiren modifier'ların formal modeli. Race passive, gear affix, buff, debuff, event bonus — hepsi aynı **Modifier** soyutlaması.

## 2. Modifier Yapısı
```go
type ModifierKind string
const (
    KindFlat        ModifierKind = "flat"        // v = base + amount
    KindMultAdd     ModifierKind = "mult_add"    // v = base * (1 + Σ amount)
    KindMultProd    ModifierKind = "mult_prod"   // v = base * Π (1 + amount)
    KindSet         ModifierKind = "set"         // v = amount (override)
    KindMin         ModifierKind = "min"         // v = max(base, amount)
    KindMax         ModifierKind = "max"         // v = min(base, amount)
)

type Modifier struct {
    Kind      ModifierKind
    Target    string      // "str", "crit_chance_pct", "damage_pct.fire", ...
    Amount    float64
    Scope     Scope
    SourceKind string     // race | class | spec | gear | affix | buff | debuff | event
    SourceID  uuid.UUID
    Priority  int32       // aynı Kind içinde sıralama
    ExpiresAt *time.Time  // null → kalıcı
    Stacks    int32       // stack count multiplier
}
```

## 3. Apply Order

Değer hesap sırası (formal):
```
1. base_value
2. flat_add (Σ KindFlat amount, priority artan)
3. mult_add — (1 + ΣA) tek çarpım (KindMultAdd)
4. mult_prod — Π (1 + Ai) her modifier ayrı çarpım (KindMultProd)
5. set — override (KindSet)
6. min/max clamp
7. hard cap (config-defined)
```

Örnek:
```
base_ap = 100
flat: +10, +5 → base += 15 → 115
mult_add: +10%, +5% → 115 * (1 + 0.15) = 132.25
mult_prod: +10%, +5% → 132.25 * 1.10 * 1.05 = 152.75
set: none
clamp: [0, ∞)
```

## 4. Additive vs Multiplicative Kararı
- **Flat** — sabit +N (stat point, gear +30 STR).
- **Mult_add** — genel %bonus'lar (race passive +%5 XP, event +%20 XP → toplanır, tek çarpım).
- **Mult_prod** — nadir; her biri bağımsız çarpım (elite mob damage taken -30% + boss damage taken -20% ayrı hesap).

Default kural: **mult_add tercih**. Mult_prod sadece istisnalar için.

## 5. Scope Kuralları

```go
type Scope struct {
    Always      bool
    InCombat    bool
    OutOfCombat bool
    Profession  string   // "mining", "" = all
    Activity    string   // "hunting", ""
    Element     string   // "fire", ""
    TargetKind  string   // "boss", "player", ""
    LocationTag string   // "forest", ""
}
```

Match logic: field non-empty ise, current context'e karşı check. All non-empty conditions match olmalı.

Örnek: dwarf passive "gathering_yield_pct +5 profession:mining":
- Aktivite `mining` ise apply.
- `foraging` ise skip.

## 6. Modifier Stack

Karakter için tüm aktif modifier'ları topla:
```go
func (e *ModifierEngine) Gather(ctx context.Context, charID uuid.UUID, sc StatContext) []Modifier {
    mods := []Modifier{}
    mods = append(mods, e.raceModifiers(ctx, charID)...)
    mods = append(mods, e.classModifiers(ctx, charID)...)
    mods = append(mods, e.specModifiers(ctx, charID)...)
    mods = append(mods, e.ascensionModifiers(ctx, charID)...)
    mods = append(mods, e.gearModifiers(ctx, charID)...)
    mods = append(mods, e.setBonusModifiers(ctx, charID)...)
    mods = append(mods, e.masteryModifiers(ctx, charID)...)
    mods = append(mods, e.buffModifiers(ctx, charID)...)
    return e.filterByScope(mods, sc)
}
```

Apply per target:
```go
func Apply(base float64, target string, mods []Modifier) float64 {
    v := base
    for _, m := range mods { if m.Target == target && m.Kind == KindFlat  { v += m.Amount * float64(m.Stacks) } }
    sumMultAdd := 0.0
    for _, m := range mods { if m.Target == target && m.Kind == KindMultAdd { sumMultAdd += m.Amount * float64(m.Stacks) } }
    v *= (1 + sumMultAdd)
    for _, m := range mods { if m.Target == target && m.Kind == KindMultProd { v *= (1 + m.Amount * float64(m.Stacks)) } }
    for _, m := range mods { if m.Target == target && m.Kind == KindSet { v = m.Amount } }  // last-writer-wins by priority
    return v
}
```

## 7. Stacking Rules

Aynı source aynı target'a birden fazla modifier eklenebilir mi?

Kural:
- **Aynı source_id + aynı target** = replace (refresh). Buff refresh senaryosu.
- **Farklı source_id + aynı target + aynı kind** = additive stack (yukarıdaki formul).
- **Farklı kind** = ayrı katman.

Buff stack cap:
- Aynı `SourceKind='buff'` grubunda max 3 stack (config).
- Aynı buff ID: `stacks` alanı; max stack config'de.

## 8. Budget

Modifier bütçe checklist:
- Sabit stat: race ±2, class +5, spec +3/+2, gear tier-based.
- Percent yield/damage: total additive %+50 civarı (out of scope: event, potion, rare buff).
- Multiplicative çarpım toplamı: cap 3.0×.

Content loader budget audit — deploy öncesi.

## 9. Persistence

- Race/class/spec modifier'lar: content data'dan direkt derive.
- Gear modifier'lar: `content.item_instance.rolled_stats` + affix + set.
- Buff modifier'lar: `character.stats_buff` tablosu.
- Ascension: `character.ascension_choice`.
- Mastery: `character.mastery` milestone-derived.

Buff table:
```sql
CREATE TABLE character.stats_buff (
  id UUID PRIMARY KEY,
  character_id UUID NOT NULL,
  buff_code VARCHAR(48) NOT NULL,
  modifiers JSONB NOT NULL,        -- Modifier list JSONB
  source_kind VARCHAR(24) NOT NULL,
  source_id UUID NULL,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NULL,
  stacks SMALLINT NOT NULL DEFAULT 1
);
CREATE INDEX ix_stats_buff__character_id__expires_at ON character.stats_buff (character_id, expires_at);
```

## 10. Expiration

- Read path: `WHERE expires_at IS NULL OR expires_at > now()`.
- Worker `buff_cleanup`: her dakika; expired row delete + cache invalidate.
- Event driven expire: skill effect explicitly calls `ExpireBuff(id)`.

## 11. Test Stratejisi

- Table-driven: (base, list of modifiers) → expected result.
- Order sensitivity: flat before mult; set overrides.
- Scope filter: irrelevant modifiers excluded.
- Stack limits.

## 12. Fail Modları
- Unknown target → warn; skip.
- Overflow (float64) → clamp.
- Corrupted modifier json → log + skip; buff row soft-marked invalid.

## 13. Metrics
- `modifier_gather_duration_seconds`.
- `modifier_apply_count{kind}`.
- `buff_active_count` gauge.

## 14. Cross-References
- [Stats](stats.md), [Skills](skills.md), [Passives](passives.md), [Active Skills](active-skills.md)
- [race/passive.md](../race/passive.md) — race modifier source
- [character/stats.md](../character/stats.md) — formul consumer

## 15. Bekleyen Kararlar
- Modifier priority ordering final — CI validation.
- Stack cap per source_kind final (buff 3 default).
- Complex conditional modifiers ("+%X while target hp < 20") — Phase 02 conditional DSL.
