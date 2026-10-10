package admin

import (
	"encoding/base32"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// RFC 6238 Appendix B vectors (SHA-1 key "12345678901234567890"), last 6 digits.
func TestTOTPCode_RFC6238Vectors(t *testing.T) {
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte("12345678901234567890"))
	for unix, want := range map[int64]string{
		59:          "287082",
		1111111109:  "081804",
		1111111111:  "050471",
		1234567890:  "005924",
		2000000000:  "279037",
		20000000000: "353130",
	} {
		got, err := TOTPCode(secret, unix/30)
		require.NoError(t, err)
		require.Equal(t, want, got, "t=%d", unix)
	}
}

func TestVerifyTOTP_DriftAndReplay(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	require.NoError(t, err)
	now := time.Unix(1_800_000_000, 0)
	step := now.Unix() / 30
	prev, _ := TOTPCode(secret, step-1)
	cur, _ := TOTPCode(secret, step)
	far, _ := TOTPCode(secret, step-3)

	got, ok := VerifyTOTP(secret, prev, now, 0)
	require.True(t, ok, "one step of clock drift accepted")
	require.Equal(t, step-1, got)

	_, ok = VerifyTOTP(secret, prev, now, step-1)
	require.False(t, ok, "same step cannot be reused")
	_, ok = VerifyTOTP(secret, cur, now, step-1)
	require.True(t, ok)
	_, ok = VerifyTOTP(secret, far, now, 0)
	require.False(t, ok, "codes outside ±1 step rejected")
	_, ok = VerifyTOTP(secret, "12345", now, 0)
	require.False(t, ok)
}

func TestPermissions_Matrix(t *testing.T) {
	viewer := Permissions([]string{"viewer"})
	require.True(t, viewer.Has(PermPlayerRead))
	require.False(t, viewer.Has(PermPlayerPIIReveal))
	require.False(t, viewer.CanReverse(1))

	l2 := Permissions([]string{"support_l2"})
	require.True(t, l2.CanReverse(ReversalLimitCopper))
	require.False(t, l2.CanReverse(ReversalLimitCopper+1))
	require.True(t, l2.CanBan(7))
	require.False(t, l2.CanBan(8))
	require.False(t, l2.CanBan(0), "permanent ban is admin-only")
	require.True(t, l2.CanMute(24))
	require.False(t, l2.CanMute(25))

	adm := Permissions([]string{"admin"})
	require.True(t, adm.CanReverse(ReversalLimitCopper*10))
	require.True(t, adm.CanBan(0))
	require.True(t, adm.Has(PermRoleGrant))

	multi := Permissions([]string{"viewer", "game_designer"})
	require.True(t, multi.Has(PermEventTrigger))
	require.False(t, ValidRole("root"))
}

func TestMaskEmailAndGoldMoved(t *testing.T) {
	require.Equal(t, "n***@example.com", MaskEmail("nazif@example.com"))
	require.Equal(t, "***", MaskEmail("broken"))
	require.Equal(t, int64(900), GoldMoved([]LedgerEntry{
		{Currency: "gold", Delta: -1000}, {Currency: "gold", Delta: 900},
		{Currency: "gold", Delta: 0}, {Currency: "gem", Delta: 50},
	}))
}
