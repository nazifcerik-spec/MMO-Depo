//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/contentpack"
	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/item"
)

// testEnv wires the real services the domain tests share, against the
// migrated database with the shipped content pack loaded.
type testEnv struct {
	t     *testing.T
	p     *pgxpool.Pool
	ctx   context.Context
	chars *character.Service
	txn   *economy.PgTxnService
	items *item.Service
	inv   *inventory.Service
}

var loadContentOnce sync.Once

func newWorld(t *testing.T) *testEnv {
	t.Helper()
	p, ctx := pool(t)
	loadContentOnce.Do(func() {
		_, err := contentpack.New(filepath.Join("..", "..", "..", "content"), p, nil).Load(ctx)
		require.NoError(t, err)
	})
	currencies, err := economy.LoadCurrencies(ctx, p)
	require.NoError(t, err)
	cat, err := item.LoadCatalog(ctx, p)
	require.NoError(t, err)
	items := item.NewService(item.NewPgRepo(p), cat, nil)
	inv := inventory.NewService(inventory.NewPgRepo(p), nil)
	items.WithSlotAssigner(inv)
	return &testEnv{
		t: t, p: p, ctx: ctx,
		chars: character.NewService(character.NewPgRepo(p), character.DefaultServiceConfig(), nil),
		txn:   economy.NewPgTxnService(p, currencies, nil),
		items: items, inv: inv,
	}
}

type player struct {
	Account uuid.UUID
	ID      uuid.UUID
	Name    string
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// randTag is a guild tag that will not collide with the guilds earlier runs
// left in the shared dev database: 5 of [A-Z0-9] (60M combinations).
func randTag() string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}

// player creates an account plus a level-`level` character with a unique name.
func (w *testEnv) player(class string, level int16) player {
	w.t.Helper()
	a := newAccount(w.t, w.ctx, auth.NewPgRepo(w.p))
	name := "P" + randHex(6)
	now := time.Now().UTC()
	c := character.Character{
		ID: uuid.New(), AccountID: a.ID, Name: name, NameNormalized: strings.ToLower(name),
		RaceCode: "human", ClassCode: class, Level: level, State: character.StateActive,
		CreatedAt: now, UpdatedAt: now, Version: 1,
	}
	require.NoError(w.t, character.NewPgRepo(w.p).Create(w.ctx, c))
	return player{Account: a.ID, ID: c.ID, Name: name}
}

// gold mints copper to a character through the ledger (dev mint intent).
func (w *testEnv) gold(p player, copper int64) {
	w.t.Helper()
	_, err := w.txn.Write(w.ctx, economy.Intent{
		Kind: "dev_mint", RefKind: "test", IdempotencyKey: "it-mint-" + uuid.NewString(),
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMintSource, Delta: -copper},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: p.ID, Delta: copper},
		},
	})
	require.NoError(w.t, err)
}

func (w *testEnv) balance(p player) int64 {
	w.t.Helper()
	b, err := w.txn.Balance(w.ctx, economy.AccountCharacter, p.ID, "gold")
	require.NoError(w.t, err)
	return b
}

// give mints one item instance straight into the character's bag.
func (w *testEnv) give(p player, template string, qty int32) *item.Instance {
	w.t.Helper()
	in, err := w.items.MintExact(w.ctx, p.ID, template, qty, "test", uuid.New())
	require.NoError(w.t, err)
	return in
}

// owner reports who holds an item instance and whether it is locked.
func (w *testEnv) owner(itemID uuid.UUID) (uuid.UUID, bool) {
	w.t.Helper()
	var owner uuid.UUID
	var locked bool
	require.NoError(w.t, w.p.QueryRow(w.ctx,
		`SELECT owner_character_id, is_locked FROM inventory.item_instance WHERE id = $1`, itemID).Scan(&owner, &locked))
	return owner, locked
}

// slotted reports whether the item sits in a bag slot (the bag view only lists those).
func (w *testEnv) slotted(itemID uuid.UUID) bool {
	w.t.Helper()
	var ok bool
	require.NoError(w.t, w.p.QueryRow(w.ctx,
		`SELECT slot IS NOT NULL FROM inventory.item_instance WHERE id = $1`, itemID).Scan(&ok))
	return ok
}

// events records WS publishes for assertions.
type events struct {
	mu    sync.Mutex
	kinds []string
}

func (e *events) Publish(_ uuid.UUID, kind string, _ any) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.kinds = append(e.kinds, kind)
	return 1
}

func (e *events) BroadcastKind(kind string, payload any) int {
	return e.Publish(uuid.Nil, kind, payload)
}

func (e *events) has(kind string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, k := range e.kinds {
		if k == kind {
			return true
		}
	}
	return false
}

func bagIDs(items []inventory.BagItem) []uuid.UUID {
	out := make([]uuid.UUID, len(items))
	for i, it := range items {
		out[i] = it.InstanceID
	}
	return out
}

// placeAt moves characters to a region, as a finished journey would.
func (w *testEnv) placeAt(t *testing.T, region string, chars ...string) {
	t.Helper()
	for _, c := range chars {
		_, err := w.p.Exec(w.ctx, `UPDATE character.character SET location_region_code = $2 WHERE id = $1`, c, region)
		require.NoError(t, err)
	}
}

func (e *events) reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.kinds = nil
}
