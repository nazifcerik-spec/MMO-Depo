package raid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otter/server/internal/action"
	"github.com/otter/server/internal/activity"
	"github.com/otter/server/internal/boost"
	"github.com/otter/server/internal/combat"
	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/dungeon"
	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/item"
	"github.com/otter/server/internal/notification"
	"github.com/otter/server/internal/partyxp"
	"github.com/otter/server/internal/progression"
)

// Progression events of a clear (also "raid_cleared:<code>").
const EventRaidCleared = "raid_cleared"

// EggMinter puts a pet egg in the bag (item.Service).
type EggMinter interface {
	MintExact(ctx context.Context, characterID uuid.UUID, templateCode string, qty int32, sourceKind string,
		sourceRefID uuid.UUID) (*item.Instance, error)
}

// WithEggs mints the pet eggs cleared raids drop (DD-328).
func (s *Service) WithEggs(m EggMinter) *Service { s.eggs = m; return s }

// Publisher fans raid lobby changes out to live sessions (wsgateway adapter).
type Publisher interface {
	Publish(accountID uuid.UUID, kind string, payload any) int
}

// Service runs raids. Templates and mobs come from the dungeon and combat
// catalogs; the ports are the dungeon service's.
type Service struct {
	boost    boost.Applier // timed XP boosts (DD-320); nil = none
	pool     *pgxpool.Pool
	raids    dungeon.DungeonLookup
	mobs     dungeon.MobsLookup
	chars    dungeon.SheetFetcher
	xp       dungeon.XPGranter
	items    dungeon.ItemGranter
	txn      economy.TxnService
	hp       dungeon.HPWriter
	skills   dungeon.SkillSource
	bag      dungeon.BagRoom
	progress progression.Recorder
	notifier notification.Notifier
	pub      Publisher
	parties  Parties            // raid finding (DD-299); nil = no queue or board
	eggs     EggMinter          // pet eggs (DD-328); nil = none
	regions  dungeon.RegionGate // nil = every raid in reach (DD-302)
	now      func() time.Time
	log      *slog.Logger
	rngNew   func(int64) combat.RNG
}

func NewService(pool *pgxpool.Pool, raids dungeon.DungeonLookup, mobs dungeon.MobsLookup, chars dungeon.SheetFetcher,
	xp dungeon.XPGranter, items dungeon.ItemGranter, txn economy.TxnService, logger *slog.Logger,
	rngNew func(int64) combat.RNG) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{pool: pool, raids: raids, mobs: mobs, chars: chars, xp: xp, items: items, txn: txn,
		now: time.Now, log: logger, rngNew: rngNew}
}

func (s *Service) WithHPWriter(w dungeon.HPWriter) *Service      { s.hp = w; return s }
func (s *Service) WithSkills(src dungeon.SkillSource) *Service   { s.skills = src; return s }
func (s *Service) WithBagRoom(b dungeon.BagRoom) *Service        { s.bag = b; return s }
func (s *Service) WithRecorder(r progression.Recorder) *Service  { s.progress = r; return s }
func (s *Service) WithNotifier(n notification.Notifier) *Service { s.notifier = n; return s }
func (s *Service) WithPublisher(p Publisher) *Service            { s.pub = p; return s }

// WithRegionGate makes raids regional (DD-302): a lobby is opened, a queue
// joined and a party advertised from the raid's region only.
func (s *Service) WithRegionGate(g dungeon.RegionGate) *Service { s.regions = g; return s }
func (s *Service) WithClock(now func() time.Time) *Service      { s.now = now; return s }

// Catalog lists the raid templates, lowest level first.
func (s *Service) Catalog() []dungeon.Template {
	var out []dungeon.Template
	for _, t := range s.raids.Catalog().All() {
		if t.IsRaid() {
			out = append(out, t)
		}
	}
	return out
}

func (s *Service) template(code string) (dungeon.Template, error) {
	t, ok := s.raids.Catalog().Get(code)
	if !ok || !t.IsRaid() {
		return t, ErrUnknownRaid
	}
	return t, nil
}

// --- Lobby -------------------------------------------------------------------

// Create opens a lobby for the actor's party; the actor (its leader) joins.
func (s *Service) Create(ctx context.Context, actor, account uuid.UUID, code string) (*Run, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	tmpl, err := s.template(code)
	if err != nil {
		return nil, err
	}
	partyID, leader, err := partyOf(ctx, s.pool, actor)
	if err != nil {
		return nil, err
	}
	if leader != actor {
		return nil, ErrNotLeader
	}
	if err := dungeon.InReach(ctx, s.regions, actor, tmpl.RegionCode); err != nil {
		return nil, err
	}
	snap, err := s.snapshot(ctx, actor, account, tmpl)
	if err != nil {
		return nil, err
	}
	runID, err := s.openLobby(ctx, partyID, tmpl, []lobbyMember{{id: actor, snap: snap}})
	if err != nil {
		return nil, err
	}
	return s.load(ctx, runID)
}

// lobbyMember is a character entering a lobby with its frozen snapshot.
type lobbyMember struct {
	id   uuid.UUID
	snap []byte
}

// openLobby opens a party's raid lobby, the first member leading, the rest
// in their order.
func (s *Service) openLobby(ctx context.Context, partyID uuid.UUID, tmpl dungeon.Template, members []lobbyMember) (uuid.UUID, error) {
	runID := uuid.New()
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
			INSERT INTO raid.run (id, party_id, raid_code, leader_character_id, state, duration_seconds)
			VALUES ($1, $2, $3, $4, 'forming', $5)`, runID, partyID, tmpl.Code, members[0].id, tmpl.DurationSecs); err != nil {
			return err
		}
		for i, m := range members {
			if _, err := tx.Exec(ctx, `
				INSERT INTO raid.member (id, run_id, character_id, slot, snapshot) VALUES ($1, $2, $3, $4, $5)`,
				uuid.New(), runID, m.id, i+1, m.snap); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, uniqueErr(err)
	}
	s.notifyParty(ctx, partyID, runID, StateForming)
	return runID, nil
}

// Join adds a party member to the lobby.
func (s *Service) Join(ctx context.Context, actor, account, runID uuid.UUID) (*Run, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	run, err := s.load(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run.State != StateForming {
		return nil, ErrNotForming
	}
	if partyID, _, err := partyOf(ctx, s.pool, actor); err != nil || partyID != run.PartyID {
		return nil, ErrNotInParty
	}
	tmpl, err := s.template(run.RaidCode)
	if err != nil {
		return nil, err
	}
	snap, err := s.snapshot(ctx, actor, account, tmpl)
	if err != nil {
		return nil, err
	}
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := lockForming(ctx, tx, runID); err != nil {
			return err
		}
		var n, last int16
		if err := tx.QueryRow(ctx, `SELECT count(*), COALESCE(max(slot), 0) FROM raid.member WHERE run_id = $1`,
			runID).Scan(&n, &last); err != nil {
			return err
		}
		if n >= tmpl.GroupMax {
			return ErrRaidFull
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO raid.member (id, run_id, character_id, slot, snapshot) VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(), runID, actor, last+1, snap)
		return err
	})
	if err != nil {
		return nil, uniqueErr(err)
	}
	s.notifyParty(ctx, run.PartyID, runID, StateForming)
	return s.load(ctx, runID)
}

// Leave takes a member out of the lobby; the leader leaving calls it off.
func (s *Service) Leave(ctx context.Context, actor, account, runID uuid.UUID) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	var partyID uuid.UUID
	var state string
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var leader uuid.UUID
		if err := tx.QueryRow(ctx, `SELECT party_id, leader_character_id FROM raid.run WHERE id = $1 AND state = 'forming' FOR UPDATE`,
			runID).Scan(&partyID, &leader); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotForming
			}
			return err
		}
		if leader == actor {
			state = StateCancelled
			if _, err := tx.Exec(ctx, `UPDATE raid.run SET state = 'cancelled', updated_at = now() WHERE id = $1`, runID); err != nil {
				return err
			}
			_, err := tx.Exec(ctx, `UPDATE raid.member SET active = false WHERE run_id = $1`, runID)
			return err
		}
		state = StateForming
		tag, err := tx.Exec(ctx, `DELETE FROM raid.member WHERE run_id = $1 AND character_id = $2`, runID, actor)
		if err == nil && tag.RowsAffected() == 0 {
			err = ErrNotMember
		}
		return err
	})
	if err != nil {
		return err
	}
	s.notifyParty(ctx, partyID, runID, state)
	return nil
}

// Start sets the raid off. Members who left the party meanwhile drop out.
func (s *Service) Start(ctx context.Context, actor, account, runID uuid.UUID) (*Run, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	var partyID uuid.UUID
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var leader uuid.UUID
		var code string
		var dur int32
		err := tx.QueryRow(ctx, `
			SELECT party_id, leader_character_id, raid_code, duration_seconds
			  FROM raid.run WHERE id = $1 AND state = 'forming' FOR UPDATE`, runID).Scan(&partyID, &leader, &code, &dur)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotForming
		}
		if err != nil {
			return err
		}
		if leader != actor {
			return ErrNotLeader
		}
		tmpl, err := s.template(code)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			DELETE FROM raid.member m
			 WHERE m.run_id = $1
			   AND NOT EXISTS (SELECT 1 FROM social.party_member pm JOIN social.party p ON p.id = pm.party_id
			                    WHERE pm.party_id = $2 AND pm.character_id = m.character_id
			                      AND pm.left_at IS NULL AND p.disbanded_at IS NULL)`, runID, partyID); err != nil {
			return err
		}
		var n int16
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM raid.member WHERE run_id = $1`, runID).Scan(&n); err != nil {
			return err
		}
		if n < tmpl.GroupMin {
			return ErrNotEnough
		}
		now := s.now().UTC()
		if _, err := tx.Exec(ctx, `DELETE FROM raid.listing WHERE party_id = $1`, partyID); err != nil {
			return err // off the board once under way (DD-299)
		}
		_, err = tx.Exec(ctx, `
			UPDATE raid.run SET state = 'running', started_at = $2, expires_at = $3, seed = $4, updated_at = $2
			 WHERE id = $1`, runID, now, now.Add(time.Duration(dur)*time.Second), seedOf(runID, uuid.New()))
		return err
	})
	if err != nil {
		return nil, err
	}
	s.notifyParty(ctx, partyID, runID, StateRunning)
	return s.load(ctx, runID)
}

// --- Reveal ------------------------------------------------------------------

// RevealResult is a member's collection: the shared fight and their rewards.
type RevealResult struct {
	Result      *Result                  `json:"result"`
	Reward      Reward                   `json:"reward"`
	Members     []Member                 `json:"members"`
	Items       []activity.GrantedItem   `json:"items,omitempty"`
	XPGain      *character.XPGrantResult `json:"xp,omitempty"`
	LedgerTxnID uuid.UUID                `json:"ledger_txn_id"`
	Lockout     bool                     `json:"lockout_written"`
	Replayed    bool                     `json:"replayed"`
}

type stored struct {
	Result *Result `json:"result"`
	Reward Reward  `json:"reward"`
}

// WithBoost applies timed XP boosts (DD-320).
func (s *Service) WithBoost(b boost.Applier) *Service { s.boost = b; return s }

// Reveal collects the actor's rewards once the raid has ended. Idempotent:
// a second call returns the stored result.
func (s *Service) Reveal(ctx context.Context, actor, account, runID uuid.UUID) (*RevealResult, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	run, err := s.load(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run.State != StateRunning {
		return nil, ErrNotRunning
	}
	snaps, err := s.snapshots(ctx, runID)
	if err != nil {
		return nil, err
	}
	pos := -1
	for i, m := range run.Members {
		if m.CharacterID == actor {
			pos = i
		}
	}
	if pos < 0 {
		return nil, ErrNotMember
	}
	me := snaps[pos]
	if me.revealed != nil {
		return replay(run, me.revealed)
	}
	now := s.now().UTC()
	if now.Before(*run.ExpiresAt) {
		return nil, ErrNotReady
	}
	tmpl, err := s.template(run.RaidCode)
	if err != nil {
		return nil, err
	}
	players := make([]combat.Combatant, len(snaps))
	for i := range snaps {
		players[i] = snaps[i].combatant
	}
	mobs := s.mobs.Catalog()
	res := Roll(run.Seed, players, tmpl, mobs, s.rngNew)
	res.Report = res.Report.For(me.combatant) // this member's side (DD-308)
	level := int(me.combatant.Level)
	rw := MemberReward(res, int16(pos+1), level, run.Members[pos].ClassCode, run.Seed, tmpl, mobs, s.rngNew)
	// The whole party runs it together: +3% XP per member (DD-305).
	rw.XP = int64(float64(rw.XP)*partyxp.Factor(len(run.Members)) + 0.5)
	rw.XP = boost.ApplyIf(ctx, s.boost, actor, boost.CatAdventure, *run.StartedAt, *run.ExpiresAt, rw.XP) // DD-320
	memberID := run.Members[pos].ID

	var granted []activity.GrantedItem
	if s.items != nil && rw.Common+rw.Rare+rw.Legendary+rw.Gear.Rolls+rw.Gear.Guaranteed > 0 {
		boss, _ := mobs.Get(tmpl.BossMobCode)
		granted, err = s.items.Grant(ctx, actor, memberID, run.Seed+int64(pos), activity.Drops{
			Kind: activity.KindHunting, MobLevel: activity.DropLevel(int(boss.Level), level),
			Common: rw.Common, Rare: rw.Rare, Legendary: rw.Legendary, Gear: rw.Gear,
		})
		if err != nil {
			return nil, fmt.Errorf("raid reveal item grant: %w", err)
		}
	}
	var txnID uuid.UUID
	if rw.GoldCopper > 0 {
		out, err := s.txn.Write(ctx, economy.Intent{
			Kind: "raid_reveal", RefKind: "raid_member", RefID: memberID,
			IdempotencyKey: "raid-reveal-" + memberID.String(),
			Entries: []economy.Entry{
				{Currency: "gold", AccountKind: economy.AccountSystem, AccountRef: economy.SystemRewardSource, Delta: -rw.GoldCopper},
				{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: rw.GoldCopper},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("raid reveal ledger: %w", err)
		}
		txnID = out.TxnID
	}
	if s.hp != nil {
		if err := s.hp.SetHP(ctx, actor, res.HPFinal[pos], *run.ExpiresAt); err != nil {
			return nil, fmt.Errorf("raid reveal hp: %w", err)
		}
	}
	body, err := json.Marshal(stored{Result: res, Reward: rw})
	if err != nil {
		return nil, err
	}
	first := false
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE raid.member SET active = false, revealed_at = $2, reveal_result = $3,
			       reveal_txn_id = NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid)
			 WHERE id = $1 AND revealed_at IS NULL`, memberID, now, body, txnID)
		if err != nil || tag.RowsAffected() == 0 {
			return err
		}
		first = true
		if !res.Cleared {
			return nil
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO dungeon.weekly_lockout (character_id, dungeon_code, week_key, cleared_at, contract_id)
			VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
			actor, run.RaidCode, dungeon.ISOWeek(*run.StartedAt), now, memberID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if !first { // a parallel call collected first
		snaps, err := s.snapshots(ctx, runID)
		if err != nil {
			return nil, err
		}
		return replay(run, snaps[pos].revealed)
	}

	var xp *character.XPGrantResult
	if s.xp != nil && rw.XP > 0 {
		xp, err = s.xp.GrantXP(ctx, actor, rw.XP, character.XPSource{Kind: "raid", ID: memberID.String(),
			Metadata: map[string]any{"raid": run.RaidCode, "cleared": res.Cleared}})
		if err != nil {
			s.log.WarnContext(ctx, "raid: xp grant failed (rewards already collected)",
				slog.String("member", memberID.String()), slog.String("err", err.Error()))
		}
	}
	events := []progression.Event{{Kind: progression.EventMobKilled, Amount: int64(res.Victories)}}
	if res.BossKilled {
		events[0].Amount--
		events = append(events, progression.Event{Kind: progression.EventBossKilled, Amount: 1},
			progression.Event{Kind: EventRaidCleared, Amount: 1},
			progression.Event{Kind: EventRaidCleared + ":" + run.RaidCode, Amount: 1})
	}
	progression.RecordIf(ctx, s.progress, actor, events...)
	// The egg is minted once, after the collection commits (DD-328).
	if rw.PetEgg != "" && s.eggs != nil {
		inst, err := s.eggs.MintExact(ctx, actor, rw.PetEgg, 1, "raid_egg", memberID)
		if err != nil {
			s.log.ErrorContext(ctx, "raid: pet egg mint failed", slog.String("member", memberID.String()),
				slog.String("egg", rw.PetEgg), slog.String("err", err.Error()))
		} else {
			granted = append(granted, activity.GrantedItem{InstanceID: inst.ID, TemplateCode: inst.TemplateCode,
				Rarity: string(inst.Rarity)})
		}
	}
	newLevel := 0
	if xp != nil && xp.LeveledUp && xp.Character != nil {
		newLevel = int(xp.Character.Level)
	}
	notification.Send(ctx, s.notifier, s.log, notification.Reward{
		AccountID: account, CharacterID: actor, Ref: memberID, Source: "raid", Target: run.RaidCode,
		GoldCopper: rw.GoldCopper, XP: rw.XP, Items: activity.RewardItems(granted), NewLevel: newLevel,
		Extra: map[string]any{"cleared": res.Cleared, "boss_killed": res.BossKilled},
	}.Notice())
	run.Members[pos].Collected = true
	return &RevealResult{Result: res, Reward: rw, Members: run.Members, Items: granted, XPGain: xp,
		LedgerTxnID: txnID, Lockout: res.Cleared}, nil
}

func replay(run *Run, body []byte) (*RevealResult, error) {
	var st stored
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, fmt.Errorf("raid: stored result: %w", err)
	}
	return &RevealResult{Result: st.Result, Reward: st.Reward, Members: run.Members, Replayed: true,
		Lockout: st.Result != nil && st.Result.Cleared}, nil
}

// --- Views -------------------------------------------------------------------

// Current is what the raid screen shows a character.
type Current struct {
	Raid     *Run     `json:"raid"`  // the raid the character is in (until collected)
	Lobby    *Run     `json:"lobby"` // an open lobby of its party it has not joined
	InParty  bool     `json:"in_party"`
	IsLeader bool     `json:"is_leader"`
	Lockouts []string `json:"lockouts"` // raids cleared this week
	WeekKey  string   `json:"week_key"`
	// Raid finding (DD-299): the character's queue place, its party's listing.
	Queue   *QueueStatus `json:"queue,omitempty"`
	Listing *Listing     `json:"listing,omitempty"`
}

func (s *Service) Current(ctx context.Context, actor, account uuid.UUID) (*Current, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	out := &Current{WeekKey: dungeon.ISOWeek(s.now()), Lockouts: []string{}}
	var runID uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT run_id FROM raid.member WHERE character_id = $1 AND active LIMIT 1`, actor).Scan(&runID)
	switch {
	case err == nil:
		if out.Raid, err = s.load(ctx, runID); err != nil {
			return nil, err
		}
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}
	partyID, leader, err := partyOf(ctx, s.pool, actor)
	switch {
	case err == nil:
		out.InParty, out.IsLeader = true, leader == actor
		if out.Listing, err = s.partyListing(ctx, partyID); err != nil {
			return nil, err
		}
		if out.Raid == nil {
			err := s.pool.QueryRow(ctx, `SELECT id FROM raid.run WHERE party_id = $1 AND state = 'forming'`, partyID).Scan(&runID)
			if err == nil {
				if out.Lobby, err = s.load(ctx, runID); err != nil {
					return nil, err
				}
			} else if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
		}
	case !errors.Is(err, ErrNoParty):
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT dungeon_code FROM dungeon.weekly_lockout
		 WHERE character_id = $1 AND week_key = $2`, actor, out.WeekKey)
	if err != nil {
		return nil, err
	}
	codes, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	for _, c := range codes {
		if _, err := s.template(c); err == nil {
			out.Lockouts = append(out.Lockouts, c)
		}
	}
	// A look at the queue is also when a waiting queue forms its party.
	var queuedFor string
	if err := s.pool.QueryRow(ctx, `SELECT raid_code FROM raid.queue WHERE character_id = $1`, actor).
		Scan(&queuedFor); err == nil && !out.InParty {
		if tmpl, err := s.template(queuedFor); err == nil {
			s.match(ctx, tmpl)
		}
		if out.Queue, err = s.queueStatus(ctx, actor); err != nil {
			return nil, err
		}
		if out.Queue == nil { // matched just now: the new party and lobby
			return s.Current(ctx, actor, account)
		}
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return out, nil
}

// --- Helpers -----------------------------------------------------------------

// snapshot checks a member's eligibility and freezes its combat snapshot
// (as a dungeon start does, DD-229).
func (s *Service) snapshot(ctx context.Context, id, account uuid.UUID, tmpl dungeon.Template) ([]byte, error) {
	c, err := s.chars.GetOwnedInternal(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.Level < tmpl.MinLevel {
		return nil, ErrLevelRequirement
	}
	var locked bool
	if err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM dungeon.weekly_lockout WHERE character_id = $1 AND dungeon_code = $2 AND week_key = $3)`,
		id, tmpl.Code, dungeon.ISOWeek(s.now())).Scan(&locked); err != nil {
		return nil, err
	}
	if locked {
		return nil, ErrWeeklyLockout
	}
	if s.bag != nil {
		ok, err := s.bag.HasRoom(ctx, id, "", 1)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrBagFull
		}
	}
	sheet, err := s.chars.Sheet(ctx, id, account)
	if err != nil {
		return nil, err
	}
	if sheet.Resources.HPCurrent < 1 {
		return nil, ErrTooWounded
	}
	snap := combat.FromSheet(sheet)
	snap.RegenPerEncounter = int32(float64(tmpl.DurationSecs/int32(tmpl.TotalEncounters())/6) *
		character.HPRegenPer2s(sheet.Primary.Vitality))
	if s.skills != nil {
		snap.Skills = s.skills(sheet)
	}
	return json.Marshal(snap)
}

func (s *Service) load(ctx context.Context, runID uuid.UUID) (*Run, error) {
	var r Run
	err := s.pool.QueryRow(ctx, `
		SELECT id, party_id, raid_code, leader_character_id, state, duration_seconds, seed,
		       started_at, expires_at, created_at
		  FROM raid.run WHERE id = $1`, runID).Scan(&r.ID, &r.PartyID, &r.RaidCode, &r.LeaderID, &r.State,
		&r.Duration, &r.Seed, &r.StartedAt, &r.ExpiresAt, &r.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t, ok := s.raids.Catalog().Get(r.RaidCode); ok {
		r.GroupMin, r.GroupMax = t.GroupMin, t.GroupMax
	}
	rows, err := s.pool.Query(ctx, `
		SELECT m.id, m.character_id, c.name, c.class_code, c.level, m.slot, m.revealed_at IS NOT NULL
		  FROM raid.member m JOIN character.character c ON c.id = m.character_id
		 WHERE m.run_id = $1 ORDER BY m.slot`, runID)
	if err != nil {
		return nil, err
	}
	r.Members, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (Member, error) {
		var m Member
		err := row.Scan(&m.ID, &m.CharacterID, &m.Name, &m.ClassCode, &m.Level, &m.Slot, &m.Collected)
		return m, err
	})
	if err != nil {
		return nil, err
	}
	for i := range r.Members {
		r.Members[i].Label = Label(int16(i + 1))
	}
	return &r, nil
}

type memberSnap struct {
	combatant combat.Combatant
	revealed  []byte
}

// snapshots are the members' frozen combatants in fight (slot) order.
func (s *Service) snapshots(ctx context.Context, runID uuid.UUID) ([]memberSnap, error) {
	rows, err := s.pool.Query(ctx, `SELECT snapshot, reveal_result FROM raid.member WHERE run_id = $1 ORDER BY slot`, runID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (memberSnap, error) {
		var raw []byte
		var m memberSnap
		if err := row.Scan(&raw, &m.revealed); err != nil {
			return m, err
		}
		return m, json.Unmarshal(raw, &m.combatant)
	})
}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// partyOf is the character's live party and its leader.
func partyOf(ctx context.Context, q querier, id uuid.UUID) (party, leader uuid.UUID, err error) {
	err = q.QueryRow(ctx, `
		SELECT p.id, p.leader_character_id
		  FROM social.party_member m JOIN social.party p ON p.id = m.party_id
		 WHERE m.character_id = $1 AND m.left_at IS NULL AND p.disbanded_at IS NULL`, id).Scan(&party, &leader)
	if errors.Is(err, pgx.ErrNoRows) {
		err = ErrNoParty
	}
	return party, leader, err
}

func lockForming(ctx context.Context, tx pgx.Tx, runID uuid.UUID) error {
	var state string
	err := tx.QueryRow(ctx, `SELECT state FROM raid.run WHERE id = $1 FOR UPDATE`, runID).Scan(&state)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err == nil && state != StateForming {
		err = ErrNotForming
	}
	return err
}

func uniqueErr(err error) error {
	if b := action.FromDB(err); b != err {
		return b // busy with another action (DD-263)
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return err
	}
	if pgErr.ConstraintName == "ux_raid_run_forming" {
		return ErrLobbyExists
	}
	return ErrAlreadyInRaid
}

// notifyParty tells the party's live sessions the raid changed.
func (s *Service) notifyParty(ctx context.Context, partyID, runID uuid.UUID, state string) {
	if s.pub == nil {
		return
	}
	rows, err := s.pool.Query(ctx, `
		SELECT DISTINCT c.account_id FROM social.party_member m JOIN character.character c ON c.id = m.character_id
		 WHERE m.party_id = $1 AND m.left_at IS NULL`, partyID)
	if err != nil {
		s.log.WarnContext(ctx, "raid: notify party", slog.String("err", err.Error()))
		return
	}
	accounts, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return
	}
	for _, a := range accounts {
		s.pub.Publish(a, "raid.changed", map[string]any{"raid_id": runID.String(), "state": state})
	}
}

func seedOf(a, b uuid.UUID) int64 {
	var acc uint64 = 1469598103934665603
	for _, x := range append(a[:], b[:]...) {
		acc ^= uint64(x)
		acc *= 1099511628211
	}
	return int64(acc)
}
