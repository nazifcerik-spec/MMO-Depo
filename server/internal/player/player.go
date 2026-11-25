// Package player serves what a character shows others (DD-371) and what it is
// doing now (DD-370). Read only.
package player

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/shop"
)

// ErrNotFound: no such character, or one that is not shown.
var ErrNotFound = errors.New("player: not found")

// Action is a character's running action (DD-263): its kind and what it is
// about, as party members see it (DD-301). A raid lobby has no times yet.
type Action struct {
	Kind      string     `json:"kind"`
	Detail    string     `json:"detail"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndsAt    *time.Time `json:"ends_at,omitempty"`
}

// Gear is a worn item as others see it: no instance id, no enchant.
type Gear struct {
	TemplateCode string `json:"template_code"`
	ItemClass    string `json:"item_class"`
	ItemSubclass string `json:"item_subclass"`
	Rarity       string `json:"rarity"`
	Prefix       string `json:"prefix,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
}

// Skill is a profession level.
type Skill struct {
	Code  string `json:"code"`
	Level int    `json:"level"`
}

// Pet is the active pet.
type Pet struct {
	Species string `json:"species"`
	Level   int    `json:"level"`
}

// Profile is a character's public page. Nothing of the account is in it but
// the look the player chose to show (avatar and background).
type Profile struct {
	CharacterID       uuid.UUID       `json:"character_id"`
	Name              string          `json:"name"`
	RaceCode          string          `json:"race_code"`
	ClassCode         string          `json:"class_code"`
	SpecCode          *string         `json:"spec_code,omitempty"`
	Level             int             `json:"level"`
	Title             *string         `json:"title,omitempty"`
	Avatar            string          `json:"avatar"`
	Background        string          `json:"background"`
	GuildName         *string         `json:"guild_name,omitempty"`
	GuildTag          *string         `json:"guild_tag,omitempty"`
	Score             int64           `json:"score"`
	Rank              int             `json:"rank"` // 0: not counted yet
	AchievementPoints int             `json:"achievement_points"`
	MasteryTiers      int             `json:"mastery_tiers"`
	DungeonsCleared   int             `json:"dungeons_cleared"`
	RaidsCleared      int             `json:"raids_cleared"`
	Professions       []Skill         `json:"professions"`
	Equipment         map[string]Gear `json:"equipment"`
	Pet               *Pet            `json:"pet,omitempty"`
}

type Service struct{ pool *pgxpool.Pool }

func NewService(pool *pgxpool.Pool) *Service { return &Service{pool: pool} }

// Action returns what the character is doing; nil when idle.
func (s *Service) Action(ctx context.Context, id uuid.UUID) (*Action, error) {
	var a Action
	err := s.pool.QueryRow(ctx, `
		SELECT 'gathering' AS kind, kind::text, started_at, expires_at FROM activity.contract
		 WHERE character_id = $1 AND state = 'active'
		UNION ALL SELECT 'hunt', mob_code::text, started_at, expires_at FROM hunt.contract
		 WHERE character_id = $1 AND state = 'active'
		UNION ALL SELECT 'dungeon', dungeon_code::text, started_at, expires_at FROM dungeon.contract
		 WHERE character_id = $1 AND state = 'active'
		UNION ALL SELECT 'raid', r.raid_code::text, r.started_at, r.expires_at
		  FROM raid.member m JOIN raid.run r ON r.id = m.run_id
		 WHERE m.character_id = $1 AND m.active
		UNION ALL SELECT 'travel', to_region_code::text, started_at, arrives_at FROM world.travel
		 WHERE character_id = $1 AND state = 'in_progress'
		UNION ALL SELECT 'crafting', COALESCE(recipe->>'ProfessionCode', ''), started_at, ready_at FROM crafting.job
		 WHERE character_id = $1 AND state = 'active'
		LIMIT 1`, id).Scan(&a.Kind, &a.Detail, &a.StartedAt, &a.EndsAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Profile returns the public page of the character called name. A character
// being deleted or on a locked account has none, as in the ranking.
func (s *Service) Profile(ctx context.Context, name string) (*Profile, error) {
	p := Profile{Professions: []Skill{}, Equipment: map[string]Gear{}}
	err := s.pool.QueryRow(ctx, `
		SELECT c.id, c.name, c.race_code, c.class_code, c.spec_code, c.level, c.active_title_code,
		       COALESCE(f.avatar, $2), COALESCE(f.background, $3), g.name, g.tag, c.achievement_points,
		       COALESCE(s.score, 0), COALESCE(s.mastery_tiers, 0), COALESCE(s.dungeons_cleared, 0),
		       COALESCE(s.raids_cleared, 0),
		       CASE WHEN s.character_id IS NULL THEN 0 ELSE (
		         SELECT count(*) + 1 FROM progression.character_score o
		          WHERE o.score > s.score OR (o.score = s.score AND o.character_id < s.character_id)) END
		  FROM character.character c
		  JOIN identity.account a ON a.id = c.account_id AND a.state = 'active'
		  LEFT JOIN cosmetic.profile f ON f.account_id = a.id
		  LEFT JOIN social.guild_member m ON m.character_id = c.id
		  LEFT JOIN social.guild g ON g.id = m.guild_id AND g.disbanded_at IS NULL
		  LEFT JOIN progression.character_score s ON s.character_id = c.id
		 WHERE c.name_normalized = $1 AND c.state = 'active'`,
		character.NameKey(name), shop.AvatarCode(1), shop.BackgroundCode(1)).
		Scan(&p.CharacterID, &p.Name, &p.RaceCode, &p.ClassCode, &p.SpecCode, &p.Level, &p.Title,
			&p.Avatar, &p.Background, &p.GuildName, &p.GuildTag, &p.AchievementPoints,
			&p.Score, &p.MasteryTiers, &p.DungeonsCleared, &p.RaidsCleared, &p.Rank)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT i.equipment_slot, i.template_code, t.item_class, COALESCE(t.item_subclass, ''), i.rarity,
		       COALESCE(i.prefix_code, ''), COALESCE(i.suffix_code, '')
		  FROM inventory.item_instance i JOIN content.item_template t ON t.code = i.template_code
		 WHERE i.owner_character_id = $1 AND i.equipment_slot IS NOT NULL`, p.CharacterID)
	if err != nil {
		return nil, err
	}
	var slot string
	var g Gear
	if _, err = pgx.ForEachRow(rows, []any{&slot, &g.TemplateCode, &g.ItemClass, &g.ItemSubclass, &g.Rarity, &g.Prefix, &g.Suffix},
		func() error { p.Equipment[slot] = g; return nil }); err != nil {
		return nil, err
	}
	rows, err = s.pool.Query(ctx, `
		SELECT profession_code, level FROM crafting.character_profession
		 WHERE character_id = $1 ORDER BY level DESC, profession_code`, p.CharacterID)
	if err != nil {
		return nil, err
	}
	var k Skill
	if _, err = pgx.ForEachRow(rows, []any{&k.Code, &k.Level},
		func() error { p.Professions = append(p.Professions, k); return nil }); err != nil {
		return nil, err
	}
	var pet Pet
	err = s.pool.QueryRow(ctx, `SELECT species_code, level FROM pet.pet WHERE character_id = $1 AND active`, p.CharacterID).
		Scan(&pet.Species, &pet.Level)
	switch {
	case err == nil:
		p.Pet = &pet
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}
	return &p, nil
}

// Handler serves GET /v1/players/{name} and GET /v1/character/{id}/action.
type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Action(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "invalid id")
		return
	}
	a, err := h.svc.Action(r.Context(), id)
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"action": a})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if dec, err := url.PathUnescape(name); err == nil { // chi matches the escaped path
		name = dec
	}
	p, err := h.svc.Profile(r.Context(), name)
	switch {
	case errors.Is(err, ErrNotFound):
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "player not found")
	case err != nil:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
	default:
		apierr.WriteJSON(w, http.StatusOK, p)
	}
}
