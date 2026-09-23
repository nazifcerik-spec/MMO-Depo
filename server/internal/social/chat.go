package social

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ChatMessage struct {
	ID       int64      `json:"id"`
	Channel  string     `json:"channel"`
	FromID   uuid.UUID  `json:"from_character_id"`
	FromName string     `json:"from_name"`
	TargetID *uuid.UUID `json:"target_character_id,omitempty"`
	Body     string     `json:"body"`
	SentAt   time.Time  `json:"sent_at"`
}

var chatChannels = map[string]bool{"global": true, "guild": true, "party": true, "whisper": true}

// channelRef resolves the guild/party the actor chats into, or nil for
// global/whisper.
func channelRef(ctx context.Context, q querier, actor uuid.UUID, channel string) (*uuid.UUID, error) {
	var sql string
	switch channel {
	case "global", "whisper":
		return nil, nil
	case "guild":
		sql = `SELECT m.guild_id FROM social.guild_member m
		         JOIN social.guild g ON g.id = m.guild_id AND g.disbanded_at IS NULL
		        WHERE m.character_id = $1`
	case "party":
		sql = `SELECT m.party_id FROM social.party_member m
		         JOIN social.party p ON p.id = m.party_id AND p.disbanded_at IS NULL
		        WHERE m.character_id = $1 AND m.left_at IS NULL`
	default:
		return nil, ErrBadChannel
	}
	var id uuid.UUID
	err := q.QueryRow(ctx, sql, actor).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotInChannel
	}
	return &id, err
}

// SendChat stores and fans out one message.
func (s *Service) SendChat(ctx context.Context, actor, account uuid.UUID, channel, targetName, body string) (*ChatMessage, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	if !chatChannels[channel] {
		return nil, ErrBadChannel
	}
	body, err := CleanText(body, MaxChatBody)
	if err != nil {
		return nil, err
	}
	var muted bool
	if err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM ops.account_sanction
		                WHERE account_id = $1 AND kind = 'mute' AND revoked_at IS NULL
		                  AND (until IS NULL OR until > now()))`, account).Scan(&muted); err != nil {
		return nil, err
	}
	if muted {
		return nil, ErrMuted
	}
	ref, err := channelRef(ctx, s.pool, actor, channel)
	if err != nil {
		return nil, err
	}
	msg := &ChatMessage{Channel: channel, FromID: actor, Body: body}
	var target charRef
	if channel == "whisper" {
		if target, err = characterByName(ctx, s.pool, targetName); err != nil {
			return nil, err
		}
		if target.ID == actor {
			return nil, ErrSelfTarget
		}
		if blocked, err := blockedEither(ctx, s.pool, actor, target.ID); err != nil {
			return nil, err
		} else if blocked {
			return nil, ErrBlocked
		}
		msg.TargetID = &target.ID
	}
	if err := s.pool.QueryRow(ctx, `
		INSERT INTO social.chat_message (channel_kind, channel_ref_id, sender_character_id, target_character_id, body)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sent_at, (SELECT name FROM character.character WHERE id = $3)`,
		channel, ref, actor, msg.TargetID, body).Scan(&msg.ID, &msg.SentAt, &msg.FromName); err != nil {
		return nil, err
	}
	s.fanoutChat(ctx, msg, ref, account, target.Account)
	return msg, nil
}

func (s *Service) fanoutChat(ctx context.Context, msg *ChatMessage, ref *uuid.UUID, sender, target uuid.UUID) {
	if s.pub == nil {
		return
	}
	switch msg.Channel {
	case "global":
		s.pub.BroadcastKind("chat.message", msg)
		return
	case "whisper":
		s.pub.Publish(target, "chat.message", msg)
		s.pub.Publish(sender, "chat.message", msg)
		return
	}
	sql := `SELECT DISTINCT c.account_id FROM social.guild_member m
	          JOIN character.character c ON c.id = m.character_id WHERE m.guild_id = $1`
	if msg.Channel == "party" {
		sql = `SELECT DISTINCT c.account_id FROM social.party_member m
		         JOIN character.character c ON c.id = m.character_id
		        WHERE m.party_id = $1 AND m.left_at IS NULL`
	}
	rows, err := s.pool.Query(ctx, sql, *ref)
	if err != nil {
		s.log.WarnContext(ctx, "social: chat fanout lookup failed", "err", err.Error())
		return
	}
	accounts, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return
	}
	for _, a := range accounts {
		s.pub.Publish(a, "chat.message", msg)
	}
}

// ChatHistory returns the latest messages of a channel (oldest first),
// hiding senders the actor has blocked.
func (s *Service) ChatHistory(ctx context.Context, actor, account uuid.UUID, channel, targetName string) ([]ChatMessage, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	if !chatChannels[channel] {
		return nil, ErrBadChannel
	}
	ref, err := channelRef(ctx, s.pool, actor, channel)
	if err != nil {
		return nil, err
	}
	const cols = `SELECT m.id, m.channel_kind, m.sender_character_id, c.name, m.target_character_id, m.body, m.sent_at
	                FROM social.chat_message m JOIN character.character c ON c.id = m.sender_character_id `
	const notBlocked = ` AND m.sender_character_id NOT IN
	      (SELECT target_character_id FROM social.ignore_list WHERE owner_character_id = $1) `
	var rows pgx.Rows
	switch channel {
	case "global":
		rows, err = s.pool.Query(ctx, cols+`WHERE m.channel_kind = 'global'`+notBlocked+
			`ORDER BY m.sent_at DESC LIMIT $2`, actor, ChatHistoryLimit)
	case "whisper":
		target, terr := characterByName(ctx, s.pool, targetName)
		if terr != nil {
			return nil, terr
		}
		rows, err = s.pool.Query(ctx, cols+`
			WHERE m.channel_kind = 'whisper'
			  AND ((m.sender_character_id = $1 AND m.target_character_id = $2)
			    OR (m.sender_character_id = $2 AND m.target_character_id = $1))
			ORDER BY m.sent_at DESC LIMIT $3`, actor, target.ID, ChatHistoryLimit)
	default:
		rows, err = s.pool.Query(ctx, cols+`WHERE m.channel_kind = $3 AND m.channel_ref_id = $4`+notBlocked+
			`ORDER BY m.sent_at DESC LIMIT $2`, actor, ChatHistoryLimit, channel, *ref)
	}
	if err != nil {
		return nil, err
	}
	msgs, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (ChatMessage, error) {
		var m ChatMessage
		err := r.Scan(&m.ID, &m.Channel, &m.FromID, &m.FromName, &m.TargetID, &m.Body, &m.SentAt)
		return m, err
	})
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}
