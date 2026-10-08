package item

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Every affix the roller can produce has a display name in every language
// (content/i18n, DD-230).
func TestAffixesHaveNames(t *testing.T) {
	for _, lang := range []string{"en", "tr"} {
		raw, err := os.ReadFile("../../../content/i18n/" + lang + ".yaml")
		require.NoError(t, err)
		var names map[string]string
		require.NoError(t, yaml.Unmarshal(raw, &names))
		for _, pools := range []map[string][]string{prefixPool, suffixPool} {
			for _, codes := range pools {
				for _, c := range codes {
					require.NotEmpty(t, names["affix."+c], "%s: affix.%s", lang, c)
				}
			}
		}
		require.Contains(t, names["affix.format"], "{name}", lang)
	}
}
