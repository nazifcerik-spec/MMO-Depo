package itemset

import (
	"net/http"

	"github.com/otter/server/internal/http/apierr"
)

// CatalogHandler is GET /v1/item-sets: every set with its pieces and bonuses
// (public, cacheable); the client shows them on a piece's details.
func CatalogHandler(w http.ResponseWriter, _ *http.Request) {
	sets := Current().All()
	if sets == nil {
		sets = []Set{}
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"sets": sets})
}
