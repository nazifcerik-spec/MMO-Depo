# Reputation

## 1. Amaç
Karakterin fraksiyonlarla ilişkisini ölçen sayısal değer. Vendor discount, unique item access, storyline gate. **v1'de sınırlı** (Phase 02'de tam faction sistemi); ancak schema + baseline design burada.

## 2. Neden Var
- Karakteri dünyaya bağlar (yaptığın işler kime yarar).
- Uzun soluklu ödül eğrisi (yeni content ile fresh grind loop).
- Kimlik: "Emerald Wardens'da Exalted" prestij.

## 3. Model
Kaynak: [database/character.md §15](../database/character.md), [database/progression-quests-achievements.md §9](../database/progression-quests-achievements.md).

Tablolar:
- `content.faction` — faction tanımı (kod, ad, tier eşikleri).
- `character.reputation` — karakter × faction × value + tier.

Şema (v1 baseline):
```sql
CREATE TABLE content.faction (
  id UUID PRIMARY KEY,
  code VARCHAR(32) UNIQUE NOT NULL,
  display_name_key VARCHAR(64) NOT NULL,
  description_key VARCHAR(96) NULL,
  tier_thresholds JSONB NOT NULL DEFAULT
    '{"hated":-6000,"hostile":-3000,"unfriendly":-1000,"neutral":0,"friendly":3000,"honored":9000,"revered":21000,"exalted":42000}',
  parent_faction_id UUID NULL,  -- alt/üst faction ilişkileri
  is_at_war_with_ids UUID[] NOT NULL DEFAULT '{}',
  is_allied_with_ids UUID[] NOT NULL DEFAULT '{}',
  version INT NOT NULL DEFAULT 1
);

CREATE TABLE character.reputation (
  character_id UUID NOT NULL,
  faction_id UUID NOT NULL,
  value INT NOT NULL DEFAULT 0,
  tier VARCHAR(16) NOT NULL DEFAULT 'neutral',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (character_id, faction_id)
);
```

## 4. Tiers
| Kod | Aralık | Sunum |
|---|---|---|
| `hated` | ≤ -6000 | Kırmızı |
| `hostile` | -6000..-3001 | Turuncu |
| `unfriendly` | -3000..-1001 | Sarı |
| `neutral` | -1000..2999 | Gri (start) |
| `friendly` | 3000..8999 | Yeşil |
| `honored` | 9000..20999 | Mavi |
| `revered` | 21000..41999 | Mor |
| `exalted` | 42000+ | Altın |

Cap: 42000 (exalted); üstü artmaz.

## 5. Kaynaklar (Reputation Grant)
| Kaynak | Delta |
|---|---|
| Quest turn-in (faction) | +100..+500 |
| Repeatable daily | +25..+100 |
| Faction dungeon kill | +50 per boss |
| Faction adventure | +30..+150 |
| Championship item turn-in | +250 per item |
| Reputation potion (buff) | +10% for 30 min |

Anti-faction: bir faction'a +X ilişkili opposing faction'a -X.

## 6. Faction Interactions (v2 Detay)
- **Ally chain**: A'ya rep +100 → allied B'ye +25 (soft chain).
- **War chain**: A'ya +100 → at war C'ye -50.
- **Neutral chain**: no effect.

v1'de: sadece direct grant (chain YOK); design için schema hazır.

## 7. Tier Etkileri

### 7.1 Vendor Discount
```
friendly:  5% discount
honored:   10%
revered:   15%
exalted:   20%
```

### 7.2 Item / Recipe Unlock
- Bazı vendor item'lar `min_reputation_tier = 'honored'`.
- Bazı recipe'ler `min_reputation_tier = 'revered'`.

### 7.3 Quest Unlock
- Faction quest chain'ler tier gate.
- Exalted quest chain: unique mount/pet/title reward.

### 7.4 Guild Reputation (Phase 02)
- Guild-based faction: guild üyeleri kolektif rep artırır.

## 8. Endpoint

### 8.1 Karakter Reputation Listesi

**GET `/v1/characters/{id}/reputation`**
Response:
```json
{
  "factions": [
    {
      "faction_id": "f-emerald-wardens",
      "code": "emerald_wardens",
      "display_name_key": "faction.emerald_wardens",
      "value": 15420,
      "tier": "honored",
      "next_tier": { "code": "revered", "at": 21000, "delta_needed": 5580 }
    }
  ]
}
```

### 8.2 Reputation Grant (Internal)
Domain call:
```go
func (s *Service) GrantReputation(ctx context.Context, charID, factionID uuid.UUID, delta int32, source string) error
```
- Transaction; opposing faction penalty computation.
- Tier transition detection → event 'character.reputation_tier_changed'.
- Audit event.

## 9. Grant Path
```
Quest turn-in:
  reward includes faction_id + delta
  → GrantReputation(char, faction, delta, "quest:xxx")
  → +delta with buff multipliers
  → tier compute
  → opposing faction penalty (Phase 02)
  → audit + event
```

## 10. Config
```yaml
reputation:
  cap: 42000
  buff_potion_pct: 10
  chain_ally_ratio: 0.25   # v2
  chain_war_ratio: -0.5    # v2
```

## 11. Event
```json
{
  "_v": 1,
  "character_id": "c-777",
  "faction_id": "f-...",
  "old_tier": "friendly",
  "new_tier": "honored",
  "value": 9200
}
```

Subscribers:
- `achievement.reputation_milestone`.
- `notification.tier_up`.

## 12. Fail Modları
- Faction unknown → 404.
- Cap overflow → clamp; log.
- Concurrent grant → optimistic (opsiyonel; low collision).

## 13. Anti-Cheat
- Repeatable quest reputation'a per-hour cap (bot prevention).

## 14. Cross-Platform
Standart character sheet; her platformda tutarlı.

## 15. Metrics
- `character_reputation_granted_total{faction,source}`.
- `character_reputation_tier_changed_total{faction,tier}`.

## 16. Audit
- `character.reputation_granted` (severity=debug).
- `character.reputation_tier_changed` (info).

## 17. Test Stratejisi
- Tier transition boundary test.
- Cap clamp.
- Grant idempotency (quest turn-in double-fire koruma).

## 18. İlişkiler
- [Titles](titles.md), [Level](level.md), [game-design/quests.md](../game-design/quests.md), [game-design/economy.md](../game-design/economy.md), [database/character.md](../database/character.md)

## 19. Bekleyen Kararlar
- v1'de aktif faction seti (1-2 örnek faction ile başla).
- Chain ally/war Phase 02'de gerçek matematik.
- Guild reputation (Phase 02+).
- Paragon (post-exalted grinding) Phase 03.
