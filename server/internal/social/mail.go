package social

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/notification"
)

type MailSummary struct {
	ID         uuid.UUID  `json:"id"`
	SenderID   *uuid.UUID `json:"sender_character_id,omitempty"`
	SenderName string     `json:"sender_name"`
	Subject    string     `json:"subject"`
	GoldCopper int64      `json:"gold_attachment_copper"`
	ItemCount  int16      `json:"item_attachment_count"`
	State      string     `json:"state"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
}

type MailAttachment struct {
	SlotIndex     int16     `json:"slot_index"`
	ItemID        uuid.UUID `json:"item_instance_id"`
	TemplateCode  string    `json:"template_code"`
	Rarity        string    `json:"rarity"`
	StackQuantity int32     `json:"stack_quantity"`
}

type Mail struct {
	MailSummary
	Body        string           `json:"body"`
	Attachments []MailAttachment `json:"attachments"`
}

type SendRequest struct {
	RecipientName  string
	Subject        string
	Body           string
	ItemIDs        []uuid.UUID
	GoldCopper     int64
	IdempotencyKey string
}

type SendResult struct {
	MailID  uuid.UUID `json:"mail_id"`
	Postage int64     `json:"postage_copper"`
}

// SendMail escrows gold (ledger → mail_escrow) and items (locked, lock_owner
// 'mail') and charges postage to fee_sink. If the recipient blocked the
// sender the mail is silently dropped (docs/social/mail.md §5 v1 choice).
func (s *Service) SendMail(ctx context.Context, actor, account uuid.UUID, req SendRequest) (*SendResult, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	if req.IdempotencyKey == "" {
		return nil, ErrIdempotencyMissing
	}
	subject, err := CleanText(req.Subject, MaxSubject)
	if err != nil {
		return nil, err
	}
	body := req.Body
	if len([]rune(body)) > MaxMailBody {
		return nil, ErrTooLong
	}
	if req.GoldCopper < 0 {
		return nil, ErrAttachmentInvalid
	}
	if len(req.ItemIDs) > MaxAttachments {
		return nil, ErrTooManyAttachments
	}
	seen := map[uuid.UUID]bool{}
	for _, id := range req.ItemIDs {
		if seen[id] {
			return nil, ErrAttachmentInvalid
		}
		seen[id] = true
	}
	to, err := characterByName(ctx, s.pool, req.RecipientName)
	if err != nil {
		return nil, err
	}
	if to.ID == actor {
		return nil, ErrSelfTarget
	}
	var blockedBySender bool
	if err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM social.ignore_list WHERE owner_character_id = $1 AND target_character_id = $2)`,
		to.ID, actor).Scan(&blockedBySender); err != nil {
		return nil, err
	}
	if blockedBySender {
		return &SendResult{}, nil
	}

	ledgerKey := "mail-send-" + req.IdempotencyKey
	if id, ok := s.mailByLedgerKey(ctx, ledgerKey); ok {
		return &SendResult{MailID: id, Postage: Postage(len(req.ItemIDs), req.GoldCopper)}, nil
	}

	var inbox int
	if err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM social.mail
		 WHERE recipient_character_id = $1 AND deleted_at IS NULL AND state IN ('sent','read')`, to.ID).
		Scan(&inbox); err != nil {
		return nil, err
	}
	if inbox >= InboxCap {
		return nil, ErrInboxFull
	}
	if len(req.ItemIDs) > 0 {
		var ok int
		if err := s.pool.QueryRow(ctx, `
			SELECT count(*) FROM inventory.item_instance
			 WHERE id = ANY($1) AND owner_character_id = $2
			   AND is_locked = false AND is_bound = false AND equipment_slot IS NULL`,
			req.ItemIDs, actor).Scan(&ok); err != nil {
			return nil, err
		}
		if ok != len(req.ItemIDs) {
			return nil, ErrAttachmentInvalid
		}
	}

	postage := Postage(len(req.ItemIDs), req.GoldCopper)
	entries := []economy.Entry{
		{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: -(postage + req.GoldCopper)},
		{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemFeeSink, Delta: postage},
	}
	if req.GoldCopper > 0 {
		entries = append(entries, economy.Entry{
			Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMailEscrow, Delta: req.GoldCopper,
		})
	}
	res, err := s.txn.Write(ctx, economy.Intent{
		Kind: "mail_send", RefKind: "mail", IdempotencyKey: ledgerKey, Entries: entries,
	})
	if err != nil {
		return nil, err
	}

	var mailID uuid.UUID
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `
			INSERT INTO social.mail (sender_character_id, recipient_character_id, subject, body,
			                         gold_attachment_copper, item_attachment_count, postage_copper,
			                         send_txn_id, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			actor, to.ID, subject, body, req.GoldCopper, len(req.ItemIDs), postage,
			res.TxnID, s.now().UTC().Add(MailTTL)).Scan(&mailID); err != nil {
			return err
		}
		if len(req.ItemIDs) == 0 {
			return nil
		}
		tag, err := tx.Exec(ctx, `
			UPDATE inventory.item_instance
			   SET is_locked = true, lock_owner = 'mail', locked_ref_id = $3, slot = NULL,
			       updated_at = now(), version = version + 1
			 WHERE id = ANY($1) AND owner_character_id = $2
			   AND is_locked = false AND is_bound = false AND equipment_slot IS NULL`,
			req.ItemIDs, actor, mailID)
		if err != nil {
			return err
		}
		if int(tag.RowsAffected()) != len(req.ItemIDs) {
			return ErrAttachmentInvalid
		}
		for i, id := range req.ItemIDs {
			if _, err := tx.Exec(ctx, `
				INSERT INTO social.mail_attachment (mail_id, slot_index, item_instance_id) VALUES ($1, $2, $3)`,
				mailID, i, id); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if id, ok := s.mailByLedgerKey(ctx, ledgerKey); ok {
			return &SendResult{MailID: id, Postage: postage}, nil
		}
		s.reverse(ctx, "mail-send-refund-"+req.IdempotencyKey, entries)
		return nil, err
	}
	s.publish(to.Account, "mail.received", map[string]any{"mail_id": mailID, "subject": subject})
	notification.Send(ctx, s.notify, s.log, notification.Notice{
		AccountID: to.Account, CharacterID: &to.ID, Kind: "mail_received",
		DedupKey: "mail_received:" + mailID.String(),
		TitleKey: "notif.mail_received.title", BodyKey: "notif.mail_received.body",
		Vars:  map[string]any{"subject": subject, "has_attachments": req.GoldCopper > 0 || len(req.ItemIDs) > 0},
		Route: "mail/" + mailID.String(), Category: "social", High: req.GoldCopper > 0 || len(req.ItemIDs) > 0,
	})
	return &SendResult{MailID: mailID, Postage: postage}, nil
}

func (s *Service) mailByLedgerKey(ctx context.Context, key string) (uuid.UUID, bool) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `
		SELECT m.id FROM social.mail m
		 WHERE m.send_txn_id IN (SELECT txn_id FROM economy.ledger_entry WHERE idempotency_key = $1)`, key).
		Scan(&id)
	return id, err == nil
}

// reverse writes a compensating txn with every entry negated.
func (s *Service) reverse(ctx context.Context, key string, entries []economy.Entry) {
	rev := make([]economy.Entry, len(entries))
	for i, e := range entries {
		e.Delta = -e.Delta
		rev[i] = e
	}
	if _, err := s.txn.Write(ctx, economy.Intent{Kind: "mail_refund", RefKind: "mail", IdempotencyKey: key, Entries: rev}); err != nil {
		s.log.ErrorContext(ctx, "social: mail refund failed — manual reconcile needed", "key", key, "err", err.Error())
	}
}

func (s *Service) Inbox(ctx context.Context, actor, account uuid.UUID) ([]MailSummary, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT m.id, m.sender_character_id, COALESCE(c.name, 'System'), m.subject,
		       m.gold_attachment_copper, m.item_attachment_count, m.state, m.created_at, m.expires_at
		  FROM social.mail m LEFT JOIN character.character c ON c.id = m.sender_character_id
		 WHERE m.recipient_character_id = $1 AND m.deleted_at IS NULL
		 ORDER BY m.created_at DESC LIMIT $2`, actor, InboxCap)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanSummary)
}

func scanSummary(r pgx.CollectableRow) (MailSummary, error) {
	var m MailSummary
	err := r.Scan(&m.ID, &m.SenderID, &m.SenderName, &m.Subject, &m.GoldCopper, &m.ItemCount,
		&m.State, &m.CreatedAt, &m.ExpiresAt)
	return m, err
}

// ReadMail returns the full mail and marks it read.
func (s *Service) ReadMail(ctx context.Context, actor, account, mailID uuid.UUID) (*Mail, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var m Mail
	err := s.pool.QueryRow(ctx, `
		SELECT m.id, m.sender_character_id, COALESCE(c.name, 'System'), m.subject,
		       m.gold_attachment_copper, m.item_attachment_count, m.state, m.created_at, m.expires_at, m.body
		  FROM social.mail m LEFT JOIN character.character c ON c.id = m.sender_character_id
		 WHERE m.id = $1 AND m.recipient_character_id = $2 AND m.deleted_at IS NULL`, mailID, actor).
		Scan(&m.ID, &m.SenderID, &m.SenderName, &m.Subject, &m.GoldCopper, &m.ItemCount,
			&m.State, &m.CreatedAt, &m.ExpiresAt, &m.Body)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT a.slot_index, i.id, i.template_code, i.rarity, i.stack_quantity
		  FROM social.mail_attachment a JOIN inventory.item_instance i ON i.id = a.item_instance_id
		 WHERE a.mail_id = $1 AND i.locked_ref_id = $1
		 ORDER BY a.slot_index`, mailID)
	if err != nil {
		return nil, err
	}
	if m.Attachments, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (MailAttachment, error) {
		var a MailAttachment
		err := r.Scan(&a.SlotIndex, &a.ItemID, &a.TemplateCode, &a.Rarity, &a.StackQuantity)
		return a, err
	}); err != nil {
		return nil, err
	}
	if m.State == "sent" {
		if _, err := s.pool.Exec(ctx, `
			UPDATE social.mail SET state = 'read', read_at = now() WHERE id = $1 AND state = 'sent'`, mailID); err != nil {
			return nil, err
		}
		m.State = "read"
	}
	return &m, nil
}

// ClaimMail moves escrowed gold and items to the recipient. The mail row is
// locked for the whole operation so the expiry worker cannot return the
// same escrow concurrently.
func (s *Service) ClaimMail(ctx context.Context, actor, account, mailID uuid.UUID) (*Mail, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var claimed []item.Instance
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var gold int64
		var state string
		err := tx.QueryRow(ctx, `
			SELECT gold_attachment_copper, state FROM social.mail
			 WHERE id = $1 AND recipient_character_id = $2 AND deleted_at IS NULL FOR UPDATE`,
			mailID, actor).Scan(&gold, &state)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if state != "sent" && state != "read" {
			return ErrNothingToClaim
		}
		rows, err := tx.Query(ctx, `
			UPDATE inventory.item_instance
			   SET owner_character_id = $2, is_locked = false, lock_owner = NULL, locked_ref_id = NULL,
			       updated_at = now(), version = version + 1
			 WHERE locked_ref_id = $1 AND lock_owner = 'mail'
			RETURNING id, template_code, stack_quantity`, mailID, actor)
		if err != nil {
			return err
		}
		if claimed, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (item.Instance, error) {
			var in item.Instance
			err := r.Scan(&in.ID, &in.TemplateCode, &in.StackQuantity)
			return in, err
		}); err != nil {
			return err
		}
		if gold == 0 && len(claimed) == 0 {
			return ErrNothingToClaim
		}
		if gold > 0 {
			if _, err := s.txn.Write(ctx, economy.Intent{
				Kind: "mail_claim", RefKind: "mail", RefID: mailID,
				IdempotencyKey: "mail-claim-" + mailID.String(),
				Entries: []economy.Entry{
					{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMailEscrow, Delta: -gold},
					{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: gold},
				},
			}); err != nil {
				return err
			}
		}
		_, err = tx.Exec(ctx, `
			UPDATE social.mail SET state = 'claimed', attachments_claimed_at = now(),
			       read_at = COALESCE(read_at, now())
			 WHERE id = $1`, mailID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if len(claimed) > 0 && s.slots != nil {
		if err := s.slots.AssignSlots(ctx, actor, claimed); err != nil {
			s.log.WarnContext(ctx, "social: slot assign after claim failed", "err", err.Error())
		}
	}
	return s.ReadMail(ctx, actor, account, mailID)
}

// DeleteMail hides a mail that holds no unclaimed value.
func (s *Service) DeleteMail(ctx context.Context, actor, account, mailID uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	var gold int64
	var items int16
	var state string
	err := s.pool.QueryRow(ctx, `
		SELECT gold_attachment_copper, item_attachment_count, state FROM social.mail
		 WHERE id = $1 AND recipient_character_id = $2 AND deleted_at IS NULL`, mailID, actor).
		Scan(&gold, &items, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if (state == "sent" || state == "read") && (gold > 0 || items > 0) {
		return ErrHasAttachments
	}
	_, err = s.pool.Exec(ctx, `UPDATE social.mail SET deleted_at = now() WHERE id = $1`, mailID)
	return err
}

// ExpireMail is the worker entry point: mails past expires_at that still hold
// value go back to the sender (gold via ledger, items unlocked and re-slotted);
// value-less mails are marked expired. One short tx per mail. A mail whose
// claim already paid out in the ledger (claim tx died after the ledger write)
// is finished as a claim instead, so its items are never stranded.
func ExpireMail(ctx context.Context, pool *pgxpool.Pool, txn economy.TxnService, slots SlotAssigner, batch int) (int64, error) {
	rows, err := pool.Query(ctx, `
		SELECT id FROM social.mail
		 WHERE state IN ('sent','read') AND expires_at <= now()
		 ORDER BY expires_at LIMIT $1`, batch)
	if err != nil {
		return 0, err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return 0, err
	}
	var n int64
	for _, id := range ids {
		done, err := expireOne(ctx, pool, txn, slots, id)
		if err != nil {
			return n, err
		}
		if done {
			n++
		}
	}
	return n, nil
}

func expireOne(ctx context.Context, pool *pgxpool.Pool, txn economy.TxnService, slots SlotAssigner, mailID uuid.UUID) (bool, error) {
	done := false
	var moved []item.Instance
	var newOwner uuid.UUID
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		var sender *uuid.UUID
		var recipient uuid.UUID
		var gold int64
		var items int16
		err := tx.QueryRow(ctx, `
			SELECT sender_character_id, recipient_character_id, gold_attachment_copper, item_attachment_count
			  FROM social.mail
			 WHERE id = $1 AND state IN ('sent','read') AND expires_at <= now()
			   FOR UPDATE SKIP LOCKED`, mailID).Scan(&sender, &recipient, &gold, &items)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // claimed or locked by a claim in flight
		}
		if err != nil {
			return err
		}
		done = true
		if gold == 0 && items == 0 {
			_, err := tx.Exec(ctx, `UPDATE social.mail SET state = 'expired' WHERE id = $1`, mailID)
			return err
		}
		// A claim whose ledger write committed but whose mail update did not:
		// the recipient was paid, so complete the claim for the items too.
		var claimedInLedger bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM economy.ledger_entry WHERE idempotency_key = $1)`,
			"mail-claim-"+mailID.String()).Scan(&claimedInLedger); err != nil {
			return err
		}
		state := "returned"
		switch {
		case claimedInLedger:
			newOwner, state = recipient, "claimed"
		case sender == nil:
			_, err := tx.Exec(ctx, `UPDATE social.mail SET state = 'expired' WHERE id = $1`, mailID)
			return err
		default:
			newOwner = *sender
			if gold > 0 {
				if _, err := txn.Write(ctx, economy.Intent{
					Kind: "mail_return", RefKind: "mail", RefID: mailID,
					IdempotencyKey: "mail-return-" + mailID.String(),
					Entries: []economy.Entry{
						{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemMailEscrow, Delta: -gold},
						{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: *sender, Delta: gold},
					},
				}); err != nil {
					return err
				}
			}
		}
		rows, err := tx.Query(ctx, `
			UPDATE inventory.item_instance
			   SET owner_character_id = $2, is_locked = false, lock_owner = NULL, locked_ref_id = NULL,
			       updated_at = now(), version = version + 1
			 WHERE locked_ref_id = $1 AND lock_owner = 'mail'
			RETURNING id, template_code, stack_quantity`, mailID, newOwner)
		if err != nil {
			return err
		}
		if moved, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (item.Instance, error) {
			var in item.Instance
			err := r.Scan(&in.ID, &in.TemplateCode, &in.StackQuantity)
			return in, err
		}); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE social.mail SET state = $2::text,
			       attachments_claimed_at = CASE WHEN $2::text = 'claimed' THEN now() ELSE attachments_claimed_at END
			 WHERE id = $1`, mailID, state)
		return err
	})
	// Items travel with slot = NULL while in the mail; put them back in a bag
	// slot (best effort, as on claim — a full bag keeps them unslotted).
	if err == nil && len(moved) > 0 && slots != nil {
		_ = slots.AssignSlots(ctx, newOwner, moved)
	}
	return done, err
}
