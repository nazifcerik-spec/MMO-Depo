//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/market"
)

func goldOf(t *testing.T, c *hc) int64 {
	t.Helper()
	return c.gold()
}

// DD-359: bids hold gold, the outbid get it back, a late bid extends, the top
// bid buys when time runs out, the seller is paid less the cut, the ledger balances.
func TestHTTPAuction(t *testing.T) {
	w := newWorld(t)
	seller, a, b := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage"), newHTTPPlayer(t, "rogue")
	for _, c := range []*hc{seller, a, b} {
		c.mint(50_000_000)
	}
	seller.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"character_id": seller.char, "template_code": "potion_hp_small", "quantity": 1})
	item := seller.bagItem("potion_hp_small")
	sell := func(body map[string]any) map[string]any {
		body["item_instance_id"] = item
		body["stack_quantity"] = 1
		body["duration_hours"] = 12
		body["auction"] = true
		return seller.ok("POST", seller.path("/market/sell"), body)
	}
	seller.fail("POST", seller.path("/market/sell"), map[string]any{"item_instance_id": item, "stack_quantity": 1, "duration_hours": 12,
		"auction": true, "start_bid_copper": 0}, 400, "E_PRICE_OUT_OF_RANGE")
	l := sell(map[string]any{"start_bid_copper": 100, "buyout_copper": 100000})
	id := l["id"].(string)
	require.Equal(t, "auction", l["kind"])
	require.EqualValues(t, 100, l["min_next_bid_copper"])

	// A bid is not a purchase, and the seller cannot bid.
	a.fail("POST", a.path("/market/buy/"+id), map[string]any{}, 409, "E_AUCTION_ONLY")
	seller.fail("POST", seller.path("/market/bid/"+id), map[string]any{"amount_copper": 500}, 409, "E_SELF_BUY")
	a.fail("POST", a.path("/market/bid/"+id), map[string]any{"amount_copper": 99}, 409, "E_BID_TOO_LOW")
	a.fail("POST", a.path("/market/bid/"+id), map[string]any{"amount_copper": 999_999_999}, 409, "E_INSUFFICIENT_GOLD")

	beforeA := goldOf(t, a)
	v := a.ok("POST", a.path("/market/bid/"+id), map[string]any{"amount_copper": 1000})
	require.EqualValues(t, 1000, v["current_bid_copper"])
	require.EqualValues(t, 1050, v["min_next_bid_copper"], "at least 5% more")
	require.Equal(t, beforeA-1000, goldOf(t, a), "the gold is held")
	a.fail("POST", a.path("/market/bid/"+id), map[string]any{"amount_copper": 2000}, 409, "E_ALREADY_TOP_BIDDER")
	seller.fail("POST", seller.path("/market/cancel/"+id), map[string]any{}, 409, "E_AUCTION_HAS_BIDS")

	// Outbid: the first bidder gets the gold straight back.
	b.fail("POST", b.path("/market/bid/"+id), map[string]any{"amount_copper": 1049}, 409, "E_BID_TOO_LOW")
	beforeB := goldOf(t, b)
	b.ok("POST", b.path("/market/bid/"+id), map[string]any{"amount_copper": 2000})
	require.Equal(t, beforeA, goldOf(t, a))
	require.Equal(t, beforeB-2000, goldOf(t, b))

	// A bid in the last minutes extends the end, a limited number of times.
	_, err := w.p.Exec(w.ctx, `UPDATE market.listing SET expires_at = now() + interval '1 minute' WHERE id = $1`, id)
	require.NoError(t, err)
	a.ok("POST", a.path("/market/bid/"+id), map[string]any{"amount_copper": 3000})
	var left float64
	var ext int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT EXTRACT(EPOCH FROM expires_at - now()), extensions FROM market.listing WHERE id = $1`, id).Scan(&left, &ext))
	require.Greater(t, left, 240.0)
	require.Equal(t, 1, ext)
	require.Equal(t, beforeB, goldOf(t, b), "b was outbid in turn")

	// Time runs out: the worker gives the item to the top bidder and pays the seller.
	beforeSeller := goldOf(t, seller)
	_, err = w.p.Exec(w.ctx, `UPDATE market.listing SET expires_at = now() - interval '1 second' WHERE id = $1`, id)
	require.NoError(t, err)
	b.fail("POST", b.path("/market/bid/"+id), map[string]any{"amount_copper": 9000}, 409, "E_LISTING_UNAVAILABLE")
	k := market.Keeper{Pool: w.p, Txn: w.txn}
	_, err = market.ExpireAndRecover(context.Background(), w.p, 500)
	require.NoError(t, err)
	_, err = k.Run(context.Background(), 500)
	require.NoError(t, err)
	var state, owner string
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT l.state, i.owner_character_id::text FROM market.listing l
		JOIN inventory.item_instance i ON i.id = l.item_instance_id WHERE l.id = $1`, id).Scan(&state, &owner))
	require.Equal(t, "sold", state, "not expired")
	require.Equal(t, a.char, owner)
	require.Equal(t, beforeSeller+3000-market.SaleCut(3000), goldOf(t, seller))
	require.Equal(t, beforeA-3000, goldOf(t, a))

	// A buyout ends it at once.
	seller.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"character_id": seller.char, "template_code": "potion_hp_small", "quantity": 1})
	item = seller.bagItem("potion_hp_small")
	l = sell(map[string]any{"start_bid_copper": 100, "buyout_copper": 7000})
	id = l["id"].(string)
	a.ok("POST", a.path("/market/bid/"+id), map[string]any{"amount_copper": 500})
	beforeA = goldOf(t, a)
	v = b.ok("POST", b.path("/market/bid/"+id), map[string]any{"buyout": true})
	require.Equal(t, "sold", v["state"])
	require.Equal(t, beforeA+500, goldOf(t, a), "the outbid bidder got the gold back")

	// An auction without bids just expires and the item goes home.
	seller.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"character_id": seller.char, "template_code": "potion_hp_small", "quantity": 1})
	item = seller.bagItem("potion_hp_small")
	l = sell(map[string]any{"start_bid_copper": 100})
	id = l["id"].(string)
	_, err = w.p.Exec(w.ctx, `UPDATE market.listing SET expires_at = now() - interval '1 second' WHERE id = $1`, id)
	require.NoError(t, err)
	_, err = market.ExpireAndRecover(context.Background(), w.p, 500)
	require.NoError(t, err)
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT state FROM market.listing WHERE id = $1`, id).Scan(&state))
	require.Equal(t, "expired", state)

	// A bid stuck after its charge is paid back by the worker.
	seller.ok("POST", "/v1/vendor/general_ironhold/buy", map[string]any{"character_id": seller.char, "template_code": "potion_hp_small", "quantity": 1})
	item = seller.bagItem("potion_hp_small")
	l = sell(map[string]any{"start_bid_copper": 100})
	id = l["id"].(string)
	var stuck string
	require.NoError(t, w.p.QueryRow(w.ctx, `INSERT INTO market.auction_bid (listing_id, bidder_character_id, amount_copper, created_at)
		VALUES ($1, $2, 100, now() - interval '10 minutes') RETURNING id`, id, a.char).Scan(&stuck))
	_, err = k.Run(context.Background(), 500)
	require.NoError(t, err)
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT state FROM market.auction_bid WHERE id = $1`, stuck).Scan(&state))
	require.Equal(t, "void", state)
	var refunded bool
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT refunded_at IS NOT NULL FROM market.auction_bid WHERE id = $1`, stuck).Scan(&refunded))
	require.True(t, refunded, "never charged: nothing to pay back")
}
