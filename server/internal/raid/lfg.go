package raid

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/party"
)

// Raid finding (DD-299): a queue that puts waiting characters into a random
// party with its lobby open, and a board where a party advertises for
// members anyone may take.

const (
	QueueFill = 5               // the queue forms a party as soon as this many wait …
	QueueWait = 5 * time.Minute // … or, once the first has waited this long, with the raid's minimum
	MaxListed = 50              // listings in one board page
)

var (
	ErrInParty   = errors.New("raid: leave your party to queue")
	ErrNotQueued = errors.New("raid: character is not queued")
	ErrListed    = errors.New("raid: listing is gone or full")
)

// Parties forms and joins parties for the queue and the board
// (party.Service satisfies).
type Parties interface {
	Form(ctx context.Context, members []uuid.UUID) (uuid.UUID, error)
	LeadOrCreate(ctx context.Context, actor, account uuid.UUID) (uuid.UUID, error)
	Join(ctx context.Context, actor, account, partyID uuid.UUID) error
}

// WithParties enables the raid queue and the listing board.
func (s *Service) WithParties(p Parties) *Service { s.parties = p; return s }

// QueueStatus is a queued character's place: its raid, since when, how many
// wait for it, and when the queue forms a smaller party at the latest.
type QueueStatus struct {
	RaidCode string    `json:"raid_code"`
	QueuedAt time.Time `json:"queued_at"`
	Waiting  int       `json:"waiting"`
	FormsBy  time.Time `json:"forms_by"`
}

// ListingMember is one member of an advertising party.
type ListingMember struct {
	Name      string `json:"name"`
	Level     int16  `json:"level"`
	ClassCode string `json:"class_code"`
}

// Listing is a party advertising for a raid.
type Listing struct {
	ID        uuid.UUID       `json:"id"`
	RaidCode  string          `json:"raid_code"`
	PartyID   uuid.UUID       `json:"party_id"`
	Members   []ListingMember `json:"members"` // the leader first
	CreatedAt time.Time       `json:"created_at"`
}

// --- Queue -------------------------------------------------------------------

// Enqueue puts the actor in the queue for a raid (a new raid replaces the
// old place). Only characters in no party who meet the level queue.
func (s *Service) Enqueue(ctx context.Context, actor, account uuid.UUID, code string) (*QueueStatus, error) {
	tmpl, err := s.eligible(ctx, actor, account, code)
	if err != nil {
		return nil, err
	}
	if _, _, err := partyOf(ctx, s.pool, actor); err == nil {
		return nil, ErrInParty
	} else if !errors.Is(err, ErrNoParty) {
		return nil, err
	}
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO raid.queue (character_id, account_id, raid_code, queued_at) VALUES ($1, $2, $3, $4)
		ON CONFLICT (character_id) DO UPDATE SET raid_code = EXCLUDED.raid_code, queued_at = EXCLUDED.queued_at`,
		actor, account, tmpl.Code, s.now().UTC()); err != nil {
		return nil, err
	}
	s.match(ctx, tmpl)
	return s.queueStatus(ctx, actor)
}

// Dequeue takes the actor out of the queue.
func (s *Service) Dequeue(ctx context.Context, actor, account uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM raid.queue WHERE character_id = $1`, actor)
	if err == nil && tag.RowsAffected() == 0 {
		err = ErrNotQueued
	}
	return err
}

// queueStatus is the actor's place in the queue after a match attempt for
// its raid; nil once it was matched (or never queued).
func (s *Service) queueStatus(ctx context.Context, actor uuid.UUID) (*QueueStatus, error) {
	var q QueueStatus
	err := s.pool.QueryRow(ctx, `
		SELECT q.raid_code, q.queued_at,
		       (SELECT count(*) FROM raid.queue o WHERE o.raid_code = q.raid_code),
		       (SELECT min(o.queued_at) FROM raid.queue o WHERE o.raid_code = q.raid_code)
		  FROM raid.queue q WHERE q.character_id = $1`, actor).Scan(&q.RaidCode, &q.QueuedAt, &q.Waiting, &q.FormsBy)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	q.FormsBy = q.FormsBy.Add(QueueWait)
	return &q, err
}

type queued struct {
	id, account uuid.UUID
	at          time.Time
}

// match forms a party from the raid's queue when it can: QueueFill waiting,
// or the raid's minimum once the first has waited QueueWait. The party gets
// the raid's lobby with everyone who can go already in; the leader starts.
// Failures are logged: the queue is retried at the next look.
func (s *Service) match(ctx context.Context, tmpl dungeon.Template) {
	if s.parties == nil {
		return
	}
	var picked []queued
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended('otter.raid.queue:' || $1, 0))`, tmpl.Code); err != nil {
			return err
		}
		rows, err := tx.Query(ctx, `
			SELECT q.character_id, q.account_id, q.queued_at
			  FROM raid.queue q
			 WHERE q.raid_code = $1
			   AND NOT EXISTS (SELECT 1 FROM social.party_member m JOIN social.party p ON p.id = m.party_id
			                    WHERE m.character_id = q.character_id AND m.left_at IS NULL AND p.disbanded_at IS NULL)
			 ORDER BY q.queued_at LIMIT $2`, tmpl.Code, QueueFill)
		if err != nil {
			return err
		}
		all, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (queued, error) {
			var q queued
			err := r.Scan(&q.id, &q.account, &q.at)
			return q, err
		})
		if err != nil {
			return err
		}
		min := max(int(tmpl.GroupMin), 2)
		if len(all) < QueueFill && (len(all) < min || s.now().Sub(all[0].at) < QueueWait) {
			return nil
		}
		picked = all
		ids := make([]uuid.UUID, len(all))
		for i, q := range all {
			ids[i] = q.id
		}
		_, err = tx.Exec(ctx, `DELETE FROM raid.queue WHERE character_id = ANY($1)`, ids)
		return err
	})
	if err != nil || len(picked) == 0 {
		if err != nil {
			s.log.WarnContext(ctx, "raid: queue match failed", slog.String("raid", tmpl.Code), slog.String("err", err.Error()))
		}
		return
	}
	ids := make([]uuid.UUID, len(picked))
	for i, q := range picked {
		ids[i] = q.id
	}
	partyID, err := s.parties.Form(ctx, ids)
	if err != nil { // someone joined a party meanwhile: back in line, in order
		s.log.WarnContext(ctx, "raid: queue party not formed", slog.String("raid", tmpl.Code), slog.String("err", err.Error()))
		for _, q := range picked {
			_, _ = s.pool.Exec(ctx, `
				INSERT INTO raid.queue (character_id, account_id, raid_code, queued_at) VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING`, q.id, q.account, tmpl.Code, q.at)
		}
		return
	}
	var members []lobbyMember
	for _, q := range picked {
		snap, err := s.snapshot(ctx, q.id, q.account, tmpl)
		if err != nil { // in the party, but it joins the lobby itself once it can
			continue
		}
		members = append(members, lobbyMember{id: q.id, snap: snap})
	}
	var runID uuid.UUID
	if len(members) > 0 {
		if runID, err = s.openLobby(ctx, partyID, tmpl, members); err != nil {
			s.log.WarnContext(ctx, "raid: queue lobby not opened", slog.String("err", err.Error()))
		}
	}
	for _, q := range picked {
		if s.pub != nil {
			s.pub.Publish(q.account, "raid.matched", map[string]any{
				"raid_code": tmpl.Code, "party_id": partyID, "raid_id": runID, "character_id": q.id,
			})
		}
		id := q.id
		notification.Send(ctx, s.notifier, s.log, notification.Notice{
			AccountID: q.account, CharacterID: &id, Kind: "raid_matched",
			DedupKey: "raid_matched:" + partyID.String() + ":" + q.id.String(),
			TitleKey: "notif.raid_matched.title", BodyKey: "notif.raid_matched.body",
			Vars:  map[string]any{"raid": tmpl.Code},
			Route: "character/raid", Category: "social", High: true,
		})
	}
}

// --- Board -------------------------------------------------------------------

// Advertise lists the actor's party for a raid, forming one led by the actor
// when they have none. The leader only; a new raid replaces the listing.
func (s *Service) Advertise(ctx context.Context, actor, account uuid.UUID, code string) (*Listing, error) {
	if s.parties == nil {
		return nil, ErrNotFound
	}
	tmpl, err := s.eligible(ctx, actor, account, code)
	if err != nil {
		return nil, err
	}
	if _, err := s.pool.Exec(ctx, `DELETE FROM raid.queue WHERE character_id = $1`, actor); err != nil {
		return nil, err
	}
	partyID, err := s.parties.LeadOrCreate(ctx, actor, account)
	if err != nil {
		if errors.Is(err, party.ErrNotLeader) {
			return nil, ErrNotLeader
		}
		return nil, err
	}
	var id uuid.UUID
	if err := s.pool.QueryRow(ctx, `
		INSERT INTO raid.listing (id, party_id, raid_code) VALUES ($1, $2, $3)
		ON CONFLICT (party_id) DO UPDATE SET raid_code = EXCLUDED.raid_code, created_at = now()
		RETURNING id`, uuid.New(), partyID, tmpl.Code).Scan(&id); err != nil {
		return nil, err
	}
	return s.listing(ctx, id)
}

// Unlist takes the actor's party off the board (leader only).
func (s *Service) Unlist(ctx context.Context, actor, account uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	partyID, leader, err := partyOf(ctx, s.pool, actor)
	if err != nil {
		return err
	}
	if leader != actor {
		return ErrNotLeader
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM raid.listing WHERE party_id = $1`, partyID)
	if err == nil && tag.RowsAffected() == 0 {
		err = ErrNotFound
	}
	return err
}

// Listings are the open listings (live parties with room), newest first;
// code filters by raid when set.
func (s *Service) Listings(ctx context.Context, actor, account uuid.UUID, code string) ([]Listing, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT l.id FROM raid.listing l JOIN social.party p ON p.id = l.party_id AND p.disbanded_at IS NULL
		 WHERE ($1 = '' OR l.raid_code = $1)
		   AND (SELECT count(*) FROM social.party_member m WHERE m.party_id = p.id AND m.left_at IS NULL) < p.size_cap
		 ORDER BY l.created_at DESC LIMIT $2`, code, MaxListed)
	if err != nil {
		return nil, err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, err
	}
	out := make([]Listing, 0, len(ids))
	for _, id := range ids {
		l, err := s.listing(ctx, id)
		if err != nil {
			continue // gone meanwhile
		}
		out = append(out, *l)
	}
	return out, nil
}

// JoinListing puts the actor in the advertising party, and in its open
// lobby when it has one. A party that is full comes off the board.
func (s *Service) JoinListing(ctx context.Context, actor, account, listingID uuid.UUID) (*Listing, error) {
	if s.parties == nil {
		return nil, ErrNotFound
	}
	l, err := s.listing(ctx, listingID)
	if err != nil {
		return nil, err
	}
	if _, err := s.eligible(ctx, actor, account, l.RaidCode); err != nil {
		return nil, err
	}
	if err := s.parties.Join(ctx, actor, account, l.PartyID); err != nil {
		switch {
		case errors.Is(err, party.ErrPartyFull), errors.Is(err, party.ErrNotFound):
			return nil, ErrListed
		case errors.Is(err, party.ErrAlreadyInParty):
			return nil, ErrInParty
		}
		return nil, err
	}
	_, _ = s.pool.Exec(ctx, `DELETE FROM raid.queue WHERE character_id = $1`, actor)
	var runID uuid.UUID
	if err := s.pool.QueryRow(ctx, `SELECT id FROM raid.run WHERE party_id = $1 AND state = 'forming'`, l.PartyID).
		Scan(&runID); err == nil {
		if _, err := s.Join(ctx, actor, account, runID); err != nil {
			s.log.InfoContext(ctx, "raid: listing joiner not in lobby", slog.String("err", err.Error()))
		}
	}
	if _, err := s.pool.Exec(ctx, `
		DELETE FROM raid.listing l USING social.party p
		 WHERE l.id = $1 AND p.id = l.party_id
		   AND (SELECT count(*) FROM social.party_member m WHERE m.party_id = p.id AND m.left_at IS NULL) >= p.size_cap`,
		listingID); err != nil {
		return nil, err
	}
	return l, nil
}

// listing loads one listing with its party's members (the leader first).
func (s *Service) listing(ctx context.Context, id uuid.UUID) (*Listing, error) {
	var l Listing
	err := s.pool.QueryRow(ctx, `
		SELECT l.id, l.raid_code, l.party_id, l.created_at
		  FROM raid.listing l JOIN social.party p ON p.id = l.party_id AND p.disbanded_at IS NULL
		 WHERE l.id = $1`, id).Scan(&l.ID, &l.RaidCode, &l.PartyID, &l.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrListed
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT c.name, c.level, c.class_code
		  FROM social.party_member m
		  JOIN social.party p ON p.id = m.party_id
		  JOIN character.character c ON c.id = m.character_id
		 WHERE m.party_id = $1 AND m.left_at IS NULL
		 ORDER BY (c.id = p.leader_character_id) DESC, m.joined_at`, l.PartyID)
	if err != nil {
		return nil, err
	}
	l.Members, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (ListingMember, error) {
		var m ListingMember
		err := r.Scan(&m.Name, &m.Level, &m.ClassCode)
		return m, err
	})
	return &l, err
}

// partyListing is the actor's party's own listing, if any.
func (s *Service) partyListing(ctx context.Context, partyID uuid.UUID) (*Listing, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM raid.listing WHERE party_id = $1`, partyID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.listing(ctx, id)
}

// eligible checks the actor owns the character, the raid exists and the
// character's level reaches it.
func (s *Service) eligible(ctx context.Context, actor, account uuid.UUID, code string) (dungeon.Template, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return dungeon.Template{}, err
	}
	tmpl, err := s.template(code)
	if err != nil {
		return tmpl, err
	}
	c, err := s.chars.GetOwnedInternal(ctx, actor)
	if err != nil {
		return tmpl, err
	}
	if c.Level < tmpl.MinLevel {
		return tmpl, ErrLevelRequirement
	}
	return tmpl, nil
}
