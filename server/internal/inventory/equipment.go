package inventory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/skill"
)

// --- Equipment errors ------------------------------------------------------

var (
	ErrNotEquippable    = errors.New("inventory: item is not equippable")
	ErrClassRestriction = errors.New("inventory: class restriction on item")
	ErrLevelRequirement = errors.New("inventory: character level below item requirement")
	ErrNothingEquipped  = errors.New("inventory: no item in that equipment slot")
	ErrUnknownEquipSlot = errors.New("inventory: unknown equipment slot")
)

// KnownEquipSlots is the v1 slot allowlist. Content-driven expansion Phase 02.
var KnownEquipSlots = map[string]bool{
	"mainhand": true,
	"offhand":  true,
	"helm":     true,
	"chest":    true,
	"legs":     true,
	"boots":    true,
	"gloves":   true,
	"necklace": true,
	"ring1":    true,
	"ring2":    true,
	"trinket1": true,
	"trinket2": true,
}

// CharacterFetcher is the read capability needed for class + level checks.
type CharacterFetcher interface {
	Get(ctx context.Context, id uuid.UUID) (*characterView, error)
}

// characterView is a minimal projection of what equipment needs. Kept as an
// unexported type here — real implementations return character.Character but
// we describe only what we read.
type characterView struct {
	ID        uuid.UUID
	Level     int16
	ClassCode string
}

// SimpleCharFetcher wraps any function returning (level, class) for an id.
// This decouples the inventory package from the character package while
// keeping the check semantics.
type SimpleCharFetcher func(ctx context.Context, id uuid.UUID) (level int16, classCode string, err error)

func (f SimpleCharFetcher) Get(ctx context.Context, id uuid.UUID) (*characterView, error) {
	lvl, cls, err := f(ctx, id)
	if err != nil {
		return nil, err
	}
	return &characterView{ID: id, Level: lvl, ClassCode: cls}, nil
}

// EquipmentService adds equip/unequip on top of the base inventory service.
// Kept separate to avoid growing Service's dep graph — this one needs item
// catalog + character read.
type EquipmentService struct {
	inv      *Service
	items    ItemLookup
	charView CharacterFetcher
}

// ItemLookup is the item-side capability we depend on.
type ItemLookup interface {
	Get(ctx context.Context, id uuid.UUID) (*item.Instance, error)
	Catalog() *item.Catalog
}

// EquippedRepo extends Repository with equipment slot operations. PgRepo and
// MemoryRepo satisfy this.
type EquippedRepo interface {
	Equipped(ctx context.Context, characterID uuid.UUID) ([]BagItem, error)
	EquipItem(ctx context.Context, itemID, ownerID uuid.UUID, equipSlot string) error
	UnequipItem(ctx context.Context, itemID, ownerID uuid.UUID, targetBagSlot Slot) error
	FindEquippedInSlot(ctx context.Context, ownerID uuid.UUID, equipSlot string) (uuid.UUID, error)
}

// NewEquipmentService constructs the equipment orchestrator.
func NewEquipmentService(inv *Service, items ItemLookup, charView CharacterFetcher) *EquipmentService {
	return &EquipmentService{inv: inv, items: items, charView: charView}
}

// EquipmentView is the compound response for GET /v1/character/{id}/equipment.
type EquipmentView struct {
	CharacterID uuid.UUID          `json:"character_id"`
	Slots       map[string]BagItem `json:"slots"`
}

// View returns the map of equipped items by slot. Empty slots omitted.
func (s *EquipmentService) View(ctx context.Context, characterID uuid.UUID) (*EquipmentView, error) {
	repo, ok := s.inv.repo.(EquippedRepo)
	if !ok {
		return nil, fmt.Errorf("inventory: repository does not implement EquippedRepo")
	}
	list, err := repo.Equipped(ctx, characterID)
	if err != nil {
		return nil, err
	}
	slots := map[string]BagItem{}
	for _, b := range list {
		if b.EquipmentSlot != nil {
			slots[*b.EquipmentSlot] = b
		}
	}
	return &EquipmentView{CharacterID: characterID, Slots: slots}, nil
}

// Equip moves an item from bag → equipment slot. Auto-detects target slot
// from template.equip_slot. If a different item is already equipped in that
// slot, it is unequipped to a free bag slot first (transactional in service).
func (s *EquipmentService) Equip(ctx context.Context, characterID, itemID uuid.UUID) (equipSlot string, err error) {
	repo, ok := s.inv.repo.(EquippedRepo)
	if !ok {
		return "", fmt.Errorf("inventory: repository does not implement EquippedRepo")
	}
	inst, err := s.items.Get(ctx, itemID)
	if err != nil {
		return "", ErrNotFound
	}
	if inst.OwnerCharacterID == nil || *inst.OwnerCharacterID != characterID {
		return "", ErrNotFound
	}
	if inst.IsLocked {
		return "", ErrLocked
	}
	tmpl, err := s.items.Catalog().Get(inst.TemplateCode)
	if err != nil {
		return "", err
	}
	if !tmpl.IsEquippable() {
		return "", ErrNotEquippable
	}
	slot := *tmpl.EquipSlot
	if !KnownEquipSlots[slot] {
		return "", ErrUnknownEquipSlot
	}

	// Class + level check.
	char, err := s.charView.Get(ctx, characterID)
	if err != nil {
		return "", err
	}
	if tmpl.ClassRestriction != nil && *tmpl.ClassRestriction != "" && *tmpl.ClassRestriction != char.ClassCode {
		return "", ErrClassRestriction
	}
	if int16(tmpl.MinLevel) > char.Level {
		return "", ErrLevelRequirement
	}

	// If existing item in that slot, unequip first.
	if existingID, err := repo.FindEquippedInSlot(ctx, characterID, slot); err == nil && existingID != uuid.Nil {
		if err := s.unequipToFreeBag(ctx, repo, characterID, existingID); err != nil {
			return "", fmt.Errorf("unequip existing: %w", err)
		}
	}

	// Now equip the new item: bag slot → NULL, equipment_slot → slot.
	if err := repo.EquipItem(ctx, itemID, characterID, slot); err != nil {
		return "", err
	}
	return slot, nil
}

// Unequip moves the item in `equipSlot` back to a free bag slot.
func (s *EquipmentService) Unequip(ctx context.Context, characterID uuid.UUID, equipSlot string) (uuid.UUID, error) {
	if !KnownEquipSlots[equipSlot] {
		return uuid.Nil, ErrUnknownEquipSlot
	}
	repo, ok := s.inv.repo.(EquippedRepo)
	if !ok {
		return uuid.Nil, fmt.Errorf("inventory: repository does not implement EquippedRepo")
	}
	itemID, err := repo.FindEquippedInSlot(ctx, characterID, equipSlot)
	if err != nil {
		return uuid.Nil, ErrNothingEquipped
	}
	if err := s.unequipToFreeBag(ctx, repo, characterID, itemID); err != nil {
		return uuid.Nil, err
	}
	return itemID, nil
}

// EquippedBonuses satisfies skill.EquippedFetcher — returns bonus payloads for
// every equipped item for the sheet computer.
func (s *EquipmentService) EquippedBonuses(ctx context.Context, characterID uuid.UUID) ([]skill.EquippedBonus, error) {
	repo, ok := s.inv.repo.(EquippedRepo)
	if !ok {
		return nil, fmt.Errorf("inventory: repository does not implement EquippedRepo")
	}
	list, err := repo.Equipped(ctx, characterID)
	if err != nil {
		return nil, err
	}
	out := make([]skill.EquippedBonus, 0, len(list))
	for _, b := range list {
		if len(b.StatBonusesRaw) > 0 {
			out = append(out, skill.EquippedBonus{StatBonusesRaw: b.StatBonusesRaw})
		}
	}
	return out, nil
}

func (s *EquipmentService) unequipToFreeBag(ctx context.Context, repo EquippedRepo, characterID, itemID uuid.UUID) error {
	occupied, err := s.inv.repo.OccupiedSlots(ctx, characterID)
	if err != nil {
		return err
	}
	taken := map[Slot]bool{}
	for _, s := range occupied {
		taken[s] = true
	}
	free := Slot(-1)
	for i := Slot(0); i < BagSlotCap; i++ {
		if !taken[i] {
			free = i
			break
		}
	}
	if free < 0 {
		return ErrBagFull
	}
	return repo.UnequipItem(ctx, itemID, characterID, free)
}
