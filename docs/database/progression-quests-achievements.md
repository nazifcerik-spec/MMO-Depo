# Database — `progression` Schema

## Kapsam
Quest progress, achievement progress, mastery (character bazlı ilerleme detayı character.mastery'de), event participation, title earn history.

Template'ler `content.*`; burası oyuncunun kişisel ilerleme kayıtları.

---

## 1. `progression.quest_progress`

| Kolon | Tip | Constraint |
|---|---|---|
| `character_id` | `UUID` | PK-part |
| `quest_id` | `UUID` | PK-part FK → content.quest_template |
| `state` | `progression.quest_state` | NOT NULL DEFAULT `'accepted'` |
| `objective_progress_json` | `JSONB` | NOT NULL DEFAULT `'{}'` — {obj_index: current_count} |
| `accepted_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `completed_at` | `TIMESTAMPTZ` | NULL |
| `turned_in_at` | `TIMESTAMPTZ` | NULL |
| `expires_at` | `TIMESTAMPTZ` | NULL — time-limited ise |
| `repeat_count` | `INT` | NOT NULL DEFAULT 0 — repeatable ise |
| `last_repeat_at` | `TIMESTAMPTZ` | NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enum:
```sql
CREATE TYPE progression.quest_state AS ENUM (
  'accepted','complete','turned_in','abandoned','expired'
);
```

Indexler:
- `ix_quest_progress__character_id__state`.
- `ix_quest_progress__expires_at` (worker).

---

## 2. `progression.quest_reward_claim`

Ödül dağıtım idempotency; ledger'a bir kez düşmüş mü kontrolü.

```
quest_progress_pk (character_id, quest_id) UNIQUE,
claimed_at TIMESTAMPTZ NOT NULL,
correlation_id UUID NOT NULL
```

(v1'de basit tutulur: `quest_progress.turned_in_at IS NOT NULL` ödül dağıtıldığını gösterir.)

---

## 3. `progression.achievement_progress`

| Kolon | Tip |
|---|---|
| `character_id` | `UUID` PK-part |
| `achievement_id` | `UUID` PK-part FK |
| `progress_json` | `JSONB` NOT NULL DEFAULT `'{}'` |
| `completed_at` | `TIMESTAMPTZ` NULL |
| `first_server_flag` | `BOOLEAN` NOT NULL DEFAULT false |
| `updated_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `version` | `INT` NOT NULL DEFAULT 1 |

Indexler:
- `ix_achievement_progress__character_id__completed_at`.
- `ix_achievement_progress__achievement_id__completed_at` (server-first leaderboard).

---

## 4. `progression.title_earned`

Character.title_owned zaten var; alternatif kayıt burada tutulabilir. Tutarlılık için: character.title_owned tekil kaynak; ancak kazanım event log'u:

```
id BIGSERIAL PK,
character_id UUID NOT NULL,
title_id UUID FK NOT NULL,
earned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
source VARCHAR(32) NOT NULL — achievement/event/quest
source_ref_id UUID NULL
```

---

## 5. `progression.event_participation`

Etkinlik katılım ve ödül claim.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `character_id` | `UUID` NOT NULL |
| `event_template_id` | `UUID` FK NOT NULL |
| `event_spawn_id` | `UUID` NULL — emergent event instance |
| `contribution_score` | `BIGINT` NOT NULL DEFAULT 0 |
| `rewards_claimed` | `BOOLEAN` NOT NULL DEFAULT false |
| `claimed_at` | `TIMESTAMPTZ` NULL |
| `first_seen_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `updated_at` | `TIMESTAMPTZ` NOT NULL |

Indexler: `ix_event_part__character_id__event_template_id`.

---

## 6. `progression.world_boss_participation`

World boss kill credit ve contribution.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `spawn_id` | `UUID` FK → progression.world_boss_spawn NOT NULL |
| `character_id` | `UUID` NOT NULL |
| `damage_dealt` | `BIGINT` NOT NULL DEFAULT 0 |
| `healing_done` | `BIGINT` NOT NULL DEFAULT 0 |
| `participation_ms` | `INT` NOT NULL DEFAULT 0 |
| `kill_credit` | `BOOLEAN` NOT NULL DEFAULT false |
| `rewards_claimed` | `BOOLEAN` NOT NULL DEFAULT false |

Indexler: `ix_wb_part__spawn_id`, `ix_wb_part__character_id__spawn_id UNIQUE`.

## 7. `progression.world_boss_spawn`

Boss spawn instance (contribution tablosunun parent'ı).

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `template_id` | `UUID` FK → content.world_boss_template NOT NULL |
| `spawned_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `defeated_at` | `TIMESTAMPTZ` NULL |
| `defeated_by_last_hit_character_id` | `UUID` NULL |
| `hp_max` | `BIGINT` NOT NULL |
| `hp_remaining` | `BIGINT` NOT NULL |
| `participants_count` | `INT` NOT NULL DEFAULT 0 |

Indexler: `ix_wb_spawn__template_id__spawned_at DESC`.

---

## 8. `progression.dungeon_lockout`

Weekly / daily lockout tracking.

| Kolon | Tip |
|---|---|
| `character_id` | `UUID` PK-part |
| `dungeon_template_id` | `UUID` PK-part |
| `lockout_kind` | `VARCHAR(16)` PK-part — 'daily' / 'weekly' |
| `locked_until` | `TIMESTAMPTZ` NOT NULL |
| `runs_completed_in_window` | `SMALLINT` NOT NULL DEFAULT 0 |

Cleanup job: `locked_until < now()` → row silinir veya count reset.

---

## 9. `progression.faction_reputation` (Phase 02)

`character.reputation` tablosu (Phase 02) burada da tanımlı; consistency için:

```
character_id UUID PK-part,
faction_id UUID PK-part,
value INT NOT NULL DEFAULT 0,
tier SMALLINT NOT NULL DEFAULT 0,
capped_at INT NULL,
updated_at TIMESTAMPTZ NOT NULL
```

---

## 10. `progression.ascension_history`

Ascension olaylarının kayıt log'u (character.character'da ascension_level var; buradan detay tarihçe).

```
id BIGSERIAL PK,
character_id UUID NOT NULL,
ascension_number SMALLINT NOT NULL — 1,2,3...
ascended_at TIMESTAMPTZ NOT NULL DEFAULT now(),
points_awarded SMALLINT NOT NULL,
xp_curve_multiplier NUMERIC(4,2) NOT NULL — 1.05, 1.10, ...
UNIQUE (character_id, ascension_number)
```

---

## 11. Kullanım Örnekleri

**Aktif quest listesi:**
```sql
SELECT qp.*, qt.code, qt.rewards_json, qt.objectives_json
FROM progression.quest_progress qp
JOIN content.quest_template qt ON qt.id = qp.quest_id
WHERE qp.character_id = $1 AND qp.state IN ('accepted','complete')
ORDER BY qp.accepted_at DESC;
```

**Quest objective tick (AFK reveal sonrası):**
```sql
UPDATE progression.quest_progress
SET objective_progress_json = jsonb_set(
      objective_progress_json,
      ARRAY[$obj_index::text],
      to_jsonb( COALESCE((objective_progress_json->>$obj_index::text)::int, 0) + $delta )
    ),
    version = version + 1,
    updated_at = now()
WHERE character_id = $1 AND quest_id = $2 AND version = $expected;
```

**Dungeon lockout check:**
```sql
SELECT locked_until FROM progression.dungeon_lockout
WHERE character_id = $1 AND dungeon_template_id = $2 AND lockout_kind = 'weekly';
```

---

## 12. İlişkiler
- Quest progress → quest_template (content), character.
- Achievement → content.achievement, character.
- Event participation → event_template (content), reward payload → economy.ledger_entry.
- World boss spawn → content.world_boss_template.
- Ascension history → character (parent), content.ascension_node (spent choices).

---

## 13. Notlar
- **Quest reward idempotency**: state 'turned_in' + ledger idempotency_key kombinasyonu ile double-grant engellenir.
- **Achievement first-server**: race condition için serializable transaction; sunucu-genelinde ilk complete = `first_server_flag = true`.
- **Dungeon lockout window sıfırlama**: Weekly reset Salı 06:00 UTC (DD-016); worker job.
- **Event participation cross-event contribution**: bir event için tek row; contribution artırılır.
- **Ascension**: character reset transaction'ında bu tabloya history satırı insert.

---

## 14. Bekleyen Kararlar
- Quest chain branching (Phase 02) modeli.
- Achievement point vendor (Phase 02) tablosu.
- Seasonal / rotating quests (Phase 02).
- Ranked leaderboard tabloları (Phase 02 PvP).
