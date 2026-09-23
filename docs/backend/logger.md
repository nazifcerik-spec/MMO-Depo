# Logger

## 1. Amaç
Tek slog-based logger; structured JSON; correlation-aware; tüm platform'a inject edilir.

Ana kural set: bkz. [architecture/logging.md](../architecture/logging.md). Bu doküman Go implementasyonuna odaklanır.

## 2. Sözleşme
```go
package logging

func New(cfg Config) *slog.Logger

type Config struct {
    Level      slog.Level  // parsed from env "info"/"debug"
    JSON       bool        // production=true, dev=false (pretty)
    Sampling   SamplingCfg // optional
    Service    string      // "server" | "worker"
    Env        string      // "dev" | "staging" | "prod"
    Version    string      // build sha
}
```

## 3. Handler Zinciri
```
BaseHandler (JSON) → SamplingHandler → ContextEnricher → slog.Logger
```

### 3.1 ContextEnricher
Her log record'una `ctx` içindeki standart alanları ekler:
- `request_id`
- `trace_id`, `span_id` (OTEL)
- `account_id`
- `character_id`
- `activity_id` (varsa)

```go
type ctxHandler struct{ inner slog.Handler }

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
    if v := ctxutil.RequestID(ctx); v != "" { r.AddAttrs(slog.String("request_id", v)) }
    if v := ctxutil.AccountID(ctx); v != uuid.Nil { r.AddAttrs(slog.String("account_id", v.String())) }
    if v := ctxutil.CharacterID(ctx); v != uuid.Nil { r.AddAttrs(slog.String("character_id", v.String())) }
    if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
        r.AddAttrs(
            slog.String("trace_id", sc.TraceID().String()),
            slog.String("span_id", sc.SpanID().String()),
        )
    }
    return h.inner.Handle(ctx, r)
}
```

### 3.2 SamplingHandler (opsiyonel)
- Debug seviyesi yüksek volume — sampling rate (1%, config).
- Info+ tam.
- Sampling key: `request_id` hash → deterministic sample.

## 4. Base Fields (her satır)
```
ts, level, msg, service, env, version, request_id, trace_id, span_id
```

Domain alanları event bazında (`character_id`, `activity_id`, `order_id`, ...).

## 5. Kullanım Kuralları
- `slog.InfoContext(ctx, "activity started", "activity_type", "mining", "spot_id", spot)`.
- `msg` sabit; parametreler alan olarak.
- **PII yasak** (email, IP raw). Hash veya kısaltılmış.
- Error log: `slog.ErrorContext(ctx, "reveal failed", "err", err)` — slog `err` alanını serialize eder.
- Attribute value büyük olmasın: > 1KB → truncate.

## 6. Level Yönlendirme (Env)
| Env | Default Level | JSON |
|---|---|---|
| dev | debug | false (pretty) |
| ci | info | true |
| staging | info | true |
| prod | info | true |

Runtime override yok; deploy config ile.

## 7. Domain Log Örnekleri
```go
log.InfoContext(ctx, "character created",
    "character_id", char.ID, "race_id", char.RaceID, "class_id", char.ClassID)

log.InfoContext(ctx, "market listing created",
    "listing_id", l.ID, "item_template_id", l.ItemTemplateID,
    "unit_price_copper", l.UnitPriceCopper, "qty", l.Qty)

log.WarnContext(ctx, "rate limited",
    "route", "market.list", "subject_kind", "account")

log.ErrorContext(ctx, "reveal failed",
    "contract_id", c.ID, "err", err)
```

## 8. Package Logger
- App genel logger `wire`/`fx`/DI ile inject.
- Alt paket kendi child logger'ını `log.With("component","httpx")` ile üretir.

## 9. Test Stratejisi
- `slog.NewTextHandler(io.Discard, ...)` fake logger.
- Structured assertion için `slogtest` (stdlib).

## 10. Performance Notları
- slog attributes lazy: `slog.Any` yerine typed helper (`slog.String`, `slog.Int`).
- Hot path (combat tick): sample veya level=warn'a düşür; per-tick log YASAK bulk çıkarımla.

## 11. Compliance
- **PII redaction**: helper `logsafe.Email(s)` → hashed 8 char prefix.
- **Secrets**: middleware whitelisted header set; Authorization asla log.

## 12. İlişkiler
- [Middleware](middleware.md), [Error System](error-system.md), [architecture/logging.md](../architecture/logging.md), [WebSocket](websocket.md), [Worker](worker.md)

## 13. Bekleyen Kararlar
- Logger source line inclusion (slog `AddSource`) — dev'de yes, prod'da hayır (cost).
- Sampling rate for debug logs — Phase 01 measurement sonrası.
- Trace-linked log storage (Loki + Tempo integration).
