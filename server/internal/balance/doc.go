// Package balance holds the balance simulation (tests only): every class at
// every level fights the content's own mobs with the gear the content gives
// that level, through the production sheet and combat code. The tests pin the
// balance targets (DD-229); `go test ./internal/balance -v` prints the matrix.
package balance
