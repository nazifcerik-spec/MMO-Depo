package combat

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/domain/character"
)

// DD-218: hunts and dungeon runs start at the character's current HP, not a
// fresh max pool; the level travels along for the level-gap rule (DD-229).
func TestFromSheet(t *testing.T) {
	s := &character.Sheet{Level: 7, Primary: character.PrimaryStats{Agility: 100},
		Resources: character.Resources{HPCurrent: 37, HPMax: 250}}
	c := FromSheet(s)
	require.Equal(t, int32(37), c.HPCurrent)
	require.Equal(t, int32(250), c.HPMax)
	require.Equal(t, int16(7), c.Level)
	require.Equal(t, int32(40), c.EvasionPct, "evasion is capped")
}

func TestLevelGap(t *testing.T) {
	p := Combatant{Level: 10, HitPct: 90, Defense: 100}
	applyLevelGap(&p, 10)
	require.Equal(t, Combatant{Level: 10, HitPct: 90, Defense: 100}, p, "same level: no change")
	applyLevelGap(&p, 13)
	require.Equal(t, int32(90-3*GapHitPerLevel), p.HitPct)
	require.Equal(t, int32(100-3*GapDefensePerLevel), p.Defense)

	old := Combatant{HitPct: 90, Defense: 100} // pre-DD-229 snapshot: no level
	applyLevelGap(&old, 30)
	require.Equal(t, int32(90), old.HitPct)
	big := Combatant{Level: 1, Defense: 100}
	applyLevelGap(&big, 40)
	require.Equal(t, int32(0), big.Defense, "defense never goes negative")
}
