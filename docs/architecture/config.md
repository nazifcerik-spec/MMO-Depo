# Configuration Management

## 1. Amaç
Uygulama davranışını build sonrası değiştirebilmek; ortamlar arasında güvenli, doğrulanabilir, tekrar edilebilir konfigürasyon.

## 2. Prensipler
- **12-factor**: config yalnızca environment üzerinden.
- **Immutable image**: aynı binary tüm ortamlarda; config runtime injection.
- **Strong typing**: config bir Go struct'a decode edilir + validate.
- **Fail-fast**: eksik veya geçersiz config → boot'ta hata, koşuştur değil.
- **Single source of truth per key**: aynı ayar iki yerde tanımlı olmaz.

## 3. Config Katmanları

Sıra (soldan sağa, override):
```
default (kod) → embedded config file → env variables → secret manager
```

- **Default (kod)**: struct alanları için sensible default.
- **Embedded config file** (`config/base.yaml`): geliştirici okunabilirliği; sadece non-secret + non-env-varying keys.
- **Env variable**: ortam-spesifik; primary override.
- **Secret manager**: hassas key'ler (JWT sig, DB pass, API key).

## 4. Config Struct Örneği (Go)
```go
type Config struct {
    App struct {
        Env     string `env:"APP_ENV,required"`
        Version string `env:"APP_VERSION,required"`
    }
    HTTP struct {
        Port         int           `env:"HTTP_PORT" envDefault:"8080"`
        ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
        WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"15s"`
    }
    DB struct {
        URL        string `env:"DB_URL,required"`
        MaxConn    int    `env:"DB_MAX_CONN" envDefault:"50"`
        StatementTimeout time.Duration `env:"DB_STATEMENT_TIMEOUT" envDefault:"3s"`
    }
    Redis struct {
        URL     string `env:"REDIS_URL,required"`
        PoolMax int    `env:"REDIS_POOL_MAX" envDefault:"100"`
    }
    Auth struct {
        JWTSignKeyID string `env:"JWT_SIGN_KEY_ID,required"`
        AccessTTL    time.Duration `env:"AUTH_ACCESS_TTL" envDefault:"15m"`
        RefreshTTL   time.Duration `env:"AUTH_REFRESH_TTL" envDefault:"720h"`
    }
    Observability struct {
        LogLevel       string `env:"LOG_LEVEL" envDefault:"info"`
        OTELEndpoint   string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
        SentryDSN      string `env:"SENTRY_DSN"`
        MetricsPort    int    `env:"METRICS_PORT" envDefault:"9090"`
    }
    Features struct {
        Mode string `env:"FEATURE_FLAGS_MODE" envDefault:"env"`
    }
    RateLimit struct {
        DefaultPerMinute int `env:"RATE_LIMIT_DEFAULT_PM" envDefault:"600"`
    }
}
```

## 5. Load & Validate
1. Struct decode (env parser).
2. Cross-field validation:
   - `App.Env in {dev, staging, prod}`.
   - Prod requires `OTELEndpoint`, `SentryDSN`.
   - `Auth.AccessTTL < Auth.RefreshTTL`.
3. Boot fail on error → exit 1 + structured log.

## 6. Config File Kullanımı (isteğe bağlı)
- Non-secret, non-varying defaults için `config/base.yaml` (embed).
- Örnek: rate limit endpoint mapping tablosu, feature flag names.
- Ortam-farklılığı olan her şey env variable'a taşınır.

## 7. Content Config (Data-Driven)
Kod-config'inden ayrı: **content config**.
- Location: `/content/*.yaml` veya DB `content` schema.
- Race, class, item template, recipe, skill vs.
- Version'lı; hot-reloadable.
- Loader: `internal/content` paketi; validation her load'da.
- Content patch süreci: staging validate → prod apply.

## 8. Feature Flags
- v1 mode `env`: `FF_<NAME>=on|off`.
- Phase 02: dedicated service (Unleash, OpenFeature, veya kendi minimal impl).
- Flag kullanımı:
  ```go
  if features.Enabled(ctx, "ascension_v2") {
      // ...
  }
  ```
- Flag lifespan: 30 gün max; expire → cleanup PR.

## 9. Secret Kaynakları
- **Prod**: cloud secret manager (Vault, AWS SM, GCP SM). Uygulama boot'ta çeker.
- **Staging**: aynı stratejik.
- **Dev / CI**: env + `.env` (dev pratikliği); repo'ya girmez.
- **Local**: `.env.local` (git-ignored) + optional `direnv`.

## 10. Reload / Live Config
- Kod config: reload yok; deploy ile yenilenir.
- Content config: hot-reload (running server yeni versionu okur).
- Feature flag: live reload (Phase 02 service pattern).

## 11. Audit
- Config değişimi: git log yeterli.
- Feature flag değişimi (service): audit log tut.
- Secret rotation: manager log; erişim iz.

## 12. Anti-Patterns
- Business logic ortam adına bakar (`if env == "prod" { ... }`) — feature flag kullan.
- Env var default'unun secret olması.
- Aynı ayarın hem env hem yaml'da tanımlı olması.
- Config validation olmadan boot.
- Runtime'da env var mutating.

## 13. Test Stratejisi
- Config decode + validate: table-driven test.
- Boot with missing required var → error path test.

## 14. İlişkiler
- [Environment](environment.md), [CI/CD](cicd.md), [Docker](docker.md), [Go Server](go.md), [Logging](logging.md)

## 15. Bekleyen Kararlar
- Env parser: `caarlos0/env` (öneri) vs `spf13/viper` vs custom.
- Secret manager: cloud-native seçim + Vault gerekliliği.
- Feature flag service: v2 için tool seçimi.
- Content storage: DB'de mi (JSONB) file-based mi (YAML in git) — öneri hibrit (git → build-time bundle + DB overlay).
