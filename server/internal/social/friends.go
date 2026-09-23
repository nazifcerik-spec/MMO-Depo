package social

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Friend struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	Level       int16     `json:"level"`
	ClassCode   string    `json:"class_code"`
	Since       time.Time `json:"since"`
}

type FriendRequest struct {
	ID        uuid.UUID `json:"id"`
	From      uuid.UUID `json:"from_character_id"`
	FromName  string    `json:"from_name"`
	CreatedAt time.Time `json:"created_at"`
}

type Blocked struct {
	CharacterID uuid.UUID `json:"character_id"`
	Name        string    `json:"name"`
	AddedAt     time.Time `json:"added_at"`
}

func (s *Service) Friends(ctx context.Context, actor, account uuid.UUID) ([]Friend, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT c.id, c.name, c.level, c.class_code, f.since
		  FROM social.friendship f
		  JOIN character.character c
		    ON c.id = CASE WHEN f.character_a = $1 THEN f.character_b ELSE f.character_a END
		 WHERE f.character_a = $1 OR f.character_b = $1
		 ORDER BY c.name`, actor)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Friend, error) {
		var f Friend
		err := r.Scan(&f.CharacterID, &f.Name, &f.Level, &f.ClassCode, &f.Since)
		return f, err
	})
}

func (s *Service) FriendRequests(ctx context.Context, actor, account uuid.UUID) ([]FriendRequest, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT r.id, r.requester_character_id, c.name, r.created_at
		  FROM social.friend_request r JOIN character.character c ON c.id = r.requester_character_id
		 WHERE r.target_character_id = $1 AND r.status = 'pending'
		 ORDER BY r.created_at DESC`, actor)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (FriendRequest, error) {
		var fr FriendRequest
		err := r.Scan(&fr.ID, &fr.From, &fr.FromName, &fr.CreatedAt)
		return fr, err
	})
}

// RequestFriend sends a request; if the target already asked us, the two
// requests meet and the friendship is created immediately.
func (s *Service) RequestFriend(ctx context.Context, actor, account uuid.UUID, targetName string) (accepted bool, err error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return false, err
	}
	target, err := characterByName(ctx, s.pool, targetName)
	if err != nil {
		return false, err
	}
	if target.ID == actor {
		return false, ErrSelfTarget
	}
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if blocked, err := blockedEither(ctx, tx, actor, target.ID); err != nil || blocked {
			if err == nil {
				err = ErrBlocked
			}
			return err
		}
		a, b := pair(actor, target.ID)
		var exists bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM social.friendship WHERE character_a = $1 AND character_b = $2)`,
			a, b).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return ErrAlreadyFriends
		}
		tag, err := tx.Exec(ctx, `
			UPDATE social.friend_request SET status = 'accepted'
			 WHERE requester_character_id = $1 AND target_character_id = $2 AND status = 'pending'`,
			target.ID, actor)
		if err != nil {
			return err
		}
		if tag.RowsAffected() > 0 {
			accepted = true
			return addFriendship(ctx, tx, actor, target.ID)
		}
		if err := checkFriendCap(ctx, tx, actor); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO social.friend_request (requester_character_id, target_character_id) VALUES ($1, $2)`,
			actor, target.ID)
		if isUnique(err) {
			return ErrAlreadyRequested
		}
		return err
	})
	if err != nil {
		return false, err
	}
	kind := "friend.request"
	if accepted {
		kind = "friend.accepted"
	}
	s.publish(target.Account, kind, map[string]any{"character_id": actor})
	return accepted, nil
}

func (s *Service) AcceptFriend(ctx context.Context, actor, account, requestID uuid.UUID) error {
	return s.answerFriend(ctx, actor, account, requestID, true)
}

func (s *Service) DeclineFriend(ctx context.Context, actor, account, requestID uuid.UUID) error {
	return s.answerFriend(ctx, actor, account, requestID, false)
}

func (s *Service) answerFriend(ctx context.Context, actor, account, requestID uuid.UUID, accept bool) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	var requester uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		status := "declined"
		if accept {
			status = "accepted"
		}
		err := tx.QueryRow(ctx, `
			UPDATE social.friend_request SET status = $3
			 WHERE id = $1 AND target_character_id = $2 AND status = 'pending'
			RETURNING requester_character_id`, requestID, actor, status).Scan(&requester)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil || !accept {
			return err
		}
		return addFriendship(ctx, tx, actor, requester)
	})
	if err != nil {
		return err
	}
	if accept {
		if a, err := accountOf(ctx, s.pool, requester); err == nil {
			s.publish(a, "friend.accepted", map[string]any{"character_id": actor})
		}
	}
	return nil
}

func (s *Service) RemoveFriend(ctx context.Context, actor, account, target uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	a, b := pair(actor, target)
	tag, err := s.pool.Exec(ctx, `DELETE FROM social.friendship WHERE character_a = $1 AND character_b = $2`, a, b)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFriends
	}
	return nil
}

func addFriendship(ctx context.Context, tx pgx.Tx, x, y uuid.UUID) error {
	if err := checkFriendCap(ctx, tx, x); err != nil {
		return err
	}
	if err := checkFriendCap(ctx, tx, y); err != nil {
		return err
	}
	a, b := pair(x, y)
	_, err := tx.Exec(ctx, `
		INSERT INTO social.friendship (character_a, character_b) VALUES ($1, $2) ON CONFLICT DO NOTHING`, a, b)
	return err
}

func checkFriendCap(ctx context.Context, q querier, c uuid.UUID) error {
	var n int
	if err := q.QueryRow(ctx, `
		SELECT count(*) FROM social.friendship WHERE character_a = $1 OR character_b = $1`, c).Scan(&n); err != nil {
		return err
	}
	if n >= FriendCap {
		return ErrFriendCap
	}
	return nil
}

// --- block list -------------------------------------------------------------

func (s *Service) BlockList(ctx context.Context, actor, account uuid.UUID) ([]Blocked, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT c.id, c.name, i.added_at FROM social.ignore_list i
		  JOIN character.character c ON c.id = i.target_character_id
		 WHERE i.owner_character_id = $1 ORDER BY c.name`, actor)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Blocked, error) {
		var b Blocked
		err := r.Scan(&b.CharacterID, &b.Name, &b.AddedAt)
		return b, err
	})
}

// Block adds the target to the actor's block list and severs any friendship
// or pending request between them.
func (s *Service) Block(ctx context.Context, actor, account uuid.UUID, targetName string) (uuid.UUID, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return uuid.Nil, err
	}
	target, err := characterByName(ctx, s.pool, targetName)
	if err != nil {
		return uuid.Nil, err
	}
	if target.ID == actor {
		return uuid.Nil, ErrSelfTarget
	}
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			INSERT INTO social.ignore_list (owner_character_id, target_character_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING`, actor, target.ID); err != nil {
			return err
		}
		a, b := pair(actor, target.ID)
		if _, err := tx.Exec(ctx, `DELETE FROM social.friendship WHERE character_a = $1 AND character_b = $2`, a, b); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			UPDATE social.friend_request SET status = 'cancelled'
			 WHERE status = 'pending'
			   AND ((requester_character_id = $1 AND target_character_id = $2)
			     OR (requester_character_id = $2 AND target_character_id = $1))`, actor, target.ID)
		return err
	})
	return target.ID, err
}

func (s *Service) Unblock(ctx context.Context, actor, account, target uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM social.ignore_list WHERE owner_character_id = $1 AND target_character_id = $2`, actor, target)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
