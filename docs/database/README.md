# Database Schema — Genel Bakış

Bu klasör, Phase 00 tasarım dokümanlarındaki tüm domain'lerin PostgreSQL şemasına indirgenmiş halini içerir. Alt sistem dokümanlarında (character.md, items.md, activity-engine.md, ...) tanımlanan entity'ler burada tablo/kolon seviyesine düşürülür.

Bu doküman **normatif** değildir; **rehber**dir. Gerçek DDL Phase 01'de `goose`/`atlas` migration'ları olarak yazılır. Burası "hangi tablo, hangi sütun, hangi ilişki" sorusunun cevabıdır.

---

## 1. Şema Grupları

PostgreSQL "schema" (namespace) bazlı ayrılır:

| Schema | Amaç | Doküman |
|---|---|---|
| `auth` | Users, accounts, sessions, permissions | [auth-accounts.md](auth-accounts.md) |
| `character` | Karaktere ait tüm veriler | [character.md](character.md) |
| `content` | Template'ler (item, race, class, skill, ...) | [content-templates.md](content-templates.md) |
| `world` | Kıta → şehir → POI hiyerarşisi | [world.md](world.md) |
| `activity` | Contract, reveal log, event log | [activity.md](activity.md) |
| `economy` | Ledger, marketplace, mail | [economy-market.md](economy-market.md) |
| `social` | Guild, party, chat, friendship | [social-guilds-parties.md](social-guilds-parties.md) |
| `companion` | Pet, mount instance + koleksiyon | [companions.md](companions.md) |
| `progression` | Quest, achievement, mastery | [progression-quests-achievements.md](progression-quests-achievements.md) |
| `audit` | Event log, admin, security | [audit-system.md](audit-system.md) |
| `system` | Migrations, feature flags, schema version | [audit-system.md](audit-system.md) |

---

## 2. Genel Konvansiyonlar

### 2.1 Naming
- **snake_case** her yerde: schema, table, column, index, constraint.
- Tablo isimleri **tekil** (`character`, değil `characters`).
- FK sütun adı: `<referenced_table>_id` (örn. `race_id`, `owner_character_id`).
- Junction/many-to-many tablo: `<a>_<b>` (örn. `character_skill`).
- Index adı: `ix_<table>__<cols>`; unique: `uq_<table>__<cols>`; check: `ck_<table>__<rule>`.

### 2.2 PK Stratejisi
- **UUID** çoğu tablo için (`gen_random_uuid()`; pgcrypto).
- Yüksek-throughput append tabloları (audit, ledger event, chat) için `BIGSERIAL` tartışılabilir; ancak cross-region distribution için UUID (v7 zaman-ordered) tercih.
- **UUID v7** planı: Phase 01'de kütüphane seçimi yapılana kadar v4 kullanılır.

### 2.3 Timestamps
- Her satırda `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`.
- Mutable satırda `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`; trigger değil, application layer güncellemesi (transparency).
- Soft delete: `deleted_at TIMESTAMPTZ NULL` + partial index `WHERE deleted_at IS NULL`.

### 2.4 Optimistic Locking
- Concurrency kritik tablolarda `version INT NOT NULL DEFAULT 1`.
- Update her seferinde `SET version = version + 1 WHERE version = :expected`.

### 2.5 Enums vs Lookup Table
- Küçük, kararlı set (< 20, nadir değişir): `CREATE TYPE ... AS ENUM (...)`.
- Büyük veya sık değişen: lookup tablo (`content.*_type`).

### 2.6 JSONB Kullanımı
- Esnek payload'lar (item stat roll, race passive definition, activity params).
- Sık sorgulanan JSONB alanları **generated column + index** ile hızlandırılır.
- JSONB **her şey için değil**; alan sabit ise sütun aç.

### 2.7 Currency ve Sayısal Alanlar
- Para: `BIGINT` (integer cents/copper birimi; gold = 10000 copper mantığı Phase 01 karar).
- Item stat: `INT` yeterli.
- XP: `BIGINT` (uzun ömür + ascension'lar için taşma güvenliği).

### 2.8 Text Fields
- Kısa isim (character, guild): `VARCHAR(64) NOT NULL`.
- Uzun serbest metin (chat, description): `TEXT`.
- Case-insensitive isim (character.name, guild.name): `CITEXT` (unique).
- i18n text: kod tarafında key + `content.i18n_string` tablosu (Phase 01 opsiyonel; bloklama olmasın diye şimdilik uygulama i18n resolver'ı).

### 2.9 Constraint Disiplini
- NOT NULL agresif kullanılır; nullable alan **istisnaidir** ve gerekçesi vardır.
- FK her zaman `ON DELETE` explicit: `RESTRICT` default; `CASCADE` sadece owned/subordinate ilişkide.
- Enum ve range validation için `CHECK` constraint.

### 2.10 Indexler
- Her FK sütununa index (Postgres bunu otomatik yapmaz).
- Composite index sorgu path'lerine göre.
- Partial index (soft delete filter, status filter) düşünülür.
- BRIN large append tablolar için (audit event_log).

---

## 3. Ana Entity İlişki Diyagramı (Yüksek Seviye)

```
auth.user ─┬─ auth.account ─── auth.session
           │
           └─────────────────── character.character ────────┐
                                    │                        │
                                    ├─ character.stats        │
                                    ├─ character.inventory   │
                                    ├─ character.equipment   │
                                    ├─ character.skill_learn │
                                    ├─ character.profession  │
                                    ├─ character.mastery     │
                                    ├─ character.ascension   │
                                    ├─ character.settings    │
                                    │                        │
                                    ├─ activity.contract ────┼─ activity.reveal_log
                                    ├─ activity.instance     │
                                    │                        │
                                    ├─ companion.pet_owned   │
                                    ├─ companion.mount_owned │
                                    ├─ progression.quest_progress
                                    ├─ progression.achievement_progress
                                    │                        │
                                    └─ social.friendship ────┤
                                       social.ignore         │
                                       social.chat_message   │
                                       social.mail           │
                                       economy.ledger  ──────┤
                                       economy.market_listing│
                                       economy.market_txn    │
                                       social.guild_member ──┤
                                       social.party_member   │

content.race ─┐
content.class ┼── character.character (fk)
content.spec ─┘

content.item_template ── character.inventory / equipment / mail / market_listing (via item_instance)
content.skill_template ── character.skill_learn
content.recipe ── (crafting side of activity)
content.dungeon_template ── activity.instance
content.adventure_template ── activity.instance
content.npc_template ── (combat/hunting side of activity)
content.loot_table ── (referenced by mob/node/boss/chest)
content.pet_template / mount_template ── companion.*

world.continent → region → territory → city|village|outpost → location → poi
world.resource_node (instance) ── linked to location

social.guild ─┬─ social.guild_member
              ├─ social.guild_bank
              ├─ social.guild_bank_tab
              └─ social.guild_research (Phase 02)

social.party ─── social.party_member
```

---

## 4. Migration Discipline

- **Forward-only** tercih; rollback yalnızca test edilmişse.
- **Expand-Contract** breaking change'te 2-adımlı:
  1. Yeni sütun/tablo ekle + dual write.
  2. Kod yeni yola tamamen geçince eski dropla.
- Tüm migration idempotent (retry safe).
- Migration adı: `NNNN_short_description.sql`.
- Data migration büyükse arka planda backfill job (worker).

---

## 5. Seed Data

- **Content seed** (race, class, item_template, ...): yaml/json içerik pack → loader.
- **World seed** (continent, region, city): yaml/json.
- **Dev seed** (fake user/character): script (`cmd/devseed`).
- **Prod'a asla dev seed** — pipeline check.

---

## 6. Bekleyen Kararlar

- UUID v7 kütüphanesi seçimi (Phase 01).
- Currency birimi: gold cinsinden mi copper cinsinden mi (öneri: copper = 1, gold = 10000 copper).
- i18n çözümü: `content.i18n_string` tablo mu, resource file mı (öneri: resource file başlangıçta, DB tabloya migrate opsiyonel).
- Partition stratejisi (audit, chat) — Phase 02.
- Row-level security: v1'de gerekli değil; admin panel için Phase 02.
- Read replica sayısı ve routing.

---

## 7. Bu Klasördeki Dokümanlara Girmek

Her doküman şu iskeleti izler:
- **Kapsam** — hangi domain, hangi tablolar.
- **Tablo detayları** — kolonlar, tipler, constraint'ler, indexler.
- **Kritik ilişkiler** — FK ve semantic bağlar.
- **Kullanım örnekleri** — tipik sorgu path'leri.
- **Notlar / edge cases** — özel dikkat.

Doküman dili: Türkçe; tablo/kolon isimleri İngilizce/snake_case.
