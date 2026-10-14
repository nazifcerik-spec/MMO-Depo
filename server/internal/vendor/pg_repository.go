package vendor

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepo struct {
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool) *PgRepo { return &PgRepo{pool: pool} }

func (r *PgRepo) Catalog(ctx context.Context) ([]Template, error) {
	const q = `
		SELECT code, display_name_key, vendor_type, buy_multiplier, sell_multiplier, sells
		  FROM content.vendor_template`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Template{}
	for rows.Next() {
		var t Template
		var sells []byte
		if err := rows.Scan(
			&t.Code, &t.DisplayNameKey, &t.VendorType,
			&t.BuyMultiplier, &t.SellMultiplier, &sells,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(sells, &t.Sells); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *PgRepo) InsertBuyback(ctx context.Context, b BuybackEntry) error {
	const q = `
		INSERT INTO economy.vendor_buyback
		  (id, character_id, vendor_code, item_snapshot, sold_price_copper, sold_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q,
		b.ID, b.CharacterID, b.VendorCode, []byte(b.ItemSnapshot),
		b.SoldPriceCopper, b.SoldAt, b.ExpiresAt,
	)
	return err
}

const buybackSelect = `
	SELECT id, character_id, vendor_code, item_snapshot, sold_price_copper,
	       sold_at, expires_at, claimed_at
	  FROM economy.vendor_buyback`

func (r *PgRepo) scanBuyback(row pgx.Row) (*BuybackEntry, error) {
	var b BuybackEntry
	var snap []byte
	err := row.Scan(
		&b.ID, &b.CharacterID, &b.VendorCode, &snap, &b.SoldPriceCopper,
		&b.SoldAt, &b.ExpiresAt, &b.ClaimedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBuybackNotFound
		}
		return nil, err
	}
	b.ItemSnapshot = json.RawMessage(snap)
	return &b, nil
}

func (r *PgRepo) GetBuyback(ctx context.Context, id uuid.UUID) (*BuybackEntry, error) {
	return r.scanBuyback(r.pool.QueryRow(ctx, buybackSelect+" WHERE id = $1", id))
}

func (r *PgRepo) Rebind(ctx context.Context, instanceID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE inventory.item_instance SET is_bound = true WHERE id = $1`, instanceID)
	return err
}

func (r *PgRepo) ClaimBuyback(ctx context.Context, id uuid.UUID, at time.Time) error {
	const q = `UPDATE economy.vendor_buyback
	              SET claimed_at = $2
	            WHERE id = $1 AND claimed_at IS NULL AND expires_at > $2`
	tag, err := r.pool.Exec(ctx, q, id, at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		// Distinguish expired vs claimed for the client.
		existing, gerr := r.GetBuyback(ctx, id)
		if gerr != nil {
			return gerr
		}
		if existing.ClaimedAt != nil {
			return ErrBuybackClaimed
		}
		return ErrBuybackExpired
	}
	return nil
}

func (r *PgRepo) ListBuybacks(ctx context.Context, characterID uuid.UUID, limit int) ([]BuybackEntry, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx,
		buybackSelect+` WHERE character_id = $1 AND claimed_at IS NULL AND expires_at > now()
		                ORDER BY sold_at DESC LIMIT $2`,
		characterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BuybackEntry{}
	for rows.Next() {
		b, err := r.scanBuyback(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *b)
	}
	return out, rows.Err()
}
