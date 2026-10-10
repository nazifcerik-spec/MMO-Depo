package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Account rights (DD-231, docs/auth/account.md §4.4-4.6).
const (
	// DefaultTermsVersion is the date of the published terms and privacy policy.
	DefaultTermsVersion = "2026-10-10"
	// DeletionGrace: signing in within this window cancels a deletion request.
	DeletionGrace = 30 * 24 * time.Hour
	// ExportLinkTTL bounds a one-time data export link.
	ExportLinkTTL = 15 * time.Minute
	// exportLinksPerHour caps export links per account.
	exportLinksPerHour = 5
)

// ErrExportLimit: too many export links in the last hour.
var ErrExportLimit = errors.New("auth: export link limit reached")

// RequestDeletion schedules the account for deletion after DeletionGrace and
// signs every device out. A repeated request keeps the original date.
// Returns when the data will be erased.
func (s *Service) RequestDeletion(ctx context.Context, accountID uuid.UUID, password string) (time.Time, error) {
	a, err := s.repo.AccountByID(ctx, accountID)
	if err != nil {
		return time.Time{}, err
	}
	if err := VerifyPassword(password, a.PasswordHash); err != nil {
		return time.Time{}, ErrInvalidCredentials
	}
	now := s.cfg.Now().UTC()
	at := now
	if a.DeletionRequestedAt != nil {
		at = *a.DeletionRequestedAt
	} else {
		if err := s.repo.SetDeletionRequested(ctx, accountID, &now); err != nil {
			return time.Time{}, err
		}
		if err := s.repo.QueueEmail(ctx, deletionEmail(a, now.Add(DeletionGrace))); err != nil {
			s.log.WarnContext(ctx, "auth: deletion mail not queued", "err", err.Error())
		}
	}
	if err := s.repo.RevokeAccountSessions(ctx, accountID, uuid.Nil, "account_deletion", now); err != nil {
		return time.Time{}, err
	}
	return at.Add(DeletionGrace), nil
}

func deletionEmail(a *Account, purgeAt time.Time) OutboundEmail {
	date := purgeAt.Format("2006-01-02")
	return OutboundEmail{
		To:      a.Email,
		Subject: "New World MMO — Hesap silme / Account deletion",
		Body: fmt.Sprintf("Merhaba %s,\n\nHesabınız ve oyun verileriniz %s tarihinde kalıcı olarak silinecek. "+
			"Vazgeçerseniz bu tarihten önce giriş yapmanız yeterli.\n\n"+
			"Hi %s,\n\nYour account and game data will be permanently deleted on %s. "+
			"To cancel, just sign in before then.\n", a.Username, date, a.Username, date),
	}
}

// IssueExportLink returns a one-time link to download the account's data.
func (s *Service) IssueExportLink(ctx context.Context, accountID uuid.UUID) (string, time.Time, error) {
	now := s.cfg.Now().UTC()
	n, err := s.repo.CountAccountTokens(ctx, accountID, TokenDataExport, now.Add(-time.Hour))
	if err != nil {
		return "", time.Time{}, err
	}
	if n >= exportLinksPerHour {
		return "", time.Time{}, ErrExportLimit
	}
	plain, hash, err := newRecoveryToken()
	if err != nil {
		return "", time.Time{}, err
	}
	exp := now.Add(ExportLinkTTL)
	if err := s.repo.IssueAccountToken(ctx, AccountToken{
		ID: uuid.New(), AccountID: accountID, Kind: TokenDataExport, TokenHash: hash, ExpiresAt: exp, CreatedAt: now,
	}, OutboundEmail{}); err != nil {
		return "", time.Time{}, err
	}
	return strings.TrimRight(s.cfg.PublicBaseURL, "/") + "/v1/me/export?token=" + plain, exp, nil
}

// ConsumeExportToken redeems an export link once; ErrRecoveryTokenInvalid
// when unknown, used or expired.
func (s *Service) ConsumeExportToken(ctx context.Context, token string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, ErrRecoveryTokenInvalid
	}
	return s.repo.ConsumeAccountToken(ctx, TokenDataExport, hashRecoveryToken(token), s.cfg.Now().UTC())
}
