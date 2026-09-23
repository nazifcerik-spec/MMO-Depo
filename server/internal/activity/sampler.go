package activity

import (
	"crypto/sha256"
	"encoding/binary"
	"math"

	"github.com/google/uuid"
)

// --- xoshiro256++ ------------------------------------------------------------
//
// Fast, high-quality 64-bit PRNG (Blackman & Vigna). Chosen v1 for:
//   - Deterministic given seed (reveal replay = same result).
//   - Portable across platforms (bit-exact).
//   - Cheap: ~1ns per Uint64() on modern hardware.
//
// See docs/systems/rng.md DD-OPEN-7 (final algorithm choice deferred; xoshiro
// is the v1 default).

// RNG is a deterministic PRNG. Every call to Uint64 advances state.
type RNG struct {
	s [4]uint64
}

// NewRNG seeds an RNG from a single int64 via SplitMix64 initialization.
func NewRNG(seed int64) *RNG {
	var r RNG
	sm := uint64(seed)
	for i := range r.s {
		sm += 0x9E3779B97F4A7C15
		z := sm
		z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
		z = (z ^ (z >> 27)) * 0x94D049BB133111EB
		z = z ^ (z >> 31)
		r.s[i] = z
	}
	return &r
}

// Uint64 returns the next 64-bit random value.
func (r *RNG) Uint64() uint64 {
	result := rotl(r.s[0]+r.s[3], 23) + r.s[0]

	t := r.s[1] << 17

	r.s[2] ^= r.s[0]
	r.s[3] ^= r.s[1]
	r.s[1] ^= r.s[2]
	r.s[0] ^= r.s[3]

	r.s[2] ^= t
	r.s[3] = rotl(r.s[3], 45)

	return result
}

// Float64 returns a uniform value in [0, 1).
func (r *RNG) Float64() float64 {
	// Top 53 bits — safe conversion to float64 without loss.
	return float64(r.Uint64()>>11) / float64(1<<53)
}

// IntN returns a uniform integer in [0, n). Panics if n <= 0.
func (r *RNG) IntN(n int) int {
	if n <= 0 {
		panic("activity: RNG.IntN requires n > 0")
	}
	return int(r.Uint64() % uint64(n))
}

func rotl(x uint64, k int) uint64 { return (x << k) | (x >> (64 - k)) }

// SeedFromContract derives a stable seed from (contract_id, character_id).
// SHA-256 gives a well-distributed 64-bit slice suitable for xoshiro init.
func SeedFromContract(contractID, characterID uuid.UUID) int64 {
	h := sha256.New()
	_, _ = h.Write(contractID[:])
	_, _ = h.Write(characterID[:])
	sum := h.Sum(nil)
	return int64(binary.LittleEndian.Uint64(sum[:8]))
}

// --- Poisson sampler --------------------------------------------------------
//
// Analytical model for "how many discrete events happen in T time given rate λ"
// — mining ore drops, fishing catches, hunt encounter count. Mean = λ, variance
// = λ. O(1) for small-to-moderate λ; O(1) branch for large λ via Gaussian
// approximation (λ ≥ 30) with continuity correction.

// Poisson returns a random draw from Poisson(lambda).
func Poisson(r *RNG, lambda float64) int {
	if lambda <= 0 {
		return 0
	}
	if lambda >= 30 {
		// Normal approximation: N(λ, λ), rounded, clamped ≥ 0.
		mu := lambda
		sigma := math.Sqrt(lambda)
		z := normal(r)
		return int(math.Max(0, math.Round(mu+sigma*z)))
	}
	// Knuth's algorithm — expected iterations ≈ λ.
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for {
		k++
		p *= r.Float64()
		if p <= L {
			return k - 1
		}
	}
}

// Binomial returns a random draw from Binomial(n, p).
//
// Small n → direct sum. Large n with np ≥ 10 and n(1-p) ≥ 10 → normal
// approximation with continuity correction.
func Binomial(r *RNG, n int, p float64) int {
	if n <= 0 || p <= 0 {
		return 0
	}
	if p >= 1 {
		return n
	}
	if n < 25 {
		count := 0
		for i := 0; i < n; i++ {
			if r.Float64() < p {
				count++
			}
		}
		return count
	}
	// Normal approximation.
	np := float64(n) * p
	nq := float64(n) * (1 - p)
	if np >= 10 && nq >= 10 {
		mu := np
		sigma := math.Sqrt(np * (1 - p))
		z := normal(r)
		k := int(math.Round(mu + sigma*z))
		if k < 0 {
			return 0
		}
		if k > n {
			return n
		}
		return k
	}
	// Fallback — direct.
	count := 0
	for i := 0; i < n; i++ {
		if r.Float64() < p {
			count++
		}
	}
	return count
}

// normal returns a draw from N(0, 1) using Box–Muller transform.
func normal(r *RNG) float64 {
	u1 := r.Float64()
	u2 := r.Float64()
	// Avoid log(0)
	if u1 < 1e-300 {
		u1 = 1e-300
	}
	return math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
}
