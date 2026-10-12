package worker

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/mail"
	"github.com/otter/server/internal/notification"
)

// FromConfig builds the full job set from config — shared by cmd/worker and
// the API's RUN_WORKER mode (one process on hosts without a free background
// worker, DD-239).
func FromConfig(ctx context.Context, cfg config.Config, pool *pgxpool.Pool, logger *slog.Logger) ([]Job, error) {
	currencies, err := economy.LoadCurrencies(ctx, pool)
	if err != nil {
		return nil, err
	}
	jobs := All(pool, Deps{
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
	return jobs, nil
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
