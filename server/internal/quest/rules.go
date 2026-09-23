// Package quest implements quests v1 (Sprint 41, DD-221): main / side /
// daily templates from the content pack, objectives counted from the same
// progression event stream that feeds mastery and achievements, and
// XP + gold rewards on turn-in. docs/quests/*.
package quest

import (
	"encoding/binary"
	"hash/fnv"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/progression"
)

const (
	CategoryMain  = "main"
	CategorySide  = "side"
	CategoryDaily = "daily"

	// DailySlots is how many dailies each character is offered per UTC day.
	DailySlots = 3
	// OpenCap bounds concurrently open quests (docs/quests/README.md §5).
	OpenCap = 25
)

// Objective counts one progression event kind ("mob_killed:grey_wolf").
type Objective struct {
	Event string `json:"event"`
	Count int64  `json:"count"`
}

// Template is one quest definition.
type Template struct {
	Code       string
	NameKey    string
	Category   string
	MinLevel   int16
	Requires   string // quest that must be turned in first; "" = none
	Objectives []Objective
	RewardXP   int64
	RewardGold int64
	SortOrder  int32
}

// Catalog is the immutable quest content snapshot (swapped on reload).
type Catalog struct {
	list   []Template
	byCode map[string]*Template
}

// NewCatalog indexes templates, ordered by sort_order then code.
func NewCatalog(list []Template) *Catalog {
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].SortOrder != list[j].SortOrder {
			return list[i].SortOrder < list[j].SortOrder
		}
		return list[i].Code < list[j].Code
	})
	c := &Catalog{list: list, byCode: make(map[string]*Template, len(list))}
	for i := range c.list {
		c.byCode[c.list[i].Code] = &c.list[i]
	}
	return c
}

func (c *Catalog) Get(code string) (*Template, bool) {
	t, ok := c.byCode[code]
	return t, ok
}

// DayKey is the UTC calendar day a daily belongs to.
func DayKey(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// DailyOffer picks the character's dailies for `day`: the DailySlots
// level-eligible dailies with the lowest hash(character, day, code). Pure and
// stable, so no per-day rows or scheduler are needed.
func (c *Catalog) DailyOffer(characterID uuid.UUID, level int16, day time.Time) []string {
	type scored struct {
		code  string
		score uint64
	}
	var pool []scored
	dayBytes := binary.BigEndian.AppendUint64(nil, uint64(DayKey(day).Unix()))
	for _, t := range c.list {
		if t.Category != CategoryDaily || t.MinLevel > level {
			continue
		}
		h := fnv.New64a()
		_, _ = h.Write(characterID[:])
		_, _ = h.Write(dayBytes)
		_, _ = h.Write([]byte(t.Code))
		pool = append(pool, scored{t.Code, h.Sum64()})
	}
	sort.Slice(pool, func(i, j int) bool { return pool[i].score < pool[j].score })
	out := make([]string, 0, DailySlots)
	for i := 0; i < len(pool) && i < DailySlots; i++ {
		out = append(out, pool[i].code)
	}
	sort.Strings(out)
	return out
}

// Advance applies events to a progress vector (parallel to t.Objectives).
// Each objective counts only the event kind it names, capped at its count.
func Advance(t *Template, progress []int64, events []progression.Event) (next []int64, changed, done bool) {
	next = make([]int64, len(t.Objectives))
	copy(next, progress)
	done = true
	for i, o := range t.Objectives {
		for _, e := range events {
			if e.Kind == o.Event && e.Amount > 0 && next[i] < o.Count {
				next[i] = min(o.Count, next[i]+e.Amount)
				changed = true
			}
		}
		if next[i] < o.Count {
			done = false
		}
	}
	return next, changed, done
}
