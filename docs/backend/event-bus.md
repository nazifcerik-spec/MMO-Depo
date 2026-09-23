# Event Bus

## 1. Amaç
Domain modülleri arasında **loosely-coupled** iletişim: bir modül olay yayınlar; başka modüller (aynı process içinde veya farklı) dinler. v1 in-process + Redis pub/sub combo; v2 dedicated broker (NATS/Redis Streams).

## 2. İki Katmanlı Model

### 2.1 In-process bus (synchronous / async)
- Aynı Go process içinde.
- Fan-out: interested subscribers directly.
- Kullanım: cross-module hooks (level up → achievement check).

### 2.2 Distributed bus (Redis pub/sub v1)
- Cross-node broadcasts.
- Kullanım: WebSocket fanout (chat, notification), world event announcement.

Şart: **domain state persist eden event Postgres outbox'a önce yazılır** (bkz. [transactions.md §Outbox](transactions.md)). Bus fire-and-forget değil; kalıcı olay outbox'tan gelir.

## 3. Sözleşme
```go
package bus

type Event struct {
    ID          string    // ulid
    Kind        string    // "character.leveled_up"
    OccurredAt  time.Time
    ActorID     uuid.UUID // character/account
    SubjectKind string
    SubjectID   uuid.UUID
    PayloadJSON []byte
    CorrelationID uuid.UUID
    TraceID     string
}

type Handler func(ctx context.Context, e Event) error

type Bus interface {
    Publish(ctx context.Context, e Event) error
    Subscribe(kind string, h Handler) (unsubscribe func())
    Close() error
}
```

Implementasyonlar:
- `InProcessBus` — memory queue + goroutine dispatcher.
- `RedisPubSubBus` — Redis publish + subscriber loop.
- `CompositeBus` — publish in-process ilk, sonra Redis (for cross-node).

## 4. Event Kind Registry

Kanonik isim: `<domain>.<past_tense_verb>`.

```
character.created
character.leveled_up
character.ascended
character.stat_allocated
character.died

activity.contract_started
activity.contract_paused
activity.contract_completed
activity.reveal_computed

economy.currency_credited
economy.currency_debited
economy.market_listing_created
economy.market_listing_cancelled
economy.market_transaction_completed
economy.mail_sent

item.instance_created
item.equipped
item.unequipped
item.salvaged
item.durability_lost

profession.leveled_up
profession.recipe_learned
crafting.craft_completed

combat.mob_killed
combat.boss_killed
combat.death_recorded

quest.accepted
quest.completed
quest.turned_in

achievement.completed

guild.created
guild.member_joined
guild.member_left
guild.bank_deposited

party.created
party.member_joined

world_boss.spawned
world_boss.defeated

event.started
event.ended

notification.dispatched
```

Registry doküman `docs/backend/event-registry.md` (Phase 01 canlı liste).

## 5. Publish Semantics

### 5.1 In-process
- Sync (immediate handler call) — sadece **critical fast handler** (auth check tetikleyici).
- Async (queue) — default; publisher block etmez.
- Handler error → log; publish rollback etmez (fire-and-forget domain event).

### 5.2 Distributed
- Redis `PUBLISH kind payload` (JSON envelope).
- No persistence — subscriber yoksa mesaj kaybolur.
- Kalıcı gereksinim: outbox + worker delivery.

## 6. Subscribe Örnek
```go
// Domain init'inde
bus.Subscribe("character.leveled_up", func(ctx context.Context, e bus.Event) error {
    var p CharacterLeveledUp
    if err := json.Unmarshal(e.PayloadJSON, &p); err != nil { return err }
    // achievement check
    return achievements.MaybeGrantLevelMilestone(ctx, p.CharacterID, p.NewLevel)
})
```

## 7. Outbox Discipline

Kalıcı side-effect'i olan event:
1. Domain tx içinde `audit.event_log` (veya `outbox`) insert.
2. Commit sonrası worker outbox'ı okur → `bus.Publish` çağırır → Redis/webhook fan-out.
3. Delivery idempotent; retry safe.

Fire-and-forget event (analytics only): direct publish OK.

## 8. Cross-Node Delivery Modeli (v1)
```
Node A: publish → Redis PUBLISH "evt:character.leveled_up"
Node B: subscriber → handler
Node C: subscriber → handler
```
- No exactly-once; at-most-once (Redis pub/sub).
- Missed subscribers (offline node): kaybolur; kritik değil for real-time push.

## 9. v2 Alternatifleri
- **Redis Streams**: at-least-once, consumer groups, replay.
- **NATS JetStream**: rich semantics, subject wildcard.
- Karar tetikleyicisi: outbox worker throughput yetersiz kalırsa, ya da guaranteed delivery gerekirse (analytics pipeline).

## 10. Payload Şeması
- JSON, versionlu (`payload._v` alanı).
- Küçük tut (< 4KB); büyük detail için id ver.
- Sensitive alan yok (PII, secret).

Örnek:
```json
{
  "_v": 1,
  "character_id": "c-777",
  "old_level": 24,
  "new_level": 25,
  "stat_points_awarded": 5
}
```

## 11. Testing
- Unit: `InProcessBus` sync mode + assertion.
- Integration: `RedisPubSubBus` testcontainers Redis + subscribe.
- Contract test: schema registry (Phase 02) — event JSON versionlama korunur.

## 12. Anti-Patterns
- Event içinde büyük payload / full entity dump.
- Handler'da başka event publish (loop riski) — cycle detector.
- Publish transaction içinde (Redis I/O tx'i kirletir) — outbox kullan.
- Consumer'da uzun süreli iş (block bus dispatcher) — worker'a devret.

## 13. Metrics
- `bus_published_total{kind}` counter.
- `bus_handled_total{kind,result}` counter.
- `bus_handler_duration_seconds{kind}` histogram.
- `bus_dispatch_lag_seconds` histogram (in-process queue depth).

## 14. Logging
- Handler log: `{msg="event handled", kind, event_id, handler, dur_ms}`.
- Handler fail: warn + retry policy (v1 no retry inside bus; outbox worker retry).

## 15. Güvenlik
- Cross-node payload trust: Redis network internal; TLS opsiyonel.
- Public webhook (Phase 02): HMAC signed payload.

## 16. İlişkiler
- [Database](database.md), [Transactions](transactions.md), [Worker](worker.md), [Redis](redis.md), [WebSocket](websocket.md), [Notification System](../systems/notification-system.md), [database/audit-system.md](../database/audit-system.md)

## 17. Bekleyen Kararlar
- v2 broker seçimi (Redis Streams vs NATS).
- Schema registry tool (Phase 02).
- Handler retry semantics (in-bus retry vs outbox-only).
- Event trace parent propagation (OTEL).
