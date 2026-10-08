package combat

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/httpctx"
)

// Handler exposes the dev-only combat simulator. Used by the dev console to
// exercise the pipeline without going through hunting/dungeons (which arrive
// in Sprint 14+). Real player-facing combat is embedded in those flows.
type Handler struct {
	mobs    *Service
	charSvc *character.Service
	log     *slog.Logger
}

func NewHandler(mobs *Service, charSvc *character.Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{mobs: mobs, charSvc: charSvc, log: logger}
}

// SimulateRequest is the POST body: pick a mob + a character (whose derived
// sheet becomes the player Combatant) and optionally a seed for repeatable
// replays. Empty seed uses nanoseconds so a fresh call feels stochastic.
type SimulateRequest struct {
	CharacterID string `json:"character_id"`
	MobCode     string `json:"mob_code"`
	Seed        int64  `json:"seed,omitempty"`
	MaxTurns    int    `json:"max_turns,omitempty"`
}

// Simulate resolves one player-vs-mob encounter deterministically and returns
// the outcome + full turn log. Dev-only: never persists anything and does not
// grant rewards.
func (h *Handler) Simulate(w http.ResponseWriter, r *http.Request) {
	acct := httpctx.AccountID(r.Context())
	if acct == uuid.Nil {
		apierr.Write(w, r, http.StatusUnauthorized, "E_UNAUTHENTICATED", "auth required")
		return
	}

	var req SimulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	req.MobCode = strings.TrimSpace(req.MobCode)
	if req.MobCode == "" {
		apierr.Write(w, r, http.StatusBadRequest, "E_MISSING_MOB", "mob_code required")
		return
	}
	charID, err := uuid.Parse(strings.TrimSpace(req.CharacterID))
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_ID", "character_id: "+err.Error())
		return
	}

	// Ownership check (via GetOwned — 404 on foreign character; leak-safe).
	if _, err := h.charSvc.GetOwned(r.Context(), charID, acct); err != nil {
		apierr.Write(w, r, http.StatusNotFound, "E_NOT_FOUND", "character not found")
		return
	}

	// Player Combatant is built from the real sheet — this exercises the full
	// modifier stack from Sprint 11 (equipped items apply).
	sheet, err := h.charSvc.Sheet(r.Context(), charID, acct)
	if err != nil {
		apierr.Write(w, r, http.StatusInternalServerError, "E_SHEET", err.Error())
		return
	}
	player := FromSheet(sheet)
	player.HPCurrent = player.HPMax // a simulation starts fresh

	mob, ok := h.mobs.Catalog().Get(req.MobCode)
	if !ok {
		apierr.Write(w, r, http.StatusNotFound, "E_UNKNOWN_MOB", "mob_code not in catalog")
		return
	}

	seed := req.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	maxTurns := req.MaxTurns
	if maxTurns <= 0 {
		maxTurns = DefaultMaxTurns
	}

	session := &Session{
		Player:   player,
		Mob:      mob.AsCombatant(),
		MaxTurns: maxTurns,
	}
	result, err := session.Resolve(func(s int64) RNG { return activity.NewRNG(s) }, seed)
	if err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_COMBAT_INVALID", err.Error())
		return
	}

	apierr.WriteJSON(w, http.StatusOK, map[string]any{
		"seed":      seed,
		"max_turns": maxTurns,
		"player":    summarizeCombatant(player),
		"mob":       map[string]any{"code": mob.Code, "level": mob.Level, "hp_max": mob.HPMax},
		"outcome":   result.Outcome,
		"duration":  result.Duration,
		"player_hp": result.PlayerHPEnd,
		"mob_hp":    result.MobHPEnd,
		"turns":     result.Turns,
	})
}

// FromSheet is the player combatant for a computed sheet (hunts, dungeons,
// simulate). Fights start at the current HP (DD-218); evasion is agi/2
// capped at 40; Level feeds the level-gap rule (DD-229).
func FromSheet(s *character.Sheet) Combatant {
	return Combatant{
		Name:       "player", // stable label: turn logs never carry character names
		HPMax:      int32(s.Resources.HPMax),
		HPCurrent:  int32(s.Resources.HPCurrent),
		Attack:     int32(s.Derived.Attack),
		Defense:    int32(s.Derived.Defense),
		CritPct:    int32(s.Derived.CritPct),
		HitPct:     int32(s.Derived.HitPct),
		EvasionPct: min(int32(s.Primary.Agility)/2, 40),
		MPMax:      int32(s.Resources.MPMax),
		Level:      int16(s.Level),
	}
}

func summarizeCombatant(c Combatant) map[string]any {
	return map[string]any{
		"hp_max":      c.HPMax,
		"attack":      c.Attack,
		"defense":     c.Defense,
		"crit_pct":    c.CritPct,
		"hit_pct":     c.HitPct,
		"evasion_pct": c.EvasionPct,
	}
}

// Sanity: force fmt import used only in some builds.
var _ = fmt.Sprintf

// Mobs is the public mob catalog (GET /v1/mobs) for the client's hunt
// picker, sorted by level then code. Combat stats are public by design;
// no per-player data.
func (h *Handler) Mobs(w http.ResponseWriter, r *http.Request) {
	all := h.mobs.Catalog().All()
	sort.Slice(all, func(i, j int) bool {
		if all[i].Level != all[j].Level {
			return all[i].Level < all[j].Level
		}
		return all[i].Code < all[j].Code
	})
	out := make([]map[string]any, 0, len(all))
	for _, m := range all {
		out = append(out, map[string]any{
			"code":             m.Code,
			"display_name_key": m.DisplayNameKey,
			"level":            m.Level,
			"hp_max":           m.HPMax,
			"attack":           m.Attack,
			"defense":          m.Defense,
			"xp_reward":        m.XPReward,
			"gold_copper_min":  m.GoldCopperMin,
			"gold_copper_max":  m.GoldCopperMax,
		})
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apierr.WriteJSON(w, http.StatusOK, map[string]any{"count": len(out), "mobs": out})
}
