package guild

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissionMatrix(t *testing.T) {
	require.True(t, CanInvite(RankLeader))
	require.True(t, CanInvite(RankOfficer))
	require.True(t, CanInvite(RankVeteran))
	require.False(t, CanInvite(RankMember))
	require.False(t, CanInvite(RankInitiate))

	require.True(t, CanSetMOTD(RankOfficer))
	require.False(t, CanSetMOTD(RankVeteran))

	require.True(t, CanKick(RankLeader, RankOfficer))
	require.True(t, CanKick(RankOfficer, RankVeteran))
	require.False(t, CanKick(RankOfficer, RankOfficer), "no kicking peers")
	require.False(t, CanKick(RankOfficer, RankLeader))
	require.False(t, CanKick(RankVeteran, RankInitiate), "veterans cannot kick")
}

func TestCanSetRank(t *testing.T) {
	// Leader: any lower member to 1..4.
	require.True(t, CanSetRank(RankLeader, RankInitiate, RankOfficer))
	require.True(t, CanSetRank(RankLeader, RankOfficer, RankInitiate))
	require.False(t, CanSetRank(RankLeader, RankInitiate, RankLeader), "rank 0 only via transfer")
	require.False(t, CanSetRank(RankLeader, RankInitiate, 5))

	// Officer: members below officer, up to veteran only.
	require.True(t, CanSetRank(RankOfficer, RankInitiate, RankVeteran))
	require.True(t, CanSetRank(RankOfficer, RankVeteran, RankInitiate))
	require.False(t, CanSetRank(RankOfficer, RankInitiate, RankOfficer), "cannot mint peers")
	require.False(t, CanSetRank(RankOfficer, RankOfficer, RankMember), "cannot demote peers")

	// Everyone else: nothing.
	require.False(t, CanSetRank(RankVeteran, RankInitiate, RankMember))
}

func TestNormalizeName(t *testing.T) {
	disp, norm, err := NormalizeName("  Emerald Wardens ")
	require.NoError(t, err)
	require.Equal(t, "Emerald Wardens", disp)
	require.Equal(t, "emerald wardens", norm)

	_, _, err = NormalizeName("Şövalyeler")
	require.NoError(t, err, "unicode letters allowed")

	for _, bad := range []string{"ab", "double  space", "bad!name", "", "a234567890123456789012345"} {
		_, _, err := NormalizeName(bad)
		require.ErrorIs(t, err, ErrInvalidName, bad)
	}
}

func TestNormalizeTag(t *testing.T) {
	disp, norm, err := NormalizeTag("ewd")
	require.NoError(t, err)
	require.Equal(t, "EWD", disp)
	require.Equal(t, "ewd", norm)
	for _, bad := range []string{"a", "toolong", "e w", "e-d"} {
		_, _, err := NormalizeTag(bad)
		require.ErrorIs(t, err, ErrInvalidTag, bad)
	}
}
