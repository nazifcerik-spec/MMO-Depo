package dailyreward

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRewardFor(t *testing.T) {
	r := RewardFor(1, 100)
	require.EqualValues(t, 50000, r.GoldCopper, "50 × 100^1.5")
	require.Empty(t, r.Potion)
	require.EqualValues(t, 3*r.GoldCopper, RewardFor(7, 100).GoldCopper, "a week's seventh day pays triple")
	require.EqualValues(t, 3*r.GoldCopper/2, RewardFor(8, 100).GoldCopper, "the second week pays more")
	require.EqualValues(t, 10*r.GoldCopper, RewardFor(29, 100).GoldCopper, "day 29 is the big one")
	require.Equal(t, "potion_hp_small", RewardFor(24, 1).Potion)
	require.Empty(t, RewardFor(29, 1).Potion)
	require.Equal(t, Reward{Day: 3, GoldCopper: 50, Potion: "potion_hp_small", PotionQty: 3}, RewardFor(3, 0),
		"level 0 counts as 1")
	require.Equal(t, "potion_hp_medium", RewardFor(6, 249).Potion)
	require.Equal(t, "potion_hp_supreme", RewardFor(6, 1000).Potion)
	require.EqualValues(t, 5, RewardFor(6, 1000).PotionQty)
}

func TestNextDay(t *testing.T) {
	today := time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	require.Equal(t, 1, nextDay(time.Time{}, 0, today), "first claim")
	require.Equal(t, 3, nextDay(yesterday, 2, today), "streak goes on")
	require.Equal(t, 8, nextDay(yesterday, 7, today), "week two follows week one")
	require.Equal(t, 1, nextDay(yesterday, 29, today), "the cycle wraps after day 29")
	require.Equal(t, 1, nextDay(today.AddDate(0, 0, -2), 5, today), "a missed day restarts")
}

func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for _, e := range []error{ErrClaimed, ErrNotFound} {
		rec := httptest.NewRecorder()
		writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("ctx: %w", e))
		require.Less(t, rec.Code, 500, e.Error())
	}
}
