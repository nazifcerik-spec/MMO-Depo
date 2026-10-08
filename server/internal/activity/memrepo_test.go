package activity

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

// --- MemoryRepo ------------------------------------------------------

type MemoryRepo struct {
	mu   sync.Mutex
	rows map[uuid.UUID]Contract
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{rows: make(map[uuid.UUID]Contract)}
}

func (r *MemoryRepo) Create(_ context.Context, c Contract) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Enforce concurrent slot cap (matches the partial UNIQUE index in Postgres).
	for _, existing := range r.rows {
		if existing.CharacterID == c.CharacterID && existing.State == StateActive {
			return ErrSlotBusy
		}
	}
	r.rows[c.ID] = c
	return nil
}

func (r *MemoryRepo) GetByID(_ context.Context, id uuid.UUID) (*Contract, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := c
	return &cp, nil
}

func (r *MemoryRepo) ActiveByCharacter(_ context.Context, characterID uuid.UUID) (*Contract, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.rows {
		if c.CharacterID == characterID && c.State == StateActive {
			cp := c
			return &cp, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MemoryRepo) MarkRevealed(_ context.Context, id uuid.UUID, revealResult json.RawMessage, txnID uuid.UUID, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	if c.State != StateActive {
		return ErrNotActive
	}
	c.State = StateRevealed
	c.RevealResult = append(json.RawMessage(nil), revealResult...)
	tid := txnID
	c.RevealTxnID = &tid
	revAt := at
	c.RevealedAt = &revAt
	c.UpdatedAt = at
	c.Version++
	r.rows[id] = c
	return nil
}

func (r *MemoryRepo) MarkCancelled(_ context.Context, id uuid.UUID, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	if c.State != StateActive {
		return ErrNotActive
	}
	c.State = StateCancelled
	canAt := at
	c.CancelledAt = &canAt
	c.UpdatedAt = at
	c.Version++
	r.rows[id] = c
	return nil
}

func (r *MemoryRepo) ListByCharacter(_ context.Context, characterID uuid.UUID, limit int) ([]Contract, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []Contract
	for _, c := range r.rows {
		if c.CharacterID == characterID {
			out = append(out, c)
		}
	}
	// Sort by StartedAt desc.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].StartedAt.After(out[i].StartedAt) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	if out == nil {
		out = []Contract{}
	}
	return out, nil
}

// Ensure sentinel imports are used (avoids linter unused warning in some
// downstream builds).
