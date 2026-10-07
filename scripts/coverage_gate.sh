#!/usr/bin/env sh
# Coverage gate (DD-216, DD-223). Three modes:
#
#   go test -tags integration -coverpkg=./... -coverprofile=cover.out ./...
#   scripts/coverage_gate.sh profile server/coverage-thresholds.txt cover.out
#     Per-package ratchet on the merged profile of unit AND integration tests
#     (a block counts as covered if any test binary ran it). This is what CI
#     gates on: most domain code is SQL that only the integration suite runs.
#
#   go test -cover ./... | scripts/coverage_gate.sh go server/coverage-thresholds.txt
#     Same ratchet on per-package unit-test output only (quick local check;
#     reads lower than the profile mode).
#
#   scripts/coverage_gate.sh lcov client/coverage/lcov.info 80
#     Total line coverage of an lcov file must be >= the given percent.
#
# Thresholds sit ~2 points under the last measured value (docs/testing: "fail
# on coverage drop > 2%"); raise them when coverage improves, never lower them
# without a DD.
set -eu

mode=${1:?usage: coverage_gate.sh profile <thresholds> <cover.out> | go <thresholds> | lcov <lcov.info> <min>}

# compare <thresholds> <report>: report lines are "<pkg> <percent>".
compare() {
  awk -v report="$2" '
    BEGIN { while ((getline line < report) > 0) { split(line, a, " "); cov[a[1]] = a[2] } }
    /^[[:space:]]*(#|$)/ { next }
    {
      pkg = $1; min = $2 + 0
      if (!(pkg in cov)) { printf "FAIL %-32s no coverage reported (min %.1f%%)\n", pkg, min; bad = 1; next }
      got = cov[pkg] + 0
      status = (got + 0.0001 < min) ? "FAIL" : "ok  "
      if (status == "FAIL") bad = 1
      printf "%s %-32s %5.1f%% (min %.1f%%)\n", status, pkg, got, min
    }
    END { exit bad }
  ' "$1"
}

report=$(mktemp)
trap 'rm -f "$report"' EXIT

case "$mode" in
profile)
  thresholds=${2:?thresholds file required}
  profile=${3:?coverprofile required}
  # "github.com/otter/server/internal/x/f.go:12.3,15.4 2 1" → per package
  # statements and covered statements; a block listed by several test
  # binaries is counted once, covered if any count > 0.
  awk '
    NR == 1 && /^mode:/ { next }
    {
      key = $1; n = $2 + 0; c = $3 + 0
      if (!(key in stmts)) { stmts[key] = n; order[++k] = key }
      if (c > 0) hit[key] = 1
    }
    END {
      for (i = 1; i <= k; i++) {
        key = order[i]
        file = key; sub(/:.*/, "", file)
        pkg = file; sub(/\/[^\/]*$/, "", pkg); sub(/^github\.com\/otter\/server\//, "", pkg)
        total[pkg] += stmts[key]
        if (key in hit) covered[pkg] += stmts[key]
      }
      for (p in total) printf "%s %.1f\n", p, (total[p] ? 100 * covered[p] / total[p] : 0)
    }
  ' "$profile" >"$report"
  compare "$thresholds" "$report"
  ;;
go)
  thresholds=${2:?thresholds file required}
  # "ok  github.com/otter/server/internal/x  0.5s  coverage: 57.8% of statements"
  sed -n 's#^ok[[:space:]]*github.com/otter/server/\([^[:space:]]*\).*coverage: \([0-9.]*\)%.*#\1 \2#p' >"$report"
  compare "$thresholds" "$report"
  ;;
lcov)
  file=${2:?lcov file required}
  min=${3:?minimum percent required}
  awk -F: -v min="$min" '
    /^LF:/ { lf += $2 } /^LH:/ { lh += $2 }
    END {
      pct = lf ? 100 * lh / lf : 0
      printf "%s total line coverage %.1f%% (min %.1f%%)\n", (pct + 0.0001 < min ? "FAIL" : "ok  "), pct, min
      exit (pct + 0.0001 < min)
    }
  ' "$file"
  ;;
*)
  echo "unknown mode: $mode" >&2
  exit 2
  ;;
esac
