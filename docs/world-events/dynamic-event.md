# Dynamic Event

## 1. Amaç
Küçük ölçekli, **emergent** olarak tetiklenen kısa event. Örnek: "Bandit ambush", "Caravan rescue", "Elemental rift". Yaşayan dünya hissi.

Online player katılımına ağır bağımlı; AFK sınırlı katılım (kaçırılırsa hafif teselli).

## 2. Karakteristik

- **Kapsam**: 1 location veya location cluster.
- **Süre**: 15-60 dk.
- **Trigger**:
  - **Cron rotation** (v1 default) — content-driven random rotation.
  - **Chain trigger** (v1 basit): bir event başarılı biterse ilgili event unlock (örn. "Bandit ambush success" → 30dk sonra "Caravan protector needed").
  - **Player-triggered** (Phase 02): oyuncu action bir event zincirini açar.
- **Fail/success state**: tamamlanamazsa event kaybolur (loss narrative — "caravan lost").
- **Chain reward bonus**: zincir tamamlanırsa büyük ödül (title, unique item).

## 3. Model

```sql
CREATE TABLE content.dynamic_event_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  location_id UUID NOT NULL,
  min_level SMALLINT NOT NULL,
  duration_seconds INT NOT NULL DEFAULT 1800,
  objective_kind VARCHAR(24) NOT NULL,          -- 'defend'|'kill_target'|'escort'|'collect'|'defuse'
  objective_params JSONB NOT NULL,
  success_reward JSONB NOT NULL,
  fail_reward JSONB NOT NULL,
  chain_parent_code VARCHAR(48) NULL,           -- önceki event
  chain_success_unlock JSONB NULL,              -- [{ next_event_code, delay_seconds }, ...]
  version INT DEFAULT 1
);

CREATE TABLE progression.dynamic_event_spawn (
  id UUID PRIMARY KEY,
  template_id UUID NOT NULL,
  location_id UUID NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  progress_score INT NOT NULL DEFAULT 0,        -- objective ilerleme (0..100)
  state VARCHAR(16) NOT NULL DEFAULT 'active',  -- active|succeeded|failed
  triggered_by_spawn_id UUID NULL,              -- chain parent spawn
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE progression.dynamic_event_contribution (
  spawn_id UUID NOT NULL,
  character_id UUID NOT NULL,
  contribution_score INT NOT NULL DEFAULT 0,
  reward_tier SMALLINT NULL,
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (spawn_id, character_id)
);
```

## 4. Objective Türleri

| Kind | Örnek | Progress kaynağı |
|---|---|---|
| `defend` | NPC/point koru; 5 dalga | Enemy kill sayısı |
| `kill_target` | Rare mob elite'i öldür | HP =0 |
| `escort` | NPC A→B taşı | NPC HP + waypoint reach |
| `collect` | 20 rare mat topla | Item pickup count |
| `defuse` | 5 rune deaktive et | Interact tamamlanan sayı |

Objective params content-driven JSONB.

## 5. Chain Sistemi (v1 basit)

Content:
```yaml
chain_success_unlock:
  - next_event_code: caravan_protect
    delay_seconds: 1800   # 30 dk sonra spawn
```

Success anında worker `dynamic_event_chain_trigger` çalışır:
- `delay_seconds` sonra `progression.dynamic_event_spawn` insert.
- WS "dynamic_event_chain_hint" (opt-in).

Chain **max 3 seviye** (v1); daha uzun zincir Phase 02.

## 6. AFK Katılım

Kısa süreli event → AFK contract oluşturma pahalı; **sınırlı model**:
- Karakter location'da AFK ise (`character.last_location.location_id == event.location_id`) event window boyunca minimum contribution alır (attendance tier 1).
- Ödül tier_1 sabit; tier_2+ online katılım gerektirir.
- Dynamic event AFK contract oluşturmaz (activity-engine yeni "location_presence" polling ile hesaplar).

## 7. Contribution → Reward

```yaml
contribution_tiers:
  tier_0: 0        # katılım yok
  tier_1: 1        # attendance (AFK dahil)
  tier_2: 30
  tier_3: 100
success:
  tier_1: { gold: 50c, event_token: 1 }
  tier_2: { gold: 200c, event_token: 5, loot: uncommon }
  tier_3: { gold: 800c, event_token: 15, loot: rare }
fail:
  tier_1: 0                            # attendance ödülü yok
  tier_2: { gold: 50c, event_token: 1 } # teselli
  tier_3: { gold: 200c, event_token: 3 }
chain_bonus_final:                     # zincir tamamlandığında toplam +
  { title: "Guardian_of_Glenmoor", unique_item: "storm_ring" }
```

## 8. Spawn Trigger

Worker `dynamic_event_spawner` her 5 dk:
- Content rotation'dan pick.
- Location'daki active player yoğunluğu > threshold ise spawn (v1: sabit rotation, Phase 02 dinamik density).
- Concurrent event limit per region: 3 (v1).

## 9. Endpoint

**GET `/v1/dynamic-event/active?location_id=`**
**GET `/v1/dynamic-event/{spawn_id}`**
**GET `/v1/character/{id}/dynamic-event/recent`**

## 10. Metrikler
- `dynamic_event_spawned_total{template}`
- `dynamic_event_success_total`, `dynamic_event_fail_total`
- `dynamic_event_chain_length_histogram`
- `dynamic_event_participant_total`

## 11. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md)
- [world/location.md](../world/location.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [systems/notification.md](../systems/notification.md)

## 12. Bekleyen Kararlar
- Chain max derinlik 3 v1; Phase 02'de 5.
- Player-triggered chain Phase 02.
- Dynamic density spawn Phase 02 (heat-map).
- Location AFK presence hesaplama detay Phase 02 (activity-engine polling vs snapshot).
