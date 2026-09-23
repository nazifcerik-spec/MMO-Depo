# Instance (Owned Pet)

## 1. Amaç
Karakterin sahiplendiği pet — level, XP, traits, custom name, active slot, bond level.

Kaynak: [database/companions.md §1](../database/companions.md).

## 2. Model

Phase 20 eklemeleri:
```sql
ALTER TABLE companion.pet_owned
  ADD COLUMN happiness SMALLINT NOT NULL DEFAULT 50,      -- 0..100
  ADD COLUMN stamina_current SMALLINT NOT NULL DEFAULT 100,
  ADD COLUMN stamina_max SMALLINT NOT NULL DEFAULT 100,
  ADD COLUMN stamina_last_recovery_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ADD COLUMN traits_json JSONB NOT NULL DEFAULT '[]',      -- rolled trait list
  ADD COLUMN skill_ranks_json JSONB NOT NULL DEFAULT '{}',  -- {skill_code: rank}
  ADD COLUMN last_feed_at TIMESTAMPTZ NULL,
  ADD COLUMN evolved_from_instance_id UUID NULL;
```

## 3. Instance Yaratımı

Pet acquire (loot/quest/event):
```
1. Load species template.
2. Create pet_owned:
   - owner_character_id = char.id
   - template_id = species.id
   - level = 1
   - xp = 0
   - bond_level = 0
   - happiness = 50 (baseline)
   - stamina_current = stamina_max = 100
   - Roll traits (from species.base_traits_pool)
   - skill_ranks = { each base_skill_code: 1 }
3. Insert row.
4. Emit event 'pet.acquired'.
5. Notification "New pet: <species> named 'Unnamed' — set name".
```

## 4. Traits Roll

Bkz. [traits.md](traits.md).

- Species trait pool → roll N traits (default 2).
- Rarity dependent count (common 1, epic 2-3, legendary 3-4).
- Trait seed-based deterministic.

## 5. Custom Name

**POST `/v1/characters/{id}/pets/{instance_id}/rename`**
- 1-32 char.
- Not unique (multiple pets same name OK).
- No profanity filter (Phase 02).

## 6. Active Slot

Bkz. [database/companions.md §1 partial unique constraint](../database/companions.md).

- `is_active + active_slot` combo.
- 2 slot per character (combat + utility).
- Partial unique index (aktif slot başına 1 pet).

Activation:
```
POST /v1/characters/{id}/pets/{instance_id}/activate
  body: { slot: "combat" | "utility" }

Server:
1. Load pet.
2. Deactivate current slot pet (if any).
3. Activate new pet.
4. Emit event.
```

## 7. Retire

`stable` (dormant) — inactive but preserved:
- `stable_since` timestamp set.
- Removed from active slots.
- Character sheet gizli (opsiyonel toggle).
- Restore: activate again → stable_since NULL.

Delete YOK — pet attachment koruma (DD-024 tercih).

## 8. Bond Level (Phase 02)

Bond = karakter-pet ilişki derinliği:
- Long-term active pet → bond artar.
- Milestone: bond 25 (companion), 50 (friend), 75 (bonded), 100 (soul-bound).
- Unlock: extra passive, unique cosmetic, evolution requirement.

v1'de yok (schema hazır); Phase 02 aktive.

## 9. Downed State (Combat)

Bkz. [game-design/pets.md](../game-design/pets.md).

Combat'ta HP 0:
- Pet downed → combat contribution 0.
- Cooldown 30 sec.
- Auto-revive to full HP.
- Ölmez; kalıcı hasar yok.

State Redis-based (`pet:downed:<instance_id>` TTL 30s).

## 10. Endpoint

**GET `/v1/characters/{id}/pets`** — owned pets list.
```json
{
  "pets": [
    {
      "instance_id": "...",
      "species_code": "wolf_cub",
      "custom_name": "Fenrir",
      "level": 15,
      "xp": 3200,
      "next_level_xp": 5000,
      "happiness": 75,
      "stamina": 85,
      "traits": ["keen_senses", "swift_paw"],
      "is_active": true,
      "active_slot": "combat"
    }
  ]
}
```

**GET `/v1/characters/{id}/pets/{instance_id}`** — pet detail.

## 11. Concurrency

Optimistic lock on `version`.

## 12. Metrics
- `pet_owned_total{species}`.
- `pet_active_slots_used{slot}` gauge.
- `pet_downed_total`.

## 13. Cross-Reference
- [Species](species.md), [Level](level.md), [Traits](traits.md), [Combat](combat.md)
- [database/companions.md](../database/companions.md)

## 14. Bekleyen Kararlar
- Bond level Phase 02 aktivasyon.
- Pet delete opsiyonu (v1 hayır; Phase 02 opt).
- Stable UI Phase 02.
- Cross-character transfer BoA Phase 03.
