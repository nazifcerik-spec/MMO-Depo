// Command worker runs Otter background jobs (see internal/worker). Safe to
// run as several replicas — each job takes a Postgres advisory lock per run.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // quiet-hours time zones in minimal images

	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/domain/health"
	"github.com/otter/server/internal/observability"
	"github.com/otter/server/internal/platform/db"
	"github.com/otter/server/internal/platform/log"
	"github.com/otter/server/internal/worker"
)

var CommitSHA = "dev"

func main() {
	// Container HEALTHCHECK (distroless: no wget) — DD-232.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(health.Probe(health.LocalURL(os.Getenv("WORKER_METRICS_ADDR"), ":9091", "/metrics")))
	}
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := log.New(cfg.LogLevel, cfg.Env, CommitSHA)
	slog.SetDefault(logger)
	logger.Info("worker boot", "env", cfg.Env, "commit", CommitSHA)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL, cfg.DatabaseMaxConns)
	if err != nil {
		return err
	}
	defer pool.Close()

	jobs, err := worker.FromConfig(ctx, cfg, pool, logger)
	if err != nil {
		return err
	}
	shutdownTracing, err := observability.InitTracing(ctx, "otter-worker", cfg.Env, CommitSHA, cfg.OTLPEndpoint, cfg.TraceSampleRatio)
	if err != nil {
		return err
	}
	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracing(sctx)
	}()
	metrics := observability.NewMetrics()
	metrics.RegisterPool(pool)
	metrics.Serve(ctx, cfg.WorkerMetricsAddr, logger)

	worker.NewRunner(worker.NewPgLocker(pool), logger, jobs...).WithMetrics(metrics).
		WithMaxParallel(worker.ParallelFor(cfg.DatabaseMaxConns)).Start(ctx)
	logger.Info("worker shutdown")
	return nil
}
