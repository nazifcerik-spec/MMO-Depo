//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/domain/character"
)

// XP rollover on the real rows: multi-level leaps, the cap, bad input.
func TestXPRollover(t *testing.T) {
	w := newWorld(t)
	p := w.player("warrior", 1)

	chars := character.NewPgRepo(w.p)
	_, _, err := chars.AddXP(w.ctx, p.ID, -1)
	require.ErrorIs(t, err, character.ErrInvalidXP)
	_, _, err = chars.AddXP(w.ctx, uuid.New(), 1)
	require.ErrorIs(t, err, character.ErrNotFound)
	c, up, err := chars.AddXP(w.ctx, p.ID, 1<<40)
	require.NoError(t, err)
	require.True(t, up)
	require.Equal(t, int16(character.LevelCap), c.Level, "a huge grant stops at the cap")
	again, up, err := chars.AddXP(w.ctx, p.ID, 10)
	require.NoError(t, err)
	require.False(t, up)
	require.Equal(t, c.XP+10, again.XP, "XP past the cap is kept")

	profs := crafting.NewPgRepo(w.p)
	_, _, err = profs.AddXP(w.ctx, p.ID, "alchemy", 5, 60)
	require.ErrorIs(t, err, crafting.ErrProfessionNotLearned)
	_, err = profs.UpsertLearn(w.ctx, p.ID, "alchemy", time.Now())
	require.NoError(t, err)
	_, _, err = profs.AddXP(w.ctx, p.ID, "alchemy", -1, 60)
	require.ErrorIs(t, err, crafting.ErrValidation)
	cp, up, err := profs.AddXP(w.ctx, p.ID, "alchemy", 100+200+50, 60)
	require.NoError(t, err)
	require.True(t, up)
	require.Equal(t, int16(3), cp.Level, "two levels in one grant")
	require.Equal(t, int32(50), cp.XP)
	cp, _, err = profs.AddXP(w.ctx, p.ID, "alchemy", 1<<20, 5)
	require.NoError(t, err)
	require.Equal(t, int16(5), cp.Level, "profession cap")
}
