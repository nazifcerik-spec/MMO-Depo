package hunt

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/consumable"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/partyxp"
	"github.com/otter/server/internal/progression"
	"math"
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
	Grant(ctx context.Context, characterID, sourceRefID uuid.UUID, seed int64, d activity.Drops) ([]activity.GrantedItem, error)
}

// XPGranter awards experience — character.Service satisfies via GrantXP.
type XPGranter interface {
	GrantXP(ctx context.Context, characterID uuid.UUID, delta int64, src character.XPSource) (*character.XPGrantResult, error)
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
	repo      Repository
	mobs      MobLookup
	chars     SheetFetcher
	xp        XPGranter
	items     ItemGranter
	txn       economy.TxnService
	cfg       ServiceConfig
	log       *slog.Logger
	rngNew    func(seed int64) combat.RNG
	progress  progression.Recorder
	notifier  notification.Notifier // reward history (DD-242); nil = none
	hp        HPWriter
	skills    SkillSource
	bag       BagRoom         // nil = no room check (DD-246)
	regions   RegionGate      // nil = every mob in reach (DD-273)
	templates Templates       // supply effects (DD-288); nil = no supplies
	food      FoodBuffs       // a dish eaten before the hunt (DD-288); nil = in the sheet
	party     partyxp.Booster // party XP bonus (DD-305); nil = none
}

// WithPartyBonus adds the party bonus for members hunting alongside (DD-305).
func (s *Service) WithPartyBonus(p partyxp.Booster) *Service { s.party = p; return s }

// Templates resolves item templates for supplies (item.Service).
type Templates interface {
	Catalog() *item.Catalog
}

// WithSupplies lets a hunt take potions and dishes from the bag (DD-288).
func (s *Service) WithSupplies(t Templates) *Service { s.templates = t; return s }

// FoodBuffs reports a character's running food buff (consumable.Service).
type FoodBuffs interface {
	Active(ctx context.Context, characterID uuid.UUID) (*consumable.Buff, error)
}

// WithFood makes the start snapshot leave a running food buff out of the
// sheet and carry it with its remaining time instead (DD-288).
func (s *Service) WithFood(f FoodBuffs) *Service { s.food = f; return s }

// WithRecorder feeds reveal results into progression (Sprint 24).
func (s *Service) WithRecorder(r progression.Recorder) *Service { s.progress = r; return s }

// WithNotifier records each collection in the reward history (DD-242).
func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notifier = n; return s }

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
// WithBagRoom refuses to start a run while the bag is full (DD-246).
func (s *Service) WithBagRoom(b BagRoom) *Service { s.bag = b; return s }

// BagRoom tells whether the bag can take more items (DD-246);
// inventory.Service satisfies. An empty templateCode asks for an empty slot.
type BagRoom interface {
	HasRoom(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32) (bool, error)
}

func (s *Service) WithSkills(src SkillSource) *Service { s.skills = src; return s }

// RegionGate lists the regions a character can hunt in from where it stands
// (DD-273); world.Service implements it.
type RegionGate interface {
	HuntRegions(ctx context.Context, characterID uuid.UUID) ([]string, error)
}

// WithRegionGate makes hunting regional: a mob is hunted only from its own
// region or one linked to it. Without a gate every mob is in reach.
func (s *Service) WithRegionGate(g RegionGate) *Service { s.regions = g; return s }

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
	Supplies        []SupplyRequest // potions and dishes taken from the bag (DD-288)
}

// SupplyRequest asks for qty units of a consumable from the bag.
type SupplyRequest struct {
	Code string `json:"template_code"`
	Qty  int32  `json:"quantity"`
}

// SupplyFromTemplate is what a consumable does on a hunt: a healing potion
// (heal_hp / heal_pct), a dish (healing with a buff) or a mana potion
// (mana_pct). ok is false for anything else.
func SupplyFromTemplate(t item.Template, qty int32) (Supply, bool) {
	ue := t.UseEffect
	if ue == nil || ue.Amount <= 0 {
		return Supply{}, false
	}
	s := Supply{Code: t.Code, Qty: qty}
	switch ue.Kind {
	case "mana_pct":
		s.Kind, s.MPPct = SupplyMP, int32(ue.Amount)
	case "heal_pct", "heal_hp":
		s.Kind = SupplyHP
		if ue.Kind == "heal_pct" {
			s.HealPct = int32(ue.Amount)
		} else {
			s.HealFlat = int32(ue.Amount)
		}
		if b := ue.Buff; b != nil {
			s.Kind = SupplyFood
			s.BuffAtk, s.BuffDef, s.BuffSecs = int32(b.AttackPct), int32(b.DefensePct), int32(b.Minutes)*60
		}
	default:
		return Supply{}, false
	}
	return s, true
}

// supplies resolves the requested supplies for a character of level.
func (s *Service) supplies(reqs []SupplyRequest, level int) ([]Supply, error) {
	if len(reqs) == 0 {
		return nil, nil
	}
	if s.templates == nil || len(reqs) > MaxSupplyKinds {
		return nil, ErrBadSupply
	}
	out := make([]Supply, 0, len(reqs))
	seen := map[string]bool{}
	for _, r := range reqs {
		if r.Qty < 1 || r.Qty > 999 || seen[r.Code] {
			return nil, ErrBadSupply
		}
		seen[r.Code] = true
		t, err := s.templates.Catalog().Get(r.Code)
		if err != nil || int(t.MinLevel) > level {
			return nil, ErrBadSupply
		}
		sup, ok := SupplyFromTemplate(t, r.Qty)
		if !ok {
			return nil, ErrBadSupply
		}
		out = append(out, sup)
	}
	return out, nil
}

// StartResult wraps the created contract.
type StartResult struct {
	Contract      *Contract `json:"contract"`
	MobDisplayKey string    `json:"mob_display_name_key"`
}

func (s *Service) Start(ctx context.Context, req StartRequest) (*StartResult, error) {
	if req.CharacterID == uuid.Nil || req.AccountID == uuid.Nil {
		return nil, ErrForbidden
	}
	if req.MobCode == "" {
		return nil, ErrValidation
	}
	if req.DurationSeconds < MinHuntSeconds || req.DurationSeconds > MaxHuntSeconds {
		return nil, ErrValidation
	}

	// Ownership check via SheetFetcher.Owns — leak-safe.
	if err := s.chars.Owns(ctx, req.CharacterID, req.AccountID); err != nil {
		return nil, err
	}

	mob, ok := s.mobs.Catalog().Get(req.MobCode)
	if !ok || mob.RaidBoss { // raid bosses are fought by a party (DD-261)
		return nil, ErrUnknownMob
	}
	if s.regions != nil && mob.RegionCode != "" {
		reach, err := s.regions.HuntRegions(ctx, req.CharacterID)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(reach, mob.RegionCode) {
			return nil, ErrWrongRegion
		}
	}

	// Freeze the player's Combatant snapshot at Start — mid-hunt gear swaps
	// do not warp the log (matches activity's Snapshot pattern). A running
	// food buff stays out of it and runs on the hunt's own clock (DD-288).
	var load Loadout
	sheetCtx := ctx
	if s.food != nil {
		b, err := s.food.Active(ctx, req.CharacterID)
		if err != nil {
			return nil, err
		}
		if b != nil {
			if left := int32(b.ExpiresAt.Sub(s.cfg.Now()).Seconds()); left > 0 {
				load.Buff = &Buff{AttackPct: int32(b.AttackPct), DefensePct: int32(b.DefensePct), Secs: left}
			}
		}
		sheetCtx = consumable.WithoutBuff(ctx)
	}
	sheet, err := s.chars.Sheet(sheetCtx, req.CharacterID, req.AccountID)
	if err != nil {
		return nil, err
	}
	if load.Supplies, err = s.supplies(req.Supplies, sheet.Level); err != nil {
		return nil, err
	}
	if sheet.Resources.HPCurrent < 1 {
		return nil, ErrTooWounded
	}
	if s.bag != nil {
		ok, err := s.bag.HasRoom(ctx, req.CharacterID, "", 1)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrBagFull
		}
	}
	snapshot := combat.FromSheet(sheet)
	// Fights run back to back (DD-281): the out-of-combat regen (DD-218)
	// keeps ticking while the hunter fights.
	snapshot.RegenPerSec = character.HPRegenPer2s(sheet.Primary.Vitality) / 2
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
		EncounterSecs:  LegacyEncounterSeconds,
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
	if len(load.Supplies) > 0 || load.Buff != nil {
		c.Loadout = &load
	}
	if err := s.repo.Create(ctx, c); err != nil { // holds the supplies in the same transaction
		return nil, err
	}
	return &StartResult{Contract: c, MobDisplayKey: mob.DisplayNameKey}, nil
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
	// Party is the party bonus on the XP (DD-305); nil = none.
	Party *partyxp.Outcome `json:"party_share,omitempty"`
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

	var snapshot combat.Combatant
	if err := json.Unmarshal(c.PlayerSnapshot, &snapshot); err != nil {
		return nil, fmt.Errorf("hunt: unmarshal snapshot: %w", err)
	}
	mob, ok := s.mobs.Catalog().Get(c.MobCode)
	if !ok {
		return nil, ErrUnknownMob
	}

	// Each kill pays the mob's own XP, less for a level gap either way
	// (DD-281; the snapshot's level, so a level-up mid-hunt does not change
	// the reward). Fights run back to back for the time elapsed.
	mob.XPReward = HuntXP(mob.XPReward, int(snapshot.Level), int(mob.Level))
	var load Loadout
	if c.Loadout != nil {
		load = *c.Loadout
	}
	hr := RollWith(c.Seed, snapshot, mob, elapsed, combat.DefaultMaxTurns, s.rngNew, load)
	if hr.Report.End == combat.EndTimeUp && elapsed < c.DurationSecs {
		hr.Report.End = combat.EndEarly // collected before the planned end (DD-308)
	}

	// --- Item mint (idempotent by source_ref_id = contract.ID) --------------
	// Each LootSeconds of winning has a GearDropPct chance of gear from the
	// mob's tier, usable by the character's class (DD-229, DD-281).
	var gear activity.GearDrop
	if hr.Victories > 0 {
		sheet, err := s.chars.Sheet(ctx, c.CharacterID, accountID)
		if err != nil {
			return nil, err
		}
		gear = activity.GearDrop{Class: sheet.Class, Level: activity.DropLevel(int(mob.Level), int(snapshot.Level)),
			Rolls: hr.LootRolls, Pct: GearDropPct}
	}
	var granted []activity.GrantedItem
	if s.items != nil && (hr.CommonDrops+hr.RareDrops+hr.LegendaryDrops+gear.Rolls) > 0 {
		granted, err = s.items.Grant(ctx, c.CharacterID, c.ID, c.Seed, activity.Drops{
			Kind: activity.KindHunting, MobLevel: activity.DropLevel(int(mob.Level), int(snapshot.Level)),
			Common: hr.CommonDrops, Rare: hr.RareDrops, Legendary: hr.LegendaryDrops, Gear: gear,
		})
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
		fightEnd := c.StartedAt.Add(time.Duration(hr.ElapsedSecs) * time.Second)
		if fightEnd.After(now) {
			fightEnd = now
		}
		if err := s.hp.SetHP(ctx, c.CharacterID, hr.PlayerHPFinal, fightEnd); err != nil {
			return nil, fmt.Errorf("hunt reveal hp: %w", err)
		}
	}
	// Marking it revealed uses up the supplies the roll drank and returns the
	// rest to the bag, in one transaction (DD-288).
	if err := s.repo.MarkRevealed(ctx, c.ID, resultJSON, txnID, now, hr.SuppliesUsed); err != nil {
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
	var party *partyxp.Outcome
	xp := hr.TotalXP
	if shared := partyxp.BoostIf(ctx, s.party, partyxp.Work{
		Kind: partyxp.KindHunt, Character: c.CharacterID,
		Start: c.StartedAt, End: c.StartedAt.Add(time.Duration(elapsed) * time.Second), XP: xp,
	}); len(shared.Partners) > 0 {
		party, xp = &shared, shared.Kept
	}
	if s.xp != nil && xp > 0 {
		xpResult, err = s.xp.GrantXP(ctx, c.CharacterID, xp, character.XPSource{
			Kind: "hunt", ID: c.ID.String(), Metadata: map[string]any{"mob": c.MobCode, "victories": hr.Victories}})
		if err != nil {
			s.log.WarnContext(ctx, "hunt: xp grant failed (contract already revealed, no retry)",
				slog.String("hunt", c.ID.String()), slog.String("err", err.Error()))
		}
	}

	progression.RecordIf(ctx, s.progress, c.CharacterID,
		progression.Event{Kind: progression.EventMobKilled, Amount: int64(hr.Victories)},
		progression.Event{Kind: progression.EventMobKilled + ":" + c.MobCode, Amount: int64(hr.Victories)})

	notification.Send(ctx, s.notifier, s.log, notification.Reward{
		AccountID: accountID, CharacterID: c.CharacterID, Ref: c.ID,
		Source: "hunt", Target: c.MobCode,
		GoldCopper: hr.TotalGoldCopper, XP: xp, Items: activity.RewardItems(granted),
		NewLevel: newLevel(xpResult),
		Extra:    map[string]any{"victories": hr.Victories, "defeats": hr.Defeats},
	}.Notice())
	return &RevealResult{
		Hunt: hr, Elapsed: elapsed, LedgerTxnID: txnID,
		Items: granted, XPGain: xpResult, Replayed: false, Party: party,
	}, nil
}

// --- Stop -------------------------------------------------------------------

// Stop cancels an active hunt with no reward; its supplies go back to the
// bag untouched. Idempotent for already-cancelled; returns ErrNotActive for
// revealed contracts.
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

// newLevel is the level reached when the reveal levelled up, else 0.
func newLevel(x *character.XPGrantResult) int {
	if x == nil || !x.LeveledUp || x.Character == nil {
		return 0
	}
	return int(x.Character.Level)
}

// HuntXP is what one kill of a mob pays a character: the mob's own XP
// scaled by the level gap either way (DD-281).
func HuntXP(mobXP int32, playerLevel, mobLevel int) int32 {
	return int32(math.Round(float64(mobXP) * leveling.Current().HuntLevelXP.Factor(playerLevel, mobLevel)))
}
