package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/progression"
)

// ServiceConfig tunes the activity service. Clock injection keeps tests
// deterministic across time-based branches.
type ServiceConfig struct {
	Now func() time.Time
}

// DefaultServiceConfig returns baseline values.
func DefaultServiceConfig() ServiceConfig { return ServiceConfig{Now: time.Now} }

// CharacterFetcher is the read capability the service needs from the character
// domain — kept as an interface so tests can swap.
type CharacterFetcher interface {
	GetOwned(ctx context.Context, characterID, accountID uuid.UUID) (*character.Character, error)
}

// GrantedItem is the compact reference to an item minted for a reveal.
// Populated from the item service (Sprint 8). Kept as a plain struct here so
// the activity package doesn't import item (which would cycle — item already
// imports activity for Kind).
type GrantedItem struct {
	InstanceID   uuid.UUID `json:"instance_id"`
	TemplateCode string    `json:"template_code"`
	Rarity       string    `json:"rarity"`
	Prefix       *string   `json:"prefix,omitempty"`
	Suffix       *string   `json:"suffix,omitempty"`
}

// ItemGranter is the item-minting capability the activity service needs.
// Implementations MUST be idempotent by sourceRefID — repeated calls with the
// same source (i.e. reveal replays) return the previously created list without
// creating duplicates.
type ItemGranter interface {
	Grant(
		ctx context.Context,
		characterID, sourceRefID uuid.UUID,
		seed int64, kind Kind,
		commonCount, rareCount, legendaryCount int,
	) ([]GrantedItem, error)
}

// Service is the activity business logic.
type Service struct {
	repo     Repository
	chars    CharacterFetcher
	txn      economy.TxnService
	items    ItemGranter // optional — nil during transition sprints
	progress progression.Recorder
	cfg      ServiceConfig
	log      *slog.Logger
}

func NewService(repo Repository, chars CharacterFetcher, txn economy.TxnService, cfg ServiceConfig, logger *slog.Logger) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{repo: repo, chars: chars, txn: txn, cfg: cfg, log: logger}
}

// WithRecorder feeds reveal yields into progression (Sprint 24).
func (s *Service) WithRecorder(r progression.Recorder) *Service {
	s.progress = r
	return s
}

// WithItemGranter attaches an item-minting service so reveal drops become real
// item_instance rows. Optional — before Sprint 8 the reveal returns only gold.
func (s *Service) WithItemGranter(g ItemGranter) *Service {
	s.items = g
	return s
}

// --- Requests / responses ---------------------------------------------------

type StartRequest struct {
	CharacterID     uuid.UUID
	AccountID       uuid.UUID
	Kind            Kind
	DurationSeconds int32
	Params          json.RawMessage
}

// StartResult is the "you started an activity" response — includes preview of
// what the max-duration yield would look like (client HUD).
type StartResult struct {
	Contract      *Contract    `json:"-"`               // response layer flattens the fields it wants
	FullYieldPeek YieldSummary `json:"full_yield_peek"` // yield IF fully completed at expires_at
}

// RevealResult wraps the yield summary + the ledger txn id (for audit trace).
type RevealResult struct {
	Yield       YieldSummary `json:"yield"`
	LedgerTxnID uuid.UUID    `json:"ledger_txn_id"`
	Replayed    bool         `json:"replayed"`
}

// --- Start ------------------------------------------------------------------

func (s *Service) Start(ctx context.Context, req StartRequest) (*StartResult, error) {
	if !KnownKinds[req.Kind] {
		return nil, &ValidationError{Field: "kind", Message: "unknown activity kind"}
	}
	if req.DurationSeconds < MinDurationSeconds || req.DurationSeconds > MaxDurationSeconds {
		return nil, &ValidationError{
			Field:   "duration_seconds",
			Message: fmt.Sprintf("must be between %d and %d", MinDurationSeconds, MaxDurationSeconds),
		}
	}

	// Ownership check (existence-leak-safe — returns ErrForbidden for both
	// not-found and foreign-owned).
	char, err := s.chars.GetOwned(ctx, req.CharacterID, req.AccountID)
	if err != nil {
		return nil, ErrForbidden
	}
	if char.State != character.StateActive {
		return nil, &ValidationError{Field: "character", Message: "character not playable"}
	}

	now := s.cfg.Now().UTC()
	id := uuid.New()
	snapshot := TakeSnapshot(char)

	contract := Contract{
		ID:              id,
		CharacterID:     req.CharacterID,
		Kind:            req.Kind,
		DurationSeconds: req.DurationSeconds,
		StartedAt:       now,
		ExpiresAt:       now.Add(time.Duration(req.DurationSeconds) * time.Second),
		Seed:            SeedFromContract(id, req.CharacterID),
		State:           StateActive,
		Snapshot:        MarshalSnapshot(snapshot),
		Params:          nonNilJSON(req.Params),
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         1,
	}

	if err := s.repo.Create(ctx, contract); err != nil {
		return nil, err
	}

	// Full-duration yield preview (deterministic — client can display "expected
	// ~180 iron ore, ~12 rare drops" at start).
	preview := ComputeYield(&contract, contract.DurationSeconds)

	return &StartResult{
		Contract:      &contract,
		FullYieldPeek: preview,
	}, nil
}

// --- Reveal -----------------------------------------------------------------

// Reveal is idempotent — repeated calls return the cached response (state
// stays 'revealed'; no new ledger writes). Partial reveal is supported: if
// the user reveals before expires_at, they get yield proportional to elapsed.
func (s *Service) Reveal(ctx context.Context, contractID, accountID uuid.UUID) (*RevealResult, error) {
	contract, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}

	// Ownership.
	if _, err := s.chars.GetOwned(ctx, contract.CharacterID, accountID); err != nil {
		return nil, ErrForbidden
	}

	// Idempotency: already revealed → return cached.
	if contract.State == StateRevealed {
		var y YieldSummary
		_ = json.Unmarshal(contract.RevealResult, &y)
		txn := uuid.Nil
		if contract.RevealTxnID != nil {
			txn = *contract.RevealTxnID
		}
		return &RevealResult{Yield: y, LedgerTxnID: txn, Replayed: true}, nil
	}

	if contract.State != StateActive {
		// cancelled → nothing to reveal.
		return nil, ErrNotActive
	}

	now := s.cfg.Now().UTC()
	elapsed := contract.ElapsedSeconds(now)
	yield := ComputeYield(contract, elapsed)

	// Materialize item drops (idempotent by source_ref_id = contract.ID). Done
	// BEFORE ledger write so a partial failure (items minted, ledger failed)
	// leaves contract active; retry re-uses the same items instead of
	// double-minting. items may be nil during transition sprints.
	if s.items != nil {
		granted, err := s.items.Grant(ctx,
			contract.CharacterID, contract.ID,
			contract.Seed, contract.Kind,
			yield.CommonCount, yield.RareCount, yield.LegendaryCount,
		)
		if err != nil {
			return nil, fmt.Errorf("reveal item grant: %w", err)
		}
		yield.Items = granted
	}

	// Persist ledger write BEFORE marking revealed — if TxnService fails, the
	// contract stays active and the user can retry. If ledger succeeds but our
	// mark-revealed fails, the next reveal call will attempt ledger again with
	// the same idempotency_key and get a Replayed=true (no double-credit).
	var txnID uuid.UUID
	if yield.GoldCopper > 0 {
		res, err := s.txn.Write(ctx, economy.Intent{
			Kind:           "activity_reveal",
			RefKind:        "activity_contract",
			RefID:          contract.ID,
			IdempotencyKey: "activity-reveal-" + contract.ID.String(),
			Entries: []economy.Entry{
				{
					Currency: "gold", AccountKind: economy.AccountSystem,
					AccountRef: economy.SystemRewardSource, Delta: -yield.GoldCopper,
				},
				{
					Currency: "gold", AccountKind: economy.AccountCharacter,
					AccountRef: contract.CharacterID, Delta: yield.GoldCopper,
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("reveal ledger write: %w", err)
		}
		txnID = res.TxnID
	}

	// Persist cache.
	if err := s.repo.MarkRevealed(ctx, contract.ID, MarshalSummary(yield), txnID, now); err != nil {
		// Race: someone else revealed between our reads. Refetch + return cached.
		refreshed, gerr := s.repo.GetByID(ctx, contract.ID)
		if gerr == nil && refreshed.State == StateRevealed {
			var y YieldSummary
			_ = json.Unmarshal(refreshed.RevealResult, &y)
			cached := uuid.Nil
			if refreshed.RevealTxnID != nil {
				cached = *refreshed.RevealTxnID
			}
			return &RevealResult{Yield: y, LedgerTxnID: cached, Replayed: true}, nil
		}
		return nil, err
	}

	progression.RecordIf(ctx, s.progress, contract.CharacterID, progression.Event{
		Kind:   progression.EventGatherPrefix + string(contract.Kind),
		Amount: int64(yield.CommonCount + yield.RareCount + yield.LegendaryCount),
	})
	return &RevealResult{Yield: yield, LedgerTxnID: txnID, Replayed: false}, nil
}

// --- Stop -------------------------------------------------------------------

// Stop cancels an active contract. No reward — caller forfeits any accrued
// yield. Idempotent (returns nil for already-cancelled contracts; returns
// ErrNotActive for revealed contracts).
func (s *Service) Stop(ctx context.Context, contractID, accountID uuid.UUID) (*Contract, error) {
	contract, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if _, err := s.chars.GetOwned(ctx, contract.CharacterID, accountID); err != nil {
		return nil, ErrForbidden
	}

	if contract.State == StateCancelled {
		return contract, nil // already stopped — idempotent
	}
	if contract.State != StateActive {
		return nil, ErrNotActive
	}

	now := s.cfg.Now().UTC()
	if err := s.repo.MarkCancelled(ctx, contract.ID, now); err != nil {
		return nil, err
	}
	contract.State = StateCancelled
	c := now
	contract.CancelledAt = &c
	return contract, nil
}

// --- Read helpers -----------------------------------------------------------

// Active returns the character's current active contract with a live yield
// preview at the current wall-clock time. Returns ErrNotFound when idle.
func (s *Service) Active(ctx context.Context, characterID, accountID uuid.UUID) (*Contract, *YieldSummary, error) {
	// Ownership.
	if _, err := s.chars.GetOwned(ctx, characterID, accountID); err != nil {
		return nil, nil, ErrForbidden
	}
	c, err := s.repo.ActiveByCharacter(ctx, characterID)
	if err != nil {
		return nil, nil, err
	}
	elapsed := c.ElapsedSeconds(s.cfg.Now().UTC())
	preview := ComputeYield(c, elapsed)
	return c, &preview, nil
}

// History returns the character's most recent contracts (any state).
func (s *Service) History(ctx context.Context, characterID, accountID uuid.UUID, limit int) ([]Contract, error) {
	if _, err := s.chars.GetOwned(ctx, characterID, accountID); err != nil {
		return nil, ErrForbidden
	}
	return s.repo.ListByCharacter(ctx, characterID, limit)
}

// --- helpers ----------------------------------------------------------------

func nonNilJSON(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage(`{}`)
	}
	return raw
}
