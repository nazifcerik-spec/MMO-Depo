package party

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

	"github.com/otter/server/internal/notification"
)

// Owner verifies an account owns a character. character.Service satisfies.
type Owner interface {
	Owns(ctx context.Context, characterID, accountID uuid.UUID) error
}

// EventPublisher fans events to an account's live WS sessions.
type EventPublisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

type Service struct {
	pool   *pgxpool.Pool
	chars  Owner
	pub    EventPublisher
	notify notification.Notifier
	now    func() time.Time
	log    *slog.Logger
}

func NewService(pool *pgxpool.Pool, chars Owner, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, chars: chars, now: time.Now, log: logger}
}

func (s *Service) WithPublisher(p EventPublisher) *Service { s.pub = p; return s }

func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notify = n; return s }

// --- reads ------------------------------------------------------------------

// Get returns the actor's live party, or ErrNotInParty.
func (s *Service) Get(ctx context.Context, actor, account uuid.UUID) (*Party, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	partyID, err := livePartyOf(ctx, s.pool, actor)
	if err != nil {
		return nil, err
	}
	return loadParty(ctx, s.pool, partyID)
}

// Invites lists pending, unexpired invites addressed to the actor.
func (s *Service) Invites(ctx context.Context, actor, account uuid.UUID) ([]Invite, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT i.id, i.party_id, i.invited_by_character_id, c.name, i.target_character_id, i.expires_at
		  FROM social.party_invite i
		  JOIN social.party p ON p.id = i.party_id AND p.disbanded_at IS NULL
		  JOIN character.character c ON c.id = i.invited_by_character_id
		 WHERE i.target_character_id = $1 AND i.status = 'pending' AND i.expires_at > now()
		 ORDER BY i.created_at DESC`, actor)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Invite, error) {
		var inv Invite
		err := r.Scan(&inv.ID, &inv.PartyID, &inv.InvitedBy, &inv.InvitedByName, &inv.TargetCharacter, &inv.ExpiresAt)
		return inv, err
	})
}

// --- mutations --------------------------------------------------------------

func (s *Service) Create(ctx context.Context, actor, account uuid.UUID) (*Party, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var partyID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var err error
		partyID, err = createParty(ctx, tx, actor)
		return err
	})
	if err != nil {
		return nil, err
	}
	return loadParty(ctx, s.pool, partyID)
}

// createParty opens a party led by leader.
func createParty(ctx context.Context, tx pgx.Tx, leader uuid.UUID) (uuid.UUID, error) {
	var partyID uuid.UUID
	if err := tx.QueryRow(ctx,
		`INSERT INTO social.party (leader_character_id) VALUES ($1) RETURNING id`, leader).
		Scan(&partyID); err != nil {
		return uuid.Nil, err
	}
	_, err := tx.Exec(ctx,
		`INSERT INTO social.party_member (party_id, character_id) VALUES ($1, $2)`, partyID, leader)
	if isUnique(err) {
		return uuid.Nil, ErrAlreadyInParty
	}
	return partyID, err
}

// LeadOrCreate returns the actor's party, opening one led by the actor when
// they have none (a raid listing, DD-299). Only the leader may advertise.
func (s *Service) LeadOrCreate(ctx context.Context, actor, account uuid.UUID) (uuid.UUID, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return uuid.Nil, err
	}
	var partyID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		p, err := lockActorParty(ctx, tx, actor)
		if errors.Is(err, ErrNotInParty) {
			partyID, err = createParty(ctx, tx, actor)
			return err
		}
		if err != nil {
			return err
		}
		if p.leader != actor {
			return ErrNotLeader
		}
		partyID = p.id
		return nil
	})
	return partyID, err
}

// Join adds the actor to a party without an invite (a raid listing,
// DD-299).
func (s *Service) Join(ctx context.Context, actor, account, partyID uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		return addMembers(ctx, tx, partyID, []uuid.UUID{actor})
	})
	if err != nil {
		return err
	}
	s.notifyMembers(ctx, partyID, "party.changed", map[string]any{"party_id": partyID, "joined": actor})
	return nil
}

// Form puts characters that are in no party into a new one, the first
// leading (the raid queue, DD-299). All or nothing.
func (s *Service) Form(ctx context.Context, members []uuid.UUID) (uuid.UUID, error) {
	if len(members) < SizeMin || len(members) > SizeCap {
		return uuid.Nil, ErrPartyFull
	}
	var partyID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var err error
		if partyID, err = createParty(ctx, tx, members[0]); err != nil {
			return err
		}
		return addMembers(ctx, tx, partyID, members[1:])
	})
	if err != nil {
		return uuid.Nil, err
	}
	s.notifyMembers(ctx, partyID, "party.changed", map[string]any{"party_id": partyID, "formed": true})
	return partyID, nil
}

// addMembers joins characters to a live party under its row lock, within
// its size cap.
func addMembers(ctx context.Context, tx pgx.Tx, partyID uuid.UUID, chars []uuid.UUID) error {
	var sizeCap int16
	var disbanded *time.Time
	err := tx.QueryRow(ctx, `SELECT size_cap, disbanded_at FROM social.party WHERE id = $1 FOR UPDATE`, partyID).
		Scan(&sizeCap, &disbanded)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && disbanded != nil) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	var n int
	if err := tx.QueryRow(ctx, `
		SELECT count(*) FROM social.party_member WHERE party_id = $1 AND left_at IS NULL`, partyID).Scan(&n); err != nil {
		return err
	}
	if n+len(chars) > int(sizeCap) {
		return ErrPartyFull
	}
	for _, c := range chars {
		if _, err := tx.Exec(ctx, `
			INSERT INTO social.party_member (party_id, character_id) VALUES ($1, $2)
			ON CONFLICT (party_id, character_id) DO UPDATE SET left_at = NULL, joined_at = now()`,
			partyID, c); err != nil {
			if isUnique(err) {
				return ErrAlreadyInParty
			}
			return err
		}
	}
	return nil
}

// Invite sends a 60s invite to the character named targetName.
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
	inv := &Invite{InvitedBy: actor, TargetCharacter: target, ExpiresAt: s.now().UTC().Add(InviteTTL)}
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		p, err := lockActorParty(ctx, tx, actor)
		if errors.Is(err, ErrNotInParty) { // the first invite forms the party (DD-297)
			if _, err = createParty(ctx, tx, actor); err != nil {
				return err
			}
			p, err = lockActorParty(ctx, tx, actor)
		}
		if err != nil {
			return err
		}
		if p.leader != actor {
			return ErrNotLeader
		}
		if p.members >= int(p.sizeCap) {
			return ErrPartyFull
		}
		if _, err := livePartyOf(ctx, tx, target); err == nil {
			return ErrTargetInParty
		} else if !errors.Is(err, ErrNotInParty) {
			return err
		}
		inv.PartyID = p.id
		if _, err := tx.Exec(ctx, `
			UPDATE social.party_invite SET status = 'expired'
			 WHERE party_id = $1 AND target_character_id = $2 AND status = 'pending' AND expires_at <= now()`,
			p.id, target); err != nil {
			return err
		}
		err = tx.QueryRow(ctx, `
			INSERT INTO social.party_invite (party_id, target_character_id, invited_by_character_id, expires_at)
			VALUES ($1, $2, $3, $4) RETURNING id`, p.id, target, actor, inv.ExpiresAt).Scan(&inv.ID)
		if isUnique(err) {
			return ErrAlreadyInvited
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	s.publish(targetAccount, "party.invite", map[string]any{
		"invite_id": inv.ID, "party_id": inv.PartyID, "expires_at": inv.ExpiresAt,
	})
	notification.Send(ctx, s.notify, s.log, notification.Notice{
		AccountID: targetAccount, CharacterID: &target, Kind: "party_invite",
		DedupKey: "party_invite:" + inv.ID.String(),
		TitleKey: "notif.party_invite.title", BodyKey: "notif.party_invite.body",
		Route: "character/party/invites", Category: "social", High: true,
	})
	return inv, nil
}

// Accept joins the party behind a pending invite addressed to the actor.
func (s *Service) Accept(ctx context.Context, actor, account, inviteID uuid.UUID) (*Party, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var partyID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var status string
		var target uuid.UUID
		var expires time.Time
		err := tx.QueryRow(ctx, `
			SELECT party_id, target_character_id, status, expires_at
			  FROM social.party_invite WHERE id = $1 FOR UPDATE`, inviteID).
			Scan(&partyID, &target, &status, &expires)
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && target != actor) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if status != "pending" || !expires.After(s.now()) {
			return ErrInviteInvalid
		}
		if err := addMembers(ctx, tx, partyID, []uuid.UUID{actor}); err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrInviteInvalid
			}
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE social.party_invite SET status = 'accepted' WHERE id = $1`, inviteID)
		return err
	})
	if err != nil {
		return nil, err
	}
	s.notifyMembers(ctx, partyID, "party.changed", map[string]any{"party_id": partyID, "joined": actor})
	return loadParty(ctx, s.pool, partyID)
}

func (s *Service) Decline(ctx context.Context, actor, account, inviteID uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE social.party_invite SET status = 'declined'
		 WHERE id = $1 AND target_character_id = $2 AND status = 'pending'`, inviteID, actor)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Leave removes the actor. A departing leader hands off to the earliest
// joiner; a party left with one member disbands (DD-297).
func (s *Service) Leave(ctx context.Context, actor, account uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	return s.removeMember(ctx, actor, actor, "left")
}

// Kick removes another member (leader only).
func (s *Service) Kick(ctx context.Context, actor, account, target uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	if target == actor {
		return ErrSelfTarget
	}
	return s.removeMember(ctx, actor, target, "kicked")
}

func (s *Service) removeMember(ctx context.Context, actor, target uuid.UUID, reason string) error {
	var partyID uuid.UUID
	var disbanded bool
	var targetAccount uuid.UUID
	var rest []uuid.UUID // accounts told of a disband
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		p, err := lockActorParty(ctx, tx, actor)
		if err != nil {
			return err
		}
		partyID = p.id
		if target != actor && p.leader != actor {
			return ErrNotLeader
		}
		members, err := liveMembers(ctx, tx, p.id)
		if err != nil {
			return err
		}
		found := false
		for _, m := range members {
			if m.CharacterID == target {
				found = true
			}
		}
		if !found {
			return ErrNotMember
		}
		if err := tx.QueryRow(ctx, `
			UPDATE social.party_member pm SET left_at = now()
			  FROM character.character c
			 WHERE pm.party_id = $1 AND pm.character_id = $2 AND pm.left_at IS NULL AND c.id = pm.character_id
			RETURNING c.account_id`, p.id, target).Scan(&targetAccount); err != nil {
			return err
		}
		if len(members)-1 < SizeMin && len(members) >= SizeMin { // a party of one is no party
			disbanded = true
			if rest, err = memberAccounts(ctx, tx, p.id); err != nil {
				return err
			}
			return disband(ctx, tx, p.id)
		}
		if target != p.leader {
			return nil
		}
		next, ok := NextLeader(members, target)
		if !ok {
			disbanded = true
			return disband(ctx, tx, p.id)
		}
		_, err = tx.Exec(ctx, `
			UPDATE social.party SET leader_character_id = $2, version = version + 1 WHERE id = $1`, p.id, next)
		return err
	})
	if err != nil {
		return err
	}
	if reason == "kicked" {
		s.publish(targetAccount, "party.kicked", map[string]any{"party_id": partyID})
	}
	if !disbanded {
		s.notifyMembers(ctx, partyID, "party.changed", map[string]any{"party_id": partyID, reason: target})
	}
	for _, a := range rest {
		s.publish(a, "party.disbanded", map[string]any{"party_id": partyID})
	}
	return nil
}

// Promote hands leadership to another live member.
func (s *Service) Promote(ctx context.Context, actor, account, target uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	var partyID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		p, err := lockActorParty(ctx, tx, actor)
		if err != nil {
			return err
		}
		partyID = p.id
		if p.leader != actor {
			return ErrNotLeader
		}
		tag, err := tx.Exec(ctx, `
			UPDATE social.party SET leader_character_id = $2, version = version + 1
			 WHERE id = $1 AND EXISTS (SELECT 1 FROM social.party_member
			                           WHERE party_id = $1 AND character_id = $2 AND left_at IS NULL)`,
			p.id, target)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotMember
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.notifyMembers(ctx, partyID, "party.changed", map[string]any{"party_id": partyID, "leader": target})
	return nil
}

// Disband ends the party for everyone (leader only).
func (s *Service) Disband(ctx context.Context, actor, account uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	var partyID uuid.UUID
	var accounts []uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		p, err := lockActorParty(ctx, tx, actor)
		if err != nil {
			return err
		}
		partyID = p.id
		if p.leader != actor {
			return ErrNotLeader
		}
		if accounts, err = memberAccounts(ctx, tx, p.id); err != nil {
			return err
		}
		return disband(ctx, tx, p.id)
	})
	if err != nil {
		return err
	}
	for _, a := range accounts {
		s.publish(a, "party.disbanded", map[string]any{"party_id": partyID})
	}
	return nil
}

// --- helpers ----------------------------------------------------------------

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type lockedParty struct {
	id      uuid.UUID
	leader  uuid.UUID
	sizeCap int16
	members int
}

// lockActorParty locks the actor's live party row for the rest of the tx.
func lockActorParty(ctx context.Context, tx pgx.Tx, actor uuid.UUID) (*lockedParty, error) {
	var p lockedParty
	err := tx.QueryRow(ctx, `
		SELECT p.id, p.leader_character_id, p.size_cap
		  FROM social.party p
		  JOIN social.party_member m ON m.party_id = p.id AND m.left_at IS NULL
		 WHERE m.character_id = $1 AND p.disbanded_at IS NULL
		   FOR UPDATE OF p`, actor).Scan(&p.id, &p.leader, &p.sizeCap)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotInParty
	}
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(ctx, `
		SELECT count(*) FROM social.party_member WHERE party_id = $1 AND left_at IS NULL`, p.id).
		Scan(&p.members)
	return &p, err
}

func livePartyOf(ctx context.Context, q querier, character uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx, `
		SELECT m.party_id FROM social.party_member m
		  JOIN social.party p ON p.id = m.party_id AND p.disbanded_at IS NULL
		 WHERE m.character_id = $1 AND m.left_at IS NULL`, character).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrNotInParty
	}
	return id, err
}

func liveMembers(ctx context.Context, q querier, partyID uuid.UUID) ([]Member, error) {
	rows, err := q.Query(ctx, `
		SELECT c.id, c.name, c.level, c.class_code, m.role_hint, m.joined_at
		  FROM social.party_member m JOIN character.character c ON c.id = m.character_id
		 WHERE m.party_id = $1 AND m.left_at IS NULL
		 ORDER BY m.joined_at`, partyID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Member, error) {
		var m Member
		err := r.Scan(&m.CharacterID, &m.Name, &m.Level, &m.ClassCode, &m.RoleHint, &m.JoinedAt)
		return m, err
	})
}

func loadParty(ctx context.Context, q querier, partyID uuid.UUID) (*Party, error) {
	var p Party
	err := q.QueryRow(ctx, `
		SELECT id, leader_character_id, size_cap, loot_rule, created_at
		  FROM social.party WHERE id = $1 AND disbanded_at IS NULL`, partyID).
		Scan(&p.ID, &p.LeaderCharacterID, &p.SizeCap, &p.LootRule, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.Members, err = liveMembers(ctx, q, partyID)
	return &p, err
}

func memberAccounts(ctx context.Context, q querier, partyID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.Query(ctx, `
		SELECT DISTINCT c.account_id FROM social.party_member m
		  JOIN character.character c ON c.id = m.character_id
		 WHERE m.party_id = $1 AND m.left_at IS NULL`, partyID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
}

func disband(ctx context.Context, tx pgx.Tx, partyID uuid.UUID) error {
	if _, err := tx.Exec(ctx, `
		UPDATE social.party SET disbanded_at = now(), version = version + 1 WHERE id = $1`, partyID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE social.party_member SET left_at = now() WHERE party_id = $1 AND left_at IS NULL`, partyID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `
		UPDATE social.party_invite SET status = 'expired' WHERE party_id = $1 AND status = 'pending'`, partyID)
	return err
}

// characterByName resolves an active character by (case-insensitive) name.
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

func (s *Service) notifyMembers(ctx context.Context, partyID uuid.UUID, kind string, payload any) {
	if s.pub == nil {
		return
	}
	accounts, err := memberAccounts(ctx, s.pool, partyID)
	if err != nil {
		s.log.WarnContext(ctx, "party: member fanout lookup failed", "err", err.Error())
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
