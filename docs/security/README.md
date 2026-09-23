# Security + Anti-Cheat — Katalog

> **Faz 33**: Sunucu-otoriter oyunun güvenlik ve hile önleme katmanı. Şimdiye kadar tanımlanan tüm sistemlerin **arka duvarı**.

## 1. Kapsam
Bu klasör, oyunda **kimin ne yapabileceğini kimlerin doğrulayacağını** ve **hile senaryolarına karşı savunma** stratejilerini tarif eder. Kod satırı değil; **prensipler + kural setleri + doğrulama matrisleri**.

## 2. Threat Model (Yüzey)

Ana tehdit vektörleri:

| Tehdit | Örnek | Karşı Önlem |
|---|---|---|
| **Client tampering** | Modified client — sahte damage, sahte item | Server-authoritative + input validation |
| **Replay attack** | Aynı isteği tekrar oynat — dupe/gold | Idempotency key + nonce + timestamp |
| **Race condition** | Concurrent buy/craft ile dupe | SERIALIZABLE TX + optimistic lock |
| **Bot / automation** | Otomatik farm / grinding | Rate limit + behavior anomaly |
| **Credential stuffing** | Şifre reuse ile hesap ele geçirme | 2FA + rate limit + password strength |
| **Session hijack** | Token çalma → hesap kullanma | HTTPS + rotating refresh + device fingerprint |
| **Trade fraud** | Bait-and-switch, offer değişimi | Auto-unlock + confirm hash (DD-138) |
| **RMT** | Gerçek para karşılığı gold satış | Anomaly detection + ban bundle |
| **Exploit** | Bug'ı sömürüp dupe/infinite | Reconciliation + reversal + patch |
| **DoS / DDoS** | Sunucuyu boğma | Cloudflare/AWS shield + rate limit + circuit breaker |
| **Data leak** | PII sızıntı | Encryption + audit + minimization (DD-184) |

## 3. Ana Prensipler (DD-186)

- **Server hakikat** — hiçbir client input tam güven görmez; her değer server tarafından yeniden hesap veya validate edilir.
- **Ledger-first** — tüm değer transferi audit'lenir (DD-135); post-hoc drift tespit + reversal mümkün.
- **Least privilege** — endpoint permission bazlı; JWT scope + admin RBAC.
- **Idempotent by default** — retry safe (`idempotency_key`).
- **Fail closed** — belirsiz durumda deny (unknown role → viewer, invalid token → 401).
- **Rate limit her yerde** — endpoint bazlı, açık default.
- **Defense in depth** — WAF + rate limit + auth + input validate + business rule + reconcile.
- **Zero-trust internal** — mikro servis içi çağrı bile mTLS + service identity (Phase 02).
- **Observability = güvenlik** — her anomaly metric ve alert.

## 4. Dokümanlar
- [rate-limit.md](rate-limit.md) — Endpoint bazlı rate limit + burst token bucket.
- [request-validation.md](request-validation.md) — Input schema + sanitization + size limits.
- [transaction-validation.md](transaction-validation.md) — Ledger integrity + double-entry + retry idempotent.
- [inventory-validation.md](inventory-validation.md) — Item ownership + lock + bind + slot cap.
- [currency-validation.md](currency-validation.md) — Wallet check + cap + negatif block.
- [trade-validation.md](trade-validation.md) — P2P fraud + auto-unlock + hash.
- [replay-protection.md](replay-protection.md) — Idempotency key + nonce + WS message dedup.
- [session-security.md](session-security.md) — Token lifecycle + rotation + device fingerprint.
- [audit.md](audit.md) — Security audit + anomaly detection + response.

## 5. Cross-Reference
- [auth](../auth/README.md), [admin/auth-rbac.md](../admin/auth-rbac.md)
- [backend](../backend/README.md), [economy/transaction.md](../economy/transaction.md)
- [social/block.md](../social/block.md), [marketplace/trade.md](../marketplace/trade.md)
- [database/audit.md](../database/audit.md)

## 6. Metrikler (Security KPI)
- `sec_rate_limit_hit_total{endpoint, actor_kind}`
- `sec_input_validation_reject_total{endpoint, reason}`
- `sec_idempotency_replay_total`
- `sec_auth_fail_total{reason}`
- `sec_session_revoke_total{reason}`
- `sec_anomaly_flag_total{kind}`
- `sec_reconcile_drift_total{severity}`

## 7. Bekleyen Kararlar
- Client integrity attestation (Play Integrity API / DeviceCheck) Phase 02.
- WAF sağlayıcı (Cloudflare vs AWS WAF) infra karar.
- Anomaly detection ML pipeline Phase 03.
- Bug bounty program Phase 02.
- Yasal compliance (GDPR/KVKK/COPPA) legal review Phase 02.
