//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// DD-358: treasury and bank through the ledger, rank limits, the log, and a
// guild that only disbands empty.
func TestHTTPGuildVault(t *testing.T) {
	w := newWorld(t)
	lead, mem := newHTTPPlayer(t, "warrior"), newHTTPPlayer(t, "mage")
	lead.mint(5_000_000)
	mem.mint(5_000_000)
	lead.ok("POST", lead.path("/guild"), map[string]any{"name": "Kasa " + randTag(), "tag": randTag()})
	inv := lead.ok("POST", lead.path("/guild/invite"), map[string]any{"target_name": mem.name})
	mem.ok("POST", mem.path("/guild/invites/"+inv["id"].(string)+"/accept"), map[string]any{})

	before := mem.gold()
	mem.fail("POST", mem.path("/guild/treasury/deposit"), map[string]any{"copper": 0}, 400, "E_VALIDATION")
	v := mem.ok("POST", mem.path("/guild/treasury/deposit"), map[string]any{"copper": 1_000_000})
	require.EqualValues(t, 1_000_000, v["treasury_copper"])
	require.Equal(t, before-1_000_000, mem.gold())

	// A new member (rank Initiate) gives but does not take.
	require.EqualValues(t, 4, v["rank_index"])
	require.Equal(t, false, v["can_use_items"])
	mem.fail("POST", mem.path("/guild/treasury/withdraw"), map[string]any{"copper": 1}, 403, "E_GUILD_FORBIDDEN")

	// A member's daily limit holds, and a failed withdrawal gives the allowance back.
	_, err := w.p.Exec(w.ctx, `UPDATE social.guild_member SET rank_index = 3 WHERE character_id = $1`, mem.char)
	require.NoError(t, err)
	mem.fail("POST", mem.path("/guild/treasury/withdraw"), map[string]any{"copper": 500_001}, 409, "E_GUILD_VAULT_LIMIT")
	v = mem.ok("POST", mem.path("/guild/treasury/withdraw"), map[string]any{"copper": 400_000})
	require.EqualValues(t, 600_000, v["treasury_copper"])
	require.EqualValues(t, 400_000, v["gold_today"])
	mem.fail("POST", mem.path("/guild/treasury/withdraw"), map[string]any{"copper": 100_001}, 409, "E_GUILD_VAULT_LIMIT")
	require.EqualValues(t, 400_000, mem.ok("GET", mem.path("/guild/bank"), nil)["gold_today"])

	// The leader has no limit but cannot take more than there is.
	lead.fail("POST", lead.path("/guild/treasury/withdraw"), map[string]any{"copper": 700_000}, 409, "E_INSUFFICIENT_GOLD")
	require.Zero(t, lead.ok("GET", lead.path("/guild/bank"), nil)["gold_today"])

	// Items go in and come out; bound ones stay.
	var memAccount uuid.UUID
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT account_id FROM character.character WHERE id = $1`, mem.char).Scan(&memAccount))
	item := w.give(player{Account: memAccount, ID: uuid.MustParse(mem.char), Name: mem.name}, "potion_hp_small", 5).ID
	v = mem.ok("POST", mem.path("/guild/bank/deposit"), map[string]any{"item_id": item.String(), "quantity": 2})
	bank := v["bank"].(map[string]any)
	require.EqualValues(t, 1, bank["slots_used"])
	first := bank["items"].([]any)[0].(map[string]any)
	require.EqualValues(t, 2, first["stack_quantity"])
	v = lead.ok("POST", lead.path("/guild/bank/withdraw"), map[string]any{"item_id": first["instance_id"], "quantity": 0})
	require.EqualValues(t, 0, v["bank"].(map[string]any)["slots_used"])
	lead.fail("POST", lead.path("/guild/bank/withdraw"), map[string]any{"item_id": first["instance_id"]}, 404, "E_NOT_FOUND")

	// Everything is logged, newest first.
	var actions []string
	for _, l := range v["log"].([]any) {
		actions = append(actions, l.(map[string]any)["action"].(string))
	}
	require.Equal(t, []string{"item_out", "item_in", "gold_out", "gold_in"}, actions)

	// Someone outside the guild sees nothing.
	stranger := newHTTPPlayer(t, "rogue")
	stranger.fail("GET", stranger.path("/guild/bank"), nil, 409, "E_NOT_IN_GUILD")

	// A guild disbands only when empty, and takes nothing in while it does.
	lead.fail("POST", lead.path("/guild/disband"), map[string]any{}, 409, "E_GUILD_VAULT_NOT_EMPTY")
	lead.ok("POST", lead.path("/guild/treasury/withdraw"), map[string]any{"copper": 600_000})
	lead.ok("POST", lead.path("/guild/disband"), map[string]any{})
	mem.fail("POST", mem.path("/guild/treasury/deposit"), map[string]any{"copper": 1}, 409, "E_GUILD_VAULT_CLOSING")

	// The ledger still balances.
	var bad int
	require.NoError(t, w.p.QueryRow(w.ctx, `SELECT count(*) FROM (SELECT txn_id FROM economy.ledger_entry
		WHERE kind LIKE 'guild_treasury_%' GROUP BY txn_id, currency_code HAVING sum(delta) <> 0) x`).Scan(&bad))
	require.Zero(t, bad)
}
