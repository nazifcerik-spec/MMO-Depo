package item

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// --- MemoryRepo ------------------------------------------------------

type MemoryRepo struct {
	mu   sync.Mutex
	rows map[uuid.UUID]Instance
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{rows: make(map[uuid.UUID]Instance)}
}

func (r *MemoryRepo) CreateInstancesBatch(_ context.Context, instances []Instance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, i := range instances {
		r.rows[i.ID] = i
	}
	return nil
}

func (r *MemoryRepo) Get(_ context.Context, id uuid.UUID) (*Instance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if i, ok := r.rows[id]; ok {
		cp := i
		return &cp, nil
	}
	return nil, ErrNotFound
}

func (r *MemoryRepo) ListByOwner(_ context.Context, characterID uuid.UUID, limit int) ([]Instance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []Instance{}
	for _, i := range r.rows {
		if i.OwnerCharacterID != nil && *i.OwnerCharacterID == characterID {
			out = append(out, i)
		}
	}
	// Newest first.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].CreatedAt.After(out[i].CreatedAt) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *MemoryRepo) ListBySource(_ context.Context, sourceRefID uuid.UUID) ([]Instance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []Instance{}
	for _, i := range r.rows {
		if i.SourceRefID != nil && *i.SourceRefID == sourceRefID {
			out = append(out, i)
		}
	}
	return out, nil
}

// Unused import guard.
