package dungeon

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/progression"
)

// --- Ports ------------------------------------------------------------------

// DungeonLookup exposes the dungeon template catalog. CatalogService satisfies.
type DungeonLookup interface {
	Catalog() *Catalog
}

// MobsLookup provides combat mob templates. combat.Service satisfies.
type MobsLookup interface {
	Catalog() *combat.Catalog
}

// GearDropPct is the chance (percent) of a gear drop per wave kill (DD-229).
const GearDropPct = 4

// ItemGranter — item.Service.Grant satisfies exactly this shape.
type ItemGranter interface {
	Grant(
		ctx context.Context,
		characterID uuid.UUID,
		sourceRefID uuid.UUID,
		seed int64, kind activity.Kind,
		common, rare, legendary int,
		gear activity.GearDrop,
	) ([]activity.GrantedItem, error)
}

// XPGranter — character.Service.GrantXP satisfies.
type XPGranter interface {
	GrantXP(ctx context.Context, characterID uuid.UUID, delta int64) (*character.XPGrantResult, error)
}

// SheetFetcher — character.Service satisfies.
type SheetFetcher interface {
	Sheet(ctx context.Context, characterID, accountID uuid.UUID) (*character.Sheet, error)
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
	GetOwnedInternal(ctx context.Context, characterID uuid.UUID) (*character.Character, error)
}

// ServiceConfig — clock injection for tests.
type ServiceConfig struct {
	Now func() time.Time
}

func DefaultServiceConfig() ServiceConfig { return ServiceConfig{Now: time.Now} }

// Service is the dungeon business logic.
type Service struct {
	repo     Repository
	dungeons DungeonLookup
	mobs     MobsLookup
	chars    SheetFetcher
	xp       XPGranter
	items    ItemGranter
	txn      economy.TxnService
	cfg      ServiceConfig
	log      *slog.Logger
	rngNew   func(seed int64) combat.RNG
	progress progression.Recorder
	hp       HPWriter
	skills   SkillSource
}

// WithRecorder feeds reveal results into progression (Sprint 24).
func (s *Service) WithRecorder(r progression.Recorder) *Service { s.progress = r; return s }

// HPWriter persists post-fight HP (character.Service). Optional: nil keeps
// the pre-DD-218 behaviour of HP never leaving the contract.
type HPWriter interface {
	SetHP(ctx context.Context, characterID uuid.UUID, hp int32, at time.Time) error
}

// WithHPWriter makes reveals write the final HP back to the character.
func (s *Service) WithHPWriter(w HPWriter) *Service { s.hp = w; return s }

// SkillSource yields the class actives a sheet has unlocked (DD-222).
type SkillSource func(sheet *character.Sheet) []combat.Skill

// WithSkills freezes unlocked class actives into the start snapshot.
func (s *Service) WithSkills(src SkillSource) *Service { s.skills = src; return s }

func NewService(
	repo Repository,
	dungeons DungeonLookup,
	mobs MobsLookup,
	chars SheetFetcher,
	xp XPGranter,
	items ItemGranter,
	txn economy.TxnService,
	cfg ServiceConfig,
	logger *slog.Logger,
	rngNew func(seed int64) combat.RNG,
) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo: repo, dungeons: dungeons, mobs: mobs, chars: chars,
		xp: xp, items: items, txn: txn, cfg: cfg, log: logger, rngNew: rngNew,
	}
}

// --- Start ------------------------------------------------------------------

type StartRequest struct {
	CharacterID uuid.UUID
	AccountID   uuid.UUID
	DungeonCode string
}

type StartResult struct {
	Contract          *Contract `json:"contract"`
	PlannedEncounters int       `json:"planned_encounters"`
	DisplayNameKey    string    `json:"display_name_key"`
	WeekKey           string    `json:"week_key"`
}

func (s *Service) Start(ctx context.Context, req StartRequest) (*StartResult, error) {
	if req.CharacterID == uuid.Nil || req.AccountID == uuid.Nil {
		return nil, ErrForbidden
	}
	if req.DungeonCode == "" {
		return nil, ErrValidation
	}
	if err := s.chars.Owns(ctx, req.CharacterID, req.AccountID); err != nil {
		return nil, err
	}

	tmpl, ok := s.dungeons.Catalog().Get(req.DungeonCode)
	if !ok {
		return nil, ErrUnknownDungeon
	}

	// Level gate.
	c, err := s.chars.GetOwnedInternal(ctx, req.CharacterID)
	if err != nil {
		return nil, err
	}
	if c.Level < tmpl.MinLevel {
		return nil, ErrLevelRequirement
	}

	// Weekly lockout gate — refuse if character already cleared this week.
	now := s.cfg.Now().UTC()
	weekKey := ISOWeek(now)
	locked, err := s.repo.HasLockout(ctx, req.CharacterID, tmpl.Code, weekKey)
	if err != nil {
		return nil, err
	}
	if locked {
		return nil, ErrWeeklyLockout
	}

	// Sanity: every wave + boss must resolve in the mob catalog. Fail early
	// so a broken template does not eat a contract slot.
	mobCat := s.mobs.Catalog()
	for _, w := range tmpl.Waves {
		if _, ok := mobCat.Get(w.MobCode); !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownMob, w.MobCode)
		}
	}
	if _, ok := mobCat.Get(tmpl.BossMobCode); !ok {
		return nil, fmt.Errorf("%w: boss %s", ErrUnknownMob, tmpl.BossMobCode)
	}

	sheet, err := s.chars.Sheet(ctx, req.CharacterID, req.AccountID)
	if err != nil {
		return nil, err
	}
	if sheet.Resources.HPCurrent < 1 {
		return nil, ErrTooWounded
	}
	snap := combat.FromSheet(sheet)
	// Short rests between packs: a third of the out-of-combat regen over the
	// run time per encounter (DD-229).
	snap.RegenPerEncounter = int32(float64(tmpl.DurationSecs/int32(tmpl.TotalEncounters())/6) *
		character.HPRegenPer2s(sheet.Primary.Vitality))
	if s.skills != nil {
		snap.Skills = s.skills(sheet)
	}
	snapJSON, err := json.Marshal(snap)
	if err != nil {
		return nil, fmt.Errorf("dungeon: marshal snapshot: %w", err)
	}

	contract := &Contract{
		ID:             uuid.New(),
		CharacterID:    req.CharacterID,
		DungeonCode:    tmpl.Code,
		DurationSecs:   tmpl.DurationSecs,
		StartedAt:      now,
		ExpiresAt:      now.Add(time.Duration(tmpl.DurationSecs) * time.Second),
		Seed:           seedFromIDs(uuid.New(), req.CharacterID),
		State:          StateActive,
		PlayerSnapshot: snapJSON,
		CreatedAt:      now,
		UpdatedAt:      now,
		Version:        1,
	}
	if err := s.repo.Create(ctx, contract); err != nil {
		return nil, err
	}
	return &StartResult{
		Contract:          contract,
		PlannedEncounters: tmpl.TotalEncounters(),
		DisplayNameKey:    tmpl.DisplayNameKey,
		WeekKey:           weekKey,
	}, nil
}

// --- Reveal -----------------------------------------------------------------

type RevealResult struct {
	Run         *RunResult               `json:"run"`
	Elapsed     int32                    `json:"elapsed_seconds"`
	LedgerTxnID uuid.UUID                `json:"ledger_txn_id"`
	Items       []activity.GrantedItem   `json:"items,omitempty"`
	XPGain      *character.XPGrantResult `json:"xp,omitempty"`
	Replayed    bool                     `json:"replayed"`
	WeekKey     string                   `json:"week_key"`
	Lockout     bool                     `json:"lockout_written"`
}

// Reveal resolves the run when duration elapsed. Idempotent — repeated calls
// return the cached response. Lockout row written iff Cleared.
func (s *Service) Reveal(ctx context.Context, contractID, accountID uuid.UUID) (*RevealResult, error) {
	c, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if err := s.chars.Owns(ctx, c.CharacterID, accountID); err != nil {
		return nil, err
	}

	if c.State == StateRevealed {
		var rr RunResult
		_ = json.Unmarshal(c.RevealResult, &rr)
		txn := uuid.Nil
		if c.RevealTxnID != nil {
			txn = *c.RevealTxnID
		}
		return &RevealResult{
			Run: &rr, Elapsed: c.DurationSecs, LedgerTxnID: txn,
			Replayed: true, WeekKey: ISOWeek(c.StartedAt), Lockout: rr.Cleared,
		}, nil
	}
	if c.State != StateActive {
		return nil, ErrNotActive
	}

	now := s.cfg.Now().UTC()
	if !c.IsReady(now) {
		return nil, ErrNotReady
	}

	tmpl, ok := s.dungeons.Catalog().Get(c.DungeonCode)
	if !ok {
		return nil, ErrUnknownDungeon
	}
	var snapshot combat.Combatant
	if err := json.Unmarshal(c.PlayerSnapshot, &snapshot); err != nil {
		return nil, fmt.Errorf("dungeon: unmarshal snapshot: %w", err)
	}

	run := Roll(c.Seed, snapshot, tmpl, s.mobs.Catalog(), combat.DefaultMaxTurns, s.rngNew)

	// Item mint — same idempotency pattern as hunt (source_ref_id = contract.ID).
	// Gear (DD-229): a GearDropPct chance per wave kill, and a guaranteed,
	// better piece from the boss — the boss's tier, the character's class.
	var gear activity.GearDrop
	if run.Victories > 0 {
		sheet, err := s.chars.Sheet(ctx, c.CharacterID, accountID)
		if err != nil {
			return nil, err
		}
		boss, _ := s.mobs.Catalog().Get(tmpl.BossMobCode)
		gear = activity.GearDrop{Class: sheet.Class, Level: int(boss.Level), Rolls: run.Victories, Pct: GearDropPct}
		if run.BossKilled {
			gear.Rolls--
			gear.Guaranteed, gear.Boss = 1, true
		}
	}
	var granted []activity.GrantedItem
	if s.items != nil && (run.CommonDrops+run.RareDrops+run.LegendaryDrops+gear.Rolls+gear.Guaranteed) > 0 {
		granted, err = s.items.Grant(ctx,
			c.CharacterID, c.ID, c.Seed, activity.KindHunting,
			run.CommonDrops, run.RareDrops, run.LegendaryDrops, gear,
		)
		if err != nil {
			return nil, fmt.Errorf("dungeon reveal item grant: %w", err)
		}
	}

	// Ledger gold — idempotent by contract id.
	var txnID uuid.UUID
	if run.TotalGoldCopper > 0 {
		res, err := s.txn.Write(ctx, economy.Intent{
			Kind:           "dungeon_reveal",
			RefKind:        "dungeon_contract",
			RefID:          c.ID,
			IdempotencyKey: "dungeon-reveal-" + c.ID.String(),
			Entries: []economy.Entry{
				{Currency: "gold", AccountKind: economy.AccountSystem,
					AccountRef: economy.SystemRewardSource, Delta: -run.TotalGoldCopper},
				{Currency: "gold", AccountKind: economy.AccountCharacter,
					AccountRef: c.CharacterID, Delta: run.TotalGoldCopper},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("dungeon reveal ledger: %w", err)
		}
		txnID = res.TxnID
	}

	resultJSON, err := json.Marshal(run)
	if err != nil {
		return nil, err
	}
	weekKey := ISOWeek(c.StartedAt)
	lockoutWritten := run.Cleared

	// HP write-back (deterministic, retry-safe; see hunt.Reveal). The run is
	// resolved as of its end.
	if s.hp != nil {
		if err := s.hp.SetHP(ctx, c.CharacterID, run.PlayerHPFinal, c.ExpiresAt); err != nil {
			return nil, fmt.Errorf("dungeon reveal hp: %w", err)
		}
	}

	// MarkRevealed writes contract row + (conditionally) lockout row in one tx.
	err = s.repo.MarkRevealed(ctx, c.ID, resultJSON, txnID, now, lockoutWritten, c.CharacterID, c.DungeonCode, weekKey)
	if err != nil {
		if err == ErrWeeklyLockout {
			// Race: another reveal (via retry or parallel client) already
			// wrote the lockout. Refetch and return cached.
			refreshed, gerr := s.repo.GetByID(ctx, c.ID)
			if gerr == nil && refreshed.State == StateRevealed {
				var cached RunResult
				_ = json.Unmarshal(refreshed.RevealResult, &cached)
				t := uuid.Nil
				if refreshed.RevealTxnID != nil {
					t = *refreshed.RevealTxnID
				}
				return &RevealResult{
					Run: &cached, Elapsed: c.DurationSecs, LedgerTxnID: t,
					Replayed: true, WeekKey: weekKey, Lockout: cached.Cleared,
				}, nil
			}
			return nil, err
		}
		if err == ErrNotActive {
			refreshed, gerr := s.repo.GetByID(ctx, c.ID)
			if gerr == nil && refreshed.State == StateRevealed {
				var cached RunResult
				_ = json.Unmarshal(refreshed.RevealResult, &cached)
				t := uuid.Nil
				if refreshed.RevealTxnID != nil {
					t = *refreshed.RevealTxnID
				}
				return &RevealResult{
					Run: &cached, Elapsed: c.DurationSecs, LedgerTxnID: t,
					Replayed: true, WeekKey: weekKey, Lockout: cached.Cleared,
				}, nil
			}
		}
		return nil, err
	}

	var xpResult *character.XPGrantResult
	if s.xp != nil && run.TotalXP > 0 {
		xpResult, err = s.xp.GrantXP(ctx, c.CharacterID, run.TotalXP)
		if err != nil {
			s.log.WarnContext(ctx, "dungeon: xp grant failed (contract already revealed, no retry)",
				slog.String("run", c.ID.String()), slog.String("err", err.Error()))
		}
	}

	events := []progression.Event{{Kind: progression.EventMobKilled, Amount: int64(run.Victories)}}
	if run.BossKilled {
		events[0].Amount-- // the boss is counted separately
		events = append(events, progression.Event{Kind: progression.EventBossKilled, Amount: 1})
	}
	if run.Cleared {
		events = append(events,
			progression.Event{Kind: progression.EventDungeonCleared, Amount: 1},
			progression.Event{Kind: progression.EventDungeonCleared + ":" + c.DungeonCode, Amount: 1})
	}
	progression.RecordIf(ctx, s.progress, c.CharacterID, events...)

	return &RevealResult{
		Run: run, Elapsed: c.DurationSecs, LedgerTxnID: txnID,
		Items: granted, XPGain: xpResult, Replayed: false,
		WeekKey: weekKey, Lockout: lockoutWritten,
	}, nil
}

// --- Stop -------------------------------------------------------------------

func (s *Service) Stop(ctx context.Context, contractID, accountID uuid.UUID) (*Contract, error) {
	c, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if err := s.chars.Owns(ctx, c.CharacterID, accountID); err != nil {
		return nil, err
	}
	if c.State == StateCancelled {
		return c, nil
	}
	if c.State != StateActive {
		return nil, ErrNotActive
	}
	now := s.cfg.Now().UTC()
	if err := s.repo.MarkCancelled(ctx, c.ID, now); err != nil {
		return nil, err
	}
	c.State = StateCancelled
	c.CancelledAt = &now
	return c, nil
}

// --- Active / History / Catalog listing ------------------------------------

func (s *Service) Active(ctx context.Context, characterID, accountID uuid.UUID) (*Contract, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	return s.repo.Active(ctx, characterID)
}

func (s *Service) History(ctx context.Context, characterID, accountID uuid.UUID, limit int) ([]Contract, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	return s.repo.History(ctx, characterID, limit)
}

// ListCatalog returns the loaded dungeon templates for the client picker.
func (s *Service) ListCatalog() []Template { return s.dungeons.Catalog().All() }

// --- Helpers ----------------------------------------------------------------

func seedFromIDs(runID, characterID uuid.UUID) int64 {
	var acc uint64 = 1469598103934665603
	for _, b := range runID[:] {
		acc ^= uint64(b)
		acc *= 1099511628211
	}
	for _, b := range characterID[:] {
		acc ^= uint64(b)
		acc *= 1099511628211
	}
	return int64(acc)
}
