package world

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/domain/character"
)

// --- Ports ------------------------------------------------------------------

// CatalogLookup exposes the region catalog. CatalogService satisfies.
type CatalogLookup interface {
	Catalog() *Catalog
}

// CharacterBinder is the tiny character-side surface world needs: ownership
// check, level lookup, current location, and location writer. character.Service
// satisfies via GetOwnedInternal, Owns, GetLocation, SetLocation.
type CharacterBinder interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
	GetOwnedInternal(ctx context.Context, characterID uuid.UUID) (*character.Character, error)
	GetLocation(ctx context.Context, characterID uuid.UUID) (string, error)
	SetLocation(ctx context.Context, characterID uuid.UUID, regionCode string) error
}

// ServiceConfig — clock injection.
type ServiceConfig struct {
	Now func() time.Time
}

func DefaultServiceConfig() ServiceConfig { return ServiceConfig{Now: time.Now} }

// Service is the world / travel business logic.
type Service struct {
	repo    TravelRepository
	catalog CatalogLookup
	chars   CharacterBinder
	cfg     ServiceConfig
	log     *slog.Logger
}

func NewService(repo TravelRepository, catalog CatalogLookup, chars CharacterBinder, cfg ServiceConfig, logger *slog.Logger) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{repo: repo, catalog: catalog, chars: chars, cfg: cfg, log: logger}
}

// --- Location read ----------------------------------------------------------

// LocationView is the compound response for GET /v1/character/{id}/location.
type LocationView struct {
	Region       Region         `json:"region"`
	Neighbors    []NeighborInfo `json:"neighbors"`
	HuntRegions  []string       `json:"hunt_regions"`
	Travel       *Travel        `json:"travel,omitempty"` // active in-progress travel, if any
	StartingCode string         `json:"starting_code,omitempty"`
}

// CurrentRegionCode resolves the character's current region — either the
// persisted location_region_code, or the world catalog's starting region as
// fallback. Returns ErrNoStartingRegion if neither is set (broken content).
func (s *Service) CurrentRegionCode(ctx context.Context, characterID uuid.UUID) (string, error) {
	code, err := s.chars.GetLocation(ctx, characterID)
	if err != nil {
		return "", err
	}
	if code == "" {
		code = s.catalog.Catalog().StartingCode()
		if code == "" {
			return "", ErrNoStartingRegion
		}
	}
	return code, nil
}

// HuntRegions lists the regions the character can hunt in from where it
// stands (DD-273).
func (s *Service) HuntRegions(ctx context.Context, characterID uuid.UUID) ([]string, error) {
	code, err := s.CurrentRegionCode(ctx, characterID)
	if err != nil {
		return nil, err
	}
	return s.catalog.Catalog().HuntRegions(code), nil
}

// Location returns the character's current region + neighbors + active travel
// (nil when not travelling).
func (s *Service) Location(ctx context.Context, characterID, accountID uuid.UUID) (*LocationView, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	cat := s.catalog.Catalog()
	code, err := s.CurrentRegionCode(ctx, characterID)
	if err != nil {
		return nil, err
	}
	region, ok := cat.Get(code)
	if !ok {
		// Content shrunk under us — return a stub rather than 500 so the
		// client can prompt the user to travel to a valid region via the
		// starting-region fallback.
		return nil, ErrUnknownRegion
	}
	view := &LocationView{
		Region:       region,
		Neighbors:    cat.Neighbors(code),
		HuntRegions:  cat.HuntRegions(code),
		StartingCode: cat.StartingCode(),
	}
	active, err := s.repo.Active(ctx, characterID)
	if err == nil {
		view.Travel = active
	}
	return view, nil
}

// --- Travel start -----------------------------------------------------------

type StartRequest struct {
	CharacterID  uuid.UUID
	AccountID    uuid.UUID
	ToRegionCode string
}

type StartResult struct {
	Travel   *Travel `json:"travel"`
	FromCode string  `json:"from_region_code"`
	Duration int32   `json:"duration_seconds"`
}

// Start begins a travel from the character's current region to `toRegionCode`.
// Enforces:
//   - ownership (via CharacterBinder.Owns)
//   - destination must be a direct neighbor (parent, child, or explicit link)
//   - destination's min_level <= character level
//   - no other travel already in progress (partial UNIQUE enforces at DB tier)
//   - source != destination
func (s *Service) Start(ctx context.Context, req StartRequest) (*StartResult, error) {
	if req.CharacterID == uuid.Nil || req.AccountID == uuid.Nil {
		return nil, ErrForbidden
	}
	if req.ToRegionCode == "" {
		return nil, ErrValidation
	}
	if err := s.chars.Owns(ctx, req.CharacterID, req.AccountID); err != nil {
		return nil, err
	}
	cat := s.catalog.Catalog()

	toRegion, ok := cat.Get(req.ToRegionCode)
	if !ok {
		return nil, ErrUnknownRegion
	}
	fromCode, err := s.CurrentRegionCode(ctx, req.CharacterID)
	if err != nil {
		return nil, err
	}
	if fromCode == req.ToRegionCode {
		return nil, ErrSameRegion
	}
	dur, ok := cat.TravelDuration(fromCode, req.ToRegionCode)
	if !ok {
		return nil, ErrNoRoute
	}
	// Level gate — dungeons and high-level zones stay unreachable until the
	// character has crossed the threshold.
	c, err := s.chars.GetOwnedInternal(ctx, req.CharacterID)
	if err != nil {
		return nil, err
	}
	if c.Level < toRegion.MinLevel {
		return nil, ErrLevelRequirement
	}

	now := s.cfg.Now().UTC()
	t := &Travel{
		ID:             uuid.New(),
		CharacterID:    req.CharacterID,
		FromRegionCode: fromCode,
		ToRegionCode:   req.ToRegionCode,
		StartedAt:      now,
		ArrivesAt:      now.Add(time.Duration(dur) * time.Second),
		State:          TravelInProgress,
		CreatedAt:      now,
		UpdatedAt:      now,
		Version:        1,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return &StartResult{Travel: t, FromCode: fromCode, Duration: dur}, nil
}

// --- Travel complete --------------------------------------------------------

// Complete confirms arrival — refuses if `arrives_at` is still in the future.
// On success, sets character.location_region_code to the destination and marks
// travel completed. Idempotent for already-completed travels (returns the
// row as-is; no double write).
func (s *Service) Complete(ctx context.Context, travelID, accountID uuid.UUID) (*Travel, error) {
	t, err := s.repo.GetByID(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if err := s.chars.Owns(ctx, t.CharacterID, accountID); err != nil {
		return nil, err
	}
	if t.State == TravelCompleted {
		return t, nil
	}
	if t.State != TravelInProgress {
		return nil, ErrTravelNotActive
	}
	now := s.cfg.Now().UTC()
	if now.Before(t.ArrivesAt) {
		return nil, ErrTravelNotReady
	}
	// Write character location BEFORE marking travel completed — if the write
	// fails the travel stays in-progress and retry will re-attempt. If mark
	// fails after location write, the next call sees location already set but
	// state still active; MarkCompleted's rows-affected guard flips it once.
	if err := s.chars.SetLocation(ctx, t.CharacterID, t.ToRegionCode); err != nil {
		return nil, err
	}
	if err := s.repo.MarkCompleted(ctx, t.ID, now); err != nil {
		if err == ErrTravelNotActive {
			// Refetch — most likely already completed by a parallel call.
			refreshed, gerr := s.repo.GetByID(ctx, t.ID)
			if gerr == nil && refreshed.State == TravelCompleted {
				return refreshed, nil
			}
		}
		return nil, err
	}
	t.State = TravelCompleted
	t.CompletedAt = &now
	return t, nil
}

// --- Travel cancel ----------------------------------------------------------

// Cancel aborts an in-progress travel. Character stays at from_region_code
// (no write needed since Complete is the only path that changes location).
// Idempotent for already-cancelled.
func (s *Service) Cancel(ctx context.Context, travelID, accountID uuid.UUID) (*Travel, error) {
	t, err := s.repo.GetByID(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if err := s.chars.Owns(ctx, t.CharacterID, accountID); err != nil {
		return nil, err
	}
	if t.State == TravelCancelled {
		return t, nil
	}
	if t.State != TravelInProgress {
		return nil, ErrTravelNotActive
	}
	now := s.cfg.Now().UTC()
	if err := s.repo.MarkCancelled(ctx, t.ID, now); err != nil {
		return nil, err
	}
	t.State = TravelCancelled
	t.CancelledAt = &now
	return t, nil
}

// --- Active / History ------------------------------------------------------

func (s *Service) History(ctx context.Context, characterID, accountID uuid.UUID, limit int) ([]Travel, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	return s.repo.History(ctx, characterID, limit)
}
