package notification

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MaxTokenFailures deactivates a device token after this many consecutive
// invalid-token responses (push-notification.md §3).
const MaxTokenFailures = 3

// Dispatch sends due push messages to every active device of the account.
// Each message is handled in its own short transaction holding the row
// lock, so parallel workers never send the same message twice.
func Dispatch(ctx context.Context, pool *pgxpool.Pool, senders map[string]Sender, batch int) (int64, error) {
	rows, err := pool.Query(ctx, `
		SELECT id FROM notification.push_message
		 WHERE sent_at IS NULL AND scheduled_for <= now()
		 ORDER BY scheduled_for LIMIT $1`, batch)
	if err != nil {
		return 0, err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return 0, err
	}
	var n int64
	for _, id := range ids {
		sent, err := dispatchOne(ctx, pool, senders, id)
		if err != nil {
			return n, err
		}
		if sent {
			n++
		}
	}
	return n, nil
}

type device struct {
	id       uuid.UUID
	platform string
	token    string
}

func dispatchOne(ctx context.Context, pool *pgxpool.Pool, senders map[string]Sender, msgID uuid.UUID) (bool, error) {
	handled := false
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		var account uuid.UUID
		var priority string
		var raw []byte
		err := tx.QueryRow(ctx, `
			SELECT account_id, priority, payload FROM notification.push_message
			 WHERE id = $1 AND sent_at IS NULL FOR UPDATE SKIP LOCKED`, msgID).Scan(&account, &priority, &raw)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		handled = true
		var payload Payload
		if err := json.Unmarshal(raw, &payload); err != nil {
			return markSent(ctx, tx, msgID, 0, "bad_payload")
		}
		rows, err := tx.Query(ctx, `
			SELECT id, platform, push_token FROM notification.device_push_token
			 WHERE account_id = $1 AND is_active`, account)
		if err != nil {
			return err
		}
		devices, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (device, error) {
			var d device
			err := r.Scan(&d.id, &d.platform, &d.token)
			return d, err
		})
		if err != nil {
			return err
		}
		if len(devices) == 0 {
			return markSent(ctx, tx, msgID, 0, "no_devices")
		}
		ok := 0
		reason := ""
		for _, d := range devices {
			sender := senders[d.platform]
			if sender == nil {
				reason = "no_sender"
				continue
			}
			err := sender.Send(ctx, Push{Token: d.token, Priority: priority, Payload: payload})
			switch {
			case err == nil:
				ok++
				if _, err := tx.Exec(ctx, `
					UPDATE notification.device_push_token SET failure_count = 0, last_used_at = now()
					 WHERE id = $1`, d.id); err != nil {
					return err
				}
			case errors.Is(err, ErrInvalidToken):
				reason = "invalid_token"
				if _, err := tx.Exec(ctx, `
					UPDATE notification.device_push_token
					   SET failure_count = failure_count + 1, last_failure_at = now(),
					       is_active = failure_count + 1 < $2
					 WHERE id = $1`, d.id, MaxTokenFailures); err != nil {
					return err
				}
			default:
				reason = "provider_error"
			}
		}
		if ok == 0 && reason == "" {
			reason = "no_delivery"
		}
		if ok > 0 {
			reason = ""
		}
		return markSent(ctx, tx, msgID, ok, reason)
	})
	return handled, err
}

func markSent(ctx context.Context, tx pgx.Tx, id uuid.UUID, ok int, reason string) error {
	_, err := tx.Exec(ctx, `
		UPDATE notification.push_message SET sent_at = now(), devices_ok = $2, fail_reason = NULLIF($3, '')
		 WHERE id = $1`, id, ok, reason)
	return err
}
