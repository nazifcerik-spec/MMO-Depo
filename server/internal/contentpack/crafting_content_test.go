package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The shipped crafting chain (DD-274): every random-drop piece of gear has
// a recipe at skill ceil(level/10) within its profession's cap, the gear
// professions use the bars smelting makes, and every crafting profession
// reaches its cap with recipes.
func TestShippedCraftingCoversGear(t *testing.T) {
	p := shippedPack(t)
	profs := map[string]Profession{}
	for _, pr := range p.Professions {
		profs[pr.Code] = pr
	}
	recipeFor := map[string]Recipe{}
	used := map[string]bool{}
	top := map[string]int16{}
	for _, r := range p.Recipes {
		recipeFor[r.OutputTemplateCode] = r
		for _, in := range r.Inputs {
			used[in.TemplateCode] = true
		}
		pr, ok := profs[r.ProfessionCode]
		require.True(t, ok, r.Code)
		require.LessOrEqual(t, r.MinSkill, pr.MaxLevel, r.Code)
		top[r.ProfessionCode] = max(top[r.ProfessionCode], r.MinSkill)
	}
	setPiece := map[string]bool{} // raid drops only (DD-277)
	for _, s := range p.ItemSets {
		for _, piece := range s.Pieces {
			setPiece[piece] = true
		}
	}
	for _, it := range p.Items {
		switch {
		case it.EquipSlot != "" && it.Requirements == nil && it.ItemSubclass != "shield" && !setPiece[it.Code]:
			r, ok := recipeFor[it.Code]
			require.True(t, ok, "%s has no recipe", it.Code)
			want := max(1, (int(it.MinLevel)+9)/10)
			if it.ItemClass == "jewelry" { // the jewelry ladder starts at skill 1 (DD-294)
				want = max(1, want-2)
			}
			require.Equal(t, int16(want), r.MinSkill, "%s: skill = ceil(level/10)", it.Code)
		case it.ItemSubclass == "bar":
			require.True(t, used[it.Code], "%s is used by a recipe", it.Code)
		}
	}
	for _, code := range []string{"blacksmithing", "leatherworking", "tailoring", "woodworking", "alchemy", "cooking"} {
		require.GreaterOrEqual(t, top[code], profs[code].MaxLevel*8/10, "%s has recipes near its cap", code)
	}
}
