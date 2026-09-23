package notification

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInQuietHours(t *testing.T) {
	at := func(h, m int) time.Time { return time.Date(2026, 5, 1, h, m, 0, 0, time.UTC) }
	// Wrapping window 22:30–08:00 UTC.
	require.True(t, InQuietHours(at(23, 0), "22:30", "08:00", "UTC"))
	require.True(t, InQuietHours(at(7, 59), "22:30", "08:00", "UTC"))
	require.False(t, InQuietHours(at(8, 0), "22:30", "08:00", "UTC"))
	require.False(t, InQuietHours(at(12, 0), "22:30", "08:00", "UTC"))
	// Same-day window.
	require.True(t, InQuietHours(at(13, 0), "12:00", "14:00", ""))
	// Timezone: 21:00 UTC is 00:00 in Istanbul (UTC+3).
	require.True(t, InQuietHours(at(21, 0), "23:00", "07:00", "Europe/Istanbul"))
	// Unset / invalid never silences.
	require.False(t, InQuietHours(at(23, 0), "", "08:00", "UTC"))
	require.False(t, InQuietHours(at(23, 0), "25:00", "08:00", "UTC"))
}

func TestParseClock(t *testing.T) {
	m, err := ParseClock("22:30")
	require.NoError(t, err)
	require.Equal(t, 22*60+30, m)
	for _, bad := range []string{"2230", "24:00", "7:00", "aa:bb", "12:60"} {
		_, err := ParseClock(bad)
		require.Error(t, err, bad)
	}
}

func TestRateAllowed(t *testing.T) {
	now := time.Now()
	recent := now.Add(-time.Minute)
	old := now.Add(-10 * time.Minute)
	require.True(t, RateAllowed("activity_completed", 0, nil, now))
	require.False(t, RateAllowed("activity_completed", 1, &recent, now), "5 min cooldown")
	require.True(t, RateAllowed("activity_completed", 1, &old, now))
	require.False(t, RateAllowed("activity_completed", 6, &old, now), "6/hour cap")
	require.True(t, RateAllowed("system_broadcast", 999, &recent, now), "no rule → unlimited")
}

func TestDecide(t *testing.T) {
	d := Decide("mail_received", ModePushAndInbox, false, false, false)
	require.Equal(t, Decision{Inbox: true, Push: true, Priority: "normal"}, d)

	d = Decide("mail_received", ModePushAndInbox, true, false, false)
	require.True(t, d.Inbox)
	require.False(t, d.Push, "global mute suppresses push, keeps inbox")

	d = Decide("mail_received", ModePushAndInbox, false, true, true)
	require.Equal(t, "silent", d.Priority, "quiet hours demote to silent")

	require.Equal(t, Decision{}, Decide("friend_online", ModeOff, false, false, false))
	require.Equal(t, Decision{Inbox: true}, Decide("market_sold", ModeInboxOnly, false, false, false))
	require.Equal(t, "silent", Decide("achievement_unlocked", ModeSilentPush, false, false, true).Priority)

	d = Decide(KindSystemBroadcast, ModeOff, true, true, false)
	require.Equal(t, Decision{Inbox: true, Push: true, Priority: "high"}, d, "broadcast overrides everything")
}

func TestAPNsSender_RequestShapeAndInvalidToken(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	var gotPath, gotTopic, gotType, gotAuth string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotTopic, gotType, gotAuth = r.URL.Path, r.Header.Get("apns-topic"), r.Header.Get("apns-push-type"), r.Header.Get("authorization")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		if r.URL.Path == "/3/device/dead" {
			w.WriteHeader(http.StatusGone)
			_, _ = w.Write([]byte(`{"reason":"Unregistered"}`))
		}
	}))
	defer srv.Close()

	a := &APNsSender{KeyID: "K", TeamID: "T", Topic: "com.otter.app", Endpoint: srv.URL, key: key}
	p := Push{Token: "abc", Priority: "normal", Payload: Payload{Kind: "mail_received", TitleKey: "t", BodyKey: "b", Category: "social"}}
	require.NoError(t, a.Send(context.Background(), p))
	require.Equal(t, "/3/device/abc", gotPath)
	require.Equal(t, "com.otter.app", gotTopic)
	require.Equal(t, "alert", gotType)
	require.Contains(t, gotAuth, "bearer ")
	require.Equal(t, "mail_received", gotBody["kind"])

	p.Priority = "silent"
	require.NoError(t, a.Send(context.Background(), p))
	require.Equal(t, "background", gotType)

	p.Token = "dead"
	require.ErrorIs(t, a.Send(context.Background(), p), ErrInvalidToken)
}

func TestFCMSender_SendsWithCachedToken(t *testing.T) {
	var sends int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer cached", r.Header.Get("Authorization"))
		require.Equal(t, "/v1/projects/p1/messages:send", r.URL.Path)
		sends++
		if sends == 2 {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"status":"NOT_FOUND","details":[{"errorCode":"UNREGISTERED"}]}}`))
		}
	}))
	defer srv.Close()
	f := &FCMSender{ProjectID: "p1", Endpoint: srv.URL, token: "cached", expires: time.Now().Add(time.Hour)}
	p := Push{Token: "tok", Priority: "high", Payload: Payload{Kind: "market_sold", Category: "social"}}
	require.NoError(t, f.Send(context.Background(), p))
	require.ErrorIs(t, f.Send(context.Background(), p), ErrInvalidToken)
}
