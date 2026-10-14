package item_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/item"
)

// The shipped gathering tables (DD-244): every profession starts at a
// level-1 resource, and a combat hide exists for every mob level.
func TestShippedGatheringTables(t *testing.T) {
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	var tmpls []item.Template
	for _, it := range p.Items {
		if it.Gather != nil {
			tmpls = append(tmpls, item.Template{Code: it.Code, GatherKind: it.Gather.Kind,
				GatherLevel: int16(it.Gather.Level), GatherTier: int16(it.Gather.Tier)})
		}
	}
	cat := item.NewCatalog(tmpls)
	for kind := range contentpack.GatherKinds {
		nodes := cat.Nodes(kind)
		require.NotEmpty(t, nodes, kind)
		require.Equal(t, int16(1), nodes[0].GatherLevel, "%s needs a level-1 resource", kind)
	}
	for lvl := 1; lvl <= 30; lvl++ {
		require.NotEmpty(t, item.CombatHides(cat, 1, lvl, 1, 1, 0), "mob level %d", lvl)
	}
}
