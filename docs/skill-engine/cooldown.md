# Cooldown Manager

## 1. Amaç
Skill cooldown state'inin low-latency, multi-node safe yönetimi. Global GCD, per-skill CD, category shared CD, internal CD (passive proc), CDR (cooldown reduction) uygulaması.

## 2. Sözleşme
```go
package cooldown

type Manager interface {
    RemainingMs(ctx context.Context, charID uuid.UUID, key string) (int32, error)
    Set(ctx context.Context, charID uuid.UUID, key string, ms int32) error
    Clear(ctx context.Context, charID uuid.UUID, key string) error
    BatchRemaining(ctx context.Context, charID uuid.UUID, keys []string) (map[string]int32, error)
    // ICD (internal cooldown) — passive proc için
    ICDReady(ctx context.Context, charID uuid.UUID, passiveID uuid.UUID) (bool, error)
    ICDSet(ctx context.Context, charID uuid.UUID, passiveID uuid.UUID, ms int32) error
}
```

## 3. Redis Key Convention

Bkz. [backend/redis.md](../backend/redis.md).

```
GCD:                  char:gcd:<char_id>              TTL = gcd_ms
Skill CD:             char:cd:<char_id>:<skill_code>  TTL = cd_ms
Category CD:          char:catcd:<char_id>:<key>      TTL = cat_cd_ms
Internal CD (passive): icd:<char_id>:<passive_id>     TTL = icd_ms
```

Value: not important; TTL yeterli (existence check → active).

## 4. Set / Check Örnekleri

```go
// Set skill cooldown
_ = c.rdb.Set(ctx, "char:cd:"+charID.String()+":berserker_rampage", "1", 60*time.Second).Err()

// Remaining
ttl, _ := c.rdb.PTTL(ctx, "char:cd:...:berserker_rampage").Result()
if ttl <= 0 { /* ready */ }

// Batch — Lua script
```

Batch remaining: MGET yerine PTTL Lua script (single round-trip):
```lua
-- KEYS[*] are cooldown keys
local out = {}
for i, k in ipairs(KEYS) do
  out[i] = redis.call("PTTL", k)
end
return out
```

## 5. Cooldown Reduction (CDR)

CDR uygulanma anı: **Set sırasında** (yani Redis'e giden `ms` değeri CDR sonrası).

```go
finalMs := int32(float64(skill.CooldownMs) * (1 - min(cdr_pct/100, 0.30)))
_ = c.Set(ctx, charID, skill.Code, finalMs)
```

CDR cap %30 (bkz. [stats.md](stats.md)).

Alternatif (Phase 02): CDR runtime sırasında check zamanında hesap. Trade-off: buff değişikliği aktif CD'yi değiştirir mi? v1 tercih: **not; Set anındaki değer sabit** (fair, çoğu MMO'da böyle).

## 6. GCD (Global Cooldown)

- Baz: 1000 ms.
- Haste ile azalır:
  ```
  gcd_ms = 1000 * (1 - min(haste_pct/100, 0.5))
  min gcd = 500 ms
  ```
- Skill cast sonrası set (unless `gcd_bypass`).
- Cast → CanCast, sonra GCD set.

Redis key TTL kısa (< 1 sn); Redis atomicity guaranteed.

## 7. Cooldown Cache Invalidation

Cooldown Redis-only; DB'de tutulmaz (ephemeral).

Redis kayıp durumu:
- Tüm cooldown expire (as-if fresh) → oyuncu için avantaj (rare).
- Alternatif: DB'de "last_cast_at per skill" (Phase 02 opsiyonel).

## 8. Cooldown Reset (Utility)

- Bazı skill'ler diğer skill'lerin CD'sini resetler ("finisher resets primary CD").
- Effect DSL:
  ```json
  { "kind": "cooldown_reset", "target": "self", "skill_codes": ["berserker_rampage"] }
  ```
- Implementation: `Clear(ctx, charID, skill_code)`.

- CD refund (%): kısmi
  ```json
  { "kind": "cooldown_refund", "target": "self", "skill_code": "x", "refund_pct": 50 }
  ```
- Implementation: get remaining, set = remaining × (1 - pct).

## 9. Category Cooldown

Şared grup (potion vs trinket):
```yaml
category_cooldown_key: potion_group
category_cooldown_ms: 60000
```

Cast sırasında:
- Check `char:catcd:<char>:potion_group` remaining.
- Set after cast: `SET char:catcd:<char>:potion_group ... EX 60`.

## 10. ICD (Internal Cooldown for Passives)

Passive proc rate limiting:
```
icd:<char>:<passive_id>
```
TTL = passive.icd_ms.

`ICDReady` check → `EXISTS icd:...`.

## 11. Concurrency

Redis atomicity yeterli (single-key operations).

Nadir race:
- Client eş zamanlı 2 cast → biri success, diğeri PTTL check'te CD gördü → reject.

## 12. Fail Modları

- Redis unreachable: fail-open? fail-closed?
  - Fail-open: cooldown skip → spam mümkün, kısa süreli.
  - **Fail-closed** tercih: cast reject `skill.cooldown_service_unavailable`.
- TTL negative dönerse: normalize 0.

## 13. Test Stratejisi

- Unit: CDR calculation.
- Integration: real Redis (testcontainers); set → PTTL → cast reject.
- Race: goroutine 100 concurrent cast → sadece 1 success.
- Reset/refund logic.

## 14. Metrics
- `cooldown_set_total{skill}`.
- `cooldown_reject_total{skill}`.
- `cooldown_reset_total{skill}`.
- `cooldown_service_unavailable_total`.

## 15. Persistence (Opsiyonel)

Bazı cooldown'ları uzun süreli tutmak (dungeon lockout, daily quest CD):
- Postgres tablosu (bkz. [database/progression-quests-achievements.md §8](../database/progression-quests-achievements.md)).
- Bu cooldown manager sadece **ephemeral combat/skill CD** için.

## 16. Cross-References
- [Active Skills](active-skills.md), [Passives](passives.md), [Stats](stats.md)
- [backend/redis.md](../backend/redis.md), [systems/afk-engine.md](../systems/afk-engine.md)
- [game-design/skills.md §5](../game-design/skills.md)

## 17. Bekleyen Kararlar
- CD reset skill'lerin balance sınırı.
- Persistent CD (dungeon lockout) engine sözleşmesine dahil edilsin mi (öneri hayır, ayrı domain).
- Cross-node CD sync (Redis pub/sub gerekli mi) — hayır, key başına atomic yeter.
- CDR runtime hesap alternatifi (Phase 02).
