package profanity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCensor(t *testing.T) {
	for in, want := range map[string]struct {
		out string
		n   int
	}{
		"merhaba dünya":             {"merhaba dünya", 0},
		"siktir git":                {"****** git", 1},
		"S1KT1RRR lan":              {"******** lan", 1},
		"amk ne oldu":               {"*** ne oldu", 1},
		"orospuçocuğu, bitch!":      {"************, ******", 2},
		"sıkıştım, sıkıntı var":     {"sıkıştım, sıkıntı var", 0},
		"I got a pic of the castle": {"I got a pic of the castle", 0},
		"ne piç adam":               {"ne *** adam", 1},
		"what the fuck you dick":    {"what the **** you ****", 2},
		"cockpit amazing":           {"cockpit amazing", 0},
		"100 altın":                 {"100 altın", 0},
	} {
		out, n := Censor(in)
		require.Equal(t, want.out, out, in)
		require.Equal(t, want.n, n, in)
	}
}
