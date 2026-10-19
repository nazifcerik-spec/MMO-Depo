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

// A minimal in-memory catalog: a mining ladder, a hunting ladder and a weapon.
func testCatalog() *item.Catalog {
	node := func(code, kind string, level, tier int16) item.Template {
		return item.Template{Code: code, ItemClass: "material", ItemSubclass: "ore", BindKind: "boe",
			BaseValueCopper: 5, Stackable: true, MaxStack: 200, GatherKind: kind, GatherLevel: level, GatherTier: tier}
	}
	return item.NewCatalog([]item.Template{
		node("raw_stone", "mining", 1, 0), node("iron_ore", "mining", 16, 1),
		node("copper_ore", "mining", 20, 2), node("mythril_ore", "mining", 90, 7),
		node("forest_rat_pelt", "hunting", 1, 0), node("hare_pelt", "hunting", 4, 0),
		node("boar_hide", "hunting", 22, 2), node("wolf_pelt", "hunting", 26, 2),
		node("bear_hide", "hunting", 40, 3), node("dire_bear_hide", "hunting", 60, 5),
		node("dragon_scale", "hunting", 100, 7),
		{Code: "sword_iron", ItemClass: "weapon", ItemSubclass: "sword", BindKind: "boe", Affixable: true},
	})
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

func TestCatalog_NodesSortedByLevel(t *testing.T) {
	var codes []string
	for _, n := range testCatalog().Nodes("mining") {
		codes = append(codes, n.Code)
	}
	require.Equal(t, []string{"raw_stone", "iron_ore", "copper_ore", "mythril_ore"}, codes)
	require.Empty(t, testCatalog().Nodes("unicorn"))
}

func TestTierRarity(t *testing.T) {
	want := []item.Rarity{item.RarityCommon, item.RarityCommon, item.RarityUncommon, item.RarityUncommon,
		item.RarityRare, item.RarityEpic, item.RarityLegendary, item.RarityMythic}
	for tier, r := range want {
		require.Equal(t, r, item.TierRarity(int16(tier)), "tier %d", tier)
	}
}

func TestGatherDrops_AllUnitsOfTheTarget(t *testing.T) {
	cat := testCatalog()
	drops := item.GatherDrops(cat, "mythril_ore", 3)
	require.Len(t, drops, 3)
	for _, d := range drops {
		require.Equal(t, "mythril_ore", d.TemplateCode)
		require.Equal(t, item.RarityMythic, d.Rarity)
	}
	require.Nil(t, item.GatherDrops(cat, "sword_iron", 3), "not a node")
	require.Nil(t, item.GatherDrops(cat, "nope", 3))
	require.Nil(t, item.GatherDrops(cat, "iron_ore", 0))
}

func TestCombatHides_ByMobLevel(t *testing.T) {
	cat := testCatalog()
	// A level-230 mob reaches hunting 23 (level / 10, DD-274): the best four
	// unlocked are rat, hare, boar (22) — wolf (26) is out of reach. Rare
	// finds are the best: boar.
	seen := map[string]bool{}
	for seed := int64(0); seed < 50; seed++ {
		for _, d := range item.CombatHides(cat, seed, 230, 4, 0, 0) {
			seen[d.TemplateCode] = true
		}
	}
	require.Equal(t, map[string]bool{"forest_rat_pelt": true, "hare_pelt": true, "boar_hide": true}, seen)
	rare := item.CombatHides(cat, 1, 230, 0, 2, 1)
	require.Len(t, rare, 3)
	for _, d := range rare {
		require.Equal(t, "boar_hide", d.TemplateCode)
		require.Equal(t, item.RarityUncommon, d.Rarity)
	}
	require.Equal(t, "forest_rat_pelt", item.CombatHides(cat, 1, 6, 0, 1, 0)[0].TemplateCode, "level 6 reaches 1")
	// Level 991+ reaches 100 (cap): dragon scale for the best finds.
	require.Equal(t, "bear_hide", item.CombatHides(cat, 1, 400, 0, 1, 0)[0].TemplateCode)
	require.Equal(t, "dragon_scale", item.CombatHides(cat, 1, 1000, 0, 1, 0)[0].TemplateCode)
	require.Equal(t, item.CombatHides(cat, 9, 12, 5, 1, 0), item.CombatHides(cat, 9, 12, 5, 1, 0), "deterministic")
	require.Empty(t, item.CombatHides(item.NewCatalog(nil), 1, 10, 3, 0, 0))
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

	drops := item.GatherDrops(svc.Catalog(), "iron_ore", 4)
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

	instances, err := svc.Grant(context.Background(), char, src, 1234,
		activity.Drops{Kind: activity.KindMining, Target: "copper_ore", Common: 2, Rare: 1})
	require.NoError(t, err)
	require.Len(t, instances, 3)
	for _, in := range instances {
		require.Equal(t, "copper_ore", in.TemplateCode, "every find is the chosen resource")
	}

	// Combat: hides by mob level, deterministic per seed.
	fight := activity.Drops{Kind: activity.KindHunting, MobLevel: 6, Common: 3, Rare: 1}
	instances, err = svc.Grant(context.Background(), char, uuid.New(), 1234, fight)
	require.NoError(t, err)
	require.Len(t, instances, 4)
	instances2, err := svc.Grant(context.Background(), char, uuid.New(), 1234, fight)
	require.NoError(t, err)
	for i := range instances {
		require.Equal(t, instances[i].TemplateCode, instances2[i].TemplateCode)
		require.Equal(t, instances[i].Rarity, instances2[i].Rarity)
	}

	none, err := svc.Grant(context.Background(), char, uuid.New(), 1, activity.Drops{Kind: activity.KindMining, Common: 5})
	require.NoError(t, err)
	require.Empty(t, none, "no target, no mob: nothing")
}

func TestService_ListBySource(t *testing.T) {
	svc, _ := newTestService()
	char := uuid.New()
	srcA := uuid.New()
	srcB := uuid.New()

	_, _ = svc.Grant(context.Background(), char, srcA, 1, activity.Drops{Kind: activity.KindMining, Target: "iron_ore", Common: 2})
	_, _ = svc.Grant(context.Background(), char, srcB, 2, activity.Drops{Kind: activity.KindMining, Target: "iron_ore", Common: 3})

	a, _ := svc.ListBySource(context.Background(), srcA)
	require.Len(t, a, 2)
	b, _ := svc.ListBySource(context.Background(), srcB)
	require.Len(t, b, 3)
}
