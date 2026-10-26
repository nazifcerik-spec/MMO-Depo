package inventory

import (
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/otter/server/internal/item"
)

// --- MemoryRepo ------------------------------------------------------

// Simplified memory model: keeps the essential fields for slot/weight logic.
// Not intended to mirror every column on item_instance.
type memItem struct {
	BagItem
	OwnerID uuid.UUID
	// Template snapshot fields (would come from Catalog join in Postgres).
	Stackable bool
}

type MemoryRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*memItem
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{items: make(map[uuid.UUID]*memItem)}
}

// Seed adds an item to the memory store — used by tests. `slot` may be -1 for
// unslotted.
func (r *MemoryRepo) Seed(item BagItem, ownerID uuid.UUID, stackable bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := item
	r.items[cp.InstanceID] = &memItem{BagItem: cp, OwnerID: ownerID, Stackable: stackable}
}

func (r *MemoryRepo) Bag(_ context.Context, characterID uuid.UUID) ([]BagItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []BagItem{}
	for _, it := range r.items {
		if it.OwnerID == characterID && it.Slot >= 0 {
			out = append(out, it.BagItem)
		}
	}
	// sort by slot
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Slot < out[i].Slot {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

func (r *MemoryRepo) Unslotted(_ context.Context, characterID uuid.UUID) ([]item.Instance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []item.Instance{}
	for _, it := range r.items {
		if it.OwnerID == characterID && it.Slot < 0 && !it.IsLocked && it.EquipmentSlot == nil {
			out = append(out, item.Instance{ID: it.InstanceID, TemplateCode: it.TemplateCode, StackQuantity: it.StackQuantity})
		}
	}
	return out, nil
}

func (r *MemoryRepo) Compact(_ context.Context, characterID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var bag []*memItem
	for _, it := range r.items {
		if it.OwnerID == characterID && it.Slot >= 0 {
			bag = append(bag, it)
		}
	}
	sort.Slice(bag, func(i, j int) bool { return bag[i].Slot < bag[j].Slot })
	for n, it := range bag {
		it.Slot = Slot(n)
	}
	return nil
}

func (r *MemoryRepo) OccupiedSlots(_ context.Context, characterID uuid.UUID) ([]Slot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []Slot{}
	for _, it := range r.items {
		if it.OwnerID == characterID && it.Slot >= 0 {
			out = append(out, it.Slot)
		}
	}
	return out, nil
}

func (r *MemoryRepo) StackableMergeCandidates(_ context.Context, characterID uuid.UUID, templateCode string) ([]BagItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []BagItem{}
	for _, it := range r.items {
		if it.OwnerID != characterID || !it.Stackable || it.Slot < 0 {
			continue
		}
		if it.IsLocked {
			continue
		}
		if it.TemplateCode != templateCode {
			continue
		}
		if it.StackQuantity >= it.MaxStack {
			continue
		}
		out = append(out, it.BagItem)
	}
	// order by slot
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Slot < out[i].Slot {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

func (r *MemoryRepo) AssignSlot(_ context.Context, itemID, ownerID uuid.UUID, slot Slot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok || it.OwnerID != ownerID {
		return ErrNotFound
	}
	if it.Slot >= 0 {
		return ErrSlotOccupied
	}
	// Check target slot not taken.
	for _, other := range r.items {
		if other.OwnerID == ownerID && other.Slot == slot {
			return ErrSlotOccupied
		}
	}
	it.Slot = slot
	return nil
}

func (r *MemoryRepo) AddStack(_ context.Context, itemID uuid.UUID, delta int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok {
		return ErrNotFound
	}
	if it.IsLocked {
		return ErrLocked
	}
	it.StackQuantity += delta
	return nil
}

func (r *MemoryRepo) SwapSlots(_ context.Context, ownerID, itemA, itemB uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.items[itemA]
	if !ok || a.OwnerID != ownerID {
		return ErrNotFound
	}
	b, ok := r.items[itemB]
	if !ok || b.OwnerID != ownerID {
		return ErrNotFound
	}
	if a.IsLocked || b.IsLocked {
		return ErrLocked
	}
	a.Slot, b.Slot = b.Slot, a.Slot
	return nil
}

func (r *MemoryRepo) MoveToEmpty(_ context.Context, itemID, ownerID uuid.UUID, target Slot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok || it.OwnerID != ownerID {
		return ErrNotFound
	}
	if it.IsLocked {
		return ErrLocked
	}
	for _, other := range r.items {
		if other.OwnerID == ownerID && other.InstanceID != itemID && other.Slot == target {
			return ErrSlotOccupied
		}
	}
	it.Slot = target
	return nil
}

func (r *MemoryRepo) SplitStack(_ context.Context, sourceID, ownerID uuid.UUID, amount int32, target Slot) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	src, ok := r.items[sourceID]
	if !ok || src.OwnerID != ownerID {
		return uuid.Nil, ErrNotFound
	}
	if src.IsLocked {
		return uuid.Nil, ErrLocked
	}
	if !src.Stackable {
		return uuid.Nil, ErrNotStackable
	}
	if amount <= 0 || amount >= src.StackQuantity {
		return uuid.Nil, ErrInvalidSplit
	}
	for _, other := range r.items {
		if other.OwnerID == ownerID && other.InstanceID != sourceID && other.Slot == target {
			return uuid.Nil, ErrSlotOccupied
		}
	}
	src.StackQuantity -= amount
	newID := uuid.New()
	cp := src.BagItem
	cp.InstanceID = newID
	cp.Slot = target
	cp.StackQuantity = amount
	r.items[newID] = &memItem{BagItem: cp, OwnerID: ownerID, Stackable: src.Stackable}
	return newID, nil
}

func (r *MemoryRepo) DeleteInstance(_ context.Context, itemID, ownerID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok || it.OwnerID != ownerID {
		return ErrNotFound
	}
	if it.IsLocked {
		return ErrLocked
	}
	delete(r.items, itemID)
	return nil
}

func (r *MemoryRepo) consumable(ownerID uuid.UUID, templateCode string) ([]*memItem, int32) {
	var cands []*memItem
	var total int32
	for _, it := range r.items {
		if it.OwnerID != ownerID || it.TemplateCode != templateCode || it.IsLocked ||
			it.BagItem.EquipmentSlot != nil || it.Slot < 0 {
			continue
		}
		cands = append(cands, it)
		total += it.StackQuantity
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].StackQuantity < cands[j].StackQuantity })
	return cands, total
}

func (r *MemoryRepo) Available(_ context.Context, ownerID uuid.UUID, templateCode string) (int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, total := r.consumable(ownerID, templateCode)
	return total, nil
}

func (r *MemoryRepo) ConsumeAll(_ context.Context, ownerID uuid.UUID, needs []Need) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range needs {
		if _, total := r.consumable(ownerID, n.TemplateCode); total < n.Quantity {
			return ErrInsufficient
		}
	}
	for _, n := range needs {
		cands, _ := r.consumable(ownerID, n.TemplateCode)
		remaining := n.Quantity
		for _, it := range cands {
			if remaining <= 0 {
				break
			}
			take := min(it.StackQuantity, remaining)
			it.StackQuantity -= take
			remaining -= take
			if it.StackQuantity == 0 {
				delete(r.items, it.InstanceID)
			}
		}
	}
	return nil
}

// --- Equipment methods (satisfies EquippedRepo) -----------------------------

func (r *MemoryRepo) Equipped(_ context.Context, characterID uuid.UUID) ([]BagItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []BagItem{}
	for _, it := range r.items {
		if it.OwnerID == characterID && it.BagItem.EquipmentSlot != nil {
			out = append(out, it.BagItem)
		}
	}
	return out, nil
}

func (r *MemoryRepo) EquipItem(_ context.Context, itemID, ownerID uuid.UUID, equipSlot string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok || it.OwnerID != ownerID {
		return ErrNotFound
	}
	if it.IsLocked {
		return ErrLocked
	}
	// Check no other item in that slot for this owner.
	for _, other := range r.items {
		if other.OwnerID == ownerID && other.BagItem.EquipmentSlot != nil && *other.BagItem.EquipmentSlot == equipSlot {
			return ErrSlotOccupied
		}
	}
	slot := equipSlot
	it.BagItem.EquipmentSlot = &slot
	it.BagItem.Slot = -1
	return nil
}

func (r *MemoryRepo) UnequipItem(_ context.Context, itemID, ownerID uuid.UUID, targetBagSlot Slot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok || it.OwnerID != ownerID {
		return ErrNotFound
	}
	if it.IsLocked {
		return ErrLocked
	}
	// Check target bag slot free.
	for _, other := range r.items {
		if other.OwnerID == ownerID && other.InstanceID != itemID && other.BagItem.Slot == targetBagSlot {
			return ErrSlotOccupied
		}
	}
	it.BagItem.EquipmentSlot = nil
	it.BagItem.Slot = targetBagSlot
	return nil
}

func (r *MemoryRepo) FindEquippedInSlot(_ context.Context, ownerID uuid.UUID, equipSlot string) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, it := range r.items {
		if it.OwnerID == ownerID && it.BagItem.EquipmentSlot != nil && *it.BagItem.EquipmentSlot == equipSlot {
			return it.InstanceID, nil
		}
	}
	return uuid.Nil, ErrNotFound
}

// Silence unused import warning in some build configurations.
