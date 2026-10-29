// Package shop is the credit shop and the cosmetics (DD-321). Credits are
// an account currency bought with real money; until a payment provider is
// set up they reach accounts only through operators (grants) and the dev
// mint. They buy timed XP boosts (DD-320) and premium cosmetics: 30 page
// backgrounds and 80 avatars, of which the first 10 and 40 are free.
package shop

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/boost"
	"github.com/otter/server/internal/economy"
)

// Currency is the credit currency code.
const Currency = "credit"

// Cosmetic counts and how many of each are free.
const (
	Backgrounds     = 30
	FreeBackgrounds = 10
	Avatars         = 80
	FreeAvatars     = 40

	BackgroundPrice = 30 // credits
	AvatarPrice     = 15
	BoostPrice      = 20
)

// Item is something the shop sells.
type Item struct {
	Code  string `json:"code"`
	Kind  string `json:"kind"` // boost | background | avatar
	Price int64  `json:"price_credits"`
	Owned bool   `json:"owned,omitempty"`
}

var (
	ErrUnknown  = errors.New("shop: unknown item")
	ErrOwned    = errors.New("shop: already owned")
	ErrNoCredit = errors.New("shop: not enough credits")
	ErrLocked   = errors.New("shop: cosmetic not owned")
)

// BackgroundCode / AvatarCode name cosmetic n (1-based).
func BackgroundCode(n int) string { return fmt.Sprintf("bg%02d", n) }
func AvatarCode(n int) string     { return fmt.Sprintf("av%02d", n) }

// cosmetic parses a cosmetic code: its kind, number and whether it is free.
func cosmetic(code string) (kind string, free, ok bool) {
	var n int
	switch {
	case strings.HasPrefix(code, "bg"):
		if _, err := fmt.Sscanf(code, "bg%02d", &n); err != nil || n < 1 || n > Backgrounds || code != BackgroundCode(n) {
			return "", false, false
		}
		return "background", n <= FreeBackgrounds, true
	case strings.HasPrefix(code, "av"):
		if _, err := fmt.Sscanf(code, "av%02d", &n); err != nil || n < 1 || n > Avatars || code != AvatarCode(n) {
			return "", false, false
		}
		return "avatar", n <= FreeAvatars, true
	}
	return "", false, false
}

// Catalog is everything for sale: the three boosts and the premium cosmetics.
func Catalog() []Item {
	var out []Item
	for _, c := range boost.Categories {
		out = append(out, Item{Code: "boost_" + c, Kind: "boost", Price: BoostPrice})
	}
	for n := FreeBackgrounds + 1; n <= Backgrounds; n++ {
		out = append(out, Item{Code: BackgroundCode(n), Kind: "background", Price: BackgroundPrice})
	}
	for n := FreeAvatars + 1; n <= Avatars; n++ {
		out = append(out, Item{Code: AvatarCode(n), Kind: "avatar", Price: AvatarPrice})
	}
	return out
}

func find(code string) (Item, bool) {
	for _, it := range Catalog() {
		if it.Code == code {
			return it, true
		}
	}
	return Item{}, false
}

// Booster grants a boost (boost.Service).
type Booster interface {
	Grant(ctx context.Context, char uuid.UUID, cat, source, ref string) (*boost.Status, error)
}

type Service struct {
	pool  *pgxpool.Pool
	txn   economy.TxnService
	boost Booster
}

func NewService(pool *pgxpool.Pool, txn economy.TxnService, b Booster) *Service {
	return &Service{pool: pool, txn: txn, boost: b}
}

// View is the shop for one account and character.
type View struct {
	Credits int64  `json:"credits"`
	Items   []Item `json:"items"`
	// Paid says whether credits can be bought with money yet.
	PaymentsEnabled bool `json:"payments_enabled"`
}

// Shop is the catalog with what the account owns.
func (s *Service) Shop(ctx context.Context, account uuid.UUID) (*View, error) {
	credits, err := s.txn.Balance(ctx, economy.AccountAccount, account, Currency)
	if err != nil {
		return nil, err
	}
	owned, err := s.owned(ctx, account)
	if err != nil {
		return nil, err
	}
	items := Catalog()
	for i := range items {
		items[i].Owned = owned[items[i].Code]
	}
	return &View{Credits: credits, Items: items}, nil
}

// Buy spends credits on an item: a boost starts now for the character, a
// cosmetic becomes the account's. key makes a retried request buy once.
func (s *Service) Buy(ctx context.Context, account, char uuid.UUID, code, key string) error {
	it, ok := find(code)
	if !ok {
		return ErrUnknown
	}
	if it.Kind != "boost" {
		owned, err := s.owned(ctx, account)
		if err != nil {
			return err
		}
		if owned[code] {
			return ErrOwned
		}
	}
	ref := "shop-" + account.String() + "-" + key
	if _, err := s.txn.Write(ctx, economy.Intent{Kind: "shop_buy", RefKind: "shop", IdempotencyKey: ref,
		Entries: []economy.Entry{
			{Currency: Currency, AccountKind: economy.AccountAccount, AccountRef: account, Delta: -it.Price},
			{Currency: Currency, AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: it.Price},
		}}); err != nil {
		if errors.Is(err, economy.ErrInsufficientBalance) {
			return ErrNoCredit
		}
		return err
	}
	if it.Kind == "boost" {
		_, err := s.boost.Grant(ctx, char, strings.TrimPrefix(code, "boost_"), boost.SourceCredit, ref)
		if err != nil {
			s.refund(ctx, ref, account, it.Price)
		}
		return err
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO cosmetic.owned (account_id, code) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		account, code)
	return err
}

func (s *Service) refund(ctx context.Context, ref string, account uuid.UUID, credits int64) {
	_, _ = s.txn.Write(ctx, economy.Intent{Kind: "shop_refund", RefKind: "shop", IdempotencyKey: ref + "-refund",
		Entries: []economy.Entry{
			{Currency: Currency, AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: -credits},
			{Currency: Currency, AccountKind: economy.AccountAccount, AccountRef: account, Delta: credits},
		}})
}

// Profile is the account's chosen avatar and background.
type Profile struct {
	Avatar     string `json:"avatar"`
	Background string `json:"background"`
}

// Cosmetics is the profile plus every cosmetic the account may use.
type Cosmetics struct {
	Profile
	Owned []string `json:"owned"` // bought premium cosmetics (free ones are always usable)
}

func (s *Service) Cosmetics(ctx context.Context, account uuid.UUID) (*Cosmetics, error) {
	c := &Cosmetics{Profile: Profile{Avatar: AvatarCode(1), Background: BackgroundCode(1)}, Owned: []string{}}
	err := s.pool.QueryRow(ctx, `SELECT avatar, background FROM cosmetic.profile WHERE account_id = $1`, account).
		Scan(&c.Avatar, &c.Background)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	owned, err := s.owned(ctx, account)
	if err != nil {
		return nil, err
	}
	for code := range owned {
		c.Owned = append(c.Owned, code)
	}
	return c, nil
}

// Select puts on an avatar or a background the account may use.
func (s *Service) Select(ctx context.Context, account uuid.UUID, code string) (*Cosmetics, error) {
	kind, free, ok := cosmetic(code)
	if !ok {
		return nil, ErrUnknown
	}
	if !free {
		owned, err := s.owned(ctx, account)
		if err != nil {
			return nil, err
		}
		if !owned[code] {
			return nil, ErrLocked
		}
	}
	col := "avatar"
	if kind == "background" {
		col = "background"
	}
	if _, err := s.pool.Exec(ctx, fmt.Sprintf(`
		INSERT INTO cosmetic.profile (account_id, %[1]s) VALUES ($1, $2)
		ON CONFLICT (account_id) DO UPDATE SET %[1]s = EXCLUDED.%[1]s`, col), account, code); err != nil {
		return nil, err
	}
	return s.Cosmetics(ctx, account)
}

func (s *Service) owned(ctx context.Context, account uuid.UUID) (map[string]bool, error) {
	rows, err := s.pool.Query(ctx, `SELECT code FROM cosmetic.owned WHERE account_id = $1`, account)
	if err != nil {
		return nil, err
	}
	codes, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	out := map[string]bool{}
	for _, c := range codes {
		out[c] = true
	}
	return out, nil
}
