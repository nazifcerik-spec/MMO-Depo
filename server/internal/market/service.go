package market

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/notification"
)

type Owner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

type CatalogSource interface {
	Catalog() *item.Catalog
}

type SlotAssigner interface {
	AssignSlots(ctx context.Context, ownerID uuid.UUID, instances []item.Instance) error
}

type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Service struct {
	pool    *pgxpool.Pool
	chars   Owner
	txn     economy.TxnService
	catalog CatalogSource
	slots   SlotAssigner
	pub     EventPublisher
	notify  notification.Notifier
	now     func() time.Time
	log     *slog.Logger
}

func NewService(pool *pgxpool.Pool, chars Owner, txn economy.TxnService, catalog CatalogSource, slots SlotAssigner, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, chars: chars, txn: txn, catalog: catalog, slots: slots, now: time.Now, log: logger}
}

func (s *Service) WithPublisher(p EventPublisher) *Service { s.pub = p; return s }

func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notify = n; return s }

// --- sell -------------------------------------------------------------------

type SellRequest struct {
	ItemID         uuid.UUID
	Quantity       int32
	UnitPrice      int64
	DurationHours  int16
	IdempotencyKey string
}

// Sell charges the listing fee, then escrows the item (or a split-off part of
// its stack) under a new listing. A failed escrow refunds the fee.
func (s *Service) Sell(ctx context.Context, actor, account uuid.UUID, req SellRequest) (*Listing, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	if req.IdempotencyKey == "" {
		return nil, ErrIdempotencyMissing
	}
	if _, ok := DurationMultiplier[req.DurationHours]; !ok {
		return nil, ErrBadDuration
	}
	if req.Quantity < 1 {
		return nil, ErrBadQuantity
	}
	ledgerKey := "market-list-" + req.IdempotencyKey
	if l, err := s.listingByFeeKey(ctx, ledgerKey); err == nil {
		return l, nil
	}

	var active int
	if err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM market.listing
		 WHERE seller_character_id = $1 AND state IN ('active','pending_sale')`, actor).Scan(&active); err != nil {
		return nil, err
	}
	if active >= ListingCap {
		return nil, ErrListingCap
	}

	var (
		tmplCode, rarity      string
		prefix, suffix, equip *string
		stack                 int32
		locked, bound         bool
	)
	err := s.pool.QueryRow(ctx, `
		SELECT template_code, rarity, prefix_code, suffix_code, stack_quantity, is_locked, is_bound, equipment_slot
		  FROM inventory.item_instance WHERE id = $1 AND owner_character_id = $2`, req.ItemID, actor).
		Scan(&tmplCode, &rarity, &prefix, &suffix, &stack, &locked, &bound, &equip)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrItemNotFound
	}
	if err != nil {
		return nil, err
	}
	tmpl, err := s.catalog.Catalog().Get(tmplCode)
	if err != nil {
		return nil, err
	}
	if locked || bound || equip != nil || tmpl.NoTrade {
		return nil, ErrNotTradable
	}
	if req.Quantity > stack || (!tmpl.Stackable && req.Quantity != stack) {
		return nil, ErrBadQuantity
	}
	lo, hi := PriceBounds(tmpl.BaseValueCopper)
	if req.UnitPrice < lo || req.UnitPrice > hi {
		return nil, fmt.Errorf("%w: [%d, %d]", ErrPriceOutOfRange, lo, hi)
	}

	fee := ListingFee(req.UnitPrice, req.Quantity, req.DurationHours)
	feeEntries := []economy.Entry{
		{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: -fee},
		{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: fee},
	}
	res, err := s.txn.Write(ctx, economy.Intent{
		Kind: "market_listing_fee", RefKind: "listing", IdempotencyKey: ledgerKey, Entries: feeEntries,
	})
	if err != nil {
		return nil, err
	}

	listingID := uuid.New()
	now := s.now().UTC()
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		escrowID := req.ItemID
		if req.Quantity == stack {
			tag, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET is_locked = true, lock_owner = 'market', locked_ref_id = $3, slot = NULL,
				       updated_at = now(), version = version + 1
				 WHERE id = $1 AND owner_character_id = $2
				   AND is_locked = false AND is_bound = false AND equipment_slot IS NULL`,
				req.ItemID, actor, listingID)
			if err != nil {
				return err
			}
			if tag.RowsAffected() != 1 {
				return ErrNotTradable
			}
		} else {
			tag, err := tx.Exec(ctx, `
				UPDATE inventory.item_instance
				   SET stack_quantity = stack_quantity - $3, updated_at = now(), version = version + 1
				 WHERE id = $1 AND owner_character_id = $2 AND is_locked = false AND stack_quantity > $3`,
				req.ItemID, actor, req.Quantity)
			if err != nil {
				return err
			}
			if tag.RowsAffected() != 1 {
				return ErrBadQuantity
			}
			if err := tx.QueryRow(ctx, `
				INSERT INTO inventory.item_instance
				  (template_code, owner_character_id, stack_quantity, rarity, prefix_code, suffix_code,
				   is_bound, is_locked, lock_owner, locked_ref_id, source_kind, source_ref_id)
				SELECT template_code, owner_character_id, $2, rarity, prefix_code, suffix_code,
				       false, true, 'market', $3, source_kind, source_ref_id
				  FROM inventory.item_instance WHERE id = $1
				RETURNING id`, req.ItemID, req.Quantity, listingID).Scan(&escrowID); err != nil {
				return err
			}
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO market.listing
			  (id, seller_character_id, item_instance_id, template_code, item_class, rarity, prefix_code, suffix_code,
			   stack_quantity, unit_price_copper, listing_fee_copper, fee_txn_id, duration_hours, created_at, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
			listingID, actor, escrowID, tmplCode, tmpl.ItemClass, rarity, prefix, suffix,
			req.Quantity, req.UnitPrice, fee, res.TxnID, req.DurationHours,
			now, now.Add(time.Duration(req.DurationHours)*time.Hour))
		return err
	})
	if err != nil {
		if l, lerr := s.listingByFeeKey(ctx, ledgerKey); lerr == nil {
			return l, nil
		}
		s.writeReverse(ctx, "market-list-refund-"+req.IdempotencyKey, feeEntries)
		return nil, err
	}
	return s.get(ctx, listingID)
}

func (s *Service) listingByFeeKey(ctx context.Context, key string) (*Listing, error) {
	var id uuid.UUID
	if err := s.pool.QueryRow(ctx, `
		SELECT id FROM market.listing
		 WHERE fee_txn_id IN (SELECT txn_id FROM economy.ledger_entry WHERE idempotency_key = $1)`, key).
		Scan(&id); err != nil {
		return nil, err
	}
	return s.get(ctx, id)
}

func (s *Service) writeReverse(ctx context.Context, key string, entries []economy.Entry) {
	rev := make([]economy.Entry, len(entries))
	for i, e := range entries {
		e.Delta = -e.Delta
		rev[i] = e
	}
	if _, err := s.txn.Write(ctx, economy.Intent{Kind: "market_listing_fee_refund", RefKind: "listing", IdempotencyKey: key, Entries: rev}); err != nil {
		s.log.ErrorContext(ctx, "market: fee refund failed — manual reconcile needed", "key", key, "err", err.Error())
	}
}

// --- buy --------------------------------------------------------------------

// Buy runs in three steps so the ledger write and the item hand-over cannot
// drift: reserve (pending_sale) → ledger 'market-buy-<id>' → finalise.
func (s *Service) Buy(ctx context.Context, actor, account, listingID uuid.UUID, expectedUnit *int64) (*Listing, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var seller uuid.UUID
	var total int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var state string
		var unit int64
		var buyer *uuid.UUID
		var expires time.Time
		var sellerAccount uuid.UUID
		err := tx.QueryRow(ctx, `
			SELECT l.seller_character_id, c.account_id, l.state, l.buyer_character_id,
			       l.unit_price_copper, l.total_price_copper, l.expires_at
			  FROM market.listing l JOIN character.character c ON c.id = l.seller_character_id
			 WHERE l.id = $1 FOR UPDATE OF l`, listingID).
			Scan(&seller, &sellerAccount, &state, &buyer, &unit, &total, &expires)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if state == "pending_sale" && buyer != nil && *buyer == actor {
			return nil // retry of our own reservation
		}
		if state != "active" || !expires.After(s.now()) {
			return ErrUnavailable
		}
		if seller == actor || sellerAccount == account {
			return ErrSelfBuy
		}
		if expectedUnit != nil && *expectedUnit != unit {
			return ErrPriceMismatch
		}
		_, err = tx.Exec(ctx, `
			UPDATE market.listing SET state = 'pending_sale', buyer_character_id = $2, pending_at = now()
			 WHERE id = $1`, listingID, actor)
		return err
	})
	if err != nil {
		return nil, err
	}

	cut := SaleCut(total)
	if _, err := s.txn.Write(ctx, economy.Intent{
		Kind: "market_buy", RefKind: "listing", RefID: listingID,
		IdempotencyKey: "market-buy-" + listingID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: -total},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: seller, Delta: total - cut},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: cut},
		},
	}); err != nil {
		if _, rerr := s.pool.Exec(ctx, `
			UPDATE market.listing SET state = 'active', buyer_character_id = NULL, pending_at = NULL
			 WHERE id = $1 AND state = 'pending_sale' AND buyer_character_id = $2`, listingID, actor); rerr != nil {
			s.log.ErrorContext(ctx, "market: revert pending sale failed", "listing", listingID.String(), "err", rerr.Error())
		}
		return nil, err
	}

	var inst item.Instance
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var ferr error
		inst, ferr = finalizeSale(ctx, tx, listingID)
		return ferr
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// The worker finished it first; fine as long as we are the buyer.
		l, gerr := s.get(ctx, listingID)
		if gerr == nil && l.State == "sold" {
			return l, nil
		}
	}
	if err != nil {
		return nil, err
	}
	if s.slots != nil {
		if err := s.slots.AssignSlots(ctx, actor, []item.Instance{inst}); err != nil {
			s.log.WarnContext(ctx, "market: slot assign after buy failed", "err", err.Error())
		}
	}
	var sellerAccount uuid.UUID
	if s.pool.QueryRow(ctx, `SELECT account_id FROM character.character WHERE id = $1`, seller).Scan(&sellerAccount) == nil {
		if s.pub != nil {
			s.pub.Publish(sellerAccount, "market.sold", map[string]any{
				"listing_id": listingID, "template_code": inst.TemplateCode, "net_copper": total - cut,
			})
		}
		notification.Send(ctx, s.notify, s.log, notification.Notice{
			AccountID: sellerAccount, CharacterID: &seller, Kind: "market_sold",
			DedupKey: "market_sold:" + listingID.String(),
			TitleKey: "notif.market_sold.title", BodyKey: "notif.market_sold.body",
			Vars:  map[string]any{"template_code": inst.TemplateCode, "net_copper": total - cut},
			Route: "market/listings", Category: "social",
		})
	}
	return s.get(ctx, listingID)
}

// finalizeSale hands the escrowed item to the reserved buyer and closes the
// listing. Returns pgx.ErrNoRows when the listing is not pending_sale.
func finalizeSale(ctx context.Context, tx pgx.Tx, listingID uuid.UUID) (item.Instance, error) {
	var inst item.Instance
	var buyer uuid.UUID
	var unit int64
	var qty int32
	if err := tx.QueryRow(ctx, `
		UPDATE market.listing SET state = 'sold', closed_at = now(), close_reason = 'sold'
		 WHERE id = $1 AND state = 'pending_sale'
		RETURNING buyer_character_id, item_instance_id, template_code, unit_price_copper, stack_quantity`,
		listingID).Scan(&buyer, &inst.ID, &inst.TemplateCode, &unit, &qty); err != nil {
		return inst, err
	}
	inst.StackQuantity = qty
	if _, err := tx.Exec(ctx, `
		UPDATE inventory.item_instance
		   SET owner_character_id = $2, is_locked = false, lock_owner = NULL, locked_ref_id = NULL, slot = NULL,
		       updated_at = now(), version = version + 1
		 WHERE id = $1`, inst.ID, buyer); err != nil {
		return inst, err
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO market.price_history (template_code, unit_price_copper, quantity) VALUES ($1, $2, $3)`,
		inst.TemplateCode, unit, qty)
	return inst, err
}

// --- cancel -----------------------------------------------------------------

// Cancel returns the item to the seller and refunds part of the listing fee
// depending on the listing age.
func (s *Service) Cancel(ctx context.Context, actor, account, listingID uuid.UUID) (*Listing, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var inst item.Instance
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var seller uuid.UUID
		var state string
		var fee int64
		var created time.Time
		err := tx.QueryRow(ctx, `
			SELECT seller_character_id, state, listing_fee_copper, created_at, item_instance_id, template_code, stack_quantity
			  FROM market.listing WHERE id = $1 FOR UPDATE`, listingID).
			Scan(&seller, &state, &fee, &created, &inst.ID, &inst.TemplateCode, &inst.StackQuantity)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if seller != actor {
			return ErrNotSeller
		}
		if state != "active" {
			return ErrUnavailable
		}
		if refund := CancelRefund(fee, s.now().Sub(created)); refund > 0 {
			if _, err := s.txn.Write(ctx, economy.Intent{
				Kind: "market_listing_fee_refund", RefKind: "listing", RefID: listingID,
				IdempotencyKey: "market-cancel-" + listingID.String(),
				Entries: []economy.Entry{
					{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: -refund},
					{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: refund},
				},
			}); err != nil {
				return err
			}
		}
		if err := releaseEscrow(ctx, tx, inst.ID, listingID); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE market.listing SET state = 'cancelled', closed_at = now(), close_reason = 'cancelled'
			 WHERE id = $1`, listingID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if s.slots != nil {
		if err := s.slots.AssignSlots(ctx, actor, []item.Instance{inst}); err != nil {
			s.log.WarnContext(ctx, "market: slot assign after cancel failed", "err", err.Error())
		}
	}
	return s.get(ctx, listingID)
}

// AdminCancel closes an active listing for moderation: the item goes back to
// the seller's overflow and the fee is refunded in full only when asked
// (docs/marketplace/cancel.md — admin cancel, default no refund).
func (s *Service) AdminCancel(ctx context.Context, listingID uuid.UUID, refundFee bool) (*Listing, error) {
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var seller, itemID uuid.UUID
		var state string
		var fee int64
		err := tx.QueryRow(ctx, `
			SELECT seller_character_id, state, listing_fee_copper, item_instance_id
			  FROM market.listing WHERE id = $1 FOR UPDATE`, listingID).Scan(&seller, &state, &fee, &itemID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if state != "active" {
			return ErrUnavailable
		}
		if refundFee && fee > 0 {
			if _, err := s.txn.Write(ctx, economy.Intent{
				Kind: "market_listing_fee_refund", RefKind: "listing", RefID: listingID,
				IdempotencyKey: "market-admin-cancel-" + listingID.String(),
				Entries: []economy.Entry{
					{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: -fee},
					{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: seller, Delta: fee},
				},
			}); err != nil {
				return err
			}
		}
		if err := releaseEscrow(ctx, tx, itemID, listingID); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE market.listing SET state = 'cancelled', closed_at = now(), close_reason = 'admin'
			 WHERE id = $1`, listingID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.get(ctx, listingID)
}

func releaseEscrow(ctx context.Context, tx pgx.Tx, itemID, listingID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		UPDATE inventory.item_instance
		   SET is_locked = false, lock_owner = NULL, locked_ref_id = NULL, updated_at = now(), version = version + 1
		 WHERE id = $1 AND locked_ref_id = $2 AND lock_owner = 'market'`, itemID, listingID)
	return err
}

// --- worker -----------------------------------------------------------------

// ExpireAndRecover is the worker entry point: it expires overdue listings
// (item back to the seller as overflow, no fee refund) and resolves
// pending sales abandoned by a crash — finished when the buy ledger txn
// exists, otherwise reverted to active.
func ExpireAndRecover(ctx context.Context, pool *pgxpool.Pool, batch int) (int64, error) {
	var n int64
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			WITH exp AS (
			  UPDATE market.listing SET state = 'expired', closed_at = now(), close_reason = 'expired'
			   WHERE id IN (SELECT id FROM market.listing WHERE state = 'active' AND expires_at <= now()
			                 ORDER BY expires_at LIMIT $1 FOR UPDATE SKIP LOCKED)
			  RETURNING id, item_instance_id)
			UPDATE inventory.item_instance i
			   SET is_locked = false, lock_owner = NULL, locked_ref_id = NULL, updated_at = now(), version = i.version + 1
			  FROM exp WHERE i.id = exp.item_instance_id AND i.locked_ref_id = exp.id`, batch)
		n = tag.RowsAffected()
		return err
	})
	if err != nil {
		return n, err
	}

	rows, err := pool.Query(ctx, `
		SELECT id FROM market.listing WHERE state = 'pending_sale' AND pending_at < $1 LIMIT $2`,
		time.Now().Add(-PendingSaleTTL), batch)
	if err != nil {
		return n, err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return n, err
	}
	for _, id := range ids {
		err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
			var paid bool
			if err := tx.QueryRow(ctx, `
				SELECT EXISTS (SELECT 1 FROM economy.ledger_entry WHERE idempotency_key = $1)`,
				"market-buy-"+id.String()).Scan(&paid); err != nil {
				return err
			}
			if paid {
				_, err := finalizeSale(ctx, tx, id)
				if errors.Is(err, pgx.ErrNoRows) {
					return nil
				}
				return err
			}
			_, err := tx.Exec(ctx, `
				UPDATE market.listing SET state = 'active', buyer_character_id = NULL, pending_at = NULL
				 WHERE id = $1 AND state = 'pending_sale'`, id)
			return err
		})
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// --- reads ------------------------------------------------------------------

const listingCols = `
	SELECT l.id, l.seller_character_id, c.name, l.item_instance_id, l.template_code, l.item_class, l.rarity,
	       l.prefix_code, l.suffix_code, l.stack_quantity, l.unit_price_copper, l.total_price_copper,
	       l.listing_fee_copper, l.duration_hours, l.state, l.created_at, l.expires_at, l.closed_at
	  FROM market.listing l JOIN character.character c ON c.id = l.seller_character_id `

func scanListing(r pgx.CollectableRow) (Listing, error) {
	var l Listing
	err := r.Scan(&l.ID, &l.SellerID, &l.SellerName, &l.ItemInstanceID, &l.TemplateCode, &l.ItemClass, &l.Rarity,
		&l.PrefixCode, &l.SuffixCode, &l.StackQuantity, &l.UnitPriceCopper, &l.TotalPriceCopper,
		&l.ListingFeeCopper, &l.DurationHours, &l.State, &l.CreatedAt, &l.ExpiresAt, &l.ClosedAt)
	return l, err
}

func (s *Service) get(ctx context.Context, id uuid.UUID) (*Listing, error) {
	rows, err := s.pool.Query(ctx, listingCols+`WHERE l.id = $1`, id)
	if err != nil {
		return nil, err
	}
	l, err := pgx.CollectOneRow(rows, scanListing)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &l, err
}

func (s *Service) MyListings(ctx context.Context, actor, account uuid.UUID) ([]Listing, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, listingCols+`
		WHERE l.seller_character_id = $1 AND (l.state IN ('active','pending_sale') OR l.closed_at > now() - INTERVAL '7 days')
		ORDER BY l.created_at DESC LIMIT 100`, actor)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanListing)
}

type SearchQuery struct {
	TemplateCode string
	ItemClass    string
	Rarity       string
	MinPrice     int64
	MaxPrice     int64
	Sort         string // price_asc (default) | price_desc | newest
	Limit        int
	Offset       int
}

func (s *Service) Search(ctx context.Context, q SearchQuery) ([]Listing, error) {
	where := []string{"l.state = 'active'", "l.expires_at > now()"}
	args := []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if q.TemplateCode != "" {
		add("l.template_code = $%d", q.TemplateCode)
	}
	if q.ItemClass != "" {
		add("l.item_class = $%d", q.ItemClass)
	}
	if q.Rarity != "" {
		add("l.rarity = $%d", q.Rarity)
	}
	if q.MinPrice > 0 {
		add("l.unit_price_copper >= $%d", q.MinPrice)
	}
	if q.MaxPrice > 0 {
		add("l.unit_price_copper <= $%d", q.MaxPrice)
	}
	order := "l.unit_price_copper ASC, l.created_at"
	switch q.Sort {
	case "price_desc":
		order = "l.unit_price_copper DESC, l.created_at"
	case "newest":
		order = "l.created_at DESC"
	}
	if q.Limit <= 0 || q.Limit > maxSearchSize {
		q.Limit = defaultSearchSize
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	args = append(args, q.Limit, q.Offset)
	sql := listingCols + "WHERE " + strings.Join(where, " AND ") +
		fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", order, len(args)-1, len(args))
	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanListing)
}

// History returns 30-day sale stats plus the cheapest active listing.
func (s *Service) History(ctx context.Context, templateCode string) (*PriceStats, error) {
	st := &PriceStats{TemplateCode: templateCode}
	if err := s.pool.QueryRow(ctx, `
		SELECT count(*), COALESCE(sum(quantity), 0),
		       COALESCE(round(sum(unit_price_copper * quantity)::numeric / NULLIF(sum(quantity), 0)), 0)::bigint,
		       COALESCE(min(unit_price_copper), 0), COALESCE(max(unit_price_copper), 0)
		  FROM market.price_history
		 WHERE template_code = $1 AND sold_at > now() - INTERVAL '30 days'`, templateCode).
		Scan(&st.Sales, &st.Quantity, &st.AvgUnit, &st.MinUnit, &st.MaxUnit); err != nil {
		return nil, err
	}
	if err := s.pool.QueryRow(ctx, `
		SELECT min(unit_price_copper) FROM market.listing
		 WHERE template_code = $1 AND state = 'active' AND expires_at > now()`, templateCode).
		Scan(&st.LowestActive); err != nil {
		return nil, err
	}
	return st, nil
}
