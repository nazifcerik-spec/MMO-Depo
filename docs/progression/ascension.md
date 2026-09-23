# Ascension (Rebirth)

## 1. Amaç
Karakter level cap'e ulaştıktan sonra bilinçli "reset" yaparak kalıcı prestij puan kazanma sistemi. Yeni oyun döngüsü ve minor advantage.

## 2. Karakteristik

- **Trigger**: karakter max level (v1: 60) + ek koşul (aşağıda §3).
- **Reset kapsamı** (v1): level → 1, XP sıfırlanır, class skill puanları geri iade, ekipman tutulur.
- **Reset dışı** (kalıcı): mastery, achievements, reputation, collections, mail, guild üyeliği, gold, banka, marketplace listings.
- **Prestij puanı**: ascension başına 1-3 puan; toplam `character.ascension_points` account bazında kısmen paylaşılır (opt).
- **Prestij tier**: ascension count'a göre görsel tier + title.
- **Frequency hedef**: 60-100 saat aktif oyun / ascension (DD-OPEN-5).

## 3. Koşul (v1)

```yaml
ascension_requirements:
  min_level: 60
  min_playtime_hours: 60
  quest_line_complete: "path_of_ascension_v1"
  cost:
    gold_copper: 100000000       # 10.000 g
    item_required: [{ code: 'orb_of_ascension', quantity: 1 }]
```

**Path of Ascension** questi (v1): 5-quest ana line, ilk ascension için ~2-5 saat.

## 4. Prestij Puanı Kullanımı

Ascension puanı yeni bir tree'ye harcanır — "**Prestige Tree**":

```yaml
prestige_tree:
  categories: [combat, gathering, crafting, exploration, social]
  nodes_per_category: 8
  cost_per_node: 1
  respec_allowed: true            # cost 100k gold
```

Örnek node'lar:
- +1% XP all
- +2% gather yield
- +1% craft crit
- +5% mount speed
- +1 max friend cap
- unlock cosmetic aura (permanent glow)

**Cap**: prestige puan gündelik oyuna küçük etki (+1-2% aralığı); balance-safe (WoW paragon benzeri değil, AFK MMO idle progression).

## 5. Model

```sql
CREATE TABLE progression.ascension_event (
  id BIGSERIAL PRIMARY KEY,
  character_id UUID NOT NULL,
  ascension_number SMALLINT NOT NULL,       -- 1, 2, 3, ...
  triggered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  playtime_hours_before NUMERIC(10,2) NOT NULL,
  cost_paid_copper BIGINT NOT NULL,
  prestige_points_awarded SMALLINT NOT NULL,
  version INT DEFAULT 1
);

ALTER TABLE character.character
  ADD COLUMN ascension_count SMALLINT NOT NULL DEFAULT 0,
  ADD COLUMN prestige_points_total INT NOT NULL DEFAULT 0,
  ADD COLUMN prestige_points_available INT NOT NULL DEFAULT 0;

CREATE TABLE progression.prestige_node_choice (
  character_id UUID NOT NULL,
  node_code VARCHAR(48) NOT NULL,
  chosen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, node_code)
);
```

## 6. Akış

`POST /v1/character/{id}/ascend`

Tek SERIALIZABLE TX:
1. Koşul check (level, playtime, quest, cost).
2. Cost ödeme (ledger).
3. Level → 1, XP → 0.
4. Class skill puanları iade (bkz. [class/respec.md](../class/respec.md) partial).
5. `ascension_count++, prestige_points_available += n, prestige_points_total += n`.
6. `progression.ascension_event` insert.
7. Achievement trigger (`ascension_1`, `ascension_5`, `ascension_10`).
8. Title trigger (`ascended_i`, `ascended_v`).
9. WS anons + server broadcast (5. ascension ve üstü).
10. COMMIT.

## 7. Prestij Point Node Seçim

`POST /v1/character/{id}/prestige/node/{code}/choose`
- prestige_points_available ≥ 1.
- Node yok değil.
- Insert + `prestige_points_available--`.
- Modifier stack invalidate.

Respec:
`POST /v1/character/{id}/prestige/respec` — 100k gold; tüm node choice sıfırlanır; puan geri iade.

## 8. Frequency Target ve Balance

- 1. ascension ~60 saat oyun.
- 2. ascension +80 saat (kumulative 140h).
- 10. ascension ~600 saat.

Cost eğrisi ve puan/ascension Phase 02 balance pass'te ayarlanır (DD-OPEN-5). V1 taslak sabit.

## 9. Server-Wide Anons

- 1. ascension: kişisel + friend list.
- 5+ ascension: server "regional" broadcast.
- 10+ ascension: server "global" broadcast + achievement + title `Eternal`.

## 10. Endpoint

**GET `/v1/character/{id}/ascension`** — status + requirements.
**POST `/v1/character/{id}/ascend`**.
**GET `/v1/character/{id}/prestige/tree`**.
**POST `/v1/character/{id}/prestige/node/{code}/choose`**.
**POST `/v1/character/{id}/prestige/respec`**.

## 11. Metrikler
- `ascension_total{ascension_number}`
- `prestige_node_choose_total{node_code}`
- `prestige_respec_total`
- `ascension_playtime_hours_histogram`

## 12. Cross-Reference
- [character/level.md](../character/level.md), [character/xp.md](../character/xp.md), [class/respec.md](../class/respec.md)
- [achievements.md](achievements.md), [titles.md](titles.md)
- [mastery.md](mastery.md) — mastery ascension'da sıfırlanmaz

## 13. Bekleyen Kararlar
- Ascension puan cap total (Phase 02 balance).
- Account-wide prestige pool payı Phase 02 (alt karakter benefit).
- Ascension özel content (post-first-ascension quest hattı) Phase 02.
- Cosmetic aura tier Phase 02 art pipeline.
