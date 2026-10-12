//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/privacy"
	"github.com/otter/server/internal/worker"
)

func get(t *testing.T, url string) (int, string, http.Header) {
	t.Helper()
	res, err := http.Get(url)
	require.NoError(t, err)
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b), res.Header
}

// Account rights end to end (DD-231): terms at sign-up, data export, deletion
// cancelled by signing in, then erasure after the grace.
func TestHTTPAccountRights(t *testing.T) {
	w := newWorld(t)
	anon := &hc{t: t, base: liveAPI(t)}
	anon.fail("POST", "/v1/auth/register", map[string]any{"email": "nt" + randHex(4) + "@example.com",
		"username": "nt" + randHex(4), "password": "long-enough-1"}, 400, "E_VALIDATION")

	p := newHTTPPlayer(t, "warrior")
	user := "h" + strings.ToLower(p.name[1:])
	pass := "pass-" + strings.ToLower(p.name[1:]) + "-word"
	p.mint(1234)
	p.ok("POST", p.path("/chat"), map[string]any{"channel": "global", "body": "hello export"})

	// Export: a one-time link the browser downloads.
	link := p.ok("POST", "/v1/me/export-link", map[string]any{})
	url := link["url"].(string)
	require.Contains(t, url, "/v1/me/export?token=")
	path := url[strings.Index(url, "/v1/"):]
	code, body, hdr := get(t, p.base+path)
	require.Equal(t, http.StatusOK, code, body)
	require.Contains(t, hdr.Get("Content-Disposition"), "attachment")
	var doc map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(body), &doc))
	var account []map[string]any
	require.NoError(t, json.Unmarshal(doc["account"], &account))
	require.Equal(t, user, account[0]["username"])
	require.Equal(t, auth.DefaultTermsVersion, account[0]["terms_version"])
	require.Contains(t, string(doc["characters"]), p.char)
	require.Contains(t, string(doc["wallets"]), "1234")
	require.Contains(t, string(doc["chat"]), "hello export")
	for _, s := range []string{"sessions", "items", "ledger", "mail", "market", "quests", "support_tickets"} {
		require.Contains(t, doc, s)
	}
	code, _, _ = get(t, p.base+path)
	require.Equal(t, http.StatusBadRequest, code, "the link works once")

	// Deletion: wrong password refused; request signs out; signing in cancels.
	p.fail("DELETE", "/v1/me", map[string]any{"current_password": "nope"}, 401, "E_INVALID_CREDENTIALS")
	del := p.ok("DELETE", "/v1/me", map[string]any{"current_password": pass})
	require.NotEmpty(t, del["purge_after"])
	p.fail("GET", "/v1/character", nil, 401, "")
	back := anon.ok("POST", "/v1/auth/login", map[string]any{"email_or_username": user, "password": pass})
	require.Equal(t, true, back["deletion_cancelled"])
	p.token = back["access_token"].(string)
	p.ok("GET", "/v1/character", nil)

	// Request again and let the grace pass: the account is erased.
	p.ok("DELETE", "/v1/me", map[string]any{"current_password": pass})
	_, err := w.p.Exec(w.ctx, `UPDATE identity.account SET deletion_requested_at = now() - interval '31 days'
	                            WHERE username = $1`, user)
	require.NoError(t, err)
	n, err := privacy.Purge(w.ctx, w.p, time.Now(), auth.DeletionGrace, 100)
	require.NoError(t, err)
	require.GreaterOrEqual(t, n, int64(1))
	anon.fail("POST", "/v1/auth/login", map[string]any{"email_or_username": user, "password": pass}, 401, "E_INVALID_CREDENTIALS")
	var state, email string
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT state, email FROM identity.account WHERE id = $1`, account[0]["id"]).Scan(&state, &email))
	require.Equal(t, "deleted", state)
	require.True(t, strings.HasSuffix(email, "@deleted.invalid"), email)

	// Its characters go through the character-deletion job.
	_, err = worker.CharacterDeleteExpire{Pool: w.p}.Run(w.ctx)
	require.NoError(t, err)
	var chars int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM character.character WHERE id = $1`, p.char).Scan(&chars))
	require.Zero(t, chars)

	n, err = privacy.Purge(w.ctx, w.p, time.Now(), auth.DeletionGrace, 100)
	require.NoError(t, err)
	require.Zero(t, n, "an erased account is not purged twice")
}

func TestHTTPLegalPages(t *testing.T) {
	base := liveAPI(t)
	code, body, hdr := get(t, base+"/terms?lang=tr")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, body, "Kullanım Şartları")
	require.Contains(t, body, "13 yaşında")
	require.Contains(t, body, auth.DefaultTermsVersion)
	require.Contains(t, hdr.Get("Content-Security-Policy"), "default-src 'none'")
	code, body, _ = get(t, base+"/privacy?lang=en")
	require.Equal(t, http.StatusOK, code)
	require.Contains(t, body, `id="delete"`)
	require.Contains(t, body, "Draft", "no operator configured in tests")
	req, _ := http.NewRequest(http.MethodGet, base+"/privacy", nil)
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	require.Contains(t, string(b), "Gizlilik Politikası")
}
