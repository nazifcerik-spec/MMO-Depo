package market

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/notification"
)

// Auction rules (DD-359, docs/marketplace/auction.md).
const (
	MinBidStepPct = 5               // a new bid beats the top one by at least this
	AntiSnipe     = 5 * time.Minute // a bid this close to the end extends it to this
	MaxExtensions = 6               // ... at most this many times
	bidPendingTTL = 2 * time.Minute // a bid stuck between charge and count is undone after this
)

var (
	ErrAuctionOnly = errors.New("market: this listing is an auction")
	ErrBidTooLow   = errors.New("market: bid below the minimum")
	ErrTopBidder   = errors.New("market: you are already the top bidder")
	ErrAuctionBusy = errors.New("market: another bid is being placed, try again")
	ErrHasBids     = errors.New("market: an auction with bids cannot be cancelled")
	ErrNotAuction  = errors.New("market: this listing is not an auction")
	ErrBadStartBid = errors.New("market: start bid and buyout out of range")
)

// MinNextBid is the lowest acceptable bid.
func MinNextBid(start int64, current *int64) int64 {
	if current == nil {
		return start
	}
	step := *current * MinBidStepPct / 100
	if step < 1 {
		step = 1
	}
	return *current + step
}

// Keeper settles auctions and pays back bids; the service and the worker share it.
type Keeper struct {
	Pool   *pgxpool.Pool
	Txn    economy.TxnService
	Notify notification.Notifier
	Pub    EventPublisher
	Log    *slog.Logger
}

func (k Keeper) log() *slog.Logger {
	if k.Log != nil {
		return k.Log
	}
	return slog.Default()
}

func (s *Service) keeper() Keeper {
	return Keeper{Pool: s.pool, Txn: s.txn, Notify: s.notify, Pub: s.pub, Log: s.log}
}

// Bid places a bid, or the buyout when buyout is set. The gold moves to escrow
// first; the outbid bidder gets theirs back right after the new bid counts.
func (s *Service) Bid(ctx context.Context, actor, account, listingID uuid.UUID, amount int64, buyout bool) (*Listing, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var bidID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var kind, state string
		var sellerAccount uuid.UUID
		var expires time.Time
		var start, buy, current *int64
		var bidder *uuid.UUID
		err := tx.QueryRow(ctx, `
			SELECT l.kind, l.state, c.account_id, l.expires_at, l.start_bid_copper, l.buyout_copper,
			       l.current_bid_copper, l.bidder_character_id
			  FROM market.listing l JOIN character.character c ON c.id = l.seller_character_id
			 WHERE l.id = $1 FOR UPDATE OF l`, listingID).
			Scan(&kind, &state, &sellerAccount, &expires, &start, &buy, &current, &bidder)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		switch {
		case kind != "auction":
			return ErrNotAuction
		case state != "active" || !expires.After(s.now()):
			return ErrUnavailable
		case sellerAccount == account:
			return ErrSelfBuy
		case bidder != nil && *bidder == actor && !buyout:
			return ErrTopBidder
		}
		if buyout {
			if buy == nil {
				return ErrBidTooLow
			}
			amount = *buy
		} else if amount < MinNextBid(*start, current) {
			return ErrBidTooLow
		}
		var busy bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM market.auction_bid
			WHERE listing_id = $1 AND state = 'pending' AND created_at > $2)`, listingID, s.now().Add(-bidPendingTTL)).Scan(&busy); err != nil {
			return err
		}
		if busy {
			return ErrAuctionBusy
		}
		return tx.QueryRow(ctx, `INSERT INTO market.auction_bid (listing_id, bidder_character_id, amount_copper)
			VALUES ($1, $2, $3) RETURNING id`, listingID, actor, amount).Scan(&bidID)
	})
	if err != nil {
		return nil, err
	}

	if _, err := s.txn.Write(ctx, economy.Intent{
		Kind: "auction_bid", RefKind: "listing", RefID: listingID, IdempotencyKey: "auction-bid-" + bidID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: -amount},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemAuctionEscrow, Delta: amount},
		},
	}); err != nil {
		if _, derr := s.pool.Exec(ctx, `DELETE FROM market.auction_bid WHERE id = $1 AND state = 'pending'`, bidID); derr != nil {
			s.log.ErrorContext(ctx, "market: drop unpaid bid failed", "bid", bidID.String(), "err", derr.Error())
		}
		return nil, err
	}

	var previous *uuid.UUID
	var settleNow bool
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var state string
		var expires time.Time
		var qty int32
		var buy *int64
		var extensions int16
		if err := tx.QueryRow(ctx, `SELECT state, expires_at, stack_quantity, buyout_copper, extensions
			FROM market.listing WHERE id = $1 FOR UPDATE`, listingID).Scan(&state, &expires, &qty, &buy, &extensions); err != nil {
			return err
		}
		now := s.now()
		if state != "active" || !expires.After(now) {
			return ErrUnavailable
		}
		if err := tx.QueryRow(ctx, `UPDATE market.auction_bid SET state = 'outbid'
			WHERE listing_id = $1 AND state = 'active' RETURNING id`, listingID).Scan(&previous); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE market.auction_bid SET state = 'active' WHERE id = $1`, bidID); err != nil {
			return err
		}
		if expires.Sub(now) < AntiSnipe && extensions < MaxExtensions {
			expires = now.Add(AntiSnipe)
			extensions++
		}
		_, err := tx.Exec(ctx, `
			UPDATE market.listing
			   SET current_bid_copper = $2::bigint, bidder_character_id = $3, bid_count = bid_count + 1,
			       unit_price_copper = GREATEST(1, ($2::bigint + stack_quantity - 1) / stack_quantity), expires_at = $4, extensions = $5
			 WHERE id = $1`, listingID, amount, actor, expires, extensions)
		settleNow = buy != nil && amount >= *buy
		return err
	})
	if err != nil {
		// The listing closed under us: the bid never counted, the gold goes back.
		if _, verr := s.pool.Exec(ctx, `UPDATE market.auction_bid SET state = 'void' WHERE id = $1`, bidID); verr != nil {
			s.log.ErrorContext(ctx, "market: void bid failed", "bid", bidID.String(), "err", verr.Error())
		}
		k := s.keeper()
		k.refund(ctx, bidID, actor, amount)
		return nil, err
	}
	k := s.keeper()
	if previous != nil {
		k.refundBid(ctx, *previous, true)
	}
	if settleNow {
		if err := k.settle(ctx, listingID); err != nil {
			k.log().ErrorContext(ctx, "market: buyout settle failed, the worker retries", "listing", listingID.String(), "err", err.Error())
		}
	}
	return s.get(ctx, listingID)
}

// refund pays a bid's gold back once; the ledger key makes a repeat a no-op.
func (k Keeper) refund(ctx context.Context, bidID, bidder uuid.UUID, amount int64) bool {
	_, err := k.Txn.Write(ctx, economy.Intent{
		Kind: "auction_refund", RefKind: "auction_bid", RefID: bidID, IdempotencyKey: "auction-refund-" + bidID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemAuctionEscrow, Delta: -amount},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: bidder, Delta: amount},
		},
	})
	if err != nil {
		k.log().ErrorContext(ctx, "market: bid refund failed, the worker retries", "bid", bidID.String(), "err", err.Error())
		return false
	}
	if _, err := k.Pool.Exec(ctx, `UPDATE market.auction_bid SET refunded_at = now() WHERE id = $1`, bidID); err != nil {
		k.log().ErrorContext(ctx, "market: mark refunded failed", "bid", bidID.String(), "err", err.Error())
	}
	return true
}

func (k Keeper) refundBid(ctx context.Context, bidID uuid.UUID, tellOutbid bool) {
	var bidder *uuid.UUID
	var amount int64
	var account *uuid.UUID
	var code string
	var listing uuid.UUID
	err := k.Pool.QueryRow(ctx, `
		SELECT b.bidder_character_id, b.amount_copper, c.account_id, l.template_code, l.id
		  FROM market.auction_bid b JOIN market.listing l ON l.id = b.listing_id
		  LEFT JOIN character.character c ON c.id = b.bidder_character_id
		 WHERE b.id = $1`, bidID).Scan(&bidder, &amount, &account, &code, &listing)
	if err != nil || bidder == nil {
		return // deleted character: there is no one to pay
	}
	if !k.refund(ctx, bidID, *bidder, amount) || !tellOutbid || account == nil {
		return
	}
	if k.Pub != nil {
		k.Pub.Publish(*account, "market.outbid", map[string]any{"listing_id": listing, "template_code": code})
	}
	notification.Send(ctx, k.Notify, k.log(), notification.Notice{
		AccountID: *account, CharacterID: bidder, Kind: "market_outbid", DedupKey: "market_outbid:" + bidID.String(),
		TitleKey: "notif.market_outbid.title", BodyKey: "notif.market_outbid.body",
		Vars: map[string]any{"template_code": code}, Route: "market/listings", Category: "social",
	})
}

// settle closes an auction that has a top bid: the seller is paid from escrow
// and the winner gets the item. Reuses the fixed-price pending_sale steps, so a
// crash is recovered the same way.
func (k Keeper) settle(ctx context.Context, listingID uuid.UUID) error {
	var seller, bidder, sellerAccount uuid.UUID
	var amount int64
	err := pgx.BeginFunc(ctx, k.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			SELECT l.seller_character_id, l.bidder_character_id, l.current_bid_copper, c.account_id
			  FROM market.listing l JOIN character.character c ON c.id = l.seller_character_id
			 WHERE l.id = $1 AND l.kind = 'auction' AND l.state = 'active' AND l.bidder_character_id IS NOT NULL
			 FOR UPDATE OF l`, listingID).Scan(&seller, &bidder, &amount, &sellerAccount)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE market.listing SET state = 'pending_sale', buyer_character_id = bidder_character_id,
			pending_at = now() WHERE id = $1`, listingID)
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil // already settled, or the winner's character is gone
	}
	if err != nil {
		return err
	}
	cut := SaleCut(amount)
	if _, err := k.Txn.Write(ctx, economy.Intent{
		Kind: "auction_settle", RefKind: "listing", RefID: listingID, IdempotencyKey: "market-buy-" + listingID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemAuctionEscrow, Delta: -amount},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: seller, Delta: amount - cut},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: cut},
		},
	}); err != nil {
		return err // pending_sale stays; the worker reverts and retries
	}
	var code string
	err = pgx.BeginFunc(ctx, k.Pool, func(tx pgx.Tx) error {
		inst, err := finalizeSale(ctx, tx, listingID)
		if err != nil {
			return err
		}
		code = inst.TemplateCode
		_, err = tx.Exec(ctx, `UPDATE market.auction_bid SET state = 'won' WHERE listing_id = $1 AND state = 'active'`, listingID)
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	if k.Pub != nil {
		k.Pub.Publish(sellerAccount, "market.sold", map[string]any{"listing_id": listingID, "template_code": code, "net_copper": amount - cut})
	}
	notification.Send(ctx, k.Notify, k.log(), notification.Notice{
		AccountID: sellerAccount, CharacterID: &seller, Kind: "market_sold", DedupKey: "market_sold:" + listingID.String(),
		TitleKey: "notif.market_sold.title", BodyKey: "notif.market_sold.body",
		Vars: map[string]any{"template_code": code, "net_copper": amount - cut}, Route: "market/listings", Category: "social",
	})
	var winnerAccount uuid.UUID
	if k.Pool.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, bidder).Scan(&winnerAccount) == nil {
		notification.Send(ctx, k.Notify, k.log(), notification.Notice{
			AccountID: winnerAccount, CharacterID: &bidder, Kind: "market_won", DedupKey: "market_won:" + listingID.String(),
			TitleKey: "notif.market_won.title", BodyKey: "notif.market_won.body",
			Vars: map[string]any{"template_code": code}, Route: "market/listings", Category: "social",
		})
	}
	return nil
}

// Run is the worker's auction pass: settle what is due, undo bids stuck between
// charge and count, pay back whatever is still owed.
func (k Keeper) Run(ctx context.Context, batch int) (int64, error) {
	rows, err := k.Pool.Query(ctx, `SELECT id FROM market.listing
		WHERE state = 'active' AND kind = 'auction' AND bidder_character_id IS NOT NULL AND expires_at <= now()
		ORDER BY expires_at LIMIT $1`, batch)
	if err != nil {
		return 0, err
	}
	due, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return 0, err
	}
	var n int64
	for _, id := range due {
		if err := k.settle(ctx, id); err != nil {
			k.log().ErrorContext(ctx, "market: settle failed", "listing", id.String(), "err", err.Error())
			continue
		}
		n++
	}
	// A pending bid past its time never counted: gold back if it was charged.
	if _, err := k.Pool.Exec(ctx, `UPDATE market.auction_bid SET state = 'void'
		WHERE state = 'pending' AND created_at < $1`, time.Now().Add(-bidPendingTTL)); err != nil {
		return n, err
	}
	rows, err = k.Pool.Query(ctx, `
		SELECT b.id, b.state, EXISTS (SELECT 1 FROM economy.ledger_entry e WHERE e.idempotency_key = 'auction-bid-' || b.id::text)
		  FROM market.auction_bid b
		 WHERE b.state IN ('outbid', 'void') AND b.refunded_at IS NULL AND b.bidder_character_id IS NOT NULL
		 ORDER BY b.created_at LIMIT $1`, batch)
	if err != nil {
		return n, err
	}
	type owed struct {
		id      uuid.UUID
		state   string
		charged bool
	}
	var list []owed
	for rows.Next() {
		var o owed
		if err := rows.Scan(&o.id, &o.state, &o.charged); err != nil {
			rows.Close()
			return n, err
		}
		list = append(list, o)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return n, err
	}
	for _, o := range list {
		if !o.charged {
			if _, err := k.Pool.Exec(ctx, `UPDATE market.auction_bid SET refunded_at = now() WHERE id = $1`, o.id); err != nil {
				return n, err
			}
			continue
		}
		k.refundBid(ctx, o.id, false)
		n++
	}
	return n, nil
}
