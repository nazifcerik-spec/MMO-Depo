package helper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypesCoverEveryProfession(t *testing.T) {
	seen := map[string]bool{}
	for _, ty := range Types {
		require.False(t, seen[ty.Profession], ty.Profession)
		seen[ty.Profession] = true
		got, ok := TypeOf(ty.Code)
		require.True(t, ok)
		require.Equal(t, ty, got)
	}
	require.Len(t, seen, 14)
	_, ok := TypeOf("nobody")
	require.False(t, ok)
}

// Half the character's speed and XP, then the type's own bonus.
func TestHalfEfficiencyAndBonuses(t *testing.T) {
	miner, _ := TypeOf("miner")   // yield +15
	angler, _ := TypeOf("angler") // xp +25
	smith, _ := TypeOf("smith")   // speed +25
	weaver, _ := TypeOf("weaver") // gold -30
	cook, _ := TypeOf("cook")     // double 15

	require.InDelta(t, 0.5*1.15-1, miner.YieldMult(), 1e-9)
	require.InDelta(t, -0.5, angler.YieldMult(), 1e-9)
	require.Equal(t, 50, miner.XP(100))
	require.Equal(t, 63, angler.XP(100), "50 × 1.25 = 62.5")
	require.EqualValues(t, 600, weaver.CraftSeconds(30, 10), "twice the time")
	require.EqualValues(t, 450, smith.CraftSeconds(30, 10), "less the speed bonus")
	require.EqualValues(t, 700, weaver.GoldCost(100, 10))
	require.EqualValues(t, 1000, smith.GoldCost(100, 10))
	require.Equal(t, 2, cook.Doubles(4, func(i int) int { return []int{3, 50, 14, 99}[i] }))
	require.Zero(t, smith.Doubles(4, func(int) int { return 0 }))
}

func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for _, e := range []error{ErrUnknownType, ErrNotOwned, ErrOwned, ErrBusy, ErrIdle, ErrNotReady, ErrValidation,
		ErrNoGold, ErrSkill, ErrMaterials} {
		rec := httptest.NewRecorder()
		writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("ctx: %w", e))
		require.Less(t, rec.Code, 500, e.Error())
	}
}
