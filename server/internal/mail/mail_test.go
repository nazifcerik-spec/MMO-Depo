package mail

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessage_HeadersAndEncoding(t *testing.T) {
	m := string(message("noreply@otter.gg", "p@example.com", "Otter — Şifre sıfırlama", "line1\nline2"))
	require.Contains(t, m, "From: noreply@otter.gg\r\n")
	require.Contains(t, m, "To: p@example.com\r\n")
	require.Contains(t, m, "Subject: =?utf-8?q?", "non-ASCII subject is Q-encoded")
	require.Contains(t, m, "Content-Type: text/plain; charset=utf-8\r\n")
	require.Contains(t, m, "@otter.gg>\r\n", "Message-ID uses the sender domain")
	require.True(t, strings.HasSuffix(m, "\r\n\r\nline1\r\nline2"), "body uses CRLF")
}

func TestLogSender_RefusesOutsideDev(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	require.NoError(t, LogSender{Dev: true, Log: log}.Send(context.Background(), "a@b.c", "s", "b"))
	require.Error(t, LogSender{Dev: false, Log: log}.Send(context.Background(), "a@b.c", "s", "b"),
		"reset links must never end up in production logs")
}

func TestSMTPSender_RejectsHeaderInjection(t *testing.T) {
	err := SMTPSender{Host: "localhost", Port: 587}.Send(context.Background(), "a@b.c\r\nBcc: x@y.z", "s", "b")
	require.ErrorContains(t, err, "invalid recipient")
}
