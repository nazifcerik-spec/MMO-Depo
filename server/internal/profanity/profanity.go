// Package profanity masks swear words in chat (DD-312). Words are matched
// after folding case, Turkish letters, common digit/symbol swaps and
// repeated letters, so "S1KT1RRR" and "siktir" are the same word. Roots
// match words that start with them (Turkish suffixes); short or ambiguous
// ones only match whole words. The dotless ı stays apart from i: "sıkıştı"
// (it got stuck) is not a swear word.
package profanity

import (
	"strings"
	"unicode"
)

// roots match any word starting with them.
var roots = []string{
	// Turkish
	"orospu", "orosbu", "siktir", "sikerim", "sikeyim", "sikik", "sikis", "sikim", "sikini", "sikiyim",
	"amina", "aminakoy", "amcik", "yarrak", "yarak", "pezevenk", "kahpe", "kaltak", "gavat", "ibne",
	"puust", "pust", "serefsiz", "godoş", "godos", "yavsak", "yavşak", "dallama", "dalyarak", "gotveren",
	"gotlek", "sokarim", "sokayim",
	// English
	"fuck", "motherfuck", "shit", "bitch", "bastard", "asshole", "dickhead", "cunt", "whore", "slut",
	"faggot", "nigger", "retard",
}

// words match only on their own.
var words = map[string]bool{
	"amk": true, "aq": true, "amq": true, "mk": true, "sg": true, "oc": true, "sik": true,
	"dick": true, "cock": true, "piss": true, "stfu": true, "fck": true,
}

// exact are whole words matched before folding: their folded forms are
// everyday English ("pic", "got").
var exact = map[string]bool{"piç": true, "göt": true}

var fold = strings.NewReplacer(
	"İ", "i", "ş", "s", "Ş", "s", "ğ", "g", "Ğ", "g", "ü", "u", "Ü", "u", "ö", "o", "Ö", "o",
	"ç", "c", "Ç", "c", "0", "o", "1", "i", "3", "e", "4", "a", "5", "s", "7", "t", "@", "a", "$", "s", "!", "i",
)

// norm folds a word for matching and squeezes repeated letters.
func norm(w string) string {
	w = fold.Replace(strings.ToLower(w))
	var b strings.Builder
	var last rune
	for _, r := range w {
		if !unicode.IsLetter(r) || r == last {
			continue
		}
		b.WriteRune(r)
		last = r
	}
	return b.String()
}

var normRoots, normWords = func() ([]string, map[string]bool) {
	rs := make([]string, len(roots))
	for i, r := range roots {
		rs[i] = norm(r)
	}
	ws := map[string]bool{}
	for w, on := range words {
		if on {
			ws[norm(w)] = true
		}
	}
	return rs, ws
}()

// Bad reports whether one word is a swear word.
func Bad(word string) bool {
	if exact[strings.ToLower(word)] {
		return true
	}
	n := norm(word)
	if n == "" {
		return false
	}
	if normWords[n] {
		return true
	}
	for _, r := range normRoots {
		if strings.HasPrefix(n, r) {
			return true
		}
	}
	return false
}

// Censor masks every swear word with asterisks (one per letter) and
// returns the text and how many words it masked. Spacing and punctuation
// around words stay as they were.
func Censor(text string) (string, int) {
	var b strings.Builder
	n := 0
	word := []rune{}
	flush := func() {
		if len(word) == 0 {
			return
		}
		if Bad(string(word)) {
			b.WriteString(strings.Repeat("*", len(word)))
			n++
		} else {
			b.WriteString(string(word))
		}
		word = word[:0]
	}
	for _, r := range text {
		// Digits and a few symbols stand in for letters inside a word.
		if unicode.IsLetter(r) || unicode.IsDigit(r) || (len(word) > 0 && strings.ContainsRune("@$!", r)) {
			word = append(word, r)
			continue
		}
		flush()
		b.WriteRune(r)
	}
	flush()
	return b.String(), n
}
