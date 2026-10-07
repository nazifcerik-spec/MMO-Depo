package inventory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/otter/server/internal/item"
)

// Repository is the persistence contract. Every mutation is a single SQL
// statement that respects the (owner, slot) unique constraint — no
// application-level lock needed for correctness.
type Repository interface {
	// Bag returns the character's slot-assigned items enriched with template
	// metadata for view/computation.
	Bag(ctx context.Context, characterID uuid.UUID) ([]BagItem, error)

	// OccupiedSlots returns the set of slots currently in use for O(1)
	// free-slot lookup (used by slot allocator).
	OccupiedSlots(ctx context.Context, characterID uuid.UUID) ([]Slot, error)

	// StackableMergeCandidates returns items that could accept more of the
	// given template — same template + not locked + not at max_stack + in bag.
	// Ordered by slot ascending so merge is deterministic.
	StackableMergeCandidates(ctx context.Context, characterID uuid.UUID, templateCode string) ([]BagItem, error)

	// Unslotted returns the character's overflow: owned items with no bag
	// slot that are neither locked (mail/market escrow) nor equipped.
	Unslotted(ctx context.Context, characterID uuid.UUID) ([]item.Instance, error)

	// AssignSlot atomically sets slot on the item, checking:
	//   - owner_character_id matches
	//   - slot is currently null
	// Returns ErrSlotOccupied if the target slot is already taken.
	AssignSlot(ctx context.Context, itemID, ownerID uuid.UUID, slot Slot) error

	// AddStack atomically increments stack_quantity by delta (must be > 0).
	// Rejected when item is locked.
	AddStack(ctx context.Context, itemID uuid.UUID, delta int32) error

	// SwapSlots atomically swaps slot values of two items. Both must be owned
	// by ownerID and unlocked. Uses a temporary slot value to bypass unique
	// constraint during the swap.
	SwapSlots(ctx context.Context, ownerID, itemA, itemB uuid.UUID) error

	// MoveToEmpty atomically moves an item to an empty slot (target had no
	// item). Fails if target is already occupied.
	MoveToEmpty(ctx context.Context, itemID, ownerID uuid.UUID, target Slot) error

	// SplitStack decrements the source stack by amount and inserts a new
	// instance (copy of the source's template + rarity + affix + owner) with
	// stack=amount in `target` slot. Returns the new instance id.
	SplitStack(ctx context.Context, sourceID, ownerID uuid.UUID, amount int32, target Slot) (uuid.UUID, error)

	// DeleteInstance removes an owned + unlocked item.
	DeleteInstance(ctx context.Context, itemID, ownerID uuid.UUID) error

	// ConsumeByTemplate reduces `quantity` units of `templateCode` from the
	// character's unlocked, un-equipped bag stacks in a single transaction.
	// Order-independent (fewest-first minimises fragmentation). Returns
	// ErrInsufficient when the total available is < quantity; the transaction
	// is rolled back on that path so no partial consumption occurs.
	ConsumeByTemplate(ctx context.Context, ownerID uuid.UUID, templateCode string, quantity int32) error
}

// --- MemoryRepo (tests) ------------------------------------------------------

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

func (r *MemoryRepo) ConsumeByTemplate(_ context.Context, ownerID uuid.UUID, templateCode string, quantity int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if quantity <= 0 {
		return nil
	}
	// Collect candidates deterministically (fewest first so we drain small
	// stacks before touching bigger ones — leaves the bag less fragmented).
	var cands []*memItem
	var total int32
	for _, it := range r.items {
		if it.OwnerID != ownerID || it.TemplateCode != templateCode {
			continue
		}
		if it.IsLocked || it.BagItem.EquipmentSlot != nil {
			continue
		}
		if it.Slot < 0 {
			continue
		}
		cands = append(cands, it)
		total += it.StackQuantity
	}
	if total < quantity {
		return ErrInsufficient
	}
	// Simple selection sort by StackQuantity asc — count is tiny (bag ≤60).
	for i := 0; i < len(cands); i++ {
		for j := i + 1; j < len(cands); j++ {
			if cands[j].StackQuantity < cands[i].StackQuantity {
				cands[i], cands[j] = cands[j], cands[i]
			}
		}
	}
	remaining := quantity
	for _, it := range cands {
		if remaining <= 0 {
			break
		}
		take := it.StackQuantity
		if take > remaining {
			take = remaining
		}
		it.StackQuantity -= take
		remaining -= take
		if it.StackQuantity == 0 {
			delete(r.items, it.InstanceID)
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
var _ = errors.Is
