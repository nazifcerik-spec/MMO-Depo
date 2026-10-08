package contentpack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateStrings(t *testing.T) {
	pack := &Pack{Items: []Item{{Code: "a"}}, Mobs: []Mob{{Code: "m"}}}
	ok := func() *Strings {
		return &Strings{ByLang: map[string]map[string]string{
			"en": {"item.a": "A", "mob.m": "M"},
			"tr": {"item.a": "A tr", "mob.m": "M tr"},
		}}
	}
	require.NoError(t, validateStrings(ok(), pack))
	require.NoError(t, validateStrings(&Strings{ByLang: map[string]map[string]string{}}, pack), "no i18n dir: nothing to check")

	for name, mut := range map[string]func(*Strings){
		"no default":   func(s *Strings) { delete(s.ByLang, "en") },
		"missing key":  func(s *Strings) { delete(s.ByLang["tr"], "mob.m") },
		"extra key":    func(s *Strings) { s.ByLang["tr"]["item.zz"] = "x" },
		"empty name":   func(s *Strings) { s.ByLang["tr"]["item.a"] = " " },
		"code no name": func(s *Strings) { delete(s.ByLang["en"], "item.a"); delete(s.ByLang["tr"], "item.a") },
	} {
		s := ok()
		mut(s)
		require.Error(t, validateStrings(s, pack), name)
	}
}

func TestStringsLookupAndLoad(t *testing.T) {
	var none *Strings
	require.Nil(t, none.Lookup("tr"))
	s := &Strings{ByLang: map[string]map[string]string{"en": {"k": "en"}, "tr": {"k": "tr"}}}
	require.Equal(t, "tr", s.Lookup("TR")["k"])
	require.Equal(t, "en", s.Lookup("de")["k"], "unknown language falls back to English")

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "en.yaml"), []byte("k: v\n"), 0o600))
	a, err := loadStrings(dir)
	require.NoError(t, err)
	require.Equal(t, "v", a.ByLang["en"]["k"])
	require.NoError(t, os.WriteFile(filepath.Join(dir, "en.yaml"), []byte("k: w\n"), 0o600))
	b, err := loadStrings(dir)
	require.NoError(t, err)
	require.NotEqual(t, a.Version, b.Version, "version follows the files")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "tr.yaml"), []byte("[nope"), 0o600))
	_, err = loadStrings(dir)
	require.Error(t, err)
}
