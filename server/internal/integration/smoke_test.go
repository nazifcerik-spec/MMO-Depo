//go:build integration

// Package integration holds tests that need a real, fully migrated Postgres.
// CI runs them after `migrate up` (see .github/workflows/ci.yml):
//
//	DATABASE_URL=postgres://... go test -tags integration ./internal/integration/
//
// They exercise the SQL that unit tests (memory repos) cannot: content-pack
// upserts, catalog loads and the ledger's SERIALIZABLE write path.
package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/classskill"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/platform/db"
	"github.com/otter/server/internal/world"
)

func pool(t *testing.T) (*pgxpool.Pool, context.Context) {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)
	p, err := db.Connect(ctx, url, 4)
	require.NoError(t, err)
	t.Cleanup(p.Close)
	return p, ctx
}

// The shipped content pack must upsert cleanly and every catalog must load
// from the resulting rows (twice: the loader is re-run on hot reload).
func TestContentPackAndCatalogs(t *testing.T) {
	p, ctx := pool(t)
	loader := contentpack.New(filepath.Join("..", "..", "..", "content"), p, nil)
	for i := 0; i < 2; i++ {
		_, err := loader.Load(ctx)
		require.NoError(t, err, "content load pass %d", i+1)
	}

	items, err := item.LoadCatalog(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, items)
	mobs, err := combat.LoadCatalog(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, mobs.All())
	_, err = dungeon.LoadCatalog(ctx, p)
	require.NoError(t, err)
	regions, err := world.LoadCatalog(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, regions.StartingCode())
	_, err = crafting.LoadCatalog(ctx, p)
	require.NoError(t, err)

	skills := classskill.NewService(p)
	require.NoError(t, skills.Reload(ctx))
	require.NotEmpty(t, skills.Catalog().Actives([]string{"warrior"}, 1), "every class starts with an active")
}

// A system→system transfer goes through the real ledger path; the same
// idempotency key must replay instead of writing twice (DD-135, DD-188).
func TestLedgerWriteAndReplay(t *testing.T) {
	p, ctx := pool(t)
	currencies, err := economy.LoadCurrencies(ctx, p)
	require.NoError(t, err)
	txn := economy.NewPgTxnService(p, currencies, nil)

	intent := economy.Intent{
		Kind:           "integration_smoke",
		RefKind:        "test",
		RefID:          uuid.New(),
		IdempotencyKey: "integration-smoke-" + uuid.NewString(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -100},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: 100},
		},
	}
	first, err := txn.Write(ctx, intent)
	require.NoError(t, err)
	require.False(t, first.Replayed)
	require.Len(t, first.Rows, 2)

	again, err := txn.Write(ctx, intent)
	require.NoError(t, err)
	require.True(t, again.Replayed)
	require.Equal(t, first.TxnID, again.TxnID)
}
