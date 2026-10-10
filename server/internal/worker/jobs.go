package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/auth"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/inventory"
	"github.com/otter/server/internal/mail"
	"github.com/otter/server/internal/market"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/privacy"
	"github.com/otter/server/internal/social"
	"github.com/otter/server/internal/worldevent"
)

// BatchSize bounds every job's per-run work so one tick never holds locks
// for long; the backlog drains over successive ticks.
const BatchSize = 500

// sqlJob runs a single statement taking the batch size as $1. When
// returnsCount is set the statement yields one bigint row (rows touched).
type sqlJob struct {
	name         string
	interval     time.Duration
	pool         *pgxpool.Pool
	query        string
	returnsCount bool
}

func (j *sqlJob) Name() string            { return j.name }
func (j *sqlJob) Interval() time.Duration { return j.interval }
func (j *sqlJob) Run(ctx context.Context) (int64, error) {
	if j.returnsCount {
		var n int64
		err := j.pool.QueryRow(ctx, j.query, BatchSize).Scan(&n)
		return n, err
	}
	tag, err := j.pool.Exec(ctx, j.query, BatchSize)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// Notifier queues player notifications. notification.Service satisfies.
type Notifier interface {
	Notify(ctx context.Context, n notification.Notice) error
}

// tagged is one row returned by a tagging statement: (ref id, account, character, extra).
type tagged struct {
	ref, account, character uuid.UUID
	extra                   string
}

func collectTagged(rows pgx.Rows) ([]tagged, error) {
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (tagged, error) {
		var t tagged
		err := r.Scan(&t.ref, &t.account, &t.character, &t.extra)
		return t, err
	})
}

func notifyAll(ctx context.Context, n Notifier, log *slog.Logger, rows []tagged, build func(tagged) notification.Notice) {
	if n == nil {
		return
	}
	for _, t := range rows {
		if err := n.Notify(ctx, build(t)); err != nil && log != nil {
			log.WarnContext(ctx, "worker: notify failed", "ref", t.ref.String(), "err", err.Error())
		}
	}
}

// ActivityExpire tags elapsed AFK contracts (activity, hunt, dungeon) with
// ready_at and sends the "ready to reveal" notification once per contract.
type ActivityExpire struct {
	Pool     *pgxpool.Pool
	Notifier Notifier
	Log      *slog.Logger
}

func (ActivityExpire) Name() string            { return "activity_expire" }
func (ActivityExpire) Interval() time.Duration { return 30 * time.Second }

func (j ActivityExpire) Run(ctx context.Context) (int64, error) {
	var total int64
	for _, src := range []struct{ table, kind, route string }{
		{"activity.contract", "activity_completed", "activity"},
		{"hunt.contract", "activity_completed", "hunt"},
		{"dungeon.contract", "dungeon_completed", "dungeon"},
	} {
		rows, err := j.Pool.Query(ctx, `
			UPDATE `+src.table+` t SET ready_at = now()
			  FROM character.character c
			 WHERE c.id = t.character_id
			   AND t.id IN (SELECT id FROM `+src.table+`
			                 WHERE state = 'active' AND ready_at IS NULL AND expires_at <= now()
			                 ORDER BY expires_at LIMIT $1 FOR UPDATE SKIP LOCKED)
			RETURNING t.id, c.account_id, t.character_id, ''`, BatchSize)
		if err != nil {
			return total, err
		}
		tags, err := collectTagged(rows)
		if err != nil {
			return total, err
		}
		total += int64(len(tags))
		notifyAll(ctx, j.Notifier, j.Log, tags, func(t tagged) notification.Notice {
			ch := t.character
			return notification.Notice{
				AccountID: t.account, CharacterID: &ch, Kind: src.kind,
				DedupKey: src.kind + ":" + t.ref.String(),
				TitleKey: "notif." + src.kind + ".title", BodyKey: "notif." + src.kind + ".body",
				Route: "character/" + src.route + "/" + t.ref.String(), Category: "activity",
				CollapseKey: src.route + "_" + ch.String(),
			}
		})
	}
	return total, nil
}

// TravelArrive completes travels whose arrival time has passed and moves the
// character in the same statement, so a player who never taps "Complete"
// still lands at the destination.
type TravelArrive struct {
	Pool     *pgxpool.Pool
	Notifier Notifier
	Log      *slog.Logger
}

func (TravelArrive) Name() string            { return "travel_arrive" }
func (TravelArrive) Interval() time.Duration { return 30 * time.Second }

func (j TravelArrive) Run(ctx context.Context) (int64, error) {
	rows, err := j.Pool.Query(ctx, `
WITH done AS (
  UPDATE world.travel
     SET state = 'completed', completed_at = now(), updated_at = now(), version = version + 1
   WHERE id IN (SELECT id FROM world.travel
                 WHERE state = 'in_progress' AND arrives_at <= now()
                 ORDER BY arrives_at LIMIT $1 FOR UPDATE SKIP LOCKED)
  RETURNING id, character_id, to_region_code)
UPDATE character.character c
   SET location_region_code = done.to_region_code, updated_at = now(), version = c.version + 1
  FROM done
 WHERE c.id = done.character_id
RETURNING done.id, c.account_id, c.id, done.to_region_code`, BatchSize)
	if err != nil {
		return 0, err
	}
	tags, err := collectTagged(rows)
	if err != nil {
		return 0, err
	}
	notifyAll(ctx, j.Notifier, j.Log, tags, func(t tagged) notification.Notice {
		ch := t.character
		return notification.Notice{
			AccountID: t.account, CharacterID: &ch, Kind: "travel_completed",
			DedupKey: "travel_completed:" + t.ref.String(),
			TitleKey: "notif.travel_completed.title", BodyKey: "notif.travel_completed.body",
			Vars: map[string]any{"region": t.extra}, Route: "character/location",
			Category: "activity", CollapseKey: "travel_" + ch.String(),
		}
	})
	return int64(len(tags)), nil
}

// PushDispatch delivers queued pushes through FCM/APNs.
type PushDispatch struct {
	Pool    *pgxpool.Pool
	Senders map[string]notification.Sender
}

func (PushDispatch) Name() string            { return "push_dispatch" }
func (PushDispatch) Interval() time.Duration { return 10 * time.Second }
func (j PushDispatch) Run(ctx context.Context) (int64, error) {
	return notification.Dispatch(ctx, j.Pool, j.Senders, 100)
}

// EmailDispatch delivers the transactional e-mail outbox (DD-219).
type EmailDispatch struct {
	Pool   *pgxpool.Pool
	Sender mail.Sender
}

func (EmailDispatch) Name() string            { return "email_dispatch" }
func (EmailDispatch) Interval() time.Duration { return 15 * time.Second }
func (j EmailDispatch) Run(ctx context.Context) (int64, error) {
	if j.Sender == nil {
		return 0, nil
	}
	return mail.Dispatch(ctx, j.Pool, j.Sender, 50)
}

// AccountPurge erases accounts whose deletion grace has passed (DD-231);
// their characters then go through CharacterDeleteExpire.
type AccountPurge struct{ Pool *pgxpool.Pool }

func (AccountPurge) Name() string            { return "account_purge" }
func (AccountPurge) Interval() time.Duration { return time.Hour }
func (j AccountPurge) Run(ctx context.Context) (int64, error) {
	return privacy.Purge(ctx, j.Pool, time.Now(), auth.DeletionGrace, BatchSize)
}

// WorldEvents schedules, opens, closes and pays world boss / invasion spawns.
type WorldEvents struct {
	Pool     *pgxpool.Pool
	Txn      economy.TxnService
	Notifier notification.Notifier
	Log      *slog.Logger
}

func (WorldEvents) Name() string            { return "world_event_tick" }
func (WorldEvents) Interval() time.Duration { return time.Minute }
func (j WorldEvents) Run(ctx context.Context) (int64, error) {
	log := j.Log
	if log == nil {
		log = slog.Default()
	}
	return worldevent.Tick(ctx, j.Pool, j.Txn, j.Notifier, log)
}

// SanctionExpire closes elapsed mutes/bans and reactivates accounts whose
// last temporary ban ran out (Sprint 26 admin sanctions).
func SanctionExpire(pool *pgxpool.Pool) Job {
	return &sqlJob{name: "sanction_expire", interval: time.Minute, pool: pool, returnsCount: true, query: `
WITH exp AS (
  UPDATE ops.account_sanction SET revoked_at = now()
   WHERE id IN (SELECT id FROM ops.account_sanction
                 WHERE revoked_at IS NULL AND until IS NOT NULL AND until <= now() LIMIT $1)
  RETURNING account_id, kind),
act AS (
  UPDATE auth.account a SET state = 'active', updated_at = now(), version = a.version + 1
    FROM exp
   WHERE exp.kind = 'ban' AND a.id = exp.account_id AND a.state = 'locked'
     AND NOT EXISTS (SELECT 1 FROM ops.account_sanction s
                      WHERE s.account_id = a.id AND s.kind = 'ban' AND s.revoked_at IS NULL
                        AND (s.until IS NULL OR s.until > now()))
  RETURNING 1)
SELECT count(*) FROM exp`}
}

// NotificationGC drops expired inbox rows and sent pushes older than 7 days.
func NotificationGC(pool *pgxpool.Pool) Job {
	return &sqlJob{name: "notification_gc", interval: time.Hour, pool: pool, returnsCount: true, query: `
WITH i AS (DELETE FROM notification.inbox_message
            WHERE id IN (SELECT id FROM notification.inbox_message WHERE expires_at < now() LIMIT $1)
           RETURNING 1),
p AS (DELETE FROM notification.push_message
       WHERE id IN (SELECT id FROM notification.push_message
                     WHERE sent_at < now() - INTERVAL '7 days' LIMIT $1)
      RETURNING 1)
SELECT (SELECT count(*) FROM i) + (SELECT count(*) FROM p)`}
}

// IdempotencyGC removes idempotency records past their 24h TTL (DD-188).
func IdempotencyGC(pool *pgxpool.Pool) Job {
	return &sqlJob{name: "idempotency_gc", interval: 10 * time.Minute, pool: pool, query: `
DELETE FROM api.idempotency_record
 WHERE (key, endpoint, account_id) IN (
   SELECT key, endpoint, account_id FROM api.idempotency_record
    WHERE expires_at < now() LIMIT $1)`}
}

// BuybackGC drops vendor buybacks that can no longer be claimed (1 day after
// expiry) and claimed ones older than 30 days.
func BuybackGC(pool *pgxpool.Pool) Job {
	return &sqlJob{name: "buyback_gc", interval: time.Hour, pool: pool, query: `
DELETE FROM economy.vendor_buyback
 WHERE id IN (
   SELECT id FROM economy.vendor_buyback
    WHERE (claimed_at IS NULL AND expires_at < now() - INTERVAL '1 day')
       OR claimed_at < now() - INTERVAL '30 days'
    LIMIT $1)`}
}

// GuildDisbandExpire finalises guild disbands whose 48h grace has passed
// (docs/guild/guild.md §4): mark disbanded, drop members, expire invites.
func GuildDisbandExpire(pool *pgxpool.Pool) Job {
	return &sqlJob{name: "guild_disband_expire", interval: 5 * time.Minute, pool: pool, returnsCount: true, query: `
WITH g AS (
  UPDATE social.guild SET disbanded_at = now(), version = version + 1
   WHERE id IN (SELECT id FROM social.guild
                 WHERE disband_scheduled_at <= now() AND disbanded_at IS NULL
                 ORDER BY disband_scheduled_at LIMIT $1 FOR UPDATE SKIP LOCKED)
  RETURNING id),
m AS (DELETE FROM social.guild_member WHERE guild_id IN (SELECT id FROM g)),
i AS (UPDATE social.guild_invite SET status = 'expired'
       WHERE status = 'pending' AND guild_id IN (SELECT id FROM g))
SELECT count(*) FROM g`}
}

// --- character_delete_expire ----------------------------------------------

// CharacterDeleteExpire hard-deletes characters whose delete grace has
// passed (DD-134). Their items are deleted first; every other character-owned
// table cascades. Ledger rows stay (no FK) for audit.
type CharacterDeleteExpire struct{ Pool *pgxpool.Pool }

func (CharacterDeleteExpire) Name() string            { return "character_delete_expire" }
func (CharacterDeleteExpire) Interval() time.Duration { return 5 * time.Minute }

func (j CharacterDeleteExpire) Run(ctx context.Context) (int64, error) {
	var n int64
	err := pgx.BeginFunc(ctx, j.Pool, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id FROM character.character
			 WHERE state = 'delete_grace' AND delete_grace_at <= now()
			 ORDER BY delete_grace_at LIMIT $1 FOR UPDATE SKIP LOCKED`, BatchSize)
		if err != nil {
			return err
		}
		ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
		if err != nil || len(ids) == 0 {
			return err
		}
		// Guilds led by a doomed character pass to the most senior, earliest
		// remaining member; guilds with no one else left are removed.
		if _, err := tx.Exec(ctx, `
			UPDATE social.guild g SET leader_character_id = s.character_id, version = g.version + 1
			  FROM (SELECT DISTINCT ON (m.guild_id) m.guild_id, m.character_id
			          FROM social.guild_member m
			          JOIN social.guild g2 ON g2.id = m.guild_id
			         WHERE g2.leader_character_id = ANY($1) AND NOT (m.character_id = ANY($1))
			         ORDER BY m.guild_id, m.rank_index, m.joined_at) s
			 WHERE g.id = s.guild_id`, ids); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			UPDATE social.guild_member m SET rank_index = 0
			  FROM social.guild g
			 WHERE g.id = m.guild_id AND g.leader_character_id = m.character_id AND m.rank_index <> 0`); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`DELETE FROM social.guild WHERE leader_character_id = ANY($1)`, ids); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`DELETE FROM inventory.item_instance WHERE owner_character_id = ANY($1)`, ids); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `DELETE FROM character.character WHERE id = ANY($1)`, ids)
		n = tag.RowsAffected()
		return err
	})
	return n, err
}

// --- economy_reconciler -----------------------------------------------------

// EconomyReconciler compares each character's ledger sum per currency with
// the wallet projection (docs/economy/transaction.md §6). Any drift is logged
// at error level and, when Freeze is set, the character is frozen so no
// further value moves until ops investigates. It also flags recent ledger
// txns whose entries do not sum to zero.
type EconomyReconciler struct {
	Pool   *pgxpool.Pool
	Freeze bool
	Log    *slog.Logger
}

func (EconomyReconciler) Name() string            { return "economy_reconciler" }
func (EconomyReconciler) Interval() time.Duration { return 24 * time.Hour }

const driftQuery = `
WITH ledger AS (
  SELECT account_ref AS character_id, currency_code, SUM(delta) AS total
    FROM economy.ledger_entry
   WHERE account_kind = 'character'
   GROUP BY 1, 2),
proj AS (
  SELECT character_id, 'gold'::varchar AS currency_code, gold_copper AS amount
    FROM economy.character_wallet
  UNION ALL
  SELECT character_id, currency_code, amount FROM economy.character_currency_balance)
SELECT c.id, COALESCE(l.currency_code, p.currency_code), COALESCE(l.total, 0), COALESCE(p.amount, 0)
  FROM ledger l
  FULL OUTER JOIN proj p ON p.character_id = l.character_id AND p.currency_code = l.currency_code
  JOIN character.character c ON c.id = COALESCE(l.character_id, p.character_id)
 WHERE COALESCE(l.total, 0) <> COALESCE(p.amount, 0)`

const unbalancedQuery = `
SELECT txn_id, currency_code, SUM(delta)
  FROM economy.ledger_entry
 WHERE created_at > now() - INTERVAL '48 hours'
 GROUP BY txn_id, currency_code
HAVING SUM(delta) <> 0`

func (j EconomyReconciler) Run(ctx context.Context) (int64, error) {
	log := j.Log
	if log == nil {
		log = slog.Default()
	}
	rows, err := j.Pool.Query(ctx, driftQuery)
	if err != nil {
		return 0, err
	}
	type drift struct {
		id             uuid.UUID
		currency       string
		ledger, wallet int64
	}
	drifts, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (drift, error) {
		var d drift
		err := r.Scan(&d.id, &d.currency, &d.ledger, &d.wallet)
		return d, err
	})
	if err != nil {
		return 0, err
	}

	var frozen []uuid.UUID
	for _, d := range drifts {
		log.ErrorContext(ctx, "economy drift",
			"character_id", d.id.String(), "currency", d.currency,
			"ledger_sum", d.ledger, "projection", d.wallet, "drift", d.wallet-d.ledger)
		frozen = append(frozen, d.id)
	}
	if j.Freeze && len(frozen) > 0 {
		tag, err := j.Pool.Exec(ctx, `
			UPDATE character.character
			   SET state = 'frozen', updated_at = now(), version = version + 1
			 WHERE id = ANY($1) AND state = 'active'`, frozen)
		if err != nil {
			return 0, err
		}
		log.ErrorContext(ctx, "economy reconciler froze characters", "count", tag.RowsAffected())
	}

	ubRows, err := j.Pool.Query(ctx, unbalancedQuery)
	if err != nil {
		return 0, err
	}
	var unbalanced int64
	for ubRows.Next() {
		var txn uuid.UUID
		var cur string
		var sum int64
		if err := ubRows.Scan(&txn, &cur, &sum); err != nil {
			ubRows.Close()
			return 0, err
		}
		unbalanced++
		log.ErrorContext(ctx, "ledger txn not balanced",
			"txn_id", txn.String(), "currency", cur, "sum", sum)
	}
	ubRows.Close()
	if err := ubRows.Err(); err != nil {
		return 0, err
	}
	return int64(len(drifts)) + unbalanced, nil
}

// ChatGC deletes chat messages older than 30 days.
func ChatGC(pool *pgxpool.Pool) Job {
	return &sqlJob{name: "chat_gc", interval: time.Hour, pool: pool, query: `
DELETE FROM social.chat_message
 WHERE id IN (SELECT id FROM social.chat_message
               WHERE sent_at < now() - INTERVAL '30 days' LIMIT $1)`}
}

// MailExpire returns unclaimed mail value to senders after 30 days.
type MailExpire struct {
	Pool  *pgxpool.Pool
	Txn   economy.TxnService
	Slots social.SlotAssigner // re-slots returned items (they travel unslotted)
}

func (MailExpire) Name() string            { return "mail_expire" }
func (MailExpire) Interval() time.Duration { return 10 * time.Minute }
func (j MailExpire) Run(ctx context.Context) (int64, error) {
	return social.ExpireMail(ctx, j.Pool, j.Txn, j.Slots, BatchSize)
}

// MarketExpire expires overdue listings and resolves abandoned pending sales.
type MarketExpire struct{ Pool *pgxpool.Pool }

func (MarketExpire) Name() string            { return "market_expire" }
func (MarketExpire) Interval() time.Duration { return time.Minute }
func (j MarketExpire) Run(ctx context.Context) (int64, error) {
	return market.ExpireAndRecover(ctx, j.Pool, BatchSize)
}

// All returns the default job set for cmd/worker.
// Deps bundles what the job set needs beyond the pool.
type Deps struct {
	Txn             economy.TxnService
	Notifier        Notifier
	Senders         map[string]notification.Sender
	ReconcileFreeze bool
	Mailer          mail.Sender // nil → email_dispatch is a no-op
	Log             *slog.Logger
}

func All(pool *pgxpool.Pool, d Deps) []Job {
	txn, logger, reconcileFreeze := d.Txn, d.Log, d.ReconcileFreeze
	return []Job{
		ChatGC(pool),
		MailExpire{Pool: pool, Txn: txn, Slots: inventory.NewService(inventory.NewPgRepo(pool), nil)},
		MarketExpire{Pool: pool},
		ActivityExpire{Pool: pool, Notifier: d.Notifier, Log: logger},
		TravelArrive{Pool: pool, Notifier: d.Notifier, Log: logger},
		PushDispatch{Pool: pool, Senders: d.Senders},
		EmailDispatch{Pool: pool, Sender: d.Mailer},
		AccountPurge{Pool: pool},
		WorldEvents{Pool: pool, Txn: txn, Notifier: d.Notifier, Log: logger},
		NotificationGC(pool),
		SanctionExpire(pool),
		IdempotencyGC(pool),
		BuybackGC(pool),
		GuildDisbandExpire(pool),
		CharacterDeleteExpire{Pool: pool},
		EconomyReconciler{Pool: pool, Freeze: reconcileFreeze, Log: logger},
	}
}
