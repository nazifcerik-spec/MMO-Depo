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
