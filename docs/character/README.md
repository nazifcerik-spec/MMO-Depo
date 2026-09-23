# Character Subsystem (Faz 6)

Karakterin oyun içi varlığının tüm mekanik boyutu. GDD [character.md](../game-design/character.md) yüksek seviye vizyonu; bu klasör Phase 6 implementation'a hazır formül + endpoint + edge-case detayı.

## Doküman Haritası

| Konu | Doküman |
|---|---|
| Genel bakış (bu doküman) | README.md |
| [Character](character.md) | Kimlik model + creation + rename + delete |
| [Level](level.md) | Level cap, level up, ascension bridge |
| [XP](xp.md) | XP kaynakları, formüller, curve, XP debt |
| [Attributes](attributes.md) | Primary stats (STR/DEX/INT/WIS/STA/CHA), allocation, respec |
| [Stats](stats.md) | Derived + tertiary stats, formulas, cap |
| [Health](health.md) | HP pool, regen, damage, death, revive |
| [Mana](mana.md) | Mana resource (mage/cleric/druid caster path) |
| [Stamina](stamina.md) | Endurance model (opt-in v1; formal design) |
| [Energy](energy.md) | Rogue-style hızlı regen kaynak |
| [Reputation](reputation.md) | Faction rep progression (Phase 02 baseline) |
| [Titles](titles.md) | Kazanılmış title'lar + active title |

## Resource System (Özet)

Her class'ın bir birincil kaynak türü var (bkz. [classes.md](../game-design/classes.md)):

| Class | Kaynak | Pool | Regen Modeli |
|---|---|---|---|
| Warrior | **Rage** | 0..100 | Alınan/verilen hasar üretir, hızlı boşalır |
| Cleric | **Mana** | Large | Doğrudan doldurma (regen + rest) |
| Mage | **Mana** | Very large | Doğrudan doldurma |
| Ranger | **Focus** | 0..100 | Zamanla artar, hızlı skill'ler tüketir |
| Rogue | **Energy** | 0..100 | Hızlı regen (10/sn) |
| Druid | **Nature** | 0..100 | Combat aktivitesi + zamanla |

**Universal secondary**: Health (HP), Stamina (opt-in).

Rage, Focus, Nature detayları bu klasördeki [stamina.md](stamina.md) ve [energy.md](energy.md) dokümanlarının kuzenidir; ancak burada listelenmediği için ilk pass'te ayrı doküman açılmadı. Class-specific detay [combat.md](../game-design/combat.md) ve [skills.md](../game-design/skills.md) altında Phase 6 sonu birleştirilecek.

## Ortak İlkeler

1. **Server-authoritative**: XP, level, HP, resource — hepsi server otoritesi.
2. **Optimistic locking**: `character.character.version` her mutation'da.
3. **Deterministic**: aynı input → aynı hesap (test edilebilir).
4. **Cap discipline**: her kaynağın hard cap'i tanımlı; overflow silent.
5. **Event-emitting**: her önemli değişim (level_up, death, allocate) event bus üzerinden yayınlanır (achievement/quest hook).
6. **Cross-platform**: aynı state Web + iOS + Android; UI cache local.

## Endpoint Özeti
```
GET    /v1/characters                       — list (account'un karakterleri)
POST   /v1/characters                       — create
GET    /v1/characters/{id}                  — profile (sheet)
DELETE /v1/characters/{id}                  — soft delete
POST   /v1/characters/{id}/rename           — rename (ücretli/limit)
POST   /v1/characters/{id}/enter            — session claim (aktif karakter)
POST   /v1/characters/{id}/exit             — release claim
POST   /v1/characters/{id}/allocate-stats   — stat point spend
POST   /v1/characters/{id}/respec-stats     — reset & re-allocate (cost)
POST   /v1/characters/{id}/respec-spec      — specialization respec
POST   /v1/characters/{id}/revive           — after death
POST   /v1/characters/{id}/set-title        — active title
GET    /v1/characters/{id}/reputation       — faction rep listing
GET    /v1/characters/{id}/titles           — owned titles
POST   /v1/characters/{id}/ascend           — ascension (separate flow)
```

## İlişkiler
- [game-design/character.md](../game-design/character.md) — vizyon
- [game-design/stats.md](../game-design/stats.md) — stat felsefesi
- [game-design/combat.md](../game-design/combat.md) — resource kullanımı
- [game-design/progression.md](../game-design/progression.md) — meta ilerleme ekseni
- [database/character.md](../database/character.md) — şema
- [systems/reward-system.md](../systems/reward-system.md) — XP grant
- [backend/event-bus.md](../backend/event-bus.md) — event yayını
