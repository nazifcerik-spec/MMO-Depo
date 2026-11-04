// Package leaderboard ranks characters (DD-357): a character score out of what
// the character has achieved, and a table per measure. No player-versus-player.
package leaderboard

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
)

// What a point of each measure is worth in the score.
const (
	PointsPerLevel       = 10
	PointsPerAchievement = 5
	PointsPerMasteryTier = 25
	PointsPerProfession  = 3
	PointsPerDungeon     = 20 // each different dungeon cleared
	PointsPerRaid        = 100
)

// RefreshSQL rewrites every active character's row and drops the others. Clears
// count once per dungeon, so repeating one cannot inflate the score. A locked
// account does not appear.
const RefreshSQL = `
WITH src AS (
  SELECT c.id, c.name, c.class_code, c.spec_code, c.level, c.achievement_points,
         COALESCE(m.tiers, 0) AS mastery_tiers, COALESCE(p.levels, 0) AS profession_levels,
         COALESCE(d.cleared, 0) AS dungeons_cleared, COALESCE(r.cleared, 0) AS raids_cleared
    FROM character.character c
    JOIN identity.account a ON a.id = c.account_id AND a.state = 'active'
    LEFT JOIN LATERAL (SELECT sum(current_tier) AS tiers FROM progression.character_mastery WHERE character_id = c.id) m ON true
    LEFT JOIN LATERAL (SELECT sum(level) AS levels FROM crafting.character_profession WHERE character_id = c.id) p ON true
    LEFT JOIN LATERAL (SELECT count(DISTINCT dungeon_code) AS cleared FROM dungeon.contract
                        WHERE character_id = c.id AND state = 'revealed' AND reveal_result->>'cleared' = 'true') d ON true
    LEFT JOIN LATERAL (SELECT count(DISTINCT dungeon_code) AS cleared FROM dungeon.weekly_lockout WHERE character_id = c.id) r ON true
   WHERE c.state = 'active'
), up AS (
  INSERT INTO progression.character_score (character_id, name, class_code, spec_code, level, achievement_points,
         mastery_tiers, profession_levels, dungeons_cleared, raids_cleared, score, updated_at)
  SELECT id, name, class_code, spec_code, level, achievement_points, mastery_tiers, profession_levels,
         dungeons_cleared, raids_cleared,
         level * 10 + achievement_points * 5 + mastery_tiers * 25 + profession_levels * 3
           + dungeons_cleared * 20 + raids_cleared * 100,
         now()
    FROM src
  ON CONFLICT (character_id) DO UPDATE SET name = EXCLUDED.name, class_code = EXCLUDED.class_code,
         spec_code = EXCLUDED.spec_code, level = EXCLUDED.level, achievement_points = EXCLUDED.achievement_points,
         mastery_tiers = EXCLUDED.mastery_tiers, profession_levels = EXCLUDED.profession_levels,
         dungeons_cleared = EXCLUDED.dungeons_cleared, raids_cleared = EXCLUDED.raids_cleared,
         score = EXCLUDED.score, updated_at = now()
  RETURNING 1
), gone AS (
  DELETE FROM progression.character_score s WHERE NOT EXISTS (SELECT 1 FROM src WHERE src.id = s.character_id)
)
SELECT count(*) FROM up`

// boards are the tables: code → column. Only these reach the SQL.
var boards = map[string]string{
	"score":        "score",
	"level":        "level",
	"achievements": "achievement_points",
	"mastery":      "mastery_tiers",
	"professions":  "profession_levels",
	"clears":       "dungeons_cleared + raids_cleared",
}

// Entry is one row of a table.
type Entry struct {
	Rank              int        `json:"rank"`
	CharacterID       uuid.UUID  `json:"character_id"`
	Name              string     `json:"name"`
	ClassCode         string     `json:"class_code"`
	SpecCode          *string    `json:"spec_code,omitempty"`
	Level             int16      `json:"level"`
	Value             int64      `json:"value"`
	Score             int64      `json:"score"`
	AchievementPoints int        `json:"achievement_points"`
	MasteryTiers      int        `json:"mastery_tiers"`
	ProfessionLevels  int        `json:"profession_levels"`
	DungeonsCleared   int        `json:"dungeons_cleared"`
	RaidsCleared      int        `json:"raids_cleared"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

var ErrUnknownBoard = errors.New("leaderboard: unknown board")

type Service struct{ pool *pgxpool.Pool }

func NewService(pool *pgxpool.Pool) *Service { return &Service{pool: pool} }

// Refresh recomputes the scores (the worker's job; tests call it directly).
func (s *Service) Refresh(ctx context.Context) (int64, error) {
	var n int64
	err := s.pool.QueryRow(ctx, RefreshSQL).Scan(&n)
	return n, err
}

// Board returns a page of a table, optionally for one class.
func (s *Service) Board(ctx context.Context, board, class string, limit, offset int) ([]Entry, error) {
	col, ok := boards[board]
	if !ok {
		return nil, ErrUnknownBoard
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.pool.Query(ctx, `
		SELECT character_id, name, class_code, spec_code, level, (`+col+`)::bigint AS value, score, achievement_points,
		       mastery_tiers, profession_levels, dungeons_cleared, raids_cleared, updated_at
		  FROM progression.character_score
		 WHERE $1 = '' OR class_code = $1
		 ORDER BY value DESC, score DESC, character_id
		 LIMIT $2 OFFSET $3`, class, limit, offset)
	if err != nil {
		return nil, err
	}
	out := []Entry{}
	i := 0
	defer rows.Close()
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.CharacterID, &e.Name, &e.ClassCode, &e.SpecCode, &e.Level, &e.Value, &e.Score,
			&e.AchievementPoints, &e.MasteryTiers, &e.ProfessionLevels, &e.DungeonsCleared, &e.RaidsCleared, &e.UpdatedAt); err != nil {
			return nil, err
		}
		i++
		e.Rank = offset + i
		out = append(out, e)
	}
	return out, rows.Err()
}

// Mine returns a character's row and rank on a table; nil before the worker has
// counted it.
func (s *Service) Mine(ctx context.Context, board, class string, id uuid.UUID) (*Entry, error) {
	col, ok := boards[board]
	if !ok {
		return nil, ErrUnknownBoard
	}
	var e Entry
	err := s.pool.QueryRow(ctx, `
		SELECT character_id, name, class_code, spec_code, level, (`+col+`)::bigint AS value, score, achievement_points,
		       mastery_tiers, profession_levels, dungeons_cleared, raids_cleared, updated_at
		  FROM progression.character_score WHERE character_id = $1`, id).
		Scan(&e.CharacterID, &e.Name, &e.ClassCode, &e.SpecCode, &e.Level, &e.Value, &e.Score,
			&e.AchievementPoints, &e.MasteryTiers, &e.ProfessionLevels, &e.DungeonsCleared, &e.RaidsCleared, &e.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if class != "" && class != e.ClassCode {
		return nil, nil
	}
	// Ties rank the way the page orders them.
	err = s.pool.QueryRow(ctx, `
		SELECT count(*) + 1 FROM progression.character_score
		 WHERE ($1 = '' OR class_code = $1)
		   AND ((`+col+`) > $2 OR ((`+col+`) = $2 AND score > $3) OR ((`+col+`) = $2 AND score = $3 AND character_id < $4))`,
		class, e.Value, e.Score, id).Scan(&e.Rank)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Handler serves GET /v1/leaderboard?board=&class=&limit=&offset=&character_id=.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Board(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	board := q.Get("board")
	if board == "" {
		board = "score"
	}
	class := q.Get("class")
	if class != "" && !character.KnownClass(class) {
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "unknown class")
		return
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	list, err := h.svc.Board(r.Context(), board, class, limit, offset)
	if errors.Is(err, ErrUnknownBoard) {
		apierr.Write(w, r, http.StatusBadRequest, "E_VALIDATION", "unknown board")
		return
	}
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	out := map[string]any{"board": board, "entries": list, "points": map[string]int{
		"level": PointsPerLevel, "achievement": PointsPerAchievement, "mastery_tier": PointsPerMasteryTier,
		"profession_level": PointsPerProfession, "dungeon": PointsPerDungeon, "raid": PointsPerRaid,
	}}
	if raw := q.Get("character_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "character_id")
			return
		}
		me, err := h.svc.Mine(r.Context(), board, class, id)
		if err != nil {
			apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
			return
		}
		out["me"] = me
	}
	apierr.WriteJSON(w, http.StatusOK, out)
}
