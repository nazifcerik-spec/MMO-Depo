# Database — `social` Schema

## Kapsam
Guild + party + friendship + ignore + chat + notification + report. Sosyal etkileşimin tüm kalıcı state'i.

---

## 1. `social.guild`

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `name` | `CITEXT` | UNIQUE NOT NULL, `CHECK length between 3 and 24` |
| `tag` | `CITEXT` | UNIQUE NOT NULL, `CHECK length between 2 and 5` |
| `leader_character_id` | `UUID` | NOT NULL |
| `motd` | `TEXT` | NOT NULL DEFAULT '' |
| `description` | `TEXT` | NOT NULL DEFAULT '' |
| `level` | `SMALLINT` | NOT NULL DEFAULT 1 |
| `xp` | `BIGINT` | NOT NULL DEFAULT 0 |
| `treasury_gold_copper` | `BIGINT` | NOT NULL DEFAULT 0 |
| `home_city_id` | `UUID` | FK → world.settlement NULL |
| `member_cap` | `SMALLINT` | NOT NULL DEFAULT 50 |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `disband_scheduled_at` | `TIMESTAMPTZ` | NULL — 48h grace |
| `disbanded_at` | `TIMESTAMPTZ` | NULL |

Indexler: `uq_guild__name`, `uq_guild__tag`.

---

## 2. `social.guild_rank`

Guild-özel rank tanımı (customizable Phase 02; v1 sabit set).

| Kolon | Tip |
|---|---|
| `guild_id` | `UUID` PK-part |
| `rank_index` | `SMALLINT` PK-part `CHECK between 0 and 9` |
| `name` | `VARCHAR(32)` NOT NULL |
| `permissions_json` | `JSONB` NOT NULL DEFAULT `'{}'` — invite, kick, promote, bank_tab_access, etc. |

Default rank set: 0=leader, 1=officer, 2=veteran, 3=member, 4=initiate.

---

## 3. `social.guild_member`

| Kolon | Tip | Constraint |
|---|---|---|
| `guild_id` | `UUID` | PK-part |
| `character_id` | `UUID` | PK-part |
| `rank_index` | `SMALLINT` | NOT NULL DEFAULT 4 |
| `joined_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `contribution_points` | `BIGINT` | NOT NULL DEFAULT 0 |
| `last_active_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `note_private` | `VARCHAR(128)` | NOT NULL DEFAULT '' |
| `note_officer` | `VARCHAR(256)` | NOT NULL DEFAULT '' |

Constraint: `UNIQUE (character_id)` — bir karakter tek bir guild'de.
Indexler: `ix_guild_member__guild_id`, `uq_guild_member__character_id`.

---

## 4. `social.guild_bank_tab`

```
guild_id UUID PK-part,
tab_index SMALLINT PK-part,
name VARCHAR(32) NOT NULL,
capacity SMALLINT NOT NULL DEFAULT 40,
purchased_at TIMESTAMPTZ NOT NULL,
purchase_cost_copper BIGINT NOT NULL,
permissions_json JSONB NOT NULL DEFAULT '{}'  -- per rank read/write
```

## 5. `social.guild_bank_slot`

```
guild_id UUID PK-part,
tab_index SMALLINT PK-part,
slot_index SMALLINT PK-part,
item_instance_id UUID FK NOT NULL,
stack_count INT NOT NULL,
updated_at TIMESTAMPTZ NOT NULL
```

## 6. `social.guild_bank_log`

Guild bank audit (deposit/withdraw).

```
id BIGSERIAL PK,
guild_id UUID NOT NULL,
character_id UUID NOT NULL,
action VARCHAR(16) NOT NULL — deposit/withdraw/move/gold_deposit/gold_withdraw,
tab_index SMALLINT NULL,
item_instance_id UUID NULL,
qty INT NULL,
gold_copper BIGINT NULL,
at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Index: `ix_guild_bank_log__guild_id__at DESC`.

## 7. `social.guild_invite`

```
id UUID PK,
guild_id UUID FK NOT NULL,
target_character_id UUID FK NOT NULL,
invited_by_character_id UUID NOT NULL,
status VARCHAR(16) NOT NULL DEFAULT 'pending' — pending/accepted/declined/expired
expires_at TIMESTAMPTZ NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
UNIQUE (guild_id, target_character_id) WHERE status='pending'
```

## 8. `social.guild_research` (Phase 02)

```
guild_id UUID PK-part,
research_id UUID PK-part FK → content.guild_research_template,
level SMALLINT NOT NULL DEFAULT 0,
xp BIGINT NOT NULL DEFAULT 0
```

---

## 9. `social.party`

Küçük ölçekli party (session-scope; long-lived variant Phase 02).

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `leader_character_id` | `UUID` NOT NULL |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `disbanded_at` | `TIMESTAMPTZ` NULL |
| `loot_rule` | `social.party_loot_rule` NOT NULL DEFAULT `'personal'` |
| `activity_instance_id` | `UUID` FK → activity.instance NULL |

Enum:
```sql
CREATE TYPE social.party_loot_rule AS ENUM ('personal','round_robin','need_before_greed','master');
```

## 10. `social.party_member`

```
party_id UUID PK-part,
character_id UUID PK-part,
role_hint VARCHAR(16) NULL — tank/heal/dps,
joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
left_at TIMESTAMPTZ NULL,
ready BOOLEAN NOT NULL DEFAULT false,
UNIQUE (character_id) WHERE left_at IS NULL
```

## 11. `social.lfg_listing`

Looking-for-group ilanı.

```
id UUID PK,
character_id UUID FK NOT NULL,
activity_kind VARCHAR(24) NOT NULL — dungeon/adventure/hunt/world_boss,
template_id UUID NULL — dungeon/adventure template
role_hint VARCHAR(16) NULL,
message VARCHAR(160) NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
expires_at TIMESTAMPTZ NOT NULL,
status VARCHAR(16) NOT NULL DEFAULT 'open' — open/matched/expired/cancelled
```

Index: `ix_lfg__status__activity_kind__template_id`.

---

## 12. `social.friendship`

```
character_a UUID PK-part,
character_b UUID PK-part,
since TIMESTAMPTZ NOT NULL DEFAULT now(),
note VARCHAR(64) NOT NULL DEFAULT '',
CHECK (character_a < character_b)  -- canonical ordering
```

## 13. `social.friendship_request` (opsiyonel)

```
id UUID PK,
requester_character_id UUID NOT NULL,
target_character_id UUID NOT NULL,
status VARCHAR(16) NOT NULL DEFAULT 'pending',
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
UNIQUE (requester, target) WHERE status='pending'
```

## 14. `social.ignore_list`

```
owner_character_id UUID PK-part,
target_character_id UUID PK-part,
added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
reason VARCHAR(64) NOT NULL DEFAULT ''
```

---

## 15. `social.chat_channel_membership`

Custom channel abonelikleri (v2 detay). v1'de sistem kanalları (guild, party, city, zone) implicit.

```
character_id UUID PK-part,
channel_key VARCHAR(48) PK-part,
joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
muted BOOLEAN NOT NULL DEFAULT false
```

## 16. `social.chat_message`

Chat log; büyür — partition adayı.

| Kolon | Tip |
|---|---|
| `id` | `BIGSERIAL` PK |
| `channel_kind` | `social.chat_channel_kind` NOT NULL |
| `channel_ref_id` | `UUID` NULL — guild_id / party_id / city_id / vs |
| `sender_character_id` | `UUID` NOT NULL |
| `target_character_id` | `UUID` NULL — whisper için |
| `body` | `TEXT` NOT NULL — 512 char sınırı (application check) |
| `sent_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `filtered` | `BOOLEAN` NOT NULL DEFAULT false — moderation |
| `deleted_at` | `TIMESTAMPTZ` NULL |

Enum:
```sql
CREATE TYPE social.chat_channel_kind AS ENUM (
  'say','zone','city','guild','party','whisper','trade','lfg','global'
);
```

Indexler:
- `ix_chat_message__channel_kind__channel_ref_id__sent_at DESC` (feed).
- `ix_chat_message__sender_character_id__sent_at DESC` (report review).
- BRIN on `sent_at`.

Partitioning: month-range (Phase 02).

## 17. `social.chat_mute`

```
character_id UUID PK-part,
scope VARCHAR(24) PK-part — global/guild/whisper,
muted_until TIMESTAMPTZ NOT NULL,
reason VARCHAR(128) NULL
```

---

## 18. `social.notification`

In-game inbox + delivery record.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `character_id` | `UUID` NOT NULL |
| `category` | `social.notification_category` NOT NULL |
| `priority` | `social.notification_priority` NOT NULL DEFAULT `'normal'` |
| `title_key` | `VARCHAR(64)` NOT NULL |
| `body_key` | `VARCHAR(96)` NOT NULL |
| `params_json` | `JSONB` NOT NULL DEFAULT `'{}'` |
| `channels_delivered` | `VARCHAR(16)[]` NOT NULL DEFAULT '{}' — in_game/push/email |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `read_at` | `TIMESTAMPTZ` NULL |

Enums:
```sql
CREATE TYPE social.notification_category AS ENUM (
  'system','activity','market','guild','event','boss','social','achievement'
);
CREATE TYPE social.notification_priority AS ENUM ('low','normal','high','critical');
```

Indexler: `ix_notification__character_id__created_at DESC`.

## 19. `social.notification_pref`

Aslında `auth.account.notification_prefs` JSONB'de tutuluyor; per-category tablo tercih edilirse:

```
character_id UUID PK-part,
category social.notification_category PK-part,
channel VARCHAR(16) PK-part,  -- in_game/push/email
enabled BOOLEAN NOT NULL DEFAULT true,
updated_at TIMESTAMPTZ NOT NULL
```

---

## 20. `social.report`

Kullanıcı şikayeti.

| Kolon | Tip |
|---|---|
| `id` | `UUID` PK |
| `reporter_character_id` | `UUID` NOT NULL |
| `reported_character_id` | `UUID` NOT NULL |
| `reason` | `VARCHAR(48)` NOT NULL — chat_abuse, cheating, name, other |
| `evidence_json` | `JSONB` NOT NULL DEFAULT `'{}'` — chat message ids, screenshots url etc. |
| `context_channel_kind` | `social.chat_channel_kind` NULL |
| `status` | `VARCHAR(16)` NOT NULL DEFAULT 'open' — open/reviewing/resolved/dismissed |
| `resolution_note` | `TEXT` NULL |
| `resolved_by_user_id` | `UUID` NULL |
| `created_at` | `TIMESTAMPTZ` NOT NULL DEFAULT now() |
| `resolved_at` | `TIMESTAMPTZ` NULL |

Indexler: `ix_report__status__created_at`, `ix_report__reported_character_id`.

---

## 21. Kullanım Örnekleri

**Guild üye listesi:**
```sql
SELECT gm.*, c.name, c.level, c.race_id, c.class_id, c.last_seen_at
FROM social.guild_member gm
JOIN character.character c ON c.id = gm.character_id
WHERE gm.guild_id = $1
ORDER BY gm.rank_index ASC, gm.contribution_points DESC;
```

**Guild bank deposit:**
```sql
BEGIN;
INSERT INTO social.guild_bank_slot (...) VALUES (...);
UPDATE content.item_instance SET owner_kind='guild_bank', owner_id=$guild_id WHERE id=$item;
INSERT INTO social.guild_bank_log (guild_id, character_id, action, ...) VALUES (...);
COMMIT;
```

**Chat retrieval (guild):**
```sql
SELECT id, sender_character_id, body, sent_at
FROM social.chat_message
WHERE channel_kind = 'guild' AND channel_ref_id = $guild_id AND deleted_at IS NULL
ORDER BY sent_at DESC
LIMIT 100;
```

---

## 22. Notlar
- **Guild leader change**: transaction; disband 48h grace period.
- **Party lifecycle**: session-scope; disband on last member leave veya explicit.
- **Chat retention**: 30 gün default; audit uzun tutar.
- **Report action → audit.event_log** (mod aksiyon).

---

## 23. Bekleyen Kararlar
- Guild alliance (Phase 03) tablosu şeması.
- Cross-guild chat integration.
- Voice channel meta store (Phase 03).
- LFG matching algorithm (Phase 01 basit; Phase 02 smart).
