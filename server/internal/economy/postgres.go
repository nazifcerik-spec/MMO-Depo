package economy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgTxnService is the production TxnService backed by Postgres SERIALIZABLE
// transactions with retry on 40001 (serialization_failure).
type PgTxnService struct {
	pool       *pgxpool.Pool
	currencies map[string]Currency
	log        *slog.Logger
	now        func() time.Time
	maxRetries int
}

// NewPgTxnService constructs the service with a loaded currency catalog.
func NewPgTxnService(pool *pgxpool.Pool, currencies map[string]Currency, logger *slog.Logger) *PgTxnService {
	if logger == nil {
		logger = slog.Default()
	}
	return &PgTxnService{
		pool:       pool,
		currencies: currencies,
		log:        logger,
		now:        time.Now,
		maxRetries: 8,
	}
}

// Balance reads the current projection row for the given account/currency.
// Missing row = 0 (never written to).
func (s *PgTxnService) Balance(ctx context.Context, kind AccountKind, ref uuid.UUID, currency string) (int64, error) {
	return s.readBalance(ctx, s.pool, kind, ref, currency)
}

// balanceReader is the subset of pgxpool.Pool + pgx.Tx that our helper needs.
type balanceReader interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (s *PgTxnService) readBalance(ctx context.Context, r balanceReader, kind AccountKind, ref uuid.UUID, currency string) (int64, error) {
	sql, err := balanceSelect(kind, currency)
	if err != nil {
		return 0, err
	}
	var bal int64
	err = r.QueryRow(ctx, sql, ref).Scan(&bal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return bal, nil
}

// balanceSelect returns the read-only SELECT for a (kind, currency) pair.
// Split by table: character_wallet is gold-only; the other two are (id, code).
func balanceSelect(kind AccountKind, currency string) (string, error) {
	switch kind {
	case AccountCharacter:
		if currency == "gold" {
			return `SELECT gold_copper FROM economy.character_wallet WHERE character_id = $1`, nil
		}
		return `SELECT amount FROM economy.character_currency_balance
		         WHERE character_id = $1 AND currency_code = '` + safeCurrency(currency) + `'`, nil
	case AccountAccount:
		return `SELECT amount FROM economy.account_currency_balance
		         WHERE account_id = $1 AND currency_code = '` + safeCurrency(currency) + `'`, nil
	case AccountGuild:
		// Guild treasury table lands with Faz 01 sprint 15+. Placeholder.
		return "", fmt.Errorf("guild balance table not yet implemented")
	case AccountSystem:
		// System balances are ledger-only; not a projection table.
		// Sum recent entries — for now not supported (system flow is measured
		// via metrics, not projection lookup).
		return "", fmt.Errorf("system balance projection not supported — sum ledger_entry directly")
	default:
		return "", fmt.Errorf("unknown account kind %s", kind)
	}
}

// safeCurrency guards against SQL injection since we inline the code (BALL is
// finite currency catalog). ValidateIntent already checks membership; this is
// belt-and-suspenders.
func safeCurrency(code string) string {
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return "invalid"
		}
	}
	return code
}

// Write is the SERIALIZABLE ledger-write path. Retries up to maxRetries on
// serialization_failure with jittered backoff. Returns Result with Replayed=true when the
// idempotency_key matches an existing txn.
func (s *PgTxnService) Write(ctx context.Context, intent Intent) (Result, error) {
	if err := ValidateIntent(intent, s.currencies); err != nil {
		return Result{}, err
	}

	// Fast-path idempotency check outside the write TX.
	if cached, ok, err := s.findExistingTxn(ctx, intent.IdempotencyKey); err != nil {
		return Result{}, err
	} else if ok {
		cached.Replayed = true
		return cached, nil
	}

	var lastErr error
	for attempt := 0; attempt < s.maxRetries; attempt++ {
		if attempt > 0 {
			// Jittered exponential backoff (4ms, 8ms, … ×1–2): writers that
			// collided on one account spread out instead of colliding again.
			d := time.Duration(2<<attempt) * time.Millisecond
			d += time.Duration(rand.Int64N(int64(d)))
			select {
			case <-ctx.Done():
				return Result{}, ctx.Err()
			case <-time.After(d):
			}
		}
		res, err := s.writeOnce(ctx, intent)
		if err == nil {
			return res, nil
		}
		if !isSerializationFailure(err) {
			return Result{}, err
		}
		lastErr = err
		s.log.WarnContext(ctx, "ledger serialization retry",
			slog.Int("attempt", attempt+1),
			slog.String("err", err.Error()))
	}
	return Result{}, fmt.Errorf("%w: %v", ErrSerialization, lastErr)
}

func isSerializationFailure(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "40001"
}

// findExistingTxn returns the cached Result for an idempotency_key if one
// exists. Uses a plain read (no lock) — race with a concurrent writer is fine:
// worst case both attempts fall through to writeOnce and the UNIQUE constraint
// resolves.
func (s *PgTxnService) findExistingTxn(ctx context.Context, key string) (Result, bool, error) {
	const q = `SELECT id, txn_id, entry_seq, kind, currency_code, account_kind,
	                  account_ref, delta, balance_after, ref_kind, ref_id, idempotency_key
	             FROM economy.ledger_entry
	            WHERE idempotency_key = $1
	         ORDER BY entry_seq`
	rows, err := s.pool.Query(ctx, q, key)
	if err != nil {
		return Result{}, false, err
	}
	defer rows.Close()

	var result Result
	for rows.Next() {
		row, err := scanLedgerRow(rows)
		if err != nil {
			return Result{}, false, err
		}
		result.TxnID = row.TxnID
		result.Rows = append(result.Rows, row)
	}
	if err := rows.Err(); err != nil {
		return Result{}, false, err
	}
	return result, len(result.Rows) > 0, nil
}

// writeOnce runs the ledger insertion + projection upsert in one SERIALIZABLE
// transaction. Called by Write; the caller handles retry on 40001.
func (s *PgTxnService) writeOnce(ctx context.Context, intent Intent) (Result, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return Result{}, err
	}
	// Defer rollback for the error path; committed transactions ignore Rollback.
	defer func() { _ = tx.Rollback(ctx) }()

	// Second idempotency check inside the TX — a concurrent writer may have
	// completed between our fast-path and this point.
	if cached, ok, err := s.findExistingTxnTx(ctx, tx, intent.IdempotencyKey); err != nil {
		return Result{}, err
	} else if ok {
		cached.Replayed = true
		return cached, nil
	}

	txnID := uuid.New()
	rows := make([]LedgerRow, 0, len(intent.Entries))

	// Group entries by (kind, ref, currency) so we lock and update each
	// balance once — order deterministic by seq for reproducibility.
	// For each entry we still compute a per-entry balance_after.
	perAccount := make(map[balanceKey]int64) // running balance during this txn

	for seq, e := range intent.Entries {
		key := balanceKey{Kind: e.AccountKind, Ref: e.AccountRef, Currency: e.Currency}

		// Establish initial balance the first time we see this account.
		if _, seen := perAccount[key]; !seen {
			if e.AccountKind == AccountSystem {
				perAccount[key] = 0 // unbounded; balance_after is cumulative in this txn
			} else {
				bal, err := s.lockAndReadBalance(ctx, tx, e.AccountKind, e.AccountRef, e.Currency)
				if err != nil {
					return Result{}, err
				}
				perAccount[key] = bal
			}
		}

		newBal := perAccount[key] + e.Delta

		if e.AccountKind != AccountSystem {
			cur := s.currencies[e.Currency]
			if newBal < 0 {
				return Result{}, fmt.Errorf("%w: currency=%s account=%s new=%d",
					ErrInsufficientBalance, e.Currency, e.AccountRef, newBal)
			}
			if newBal > cur.Cap {
				return Result{}, fmt.Errorf("%w: currency=%s cap=%d new=%d",
					ErrCurrencyCap, e.Currency, cur.Cap, newBal)
			}
		}
		perAccount[key] = newBal

		var refIDPtr *uuid.UUID
		if intent.RefID != uuid.Nil {
			rid := intent.RefID
			refIDPtr = &rid
		}

		row := LedgerRow{
			TxnID:        txnID,
			EntrySeq:     int16(seq),
			Kind:         intent.Kind,
			Currency:     e.Currency,
			AccountKind:  e.AccountKind,
			AccountRef:   e.AccountRef,
			Delta:        e.Delta,
			BalanceAfter: newBal,
			RefKind:      intent.RefKind,
			RefID:        refIDPtr,
			IdempKey:     intent.IdempotencyKey,
		}
		rows = append(rows, row)
	}

	// Persist ledger entries.
	if err := insertLedgerEntries(ctx, tx, rows); err != nil {
		return Result{}, err
	}

	// Upsert projections (skip system).
	for key, bal := range perAccount {
		if key.Kind == AccountSystem {
			continue
		}
		if err := upsertBalance(ctx, tx, key, bal); err != nil {
			return Result{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Result{}, err
	}
	return Result{TxnID: txnID, Rows: rows}, nil
}

// lockAndReadBalance runs SELECT ... FOR UPDATE, upserting a zero row when the
// projection doesn't exist yet. Ensures every subsequent update is safe.
func (s *PgTxnService) lockAndReadBalance(ctx context.Context, tx pgx.Tx, kind AccountKind, ref uuid.UUID, currency string) (int64, error) {
	switch kind {
	case AccountCharacter:
		if currency == "gold" {
			// Upsert zero row first — idempotent + gives us a lock.
			_, err := tx.Exec(ctx, `
				INSERT INTO economy.character_wallet (character_id, gold_copper)
				VALUES ($1, 0)
				ON CONFLICT (character_id) DO NOTHING`, ref)
			if err != nil {
				return 0, err
			}
			var bal int64
			err = tx.QueryRow(ctx,
				`SELECT gold_copper FROM economy.character_wallet WHERE character_id = $1 FOR UPDATE`,
				ref).Scan(&bal)
			return bal, err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO economy.character_currency_balance (character_id, currency_code, amount)
			VALUES ($1, $2, 0)
			ON CONFLICT (character_id, currency_code) DO NOTHING`, ref, currency)
		if err != nil {
			return 0, err
		}
		var bal int64
		err = tx.QueryRow(ctx,
			`SELECT amount FROM economy.character_currency_balance
			  WHERE character_id = $1 AND currency_code = $2 FOR UPDATE`,
			ref, currency).Scan(&bal)
		return bal, err

	case AccountAccount:
		_, err := tx.Exec(ctx, `
			INSERT INTO economy.account_currency_balance (account_id, currency_code, amount)
			VALUES ($1, $2, 0)
			ON CONFLICT (account_id, currency_code) DO NOTHING`, ref, currency)
		if err != nil {
			return 0, err
		}
		var bal int64
		err = tx.QueryRow(ctx,
			`SELECT amount FROM economy.account_currency_balance
			  WHERE account_id = $1 AND currency_code = $2 FOR UPDATE`,
			ref, currency).Scan(&bal)
		return bal, err

	case AccountGuild:
		return 0, fmt.Errorf("guild balance not yet implemented")
	}
	return 0, fmt.Errorf("unknown account kind %s", kind)
}

func upsertBalance(ctx context.Context, tx pgx.Tx, key balanceKey, newBal int64) error {
	switch key.Kind {
	case AccountCharacter:
		if key.Currency == "gold" {
			_, err := tx.Exec(ctx, `
				UPDATE economy.character_wallet
				   SET gold_copper = $2, updated_at = now(), version = version + 1
				 WHERE character_id = $1`, key.Ref, newBal)
			return err
		}
		_, err := tx.Exec(ctx, `
			UPDATE economy.character_currency_balance
			   SET amount = $3, updated_at = now(), version = version + 1
			 WHERE character_id = $1 AND currency_code = $2`, key.Ref, key.Currency, newBal)
		return err
	case AccountAccount:
		_, err := tx.Exec(ctx, `
			UPDATE economy.account_currency_balance
			   SET amount = $3, updated_at = now(), version = version + 1
			 WHERE account_id = $1 AND currency_code = $2`, key.Ref, key.Currency, newBal)
		return err
	}
	return fmt.Errorf("unknown account kind %s", key.Kind)
}

func insertLedgerEntries(ctx context.Context, tx pgx.Tx, rows []LedgerRow) error {
	const q = `
		INSERT INTO economy.ledger_entry
		  (txn_id, entry_seq, kind, currency_code, account_kind, account_ref,
		   delta, balance_after, ref_kind, ref_id, idempotency_key)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id`
	batch := &pgx.Batch{}
	for _, r := range rows {
		var refID any
		if r.RefID != nil {
			refID = *r.RefID
		}
		var refKind any
		if r.RefKind != "" {
			refKind = r.RefKind
		}
		batch.Queue(q,
			r.TxnID, r.EntrySeq, r.Kind, r.Currency,
			string(r.AccountKind), r.AccountRef, r.Delta, r.BalanceAfter,
			refKind, refID, r.IdempKey,
		)
	}
	br := tx.SendBatch(ctx, batch)
	defer br.Close()
	for i := range rows {
		var id int64
		if err := br.QueryRow().Scan(&id); err != nil {
			return err
		}
		rows[i].ID = id
	}
	return nil
}

// findExistingTxnTx is the TX-scoped variant of findExistingTxn.
func (s *PgTxnService) findExistingTxnTx(ctx context.Context, tx pgx.Tx, key string) (Result, bool, error) {
	const q = `SELECT id, txn_id, entry_seq, kind, currency_code, account_kind,
	                  account_ref, delta, balance_after, ref_kind, ref_id, idempotency_key
	             FROM economy.ledger_entry
	            WHERE idempotency_key = $1
	         ORDER BY entry_seq`
	rows, err := tx.Query(ctx, q, key)
	if err != nil {
		return Result{}, false, err
	}
	defer rows.Close()

	var result Result
	for rows.Next() {
		row, err := scanLedgerRow(rows)
		if err != nil {
			return Result{}, false, err
		}
		result.TxnID = row.TxnID
		result.Rows = append(result.Rows, row)
	}
	return result, len(result.Rows) > 0, rows.Err()
}

// scanLedgerRow is the shared scan for both TX and pool query rows.
func scanLedgerRow(row pgx.Row) (LedgerRow, error) {
	var (
		lr       LedgerRow
		acctKind string
		refKind  *string
		refID    *uuid.UUID
	)
	if err := row.Scan(
		&lr.ID, &lr.TxnID, &lr.EntrySeq, &lr.Kind, &lr.Currency,
		&acctKind, &lr.AccountRef, &lr.Delta, &lr.BalanceAfter,
		&refKind, &refID, &lr.IdempKey,
	); err != nil {
		return LedgerRow{}, err
	}
	lr.AccountKind = AccountKind(acctKind)
	if refKind != nil {
		lr.RefKind = *refKind
	}
	lr.RefID = refID
	return lr, nil
}
