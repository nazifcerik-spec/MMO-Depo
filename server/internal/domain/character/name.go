package character

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const (
	nameMaxRunes = 16
	nameMinRunes = 3
	// One CJK or Hangul syllable says as much as a few Latin letters.
	nameMinRunesCJK = 2
)

// Scripts a character name may be written in (DD-354), the game's seven
// languages. A name stays in one: mixing lets "Аli" pass for "Ali".
type nameScript int

const (
	scriptNone nameScript = iota
	scriptLatin
	scriptCyrillic
	scriptCJK // Han, Hiragana, Katakana, Hangul: names combine these
)

func scriptOf(r rune) nameScript {
	switch {
	case unicode.Is(unicode.Latin, r):
		return scriptLatin
	case unicode.Is(unicode.Cyrillic, r):
		return scriptCyrillic
	case unicode.Is(unicode.Han, r), unicode.Is(unicode.Hiragana, r), unicode.Is(unicode.Katakana, r),
		r == 'ー', r >= 0xAC00 && r <= 0xD7A3:
		return scriptCJK
	}
	return scriptNone
}

// validateName checks a character name and returns its display form and its
// lookup key.
func validateName(raw string) (string, string, error) {
	s := norm.NFC.String(strings.TrimSpace(raw))
	if s == "" {
		return "", "", &ValidationError{Field: "name", Message: "required"}
	}
	if !nameAllowed(s) {
		return "", "", &ValidationError{Field: "name", Message: "letters of one alphabet, digits and _; must start with a letter"}
	}
	return s, NameKey(s), nil
}

func nameAllowed(s string) bool {
	if norm.NFKC.String(s) != s { // full-width and other look-alike forms
		return false
	}
	script, n := scriptNone, 0
	for i, r := range s {
		n++
		if r == '_' || (r >= '0' && r <= '9') {
			if i == 0 {
				return false
			}
			continue
		}
		rs := scriptOf(r)
		if rs == scriptNone || (script != scriptNone && rs != script) {
			return false
		}
		script = rs
	}
	if script == scriptNone || n > nameMaxRunes {
		return false
	}
	if script == scriptCJK {
		return n >= nameMinRunesCJK
	}
	return n >= nameMinRunes
}

// cyrillicLookalikes fold to the Latin letter they are drawn like, so two
// names that read the same cannot both be taken.
var cyrillicLookalikes = map[rune]rune{
	'а': 'a', 'е': 'e', 'о': 'o', 'р': 'p', 'с': 'c', 'у': 'y', 'х': 'x', 'і': 'i', 'ј': 'j', 'ѕ': 's',
}

// NameKey is the case-insensitive key names are unique under and looked up
// by. It does not depend on the reader's locale: I, İ, i and ı are one letter.
func NameKey(raw string) string {
	s := norm.NFKC.String(strings.TrimSpace(raw))
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		r = unicode.ToLower(r)
		if r == 'ı' {
			r = 'i'
		}
		if l, ok := cyrillicLookalikes[r]; ok {
			r = l
		}
		b.WriteRune(r)
	}
	return b.String()
}
