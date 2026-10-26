package inventory_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
)

func newSvc() (*inventory.Service, *inventory.MemoryRepo) {
	repo := inventory.NewMemoryRepo()
	svc := inventory.NewService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo
}

func seedItem(repo *inventory.MemoryRepo, owner uuid.UUID, tmpl string, stack, max int32, slot inventory.Slot, stackable bool) uuid.UUID {
	id := uuid.New()
	repo.Seed(inventory.BagItem{
		InstanceID:     id,
		TemplateCode:   tmpl,
		Slot:           slot,
		StackQuantity:  stack,
		MaxStack:       max,
		WeightGramsPer: 100,
	}, owner, stackable)
	return id
}

func TestView_EmptyBag(t *testing.T) {
	svc, _ := newSvc()
	view, err := svc.View(context.Background(), uuid.New())
	require.NoError(t, err)
	require.Empty(t, view.Items)
	require.Equal(t, 0, view.SlotsUsed)
	require.Equal(t, int64(0), view.WeightGrams)
	require.Equal(t, inventory.BagSlotCap, view.SlotCap)
}

func TestView_ComputesWeight(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	seedItem(repo, owner, "iron_ore", 10, 99, 0, true) // 100g × 10 = 1000
	seedItem(repo, owner, "iron_ore", 20, 99, 1, true) // 100g × 20 = 2000

	view, err := svc.View(context.Background(), owner)
	require.NoError(t, err)
	require.Equal(t, 2, view.SlotsUsed)
	require.Equal(t, int64(3000), view.WeightGrams)
}

func TestAssignSlots_AutoSlotStackable(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()

	// Two new stackable instances to place.
	instances := []item.Instance{
		mkInstance(owner, "iron_ore", 5),
		mkInstance(owner, "iron_ore", 3),
	}
	// Seed the created instances into memory repo (item.MaterializeDrops would
	// have done this before AssignSlots).
	for _, inst := range instances {
		repo.Seed(inventory.BagItem{
			InstanceID:     inst.ID,
			TemplateCode:   inst.TemplateCode,
			Slot:           -1, // unslotted
			StackQuantity:  inst.StackQuantity,
			MaxStack:       99,
			WeightGramsPer: 100,
		}, owner, true)
	}

	require.NoError(t, svc.AssignSlots(context.Background(), owner, instances))

	view, _ := svc.View(context.Background(), owner)
	// With merge, both stacks might merge into one; with fresh bag they get slots 0 and 1.
	require.GreaterOrEqual(t, view.SlotsUsed, 1)
}

func TestMove_ToEmptySlot(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 5, 99, 3, true)

	require.NoError(t, svc.Move(context.Background(), owner, id, 10))

	bag, _ := repo.Bag(context.Background(), owner) // as stored: the view would close the gap
	require.Len(t, bag, 1)
	require.Equal(t, inventory.Slot(10), bag[0].Slot)
}

func TestMove_Swap(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	a := seedItem(repo, owner, "iron_ore", 5, 99, 0, true)
	b := seedItem(repo, owner, "copper_ore", 3, 99, 1, true)

	require.NoError(t, svc.Move(context.Background(), owner, a, 1))

	view, _ := svc.View(context.Background(), owner)
	slotByID := map[uuid.UUID]inventory.Slot{}
	for _, it := range view.Items {
		slotByID[it.InstanceID] = it.Slot
	}
	require.Equal(t, inventory.Slot(1), slotByID[a])
	require.Equal(t, inventory.Slot(0), slotByID[b])
}

func TestMove_SelfIsNoop(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 5, 99, 3, true)

	require.NoError(t, svc.Move(context.Background(), owner, id, 3))

	bag, _ := repo.Bag(context.Background(), owner)
	require.Equal(t, inventory.Slot(3), bag[0].Slot)
}

// The bag view closes the gaps sales and destroys leave, keeping the order
// (DD-291).
func TestView_ClosesGaps(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	c := seedItem(repo, owner, "copper_ore", 1, 99, 7, true)
	a := seedItem(repo, owner, "iron_ore", 1, 99, 2, true)
	b := seedItem(repo, owner, "tin_ore", 1, 99, 4, true)

	view, err := svc.View(context.Background(), owner)
	require.NoError(t, err)
	slotByID := map[uuid.UUID]inventory.Slot{}
	for _, it := range view.Items {
		slotByID[it.InstanceID] = it.Slot
	}
	require.Equal(t, map[uuid.UUID]inventory.Slot{a: 0, b: 1, c: 2}, slotByID)
}

func TestMove_OutOfRange(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 5, 99, 0, true)

	require.ErrorIs(t, svc.Move(context.Background(), owner, id, inventory.BagSlotCap), inventory.ErrSlotOutOfRange)
	require.ErrorIs(t, svc.Move(context.Background(), owner, id, -1), inventory.ErrSlotOutOfRange)
}

func TestMove_Locked(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := uuid.New()
	repo.Seed(inventory.BagItem{
		InstanceID: id, TemplateCode: "iron_ore", Slot: 0,
		StackQuantity: 5, MaxStack: 99, IsLocked: true,
	}, owner, true)

	err := svc.Move(context.Background(), owner, id, 5)
	require.ErrorIs(t, err, inventory.ErrLocked)
}

func TestSplit_Success(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 10, 99, 0, true)

	newID, err := svc.Split(context.Background(), owner, id, 4, 5)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, newID)

	bag, _ := repo.Bag(context.Background(), owner)
	require.Len(t, bag, 2)
	// Original decremented to 6; new stack has 4 in slot 5.
	var orig, split *inventory.BagItem
	for i := range bag {
		if bag[i].InstanceID == id {
			orig = &bag[i]
		} else if bag[i].InstanceID == newID {
			split = &bag[i]
		}
	}
	require.NotNil(t, orig)
	require.NotNil(t, split)
	require.Equal(t, int32(6), orig.StackQuantity)
	require.Equal(t, int32(4), split.StackQuantity)
	require.Equal(t, inventory.Slot(5), split.Slot)
}

func TestSplit_InvalidAmount(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 10, 99, 0, true)

	_, err := svc.Split(context.Background(), owner, id, 0, 1)
	require.ErrorIs(t, err, inventory.ErrInvalidSplit)

	_, err = svc.Split(context.Background(), owner, id, 10, 1) // == full stack
	require.ErrorIs(t, err, inventory.ErrInvalidSplit)

	_, err = svc.Split(context.Background(), owner, id, 15, 1) // > stack
	require.ErrorIs(t, err, inventory.ErrInvalidSplit)
}

func TestSplit_NotStackable(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "sword_iron", 1, 1, 0, false)

	_, err := svc.Split(context.Background(), owner, id, 1, 5)
	require.ErrorIs(t, err, inventory.ErrNotStackable)
}

func TestDelete_Success(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 5, 99, 0, true)

	require.NoError(t, svc.Delete(context.Background(), owner, id))

	view, _ := svc.View(context.Background(), owner)
	require.Empty(t, view.Items)
}

func TestDelete_Locked(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := uuid.New()
	repo.Seed(inventory.BagItem{
		InstanceID: id, TemplateCode: "iron_ore", Slot: 0,
		StackQuantity: 5, MaxStack: 99, IsLocked: true,
	}, owner, true)

	require.ErrorIs(t, svc.Delete(context.Background(), owner, id), inventory.ErrLocked)
}

func TestDelete_ForeignOwnedReturnsNotFound(t *testing.T) {
	svc, repo := newSvc()
	id := seedItem(repo, uuid.New(), "iron_ore", 5, 99, 0, true)

	err := svc.Delete(context.Background(), uuid.New(), id)
	require.ErrorIs(t, err, inventory.ErrNotFound)
}

func TestAssignSlots_BagFullOverflow(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()

	// Fill every slot with non-stackable items so no merge is possible.
	for s := inventory.Slot(0); s < inventory.BagSlotCap; s++ {
		seedItem(repo, owner, "sword_iron", 1, 1, s, false)
	}

	// Try to add one more; should log warning and leave unslotted.
	inst := mkInstance(owner, "sword_iron", 1)
	repo.Seed(inventory.BagItem{
		InstanceID:     inst.ID,
		TemplateCode:   inst.TemplateCode,
		Slot:           -1,
		StackQuantity:  1,
		MaxStack:       1,
		WeightGramsPer: 100,
	}, owner, false)

	require.NoError(t, svc.AssignSlots(context.Background(), owner, []item.Instance{inst}))

	view, _ := svc.View(context.Background(), owner)
	require.Equal(t, inventory.BagSlotCap, view.SlotsUsed, "overflow item stays unslotted")
}

// Overflow (unslotted, unlocked, unequipped) shows up in the bag as soon as a
// slot is free; escrowed and equipped items are left alone.
func TestView_SlotsOverflowWhenRoom(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	for s := inventory.Slot(0); s < inventory.BagSlotCap; s++ {
		seedItem(repo, owner, "sword_iron", 1, 1, s, false)
	}
	overflow := seedItem(repo, owner, "helm_iron", 1, 1, -1, false)
	escrow := uuid.New()
	repo.Seed(inventory.BagItem{InstanceID: escrow, TemplateCode: "boots_iron", Slot: -1, StackQuantity: 1, MaxStack: 1,
		IsLocked: true}, owner, false)
	worn := "chest"
	repo.Seed(inventory.BagItem{InstanceID: uuid.New(), TemplateCode: "chest_leather", Slot: -1, StackQuantity: 1,
		MaxStack: 1, EquipmentSlot: &worn}, owner, false)

	view, err := svc.View(context.Background(), owner)
	require.NoError(t, err)
	require.Equal(t, inventory.BagSlotCap, view.SlotsUsed, "bag full: overflow waits")

	view, err = svc.View(context.Background(), owner)
	require.NoError(t, err)
	require.NoError(t, svc.Delete(context.Background(), owner, view.Items[0].InstanceID))
	view, err = svc.View(context.Background(), owner)
	require.NoError(t, err)
	require.Equal(t, inventory.BagSlotCap, view.SlotsUsed)
	ids := map[uuid.UUID]bool{}
	for _, it := range view.Items {
		ids[it.InstanceID] = true
	}
	require.True(t, ids[overflow], "freed slot is taken by the overflow item")
	require.False(t, ids[escrow], "escrowed items stay out of the bag")
}

// A stack that a concurrent View already slotted must not merge into itself
// (that would delete the item).
func TestAssignSlots_NeverMergesIntoItself(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	id := seedItem(repo, owner, "iron_ore", 5, 99, 3, true)
	inst := mkInstance(owner, "iron_ore", 5)
	inst.ID = id

	require.NoError(t, svc.AssignSlots(context.Background(), owner, []item.Instance{inst}))
	view, err := svc.View(context.Background(), owner)
	require.NoError(t, err)
	require.Len(t, view.Items, 1, "the item still exists")
	require.Equal(t, int32(5), view.Items[0].StackQuantity, "and was not doubled")
}

// --- helpers ---------------------------------------------------------------

func mkInstance(owner uuid.UUID, tmpl string, stack int32) item.Instance {
	o := owner
	return item.Instance{
		ID:               uuid.New(),
		TemplateCode:     tmpl,
		OwnerCharacterID: &o,
		StackQuantity:    stack,
		Rarity:           item.RarityCommon,
		CreatedAt:        time.Now(),
	}
}

// HasRoom / FreeSlots (DD-246): a full bag still takes stackables that fit on
// an existing stack, nothing else.
func TestHasRoom_FullBag(t *testing.T) {
	svc, repo := newSvc()
	owner := uuid.New()
	ctx := context.Background()
	ok, err := svc.HasRoom(ctx, owner, "", 1)
	require.NoError(t, err)
	require.True(t, ok)
	for s := inventory.Slot(0); s < inventory.BagSlotCap-1; s++ {
		seedItem(repo, owner, "sword_iron", 1, 1, s, false)
	}
	free, err := svc.FreeSlots(ctx, owner)
	require.NoError(t, err)
	require.Equal(t, 1, free)
	seedItem(repo, owner, "iron_ore", 90, 99, inventory.BagSlotCap-1, true) // last slot: a stack with 9 room

	free, _ = svc.FreeSlots(ctx, owner)
	require.Zero(t, free)
	for _, c := range []struct {
		code string
		qty  int32
		want bool
	}{{"", 1, false}, {"sword_iron", 1, false}, {"iron_ore", 9, true}, {"iron_ore", 10, false}} {
		ok, err := svc.HasRoom(ctx, owner, c.code, c.qty)
		require.NoError(t, err)
		require.Equal(t, c.want, ok, "%s x%d", c.code, c.qty)
	}
}
