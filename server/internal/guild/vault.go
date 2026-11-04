package guild

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/otter/server/internal/economy"
	"github.com/otter/server/internal/http/apierr"
	"github.com/otter/server/internal/inventory"
)

// Daily withdrawal limits per rank (DD-358, docs/guild/permissions.md §5,
// scaled to this economy). -1 is no limit, 0 no access.
var (
	goldPerDay  = map[int16]int64{RankLeader: -1, RankOfficer: 10_000_000, RankVeteran: 2_000_000, RankMember: 500_000, RankInitiate: 0}
	itemsPerDay = map[int16]int{RankLeader: -1, RankOfficer: -1, RankVeteran: 30, RankMember: 10, RankInitiate: 0}
)

// CanUseBankItems: initiates only give gold; they neither see nor move items.
func CanUseBankItems(rank int16) bool { return rank < RankInitiate }

var (
	ErrVaultNotEmpty = errors.New("guild: empty the treasury and the bank first")
	ErrVaultClosing  = errors.New("guild: the guild is disbanding, nothing goes in")
	ErrVaultLimit    = errors.New("guild: daily withdrawal limit reached")
	ErrVaultAmount   = errors.New("guild: amount must be positive")
)

// ItemVault moves items in and out of the guild's slots.
type ItemVault interface {
	GuildBank(ctx context.Context, guildID uuid.UUID) (*inventory.BankView, error)
	GuildDeposit(ctx context.Context, characterID, guildID, itemID uuid.UUID, qty int32) (inventory.Moved, error)
	GuildWithdraw(ctx context.Context, characterID, guildID, itemID uuid.UUID, qty int32) (inventory.Moved, error)
}

// WithVault enables the guild bank.
func (s *Service) WithVault(v ItemVault) *Service { s.vault = v; return s }

// VaultLog is one movement.
type VaultLog struct {
	CharacterName string    `json:"character_name"`
	Action        string    `json:"action"`
	TemplateCode  *string   `json:"template_code,omitempty"`
	Rarity        *string   `json:"rarity,omitempty"`
	Quantity      *int32    `json:"quantity,omitempty"`
	GoldCopper    *int64    `json:"gold_copper,omitempty"`
	At            time.Time `json:"at"`
}

// VaultView is GET /guild/bank.
type VaultView struct {
	TreasuryCopper int64               `json:"treasury_copper"`
	Bank           *inventory.BankView `json:"bank"`
	Log            []VaultLog          `json:"log"`
	Rank           int16               `json:"rank_index"`
	CanUseItems    bool                `json:"can_use_items"`
	GoldLimit      int64               `json:"gold_limit"` // -1 = no limit
	GoldToday      int64               `json:"gold_today"`
	ItemsLimit     int                 `json:"items_limit"`
	ItemsToday     int                 `json:"items_today"`
}

func (s *Service) today() time.Time { return s.now().UTC().Truncate(24 * time.Hour) }

// Vault shows the treasury, the bank and the latest movements.
func (s *Service) Vault(ctx context.Context, actor, account uuid.UUID) (*VaultView, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return nil, err
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return nil, err
	}
	v := &VaultView{Rank: rank, CanUseItems: CanUseBankItems(rank), GoldLimit: goldPerDay[rank], ItemsLimit: itemsPerDay[rank],
		Log: []VaultLog{}, Bank: &inventory.BankView{Items: []inventory.BagItem{}}}
	if v.TreasuryCopper, err = s.txn.Balance(ctx, economy.AccountGuild, guildID, "gold"); err != nil {
		return nil, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT COALESCE(sum(copper), 0), COALESCE(sum(items), 0) FROM social.guild_withdrawal
		WHERE guild_id = $1 AND character_id = $2 AND day = $3`, guildID, actor, s.today()).Scan(&v.GoldToday, &v.ItemsToday); err != nil {
		return nil, err
	}
	if v.CanUseItems && s.vault != nil {
		if v.Bank, err = s.vault.GuildBank(ctx, guildID); err != nil {
			return nil, err
		}
	}
	rows, err := s.pool.Query(ctx, `SELECT character_name, action, template_code, rarity, quantity, gold_copper, at
		FROM social.guild_vault_log WHERE guild_id = $1 ORDER BY at DESC, id DESC LIMIT 50`, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l VaultLog
		if err := rows.Scan(&l.CharacterName, &l.Action, &l.TemplateCode, &l.Rarity, &l.Quantity, &l.GoldCopper, &l.At); err != nil {
			return nil, err
		}
		v.Log = append(v.Log, l)
	}
	return v, rows.Err()
}

func (s *Service) logVault(ctx context.Context, guildID, actor uuid.UUID, action string, m *inventory.Moved, copper *int64) {
	var tmpl, rarity *string
	var qty *int32
	if m != nil {
		tmpl, rarity, qty = &m.Template, &m.Rarity, &m.Qty
	}
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO social.guild_vault_log (guild_id, character_id, character_name, action, template_code, rarity, quantity, gold_copper)
		SELECT $1, $2, name, $3, $4, $5, $6, $7 FROM character.character WHERE id = $2`,
		guildID, actor, action, tmpl, rarity, qty, copper); err != nil {
		s.log.ErrorContext(ctx, "guild: vault log failed", "guild_id", guildID.String(), "err", err.Error())
	}
}

// open is the membership check every deposit shares: no deposits while the
// guild is disbanding, so it can only get emptier.
func (s *Service) openForDeposit(ctx context.Context, actor, account uuid.UUID) (uuid.UUID, int16, error) {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return uuid.Nil, 0, err
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return uuid.Nil, 0, err
	}
	var closing bool
	if err := s.pool.QueryRow(ctx, `SELECT disband_scheduled_at IS NOT NULL FROM social.guild WHERE id = $1`, guildID).Scan(&closing); err != nil {
		return uuid.Nil, 0, err
	}
	if closing {
		return uuid.Nil, 0, ErrVaultClosing
	}
	return guildID, rank, nil
}

// DepositGold moves gold from the character to the treasury. Every rank may give.
func (s *Service) DepositGold(ctx context.Context, actor, account uuid.UUID, copper int64, idem string) error {
	if copper <= 0 {
		return ErrVaultAmount
	}
	if idem == "" {
		return ErrIdempotencyMissing
	}
	guildID, _, err := s.openForDeposit(ctx, actor, account)
	if err != nil {
		return err
	}
	res, err := s.txn.Write(ctx, economy.Intent{
		Kind: "guild_treasury_deposit", RefKind: "guild", RefID: guildID, IdempotencyKey: "guild-treasury-in-" + idem,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: -copper},
			{Currency: "gold", AccountKind: economy.AccountGuild, AccountRef: guildID, Delta: copper},
		},
	})
	if err != nil {
		return err
	}
	if !res.Replayed {
		s.logVault(ctx, guildID, actor, "gold_in", nil, &copper)
	}
	return nil
}

// reserve counts a withdrawal against today's limit first, so two at once
// cannot both pass it; a failed move gives the allowance back.
func (s *Service) reserve(ctx context.Context, guildID, actor uuid.UUID, rank int16, copper int64, items int) error {
	if (copper > 0 && goldPerDay[rank] == 0) || (items > 0 && itemsPerDay[rank] == 0) {
		return ErrForbidden
	}
	var gotCopper int64
	var gotItems int
	err := s.pool.QueryRow(ctx, `
		INSERT INTO social.guild_withdrawal (guild_id, character_id, day, copper, items) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (guild_id, character_id, day) DO UPDATE
		   SET copper = social.guild_withdrawal.copper + $4, items = social.guild_withdrawal.items + $5
		RETURNING copper, items`, guildID, actor, s.today(), copper, items).Scan(&gotCopper, &gotItems)
	if err != nil {
		return err
	}
	if (goldPerDay[rank] >= 0 && gotCopper > goldPerDay[rank]) || (itemsPerDay[rank] >= 0 && gotItems > itemsPerDay[rank]) {
		s.release(ctx, guildID, actor, copper, items)
		return ErrVaultLimit
	}
	return nil
}

func (s *Service) release(ctx context.Context, guildID, actor uuid.UUID, copper int64, items int) {
	if _, err := s.pool.Exec(ctx, `UPDATE social.guild_withdrawal SET copper = copper - $4, items = items - $5
		WHERE guild_id = $1 AND character_id = $2 AND day = $3`, guildID, actor, s.today(), copper, items); err != nil {
		s.log.ErrorContext(ctx, "guild: withdrawal release failed", "guild_id", guildID.String(), "err", err.Error())
	}
}

// WithdrawGold moves gold from the treasury to the character, within the rank's
// daily limit.
func (s *Service) WithdrawGold(ctx context.Context, actor, account uuid.UUID, copper int64, idem string) error {
	if copper <= 0 {
		return ErrVaultAmount
	}
	if idem == "" {
		return ErrIdempotencyMissing
	}
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return err
	}
	if err := s.reserve(ctx, guildID, actor, rank, copper, 0); err != nil {
		return err
	}
	res, err := s.txn.Write(ctx, economy.Intent{
		Kind: "guild_treasury_withdraw", RefKind: "guild", RefID: guildID, IdempotencyKey: "guild-treasury-out-" + idem,
		Entries: []economy.Entry{
			{Currency: "gold", AccountKind: economy.AccountGuild, AccountRef: guildID, Delta: -copper},
			{Currency: "gold", AccountKind: economy.AccountCharacter, AccountRef: actor, Delta: copper},
		},
	})
	if err != nil {
		s.release(ctx, guildID, actor, copper, 0)
		return err
	}
	if res.Replayed {
		s.release(ctx, guildID, actor, copper, 0)
		return nil
	}
	s.logVault(ctx, guildID, actor, "gold_out", nil, &copper)
	return nil
}

// DepositItem moves a bag item into the guild's slots.
func (s *Service) DepositItem(ctx context.Context, actor, account, item uuid.UUID, qty int32) error {
	guildID, rank, err := s.openForDeposit(ctx, actor, account)
	if err != nil {
		return err
	}
	if !CanUseBankItems(rank) {
		return ErrForbidden
	}
	m, err := s.vault.GuildDeposit(ctx, actor, guildID, item, qty)
	if err != nil {
		return err
	}
	s.logVault(ctx, guildID, actor, "item_in", &m, nil)
	return nil
}

// WithdrawItem moves a banked item to the character's bag, within the rank's
// daily limit (a stack counts as one).
func (s *Service) WithdrawItem(ctx context.Context, actor, account, item uuid.UUID, qty int32) error {
	if err := s.chars.Owns(ctx, actor, account); err != nil {
		return err
	}
	guildID, rank, err := membership(ctx, s.pool, actor)
	if err != nil {
		return err
	}
	if !CanUseBankItems(rank) {
		return ErrForbidden
	}
	if err := s.reserve(ctx, guildID, actor, rank, 0, 1); err != nil {
		return err
	}
	m, err := s.vault.GuildWithdraw(ctx, actor, guildID, item, qty)
	if err != nil {
		s.release(ctx, guildID, actor, 0, 1)
		return err
	}
	s.logVault(ctx, guildID, actor, "item_out", &m, nil)
	return nil
}

// vaultEmpty: a guild only disbands empty, so nothing is left behind.
func (s *Service) vaultEmpty(ctx context.Context, q querier, guildID uuid.UUID) (bool, error) {
	var empty bool
	err := q.QueryRow(ctx, `SELECT
		NOT EXISTS (SELECT 1 FROM economy.guild_wallet WHERE guild_id = $1 AND amount > 0) AND
		NOT EXISTS (SELECT 1 FROM inventory.item_instance WHERE guild_bank_id = $1)`, guildID).Scan(&empty)
	return empty, err
}

// --- HTTP ---------------------------------------------------------------------

type vaultBody struct {
	Copper int64  `json:"copper"`
	ItemID string `json:"item_id"`
	Qty    int32  `json:"quantity"`
}

func (h *Handler) Vault(w http.ResponseWriter, r *http.Request) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	v, err := h.svc.Vault(r.Context(), actor, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, v)
}

// vaultAction runs a treasury or bank move and answers with the new vault.
func (h *Handler) vaultAction(w http.ResponseWriter, r *http.Request, do func(ctx context.Context, actor, acct uuid.UUID, b vaultBody) error) {
	actor, acct, ok := scope(w, r)
	if !ok {
		return
	}
	var b vaultBody
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&b); err != nil {
		apierr.Write(w, r, http.StatusBadRequest, "E_INVALID_JSON", err.Error())
		return
	}
	if err := do(r.Context(), actor, acct, b); err != nil {
		writeErr(w, r, err)
		return
	}
	v, err := h.svc.Vault(r.Context(), actor, acct)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	apierr.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) DepositGold(w http.ResponseWriter, r *http.Request) {
	idem := r.Header.Get("Idempotency-Key")
	h.vaultAction(w, r, func(ctx context.Context, actor, acct uuid.UUID, b vaultBody) error {
		return h.svc.DepositGold(ctx, actor, acct, b.Copper, idem)
	})
}

func (h *Handler) WithdrawGold(w http.ResponseWriter, r *http.Request) {
	idem := r.Header.Get("Idempotency-Key")
	h.vaultAction(w, r, func(ctx context.Context, actor, acct uuid.UUID, b vaultBody) error {
		return h.svc.WithdrawGold(ctx, actor, acct, b.Copper, idem)
	})
}

func (h *Handler) DepositItem(w http.ResponseWriter, r *http.Request) {
	h.vaultAction(w, r, func(ctx context.Context, actor, acct uuid.UUID, b vaultBody) error {
		id, err := uuid.Parse(b.ItemID)
		if err != nil {
			return inventory.ErrNotFound
		}
		return h.svc.DepositItem(ctx, actor, acct, id, b.Qty)
	})
}

func (h *Handler) WithdrawItem(w http.ResponseWriter, r *http.Request) {
	h.vaultAction(w, r, func(ctx context.Context, actor, acct uuid.UUID, b vaultBody) error {
		id, err := uuid.Parse(b.ItemID)
		if err != nil {
			return inventory.ErrNotFound
		}
		return h.svc.WithdrawItem(ctx, actor, acct, id, b.Qty)
	})
}
