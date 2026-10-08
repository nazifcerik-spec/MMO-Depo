package inventory

import (
	"context"

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

	// ConsumeAll removes every need from the character's unlocked,
	// un-equipped bag stacks in ONE transaction (fewest-first per template to
	// limit fragmentation). Returns ErrInsufficient when any template falls
	// short; nothing is consumed on that path.
	ConsumeAll(ctx context.Context, ownerID uuid.UUID, needs []Need) error

	// Available counts the consumable units of a template (same filter as
	// ConsumeAll). Read-only pre-check.
	Available(ctx context.Context, ownerID uuid.UUID, templateCode string) (int32, error)
}

// Need is a quantity of one template to consume.
type Need struct {
	TemplateCode string
	Quantity     int32
}
