package worker

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgLocker implements Locker with session-level advisory locks. The lock is
// tied to one pooled connection, which is held until unlock.
type PgLocker struct{ pool *pgxpool.Pool }

func NewPgLocker(pool *pgxpool.Pool) *PgLocker { return &PgLocker{pool: pool} }

func (l *PgLocker) TryLock(ctx context.Context, key int64) (func(), bool, error) {
	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		return nil, false, err
	}
	var ok bool
	if err := conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, key).Scan(&ok); err != nil {
		conn.Release()
		return nil, false, err
	}
	if !ok {
		conn.Release()
		return nil, false, nil
	}
	return func() {
		// Background ctx: shutdown must still release the lock.
		_, _ = conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, key)
		conn.Release()
	}, true, nil
}
