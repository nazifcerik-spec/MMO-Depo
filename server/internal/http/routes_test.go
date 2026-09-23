package http

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

// chi's Route mounts a subrouter on its pattern; a sibling r.Get/Post on the
// very same pattern is silently shadowed and answers 404 (the vendor GET
// catalog was unreachable this way). Register such handlers as "/" inside
// the Route instead.
func TestNoLeafRouteShadowedByMount(t *testing.T) {
	src, err := os.ReadFile("server.go")
	require.NoError(t, err)
	leaf := map[string]bool{}
	for _, m := range regexp.MustCompile(`\.(?:Get|Post|Put|Patch|Delete)\("([^"]+)"`).FindAllSubmatch(src, -1) {
		leaf[string(m[1])] = true
	}
	for _, m := range regexp.MustCompile(`\.Route\("([^"]+)"`).FindAllSubmatch(src, -1) {
		require.False(t, leaf[string(m[1])], "route %s is both mounted and a leaf", m[1])
	}
}
