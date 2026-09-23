package character

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// ServiceConfig tunes the character service.
type ServiceConfig struct {
	SlotCap     int           // DD-OPEN-8, v1 default 2
	DeleteGrace time.Duration // DD-134 pattern, default 48h
	Now         func() time.Time
}

// DefaultServiceConfig returns Faz 00 baseline.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		SlotCap:     2,
		DeleteGrace: 48 * time.Hour,
		Now:         time.Now,
	}
}

// SheetComputer overrides the placeholder Sheet computation with real modifier
// stack (equipped items, buffs, mastery bonuses). Sprint 11 skill.SheetComputer
// satisfies. Optional — nil = fall back to class-base placeholder.
type SheetComputer interface {
	ComputeSheet(ctx context.Context, characterID uuid.UUID, base Sheet) (Sheet, error)
}

// Service is the character business logic.
type Service struct {
	repo    Repository
	sheeter SheetComputer
	cfg     ServiceConfig
	log     *slog.Logger
}

func NewService(repo Repository, cfg ServiceConfig, logger *slog.Logger) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.SlotCap <= 0 {
		cfg.SlotCap = 2
	}
	if cfg.DeleteGrace <= 0 {
		cfg.DeleteGrace = 48 * time.Hour
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{repo: repo, cfg: cfg, log: logger}
}

// WithSheetComputer wires an external stat-computer (Sprint 11 skill engine).
func (s *Service) WithSheetComputer(sc SheetComputer) *Service {
	s.sheeter = sc
	return s
}

// --- Requests ---------------------------------------------------------------

type CreateRequest struct {
	AccountID uuid.UUID
	Name      string
	RaceCode  string
	ClassCode string
}

// --- Create -----------------------------------------------------------------

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Character, error) {
	if req.AccountID == uuid.Nil {
		return nil, ErrForbidden
	}
	name, normalized, err := validateName(req.Name)
	if err != nil {
		return nil, err
	}
	race, err := validateRace(req.RaceCode)
	if err != nil {
		return nil, err
	}
	class, err := validateClass(req.ClassCode)
	if err != nil {
		return nil, err
	}

	// Slot cap check.
	count, err := s.repo.CountActiveByAccount(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}
	if count >= s.cfg.SlotCap {
		return nil, ErrSlotCap
	}

	now := s.cfg.Now().UTC()
	c := Character{
		ID:             uuid.New(),
		AccountID:      req.AccountID,
		Name:           name,
		NameNormalized: normalized,
		RaceCode:       race,
		ClassCode:      class,
		Level:          1,
		XP:             0,
		State:          StateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		Version:        1,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &c, nil
}

// --- List / Get -------------------------------------------------------------

func (s *Service) ListForAccount(ctx context.Context, accountID uuid.UUID) ([]Character, error) {
	if accountID == uuid.Nil {
		return nil, ErrForbidden
	}
	return s.repo.ListByAccount(ctx, accountID)
}

// GetOwned fetches the character and asserts ownership. Returns ErrForbidden
// when the character exists but belongs to another account (does not leak
// existence).
func (s *Service) GetOwned(ctx context.Context, characterID, accountID uuid.UUID) (*Character, error) {
	c, err := s.repo.Get(ctx, characterID)
	if err != nil {
		return nil, err
	}
	if c.AccountID != accountID {
		return nil, ErrForbidden
	}
	return c, nil
}

// GetOwnedInternal fetches the character without an ownership check — used by
// internal services (skill engine, equipment) that already have a trusted
// character id (e.g. from URL parsing that ownership middleware validated).
func (s *Service) GetOwnedInternal(ctx context.Context, characterID uuid.UUID) (*Character, error) {
	return s.repo.Get(ctx, characterID)
}

// --- Enter (session bind) ---------------------------------------------------

// Enter is called when the client selects a character. Verifies ownership and
// state, touches last_played_at, and returns the sheet snapshot.
func (s *Service) Enter(ctx context.Context, characterID, accountID uuid.UUID) (*Sheet, error) {
	c, err := s.GetOwned(ctx, characterID, accountID)
	if err != nil {
		return nil, err
	}
	if !c.IsPlayable() {
		return nil, ErrNotPlayable
	}
	_ = s.repo.TouchLastPlayed(ctx, c.ID, s.cfg.Now().UTC())
	return s.liveSheet(ctx, c), nil
}

// Sheet returns the compound character view without touching last_played_at.
// Used by GET endpoints.
func (s *Service) Sheet(ctx context.Context, characterID, accountID uuid.UUID) (*Sheet, error) {
	c, err := s.GetOwned(ctx, characterID, accountID)
	if err != nil {
		return nil, err
	}
	return s.liveSheet(ctx, c), nil
}

// liveSheet applies the modifier stack (equipped bonuses) and then resolves
// current HP against the final max (DD-218: stored HP + analytical regen).
func (s *Service) liveSheet(ctx context.Context, c *Character) *Sheet {
	sheet := buildSheet(c)
	if s.sheeter != nil {
		if computed, err := s.sheeter.ComputeSheet(ctx, c.ID, sheet); err == nil {
			sheet = computed
		} else {
			s.log.WarnContext(ctx, "character sheet compute failed — falling back to placeholder",
				slog.String("err", err.Error()))
		}
	}
	sheet.Resources.HPCurrent = EffectiveHP(c.HPStored, c.HPUpdatedAt, sheet.Resources.HPMax,
		sheet.Primary.Vitality, s.cfg.Now())
	return &sheet
}

// SetHP records current HP as of `at` (hunt/dungeon reveal, consumables).
// Internal: callers have already verified ownership.
func (s *Service) SetHP(ctx context.Context, characterID uuid.UUID, hp int32, at time.Time) error {
	return s.repo.SetHP(ctx, characterID, hp, at)
}

// --- Delete grace flow ------------------------------------------------------

// RequestDelete moves the character to delete_grace. Cancellable within
// s.cfg.DeleteGrace. Worker `character_delete_expire` (Faz 01 sprint 6+) is
// the one that actually removes the row when the deadline passes.
func (s *Service) RequestDelete(ctx context.Context, characterID, accountID uuid.UUID) (*Character, error) {
	c, err := s.GetOwned(ctx, characterID, accountID)
	if err != nil {
		return nil, err
	}
	if c.State != StateActive {
		return nil, ErrNotPlayable
	}
	deadline := s.cfg.Now().UTC().Add(s.cfg.DeleteGrace)
	if err := s.repo.SetDeleteGrace(ctx, c.ID, deadline); err != nil {
		return nil, err
	}
	c.State = StateDeleteGrace
	c.DeleteGraceAt = &deadline
	return c, nil
}

// CancelDelete restores a delete_grace character back to active.
func (s *Service) CancelDelete(ctx context.Context, characterID, accountID uuid.UUID) (*Character, error) {
	c, err := s.GetOwned(ctx, characterID, accountID)
	if err != nil {
		return nil, err
	}
	if c.State != StateDeleteGrace {
		return nil, ErrNotInDeleteGrace
	}
	if err := s.repo.CancelDeleteGrace(ctx, c.ID); err != nil {
		return nil, err
	}
	c.State = StateActive
	c.DeleteGraceAt = nil
	return c, nil
}

// --- Rename -----------------------------------------------------------------

func (s *Service) Rename(ctx context.Context, characterID, accountID uuid.UUID, newName string) (*Character, error) {
	c, err := s.GetOwned(ctx, characterID, accountID)
	if err != nil {
		return nil, err
	}
	if !c.IsPlayable() {
		return nil, ErrNotPlayable
	}
	name, normalized, err := validateName(newName)
	if err != nil {
		return nil, err
	}
	if normalized == c.NameNormalized {
		// No-op rename — return current record unchanged.
		return c, nil
	}
	if err := s.repo.Rename(ctx, c.ID, name, normalized); err != nil {
		return nil, err
	}
	c.Name = name
	c.NameNormalized = normalized
	return c, nil
}

// --- XP grant --------------------------------------------------------------

// XPGrantResult reports what the grant did: the character AFTER the update
// plus whether one or more level thresholds were crossed. Used by hunt/dungeon
// reveal handlers to surface "you leveled up" to the client.
type XPGrantResult struct {
	Character *Character
	LeveledUp bool
}

// MarshalJSON keeps the reveal payload to what clients need — never the full
// character row (account id, timestamps, internal state).
func (r XPGrantResult) MarshalJSON() ([]byte, error) {
	out := struct {
		Level     int16 `json:"level"`
		XP        int64 `json:"xp"`
		LeveledUp bool  `json:"leveled_up"`
	}{LeveledUp: r.LeveledUp}
	if r.Character != nil {
		out.Level, out.XP = r.Character.Level, r.Character.XP
	}
	return json.Marshal(out)
}

// GrantXP awards experience without an ownership check — reveal handlers
// upstream have already validated that the character belongs to the caller
// (activity contract ownership, hunt contract ownership, etc). Refuses
// negative deltas via ErrInvalidXP.
func (s *Service) GrantXP(ctx context.Context, characterID uuid.UUID, delta int64) (*XPGrantResult, error) {
	if delta < 0 {
		return nil, ErrInvalidXP
	}
	c, leveled, err := s.repo.AddXP(ctx, characterID, delta)
	if err != nil {
		return nil, err
	}
	if leveled {
		s.log.InfoContext(ctx, "character leveled up",
			slog.String("character_id", characterID.String()),
			slog.Int("new_level", int(c.Level)),
		)
	}
	return &XPGrantResult{Character: c, LeveledUp: leveled}, nil
}

// --- Location (Sprint 16) --------------------------------------------------

// GetLocation returns the character's current region code, or empty when the
// column is NULL (fresh character — world.Service treats as "at starting").
// No ownership check; internal callers already validated the character id.
func (s *Service) GetLocation(ctx context.Context, characterID uuid.UUID) (string, error) {
	return s.repo.GetLocation(ctx, characterID)
}

// SetLocation writes the character's current region code. Called by
// world.Service after travel completion — the caller has already ownership-
// validated the character id.
func (s *Service) SetLocation(ctx context.Context, characterID uuid.UUID, regionCode string) error {
	return s.repo.SetLocation(ctx, characterID, regionCode)
}

// --- Character-scope authorization helper -----------------------------------

// Owns returns nil if the account owns the character (regardless of state);
// used by the CharacterContext middleware.
func (s *Service) Owns(ctx context.Context, characterID, accountID uuid.UUID) error {
	c, err := s.repo.Get(ctx, characterID)
	if err != nil {
		return err
	}
	if c.AccountID != accountID {
		return ErrForbidden
	}
	return nil
}
