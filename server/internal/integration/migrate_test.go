//go:build integration

package integration

import (
	"io/fs"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/platform/db"
	"github.com/otter/server/migrations"
)

// The embedded runner (DD-239) brings an empty database — like a fresh
// Supabase project — to the latest version, keeps account tables out of the
// `auth` schema Supabase reserves, records the version where golang-migrate
// reads it, and is a no-op the second time.
func TestMigrateFreshDatabase(t *testing.T) {
	admin, ctx := pool(t)
	name := "otter_migtest_" + randHex(4)
	_, err := admin.Exec(ctx, `CREATE DATABASE `+name)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = admin.Exec(ctx, `DROP DATABASE IF EXISTS `+name+` WITH (FORCE)`) })

	u, err := url.Parse(os.Getenv("DATABASE_URL"))
	require.NoError(t, err)
	u.Path = "/" + name
	fresh, err := db.Connect(ctx, u.String(), 2)
	require.NoError(t, err)
	defer fresh.Close()

	ups, err := fs.Glob(migrations.FS, "*.up.sql")
	require.NoError(t, err)
	n, err := db.Migrate(ctx, fresh, migrations.FS, nil)
	require.NoError(t, err)
	require.Equal(t, len(ups), n)

	var version int64
	var dirty bool
	require.NoError(t, fresh.QueryRow(ctx, `SELECT version, dirty FROM schema_migrations`).Scan(&version, &dirty))
	require.False(t, dirty)
	require.EqualValues(t, len(ups), version, "versions are numbered 1..N")

	var inIdentity, inAuth int
	require.NoError(t, fresh.QueryRow(ctx, `
		SELECT count(*) FILTER (WHERE table_schema = 'identity'), count(*) FILTER (WHERE table_schema = 'auth')
		  FROM information_schema.tables`).Scan(&inIdentity, &inAuth))
	require.Equal(t, 4, inIdentity)
	require.Zero(t, inAuth, "nothing may be created in Supabase's auth schema")
	var rls bool
	require.NoError(t, fresh.QueryRow(ctx, `SELECT relrowsecurity FROM pg_class WHERE oid = 'public.schema_migrations'::regclass`).Scan(&rls))
	require.True(t, rls, "the version table is not reachable through Supabase's Data API")

	n, err = db.Migrate(ctx, fresh, migrations.FS, nil)
	require.NoError(t, err)
	require.Zero(t, n)
}
