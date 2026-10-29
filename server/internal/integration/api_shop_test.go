//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Credit shop + cosmetics (DD-321): credits buy boosts and premium
// cosmetics; free cosmetics need no purchase; nothing is bought twice.
func TestHTTPShop(t *testing.T) {
	newWorld(t)
	p := newHTTPPlayer(t, "mage")
	shop := p.ok("GET", p.path("/shop"), nil)
	require.EqualValues(t, 0, shop["credits"])
	require.Equal(t, false, shop["payments_enabled"])
	p.fail("POST", p.path("/shop/boost_gather/buy"), map[string]any{}, 409, "E_INSUFFICIENT_CREDITS")

	p.ok("POST", "/v1/dev/mint", map[string]any{"character_id": p.char, "amount": 100, "currency": "credit"})
	shop = p.ok("POST", p.path("/shop/boost_gather/buy"), map[string]any{})
	require.EqualValues(t, 80, shop["credits"])
	boosts := p.ok("GET", p.path("/boosts"), nil)["boosts"].([]any)
	require.NotNil(t, boosts[0].(map[string]any)["ends_at"], "the bought boost runs")
	p.fail("POST", p.path("/shop/boost_gather/buy"), map[string]any{}, 409, "E_BOOST_ACTIVE")
	require.EqualValues(t, 80, p.ok("GET", p.path("/shop"), nil)["credits"], "refunded")

	cos := p.ok("GET", p.path("/cosmetics"), nil)
	require.Equal(t, "av01", cos["avatar"])
	require.Equal(t, "bg01", cos["background"])
	p.ok("POST", p.path("/cosmetics/select"), map[string]any{"code": "av07"})
	p.fail("POST", p.path("/cosmetics/select"), map[string]any{"code": "bg20"}, 409, "E_COSMETIC_LOCKED")
	p.ok("POST", p.path("/shop/bg20/buy"), map[string]any{})
	p.fail("POST", p.path("/shop/bg20/buy"), map[string]any{}, 409, "E_OWNED")
	cos = p.ok("POST", p.path("/cosmetics/select"), map[string]any{"code": "bg20"})
	require.Equal(t, "bg20", cos["background"])
	require.Equal(t, "av07", cos["avatar"])
	require.Contains(t, cos["owned"], "bg20")
	require.EqualValues(t, 50, p.ok("GET", p.path("/shop"), nil)["credits"])
	p.fail("POST", p.path("/cosmetics/select"), map[string]any{"code": "zz99"}, 404, "E_NOT_FOUND")
}
