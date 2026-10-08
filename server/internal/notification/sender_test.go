package notification

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(p, data, 0o600))
	return p
}

func TestFCMSender(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	var oauthCalls atomic.Int32
	var lastMsg map[string]any
	reply := http.StatusOK
	replyBody := `{}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/token":
			oauthCalls.Add(1)
			require.NoError(t, r.ParseForm())
			require.NotEmpty(t, r.Form.Get("assertion"))
			_, _ = io.WriteString(w, `{"access_token":"at-1","expires_in":3600}`)
		case strings.HasSuffix(r.URL.Path, "/messages:send"):
			require.Equal(t, "Bearer at-1", r.Header.Get("Authorization"))
			require.Equal(t, "/v1/projects/proj-x/messages:send", r.URL.Path)
			lastMsg = map[string]any{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&lastMsg))
			w.WriteHeader(reply)
			_, _ = io.WriteString(w, replyBody)
		}
	}))
	defer srv.Close()

	sa, _ := json.Marshal(map[string]string{"client_email": "svc@x", "private_key": string(keyPEM),
		"token_uri": srv.URL + "/token", "project_id": "proj-x"})
	f, err := NewFCMSender("", writeFile(t, "sa.json", sa))
	require.NoError(t, err)
	f.Endpoint = srv.URL
	ctx := context.Background()

	push := Push{Token: "dev-1", Priority: "high", Payload: Payload{Kind: "market_sold", TitleKey: "t", BodyKey: "b",
		Category: "economy", CollapseKey: "c1", Vars: map[string]any{"item": "sword"}}}
	require.NoError(t, f.Send(ctx, push))
	android := lastMsg["message"].(map[string]any)["android"].(map[string]any)
	require.Equal(t, "high", android["priority"])
	require.Equal(t, "c1", android["collapse_key"])
	require.NotNil(t, android["notification"])

	push.Priority = "silent"
	require.NoError(t, f.Send(ctx, push))
	require.Nil(t, lastMsg["message"].(map[string]any)["android"].(map[string]any)["notification"])
	require.Equal(t, int32(1), oauthCalls.Load(), "access token is cached")

	reply, replyBody = http.StatusNotFound, `{"error":{"status":"NOT_FOUND"}}`
	require.ErrorIs(t, f.Send(ctx, push), ErrInvalidToken)
	reply, replyBody = http.StatusBadRequest, `invalid registration token`
	require.ErrorIs(t, f.Send(ctx, push), ErrInvalidToken)
	reply, replyBody = http.StatusServiceUnavailable, `busy`
	err = f.Send(ctx, push)
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrInvalidToken)

	// OAuth failure surfaces; bad credentials fail at construction.
	f.token = ""
	f.sa.TokenURI = srv.URL + "/nope"
	require.Error(t, f.Send(ctx, push))
	_, err = NewFCMSender("p", writeFile(t, "bad.json", []byte(`{"private_key":"junk"}`)))
	require.Error(t, err)
	_, err = NewFCMSender("p", writeFile(t, "notjson.json", []byte(`{`)))
	require.Error(t, err)
	_, err = NewFCMSender("p", filepath.Join(t.TempDir(), "missing.json"))
	require.Error(t, err)
}

func TestAPNsSender(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	der, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	keyFile := writeFile(t, "key.p8", pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))

	var hdr http.Header
	var body map[string]any
	reply, replyBody := http.StatusOK, ""
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/3/device/tok-1", r.URL.Path)
		hdr = r.Header.Clone()
		body = map[string]any{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.WriteHeader(reply)
		_, _ = io.WriteString(w, replyBody)
	}))
	defer srv.Close()

	a, err := NewAPNsSender(keyFile, "KID", "TEAM", "com.otter.app", true)
	require.NoError(t, err)
	require.Contains(t, a.Endpoint, "sandbox")
	a.Endpoint = srv.URL
	ctx := context.Background()

	push := Push{Token: "tok-1", Priority: "high", Payload: Payload{Kind: "mail_received", TitleKey: "t", BodyKey: "b",
		Category: "social", CollapseKey: "m1"}}
	require.NoError(t, a.Send(ctx, push))
	require.Equal(t, "alert", hdr.Get("apns-push-type"))
	require.Equal(t, "m1", hdr.Get("apns-collapse-id"))
	require.Equal(t, "com.otter.app", hdr.Get("apns-topic"))
	require.True(t, strings.HasPrefix(hdr.Get("authorization"), "bearer "))
	first := hdr.Get("authorization")

	push.Priority = "silent"
	require.NoError(t, a.Send(ctx, push))
	require.Equal(t, "background", hdr.Get("apns-push-type"))
	require.Equal(t, float64(1), body["aps"].(map[string]any)["content-available"])
	require.Equal(t, first, hdr.Get("authorization"), "provider token is reused")

	reply, replyBody = http.StatusGone, `{"reason":"Unregistered"}`
	require.ErrorIs(t, a.Send(ctx, push), ErrInvalidToken)
	reply, replyBody = http.StatusBadRequest, `{"reason":"BadDeviceToken"}`
	require.ErrorIs(t, a.Send(ctx, push), ErrInvalidToken)
	reply, replyBody = http.StatusTooManyRequests, `{"reason":"TooManyRequests"}`
	require.NotErrorIs(t, a.Send(ctx, push), ErrInvalidToken)

	prod, err := NewAPNsSender(keyFile, "KID", "TEAM", "t", false)
	require.NoError(t, err)
	require.Equal(t, "https://api.push.apple.com", prod.Endpoint)

	// Key file problems.
	_, err = NewAPNsSender(writeFile(t, "plain.p8", []byte("nope")), "K", "T", "t", false)
	require.Error(t, err)
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	rsaDER, _ := x509.MarshalPKCS8PrivateKey(rsaKey)
	_, err = NewAPNsSender(writeFile(t, "rsa.p8", pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: rsaDER})), "K", "T", "t", false)
	require.Error(t, err)
	_, err = NewAPNsSender(writeFile(t, "junk.p8", pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("x")})), "K", "T", "t", false)
	require.Error(t, err)
	_, err = NewAPNsSender(filepath.Join(t.TempDir(), "missing.p8"), "K", "T", "t", false)
	require.Error(t, err)
}

func TestLogSender(t *testing.T) {
	require.NoError(t, LogSender{Platform: "android", Log: slog.New(slog.NewTextHandler(io.Discard, nil))}.
		Send(context.Background(), Push{Payload: Payload{Kind: "x"}}))
}
