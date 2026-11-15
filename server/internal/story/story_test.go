package story

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func testBook() *Book {
	return &Book{
		Chapters: []Chapter{{Code: "pro", MinLevel: 1}, {Code: "ch1", MinLevel: 5, Requires: "pro"}},
		Scenes: []Scene{
			{Code: "s1", Chapter: "pro", MinLevel: 1, Extras: []Extra{{Key: "dwarf", Race: []string{"dwarf"}}},
				Options: []Option{{Code: "a", Set: map[string]string{"f": "x"}, Add: map[string]int{"m": 1}},
					{Code: "b", Set: map[string]string{"f": "y"}}}},
			{Code: "s2", Chapter: "pro", MinLevel: 1, When: Cond{"f": "x"}, Quests: []string{"q1"}, After: true},
			{Code: "s3", Chapter: "ch1", MinLevel: 5, OptionalQuests: []string{"q2"},
				Extras:  []Extra{{Key: "kind", When: Cond{"m": ">=1"}}, {Key: "late", After: true, When: Cond{"f": "!x"}}},
				Options: []Option{{Code: "save", RequiresQuest: "q2"}, {Code: "y_only", When: Cond{"f": "y"}}, {Code: "go"}}},
			{Code: "s4", Chapter: "ch1", MinLevel: 9, Options: []Option{{Code: "hidden", When: Cond{"f": "z"}}}},
		},
	}
}

var testQuests = map[string]bool{"q1": true, "q2": true}

func testText(b *Book) *Text {
	en := map[string]string{}
	for _, k := range b.TextKeys() {
		en[k] = "text " + k
	}
	return &Text{ByLang: map[string]map[string]string{"en": en, "tr": en, "de": {"s1": "Szene"}}, Version: "v1"}
}

func TestCondMatch(t *testing.T) {
	st := State{Flags: map[string]string{"mirel": "yoldas"}, Counters: map[string]int{"merhamet": 8}}
	for cond, want := range map[string]bool{
		"yoldas": true, "serbest|yoldas": true, "serbest": false, "!tutuklu": true, "!yoldas|serbest": false,
	} {
		require.Equal(t, want, Cond{"mirel": cond}.Match(st), cond)
	}
	require.True(t, Cond{"merhamet": ">=8"}.Match(st))
	require.False(t, Cond{"merhamet": "<8"}.Match(st))
	require.True(t, Cond{"yut": "<1"}.Match(st))
	require.True(t, Cond{"pip": "!yok"}.Match(st), "an unset flag is \"\"")
	require.False(t, Cond{"pip": "yok"}.Match(st))
}

func TestValidate(t *testing.T) {
	require.NoError(t, Validate(testBook(), testQuests))
	require.NoError(t, Validate(nil, nil))
	for name, mut := range map[string]func(b *Book) map[string]bool{
		"bad code":         func(b *Book) map[string]bool { b.Scenes[0].Code = "S 1"; return testQuests },
		"dup with chapter": func(b *Book) map[string]bool { b.Scenes[0].Code = "pro"; return testQuests },
		"unknown chapter":  func(b *Book) map[string]bool { b.Scenes[0].Chapter = "x"; return testQuests },
		"chapter order":    func(b *Book) map[string]bool { b.Scenes[3].Chapter = "pro"; return testQuests },
		"below chapter":    func(b *Book) map[string]bool { b.Scenes[2].MinLevel = 3; return testQuests },
		"level drops":      func(b *Book) map[string]bool { b.Scenes[3].MinLevel = 4; return testQuests },
		"not story quest":  func(b *Book) map[string]bool { return map[string]bool{"q1": true} },
		"unused quest":     func(b *Book) map[string]bool { return map[string]bool{"q1": true, "q2": true, "q3": true} },
		"quest twice":      func(b *Book) map[string]bool { b.Scenes[2].OptionalQuests = []string{"q1"}; return testQuests },
		"after no quests":  func(b *Book) map[string]bool { b.Scenes[0].After = true; return testQuests },
		"dup extra": func(b *Book) map[string]bool {
			b.Scenes[0].Extras = append(b.Scenes[0].Extras, Extra{Key: "dwarf"})
			return testQuests
		},
		"option = extra": func(b *Book) map[string]bool { b.Scenes[0].Options[0].Code = "dwarf"; return testQuests },
		"option = body":  func(b *Book) map[string]bool { b.Scenes[0].Options[0].Code = "body"; return testQuests },
		"empty flag":     func(b *Book) map[string]bool { b.Scenes[0].Options[0].Set["f"] = ""; return testQuests },
		"zero counter":   func(b *Book) map[string]bool { b.Scenes[0].Options[0].Add["m"] = 0; return testQuests },
		"bad reputation": func(b *Book) map[string]bool {
			b.Scenes[0].Options[0].Reputation = map[string]int{"x": -1}
			return testQuests
		},
		"bad condition":    func(b *Book) map[string]bool { b.Scenes[1].When = Cond{"f": "!"}; return testQuests },
		"later quest":      func(b *Book) map[string]bool { b.Scenes[0].Options[0].RequiresQuest = "q2"; return testQuests },
		"chapter requires": func(b *Book) map[string]bool { b.Chapters[0].Requires = "ch1"; return testQuests },
	} {
		b := testBook()
		require.Error(t, Validate(b, mut(b)), name)
	}
}

func TestWalkSkipsAndReplays(t *testing.T) {
	b := testBook()
	st, path, cur := b.Walk(nil)
	require.Empty(t, path)
	require.Equal(t, "s1", cur.Code)
	require.Empty(t, st.Flags)

	// Answer b: s2 (needs f=x) is skipped.
	st, path, cur = b.Walk(map[string]string{"s1": "b"})
	require.Len(t, path, 1)
	require.Equal(t, "s3", cur.Code)
	require.Equal(t, "y", st.Flags["f"])
	require.Equal(t, []string{"save", "y_only", "go"}, codes(cur.Visible(st)))

	st, _, cur = b.Walk(map[string]string{"s1": "a", "s2": ""})
	require.Equal(t, "s3", cur.Code)
	require.Equal(t, 1, st.Counters["m"])
	require.Equal(t, []string{"save", "go"}, codes(cur.Visible(st)))

	_, path, cur = b.Walk(map[string]string{"s1": "a", "s2": "", "s3": "go", "s4": ""})
	require.Nil(t, cur)
	require.Len(t, path, 4)
}

func codes(os []*Option) []string {
	var out []string
	for _, o := range os {
		out = append(out, o.Code)
	}
	return out
}

func TestBuildJournal(t *testing.T) {
	b := testBook()
	require.NoError(t, Install(b, testQuests, testText(b)))
	in := load()

	j := build(in, nil, nil, 1, "dwarf")
	require.Equal(t, []string{"open", "locked"}, chapterStates(j))
	require.Len(t, j.Scenes, 1)
	require.Equal(t, "current", j.Scenes[0].State)
	require.Equal(t, []string{"dwarf"}, j.Scenes[0].Extras)
	require.True(t, j.Scenes[0].Ready)
	require.Equal(t, "v1", j.TextVersion)

	j = build(in, map[string]string{"s1": "a"}, map[string]string{"q1": "completed"}, 1, "human")
	require.Equal(t, "done", j.Scenes[0].State)
	require.Equal(t, "a", j.Scenes[0].Chosen)
	require.Empty(t, j.Scenes[0].Extras, "extras follow the character's race")
	cur := j.Scenes[1]
	require.Equal(t, "s2", cur.Code)
	require.False(t, cur.Ready, "a completed quest is not turned in yet")
	require.Equal(t, []QuestView{{Code: "q1", State: "completed"}}, cur.Quests)

	j = build(in, map[string]string{"s1": "a", "s2": ""}, map[string]string{"q1": "turned_in"}, 4, "human")
	cur = j.Scenes[2]
	require.Equal(t, "locked", cur.State, "level 4 < 5")
	require.Equal(t, []string{"done", "open"}, chapterStates(j))
	require.True(t, cur.Ready, "only optional quests")
	require.Equal(t, []string{"kind"}, cur.Extras)
	require.Empty(t, cur.AfterExtras)
	require.Equal(t, []OptionView{{Code: "save", RequiresQuest: "q2"}, {Code: "go"}}, cur.Options)

	j = build(in, map[string]string{"s1": "a", "s2": "", "s3": "go", "s4": ""}, nil, 20, "human")
	require.True(t, j.Finished)
	require.Equal(t, []string{"done", "done"}, chapterStates(j))
}

func chapterStates(j *Journal) []string {
	var out []string
	for _, c := range j.Chapters {
		out = append(out, c.State)
	}
	return out
}

func TestTextValidateAndHandler(t *testing.T) {
	b := testBook()
	require.NoError(t, testText(b).Validate(b))
	require.Error(t, (&Text{}).Validate(b), "no text")
	require.NoError(t, (*Text)(nil).Validate(&Book{}))

	tx := testText(b)
	delete(tx.ByLang["tr"], "s3.save.result")
	tx.ByLang["tr"] = copyMap(tx.ByLang["en"])
	delete(tx.ByLang["tr"], "s3.save.result")
	require.ErrorContains(t, tx.Validate(b), "missing s3.save.result")

	tx = testText(b)
	tx.ByLang["de"]["s9"] = "x"
	require.ErrorContains(t, tx.Validate(b), "not used")

	tx = testText(b)
	tx.ByLang["de"]["s1"] = " "
	require.ErrorContains(t, tx.Validate(b), "empty")

	require.NoError(t, Install(b, testQuests, testText(b)))
	rec := httptest.NewRecorder()
	TextHandler(rec, httptest.NewRequest("GET", "/v1/story/text?lang=de-AT", nil))
	require.Equal(t, 200, rec.Code)
	require.Contains(t, rec.Body.String(), `"s1":"Szene"`, "the language's own text")
	require.Contains(t, rec.Body.String(), `"s1.body":"text s1.body"`, "English fills the rest")
	etag := rec.Header().Get("ETag")
	req := httptest.NewRequest("GET", "/v1/story/text?lang=de-AT", nil)
	req.Header.Set("If-None-Match", etag)
	rec = httptest.NewRecorder()
	TextHandler(rec, req)
	require.Equal(t, 304, rec.Code)

	rec = httptest.NewRecorder()
	CatalogHandler(rec, httptest.NewRequest("GET", "/v1/story", nil))
	require.Contains(t, rec.Body.String(), `"code":"ch1"`)

	require.Error(t, Install(b, map[string]bool{}, testText(b)), "invalid book")
	require.Error(t, Install(b, testQuests, &Text{}), "invalid text")
}

func copyMap(m map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range m {
		out[k] = v
	}
	return out
}
