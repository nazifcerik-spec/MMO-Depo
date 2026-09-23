# Database — `companion` Schema

## Kapsam
Pet ve mount **instance**'ları (template'ler `content.pet_template` / `content.mount_template`). Sahiplenilmiş companion'lar, level/bond ilerlemesi, koleksiyon.

---

## 1. `companion.pet_owned`

Karakter sahipli pet instance.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `owner_character_id` | `UUID` | FK → character.character NOT NULL |
| `template_id` | `UUID` | FK → content.pet_template NOT NULL |
| `custom_name` | `VARCHAR(32)` | NULL |
| `level` | `SMALLINT` | NOT NULL DEFAULT 1 |
| `xp` | `BIGINT` | NOT NULL DEFAULT 0 |
| `bond_level` | `SMALLINT` | NOT NULL DEFAULT 0 — Phase 02 |
| `is_active` | `BOOLEAN` | NOT NULL DEFAULT false — bir seferde birden fazla true olabilir (combat+utility slot ayrı) |
| `active_slot` | `companion.pet_slot` | NULL — combat / utility |
| `stable_since` | `TIMESTAMPTZ` | NULL — retired ise |
| `acquired_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `acquired_source` | `VARCHAR(32)` | NOT NULL — quest/drop/event/purchase/starter |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Enum:
```sql
CREATE TYPE companion.pet_slot AS ENUM ('combat','utility');
```

Constraint (application enforce; opsiyonel partial unique):
```sql
CREATE UNIQUE INDEX uq_pet_active__owner_slot
ON companion.pet_owned (owner_character_id, active_slot)
WHERE is_active = true;
```

Indexler:
- `ix_pet_owned__owner_character_id`.
- `ix_pet_owned__template_id`.

Notlar:
- Bound to character (BoP, DD-024); trade edilemez.
- Pet ölmez; combat'ta downed → cooldown; state Redis (ephemeral).

---

## 2. `companion.pet_food_buff` (opsiyonel)

Pet food consumable ile aktif buff.

```
character_id UUID PK-part,
pet_owned_id UUID PK-part,
buff_id UUID NOT NULL,
expires_at TIMESTAMPTZ NOT NULL
```

---

## 3. `companion.mount_owned`

Karakter sahipli mount.

| Kolon | Tip | Constraint |
|---|---|---|
| `id` | `UUID` | PK |
| `owner_character_id` | `UUID` | FK NOT NULL |
| `template_id` | `UUID` | FK → content.mount_template NOT NULL |
| `is_active` | `BOOLEAN` | NOT NULL DEFAULT false |
| `custom_name` | `VARCHAR(32)` | NULL |
| `acquired_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `acquired_source` | `VARCHAR(32)` | NOT NULL |
| `version` | `INT` | NOT NULL DEFAULT 1 |

Constraint:
```sql
CREATE UNIQUE INDEX uq_mount_active__owner
ON companion.mount_owned (owner_character_id)
WHERE is_active = true;
```

Notlar: BoP.

---

## 4. `companion.mount_barding` (Phase 02)

Mount equipment (speed, cargo mods).

```
mount_owned_id UUID PK,
barding_item_instance_id UUID FK → content.item_instance NULL
```

---

## 5. `companion.pet_compendium` (Phase 02)

Karakterin sahiplenmiş olduğu tüm pet template'leri (koleksiyon achievement).

```
character_id UUID PK-part,
pet_template_id UUID PK-part,
first_owned_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 6. `companion.mount_compendium` (Phase 02)

```
character_id UUID PK-part,
mount_template_id UUID PK-part,
first_owned_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

---

## 7. `companion.pet_downed_state`

Live-only preferred (Redis); persist yalnızca cross-session karma önlemek için.

```
pet_owned_id UUID PK,
downed_at TIMESTAMPTZ NULL,
revive_at TIMESTAMPTZ NULL
```

v1'de tercih: Redis-only; DB'ye yazmayabiliriz.

---

## 8. Kullanım Örnekleri

**Aktif pet + mount karakter için:**
```sql
SELECT
  (SELECT id FROM companion.pet_owned WHERE owner_character_id = $1 AND is_active AND active_slot='combat') AS combat_pet_id,
  (SELECT id FROM companion.pet_owned WHERE owner_character_id = $1 AND is_active AND active_slot='utility') AS utility_pet_id,
  (SELECT id FROM companion.mount_owned WHERE owner_character_id = $1 AND is_active) AS active_mount_id;
```

**Pet activation swap:**
```sql
BEGIN;
UPDATE companion.pet_owned SET is_active=false, active_slot=NULL, version=version+1
WHERE owner_character_id = $1 AND active_slot = 'combat' AND is_active;

UPDATE companion.pet_owned SET is_active=true, active_slot='combat', version=version+1
WHERE id = $2 AND owner_character_id = $1;
COMMIT;
```

**Pet leveling (activity yield beslendiğinde):**
```sql
UPDATE companion.pet_owned
SET xp = xp + $delta, updated_at = now()
WHERE id = $1;
-- Level up check application side; threshold triggered UPDATE ile level += 1.
```

---

## 9. İlişkiler
- Pet/mount owned → character (owner), content template.
- Aktif pet AFK combat yield hesabına ekleniyor: `activity.contract` snapshot içinde pet stats bulunur.
- Aktif mount travel aktivitesi hız modifier'ı olarak kullanılır.

---

## 10. Notlar
- **Ölmez pet**: DB'de kalıcı state değişimi ölümde yok; sadece downed timer (Redis).
- **Rename**: `custom_name` unique değil (aynı isimli farklı pet olabilir); trademark filter application-side.
- **Retirement** vs delete: retirement (stable) tercih; delete yok. Player attachment koruma.
- **Cosmetic pets** stat vermez; sadece görsel.

---

## 11. Bekleyen Kararlar
- Bond level formulü ve unlock ödülleri (Phase 02).
- Pet battle mini-game data modeli (Phase 03).
- Account-wide pet/mount BoA (Phase 03) — `companion.pet_owned.owner_kind` genişlemesi.
