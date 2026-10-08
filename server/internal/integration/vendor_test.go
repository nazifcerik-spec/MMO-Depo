//go:build integration

package integration

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/vendor"
)

func TestVendorRules(t *testing.T) {
	w := newWorld(t)
	svc, err := vendor.NewService(w.ctx, vendor.NewPgRepo(w.p), w.items, w.inv, w.txn, vendor.DefaultServiceConfig(), nil)
	require.NoError(t, err)
	a, b := w.player("warrior", 1), w.player("mage", 1)
	const shop = "general_ironhold"

	buy := func(p player, code string, price *int64) (*vendor.BuyResult, error) {
		return svc.Buy(w.ctx, vendor.BuyRequest{VendorCode: shop, CharacterID: p.ID, TemplateCode: code,
			ExpectedPriceCopper: price, IdempotencyKey: uuid.NewString()})
	}
	_, err = buy(a, "potion_hp_small", nil)
	require.ErrorIs(t, err, economy.ErrInsufficientBalance)
	w.gold(a, 50_000)
	_, err = svc.Buy(w.ctx, vendor.BuyRequest{VendorCode: "nope", CharacterID: a.ID, TemplateCode: "potion_hp_small", IdempotencyKey: uuid.NewString()})
	require.ErrorIs(t, err, vendor.ErrVendorNotFound)
	_, err = buy(a, "iron_ore", nil)
	require.ErrorIs(t, err, vendor.ErrItemNotSold)
	_, err = svc.Buy(w.ctx, vendor.BuyRequest{VendorCode: shop, CharacterID: a.ID, TemplateCode: "potion_hp_small"})
	require.Error(t, err, "idempotency key is mandatory")
	key := uuid.NewString()
	first, err := svc.Buy(w.ctx, vendor.BuyRequest{VendorCode: shop, CharacterID: a.ID, TemplateCode: "sword_training", IdempotencyKey: key})
	require.NoError(t, err)
	balance := w.balance(a)
	again, err := svc.Buy(w.ctx, vendor.BuyRequest{VendorCode: shop, CharacterID: a.ID, TemplateCode: "sword_training", IdempotencyKey: key})
	require.NoError(t, err)
	require.Equal(t, first.InstanceID, again.InstanceID, "replay mints nothing new")
	require.Equal(t, balance, w.balance(a), "and charges nothing")

	// Rarity raises the vendor sell price (single table: item.ValueMultiplier).
	common := w.give(a, "sword_iron", 1)
	rare := w.give(a, "sword_iron", 1)
	_, err = w.p.Exec(w.ctx, `UPDATE inventory.item_instance SET rarity = 'rare' WHERE id = $1`, rare.ID)
	require.NoError(t, err)
	sell := func(p player, id uuid.UUID) (*vendor.SellResult, error) {
		return svc.Sell(w.ctx, vendor.SellRequest{VendorCode: shop, CharacterID: p.ID, ItemInstanceID: id, IdempotencyKey: uuid.NewString()})
	}
	c, err := sell(a, common.ID)
	require.NoError(t, err)
	r, err := sell(a, rare.ID)
	require.NoError(t, err)
	require.InDelta(t, float64(c.Received)*item.ValueMultiplier(item.RarityRare), float64(r.Received), 1)

	_, err = sell(b, first.InstanceID)
	require.ErrorIs(t, err, vendor.ErrItemNotOwned)
	_, err = w.p.Exec(w.ctx, `UPDATE inventory.item_instance SET is_bound = true WHERE id = $1`, first.InstanceID)
	require.NoError(t, err)
	_, err = sell(a, first.InstanceID)
	require.ErrorIs(t, err, vendor.ErrItemBound)

	// Buyback: own, once, within the TTL.
	list, err := svc.ListBuybacks(w.ctx, a.ID)
	require.NoError(t, err)
	require.Len(t, list, 2)
	_, err = svc.Buyback(w.ctx, b.ID, c.BuybackID, uuid.NewString())
	require.Error(t, err, "someone else's buyback")
	_, err = svc.Buyback(w.ctx, a.ID, c.BuybackID, uuid.NewString())
	require.NoError(t, err)
	_, err = svc.Buyback(w.ctx, a.ID, c.BuybackID, uuid.NewString())
	require.ErrorIs(t, err, vendor.ErrBuybackClaimed)
	_, err = w.p.Exec(w.ctx, `UPDATE economy.vendor_buyback SET expires_at = now() - interval '1 minute' WHERE id = $1`, r.BuybackID)
	require.NoError(t, err)
	_, err = svc.Buyback(w.ctx, a.ID, r.BuybackID, uuid.NewString())
	require.ErrorIs(t, err, vendor.ErrBuybackExpired)
	_, err = svc.Buyback(w.ctx, a.ID, uuid.New(), uuid.NewString())
	require.ErrorIs(t, err, vendor.ErrBuybackNotFound)
}

// Concurrent twins must never mint a free copy or charge twice (migration
// 0028 + per-entry buyback ledger key).
func TestVendorConcurrentMintsOnce(t *testing.T) {
	w := newWorld(t)
	svc, err := vendor.NewService(w.ctx, vendor.NewPgRepo(w.p), w.items, w.inv, w.txn, vendor.DefaultServiceConfig(), nil)
	require.NoError(t, err)
	p := w.player("warrior", 1)
	w.gold(p, 100_000)
	count := func(source uuid.UUID) int {
		var n int
		require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM inventory.item_instance WHERE source_ref_id = $1`, source).Scan(&n))
		return n
	}
	parallel := func(n int, fn func() error) []error {
		errs := make([]error, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) { defer wg.Done(); errs[i] = fn() }(i)
		}
		wg.Wait()
		return errs
	}

	key := uuid.NewString()
	before := w.balance(p)
	var txn uuid.UUID
	var mu sync.Mutex
	parallel(8, func() error {
		r, err := svc.Buy(w.ctx, vendor.BuyRequest{VendorCode: "general_ironhold", CharacterID: p.ID,
			TemplateCode: "sword_training", IdempotencyKey: key})
		if err == nil {
			mu.Lock()
			txn = r.LedgerTxnID
			mu.Unlock()
		}
		return err
	})
	require.NotEqual(t, uuid.Nil, txn)
	require.Equal(t, 1, count(txn), "one item for one purchase")
	price := before - w.balance(p)
	require.Positive(t, price)

	sold, err := svc.Sell(w.ctx, vendor.SellRequest{VendorCode: "general_ironhold", CharacterID: p.ID,
		ItemInstanceID: w.give(p, "sword_iron", 1).ID, IdempotencyKey: uuid.NewString()})
	require.NoError(t, err)
	before = w.balance(p)
	errs := parallel(8, func() error {
		_, err := svc.Buyback(w.ctx, p.ID, sold.BuybackID, uuid.NewString()) // different keys
		return err
	})
	ok := 0
	for _, e := range errs {
		if e == nil {
			ok++
		}
	}
	require.Positive(t, ok)
	require.Equal(t, 1, count(sold.BuybackID), "one item for one buyback")
	require.Equal(t, before-sold.Received, w.balance(p), "charged exactly once")
}
