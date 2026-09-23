# Calculate Rewards

## 1. Amaç
`Yield` struct → concrete ledger + inventory + XP dispatch. Idempotent, atomic, event-emitting.

## 2. Yield → Reward Dispatcher

```go
package rewarddispatcher

func Dispatch(ctx context.Context, tx db.TX, char Character, contract Contract, yield Yield, idempotencyKey string) error {
    // 1. XP
    if yield.XPCharacter > 0 {
        _ = characterService.GrantXP(ctx, tx, char.ID, XPSource{
            Kind: "activity_reveal",
            RefID: contract.ID,
            BaseAmount: yield.XPCharacter,
            IdempotencyKey: idempotencyKey + ":char_xp",
        })
    }
    for profID, xp := range yield.XPProfession {
        _ = professionService.GrantXP(ctx, tx, char.ID, profID, xp, idempotencyKey+":prof_"+profID.String())
    }
    for trackID, xp := range yield.XPMastery {
        _ = masteryService.GrantXP(ctx, tx, char.ID, trackID, xp, idempotencyKey+":mastery_"+trackID.String())
    }

    // 2. Currency
    for currency, delta := range yield.CurrencyDeltas {
        _ = economyService.CreditWallet(ctx, tx, char.ID, currency, delta, idempotencyKey+":cur_"+currency)
    }

    // 3. Reputation
    for factionID, delta := range yield.ReputationDeltas {
        _ = reputationService.Grant(ctx, tx, char.ID, factionID, delta, idempotencyKey+":rep_"+factionID.String())
    }

    // 4. Item drops
    for _, drop := range yield.ItemDrops {
        _ = distributeItemDrop(ctx, tx, char, drop, idempotencyKey)
    }

    // 5. Item consumptions (crafting inputs, bait)
    for _, consume := range yield.ItemConsumptions {
        _ = consumeInventoryItem(ctx, tx, char, consume, idempotencyKey)
    }

    // 6. Death handling
    if yield.DeathCount > 0 {
        _ = characterService.RegisterDeaths(ctx, tx, char.ID, contract.ID, yield.DeathCount, yield.XPDebt, idempotencyKey)
    }

    return nil
}
```

Tüm dispatch **single transaction** (reveal wrapper); partial fail → rollback.

## 3. Item Drop Distribution

```go
func distributeItemDrop(ctx, tx, char, drop, idemKey) error {
    tmpl, _ := content.ItemTemplates.ByCode(drop.TemplateCode)

    // Non-stackable (equipment, unique roll) → her adet ayrı instance
    if tmpl.StackMax == 1 {
        for i := int32(0); i < drop.Quantity; i++ {
            createReq := item.CreateRequest{
                TemplateID: tmpl.ID,
                Owner: item.OwnerRef{Kind: "character_inventory", ID: char.ID},
                Seed: seed.Sub(idemKey + fmt.Sprintf(":drop-%s-%d", drop.TemplateCode, i)),
                Source: item.SourceLoot,
            }
            instance, err := itemEngine.Create(ctx, tx, createReq)
            if err != nil {
                if errors.Is(err, ErrInventoryFull) {
                    // Overflow to mailbox
                    _ = overflowMailbox.Add(ctx, tx, char.ID, instance)
                } else {
                    return err
                }
            }
        }
        return nil
    }

    // Stackable material — inventory merge or create
    remaining := drop.Quantity
    for remaining > 0 {
        placed, err := inventoryService.PutMaterial(ctx, tx, char.ID, tmpl.ID, remaining)
        if err != nil {
            return err
        }
        if placed == 0 {
            // Inventory full — overflow
            _ = overflowMailbox.AddMaterial(ctx, tx, char.ID, tmpl.ID, remaining)
            break
        }
        remaining -= placed
    }
    return nil
}
```

Inventory full handling: overflow mailbox (bkz. [inventory/inventory.md §7](../inventory/inventory.md)).

## 4. Item Instance vs Material Stack

- **Equipment / instance-özel**: her drop yeni instance (stack_max=1).
- **Material**: aynı template stack merge.

Bkz. [inventory/stacking.md §7](../inventory/stacking.md).

## 5. Currency Dispatch

Ledger append-only:
```
INSERT INTO economy.ledger_entry (character_id, currency_code, delta, balance_after, reason, source_kind, source_ref_id, idempotency_key)
VALUES (char_id, 'gold', +N, new_balance, 'credit_activity', 'activity_reveal', contract_id, idem_key);
```

Aynı `idempotency_key` unique → duplicate no-op safe.

Bkz. [database/economy-market.md §1](../database/economy-market.md).

## 6. Reputation Dispatch

```
UPDATE character.reputation SET value = value + delta ...
Tier transition check → event 'character.reputation_tier_changed'.
Cap enforce (42000).
```

Bkz. [character/reputation.md §8](../character/reputation.md).

## 7. XP Debt (Death)

Death event → xp_debt accumulate:
```
if char.level >= 20:
    char.xp_debt += level_xp_required(char.level) × 0.05 × death_count
    cap: xp_debt <= level_xp_required(char.level)
```

Bkz. [character/xp.md §7](../character/xp.md).

## 8. Consumables Consumption (Crafting, Bait)

```go
func consumeInventoryItem(ctx, tx, char, consume, idemKey) error {
    remaining := consume.Quantity
    // Iterate inventory slots with this template
    slots := inventoryService.FindMaterialSlots(ctx, char.ID, consume.TemplateCode)
    for _, s := range slots {
        if remaining <= 0 { break }
        take := min(s.StackCount, remaining)
        _ = inventoryService.TakeMaterial(ctx, tx, char.ID, s.SlotIndex, take)
        remaining -= take
    }
    if remaining > 0 {
        return errs.Wrap(ErrConsumeShortage, fmt.Sprintf("%s short by %d", consume.TemplateCode, remaining))
    }
    return nil
}
```

Shortage handling: reveal öncesi threshold detect etti; shortage rare edge case; error → contract paused.

## 9. Reveal Payload Response

Client'a döner:
```json
{
  "reveal_id": "...",
  "contract_id": "...",
  "elapsed_ms": 21600000,        // 6h
  "capped_ms": 21600000,          // full (under soft cap)
  "summary": {
    "character_xp_gained": 40545,
    "profession_xp_gained": { "hunting": 8109 },
    "mastery_xp_gained": { "bestiary": 8109 },
    "gold_gained_copper": 12500,
    "items_gained": [
      { "template_code": "leather_scrap", "quantity": 2200 },
      { "template_code": "wolf_fang", "quantity": 600 },
      { "template_code": "alpha_pelt", "quantity": 5 }
    ],
    "deaths": 0,
    "notable_events": [
      { "kind": "rare_drop", "item": "ancient_wolf_fang" }
    ]
  },
  "new_status": "active",
  "notification_hint": "6 hours of hunting revealed."
}
```

Client "reveal panel" UI'yı bu payload'la doldurur.

## 10. Overflow Mailbox

Envanter dolduğunda:
- Rare/BoP: system mail'e attach (7 gün TTL).
- Notification "Rare item awaits in mail".

Detay: [inventory/inventory.md §7](../inventory/inventory.md).

## 11. Idempotency

Reveal reveal_id per-reveal unique. Idempotency key = correlation ID.

Per-source key:
- `reveal:<contract_id>:<reveal_at>:char_xp`
- `reveal:<contract_id>:<reveal_at>:prof_<id>`
- ...

Aynı reveal 2 kez → hepsi no-op (ledger unique + xp_ledger unique).

## 12. Fail Modları

- Insufficient inventory + overflow disabled → 500 rare; retry.
- Ledger conflict → 40001 retry.
- Item template deprecated (rare) → drop yerine safe fallback (junk).
- Wallet overflow (unlikely BIGINT) → clamp + audit.

## 13. Rate Limit

Reveal max 1 per second per character (bkz. [systems/afk-engine.md §12](../systems/afk-engine.md)).

## 14. Test Stratejisi

- Full dispatch happy path (XP + gold + item + rep).
- Overflow mailbox trigger.
- Consumable shortage.
- Idempotent double-reveal.
- Rollback on partial fail.

## 15. Metrics

- `reward_dispatched_total{kind}`.
- `reward_dispatch_duration_seconds`.
- `reward_overflow_mailbox_total`.
- `reward_dispatch_rollback_total`.

## 16. Audit

- `activity.reveal_dispatched` (severity=debug; high volume).
- `activity.overflow_mailbox` (severity=info).
- `reward.item_created` (via item engine).

## 17. Cross-Reference

- [Activity](activity.md), [Offline Progress](offline-progress.md), [XP](xp.md), [Lifecycle](lifecycle.md)
- [systems/reward-system.md](../systems/reward-system.md), [systems/loot-system.md](../systems/loot-system.md)
- [item-engine/README.md](../item-engine/README.md)
- [inventory/README.md](../inventory/README.md)
- [character/xp.md](../character/xp.md), [character/reputation.md](../character/reputation.md), [character/health.md](../character/health.md)
- [database/economy-market.md §1](../database/economy-market.md)

## 18. Bekleyen Kararlar

- Overflow mailbox cap tam davranış (mail cap dolarsa → void with alarm).
- Bulk item instance insert batch (Phase 02 optimization).
- Reveal payload versioning (`_v` alanı client compat).
