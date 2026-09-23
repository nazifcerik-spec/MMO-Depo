# World Boss

## 1. Amaç
Global schedule ile spawn olan, çok sayıda oyuncunun aynı window'da katıldığı büyük boss encounter. AFK dostu — attendance contract ile offline oyuncu da paylaşabilir.

## 2. Karakteristik

- **Spawn window**: schedule veya rotation (haftada 2-4 kez).
- **Aktif encounter süresi**: 30-90 dk (varsayılan 60 dk).
- **Global visibility**: WS broadcast + world map marker + inbox push.
- **Attendance contract** (AFK): oyuncu event window'dan önce "attend" kaydı bırakır; window boyunca AFK contribution otomatik.
- **Personal contribution ledger**: her karakter için ayrı sayaç; overkill rekabeti yok.
- **Phase-based encounter**: HP eşiklerinde mekanik değişir (bkz. [combat/pve-encounters.md](../combat/pve-encounters.md) §7).
- **Reward tier**: contribution eşiklerine göre 4 tier.

## 3. Model

```sql
CREATE TABLE content.world_boss_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  min_level SMALLINT NOT NULL,
  hp_base BIGINT NOT NULL,
  phase_definitions JSONB NOT NULL,      -- HP threshold phases + mechanics
  enrage_timer_ms INT NOT NULL,
  location_id UUID NOT NULL,             -- spawn location
  duration_seconds INT NOT NULL DEFAULT 3600,
  loot_table_id UUID NOT NULL,
  event_token_reward JSONB NOT NULL,     -- per tier
  version INT DEFAULT 1
);

CREATE TABLE progression.world_boss_spawn (
  id UUID PRIMARY KEY,
  boss_template_id UUID NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'scheduled',  -- scheduled|active|defeated|expired
  defeated_at TIMESTAMPTZ NULL,
  hp_remaining BIGINT NOT NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE progression.world_boss_attendance (
  world_boss_spawn_id UUID NOT NULL,
  character_id UUID NOT NULL,
  registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  contribution_dmg BIGINT NOT NULL DEFAULT 0,
  contribution_heal BIGINT NOT NULL DEFAULT 0,
  contribution_afk_ticks INT NOT NULL DEFAULT 0,
  contribution_score BIGINT NOT NULL DEFAULT 0,   -- normalize edilmiş toplam
  reward_tier SMALLINT NULL,                       -- 0..4 (0=no-reward)
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (world_boss_spawn_id, character_id)
);
```

## 4. Schedule Modeli

- Cron rotation `world_boss_rotation` (content-driven): `[boss_code, day_of_week, hour_utc]`.
- Worker `world_boss_scheduler` her 5 dk'da bakıp bir sonraki spawn'ı `progression.world_boss_spawn` state='scheduled' olarak insert eder (24h öncesinden).
- Attendance registration açılır (event `world_boss.registration_open`).
- `starts_at` geldiğinde worker state='active' + `hp_remaining=hp_base` set eder, WS broadcast.
- `ends_at` veya `hp_remaining=0` olduğunda closer worker devreye girer.

## 5. Attendance Contract (AFK Katılım)

```
POST /v1/world-boss/{spawn_id}/attend
```
- Registration window: `starts_at - 24h → starts_at`.
- Karakter online değil de attendance kaydı bıraktıysa (aynı account başka session'dan veya event notification tıklamasından), event window boyunca `activity-engine` "world_boss_attend" contract'ı çalıştırır (bkz. [systems/activity-engine.md](../systems/activity-engine.md)).

Contract snapshot (DD-088):
```json
{
  "kind": "world_boss_attend",
  "world_boss_spawn_id": "...",
  "character_snapshot": { "level": 42, "gear_score": 1250, "class": "warrior" },
  "started_at": "starts_at",
  "expires_at": "ends_at"
}
```

Contract reveal (event kapanışında):
- AFK tick sayısı × character DPS proxy → `contribution_afk_ticks` + `contribution_score`.
- Online katılım (gerçek damage/heal) reveal ekler; online > AFK weight (AFK karşılığı yaklaşık %30-60 aynı seviye online oyuncunun; balance pass Phase 02).
- **Cap**: AFK contribution puanı `min(afk_score, tier_3_threshold)` — sadece attend etti AFK'de tier 4 alamaz (online katılım tavan için gerekli).

## 6. Contribution → Reward Tier

Content-driven eşikler:
```yaml
world_boss:
  contribution_tiers:
    tier_1: 1            # attended (min)
    tier_2: 1000
    tier_3: 10000
    tier_4: 50000
```

Her tier:
```yaml
tier_1: { gold: 100c, event_token: 1, loot_roll: personal_common }
tier_2: { gold: 500c, event_token: 5, loot_roll: personal_uncommon }
tier_3: { gold: 2000c, event_token: 15, loot_roll: personal_rare }
tier_4: { gold: 10000c, event_token: 40, loot_roll: personal_epic + chance_legendary }
```

Loot **personal roll** (bkz. [systems/loot.md](../systems/loot.md)); rekabet yok.

## 7. Reveal / Kapanış Akışı

Closer worker (event `ends_at` veya HP=0):
1. `state='defeated' | 'expired'`.
2. Her `world_boss_attendance` için:
   - AFK contract reveal → contribution_afk_ticks + score.
   - Total `contribution_score = dmg + heal×0.5 + afk_score`.
   - Reward tier map.
3. Reward grant: gold + event_token + loot roll (bkz. reward-engine).
4. Ledger + audit; player inbox notification.
5. First-kill bonus: HP=0 anında son 3 vurucu bonus + server anons (Achievement referansı).
6. Metrikler (`world_boss_defeat_latency_ms`, `world_boss_attendance_total`).

## 8. Contract Failure

- Attend etmedi + hiç dmg/heal atmadı → attendance satırı yok → hiç ödül yok.
- Registered + hiç aktivite yok (character offline & AFK contract active) → AFK ticks minimum → tier_1.
- Registered + event window boyunca farklı activity (mining vs) çalıştırıyor → **contract çakışması** → max concurrent AFK (DD-OPEN-3) sınırı devrede.

## 9. Endpoint

**GET `/v1/world-boss/next`** — sıradaki world boss'lar (24h horizon).
**POST `/v1/world-boss/{spawn_id}/attend`**
**POST `/v1/world-boss/{spawn_id}/unattend`**
**GET `/v1/world-boss/{spawn_id}`** — state + HP + kalan süre.
**GET `/v1/character/{id}/world-boss/history`**

## 10. Metrikler
- `world_boss_scheduled_total{boss_code}`
- `world_boss_active_gauge`
- `world_boss_defeat_latency_ms`
- `world_boss_attendance_total`
- `world_boss_tier_grant_total{tier}`
- `world_boss_afk_vs_online_dps_ratio`

## 11. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md), [systems/reward-engine.md](../systems/reward-engine.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [economy/currency.md](../economy/currency.md) — event_token
- [guild/boss.md](../guild/boss.md) — guild-scoped kuzeni

## 12. Bekleyen Kararlar
- AFK contribution weight (%30 vs %60 online) Phase 02 balance.
- Registration window 24h v1; Phase 02 dinamik (7d rotation preview).
- First-kill server anons + achievement Phase 02.
- Cross-shard world boss Phase 03.
- Tier eşikleri v1 taslak; balance pass'te tune.
