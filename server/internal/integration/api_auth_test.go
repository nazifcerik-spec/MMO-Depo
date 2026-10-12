//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Auth error paths through the real router: every one must be a precise 4xx.
func TestHTTPAuthEdges(t *testing.T) {
	w := newWorld(t)
	anon := &hc{t: t, base: liveAPI(t)}
	suffix := randHex(5)
	email, user, pass := "e"+suffix+"@example.com", "e"+suffix, "pass-"+suffix+"-word"
	reg := anon.ok("POST", "/v1/auth/register", map[string]any{"accept_terms": true, "email": email, "username": user, "password": pass})
	acct := reg["account_id"].(string)

	anon.fail("POST", "/v1/auth/register", map[string]any{"accept_terms": true, "email": email, "username": "x" + suffix, "password": pass}, 409, "E_ACCOUNT_EXISTS")
	anon.fail("POST", "/v1/auth/register", map[string]any{"accept_terms": true, "email": "not-an-email", "username": "y" + suffix, "password": pass}, 400, "E_VALIDATION")
	anon.fail("POST", "/v1/auth/register", map[string]any{"accept_terms": true, "email": "z" + suffix + "@example.com", "username": "z" + suffix, "password": "short"}, 400, "E_VALIDATION")
	anon.fail("POST", "/v1/auth/refresh", map[string]any{"refresh_token": "unknown-" + suffix}, 401, "E_REFRESH_TOKEN_UNKNOWN")
	anon.fail("POST", "/v1/auth/refresh", "{", 400, "")
	anon.fail("POST", "/v1/auth/login", "", 400, "")
	anon.fail("POST", "/v1/auth/logout", map[string]any{}, 401, "")

	// Expired refresh token and expired session.
	s := anon.ok("POST", "/v1/auth/login", map[string]any{"email_or_username": user, "password": pass})
	_, err := w.p.Exec(w.ctx, `UPDATE identity.refresh_token SET expires_at = now() - interval '1 minute' WHERE session_id = $1`, s["session_id"])
	require.NoError(t, err)
	anon.fail("POST", "/v1/auth/refresh", map[string]any{"refresh_token": s["refresh_token"]}, 401, "E_TOKEN_EXPIRED")
	_, err = w.p.Exec(w.ctx, `UPDATE identity.session SET expires_at = now() - interval '1 minute' WHERE id = $1`, s["session_id"])
	require.NoError(t, err)
	(&hc{t: t, base: anon.base, token: s["access_token"].(string)}).fail("GET", "/v1/character", nil, 401, "E_SESSION_EXPIRED")
	anon.fail("GET", "/v1/character", nil, 401, "")
	(&hc{t: t, base: anon.base, token: "not.a.jwt"}).fail("GET", "/v1/character", nil, 401, "")

	// Session housekeeping: e-mail verified, logout ends the session.
	c := &hc{t: t, base: anon.base, token: anon.ok("POST", "/v1/auth/login", map[string]any{"email_or_username": email, "password": pass})["access_token"].(string)}
	_, err = w.p.Exec(w.ctx, `UPDATE identity.account SET email_verified_at = now() WHERE id = $1`, acct)
	require.NoError(t, err)
	c.fail("POST", "/v1/auth/resend-verification", map[string]any{}, 409, "E_ALREADY_VERIFIED")
	c.fail("POST", "/v1/auth/password/change", "{", 400, "E_INVALID_JSON")
	c.fail("POST", "/v1/auth/password/change", map[string]any{"current_password": pass, "new_password": "short"}, 400, "E_VALIDATION")
	c.ok("POST", "/v1/auth/logout", map[string]any{})
	c.fail("POST", "/v1/auth/logout", map[string]any{}, 401, "E_SESSION_REVOKED")

	// Locked and banned accounts can't sign in.
	for state, code := range map[string]string{"locked": "E_ACCOUNT_LOCKED", "banned": "E_ACCOUNT_BANNED"} {
		_, err = w.p.Exec(w.ctx, `UPDATE identity.account SET state = $2 WHERE id = $1`, acct, state)
		require.NoError(t, err)
		anon.fail("POST", "/v1/auth/login", map[string]any{"email_or_username": user, "password": pass}, 403, code)
	}
}
