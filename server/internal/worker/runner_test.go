package worker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeJob struct {
	name  string
	every time.Duration
	runs  atomic.Int32
	err   error
}

func (f *fakeJob) Name() string            { return f.name }
func (f *fakeJob) Interval() time.Duration { return f.every }
func (f *fakeJob) Run(context.Context) (int64, error) {
	f.runs.Add(1)
	return 1, f.err
}

// memLocker is a process-local stand-in for the advisory lock.
type memLocker struct {
	mu      sync.Mutex
	held    map[int64]bool
	unlocks int
}

func newMemLocker() *memLocker { return &memLocker{held: map[int64]bool{}} }

func (l *memLocker) TryLock(_ context.Context, key int64) (func(), bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.held[key] {
		return nil, false, nil
	}
	l.held[key] = true
	return func() {
		l.mu.Lock()
		delete(l.held, key)
		l.unlocks++
		l.mu.Unlock()
	}, true, nil
}

func TestRunOnce_RunsAndReleasesLock(t *testing.T) {
	l := newMemLocker()
	j := &fakeJob{name: "a", every: time.Hour}
	r := NewRunner(l, nil, j)
	r.RunOnce(context.Background(), j)
	r.RunOnce(context.Background(), j)
	require.Equal(t, int32(2), j.runs.Load())
	require.Equal(t, 2, l.unlocks)
	require.Empty(t, l.held)
}

func TestRunOnce_SkipsWhenLockHeldElsewhere(t *testing.T) {
	l := newMemLocker()
	j := &fakeJob{name: "a", every: time.Hour}
	l.held[LockKey("a")] = true // another replica owns it
	NewRunner(l, nil, j).RunOnce(context.Background(), j)
	require.Zero(t, j.runs.Load())
}

func TestRunOnce_JobErrorStillReleasesLock(t *testing.T) {
	l := newMemLocker()
	j := &fakeJob{name: "boom", every: time.Hour, err: errors.New("db down")}
	NewRunner(l, nil, j).RunOnce(context.Background(), j)
	require.Equal(t, int32(1), j.runs.Load())
	require.Empty(t, l.held)
}

type errLocker struct{}

func (errLocker) TryLock(context.Context, int64) (func(), bool, error) {
	return nil, false, errors.New("no conn")
}

func TestRunOnce_LockErrorSkipsJob(t *testing.T) {
	j := &fakeJob{name: "a", every: time.Hour}
	NewRunner(errLocker{}, nil, j).RunOnce(context.Background(), j)
	require.Zero(t, j.runs.Load())
}

func TestStart_RunsImmediatelyAndStopsOnCancel(t *testing.T) {
	l := newMemLocker()
	fast := &fakeJob{name: "fast", every: 10 * time.Millisecond}
	slow := &fakeJob{name: "slow", every: time.Hour}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		NewRunner(l, nil, fast, slow).Start(ctx)
		close(done)
	}()
	require.Eventually(t, func() bool { return fast.runs.Load() >= 3 }, time.Second, 5*time.Millisecond)
	require.Eventually(t, func() bool { return slow.runs.Load() == 1 }, time.Second, 5*time.Millisecond,
		"every job runs once at startup")
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop after cancel")
	}
}

func TestLockKey_StableAndDistinct(t *testing.T) {
	require.Equal(t, LockKey("activity_expire"), LockKey("activity_expire"))
	require.NotEqual(t, LockKey("activity_expire"), LockKey("travel_arrive"))
}

func TestAll_UniqueNamesAndPositiveIntervals(t *testing.T) {
	seen := map[string]bool{}
	for _, j := range All(nil, Deps{}) {
		require.False(t, seen[j.Name()], "duplicate job %s", j.Name())
		seen[j.Name()] = true
		require.Positive(t, j.Interval())
	}
	require.Len(t, seen, 16)
}
