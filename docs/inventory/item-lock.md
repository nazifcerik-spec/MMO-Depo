# Item Lock

## 1. Amaç
Kritik item'ları yanlışlıkla sat/parçala/mail/trade etmekten koruma. Kullanıcı manuel lock/unlock.

## 2. Model

Schema ekleme (Phase 11):
```sql
ALTER TABLE content.item_instance
  ADD COLUMN locked BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN locked_at TIMESTAMPTZ NULL,
  ADD COLUMN locked_reason VARCHAR(48) NULL;
```

## 3. Kilit Etkisi

Locked item için engellenen aksiyonlar:
- Salvage
- Vendor sell
- Destroy
- Mail send (attach)
- Marketplace list
- Guild bank deposit
- Bank deposit (opsiyonel; default allowed — güvenli storage)
- Equipment change:
  - Unequip: BLOKLU (equipped locked item).
  - Equip başka item (aynı slot): BLOKLU.

İzin verilen aksiyonlar:
- Inventory içi move
- Character bank deposit/withdraw (opt-in; config)
- Consume (potion, food) — opsiyonel (kullanıcı istese de food consume edebilsin; config)

Config:
```yaml
item_lock:
  bank_deposit_allowed: true
  consume_allowed: false
  equip_lock_prevents_unequip: true
```

## 4. Endpoint

### 4.1 Toggle Lock

**POST `/v1/characters/{id}/inventory/lock`**
Auth + CharacterScope.
```json
{
  "actions": [
    { "item_instance_id": "ii-...", "lock": true, "reason": "trophy" }
  ]
}
```

Server:
- Validate ownership.
- Bulk lock/unlock atomic.
- `locked=true`, `locked_at=now()`, `locked_reason=input`.
- Emit event `item.lock_toggled`.

### 4.2 Query Locked Items

**GET `/v1/characters/{id}/inventory/locked`**
```json
{ "items": [ { "item_instance_id": "...", "template_code": "...", "locked_at": "...", "reason": "..." } ] }
```

## 5. Auto-Lock (Content-Driven)

Bazı item'lar creation'da otomatik lock:
- `content.item_template.auto_lock_on_pickup BOOLEAN`.
- Örnek: quest item, unique legendary (first drop).

Yeni schema alanı (opsiyonel):
```sql
ALTER TABLE content.item_template
  ADD COLUMN auto_lock_on_pickup BOOLEAN NOT NULL DEFAULT false;
```

Item Engine `Create` sırasında bakar; auto-lock true ise `locked=true` başlangıçtan.

Kullanıcı manual unlock edebilir.

## 6. Reason Kataloğu

- `trophy` — kullanıcı tercih (keepsake).
- `quest_item` — quest için gerekli.
- `unique_first` — ilk drop legendary (auto).
- `expensive` — pahalı item, aksidental sale koruma.
- `mount_component` — mount unlock item.
- `set_progress` — set collection için.

## 7. UI

Inventory tooltip'te "🔒 Locked" badge.
Right-click context menu → Unlock.
Bulk lock/unlock UI (multi-select).

## 8. Auto-Lock Config

Client-side opt-in patterns (v1'de basit):
- "Auto-lock rare+ items" toggle.
- "Auto-lock unique effect items" toggle.

Client bu preference'i lock endpoint ile server'a yansıtır (item drop notification flow).

Phase 02 opsiyonu: server-side auto-lock rule engine (persist in account settings).

## 9. Lock Persistence

- Ownership transfer edilirse (mail, trade v2) lock durumu **kaldırılır**.
- Character silme + restore: lock durumu korunur.
- Content patch: mevcut locked item locked kalır.

## 10. Fail Modları

Locked item aksiyonu:
```
400 code: item.locked
message: "This item is locked. Unlock it first."
details: { instance_id: "...", locked_reason: "..." }
```

## 11. Race Conditions

- Concurrent lock toggle → optimistic lock version.
- Lock + concurrent sell → sell reject with `item.locked`.

## 12. Rate Limit
- Lock/unlock: 60/dk.
- Bulk (>10 item): 10/dk.

## 13. Metrics
- `item_locked_count` gauge (aktif locked item sayısı).
- `item_lock_toggle_total{action}` counter.
- `item_action_blocked_by_lock_total{action}` counter.

## 14. Audit
- `item.locked` (severity=debug).
- `item.unlocked`.
- `item.action_blocked_by_lock` (severity=info; user awareness).

## 15. Test Stratejisi
- Lock + attempt sell → reject.
- Lock + attempt destroy → reject.
- Lock + move within inventory → allowed.
- Bulk toggle atomic.
- Auto-lock on pickup.

## 16. Cross-References
- [Inventory](inventory.md), [Item Transfer](item-transfer.md), [Equipment](equipment.md)
- [item-engine/item-instance.md](../item-engine/item-instance.md)
- [database/character.md](../database/character.md)

## 17. Bekleyen Kararlar
- Auto-lock rule engine server-side (Phase 02).
- Lock inherit on trade (v1 kaldır; Phase 02 opt-in).
- Special "loot filter" (server-side; auto-destroy junk based on rules) — Phase 02.
- Guild bank locked item (guild-scope lock) Phase 02.
