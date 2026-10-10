# Load Test: Login

## 1. Amaç
Login endpoint'in burst yükü altında SLO'yu koruması. Server restart / maintenance return / patch day peak senaryoları kritik.

## 2. Senaryolar

### 2.1 Steady Login
- Yük: 5-10 login/sn (background rate CCU T3 için).
- Süre: 60 dk.
- Beklenti: p95 < 300ms, error < 0.5%.

### 2.2 Login Burst (Maintenance Return)
- Yük: 5 dk içinde 500-1000 login/sn (T3-T4 scale).
- Süre: 5 dk burst + 15 dk normalize.
- Beklenti: p95 < 800ms (relaxed), error < 2%, rate limit engage.

### 2.3 Login Storm (Patch Day)
- Yük: 1 dk içinde 2000+ login/sn.
- Süre: 60 sn.
- Beklenti: rate limit koruma engage; error rate elevated but recovery immediate; DB pool sağlam.

### 2.4 Login + OTP (Yeni Cihaz)
- Yük: %10 login yeni device (OTP challenge zorunlu).
- Beklenti: OTP flow < 500ms p95 additional.

## 3. Kritik Path

Login flow:
1. HTTP receive.
2. Rate limit check (Redis).
3. Body parse + validate.
4. Password verify (Argon2id — CPU intensive).
5. 2FA check.
6. Session insert (DB).
7. Refresh token generate + insert.
8. JWT sign (EdDSA).
9. Response.

**Bottleneck adayları:**
- **Argon2id** CPU (memory 64MB × concurrent count → server RAM basıncı).
- **DB session insert** — write throughput.
- **Redis rate limit** — network RTT.

## 4. Test Metrics

- `login_p50/p95/p99_ms`.
- `login_error_rate{reason}` (invalid_credentials, rate_limit, lock, otp_required).
- `argon2_hash_duration_ms`.
- `argon2_concurrent_count`.
- `db_session_insert_latency_ms`.

## 5. k6 Script Kalıbı (Kavramsal)

```js
scenario 'login_burst':
  executor: ramping-arrival-rate
  startRate: 10
  stages:
    - target: 500, duration: '2m'
    - target: 500, duration: '5m'
    - target: 10,  duration: '5m'
  script:
    creds = randomCredential()
    r = http.post('/v1/auth/login', {email, password})
    check(r.status == 200)
```

## 6. Argon2id Tuning

Password verify CPU-bound; concurrent limit:
- **memory 64 MB × 500 concurrent = 32 GB** → dev server RAM aşabilir.
- Çözüm: **thread pool** (max 100 concurrent Argon2 verify); queue overflow → 503 briefly veya rate limit friendly reject.
- Parametreleri Phase 01 benchmark ile ayarla (memory/iterations trade-off).
- **Uygulama (DD-238):** slot sınırı `AUTH_HASH_CONCURRENCY` (varsayılan CPU sayısı), 500 ms bekleme → 503 `E_BUSY` + `Retry-After`. Parametreler m=19 MiB, t=2 (hash ~27 ms, 64 MiB×3'te ~140 ms idi). Ölçüm: `scripts/loadtest.sh login` (`perf/k6/login.js`).

## 7. DB Pool Sağlığı

- Max pool 50 v1 (küçük dev); T3 için 200; T4 için 500.
- Login DB write per → 2 insert (session, refresh_token).
- Peak login 1000/sn × 2 = 2000 write/sn. DB write cap kontrol.

## 8. Rate Limit Engage

Bkz. [security/rate-limit.md](../security/rate-limit.md):
- `/v1/auth/login`: 10/dk / IP.
- 5 fail sonrası +15dk lock.
- Storm test: rate limit hit ratio ölçü; false positive kontrol.

## 9. Fail Mode Handling

- 5xx altında client backoff (exp) beklenir.
- Cascade: login fail → refresh fail → user retry loop; monitoring gerekli.

## 10. Beklenen Sonuç Örnek

T3 login burst 500/sn:
- p50: ~120ms.
- p95: ~350ms.
- p99: ~800ms.
- Error rate < 1%.
- Argon2 pool utilization ~70%.
- DB session insert p95 < 30ms.

## 11. Cross-Reference
- [slo-slis.md](slo-slis.md), [ccu-scenarios.md](ccu-scenarios.md)
- [auth/login.md](../auth/login.md), [security/session-security.md](../security/session-security.md)

## 12. Bekleyen Kararlar
- Argon2 CPU / concurrency final ayar Phase 01 benchmark.
- Session insert async (fire-and-forget) opt Phase 02.
- Login queue buffer (503 vs slow response) Phase 02.
