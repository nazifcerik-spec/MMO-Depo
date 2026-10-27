package partyxp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Granter pays a partner's share: profession XP for gathering and crafting
// kinds, character XP for hunts and dungeons. Must be idempotent per
// (character, ref) or tolerate being called once.
type Granter interface {
	GrantShare(ctx context.Context, kind string, characterID uuid.UUID, xp int64, ref uuid.UUID) error
}

// Publisher tells a partner what they received (wsgateway adapter satisfies).
type Publisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

// Service finds the earner's partners and pays their shares.
type Service struct {
	pool  *pgxpool.Pool
	grant Granter
	pub   Publisher
	log   *slog.Logger
}

func NewService(pool *pgxpool.Pool, grant Granter, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, grant: grant, log: logger}
}

func (s *Service) WithPublisher(p Publisher) *Service { s.pub = p; return s }

type partner struct {
	id      uuid.UUID
	account uuid.UUID
	name    string
	windows []Interval
}

// Share splits w among the earner's live party members whose jobs of the
// same kind overlap it, pays them and returns the earner's side. A failure
// to look partners up never costs the earner: the job then stays unshared.
func (s *Service) Share(ctx context.Context, w Work) Outcome {
	if w.XP <= 0 || !w.End.After(w.Start) {
		return Solo(w.XP)
	}
	partners, err := s.partners(ctx, w)
	if err != nil {
		s.log.WarnContext(ctx, "partyxp: partner lookup failed", "ref", w.Ref.String(), "err", err.Error())
		return Solo(w.XP)
	}
	if len(partners) == 0 {
		return Solo(w.XP)
	}
	windows := make([][]Interval, len(partners))
	for i, p := range partners {
		windows[i] = p.windows
	}
	keep, gifts := Split(w, windows)
	out := Outcome{Earned: w.XP, Kept: keep}
	for i, p := range partners {
		if gifts[i] <= 0 {
			continue
		}
		paid, err := s.pay(ctx, w, p, gifts[i])
		if err != nil {
			s.log.WarnContext(ctx, "partyxp: share not paid", "ref", w.Ref.String(),
				"partner", p.id.String(), "err", err.Error())
			continue
		}
		out.Partners = append(out.Partners, Share{CharacterID: p.id, Name: p.name, XP: gifts[i]})
		if paid && s.pub != nil {
			s.pub.Publish(p.account, "party.xp_share", map[string]any{
				"character_id": p.id, "from_character_id": w.Character, "kind": w.Kind, "xp": gifts[i],
			})
		}
	}
	return out
}

// pay records a share once and grants it; false = already paid.
func (s *Service) pay(ctx context.Context, w Work, p partner, xp int64) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		INSERT INTO social.party_xp_share (source_ref, recipient_character_id, giver_character_id, work_kind, xp)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`, w.Ref, p.id, w.Character, w.Kind, xp)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() == 0 {
		return false, nil
	}
	if s.grant == nil {
		return true, nil
	}
	return true, s.grant.GrantShare(ctx, w.Kind, p.id, xp, w.Ref)
}

// partners are the earner's live party members (not the earner) with their
// windows of the same kind of work overlapping w; members without one drop.
func (s *Service) partners(ctx context.Context, w Work) ([]partner, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT c.id, c.account_id, c.name
		  FROM social.party_member me
		  JOIN social.party p ON p.id = me.party_id AND p.disbanded_at IS NULL
		  JOIN social.party_member m ON m.party_id = p.id AND m.left_at IS NULL AND m.character_id <> me.character_id
		  JOIN character.character c ON c.id = m.character_id
		 WHERE me.character_id = $1 AND me.left_at IS NULL
		 ORDER BY m.joined_at`, w.Character)
	if err != nil {
		return nil, err
	}
	members, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (partner, error) {
		var p partner
		err := r.Scan(&p.id, &p.account, &p.name)
		return p, err
	})
	if err != nil || len(members) == 0 {
		return nil, err
	}
	ids := make([]uuid.UUID, len(members))
	for i, m := range members {
		ids[i] = m.id
	}
	sql, arg, err := windowsQuery(w.Kind)
	if err != nil {
		return nil, err
	}
	args := []any{ids, w.Start, w.End}
	if arg != "" {
		args = append(args, arg)
	}
	rows, err = s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	byID := map[uuid.UUID][]Interval{}
	for rows.Next() {
		var id uuid.UUID
		var iv Interval
		if err := rows.Scan(&id, &iv.Start, &iv.End); err != nil {
			rows.Close()
			return nil, err
		}
		byID[id] = append(byID[id], iv)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var out []partner
	for _, m := range members {
		if ws := byID[m.id]; len(ws) > 0 {
			m.windows = ws
			out = append(out, m)
		}
	}
	return out, nil
}

// windowsQuery selects (character, start, end) of the jobs of a kind for
// the characters $1 overlapping [$2, $3]; $4 is the gathering kind or the
// profession. Running jobs count until their planned end, finished ones
// until they were collected; stopped jobs earned nothing and do not count.
func windowsQuery(kind string) (string, string, error) {
	const contract = `
		SELECT character_id, started_at, LEAST(COALESCE(revealed_at, expires_at), expires_at) AS ended
		  FROM %s
		 WHERE character_id = ANY($1) AND state IN ('active', 'revealed') %s
		   AND started_at < $3 AND LEAST(COALESCE(revealed_at, expires_at), expires_at) > $2`
	switch {
	case kind == KindHunt:
		return fmt.Sprintf(contract, "hunt.contract", ""), "", nil
	case kind == KindDungeon:
		return fmt.Sprintf(contract, "dungeon.contract", ""), "", nil
	case strings.HasPrefix(kind, GatherPrefix):
		return fmt.Sprintf(contract, "activity.contract", "AND kind = $4"), strings.TrimPrefix(kind, GatherPrefix), nil
	case strings.HasPrefix(kind, CraftPrefix):
		return `
		SELECT character_id, started_at, LEAST(COALESCE(collected_at, ready_at), ready_at)
		  FROM crafting.job
		 WHERE character_id = ANY($1) AND recipe->>'ProfessionCode' = $4
		   AND started_at < $3 AND LEAST(COALESCE(collected_at, ready_at), ready_at) > $2`,
			strings.TrimPrefix(kind, CraftPrefix), nil
	}
	return "", "", fmt.Errorf("partyxp: unknown work kind %q", kind)
}
