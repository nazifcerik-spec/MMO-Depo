# Raid (Guild Raid)

> ⚠️ **Phase 02** — 10-25 kişi raid content.
>
> **Faz 01 v1 (DD-261):** 3–5 kişilik *parti* baskınları uygulandı — `group_max > 1` zindan şablonları, `internal/raid`, kişisel ganimet, zindanlarla ortak haftalık kilit. Bu belge 10–25 kişilik lonca baskını tasarımı olarak kalır.

## 1. Amaç
Party'nin (2-5 kişi) ötesinde büyük ölçekli coop encounter. Guild-organized.

## 2. Karakteristik

- **Size**: 10-25 üye (v1 Phase 02: 10; Phase 03 raid 20-25).
- **Weekly lockout** shared.
- **Difficulty tiers**: normal / heroic / mythic.
- **Loot**: master loot rule available (raid leader distribute).
- **Content**: multi-phase boss + adds + mechanics.

## 3. Model (Phase 02)

```sql
CREATE TABLE content.raid_template (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  min_guild_level SMALLINT NOT NULL,
  raid_size_min SMALLINT NOT NULL DEFAULT 10,
  raid_size_max SMALLINT NOT NULL DEFAULT 25,
  boss_encounters JSONB NOT NULL,
  loot_table_ids UUID[] NOT NULL,
  weekly_lockout BOOLEAN DEFAULT true,
  version INT DEFAULT 1
);

CREATE TABLE activity.raid_instance (
  id UUID PRIMARY KEY,
  template_id UUID NOT NULL,
  guild_id UUID NULL,
  raid_leader_character_id UUID NOT NULL,
  participants UUID[] NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'active'
);
```

## 4. Raid Formation

Officer+ organize raid:
```
POST /v1/guilds/{id}/raid/create
  body: { raid_template_code, difficulty }
```

Members opt-in / invite based.
Raid ready check → enter.

## 5. Combat

Raid combat session:
- Party-like combat but 10-25 participants.
- Boss encounter multi-phase.
- Adds spawn mechanic.
- Positional mechanics (Phase 02+).

## 6. Loot

Master loot rule:
- Raid leader distributes boss drops.
- Alternate: personal loot (fair; v1 default).
- Round-robin.

Weekly lockout: raid character-özel.

## 7. Endpoint (Phase 02)

**POST `/v1/guilds/{id}/raid/create`**
**POST `/v1/raid/{id}/enter`**
**POST `/v1/raid/{id}/loot/distribute`** — master loot only.

## 8. Cross-Reference
- [Guild](guild.md), [Boss](boss.md), [Party](../party/party.md)
- [combat/pve-encounters.md](../combat/pve-encounters.md)
- [dungeon/README.md](../dungeon/README.md) — raid alternate to dungeon

## 9. Bekleyen Kararlar
- Aktivasyon Phase 02.
- Raid size final (10 v1 Phase 02; 25 later).
- Difficulty mythic affix Phase 02+.
- Positional mechanics Phase 03.
- Cross-guild raid (public queue) Phase 03.
