package inventory

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/otter/server/internal/item"
)

// Service is the inventory business logic. Slot allocation + stack merge live
// here; the repository stays a thin persistence layer.
type Service struct {
	repo Repository
	log  *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{repo: repo, log: logger}
}

// --- View -------------------------------------------------------------------

// View returns the compound bag snapshot for GET /inventory. Overflow items
// (unslotted: bag was full on arrival, or returned from mail/market escrow)
// are slotted first, so they show up as soon as there is room.
func (s *Service) View(ctx context.Context, characterID uuid.UUID) (*View, error) {
	if overflow, err := s.repo.Unslotted(ctx, characterID); err != nil {
		return nil, err
	} else if len(overflow) > 0 {
		if err := s.AssignSlots(ctx, characterID, overflow); err != nil {
			return nil, err
		}
	}
	items, err := s.repo.Bag(ctx, characterID)
	if err != nil {
		return nil, err
	}
	var weight int64
	for _, b := range items {
		weight += b.TotalWeight()
	}
	if items == nil {
		items = []BagItem{}
	}
	return &View{
		CharacterID:    characterID,
		Items:          items,
		SlotCap:        BagSlotCap,
		SlotsUsed:      len(items),
		WeightGrams:    weight,
		WeightCapGrams: WeightCapGrams,
	}, nil
}

// --- AssignSlots (item.SlotAssigner impl) -----------------------------------

// AssignSlots satisfies item.SlotAssigner. Called by item.Service after batch
// insert of new instances. For each instance:
//  1. If stackable + existing stack has room → merge (add to stack, delete
//     the newly-created instance).
//  2. Else → find first free slot and assign.
//  3. If bag full → leave slot=NULL (overflow — Phase 02 mail fallback).
//
// Returns nil on best-effort success; per-item failures logged but not fatal.
func (s *Service) AssignSlots(ctx context.Context, ownerID uuid.UUID, instances []item.Instance) error {
	occupied, err := s.repo.OccupiedSlots(ctx, ownerID)
	if err != nil {
		return err
	}
	takenSet := make(map[Slot]bool, len(occupied))
	for _, s := range occupied {
		takenSet[s] = true
	}
	nextSlot := func() (Slot, bool) {
		for i := Slot(0); i < BagSlotCap; i++ {
			if !takenSet[i] {
				takenSet[i] = true
				return i, true
			}
		}
		return -1, false
	}

	for _, inst := range instances {
		if err := s.assignOne(ctx, ownerID, inst, nextSlot); err != nil {
			s.log.WarnContext(ctx, "inventory: assign slot failed",
				slog.String("item_id", inst.ID.String()),
				slog.String("template", inst.TemplateCode),
				slog.String("err", err.Error()))
			// Best-effort: continue. Item stays slot=NULL (overflow).
		}
	}
	return nil
}

func (s *Service) assignOne(ctx context.Context, ownerID uuid.UUID, inst item.Instance, nextSlot func() (Slot, bool)) error {
	// Try merge for stackable items.
	candidates, err := s.repo.StackableMergeCandidates(ctx, ownerID, inst.TemplateCode)
	if err == nil && len(candidates) > 0 {
		remaining := inst.StackQuantity
		for _, c := range candidates {
			if c.InstanceID == inst.ID {
				continue // already slotted (e.g. by a concurrent View): never merge into itself
			}
			space := c.MaxStack - c.StackQuantity
			if space <= 0 {
				continue
			}
			take := remaining
			if take > space {
				take = space
			}
			if err := s.repo.AddStack(ctx, c.InstanceID, take); err == nil {
				remaining -= take
			}
			if remaining <= 0 {
				break
			}
		}
		if remaining <= 0 {
			// Fully merged — the newly-created instance is redundant. Delete it.
			return s.repo.DeleteInstance(ctx, inst.ID, ownerID)
		}
		// Partial merge: shrink the new instance to the leftover, then slot it.
		delta := remaining - inst.StackQuantity // negative
		if delta < 0 {
			// Best-effort adjust — memory repo respects; pg repo uses AddStack which
			// accepts any int32 (positive or negative).
			_ = s.repo.AddStack(ctx, inst.ID, delta)
		}
	}
	// Find free slot.
	slot, ok := nextSlot()
	if !ok {
		return ErrBagFull
	}
	return s.repo.AssignSlot(ctx, inst.ID, ownerID, slot)
}

// --- Move -------------------------------------------------------------------

// Move relocates an item to `targetSlot`. Behavior:
//   - target empty → simple move
//   - target occupied by another unlocked item → atomic swap
//   - target occupied by a locked item → ErrSlotOccupied
//   - self-target (source == target) → no-op
func (s *Service) Move(ctx context.Context, characterID, itemID uuid.UUID, targetSlot Slot) error {
	if targetSlot < 0 || targetSlot >= BagSlotCap {
		return ErrSlotOutOfRange
	}
	bag, err := s.repo.Bag(ctx, characterID)
	if err != nil {
		return err
	}

	var source *BagItem
	var target *BagItem
	for i := range bag {
		if bag[i].InstanceID == itemID {
			source = &bag[i]
		}
		if bag[i].Slot == targetSlot {
			target = &bag[i]
		}
	}
	if source == nil {
		return ErrNotFound
	}
	if source.IsLocked {
		return ErrLocked
	}
	if source.Slot == targetSlot {
		return nil
	}
	if target == nil {
		return s.repo.MoveToEmpty(ctx, itemID, characterID, targetSlot)
	}
	if target.IsLocked {
		return ErrSlotOccupied
	}
	return s.repo.SwapSlots(ctx, characterID, itemID, target.InstanceID)
}

// --- Split ------------------------------------------------------------------

// Split takes `amount` from a stackable item and places it in `targetSlot`.
// Target must be empty. Amount must be in (0, source_stack).
func (s *Service) Split(ctx context.Context, characterID, itemID uuid.UUID, amount int32, targetSlot Slot) (uuid.UUID, error) {
	if targetSlot < 0 || targetSlot >= BagSlotCap {
		return uuid.Nil, ErrSlotOutOfRange
	}
	return s.repo.SplitStack(ctx, itemID, characterID, amount, targetSlot)
}

// --- Delete -----------------------------------------------------------------

// Delete removes an owned + unlocked item. Idempotent: repeated calls after
// the first success return ErrNotFound.
func (s *Service) Delete(ctx context.Context, characterID, itemID uuid.UUID) error {
	return s.repo.DeleteInstance(ctx, itemID, characterID)
}

// ConsumeAll removes all needs atomically (crafting inputs).
func (s *Service) ConsumeAll(ctx context.Context, characterID uuid.UUID, needs []Need) error {
	return s.repo.ConsumeAll(ctx, characterID, needs)
}

// Available counts consumable units of a template.
func (s *Service) Available(ctx context.Context, characterID uuid.UUID, templateCode string) (int32, error) {
	return s.repo.Available(ctx, characterID, templateCode)
}
