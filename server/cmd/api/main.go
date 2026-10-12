// Command api is the Otter HTTP entrypoint.
//
// Wires config → platform (log, db, redis) → domain modules → chi router → HTTP server
// with graceful shutdown on SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // quiet-hours time zones in minimal images

	"github.com/otter/server/internal/app"
	"github.com/otter/server/internal/config"
	"github.com/otter/server/internal/domain/health"
	"github.com/otter/server/internal/observability"
	"github.com/otter/server/internal/platform/db"
	"github.com/otter/server/internal/platform/log"
	"github.com/otter/server/internal/platform/redis"
	"github.com/otter/server/internal/ratelimit"
	"github.com/otter/server/internal/worker"
	"github.com/otter/server/migrations"
)

var CommitSHA = "dev"

func main() {
	// Container HEALTHCHECK (distroless: no wget) — DD-232.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(health.Probe(health.LocalURL(os.Getenv("HTTP_ADDR"), ":8080", "/live")))
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
	logger.Info("boot", "env", cfg.Env, "commit", CommitSHA)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(rootCtx, cfg.DatabaseURL, cfg.DatabaseMaxConns)
	if err != nil {
		return err
	}
	defer pool.Close()
	logger.Info("postgres connected", "max_conns", cfg.DatabaseMaxConns)

	if cfg.MigrateOnStart {
		n, err := db.Migrate(rootCtx, pool, migrations.FS, logger)
		if err != nil {
			return err
		}
		logger.Info("migrations up to date", "applied", n)
	}

	// --- Observability (Sprint 27, DD-214) ---
	shutdownTracing, err := observability.InitTracing(rootCtx, "otter-api", cfg.Env, CommitSHA, cfg.OTLPEndpoint, cfg.TraceSampleRatio)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracing(ctx)
	}()
	metrics := observability.NewMetrics()
	metrics.RegisterPool(pool)
	metrics.Serve(rootCtx, cfg.MetricsAddr, logger)

	rdb, err := redis.Connect(rootCtx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		return err
	}
	defer func() { _ = rdb.Close() }()
	logger.Info("redis connected", "addr", cfg.RedisAddr)

	router, err := app.Build(rootCtx, app.Options{
		Config: cfg, Pool: pool, Redis: rdb, RateLimit: ratelimit.NewRedisStore(rdb),
		Metrics: metrics, Logger: logger, Commit: CommitSHA,
	})
	if err != nil {
		return err
	}

	// Worker jobs in this process (RUN_WORKER, DD-239); they stop with the
	// server and finish before the pool closes.
	workerDone := make(chan struct{})
	if cfg.RunWorker {
		jobs, err := worker.FromConfig(rootCtx, cfg, pool, logger)
		if err != nil {
			return err
		}
		go func() {
			defer close(workerDone)
			// One job at a time: the pool is shared with live requests.
			worker.NewRunner(worker.NewPgLocker(pool), logger, jobs...).WithMetrics(metrics).
				WithMaxParallel(1).Start(rootCtx)
		}()
	} else {
		close(workerDone)
	}
	defer func() { stop(); <-workerDone }()

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("http listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case <-rootCtx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown", "err", err)
	} else {
		logger.Info("http shutdown clean")
	}
	return nil
}
