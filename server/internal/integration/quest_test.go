//go:build integration

package integration

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/progression"
	"github.com/otter/server/internal/quest"
)

// DD-221 SQL: accept gates, event ticks, single-payout turn-in, daily
// once-per-day and past-day expiry.
func TestQuestSQL(t *testing.T) {
	p, ctx := pool(t)
	loader := contentpack.New(filepath.Join("..", "..", "..", "content"), p, nil)
	_, err := loader.Load(ctx)
	require.NoError(t, err)

	a := newAccount(t, ctx, auth.NewPgRepo(p))
	repo := character.NewPgRepo(p)
	now := time.Now().UTC()
	c := character.Character{
		ID: uuid.New(), AccountID: a.ID, Name: "Qx" + a.Username[3:9], NameNormalized: strings.ToLower("qx" + a.Username[3:9]),
		RaceCode: "human", ClassCode: "warrior", Level: 1, State: character.StateActive,
		CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	require.NoError(t, repo.Create(ctx, c))
	chars := character.NewService(repo, character.DefaultServiceConfig(), nil)
	txn := economy.NewPgTxnService(p, economy.StaticCurrencies(), nil)
	clock := now
	svc := quest.NewService(p, chars, txn, chars, nil).WithClock(func() time.Time { return clock })
	require.NoError(t, svc.Reload(ctx))

	codes := func(vs []quest.View) []string {
		out := []string{}
		for _, v := range vs {
			out = append(out, v.Code)
		}
		return out
	}
	board, err := svc.List(ctx, c.ID, a.ID)
	require.NoError(t, err)
	require.Contains(t, codes(board.Available), "main_01_first_hunt")
	require.NotContains(t, codes(board.Available), "main_02_supplies", "locked behind main_01")
	require.NotContains(t, codes(board.Available), "side_bandit_trouble", "level 5")

	_, err = svc.Accept(ctx, c.ID, a.ID, "main_02_supplies")
	require.ErrorIs(t, err, quest.ErrLocked)
	_, err = svc.Accept(ctx, c.ID, uuid.New(), "main_01_first_hunt")
	require.Error(t, err, "foreign account")
	_, err = svc.Accept(ctx, c.ID, a.ID, "main_01_first_hunt")
	require.NoError(t, err)
	_, err = svc.Accept(ctx, c.ID, a.ID, "main_01_first_hunt")
	require.ErrorIs(t, err, quest.ErrTaken)
	_, err = svc.TurnIn(ctx, c.ID, a.ID, "main_01_first_hunt")
	require.ErrorIs(t, err, quest.ErrIncomplete)

	svc.Record(ctx, c.ID, progression.Event{Kind: "mob_killed", Amount: 9}, progression.Event{Kind: "mob_killed:young_boar", Amount: 3})
	board, err = svc.List(ctx, c.ID, a.ID)
	require.NoError(t, err)
	require.Equal(t, int64(3), board.Open[0].Objectives[0].Progress)
	svc.Record(ctx, c.ID, progression.Event{Kind: "mob_killed:young_boar", Amount: 7})
	board, err = svc.List(ctx, c.ID, a.ID)
	require.NoError(t, err)
	require.Equal(t, "completed", board.Open[0].State)
	require.Equal(t, int64(5), board.Open[0].Objectives[0].Progress)

	res, err := svc.TurnIn(ctx, c.ID, a.ID, "main_01_first_hunt")
	require.NoError(t, err)
	require.Equal(t, int64(200), res.RewardGold)
	require.NotNil(t, res.XP)
	_, err = svc.TurnIn(ctx, c.ID, a.ID, "main_01_first_hunt")
	require.ErrorIs(t, err, quest.ErrNotActive)
	gold, err := txn.Balance(ctx, economy.AccountCharacter, c.ID, "gold")
	require.NoError(t, err)
	require.Equal(t, int64(200), gold, "paid exactly once")
	_, err = svc.Accept(ctx, c.ID, a.ID, "main_01_first_hunt")
	require.ErrorIs(t, err, quest.ErrTaken, "main quests are one-shot")

	board, err = svc.List(ctx, c.ID, a.ID)
	require.NoError(t, err)
	require.Contains(t, codes(board.Available), "main_02_supplies")

	// Dailies: only today's offer, once per day, expire at the UTC day end.
	offer := []string{}
	for _, v := range board.Available {
		if v.Category == quest.CategoryDaily {
			offer = append(offer, v.Code)
		}
	}
	require.Len(t, offer, quest.DailySlots)
	for _, code := range []string{"daily_hunter", "daily_miner", "daily_angler", "daily_lumberjack", "daily_herbalist"} {
		if !contains(offer, code) {
			_, err = svc.Accept(ctx, c.ID, a.ID, code)
			require.ErrorIs(t, err, quest.ErrNotOffered)
			break
		}
	}
	_, err = svc.Accept(ctx, c.ID, a.ID, offer[0])
	require.NoError(t, err)
	require.NoError(t, svc.Abandon(ctx, c.ID, a.ID, offer[0]))
	_, err = svc.Accept(ctx, c.ID, a.ID, offer[0])
	require.ErrorIs(t, err, quest.ErrTaken, "abandoning does not re-roll today's daily")
	_, err = svc.Accept(ctx, c.ID, a.ID, offer[1])
	require.NoError(t, err)

	clock = clock.AddDate(0, 0, 1)
	svc.Record(ctx, c.ID, progression.Event{Kind: "mob_killed", Amount: 1000}, progression.Event{Kind: "gather_mining", Amount: 1000},
		progression.Event{Kind: "gather_fishing", Amount: 1000}, progression.Event{Kind: "gather_woodcutting", Amount: 1000},
		progression.Event{Kind: "gather_herbalism", Amount: 1000})
	board, err = svc.List(ctx, c.ID, a.ID)
	require.NoError(t, err)
	require.NotContains(t, codes(board.Open), offer[1], "yesterday's daily is gone")
	_, err = svc.Accept(ctx, c.ID, a.ID, "main_02_supplies")
	require.NoError(t, err)
	var state string
	require.NoError(t, p.QueryRow(ctx, `SELECT state FROM progression.character_quest
		WHERE character_id = $1 AND quest_code = $2 AND day_key IS NOT NULL AND state <> 'abandoned'`, c.ID, offer[1]).Scan(&state))
	require.Equal(t, "expired", state, "past-day open daily expires, never ticked after midnight")
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
