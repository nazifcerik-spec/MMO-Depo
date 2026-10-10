// Command worker runs Otter background jobs (see internal/worker). Safe to
// run as several replicas — each job takes a Postgres advisory lock per run.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // quiet-hours time zones in minimal images

	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/domain/health"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/mail"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/observability"
	"github.com/otter/server/internal/platform/db"
	"github.com/otter/server/internal/platform/log"
	"github.com/otter/server/internal/worker"
)

var CommitSHA = "dev"

func main() {
	// Container HEALTHCHECK (distroless: no wget) — DD-232.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(health.Probe(health.LocalURL(os.Getenv("WORKER_METRICS_ADDR"), ":9091", "/metrics")))
	}
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := log.New(cfg.LogLevel, cfg.Env, CommitSHA)
	slog.SetDefault(logger)
	logger.Info("worker boot", "env", cfg.Env, "commit", CommitSHA)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL, cfg.DatabaseMaxConns)
	if err != nil {
		return err
	}
	defer pool.Close()

	currencies, err := economy.LoadCurrencies(ctx, pool)
	if err != nil {
		return err
	}

	jobs := worker.All(pool, worker.Deps{
		Txn:             economy.NewPgTxnService(pool, currencies, logger),
		Notifier:        notification.NewService(pool, logger),
		Senders:         pushSenders(cfg, logger),
		ReconcileFreeze: cfg.WorkerReconcileFreeze,
		Mailer:          mailSender(cfg, logger),
		Log:             logger,
	})
	for _, j := range jobs {
		logger.Info("worker job registered", "job", j.Name(), "interval", j.Interval().String())
	}
	shutdownTracing, err := observability.InitTracing(ctx, "otter-worker", cfg.Env, CommitSHA, cfg.OTLPEndpoint, cfg.TraceSampleRatio)
	if err != nil {
		return err
	}
	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracing(sctx)
	}()
	metrics := observability.NewMetrics()
	metrics.RegisterPool(pool)
	metrics.Serve(ctx, cfg.WorkerMetricsAddr, logger)

	worker.NewRunner(worker.NewPgLocker(pool), logger, jobs...).WithMetrics(metrics).Start(ctx)
	logger.Info("worker shutdown")
	return nil
}

// pushSenders builds FCM/APNs senders from config, falling back to a log
// sender for any provider without credentials.
func pushSenders(cfg config.Config, logger *slog.Logger) map[string]notification.Sender {
	senders := map[string]notification.Sender{
		"android": notification.LogSender{Platform: "android", Log: logger},
		"ios":     notification.LogSender{Platform: "ios", Log: logger},
	}
	if cfg.FCMCredentialsFile != "" {
		if s, err := notification.NewFCMSender(cfg.FCMProjectID, cfg.FCMCredentialsFile); err != nil {
			logger.Error("fcm sender disabled", "err", err.Error())
		} else {
			senders["android"] = s
			logger.Info("fcm sender enabled", "project", s.ProjectID)
		}
	}
	if cfg.APNsKeyFile != "" {
		if s, err := notification.NewAPNsSender(cfg.APNsKeyFile, cfg.APNsKeyID, cfg.APNsTeamID, cfg.APNsTopic, cfg.APNsSandbox); err != nil {
			logger.Error("apns sender disabled", "err", err.Error())
		} else {
			senders["ios"] = s
			logger.Info("apns sender enabled", "sandbox", cfg.APNsSandbox)
		}
	}
	return senders
}

// mailSender picks SMTP when configured; otherwise mail is only logged in dev
// and refused (kept in the outbox, retried) elsewhere.
func mailSender(cfg config.Config, logger *slog.Logger) mail.Sender {
	if cfg.SMTPHost == "" {
		logger.Warn("SMTP_HOST not set — e-mails are only logged (dev) or held in the outbox")
		return mail.LogSender{Dev: cfg.Env == "dev", Log: logger}
	}
	if cfg.SMTPFrom == "" {
		logger.Error("SMTP_FROM is required when SMTP_HOST is set — e-mail disabled")
		return mail.LogSender{Log: logger}
	}
	logger.Info("smtp mail enabled", "host", cfg.SMTPHost, "port", cfg.SMTPPort)
	return mail.SMTPSender{
		Host: cfg.SMTPHost, Port: cfg.SMTPPort,
		Username: cfg.SMTPUsername, Password: cfg.SMTPPassword, From: cfg.SMTPFrom,
	}
}
