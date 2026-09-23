# Character (Kimlik & Yaşam Döngüsü)

## 1. Amaç
Karakter entity'sinin oluşturulması, yaşam döngüsü, sahiplenmesi ve silinmesi. Diğer tüm sistemlerin taşıyıcısı.

## 2. Model
Kaynak: [database/character.md §1](../database/character.md).

Ana alanlar (özet):
- Kimlik: `id`, `account_id`, `name` (CITEXT unique).
- Kimlik seçimi: `race_id`, `class_id`, `spec_tier1_id`, `spec_tier2_id`.
- İlerleme: `level`, `xp`, `ascension_level`.
- Konum: `home_city_id`, `current_location_id`.
- Aktif runtime: `current_activity_id`, `active_pet_id`, `active_mount_id`.
- Meta: `cosmetic_json`, `settings_json`, `flags_json`.
- Sistem: `seed` (RNG kökü), `version` (optimistic lock).

## 3. Yaşam Döngüsü
```
CREATE (account slot check)
  ↓
ACTIVE (login sessions arasında)
  ↓ (isteğe bağlı)
SOFT_DELETE (7 gün grace)
  ↓
HARD_DELETE (worker cleanup)
```

State transitions:
- `soft_deleted_at IS NULL` → aktif.
- `soft_deleted_at IS NOT NULL AND soft_deleted_at > now() - 7d` → grace (restore mümkün).
- `soft_deleted_at < now() - 7d` → hard delete eligible.

## 4. Character Creation

### 4.1 Endpoint
**POST `/v1/characters`**
Auth required (account'a bağlı).

Request:
```json
{
  "name": "Rhaegar",
  "race_id": "r-elf",
  "class_id": "c-mage",
  "cosmetic": {
    "hair_color": "silver",
    "skin_tone": "pale",
    "face_id": 3
  },
  "home_city_id": "city-ashfield"  // opsiyonel; default: race'in başlangıç şehri
}
```

Response 201:
```json
{
  "character": {
    "id": "c-...",
    "name": "Rhaegar",
    "race_id": "r-elf",
    "class_id": "c-mage",
    "level": 1,
    "xp": 0,
    "created_at": "..."
  }
}
```

### 4.2 Validation
- `name`:
  - 3..16 char.
  - alfanumeric + apostrophe (`'`) + hyphen; başta/sonda hyphen yasak; ardışık `--`/`''` yasak.
  - Reserved keyword blacklist (`admin`, `gm`, `system`, `mod`, `mmo`, ...).
  - Profane filter (opsiyonel Phase 02).
  - Unique (CITEXT).
- `race_id` var, aktif.
- `class_id` var, aktif.
- Cosmetic seçenekleri race allowed set içinde.
- Account slot: `COUNT(active characters) < account.character_slots`.

### 4.3 Server Akışı
```
BEGIN;
1. Slot cap check (account.character_slots).
2. Name unique reserve (attempt insert; catch unique violation).
3. Load race + class → derive starter stats + starter skills.
4. Insert character row (with seed = crypto/rand 32 byte).
5. Insert character.stats (allocated_points = 0; primary stats = race+class bias).
6. Insert character.wallet (gold = 0, starting purse config'ten).
7. Grant starter items:
     - Weapon (class starter weapon template)
     - Basic armor set
     - Starter potion x5
     - Starter food x3
     Each item = insert content.item_instance + character.inventory_slot.
8. Grant starter skills (class.starter_skill_ids) → character.skill_learn.
9. Grant starter equipment → character.equipment.
10. Grant starter pet (opsiyonel; class-driven, e.g. ranger'a wolf pup).
11. Grant starter mount (Phase 02; v1'de walk).
12. Grant starter home city (spawn at spawn_point).
13. Grant starter quest (main story first).
14. Insert character.rotation (default per-class rotation).
15. Insert into audit.event_log 'character.created'.
16. Publish event 'character.created'.
COMMIT;
```

Transaction Serializable (name unique + slot cap race safety).

## 5. Character List

**GET `/v1/characters`**
Auth required.

Response:
```json
{
  "characters": [
    {
      "id": "c-1", "name": "Rhaegar",
      "race_code": "elf", "class_code": "mage",
      "level": 25, "ascension_level": 0,
      "last_seen_at": "...",
      "location_summary": "Ashfield"
    }
  ],
  "slots": { "used": 1, "total": 2 }
}
```

## 6. Character Profile (Sheet)

**GET `/v1/characters/{id}`**
Auth + CharacterScope middleware (ownership check).

Response 200 (full sheet):
```json
{
  "character": { ...core fields... },
  "stats": { "primary": {...}, "derived": {...}, "tertiary": {...} },
  "equipment": { "slot": item_instance_id, ... },
  "inventory_summary": { "used_slots": 40, "total_slots": 40 },
  "active_activity": { "id": "...", "type": "mining", "started_at": "..." },
  "professions": [...],
  "active_pet": { "id": "...", "name": "..." },
  "active_mount": null,
  "titles": { "owned_count": 3, "active": "Grandmaster Miner" },
  "reputation": [...],
  "wallet": { "gold_copper": 12345, ... }
}
```

Cache: 30 sn ETag; mutation invalidate.

## 7. Rename

**POST `/v1/characters/{id}/rename`**
Auth + CharacterScope.

Request:
```json
{ "new_name": "NewName" }
```

Server:
1. Validate name (aynı creation kuralı).
2. Cost check: **first rename in 30 days FREE**; sonrası 500 gold veya rename token (Phase 02 monetization).
3. Cooldown: 30 gün per character.
4. Transaction:
   - Debit currency (varsa).
   - Update `character.name` (CITEXT unique violation → 409).
   - `updated_at`, `version+1`.
   - Audit event `character.renamed`.
5. Publish `character.renamed` event (guild/party/friend chat cache invalidate).

Old name reservation:
- Immediate release (yeni bir karakter alabilir).
- Alternatif (Phase 02): 90 gün kilit ("name squatting" önleme).

Response 200: `{ "character": { ...updated... } }`.

## 8. Delete (Soft)

**DELETE `/v1/characters/{id}`**
Auth + CharacterScope.

Body (confirmation):
```json
{ "confirm_name": "Rhaegar" }
```
(kullanıcı ismi elle yazarak onaylar — accidental delete koruma).

Server:
1. Verify `confirm_name == character.name`.
2. Guild leader check → 409 `character.guild_leader_transfer_required`.
3. Active market listing count > 0 → uyarı; kabul edilirse listeler iade (mail).
4. Active party/session cleanup.
5. Set `soft_deleted_at = now()`.
6. Kill active session claim; log out from any character-linked WebSocket.
7. Audit `character.soft_deleted`.
8. Enqueue worker `character_hard_delete_countdown` (30 gün grace, ancak öneri: **7 gün** consistent with database/character.md §20).

Response 204.

## 9. Restore (Cancel Deletion)

**POST `/v1/characters/{id}/restore`**
Auth + CharacterScope.

Server:
1. Character `soft_deleted_at IS NOT NULL AND soft_deleted_at > now() - 7d`.
2. Clear `soft_deleted_at`.
3. Cancel countdown job.
4. Audit.
5. Response 200.

Post 7-gün: 404 (hard deleted).

## 10. Hard Delete (Worker)
Worker `character_hard_delete`:
```
1. Load character (soft_deleted_at < now() - 7d).
2. Anonymize name → "Deleted-<random8>" (unique preserve, isim havuza döner).
3. Delete related rows:
   - character.stats, inventory_slot, bank_slot, equipment, skill_learn, skill_bar,
     rotation, profession, recipe_known, mastery, ascension_choice,
     wallet, reputation, death_log, title_owned.
   - companion.pet_owned, mount_owned (bound to char).
   - progression.quest_progress, achievement_progress.
   - activity.contract (status='archived'), snapshot.
   - social.friendship (both sides), guild_member (already handled at soft).
4. Redirect item_instance owner_kind='void' (owner_id=NULL).
5. Historical rows kalır (audit.event_log, economy.ledger_entry, market_transaction) — ekonomi bütünlüğü.
6. Hard delete character row.
7. Audit 'character.hard_deleted'.
```

## 11. Session Claim (Enter/Exit Karakter)

**POST `/v1/characters/{id}/enter`**
Auth + CharacterScope.

Server:
1. Check current claim: `character.current_session_id`.
   - Aynı session → 200 no-op.
   - Farklı session → kick eski (Redis pub `pub:char:claim:<char_id>` → hub eski WS kick).
2. Set `current_session_id`.
3. Response 200 + character sheet.

**POST `/v1/characters/{id}/exit`**
Auth + CharacterScope.

Server:
1. Clear `current_session_id`.
2. WS unsubscribe character-specific channels.
3. Response 204.

## 12. Ownership Kontrolü (Middleware)
`CharacterScope` middleware (bkz. [backend/middleware.md §3.11](../backend/middleware.md)):
- Route param `{id}` load.
- `character.account_id == ctx.AccountID`.
- `soft_deleted_at IS NULL`.
- Mismatch → 403 `character.forbidden` veya 404.

## 13. Cross-Platform Notlar
- Karakter state cloud-authoritative; her platform aynı sheet.
- UI cache local; ETag ile invalidate.
- Cosmetic seçim server tarafından validate (kilitli cosmetic'ler için).

## 14. Fail Modları
- Slot cap → 403 `character.slot_cap_reached`.
- Isim çakışması → 409 `character.name_taken`.
- Guild leader delete → 409.
- Race/Class not found → 400 `character.invalid_choice`.
- Rename cooldown → 429 `character.rename_cooldown`.

## 15. Rate Limit
- Create: 3/day/account (spam koruma).
- Delete: 5/day/account.
- Restore: 5/day/account.
- Rename: 1 request/hour/character.
- Enter: 20/hour/character.

## 16. Metrics
- `character_created_total{race,class}`.
- `character_deleted_total`.
- `character_restored_total`.
- `character_rename_total{paid}`.
- `character_active_count` gauge (last_seen_at within 24h).

## 17. Audit Events
- `character.created`
- `character.renamed`
- `character.soft_deleted`
- `character.restored`
- `character.hard_deleted`
- `character.entered_session`
- `character.forced_exit` (kick)

## 18. Test Stratejisi
- Unit: name validation matrix.
- Integration: create → sheet → rename → delete → restore.
- Race: concurrent same-name creation → biri success, biri 409.
- Concurrency: two sessions enter same character → second kicks first.

## 19. İlişkiler
- [Level](level.md), [XP](xp.md), [Attributes](attributes.md), [Stats](stats.md), [Health](health.md), [Titles](titles.md), [Reputation](reputation.md)
- [game-design/character.md](../game-design/character.md), [database/character.md](../database/character.md), [auth/session.md](../auth/session.md)

## 20. Bekleyen Kararlar
- Rename cost final (500g v1; token Phase 02).
- Soft-delete grace: 7 gün (character) vs 30 gün (account) tutarlılık.
- Name squatting koruma (delete sonrası isim kilidi).
- Race/Class change ücretli işlem — Phase 02.
