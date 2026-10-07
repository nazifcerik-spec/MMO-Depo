//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/worldevent"
)

func spawn(t *testing.T, w *testEnv, code string, start, end time.Time) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	require.NoError(t, w.p.QueryRow(w.ctx, `
		INSERT INTO worldevent.spawn (template_code, starts_at, ends_at) VALUES ($1, $2, $3) RETURNING id`,
		code, start, end).Scan(&id))
	return id
}

func TestWorldEventEnlistCloseAndPay(t *testing.T) {
	w := newWorld(t)
	svc := worldevent.NewService(w.p, w.chars, w.txn, func() string { return "ironhold_city" }, nil)
	hero, low := w.player("warrior", 12), w.player("mage", 2)
	// Unique start minutes keep the (template, starts_at) key free across runs.
	base := time.Now().UTC().Truncate(time.Minute).Add(time.Duration(time.Now().UnixNano()%997) * time.Minute)

	_, err := svc.Attend(w.ctx, hero.ID, hero.Account, uuid.New())
	require.ErrorIs(t, err, worldevent.ErrNotFound)
	far := spawn(t, w, "ashen_colossus", base.Add(72*time.Hour), base.Add(73*time.Hour))
	_, err = svc.Attend(w.ctx, hero.ID, hero.Account, far)
	require.ErrorIs(t, err, worldevent.ErrNotOpen, "registration opens 24h before start")
	inv := spawn(t, w, "greenwood_incursion", base.Add(time.Hour), base.Add(3*time.Hour))
	_, err = svc.Attend(w.ctx, hero.ID, hero.Account, inv)
	require.ErrorIs(t, err, worldevent.ErrWrongRegion)

	boss := spawn(t, w, "ashen_colossus", base.Add(30*time.Minute), base.Add(90*time.Minute))
	_, err = svc.Attend(w.ctx, low.ID, low.Account, boss)
	require.ErrorIs(t, err, worldevent.ErrLevelRequirement)
	view, err := svc.Attend(w.ctx, hero.ID, hero.Account, boss)
	require.NoError(t, err)
	require.True(t, view.Enlisted)
	require.NotNil(t, view.MyPower)
	_, err = svc.Attend(w.ctx, hero.ID, hero.Account, boss)
	require.ErrorIs(t, err, worldevent.ErrAlreadyAttending)
	list, err := svc.List(w.ctx, hero.ID, hero.Account)
	require.NoError(t, err)
	require.NotEmpty(t, list)

	// Fast-forward: the event ran for an hour with the hero there from the start.
	_, err = w.p.Exec(w.ctx, `
		UPDATE worldevent.spawn SET starts_at = now() - interval '61 minutes', ends_at = now() - interval '1 minute'
		 WHERE id = $1`, boss)
	require.NoError(t, err)
	_, err = w.p.Exec(w.ctx, `
		UPDATE worldevent.attendance SET registered_at = now() - interval '2 hours' WHERE spawn_id = $1`, boss)
	require.NoError(t, err)
	_, err = worldevent.Tick(w.ctx, w.p, w.txn, nil, nil)
	require.NoError(t, err)

	var state string
	var tier *int16
	var rewarded *time.Time
	require.NoError(t, w.p.QueryRow(w.ctx, `
		SELECT sp.state, a.reward_tier, a.rewarded_at FROM worldevent.spawn sp
		  JOIN worldevent.attendance a ON a.spawn_id = sp.id WHERE sp.id = $1`, boss).Scan(&state, &tier, &rewarded))
	require.Contains(t, []string{"succeeded", "failed"}, state)
	require.NotNil(t, tier)
	require.GreaterOrEqual(t, *tier, int16(1), "an hour of attendance earns a tier")
	require.NotNil(t, rewarded)
	require.Positive(t, w.balance(hero), "tier gold paid")
	tokens, err := w.txn.Balance(w.ctx, economy.AccountAccount, hero.Account, "event_token")
	require.NoError(t, err)
	require.Positive(t, tokens)

	before := w.balance(hero)
	_, err = worldevent.Tick(w.ctx, w.p, w.txn, nil, nil)
	require.NoError(t, err)
	require.Equal(t, before, w.balance(hero), "rewards pay once")
}

func TestSeasonalDailyClaim(t *testing.T) {
	w := newWorld(t)
	svc := worldevent.NewService(w.p, w.chars, w.txn, nil, nil)
	p := w.player("warrior", 1)
	code := "it_season_" + randHex(3)
	_, err := w.p.Exec(w.ctx, `
		INSERT INTO content.seasonal_event (code, display_name_key, starts_at, ends_at, daily_login_tokens)
		VALUES ($1, 'seasonal.it', now() - interval '1 day', now() + interval '1 day', 7)`, code)
	require.NoError(t, err)

	got, err := svc.ClaimDaily(w.ctx, p.Account, code)
	require.NoError(t, err)
	require.Equal(t, 7, got)
	_, err = svc.ClaimDaily(w.ctx, p.Account, code)
	require.ErrorIs(t, err, worldevent.ErrAlreadyClaimed)
	_, err = svc.ClaimDaily(w.ctx, p.Account, "nope_"+randHex(3))
	require.ErrorIs(t, err, worldevent.ErrNotFound)

	list, err := svc.Seasonal(w.ctx, p.Account)
	require.NoError(t, err)
	var mine *worldevent.SeasonalView
	for i := range list {
		if list[i].Code == code {
			mine = &list[i]
		}
	}
	require.NotNil(t, mine)
	require.True(t, mine.ClaimedToday)
	require.GreaterOrEqual(t, mine.TokenBalance, int64(7))

	ended := "it_ended_" + randHex(3)
	_, err = w.p.Exec(w.ctx, `
		INSERT INTO content.seasonal_event (code, display_name_key, starts_at, ends_at, daily_login_tokens)
		VALUES ($1, 'seasonal.it', now() - interval '3 days', now() - interval '1 day', 7)`, ended)
	require.NoError(t, err)
	_, err = svc.ClaimDaily(w.ctx, p.Account, ended)
	require.ErrorIs(t, err, worldevent.ErrEventNotActive)
}
