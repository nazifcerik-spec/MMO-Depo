// Package story is the story mode (DD-313, DD-367): chapters made of scenes,
// read in order from content/story.yaml. A scene may ask for story quests
// (category "story", turned in on the quest board) and end in a choice. The
// choices set flags and counters that later scenes, extra paragraphs and
// options test; the state is never stored, it is replayed from the chosen
// options. Titles of chapters come from content/i18n (story.<chapter>,
// story.<chapter>.body); scene text from content/story/<lang>.yaml.
package story

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/otter/server/internal/http/apierr"
)

// Chapter is one story chapter; its scenes say what it holds.
type Chapter struct {
	Code     string `yaml:"code" json:"code"`
	MinLevel int    `yaml:"min_level" json:"min_level"`
	Requires string `yaml:"requires,omitempty" json:"requires,omitempty"`
}

// Cond tests the story state: flag → "a" / "a|b" (one of) / "!a" (not);
// counter → ">=N" or "<N". A missing flag is "".
type Cond map[string]string

// Extra is a paragraph shown only when its condition (and race) holds;
// After puts it below the quests instead of under the scene's text.
type Extra struct {
	Key   string   `yaml:"key"`
	When  Cond     `yaml:"when,omitempty"`
	Race  []string `yaml:"race,omitempty"`
	After bool     `yaml:"after,omitempty"`
}

// Option is one answer of a scene's choice.
type Option struct {
	Code string            `yaml:"code"`
	Set  map[string]string `yaml:"set,omitempty"` // flags it sets
	Add  map[string]int    `yaml:"add,omitempty"` // counters it raises
	When Cond              `yaml:"when,omitempty"`
	// RequiresQuest: a story quest that must be turned in first.
	RequiresQuest string `yaml:"requires_quest,omitempty"`
	// Reputation: faction → reputation the choice grants once.
	Reputation map[string]int `yaml:"reputation,omitempty"`
}

// Scene is one step of the story. Quests must be turned in before it can be
// passed; optional quests are offered alongside (an option may need one).
// With no visible option the scene is passed with "continue".
type Scene struct {
	Code           string   `yaml:"code"`
	Chapter        string   `yaml:"chapter"`
	MinLevel       int      `yaml:"min_level"`
	When           Cond     `yaml:"when,omitempty"`
	Quests         []string `yaml:"quests,omitempty"`
	OptionalQuests []string `yaml:"optional_quests,omitempty"`
	After          bool     `yaml:"after,omitempty"` // has an after-quests text
	Extras         []Extra  `yaml:"extras,omitempty"`
	Options        []Option `yaml:"options,omitempty"`
}

// Book is the whole story.
type Book struct {
	Chapters []Chapter `yaml:"chapters"`
	Scenes   []Scene   `yaml:"scenes"`
}

// State is the replayed story state of one character.
type State struct {
	Flags    map[string]string `json:"flags"`
	Counters map[string]int    `json:"counters"`
}

func newState() State { return State{Flags: map[string]string{}, Counters: map[string]int{}} }

var cmpRe = regexp.MustCompile(`^(>=|<)(-?\d+)$`)

// Match reports whether s satisfies every entry of c.
func (c Cond) Match(s State) bool {
	for k, want := range c {
		if m := cmpRe.FindStringSubmatch(want); m != nil {
			n, _ := strconv.Atoi(m[2])
			if (m[1] == ">=") != (s.Counters[k] >= n) {
				return false
			}
			continue
		}
		neg := strings.HasPrefix(want, "!")
		in := slices.Contains(strings.Split(strings.TrimPrefix(want, "!"), "|"), s.Flags[k])
		if in == neg {
			return false
		}
	}
	return true
}

func (c Cond) valid() error {
	for k, v := range c {
		if k == "" || strings.Trim(v, "!|") == "" {
			return fmt.Errorf("bad condition %q: %q", k, v)
		}
	}
	return nil
}

// apply records a chosen option in s.
func (o *Option) apply(s State) {
	for k, v := range o.Set {
		s.Flags[k] = v
	}
	for k, n := range o.Add {
		s.Counters[k] += n
	}
}

// Option returns the scene's option with code.
func (sc *Scene) Option(code string) *Option {
	for i := range sc.Options {
		if sc.Options[i].Code == code {
			return &sc.Options[i]
		}
	}
	return nil
}

// Visible is the options whose condition holds.
func (sc *Scene) Visible(s State) []*Option {
	var out []*Option
	for i := range sc.Options {
		if sc.Options[i].When.Match(s) {
			out = append(out, &sc.Options[i])
		}
	}
	return out
}

// Walk replays the chosen options (scene → option, "" = continued) in book
// order. It returns the state, the scenes passed and the current scene (nil
// once the story is done). A scene whose condition fails is skipped; a passed
// scene counts even if content changed its condition since.
func (b *Book) Walk(done map[string]string) (State, []*Scene, *Scene) {
	var path []*Scene
	st, cur := b.walk(done, func(sc *Scene, _ State) { path = append(path, sc) })
	return st, path, cur
}

// walk is Walk calling visit with each passed scene and the state it was
// played in (before its option applied).
func (b *Book) walk(done map[string]string, visit func(sc *Scene, before State)) (State, *Scene) {
	st := newState()
	for i := range b.Scenes {
		sc := &b.Scenes[i]
		opt, passed := done[sc.Code]
		if !passed {
			if !sc.When.Match(st) {
				continue
			}
			return st, sc
		}
		visit(sc, st)
		if o := sc.Option(opt); o != nil {
			o.apply(st)
		}
	}
	return st, nil
}

var codeRe = regexp.MustCompile(`^[a-z0-9_]{1,64}$`)

// Validate checks the book against the story quests (code → used): codes
// unique and well formed, chapters and levels in order, every story quest
// used by exactly one scene, options well formed.
func Validate(b *Book, storyQuests map[string]bool) error {
	if b == nil {
		b = &Book{}
	}
	seen := map[string]bool{}
	chIndex := map[string]int{}
	chMin := map[string]int{}
	for i, c := range b.Chapters {
		switch {
		case !codeRe.MatchString(c.Code):
			return fmt.Errorf("story: chapter code %q", c.Code)
		case seen[c.Code]:
			return fmt.Errorf("story %s: duplicate code", c.Code)
		case c.MinLevel < 1 || c.MinLevel > 1000:
			return fmt.Errorf("story %s: min_level must be 1..1000", c.Code)
		case c.Requires != "" && !seen[c.Requires]:
			return fmt.Errorf("story %s: requires %q, not an earlier chapter", c.Code, c.Requires)
		}
		seen[c.Code] = true
		chIndex[c.Code], chMin[c.Code] = i, c.MinLevel
	}
	used := map[string]string{}
	lastCh, lastLevel := -1, 1
	for _, sc := range b.Scenes {
		ci, ok := chIndex[sc.Chapter]
		switch {
		case !codeRe.MatchString(sc.Code):
			return fmt.Errorf("story: scene code %q", sc.Code)
		case seen[sc.Code]:
			return fmt.Errorf("story %s: duplicate code", sc.Code)
		case !ok:
			return fmt.Errorf("story %s: unknown chapter %q", sc.Code, sc.Chapter)
		case ci < lastCh:
			return fmt.Errorf("story %s: chapter %s comes earlier", sc.Code, sc.Chapter)
		case sc.MinLevel < chMin[sc.Chapter] || sc.MinLevel > 1000:
			return fmt.Errorf("story %s: min_level must be %d..1000", sc.Code, chMin[sc.Chapter])
		case sc.MinLevel < lastLevel:
			return fmt.Errorf("story %s: min_level below the scene before", sc.Code)
		}
		if err := sc.When.valid(); err != nil {
			return fmt.Errorf("story %s: %w", sc.Code, err)
		}
		seen[sc.Code] = true
		lastCh, lastLevel = ci, sc.MinLevel
		for _, q := range append(slices.Clone(sc.Quests), sc.OptionalQuests...) {
			if !storyQuests[q] {
				return fmt.Errorf("story %s: %q is not a story quest", sc.Code, q)
			}
			if prev, dup := used[q]; dup {
				return fmt.Errorf("story %s: quest %s already in %s", sc.Code, q, prev)
			}
			used[q] = sc.Code
		}
		if sc.After && len(sc.Quests)+len(sc.OptionalQuests) == 0 {
			return fmt.Errorf("story %s: after text without quests", sc.Code)
		}
		keys := map[string]bool{}
		for _, e := range sc.Extras {
			if !codeRe.MatchString(e.Key) || keys[e.Key] {
				return fmt.Errorf("story %s: extra key %q", sc.Code, e.Key)
			}
			if err := e.When.valid(); err != nil {
				return fmt.Errorf("story %s extra %s: %w", sc.Code, e.Key, err)
			}
			keys[e.Key] = true
		}
		for _, o := range sc.Options {
			if !codeRe.MatchString(o.Code) || keys[o.Code] || o.Code == "body" || o.Code == "after" || o.Code == "x" {
				return fmt.Errorf("story %s: option code %q", sc.Code, o.Code)
			}
			keys[o.Code] = true
			if err := o.When.valid(); err != nil {
				return fmt.Errorf("story %s option %s: %w", sc.Code, o.Code, err)
			}
			for k, v := range o.Set {
				if k == "" || v == "" {
					return fmt.Errorf("story %s option %s: empty flag", sc.Code, o.Code)
				}
			}
			for k, n := range o.Add {
				if k == "" || n == 0 {
					return fmt.Errorf("story %s option %s: empty counter", sc.Code, o.Code)
				}
			}
			for k, n := range o.Reputation {
				if k == "" || n <= 0 {
					return fmt.Errorf("story %s option %s: reputation %q", sc.Code, o.Code, k)
				}
			}
			if o.RequiresQuest != "" {
				if _, ok := used[o.RequiresQuest]; !ok {
					return fmt.Errorf("story %s option %s: requires %q, not a quest of this or an earlier scene",
						sc.Code, o.Code, o.RequiresQuest)
				}
			}
		}
	}
	for q := range storyQuests {
		if _, ok := used[q]; !ok {
			return fmt.Errorf("story: quest %s is in no scene", q)
		}
	}
	return nil
}

// TextKeys are the scene text keys every complete language must have.
func (b *Book) TextKeys() []string {
	var out []string
	if b == nil {
		return nil
	}
	for _, sc := range b.Scenes {
		out = append(out, sc.Code, sc.Code+".body")
		if sc.After {
			out = append(out, sc.Code+".after")
		}
		for _, e := range sc.Extras {
			out = append(out, sc.Code+".x."+e.Key)
		}
		for _, o := range sc.Options {
			out = append(out, sc.Code+"."+o.Code, sc.Code+"."+o.Code+".result")
		}
	}
	return out
}

// installed is one content load: the book, its scene index and text.
type installed struct {
	book    *Book
	byScene map[string]*Scene
	quests  map[string]string // story quest → scene
	text    *Text
}

var current atomic.Pointer[installed]

// Install swaps in the story of a content load (validated by the loader).
func Install(b *Book, storyQuests map[string]bool, text *Text) error {
	if b == nil {
		b = &Book{}
	}
	if err := Validate(b, storyQuests); err != nil {
		return err
	}
	if err := text.Validate(b); err != nil {
		return err
	}
	in := &installed{book: b, byScene: map[string]*Scene{}, quests: map[string]string{}, text: text}
	for i := range b.Scenes {
		sc := &b.Scenes[i]
		in.byScene[sc.Code] = sc
		for _, q := range append(slices.Clone(sc.Quests), sc.OptionalQuests...) {
			in.quests[q] = sc.Code
		}
	}
	current.Store(in)
	return nil
}

func load() *installed {
	if p := current.Load(); p != nil {
		return p
	}
	return &installed{book: &Book{}, byScene: map[string]*Scene{}, quests: map[string]string{}}
}

// All is the installed chapters in order.
func All() []Chapter {
	if ch := load().book.Chapters; ch != nil {
		return ch
	}
	return []Chapter{}
}

// CatalogHandler is GET /v1/story (public, cacheable).
func CatalogHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"chapters": All()})
}
