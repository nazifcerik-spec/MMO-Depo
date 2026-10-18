package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
	"github.com/otter/server/internal/leveling"
	"github.com/otter/server/internal/skill"
)

// busyChecker: a character is busy while it has an active AFK activity,
// hunt or dungeon run (respec is refused then, DD-250 §35).
type busyChecker struct{ pool *pgxpool.Pool }

func (b busyChecker) Busy(ctx context.Context, id uuid.UUID) (bool, error) {
	var busy bool
	err := b.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM activity.contract WHERE character_id = $1 AND state = 'active')
		    OR EXISTS (SELECT 1 FROM hunt.contract WHERE character_id = $1 AND state = 'active')
		    OR EXISTS (SELECT 1 FROM dungeon.contract WHERE character_id = $1 AND state = 'active')`, id).Scan(&busy)
	return busy, err
}

// soloAdaptation is the Solo Adaptation modifier source: without a party of
// two or more, a support path's support power partly becomes attack.
type soloAdaptation struct{ pool *pgxpool.Pool }

func (s soloAdaptation) Modifiers(ctx context.Context, id uuid.UUID, sheet character.Sheet) ([]skill.Modifier, error) {
	cfg := leveling.Current()
	spec := sheet.Class
	if sheet.Spec != nil && *sheet.Spec != "" {
		spec = *sheet.Spec
	}
	if cfg.SupportPower(spec) <= 0 {
		return nil, nil // no query for non-support paths
	}
	inParty, err := s.grouped(ctx, id)
	if err != nil {
		return nil, err
	}
	return skill.FromLeveling(cfg.SoloAdaptationModifiers(spec, inParty)), nil
}

// grouped: an active party with at least one other active member.
func (s soloAdaptation) grouped(ctx context.Context, id uuid.UUID) (bool, error) {
	var grouped bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (
		  SELECT 1 FROM social.party_member me
		    JOIN social.party_member other ON other.party_id = me.party_id
		         AND other.character_id <> me.character_id AND other.left_at IS NULL
		   WHERE me.character_id = $1 AND me.left_at IS NULL)`, id).Scan(&grouped)
	return grouped, err
}

// progressionDev are dev-only tools: grant XP through the canonical path and
// change specialization (the respec boundary has no player route yet).
type progressionDev struct {
	chars *character.Service
	busy  character.BusyChecker
	gear  character.GearValidator
}

func (d *progressionDev) owned(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "invalid id")
		return uuid.Nil, false
	}
	if _, err := d.chars.GetOwned(r.Context(), id, httpctx.AccountID(r.Context())); err != nil {
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
		return uuid.Nil, false
	}
	return id, true
}

func (d *progressionDev) GrantXP(w http.ResponseWriter, r *http.Request) {
	id, ok := d.owned(w, r)
	if !ok {
		return
	}
	var body struct {
		Amount int64 `json:"amount"`
	}
	if err := decodeBody(r, &body); err != nil || body.Amount <= 0 {
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", "amount must be > 0", "amount")
		return
	}
	res, err := d.chars.GrantXP(r.Context(), id, body.Amount, character.XPSource{Kind: "dev"})
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	apierr.WriteJSON(w, http.StatusOK, res)
}

func (d *progressionDev) Respec(w http.ResponseWriter, r *http.Request) {
	id, ok := d.owned(w, r)
	if !ok {
		return
	}
	var body struct {
		Spec string `json:"spec"`
	}
	if err := decodeBody(r, &body); err != nil || body.Spec == "" {
		apierr.WriteField(w, r, http.StatusBadRequest, "E_VALIDATION", "spec required", "spec")
		return
	}
	c, removed, err := d.chars.ChangeSpecialization(r.Context(), id, body.Spec, d.busy, d.gear)
	switch {
	case errors.Is(err, character.ErrBusy):
		apierr.Write(w, r, http.StatusConflict, "E_BUSY_COMBAT", "not while in combat, AFK or a dungeon")
		return
	case errors.Is(err, character.ErrInvalidSpec):
		apierr.Write(w, r, http.StatusConflict, "E_SPEC_UNAVAILABLE", "that specialization is not available")
		return
	case err != nil:
		apierr.Write(w, r, http.StatusInternalServerError, "E_INTERNAL", err.Error())
		return
	}
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"spec_code": c.Spec(), "unequipped": removed})
}

func decodeBody(r *http.Request, dst any) error {
	return json.NewDecoder(io.LimitReader(r.Body, 4<<10)).Decode(dst)
}
