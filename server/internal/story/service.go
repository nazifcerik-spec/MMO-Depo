package story

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
)

var (
	ErrNotCurrent  = errors.New("story: not the current scene")
	ErrLevel       = errors.New("story: level too low")
	ErrQuests      = errors.New("story: story quests not turned in")
	ErrOption      = errors.New("story: no such option")
	ErrLocked      = errors.New("story: option needs a quest first")
	ErrChosen      = errors.New("story: already answered differently")
	ErrNotPlayable = errors.New("story: character not playable")
)

// CharacterFetcher is the slice of character.Service the story needs.
type CharacterFetcher interface {
	GetOwned(ctx context.Context, characterID, accountID uuid.UUID) (*character.Character, error)
}

// ReputationGranter adds faction reputation in the caller's tx
// (progression.Service).
type ReputationGranter interface {
	GrantReputation(ctx context.Context, tx pgx.Tx, char uuid.UUID, faction string, delta int) error
}

type Service struct {
	pool *pgxpool.Pool
	char CharacterFetcher
	rep  ReputationGranter
}

func NewService(pool *pgxpool.Pool, chars CharacterFetcher, rep ReputationGranter) *Service {
	return &Service{pool: pool, char: chars, rep: rep}
}

// --- views -------------------------------------------------------------------

// Quest states of a scene's quest: "" (not taken), active, completed, turned_in.
type QuestView struct {
	Code     string `json:"code"`
	Optional bool   `json:"optional,omitempty"`
	State    string `json:"state"`
}

type OptionView struct {
	Code string `json:"code"`
	// RequiresQuest is set while the quest it needs is not turned in.
	RequiresQuest string `json:"requires_quest,omitempty"`
}

type SceneView struct {
	Code     string `json:"code"`
	Chapter  string `json:"chapter"`
	MinLevel int    `json:"min_level"`
	// State: done | current | locked (current, level too low).
	State       string       `json:"state"`
	Chosen      string       `json:"chosen,omitempty"`
	Extras      []string     `json:"extras"`
	AfterExtras []string     `json:"after_extras"`
	After       bool         `json:"after"`
	Quests      []QuestView  `json:"quests"`
	Options     []OptionView `json:"options"`
	// Ready: the required quests are turned in (the scene can be passed).
	Ready bool `json:"ready"`
}

type ChapterView struct {
	Code     string `json:"code"`
	MinLevel int    `json:"min_level"`
	State    string `json:"state"` // done | open | locked
}

type Journal struct {
	Chapters    []ChapterView `json:"chapters"`
	Scenes      []SceneView   `json:"scenes"`
	State       State         `json:"state"`
	Finished    bool          `json:"finished"`
	TextVersion string        `json:"text_version"`
}

// q is satisfied by the pool and a tx.
type q interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func doneScenes(ctx context.Context, db q, char uuid.UUID) (map[string]string, error) {
	rows, err := db.Query(ctx, `SELECT scene_code, option_code FROM progression.story_scene WHERE character_id = $1`, char)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var sc, opt string
		if err := rows.Scan(&sc, &opt); err != nil {
			return nil, err
		}
		out[sc] = opt
	}
	return out, rows.Err()
}

var questRank = map[string]int{"active": 1, "completed": 2, "turned_in": 3}

// questStates is the best state of each story quest the character took
// (abandoned / expired count as not taken).
func questStates(ctx context.Context, db q, char uuid.UUID, codes []string) (map[string]string, error) {
	out := map[string]string{}
	if len(codes) == 0 {
		return out, nil
	}
	rows, err := db.Query(ctx, `SELECT quest_code, state FROM progression.character_quest
		WHERE character_id = $1 AND quest_code = ANY($2)`, char, codes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var code, st string
		if err := rows.Scan(&code, &st); err != nil {
			return nil, err
		}
		if questRank[st] > questRank[out[code]] {
			out[code] = st
		}
	}
	return out, rows.Err()
}

func sceneQuests(sc *Scene) []string { return append(slices.Clone(sc.Quests), sc.OptionalQuests...) }

func extras(sc *Scene, st State, race string) (before, after []string) {
	before, after = []string{}, []string{}
	for _, e := range sc.Extras {
		if !e.When.Match(st) || (len(e.Race) > 0 && !slices.Contains(e.Race, race)) {
			continue
		}
		if e.After {
			after = append(after, e.Key)
		} else {
			before = append(before, e.Key)
		}
	}
	return before, after
}

func ready(sc *Scene, qs map[string]string) bool {
	for _, code := range sc.Quests {
		if qs[code] != "turned_in" {
			return false
		}
	}
	return true
}

// build is the journal of a character of level and race with these passed
// scenes and quest states.
func build(in *installed, done, qs map[string]string, level int, race string) *Journal {
	j := &Journal{Chapters: []ChapterView{}, Scenes: []SceneView{}}
	if in.text != nil {
		j.TextVersion = in.text.Version
	}
	view := func(sc *Scene, st State, state string) SceneView {
		v := SceneView{Code: sc.Code, Chapter: sc.Chapter, MinLevel: sc.MinLevel, State: state, After: sc.After,
			Quests: []QuestView{}, Options: []OptionView{}, Ready: ready(sc, qs)}
		v.Extras, v.AfterExtras = extras(sc, st, race)
		for _, code := range sc.Quests {
			v.Quests = append(v.Quests, QuestView{Code: code, State: qs[code]})
		}
		for _, code := range sc.OptionalQuests {
			v.Quests = append(v.Quests, QuestView{Code: code, Optional: true, State: qs[code]})
		}
		return v
	}
	st, cur := in.book.walk(done, func(sc *Scene, before State) {
		v := view(sc, before, "done")
		v.Chosen = done[sc.Code]
		j.Scenes = append(j.Scenes, v)
	})
	j.State = st
	curChapter := len(in.book.Chapters)
	if cur != nil {
		state := "current"
		if level < cur.MinLevel {
			state = "locked"
		}
		v := view(cur, st, state)
		for _, o := range cur.Visible(st) {
			ov := OptionView{Code: o.Code}
			if o.RequiresQuest != "" && qs[o.RequiresQuest] != "turned_in" {
				ov.RequiresQuest = o.RequiresQuest
			}
			v.Options = append(v.Options, ov)
		}
		j.Scenes = append(j.Scenes, v)
		curChapter = slices.IndexFunc(in.book.Chapters, func(c Chapter) bool { return c.Code == cur.Chapter })
	} else {
		j.Finished = len(in.book.Scenes) > 0
	}
	for i, c := range in.book.Chapters {
		state := "locked"
		switch {
		case i < curChapter:
			state = "done"
		case i == curChapter:
			state = "open"
		}
		j.Chapters = append(j.Chapters, ChapterView{Code: c.Code, MinLevel: c.MinLevel, State: state})
	}
	return j
}

func allQuests(in *installed) []string {
	out := make([]string, 0, len(in.quests))
	for code := range in.quests {
		out = append(out, code)
	}
	return out
}

// Journal is the character's story so far: passed scenes, the current one
// and the replayed flags.
func (s *Service) Journal(ctx context.Context, char, account uuid.UUID) (*Journal, error) {
	c, err := s.char.GetOwned(ctx, char, account)
	if err != nil {
		return nil, err
	}
	in := load()
	done, err := doneScenes(ctx, s.pool, char)
	if err != nil {
		return nil, err
	}
	qs, err := questStates(ctx, s.pool, char, allQuests(in))
	if err != nil {
		return nil, err
	}
	return build(in, done, qs, int(c.Level), c.RaceCode), nil
}

// Advance passes the current scene with option ("" when it has no visible
// option). Replaying the same answer is a no-op; a different one is refused.
func (s *Service) Advance(ctx context.Context, char, account uuid.UUID, scene, option string) (*Journal, error) {
	c, err := s.char.GetOwned(ctx, char, account)
	if err != nil {
		return nil, err
	}
	if c.State != character.StateActive {
		return nil, ErrNotPlayable
	}
	in := load()
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		// One story step at a time per character.
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1::text, 367))`, char); err != nil {
			return err
		}
		done, err := doneScenes(ctx, tx, char)
		if err != nil {
			return err
		}
		if prev, ok := done[scene]; ok {
			if prev == option {
				return nil
			}
			return ErrChosen
		}
		st, cur := in.book.walk(done, func(*Scene, State) {})
		if cur == nil || cur.Code != scene {
			return ErrNotCurrent
		}
		if int(c.Level) < cur.MinLevel {
			return ErrLevel
		}
		qs, err := questStates(ctx, tx, char, sceneQuests(cur))
		if err != nil {
			return err
		}
		// An option may need a quest of an earlier scene too.
		var opt *Option
		visible := cur.Visible(st)
		switch {
		case len(visible) == 0 && option != "":
			return ErrOption
		case len(visible) > 0:
			i := slices.IndexFunc(visible, func(o *Option) bool { return o.Code == option })
			if i < 0 {
				return ErrOption
			}
			opt = visible[i]
		}
		if !ready(cur, qs) {
			return ErrQuests
		}
		if opt != nil && opt.RequiresQuest != "" {
			need, err := questStates(ctx, tx, char, []string{opt.RequiresQuest})
			if err != nil {
				return err
			}
			if need[opt.RequiresQuest] != "turned_in" {
				return ErrLocked
			}
		}
		if _, err := tx.Exec(ctx, `INSERT INTO progression.story_scene (character_id, scene_code, option_code)
			VALUES ($1, $2, $3)`, char, scene, option); err != nil {
			return err
		}
		if opt != nil && s.rep != nil {
			for faction, n := range opt.Reputation {
				if err := s.rep.GrantReputation(ctx, tx, char, faction, n); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.Journal(ctx, char, account)
}

// StoryQuests is the quest gate (quest.StoryGate): the story quests the
// character's current scene offers, if its level reaches the scene.
func (s *Service) StoryQuests(ctx context.Context, char uuid.UUID, level int16) (map[string]bool, error) {
	in := load()
	out := map[string]bool{}
	if len(in.quests) == 0 {
		return out, nil
	}
	done, err := doneScenes(ctx, s.pool, char)
	if err != nil {
		return nil, err
	}
	_, cur := in.book.walk(done, func(*Scene, State) {})
	if cur == nil || int(level) < cur.MinLevel {
		return out, nil
	}
	for _, code := range sceneQuests(cur) {
		out[code] = true
	}
	return out, nil
}
