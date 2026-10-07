//go:build integration

package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/notification"
)

type fakeSender struct {
	err   error
	sends int
}

func (f *fakeSender) Send(context.Context, notification.Push) error {
	f.sends++
	return f.err
}

func notice(acct uuid.UUID, kind string) notification.Notice {
	return notification.Notice{AccountID: acct, Kind: kind, DedupKey: kind + ":" + uuid.NewString(),
		TitleKey: "t", BodyKey: "b", Category: "social"}
}

func pushState(t *testing.T, w *testEnv, acct uuid.UUID, dedup string) (queued bool, ok int, reason *string) {
	t.Helper()
	err := w.p.QueryRow(w.ctx, `
		SELECT true, COALESCE(devices_ok, 0), fail_reason FROM notification.push_message
		 WHERE account_id = $1 AND dedup_key = $2`, acct, dedup).Scan(&queued, &ok, &reason)
	if err != nil {
		return false, 0, nil
	}
	return queued, ok, reason
}

func TestNotificationQueueAndDispatch(t *testing.T) {
	w := newWorld(t)
	svc := notification.NewService(w.p, nil)
	p := w.player("warrior", 1)

	require.ErrorIs(t, svc.Notify(w.ctx, notice(p.Account, "no_such_kind")), notification.ErrUnknownKind)

	// No device yet: the push is consumed as "no_devices", the inbox keeps it.
	n1 := notice(p.Account, "market_sold")
	require.NoError(t, svc.Notify(w.ctx, n1))
	require.NoError(t, svc.Notify(w.ctx, n1), "same dedup key is a no-op")
	items, unread, err := svc.Inbox(w.ctx, p.Account, false)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, 1, unread)
	ok := &fakeSender{}
	_, err = notification.Dispatch(w.ctx, w.p, map[string]notification.Sender{"android": ok}, 500)
	require.NoError(t, err)
	_, _, reason := pushState(t, w, p.Account, n1.DedupKey)
	require.Equal(t, "no_devices", *reason)

	// With a device the push is delivered.
	dev := uuid.New()
	require.ErrorIs(t, svc.RegisterDevice(w.ctx, p.Account, notification.DeviceRequest{DeviceID: dev, Platform: "web", PushToken: "x"}),
		notification.ErrBadPlatform)
	require.NoError(t, svc.RegisterDevice(w.ctx, p.Account, notification.DeviceRequest{DeviceID: dev, Platform: "android", PushToken: "tok-" + randHex(4)}))
	n2 := notice(p.Account, "mail_received")
	require.NoError(t, svc.Notify(w.ctx, n2))
	_, err = notification.Dispatch(w.ctx, w.p, map[string]notification.Sender{"android": ok}, 500)
	require.NoError(t, err)
	_, delivered, _ := pushState(t, w, p.Account, n2.DedupKey)
	require.Equal(t, 1, delivered)
	require.Positive(t, ok.sends)

	// Inbox-only kinds and global mute skip the push but keep the inbox.
	_, err = svc.UpdatePreferences(w.ctx, p.Account, notification.PreferenceUpdate{
		Kinds: map[string]notification.Mode{"guild_invite": notification.ModeInboxOnly}})
	require.NoError(t, err)
	n3 := notice(p.Account, "guild_invite")
	require.NoError(t, svc.Notify(w.ctx, n3))
	queued, _, _ := pushState(t, w, p.Account, n3.DedupKey)
	require.False(t, queued)
	_, err = svc.SetGlobalMute(w.ctx, p.Account, 2)
	require.NoError(t, err)
	n4 := notice(p.Account, "party_invite")
	require.NoError(t, svc.Notify(w.ctx, n4))
	queued, _, _ = pushState(t, w, p.Account, n4.DedupKey)
	require.False(t, queued, "global mute holds pushes")
	_, err = svc.SetGlobalMute(w.ctx, p.Account, 0)
	require.NoError(t, err)
	_, err = svc.SetGlobalMute(w.ctx, p.Account, 999)
	require.ErrorIs(t, err, notification.ErrBadMuteHours)
	_, err = svc.UpdatePreferences(w.ctx, p.Account, notification.PreferenceUpdate{
		Kinds: map[string]notification.Mode{"guild_invite": "loud"}})
	require.ErrorIs(t, err, notification.ErrBadMode)
	bad := "25:00"
	_, err = svc.UpdatePreferences(w.ctx, p.Account, notification.PreferenceUpdate{QuietHoursStart: &bad, QuietHoursEnd: &bad})
	require.ErrorIs(t, err, notification.ErrBadClock)

	// Three invalid-token answers deactivate the device.
	invalid := &fakeSender{err: notification.ErrInvalidToken}
	for _, kind := range []string{"activity_completed", "dungeon_completed", "travel_completed"} {
		require.NoError(t, svc.Notify(w.ctx, notice(p.Account, kind)))
		_, err = notification.Dispatch(w.ctx, w.p, map[string]notification.Sender{"android": invalid}, 500)
		require.NoError(t, err)
	}
	var active bool
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT is_active FROM notification.device_push_token WHERE device_id = $1`, dev).Scan(&active))
	require.False(t, active)

	// Provider errors are recorded, not retried forever.
	require.NoError(t, svc.RegisterDevice(w.ctx, p.Account, notification.DeviceRequest{DeviceID: dev, Platform: "android", PushToken: "tok2"}))
	n5 := notice(p.Account, "world_event_active")
	require.NoError(t, svc.Notify(w.ctx, n5))
	_, err = notification.Dispatch(w.ctx, w.p, map[string]notification.Sender{"android": &fakeSender{err: errors.New("503")}}, 500)
	require.NoError(t, err)
	_, _, reason = pushState(t, w, p.Account, n5.DedupKey)
	require.Equal(t, "provider_error", *reason)

	marked, err := svc.MarkRead(w.ctx, p.Account, uuid.Nil)
	require.NoError(t, err)
	require.Positive(t, marked)
	require.NoError(t, svc.UnregisterDevice(w.ctx, p.Account, dev))
}
