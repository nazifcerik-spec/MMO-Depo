package crafting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/boost"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/partyxp"
	"github.com/otter/server/internal/progression"
)

// --- Ports (narrow — crafting is a coordinator of existing services) --------

// CatalogLookup is the crafting catalog. CatalogService satisfies.
type CatalogLookup interface {
	Catalog() *Catalog
}

// InventoryConsumer is the material-eating capability. inventory.Service.ConsumeByTemplate satisfies.
type InventoryConsumer interface {
	Available(ctx context.Context, characterID uuid.UUID, templateCode string) (int32, error)
	ConsumeAll(ctx context.Context, characterID uuid.UUID, needs []inventory.Need) error
	// HasRoom: is there bag room for the output (DD-246)?
	HasRoom(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32) (bool, error)
}

// ItemMinter is the output-producing capability. item.Service.MintExact satisfies.
type ItemMinter interface {
	MintExact(ctx context.Context, characterID uuid.UUID, templateCode string, stackQuantity int32, sourceKind string, sourceRefID uuid.UUID) (*item.Instance, error)
}

// CharacterOwner is ownership + label lookup. character.Service satisfies.
type CharacterOwner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

// ServiceConfig — clock injection.
type ServiceConfig struct {
	Now func() time.Time
}

func DefaultServiceConfig() ServiceConfig { return ServiceConfig{Now: time.Now} }

// Service is the crafting business logic — coordinates inventory, item, ledger,
// and profession XP writes in one server-authoritative flow.
type Service struct {
	repo     Repository
	catalog  CatalogLookup
	inv      InventoryConsumer
	items    ItemMinter
	txn      economy.TxnService
	chars    CharacterOwner
	cfg      ServiceConfig
	log      *slog.Logger
	progress progression.Recorder
	speed    SpeedSource     // worn enchants' crafting speed (DD-293); nil = none
	party    partyxp.Booster // party XP bonus (DD-305); nil = none
	boost    boost.Applier   // timed XP boosts (DD-320); nil = none
}

// WithBoost applies timed XP boosts (DD-320).
func (s *Service) WithBoost(b boost.Applier) *Service { s.boost = b; return s }

// WithPartyBonus adds the party bonus for members crafting in the same
// profession alongside (DD-305).
func (s *Service) WithPartyBonus(p partyxp.Booster) *Service { s.party = p; return s }

// WithRecorder feeds crafts into progression (Sprint 24).
func (s *Service) WithRecorder(r progression.Recorder) *Service { s.progress = r; return s }

// SpeedSource reports a profession bonus from worn gear
// (inventory.EquipmentService): here the "craft_speed" share.
type SpeedSource interface {
	ProfessionBonus(ctx context.Context, characterID uuid.UUID, target string) (float64, error)
}

// WithSpeed makes worn crafting-speed enchants shorten new jobs (DD-293).
func (s *Service) WithSpeed(src SpeedSource) *Service { s.speed = src; return s }

// GrantXP adds profession XP from outside a craft: a party member's share
// (DD-296). The profession's level cap applies.
func (s *Service) GrantXP(ctx context.Context, characterID uuid.UUID, profession string, xp int32) error {
	maxLevel := int16(100)
	if p, ok := s.catalog.Catalog().Profession(profession); ok {
		maxLevel = p.MaxLevel
	}
	_, _, err := s.repo.AddXP(ctx, characterID, profession, xp, maxLevel)
	return err
}

// MaxCraftSpeed caps how much enchants shorten a craft.
const MaxCraftSpeed = 0.5

func NewService(
	repo Repository,
	catalog CatalogLookup,
	inv InventoryConsumer,
	items ItemMinter,
	txn economy.TxnService,
	chars CharacterOwner,
	cfg ServiceConfig,
	logger *slog.Logger,
) *Service {
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo: repo, catalog: catalog, inv: inv, items: items, txn: txn,
		chars: chars, cfg: cfg, log: logger,
	}
}

// --- Learn / List ----------------------------------------------------------

func (s *Service) Learn(ctx context.Context, characterID, accountID uuid.UUID, professionCode string) (*CharacterProfession, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	if p, ok := s.catalog.Catalog().Profession(professionCode); !ok || p.Kind == KindGathering {
		return nil, ErrUnknownProfession // gathering levels come from activities (DD-244)
	}
	return s.repo.UpsertLearn(ctx, characterID, professionCode, s.cfg.Now().UTC())
}

func (s *Service) List(ctx context.Context, characterID, accountID uuid.UUID) ([]CharacterProfession, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, characterID)
}

// --- Timed crafting (DD-264) -------------------------------------------------

// Start begins a job of count units. Everything the whole job needs is
// checked first; the job row claims the character's one action (DD-263),
// then the gold and inputs are taken. A start that cannot pay undoes itself.
func (s *Service) Start(ctx context.Context, characterID, accountID uuid.UUID, recipeCode string, count int) (*Job, error) {
	if count < 1 || count > MaxCraftCount {
		return nil, fmt.Errorf("%w: count must be 1..%d", ErrValidation, MaxCraftCount)
	}
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	cat := s.catalog.Catalog()
	recipe, ok := cat.Recipe(recipeCode)
	if !ok {
		return nil, ErrUnknownRecipe
	}
	if _, ok := cat.Profession(recipe.ProfessionCode); !ok {
		return nil, fmt.Errorf("%w: recipe %s references profession %q not in catalog",
			ErrUnknownProfession, recipe.Code, recipe.ProfessionCode)
	}
	cp, err := s.repo.Get(ctx, characterID, recipe.ProfessionCode)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrProfessionNotLearned
	}
	if err != nil {
		return nil, err
	}
	if cp.Level < recipe.MinSkill {
		return nil, ErrSkillTooLow
	}
	// Only what the character could use (DD-316).
	if lv, ok := s.chars.(interface {
		Level(ctx context.Context, characterID uuid.UUID) (int, error)
	}); ok && recipe.OutputMinLevel > 1 {
		level, err := lv.Level(ctx, characterID)
		if err != nil {
			return nil, err
		}
		if level < int(recipe.OutputMinLevel) {
			return nil, ErrCharLevelTooLow
		}
	}
	// The XP each unit pays follows the recipe's difficulty for the skill now
	// (DD-294); the job freezes it with the recipe.
	_, factor := Difficulty(recipe.MinSkill, cp.Level)
	recipe.ProfessionXPReward = int32(math.Round(float64(recipe.ProfessionXPReward) * factor))
	unit := recipe.CraftSeconds
	if unit < 1 {
		unit = DefaultCraftSeconds
	}
	if s.speed != nil {
		speed, err := s.speed.ProfessionBonus(ctx, characterID, "craft_speed")
		if err != nil {
			return nil, err
		}
		unit = max(1, int32(math.Round(float64(unit)*(1-min(speed, MaxCraftSpeed)))))
	}
	if int64(unit)*int64(count) > MaxJobSeconds {
		return nil, ErrJobTooLong
	}
	needs := make([]inventory.Need, 0, len(recipe.Inputs))
	for _, in := range recipe.Inputs {
		need := in.Quantity * int32(count)
		have, err := s.inv.Available(ctx, characterID, in.TemplateCode)
		if err != nil {
			return nil, err
		}
		if have < need {
			return nil, fmt.Errorf("craft consume %s x%d (have %d): %w", in.TemplateCode, need, have, inventory.ErrInsufficient)
		}
		needs = append(needs, inventory.Need{TemplateCode: in.TemplateCode, Quantity: need})
	}
	// A full bag refuses the job before anything is spent (DD-246).
	if ok, err := s.inv.HasRoom(ctx, characterID, recipe.OutputTemplateCode, recipe.OutputStack); err != nil {
		return nil, err
	} else if !ok {
		return nil, inventory.ErrBagFull
	}

	now := s.cfg.Now().UTC()
	job := &Job{ID: uuid.New(), CharacterID: characterID, Recipe: recipe, Count: count, UnitSeconds: unit,
		StartedAt: now, ReadyAt: now.Add(time.Duration(unit) * time.Duration(count) * time.Second), State: JobActive}
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}
	cost := recipe.GoldCostCopper * int64(count)
	var paid bool
	if cost > 0 {
		_, err := s.txn.Write(ctx, economy.Intent{
			Kind: "crafting_craft", RefKind: "crafting_job", RefID: job.ID,
			IdempotencyKey: "craft-job-" + job.ID.String(),
			Entries:        goldMove(characterID, -cost),
		})
		if err != nil {
			s.undoStart(ctx, job, false, cost)
			if errors.Is(err, economy.ErrInsufficientBalance) {
				return nil, ErrInsufficientGold
			}
			return nil, fmt.Errorf("craft gold debit: %w", err)
		}
		paid = true
	}
	// One transaction for every input; losing a race with another spend
	// after the pre-check undoes the job and refunds the gold.
	if err := s.inv.ConsumeAll(ctx, characterID, needs); err != nil {
		s.undoStart(ctx, job, paid, cost)
		return nil, fmt.Errorf("craft consume: %w", err)
	}
	return job, nil
}

// goldMove moves delta copper into (negative: out of) the character's purse
// against the vendor sink.
func goldMove(characterID uuid.UUID, delta int64) []economy.Entry {
	return []economy.Entry{
		{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: characterID, Delta: delta},
		{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: -delta},
	}
}

func (s *Service) undoStart(ctx context.Context, job *Job, paid bool, cost int64) {
	if paid {
		if _, err := s.txn.Write(ctx, economy.Intent{
			Kind: "crafting_refund", RefKind: "crafting_job", RefID: job.ID,
			IdempotencyKey: "craft-start-refund-" + job.ID.String(), Entries: goldMove(job.CharacterID, cost),
		}); err != nil {
			s.log.ErrorContext(ctx, "crafting: gold refund failed — manual reconcile needed",
				slog.String("job", job.ID.String()), slog.String("err", err.Error()))
		}
	}
	if err := s.repo.DeleteJob(ctx, job.ID); err != nil {
		s.log.ErrorContext(ctx, "crafting: undo job failed", slog.String("job", job.ID.String()), slog.String("err", err.Error()))
	}
}

// Active is the character's running job (nil when idle).
func (s *Service) Active(ctx context.Context, characterID, accountID uuid.UUID) (*Job, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	j, err := s.repo.ActiveJob(ctx, characterID)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return j, err
}

// Collect ends the job: finished units are minted and pay XP, the rest is
// refunded. A second call returns the stored result.
func (s *Service) Collect(ctx context.Context, characterID, accountID, jobID uuid.UUID) (*CraftResult, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job.CharacterID != characterID {
		return nil, ErrNotFound
	}
	if job.State != JobActive {
		return replay(job), nil
	}
	now := s.cfg.Now().UTC()
	made := job.Done(now)
	first, err := s.repo.FinishJob(ctx, job.ID, made, now)
	if err != nil {
		return nil, err
	}
	if !first { // a parallel collect won
		if job, err = s.repo.GetJob(ctx, jobID); err != nil {
			return nil, err
		}
		return replay(job), nil
	}

	rec := job.Recipe
	res := &CraftResult{OutputTemplateCode: rec.OutputTemplateCode, ProfessionCode: rec.ProfessionCode,
		Count: job.Count, Made: made, Refunded: job.Count - made}
	for i := 0; i < made; i++ {
		out, err := s.items.MintExact(ctx, characterID, rec.OutputTemplateCode, rec.OutputStack, "craft", job.ID)
		if err != nil {
			return nil, fmt.Errorf("craft mint: %w", err)
		}
		res.OutputInstanceID = out.ID
		res.OutputStack += out.StackQuantity
	}
	if rest := job.Count - made; rest > 0 {
		if err := s.refund(ctx, job, rest); err != nil {
			return nil, err
		}
	}
	if made > 0 {
		maxLevel := int16(100)
		if p, ok := s.catalog.Catalog().Profession(rec.ProfessionCode); ok {
			maxLevel = p.MaxLevel
		}
		before, err := s.repo.Get(ctx, characterID, rec.ProfessionCode)
		if err == nil {
			res.SkillBefore, res.SkillAfter = before.Level, before.Level
		}
		res.XPGained = rec.ProfessionXPReward * int32(made)
		end := job.StartedAt.Add(time.Duration(made*int(job.UnitSeconds)) * time.Second)
		if shared := partyxp.BoostIf(ctx, s.party, partyxp.Work{
			Kind: partyxp.CraftPrefix + rec.ProfessionCode, Character: characterID,
			Start: job.StartedAt, End: end, XP: int64(res.XPGained),
		}); len(shared.Partners) > 0 {
			res.Party, res.XPGained = &shared, int32(shared.Kept)
		}
		res.XPGained = int32(boost.ApplyIf(ctx, s.boost, characterID, boost.CatCraft, job.StartedAt, end,
			int64(res.XPGained))) // DD-320
		updated, leveled, err := s.repo.AddXP(ctx, characterID, rec.ProfessionCode, res.XPGained, maxLevel)
		if err != nil {
			s.log.WarnContext(ctx, "crafting: xp grant failed after a collected job",
				slog.String("job", job.ID.String()), slog.String("err", err.Error()))
		} else {
			res.SkillAfter, res.LeveledUp = updated.Level, leveled
		}
		out := res.OutputInstanceID
		if err := s.repo.InsertCraftLog(ctx, characterID, rec.Code, &out, nil, now); err != nil {
			s.log.WarnContext(ctx, "crafting: audit log insert failed", slog.String("err", err.Error()))
		}
		progression.RecordIf(ctx, s.progress, characterID,
			progression.Event{Kind: progression.EventItemCrafted, Amount: int64(made)},
			progression.Event{Kind: progression.EventItemCrafted + ":" + rec.OutputTemplateCode, Amount: int64(made)})
	}
	if body, err := json.Marshal(res); err == nil {
		if err := s.repo.SaveJobResult(ctx, job.ID, body); err != nil {
			s.log.WarnContext(ctx, "crafting: result not stored", slog.String("err", err.Error()))
		}
	}
	return res, nil
}

// refund returns the gold and inputs of rest unfinished units.
func (s *Service) refund(ctx context.Context, job *Job, rest int) error {
	rec := job.Recipe
	if cost := rec.GoldCostCopper * int64(rest); cost > 0 {
		if _, err := s.txn.Write(ctx, economy.Intent{
			Kind: "crafting_refund", RefKind: "crafting_job", RefID: job.ID,
			IdempotencyKey: "craft-refund-" + job.ID.String(), Entries: goldMove(job.CharacterID, cost),
		}); err != nil {
			return fmt.Errorf("craft refund gold: %w", err)
		}
	}
	for _, in := range rec.Inputs {
		for left := in.Quantity * int32(rest); left > 0; {
			got, err := s.items.MintExact(ctx, job.CharacterID, in.TemplateCode, left, "craft_refund", job.ID)
			if err != nil {
				return fmt.Errorf("craft refund %s: %w", in.TemplateCode, err)
			}
			left -= max(got.StackQuantity, 1)
		}
	}
	return nil
}

func replay(job *Job) *CraftResult {
	res := &CraftResult{OutputTemplateCode: job.Recipe.OutputTemplateCode, ProfessionCode: job.Recipe.ProfessionCode,
		Count: job.Count}
	_ = json.Unmarshal(job.Result, res)
	res.Replayed = true
	return res
}

// --- Public catalog readers ------------------------------------------------

func (s *Service) ListProfessions() []Profession { return s.catalog.Catalog().AllProfessions() }
func (s *Service) ListRecipes() []Recipe         { return s.catalog.Catalog().AllRecipes() }
