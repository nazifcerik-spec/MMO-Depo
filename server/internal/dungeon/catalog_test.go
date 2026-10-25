package dungeon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The catalog lists lowest level first, then by code, however it was built:
// the dungeon and raid lists show in that order.
func TestCatalogAllSortedByLevel(t *testing.T) {
	c := NewCatalog([]Template{
		{Code: "c_high", MinLevel: 40},
		{Code: "b_low", MinLevel: 5},
		{Code: "a_mid", MinLevel: 20},
		{Code: "a_low", MinLevel: 5},
	})
	for range 5 { // map order varies run to run
		var codes []string
		for _, t := range c.All() {
			codes = append(codes, t.Code)
		}
		require.Equal(t, []string{"a_low", "b_low", "a_mid", "c_high"}, codes)
	}
}
