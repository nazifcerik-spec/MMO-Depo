//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/crafting"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/leveling"
)

// XP rollover on the real rows: multi-level leaps, the cap, bad input.
func TestXPRollover(t *testing.T) {
	w := newWorld(t)
	p := w.player("warrior", 1)

	chars := character.NewPgRepo(w.p)
	cfg := leveling.Current()
	src := character.XPSource{Kind: "test"}
	_, err := chars.GrantXP(w.ctx, p.ID, -1, src, cfg)
	require.ErrorIs(t, err, character.ErrInvalidXP)
	_, err = chars.GrantXP(w.ctx, uuid.New(), 1, src, cfg)
	require.ErrorIs(t, err, character.ErrNotFound)
	// A final specialization has no gate left and runs straight to the cap.
	m := w.player("mage", 300)
	_, err = w.p.Exec(w.ctx, `UPDATE character.character SET spec_code = 'pyromancer' WHERE id = $1`, m.ID)
	require.NoError(t, err)
	g, err := chars.GrantXP(w.ctx, m.ID, 1<<40, src, cfg)
	require.NoError(t, err)
	require.NotEmpty(t, g.LevelsGained)
	require.Equal(t, int16(character.LevelCap), g.Character.Level, "a huge grant stops at the cap")
	again, err := chars.GrantXP(w.ctx, m.ID, 10, src, cfg)
	require.NoError(t, err)
	require.Empty(t, again.LevelsGained)
	require.Equal(t, g.Character.PostCapXP+10, again.Character.PostCapXP, "XP past the cap is kept")
	// A warrior stops at its first promotion gate instead.
	g, err = chars.GrantXP(w.ctx, p.ID, 1<<40, src, cfg)
	require.NoError(t, err)
	require.Equal(t, int16(100), g.Character.Level)
	require.True(t, g.PromotionAvailable)
	require.Positive(t, g.Character.OverflowXP)

	profs := crafting.NewPgRepo(w.p)
	_, _, err = profs.AddXP(w.ctx, p.ID, "alchemy", 5, 60)
	require.ErrorIs(t, err, crafting.ErrProfessionNotLearned)
	_, err = profs.UpsertLearn(w.ctx, p.ID, "alchemy", time.Now())
	require.NoError(t, err)
	_, _, err = profs.AddXP(w.ctx, p.ID, "alchemy", -1, 60)
	require.ErrorIs(t, err, crafting.ErrValidation)
	cp, up, err := profs.AddXP(w.ctx, p.ID, "alchemy", 76+81+50, 60) // levels 1→3 (DD-309) + 50
	require.NoError(t, err)
	require.True(t, up)
	require.Equal(t, int16(3), cp.Level, "two levels in one grant")
	require.Equal(t, int32(50), cp.XP)
	cp, _, err = profs.AddXP(w.ctx, p.ID, "alchemy", 1<<20, 5)
	require.NoError(t, err)
	require.Equal(t, int16(5), cp.Level, "profession cap")
}
