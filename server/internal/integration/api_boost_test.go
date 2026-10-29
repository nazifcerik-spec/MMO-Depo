//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/boost"
)

// Timed XP boosts (DD-320): +100% on the share of a job a boost covered;
// three ad boosts a day per category; one at a time.
func TestBoosts(t *testing.T) {
	w := newWorld(t)
	p := newHTTPPlayer(t, "warrior")
	char := uuid.MustParse(p.char)
	ctx := context.Background()
	now := time.Now().UTC()
	svc := boost.NewService(w.p).WithClock(func() time.Time { return now })

	list := p.ok("GET", p.path("/boosts"), nil)["boosts"].([]any)
	require.Len(t, list, 3)
	require.EqualValues(t, boost.AdsPerDay, list[0].(map[string]any)["ads_left"])

	st, err := svc.Grant(ctx, char, boost.CatGather, boost.SourceAd, "t1-"+p.char)
	require.NoError(t, err)
	require.NotNil(t, st.EndsAt)
	_, err = svc.Grant(ctx, char, boost.CatGather, boost.SourceAd, "t2-"+p.char)
	require.ErrorIs(t, err, boost.ErrActive)
	_, err = svc.Grant(ctx, char, boost.CatGather, boost.SourceAd, "t1-"+p.char)
	require.NoError(t, err, "the same reward twice grants once")
	_, err = svc.Grant(ctx, char, "fishing", boost.SourceAd, "t3-"+p.char)
	require.ErrorIs(t, err, boost.ErrCategory)

	// A two-hour boost over the second half of a four-hour job: +50%.
	require.EqualValues(t, 150, svc.Apply(ctx, char, boost.CatGather, now.Add(-2*time.Hour), now.Add(2*time.Hour), 100))
	require.EqualValues(t, 200, svc.Apply(ctx, char, boost.CatGather, now, now.Add(time.Hour), 100), "fully covered")
	require.EqualValues(t, 100, svc.Apply(ctx, char, boost.CatAdventure, now, now.Add(time.Hour), 100), "other category")
	require.EqualValues(t, 100, svc.Apply(ctx, char, boost.CatGather, now.Add(-3*time.Hour), now.Add(-time.Hour), 100),
		"before the boost")
	require.EqualValues(t, 200, svc.Apply(ctx, char, boost.CatGather, now.Add(time.Minute), now.Add(time.Minute), 100),
		"an instant job while it runs")

	// Ads per day: three, the others wait for tomorrow.
	later := now
	for i := 0; i < boost.AdsPerDay-1; i++ {
		later = later.Add(boost.Duration + time.Minute)
		at := later
		_, err := boost.NewService(w.p).WithClock(func() time.Time { return at }).
			Grant(ctx, char, boost.CatGather, boost.SourceAd, uuid.NewString())
		if at.Day() != now.Day() {
			break // the run crossed midnight: a new day's allowance
		}
		require.NoError(t, err)
	}
	if later.Add(boost.Duration+time.Minute).Day() == now.Day() {
		at := later.Add(boost.Duration + time.Minute)
		_, err := boost.NewService(w.p).WithClock(func() time.Time { return at }).
			Grant(ctx, char, boost.CatGather, boost.SourceAd, uuid.NewString())
		require.ErrorIs(t, err, boost.ErrAdLimit)
	}

	// Development builds watch a stand-in ad; production needs the network.
	res := p.ok("POST", p.path("/boosts/adventure/dev-ad"), map[string]any{})
	require.Equal(t, "adventure", res["category"])
	p.fail("GET", "/v1/ads/ssv?user_id=x&signature=abc&key_id=1", nil, 400, "E_BAD_REWARD")
}
