# Inventory Validation

## 1. Amaç
Item-etkileşen her aksiyon (equip, sell, list, mail, trade, delete, consume) için ownership + lock + bind + cap kontrolü.

## 2. Check Matrisi (Her Aksiyon Öncesi)

| Check | Amaç |
|---|---|
| **Ownership** | `item_instance.owner_character_id == self` |
| **Lock state** | `is_locked=false` (aksi taksrif `lock_owner` göre error) |
| **Bind rule** | `bind_kind` uygunsa aksiyona (BoP → trade edilemez, BoA → sadece account içi) |
| **No-sell flag** | `is_no_sell` → vendor sell yasak |
| **No-trade flag** | `is_no_trade` → trade/market yasak |
| **Level req** | Equip için `character.level >= item.level_req` |
| **Class req** | Equip için class uygun (`item.class_restriction`) |
| **Slot valid** | Equip için `slot == item.slot_kind` |
| **Slot cap** | Envanter transfer için hedef slot boş / stack birleşimi mümkün |
| **Weight cap** | Envanter kabul için toplam weight ≤ karakter cap |
| **Stack rule** | Stackable ise adet ≥ 1, ≤ max_stack; non-stackable → 1 |
| **Instance existence** | Not soft-deleted, ledger tutarlı |

## 3. Lock Types

Bkz. [inventory/item-lock.md](../inventory/item-lock.md):
- `market` — listeleme escrow.
- `trade` — aktif trade session.
- `mail` — mail attachment escrow.
- `craft` — crafting input rezerve.
- `admin` — freeze/forensic.

Lock varken mutasyon → `E_ITEM_LOCKED{owner}`.

## 4. Race Condition Guard

- İki taraf aynı anda aynı item için işlem → SERIALIZABLE TX + `SELECT ... FOR UPDATE`.
- Item movement TX atomic; ownership değişir ve eski slot boşalır aynı commit'te.

## 5. Instance Fingerprint (Dupe Detection)

Her item instance:
- `id UUID` unique (dupe UUID insert impossible → PK).
- `crafted_signature` (crafter_id + timestamp + recipe) opsiyonel (mastery signature).
- `looted_from` (source_ref) — forensic.

Bir instance UUID iki karakter envanterinde asla olmamalı (partial index / unique constraint). Reconciler kontrol eder.

## 6. Batch Aksiyonlar

Batch sell / batch delete → tek TX; her instance için tek tek check + rollback all-or-nothing.

## 7. Envanter Cap Rules

- Slot count: karakter bazlı (level ile artabilir Phase 02).
- Weight cap: karakter STR + gear ile artabilir.
- Cap dolarken pickup → overflow (mail queue veya destroy prompt).

## 8. Fraud Signal

- **Impossibly rapid item flow**: aynı karakter 1 saat içinde 100+ item transfer → RMT flag.
- **Impossibly rare loot rate**: characterin loot geçmişinde legendary rate 10× beklenen → potential exploit.
- **Bind rule bypass attempt**: BoP item trade endpoint'e gelirse → E_ITEM_NOT_TRADABLE + `sec_bind_bypass_attempt_total`.

## 9. Metrikler
- `sec_inventory_check_reject_total{check_kind}`
- `sec_bind_bypass_attempt_total`
- `sec_item_lock_conflict_total{lock_owner}`
- `sec_dupe_instance_detected_total` (reconciler)

## 10. Cross-Reference
- [inventory](../inventory/README.md), [inventory/item-lock.md](../inventory/item-lock.md), [inventory/item-transfer.md](../inventory/item-transfer.md)
- [item-engine](../item-engine/README.md)

## 11. Bekleyen Kararlar
- Rate-based RMT signal Phase 02 balance.
- Reconciler cycle sıklığı (v1 nightly; Phase 02 hourly opt).
- Overflow queue (envanter dolu) tam politika Phase 02.
