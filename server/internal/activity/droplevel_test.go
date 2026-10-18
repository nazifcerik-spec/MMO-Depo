package activity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// DD-255: loot never comes from above the character's level.
func TestDropLevel(t *testing.T) {
	require.Equal(t, 40, DropLevel(50, 40), "stronger mob: capped at the character")
	require.Equal(t, 30, DropLevel(30, 40), "weaker mob: its own level")
	require.Equal(t, 50, DropLevel(50, 0), "unknown character level (old snapshot): the mob's")
}
