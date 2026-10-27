//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/app"
	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/platform/db"
	"github.com/otter/server/internal/ratelimit"
)

// allowAll keeps rate limits out of the way of functional HTTP tests (the
// limiter itself has unit tests).
type allowAll struct{}

func (allowAll) Consume(context.Context, string, ratelimit.Bucket, int64) (ratelimit.Result, error) {
	return ratelimit.Result{Allowed: true, Remaining: 1}, nil
}

var (
	apiOnce sync.Once
	apiURL  string
	apiErr  error
)

// liveAPI starts the production router (app.Build) in-process once per test
// binary, on its own pool so it outlives individual tests.
func liveAPI(t *testing.T) string {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	apiOnce.Do(func() {
		ctx := context.Background()
		p, err := db.Connect(ctx, url, 8)
		if err != nil {
			apiErr = err
			return
		}
		cfg, err := config.Load()
		if err != nil {
			apiErr = err
			return
		}
		cfg.Env = "dev"
		cfg.ContentPath = filepath.Join("..", "..", "..", "content")
		h, err := app.Build(ctx, app.Options{Config: cfg, Pool: p, RateLimit: allowAll{},
			Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
		if err != nil {
			apiErr = err
			return
		}
		apiURL = httptest.NewServer(h).URL
	})
	require.NoError(t, apiErr)
	return apiURL
}

// hc is a tiny API client: bearer token, Idempotency-Key on writes,
// X-Character-Id once a character is entered.
type hc struct {
	t     *testing.T
	base  string
	token string
	char  string
	name  string
}

func (c *hc) call(method, path string, body any) (int, map[string]any) {
	c.t.Helper()
	var rd io.Reader
	switch b := body.(type) {
	case nil:
	case string:
		rd = strings.NewReader(b) // raw (malformed JSON tests)
	default:
		raw, err := json.Marshal(b)
		require.NoError(c.t, err)
		rd = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, c.base+path, rd)
	require.NoError(c.t, err)
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.char != "" {
		req.Header.Set("X-Character-Id", c.char)
	}
	if method != http.MethodGet {
		req.Header.Set("Idempotency-Key", uuid.NewString())
	}
	res, err := http.DefaultClient.Do(req)
	require.NoError(c.t, err)
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	out := map[string]any{}
	_ = json.Unmarshal(raw, &out)
	if len(out) == 0 && len(raw) > 0 {
		out["_raw"] = string(raw)
	}
	return res.StatusCode, out
}

// ok asserts a 2xx and returns the JSON body.
func (c *hc) ok(method, path string, body any) map[string]any {
	c.t.Helper()
	code, out := c.call(method, path, body)
	require.True(c.t, code >= 200 && code < 300, "%s %s → %d %v", method, path, code, out)
	return out
}

// fail asserts the given status (and error code when non-empty).
func (c *hc) fail(method, path string, body any, status int, code string) {
	c.t.Helper()
	got, out := c.call(method, path, body)
	require.Equal(c.t, status, got, "%s %s → %v", method, path, out)
	if code != "" {
		require.Equal(c.t, code, errCode(out), "%s %s", method, path)
	}
}

func errCode(out map[string]any) string {
	if e, ok := out["error"].(map[string]any); ok {
		if c, ok := e["code"].(string); ok {
			return c
		}
	}
	if c, ok := out["code"].(string); ok {
		return c
	}
	return ""
}

// newHTTPPlayer registers an account, creates a character and enters it.
func newHTTPPlayer(t *testing.T, class string) *hc {
	t.Helper()
	c := &hc{t: t, base: liveAPI(t)}
	suffix := randHex(5)
	s := c.ok("POST", "/v1/auth/register", map[string]any{"accept_terms": true,
		"email": "h" + suffix + "@example.com", "username": "h" + suffix, "password": "pass-" + suffix + "-word",
	})
	c.token = s["access_token"].(string)
	c.name = "H" + suffix
	ch := c.ok("POST", "/v1/character", map[string]any{"name": c.name, "race_code": "human", "class_code": class})
	c.char = ch["id"].(string)
	c.ok("POST", "/v1/character/"+c.char+"/enter", map[string]any{})
	return c
}

func (c *hc) path(p string) string { return "/v1/character/" + c.char + p }

func items(out map[string]any, key string) []map[string]any {
	raw, _ := out[key].([]any)
	res := make([]map[string]any, 0, len(raw))
	for _, r := range raw {
		if m, ok := r.(map[string]any); ok {
			res = append(res, m)
		}
	}
	return res
}

func TestHTTPAuthAndCharacters(t *testing.T) {
	base := liveAPI(t)
	anon := &hc{t: t, base: base}
	anon.ok("GET", "/live", nil)
	anon.ok("GET", "/v1/version", nil)
	anon.fail("GET", "/v1/character", nil, 401, "")
	anon.fail("POST", "/v1/auth/login", map[string]any{"email_or_username": "nobody" + randHex(3), "password": "x"}, 401, "E_INVALID_CREDENTIALS")
	anon.fail("POST", "/v1/auth/register", `{"email":`, 400, "")
	anon.ok("POST", "/v1/auth/password/reset-request", map[string]any{"email": "nobody@example.com"})
	anon.fail("POST", "/v1/auth/password/reset-confirm", map[string]any{"token": "nope", "new_password": "longenough1"}, 400, "")
	anon.fail("POST", "/v1/auth/verify-email", map[string]any{"token": "nope"}, 400, "")
	for _, page := range []string{"/verify-email?token=x", "/reset-password?token=x"} {
		code, _ := anon.call("GET", page, nil)
		require.Less(t, code, 500, page)
	}

	a := newHTTPPlayer(t, "warrior")
	suffix := randHex(5)
	reg := anon.ok("POST", "/v1/auth/register", map[string]any{"accept_terms": true,
		"email": "l" + suffix + "@example.com", "username": "l" + suffix, "password": "pass-" + suffix + "-word",
	})
	login := anon.ok("POST", "/v1/auth/login", map[string]any{"email_or_username": "l" + suffix, "password": "pass-" + suffix + "-word"})
	require.NotEmpty(t, login["refresh_token"])
	refreshed := anon.ok("POST", "/v1/auth/refresh", map[string]any{"refresh_token": login["refresh_token"]})
	user := &hc{t: t, base: base, token: refreshed["access_token"].(string)}
	// Reusing a rotated refresh token revokes the whole session (DD-189).
	other := anon.ok("POST", "/v1/auth/login", map[string]any{"email_or_username": "l" + suffix, "password": "pass-" + suffix + "-word"})
	stolen := &hc{t: t, base: base, token: anon.ok("POST", "/v1/auth/refresh", map[string]any{"refresh_token": other["refresh_token"]})["access_token"].(string)}
	anon.fail("POST", "/v1/auth/refresh", map[string]any{"refresh_token": other["refresh_token"]}, 401, "")
	stolen.fail("GET", "/v1/character", nil, 401, "E_SESSION_REVOKED")
	user.ok("POST", "/v1/auth/resend-verification", map[string]any{})
	user.fail("POST", "/v1/auth/password/change", map[string]any{"current_password": "wrong", "new_password": "another-long-1"}, 401, "")
	user.ok("POST", "/v1/auth/password/change", map[string]any{"current_password": "pass-" + suffix + "-word", "new_password": "another-long-1"})
	_ = reg

	// Characters: validation, list, sheet, rename, delete grace.
	user.fail("POST", "/v1/character", map[string]any{"name": "1bad", "race_code": "human", "class_code": "mage"}, 400, "")
	user.fail("POST", "/v1/character", map[string]any{"name": a.name, "race_code": "human", "class_code": "mage"}, 409, "")
	ch := user.ok("POST", "/v1/character", map[string]any{"name": "R" + randHex(5), "race_code": "elf", "class_code": "mage"})
	id := ch["id"].(string)
	list := user.ok("GET", "/v1/character", nil)
	require.Len(t, items(list, "characters"), 1)
	user.ok("GET", "/v1/character/"+id+"/sheet", nil)
	user.ok("GET", "/v1/character/"+id+"/wallet", nil)
	renamed := "N" + randHex(5)
	user.ok("POST", "/v1/character/"+id+"/rename", map[string]any{"name": renamed})
	user.ok("POST", "/v1/character/"+id+"/delete", map[string]any{"confirm_name": renamed})
	require.Empty(t, items(user.ok("GET", "/v1/character", nil), "characters"))
	// Someone else's character is invisible.
	user.fail("GET", "/v1/character/"+a.char+"/sheet", nil, 404, "")
	user.fail("GET", "/v1/character/not-a-uuid/sheet", nil, 400, "")
	user.ok("POST", "/v1/auth/logout", map[string]any{})
}

func TestHTTPGameplayLoop(t *testing.T) {
	p := newHTTPPlayer(t, "warrior")
	dev := func(copper int) {
		p.ok("POST", "/v1/dev/mint", map[string]any{"character_id": p.char, "amount": copper, "currency": "gold"})
	}

	// Gathering.
	p.fail("POST", p.path("/activity"), map[string]any{"kind": "dancing", "duration_seconds": 60}, 400, "")
	act := p.ok("POST", p.path("/activity"), map[string]any{"kind": "mining", "duration_seconds": 60})
	actID := act["contract"].(map[string]any)["id"].(string)
	p.ok("GET", p.path("/activity/active"), nil)
	p.ok("POST", "/v1/activity/"+actID+"/reveal", map[string]any{})
	p.ok("GET", p.path("/activity/history"), nil)
	act2 := p.ok("POST", p.path("/activity"), map[string]any{"kind": "fishing", "duration_seconds": 60})
	p.ok("POST", "/v1/activity/"+act2["contract"].(map[string]any)["id"].(string)+"/stop", map[string]any{})

	// Hunt + combat sim.
	p.ok("GET", "/v1/mobs", nil)
	p.fail("POST", p.path("/hunt"), map[string]any{"mob_code": "dragon", "duration_seconds": 600}, 400, "E_UNKNOWN_MOB")
	h := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
	hID := h["contract"].(map[string]any)["id"].(string)
	p.ok("GET", p.path("/hunt/active"), nil)
	p.ok("POST", "/v1/hunt/"+hID+"/reveal", map[string]any{})
	p.ok("GET", p.path("/hunt/history"), nil)
	h2 := p.ok("POST", p.path("/hunt"), map[string]any{"mob_code": "young_boar", "duration_seconds": 600})
	p.ok("POST", "/v1/hunt/"+h2["contract"].(map[string]any)["id"].(string)+"/stop", map[string]any{})
	p.ok("POST", "/v1/dev/combat/simulate", map[string]any{"character_id": p.char, "mob_code": "young_boar", "seed": 7})

	// Dungeon gate, world travel.
	cat := p.ok("GET", "/v1/dungeon/catalog", nil)
	require.NotEmpty(t, items(cat, "dungeons"))
	p.fail("POST", p.path("/dungeon"), map[string]any{"dungeon_code": "crypt_of_ashes"}, 403, "")
	p.ok("GET", p.path("/dungeon/active"), nil)
	p.ok("GET", p.path("/dungeon/history"), nil)
	p.ok("GET", "/v1/world/regions", nil)
	loc := p.ok("GET", p.path("/location"), nil)
	var dest string
	for _, n := range items(loc, "neighbors") {
		if lv, _ := n["min_level"].(float64); lv <= 1 {
			dest = n["code"].(string)
			break
		}
	}
	require.NotEmpty(t, dest)
	tr := p.ok("POST", p.path("/travel"), map[string]any{"to_region_code": dest})
	trID := tr["travel"].(map[string]any)["id"].(string)
	p.fail("POST", "/v1/travel/"+trID+"/complete", map[string]any{}, 409, "")
	p.ok("POST", "/v1/travel/"+trID+"/cancel", map[string]any{})
	p.ok("GET", p.path("/travel/history"), nil)

	// Vendor, bag, equipment, consumables.
	dev(200_000)
	shop := p.ok("GET", "/v1/vendor/general_ironhold", nil)
	require.NotEmpty(t, shop)
	priceOf := func(code string) float64 {
		for _, s := range items(shop, "sells") {
			if s["template_code"] == code {
				return s["buy_price_copper"].(float64)
			}
		}
		t.Fatalf("no %s in stock: %v", code, shop)
		return 0
	}
	p.fail("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"template_code": "chest_leather",
		"expected_price_copper": 1, "character_id": p.char}, 400, "E_VALIDATION")
	for _, code := range []string{"chest_novice", "potion_hp_small", "potion_hp_small"} { // warriors wear plate (DD-271)
		p.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"template_code": code,
			"expected_price_copper": priceOf(code), "character_id": p.char})
	}
	bag := p.ok("GET", p.path("/inventory"), nil)
	find := func(code string) string {
		for _, it := range items(p.ok("GET", p.path("/inventory"), nil), "items") {
			if it["template_code"] == code {
				return it["instance_id"].(string)
			}
		}
		t.Fatalf("no %s in bag: %v", code, bag)
		return ""
	}
	chest := find("chest_novice")
	p.ok("POST", p.path("/equip/"+chest), map[string]any{})
	p.ok("GET", p.path("/equipment"), nil)
	p.ok("POST", p.path("/unequip/chest"), map[string]any{})
	p.fail("POST", p.path("/unequip/chest"), map[string]any{}, 404, "")
	potion := find("potion_hp_small")
	p.fail("POST", p.path("/items/"+potion+"/use"), map[string]any{}, 409, "E_HP_FULL")
	p.ok("POST", p.path("/inventory/"+potion+"/move"), map[string]any{"target_slot": 40})
	p.ok("POST", p.path("/inventory/"+potion+"/split"), map[string]any{"amount": 1, "target_slot": 41})
	p.ok("GET", p.path("/items"), nil)
	p.ok("POST", "/v1/vendor/general_ironhold/sell", map[string]any{"item_instance_id": chest, "character_id": p.char}) // worn once (bound), still sells (DD-245)
	sold := p.ok("POST", "/v1/vendor/general_ironhold/sell", map[string]any{"item_instance_id": potion, "character_id": p.char})
	require.NotZero(t, sold["received_copper"])
	bbs := items(p.ok("GET", p.path("/buybacks"), nil), "buybacks")
	require.NotEmpty(t, bbs)
	var chestBB string
	for _, bb := range bbs {
		if bb["template_code"] == "chest_novice" {
			chestBB = bb["id"].(string)
		}
	}
	back := p.ok("POST", "/v1/vendor/general_ironhold/buyback/"+chestBB, map[string]any{"character_id": p.char})
	p.ok("DELETE", p.path("/inventory/"+back["instance_id"].(string)), nil)

	// Crafting.
	p.ok("GET", "/v1/professions", nil)
	p.ok("GET", "/v1/recipes?profession=alchemy", nil)
	p.ok("POST", p.path("/professions/learn"), map[string]any{"profession_code": "alchemy"})
	p.ok("GET", p.path("/professions"), nil)
	p.fail("POST", p.path("/craft"), map[string]any{"recipe_code": "recipe_potion_hp_small"}, 409, "")

	// Progression, quests, events, class skills.
	p.ok("GET", p.path("/progression"), nil)
	p.ok("GET", p.path("/achievements"), nil)
	p.ok("GET", p.path("/titles"), nil)
	p.fail("POST", p.path("/titles/active"), map[string]any{"title_code": "bossbane"}, 409, "E_TITLE_NOT_OWNED")
	p.ok("GET", p.path("/quests"), nil)
	p.ok("POST", p.path("/quests/main_01_first_hunt/accept"), map[string]any{})
	p.fail("POST", p.path("/quests/main_01_first_hunt/turn-in"), map[string]any{}, 409, "E_QUEST_INCOMPLETE")
	p.ok("POST", p.path("/quests/main_01_first_hunt/abandon"), map[string]any{})
	p.ok("GET", p.path("/world-events"), nil)
	p.fail("POST", p.path("/world-events/"+uuid.NewString()+"/attend"), map[string]any{}, 404, "")
	p.ok("GET", "/v1/events/seasonal", nil)
	p.ok("GET", "/v1/classes/warrior/skills", nil)
	p.ok("POST", "/v1/dev/content/reload", map[string]any{})
	p.ok("POST", "/v1/echo", map[string]any{"hello": "world"})
}

func TestHTTPSocialAndEconomy(t *testing.T) {
	a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	a.ok("POST", "/v1/dev/mint", map[string]any{"character_id": a.char, "amount": 500_000, "currency": "gold"})

	// Friends + blocks.
	a.ok("POST", a.path("/friends/request"), map[string]any{"target_name": b.name})
	reqs := items(b.ok("GET", b.path("/friends/requests"), nil), "requests")
	require.Len(t, reqs, 1)
	b.ok("POST", b.path("/friends/requests/"+reqs[0]["id"].(string)+"/accept"), map[string]any{})
	friends := items(a.ok("GET", a.path("/friends"), nil), "friends")
	require.Len(t, friends, 1)
	a.ok("POST", a.path("/friends/remove"), map[string]any{"target_character_id": b.char})
	a.ok("POST", a.path("/friends/request"), map[string]any{"target_name": b.name})
	reqs = items(b.ok("GET", b.path("/friends/requests"), nil), "requests")
	b.ok("POST", b.path("/friends/requests/"+reqs[0]["id"].(string)+"/decline"), map[string]any{})
	a.ok("POST", a.path("/blocks"), map[string]any{"target_name": b.name})
	require.Len(t, items(a.ok("GET", a.path("/blocks"), nil), "blocked"), 1)
	a.ok("POST", a.path("/blocks/remove"), map[string]any{"target_character_id": b.char})

	// Chat.
	a.ok("POST", a.path("/chat"), map[string]any{"channel": "global", "body": "hi all"})
	a.ok("POST", a.path("/chat"), map[string]any{"channel": "whisper", "body": "psst", "target_name": b.name})
	b.ok("GET", b.path("/chat/global"), nil)
	b.ok("GET", b.path("/chat/whisper?target="+a.name), nil)
	a.fail("POST", a.path("/chat"), map[string]any{"channel": "guild", "body": "x"}, 409, "E_NOT_IN_CHANNEL")

	// Mail with gold.
	a.ok("POST", a.path("/mail"), map[string]any{"recipient_name": b.name, "subject": "gift", "body": "for you", "gold_copper": 1000})
	inbox := items(b.ok("GET", b.path("/mail"), nil), "mail")
	require.NotEmpty(t, inbox)
	mailID := inbox[0]["id"].(string)
	b.ok("GET", b.path("/mail/"+mailID), nil)
	b.fail("POST", b.path("/mail/"+mailID+"/delete"), map[string]any{}, 409, "E_HAS_ATTACHMENTS")
	b.ok("POST", b.path("/mail/"+mailID+"/claim"), map[string]any{})
	b.ok("POST", b.path("/mail/"+mailID+"/delete"), map[string]any{})

	// Market.
	var swordPrice float64
	for _, o := range items(a.ok("GET", "/v1/vendor/general_ironhold", nil), "sells") {
		if o["template_code"] == "sword_training" {
			swordPrice = o["buy_price_copper"].(float64)
		}
	}
	a.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"template_code": "sword_training",
		"expected_price_copper": swordPrice, "character_id": a.char})
	var sword string
	for _, it := range items(a.ok("GET", a.path("/inventory"), nil), "items") {
		if it["template_code"] == "sword_training" {
			sword = it["instance_id"].(string)
		}
	}
	require.NotEmpty(t, sword)
	l := a.ok("POST", a.path("/market/sell"), map[string]any{"item_instance_id": sword, "stack_quantity": 1,
		"unit_price_copper": 5000, "duration_hours": 12})
	lid := l["id"].(string)
	a.ok("GET", a.path("/market/listings"), nil)
	b.ok("GET", "/v1/market/listings?template_code=sword_training&sort=price_asc", nil)
	b.ok("GET", "/v1/market/history/sword_training", nil)
	a.fail("POST", a.path("/market/buy/"+lid), map[string]any{"expected_unit_price_copper": 5000}, 409, "")
	a.ok("POST", a.path("/market/cancel/"+lid), map[string]any{})

	// Party + guild over HTTP.
	a.ok("POST", a.path("/party"), map[string]any{})
	a.ok("POST", a.path("/party/invite"), map[string]any{"target_name": b.name})
	pinv := items(b.ok("GET", b.path("/party/invites"), nil), "invites")
	b.ok("POST", b.path("/party/invites/"+pinv[0]["id"].(string)+"/accept"), map[string]any{})
	a.ok("GET", a.path("/party"), nil)
	a.ok("POST", a.path("/party/promote"), map[string]any{"target_character_id": b.char})
	b.ok("POST", b.path("/party/kick"), map[string]any{"target_character_id": a.char})
	// Down to one member the party is gone (DD-297).
	b.fail("POST", b.path("/party/disband"), map[string]any{}, 409, "E_NOT_IN_PARTY")

	g := a.ok("POST", a.path("/guild"), map[string]any{"name": "Http " + randHex(3), "tag": randTag()})
	gid := g["guild"].(map[string]any)["id"].(string)
	a.ok("GET", "/v1/guilds/"+gid, nil)
	a.ok("POST", a.path("/guild/invite"), map[string]any{"target_name": b.name})
	ginv := items(b.ok("GET", b.path("/guild/invites"), nil), "invites")
	b.ok("POST", b.path("/guild/invites/"+ginv[0]["id"].(string)+"/accept"), map[string]any{})
	a.ok("POST", a.path("/guild/rank"), map[string]any{"target_character_id": b.char, "rank_index": 1})
	a.ok("POST", a.path("/guild/motd"), map[string]any{"motd": "raid"})
	a.ok("POST", a.path("/guild/transfer"), map[string]any{"target_character_id": b.char})
	b.ok("POST", b.path("/guild/kick"), map[string]any{"target_character_id": a.char})
	b.ok("POST", b.path("/guild/disband"), map[string]any{})
	b.ok("POST", b.path("/guild/disband/cancel"), map[string]any{})
	b.ok("GET", b.path("/guild"), nil)

	// Notifications (b got mail/party/guild notices).
	n := b.ok("GET", "/v1/notifications", nil)
	ns := items(n, "notifications")
	require.NotEmpty(t, ns)
	b.ok("POST", "/v1/notifications/"+ns[0]["id"].(string)+"/read", map[string]any{})
	b.ok("POST", "/v1/notifications/read-all", map[string]any{})
	prefs := b.ok("GET", "/v1/notification/preference", nil)
	require.NotNil(t, prefs["kinds"])
	b.ok("PUT", "/v1/notification/preference", map[string]any{"kinds": map[string]any{"market_sold": "inbox_only"}})
	b.fail("PUT", "/v1/notification/preference", map[string]any{"kinds": map[string]any{"market_sold": "loud"}}, 400, "")
	b.ok("POST", "/v1/notification/preference/global-mute", map[string]any{"hours": 2})
	b.ok("DELETE", "/v1/notification/preference/global-mute", nil)
	dev := uuid.NewString()
	b.ok("POST", "/v1/device/push-token", map[string]any{"device_id": dev, "platform": "android",
		"push_token": "tok-" + randHex(8), "locale": "tr", "app_version": "1.0.0"})
	b.ok("DELETE", "/v1/device/push-token/"+dev, nil)
	b.ok("GET", "/v1/ws/presence", nil)
	_ = time.Now
}
