// Package config loads environment-based configuration (12-factor).
// Every field must be either required or have a safe default; secrets never appear in code.
package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config is the process-wide immutable configuration snapshot.
type Config struct {
	Env      string `envconfig:"ENV" default:"dev"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	HTTPAddr         string        `envconfig:"HTTP_ADDR" default:":8080"`
	HTTPReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"15s"`
	HTTPWriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"30s"`

	DatabaseURL      string `envconfig:"DATABASE_URL" required:"true"`
	DatabaseMaxConns int32  `envconfig:"DATABASE_MAX_CONNS" default:"25"`

	RedisAddr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	ContentPath string `envconfig:"CONTENT_PATH" default:"./content"`

	// --- Auth (DD-189) -------------------------------------------------------
	// JWTIssuer identifies this server in the JWT `iss` claim.
	JWTIssuer string `envconfig:"JWT_ISSUER" default:"otter"`

	// JWTPrivateKeySeed is the base64-encoded 32-byte Ed25519 seed. When empty
	// (dev only), the server generates an ephemeral key at boot and logs a
	// warning. Production MUST provide the seed via secret manager.
	JWTPrivateKeySeed string `envconfig:"JWT_PRIVATE_KEY_SEED" default:""`

	AuthAccessTTL      time.Duration `envconfig:"AUTH_ACCESS_TTL" default:"15m"`
	AuthRefreshTTL     time.Duration `envconfig:"AUTH_REFRESH_TTL" default:"720h"` // 30 days
	AuthSessionTTL     time.Duration `envconfig:"AUTH_SESSION_TTL" default:"720h"`
	AuthMinPasswordLen int           `envconfig:"AUTH_MIN_PASSWORD_LEN" default:"8"`

	// --- Character (DD-OPEN-8 slot cap; DD-134 pattern for delete grace) ---
	CharacterSlotCap     int           `envconfig:"CHARACTER_SLOT_CAP" default:"2"`
	CharacterDeleteGrace time.Duration `envconfig:"CHARACTER_DELETE_GRACE" default:"48h"`

	// --- Worker (Sprint 19) ---
	// Freeze characters whose wallet projection drifts from the ledger
	// (docs/economy/transaction.md §6).
	WorkerReconcileFreeze bool `envconfig:"WORKER_RECONCILE_FREEZE" default:"true"`

	// --- Push (Sprint 23). Unset → pushes are only logged. ---
	FCMProjectID       string `envconfig:"FCM_PROJECT_ID" default:""`
	FCMCredentialsFile string `envconfig:"FCM_CREDENTIALS_FILE" default:""`
	APNsKeyFile        string `envconfig:"APNS_KEY_FILE" default:""`
	APNsKeyID          string `envconfig:"APNS_KEY_ID" default:""`
	APNsTeamID         string `envconfig:"APNS_TEAM_ID" default:""`
	APNsTopic          string `envconfig:"APNS_TOPIC" default:""`
	APNsSandbox        bool   `envconfig:"APNS_SANDBOX" default:"true"`

	// Browser origins allowed to call the API (Flutter web). Comma-separated;
	// dev additionally accepts any http://localhost:<port>.
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:""`

	// --- Account recovery (Sprint 37, DD-219) ---
	// Base URL of the links in recovery e-mails (the API host serving
	// /verify-email and /reset-password).
	PublicBaseURL string `envconfig:"PUBLIC_BASE_URL" default:"http://localhost:8080"`
	// When true, characters can only be created after the e-mail is verified
	// (docs/auth/register.md §7). Enable once SMTP delivery is configured.
	AuthRequireVerifiedEmail bool `envconfig:"AUTH_REQUIRE_VERIFIED_EMAIL" default:"false"`
	// SMTP for the email_dispatch worker. Empty host = e-mails are only logged
	// (dev). Port 587 uses STARTTLS; 465 implicit TLS.
	SMTPHost     string `envconfig:"SMTP_HOST" default:""`
	SMTPPort     int    `envconfig:"SMTP_PORT" default:"587"`
	SMTPUsername string `envconfig:"SMTP_USERNAME" default:""`
	SMTPPassword string `envconfig:"SMTP_PASSWORD" default:""`
	SMTPFrom     string `envconfig:"SMTP_FROM" default:""`

	// --- Observability (Sprint 27, DD-214) ---
	// /metrics listens on its own port so it is never exposed through the
	// public reverse proxy. Empty disables it.
	MetricsAddr       string `envconfig:"METRICS_ADDR" default:":9090"`
	WorkerMetricsAddr string `envconfig:"WORKER_METRICS_ADDR" default:":9091"`
	// Tracing is on only when an OTLP/HTTP endpoint is configured (standard
	// OTEL_EXPORTER_OTLP_ENDPOINT, e.g. http://tempo:4318).
	OTLPEndpoint     string  `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:""`
	TraceSampleRatio float64 `envconfig:"TRACE_SAMPLE_RATIO" default:"1.0"`
}

// Load reads env vars into Config; returns friendly error on missing required.
func Load() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return Config{}, fmt.Errorf("config load: %w", err)
	}
	return c, nil
}
