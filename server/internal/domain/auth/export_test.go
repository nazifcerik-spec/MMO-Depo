package auth

import (
	"testing"
	"time"
)

// DummyHash exposes the unknown-account timing dummy.
func (s *Service) DummyHash() string { return s.dummyHash }

// SaturateKDF takes every key-derivation slot and shortens the wait, so the
// next hash or verify returns ErrBusy. Restored on cleanup.
func SaturateKDF(t testing.TB) {
	slots, wait := kdfSlots, kdfWait
	full := make(chan struct{}, 1)
	full <- struct{}{}
	kdfSlots, kdfWait = full, 10*time.Millisecond
	t.Cleanup(func() { kdfSlots, kdfWait = slots, wait })
}
