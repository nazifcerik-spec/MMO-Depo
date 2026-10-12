package worker

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/mail"
	"github.com/otter/server/internal/notification"
)

func TestSendersFromConfig(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// No credentials: pushes are logged, never sent.
	s := pushSenders(config.Config{}, logger)
	require.IsType(t, notification.LogSender{}, s["android"])
	require.IsType(t, notification.LogSender{}, s["ios"])
	// Broken credentials fall back to logging instead of failing the worker.
	s = pushSenders(config.Config{FCMCredentialsFile: "/nope.json", APNsKeyFile: "/nope.p8"}, logger)
	require.IsType(t, notification.LogSender{}, s["android"])
	require.IsType(t, notification.LogSender{}, s["ios"])

	require.Equal(t, mail.LogSender{Dev: true, Log: logger}, mailSender(config.Config{Env: "dev"}, logger))
	require.Equal(t, mail.LogSender{Log: logger}, mailSender(config.Config{Env: "prod", SMTPHost: "smtp.example.com"}, logger),
		"SMTP without a From address stays disabled")
	require.IsType(t, mail.SMTPSender{}, mailSender(config.Config{SMTPHost: "smtp.example.com", SMTPFrom: "game@example.com", SMTPPort: 587}, logger))
}
