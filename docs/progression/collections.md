# Collections

## 1. Amaç
Kalıcı unlock galerileri — pet, mount, appearance (transmog), mount skin, tabard, cosmetic. Long-term "gotta catch 'em all" motivasyon.

## 2. Collection Türleri (v1)

| Tür | Kaynak | Set Bonus? |
|---|---|---|
| **Pet species** | Hunt/tame, quest, event, achievement | ❌ (v1); Phase 02 set-of-species bonus |
| **Mount species** | Quest, event, faction, achievement | ❌ (v1) |
| **Appearance (transmog)** | Item pickup — item template ilk kez looted / equiped | Set unlock (tam set giyme cosmetic) |
| **Mount skin** | Store, seasonal, achievement | ❌ |
| **Tabard** | Guild, faction, achievement | ❌ |
| **Weapon effect** (Phase 02) | Enchant proc visual | ❌ |
| **Housing decoration** (Phase 03) | Quest, craft, event | ❌ |

## 3. Model

```sql
CREATE TABLE content.collection_category (
  id UUID PRIMARY KEY,
  code VARCHAR(48) UNIQUE NOT NULL,      -- pet|mount|appearance|mount_skin|tabard|...
  display_name_key VARCHAR(64) NOT NULL,
  scope VARCHAR(16) NOT NULL,            -- 'account'|'character'
  set_bonus_supported BOOLEAN NOT NULL DEFAULT false,
  version INT DEFAULT 1
);

CREATE TABLE content.collection_entry (
  id UUID PRIMARY KEY,
  category_id UUID NOT NULL,
  code VARCHAR(64) UNIQUE NOT NULL,      -- "pet_wolf_common", "appearance_sword_iron"
  display_name_key VARCHAR(64) NOT NULL,
  rarity VARCHAR(16) NOT NULL,
  source_kind VARCHAR(24) NOT NULL,      -- hunt|quest|event|achievement|store|craft
  source_hint JSONB NULL,                -- ör. { region: 'glenmoor', season: 'winter' }
  set_code VARCHAR(48) NULL,             -- appearance set için
  cosmetic_ref VARCHAR(128) NULL,        -- art asset
  is_retired BOOLEAN NOT NULL DEFAULT false,
  version INT DEFAULT 1
);

CREATE TABLE progression.account_collection (
  account_id UUID NOT NULL,
  collection_entry_id UUID NOT NULL,
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  source_ref_kind VARCHAR(24) NULL,
  source_ref_id UUID NULL,
  PRIMARY KEY (account_id, collection_entry_id)
);

CREATE TABLE progression.character_collection (
  character_id UUID NOT NULL,
  collection_entry_id UUID NOT NULL,
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, collection_entry_id)
);
```

Scope kararı category başına; pet/mount/mount_skin/tabard → **account** (paylaşımlı); appearance → **character** (rol içi immersion) v1 seçim; Phase 02 review.

## 4. Unlock Akışı

Trigger event bus:
```
event.progression.collection_unlock_hint
  { account_id or character_id, entry_code, source_ref_kind, source_ref_id }
```

Worker `collection_apply`:
1. Zaten unlocked → no-op.
2. `progression.*_collection` insert.
3. Achievement trigger (collection sayısı milestone → achievement).
4. Set completion check (tabana set entry'nin tümü unlocked → set bonus / cosmetic reward).
5. WS "collection_unlocked" push + inbox notification.

## 5. Appearance (Transmog v1 Temel)

- Karakter item **equip veya loot** ilk kez → `character_collection` insert appearance entry.
- Transmog UI karakter için: kilit açık appearance'lardan seçim yaparak slot skin change (görsel değiştir, stat aynı kalır — Phase 02 tam feature).
- V1'de sadece **unlock tracking** var; transmog apply Phase 02.

## 6. Set Bonus (Appearance)

Content:
```yaml
- set_code: "iron_battlegear"
  entries: [ chest_iron, legs_iron, gloves_iron, boots_iron, helm_iron ]
  reward:
    - { title_code: "the_ironclad" }
    - { cosmetic_sku: "iron_aura" }
    - { achievement_code: "set_iron_complete" }
```

Set complete → reward apply.

## 7. Rarity Ekonomisi

- Rare/epic entry kaynağı hunt / low-drop; event limited.
- Retired entry: yeni oyuncu unlock edemez (seasonal, alpha).

## 8. Endpoint

**GET `/v1/account/{id}/collection?category=&rarity=&unlocked=`**
**GET `/v1/character/{id}/collection?category=&rarity=&unlocked=`**
**GET `/v1/collection/entry/{code}`** — detail + source hint (kullanıcının yönlendirilmesi için).

## 9. Metrikler
- `collection_unlock_total{category, rarity}`
- `collection_set_complete_total{set_code}`
- `collection_retired_gauge`

## 10. Cross-Reference
- [pets/README.md](../pets/README.md), [mounts/README.md](../mounts/README.md), [item-engine/README.md](../item-engine/README.md)
- [achievements.md](achievements.md), [titles.md](titles.md)

## 11. Bekleyen Kararlar
- Transmog full feature Phase 02.
- Housing decorations Phase 03.
- Weapon effect proc Phase 02.
- Set bonus non-appearance kategoriler (pet set-of-5 bonus vs) Phase 02.
