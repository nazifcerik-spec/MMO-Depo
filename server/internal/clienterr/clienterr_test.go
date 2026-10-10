package clienterr

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanRedactsAndBounds(t *testing.T) {
	r, err := clean(Report{Kind: "flutter", Platform: "android", AppVersion: "1.0.0+3", Flavor: "prod",
		Message: "  Failed for bob@example.com with Bearer abc.def and eyJhbGc.eyJzdWI.sig?token=s3cret&x=1  ",
		Stack:   strings.Repeat("#0 frame\n", 100)})
	require.NoError(t, err)
	require.Equal(t, "Failed for ***@*** with Bearer *** and ***?token=***&x=1", r.Message)
	require.Len(t, strings.Split(r.Stack, "\n"), MaxStackRows)

	r, err = clean(Report{Kind: "zone", Message: strings.Repeat("m", MaxMessage+50), Platform: "beos",
		AppVersion: "1.0 <script>", Flavor: strings.Repeat("f", 17)})
	require.NoError(t, err)
	require.Len(t, r.Message, MaxMessage)
	require.Equal(t, "other", r.Platform)
	require.Empty(t, r.AppVersion)
	require.Empty(t, r.Flavor)

	for _, bad := range []Report{{Kind: "flutter", Message: "  "}, {Kind: "native", Message: "x"}} {
		_, err := clean(bad)
		require.ErrorIs(t, err, ErrInvalid)
	}
}

// The same bug on another device/build lands in the same group; a
// different top frame or kind does not.
func TestFingerprint(t *testing.T) {
	a := Report{Kind: "flutter", Message: "RangeError (index): Invalid value: 5\nmore",
		Stack: "#0      List.[] (dart:core-patch/growable_array.dart:264:36)\n" +
			"#1      BagGrid.build (package:otter/bag.dart:88:12)\n#2 x\n#3 tail-a"}
	b := Report{Kind: "flutter", Message: "RangeError (index): Invalid value: 12",
		Stack: "#0      List.[] (dart:core-patch/growable_array.dart:270:1)\n" +
			"#1      BagGrid.build (package:otter/bag.dart:91:3)\n#2 x\n#3 tail-b"}
	require.Equal(t, fingerprint(a), fingerprint(b))
	c := b
	c.Stack = "#0      List.[] (dart:core-patch/growable_array.dart:270:1)\n#1      ShopList.build (package:otter/shop.dart:9:3)"
	require.NotEqual(t, fingerprint(a), fingerprint(c))
	d := a
	d.Kind = "platform"
	require.NotEqual(t, fingerprint(a), fingerprint(d))
}

func TestHandlerRejects(t *testing.T) {
	h := NewHandler(NewService(nil, nil))
	for body, code := range map[string]int{"{": 400, `{"kind":"nope","message":"x"}`: 400} {
		rec := httptest.NewRecorder()
		h.Report(rec, httptest.NewRequest(http.MethodPost, "/v1/client/errors", strings.NewReader(body)))
		require.Equal(t, code, rec.Code, body)
	}
}
