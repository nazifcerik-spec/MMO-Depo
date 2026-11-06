package character

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-354: names in the game's seven languages.
func TestValidateNameAcceptsEveryLanguage(t *testing.T) {
	for _, name := range []string{
		"Grim", "Hero_42", "Şule", "Çağrı", "İpek", "Işıl", "Gül", // Turkish
		"Niño", "José", // Spanish
		"Иван", "Наташа", // Russian
		"민수", "김철수", // Korean
		"李明", "王小明", // Chinese
		"さくら", "サクラ", "山田太郎", "コーヒー", // Japanese
	} {
		display, _, err := validateName(name)
		require.NoError(t, err, name)
		require.Equal(t, name, display)
	}
}

func TestValidateNameRefuses(t *testing.T) {
	for _, bad := range []string{
		"", " ", "ab", "明", // too short
		"abcdefghijklmnopq",    // 17
		"1abc", "_abc", "1234", // must start with a letter
		"a b", "a-b", "a.b", "a@b", "<b>x", "a\u200bb", // spaces, punctuation, invisible
		"Аli", "Иvan", "민su", "山田Taro", // mixed alphabets
		"Ａｂｃ", "ｻｸﾗ", // full and half width
		"ㅋㅋㅋ",   // bare Hangul letters
		"Νίκος", // Greek is not one of the game's alphabets
		"😀😀😀",   // emoji
		"q́bc",  // combining mark that does not compose
	} {
		_, _, err := validateName(bad)
		require.Error(t, err, "%q", bad)
	}
}

func TestValidateNameComposes(t *testing.T) {
	display, key, err := validateName("José") // e + combining acute
	require.NoError(t, err)
	require.Equal(t, "José", display)
	require.Equal(t, "josé", key)
}

func TestNameKeyFoldsAcrossLocales(t *testing.T) {
	for _, same := range [][]string{
		{"ipek", "İPEK", "IPEK", "ıpek", " Ipek "},
		{"ali", "ALI", "Alı", "ALİ"},
		{"şule", "ŞULE"},
		{"иван", "ИВАН"},
		{"sasha", "ѕasha"},
		{"opa", "ора"}, // Cyrillic о, р, а
	} {
		for _, s := range same {
			require.Equal(t, NameKey(same[0]), NameKey(s), s)
		}
	}
	require.NotEqual(t, NameKey("gul"), NameKey("gül"))
	// Old ASCII keys were lowercase: they stay valid without a migration.
	require.Equal(t, "grim_42", NameKey("Grim_42"))
}
