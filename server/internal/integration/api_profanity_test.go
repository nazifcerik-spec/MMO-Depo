//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Chat profanity (DD-312): swear words are masked; seven within a day mute
// the account for five minutes and global chat hears about it.
func TestHTTPChatProfanity(t *testing.T) {
	newWorld(t)
	p := newHTTPPlayer(t, "rogue")
	msg := p.ok("POST", p.path("/chat"), map[string]any{"channel": "global", "body": "siktir git amk"})
	require.Equal(t, "****** git ***", msg["body"])

	p.ok("POST", p.path("/chat"), map[string]any{"channel": "global", "body": "fuck shit bitch"})       // 5
	p.ok("POST", p.path("/chat"), map[string]any{"channel": "global", "body": "orospu ve pezevenk"}) // 7: muted
	p.fail("POST", p.path("/chat"), map[string]any{"channel": "global", "body": "merhaba"}, 403, "E_MUTED")

	hist := p.ok("GET", p.path("/chat/global"), nil)["messages"].([]any)
	last := hist[len(hist)-1].(map[string]any)
	require.Equal(t, "profanity_mute", last["system"])
	require.Equal(t, p.name, last["from_name"])
}
