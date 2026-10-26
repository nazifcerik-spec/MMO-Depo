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
	jobs map[uuid.UUID]*Job
}

func newMemRepo() *memRepo {
	return &memRepo{rows: map[string]*CharacterProfession{}, jobs: map[uuid.UUID]*Job{}}
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

func (m *memRepo) CreateJob(_ context.Context, j *Job) error {
	for _, x := range m.jobs {
		if x.CharacterID == j.CharacterID && x.State == JobActive {
			return ErrJobActive
		}
	}
	c := *j
	m.jobs[j.ID] = &c
	return nil
}

func (m *memRepo) ActiveJob(_ context.Context, id uuid.UUID) (*Job, error) {
	for _, x := range m.jobs {
		if x.CharacterID == id && x.State == JobActive {
			c := *x
			return &c, nil
		}
	}
	return nil, ErrNotFound
}

func (m *memRepo) GetJob(_ context.Context, id uuid.UUID) (*Job, error) {
	if x, ok := m.jobs[id]; ok {
		c := *x
		return &c, nil
	}
	return nil, ErrNotFound
}

func (m *memRepo) DeleteJob(_ context.Context, id uuid.UUID) error { delete(m.jobs, id); return nil }

func (m *memRepo) FinishJob(_ context.Context, id uuid.UUID, _ int, _ time.Time) (bool, error) {
	x, ok := m.jobs[id]
	if !ok || x.State != JobActive {
		return false, nil
	}
	x.State = JobCollected
	return true, nil
}

func (m *memRepo) SaveJobResult(_ context.Context, id uuid.UUID, result []byte) error {
	m.jobs[id].Result = result
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
	have     int32 // 0 = plenty
	full     bool  // HasRoom says no
	failAt   int   // ConsumeAll fails once this many calls were made (0 = never)
}

func (f *fakeInv) Available(_ context.Context, _ uuid.UUID, tmpl string) (int32, error) {
	if f.shortage[tmpl] {
		return 0, nil
	}
	if f.have > 0 {
		return f.have, nil
	}
	return 1 << 20, nil
}

func (f *fakeInv) HasRoom(context.Context, uuid.UUID, string, int32) (bool, error) {
	return !f.full, nil
}

func (f *fakeInv) ConsumeAll(_ context.Context, _ uuid.UUID, needs []inventory.Need) error {
	if f.failAt < 0 || (f.failAt > 0 && len(f.calls) >= f.failAt) {
		return inventory.ErrInsufficient
	}
	for _, n := range needs {
		f.calls = append(f.calls, consumeCall{n.TemplateCode, n.Quantity})
	}
	return nil
}

type fakeItemMinter struct {
	lastMint *item.Instance
	minted   map[string]int32
}

func (f *fakeItemMinter) MintExact(_ context.Context, ownerID uuid.UUID, tmpl string, qty int32, srcKind string, srcRef uuid.UUID) (*item.Instance, error) {
	inst := &item.Instance{
		ID: uuid.New(), TemplateCode: tmpl, OwnerCharacterID: &ownerID,
		StackQuantity: qty, Rarity: "common",
	}
	f.lastMint = inst
	if f.minted == nil {
		f.minted = map[string]int32{}
	}
	f.minted[tmpl] += qty
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
			Code: "recipe_potion", ProfessionCode: "alchemy", MinSkill: 1, CraftSeconds: 30,
			Inputs:             []RecipeInput{{TemplateCode: "linseed", Quantity: 2}},
			GoldCostCopper:     0, // skip ledger in tests
			OutputTemplateCode: "potion_hp_small", OutputStack: 1, ProfessionXPReward: 100,
		}, {
			Code: "recipe_forge", ProfessionCode: "alchemy", MinSkill: 10, // above starting level
			Inputs:             []RecipeInput{{TemplateCode: "iron_ore", Quantity: 3}},
			OutputTemplateCode: "sword_iron", OutputStack: 1, ProfessionXPReward: 40,
		}, {
			Code: "recipe_slow", ProfessionCode: "alchemy", MinSkill: 1, CraftSeconds: 4 * 3600,
			Inputs:             []RecipeInput{{TemplateCode: "linseed", Quantity: 1}},
			OutputTemplateCode: "potion_hp_small", OutputStack: 1,
		}},
	)
	repo := newMemRepo()
	inv := &fakeInv{}
	minter := &fakeItemMinter{}
	svc := NewService(repo, &fakeCatalog{cat}, inv, minter,
		economy.NewMemoryTxnService(map[string]economy.Currency{}),
		fakeOwner{},
		ServiceConfig{Now: func() time.Time { return clock }}, nil)
	clock = now
	return svc, repo, inv, minter
}

// now is the tests' start time; clock is what the service reads.
var now = time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
var clock = now

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

// learned starts a character with alchemy learned.
func learned(t *testing.T, svc *Service) (uuid.UUID, uuid.UUID) {
	t.Helper()
	charID, acct := uuid.New(), uuid.New()
	_, err := svc.Learn(context.Background(), charID, acct, "alchemy")
	require.NoError(t, err)
	return charID, acct
}

func TestStart_ValidatesBeforeSpending(t *testing.T) {
	svc, _, inv, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := learned(t, svc)

	_, err := svc.Start(ctx, charID, acct, "recipe_potion", 0)
	require.ErrorIs(t, err, ErrValidation)
	_, err = svc.Start(ctx, charID, acct, "recipe_potion", MaxCraftCount+1)
	require.ErrorIs(t, err, ErrValidation)
	_, err = svc.Start(ctx, charID, acct, "recipe_ghost", 1)
	require.ErrorIs(t, err, ErrUnknownRecipe)
	_, err = svc.Start(ctx, charID, acct, "recipe_forge", 1) // min_skill 10 vs level 1
	require.ErrorIs(t, err, ErrSkillTooLow)
	_, err = svc.Start(ctx, uuid.New(), acct, "recipe_potion", 1)
	require.ErrorIs(t, err, ErrProfessionNotLearned)
	_, err = svc.Start(ctx, charID, acct, "recipe_slow", 2) // 2 × 4 h > 6 h
	require.ErrorIs(t, err, ErrJobTooLong)

	inv.have = 5 // 2 per unit: 2 units possible, 3 asked
	_, err = svc.Start(ctx, charID, acct, "recipe_potion", 3)
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	inv.have, inv.full = 0, true
	_, err = svc.Start(ctx, charID, acct, "recipe_potion", 1)
	require.ErrorIs(t, err, inventory.ErrBagFull)
	require.Empty(t, inv.calls, "nothing consumed")
}

func TestJob_CollectDeliversFinishedAndRefundsTheRest(t *testing.T) {
	svc, repo, inv, minter := newTestSvc(t)
	ctx := context.Background()
	charID, acct := learned(t, svc)

	job, err := svc.Start(ctx, charID, acct, "recipe_potion", 4) // 30 s each
	require.NoError(t, err)
	require.Equal(t, now.Add(2*time.Minute), job.ReadyAt)
	require.Equal(t, []consumeCall{{Template: "linseed", Qty: 8}}, inv.calls, "the whole job's inputs at start")
	_, err = svc.Start(ctx, charID, acct, "recipe_potion", 1)
	require.ErrorIs(t, err, ErrJobActive)
	active, err := svc.Active(ctx, charID, acct)
	require.NoError(t, err)
	require.Equal(t, job.ID, active.ID)

	clock = now.Add(65 * time.Second) // two units done
	require.Equal(t, 2, job.Done(clock))
	res, err := svc.Collect(ctx, charID, acct, job.ID)
	require.NoError(t, err)
	require.Equal(t, 2, res.Made)
	require.Equal(t, 2, res.Refunded)
	require.Equal(t, int32(2), res.OutputStack)
	require.Equal(t, int32(200), res.XPGained)
	require.True(t, res.LeveledUp)
	require.Equal(t, int16(1), res.SkillBefore)
	require.Equal(t, 1, repo.logs)
	require.Equal(t, map[string]int32{"potion_hp_small": 2, "linseed": 4}, minter.minted, "two potions; four linseed back")

	again, err := svc.Collect(ctx, charID, acct, job.ID)
	require.NoError(t, err)
	require.True(t, again.Replayed)
	require.Equal(t, 2, again.Made)
	require.Equal(t, int32(2), minter.minted["potion_hp_small"], "collected once")
	idle, err := svc.Active(ctx, charID, acct)
	require.NoError(t, err)
	require.Nil(t, idle)

	_, err = svc.Collect(ctx, uuid.New(), acct, job.ID)
	require.ErrorIs(t, err, ErrNotFound, "someone else's job")
}

func TestJob_ConsumeRaceUndoesTheStart(t *testing.T) {
	svc, repo, inv, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := learned(t, svc)
	inv.failAt = -1
	_, err := svc.Start(ctx, charID, acct, "recipe_potion", 1)
	require.ErrorIs(t, err, inventory.ErrInsufficient)
	require.Empty(t, repo.jobs, "the job row is gone")
	require.Equal(t, 3, (&Job{Count: 3}).Done(now), "no unit time: all done")
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

func TestDifficulty(t *testing.T) {
	for _, c := range []struct {
		min, skill int16
		color      string
		factor     float64
	}{
		{10, 10, DifficultyOrange, 1}, {10, 14, DifficultyOrange, 1},
		{10, 15, DifficultyYellow, 0.6}, {10, 19, DifficultyYellow, 0.6},
		{10, 20, DifficultyGreen, 0.25}, {10, 29, DifficultyGreen, 0.25},
		{10, 30, DifficultyGrey, 0}, {1, 100, DifficultyGrey, 0},
	} {
		color, f := Difficulty(c.min, c.skill)
		require.Equal(t, c.color, color, "%d at %d", c.min, c.skill)
		require.Equal(t, c.factor, f)
	}
}

type fakeSpeed float64

func (f fakeSpeed) ProfessionBonus(_ context.Context, _ uuid.UUID, target string) (float64, error) {
	if target != "craft_speed" {
		return 0, nil
	}
	return float64(f), nil
}

// A job freezes the XP the recipe's difficulty pays at the skill it starts
// at (DD-294) and runs faster with crafting-speed enchants (DD-293).
func TestStart_DifficultyAndSpeed(t *testing.T) {
	svc, repo, _, _ := newTestSvc(t)
	ctx := context.Background()
	charID, acct := learned(t, svc)

	job, err := svc.Start(ctx, charID, acct, "recipe_potion", 1)
	require.NoError(t, err)
	require.Equal(t, int32(100), job.Recipe.ProfessionXPReward, "orange: full XP")
	require.Equal(t, int32(30), job.UnitSeconds)
	require.NoError(t, repo.DeleteJob(ctx, job.ID))

	repo.rows[key(charID, "alchemy")].Level = 8 // 7 past the recipe: yellow
	job, err = svc.Start(ctx, charID, acct, "recipe_potion", 1)
	require.NoError(t, err)
	require.Equal(t, int32(60), job.Recipe.ProfessionXPReward)
	require.NoError(t, repo.DeleteJob(ctx, job.ID))

	repo.rows[key(charID, "alchemy")].Level = 40 // grey: teaches nothing
	svc.WithSpeed(fakeSpeed(0.25))
	job, err = svc.Start(ctx, charID, acct, "recipe_potion", 2)
	require.NoError(t, err)
	require.Zero(t, job.Recipe.ProfessionXPReward)
	require.Equal(t, int32(23), job.UnitSeconds, "30 s less a quarter")
	require.Equal(t, now.Add(46*time.Second), job.ReadyAt)
	require.NoError(t, repo.DeleteJob(ctx, job.ID))

	svc.WithSpeed(fakeSpeed(0.9))
	job, err = svc.Start(ctx, charID, acct, "recipe_potion", 1)
	require.NoError(t, err)
	require.Equal(t, int32(15), job.UnitSeconds, "capped at half")
}
