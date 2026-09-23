# Achievements

## 1. Amaç
Oyuncunun gerçekleştirdiği notable eylemleri kalıcı kayıt altına alan completion sistemi. Puan + title + cosmetic ödül.

## 2. Kategoriler (v1)

| Kategori | Örnek |
|---|---|
| `combat` | "1000 enemy kill", "solo dungeon boss" |
| `gathering` | "10k ore mined", "unique fish caught" |
| `crafting` | "100 rare item crafted", "first legendary craft" |
| `exploration` | "region 100% discovered", "10 hidden location" |
| `social` | "guild founded", "100 friends" |
| `dungeon` | "5-man X on hard", "no-death dungeon" |
| `pvp` | Phase 02 |
| `event` | "world boss defeat", "seasonal all-quest" |
| `progression` | "level 60", "1st ascension" |
| `collection` | "all pet species", "all mount rarities" |
| `misc` | Easter egg, meta |

## 3. Karakteristik

- **Trigger-based**: event bus dinleyen worker; anlık award.
- **Threshold-based**: sayaç (kill count) achievement → periyodik update.
- **Point value**: 5, 10, 25, 50, 100 (rarity göre).
- **Meta-achievement**: alt achievement setinin tamamlanması → meta unlock.
- **Retroactive**: yeni achievement content release'i eski data ile grant (opt-in flag).
- **Hidden**: bazı achievement UI'da gizli (spoiler / secret).

## 4. Model

```sql
CREATE TABLE content.achievement_template (
  id UUID PRIMARY KEY,
  code VARCHAR(64) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  description_key VARCHAR(128) NOT NULL,
  category VARCHAR(24) NOT NULL,
  point_value SMALLINT NOT NULL,
  icon_ref VARCHAR(128),
  is_hidden BOOLEAN NOT NULL DEFAULT false,
  is_meta BOOLEAN NOT NULL DEFAULT false,
  progress_kind VARCHAR(24) NOT NULL,       -- 'boolean'|'counter'|'checklist'|'meta'
  progress_target INT NOT NULL DEFAULT 1,
  checklist_items JSONB NULL,               -- checklist_kind için
  meta_child_codes JSONB NULL,              -- meta_kind için
  trigger_kinds JSONB NOT NULL,             -- ['enemy_killed', 'ore_gathered', ...]
  trigger_predicate JSONB NULL,             -- filter (ör. enemy_type='dragon')
  reward JSONB NOT NULL,                    -- { title_code?, cosmetic_sku?, currency? }
  scope VARCHAR(16) NOT NULL DEFAULT 'character', -- 'character'|'account'
  version INT DEFAULT 1
);

CREATE TABLE progression.achievement_progress (
  character_id UUID NOT NULL,
  achievement_id UUID NOT NULL,
  progress_current INT NOT NULL DEFAULT 0,
  checklist_state JSONB NOT NULL DEFAULT '{}',
  unlocked_at TIMESTAMPTZ NULL,
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (character_id, achievement_id)
);

CREATE TABLE progression.achievement_account_progress (
  account_id UUID NOT NULL,
  achievement_id UUID NOT NULL,
  progress_current INT NOT NULL DEFAULT 0,
  unlocked_at TIMESTAMPTZ NULL,
  claimed_at TIMESTAMPTZ NULL,
  PRIMARY KEY (account_id, achievement_id)
);

ALTER TABLE character.character
  ADD COLUMN achievement_points INT NOT NULL DEFAULT 0;
```

## 5. Trigger Akışı

Event bus:
```
event.progression.achievement_hint
  { character_id, kind, ref_id, payload }  -- ör. { kind:'enemy_killed', payload:{enemy_code:'goblin', is_elite:false} }
```

Worker `achievement_evaluator`:
1. Karakter için ilgili achievement template'lerini match (trigger_kinds + predicate).
2. Boolean: `unlocked_at IS NULL` ise mark unlocked.
3. Counter: `progress_current++`; hedef ulaşırsa unlock.
4. Checklist: `checklist_state[item_key]=true`; tüm items true ise unlock.
5. Meta: child achievement unlocked ise parent'ı reevaluate.
6. Unlock anında: reward apply (title, cosmetic, currency), `character.achievement_points += point_value`, WS "achievement_unlocked", audit log.

## 6. Meta Achievement

```yaml
- code: 'gathering_grandmaster'
  is_meta: true
  point_value: 100
  meta_child_codes: [ 'mining_master', 'woodcutting_master', 'fishing_master', 'herbalism_master', 'hunting_master' ]
  reward: { title_code: 'grandmaster_gatherer', cosmetic_sku: 'backpack_masterhand' }
```

Bir child unlock olunca meta reevaluate.

## 7. Retroactive Grant

Yeni achievement content release'te opsiyonel `retroactive=true` — worker `achievement_retroactive` `content.achievement_template.retroactive_since` timestamp ile geçmiş event log tarar (ör. `analytics.enemy_killed`).

V1'de retroactive **manuel opt-in** (ops decision content bazında); default false.

## 8. Reward Delivery

Unlock → reward JSONB parse → ilgili domain servis:
- `title_code` → [titles.md](titles.md) unlock.
- `cosmetic_sku` → [collections.md](collections.md) unlock.
- `currency` → ledger grant via TxnService.
- `mail_item` (Phase 02) → mail sender=SYSTEM attach item.

## 9. Endpoint

**GET `/v1/character/{id}/achievements?category=&unlocked=`**
**GET `/v1/character/{id}/achievements/{code}`**
**POST `/v1/character/{id}/achievements/{code}/claim`** — opt (auto-claim v1 default).

## 10. Rate & Batching

- Trigger event'ler batched (500ms window) → evaluator batch.
- High-volume trigger (enemy_killed) hem char-level hem account-level match count.
- `progression.achievement_progress` upsert `ON CONFLICT ... DO UPDATE ... progress_current = LEAST(target, progress_current + delta)`.

## 11. Metrikler
- `achievement_unlocked_total{category, code}`
- `achievement_progress_event_total{kind}`
- `achievement_evaluator_batch_latency_ms`
- `achievement_points_gauge_p50/p95`

## 12. Cross-Reference
- [titles.md](titles.md), [collections.md](collections.md), [milestones.md](milestones.md)
- [systems/reward-engine.md](../systems/reward-engine.md)
- [database/progression.md](../database/progression.md)

## 13. Bekleyen Kararlar
- Achievement puan leaderboard (Phase 02).
- Hidden achievement policy (spoiler warning) Phase 02.
- Cross-character (account) achievement sharing hangi ölçüde (v1: account scope explicit tanımlı; character scope default).
- Retroactive default policy Phase 02.
