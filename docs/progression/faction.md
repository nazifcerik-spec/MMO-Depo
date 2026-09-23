# Faction

## 1. Amaç
Karakterin ilişkilendiği/karşı olduğu organize NPC gruplarının taksonomisi. Reputation kaynağı + story branch + PvP dayanak (Phase 02+).

## 2. Karakteristik

- **NPC-driven** (v1): faction = NPC organizasyon. Oyuncu "faction seç" değil (WoW Horde/Alliance modeli yok v1); tüm oyuncular tüm faction ile rep işleyebilir.
- **Faction relationship graph**: faction A ↔ faction B ilişkisi (ally / neutral / enemy). Rival faction'a rep kazanmak diğerinden rep kaybettirir.
- **Story branching**: bazı ana quest hat faction seçim (yerel karar) — rep global.

## 3. Faction Türleri (v1 örnek)

| Kategori | Örnek |
|---|---|
| **City / Kingdom** | Ironhold Kingdom, Greenfield Village Council |
| **Guild (NPC)** | Miners' Consortium, Alchemy Guild |
| **Race-tied** | Elven Council, Dwarven Clan |
| **Religion / Order** | Sun Temple, Storm Cult |
| **Outlaw** | Bandit Coalition, Pirate Fleet |
| **Wildlife (opt)** | Wolf Pack (rare rep track for hunters) |

V1 hedef: ~15 faction template.

## 4. Model

```sql
CREATE TABLE content.faction (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  category VARCHAR(24) NOT NULL,
  home_location_id UUID NULL,
  is_hostile_by_default BOOLEAN NOT NULL DEFAULT false,
  starter_rep_xp INT NOT NULL DEFAULT 0,
  daily_cap_override INT NULL,             -- rep gain daily cap özel
  version INT DEFAULT 1
);

CREATE TABLE content.faction_relationship (
  faction_a_id UUID NOT NULL,
  faction_b_id UUID NOT NULL,
  kind VARCHAR(16) NOT NULL,               -- 'ally'|'neutral'|'enemy'
  rep_cross_effect_ratio NUMERIC(5,2) NOT NULL DEFAULT 0.0,
  -- ör. ally=+0.25 (A'ya kazan → B'ye +25% gain)
  -- enemy=-0.5 (A'ya kazan → B'den %50 kayıp)
  PRIMARY KEY (faction_a_id, faction_b_id),
  CHECK (faction_a_id <> faction_b_id)
);
```

## 5. Cross-Effect

Faction A'ya rep kazanma → faction B (ilişkiye göre) rep etkilenir:
```
delta_B = round(delta_A × rep_cross_effect_ratio(A→B))
```

Ör: Sun Temple → +100 rep. Sun Temple enemy Storm Cult (`ratio = -0.5`) → Storm Cult -50 rep.

Worker rep_gain applier faction ilişki grafiğini traverse eder (depth 1 v1; Phase 02 chain).

## 6. Hostile Faction Etkisi

`character_reputation.current_tier <= 1 (hostile)` → faction NPC saldırgan:
- Faction NPC'lerin görülür range içine girdiğinde threat kaydı (bkz. [combat/threat.md](../combat/threat.md)).
- Faction city içinde flag → guard NPC saldırısı.
- Faction vendor/quest access reddi.

## 7. Story Branch

Bazı quest hat "faction pick" sunar:
- "Yol A: Iron Kingdom'a yardım et → +500 Iron rep, -300 Bandit rep."
- "Yol B: Bandit'lere teslim et → +500 Bandit, -300 Iron."

Karar `progression.quest_choice_log`'a yazılır (bkz. [quests](../quests/README.md)). Retro değiştirilmez v1; Phase 02'de "faction reconciliation quest" ile silinip alternatif açılabilir.

## 8. Endpoint

**GET `/v1/faction`** — tüm faction katalog + karakter rep summary.
**GET `/v1/faction/{code}`** — detay + relationship graph + vendor + quest.

## 9. Metrikler
- `faction_active_count_gauge`
- `faction_relationship_traverse_latency_ms`
- `faction_hostile_gate_reject_total`

## 10. Cross-Reference
- [reputation.md](reputation.md), [character/reputation.md](../character/reputation.md)
- [quests](../quests/README.md), [economy/vendor.md](../economy/vendor.md)
- [combat/threat.md](../combat/threat.md)

## 11. Bekleyen Kararlar
- Faction pick sistem (WoW Horde/Alliance) — v1'de yok; Phase 02 story-driven soft pick (call event).
- PvP faction hostility (open-world PvP) Phase 02+.
- Faction leader NPC storyline Phase 02.
- Rep cross-effect chain derinlik (v1 depth 1; Phase 02 depth 2 opt).
- Faction switch mekaniği Phase 03.
