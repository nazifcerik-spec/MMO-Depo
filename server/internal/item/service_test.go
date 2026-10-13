package item_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/item"
)

// A minimal in-memory catalog with the same templates seeded by 0007 migration
// (subset — enough to cover mining + weapons for the tests).
func testCatalog() *item.Catalog {
	return item.NewCatalog(append([]item.Template{
		{Code: "iron_ore", ItemClass: "material", ItemSubclass: "ore", BindKind: "boe", BaseValueCopper: 5},
		{Code: "copper_ore", ItemClass: "material", ItemSubclass: "ore", BindKind: "boe", BaseValueCopper: 8},
		{Code: "tin_ore", ItemClass: "material", ItemSubclass: "ore", BindKind: "boe", BaseValueCopper: 7},
		{Code: "silver_ore", ItemClass: "material", ItemSubclass: "ore", BindKind: "boe", BaseValueCopper: 25},
		{Code: "gold_ore", ItemClass: "material", ItemSubclass: "ore", BindKind: "boe", BaseValueCopper: 40},
		{Code: "mythril_ore", ItemClass: "material", ItemSubclass: "ore", BindKind: "boe", BaseValueCopper: 250},
		{Code: "sword_iron", ItemClass: "weapon", ItemSubclass: "sword", BindKind: "boe", Affixable: true},
	}, poolTemplates()...))
}

// poolTemplates covers every code the loot tables can roll, so no draw is
// skipped as unknown in these tests.
func poolTemplates() []item.Template {
	seen := map[string]bool{"iron_ore": true, "copper_ore": true, "tin_ore": true, "silver_ore": true, "gold_ore": true, "mythril_ore": true}
	var out []item.Template
	for _, p := range item.PoolPerKind {
		for _, tier := range [][]string{p.Common, p.Uncommon, p.Rare, p.Epic, p.Legendary, p.Mythic} {
			for _, c := range tier {
				if !seen[c] {
					seen[c] = true
					out = append(out, item.Template{Code: c, ItemClass: "material", BindKind: "boe"})
				}
			}
		}
	}
	return out
}

func newTestService() (*item.Service, *item.MemoryRepo) {
	repo := item.NewMemoryRepo()
	return item.NewService(repo, testCatalog(), slog.New(slog.NewTextHandler(io.Discard, nil))), repo
}

func TestCatalog_Lookup(t *testing.T) {
	cat := testCatalog()
	tmpl, err := cat.Get("iron_ore")
	require.NoError(t, err)
	require.Equal(t, "material", tmpl.ItemClass)

	_, err = cat.Get("unknown")
	require.ErrorIs(t, err, item.ErrUnknownTemplate)
}

func TestRollDrops_DeterministicGivenSeed(t *testing.T) {
	cat := testCatalog()
	d1 := item.RollDrops(cat, 4242, activity.KindMining, 5, 2, 1)
	d2 := item.RollDrops(cat, 4242, activity.KindMining, 5, 2, 1)
	require.Equal(t, d1, d2, "same seed must yield same drops")
}

func TestRollDrops_DifferentSeedsDiverge(t *testing.T) {
	cat := testCatalog()
	d1 := item.RollDrops(cat, 1, activity.KindMining, 10, 0, 0)
	d2 := item.RollDrops(cat, 2, activity.KindMining, 10, 0, 0)
	// With single-template pool for common, template_code is same; but the RNG
	// sequence differs — assert the SEQUENCE isn't identical for larger pools.
	d3 := item.RollDrops(cat, 1, activity.KindMining, 0, 10, 0)
	d4 := item.RollDrops(cat, 2, activity.KindMining, 0, 10, 0)
	require.Equal(t, len(d1), len(d2))
	require.NotEqual(t, d3, d4, "rare-tier pool has 2 templates — different seed → different sequence")
}

func TestRollDrops_RarityCountsExact(t *testing.T) {
	cat := testCatalog()
	drops := item.RollDrops(cat, 42, activity.KindMining, 3, 2, 1)
	require.Len(t, drops, 6)

	byRarity := map[item.Rarity]int{}
	for _, d := range drops {
		byRarity[d.Rarity]++
	}
	// Each yield tier keeps its count; part of it may step up one rarity.
	require.Equal(t, 3, byRarity[item.RarityCommon]+byRarity[item.RarityUncommon])
	require.Equal(t, 2, byRarity[item.RarityRare]+byRarity[item.RarityEpic])
	require.Equal(t, 1, byRarity[item.RarityLegendary]+byRarity[item.RarityMythic])
}

// Every rarity of the mining ladder drops, at roughly the configured split
// (DD-243) — before it, uncommon/epic/mythic tables were never rolled.
func TestRollDrops_UpgradeSplitReachesEveryTier(t *testing.T) {
	cat := testCatalog()
	byRarity := map[item.Rarity]int{}
	codes := map[string]bool{}
	for seed := int64(0); seed < 400; seed++ {
		for _, d := range item.RollDrops(cat, seed, activity.KindMining, 10, 10, 10) {
			byRarity[d.Rarity]++
			codes[d.TemplateCode] = true
		}
	}
	for _, r := range []item.Rarity{item.RarityCommon, item.RarityUncommon, item.RarityRare, item.RarityEpic, item.RarityLegendary, item.RarityMythic} {
		require.Positive(t, byRarity[r], r)
	}
	share := func(up, base item.Rarity) float64 {
		return float64(byRarity[up]) / float64(byRarity[up]+byRarity[base])
	}
	require.InDelta(t, 0.25, share(item.RarityUncommon, item.RarityCommon), 0.03)
	require.InDelta(t, 0.25, share(item.RarityEpic, item.RarityRare), 0.03)
	require.InDelta(t, 0.15, share(item.RarityMythic, item.RarityLegendary), 0.03)
	require.True(t, codes["copper_ore"] && codes["adamantite_ore"], "both ends of the ladder drop")
	// Weight by repetition: coal is listed twice, lead once.
	coal, lead := 0, 0
	for seed := int64(0); seed < 400; seed++ {
		for _, d := range item.RollDrops(cat, seed, activity.KindMining, 20, 0, 0) {
			switch d.TemplateCode {
			case "coal":
				coal++
			case "lead_ore":
				lead++
			}
		}
	}
	require.InDelta(t, 2.0, float64(coal)/float64(lead), 0.35)
}

func TestRollDrops_UnknownKindReturnsNil(t *testing.T) {
	cat := testCatalog()
	require.Nil(t, item.RollDrops(cat, 1, "unicorn", 5, 0, 0))
}

func TestRarityAffixRules(t *testing.T) {
	// Common/Uncommon get NO affix (rule)
	require.False(t, item.AffixRuleFor(item.RarityCommon).HasPrefix)
	require.False(t, item.AffixRuleFor(item.RarityCommon).OptionalPrefix)
	require.False(t, item.AffixRuleFor(item.RarityUncommon).HasPrefix)

	// Legendary+ guaranteed both
	require.True(t, item.AffixRuleFor(item.RarityLegendary).HasPrefix)
	require.True(t, item.AffixRuleFor(item.RarityLegendary).HasSuffix)
	require.True(t, item.AffixRuleFor(item.RarityMythic).HasPrefix)
}

func TestValueMultiplier_Monotonic(t *testing.T) {
	last := 0.0
	for _, r := range []item.Rarity{
		item.RarityCommon, item.RarityUncommon, item.RarityRare,
		item.RarityEpic, item.RarityLegendary, item.RarityMythic,
	} {
		v := item.ValueMultiplier(r)
		require.Greater(t, v, last, "rarity %s value must exceed lower tier", r)
		last = v
	}
}

func TestService_MaterializeDrops_PersistsInstances(t *testing.T) {
	svc, repo := newTestService()
	char := uuid.New()
	src := uuid.New()

	drops := item.RollDrops(svc.Catalog(), 999, activity.KindMining, 3, 1, 0)
	require.Len(t, drops, 4)

	instances, err := svc.MaterializeDrops(context.Background(), char, "activity_reveal", src, drops)
	require.NoError(t, err)
	require.Len(t, instances, 4)

	// Persisted.
	got, err := repo.ListByOwner(context.Background(), char, 100)
	require.NoError(t, err)
	require.Len(t, got, 4)

	// All owned by the character + tagged with source.
	for _, i := range got {
		require.NotNil(t, i.OwnerCharacterID)
		require.Equal(t, char, *i.OwnerCharacterID)
		require.NotNil(t, i.SourceRefID)
		require.Equal(t, src, *i.SourceRefID)
	}
}

func TestService_Grant_EndToEnd(t *testing.T) {
	svc, _ := newTestService()
	char := uuid.New()
	src := uuid.New()

	instances, err := svc.Grant(context.Background(), char, src,
		1234, activity.KindMining, 2, 1, 0, activity.GearDrop{})
	require.NoError(t, err)
	require.Len(t, instances, 3)

	// Deterministic — same seed => same TemplateCode / Rarity sequence.
	instances2, err := svc.Grant(context.Background(), char, uuid.New(),
		1234, activity.KindMining, 2, 1, 0, activity.GearDrop{})
	require.NoError(t, err)
	for i := range instances {
		require.Equal(t, instances[i].TemplateCode, instances2[i].TemplateCode)
		require.Equal(t, instances[i].Rarity, instances2[i].Rarity)
	}
}

func TestService_ListBySource(t *testing.T) {
	svc, _ := newTestService()
	char := uuid.New()
	srcA := uuid.New()
	srcB := uuid.New()

	_, _ = svc.Grant(context.Background(), char, srcA, 1, activity.KindMining, 2, 0, 0, activity.GearDrop{})
	_, _ = svc.Grant(context.Background(), char, srcB, 2, activity.KindMining, 3, 0, 0, activity.GearDrop{})

	a, _ := svc.ListBySource(context.Background(), srcA)
	require.Len(t, a, 2)
	b, _ := svc.ListBySource(context.Background(), srcB)
	require.Len(t, b, 3)
}
