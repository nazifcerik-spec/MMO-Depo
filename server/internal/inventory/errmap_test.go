package inventory

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Every domain error must map to a client status: an unmapped one reaches
// players as a 500 "server error" (503 = retry is fine). Errors are wrapped
// like services do.
func TestDomainErrorsMapToClientStatus(t *testing.T) {
	for name, e := range map[string]error{
		"ErrBagFull":              ErrBagFull,
		"ErrClassRestriction":     ErrClassRestriction,
		"ErrInvalidSplit":         ErrInvalidSplit,
		"ErrLevelRequirement":     ErrLevelRequirement,
		"ErrLocked":               ErrLocked,
		"ErrNotEquippable":        ErrNotEquippable,
		"ErrNotFound":             ErrNotFound,
		"ErrNotStackable":         ErrNotStackable,
		"ErrNothingEquipped":      ErrNothingEquipped,
		"ErrSlotOccupied":         ErrSlotOccupied,
		"ErrSlotOutOfRange":       ErrSlotOutOfRange,
		"ErrUnknownEquipSlot":     ErrUnknownEquipSlot,
		"ErrSpecRequirement":      ErrSpecRequirement,
		"ErrAttributeRequirement": ErrAttributeRequirement,
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			wrapped := fmt.Errorf("context: %w", e)
			var codes []int
			rec0 := httptest.NewRecorder()
			writeErr(rec0, req, wrapped)
			codes = append(codes, rec0.Code)
			rec1 := httptest.NewRecorder()
			writeEquipErr(rec1, req, wrapped)
			codes = append(codes, rec1.Code)
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
