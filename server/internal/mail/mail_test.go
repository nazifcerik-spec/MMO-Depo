package mail

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

func TestBrevoSender(t *testing.T) {
	var got map[string]any
	var key string
	status := http.StatusCreated
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key = r.Header.Get("api-key")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"message":"sender not verified"}`))
	}))
	defer srv.Close()
	s := BrevoSender{APIKey: "k-1", From: "noreply@otter.gg", FromName: "New World MMO", Endpoint: srv.URL}

	require.NoError(t, s.Send(context.Background(), "p@example.com", "Şifre sıfırlama", "link"))
	require.Equal(t, "k-1", key)
	require.Equal(t, map[string]any{"email": "noreply@otter.gg", "name": "New World MMO"}, got["sender"])
	require.Equal(t, []any{map[string]any{"email": "p@example.com"}}, got["to"])
	require.Equal(t, "Şifre sıfırlama", got["subject"])
	require.Equal(t, "link", got["textContent"])

	status = http.StatusBadRequest
	err := s.Send(context.Background(), "p@example.com", "s", "b")
	require.ErrorContains(t, err, "status 400")
	require.ErrorContains(t, err, "sender not verified")
	require.ErrorContains(t, s.Send(context.Background(), "a@b.c\r\nBcc: x@y.z", "s", "b"), "invalid recipient")
}
