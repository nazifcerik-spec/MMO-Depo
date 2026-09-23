package auth

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Repository is the persistence contract for the auth domain. One aggregated
// interface (rather than three fine-grained) keeps the service wiring simple;
// consumers can still adapt / spy specific methods.
type Repository interface {
	// --- Account
	CreateAccount(ctx context.Context, a Account) error
	AccountByEmailNormalized(ctx context.Context, emailN string) (*Account, error)
	AccountByUsernameNormalized(ctx context.Context, userN string) (*Account, error)
	AccountByID(ctx context.Context, id uuid.UUID) (*Account, error)
	UpdateAccountLastLogin(ctx context.Context, id uuid.UUID, ts time.Time) error

	// --- Session
	CreateSession(ctx context.Context, s Session) error
	Session(ctx context.Context, id uuid.UUID) (*Session, error)
	RevokeSession(ctx context.Context, id uuid.UUID, reason string, at time.Time) error
	TouchSession(ctx context.Context, id uuid.UUID, at time.Time) error

	// --- Refresh token
	CreateRefreshToken(ctx context.Context, t RefreshToken) error
	RefreshTokenByHash(ctx context.Context, hash []byte) (*RefreshToken, error)
	MarkRefreshTokenUsed(ctx context.Context, id uuid.UUID, at time.Time) error
	RevokeRefreshTokensForSession(ctx context.Context, sessionID uuid.UUID, reason string, at time.Time) error
}

// --- MemoryRepo (tests) ------------------------------------------------------

// MemoryRepo is a goroutine-safe in-memory Repository. Test-only.
type MemoryRepo struct {
	mu       sync.Mutex
	accounts map[uuid.UUID]Account
	sessions map[uuid.UUID]Session
	tokens   map[uuid.UUID]RefreshToken
	tokIdx   map[string]uuid.UUID // sha-hex → token id
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		accounts: make(map[uuid.UUID]Account),
		sessions: make(map[uuid.UUID]Session),
		tokens:   make(map[uuid.UUID]RefreshToken),
		tokIdx:   make(map[string]uuid.UUID),
	}
}

func (r *MemoryRepo) CreateAccount(_ context.Context, a Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.accounts {
		if existing.EmailNormalized == a.EmailNormalized || existing.UsernameNormalized == a.UsernameNormalized {
			return ErrAccountExists
		}
	}
	r.accounts[a.ID] = a
	return nil
}

func (r *MemoryRepo) AccountByEmailNormalized(_ context.Context, emailN string) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, a := range r.accounts {
		if a.EmailNormalized == emailN {
			cp := a
			return &cp, nil
		}
	}
	return nil, ErrInvalidCredentials
}

func (r *MemoryRepo) AccountByUsernameNormalized(_ context.Context, userN string) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, a := range r.accounts {
		if a.UsernameNormalized == userN {
			cp := a
			return &cp, nil
		}
	}
	return nil, ErrInvalidCredentials
}

func (r *MemoryRepo) AccountByID(_ context.Context, id uuid.UUID) (*Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a, ok := r.accounts[id]; ok {
		cp := a
		return &cp, nil
	}
	return nil, ErrInvalidCredentials
}

func (r *MemoryRepo) UpdateAccountLastLogin(_ context.Context, id uuid.UUID, ts time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a, ok := r.accounts[id]; ok {
		a.LastLoginAt = &ts
		a.UpdatedAt = ts
		r.accounts[id] = a
	}
	return nil
}

func (r *MemoryRepo) CreateSession(_ context.Context, s Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID] = s
	return nil
}

func (r *MemoryRepo) Session(_ context.Context, id uuid.UUID) (*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[id]; ok {
		cp := s
		return &cp, nil
	}
	return nil, ErrSessionRevoked
}

func (r *MemoryRepo) RevokeSession(_ context.Context, id uuid.UUID, reason string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[id]; ok {
		s.RevokedAt = &at
		s.RevokeReason = reason
		r.sessions[id] = s
	}
	return nil
}

func (r *MemoryRepo) TouchSession(_ context.Context, id uuid.UUID, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[id]; ok {
		s.LastActiveAt = at
		r.sessions[id] = s
	}
	return nil
}

func hashKey(hash []byte) string {
	// Fast stringification for map key; not a security operation.
	const hex = "0123456789abcdef"
	out := make([]byte, len(hash)*2)
	for i, b := range hash {
		out[i*2] = hex[b>>4]
		out[i*2+1] = hex[b&0x0f]
	}
	return string(out)
}

func (r *MemoryRepo) CreateRefreshToken(_ context.Context, t RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[t.ID] = t
	r.tokIdx[hashKey(t.TokenHash)] = t.ID
	return nil
}

func (r *MemoryRepo) RefreshTokenByHash(_ context.Context, hash []byte) (*RefreshToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.tokIdx[hashKey(hash)]
	if !ok {
		return nil, ErrRefreshTokenNotFound
	}
	t := r.tokens[id]
	cp := t
	return &cp, nil
}

func (r *MemoryRepo) MarkRefreshTokenUsed(_ context.Context, id uuid.UUID, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.tokens[id]; ok {
		t.UsedAt = &at
		r.tokens[id] = t
	}
	return nil
}

func (r *MemoryRepo) RevokeRefreshTokensForSession(_ context.Context, sessionID uuid.UUID, reason string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, t := range r.tokens {
		if t.SessionID == sessionID && t.RevokedAt == nil {
			t.RevokedAt = &at
			t.RevokeReason = reason
			r.tokens[id] = t
		}
	}
	return nil
}
