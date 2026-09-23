# Invasion (Region Invasion)

## 1. Amaç
Bir region boyunca yayılan çoklu-dalga (multi-wave) düşman istilası. Bölge savunması metaforu; tüm active/AFK karakterler otomatik eligible.

## 2. Karakteristik

- **Kapsam**: 1 region + tüm territory'ler.
- **Süre**: 2-6 saat (varsayılan 3 saat).
- **Wave sistemi**: 5-10 dalga; her dalga güçlenerek gelir; boss dalga son.
- **Region-active auto-eligible**: event tetiklendiğinde region'da bulunan / son 24h içinde region'da aktif olan karakter otomatik katılımcı.
- **Personal contribution instance**: shared timer + shared wave state ama personal kill contribution.
- **Fail state**: son wave'e ulaşılmadan süre biterse **başarısız** → yalnızca participation tier ödülü.
- **Success**: tüm dalgalar temizlenirse → full ödül + territory bonus buff (48h).

## 3. Model

```sql
CREATE TABLE content.invasion_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  min_level SMALLINT NOT NULL,
  region_id UUID NOT NULL,
  duration_seconds INT NOT NULL DEFAULT 10800,   -- 3h
  wave_definitions JSONB NOT NULL,               -- [{index, enemies, mini_boss}, ...]
  boss_enemy_template_id UUID NULL,
  success_reward JSONB NOT NULL,
  fail_reward JSONB NOT NULL,
  version INT DEFAULT 1
);

CREATE TABLE progression.invasion_spawn (
  id UUID PRIMARY KEY,
  invasion_template_id UUID NOT NULL,
  region_id UUID NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  current_wave SMALLINT NOT NULL DEFAULT 0,
  waves_total SMALLINT NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'scheduled', -- scheduled|active|succeeded|failed
  succeeded_at TIMESTAMPTZ NULL,
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE progression.invasion_participation (
  invasion_spawn_id UUID NOT NULL,
  character_id UUID NOT NULL,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  kill_count INT NOT NULL DEFAULT 0,
  wave_boss_participations INT NOT NULL DEFAULT 0,
  contribution_score BIGINT NOT NULL DEFAULT 0,
  reward_tier SMALLINT NULL,
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (invasion_spawn_id, character_id)
);
```

## 4. Trigger

İki yol:
- **Cron rotation** (v1 default): worker `invasion_scheduler` haftada X kez region rotation.
- **Dynamic trigger** (Phase 02): region PvE aktivitesi yüksek olduğunda otomatik (sıcaklık metriği).

Trigger anı worker state='active' + wave scheduler başlatır (5-10 dalga interval'li).

## 5. Region-Active Eligibility

Karakter eligible kabul edilir:
- Şu an region'da (`character.last_location.region_id`), veya
- Son 24 saat içinde region'da bulunmuş (`character.recent_regions` history).

Kayıt `progression.invasion_participation` insert (joined_at=now()) — event window'da ilk kill/heal/AFK tick attığında otomatik.

## 6. Wave Sistemi

Content:
```yaml
waves:
  - index: 1
    enemy_pack: goblin_scout×20
    duration_seconds: 900
    threshold_kills_to_advance: 15   # yeterince öldür → sonraki dalga
  - index: 2
    enemy_pack: goblin_raider×30
    duration_seconds: 1200
  ...
  - index: N
    boss: goblin_warlord
    duration_seconds: 1800
```

Her dalga:
- Wave başlar → enemy spawn (personal instance içinde, bkz. DD-091/115).
- Threshold aşıldığında sonraki dalga.
- Süre dolarsa yine sonraki dalga (fail contribution eksik).
- Son dalga boss: HP=0 → event success.

## 7. AFK Katılım

Region-active AFK karakter otomatik AFK contract oluşturur (event trigger anında). Contract snapshot:
```json
{
  "kind": "invasion_defend",
  "invasion_spawn_id": "...",
  "character_snapshot": { ... },
  "started_at": "...",
  "expires_at": "ends_at"
}
```

Reveal: AFK tick × dps proxy × wave difficulty scaling = kill_count contribution.

**Cap**: AFK boss dalgada contribution %50 → boss son vurucu bonus alamaz AFK.

## 8. Reward Tier

Success case:
```yaml
tiers:
  1: score >= 1        # attended (min)
  2: score >= 500
  3: score >= 2000
  4: score >= 8000
success_reward_per_tier:
  1: { gold: 200c, event_token: 2 }
  2: { gold: 1000c, event_token: 8, loot: uncommon }
  3: { gold: 3000c, event_token: 20, loot: rare }
  4: { gold: 10000c, event_token: 50, loot: epic }
region_success_bonus:
  buff: "invasion_survivor" duration=48h  # +5% XP/gathering in region
```

Fail case: her tier'ın **yarısı** (attendance ve teselli).

## 9. Region Buff (48h)

Success sonrası region'da bulunan tüm karakterler `invasion_survivor` buff alır (48 saat). Bkz. [combat/status-effects.md](../combat/status-effects.md).

## 10. Endpoint

**GET `/v1/invasion/active?region_id=`**
**GET `/v1/invasion/{spawn_id}`**
**GET `/v1/character/{id}/invasion/history`**

Explicit attendance yok — auto-eligible. Reveal window'da katılım gösterilir.

## 11. Metrikler
- `invasion_scheduled_total{region_id}`
- `invasion_success_total`, `invasion_fail_total`
- `invasion_participant_total`
- `invasion_wave_advance_latency_ms`
- `invasion_afk_share_ratio`

## 12. Cross-Reference
- [systems/activity-engine.md](../systems/activity-engine.md)
- [world/region.md](../world/region.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [combat/status-effects.md](../combat/status-effects.md)

## 13. Bekleyen Kararlar
- Auto-eligible 24h horizon tuning (12h vs 48h) Phase 02.
- Dynamic trigger sıcaklık metriği Phase 02.
- Region buff çakışması (concurrent invasion) — v1 yeni buff eskisini refresh eder.
- Personal instance vs shared spawn — v1 personal (DD-091 uyumlu); Phase 02 shared world boss kolektif hissiyat opt.
