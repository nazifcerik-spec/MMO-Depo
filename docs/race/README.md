# Race System (Faz 7)

Karakterin köken kimliği. Class'ı kilitlemez, tamamlayıcı bir mekanik katmandır. **Data-driven**: yeni race eklemek kod değişikliği gerektirmez — content pack + hot reload.

Yüksek seviye vizyon: [game-design/races.md](../game-design/races.md). Bu klasör Phase 7'de implementation-hazır detay.

## 1. Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış + data-driven prensipleri | [README.md](README.md) |
| [Race](race.md) | Entity model + creation + change + lifecycle |
| [Base Stats](base-stats.md) | Stat modifier data + budget rules |
| [Passive](passive.md) | Racial passive ability tanımı + effect DSL |
| [Racial Skills](racial-skills.md) | Aktif racial ability + cooldown + skill bar |
| [Racial Bonuses](racial-bonuses.md) | Profession affinity + combat/exploration bonuses + unique mechanic |

## 2. Data-Driven İlkeleri

1. **Kod race bilmez** — race code'ları (`human`, `dwarf`, ...) hiçbir Go switch/if'inde geçmez. Tüm davranış `content.race` + ilgili tablolara bakar.
2. **Effect DSL** — passive/active davranışlar JSON-encoded "effect list" ile ifade edilir; effect combiner (bkz. [game-design/skills.md](../game-design/skills.md)) çözer.
3. **Hot reload** — content patch: yeni race satırı + skill + bonus → server-side loader `content_version` bump; running server yeni versionu benimser (bkz. [database/postgresql.md §Content Data Loading](../architecture/postgresql.md)).
4. **Backward compatible** — mevcut karakterler snapshot ile korunur; nerf durumunda compensation token (DD-open).
5. **Testability** — race stat modifier + effect toplamı **saf fonksiyon**; golden test suite race sayısı arttıkça büyür.

## 3. Content Pack Format

Race pack tek YAML dosyası:

```yaml
# content/races/human.yaml
code: human
version: 1
display_name_key: race.human.name
lore_key: race.human.lore
art_ref: races/human/portrait.png

stat_modifiers:
  str: 1
  dex: 1
  int_: 1
  wis: 1
  sta: 1
  cha: 1

passive_ability: human_versatility
active_ability: human_rally

profession_affinity:
  profession_code: trading
  bonus_pct: 5

unique_mechanic:
  code: adaptable_xp_bonus
  params:
    xp_bonus_pct: 5

allowed_cosmetics:
  hair_colors: [brown, blonde, black, red, gray, white]
  skin_tones: [pale, tan, dark]
  face_ids: [1, 2, 3, 4, 5, 6, 7, 8]

starter_location_city_code: ashfield

tags: [balanced, humanoid, versatile]
```

Loader (content build/deploy):
1. Validate YAML → schema check.
2. Resolve references (ability codes → skill_template ids; city code → world.settlement id).
3. Upsert `content.race` row.
4. Bump `content.content_version`.
5. Publish event `content.race_updated`.

## 4. Runtime Yapısı

Server startup:
```go
racesCache := content.LoadAllRaces(ctx)   // in-memory read-through
```
Race'e erişim:
```go
r, err := content.Races.ByID(ctx, character.RaceID)
```
Cache invalidate: content_version bump event → in-process reload.

## 5. Tercih Edilen v1 Race Seti

Placeholder isimler (final lore Phase 01):

| Kod | Kısa Kimlik | Profession Affinity |
|---|---|---|
| `human` | Balanced/versatile | Trading |
| `dwarf` | Endurance/craft | Mining |
| `elf` | Speed/precision | Woodcutting |
| `orc` | Raw power | Hunting |
| `merfolk` | Coastal/agile | Fishing |
| `arcanborn` | Magic-affine | Enchanting |

Detaylı stat/passive/active: [base-stats.md](base-stats.md), [passive.md](passive.md), [racial-skills.md](racial-skills.md).

## 6. Ortak İlkeler

1. **Race class'ı kilitlemez** (DD-005).
2. **Stat delta bütçesi** ±5 civarında (bkz. [base-stats.md](base-stats.md)).
3. **Passive bonus küçük** (%5 civarı yield / %5 XP / vs.).
4. **Active ability class ability bütçesinden küçük** (cooldown 3-5 dk).
5. **Unique mechanic** kimlik verir ama meta-changer değil.

## 7. İlişkiler
- [game-design/races.md](../game-design/races.md), [game-design/classes.md](../game-design/classes.md)
- [character/character.md](../character/character.md), [character/attributes.md](../character/attributes.md)
- [database/content-templates.md §1](../database/content-templates.md)
- [systems/afk-engine.md](../systems/afk-engine.md) (racial passive yield hesabı)
- [backend/event-bus.md](../backend/event-bus.md) (content update event)

## 8. Ortak Bekleyen Kararlar
- Race change (Phase 02+) ücretli mi ücretsiz mi (nerf compensation).
- Sub-race / clan varyantları (Phase 02).
- Race-locked cosmetic exclusivity katmanı (Phase 02).
- Race compendium (koleksiyon: tüm race'lerle karakter tutmak) — Phase 03.
