package dungeon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
)

// DD-218: a dungeon run starts at the character's current HP, not a fresh max pool.
func TestCombatantFromSheet_UsesCurrentHP(t *testing.T) {
	s := &character.Sheet{Resources: character.Resources{HPCurrent: 37, HPMax: 250}}
	c := combatantFromSheet(s)
	require.Equal(t, int32(37), c.HPCurrent)
	require.Equal(t, int32(250), c.HPMax)
}
