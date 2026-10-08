//go:build integration

package integration

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/market"
)

func newMarket(w *testEnv, ev *events) *market.Service {
	return market.NewService(w.p, w.chars, w.txn, w.items, w.inv, nil).WithPublisher(ev)
}

func sell(t *testing.T, svc *market.Service, w *testEnv, p player, itemID uuid.UUID, qty int32, unit int64) *market.Listing {
	t.Helper()
	l, err := svc.Sell(w.ctx, p.ID, p.Account, market.SellRequest{
		ItemID: itemID, Quantity: qty, UnitPrice: unit, DurationHours: 12, IdempotencyKey: uuid.NewString(),
	})
	require.NoError(t, err)
	return l
}

func TestMarketSellValidation(t *testing.T) {
	w := newWorld(t)
	svc := newMarket(w, &events{})
	a := w.player("warrior", 1)
	w.gold(a, 100_000)
	potions := w.give(a, "potion_hp_small", 10)
	sword := w.give(a, "sword_iron", 1)

	req := func(id uuid.UUID, qty int32, unit int64, hours int16) market.SellRequest {
		return market.SellRequest{ItemID: id, Quantity: qty, UnitPrice: unit, DurationHours: hours, IdempotencyKey: uuid.NewString()}
	}
	cases := map[string]struct {
		r   market.SellRequest
		err error
	}{
		"duration":      {req(potions.ID, 1, 100, 6), market.ErrBadDuration},
		"zero qty":      {req(potions.ID, 0, 100, 12), market.ErrBadQuantity},
		"too many":      {req(potions.ID, 11, 100, 12), market.ErrBadQuantity},
		"partial sword": {req(sword.ID, 2, 1000, 12), market.ErrBadQuantity},
		"too cheap":     {req(potions.ID, 1, 1, 12), market.ErrPriceOutOfRange},
		"not mine":      {req(uuid.New(), 1, 100, 12), market.ErrItemNotFound},
	}
	for name, c := range cases {
		_, err := svc.Sell(w.ctx, a.ID, a.Account, c.r)
		require.ErrorIs(t, err, c.err, name)
	}
	_, err := svc.Sell(w.ctx, a.ID, a.Account, market.SellRequest{ItemID: potions.ID, Quantity: 1, UnitPrice: 100, DurationHours: 12})
	require.ErrorIs(t, err, market.ErrIdempotencyMissing)

	_, err = w.p.Exec(w.ctx, `UPDATE inventory.item_instance SET is_bound = true WHERE id = $1`, sword.ID)
	require.NoError(t, err)
	_, err = svc.Sell(w.ctx, a.ID, a.Account, req(sword.ID, 1, 1000, 12))
	require.ErrorIs(t, err, market.ErrNotTradable)
	require.Equal(t, int64(100_000), w.balance(a), "rejected listings cost nothing")
}

func TestMarketSellBuyFlow(t *testing.T) {
	w := newWorld(t)
	ev := &events{}
	svc := newMarket(w, ev)
	seller, buyer := w.player("warrior", 1), w.player("mage", 1)
	w.gold(seller, 10_000)
	w.gold(buyer, 50_000)
	potions := w.give(seller, "potion_hp_small", 10)

	// List 4 of 10: the listing escrows a split-off stack; 6 stay in the bag.
	key := uuid.NewString()
	l, err := svc.Sell(w.ctx, seller.ID, seller.Account, market.SellRequest{
		ItemID: potions.ID, Quantity: 4, UnitPrice: 500, DurationHours: 24, IdempotencyKey: key,
	})
	require.NoError(t, err)
	fee := market.ListingFee(500, 4, 24)
	require.Equal(t, fee, l.ListingFeeCopper)
	require.Equal(t, int64(10_000)-fee, w.balance(seller))
	require.Equal(t, int64(2_000), l.TotalPriceCopper)
	replay, err := svc.Sell(w.ctx, seller.ID, seller.Account, market.SellRequest{
		ItemID: potions.ID, Quantity: 4, UnitPrice: 500, DurationHours: 24, IdempotencyKey: key,
	})
	require.NoError(t, err)
	require.Equal(t, l.ID, replay.ID, "replay returns the same listing")
	require.Equal(t, int64(10_000)-fee, w.balance(seller), "and charges once")
	var left int32
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT stack_quantity FROM inventory.item_instance WHERE id = $1`, potions.ID).Scan(&left))
	require.Equal(t, int32(6), left)

	found, err := svc.Search(w.ctx, market.SearchQuery{TemplateCode: "potion_hp_small", MinPrice: 500, MaxPrice: 500, Sort: "newest"})
	require.NoError(t, err)
	require.Contains(t, listingIDs(found), l.ID)
	mine, err := svc.MyListings(w.ctx, seller.ID, seller.Account)
	require.NoError(t, err)
	require.Equal(t, l.ID, mine[0].ID)

	_, err = svc.Buy(w.ctx, seller.ID, seller.Account, l.ID, nil)
	require.ErrorIs(t, err, market.ErrSelfBuy)
	wrong := int64(499)
	_, err = svc.Buy(w.ctx, buyer.ID, buyer.Account, l.ID, &wrong)
	require.ErrorIs(t, err, market.ErrPriceMismatch)

	right := int64(500)
	bought, err := svc.Buy(w.ctx, buyer.ID, buyer.Account, l.ID, &right)
	require.NoError(t, err)
	require.Equal(t, "sold", bought.State)
	cut := market.SaleCut(2_000)
	require.Equal(t, int64(48_000), w.balance(buyer))
	require.Equal(t, int64(10_000)-fee+2_000-cut, w.balance(seller), "seller nets total minus the 5% cut")
	owner, locked := w.owner(*l.ItemInstanceID)
	require.Equal(t, buyer.ID, owner)
	require.False(t, locked)
	require.True(t, w.slotted(*l.ItemInstanceID), "bought item is visible in the bag")
	require.True(t, ev.has("market.sold"))

	_, err = svc.Buy(w.ctx, buyer.ID, buyer.Account, l.ID, &right)
	require.ErrorIs(t, err, market.ErrUnavailable, "sold listings cannot be bought again")
	stats, err := svc.History(w.ctx, "potion_hp_small")
	require.NoError(t, err)
	require.GreaterOrEqual(t, stats.Sales, int64(1))
}

func listingIDs(ls []market.Listing) []uuid.UUID {
	out := make([]uuid.UUID, len(ls))
	for i, l := range ls {
		out[i] = l.ID
	}
	return out
}

func TestMarketCancelExpireAndAdmin(t *testing.T) {
	w := newWorld(t)
	svc := newMarket(w, &events{})
	a, b := w.player("warrior", 1), w.player("mage", 1)
	w.gold(a, 100_000)

	// Cancel within the first hour refunds half the fee and returns the item.
	sword := w.give(a, "sword_iron", 1)
	l := sell(t, svc, w, a, sword.ID, 1, 10_000)
	before := w.balance(a)
	_, err := svc.Cancel(w.ctx, b.ID, b.Account, l.ID)
	require.ErrorIs(t, err, market.ErrNotSeller)
	cancelled, err := svc.Cancel(w.ctx, a.ID, a.Account, l.ID)
	require.NoError(t, err)
	require.Equal(t, "cancelled", cancelled.State)
	require.Equal(t, before+market.CancelRefund(l.ListingFeeCopper, 0), w.balance(a))
	owner, locked := w.owner(sword.ID)
	require.Equal(t, a.ID, owner)
	require.False(t, locked)
	require.True(t, w.slotted(sword.ID))
	_, err = svc.Cancel(w.ctx, a.ID, a.Account, l.ID)
	require.ErrorIs(t, err, market.ErrUnavailable)

	// Expiry returns the item unslotted; the next bag view slots it again.
	l2 := sell(t, svc, w, a, sword.ID, 1, 10_000)
	_, err = w.p.Exec(w.ctx, `UPDATE market.listing SET expires_at = now() - interval '1 minute' WHERE id = $1`, l2.ID)
	require.NoError(t, err)
	_, err = market.ExpireAndRecover(w.ctx, w.p, 500)
	require.NoError(t, err)
	_, locked = w.owner(sword.ID)
	require.False(t, locked, "expired listing releases the item")
	view, err := w.inv.View(w.ctx, a.ID)
	require.NoError(t, err)
	require.Contains(t, bagIDs(view.Items), sword.ID, "and it shows up in the bag again")

	// A pending sale abandoned before payment reverts to active.
	l3 := sell(t, svc, w, a, sword.ID, 1, 10_000)
	_, err = w.p.Exec(w.ctx, `
		UPDATE market.listing SET state = 'pending_sale', buyer_character_id = $2, pending_at = now() - interval '1 hour'
		 WHERE id = $1`, l3.ID, b.ID)
	require.NoError(t, err)
	_, err = market.ExpireAndRecover(w.ctx, w.p, 500)
	require.NoError(t, err)
	var state string
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT state FROM market.listing WHERE id = $1`, l3.ID).Scan(&state))
	require.Equal(t, "active", state)

	// A pending sale whose payment landed is finished for the buyer.
	w.gold(b, 50_000)
	_, err = w.txn.Write(w.ctx, economy.Intent{
		Kind: "market_buy", RefKind: "listing", RefID: l3.ID, IdempotencyKey: "market-buy-" + l3.ID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: b.ID, Delta: -10_000},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: a.ID, Delta: 9_500},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: 500},
		},
	})
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `
		UPDATE market.listing SET state = 'pending_sale', buyer_character_id = $2, pending_at = now() - interval '1 hour'
		 WHERE id = $1`, l3.ID, b.ID)
	require.NoError(t, err)
	_, err = market.ExpireAndRecover(w.ctx, w.p, 500)
	require.NoError(t, err)
	owner, _ = w.owner(sword.ID)
	require.Equal(t, b.ID, owner, "paid pending sale completes")

	// Admin cancel with full refund.
	shield := w.give(a, "sword_iron", 1)
	l4 := sell(t, svc, w, a, shield.ID, 1, 10_000)
	before = w.balance(a)
	got, err := svc.AdminCancel(w.ctx, l4.ID, true)
	require.NoError(t, err)
	require.Equal(t, "cancelled", got.State)
	require.Equal(t, before+l4.ListingFeeCopper, w.balance(a))
	_, err = svc.AdminCancel(w.ctx, l4.ID, true)
	require.ErrorIs(t, err, market.ErrUnavailable)
	_, err = svc.AdminCancel(w.ctx, uuid.New(), false)
	require.ErrorIs(t, err, market.ErrNotFound)
}

// Racing listings of the same item: no duplicated units, and every listing
// that lost the race gets its fee back.
func TestMarketConcurrentSellSameItem(t *testing.T) {
	w := newWorld(t)
	svc := newMarket(w, &events{})
	a := w.player("warrior", 1)
	w.gold(a, 100_000)
	sword := w.give(a, "sword_iron", 1)
	potions := w.give(a, "potion_hp_small", 10)

	race := func(itemID uuid.UUID, qty int32, unit int64) int {
		const n = 8
		var wg sync.WaitGroup
		var ok atomic.Int32
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if _, err := svc.Sell(w.ctx, a.ID, a.Account, market.SellRequest{ItemID: itemID, Quantity: qty,
					UnitPrice: unit, DurationHours: 12, IdempotencyKey: uuid.NewString()}); err == nil {
					ok.Add(1)
				}
			}()
		}
		wg.Wait()
		return int(ok.Load())
	}

	require.Equal(t, 1, race(sword.ID, 1, 1000), "a single item lists once")
	won := race(potions.ID, 4, 100)
	require.LessOrEqual(t, won, 2, "10 potions hold at most two lots of 4")
	require.GreaterOrEqual(t, won, 1)

	var units int32
	require.NoError(t, w.p.QueryRow(w.ctx, `
		SELECT coalesce(sum(stack_quantity), 0) FROM inventory.item_instance
		 WHERE owner_character_id = $1 AND template_code = 'potion_hp_small'`, a.ID).Scan(&units))
	require.Equal(t, int32(10), units, "no units created or lost")
	fees := market.ListingFee(1000, 1, 12) + int64(won)*market.ListingFee(100, 4, 12)
	require.Equal(t, int64(100_000)-fees, w.balance(a), "losers are refunded")

	// Leave no cheap listings behind for other tests' searches.
	mine, err := svc.MyListings(w.ctx, a.ID, a.Account)
	require.NoError(t, err)
	for _, l := range mine {
		_, err := svc.Cancel(w.ctx, a.ID, a.Account, l.ID)
		require.NoError(t, err)
	}
}
