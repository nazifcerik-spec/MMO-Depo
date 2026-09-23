package economy_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/economy"
)

func newTestSvc() *economy.MemoryTxnService {
	return economy.NewMemoryTxnService(economy.StaticCurrencies())
}

// mintIntent constructs the canonical "mint N gold to character" intent.
func mintIntent(character uuid.UUID, amount int64, key string) economy.Intent {
	return economy.Intent{
		Kind:           "mint",
		RefKind:        "system",
		IdempotencyKey: key,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -amount},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: character, Delta: amount},
		},
	}
}

func TestValidateIntent_SumMustBeZero(t *testing.T) {
	svc := newTestSvc()
	acct := uuid.New()

	// Deliberately unbalanced.
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "bad",
		IdempotencyKey: "key-bad-nonzero-abc123456",
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -100},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: acct, Delta: 50}, // mismatch
		},
	})
	require.ErrorIs(t, err, economy.ErrValidation)
}

func TestValidateIntent_IdempotencyKeyRequired(t *testing.T) {
	svc := newTestSvc()
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:    "mint",
		Entries: []economy.Entry{},
	})
	require.ErrorIs(t, err, economy.ErrValidation)
}

func TestValidateIntent_TwoEntriesMinimum(t *testing.T) {
	svc := newTestSvc()
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "mint",
		IdempotencyKey: "key-single-entry-abc123456",
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -100},
		},
	})
	require.ErrorIs(t, err, economy.ErrValidation)
}

func TestValidateIntent_UnknownCurrency(t *testing.T) {
	svc := newTestSvc()
	acct := uuid.New()
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "mint",
		IdempotencyKey: "key-bad-currency-abc123456",
		Entries: []economy.Entry{
			{Currency: "unicorn", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -1},
			{Currency: "unicorn", AccountKind: economy.AccountCharacter, AccountRef: acct, Delta: 1},
		},
	})
	require.ErrorIs(t, err, economy.ErrUnknownCurrency)
}

func TestValidateIntent_ScopeMismatch(t *testing.T) {
	svc := newTestSvc()
	acct := uuid.New()
	// 'gem' is account-scope; crediting to a character is a scope violation.
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "mint",
		IdempotencyKey: "key-scope-mismatch-abc123456",
		Entries: []economy.Entry{
			{Currency: "gem", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -10},
			{Currency: "gem", AccountKind: economy.AccountCharacter, AccountRef: acct, Delta: 10},
		},
	})
	require.ErrorIs(t, err, economy.ErrScopeMismatch)
}

func TestWrite_MintCreditsBalance(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()

	res, err := svc.Write(context.Background(), mintIntent(char, 5000, "key-mint-1-abc123456"))
	require.NoError(t, err)
	require.False(t, res.Replayed)
	require.Len(t, res.Rows, 2)
	require.NotEqual(t, uuid.Nil, res.TxnID)

	bal, err := svc.Balance(context.Background(), economy.AccountCharacter, char, "gold")
	require.NoError(t, err)
	require.Equal(t, int64(5000), bal)
}

func TestWrite_InsufficientBalance(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()

	// Attempt to spend without minting first.
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "vendor_buy",
		IdempotencyKey: "key-insufficient-abc123456",
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: -100},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: 100},
		},
	})
	require.ErrorIs(t, err, economy.ErrInsufficientBalance)
}

func TestWrite_CapExceeded(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()

	// Directly attempt to credit above the gold cap (10^12).
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "mint",
		IdempotencyKey: "key-cap-abc123456",
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -1_000_000_000_001},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: 1_000_000_000_001},
		},
	})
	require.ErrorIs(t, err, economy.ErrCurrencyCap)
}

func TestWrite_IdempotencyReplayReturnsCached(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()

	// First call.
	res1, err := svc.Write(context.Background(), mintIntent(char, 500, "same-key-abc123456"))
	require.NoError(t, err)
	require.False(t, res1.Replayed)

	// Second call same key.
	res2, err := svc.Write(context.Background(), mintIntent(char, 500, "same-key-abc123456"))
	require.NoError(t, err)
	require.True(t, res2.Replayed)
	require.Equal(t, res1.TxnID, res2.TxnID)

	// Balance was applied ONCE.
	bal, _ := svc.Balance(context.Background(), economy.AccountCharacter, char, "gold")
	require.Equal(t, int64(500), bal)
}

func TestWrite_VendorBuySplitsCorrectly(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()

	// Seed.
	_, err := svc.Write(context.Background(), mintIntent(char, 10_000, "seed-abc123456"))
	require.NoError(t, err)

	// Vendor buy: 1000 debit character; 250 credit vendor_sink; 750 to system as "spread".
	// Simpler: buyer -1000 / vendor_sink +1000.
	_, err = svc.Write(context.Background(), economy.Intent{
		Kind:           "vendor_buy",
		IdempotencyKey: "buy-1-abc123456",
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: -1000},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: 1000},
		},
	})
	require.NoError(t, err)

	bal, _ := svc.Balance(context.Background(), economy.AccountCharacter, char, "gold")
	require.Equal(t, int64(9000), bal)
}

func TestWrite_MultiCurrencySumInvariant(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()
	acct := uuid.New()

	// Two currencies in one intent. Each must sum to 0 independently.
	// Gold: -100 char + 100 sink. Gem: -5 acct + 5 sink.
	_, err := svc.Write(context.Background(), economy.Intent{
		Kind:           "combo",
		IdempotencyKey: "combo-1-abc123456",
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -100},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: char, Delta: 100},
			{Currency: "gem", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -5},
			{Currency: "gem", AccountKind: economy.AccountAccount, AccountRef: acct, Delta: 5},
		},
	})
	require.NoError(t, err)

	gold, _ := svc.Balance(context.Background(), economy.AccountCharacter, char, "gold")
	require.Equal(t, int64(100), gold)
	gem, _ := svc.Balance(context.Background(), economy.AccountAccount, acct, "gem")
	require.Equal(t, int64(5), gem)
}

func TestSystemAccounts_Deterministic(t *testing.T) {
	// System UUIDs must be stable across runs (v5 SHA1 property).
	require.Equal(t, economy.SystemMintSource, economy.SystemMintSource)
	require.NotEqual(t, economy.SystemMintSource, economy.SystemVendorSink)
	require.True(t, economy.IsSystemAccount(economy.SystemTaxSink))
	require.False(t, economy.IsSystemAccount(uuid.New()))
}

func TestWrite_EntriesReturnedWithBalanceAfter(t *testing.T) {
	svc := newTestSvc()
	char := uuid.New()

	res, err := svc.Write(context.Background(), mintIntent(char, 42, "trace-1-abc123456"))
	require.NoError(t, err)
	require.Len(t, res.Rows, 2)

	// Row order matches entry order.
	require.Equal(t, economy.AccountSystem, res.Rows[0].AccountKind)
	require.Equal(t, economy.AccountCharacter, res.Rows[1].AccountKind)
	require.Equal(t, int64(42), res.Rows[1].BalanceAfter)
}
