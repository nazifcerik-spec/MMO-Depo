# Rarity

## 1. Amaç
Rarity band tanımı; drop distribution; affix budget; socket count; visual (renk + parıltı); ekonomi tier'ı.

## 2. Rarity Bandları

| Kod | Ad (i18n) | Renk (öneri) | Affix count | Socket | Set üye | Unique effect | Drop rate baz |
|---|---|---|---|---|---|---|---|
| `poor` | Bozuk | #9d9d9d gri | 0 | 0 | hayır | hayır | %5 |
| `common` | Sıradan | #ffffff beyaz | 0-1 | 0 | hayır | hayır | %60 |
| `uncommon` | Nadir | #1eff00 yeşil | 1-2 | 0 | rare | hayır | %25 |
| `rare` | Ender | #0070dd mavi | 2-3 | 0 | evet | hayır | %8 |
| `epic` | Destansı | #a335ee mor | 3-4 | 0-1 | evet | rare | %1.8 |
| `legendary` | Efsanevi | #ff8000 turuncu | 4 | 1-2 | evet | evet (required) | %0.18 |
| `mythic` | Mit | #e6cc80 kırmızı-altın | 4-5 | 2 | evet | evet + progression gated | %0.02 |

Toplam: %100.

Not: yield tablo `content.loot_table.rarity_distribution` override eder (boss/world-boss/chest için farklı).

## 3. Drop Distribution

Loot table content:
```yaml
# content/loot_tables/ashen_wolf.yaml
code: ashen_wolf_loot
source_type: mob
rarity_distribution:
  poor: 5
  common: 60
  uncommon: 25
  rare: 8
  epic: 1.8
  legendary: 0.18
  mythic: 0.02
meta:
  luck_scale: 1.0
```

Meta `luck_scale` gear/potion "magic find" ile çarpılır — rarity daha yüksek band'a düşme şansı artar.

Boss loot table örneği:
```yaml
code: ashen_wyrm_boss_loot
source_type: boss
rarity_distribution:
  rare: 40
  epic: 40
  legendary: 15
  mythic: 5
```

## 4. Rarity Roll Algoritması

```go
func rollRarity(dist RarityDist, luckMult float64, seed Seed) Rarity {
    // Adjust distribution by luck (shift higher rarities up)
    adjusted := adjustForLuck(dist, luckMult)
    // Roll 0..100
    r := seed.NextFloat() * 100
    cumulative := 0.0
    for _, band := range rarityBandsOrdered {  // sorted lowest → highest
        cumulative += adjusted[band]
        if r <= cumulative { return band }
    }
    return highestBand  // fallback
}
```

Luck adjustment (soft cap +100% luck):
```
each rarity above 'common' gets weight *= (1 + luck_pct/100 * band_bias)
poor + common decreased proportionally
```

## 5. Affix Budget

Rarity başına affix count range:
```
poor:       0 affix
common:     0-1 (roll: %70 no, %30 one)
uncommon:   1-2
rare:       2-3
epic:       3-4  (min 1 socket)
legendary:  4    (1-2 socket, unique effect required)
mythic:     4-5  (2 socket, unique effect + progression gate)
```

Affix distribution:
- Odd number toplam → prefix priority.
- Örnek rare (2-3 affix): min 1 prefix + 1 suffix; ek affix rasgele bir tarafa.

## 6. Rarity Effect Multipliers

Stat roll'unda rarity band multiplier:
```
poor:        0.6×
common:      0.8×
uncommon:    1.0×
rare:        1.1×
epic:        1.25×
legendary:   1.5×
mythic:      1.75×
```

Yani aynı `steel_longsword` template mythic olarak drop olursa `weapon_damage` %75 fazla.

## 7. Rarity Rendering (Client)
- Renk (yukarıdaki tablo).
- Loot beam / parıltı (rare+).
- Icon frame (rarity color).
- Nameplate tooltip glow.

## 8. Rarity Ekonomi Etkisi
- Vendor sell price ×rarity multiplier.
- Market bid asking rarity-based.
- Rarity band'ı içinde salvage material tier'ı farklı (Phase 02 crafting).

## 9. Bad Luck Protection (Phase 02)

Rare+ drop'ta pity system:
- Karakter için per rarity `dry_streak` counter.
- Streak > N → drop chance geometric artar.
- Drop olduğunda counter reset.

Config:
```yaml
bad_luck_protection:
  enabled: false          # v1'de kapalı
  legendary_pity_after: 500
  mythic_pity_after: 10000
```

## 10. Loot Roll Server Audit

Her drop:
```
audit.event_log kind='item.dropped'
payload={template_id, rarity, seed, luck_mult, source_kind, source_ref}
```

Rare+ drop → security_signal düşük şiddet (anomaly baseline).

## 11. Cross-Platform
Client tüm rarity metadata'yı content endpoint'inden alır; rarity string standard.

## 12. Metrics
- `item_dropped_total{rarity,source}`.
- `item_dropped_rarity_ratio` gauge.
- `luck_multiplier_active_gauge`.

## 13. Fail Modları
- Distribution toplam ≠ 100 → loader warn; normalize.
- Empty pool for band → fallback next-lower band.

## 14. Test Stratejisi
- Deterministic rarity: same seed → same rarity.
- Distribution over 100k rolls matches expected within ±%2.
- Luck scale effect.

## 15. Cross-References
- [Item Template](item-template.md), [Item Instance](item-instance.md), [Affix](affix.md), [Prefix](prefix.md), [Suffix](suffix.md), [Sockets](sockets.md), [Unique Effect](unique-effect.md)
- [systems/loot-system.md](../systems/loot-system.md), [systems/random-system.md](../systems/random-system.md)
- [game-design/items.md §6](../game-design/items.md)

## 16. Bekleyen Kararlar
- Rarity distribution final (playtest-derived tune).
- Bad luck protection v1 vs v2 aktivasyon.
- Rarity beyond mythic (v3+ prestige).
- Rarity name localization ton (fantasy vs teknik).
