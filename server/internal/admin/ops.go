package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/otter/server/internal/economy"
)

// --- players ----------------------------------------------------------------

type Sanction struct {
	ID        uuid.UUID  `json:"id"`
	Kind      string     `json:"kind"`
	Until     *time.Time `json:"until"`
	Reason    string     `json:"reason"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type PlayerRow struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	State       string     `json:"state"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type CharacterRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Level     int16     `json:"level"`
	ClassCode string    `json:"class_code"`
	State     string    `json:"state"`
}

type PlayerDetail struct {
	PlayerRow
	Characters []CharacterRow `json:"characters"`
	Sanctions  []Sanction     `json:"sanctions"`
}

// MaskEmail keeps the first character of the local part: a***@example.com.
func MaskEmail(email string) string {
	at := strings.IndexByte(email, '@')
	if at <= 0 {
		return "***"
	}
	return email[:1] + "***" + email[at:]
}

func (s *Service) mask(a *Admin, p *PlayerRow) {
	if !a.Perms.Has(PermPlayerPIIReveal) {
		p.Email = MaskEmail(p.Email)
	}
}

// SearchPlayers matches username/email prefix or an exact account id.
func (s *Service) SearchPlayers(ctx context.Context, a *Admin, q string) ([]PlayerRow, error) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil, ErrInvalid
	}
	id, _ := uuid.Parse(q)
	rows, err := s.pool.Query(ctx, `
		SELECT id, username, email, state, created_at, last_login_at FROM identity.account
		 WHERE id = $1 OR username_normalized LIKE $2 OR email_normalized LIKE $2
		 ORDER BY username LIMIT 50`, id, strings.NewReplacer("%", `\%`, "_", `\_`).Replace(q)+"%")
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (PlayerRow, error) {
		var p PlayerRow
		err := r.Scan(&p.ID, &p.Username, &p.Email, &p.State, &p.CreatedAt, &p.LastLoginAt)
		return p, err
	})
	for i := range out {
		s.mask(a, &out[i])
	}
	s.Audit(ctx, a, AuditEntry{Action: "player.search", After: map[string]any{"q": q, "pii": a.Perms.Has(PermPlayerPIIReveal)}})
	return out, err
}

func (s *Service) Player(ctx context.Context, a *Admin, id uuid.UUID) (*PlayerDetail, error) {
	d := &PlayerDetail{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, username, email, state, created_at, last_login_at FROM identity.account WHERE id = $1`, id).
		Scan(&d.ID, &d.Username, &d.Email, &d.State, &d.CreatedAt, &d.LastLoginAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	s.mask(a, &d.PlayerRow)
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, level, class_code, state FROM character.character WHERE account_id = $1 ORDER BY created_at`, id)
	if err != nil {
		return nil, err
	}
	if d.Characters, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (CharacterRow, error) {
		var c CharacterRow
		err := r.Scan(&c.ID, &c.Name, &c.Level, &c.ClassCode, &c.State)
		return c, err
	}); err != nil {
		return nil, err
	}
	rows, err = s.pool.Query(ctx, `
		SELECT id, kind, until, reason, created_at, revoked_at FROM ops.account_sanction
		 WHERE account_id = $1 ORDER BY created_at DESC LIMIT 50`, id)
	if err != nil {
		return nil, err
	}
	d.Sanctions, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (Sanction, error) {
		var x Sanction
		err := r.Scan(&x.ID, &x.Kind, &x.Until, &x.Reason, &x.CreatedAt, &x.RevokedAt)
		return x, err
	})
	s.Audit(ctx, a, AuditEntry{Action: "player.view", TargetKind: "account", TargetID: &id})
	return d, err
}

// Sanction mutes or bans an account. duration 0 = permanent. Bans lock the
// account (temporary) or ban it (permanent) and revoke all player sessions.
func (s *Service) Sanction(ctx context.Context, a *Admin, account uuid.UUID, kind string, duration time.Duration, reason string) (*Sanction, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, ErrReasonRequired
	}
	var until *time.Time
	if duration > 0 {
		t := s.now().UTC().Add(duration)
		until = &t
	}
	out := &Sanction{Kind: kind, Until: until, Reason: reason}
	var before string
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `SELECT state FROM identity.account WHERE id = $1 FOR UPDATE`, account).Scan(&before); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotFound
			}
			return err
		}
		if err := tx.QueryRow(ctx, `
			INSERT INTO ops.account_sanction (account_id, kind, until, reason, created_by_admin_id)
			VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`, account, kind, until, reason, a.ID).
			Scan(&out.ID, &out.CreatedAt); err != nil {
			return err
		}
		if kind != "ban" {
			return nil
		}
		state := "locked"
		if until == nil {
			state = "banned"
		}
		if before == "banned" {
			state = "banned" // never downgrade a permanent ban
		}
		if _, err := tx.Exec(ctx, `UPDATE identity.account SET state = $2, updated_at = now(), version = version + 1 WHERE id = $1`,
			account, state); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			UPDATE identity.session SET revoked_at = now(), revoke_reason = 'admin_ban'
			 WHERE account_id = $1 AND revoked_at IS NULL`, account)
		return err
	})
	if err != nil {
		return nil, err
	}
	s.Audit(ctx, a, AuditEntry{Action: "player." + kind, TargetKind: "account", TargetID: &account,
		Before: map[string]any{"state": before}, After: out, Reason: reason, StepUp: kind == "ban" && until == nil})
	return out, nil
}

// LiftSanctions revokes active sanctions of a kind and reactivates a
// locked/banned account when lifting bans.
func (s *Service) LiftSanctions(ctx context.Context, a *Admin, account uuid.UUID, kind, reason string) (int64, error) {
	if strings.TrimSpace(reason) == "" {
		return 0, ErrReasonRequired
	}
	var n int64
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE ops.account_sanction SET revoked_at = now(), revoked_by_admin_id = $3
			 WHERE account_id = $1 AND kind = $2 AND revoked_at IS NULL`, account, kind, a.ID)
		if err != nil {
			return err
		}
		n = tag.RowsAffected()
		if kind == "ban" {
			_, err = tx.Exec(ctx, `
				UPDATE identity.account SET state = 'active', updated_at = now(), version = version + 1
				 WHERE id = $1 AND state IN ('locked','banned')`, account)
		}
		return err
	})
	if err != nil {
		return 0, err
	}
	s.Audit(ctx, a, AuditEntry{Action: "player.un" + kind, TargetKind: "account", TargetID: &account,
		After: map[string]any{"revoked": n}, Reason: reason})
	return n, nil
}

// --- characters -------------------------------------------------------------

type CharacterDetail struct {
	CharacterRow
	AccountID       uuid.UUID `json:"account_id"`
	XP              int64     `json:"xp"`
	GoldCopper      int64     `json:"gold_copper"`
	ItemCount       int       `json:"item_count"`
	ActiveContracts int       `json:"active_contracts"`
	Location        *string   `json:"location_region_code,omitempty"`
}

func (s *Service) Character(ctx context.Context, a *Admin, id uuid.UUID) (*CharacterDetail, error) {
	d := &CharacterDetail{}
	err := s.pool.QueryRow(ctx, `
		SELECT c.id, c.name, c.level, c.class_code, c.state, c.account_id, c.xp, c.location_region_code,
		       COALESCE((SELECT gold_copper FROM economy.character_wallet w WHERE w.character_id = c.id), 0),
		       (SELECT count(*) FROM inventory.item_instance i WHERE i.owner_character_id = c.id),
		       (SELECT count(*) FROM activity.contract x WHERE x.character_id = c.id AND x.state = 'active')
		     + (SELECT count(*) FROM hunt.contract x WHERE x.character_id = c.id AND x.state = 'active')
		     + (SELECT count(*) FROM dungeon.contract x WHERE x.character_id = c.id AND x.state = 'active')
		  FROM character.character c WHERE c.id = $1`, id).
		Scan(&d.ID, &d.Name, &d.Level, &d.ClassCode, &d.State, &d.AccountID, &d.XP, &d.Location,
			&d.GoldCopper, &d.ItemCount, &d.ActiveContracts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err == nil {
		s.Audit(ctx, a, AuditEntry{Action: "character.view", TargetKind: "character", TargetID: &id})
	}
	return d, err
}

// SetFrozen freezes (active → frozen) or unfreezes (frozen → active).
func (s *Service) SetFrozen(ctx context.Context, a *Admin, id uuid.UUID, freeze bool, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrReasonRequired
	}
	from, to := "active", "frozen"
	if !freeze {
		from, to = "frozen", "active"
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE character.character SET state = $3, updated_at = now(), version = version + 1
		 WHERE id = $1 AND state = $2`, id, from, to)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInvalid
	}
	s.Audit(ctx, a, AuditEntry{Action: "character." + map[bool]string{true: "freeze", false: "unfreeze"}[freeze],
		TargetKind: "character", TargetID: &id, Before: map[string]any{"state": from}, After: map[string]any{"state": to}, Reason: reason})
	return nil
}

// --- ledger -----------------------------------------------------------------

type LedgerEntry struct {
	TxnID          uuid.UUID  `json:"txn_id"`
	Seq            int16      `json:"entry_seq"`
	Kind           string     `json:"kind"`
	Currency       string     `json:"currency"`
	AccountKind    string     `json:"account_kind"`
	AccountRef     uuid.UUID  `json:"account_ref"`
	Delta          int64      `json:"delta"`
	BalanceAfter   int64      `json:"balance_after"`
	RefKind        *string    `json:"ref_kind,omitempty"`
	RefID          *uuid.UUID `json:"ref_id,omitempty"`
	IdempotencyKey string     `json:"idempotency_key"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (s *Service) Ledger(ctx context.Context, account, txn *uuid.UUID, limit int) ([]LedgerEntry, error) {
	if account == nil && txn == nil {
		return nil, ErrInvalid
	}
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := s.pool.Query(ctx, `
		SELECT txn_id, entry_seq, kind, currency_code, account_kind, account_ref, delta, balance_after,
		       ref_kind, ref_id, idempotency_key, created_at
		  FROM economy.ledger_entry
		 WHERE ($1::uuid IS NULL OR account_ref = $1) AND ($2::uuid IS NULL OR txn_id = $2)
		 ORDER BY created_at DESC, entry_seq LIMIT $3`, account, txn, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (LedgerEntry, error) {
		var e LedgerEntry
		err := r.Scan(&e.TxnID, &e.Seq, &e.Kind, &e.Currency, &e.AccountKind, &e.AccountRef, &e.Delta,
			&e.BalanceAfter, &e.RefKind, &e.RefID, &e.IdempotencyKey, &e.CreatedAt)
		return e, err
	})
}

// GoldMoved is the gross gold (copper) a txn credited — the size used for
// the reversal permission limit.
func GoldMoved(entries []LedgerEntry) int64 {
	var n int64
	for _, e := range entries {
		if e.Currency == "gold" && e.Delta > 0 {
			n += e.Delta
		}
	}
	return n
}

// Reverse writes the negation of a ledger txn (docs/economy/transaction.md
// §5). Idempotent per original txn. Fails with insufficient balance when the
// credited party already spent the funds.
func (s *Service) Reverse(ctx context.Context, a *Admin, txnID uuid.UUID, reason string) (uuid.UUID, error) {
	if strings.TrimSpace(reason) == "" {
		return uuid.Nil, ErrReasonRequired
	}
	orig, err := s.Ledger(ctx, nil, &txnID, 500)
	if err != nil {
		return uuid.Nil, err
	}
	if len(orig) == 0 {
		return uuid.Nil, ErrNotFound
	}
	if strings.HasPrefix(orig[0].IdempotencyKey, "admin-reversal-") {
		return uuid.Nil, fmt.Errorf("%w: cannot reverse a reversal", ErrInvalid)
	}
	if !a.Perms.CanReverse(GoldMoved(orig)) {
		return uuid.Nil, ErrForbidden
	}
	entries := make([]economy.Entry, 0, len(orig))
	for _, e := range orig {
		entries = append(entries, economy.Entry{
			Currency: e.Currency, AccountKind: economy.AccountKind(e.AccountKind), AccountRef: e.AccountRef, Delta: -e.Delta,
		})
	}
	res, err := s.txn.Write(ctx, economy.Intent{
		Kind: "admin_reversal", RefKind: "ledger_txn", RefID: txnID,
		IdempotencyKey: "admin-reversal-" + txnID.String(), Entries: entries,
	})
	if err != nil {
		return uuid.Nil, err
	}
	s.Audit(ctx, a, AuditEntry{Action: "ledger.reversal", TargetKind: "ledger_txn", TargetID: &txnID,
		Before: orig, After: map[string]any{"reversal_txn_id": res.TxnID}, Reason: reason, StepUp: true})
	return res.TxnID, nil
}

// --- guild / world events ---------------------------------------------------

// DisbandGuild immediately disbands a guild (override of the 48h grace).
func (s *Service) DisbandGuild(ctx context.Context, a *Admin, guildID uuid.UUID, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrReasonRequired
	}
	var name string
	var members int
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			UPDATE social.guild SET disbanded_at = now(), version = version + 1
			 WHERE id = $1 AND disbanded_at IS NULL RETURNING name`, guildID).Scan(&name)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `DELETE FROM social.guild_member WHERE guild_id = $1`, guildID)
		if err != nil {
			return err
		}
		members = int(tag.RowsAffected())
		_, err = tx.Exec(ctx, `
			UPDATE social.guild_invite SET status = 'expired' WHERE guild_id = $1 AND status = 'pending'`, guildID)
		return err
	})
	if err != nil {
		return err
	}
	s.Audit(ctx, a, AuditEntry{Action: "guild.disband_override", TargetKind: "guild", TargetID: &guildID,
		Before: map[string]any{"name": name, "members": members}, Reason: reason, StepUp: true})
	return nil
}

// TriggerWorldEvent schedules an extra spawn of a template.
func (s *Service) TriggerWorldEvent(ctx context.Context, a *Admin, template string, startsIn time.Duration, reason string) (uuid.UUID, error) {
	if strings.TrimSpace(reason) == "" {
		return uuid.Nil, ErrReasonRequired
	}
	if startsIn < 0 || startsIn > 24*time.Hour {
		return uuid.Nil, ErrInvalid
	}
	start := s.now().UTC().Add(startsIn).Truncate(time.Minute)
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `
		INSERT INTO worldevent.spawn (template_code, starts_at, ends_at)
		SELECT code, $2, $2 + make_interval(secs => duration_seconds)
		  FROM content.world_event_template WHERE code = $1
		RETURNING id`, template, start).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrNotFound
	}
	if err != nil {
		return uuid.Nil, err
	}
	s.Audit(ctx, a, AuditEntry{Action: "content.event_trigger", TargetKind: "world_event_spawn", TargetID: &id,
		After: map[string]any{"template": template, "starts_at": start}, Reason: reason})
	return id, nil
}
