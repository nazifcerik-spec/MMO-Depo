package hunt

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

// --- Ports (kept narrow — hunt does not "own" the underlying services) ------

// MobLookup exposes the combat mob catalog. combat.Service satisfies.
type MobLookup interface {
	Catalog() *combat.Catalog
}

// GearDropPct is the chance (percent) of a gear drop per victory (DD-229).
const GearDropPct = 4

// ItemGranter is the drop-side capability. item.Service.Grant already
// satisfies exactly this signature (used by activity too).
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

// XPGranter awards experience — character.Service satisfies via GrantXP.
type XPGranter interface {
	GrantXP(ctx context.Context, characterID uuid.UUID, delta int64) (*character.XPGrantResult, error)
}

// SheetFetcher builds a live sheet (with equipped bonuses). character.Service
// satisfies via Sheet.
type SheetFetcher interface {
	Sheet(ctx context.Context, characterID, accountID uuid.UUID) (*character.Sheet, error)
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

// ServiceConfig — clock injection for tests.
type ServiceConfig struct {
	Now func() time.Time
}

func DefaultServiceConfig() ServiceConfig { return ServiceConfig{Now: time.Now} }

// Service is the hunt business logic.
type Service struct {
	repo     Repository
	mobs     MobLookup
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

// NewService wires the service. rngNew is the RNG factory (activity.NewRNG in
// prod; a stub in tests).
func NewService(
	repo Repository,
	mobs MobLookup,
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
		repo: repo, mobs: mobs, chars: chars, xp: xp, items: items, txn: txn,
		cfg: cfg, log: logger, rngNew: rngNew,
	}
}

// --- Start ------------------------------------------------------------------

// StartRequest is the input for POST /v1/character/{id}/hunt.
type StartRequest struct {
	CharacterID     uuid.UUID
	AccountID       uuid.UUID
	MobCode         string
	DurationSeconds int32
	EncounterSecs   int32
}

// StartResult wraps the created contract + a full-duration peek so the client
// can preview expected encounter count.
type StartResult struct {
	Contract          *Contract `json:"contract"`
	PlannedEncounters int       `json:"planned_encounters_full"`
	MobDisplayKey     string    `json:"mob_display_name_key"`
}

func (s *Service) Start(ctx context.Context, req StartRequest) (*StartResult, error) {
	if req.CharacterID == uuid.Nil || req.AccountID == uuid.Nil {
		return nil, ErrForbidden
	}
	if req.MobCode == "" {
		return nil, ErrValidation
	}
	if req.DurationSeconds <= 0 || req.DurationSeconds > MaxHuntSeconds {
		return nil, ErrValidation
	}
	if req.EncounterSecs <= 0 {
		req.EncounterSecs = DefaultEncounterSeconds
	}
	if req.EncounterSecs < MinEncounterSeconds || req.DurationSeconds < req.EncounterSecs {
		return nil, ErrValidation // a hunt shorter than one encounter can't yield anything
	}

	// Ownership check via SheetFetcher.Owns — leak-safe.
	if err := s.chars.Owns(ctx, req.CharacterID, req.AccountID); err != nil {
		return nil, err
	}

	mob, ok := s.mobs.Catalog().Get(req.MobCode)
	if !ok {
		return nil, ErrUnknownMob
	}

	// Freeze the player's Combatant snapshot at Start — mid-hunt gear swaps
	// do not warp the log (matches activity's Snapshot pattern).
	sheet, err := s.chars.Sheet(ctx, req.CharacterID, req.AccountID)
	if err != nil {
		return nil, err
	}
	if sheet.Resources.HPCurrent < 1 {
		return nil, ErrTooWounded
	}
	snapshot := combat.FromSheet(sheet)
	// Resting between encounters: the out-of-combat regen (DD-218) over the
	// encounter interval (DD-229).
	snapshot.RegenPerEncounter = int32(float64(req.EncounterSecs/2) * character.HPRegenPer2s(sheet.Primary.Vitality))
	if s.skills != nil {
		snapshot.Skills = s.skills(sheet)
	}
	snapJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("hunt: marshal snapshot: %w", err)
	}

	now := s.cfg.Now().UTC()
	c := &Contract{
		ID:             uuid.New(),
		CharacterID:    req.CharacterID,
		MobCode:        req.MobCode,
		EncounterSecs:  req.EncounterSecs,
		DurationSecs:   req.DurationSeconds,
		StartedAt:      now,
		ExpiresAt:      now.Add(time.Duration(req.DurationSeconds) * time.Second),
		Seed:           seedFromIDs(uuid.New(), req.CharacterID),
		State:          StateActive,
		PlayerSnapshot: snapJSON,
		CreatedAt:      now,
		UpdatedAt:      now,
		Version:        1,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &StartResult{
		Contract:          c,
		PlannedEncounters: c.EncountersPlanned(c.DurationSecs),
		MobDisplayKey:     mob.DisplayNameKey,
	}, nil
}

// --- Reveal -----------------------------------------------------------------

// RevealResult is the reveal HTTP response.
type RevealResult struct {
	Hunt        *HuntResult              `json:"hunt"`
	Elapsed     int32                    `json:"elapsed_seconds"`
	LedgerTxnID uuid.UUID                `json:"ledger_txn_id"`
	Items       []activity.GrantedItem   `json:"items,omitempty"`
	XPGain      *character.XPGrantResult `json:"xp,omitempty"`
	Replayed    bool                     `json:"replayed"`
}

// Reveal is idempotent — repeated calls return the cached response. Rewards
// are granted in-order (item mint → ledger → xp) with the same "ledger
// succeeded but mark-revealed race lost" recovery pattern used in activity.
func (s *Service) Reveal(ctx context.Context, contractID, accountID uuid.UUID) (*RevealResult, error) {
	c, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if err := s.chars.Owns(ctx, c.CharacterID, accountID); err != nil {
		return nil, err
	}

	if c.State == StateRevealed {
		var hr HuntResult
		_ = json.Unmarshal(c.RevealResult, &hr)
		txn := uuid.Nil
		if c.RevealTxnID != nil {
			txn = *c.RevealTxnID
		}
		return &RevealResult{
			Hunt: &hr, Elapsed: c.DurationSecs, LedgerTxnID: txn, Replayed: true,
		}, nil
	}
	if c.State != StateActive {
		return nil, ErrNotActive
	}

	now := s.cfg.Now().UTC()
	elapsed := c.ElapsedSeconds(now)
	planned := c.EncountersPlanned(elapsed)

	var snapshot combat.Combatant
	if err := json.Unmarshal(c.PlayerSnapshot, &snapshot); err != nil {
		return nil, fmt.Errorf("hunt: unmarshal snapshot: %w", err)
	}
	mob, ok := s.mobs.Catalog().Get(c.MobCode)
	if !ok {
		return nil, ErrUnknownMob
	}

	hr := Roll(c.Seed, snapshot, mob, planned, combat.DefaultMaxTurns, s.rngNew)

	// --- Item mint (idempotent by source_ref_id = contract.ID) --------------
	// Each victory has a GearDropPct chance of gear from the mob's tier,
	// usable by the character's class (DD-229).
	var gear activity.GearDrop
	if hr.Victories > 0 {
		sheet, err := s.chars.Sheet(ctx, c.CharacterID, accountID)
		if err != nil {
			return nil, err
		}
		gear = activity.GearDrop{Class: sheet.Class, Level: int(mob.Level), Rolls: hr.Victories, Pct: GearDropPct}
	}
	var granted []activity.GrantedItem
	if s.items != nil && (hr.CommonDrops+hr.RareDrops+hr.LegendaryDrops+gear.Rolls) > 0 {
		granted, err = s.items.Grant(ctx,
			c.CharacterID, c.ID, c.Seed, activity.KindHunting,
			hr.CommonDrops, hr.RareDrops, hr.LegendaryDrops, gear,
		)
		if err != nil {
			return nil, fmt.Errorf("hunt reveal item grant: %w", err)
		}
	}

	// --- Ledger gold credit (idempotent by key = "hunt-reveal-<contract>") ---
	var txnID uuid.UUID
	if hr.TotalGoldCopper > 0 {
		res, err := s.txn.Write(ctx, economy.Intent{
			Kind:           "hunt_reveal",
			RefKind:        "hunt_contract",
			RefID:          c.ID,
			IdempotencyKey: "hunt-reveal-" + c.ID.String(),
			Entries: []economy.Entry{
				{Currency: "gold", AccountKind: economy.AccountSystem,
					AccountRef: economy.SystemRewardSource, Delta: -hr.TotalGoldCopper},
				{Currency: "gold", AccountKind: economy.AccountCharacter,
					AccountRef: c.CharacterID, Delta: hr.TotalGoldCopper},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("hunt reveal ledger: %w", err)
		}
		txnID = res.TxnID
	}

	// --- XP grant (also idempotent-safe: AddXP is authoritative per row +
	// happens AFTER MarkRevealed so a retry sees state=revealed and skips
	// double-grant via the replay path above). ------------------------------
	resultJSON, err := json.Marshal(hr)
	if err != nil {
		return nil, err
	}
	// HP write-back before MarkRevealed: the value and timestamp are
	// deterministic (same roll on retry), so repeating it is harmless, and a
	// crash in between never leaves a revealed hunt with stale HP.
	if s.hp != nil {
		fightEnd := c.StartedAt.Add(time.Duration(hr.Completed) * time.Duration(c.EncounterSecs) * time.Second)
		if fightEnd.After(now) {
			fightEnd = now
		}
		if err := s.hp.SetHP(ctx, c.CharacterID, hr.PlayerHPFinal, fightEnd); err != nil {
			return nil, fmt.Errorf("hunt reveal hp: %w", err)
		}
	}
	if err := s.repo.MarkRevealed(ctx, c.ID, resultJSON, txnID, now); err != nil {
		// Race: someone else revealed. Refetch + return cached.
		refreshed, gerr := s.repo.GetByID(ctx, c.ID)
		if gerr == nil && refreshed.State == StateRevealed {
			var cached HuntResult
			_ = json.Unmarshal(refreshed.RevealResult, &cached)
			t := uuid.Nil
			if refreshed.RevealTxnID != nil {
				t = *refreshed.RevealTxnID
			}
			return &RevealResult{Hunt: &cached, Elapsed: elapsed, LedgerTxnID: t, Replayed: true}, nil
		}
		return nil, err
	}

	var xpResult *character.XPGrantResult
	if s.xp != nil && hr.TotalXP > 0 {
		xpResult, err = s.xp.GrantXP(ctx, c.CharacterID, hr.TotalXP)
		if err != nil {
			s.log.WarnContext(ctx, "hunt: xp grant failed (contract already revealed, no retry)",
				slog.String("hunt", c.ID.String()), slog.String("err", err.Error()))
		}
	}

	progression.RecordIf(ctx, s.progress, c.CharacterID,
		progression.Event{Kind: progression.EventMobKilled, Amount: int64(hr.Victories)},
		progression.Event{Kind: progression.EventMobKilled + ":" + c.MobCode, Amount: int64(hr.Victories)})

	return &RevealResult{
		Hunt: hr, Elapsed: elapsed, LedgerTxnID: txnID,
		Items: granted, XPGain: xpResult, Replayed: false,
	}, nil
}

// --- Stop -------------------------------------------------------------------

// Stop cancels an active hunt with no reward. Idempotent for already-cancelled;
// returns ErrNotActive for revealed contracts.
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

// --- Active / History -------------------------------------------------------

func (s *Service) Active(ctx context.Context, characterID, accountID uuid.UUID) (*Contract, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	c, err := s.repo.Active(ctx, characterID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) History(ctx context.Context, characterID, accountID uuid.UUID, limit int) ([]Contract, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	return s.repo.History(ctx, characterID, limit)
}

// --- Helpers ----------------------------------------------------------------

// seedFromIDs derives a deterministic 64-bit seed from a fresh hunt id + the
// character id. Same shape as activity.SeedFromContract — namespaces the
// combat RNG from the ledger idempotency-key so tampering with one does not
// reveal the other's stream.
func seedFromIDs(huntID, characterID uuid.UUID) int64 {
	h := huntID
	c := characterID
	// FNV-1a-ish fold — cheap, no allocations, no crypto pretense.
	var acc uint64 = 1469598103934665603
	for _, b := range h[:] {
		acc ^= uint64(b)
		acc *= 1099511628211
	}
	for _, b := range c[:] {
		acc ^= uint64(b)
		acc *= 1099511628211
	}
	return int64(acc)
}
