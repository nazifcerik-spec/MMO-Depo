//go:build integration

package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/admin"
)

// Support desk end to end (DD-234): a player files tickets, an operator
// answers from /internal/ops, the player sees the reply and gets notified.
func TestHTTPSupportTickets(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	other := newHTTPPlayer(t, "mage")

	ticket := map[string]any{"category": "bug", "subject": "Sword vanished",
		"body": "After the hunt my sword was gone from the bag.", "character_id": p.char,
		"client_info": map[string]any{"platform": "android", "locale": "tr", "secret": "dropped"}}
	created := p.ok("POST", "/v1/support/tickets", ticket)
	require.Equal(t, "open", created["state"])
	require.Equal(t, map[string]any{"platform": "android", "locale": "tr"}, created["client_info"], "unknown keys are dropped")

	p.fail("POST", "/v1/support/tickets", map[string]any{"category": "bug", "subject": "x", "body": "too short"}, 400, "E_VALIDATION")
	p.fail("POST", "/v1/support/tickets", map[string]any{"category": "nope", "subject": "Valid subject",
		"body": "A long enough message body."}, 400, "E_VALIDATION")
	p.fail("POST", "/v1/support/tickets", map[string]any{"category": "bug", "subject": "Valid subject",
		"body": "A long enough message body.", "character_id": other.char}, 404, "E_NOT_FOUND")
	p.fail("POST", "/v1/support/tickets", "{", 400, "E_INVALID_JSON")
	require.Empty(t, items(other.ok("GET", "/v1/support/tickets", nil), "tickets"), "players only see their own tickets")

	// Operators: a viewer may not, a support agent may.
	viewer := opsLogin(t, w, p.base, "viewer")
	viewer.fail("GET", "/internal/ops/support/tickets", nil, 403, "")
	agent := opsLogin(t, w, p.base, "support_l1")
	// The queue is oldest-first: earlier local runs' leftovers must not push
	// this ticket past the page.
	_, err := w.p.Exec(w.ctx, `UPDATE ops.support_ticket SET state = 'closed' WHERE state = 'open' AND id <> $1`, created["id"])
	require.NoError(t, err)
	queue := items(agent.ok("GET", "/internal/ops/support/tickets?state=open&limit=200", nil), "tickets")
	var found bool
	for _, q := range queue {
		found = found || q["id"] == created["id"]
	}
	require.True(t, found, "the new ticket is in the open queue")

	id := created["id"].(string)
	agent.fail("POST", "/internal/ops/support/tickets/"+id+"/reply", map[string]any{"reply": " "}, 400, "E_VALIDATION")
	agent.fail("POST", "/internal/ops/support/tickets/"+uuid.NewString()+"/reply", map[string]any{"reply": "hi"}, 404, "")
	ans := agent.ok("POST", "/internal/ops/support/tickets/"+id+"/reply", map[string]any{"reply": "Restored it, sorry!"})
	require.Equal(t, "answered", ans["state"])
	mine := items(p.ok("GET", "/v1/support/tickets", nil), "tickets")
	require.Equal(t, "Restored it, sorry!", mine[0]["reply"])
	inbox := p.ok("GET", "/v1/notifications", nil)
	require.True(t, strings.Contains(fmt.Sprint(inbox), "support_reply"), "the player is notified: %v", inbox)

	agent.ok("POST", "/internal/ops/support/tickets/"+id+"/reply", map[string]any{"reply": "Closing.", "close": true})
	agent.fail("POST", "/internal/ops/support/tickets/"+id+"/reply", map[string]any{"reply": "again"}, 409, "E_TICKET_CLOSED")
	var audited int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM audit.admin_action WHERE action = 'support.reply' AND target_id = $1`,
		id).Scan(&audited))
	require.Equal(t, 2, audited)

	// Open-ticket cap (the closed one does not count).
	for i := 0; i < 10; i++ {
		p.ok("POST", "/v1/support/tickets", map[string]any{"category": "feedback", "subject": "Idea " + randHex(2),
			"body": "Please add more dungeons."})
	}
	p.fail("POST", "/v1/support/tickets", map[string]any{"category": "feedback", "subject": "One more",
		"body": "Please add more dungeons."}, 409, "E_TOO_MANY_TICKETS")
}

// opsLogin creates an operator with the given roles and signs in (password + TOTP).
func opsLogin(t *testing.T, w *testEnv, base string, roles ...string) *hc {
	t.Helper()
	email, pass := "ops-"+randHex(4)+"@example.com", "admin-pass-"+randHex(4)+"-x"
	_, secret, err := admin.NewService(w.p, nil, nil).CreateAdmin(w.ctx, email, "Ops", pass, roles)
	require.NoError(t, err)
	code, err := admin.TOTPCode(secret, time.Now().Unix()/30)
	require.NoError(t, err)
	c := &hc{t: t, base: base}
	c.token = c.ok("POST", "/internal/admin/login", map[string]any{"email": email, "password": pass, "code": code})["token"].(string)
	return c
}
