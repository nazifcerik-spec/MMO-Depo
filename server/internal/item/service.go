package item

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/activity"
)

// Ensure context is referenced when SlotAssigner interface is defined.
var _ context.Context = context.Background()

// SlotAssigner assigns bag slots (with merge for stackables) to newly created
// item instances. Satisfied by inventory.Service. Optional — if not wired,
// items are created without slot and remain "detached" (visible via
// GET /v1/character/{id}/items but not in inventory grid).
type SlotAssigner interface {
	AssignSlots(ctx context.Context, ownerID uuid.UUID, instances []Instance) error
}

// Service materializes DropSpecs (from RollDrops) into persisted Instances +
// exposes read helpers for the inventory view.
//
// catPtr is atomically swapped by ReloadCatalog so Sprint 12 hot-reload does
// not need to lock every read. Store the address of a *Catalog inside so we
// never write nil or narrow interfaces through the pointer.
type Service struct {
	repo    Repository
	catPtr  atomic.Pointer[Catalog]
	slotter SlotAssigner
	log     *slog.Logger
	now     func() time.Time
}

func NewService(repo Repository, cat *Catalog, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Service{repo: repo, log: logger, now: time.Now}
	s.catPtr.Store(cat)
	return s
}

// ReloadCatalog re-reads the item template catalog from Postgres and atomically
// swaps it in. Called by contentpack loader after YAML upsert so downstream
// callers (item.Service.Catalog, vendor.Info, etc.) see the fresh templates
// without a restart.
func (s *Service) ReloadCatalog(ctx context.Context, pool *pgxpool.Pool) error {
	fresh, err := LoadCatalog(ctx, pool)
	if err != nil {
		return err
	}
	s.catPtr.Store(fresh)
	s.log.InfoContext(ctx, "item catalog reloaded", slog.Int("templates", len(fresh.All())))
	return nil
}

// WithSlotAssigner wires a SlotAssigner. Returns the service for chaining.
func (s *Service) WithSlotAssigner(a SlotAssigner) *Service {
	s.slotter = a
	return s
}

// Catalog returns the loaded template catalog for callers that need template
// metadata (e.g. HTTP handlers rendering display names, vendors resolving
// price).
func (s *Service) Catalog() *Catalog { return s.catPtr.Load() }

// Get returns a single item instance by id. Delegates to the repo — no
// ownership check here (callers with authenticated character context validate).
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Instance, error) {
	return s.repo.Get(ctx, id)
}

// MaterializeDrops persists a slice of DropSpecs as Instances owned by the
// given character with a shared source ref. Uses batch INSERT. Idempotency is
// the caller's job — activity reveal already dedupes at contract level.
//
// Returns the created Instance list (with fresh IDs + timestamps).
// BySource returns the instances already minted for a source reference
// (idempotent mints check it before creating anything).
func (s *Service) BySource(ctx context.Context, sourceRefID uuid.UUID) ([]Instance, error) {
	return s.repo.ListBySource(ctx, sourceRefID)
}

func (s *Service) MaterializeDrops(
	ctx context.Context,
	characterID uuid.UUID,
	sourceKind string,
	sourceRefID uuid.UUID,
	drops []DropSpec,
) ([]Instance, error) {
	if len(drops) == 0 {
		return nil, nil
	}
	now := s.now().UTC()
	out := make([]Instance, 0, len(drops))
	src := sourceKind
	ref := sourceRefID
	owner := characterID
	for _, d := range drops {
		tmpl, err := s.catPtr.Load().Get(d.TemplateCode)
		if err != nil {
			return nil, err
		}
		out = append(out, Instance{
			ID:               uuid.New(),
			TemplateCode:     d.TemplateCode,
			OwnerCharacterID: &owner,
			StackQuantity:    1,
			Rarity:           d.Rarity,
			PrefixCode:       d.Prefix,
			SuffixCode:       d.Suffix,
			IsBound:          tmpl.BindKind == BindOnPickup,
			SourceKind:       &src,
			SourceRefID:      &ref,
			CreatedAt:        now,
			UpdatedAt:        now,
			Version:          1,
		})
	}
	if err := s.repo.CreateInstancesBatch(ctx, out); err != nil {
		return nil, err
	}
	// Best-effort slot assignment (merge stackable / find free slot). Failures
	// logged inside SlotAssigner, items remain slot=NULL (overflow).
	if s.slotter != nil {
		if err := s.slotter.AssignSlots(ctx, characterID, out); err != nil {
			s.log.WarnContext(ctx, "item: slot assign failed", slog.String("err", err.Error()))
		}
	}
	return out, nil
}

// MintExact creates ONE item instance owned by `characterID` with the given
// stack quantity and source metadata. Used by crafting (Sprint 17): the caller
// has already validated that the recipe output template exists and that the
// character is entitled to it.
//
// Unlike MaterializeDrops (which is for rolled loot with per-drop rarity),
// this path takes a plain template code and a stack size — no rarity roll,
// no affix, no drop-spec plumbing. Slot assignment happens via the wired
// SlotAssigner just like the drop path so the output lands in the bag grid.
func (s *Service) MintExact(
	ctx context.Context,
	characterID uuid.UUID,
	templateCode string,
	stackQuantity int32,
	sourceKind string,
	sourceRefID uuid.UUID,
) (*Instance, error) {
	tmpl, err := s.catPtr.Load().Get(templateCode)
	if err != nil {
		return nil, err
	}
	if stackQuantity < 1 {
		stackQuantity = 1
	}
	// Cap stack at template max so a broken recipe cannot mint an over-sized
	// stack (would then violate inventory move/split rules).
	if tmpl.MaxStack > 0 && stackQuantity > tmpl.MaxStack {
		stackQuantity = tmpl.MaxStack
	}
	now := s.now().UTC()
	src := sourceKind
	ref := sourceRefID
	owner := characterID
	inst := Instance{
		ID:               uuid.New(),
		TemplateCode:     templateCode,
		OwnerCharacterID: &owner,
		StackQuantity:    stackQuantity,
		Rarity:           RarityCommon,
		IsBound:          tmpl.BindKind == BindOnPickup,
		SourceKind:       &src,
		SourceRefID:      &ref,
		CreatedAt:        now,
		UpdatedAt:        now,
		Version:          1,
	}
	if err := s.repo.CreateInstancesBatch(ctx, []Instance{inst}); err != nil {
		return nil, err
	}
	if s.slotter != nil {
		if err := s.slotter.AssignSlots(ctx, characterID, []Instance{inst}); err != nil {
			s.log.WarnContext(ctx, "item: mint slot assign failed", slog.String("err", err.Error()))
		}
	}
	return &inst, nil
}

// ListByOwner returns the character's owned instances, newest first.
func (s *Service) ListByOwner(ctx context.Context, characterID uuid.UUID, limit int) ([]Instance, error) {
	return s.repo.ListByOwner(ctx, characterID, limit)
}

// ListBySource is used by activity reveal to expose "what dropped from this
// contract" — the reveal response can include actual instance IDs.
func (s *Service) ListBySource(ctx context.Context, sourceRefID uuid.UUID) ([]Instance, error) {
	return s.repo.ListBySource(ctx, sourceRefID)
}

// Grant satisfies activity.ItemGranter. Idempotent by sourceRefID — if any
// items are already tagged with this source_ref_id, they are returned as-is
// (no new mint). This guarantees reveal replay never double-creates loot.
func (s *Service) Grant(
	ctx context.Context,
	characterID, sourceRefID uuid.UUID,
	seed int64, kind activity.Kind,
	commonCount, rareCount, legendaryCount int,
	gear activity.GearDrop,
) ([]activity.GrantedItem, error) {
	// Idempotency check: if this source already produced items, return them.
	existing, err := s.repo.ListBySource(ctx, sourceRefID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return toGrantedItems(existing), nil
	}
	// Fresh mint.
	cat := s.catPtr.Load()
	drops := append(RollDrops(cat, seed, kind, commonCount, rareCount, legendaryCount), RollGear(cat, seed, gear)...)
	instances, err := s.MaterializeDrops(ctx, characterID, "activity_reveal", sourceRefID, drops)
	if err != nil {
		return nil, err
	}
	return toGrantedItems(instances), nil
}

func toGrantedItems(instances []Instance) []activity.GrantedItem {
	out := make([]activity.GrantedItem, 0, len(instances))
	for _, i := range instances {
		out = append(out, activity.GrantedItem{
			InstanceID:   i.ID,
			TemplateCode: i.TemplateCode,
			Rarity:       string(i.Rarity),
			Prefix:       i.PrefixCode,
			Suffix:       i.SuffixCode,
		})
	}
	return out
}
