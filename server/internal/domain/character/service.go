package character

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// ServiceConfig tunes the character service.
type ServiceConfig struct {
	SlotCap int // DD-OPEN-8, v1 default 2
	Now     func() time.Time
}

// DefaultServiceConfig returns Faz 00 baseline.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		SlotCap: 5, // DD-251
		Now:     time.Now,
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
	repo     Repository
	sheeter  SheetComputer
	verified func(ctx context.Context, accountID uuid.UUID) (bool, error) // nil = no e-mail gate
	cfg      ServiceConfig
	log      *slog.Logger
}

func NewService(repo Repository, cfg ServiceConfig, logger *slog.Logger) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.SlotCap <= 0 {
		cfg.SlotCap = 5
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

// WithEmailVerification makes Create refuse accounts whose e-mail is not
// verified yet (docs/auth/register.md §7, DD-219).
func (s *Service) WithEmailVerification(check func(ctx context.Context, accountID uuid.UUID) (bool, error)) *Service {
	s.verified = check
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

	if s.verified != nil {
		ok, err := s.verified(ctx, req.AccountID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrEmailNotVerified
		}
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

// --- Delete -----------------------------------------------------------------

// Delete removes the character at once, with its items (DD-283). confirmName
// must be the character's name (case-insensitive) so a stray call cannot
// delete the wrong one. Frozen characters stay for ops.
func (s *Service) Delete(ctx context.Context, characterID, accountID uuid.UUID, confirmName string) error {
	c, err := s.GetOwned(ctx, characterID, accountID)
	if err != nil {
		return err
	}
	if c.State != StateActive && c.State != StateDeleteGrace {
		return ErrNotPlayable
	}
	if _, normalized, err := validateName(confirmName); err != nil || normalized != c.NameNormalized {
		return ErrNameMismatch
	}
	return s.repo.Purge(ctx, c.ID)
}

// CancelDelete restores a delete_grace character (one left pending by an
// account deletion or from before direct deletion) back to active.
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
// Level is the character's level (crafting gates on it, DD-316).
func (s *Service) Level(ctx context.Context, characterID uuid.UUID) (int, error) {
	c, err := s.repo.Get(ctx, characterID)
	if err != nil {
		return 0, err
	}
	return int(c.Level), nil
}

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
