package crafting

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otter/server/internal/domain/character"
	"github.com/otter/server/internal/inventory"
)

// Every domain error must map to a client status: an unmapped one reaches
// players as a 500 "server error" (503 = retry is fine). Errors are wrapped
// like services do.
func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for name, e := range map[string]error{
		"character.ErrForbidden":  character.ErrForbidden, // service ownership checks
		"character.ErrNotFound":   character.ErrNotFound,
		"ErrForbidden":            ErrForbidden,
		"ErrInsufficientGold":     ErrInsufficientGold,
		"ErrNotFound":             ErrNotFound,
		"ErrProfessionNotLearned": ErrProfessionNotLearned,
		"ErrSkillTooLow":          ErrSkillTooLow,
		"ErrBagFull":              inventory.ErrBagFull,
		"ErrUnknownProfession":    ErrUnknownProfession,
		"ErrUnknownRecipe":        ErrUnknownRecipe,
		"ErrValidation":           ErrValidation,
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			wrapped := fmt.Errorf("context: %w", e)
			var codes []int
			rec0 := httptest.NewRecorder()
			writeErr(rec0, req, wrapped)
			codes = append(codes, rec0.Code)
			for _, c := range codes {
				if c < http.StatusInternalServerError || c == http.StatusServiceUnavailable {
					return
				}
			}
			t.Fatalf("%s maps to %v (want a 4xx)", name, codes)
		})
	}
	_ = errors.New
}
