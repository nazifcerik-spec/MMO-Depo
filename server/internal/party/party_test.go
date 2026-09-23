package party

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNextLeader_EarliestJoinerAmongRest(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	leader, a, b := uuid.New(), uuid.New(), uuid.New()
	members := []Member{
		{CharacterID: b, JoinedAt: t0.Add(2 * time.Minute)},
		{CharacterID: leader, JoinedAt: t0},
		{CharacterID: a, JoinedAt: t0.Add(time.Minute)},
	}
	next, ok := NextLeader(members, leader)
	require.True(t, ok)
	require.Equal(t, a, next)
}

func TestNextLeader_NobodyLeft(t *testing.T) {
	only := uuid.New()
	_, ok := NextLeader([]Member{{CharacterID: only}}, only)
	require.False(t, ok)
	_, ok = NextLeader(nil, only)
	require.False(t, ok)
}
