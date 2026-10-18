//go:build race

package balance

// raceDetector: built with -race. The simulations are single-goroutine and
// deterministic, and run without the detector in the coverage job, so the
// race job skips them (they would exceed the package timeout there).
const raceDetector = true
