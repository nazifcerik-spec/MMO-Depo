#!/usr/bin/env sh
# Coverage gate (DD-216). Two modes:
#
#   go test -cover ./... | scripts/coverage_gate.sh go server/coverage-thresholds.txt
#     Per-package ratchet: every package listed in the thresholds file must
#     report coverage >= its minimum. Thresholds sit ~2 points under the last
#     measured value (docs/testing: "fail on coverage drop > 2%"); raise them
#     when coverage improves, never lower them without a DD.
#
#   scripts/coverage_gate.sh lcov client/coverage/lcov.info 80
#     Total line coverage of an lcov file must be >= the given percent.
set -eu

mode=${1:?usage: coverage_gate.sh go <thresholds> | lcov <lcov.info> <min>}

case "$mode" in
go)
  thresholds=${2:?thresholds file required}
  report=$(mktemp)
  trap 'rm -f "$report"' EXIT
  # "ok  github.com/otter/server/internal/x  0.5s  coverage: 57.8% of statements"
  sed -n 's#^ok[[:space:]]*github.com/otter/server/\([^[:space:]]*\).*coverage: \([0-9.]*\)%.*#\1 \2#p' >"$report"
  awk -v report="$report" '
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
  ' "$thresholds"
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
