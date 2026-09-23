# Request Validation

## 1. Amaç
Server'a gelen her istek payload'un şema, boyut, tip, semantik açıdan doğrulanması. Malformed/malicious input hemen reddedilir.

## 2. Katmanlar

Sıralı middleware:
1. **HTTP boyut limiti** — body max 256 KB (avatar/asset upload hariç).
2. **Content-Type check** — `application/json` (aksi 415).
3. **JSON parse** — malformed → 400.
4. **Schema validation** — struct + tag'ler (Go `validator` benzeri) veya JSON Schema.
5. **Sanitization** — string trim + normalize (NFC unicode).
6. **Semantic validation** — business rule (level range, enum, referential integrity).

## 3. Field Kuralları (Örnek)

```go
type StartActivityReq struct {
  Kind           string  `validate:"required,oneof=mining fishing hunting herbalism woodcutting farming skinning archaeology combat travel dungeon craft"`
  Params         Params  `validate:"required"`
  DurationSeconds int    `validate:"required,gte=60,lte=43200"`  // 1dk-12s
  IdempotencyKey string  `validate:"required,len=32,alphanum"`
}
```

Reject → `400 { code: 'E_VALIDATION', field: 'kind', message: '...' }`.

## 4. String Kuralları

- Max length her string field'da (256 default; long body 4096; chat 512).
- **UTF-8** enforce.
- **Unicode normalize NFC** (Türkçe karakter tutarlı).
- **Control char** strip (`\x00-\x1F` hariç `\n \t \r`).
- **HTML/JS injection**: chat/mail body plaintext — HTML sanitize (bleach benzeri).
- **Zero-width chars** (U+200B, U+200C, U+FEFF) strip (impersonation koruma).

## 5. Numeric Kuralları

- Integer overflow guard (Go int64 cap check).
- Currency amount **BIGINT copper** (DD-136); negatif reddedilir (delta üretimi ledger service içindedir).
- Coordinate / percentage range enforce (0..100 vs).

## 6. UUID / Referential

- UUID format (v4/v7) validation.
- Referential: `character_id` istek yapan user'a mı ait? Ait değilse 403.
- Item instance owner check (bkz. [inventory-validation.md](inventory-validation.md)).

## 7. Enum

Whitelist enum liste; unknown → `E_INVALID_ENUM`. Enum listesi content/schema versioned.

## 8. Time Fields

- Timestamp UTC (RFC3339).
- Clock skew tolerance ±5 dk (client saati güvenilmez).
- Future timestamp `E_TIME_FUTURE` (bazı endpoint).

## 9. Sanitization Örnek (Chat)

```
Original: "  Hello​ world <script>alert(1)</script>  "
After:    "Hello world <script>alert(1)</script>"    # trim + zero-width strip
After sanitize (plaintext channel): text kept but escaped on render
```

**Not:** Chat plaintext; HTML escape client-side render'da. Server'da tag STRIP değil ESCAPE.

## 10. File Upload (Phase 02)

Avatar / cosmetic upload:
- Content-Type whitelist (image/png, image/jpeg, image/webp).
- Magic byte check (header signature).
- Boyut < 2 MB.
- Dimension cap (2048×2048).
- Malware scan (ClamAV Phase 02).

## 11. Error Response Formatı

```json
{
  "error": {
    "code": "E_VALIDATION",
    "message": "duration_seconds must be between 60 and 43200",
    "field": "duration_seconds",
    "correlation_id": "..."
  }
}
```

Client localization key: `error.<code>`; message default English fallback.

## 12. Metrikler
- `sec_input_validation_reject_total{endpoint, reason, field}`
- `sec_input_size_reject_total{endpoint}`
- `sec_input_utf8_normalize_total`

## 13. Cross-Reference
- [backend/error-system.md](../backend/error-system.md), [backend/middleware.md](../backend/middleware.md)
- [social/chat.md](../social/chat.md) — chat filter ek katman

## 14. Bekleyen Kararlar
- JSON Schema vs struct tag (Go validator) — Phase 01 backend team.
- File upload MVP scope Phase 02.
- Client SDK auto-validation (server schema paylaşımı) Phase 02.
