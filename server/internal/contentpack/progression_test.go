package contentpack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-250: progression strings exist in en, tr, zh-CN and es; partial
// languages fall back to English for everything else.
func TestProgressionLocalization(t *testing.T) {
	p, err := Read("../../../content")
	require.NoError(t, err)
	keys := progressionKeys(p)
	require.Contains(t, keys, "class.iron_fortress")
	require.Contains(t, keys, "skill.warrior_execute.desc")
	require.Contains(t, keys, "talent.warrior_blood_frenzy.desc")
	require.Contains(t, keys, "title.eternal")
	require.Contains(t, keys, "attribute.wisdom")
	for _, lang := range []string{"en", "tr", "zh-cn", "es"} {
		table := p.Strings.ByLang[lang]
		require.NotNil(t, table, lang)
		for _, k := range keys {
			require.NotEmpty(t, table[k], "%s: %s", lang, k)
		}
	}
	require.Equal(t, "战士", p.Strings.Lookup("zh")["class.warrior"], "zh finds zh-CN")
	require.Equal(t, "Guerrero", p.Strings.Lookup("es-ES")["class.warrior"])
	require.Equal(t, p.Strings.Lookup("en")["item.iron_ore"], p.Strings.Lookup("zh")["item.iron_ore"], "fallback to en")
	require.Equal(t, "Savaşçı", p.Strings.Lookup("tr")["class.warrior"])

	// A partial language missing a progression key is refused.
	delete(p.Strings.ByLang["es"], "class.guardian")
	require.ErrorContains(t, validateStrings(p.Strings, p), "missing progression key class.guardian")
}
