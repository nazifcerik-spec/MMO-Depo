package activity_test

import (
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/otter/server/internal/activity"
)

func TestRNG_Deterministic(t *testing.T) {
	r1 := activity.NewRNG(42)
	r2 := activity.NewRNG(42)
	for i := 0; i < 100; i++ {
		require.Equal(t, r1.Uint64(), r2.Uint64(), "same seed must produce same sequence at step %d", i)
	}
}

func TestRNG_DifferentSeedsDiffer(t *testing.T) {
	r1 := activity.NewRNG(42)
	r2 := activity.NewRNG(43)
	same := 0
	for i := 0; i < 100; i++ {
		if r1.Uint64() == r2.Uint64() {
			same++
		}
	}
	require.Less(t, same, 5, "different seeds should almost never collide")
}

func TestSeedFromContract_Stable(t *testing.T) {
	c := uuid.New()
	ch := uuid.New()
	require.Equal(t, activity.SeedFromContract(c, ch), activity.SeedFromContract(c, ch))
}

func TestSeedFromContract_ChangesWithInputs(t *testing.T) {
	c1, c2 := uuid.New(), uuid.New()
	ch := uuid.New()
	require.NotEqual(t, activity.SeedFromContract(c1, ch), activity.SeedFromContract(c2, ch))
}

func TestPoisson_MeanCloseToLambda(t *testing.T) {
	// 10k samples — sample mean should be within a few sigma of λ.
	cases := []float64{0.5, 2, 10, 40, 200}
	for _, lambda := range cases {
		r := activity.NewRNG(int64(lambda * 1000))
		sum := 0
		n := 10_000
		for i := 0; i < n; i++ {
			sum += activity.Poisson(r, lambda)
		}
		mean := float64(sum) / float64(n)
		// 3σ tolerance: σ² = λ, σ_of_mean = √(λ/n)
		tol := 3 * math.Sqrt(lambda/float64(n))
		require.InDelta(t, lambda, mean, tol,
			"Poisson(λ=%v) mean %v not within 3σ tolerance %v", lambda, mean, tol)
	}
}

func TestPoisson_ZeroLambdaReturnsZero(t *testing.T) {
	r := activity.NewRNG(1)
	for i := 0; i < 100; i++ {
		require.Equal(t, 0, activity.Poisson(r, 0))
	}
}

func TestBinomial_MeanCloseToNP(t *testing.T) {
	cases := []struct {
		n int
		p float64
	}{
		{10, 0.5}, {100, 0.01}, {1000, 0.05}, {10000, 0.01},
	}
	for _, c := range cases {
		r := activity.NewRNG(int64(c.n*1000) + int64(c.p*1000))
		sum := 0
		reps := 5000
		for i := 0; i < reps; i++ {
			sum += activity.Binomial(r, c.n, c.p)
		}
		mean := float64(sum) / float64(reps)
		expected := float64(c.n) * c.p
		variance := expected * (1 - c.p)
		tol := 3 * math.Sqrt(variance/float64(reps))
		require.InDelta(t, expected, mean, tol,
			"Binomial(n=%d, p=%v) mean %v not within tolerance", c.n, c.p, mean)
	}
}

func TestBinomial_EdgeCases(t *testing.T) {
	r := activity.NewRNG(1)
	require.Equal(t, 0, activity.Binomial(r, 100, 0))
	require.Equal(t, 100, activity.Binomial(r, 100, 1))
	require.Equal(t, 0, activity.Binomial(r, 0, 0.5))
}

func TestPoisson_Reproducible(t *testing.T) {
	// Same seed → same sequence of Poisson draws.
	r1 := activity.NewRNG(999)
	r2 := activity.NewRNG(999)
	for i := 0; i < 100; i++ {
		require.Equal(t, activity.Poisson(r1, 5), activity.Poisson(r2, 5))
	}
}
