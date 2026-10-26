package vendor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
)

// InventoryDeleter is the delete capability the vendor needs from inventory.
// Kept as a narrow interface — inventory.Service satisfies it.
type InventoryDeleter interface {
	Delete(ctx context.Context, characterID, itemID uuid.UUID) error
	// HasRoom: can the bag take the item (DD-246)?
	HasRoom(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32) (bool, error)
}

// ensureRoom refuses a purchase the bag can't hold, before any gold moves.
func (s *Service) ensureRoom(ctx context.Context, characterID uuid.UUID, templateCode string) error {
	ok, err := s.invDel.HasRoom(ctx, characterID, templateCode, 1)
	if err != nil {
		return err
	}
	if !ok {
		return inventory.ErrBagFull
	}
	return nil
}

// ItemFetcher exposes the item read+catalog needed by vendor.
// item.Service satisfies via Get, Catalog, and MaterializeDrops (the latter
// method name is aliased here as Materialize on the interface — Go interfaces
// are satisfied by method name, so we rename to match).
type ItemFetcher interface {
	Get(ctx context.Context, id uuid.UUID) (*item.Instance, error)
	Catalog() *item.Catalog
	BySource(ctx context.Context, sourceRefID uuid.UUID) ([]item.Instance, error)
	MaterializeDrops(
		ctx context.Context,
		characterID uuid.UUID,
		sourceKind string,
		sourceRefID uuid.UUID,
		drops []item.DropSpec,
	) ([]item.Instance, error)
}

// ServiceConfig — clock injection for tests, buyback TTL override.
type ServiceConfig struct {
	Now        func() time.Time
	BuybackTTL time.Duration
}

func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{Now: time.Now, BuybackTTL: 2 * time.Hour}
}

// Service is the vendor business logic. Loads the vendor catalog at boot;
// Sprint 12 contentpack reloads it in place under catMu.
type Service struct {
	repo    Repository
	items   ItemFetcher
	invDel  InventoryDeleter
	txn     economy.TxnService
	catMu   sync.RWMutex
	catalog map[string]Template
	cfg     ServiceConfig
	log     *slog.Logger
}

// NewService constructs the service after loading the vendor catalog.
func NewService(
	ctx context.Context,
	repo Repository,
	items ItemFetcher,
	invDel InventoryDeleter,
	txn economy.TxnService,
	cfg ServiceConfig,
	logger *slog.Logger,
) (*Service, error) {
	templates, err := repo.Catalog(ctx)
	if err != nil {
		return nil, fmt.Errorf("vendor catalog load: %w", err)
	}
	m := make(map[string]Template, len(templates))
	for _, t := range templates {
		m[t.Code] = t
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.BuybackTTL == 0 {
		cfg.BuybackTTL = 2 * time.Hour
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo: repo, items: items, invDel: invDel, txn: txn,
		catalog: m, cfg: cfg, log: logger,
	}, nil
}

// lookupTemplate is the read-locked accessor for the vendor catalog. Kept
// short so callers stay concise.
func (s *Service) lookupTemplate(code string) (Template, bool) {
	s.catMu.RLock()
	t, ok := s.catalog[code]
	s.catMu.RUnlock()
	return t, ok
}

// ReloadCatalog re-reads content.vendor_template rows and atomically swaps the
// in-memory map. Called by contentpack after YAML upsert.
func (s *Service) ReloadCatalog(ctx context.Context) error {
	templates, err := s.repo.Catalog(ctx)
	if err != nil {
		return err
	}
	m := make(map[string]Template, len(templates))
	for _, t := range templates {
		m[t.Code] = t
	}
	s.catMu.Lock()
	s.catalog = m
	s.catMu.Unlock()
	s.log.InfoContext(ctx, "vendor catalog reloaded", slog.Int("vendors", len(m)))
	return nil
}

// --- Info -------------------------------------------------------------------

// VendorInfo is the compound response for GET /v1/vendor/{code}.
type VendorInfo struct {
	Code           string      `json:"code"`
	DisplayNameKey string      `json:"display_name_key"`
	VendorType     string      `json:"vendor_type"`
	Sells          []SellOffer `json:"sells"`
	SellMultiplier float64     `json:"sell_multiplier"` // player→vendor rate
}

// SellOffer is one item in vendor's inventory + computed price.
type SellOffer struct {
	TemplateCode   string `json:"template_code"`
	DisplayNameKey string `json:"display_name_key,omitempty"`
	ItemClass      string `json:"item_class,omitempty"`
	ItemSubclass   string `json:"item_subclass,omitempty"`
	BuyPriceCopper int64  `json:"buy_price_copper"`
}

// Info returns vendor + sell list with computed prices.
func (s *Service) Info(_ context.Context, code string) (*VendorInfo, error) {
	tmpl, ok := s.lookupTemplate(code)
	if !ok {
		return nil, ErrVendorNotFound
	}
	cat := s.items.Catalog()
	offers := make([]SellOffer, 0, len(tmpl.Sells))
	for _, entry := range tmpl.Sells {
		itemTmpl, err := cat.Get(entry.TemplateCode)
		if err != nil {
			continue // skip unknown template (content pack drift)
		}
		offers = append(offers, SellOffer{
			TemplateCode:   entry.TemplateCode,
			DisplayNameKey: itemTmpl.DisplayNameKey,
			ItemClass:      itemTmpl.ItemClass,
			ItemSubclass:   itemTmpl.ItemSubclass,
			BuyPriceCopper: BuyPrice(itemTmpl.BaseValueCopper, tmpl.BuyMultiplier, entry.PriceOverrideCopper),
		})
	}
	return &VendorInfo{
		Code:           tmpl.Code,
		DisplayNameKey: tmpl.DisplayNameKey,
		VendorType:     tmpl.VendorType,
		Sells:          offers,
		SellMultiplier: tmpl.SellMultiplier,
	}, nil
}

// --- Buy --------------------------------------------------------------------

// BuyRequest is the shape from the client.
type BuyRequest struct {
	VendorCode          string
	CharacterID         uuid.UUID
	TemplateCode        string
	ExpectedPriceCopper *int64 // optional price-race guard
	IdempotencyKey      string // must be non-empty; also used for ledger idem
}

// BuyResult wraps the purchase outcome.
type BuyResult struct {
	InstanceID   uuid.UUID `json:"instance_id"`
	TemplateCode string    `json:"template_code"`
	PricePaid    int64     `json:"price_paid_copper"`
	LedgerTxnID  uuid.UUID `json:"ledger_txn_id"`
}

// Buy: v1 quantity = 1. Player pays vendor's BuyPrice; new item minted with
// auto-slot (inventory.Service via item.SlotAssigner wiring).
func (s *Service) Buy(ctx context.Context, req BuyRequest) (*BuyResult, error) {
	if req.IdempotencyKey == "" {
		return nil, &ValidationError{Field: "idempotency_key", Message: "required"}
	}
	tmpl, ok := s.lookupTemplate(req.VendorCode)
	if !ok {
		return nil, ErrVendorNotFound
	}
	idx := tmpl.SellsIndex()
	entry, ok := idx[req.TemplateCode]
	if !ok {
		return nil, ErrItemNotSold
	}
	cat := s.items.Catalog()
	itemTmpl, err := cat.Get(req.TemplateCode)
	if err != nil {
		return nil, err
	}

	price := BuyPrice(itemTmpl.BaseValueCopper, tmpl.BuyMultiplier, entry.PriceOverrideCopper)
	if req.ExpectedPriceCopper != nil && *req.ExpectedPriceCopper != price {
		return nil, &ValidationError{
			Field:   "expected_price_copper",
			Message: fmt.Sprintf("price is %d, client sent %d", price, *req.ExpectedPriceCopper),
		}
	}
	if price <= 0 {
		return nil, &ValidationError{Field: "price", Message: "vendor price must be positive"}
	}
	if err := s.ensureRoom(ctx, req.CharacterID, req.TemplateCode); err != nil {
		return nil, err
	}

	// Debit ledger FIRST (idempotent via key). Item mint follows; if it fails,
	// the retry hits ledger idempotency (no double debit) + re-attempts mint.
	ledgerRes, err := s.txn.Write(ctx, economy.Intent{
		Kind:           "vendor_buy",
		RefKind:        "vendor",
		RefID:          uuid.Nil, // vendor is a code, not a UUID; use txn key for trace
		IdempotencyKey: "vendor-buy-" + req.IdempotencyKey,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: req.CharacterID, Delta: -price},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: price},
		},
	})
	if err != nil {
		return nil, err
	}

	// Materialize the item once per ledger txn: a replay of the same key gets
	// the same txn back and must not mint a second (free) item.
	drops := []item.DropSpec{{
		TemplateCode: req.TemplateCode,
		Rarity:       item.RarityCommon, // vendor stock is base rarity in v1
	}}
	instances, err := s.mintOnce(ctx, req.CharacterID, "vendor_buy", ledgerRes.TxnID, drops)
	if err != nil {
		return nil, fmt.Errorf("vendor buy materialize: %w", err)
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("vendor buy: item materialization returned no instances")
	}
	return &BuyResult{
		InstanceID:   instances[0].ID,
		TemplateCode: instances[0].TemplateCode,
		PricePaid:    price,
		LedgerTxnID:  ledgerRes.TxnID,
	}, nil
}

// mintOnce materializes the drops for a source reference exactly once: a
// replay (or a concurrent twin, which loses on the partial unique index of
// migration 0028) gets the instances already minted for that reference.
func (s *Service) mintOnce(ctx context.Context, characterID uuid.UUID, kind string, ref uuid.UUID, drops []item.DropSpec) ([]item.Instance, error) {
	if existing, err := s.items.BySource(ctx, ref); err != nil {
		return nil, err
	} else if len(existing) > 0 {
		return existing, nil
	}
	instances, err := s.items.MaterializeDrops(ctx, characterID, kind, ref, drops)
	if err != nil {
		if existing, lerr := s.items.BySource(ctx, ref); lerr == nil && len(existing) > 0 {
			return existing, nil
		}
		return nil, err
	}
	return instances, nil
}

// --- Sell -------------------------------------------------------------------

type SellRequest struct {
	VendorCode     string
	CharacterID    uuid.UUID
	ItemInstanceID uuid.UUID
	IdempotencyKey string
}

type SellResult struct {
	InstanceID       uuid.UUID `json:"instance_id"`
	TemplateCode     string    `json:"template_code"`
	Received         int64     `json:"received_copper"`
	BuybackID        uuid.UUID `json:"buyback_id"`
	BuybackExpiresAt time.Time `json:"buyback_expires_at"`
	LedgerTxnID      uuid.UUID `json:"ledger_txn_id"`
}

// Sell removes the (whole) item from inventory, snapshots it into buyback
// (2h TTL), and credits gold to the character wallet.
func (s *Service) Sell(ctx context.Context, req SellRequest) (*SellResult, error) {
	if req.IdempotencyKey == "" {
		return nil, &ValidationError{Field: "idempotency_key", Message: "required"}
	}
	tmpl, ok := s.lookupTemplate(req.VendorCode)
	if !ok {
		return nil, ErrVendorNotFound
	}
	inst, err := s.items.Get(ctx, req.ItemInstanceID)
	if err != nil {
		return nil, ErrItemNotOwned
	}
	if inst.OwnerCharacterID == nil || *inst.OwnerCharacterID != req.CharacterID {
		return nil, ErrItemNotOwned
	}
	if inst.IsLocked {
		return nil, ErrItemLocked
	}
	if inst.IsEquipped() {
		return nil, ErrItemLocked // can't sell equipped items — unequip first
	}
	// Bound items (e.g. worn once) still sell to a vendor (DD-245); only
	// player-to-player trade stays closed to them.

	itemTmpl, err := s.items.Catalog().Get(inst.TemplateCode)
	if err != nil {
		return nil, err
	}
	if itemTmpl.NoSell {
		return nil, ErrItemNoSell
	}

	// Price per unit × stack.
	pricePerUnit := SellPrice(itemTmpl.BaseValueCopper, tmpl.SellMultiplier, item.ValueMultiplier(inst.Rarity))
	totalPrice := pricePerUnit * int64(inst.StackQuantity)
	if totalPrice <= 0 {
		totalPrice = 1 // minimum 1 copper for non-worthless items
	}

	now := s.cfg.Now().UTC()

	// Snapshot BEFORE delete for buyback.
	snap := ItemSnapshot{
		TemplateCode:  inst.TemplateCode,
		Rarity:        string(inst.Rarity),
		Prefix:        inst.PrefixCode,
		Suffix:        inst.SuffixCode,
		StackQuantity: inst.StackQuantity,
		IsBound:       inst.IsBound,
		Enchant:       inst.EnchantCode,
	}
	snapJSON, _ := json.Marshal(snap)

	buybackID := uuid.New()
	if err := s.repo.InsertBuyback(ctx, BuybackEntry{
		ID:              buybackID,
		CharacterID:     req.CharacterID,
		VendorCode:      req.VendorCode,
		ItemSnapshot:    snapJSON,
		SoldPriceCopper: totalPrice,
		SoldAt:          now,
		ExpiresAt:       now.Add(s.cfg.BuybackTTL),
	}); err != nil {
		return nil, err
	}

	// Delete item from inventory.
	if err := s.invDel.Delete(ctx, req.CharacterID, inst.ID); err != nil {
		return nil, err
	}

	// Credit gold.
	ledgerRes, err := s.txn.Write(ctx, economy.Intent{
		Kind:           "vendor_sell",
		RefKind:        "vendor_buyback",
		RefID:          buybackID,
		IdempotencyKey: "vendor-sell-" + req.IdempotencyKey,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -totalPrice},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: req.CharacterID, Delta: totalPrice},
		},
	})
	if err != nil {
		return nil, err
	}

	return &SellResult{
		InstanceID:       inst.ID,
		TemplateCode:     inst.TemplateCode,
		Received:         totalPrice,
		BuybackID:        buybackID,
		BuybackExpiresAt: now.Add(s.cfg.BuybackTTL),
		LedgerTxnID:      ledgerRes.TxnID,
	}, nil
}

// --- Buyback ----------------------------------------------------------------

type BuybackResult struct {
	InstanceID   uuid.UUID `json:"instance_id"`
	TemplateCode string    `json:"template_code"`
	PricePaid    int64     `json:"price_paid_copper"`
	LedgerTxnID  uuid.UUID `json:"ledger_txn_id"`
}

// Buyback restores a previously sold item (within TTL) at the same price paid.
func (s *Service) Buyback(ctx context.Context, characterID, buybackID uuid.UUID, idemKey string) (*BuybackResult, error) {
	if idemKey == "" {
		return nil, &ValidationError{Field: "idempotency_key", Message: "required"}
	}
	entry, err := s.repo.GetBuyback(ctx, buybackID)
	if err != nil {
		return nil, err
	}
	if entry.CharacterID != characterID {
		return nil, ErrBuybackNotFound
	}
	if entry.ClaimedAt != nil {
		return nil, ErrBuybackClaimed
	}
	now := s.cfg.Now().UTC()
	if now.After(entry.ExpiresAt) {
		return nil, ErrBuybackExpired
	}
	var snapCode struct {
		TemplateCode string `json:"template_code"`
	}
	_ = json.Unmarshal(entry.ItemSnapshot, &snapCode)
	if err := s.ensureRoom(ctx, characterID, snapCode.TemplateCode); err != nil {
		return nil, err
	}

	// Debit gold once per buyback entry (not per request): two concurrent
	// requests with different keys must not both charge.
	ledgerRes, err := s.txn.Write(ctx, economy.Intent{
		Kind:           "vendor_buyback",
		RefKind:        "vendor_buyback",
		RefID:          entry.ID,
		IdempotencyKey: "vendor-buyback-" + entry.ID.String(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: characterID, Delta: -entry.SoldPriceCopper},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: entry.SoldPriceCopper},
		},
	})
	if err != nil {
		return nil, err
	}

	// Recreate the item from its snapshot (once per entry; auto-slots via the
	// inventory SlotAssigner wired to item.Service), then mark the entry
	// claimed last so a failed mint can be retried.
	var snap ItemSnapshot
	if err := json.Unmarshal(entry.ItemSnapshot, &snap); err != nil {
		return nil, fmt.Errorf("buyback snapshot decode: %w", err)
	}
	drops := []item.DropSpec{{
		TemplateCode: snap.TemplateCode,
		Rarity:       item.Rarity(snap.Rarity),
		Prefix:       snap.Prefix,
		Suffix:       snap.Suffix,
	}}
	instances, err := s.mintOnce(ctx, characterID, "vendor_buyback", entry.ID, drops)
	if err != nil {
		return nil, fmt.Errorf("buyback materialize: %w", err)
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("buyback: item materialization returned no instances")
	}
	if snap.IsBound {
		if err := s.repo.Rebind(ctx, instances[0].ID); err != nil {
			return nil, fmt.Errorf("buyback rebind: %w", err)
		}
	}
	if snap.Enchant != nil {
		if err := s.repo.SetEnchant(ctx, instances[0].ID, *snap.Enchant); err != nil {
			return nil, fmt.Errorf("buyback enchant: %w", err)
		}
	}
	if err := s.repo.ClaimBuyback(ctx, entry.ID, now); err != nil && !errors.Is(err, ErrBuybackClaimed) {
		return nil, err
	}
	return &BuybackResult{
		InstanceID:   instances[0].ID,
		TemplateCode: instances[0].TemplateCode,
		PricePaid:    entry.SoldPriceCopper,
		LedgerTxnID:  ledgerRes.TxnID,
	}, nil
}

// --- List buybacks ----------------------------------------------------------

func (s *Service) ListBuybacks(ctx context.Context, characterID uuid.UUID) ([]BuybackEntry, error) {
	return s.repo.ListBuybacks(ctx, characterID, 20)
}
