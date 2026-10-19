// Package itemset is the item set catalog (DD-277): named groups of gear
// whose worn pieces add bonuses at piece-count thresholds. Content
// (content/sets.yaml) installs it on every load, like the progression
// config; readers take Current().
package itemset

import (
	"fmt"
	"sort"
	"sync/atomic"
)

// Bonus is what wearing at least Pieces items of a set adds.
type Bonus struct {
	Pieces    int        `yaml:"pieces"    json:"pieces"`
	Modifiers []Modifier `yaml:"modifiers" json:"modifiers"`
}

// Modifier mirrors skill.Modifier (kind, target, amount) without the import.
type Modifier struct {
	Kind   string  `yaml:"kind"   json:"kind"`
	Target string  `yaml:"target" json:"target"`
	Amount float64 `yaml:"amount" json:"amount"`
}

// Set is one item set. Source is the raid (dungeon template code) whose boss
// drops its pieces.
type Set struct {
	Code           string   `yaml:"code"             json:"code"`
	DisplayNameKey string   `yaml:"display_name_key" json:"display_name_key"`
	Source         string   `yaml:"source"           json:"source"`
	Pieces         []string `yaml:"pieces"           json:"pieces"`
	Bonuses        []Bonus  `yaml:"bonuses"          json:"bonuses"`
}

// Catalog indexes the sets by code and by piece.
type Catalog struct {
	sets    []Set
	byPiece map[string]int
}

var current atomic.Pointer[Catalog]

// Current is the installed catalog; empty before content loads.
func Current() *Catalog {
	if c := current.Load(); c != nil {
		return c
	}
	return &Catalog{}
}

// New validates sets: unique codes, pieces in one set only, bonus
// thresholds ascending within 1..len(pieces).
func New(sets []Set) (*Catalog, error) {
	c := &Catalog{sets: sets, byPiece: map[string]int{}}
	codes := map[string]bool{}
	for i, s := range sets {
		if s.Code == "" || codes[s.Code] {
			return nil, fmt.Errorf("itemset: empty or duplicate code %q", s.Code)
		}
		codes[s.Code] = true
		if len(s.Pieces) < 2 || len(s.Bonuses) == 0 {
			return nil, fmt.Errorf("itemset %s: needs two pieces and a bonus", s.Code)
		}
		for _, p := range s.Pieces {
			if _, dup := c.byPiece[p]; dup {
				return nil, fmt.Errorf("itemset %s: piece %s is in two sets", s.Code, p)
			}
			c.byPiece[p] = i
		}
		last := 0
		for _, b := range s.Bonuses {
			if b.Pieces <= last || b.Pieces > len(s.Pieces) || len(b.Modifiers) == 0 {
				return nil, fmt.Errorf("itemset %s: bonus thresholds must ascend within 1..%d", s.Code, len(s.Pieces))
			}
			last = b.Pieces
		}
	}
	return c, nil
}

// Install validates and swaps in a catalog.
func Install(sets []Set) error {
	c, err := New(sets)
	if err != nil {
		return err
	}
	current.Store(c)
	return nil
}

// All lists the sets in content order.
func (c *Catalog) All() []Set { return c.sets }

// Of is the set a template belongs to.
func (c *Catalog) Of(template string) (Set, bool) {
	i, ok := c.byPiece[template]
	if !ok {
		return Set{}, false
	}
	return c.sets[i], true
}

// FromSource lists the sets a raid's boss drops.
func (c *Catalog) FromSource(source string) []Set {
	var out []Set
	for _, s := range c.sets {
		if s.Source == source {
			out = append(out, s)
		}
	}
	return out
}

// Modifiers are the bonuses the worn templates earn: per set, every bonus
// whose piece threshold the count of distinct worn pieces reaches.
func (c *Catalog) Modifiers(worn []string) []Modifier {
	count := map[int]map[string]bool{}
	for _, t := range worn {
		if i, ok := c.byPiece[t]; ok {
			if count[i] == nil {
				count[i] = map[string]bool{}
			}
			count[i][t] = true
		}
	}
	idx := make([]int, 0, len(count))
	for i := range count {
		idx = append(idx, i)
	}
	sort.Ints(idx) // deterministic order
	var out []Modifier
	for _, i := range idx {
		for _, b := range c.sets[i].Bonuses {
			if len(count[i]) >= b.Pieces {
				out = append(out, b.Modifiers...)
			}
		}
	}
	return out
}
