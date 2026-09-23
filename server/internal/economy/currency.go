package economy

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Currency mirrors the economy.currency catalog row.
type Currency struct {
	Code           string
	DisplayNameKey string
	Scope          string // "character" | "account" | "guild"
	Tradable       bool
	Cap            int64
	BaseUnitRatio  int32
}

// LoadCurrencies reads the full currency catalog from Postgres. Called once at
// boot; the map is process-cached (currencies change only via migration).
func LoadCurrencies(ctx context.Context, pool *pgxpool.Pool) (map[string]Currency, error) {
	const q = `SELECT code, display_name_key, scope, tradable, cap, base_unit_ratio
	             FROM economy.currency`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("load currencies: %w", err)
	}
	defer rows.Close()

	out := make(map[string]Currency)
	for rows.Next() {
		var c Currency
		if err := rows.Scan(&c.Code, &c.DisplayNameKey, &c.Scope, &c.Tradable, &c.Cap, &c.BaseUnitRatio); err != nil {
			return nil, fmt.Errorf("scan currency: %w", err)
		}
		out[c.Code] = c
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty currency catalog — did 0005_economy migration run?")
	}
	return out, nil
}

// StaticCurrencies returns the v1 baseline currencies for tests / boot without
// a database. Must stay in sync with 0005_economy.up.sql seed rows.
func StaticCurrencies() map[string]Currency {
	return map[string]Currency{
		"gold": {
			Code: "gold", DisplayNameKey: "currency.gold",
			Scope: "character", Tradable: true,
			Cap: 1_000_000_000_000, BaseUnitRatio: 10_000,
		},
		"gem": {
			Code: "gem", DisplayNameKey: "currency.gem",
			Scope: "account", Tradable: false,
			Cap: 10_000_000, BaseUnitRatio: 1,
		},
		"honor": {
			Code: "honor", DisplayNameKey: "currency.honor",
			Scope: "character", Tradable: false,
			Cap: 100_000, BaseUnitRatio: 1,
		},
		"guild_contribution": {
			Code: "guild_contribution", DisplayNameKey: "currency.gcontribution",
			Scope: "guild", Tradable: false,
			Cap: 1_000_000, BaseUnitRatio: 1,
		},
	}
}
