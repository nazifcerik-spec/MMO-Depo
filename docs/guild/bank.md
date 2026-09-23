# Bank (Guild Bank)

## 1. Amaç
Shared item + gold storage. Tab-based; rank permissions per tab.

Foundation: [database/social-guilds-parties.md §4-6](../database/social-guilds-parties.md).

## 2. Model

```sql
social.guild_bank_tab
  guild_id, tab_index (PK)
  name, capacity (default 100)
  purchased_at, purchase_cost_copper
  permissions_json

social.guild_bank_slot
  guild_id, tab_index, slot_index (PK)
  item_instance_id, stack_count

social.guild_bank_log
  id (bigserial), guild_id, character_id
  action (deposit/withdraw/move/gold_deposit/gold_withdraw)
  tab_index, item_instance_id, qty, gold_copper
  at
```

## 3. Tab Unlock

Guild level → tab unlock:
- Level 1: tab 0 unlocked (free).
- Level 3: tab 1 available (25000 gold).
- Level 5: tab 2 (100000 gold).
- Level 7: tab 3 (500000 gold).
- Level 10: tab 4 (2000000 gold).
- Level 15+ (Phase 02): tab 5+.

## 4. Endpoints

### 4.1 Deposit

**POST `/v1/guilds/{id}/bank/deposit`**
```json
{
  "moves": [
    { "from_inv_slot": 5, "to_tab": 1, "to_slot": 10, "count": 1 }
  ]
}
```

Server:
1. Permission check: `bank_tab_write[1]`.
2. Validate item eligible (BoP block).
3. Transaction:
   - Item transfer inv → guild_bank_slot.
   - Log guild_bank_log.
   - Emit event.

### 4.2 Withdraw

**POST `/v1/guilds/{id}/bank/withdraw`**
```json
{ "moves": [ ... ] }
```

Server:
1. Permission check.
2. Withdraw cap check (per-rank daily; bkz. [permissions.md §5](permissions.md)).
3. Transaction:
   - Item transfer guild_bank_slot → inv.
   - Log.

### 4.3 Gold Deposit / Withdraw

Similar; gold via ledger.

Gold sink: guild treasury (shared).

## 5. BoP Bank Rule

`bind_type='bop'` items:
- **v1**: cannot deposit BoP items to guild bank.
- Phase 02: BoP allowed but only original owner can withdraw.

## 6. Purchase Bank Tab

**POST `/v1/guilds/{id}/bank/tab/purchase`**
Officer+ only.

Server:
1. Validate next tab available.
2. Debit treasury (or leader personal — v1 treasury).
3. Insert guild_bank_tab.
4. Emit event.

## 7. Tab Rename & Permissions

Leader/officer:
- Rename tab (32 char).
- Edit tab permissions (rank read/write per tab).

## 8. Audit Log

`social.guild_bank_log` append-only:
- Every deposit/withdraw/move logged.
- Officer can view; audit forensic.
- Retention 12 months.

## 9. Bank Contents Query

**GET `/v1/guilds/{id}/bank`**
```json
{
  "tabs": [
    {
      "tab_index": 0,
      "name": "General",
      "capacity": 100,
      "used": 45,
      "read_allowed": true,
      "write_allowed": true,
      "slots": [ ... ]
    }
  ],
  "treasury_gold_copper": 5000000,
  "next_tab_available": true,
  "next_tab_cost_copper": 100000000
}
```

## 10. Concurrent Access

Optimistic lock per slot (updated_at).
Concurrent deposit same slot → 409 retry.

## 11. Guild Treasury

Shared gold pool:
- Members deposit / officer withdraw (cap-based).
- Used for: bank tab purchase, guild research funding, guild boss summon fee, guild quest reward pool.

Treasury drain:
- Auto: guild boss summon (10000 gold per boss).
- Manuel: officer withdraw for guild event.

## 12. Metrics
- `guild_bank_deposit_total{tab}`.
- `guild_bank_withdraw_total{tab}`.
- `guild_treasury_gold_avg` gauge.

## 13. Audit
- `guild.bank_deposit`, `guild.bank_withdraw`, `guild.bank_tab_purchased`, `guild.treasury_transaction`.

## 14. Cross-Reference
- [Guild](guild.md), [Permissions](permissions.md), [Research](research.md), [Boss](boss.md)
- [database/social-guilds-parties.md §4-6](../database/social-guilds-parties.md)

## 15. Bekleyen Kararlar
- BoP allow with owner-lock Phase 02.
- Bank slot capacity per tab (v1: 100).
- Treasury audit UI (Phase 02).
- Cross-alliance shared bank Phase 03.
