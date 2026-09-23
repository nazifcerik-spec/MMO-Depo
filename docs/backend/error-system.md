# Error System

## 1. Amaç
Domain hatalarını sentinel + wrap ile ifade etmek; HTTP/WS/worker katmanlarında tutarlı yanıt üretmek; log/metric/trace'e uygun context taşımak.

## 2. Sentinel Set (Kategori)

```go
package errors

import "errors"

// Kategori sentinel'ları — domain paketleri kendi altında specialize edebilir.
var (
    ErrValidation    = errors.New("validation")
    ErrNotFound      = errors.New("not_found")
    ErrConflict      = errors.New("conflict")
    ErrForbidden     = errors.New("forbidden")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrRateLimited   = errors.New("rate_limited")
    ErrUnavailable   = errors.New("unavailable")
    ErrInternal      = errors.New("internal")
    ErrTimeout       = errors.New("timeout")
    ErrIdempotent    = errors.New("idempotent_replay")  // OK / no-op
)
```

Domain-specific:
```go
package character

var (
    ErrNameTaken      = fmt.Errorf("%w: name taken", errs.ErrConflict)
    ErrInvalidClass   = fmt.Errorf("%w: class not found", errs.ErrValidation)
    ErrSlotExhausted  = fmt.Errorf("%w: character slot cap reached", errs.ErrForbidden)
)
```

## 3. Error Wrap Kuralı
```go
if err := db.Query(...); err != nil {
    return fmt.Errorf("character.load: %w", err)
}
```
- `fmt.Errorf("%w", ...)` ile cause zinciri korunur.
- Sentinel yakalama: `errors.Is(err, ErrNotFound)`.
- Custom detail: özel struct + `Unwrap()`.

## 4. AppError Struct (opsiyonel enrichment)
```go
type AppError struct {
    Code       string            // "character.name_taken"
    Sentinel   error             // errs.ErrConflict
    Message    string            // user-safe (i18n key veya default)
    Details    map[string]any    // opsiyonel
    Cause      error             // wrap chain
}

func (e *AppError) Error() string { return e.Code + ": " + e.Message }
func (e *AppError) Unwrap() error { return e.Cause }
func (e *AppError) Is(target error) bool { return errors.Is(e.Sentinel, target) }
```

Builder:
```go
func NameTaken(name string) *AppError {
    return &AppError{
        Code:     "character.name_taken",
        Sentinel: errs.ErrConflict,
        Message:  "This name is already in use.",
        Details:  map[string]any{"name": name},
    }
}
```

## 5. HTTP Mapping
```go
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
    code, status, msg, details := classify(err)
    log := slog.FromContext(r.Context())

    switch {
    case status >= 500:
        log.ErrorContext(r.Context(), "http error", "code", code, "err", err)
    case status == 429:
        log.WarnContext(r.Context(), "rate limited", "code", code)
    default:
        log.DebugContext(r.Context(), "handled error", "code", code)
    }

    metrics.HTTPErrorTotal.WithLabelValues(code, strconv.Itoa(status)).Inc()

    resp := ErrorResponse{
        Error: ErrorBody{
            Code:      code,
            Message:   msg,
            RequestID: ctxutil.RequestID(r.Context()),
            Details:   details,
        },
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(resp)
}
```

Sentinel → HTTP status mapping:
| Sentinel | Status | Notes |
|---|---|---|
| `ErrValidation` | 400 | Detail with field errors |
| `ErrUnauthorized` | 401 | `WWW-Authenticate` header |
| `ErrForbidden` | 403 | |
| `ErrNotFound` | 404 | |
| `ErrConflict` | 409 | |
| `ErrRateLimited` | 429 | `Retry-After` |
| `ErrTimeout` | 504 | |
| `ErrUnavailable` | 503 | Ready check yansıması |
| `ErrIdempotent` | 200/204 | Actually success; safe replay |
| `ErrInternal` | 500 | Generic message; no leak |
| Diğer (unknown) | 500 | Generic |

## 6. Response Body
```json
{
  "error": {
    "code": "character.name_taken",
    "message": "This name is already in use.",
    "request_id": "01J8T2P...",
    "details": { "name": "Rhaegar" }
  }
}
```

Kurallar:
- `code` client-parseable, versionlu (breaking change'de yeni code).
- `message` user-facing (i18n resolver client tarafında).
- `request_id` support için.
- `details` opsiyonel, structured.

## 7. Field Validation Errors
```json
{
  "error": {
    "code": "character.validation",
    "message": "Some fields are invalid.",
    "details": {
      "fields": {
        "name": "must be 3-16 chars",
        "race_id": "unknown race"
      }
    }
  }
}
```

## 8. WebSocket Mapping
- Server → client mesaj:
```json
{ "type": "error", "code": "activity.contract_active", "message": "...", "request_id": "..." }
```
- Persistent error (auth invalid) → close frame kod + reason.

## 9. Worker / Job Mapping
- Job error → `last_error` alanına yazılır; `attempts++`.
- `max_attempts` aşılınca `state='dead'`.
- Sentinel `ErrIdempotent` → job success (already applied).

## 10. Security / Info Leak
- Internal error message asla external'a taşınmaz.
- 500 message default: `"Something went wrong."` (i18n).
- Stack trace yalnız log.

## 11. Sentry (Client + Server)
- Server: `ErrInternal` ve `panic` Sentry'ye gönderilir; `request_id` tag.
- Client: crash + boundary errors.

## 12. Test Stratejisi
- Table-driven: her sentinel → expected status/body.
- Wrap chain test: `errors.Is(fmt.Errorf("wrap: %w", ...), ErrNotFound)` doğrulanır.

## 13. Anti-Patterns
- `panic()` iş akışında.
- Error yutmak (`_ = err`).
- Error mesajını user-visible olarak internal detail'la doldurmak.
- Wrap zincirini `err.Error()` string parsing ile kırmak.

## 14. İlişkiler
- [Middleware](middleware.md), [Logger](logger.md), [WebSocket](websocket.md), [Worker](worker.md), [Router](router.md)

## 15. Bekleyen Kararlar
- Error code registry: enum-like sabitler yerine dokümante edilmiş code listesi (Phase 01 CI check).
- i18n error message resolver (server tarafında mı client mı) — öneri client.
- Detail schema (fields object structure) versionlama.
