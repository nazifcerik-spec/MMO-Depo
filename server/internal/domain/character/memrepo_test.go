package character

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// --- MemoryRepo ------------------------------------------------------

type MemoryRepo struct {
	mu    sync.Mutex
	rows  map[uuid.UUID]Character
	names map[string]uuid.UUID // name_normalized → id
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		rows:  make(map[uuid.UUID]Character),
		names: make(map[string]uuid.UUID),
	}
}

func (r *MemoryRepo) Create(_ context.Context, c Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.names[c.NameNormalized]; exists {
		return ErrNameTaken
	}
	r.rows[c.ID] = c
	r.names[c.NameNormalized] = c.ID
	return nil
}

func (r *MemoryRepo) Get(_ context.Context, id uuid.UUID) (*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.rows[id]; ok {
		cp := c
		return &cp, nil
	}
	return nil, ErrNotFound
}

func (r *MemoryRepo) ListByAccount(_ context.Context, accountID uuid.UUID) ([]Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []Character{}
	for _, c := range r.rows {
		if c.AccountID == accountID && (c.State == StateActive || c.State == StateDeleteGrace) {
			out = append(out, c)
		}
	}
	return out, nil
}

func (r *MemoryRepo) CountActiveByAccount(_ context.Context, accountID uuid.UUID) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, c := range r.rows {
		if c.AccountID == accountID && (c.State == StateActive || c.State == StateDeleteGrace) {
			n++
		}
	}
	return n, nil
}

func (r *MemoryRepo) Rename(_ context.Context, id uuid.UUID, name, normalized string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	if existing, taken := r.names[normalized]; taken && existing != id {
		return ErrNameTaken
	}
	delete(r.names, c.NameNormalized)
	c.Name = name
	c.NameNormalized = normalized
	c.UpdatedAt = time.Now().UTC()
	c.Version++
	r.rows[id] = c
	r.names[normalized] = id
	return nil
}

func (r *MemoryRepo) TouchLastPlayed(_ context.Context, id uuid.UUID, ts time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok || c.State != StateActive {
		return nil
	}
	c.LastPlayedAt = &ts
	r.rows[id] = c
	return nil
}

func (r *MemoryRepo) SetDeleteGrace(_ context.Context, id uuid.UUID, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	if c.State != StateActive {
		return ErrNotPlayable
	}
	c.State = StateDeleteGrace
	c.DeleteGraceAt = &expiresAt
	c.UpdatedAt = time.Now().UTC()
	c.Version++
	r.rows[id] = c
	return nil
}

// AddXP is the in-memory implementation of the level-rollover behaviour.
// The Postgres impl uses a transaction to make xp update + level-up atomic.
func (r *MemoryRepo) AddXP(_ context.Context, id uuid.UUID, delta int64) (*Character, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if delta < 0 {
		return nil, false, ErrInvalidXP
	}
	c, ok := r.rows[id]
	if !ok {
		return nil, false, ErrNotFound
	}
	leveledUp := false
	c.XP += delta
	for {
		needed := XPToNext(int(c.Level))
		if c.XP < needed || c.Level >= LevelCap {
			break
		}
		c.XP -= needed
		c.Level++
		leveledUp = true
	}
	c.UpdatedAt = time.Now().UTC()
	c.Version++
	r.rows[id] = c
	cp := c
	return &cp, leveledUp, nil
}

// SetLocation / GetLocation — in-memory. LocationRegionCode lives alongside
// the Character; we widen Character with a text field when we add it (see
// character.go) and just mirror it here.
func (r *MemoryRepo) SetLocation(_ context.Context, id uuid.UUID, regionCode string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	c.LocationRegionCode = regionCode
	c.UpdatedAt = time.Now().UTC()
	c.Version++
	r.rows[id] = c
	return nil
}

func (r *MemoryRepo) GetLocation(_ context.Context, id uuid.UUID) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return "", ErrNotFound
	}
	return c.LocationRegionCode, nil
}

func (r *MemoryRepo) CancelDeleteGrace(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	if c.State != StateDeleteGrace {
		return ErrNotInDeleteGrace
	}
	c.State = StateActive
	c.DeleteGraceAt = nil
	c.UpdatedAt = time.Now().UTC()
	c.Version++
	r.rows[id] = c
	return nil
}

func (r *MemoryRepo) SetHP(_ context.Context, id uuid.UUID, hp int32, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.rows[id]
	if !ok {
		return ErrNotFound
	}
	c.HPStored, c.HPUpdatedAt = &hp, &at
	r.rows[id] = c
	return nil
}
