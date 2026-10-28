// Package story is the story mode's chapter list (DD-313): read from
// content/story.yaml, names and text from content/i18n (story.<code> is the
// title, story.<code>.body the text). The content comes later; an empty
// list shows the client's "coming soon" page.
package story

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/otter/server/internal/http/apierr"
)

// Chapter is one story chapter, opened by level and, optionally, by the
// chapter before it.
type Chapter struct {
	Code     string `yaml:"code" json:"code"`
	MinLevel int    `yaml:"min_level" json:"min_level"`
	Requires string `yaml:"requires,omitempty" json:"requires,omitempty"`
}

var current atomic.Pointer[[]Chapter]

// Validate checks codes are unique, levels in 1..1000 and every required
// chapter comes earlier in the list.
func Validate(chs []Chapter) error {
	seen := map[string]bool{}
	for _, c := range chs {
		switch {
		case c.Code == "":
			return fmt.Errorf("story: chapter without code")
		case seen[c.Code]:
			return fmt.Errorf("story %s: duplicate code", c.Code)
		case c.MinLevel < 1 || c.MinLevel > 1000:
			return fmt.Errorf("story %s: min_level must be 1..1000", c.Code)
		case c.Requires != "" && !seen[c.Requires]:
			return fmt.Errorf("story %s: requires %q, not an earlier chapter", c.Code, c.Requires)
		}
		seen[c.Code] = true
	}
	return nil
}

// Install swaps in the chapters of a content load.
func Install(chs []Chapter) error {
	if err := Validate(chs); err != nil {
		return err
	}
	cp := append([]Chapter{}, chs...)
	current.Store(&cp)
	return nil
}

// All is the installed chapters in order.
func All() []Chapter {
	if p := current.Load(); p != nil {
		return *p
	}
	return []Chapter{}
}

// CatalogHandler is GET /v1/story (public, cacheable).
func CatalogHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"chapters": All()})
}
