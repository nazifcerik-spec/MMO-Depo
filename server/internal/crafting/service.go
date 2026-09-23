package crafting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/progression"
)

// --- Ports (narrow — crafting is a coordinator of existing services) --------

// CatalogLookup is the crafting catalog. CatalogService satisfies.
type CatalogLookup interface {
	Catalog() *Catalog
}

// InventoryConsumer is the material-eating capability. inventory.Service.ConsumeByTemplate satisfies.
type InventoryConsumer interface {
	ConsumeByTemplate(ctx context.Context, characterID uuid.UUID, templateCode string, quantity int32) error
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
}

// WithRecorder feeds crafts into progression (Sprint 24).
func (s *Service) WithRecorder(r progression.Recorder) *Service { s.progress = r; return s }

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
	if _, ok := s.catalog.Catalog().Profession(professionCode); !ok {
		return nil, ErrUnknownProfession
	}
	return s.repo.UpsertLearn(ctx, characterID, professionCode, s.cfg.Now().UTC())
}

func (s *Service) List(ctx context.Context, characterID, accountID uuid.UUID) ([]CharacterProfession, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, characterID)
}

// --- Craft (the reveal-shaped orchestrator) --------------------------------

func (s *Service) Craft(ctx context.Context, characterID, accountID uuid.UUID, recipeCode string) (*CraftResult, error) {
	if err := s.chars.Owns(ctx, characterID, accountID); err != nil {
		return nil, err
	}
	cat := s.catalog.Catalog()
	recipe, ok := cat.Recipe(recipeCode)
	if !ok {
		return nil, ErrUnknownRecipe
	}
	prof, ok := cat.Profession(recipe.ProfessionCode)
	if !ok {
		// Content bug: recipe references a profession that no longer exists.
		return nil, fmt.Errorf("%w: recipe %s references profession %q not in catalog",
			ErrUnknownProfession, recipe.Code, recipe.ProfessionCode)
	}

	// Skill gate.
	cp, err := s.repo.Get(ctx, characterID, recipe.ProfessionCode)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrProfessionNotLearned
		}
		return nil, err
	}
	if cp.Level < recipe.MinSkill {
		return nil, ErrSkillTooLow
	}

	// --- Debit gold first (idempotent ledger key = "craft-<recipe>-<char>-<now>")
	// v1 uses a fresh key per attempt; retries within the same HTTP request are
	// deduped one layer higher via the Idempotency middleware.
	now := s.cfg.Now().UTC()
	var goldTxnID uuid.UUID
	if recipe.GoldCostCopper > 0 {
		res, err := s.txn.Write(ctx, economy.Intent{
			Kind:           "crafting_craft",
			RefKind:        "recipe",
			RefID:          uuid.New(), // per-craft unique ref
			IdempotencyKey: fmt.Sprintf("craft-%s-%s-%d", recipe.Code, characterID.String(), now.UnixNano()),
			Entries: []economy.Entry{
				{Currency: "gold", AccountKind: economy.AccountCharacter,
					AccountRef: characterID, Delta: -recipe.GoldCostCopper},
				{Currency: "gold", AccountKind: economy.AccountSystem,
					AccountRef: economy.SystemVendorSink, Delta: recipe.GoldCostCopper},
			},
		})
		if err != nil {
			if errors.Is(err, economy.ErrInsufficientBalance) {
				return nil, ErrInsufficientGold
			}
			return nil, fmt.Errorf("craft gold debit: %w", err)
		}
		goldTxnID = res.TxnID
	}

	// --- Consume inputs (atomic; refuses if any shortfall).
	// If a shortfall triggers after gold was already debited, we have a stuck
	// state: gold gone, no output. Sprint 17 v1 accepts this as a known edge
	// (the ledger keeps the audit trail; support tools can refund). Phase 02
	// wraps the whole flow in a single transaction across services.
	for _, in := range recipe.Inputs {
		if err := s.inv.ConsumeByTemplate(ctx, characterID, in.TemplateCode, in.Quantity); err != nil {
			return nil, fmt.Errorf("craft consume %s x%d: %w", in.TemplateCode, in.Quantity, err)
		}
	}

	// --- Mint output.
	out, err := s.items.MintExact(ctx, characterID, recipe.OutputTemplateCode, recipe.OutputStack, "craft", uuid.New())
	if err != nil {
		return nil, fmt.Errorf("craft mint: %w", err)
	}

	// --- Grant profession XP.
	skillBefore := cp.Level
	updated, leveled, err := s.repo.AddXP(ctx, characterID, recipe.ProfessionCode, recipe.ProfessionXPReward, prof.MaxLevel)
	if err != nil {
		s.log.WarnContext(ctx, "crafting: xp grant failed after successful craft",
			slog.String("recipe", recipe.Code), slog.String("err", err.Error()))
		updated = cp // best-effort — show pre-XP state so client still gets a result
	}

	// --- Audit log (best-effort).
	var goldPtr *uuid.UUID
	if goldTxnID != uuid.Nil {
		goldPtr = &goldTxnID
	}
	if err := s.repo.InsertCraftLog(ctx, characterID, recipe.Code, &out.ID, goldPtr, now); err != nil {
		s.log.WarnContext(ctx, "crafting: audit log insert failed", slog.String("err", err.Error()))
	}

	progression.RecordIf(ctx, s.progress, characterID,
		progression.Event{Kind: progression.EventItemCrafted, Amount: 1})

	return &CraftResult{
		OutputInstanceID:   out.ID,
		OutputTemplateCode: out.TemplateCode,
		OutputStack:        out.StackQuantity,
		GoldTxnID:          goldTxnID,
		ProfessionCode:     recipe.ProfessionCode,
		SkillBefore:        skillBefore,
		SkillAfter:         updated.Level,
		XPGained:           recipe.ProfessionXPReward,
		LeveledUp:          leveled,
	}, nil
}

// --- Public catalog readers ------------------------------------------------

func (s *Service) ListProfessions() []Profession { return s.catalog.Catalog().AllProfessions() }
func (s *Service) ListRecipes() []Recipe         { return s.catalog.Catalog().AllRecipes() }
