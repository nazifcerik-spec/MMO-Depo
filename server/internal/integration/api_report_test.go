//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/report"
)

// DD-355: a player reports another; repeats merge; operators resolve; the
// excerpt comes from the server's log.
func TestHTTPPlayerReport(t *testing.T) {
	w := newWorld(t)
	a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	send := func(c *hc, body map[string]any, status int, code string) {
		c.fail("POST", "/v1/reports", body, status, code)
	}
	send(a, map[string]any{"target_name": b.name, "category": "nope"}, 400, "E_VALIDATION")
	send(a, map[string]any{"target_name": "NoSuchPlayerHere", "category": "spam"}, 404, "E_NOT_FOUND")
	send(a, map[string]any{"target_name": a.name, "category": "spam"}, 409, "E_REPORT_SELF")

	var msg int64
	require.NoError(t, w.p.QueryRow(w.ctx, `INSERT INTO social.chat_message (channel_kind, sender_character_id, body)
		VALUES ('global', $1, 'rude words') RETURNING id`, b.char).Scan(&msg))
	for range 2 {
		send(a, map[string]any{"character_id": a.char, "target_name": b.name, "category": "harassment", "note": " bad ",
			"chat_message_id": msg}, 202, "")
	}
	var n int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM ops.player_report WHERE target_character = $1`, b.char).Scan(&n))
	require.Equal(t, 1, n, "a repeat while open changes nothing")

	svc := report.NewService(w.p, nil, nil)
	list, err := svc.Queue(context.Background(), "open", 200)
	require.NoError(t, err)
	var mine *report.Report
	for i := range list {
		if list[i].TargetName == b.name {
			mine = &list[i]
		}
	}
	require.NotNil(t, mine)
	require.Equal(t, "bad", mine.Note)
	require.Equal(t, "rude words", *mine.ChatExcerpt)
	require.Equal(t, 1, mine.OpenAgainst)

	admin := uuid.Nil
	require.ErrorIs(t, svc.Resolve(w.ctx, admin, mine.ID, "maybe", ""), report.ErrInvalid)
	require.ErrorIs(t, svc.Resolve(w.ctx, admin, uuid.New(), "dismissed", ""), report.ErrNotFound)
}
