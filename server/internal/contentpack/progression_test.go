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
	for _, lang := range []string{"en", "tr", "zh-cn", "es", "ru", "ko", "ja"} {
		table := p.Strings.ByLang[lang]
		require.NotNil(t, table, lang)
		for _, k := range keys {
			require.NotEmpty(t, table[k], "%s: %s", lang, k)
		}
	}
	require.Equal(t, "战士", p.Strings.Lookup("zh")["class.warrior"], "zh finds zh-CN")
	require.Equal(t, "Guerrero", p.Strings.Lookup("es-ES")["class.warrior"])
	// DD-361: every language names the items too, in its own word order.
	require.Equal(t, "铁矿", p.Strings.Lookup("zh")["item.iron_ore"])
	require.Equal(t, "Шлем из железа", p.Strings.Lookup("ru")["item.helm_iron"])
	require.Equal(t, "철 투구", p.Strings.Lookup("ko")["item.helm_iron"])
	require.Equal(t, "鉄の兜", p.Strings.Lookup("ja")["item.helm_iron"])
	for lang, table := range p.Strings.ByLang {
		require.Contains(t, table["affix.format"], "{name}", lang)
		require.Len(t, table, len(p.Strings.ByLang["en"]), lang)
	}
	require.Equal(t, "Savaşçı", p.Strings.Lookup("tr")["class.warrior"])

	// A partial language missing a progression key is refused.
	delete(p.Strings.ByLang["es"], "class.guardian")
	require.ErrorContains(t, validateStrings(p.Strings, p), "missing progression key class.guardian")
}
