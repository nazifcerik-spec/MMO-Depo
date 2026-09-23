# Mastery (Cross-System)

## 1. Amaç
Karakterin belirli bir alandaki derinliğini ölçen ve pasif bonus veren horizontal progression.

## 2. Mastery Track'leri

| Track Kategorisi | Örnekler | Kaynak |
|---|---|---|
| **Combat** | Warrior mastery, Fire mastery, Melee mastery | Enemy kill + skill kullanım |
| **Gathering** | Mining, Woodcutting, Fishing, Herbalism, Hunting | Gather event |
| **Crafting** | Blacksmith, Alchemy, Cooking, Leatherworking | Craft complete |
| **Exploration** | Region discovery, Territory scout | Travel/scout complete |
| **Social** | Party leader, Guild contribution, Trade volume | Social event |
| **Meta** | Time played, Login streak | Login/session |

## 3. Model

```sql
CREATE TABLE content.mastery_track (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  category VARCHAR(24) NOT NULL,            -- combat|gathering|crafting|exploration|social|meta
  max_tier SMALLINT NOT NULL DEFAULT 10,
  tier_xp_curve JSONB NOT NULL,             -- [xp_to_tier_1, xp_to_tier_2, ...]
  tier_bonus_definitions JSONB NOT NULL,    -- [{tier: N, bonus_effects: [...]}, ...]
  source_kinds JSONB NOT NULL,              -- ['mine_ore','craft_complete',...]
  version INT DEFAULT 1
);

CREATE TABLE progression.character_mastery (
  character_id UUID NOT NULL,
  mastery_track_id UUID NOT NULL,
  xp BIGINT NOT NULL DEFAULT 0,
  current_tier SMALLINT NOT NULL DEFAULT 0,
  last_gain_at TIMESTAMPTZ NULL,
  version INT NOT NULL DEFAULT 1,
  PRIMARY KEY (character_id, mastery_track_id)
);
```

## 4. XP Curve

Tier XP `tier_xp_curve` content'de tanımlı. Örnek:
```yaml
tier_xp_curve:
  - 100     # tier 1
  - 500
  - 2000
  - 8000
  - 30000
  - 100000
  - 300000
  - 800000
  - 2000000
  - 5000000
```

Toplam tier 1→10 örnekte ~8.2M mastery XP.

## 5. Bonus Türleri

Bkz. [skill-engine/modifiers.md](../skill-engine/modifiers.md) modifier stack:

```yaml
tier_bonus_definitions:
  - tier: 1
    bonus_effects:
      - { kind: 'flat', target: 'gather_yield_mining', amount: 1 }
  - tier: 3
    bonus_effects:
      - { kind: 'mult_add', target: 'gather_speed_mining', amount: 0.05 }
      - { kind: 'flat', target: 'ore_rare_chance_permil', amount: 5 }
  - tier: 5
    bonus_effects:
      - { kind: 'unlock_recipe', recipe_code: 'ingot_masterwork_iron' }
  - tier: 10
    bonus_effects:
      - { kind: 'mult_add', target: 'gather_yield_mining', amount: 0.20 }
      - { kind: 'title_unlock', title_code: 'ironmaster' }
      - { kind: 'cosmetic_unlock', sku: 'pickaxe_prestige_iron' }
```

Bonus sistem modifier stack'ine ekleniyor — passive bonus (cooldown yok, aktif tetik yok).

## 6. XP Kaynağı — Aktivite Entegrasyonu

Activity Engine reveal sırasında contract complete/tick olay yayınlar:
```
event.progression.mastery_gain
  { character_id, track_code, xp_delta, source_ref_id, source_kind }
```

Progression worker `mastery_apply` batch tüketir:
1. `character_mastery.xp += delta`.
2. Tier threshold aşıldı mı? Aşıldıysa `current_tier++` + bonus_effects trigger (unlock event).
3. Modifier stack cache invalidate → karakter reload.
4. Notify (WS "mastery_tier_up") + achievement trigger (bkz. [achievements.md](achievements.md)).

## 7. AFK Uyumluluk

- Mining/fishing/crafting AFK ticks → doğal mastery XP.
- Combat mastery AFK'de PvE hunting/dungeon AFK contract reveal içinde XP hesap.
- Meta track (time played) her session tick += 1 (ephemeral); daily consolidate.

## 8. Cap ve Reset

- Tier 10 hard cap v1 (Phase 02'de "mastery paragon" tier 11-20 opt).
- Ascension mastery XP'yi **sıfırlamaz** (yatay progression korunur).
- Ekstra respec / reset yok — mastery yönü seçim değil, aktivite yansıması.

## 9. Endpoint

**GET `/v1/character/{id}/mastery`** — tüm track özeti.
**GET `/v1/character/{id}/mastery/{track_code}`** — detay + next tier XP.

## 10. Metrikler
- `mastery_xp_gain_total{track}`
- `mastery_tier_up_total{track, tier}`
- `mastery_apply_batch_latency_ms`

## 11. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)
- [skill-engine/modifiers.md](../skill-engine/modifiers.md), [skill-engine/mastery.md](../skill-engine/mastery.md) — skill-özgü mastery
- [professions](../professions/README.md), [crafting](../crafting/README.md)
- [achievements.md](achievements.md), [titles.md](titles.md)

## 12. Bekleyen Kararlar
- Mastery paragon (tier 11+) Phase 02.
- Track sayısı content-driven; v1 hedef ~25 track.
- Cross-mastery synergy (ör. mining + smithing) Phase 02.
- Mastery loss (ascension opt) — v1'de yok.
