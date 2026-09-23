# Loot (Party Loot Distribution)

## 1. Amaç
Party içi loot dağılım kuralları — personal loot default, alternatif rules.

Foundation: [systems/loot-system.md](../systems/loot-system.md), DD-012 personal loot.

## 2. Loot Rules

`social.party.loot_rule` enum:
- `personal` — her karakter kendi roll (default; DD-012).
- `round_robin` — sırayla loot alan üye.
- `need_before_greed` — need/greed prompt (Phase 02).
- `master` — leader dağıtım kararı (raid Phase 02).

## 3. Personal Loot (v1 Default)

Her karakter kill'de loot table roll:
```go
for _, character := range party.Members {
    lootRoll := lootSystem.Roll(mob.LootTableID, character.Seed.Sub(...))
    grantToCharacter(character, lootRoll)
}
```

No competition; no share.
Boss unique drop: dedup enforce (weekly first kill only).

## 4. Round-Robin (Phase 02)

- Sıralı order: A → B → C → A ...
- Her drop → next in order.
- Guaranteed drop rotation.

## 5. Need Before Greed (Phase 02)

- Drop occurred → prompt to eligible members.
- Members respond "Need" / "Greed" / "Pass".
- Need: class + role eligible; roll among all "Need".
- Greed: any; roll among all "Greed" if no Need.
- Timeout 30 sec → auto-pass.

Winner rolls 1-100; highest wins.

## 6. Master Loot (Phase 02+; Raid)

- Leader assigns each drop to member manually.
- Ledger audit for fairness.
- Raid content default.

## 7. Loot Rule Change

**POST `/v1/parties/{id}/loot-rule`**
Leader only.
```json
{ "loot_rule": "need_before_greed" }
```

Server:
- Update party.loot_rule.
- Emit event `party.loot_rule_changed`.

Effect immediate for next drop.

## 8. Boss Loot Priority

Boss loot priority (bkz. [dungeon/loot.md §4](../dungeon/loot.md)):
1. Guaranteed items (weekly first kill).
2. Set item pool.
3. Standard loot table.

Party personal loot: each character independent.

## 9. Party Bonus Loot

Party bonus (v1):
- Rare drop chance +5%.
- Applied per character loot roll.

## 10. Contribution Reward (Phase 02)

Contribution-based bonus:
- Damage dealer (topluluğa >20%) → +5% loot chance.
- Healer (heal done proportional) → same bonus.
- Tank (damage taken proportional) → same bonus.

Fair reward for role play.

## 11. Loot Dispute

Personal loot: no dispute (each independent).
Master/RR/NBG: leader arbitration; report system.

Loot audit (`economy.item_transfer_log`):
- Party context recorded.
- Post-hoc review possible.

## 12. Endpoint

**GET `/v1/parties/{id}/loot-log`** — recent loot log.

## 13. Metrics
- `party_loot_rule_distribution{rule}` gauge.
- `party_loot_dispute_total`.
- `party_loot_bonus_applied_total`.

## 14. Cross-Reference
- [Party](party.md), [Combat](combat.md), [Dungeon](dungeon.md), [Hunt](hunt.md)
- [systems/loot-system.md](../systems/loot-system.md)
- [database/economy-market.md §2](../database/economy-market.md) — item transfer log

## 15. Bekleyen Kararlar
- Need Before Greed Phase 02.
- Master Loot raid Phase 02+.
- Contribution reward formula Phase 02.
- Loot trade window (drop trade to teammate) Phase 02.
