package vendor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Repository is the persistence contract for vendor templates and buybacks.
// Vendor templates load once at boot (like item.Catalog).
type Repository interface {
	// Catalog loads all vendor templates (called at boot).
	Catalog(ctx context.Context) ([]Template, error)

	// InsertBuyback persists a fresh buyback entry.
	InsertBuyback(ctx context.Context, b BuybackEntry) error

	// GetBuyback returns a buyback for restore. Ownership + expiry check is
	// the service's job.
	GetBuyback(ctx context.Context, id uuid.UUID) (*BuybackEntry, error)

	// ClaimBuyback atomically marks a buyback claimed (state transition
	// active → claimed). Fails on already-claimed / expired.
	ClaimBuyback(ctx context.Context, id uuid.UUID, at time.Time) error

	// ListBuybacks returns character's active (unclaimed + not expired)
	// buybacks, newest first.
	ListBuybacks(ctx context.Context, characterID uuid.UUID, limit int) ([]BuybackEntry, error)
}

// --- MemoryRepo (tests) ------------------------------------------------------

type MemoryRepo struct {
	mu        sync.Mutex
	templates map[string]Template
	buybacks  map[uuid.UUID]BuybackEntry
}

func NewMemoryRepo(templates []Template) *MemoryRepo {
	m := make(map[string]Template, len(templates))
	for _, t := range templates {
		m[t.Code] = t
	}
	return &MemoryRepo{templates: m, buybacks: make(map[uuid.UUID]BuybackEntry)}
}

func (r *MemoryRepo) Catalog(_ context.Context) ([]Template, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Template, 0, len(r.templates))
	for _, t := range r.templates {
		out = append(out, t)
	}
	return out, nil
}

func (r *MemoryRepo) InsertBuyback(_ context.Context, b BuybackEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := b
	cp.ItemSnapshot = append(json.RawMessage(nil), b.ItemSnapshot...)
	r.buybacks[cp.ID] = cp
	return nil
}

func (r *MemoryRepo) GetBuyback(_ context.Context, id uuid.UUID) (*BuybackEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.buybacks[id]
	if !ok {
		return nil, ErrBuybackNotFound
	}
	cp := b
	return &cp, nil
}

func (r *MemoryRepo) ClaimBuyback(_ context.Context, id uuid.UUID, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.buybacks[id]
	if !ok {
		return ErrBuybackNotFound
	}
	if b.ClaimedAt != nil {
		return ErrBuybackClaimed
	}
	if at.After(b.ExpiresAt) {
		return ErrBuybackExpired
	}
	claim := at
	b.ClaimedAt = &claim
	r.buybacks[id] = b
	return nil
}

func (r *MemoryRepo) ListBuybacks(_ context.Context, characterID uuid.UUID, limit int) ([]BuybackEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	out := []BuybackEntry{}
	for _, b := range r.buybacks {
		if b.CharacterID != characterID || b.ClaimedAt != nil || now.After(b.ExpiresAt) {
			continue
		}
		out = append(out, b)
	}
	// Sort by SoldAt desc.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].SoldAt.After(out[i].SoldAt) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// Silence unused warnings.
var _ = fmt.Sprintf
var _ = errors.Is
