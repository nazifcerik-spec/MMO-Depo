package contentpack

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-363: a guardian shield must be wearable when it is on offer. The branch
// opens at level 300, so the ladder starts one tier below it and the first
// shield is in the shop: nobody reaches the branch without a way to get one.
func TestGuardianShieldLadderStartsAtThePromotion(t *testing.T) {
	p, err := Read("../../../content")
	require.NoError(t, err)

	shields := map[string]bool{}
	first := int16(0)
	for _, i := range p.Items {
		if i.ItemSubclass != "shield" {
			continue
		}
		shields[i.Code] = true
		if first == 0 || i.MinLevel < first {
			first = i.MinLevel
		}
	}
	require.Equal(t, int16(285), first)

	var shop []string
	for _, v := range p.Vendors {
		if v.Code == "general_ironhold" {
			for _, s := range v.Sells {
				shop = append(shop, s.TemplateCode)
			}
		}
	}
	require.Contains(t, shop, "shield_sovereign")

	for _, r := range p.Recipes {
		require.False(t, slices.Contains([]string{"shield_royal", "shield_phoenix", "shield_dreadforged"}, r.OutputTemplateCode), r.Code)
	}
	for _, code := range []string{"shield_royal", "shield_phoenix", "shield_dreadforged"} {
		require.False(t, shields[code], code)
		require.NotContains(t, p.Strings.ByLang["en"], "item."+code)
	}
}
