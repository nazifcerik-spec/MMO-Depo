# Environments (Ortamlar)

## 1. Amaç
Kod ve içeriğin production'a giden yolda farklı ortamlarda güvenli olarak test edilmesi. Her ortamın rolü, veri politikası ve erişimi net.

## 2. Ortam Listesi

| Ortam | Amaç | Data Kaynağı | Erişim | Güncel Kod |
|---|---|---|---|---|
| **local** | Geliştirici makinesi | Seed script | Dev bireysel | branch |
| **ci** | Test pipeline | Testcontainers | CI runner | PR HEAD |
| **dev** | Sürekli entegrasyon deploy | Fresh seed, günlük reset | Takım | `main` |
| **staging** | Prod-like doğrulama | Anonymized prod snapshot | Takım + QA | release candidate tag |
| **prod** | Canlı oyuncu | Real | Sınırlı (SRE) | released tag |
| **preview** (Phase 02) | PR başına ephemeral | Sub-seed | Reviewer | PR |

## 3. Ortamlar Arasındaki Parity
- Aynı Docker image dev → staging → prod'a promote.
- Aynı schema; migration önce staging'de doğrulanır.
- Sadece config farklı (env variable + secret).
- **Prod-only kod path yok**; feature flag ile şartlı davranış.

## 4. Data Politikaları
- **local**: fake data, seed script; herhangi bir gerçek PII yok.
- **dev**: seed + günlük reset (haftalık full wipe).
- **staging**: anonymized prod snapshot; PII masked (email → hash, gold ölçek korunur).
- **prod**: gerçek veri; backup + PITR (bkz. [PostgreSQL](postgresql.md)).

## 5. Erişim Kontrolü
- Prod DB: sadece SRE/on-call; audit log her sorguya.
- Staging: takım geneli read; write kısıtlı.
- Dev: takım geneli tam erişim.
- Local: bireysel.

## 6. Env Variable Kataloğu (özet)
| Değişken | Zorunlu | Örnek |
|---|---|---|
| `APP_ENV` | evet | `dev`, `staging`, `prod` |
| `APP_VERSION` | evet | git sha veya semver |
| `LOG_LEVEL` | hayır (default `info`) | `debug` |
| `DB_URL` | evet | `postgres://user:pass@host:5432/mmo?sslmode=require` |
| `DB_MAX_CONN` | hayır | `50` |
| `REDIS_URL` | evet | `redis://:pass@host:6379/0` |
| `HTTP_PORT` | hayır | `8080` |
| `WS_PORT` | hayır | `8081` |
| `JWT_SIGNING_KEY_ID` | evet | key id (secret'a lookup) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | staging+ | `https://otel...` |
| `SENTRY_DSN` | staging+ | dsn |
| `RATE_LIMIT_DEFAULT` | hayır | `100/60s` |
| `FEATURE_FLAGS_MODE` | hayır | `env` / `service` |

Detaylı validasyon: bkz. [Config](config.md).

## 7. Secret Management
- Provider-native secret manager (Vault, AWS Secrets Manager, GCP Secret Manager). Env'e enjekte edilir.
- Local: `.env.local` (git-ignored); dev seed script placeholder.
- Rotate: quarterly; auto-warn `last_used > 90d`.
- JWT signing key: rotation destekli (key id ile version'lı).

## 8. Domain / URL Konvansiyonu
```
dev.mmo.example       (API)
dev-play.mmo.example  (client)
staging.mmo.example
staging-play.mmo.example
api.mmo.example       (prod)
play.mmo.example
```
TLS: her ortamda; cert managed (Let's Encrypt veya provider).

## 9. Feature Flag Ortam Politikası
- **dev**: yeni feature default ON.
- **staging**: opt-in per test; gerekiyorsa ON.
- **prod**: default OFF; gradual rollout (Phase 02 flag service).

## 10. Ortam-Spesifik Ayarlar
- Rate limit prod'da daha sıkı.
- Log level: dev `debug`, staging `info`, prod `info` (isteğe bağlı `warn`).
- Notification: prod'da e-mail aktif; dev'de dry-run (log only).
- Push provider: dev/staging test project; prod ayrı.
- Payment provider (Phase 02): sandbox dev/staging; live prod.

## 11. Reset & Refresh
- **dev**: haftalık full wipe; günlük seed script re-run.
- **staging**: manual refresh from prod snapshot (ayda 1 planlı).
- **prod**: reset yok.

## 12. Health / Readiness
- Health endpoint her ortamda aynı.
- Monitoring: prod öncelik; staging ikincil; dev best-effort.

## 13. Ortamlar Arası Data Taşıma
- prod → staging: managed snapshot restore + anonymization job.
- staging → prod: **yok**; asla staging veri'si prod'a itilmez.
- prod → local: yok; local her zaman synthetic.

## 14. Anti-Patterns
- Prod DB'yi dev'e mount etmek.
- Aynı secret'ı dev + prod'da paylaşmak.
- Ortam URL'sini hardcode etmek.
- `.env` dosyasını commit etmek.
- Prod'da manual DB değişikliği (migration bypass).

## 15. İlişkiler
- [Config](config.md), [CI/CD](cicd.md), [Docker](docker.md), [PostgreSQL](postgresql.md), [Redis](redis.md), [Logging](logging.md)

## 16. Bekleyen Kararlar
- Cloud provider (AWS / GCP / hybrid).
- Managed vs self-hosted DB/cache.
- Preview env pipeline (Phase 02).
- Anonymization script öncelik ve tool seçimi.
