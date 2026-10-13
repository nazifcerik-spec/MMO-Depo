package item_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/item"
)

// Every code the loot tables roll must exist in the shipped content: an
// unknown code is skipped silently at drop time, so it would simply never
// drop (DD-243).
func TestLootPoolsExistInContent(t *testing.T) {
	p, err := contentpack.Read("../../../content")
	require.NoError(t, err)
	known := map[string]bool{}
	for _, it := range p.Items {
		known[it.Code] = true
	}
	for kind, pool := range item.PoolPerKind {
		for _, tier := range [][]string{pool.Common, pool.Uncommon, pool.Rare, pool.Epic, pool.Legendary, pool.Mythic} {
			for _, code := range tier {
				require.True(t, known[code], "%s pool rolls unknown item %q", kind, code)
			}
		}
	}
}
