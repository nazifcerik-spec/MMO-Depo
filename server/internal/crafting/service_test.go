package crafting

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
)

// --- Fakes ------------------------------------------------------------------

type memRepo struct {
	rows map[string]*CharacterProfession // "charID|profession" → row
	logs int
}

func newMemRepo() *memRepo {
	return &memRepo{rows: map[string]*CharacterProfession{}}
}

func key(id uuid.UUID, code string) string { return id.String() + "|" + code }

func (m *memRepo) UpsertLearn(_ context.Context, id uuid.UUID, code string, now time.Time) (*CharacterProfession, error) {
	k := key(id, code)
	if cp, ok := m.rows[k]; ok {
		return cp, nil
	}
	cp := &CharacterProfession{
		CharacterID: id, ProfessionCode: code, Level: 1, XP: 0,
		LearnedAt: now, UpdatedAt: now, Version: 1,
	}
	m.rows[k] = cp
	return cp, nil
}

func (m *memRepo) Get(_ context.Context, id uuid.UUID, code string) (*CharacterProfession, error) {
	if cp, ok := m.rows[key(id, code)]; ok {
		return cp, nil
	}
	return nil, ErrNotFound
}

func (m *memRepo) List(_ context.Context, id uuid.UUID) ([]CharacterProfession, error) {
	out := []CharacterProfession{}
	for _, cp := range m.rows {
		if cp.CharacterID == id {
			out = append(out, *cp)
		}
	}
	return out, nil
}

func (m *memRepo) AddXP(_ context.Context, id uuid.UUID, code string, delta int32, maxLevel int16) (*CharacterProfession, bool, error) {
	cp, ok := m.rows[key(id, code)]
	if !ok {
		return nil, false, ErrProfessionNotLearned
	}
	leveled := false
	cp.XP += delta
	for {
		needed := ProfessionXPToNext(cp.Level)
		if cp.XP < needed || cp.Level >= maxLevel {
			break
		}
		cp.XP -= needed
		cp.Level++
		leveled = true
	}
	return cp, leveled, nil
}

func (m *memRepo) InsertCraftLog(_ context.Context, _ uuid.UUID, _ string, _ *uuid.UUID, _ *uuid.UUID, _ time.Time) error {
	m.logs++
	return nil
}

type fakeCatalog struct{ inner *Catalog }

func (f *fakeCatalog) Catalog() *Catalog { return f.inner }

type consumeCall struct {
	Template string
	Qty      int32
}
type fakeInv struct {
	calls    []consumeCall
	shortage map[string]bool
}

func (f *fakeInv) Available(_ context.Context, _ uuid.UUID, tmpl string) (int32, error) {
	if f.shortage[tmpl] {
		return 0, nil
	}
	return 1 << 20, nil
}

func (f *fakeInv) ConsumeAll(_ context.Context, _ uuid.UUID, needs []inventory.Need) error {
	for _, n := range needs {
		f.calls = append(f.calls, consumeCall{n.TemplateCode, n.Quantity})
	}
	return nil
}

type fakeItemMinter struct {
	lastMint *item.Instance
}

func (f *fakeItemMinter) MintExact(_ context.Context, ownerID uuid.UUID, tmpl string, qty int32, srcKind string, srcRef uuid.UUID) (*item.Instance, error) {
	inst := &item.Instance{
		ID: uuid.New(), TemplateCode: tmpl, OwnerCharacterID: &ownerID,
		StackQuantity: qty, Rarity: "common",
	}
	f.lastMint = inst
	return inst, nil
}

type fakeOwner struct{}

func (fakeOwner) Owns(_ context.Context, _ uuid.UUID, _ uuid.UUID) error { return nil }

// --- Tests ------------------------------------------------------------------

func newTestSvc(t *testing.T) (*Service, *memRepo, *fakeInv, *fakeItemMinter) {
	t.Helper()
	cat := NewCatalog(
		[]Profession{
			{Code: "alchemy", DisplayNameKey: "p.a", MaxLevel: 60, Kind: KindCrafting},
			{Code: "mining", DisplayNameKey: "p.m", MaxLevel: 100, Kind: KindGathering},
		},
		[]Recipe{{
			Code: "recipe_potion", ProfessionCode: "alchemy", MinSkill: 1,
			Inputs:             []RecipeInput{{TemplateCode: "linseed", Quantity: 2}},
			GoldCostCopper:     0, // skip ledger in tests
			OutputTemplateCode: "potion_hp_small", OutputStack: 1, ProfessionXPReward: 100,
		}, {
			Code: "recipe_forge", ProfessionCode: "alchemy", MinSkill: 10, // above starting level
			Inputs:             []RecipeInput{{TemplateCode: "iron_ore", Quantity: 3}},
			OutputTemplateCode: "sword_iron", OutputStack: 1, ProfessionXPReward: 40,
		}},
	)
	repo := newMemRepo()
	inv := &fakeInv{}
	minter := &fakeItemMinter{}
	svc := NewService(repo, &fakeCatalog{cat}, inv, minter,
		economy.NewMemoryTxnService(map[string]economy.Currency{}),
		fakeOwner{},
		DefaultServiceConfig(), nil)
	return svc, repo, inv, minter
}

func TestLearn_IdempotentReturnsRowUnchanged(t *testing.T) {
	svc, _, _, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := uuid.New(), uuid.New()
	first, err := svc.Learn(ctx, charID, acct, "alchemy")
	require.NoError(t, err)
	require.Equal(t, int16(1), first.Level)
	second, err := svc.Learn(ctx, charID, acct, "alchemy")
	require.NoError(t, err)
	require.Equal(t, first.LearnedAt, second.LearnedAt, "re-learn must not reset learned_at")
}

func TestLearn_GatheringProfessionIsNotLearnable(t *testing.T) {
	svc, _, _, _ := newTestSvc(t)
	_, err := svc.Learn(context.Background(), uuid.New(), uuid.New(), "mining")
	require.ErrorIs(t, err, ErrUnknownProfession, "gathering levels come from activities")
}

func TestCraft_HappyPathConsumesInputsAndGrantsXP(t *testing.T) {
	svc, repo, inv, minter := newTestSvc(t)
	ctx := context.Background()
	charID, acct := uuid.New(), uuid.New()
	_, err := svc.Learn(ctx, charID, acct, "alchemy")
	require.NoError(t, err)

	res, err := svc.Craft(ctx, charID, acct, "recipe_potion")
	require.NoError(t, err)
	require.Equal(t, "potion_hp_small", res.OutputTemplateCode)
	require.Equal(t, int32(1), res.OutputStack)
	require.NotZero(t, res.OutputInstanceID)
	require.Equal(t, int32(100), res.XPGained)
	require.True(t, res.LeveledUp, "100 xp at level 1 crosses the 64 threshold, not 64+67")
	require.Equal(t, int16(2), res.SkillAfter)

	require.Equal(t, []consumeCall{{Template: "linseed", Qty: 2}}, inv.calls)
	require.NotNil(t, minter.lastMint)
	require.Equal(t, 1, repo.logs, "audit log written")
}

func TestCraft_RefusesUnlearnedProfession(t *testing.T) {
	svc, _, _, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := uuid.New(), uuid.New()
	_, err := svc.Craft(ctx, charID, acct, "recipe_potion")
	require.ErrorIs(t, err, ErrProfessionNotLearned)
}

func TestCraft_RefusesLowSkill(t *testing.T) {
	svc, _, _, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := uuid.New(), uuid.New()
	_, err := svc.Learn(ctx, charID, acct, "alchemy")
	require.NoError(t, err)
	_, err = svc.Craft(ctx, charID, acct, "recipe_forge") // min_skill 10 vs level 1
	require.ErrorIs(t, err, ErrSkillTooLow)
}

func TestCraft_RefusesUnknownRecipe(t *testing.T) {
	svc, _, _, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := uuid.New(), uuid.New()
	_, err := svc.Learn(ctx, charID, acct, "alchemy")
	require.NoError(t, err)
	_, err = svc.Craft(ctx, charID, acct, "recipe_ghost")
	require.ErrorIs(t, err, ErrUnknownRecipe)
}

func TestCraft_InventoryShortageBubblesUpBeforeMint(t *testing.T) {
	svc, _, inv, minter := newTestSvc(t)
	ctx := context.Background()
	charID, acct := uuid.New(), uuid.New()
	_, err := svc.Learn(ctx, charID, acct, "alchemy")
	require.NoError(t, err)
	inv.shortage = map[string]bool{"linseed": true}

	_, err = svc.Craft(ctx, charID, acct, "recipe_potion")
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	require.Nil(t, minter.lastMint, "mint must not run when inputs failed")
	require.Empty(t, inv.calls, "the pre-check stops before anything is consumed (or charged)")
}

func TestProfessionXPToNext_MonotonicByLevel(t *testing.T) {
	require.Equal(t, int32(64), ProfessionXPToNext(1))
	require.Equal(t, int32(67), ProfessionXPToNext(2))
	require.Equal(t, int32(107), ProfessionXPToNext(10))
	require.Equal(t, int32(1105), ProfessionXPToNext(50))
	require.Equal(t, int32(19206), ProfessionXPToNext(99))
	for l := int16(1); l < 100; l++ {
		require.Greater(t, ProfessionXPToNext(l+1), ProfessionXPToNext(l))
	}
	// Clamps sub-1 levels to 1's cost.
	require.Equal(t, int32(64), ProfessionXPToNext(0))
	require.Equal(t, int32(64), ProfessionXPToNext(-5))
}
