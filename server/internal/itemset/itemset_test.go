package itemset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testSets() []Set {
	hp := Modifier{Kind: "mult_add", Target: "hp_max", Amount: 0.04}
	atk := Modifier{Kind: "mult_add", Target: "attack", Amount: 0.05}
	return []Set{
		{Code: "keep_plate", Source: "keep", Pieces: []string{"a", "b", "c", "d"},
			Bonuses: []Bonus{{Pieces: 2, Modifiers: []Modifier{hp}}, {Pieces: 4, Modifiers: []Modifier{atk}}}},
		{Code: "keep_cloth", Source: "keep", Pieces: []string{"x", "y"}, Bonuses: []Bonus{{Pieces: 2, Modifiers: []Modifier{atk}}}},
		{Code: "tower_plate", Source: "tower", Pieces: []string{"t1", "t2"}, Bonuses: []Bonus{{Pieces: 2, Modifiers: []Modifier{hp}}}},
	}
}

// Worn pieces earn every bonus whose threshold their distinct count reaches
// (DD-277).
func TestModifiers(t *testing.T) {
	c, err := New(testSets())
	require.NoError(t, err)
	require.Empty(t, c.Modifiers([]string{"a", "x", "t1", "sword"}), "one piece of each: nothing")
	require.Len(t, c.Modifiers([]string{"a", "b"}), 1)
	require.Len(t, c.Modifiers([]string{"a", "a", "b"}), 1, "a duplicate is one piece")
	got := c.Modifiers([]string{"d", "c", "b", "a", "x", "y"})
	require.Equal(t, []Modifier{{"mult_add", "hp_max", 0.04}, {"mult_add", "attack", 0.05}, {"mult_add", "attack", 0.05}}, got)

	s, ok := c.Of("c")
	require.True(t, ok)
	require.Equal(t, "keep_plate", s.Code)
	_, ok = c.Of("sword")
	require.False(t, ok)
	require.Len(t, c.FromSource("keep"), 2)
	require.Empty(t, c.FromSource("nowhere"))
}

func TestNewRejectsBadSets(t *testing.T) {
	for name, mut := range map[string]func([]Set){
		"duplicate code":   func(s []Set) { s[1].Code = s[0].Code },
		"piece in two":     func(s []Set) { s[1].Pieces[0] = "a" },
		"one piece":        func(s []Set) { s[1].Pieces = s[1].Pieces[:1] },
		"no bonus":         func(s []Set) { s[1].Bonuses = nil },
		"threshold order":  func(s []Set) { s[0].Bonuses[0].Pieces = 4 },
		"threshold > size": func(s []Set) { s[1].Bonuses[0].Pieces = 3 },
		"empty bonus":      func(s []Set) { s[1].Bonuses[0].Modifiers = nil },
	} {
		sets := testSets()
		mut(sets)
		_, err := New(sets)
		require.Error(t, err, name)
	}
	require.NoError(t, Install(testSets()))
	require.Len(t, Current().All(), 3)
	require.NoError(t, Install(nil))
	require.Empty(t, Current().All())
}
