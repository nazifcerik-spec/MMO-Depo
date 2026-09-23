// Package worker runs periodic background jobs for cmd/worker (Sprint 19).
//
// Each job runs on its own ticker. Before a run the runner takes a
// Postgres session advisory lock keyed by the job name, so several worker
// replicas can run side by side without double-executing a job. Jobs are
// idempotent SQL batches — a crashed run is simply retried on the next tick.
package worker

import (
	"go.opentelemetry.io/otel/codes"

	"github.com/otter/server/internal/observability"

	"context"
	"hash/fnv"
	"log/slog"
	"sync"
	"time"
)

// Job is one periodic task. Run returns how many rows it touched.
type Job interface {
	Name() string
	Interval() time.Duration
	Run(ctx context.Context) (int64, error)
}

// Locker guards a job against concurrent execution across replicas.
// TryLock returns (unlock, true) when this process owns the lock.
type Locker interface {
	TryLock(ctx context.Context, key int64) (unlock func(), ok bool, err error)
}

type Runner struct {
	jobs    []Job
	locker  Locker
	log     *slog.Logger
	metrics *observability.Metrics // nil → no metrics
}

// WithMetrics records per-job run count, latency and rows (Sprint 27).
func (r *Runner) WithMetrics(m *observability.Metrics) *Runner {
	r.metrics = m
	return r
}

func NewRunner(locker Locker, logger *slog.Logger, jobs ...Job) *Runner {
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{jobs: jobs, locker: locker, log: logger}
}

// Start runs every job once immediately, then on its interval, until ctx is
// cancelled. Blocks until all job loops have returned.
func (r *Runner) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for _, j := range r.jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			r.loop(ctx, j)
		}(j)
	}
	wg.Wait()
}

func (r *Runner) loop(ctx context.Context, j Job) {
	t := time.NewTicker(j.Interval())
	defer t.Stop()
	for {
		r.RunOnce(ctx, j)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// RunOnce executes a single job under its advisory lock. Exported for tests
// and for a future ops "run now" trigger.
func (r *Runner) RunOnce(ctx context.Context, j Job) {
	if ctx.Err() != nil {
		return
	}
	unlock, ok, err := r.locker.TryLock(ctx, LockKey(j.Name()))
	if err != nil {
		r.log.WarnContext(ctx, "worker: lock failed", "job", j.Name(), "err", err.Error())
		return
	}
	if !ok {
		r.log.DebugContext(ctx, "worker: job held by another replica", "job", j.Name())
		return
	}
	defer unlock()

	ctx, span := observability.StartSpan(ctx, "worker."+j.Name())
	defer span.End()
	start := time.Now()
	n, err := j.Run(ctx)
	dur := time.Since(start)
	r.metrics.ObserveJob(j.Name(), n, dur, err)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		r.log.ErrorContext(ctx, "worker: job failed",
			"job", j.Name(), "err", err.Error(), "duration_ms", dur.Milliseconds())
		return
	}
	level := slog.LevelDebug
	if n > 0 {
		level = slog.LevelInfo
	}
	r.log.Log(ctx, level, "worker: job done",
		"job", j.Name(), "rows", n, "duration_ms", dur.Milliseconds())
}

// LockKey maps a job name to a stable advisory-lock key.
func LockKey(name string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte("otter.worker." + name))
	return int64(h.Sum64())
}
