# Reputation (Extended)

> [character/reputation.md](../character/reputation.md) v1 temel; bu doküman genişletir.

## 1. Amaç
Karakterin belirli bir faction ile ilişki seviyesi. Vendor access, quest unlock, faction cosmetic, faction gate.

## 2. Karakteristik

- **Scope**: karakter (v1); account paylaşımı yok (RP korunur).
- **Faction başına ayrı rep**.
- **Tier sistem**: 8 tier (v1).
- **Gain/loss**: quest, kill, delivery, event.
- **Decay yok v1** (opsiyonel Phase 02 slow decay).
- **Neutral start**: yeni karakter tüm faction ile tier `neutral`.

## 3. Tier Sistemi

| Tier | Adı | XP Threshold (kümülatif) |
|---|---|---|
| 0 | Hated | -6000 |
| 1 | Hostile | -3000 |
| 2 | Unfriendly | -1000 |
| 3 | Neutral | 0 (start) |
| 4 | Friendly | 3000 |
| 5 | Honored | 9000 |
| 6 | Revered | 21000 |
| 7 | Exalted | 42000 |

Negatif rep'te hostile: faction NPC saldıracak; faction vendor kapanacak; faction quest'ler kilitli.

## 4. Model

```sql
CREATE TABLE progression.character_reputation (
  character_id UUID NOT NULL,
  faction_id UUID NOT NULL,
  rep_xp INT NOT NULL DEFAULT 0,           -- signed, negatif olabilir
  current_tier SMALLINT NOT NULL DEFAULT 3, -- neutral start
  last_gain_at TIMESTAMPTZ NULL,
  daily_gain_used INT NOT NULL DEFAULT 0,   -- daily cap tracker
  daily_reset_at TIMESTAMPTZ NULL,
  version INT NOT NULL DEFAULT 1,
  PRIMARY KEY (character_id, faction_id)
);
```

## 5. Gain / Loss Kaynakları

| Source | Delta | Notlar |
|---|---|---|
| Faction quest completion | +25 / +100 / +500 (rarity göre) | Ana kaynak |
| Faction kill (opposite hostile) | +5 | Elite +20 |
| Faction delivery | +50 | Repeatable, daily cap altında |
| World event (faction) | +50 - +500 | Contribution tier ile |
| Enemy faction quest complete | -100 | Opposite rep loss |
| Faction diplomat kill | -500 | Explicit hostility |

**Daily cap**: rep_xp gain / gün / faction = **1500** (grinding koruma). Cap dolarsa fazla gain 0. Reset 00:00 UTC (DD-016 uyumlu).

## 6. Etki

- **Vendor discount**: [economy/vendor.md](../economy/vendor.md) — `reputation_discount = clamp(1 - repTier × 0.02, 0.85, 1.0)`.
- **Fee discount**: [economy/fee.md](../economy/fee.md) — Phase 02 opt.
- **Faction quest unlock**: tier threshold (friendly+ ana quest, revered+ epic).
- **Faction cosmetic vendor**: tier gated item (mount, tabard, title).
- **Hostility gate**: negative tier → NPC saldırganlık, city ban.

## 7. Tier Change Akışı

Rep gain event → worker:
1. `character_reputation.rep_xp += delta` (with sign).
2. Yeni tier hesap (tier_thresholds lookup).
3. Tier değişti mi?
   - Yükseldi → WS "reputation_tier_up" + achievement trigger + potansiyel title.
   - Düştü → WS "reputation_tier_down" (uyarı).
4. Faction cosmetic/quest unlock/lock kontrol.
5. Modifier stack invalidate (vendor discount cache).

## 8. Endpoint

**GET `/v1/character/{id}/reputation`** — tüm faction rep.
**GET `/v1/character/{id}/reputation/{faction_code}`** — detay + next tier XP.

## 9. Metrikler
- `reputation_gain_total{faction, source}`
- `reputation_loss_total{faction, source}`
- `reputation_daily_cap_hit_total`
- `reputation_tier_change_total{faction, direction}`

## 10. Cross-Reference
- [character/reputation.md](../character/reputation.md) — v1 temel
- [faction.md](faction.md)
- [economy/vendor.md](../economy/vendor.md), [economy/fee.md](../economy/fee.md)
- [quests/README.md](../quests/README.md)

## 11. Bekleyen Kararlar
- Rep decay Phase 02 opt (uzun süre inaktif → rep -1 tier / 90 gün).
- Account-wide rep paylaşımı Phase 02.
- Faction switch (rep transfer) Phase 03 mekaniği.
- Daily cap tune Phase 02 balance.
