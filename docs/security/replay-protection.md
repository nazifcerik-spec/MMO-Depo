# Replay Protection

## 1. Amaç
Aynı isteğin (kasıtlı veya yanlışlıkla) ikinci kez işlenerek dupe/gold/state bozukluğuna yol açmaması.

## 2. Katmanlar

- **Idempotency key** — client-provided; server ilk çağrıda sonucu cache eder; ikinci çağrı aynı response.
- **Nonce** — bazı kritik akışlarda (ör. auth challenge) tek kullanımlık; server tarafta consume.
- **Timestamp + window** — istek timestamp; ±5 dk dışında reject (opsiyonel, replay bot koruma).
- **Version / sequence number** — websocket message id monotonic; missed events resume.

## 3. Idempotency Key

**Format**: `alphanumeric 32 karakter`, UUID v4 önerilir.

**Kapsamı zorunlu endpoint'ler** (mutating + reveal):
- Activity start / reveal / stop
- Vendor buy / sell
- Market listing / buy / cancel
- Trade offer update / confirm
- Mail send / claim
- Bank deposit / withdraw
- Guild bank ops
- Craft enqueue
- Prestige node choose
- Any endpoint invoking `TxnService.Write`

**Model** (per endpoint namespace):
```sql
CREATE TABLE api.idempotency_record (
  key VARCHAR(64) NOT NULL,
  endpoint VARCHAR(96) NOT NULL,
  account_id UUID NOT NULL,
  request_hash BYTEA NOT NULL,
  response_body JSONB NOT NULL,
  response_status SMALLINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '24 hours',
  PRIMARY KEY (key, endpoint, account_id)
);
CREATE INDEX idx_idem_expire ON api.idempotency_record (expires_at);
```

Retention 24 saat hot (retry window); worker cleanup > 24h.

## 4. Lookup Flow

```
POST /v1/vendor/buy { idempotency_key: "abc..." }
  ↓
Middleware:
  1. Lookup (key, endpoint, account_id) → hit?
     - HIT + request_hash match: return cached response.
     - HIT + request_hash mismatch: 409 E_IDEMPOTENT_MISMATCH (aynı key farklı body suspicious).
     - MISS: proceed to handler.
  2. Handler success → insert idempotency_record.
  3. Response.
```

**Request hash**: SHA-256(body + endpoint + account_id) — client'ın aynı istek olduğunu doğrular.

## 5. Ledger Entry Idempotency

`economy.ledger_entry` içinde `idempotency_key` UNIQUE (idempotency_key, entry_seq) constraint (bkz. [economy/transaction.md](../economy/transaction.md) §Model). Retry TX aynı ledger entry'yi ikinci kez yazamaz.

## 6. Websocket Message Dedup

- Client → server WS message'lara `client_msg_id` (monotonic per session).
- Server tarafta son 100 message id Redis'te; duplicate reject.
- Server → client push'ta `server_msg_id`; client resume `resume_from:<last_id>`.

## 7. Auth Challenge Nonce

- Login challenge (2FA) `nonce` server generate + Redis TTL 5 dk.
- Verify → nonce consume (delete).
- Reuse attempt → `E_NONCE_USED`.

## 8. Timestamp Window (Opsiyonel v1)

Bazı yüksek riskli endpoint'ler (`/internal/ops/*`, admin step-up) request timestamp header:
```
X-Request-Timestamp: 2026-09-09T14:32:15Z
```
Server tarafta ±5 dk pencere dışı reject.

Public API'da timestamp header **zorunlu değil** (mobil client clock drift); idempotency + rate limit yeterli.

## 9. Replay Attack Response

Suspicious pattern (idempotency mismatch + peş peşe attempt):
- `sec_replay_attempt_total` metric.
- 5+ mismatch / 5dk → session revoke + inbox notification "unusual activity".
- Repeat → account temporary lock + support review.

## 10. Metrikler
- `sec_idempotency_replay_total{endpoint}`
- `sec_idempotency_mismatch_total{endpoint}` (suspicious)
- `sec_ws_msg_dedup_total`
- `sec_nonce_reuse_attempt_total`
- `sec_timestamp_window_reject_total`

## 11. Cross-Reference
- [economy/transaction.md](../economy/transaction.md)
- [backend/middleware.md](../backend/middleware.md)
- [auth/session.md](../auth/session.md)

## 12. Bekleyen Kararlar
- Idempotency retention (24h v1; Phase 02 endpoint bazında tune).
- Timestamp window zorunlu public API Phase 02.
- Client SDK auto-idempotency key generate Phase 02.
