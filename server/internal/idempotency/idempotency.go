// Package idempotency implements DD-188: every mutating endpoint requires a client
// Idempotency-Key; the same (key, endpoint, account_id) with the same request body
// returns the cached response; with a different body it is rejected with
// E_IDEMPOTENT_MISMATCH. Retention 24h; cleanup by a periodic worker.
//
// The Store interface separates the wire from persistence so unit tests can drop
// in MemoryStore and integration/production wires up PgStore.
package idempotency

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound is returned by Store.Get when no record matches. It is not an
// error condition — a MISS means the middleware should execute the handler.
var ErrNotFound = errors.New("idempotency: not found")

// Record is a single cached response.
type Record struct {
	Key            string
	Endpoint       string
	AccountID      uuid.UUID
	RequestHash    []byte // sha256 of request body
	ResponseStatus int16
	ResponseBody   []byte
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

// Store is the persistence contract. Get returns ErrNotFound on MISS; other
// errors are transport failures. Set is best-effort — on unique-key conflict
// (race) it returns nil, matching the ON CONFLICT DO NOTHING policy in SQL.
type Store interface {
	Get(ctx context.Context, key, endpoint string, accountID uuid.UUID) (*Record, error)
	Set(ctx context.Context, rec Record) error
}

// --- MemoryStore (tests) -----------------------------------------------------

// MemoryStore is a goroutine-safe in-memory Store. Tests only; not exported for
// production use.
type MemoryStore struct {
	mu   sync.Mutex
	rows map[string]Record
	now  func() time.Time
}

// NewMemoryStore returns an empty store with real time.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rows: make(map[string]Record),
		now:  time.Now,
	}
}

// NewMemoryStoreWithClock lets tests freeze time to assert expiration behavior.
func NewMemoryStoreWithClock(clock func() time.Time) *MemoryStore {
	return &MemoryStore{
		rows: make(map[string]Record),
		now:  clock,
	}
}

func memKey(key, endpoint string, accountID uuid.UUID) string {
	return key + "|" + endpoint + "|" + accountID.String()
}

func (m *MemoryStore) Get(_ context.Context, key, endpoint string, accountID uuid.UUID) (*Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.rows[memKey(key, endpoint, accountID)]
	if !ok || rec.ExpiresAt.Before(m.now()) {
		return nil, ErrNotFound
	}
	// Return a copy so callers can't mutate the stored bytes.
	cp := rec
	cp.RequestHash = append([]byte(nil), rec.RequestHash...)
	cp.ResponseBody = append([]byte(nil), rec.ResponseBody...)
	return &cp, nil
}

func (m *MemoryStore) Set(_ context.Context, rec Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := memKey(rec.Key, rec.Endpoint, rec.AccountID)
	if _, exists := m.rows[k]; exists {
		// ON CONFLICT DO NOTHING — race loser is silent.
		return nil
	}
	now := m.now()
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = now
	}
	if rec.ExpiresAt.IsZero() {
		rec.ExpiresAt = now.Add(24 * time.Hour)
	}
	rec.RequestHash = append([]byte(nil), rec.RequestHash...)
	rec.ResponseBody = append([]byte(nil), rec.ResponseBody...)
	m.rows[k] = rec
	return nil
}
