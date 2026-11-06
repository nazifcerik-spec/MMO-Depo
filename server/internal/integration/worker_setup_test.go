//go:build integration

package integration

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/worker"
)

// The API's RUN_WORKER mode and cmd/worker build the same job set (DD-239).
func TestWorkerFromConfig(t *testing.T) {
	p, ctx := pool(t)
	jobs, err := worker.FromConfig(ctx, config.Config{Env: "dev"}, p, slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, err)
	require.Len(t, jobs, len(worker.All(nil, worker.Deps{})))
}

// The advisory lock keeps a job on one replica (Sprint 19): a second holder is
// refused until the first unlocks, and other keys stay free.
func TestPgLockerIsExclusive(t *testing.T) {
	p, ctx := pool(t)
	l := worker.NewPgLocker(p)
	const key = int64(0x6f74746572) // a key no job name hashes to in this test run

	unlock, ok, err := l.TryLock(ctx, key)
	require.NoError(t, err)
	require.True(t, ok)

	_, ok, err = l.TryLock(ctx, key)
	require.NoError(t, err)
	require.False(t, ok, "held by the first caller")

	other, ok, err := l.TryLock(ctx, key+1)
	require.NoError(t, err)
	require.True(t, ok, "another key is free")
	other()

	unlock()
	again, ok, err := l.TryLock(ctx, key)
	require.NoError(t, err)
	require.True(t, ok, "free after unlock")
	again()
}

// Every job's statements run against the migrated schema: a column renamed by
// a migration must not reach a worker unnoticed.
func TestEveryWorkerJobRuns(t *testing.T) {
	p, ctx := pool(t)
	jobs, err := worker.FromConfig(ctx, config.Config{Env: "dev"}, p, slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, err)
	require.NotEmpty(t, jobs)
	for _, j := range jobs {
		_, err := j.Run(ctx)
		require.NoError(t, err, j.Name())
	}
}
