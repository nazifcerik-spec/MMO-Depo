// Package privacy serves the account-rights side of DD-231: the personal-data
// export, erasure after the deletion grace, and the legal pages.
package privacy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExportFormat versions the export document.
const ExportFormat = 1

const chars = `(SELECT id FROM character.character WHERE account_id = $1)`

// exportSections are the read-only queries behind an export, keyed by the
// section name in the document; $1 is the account id. Bounded where a table
// can grow without limit.
var exportSections = []struct{ name, sql string }{
	{"account", `SELECT id, email, username, created_at, email_verified_at, last_login_at,
	        terms_version, terms_accepted_at, deletion_requested_at FROM identity.account WHERE id = $1`},
	{"sessions", `SELECT created_at, last_active_at, expires_at, revoked_at, user_agent
	        FROM identity.session WHERE account_id = $1 ORDER BY created_at DESC LIMIT 200`},
	{"devices", `SELECT platform, locale, app_version, registered_at, last_used_at, is_active
	        FROM notification.device_push_token WHERE account_id = $1`},
	{"characters", `SELECT id, name, race_code, class_code, level, xp, state, location_region_code,
	        active_title_code, created_at, last_played_at FROM character.character WHERE account_id = $1`},
	{"wallets", `SELECT character_id, gold_copper FROM economy.character_wallet WHERE character_id IN ` + chars},
	{"character_currencies", `SELECT character_id, currency_code, amount
	        FROM economy.character_currency_balance WHERE character_id IN ` + chars},
	{"account_currencies", `SELECT currency_code, amount FROM economy.account_currency_balance WHERE account_id = $1`},
	{"items", `SELECT owner_character_id AS character_id, template_code, stack_quantity, rarity, prefix_code,
	        suffix_code, is_bound, slot, equipment_slot, created_at
	        FROM inventory.item_instance WHERE owner_character_id IN ` + chars},
	{"professions", `SELECT character_id, profession_code, level, xp, learned_at
	        FROM crafting.character_profession WHERE character_id IN ` + chars},
	{"quests", `SELECT character_id, quest_code, state, progress, accepted_at, completed_at
	        FROM progression.character_quest WHERE character_id IN ` + chars + ` ORDER BY accepted_at DESC LIMIT 2000`},
	{"achievements", `SELECT character_id, achievement_code, progress_current, unlocked_at
	        FROM progression.achievement_progress WHERE character_id IN ` + chars},
	{"titles", `SELECT character_id, title_code, unlocked_at FROM progression.character_title WHERE character_id IN ` + chars},
	{"reputation", `SELECT character_id, faction_code, rep_xp, current_tier
	        FROM progression.character_reputation WHERE character_id IN ` + chars},
	{"mastery", `SELECT character_id, track_code, xp, current_tier FROM progression.character_mastery WHERE character_id IN ` + chars},
	{"friends", `SELECT character_a, character_b, since FROM social.friendship
	        WHERE character_a IN ` + chars + ` OR character_b IN ` + chars},
	{"ignored", `SELECT owner_character_id AS character_id, target_character_id, added_at
	        FROM social.ignore_list WHERE owner_character_id IN ` + chars},
	{"guilds", `SELECT character_id, guild_id, rank_index, joined_at FROM social.guild_member WHERE character_id IN ` + chars},
	{"mail", `SELECT id, sender_character_id, recipient_character_id, subject, body, gold_attachment_copper,
	        item_attachment_count, state, created_at, read_at FROM social.mail
	        WHERE sender_character_id IN ` + chars + ` OR recipient_character_id IN ` + chars + `
	        ORDER BY created_at DESC LIMIT 1000`},
	{"chat", `SELECT sender_character_id AS character_id, channel_kind, target_character_id, body, sent_at
	        FROM social.chat_message WHERE sender_character_id IN ` + chars + ` ORDER BY sent_at DESC LIMIT 2000`},
	{"market", `SELECT id, seller_character_id, buyer_character_id, template_code, stack_quantity,
	        unit_price_copper, state, created_at, closed_at FROM market.listing
	        WHERE seller_character_id IN ` + chars + ` OR buyer_character_id IN ` + chars + `
	        ORDER BY created_at DESC LIMIT 2000`},
	{"ledger", `SELECT created_at, kind, currency_code, account_kind, account_ref, delta, balance_after, ref_kind
	        FROM economy.ledger_entry
	        WHERE (account_kind = 'character' AND account_ref IN ` + chars + `)
	           OR (account_kind = 'account' AND account_ref = $1)
	        ORDER BY created_at DESC LIMIT 5000`},
	{"notifications", `SELECT kind, body, created_at, read_at FROM notification.inbox_message
	        WHERE account_id = $1 ORDER BY created_at DESC LIMIT 500`},
	{"world_events", `SELECT spawn_id, character_id, registered_at, contribution_score, reward_tier
	        FROM worldevent.attendance WHERE account_id = $1`},
	{"support_tickets", `SELECT id, category, subject, body, state, reply, created_at, replied_at, closed_at
	        FROM ops.support_ticket WHERE account_id = $1 ORDER BY created_at DESC`},
}

// Export assembles everything stored about an account into one document
// (KVKK/GDPR access right). Read-only; one snapshot transaction.
func Export(ctx context.Context, pool *pgxpool.Pool, accountID uuid.UUID, now time.Time) (map[string]any, error) {
	doc := map[string]any{"format": ExportFormat, "generated_at": now.UTC()}
	err := pgx.BeginTxFunc(ctx, pool, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly},
		func(tx pgx.Tx) error {
			for _, s := range exportSections {
				var raw json.RawMessage
				if err := tx.QueryRow(ctx, `SELECT coalesce(json_agg(t), '[]'::json) FROM (`+s.sql+`) t`, accountID).
					Scan(&raw); err != nil {
					return fmt.Errorf("export %s: %w", s.name, err)
				}
				doc[s.name] = raw
			}
			return nil
		})
	return doc, err
}

// Purge erases accounts whose deletion grace has passed: characters go to
// the character-deletion job (which handles items and guild leadership),
// personal data is anonymised, devices, tokens, notifications and unsent
// mail are deleted. Economy and audit records stay, attached to the
// anonymous account (docs/database/auth-accounts.md). Returns how many.
func Purge(ctx context.Context, pool *pgxpool.Pool, now time.Time, grace time.Duration, batch int) (int64, error) {
	var n int64
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id, email FROM identity.account
			 WHERE deletion_requested_at <= $1 AND state <> 'deleted'
			 ORDER BY deletion_requested_at LIMIT $2 FOR UPDATE SKIP LOCKED`, now.Add(-grace), batch)
		if err != nil {
			return err
		}
		type due struct {
			id    uuid.UUID
			email string
		}
		var list []due
		for rows.Next() {
			var d due
			if err := rows.Scan(&d.id, &d.email); err != nil {
				rows.Close()
				return err
			}
			list = append(list, d)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		for _, d := range list {
			for _, q := range purgeStamped {
				if _, err := tx.Exec(ctx, q, d.id, now); err != nil {
					return fmt.Errorf("purge %s: %w", d.id, err)
				}
			}
			for _, q := range purgeDeleted {
				if _, err := tx.Exec(ctx, q, d.id); err != nil {
					return fmt.Errorf("purge %s: %w", d.id, err)
				}
			}
			if _, err := tx.Exec(ctx, `DELETE FROM ops.email_outbox WHERE to_address = $1 AND sent_at IS NULL`, d.email); err != nil {
				return err
			}
			n++
		}
		return nil
	})
	return n, err
}

// purgeStamped run with ($1 account id, $2 now); purgeDeleted with ($1).
var (
	purgeStamped = []string{
		`UPDATE character.character SET state = 'delete_grace', delete_grace_at = $2, updated_at = $2,
		        version = version + 1 WHERE account_id = $1 AND state <> 'delete_grace'`,
		`UPDATE character.character SET delete_grace_at = $2 WHERE account_id = $1 AND delete_grace_at > $2`,
		`UPDATE identity.refresh_token rt SET revoked_at = $2, revoke_reason = 'account_deleted'
		   FROM identity.session s WHERE rt.session_id = s.id AND s.account_id = $1 AND rt.revoked_at IS NULL`,
		`UPDATE identity.session SET revoked_at = $2, revoke_reason = 'account_deleted'
		  WHERE account_id = $1 AND revoked_at IS NULL`,
		`UPDATE identity.account
		    SET email = 'deleted+' || id || '@deleted.invalid', email_normalized = 'deleted+' || id || '@deleted.invalid',
		        username = 'deleted_' || left(replace(id::text, '-', ''), 16),
		        username_normalized = 'deleted_' || left(replace(id::text, '-', ''), 16),
		        password_hash = '!', state = 'deleted', updated_at = $2, version = version + 1
		  WHERE id = $1`,
	}
	purgeDeleted = []string{
		`DELETE FROM identity.account_token WHERE account_id = $1`,
		`DELETE FROM notification.device_push_token WHERE account_id = $1`,
		`DELETE FROM notification.push_message WHERE account_id = $1`,
		`DELETE FROM notification.inbox_message WHERE account_id = $1`,
		`DELETE FROM notification.preference_kind WHERE account_id = $1`,
		`DELETE FROM notification.preference WHERE account_id = $1`,
		`DELETE FROM ops.support_ticket WHERE account_id = $1`,
	}
)
