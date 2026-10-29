package boost

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type fixedKeys map[string]*ecdsa.PublicKey

func (f fixedKeys) Key(_ context.Context, id string) (*ecdsa.PublicKey, error) {
	if k, ok := f[id]; ok {
		return k, nil
	}
	return nil, ErrBadReward
}

func signed(t *testing.T, key *ecdsa.PrivateKey, content string) string {
	sum := sha256.Sum256([]byte(content))
	sig, err := ecdsa.SignASN1(rand.Reader, key, sum[:])
	require.NoError(t, err)
	return content + "&signature=" + base64.RawURLEncoding.EncodeToString(sig) + "&key_id=7"
}

func TestVerifySSV(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	keys := fixedKeys{"7": &key.PublicKey}
	content := "ad_network=5450213213286189855&ad_unit=123&custom_data=gather&key_id=7&reward_amount=1" +
		"&reward_item=boost&timestamp=1761700000000&transaction_id=tx1&user_id=c0ffee"
	r, err := VerifySSV(context.Background(), signed(t, key, content), keys)
	require.NoError(t, err)
	require.Equal(t, Reward{TransactionID: "tx1", CharacterID: "c0ffee", Category: "gather", At: r.At}, *r)
	require.Equal(t, int64(1761700000000), r.At.UnixMilli())

	// A changed reward no longer matches its signature.
	forged := signed(t, key, content)
	forged = "custom_data=craft" + forged[len("ad_network=5450213213286189855&ad_unit=123&custom_data=gather"):]
	_, err = VerifySSV(context.Background(), forged, keys)
	require.ErrorIs(t, err, ErrBadReward)
	_, err = VerifySSV(context.Background(), content, keys)
	require.ErrorIs(t, err, ErrBadReward, "no signature")
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	_, err = VerifySSV(context.Background(), signed(t, other, content), keys)
	require.ErrorIs(t, err, ErrBadReward, "someone else's key")
}

func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for _, e := range []error{ErrCategory, ErrActive, ErrAdLimit, ErrNotFound, ErrAdsOff, ErrBadReward} {
		rec := httptest.NewRecorder()
		writeErr(rec, httptest.NewRequest(http.MethodGet, "/", nil), fmt.Errorf("ctx: %w", e))
		require.Less(t, rec.Code, 500, e.Error())
	}
}
