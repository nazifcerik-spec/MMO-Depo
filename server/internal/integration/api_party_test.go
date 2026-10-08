//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPPartyLifecycle(t *testing.T) {
	ps := make([]*hc, 7)
	for i := range ps {
		ps[i] = newHTTPPlayer(t, "warrior")
	}
	lead, b, c := ps[0], ps[1], ps[2]
	pp := func(p *hc, s string) string { return p.path("/party" + s) }
	invite := func(from, to *hc) string {
		return from.ok("POST", pp(from, "/invite"), map[string]any{"target_name": to.name})["id"].(string)
	}

	require.Nil(t, lead.ok("GET", pp(lead, ""), nil)["party"])
	lead.fail("POST", pp(lead, "/invite"), map[string]any{"target_name": b.name}, 409, "E_NOT_IN_PARTY")
	lead.fail("POST", pp(lead, "/leave"), map[string]any{}, 409, "E_NOT_IN_PARTY")
	lead.ok("POST", pp(lead, ""), map[string]any{})
	lead.fail("POST", pp(lead, ""), map[string]any{}, 409, "E_ALREADY_IN_PARTY")

	// Invites: decline, bad ids, self, unknown, duplicate.
	inv := invite(lead, b)
	b.ok("POST", pp(b, "/invites/"+inv+"/decline"), map[string]any{})
	b.fail("POST", pp(b, "/invites/"+inv+"/accept"), map[string]any{}, 409, "E_INVITE_INVALID")
	b.fail("POST", pp(b, "/invites/nope/decline"), map[string]any{}, 400, "E_INVALID_ID")
	b.fail("POST", pp(b, "/invites/nope/accept"), map[string]any{}, 400, "E_INVALID_ID")
	lead.fail("POST", pp(lead, "/invite"), map[string]any{"target_name": lead.name}, 400, "E_SELF_TARGET")
	lead.fail("POST", pp(lead, "/invite"), map[string]any{"target_name": "Nobody" + randHex(3)}, 404, "E_TARGET_NOT_FOUND")
	lead.fail("POST", pp(lead, "/invite"), "{", 400, "E_INVALID_JSON")
	inv = invite(lead, b)
	lead.fail("POST", pp(lead, "/invite"), map[string]any{"target_name": b.name}, 409, "E_ALREADY_INVITED")
	require.Equal(t, 1, len(items(b.ok("GET", pp(b, "/invites"), nil), "invites")))
	b.ok("POST", pp(b, "/invites/"+inv+"/accept"), map[string]any{})
	lead.fail("POST", pp(lead, "/invite"), map[string]any{"target_name": b.name}, 409, "E_TARGET_IN_PARTY")

	// Only the leader manages; non-members can't be kicked or promoted.
	b.fail("POST", pp(b, "/invite"), map[string]any{"target_name": c.name}, 403, "E_NOT_LEADER")
	b.fail("POST", pp(b, "/kick"), map[string]any{"target_character_id": lead.char}, 403, "E_NOT_LEADER")
	lead.fail("POST", pp(lead, "/kick"), map[string]any{"target_character_id": c.char}, 404, "E_NOT_MEMBER")
	lead.fail("POST", pp(lead, "/promote"), map[string]any{"target_character_id": c.char}, 404, "E_NOT_MEMBER")
	lead.fail("POST", pp(lead, "/kick"), map[string]any{"target_character_id": lead.char}, 400, "E_SELF_TARGET")
	lead.fail("POST", pp(lead, "/kick"), "{", 400, "E_INVALID_JSON")
	lead.fail("POST", pp(lead, "/promote"), "{", 400, "E_INVALID_JSON")

	// Size cap: an invite sent while there was room can't be accepted once
	// the party has filled up; with five members no new invite goes out.
	late := invite(lead, ps[5])
	for _, p := range ps[2:5] {
		p.ok("POST", pp(p, "/invites/"+invite(lead, p)+"/accept"), map[string]any{})
	}
	lead.fail("POST", pp(lead, "/invite"), map[string]any{"target_name": ps[6].name}, 409, "E_PARTY_FULL")
	ps[5].fail("POST", pp(ps[5], "/invites/"+late+"/accept"), map[string]any{}, 409, "E_PARTY_FULL")
	party := lead.ok("GET", pp(lead, ""), nil)["party"].(map[string]any)
	require.Len(t, party["members"], 5)

	// Leaving hands over leadership; the last one out disbands.
	lead.ok("POST", pp(lead, "/leave"), map[string]any{})
	require.Nil(t, lead.ok("GET", pp(lead, ""), nil)["party"])
	b.ok("POST", pp(b, "/disband"), map[string]any{})
	require.Nil(t, b.ok("GET", pp(b, ""), nil)["party"])
	ps[6].fail("GET", "/v1/character/"+b.char+"/party", nil, 404, "")
}
