package guild

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/notification"
)

type Owner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Service struct {
	pool   *pgxpool.Pool
	chars  Owner
	txn    economy.TxnService
	pub    EventPublisher
	notify notification.Notifier
	now    func() time.Time
	log    *slog.Logger
}

func NewService(pool *pgxpool.Pool, chars Owner, txn economy.TxnService, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, chars: chars, txn: txn, now: time.Now, log: logger}
}

func (s *Service) WithPublisher(p EventPublisher) *Service { s.pub = p; return s }

func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notify = n; return s }

// --- reads ------------------------------------------------------------------

// Mine returns the actor's guild with member list and the actor's rank.
func (s *Service) Mine(ctx context.Context, actor, account uuid.UUID) (*Guild, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return nil, err
	}
	g, err := loadGuild(ctx, s.pool, guildID)
	if err != nil {
		return nil, err
	}
	if g.Members, err = members(ctx, s.pool, guildID); err != nil {
		return nil, err
	}
	g.MyRank = &rank
	return g, nil
}

// Public returns the guild summary without the member list.
func (s *Service) Public(ctx context.Context, guildID uuid.UUID) (*Guild, error) {
	return loadGuild(ctx, s.pool, guildID)
}

func (s *Service) Invites(ctx context.Context, actor, account uuid.UUID) ([]Invite, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT i.id, g.id, g.name, g.tag, c.name, i.expires_at
		  FROM social.guild_invite i
		  JOIN social.guild g ON g.id = i.guild_id AND g.disbanded_at IS NULL
		  JOIN character.character c ON c.id = i.invited_by_character_id
		 WHERE i.target_character_id = $1 AND i.status = 'pending' AND i.expires_at > now()
		 ORDER BY i.created_at DESC`, actor)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Invite, error) {
		var inv Invite
		err := r.Scan(&inv.ID, &inv.GuildID, &inv.GuildName, &inv.GuildTag, &inv.InvitedByName, &inv.ExpiresAt)
		return inv, err
	})
}

// --- create -----------------------------------------------------------------

type CreateRequest struct {
	Name           string
	Tag            string
	Description    string
	IdempotencyKey string
}

// Create charges the creation fee through the ledger, then inserts the guild
// with the actor as leader. If the insert loses a name/tag race after the
// debit, the fee is refunded with a compensating ledger txn.
func (s *Service) Create(ctx context.Context, actor, account uuid.UUID, req CreateRequest) (*Guild, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	if req.IdempotencyKey == "" {
		return nil, ErrIdempotencyMissing
	}
	name, nameNorm, err := NormalizeName(req.Name)
	if err != nil {
		return nil, err
	}
	tag, tagNorm, err := NormalizeTag(req.Tag)
	if err != nil {
		return nil, err
	}
	desc := strings.TrimSpace(req.Description)
	if len([]rune(desc)) > MaxDescription {
		return nil, ErrTooLong
	}
	ledgerKey := "guild-create-" + req.IdempotencyKey

	// Replay of an attempt that committed but never got its response cached.
	if g, err := s.guildCreatedWith(ctx, actor, ledgerKey); err == nil {
		return g, nil
	}
	if _, _, err := membership(ctx, s.pool, actor); err == nil {
		return nil, ErrAlreadyInGuild
	} else if !errors.Is(err, ErrNotInGuild) {
		return nil, err
	}
	var nameTaken, tagTaken bool
	if err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM social.guild WHERE name_normalized = $1 AND disbanded_at IS NULL),
		       EXISTS (SELECT 1 FROM social.guild WHERE tag_normalized  = $2 AND disbanded_at IS NULL)`,
		nameNorm, tagNorm).Scan(&nameTaken, &tagTaken); err != nil {
		return nil, err
	}
	if nameTaken {
		return nil, ErrNameTaken
	}
	if tagTaken {
		return nil, ErrTagTaken
	}

	res, err := s.txn.Write(ctx, economy.Intent{
		Kind:           "guild_create",
		RefKind:        "guild",
		IdempotencyKey: ledgerKey,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: -CreationFeeCopper},
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: CreationFeeCopper},
		},
	})
	if err != nil {
		return nil, err
	}

	var guildID uuid.UUID
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `
			INSERT INTO social.guild (name, name_normalized, tag, tag_normalized,
			                          leader_character_id, description, member_cap, creation_txn_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			name, nameNorm, tag, tagNorm, actor, desc, MemberCap, res.TxnID).Scan(&guildID); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO social.guild_member (guild_id, character_id, rank_index) VALUES ($1, $2, 0)`,
			guildID, actor)
		return err
	})
	if err != nil {
		if g, gerr := s.guildCreatedWith(ctx, actor, ledgerKey); gerr == nil {
			return g, nil
		}
		s.refund(ctx, actor, req.IdempotencyKey)
		return nil, mapUnique(err)
	}
	return s.Mine(ctx, actor, account)
}

// guildCreatedWith finds the actor-led guild whose creation fee used ledgerKey.
func (s *Service) guildCreatedWith(ctx context.Context, actor uuid.UUID, ledgerKey string) (*Guild, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `
		SELECT g.id FROM social.guild g
		 WHERE g.leader_character_id = $1 AND g.disbanded_at IS NULL
		   AND g.creation_txn_id IN (SELECT txn_id FROM economy.ledger_entry WHERE idempotency_key = $2)`,
		actor, ledgerKey).Scan(&id)
	if err != nil {
		return nil, err
	}
	return loadGuild(ctx, s.pool, id)
}

func (s *Service) refund(ctx context.Context, actor uuid.UUID, idem string) {
	_, err := s.txn.Write(ctx, economy.Intent{
		Kind:           "guild_create_refund",
		RefKind:        "guild",
		IdempotencyKey: "guild-create-refund-" + idem,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemVendorSink, Delta: -CreationFeeCopper},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: CreationFeeCopper},
		},
	})
	if err != nil {
		s.log.ErrorContext(ctx, "guild: creation fee refund failed — manual reconcile needed",
			"character_id", actor.String(), "idempotency_key", idem, "err", err.Error())
	}
}

// --- invites ----------------------------------------------------------------

func (s *Service) Invite(ctx context.Context, actor, account uuid.UUID, targetName string) (*Invite, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	target, targetAccount, err := characterByName(ctx, s.pool, targetName)
	if err != nil {
		return nil, err
	}
	if target == actor {
		return nil, ErrSelfTarget
	}
	inv := &Invite{ExpiresAt: s.now().UTC().Add(InviteTTL)}
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		g, err := lockMembership(ctx, tx, actor)
		if err != nil {
			return err
		}
		if !CanInvite(g.rank) {
			return ErrForbidden
		}
		if g.members >= int(g.cap) {
			return ErrGuildFull
		}
		if _, _, err := membership(ctx, tx, target); err == nil {
			return ErrTargetInGuild
		} else if !errors.Is(err, ErrNotInGuild) {
			return err
		}
		inv.GuildID, inv.GuildName, inv.GuildTag = g.id, g.name, g.tag
		if _, err := tx.Exec(ctx, `
			UPDATE social.guild_invite SET status = 'expired'
			 WHERE guild_id = $1 AND target_character_id = $2 AND status = 'pending' AND expires_at <= now()`,
			g.id, target); err != nil {
			return err
		}
		err = tx.QueryRow(ctx, `
			INSERT INTO social.guild_invite (guild_id, target_character_id, invited_by_character_id, expires_at)
			VALUES ($1, $2, $3, $4) RETURNING id`, g.id, target, actor, inv.ExpiresAt).Scan(&inv.ID)
		if isUnique(err) {
			return ErrAlreadyInvited
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	s.publish(targetAccount, "guild.invite", map[string]any{
		"invite_id": inv.ID, "guild_id": inv.GuildID, "guild_name": inv.GuildName, "expires_at": inv.ExpiresAt,
	})
	notification.Send(ctx, s.notify, s.log, notification.Notice{
		AccountID: targetAccount, CharacterID: &target, Kind: "guild_invite",
		DedupKey: "guild_invite:" + inv.ID.String(),
		TitleKey: "notif.guild_invite.title", BodyKey: "notif.guild_invite.body",
		Vars:  map[string]any{"guild_name": inv.GuildName},
		Route: "character/guild/invites", Category: "social",
	})
	return inv, nil
}

func (s *Service) Accept(ctx context.Context, actor, account, inviteID uuid.UUID) (*Guild, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var guildID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var target uuid.UUID
		var status string
		var expires time.Time
		err := tx.QueryRow(ctx, `
			SELECT guild_id, target_character_id, status, expires_at
			  FROM social.guild_invite WHERE id = $1 FOR UPDATE`, inviteID).
			Scan(&guildID, &target, &status, &expires)
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && target != actor) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if status != "pending" || !expires.After(s.now()) {
			return ErrInviteInvalid
		}
		var capN int16
		var scheduled, disbanded *time.Time
		if err := tx.QueryRow(ctx, `
			SELECT member_cap, disband_scheduled_at, disbanded_at FROM social.guild WHERE id = $1 FOR UPDATE`,
			guildID).Scan(&capN, &scheduled, &disbanded); err != nil {
			return err
		}
		if scheduled != nil || disbanded != nil {
			return ErrInviteInvalid
		}
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM social.guild_member WHERE guild_id = $1`, guildID).
			Scan(&n); err != nil {
			return err
		}
		if n >= int(capN) {
			return ErrGuildFull
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO social.guild_member (guild_id, character_id, rank_index) VALUES ($1, $2, $3)`,
			guildID, actor, RankInitiate); err != nil {
			if isUnique(err) {
				return ErrAlreadyInGuild
			}
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE social.guild_invite SET status = 'accepted' WHERE id = $1`, inviteID)
		return err
	})
	if err != nil {
		return nil, err
	}
	s.notifyMembers(ctx, guildID, "guild.changed", map[string]any{"guild_id": guildID, "joined": actor})
	return s.Mine(ctx, actor, account)
}

func (s *Service) Decline(ctx context.Context, actor, account, inviteID uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE social.guild_invite SET status = 'declined'
		 WHERE id = $1 AND target_character_id = $2 AND status = 'pending'`, inviteID, actor)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- membership changes -----------------------------------------------------

func (s *Service) Leave(ctx context.Context, actor, account uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	var guildID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		g, err := lockMembership(ctx, tx, actor)
		if err != nil {
			return err
		}
		if g.rank == RankLeader {
			return ErrLeaderCannotLeave
		}
		guildID = g.id
		_, err = tx.Exec(ctx, `DELETE FROM social.guild_member WHERE guild_id = $1 AND character_id = $2`, g.id, actor)
		return err
	})
	if err != nil {
		return err
	}
	s.notifyMembers(ctx, guildID, "guild.changed", map[string]any{"guild_id": guildID, "left": actor})
	return nil
}

func (s *Service) Kick(ctx context.Context, actor, account, target uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	if target == actor {
		return ErrSelfTarget
	}
	var guildID, targetAccount uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		g, err := lockMembership(ctx, tx, actor)
		if err != nil {
			return err
		}
		guildID = g.id
		targetRank, err := rankIn(ctx, tx, g.id, target)
		if err != nil {
			return err
		}
		if !CanKick(g.rank, targetRank) {
			return ErrForbidden
		}
		return tx.QueryRow(ctx, `
			DELETE FROM social.guild_member m USING character.character c
			 WHERE m.guild_id = $1 AND m.character_id = $2 AND c.id = m.character_id
			RETURNING c.account_id`, g.id, target).Scan(&targetAccount)
	})
	if err != nil {
		return err
	}
	s.publish(targetAccount, "guild.kicked", map[string]any{"guild_id": guildID})
	s.notifyMembers(ctx, guildID, "guild.changed", map[string]any{"guild_id": guildID, "kicked": target})
	return nil
}

func (s *Service) SetRank(ctx context.Context, actor, account, target uuid.UUID, newRank int16) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	if target == actor {
		return ErrSelfTarget
	}
	var guildID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		g, err := lockMembership(ctx, tx, actor)
		if err != nil {
			return err
		}
		guildID = g.id
		targetRank, err := rankIn(ctx, tx, g.id, target)
		if err != nil {
			return err
		}
		if !CanSetRank(g.rank, targetRank, newRank) {
			return ErrForbidden
		}
		_, err = tx.Exec(ctx, `
			UPDATE social.guild_member SET rank_index = $3 WHERE guild_id = $1 AND character_id = $2`,
			g.id, target, newRank)
		return err
	})
	if err != nil {
		return err
	}
	s.notifyMembers(ctx, guildID, "guild.changed", map[string]any{"guild_id": guildID, "rank": target, "rank_index": newRank})
	return nil
}

// TransferLeadership makes target the leader; the old leader becomes officer.
func (s *Service) TransferLeadership(ctx context.Context, actor, account, target uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	if target == actor {
		return ErrSelfTarget
	}
	var guildID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		g, err := lockMembership(ctx, tx, actor)
		if err != nil {
			return err
		}
		guildID = g.id
		if g.rank != RankLeader {
			return ErrForbidden
		}
		if _, err := rankIn(ctx, tx, g.id, target); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			UPDATE social.guild_member SET rank_index = CASE WHEN character_id = $2 THEN 0 ELSE 1 END
			 WHERE guild_id = $1 AND character_id IN ($2, $3)`, g.id, target, actor); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE social.guild SET leader_character_id = $2, version = version + 1 WHERE id = $1`, g.id, target)
		return err
	})
	if err != nil {
		return err
	}
	s.notifyMembers(ctx, guildID, "guild.changed", map[string]any{"guild_id": guildID, "leader": target})
	return nil
}

func (s *Service) SetMOTD(ctx context.Context, actor, account uuid.UUID, motd string) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	motd = strings.TrimSpace(motd)
	if len([]rune(motd)) > MaxMOTD {
		return ErrTooLong
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return err
	}
	if !CanSetMOTD(rank) {
		return ErrForbidden
	}
	if _, err := s.pool.Exec(ctx, `
		UPDATE social.guild SET motd = $2, version = version + 1 WHERE id = $1`, guildID, motd); err != nil {
		return err
	}
	s.notifyMembers(ctx, guildID, "guild.motd", map[string]any{"guild_id": guildID, "motd": motd})
	return nil
}

// ScheduleDisband starts the 48h grace; the worker finalises it.
func (s *Service) ScheduleDisband(ctx context.Context, actor, account uuid.UUID) (time.Time, error) {
	at := s.now().UTC().Add(DisbandGrace)
	err := s.setDisband(ctx, actor, account, &at)
	return at, err
}

func (s *Service) CancelDisband(ctx context.Context, actor, account uuid.UUID) error {
	return s.setDisband(ctx, actor, account, nil)
}

func (s *Service) setDisband(ctx context.Context, actor, account uuid.UUID, at *time.Time) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return err
	}
	if rank != RankLeader {
		return ErrForbidden
	}
	cond := `disband_scheduled_at IS NULL`
	errNone := ErrDisbandScheduled
	if at == nil {
		cond = `disband_scheduled_at IS NOT NULL`
		errNone = ErrDisbandNotPending
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE social.guild SET disband_scheduled_at = $2, version = version + 1
		 WHERE id = $1 AND disbanded_at IS NULL AND `+cond, guildID, at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNone
	}
	s.notifyMembers(ctx, guildID, "guild.disband", map[string]any{"guild_id": guildID, "disband_scheduled_at": at})
	return nil
}

// --- helpers ----------------------------------------------------------------

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type lockedGuild struct {
	id      uuid.UUID
	name    string
	tag     string
	rank    int16
	cap     int16
	members int
}

// lockMembership locks the actor's guild row and returns the actor's rank.
func lockMembership(ctx context.Context, tx pgx.Tx, actor uuid.UUID) (*lockedGuild, error) {
	var g lockedGuild
	err := tx.QueryRow(ctx, `
		SELECT g.id, g.name, g.tag, m.rank_index, g.member_cap
		  FROM social.guild g
		  JOIN social.guild_member m ON m.guild_id = g.id
		 WHERE m.character_id = $1 AND g.disbanded_at IS NULL
		   FOR UPDATE OF g`, actor).Scan(&g.id, &g.name, &g.tag, &g.rank, &g.cap)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotInGuild
	}
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(ctx, `SELECT count(*) FROM social.guild_member WHERE guild_id = $1`, g.id).Scan(&g.members)
	return &g, err
}

func membership(ctx context.Context, q querier, character uuid.UUID) (uuid.UUID, int16, error) {
	var id uuid.UUID
	var rank int16
	err := q.QueryRow(ctx, `
		SELECT m.guild_id, m.rank_index FROM social.guild_member m
		  JOIN social.guild g ON g.id = m.guild_id AND g.disbanded_at IS NULL
		 WHERE m.character_id = $1`, character).Scan(&id, &rank)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, 0, ErrNotInGuild
	}
	return id, rank, err
}

func rankIn(ctx context.Context, q querier, guildID, character uuid.UUID) (int16, error) {
	var rank int16
	err := q.QueryRow(ctx, `
		SELECT rank_index FROM social.guild_member WHERE guild_id = $1 AND character_id = $2`,
		guildID, character).Scan(&rank)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotMember
	}
	return rank, err
}

func loadGuild(ctx context.Context, q querier, guildID uuid.UUID) (*Guild, error) {
	var g Guild
	err := q.QueryRow(ctx, `
		SELECT g.id, g.name, g.tag, g.leader_character_id, g.motd, g.description, g.level, g.xp,
		       (SELECT count(*) FROM social.guild_member m WHERE m.guild_id = g.id),
		       g.member_cap, g.created_at, g.disband_scheduled_at
		  FROM social.guild g WHERE g.id = $1 AND g.disbanded_at IS NULL`, guildID).
		Scan(&g.ID, &g.Name, &g.Tag, &g.LeaderCharacterID, &g.MOTD, &g.Description, &g.Level, &g.XP,
			&g.MemberCount, &g.MemberCap, &g.CreatedAt, &g.DisbandScheduledAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &g, err
}

func members(ctx context.Context, q querier, guildID uuid.UUID) ([]Member, error) {
	rows, err := q.Query(ctx, `
		SELECT c.id, c.name, c.level, c.class_code, m.rank_index, m.joined_at
		  FROM social.guild_member m JOIN character.character c ON c.id = m.character_id
		 WHERE m.guild_id = $1
		 ORDER BY m.rank_index, m.joined_at`, guildID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Member, error) {
		var m Member
		err := r.Scan(&m.CharacterID, &m.Name, &m.Level, &m.ClassCode, &m.Rank, &m.JoinedAt)
		m.RankName = RankNames[m.Rank]
		return m, err
	})
}

func characterByName(ctx context.Context, q querier, name string) (uuid.UUID, uuid.UUID, error) {
	var id, account uuid.UUID
	err := q.QueryRow(ctx, `
		SELECT id, account_id FROM character.character
		 WHERE name_normalized = $1 AND state = 'active'`, strings.ToLower(strings.TrimSpace(name))).
		Scan(&id, &account)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, uuid.Nil, ErrTargetNotFound
	}
	return id, account, err
}

func (s *Service) publish(account uuid.UUID, kind string, payload any) {
	if s.pub != nil && account != uuid.Nil {
		s.pub.Publish(account, kind, payload)
	}
}

func (s *Service) notifyMembers(ctx context.Context, guildID uuid.UUID, kind string, payload any) {
	if s.pub == nil {
		return
	}
	rows, err := s.pool.Query(ctx, `
		SELECT DISTINCT c.account_id FROM social.guild_member m
		  JOIN character.character c ON c.id = m.character_id WHERE m.guild_id = $1`, guildID)
	if err != nil {
		s.log.WarnContext(ctx, "guild: member fanout lookup failed", "err", err.Error())
		return
	}
	accounts, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return
	}
	for _, a := range accounts {
		s.pub.Publish(a, kind, payload)
	}
}

func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// mapUnique turns guild-insert unique violations into domain errors.
func mapUnique(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return err
	}
	switch pgErr.ConstraintName {
	case "ux_guild_name":
		return ErrNameTaken
	case "ux_guild_tag":
		return ErrTagTaken
	default:
		return ErrAlreadyInGuild
	}
}
