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
