package character

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/otter/server/internal/leveling"
)

// --- Transactional progression (DD-250) ------------------------------------
//
// Every write locks the character row (FOR UPDATE) first, so concurrent
// grants, promotions and allocations serialize; rewards are either derived
// from the level (points) or recorded with ON CONFLICT DO NOTHING keys
// (events, titles), so nothing can be granted twice.

type progState struct {
	class, spec       string
	level             int16
	xp                int64
	overflow, postcap int64
	alloc             leveling.Stats
}

func lockState(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*progState, error) {
	var st progState
	var spec *string
	err := tx.QueryRow(ctx, `
		SELECT class_code, spec_code, level, xp, promotion_overflow_xp, postcap_xp,
		       attr_str, attr_vit, attr_dex, attr_int, attr_wis, attr_luck
		  FROM character.character WHERE id = $1 FOR UPDATE`, id).
		Scan(&st.class, &spec, &st.level, &st.xp, &st.overflow, &st.postcap,
			&st.alloc.Strength, &st.alloc.Vitality, &st.alloc.Dexterity, &st.alloc.Intelligence,
			&st.alloc.Wisdom, &st.alloc.Luck)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	st.spec = st.class
	if spec != nil && *spec != "" {
		st.spec = *spec
	}
	return &st, nil
}

// event records one progression event once per (kind, ref).
func event(ctx context.Context, tx pgx.Tx, id uuid.UUID, kind, ref string, level int, meta map[string]any) (bool, error) {
	if meta == nil {
		meta = map[string]any{}
	}
	raw, err := json.Marshal(meta)
	if err != nil {
		return false, err
	}
	tag, err := tx.Exec(ctx, `
		INSERT INTO character.progression_event (character_id, kind, ref, level, metadata)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (character_id, kind, ref) DO NOTHING`,
		id, kind, ref, level, raw)
	return err == nil && tag.RowsAffected() == 1, err
}

func grantTitle(ctx context.Context, tx pgx.Tx, id uuid.UUID, title, source string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO progression.character_title (character_id, title_code, source_ref_kind)
		SELECT $1::uuid, $2::varchar, $3::varchar WHERE EXISTS (SELECT 1 FROM content.title_template WHERE code = $2::varchar)
		ON CONFLICT (character_id, title_code) DO NOTHING`, id, title, source)
	return err
}

var unlockEvent = map[string]string{
	"skill": EventSkillUnlocked, "passive": EventPassiveUnlocked,
	"awakening": EventAwakening, "ultimate": EventUltimate,
}

// onLevelUp is the one level-up pipeline (DD-250 §31): for every level
// gained it records the registry rewards (points, titles, promotion, capstone,
// mastery) and the unlocks of the class path at that level. Each record is
// keyed, so replays and retries add nothing.
// grantAchievement unlocks a mastery achievement once and adds its points
// (no gold reward: the ledger is not part of this transaction). Unknown
// codes are skipped (content not loaded yet).
func grantAchievement(ctx context.Context, tx pgx.Tx, id uuid.UUID, code string) error {
	if code == "" {
		return nil
	}
	tag, err := tx.Exec(ctx, `
		INSERT INTO progression.achievement_progress (character_id, achievement_code, progress_current, unlocked_at)
		SELECT $1, code, progress_target, now() FROM content.achievement_template WHERE code = $2::varchar
		ON CONFLICT (character_id, achievement_code) DO UPDATE
		   SET progress_current = EXCLUDED.progress_current, unlocked_at = EXCLUDED.unlocked_at
		 WHERE progression.achievement_progress.unlocked_at IS NULL`, id, code)
	if err != nil || tag.RowsAffected() == 0 {
		return err
	}
	_, err = tx.Exec(ctx, `
		UPDATE character.character SET achievement_points = achievement_points +
		       (SELECT point_value FROM content.achievement_template WHERE code = $2::varchar)
		 WHERE id = $1`, id, code)
	return err
}

func onLevelUp(ctx context.Context, tx pgx.Tx, id uuid.UUID, spec string, levels []int, cfg *leveling.Config, rec *[]string) error {
	path := cfg.Path(spec)
	note := func(kind string, fresh bool, err error) error {
		if err == nil && fresh {
			*rec = append(*rec, kind)
		}
		return err
	}
	for _, l := range levels {
		ref := strconv.Itoa(l)
		ok, err := event(ctx, tx, id, EventLevelUp, ref, l, nil)
		if err := note(EventLevelUp, ok, err); err != nil {
			return err
		}
		for _, r := range cfg.RewardsAt(l) {
			switch r.Type {
			case leveling.RewardAttributePoints:
				ok, err = event(ctx, tx, id, EventAttributePoints, ref, l, map[string]any{"amount": r.Amount})
				err = note(EventAttributePoints, ok, err)
			case leveling.RewardTalentPoint:
				ok, err = event(ctx, tx, id, EventTalentPoint, ref, l, map[string]any{"amount": r.Amount})
				err = note(EventTalentPoint, ok, err)
			case leveling.RewardTitle:
				if err = grantTitle(ctx, tx, id, r.Ref, "level"); err == nil {
					ok, err = event(ctx, tx, id, EventTitleUnlocked, r.Ref, l, nil)
					err = note(EventTitleUnlocked, ok, err)
				}
			case leveling.RewardPromotion:
				if r.Ref == spec {
					ok, err = event(ctx, tx, id, EventPromotionAvail, "tier:"+spec, l, nil)
					err = note(EventPromotionAvail, ok, err)
				}
			case leveling.RewardCapstone:
				ok, err = event(ctx, tx, id, EventCapstone, "capstone", l, nil)
				err = note(EventCapstone, ok, err)
			case leveling.RewardMastery:
				if k, found := cfg.Class(spec); found && k.Mastery != nil {
					if err = grantTitle(ctx, tx, id, k.Mastery.Title, "class_mastery"); err == nil {
						err = grantAchievement(ctx, tx, id, k.Mastery.Achievement)
					}
					if err == nil {
						ok, err = event(ctx, tx, id, EventMastery, spec, l, map[string]any{
							"title": k.Mastery.Title, "badge": k.Mastery.Badge, "frame": k.Mastery.Frame,
							"achievement": k.Mastery.Achievement, "prestige_eligible": k.Mastery.PrestigeEligible})
						err = note(EventMastery, ok, err)
					}
				}
			}
			if err != nil {
				return err
			}
		}
		for _, u := range cfg.UnlocksAt(path, l) {
			kind := unlockEvent[u.Category]
			if kind == "" {
				kind = EventSkillUnlocked
			}
			ok, err := event(ctx, tx, id, kind, u.Code, l, map[string]any{"spec": u.Spec})
			if err := note(kind, ok, err); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *PgRepo) GrantXP(ctx context.Context, id uuid.UUID, amount int64, src XPSource, cfg *leveling.Config) (*Grant, error) {
	if amount < 0 {
		return nil, ErrInvalidXP
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	st, err := lockState(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if src.ID != "" {
		// Under the row lock: a concurrent twin waits here, then sees the row.
		tag, err := tx.Exec(ctx, `
			INSERT INTO character.xp_grant (character_id, source, source_id, amount)
			VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`, id, src.Kind, src.ID, amount)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			_ = tx.Rollback(ctx)
			c, err := r.Get(ctx, id)
			if err != nil {
				return nil, err
			}
			return &Grant{Character: c, FromLevel: int(c.Level), Replayed: true,
				PromotionAvailable: cfg.PromotionPending(c.Spec(), int(c.Level))}, nil
		}
	}
	from := leveling.State{Spec: st.spec, Level: int(st.level), XP: st.xp, Overflow: st.overflow, PostCap: st.postcap}
	to, gained := cfg.Advance(from, amount)
	if _, err := tx.Exec(ctx, `
		UPDATE character.character
		   SET level = $2, xp = $3, promotion_overflow_xp = $4, postcap_xp = $5,
		       updated_at = now(), version = version + 1
		 WHERE id = $1`, id, to.Level, to.XP, to.Overflow, to.PostCap); err != nil {
		return nil, err
	}
	var recorded []string
	if err := onLevelUp(ctx, tx, id, st.spec, gained, cfg, &recorded); err != nil {
		return nil, err
	}
	ref := src.Kind + ":" + src.ID
	if src.ID == "" {
		ref = src.Kind + ":" + uuid.NewString()
	}
	meta := map[string]any{"amount": amount, "source": src.Kind}
	for k, v := range src.Metadata {
		meta[k] = v
	}
	if _, err := event(ctx, tx, id, EventXPGained, ref, to.Level, meta); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Grant{Character: c, FromLevel: int(st.level), LevelsGained: gained, Events: recorded,
		PromotionAvailable: cfg.PromotionPending(st.spec, to.Level)}, nil
}

// Promote applies a promotion choice at a reached gate, records the unlocks
// of the new class at the current level, then pays out the held overflow.
func (r *PgRepo) Promote(ctx context.Context, id uuid.UUID, choice string, cfg *leveling.Config) (*Grant, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	st, err := lockState(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	valid := false
	for _, k := range cfg.Choices(st.spec, int(st.level)) {
		valid = valid || k.Code == choice
	}
	if !valid {
		return nil, ErrInvalidPromotion
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO character.class_history (character_id, from_spec, to_spec, level, reason)
		VALUES ($1, $2, $3, $4, 'promotion')`, id, st.spec, choice, st.level); err != nil {
		return nil, err
	}
	var recorded []string
	if ok, err := event(ctx, tx, id, EventPromotionSelected, choice, int(st.level), map[string]any{"from": st.spec}); err != nil {
		return nil, err
	} else if ok {
		recorded = append(recorded, EventPromotionSelected)
	}
	path := cfg.Path(choice)
	for _, u := range cfg.UnlockedBy(path, int(st.level)) {
		if u.Spec != choice {
			continue
		}
		kind := unlockEvent[u.Category]
		if kind == "" {
			kind = EventSkillUnlocked
		}
		if ok, err := event(ctx, tx, id, kind, u.Code, int(st.level), map[string]any{"spec": u.Spec}); err != nil {
			return nil, err
		} else if ok {
			recorded = append(recorded, kind)
		}
	}
	// The held XP now flows through the same engine (it may reach the next gate).
	to, gained := cfg.Advance(leveling.State{Spec: choice, Level: int(st.level), XP: st.xp, PostCap: st.postcap}, st.overflow)
	if _, err := tx.Exec(ctx, `
		UPDATE character.character
		   SET spec_code = $2, level = $3, xp = $4, promotion_overflow_xp = $5, postcap_xp = $6,
		       updated_at = now(), version = version + 1
		 WHERE id = $1`, id, choice, to.Level, to.XP, to.Overflow, to.PostCap); err != nil {
		return nil, err
	}
	if err := onLevelUp(ctx, tx, id, choice, gained, cfg, &recorded); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Grant{Character: c, FromLevel: int(st.level), LevelsGained: gained, Events: recorded,
		PromotionAvailable: cfg.PromotionPending(choice, to.Level)}, nil
}

func (r *PgRepo) AllocateAttributes(ctx context.Context, id uuid.UUID, add leveling.Stats, cfg *leveling.Config) (*Character, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	st, err := lockState(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if st.alloc.Total()+add.Total() > cfg.AttributePointsEarned(int(st.level)) {
		return nil, ErrNotEnoughPoints
	}
	n := st.alloc
	for _, a := range leveling.Attributes {
		n.Add(a, add.Get(a))
	}
	if _, err := tx.Exec(ctx, `
		UPDATE character.character
		   SET attr_str = $2, attr_vit = $3, attr_dex = $4, attr_int = $5, attr_wis = $6, attr_luck = $7,
		       updated_at = now(), version = version + 1
		 WHERE id = $1`, id, n.Strength, n.Vitality, n.Dexterity, n.Intelligence, n.Wisdom, n.Luck); err != nil {
		return nil, err
	}
	meta := map[string]any{}
	for _, a := range leveling.Attributes {
		if v := add.Get(a); v > 0 {
			meta[a] = v
		}
	}
	if _, err := event(ctx, tx, id, EventAttributesSpent, uuid.NewString(), int(st.level), meta); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *PgRepo) Talents(ctx context.Context, id uuid.UUID) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT talent_code, rank FROM character.talent WHERE character_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var code string
		var rank int
		if err := rows.Scan(&code, &rank); err != nil {
			return nil, err
		}
		out[code] = rank
	}
	return out, rows.Err()
}

// LearnTalent spends one point: points available, rank below max, the tree
// belongs to the character's class, capstones need the capstone level and a
// full tree.
func (r *PgRepo) LearnTalent(ctx context.Context, id uuid.UUID, code string, cfg *leveling.Config) (int, error) {
	t, tree, ok := cfg.Talent(code)
	if !ok {
		return 0, ErrUnknownTalent
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	st, err := lockState(ctx, tx, id)
	if err != nil {
		return 0, err
	}
	if !contains(cfg.Path(st.spec), tree.Class) {
		return 0, ErrUnknownTalent
	}
	rows, err := tx.Query(ctx, `SELECT talent_code, rank FROM character.talent WHERE character_id = $1`, id)
	if err != nil {
		return 0, err
	}
	spent, inTree, rank := 0, 0, 0
	for rows.Next() {
		var c string
		var n int
		if err := rows.Scan(&c, &n); err != nil {
			rows.Close()
			return 0, err
		}
		spent += n
		if other, ot, ok := cfg.Talent(c); ok && ot.Code == tree.Code && !other.Capstone {
			inTree += n
		}
		if c == code {
			rank = n
		}
	}
	rows.Close()
	if spent+1 > cfg.TalentPointsEarned(int(st.level)) {
		return 0, ErrNotEnoughPoints
	}
	if rank >= t.MaxRank {
		return 0, ErrTalentMaxed
	}
	if t.Capstone && (int(st.level) < cfg.Milestones.Capstone || inTree < cfg.TalentTreePoints) {
		return 0, ErrCapstoneLocked
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO character.talent (character_id, talent_code, rank) VALUES ($1, $2, 1)
		ON CONFLICT (character_id, talent_code) DO UPDATE SET rank = character.talent.rank + 1, updated_at = now()`,
		id, code); err != nil {
		return 0, err
	}
	if _, err := event(ctx, tx, id, EventTalentLearned, fmt.Sprintf("%s:%d", code, rank+1), int(st.level), nil); err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return rank + 1, nil
}

// ChangeSpec moves to another class node of the same tier the level allows
// (respec). Level, XP, attributes, talents and earned titles stay.
func (r *PgRepo) ChangeSpec(ctx context.Context, id uuid.UUID, spec string, cfg *leveling.Config) (*Character, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	st, err := lockState(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	cur, _ := cfg.Class(st.spec)
	next, ok := cfg.Class(spec)
	if !ok || cur == nil || next.Root() != cur.Root() || next.Tier() != cur.Tier() || next.Code == cur.Code ||
		next.RequiredLevel > int(st.level) {
		return nil, ErrInvalidSpec
	}
	if _, err := tx.Exec(ctx, `
		UPDATE character.character SET spec_code = $2, updated_at = now(), version = version + 1 WHERE id = $1`,
		id, spec); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO character.class_history (character_id, from_spec, to_spec, level, reason)
		VALUES ($1, $2, $3, $4, 'respec')`, id, st.spec, spec, st.level); err != nil {
		return nil, err
	}
	if _, err := event(ctx, tx, id, EventRespec, uuid.NewString(), int(st.level),
		map[string]any{"from": st.spec, "to": spec}); err != nil {
		return nil, err
	}
	for _, u := range cfg.UnlockedBy(cfg.Path(spec), int(st.level)) {
		kind := unlockEvent[u.Category]
		if kind == "" {
			kind = EventSkillUnlocked
		}
		if _, err := event(ctx, tx, id, kind, u.Code, int(st.level), map[string]any{"spec": u.Spec}); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *PgRepo) Events(ctx context.Context, id uuid.UUID, limit int) ([]ProgressionEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT kind, ref, level, metadata, created_at FROM character.progression_event
		 WHERE character_id = $1 ORDER BY id DESC LIMIT $2`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ProgressionEvent{}
	for rows.Next() {
		var e ProgressionEvent
		var raw []byte
		var at time.Time
		if err := rows.Scan(&e.Kind, &e.Ref, &e.Level, &raw, &at); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(raw, &e.Metadata)
		e.At = at.UTC().Format(time.RFC3339)
		out = append(out, e)
	}
	return out, rows.Err()
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
