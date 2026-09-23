# Base Stats (Race Stat Modifiers)

## 1. Amaç
Her race'in altı primary stat üzerindeki delta'sı. Kimlik ve mekanik farkın **düşük etkili** ama **hissedilebilir** ifadesi.

## 2. Bütçe Kuralları
- **Net toplam**: +3 stat (5 pozitif + 2 negatif = +3 net) hedef; ±%15 sapma tolere.
- **Tek stat max**: +3 (bir race tek başına dev bir stat'a sahip olamaz).
- **Tek stat min**: -2 (dezavantajlar hafif).
- **Sıfır stat**: kabul (Human "her şeyden +1" tarzı denge tercihi).
- **Toplam mutlak değer** ≤ 8 (kimlik yeterince belirgin).

## 3. Effective Formulle Etkileşim
Primary stat hesabı (bkz. [character/attributes.md §5](../character/attributes.md)):
```
effective[i] = character.stats.base[i]  (=level allocations)
             + race.stat_modifiers[i]
             + class.base_stat_bias[i]
             + ...
```

Bu doküman **race.stat_modifiers[i]** kısmına odaklanır.

## 4. v1 Race Setinin Modifier'ları (Placeholder)

| Race | STR | DEX | INT | WIS | STA | CHA | Net |
|---|---|---|---|---|---|---|---|
| `human` | +1 | +1 | +1 | +1 | +1 | +1 | **+6** |
| `dwarf` | +1 | -1 | 0 | +1 | +2 | 0 | **+3** |
| `elf` | -1 | +2 | +1 | +1 | 0 | 0 | **+3** |
| `orc` | +2 | 0 | -1 | 0 | +1 | -1 | **+1** |
| `merfolk` | -1 | +2 | +1 | +1 | 0 | 0 | **+3** |
| `arcanborn` | -1 | 0 | +2 | +1 | -1 | +1 | **+2** |

Not: Human "her şeyden +1 = net +6" görünse de tüm stat'ları dağıtıp bir yerde uzman değil; "generalist bonus" olarak dengeli.

Elf ve Merfolk aynı modifier'a sahip ancak passive/active/affinity ile ayrışır.

## 5. Race + Class Kombine Bütçe

Race modifier + class base bias toplam bütçesi ±%5 tutulur. Örnek:
```
human + warrior:
  race: STR+1, DEX+1, INT+1, WIS+1, STA+1, CHA+1  (+6 net)
  class: STR+3, STA+2 (+5 net)
  combined: STR+4, DEX+1, INT+1, WIS+1, STA+3, CHA+1  (+11 net)

orc + rogue:
  race: STR+2, STA+1, CHA-1, INT-1  (+1 net)
  class: DEX+3, STR+1 (+4 net)
  combined: STR+3, DEX+3, INT-1, STA+1, CHA-1  (+5 net)
```

Ölçüm: net toplam sapma 6-11 arası; oyuncu için "yanlış race seçtim" hissi minimal.

## 6. Data Format

`content.race.stat_modifiers` JSONB:
```json
{ "str": 1, "dex": 1, "int_": 1, "wis": 1, "sta": 1, "cha": 1 }
```

`int` reserved keyword olduğu için trailing underscore (`int_`). Loader normalize eder.

## 7. Content Loader Validation

Loader kuralları (build time — CI check):
- Tüm 6 anahtar var (eksik = 0 kabul edilir; ancak explicit yaz).
- Her değer -3..+3 integer.
- Net toplam: -2..+6.
- Toplam abs value ≤ 10.
- Warning: net toplam > +4 veya < 0 (unusual).

## 8. Nerf/Buff Yönetimi

Race stat modifier değişimi:
1. Content patch draft (yaml).
2. Simülasyon: mevcut karakterlerin effective stat delta'sı.
3. Impact assessment: delta > ±3 stat toplamda ise **compensation gerekli**.
4. Rollout:
   - Content_version bump.
   - Publish `content.race_updated` event.
   - Worker: mevcut karakter için compensating stat token (Phase 02) veya free respec.
   - Notification kullanıcıya.

## 9. Race + Ascension Etkileşimi

Ascension level artışı race modifier'ı **değiştirmez**. Ascension bonuses ayrı katmandan gelir. Bir Human ascension 5'te de Human'ın +1 STR/DEX/... bonusunu alır; artı ascension'dan gelen ek stat.

## 10. Endpoint Yansıması

Character sheet response (`GET /v1/characters/{id}/stats`):
```json
{
  "primary": {
    "str": 45,
    "components": {
      "base_allocations": 40,
      "race_bonus": 1,
      "class_bias": 3,
      "spec_bonus": 0,
      "ascension": 1,
      "gear": 0,
      "buffs": 0
    }
  },
  ...
}
```

Client debug UI için component breakdown gösterilir (Phase 02 opsiyonel).

## 11. Determinism

Race stat modifier saf integer; formul deterministic; test edilebilir.

## 12. Test Stratejisi

- Table-driven: her race için modifier + expected effective stat.
- Race + Class kombinasyonu 36 case (6×6) — golden test.
- Nerf senaryosu: modifier change → compensation trigger.

## 13. Metrics
- `race_stat_modifier_impact` histogram — nadir; nerf öncesi/sonrası kıyas.

## 14. İlişkiler
- [Race](race.md), [Passive](passive.md), [character/attributes.md](../character/attributes.md), [character/stats.md](../character/stats.md), [game-design/races.md](../game-design/races.md)

## 15. Bekleyen Kararlar
- Race modifier magnitude final: ±3 vs ±5 (DD-OPEN-4).
- Human "generalist +6" özel case korunur mu revize edilir mi (aşırı avantajlı görülebilir).
- Level scaling opsiyonu: race passive stat modifier level ile scale etsin mi (bkz. game-design/races.md §9).
