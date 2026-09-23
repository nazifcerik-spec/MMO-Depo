package economy

import (
	"context"

	"github.com/google/uuid"
)

// Wallet is the compound response for GET /v1/character/{id}/wallet. Extends
// as new currencies land — v1 exposes gold + gem + honor.
type Wallet struct {
	CharacterID string `json:"character_id"`
	GoldCopper  int64  `json:"gold_copper"`
	Honor       int64  `json:"honor,omitempty"`
	Gem         int64  `json:"gem,omitempty"` // account-scope; requires account_id
}

// ReadWallet fetches a character's balances via the TxnService's Balance
// interface — same read path production and tests use. account is optional
// for account-scope currency lookup (nil UUID → skipped).
func ReadWallet(ctx context.Context, svc BalanceReader, characterID, accountID uuid.UUID) (*Wallet, error) {
	w := &Wallet{CharacterID: characterID.String()}
	var err error
	if w.GoldCopper, err = svc.Balance(ctx, AccountCharacter, characterID, "gold"); err != nil {
		return nil, err
	}
	if w.Honor, err = svc.Balance(ctx, AccountCharacter, characterID, "honor"); err != nil {
		return nil, err
	}
	if accountID != uuid.Nil {
		if w.Gem, err = svc.Balance(ctx, AccountAccount, accountID, "gem"); err != nil {
			return nil, err
		}
	}
	return w, nil
}

// BalanceReader is the minimal read contract Wallet needs.
type BalanceReader interface {
	Balance(ctx context.Context, kind AccountKind, ref uuid.UUID, currency string) (int64, error)
}
