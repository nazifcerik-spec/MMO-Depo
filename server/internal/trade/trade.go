// Package trade is player-to-player item trading (DD-315). Two characters
// open a trade, each puts items up, and when both are ready the items
// change hands in one go. Each side pays a tax on what it gives: TaxPct of
// the item's average market price (last 30 days, else its base value) per
// unit — 10 ore at 100 gold average cost 1 000 × 8% = 80 gold.
package trade

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
)

const (
	TaxPct = 8
	// Idle is how long an open trade lives without a change.
	Idle = 15 * time.Minute
	// MaxItems is the most items one side can put up.
	MaxItems = 12
)

var (
	ErrNotFound      = errors.New("trade: not found")
	ErrSelf          = errors.New("trade: cannot trade with yourself")
	ErrBusy          = errors.New("trade: a trade is already open")
	ErrClosed        = errors.New("trade: the trade is over")
	ErrItem          = errors.New("trade: the item cannot be traded")
	ErrTooMany       = errors.New("trade: too many items")
	ErrBagFull       = errors.New("trade: no room in a bag")
	ErrNoGold        = errors.New("trade: not enough gold for the tax")
	ErrNotReadyToGo  = errors.New("trade: nothing to trade")
	errChangedMidway = errors.New("trade: an item changed meanwhile")
)

// Bags places arrived items and counts room (inventory.Service satisfies).
type Bags interface {
	FreeSlots(ctx context.Context, characterID uuid.UUID) (int, error)
	PlaceOverflow(ctx context.Context, characterID uuid.UUID) error
}

// Publisher tells both sides a trade changed (wsgateway adapter satisfies).
type Publisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Service struct {
	pool *pgxpool.Pool
	txn  economy.TxnService
	bags Bags
	pub  Publisher
	log  *slog.Logger
	now  func() time.Time
}

func NewService(pool *pgxpool.Pool, txn economy.TxnService, bags Bags, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, txn: txn, bags: bags, log: logger, now: time.Now}
}

func (s *Service) WithPublisher(p Publisher) *Service { s.pub = p; return s }

// Item is one item put up, with its tax.
type Item struct {
	InstanceID   uuid.UUID `json:"instance_id"`
	TemplateCode string    `json:"template_code"`
	Quantity     int32     `json:"quantity"`
	Rarity       string    `json:"rarity"`
	Prefix       *string   `json:"prefix,omitempty"`
	Suffix       *string   `json:"suffix,omitempty"`
	UnitValue    int64     `json:"unit_value_copper"` // the average market price used
	TaxCopper    int64     `json:"tax_copper"`
}

// Side is one trader's half.
type Side struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	Items       []Item    `json:"items"`
	TaxCopper   int64     `json:"tax_copper"`
	Ready       bool      `json:"ready"`
	account     uuid.UUID
}

// View is an open trade seen by one side.
type View struct {
	ID        uuid.UUID `json:"id"`
	State     string    `json:"state"`
	Me        Side      `json:"me"`
	Them      Side      `json:"them"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Tax is TaxPct of qty units at unitValue copper.
func Tax(unitValue int64, qty int32) int64 {
	return int64(math.Round(float64(unitValue) * float64(qty) * TaxPct / 100))
}

func (s *Service) owns(ctx context.Context, actor, account uuid.UUID) error {
	var ok bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM character.character
		WHERE id = $1 AND account_id = $2)`, actor, account).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// expireIdle closes the actor's open trades past their idle time.
func (s *Service) expireIdle(ctx context.Context, actor uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE trade.session SET state = 'cancelled', closed_at = now(), updated_at = now()
		 WHERE state = 'open' AND expires_at < $2 AND (a_character_id = $1 OR b_character_id = $1)`,
		actor, s.now().UTC())
	return err
}

// Open starts a trade with the character named targetName.
func (s *Service) Open(ctx context.Context, actor, account uuid.UUID, targetName string) (*View, error) {
	if err := s.owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var target uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM character.character WHERE name_normalized = $1
		AND state = 'active'`, strings.ToLower(strings.TrimSpace(targetName))).Scan(&target)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if target == actor {
		return nil, ErrSelf
	}
	for _, c := range []uuid.UUID{actor, target} {
		if err := s.expireIdle(ctx, c); err != nil {
			return nil, err
		}
	}
	var id uuid.UUID
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// One open trade per character: both locked in a fixed order.
		lo, hi := actor, target
		if hi.String() < lo.String() {
			lo, hi = hi, lo
		}
		for _, c := range []uuid.UUID{lo, hi} {
			if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('trade:' || $1::text))`, c); err != nil {
				return err
			}
		}
		var busy bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM trade.session WHERE state = 'open'
			AND (a_character_id = ANY($1) OR b_character_id = ANY($1)))`, []uuid.UUID{actor, target}).Scan(&busy); err != nil {
			return err
		}
		if busy {
			return ErrBusy
		}
		return tx.QueryRow(ctx, `
			INSERT INTO trade.session (a_character_id, b_character_id, expires_at)
			VALUES ($1, $2, $3) RETURNING id`, actor, target, s.now().UTC().Add(Idle)).Scan(&id)
	})
	if err != nil {
		return nil, err
	}
	v, err := s.view(ctx, id, actor)
	if err == nil {
		s.notify(ctx, v)
	}
	return v, err
}

// Current is the actor's open trade, or nil.
func (s *Service) Current(ctx context.Context, actor, account uuid.UUID) (*View, error) {
	if err := s.owns(ctx, actor, account); err != nil {
		return nil, err
	}
	if err := s.expireIdle(ctx, actor); err != nil {
		return nil, err
	}
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM trade.session WHERE state = 'open'
		AND (a_character_id = $1 OR b_character_id = $1)`, actor).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.view(ctx, id, actor)
}

// Offer puts qty of an item up (replacing an earlier qty of it); both
// sides become not ready.
func (s *Service) Offer(ctx context.Context, actor, account, id, itemID uuid.UUID, qty int32) (*View, error) {
	return s.change(ctx, actor, account, id, func(tx pgx.Tx) error {
		var stack int32
		var locked, bound bool
		var equipped *string
		err := tx.QueryRow(ctx, `SELECT stack_quantity, is_locked, is_bound, equipment_slot
			FROM inventory.item_instance WHERE id = $1 AND owner_character_id = $2`, itemID, actor).
			Scan(&stack, &locked, &bound, &equipped)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrItem
		}
		if err != nil {
			return err
		}
		if locked || bound || equipped != nil {
			return ErrItem
		}
		if qty <= 0 || qty > stack {
			qty = stack
		}
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM trade.offer WHERE session_id = $1 AND character_id = $2
			AND item_instance_id <> $3`, id, actor, itemID).Scan(&n); err != nil {
			return err
		}
		if n >= MaxItems {
			return ErrTooMany
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO trade.offer (session_id, character_id, item_instance_id, quantity) VALUES ($1, $2, $3, $4)
			ON CONFLICT (session_id, item_instance_id) DO UPDATE SET quantity = EXCLUDED.quantity`,
			id, actor, itemID, qty)
		return err
	})
}

// Remove takes an item back off the trade.
func (s *Service) Remove(ctx context.Context, actor, account, id, itemID uuid.UUID) (*View, error) {
	return s.change(ctx, actor, account, id, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `DELETE FROM trade.offer WHERE session_id = $1 AND character_id = $2
			AND item_instance_id = $3`, id, actor, itemID)
		return err
	})
}

// Cancel ends the trade; nothing changes hands.
func (s *Service) Cancel(ctx context.Context, actor, account, id uuid.UUID) error {
	if err := s.owns(ctx, actor, account); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `UPDATE trade.session SET state = 'cancelled', closed_at = now(), updated_at = now()
		WHERE id = $1 AND state = 'open' AND (a_character_id = $2 OR b_character_id = $2)`, id, actor)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrClosed
	}
	if v, err := s.view(ctx, id, actor); err == nil {
		s.notify(ctx, v)
	}
	return nil
}

// change runs an offer edit on an open trade of the actor's and clears
// both ready marks.
func (s *Service) change(ctx context.Context, actor, account, id uuid.UUID, edit func(pgx.Tx) error) (*View, error) {
	if err := s.owns(ctx, actor, account); err != nil {
		return nil, err
	}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := s.lockOpen(ctx, tx, id, actor); err != nil {
			return err
		}
		if err := edit(tx); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE trade.session SET a_ready = false, b_ready = false, updated_at = now(),
			expires_at = $2 WHERE id = $1`, id, s.now().UTC().Add(Idle))
		return err
	})
	if err != nil {
		return nil, err
	}
	v, err := s.view(ctx, id, actor)
	if err == nil {
		s.notify(ctx, v)
	}
	return v, err
}

func (s *Service) lockOpen(ctx context.Context, tx pgx.Tx, id, actor uuid.UUID) error {
	var state string
	var expires time.Time
	err := tx.QueryRow(ctx, `SELECT state, expires_at FROM trade.session WHERE id = $1
		AND (a_character_id = $2 OR b_character_id = $2) FOR UPDATE`, id, actor).Scan(&state, &expires)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if state != "open" || expires.Before(s.now().UTC()) {
		return ErrClosed
	}
	return nil
}

// Ready marks the actor's side ready; the second ready side completes the
// trade: taxes paid, items moved, both bags placed.
func (s *Service) Ready(ctx context.Context, actor, account, id uuid.UUID) (*View, error) {
	if err := s.owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var both bool
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := s.lockOpen(ctx, tx, id, actor); err != nil {
			return err
		}
		return tx.QueryRow(ctx, `
			UPDATE trade.session SET
			  a_ready = a_ready OR a_character_id = $2, b_ready = b_ready OR b_character_id = $2, updated_at = now()
			 WHERE id = $1 RETURNING a_ready AND b_ready`, id, actor).Scan(&both)
	})
	if err != nil {
		return nil, err
	}
	if both {
		if err := s.complete(ctx, id); err != nil {
			// The trade stays open, both sides not ready: they fix it and retry.
			_, _ = s.pool.Exec(ctx, `UPDATE trade.session SET a_ready = false, b_ready = false WHERE id = $1`, id)
			if v, verr := s.view(ctx, id, actor); verr == nil {
				s.notify(ctx, v)
			}
			return nil, err
		}
	}
	v, err := s.view(ctx, id, actor)
	if err == nil {
		s.notify(ctx, v)
	}
	return v, err
}

// complete moves the items and charges the taxes of a trade both sides
// are ready for.
func (s *Service) complete(ctx context.Context, id uuid.UUID) error {
	v, err := s.view(ctx, id, uuid.Nil)
	if err != nil {
		return err
	}
	sides := []Side{v.Me, v.Them}
	if len(sides[0].Items)+len(sides[1].Items) == 0 {
		return ErrNotReadyToGo
	}
	// Room in each bag for what arrives.
	for i, sd := range sides {
		other := sides[1-i]
		free, err := s.bags.FreeSlots(ctx, sd.CharacterID)
		if err != nil {
			return err
		}
		if free < len(other.Items) {
			return ErrBagFull
		}
	}
	// Taxes first (idempotent per side); undone if the items cannot move.
	var paid []Side
	refund := func() {
		for _, sd := range paid {
			s.ledger(ctx, fmt.Sprintf("trade-refund-%s-%s", id, sd.CharacterID), sd.CharacterID, sd.TaxCopper, true)
		}
	}
	for _, sd := range sides {
		if sd.TaxCopper <= 0 {
			continue
		}
		if err := s.ledger(ctx, fmt.Sprintf("trade-%s-%s", id, sd.CharacterID), sd.CharacterID, sd.TaxCopper, false); err != nil {
			refund()
			if errors.Is(err, economy.ErrInsufficientBalance) {
				return ErrNoGold
			}
			return err
		}
		paid = append(paid, sd)
	}
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var state string
		if err := tx.QueryRow(ctx, `SELECT state FROM trade.session WHERE id = $1 FOR UPDATE`, id).Scan(&state); err != nil {
			return err
		}
		if state != "open" {
			return ErrClosed
		}
		for i, sd := range sides {
			to := sides[1-i].CharacterID
			for _, it := range sd.Items {
				if err := give(ctx, tx, it, sd.CharacterID, to); err != nil {
					return err
				}
			}
		}
		_, err := tx.Exec(ctx, `UPDATE trade.session SET state = 'done', closed_at = now(), updated_at = now(),
			a_tax_copper = $2, b_tax_copper = $3 WHERE id = $1`, id, sides[0].TaxCopper, sides[1].TaxCopper)
		return err
	})
	if err != nil {
		refund()
		return err
	}
	for _, sd := range sides {
		if err := s.bags.PlaceOverflow(ctx, sd.CharacterID); err != nil {
			s.log.WarnContext(ctx, "trade: placing arrived items failed", "trade", id.String(), "err", err.Error())
		}
	}
	return nil
}

// give moves qty of an item from one character to another (a part of a
// stack splits off); the item must still be theirs and tradeable.
func give(ctx context.Context, tx pgx.Tx, it Item, from, to uuid.UUID) error {
	var stack int32
	var locked, bound bool
	var equipped *string
	err := tx.QueryRow(ctx, `SELECT stack_quantity, is_locked, is_bound, equipment_slot
		FROM inventory.item_instance WHERE id = $1 AND owner_character_id = $2 FOR UPDATE`, it.InstanceID, from).
		Scan(&stack, &locked, &bound, &equipped)
	if errors.Is(err, pgx.ErrNoRows) {
		return errChangedMidway
	}
	if err != nil {
		return err
	}
	if locked || bound || equipped != nil || stack < it.Quantity {
		return errChangedMidway
	}
	if it.Quantity == stack {
		_, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET owner_character_id = $2, slot = NULL,
			updated_at = now(), version = version + 1 WHERE id = $1`, it.InstanceID, to)
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE inventory.item_instance SET stack_quantity = stack_quantity - $2,
		updated_at = now(), version = version + 1 WHERE id = $1`, it.InstanceID, it.Quantity); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO inventory.item_instance
		  (template_code, owner_character_id, stack_quantity, rarity, prefix_code, suffix_code, is_bound, source_kind)
		SELECT template_code, $2, $3, rarity, prefix_code, suffix_code, false, 'trade'
		  FROM inventory.item_instance WHERE id = $1`, it.InstanceID, to, it.Quantity)
	return err
}

// ledger charges (or with back, refunds) a side's tax to the tax sink.
func (s *Service) ledger(ctx context.Context, key string, char uuid.UUID, copper int64, back bool) error {
	sign := int64(1)
	if back {
		sign = -1
	}
	_, err := s.txn.Write(ctx, economy.Intent{
		Kind: "trade_tax", RefKind: "trade", IdempotencyKey: key,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: -sign * copper},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemTaxSink, Delta: sign * copper},
		},
	})
	if err != nil && back {
		s.log.ErrorContext(ctx, "trade: tax refund failed", "key", key, "err", err.Error())
	}
	return err
}

// view loads a trade seen by viewer, who must be in it (uuid.Nil: any
// viewer, Me is the opener).
func (s *Service) view(ctx context.Context, id, viewer uuid.UUID) (*View, error) {
	v := &View{ID: id}
	var a, b Side
	err := s.pool.QueryRow(ctx, `
		SELECT s.state, s.expires_at, s.a_character_id, ca.name, ca.account_id, s.a_ready,
		       s.b_character_id, cb.name, cb.account_id, s.b_ready
		  FROM trade.session s
		  JOIN character.character ca ON ca.id = s.a_character_id
		  JOIN character.character cb ON cb.id = s.b_character_id
		 WHERE s.id = $1`, id).Scan(&v.State, &v.ExpiresAt, &a.CharacterID, &a.Name, &a.account, &a.Ready,
		&b.CharacterID, &b.Name, &b.account, &b.Ready)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if viewer != uuid.Nil && viewer != a.CharacterID && viewer != b.CharacterID {
		return nil, ErrNotFound
	}
	rows, err := s.pool.Query(ctx, `
		SELECT o.character_id, o.item_instance_id, i.template_code, o.quantity, i.rarity, i.prefix_code, i.suffix_code,
		       COALESCE((SELECT round(sum(h.unit_price_copper::numeric * h.quantity) / NULLIF(sum(h.quantity), 0))
		                   FROM market.price_history h
		                  WHERE h.template_code = i.template_code AND h.sold_at > now() - interval '30 days'),
		                t.base_value_copper)::bigint
		  FROM trade.offer o
		  JOIN inventory.item_instance i ON i.id = o.item_instance_id
		  JOIN content.item_template t ON t.code = i.template_code
		 WHERE o.session_id = $1 ORDER BY i.template_code`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var owner uuid.UUID
		var it Item
		if err := rows.Scan(&owner, &it.InstanceID, &it.TemplateCode, &it.Quantity, &it.Rarity, &it.Prefix,
			&it.Suffix, &it.UnitValue); err != nil {
			return nil, err
		}
		it.TaxCopper = Tax(it.UnitValue, it.Quantity)
		sd := &a
		if owner == b.CharacterID {
			sd = &b
		}
		sd.Items = append(sd.Items, it)
		sd.TaxCopper += it.TaxCopper
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, sd := range []*Side{&a, &b} {
		if sd.Items == nil {
			sd.Items = []Item{}
		}
	}
	v.Me, v.Them = a, b
	if viewer == b.CharacterID {
		v.Me, v.Them = b, a
	}
	return v, nil
}

// notify tells both traders' accounts the trade changed.
func (s *Service) notify(_ context.Context, v *View) {
	if s.pub == nil || v == nil {
		return
	}
	for _, sd := range []Side{v.Me, v.Them} {
		s.pub.Publish(sd.account, "trade.changed", map[string]any{"trade_id": v.ID, "state": v.State})
	}
}
