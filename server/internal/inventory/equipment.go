package inventory

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/itemset"
	"github.com/otter/server/internal/leveling"
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

// PairedSlots: a ring or trinket template fills whichever of its two
// equipment slots is free (DD-276).
var PairedSlots = map[string][2]string{
	"ring":    {"ring1", "ring2"},
	"trinket": {"trinket1", "trinket2"},
}

// ValidTemplateSlot reports whether content may give a template this
// equip_slot: an equipment slot, or "ring"/"trinket" for the paired ones.
func ValidTemplateSlot(slot string) bool {
	_, paired := PairedSlots[slot]
	return paired || (KnownEquipSlots[slot] && slot != "ring1" && slot != "ring2" &&
		slot != "trinket1" && slot != "trinket2")
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
	Path      []string       // class path (base → specialization), DD-250
	Raw       leveling.Stats // raw attributes (requirements ignore diminishing returns)
}

// Profile is what equip checks read about a character (DD-250 §34).
type Profile struct {
	Level int16
	Class string
	Path  []string
	Raw   leveling.Stats
}

// ProfileFetcher returns a character's equip profile.
type ProfileFetcher func(ctx context.Context, id uuid.UUID) (Profile, error)

func (f ProfileFetcher) Get(ctx context.Context, id uuid.UUID) (*characterView, error) {
	p, err := f(ctx, id)
	if err != nil {
		return nil, err
	}
	return &characterView{ID: id, Level: p.Level, ClassCode: p.Class, Path: p.Path, Raw: p.Raw}, nil
}

// meetsRequirements checks a template against a character: class, level,
// class branch and raw attribute minimums.
func meetsRequirements(tmpl item.Template, c *characterView) error {
	// What the class can never wear comes first, then what it cannot wear yet.
	if tmpl.ClassRestriction != nil && *tmpl.ClassRestriction != "" && *tmpl.ClassRestriction != c.ClassCode {
		return ErrClassRestriction
	}
	if len(c.Path) > 0 && !leveling.Current().CanUse(c.Path[len(c.Path)-1], tmpl.ItemSubclass) {
		return ErrEquipmentPermission
	}
	rq := tmpl.Requirements
	if rq != nil && rq.Spec != "" && !slices.Contains(c.Path, rq.Spec) {
		return ErrSpecRequirement
	}
	if int16(tmpl.MinLevel) > c.Level {
		return ErrLevelRequirement
	}
	if rq != nil && (c.Raw.Strength < rq.Strength || c.Raw.Vitality < rq.Vitality || c.Raw.Dexterity < rq.Dexterity) {
		return ErrAttributeRequirement
	}
	return nil
}

// Wear verdicts for a bag item (DD-290).
const (
	WearOK        = "ok"
	WearClass     = "class"     // the class (or its branch) never wears it
	WearLevel     = "level"     // not yet: the level is too low
	WearAttribute = "attribute" // not yet: attributes too low
)

// MarkWearable sets Wear on every equippable item: whether the character
// can wear it now and, if not, why (the same checks as Equip).
func (s *EquipmentService) MarkWearable(ctx context.Context, characterID uuid.UUID, items []BagItem) error {
	var char *characterView
	for i := range items {
		it := &items[i]
		if it.EquipSlotOnTmpl == nil || *it.EquipSlotOnTmpl == "" {
			continue
		}
		tmpl, err := s.items.Catalog().Get(it.TemplateCode)
		if err != nil {
			continue
		}
		if char == nil {
			if char, err = s.charView.Get(ctx, characterID); err != nil {
				return err
			}
		}
		switch err := meetsRequirements(tmpl, char); {
		case err == nil:
			it.Wear = WearOK
		case errors.Is(err, ErrLevelRequirement):
			it.Wear = WearLevel
		case errors.Is(err, ErrAttributeRequirement):
			it.Wear = WearAttribute
		default:
			it.Wear = WearClass
		}
	}
	return nil
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
	if pair, ok := PairedSlots[slot]; ok {
		// The first free one of the pair; both taken: the first is swapped.
		slot = pair[0]
		if id, err := repo.FindEquippedInSlot(ctx, characterID, pair[0]); err == nil && id != uuid.Nil {
			if id2, err := repo.FindEquippedInSlot(ctx, characterID, pair[1]); err != nil || id2 == uuid.Nil {
				slot = pair[1]
			}
		}
	}
	if !KnownEquipSlots[slot] {
		return "", ErrUnknownEquipSlot
	}

	// Class, level, branch and attribute requirements.
	char, err := s.charView.Get(ctx, characterID)
	if err != nil {
		return "", err
	}
	if err := meetsRequirements(tmpl, char); err != nil {
		return "", err
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
	out := make([]skill.EquippedBonus, 0, len(list)+1)
	worn := make([]string, 0, len(list))
	for _, b := range list {
		worn = append(worn, b.TemplateCode)
		if mods := EffectiveModifiers(b); len(mods) > 0 {
			out = append(out, skill.EquippedBonus{Modifiers: mods})
		}
	}
	// Set bonuses for the worn pieces (DD-277).
	var set []skill.Modifier
	for _, m := range itemset.Current().Modifiers(worn) {
		set = append(set, skill.Modifier{Kind: skill.ModifierKind(m.Kind), Target: skill.Target(m.Target), Amount: m.Amount})
	}
	if len(set) > 0 {
		out = append(out, skill.EquippedBonus{Modifiers: set})
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

// UnequipInvalid moves every equipped item the character may no longer wear
// (e.g. after a specialization change) back to the bag. Returns how many
// moved; a full bag stops it with ErrBagFull.
func (s *EquipmentService) UnequipInvalid(ctx context.Context, characterID uuid.UUID) (int, error) {
	repo, ok := s.inv.repo.(EquippedRepo)
	if !ok {
		return 0, fmt.Errorf("inventory: repository does not implement EquippedRepo")
	}
	char, err := s.charView.Get(ctx, characterID)
	if err != nil {
		return 0, err
	}
	worn, err := repo.Equipped(ctx, characterID)
	if err != nil {
		return 0, err
	}
	moved := 0
	for _, it := range worn {
		tmpl, err := s.items.Catalog().Get(it.TemplateCode)
		if err != nil || meetsRequirements(tmpl, char) == nil {
			continue
		}
		if err := s.unequipToFreeBag(ctx, repo, characterID, it.InstanceID); err != nil {
			return moved, err
		}
		moved++
	}
	return moved, nil
}
