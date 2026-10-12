package db

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"regexp"
	"sort"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// migrateLockKey serialises Migrate across instances starting together.
const migrateLockKey = 7_165_243_001

var upFile = regexp.MustCompile(`^(\d+)_[^/]+\.up\.sql$`)

// Migrate applies the up-migrations in files that are newer than the version
// recorded in schema_migrations — the table golang-migrate keeps (dev, CI and
// the prod compose use that CLI), so either tool can take over (DD-239). Each
// file runs in one transaction with its version bump, so a failure leaves the
// database at the previous version instead of "dirty". Returns how many were
// applied.
func Migrate(ctx context.Context, pool *pgxpool.Pool, files fs.FS, log *slog.Logger) (int, error) {
	type migration struct {
		version int64
		name    string
	}
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return 0, fmt.Errorf("migrate: read files: %w", err)
	}
	var all []migration
	for _, e := range entries {
		if m := upFile.FindStringSubmatch(e.Name()); m != nil {
			v, _ := strconv.ParseInt(m[1], 10, 64)
			all = append(all, migration{v, e.Name()})
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].version < all[j].version })

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("migrate: acquire: %w", err)
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, migrateLockKey); err != nil {
		return 0, fmt.Errorf("migrate: lock: %w", err)
	}
	defer func() { _, _ = conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, migrateLockKey) }()

	if _, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL PRIMARY KEY, dirty boolean NOT NULL)`); err != nil {
		return 0, fmt.Errorf("migrate: version table: %w", err)
	}
	// Supabase serves the public schema over its Data API with a published
	// key; RLS without policies shuts that out (the owning role is exempt).
	if _, err := conn.Exec(ctx, `ALTER TABLE schema_migrations ENABLE ROW LEVEL SECURITY`); err != nil {
		return 0, fmt.Errorf("migrate: version table rls: %w", err)
	}
	var current int64
	var dirty bool
	err = conn.QueryRow(ctx, `SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&current, &dirty)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		current = 0
	case err != nil:
		return 0, fmt.Errorf("migrate: read version: %w", err)
	case dirty:
		return 0, fmt.Errorf("migrate: database is dirty at version %d — fix it by hand, then clear the flag", current)
	}

	applied := 0
	for _, m := range all {
		if m.version <= current {
			continue
		}
		body, err := fs.ReadFile(files, m.name)
		if err != nil {
			return applied, fmt.Errorf("migrate: read %s: %w", m.name, err)
		}
		err = pgx.BeginFunc(ctx, conn.Conn(), func(tx pgx.Tx) error {
			// Simple protocol: a file holds several statements (and DO blocks).
			if _, err := tx.Conn().PgConn().Exec(ctx, string(body)).ReadAll(); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `DELETE FROM schema_migrations`); err != nil {
				return err
			}
			_, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)`, m.version)
			return err
		})
		if err != nil {
			return applied, fmt.Errorf("migrate: %s: %w", m.name, err)
		}
		applied++
		if log != nil {
			log.Info("migration applied", "file", m.name)
		}
	}
	return applied, nil
}
