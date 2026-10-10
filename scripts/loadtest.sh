#!/bin/sh
# Runs a k6 load scenario against the dev stack (DD-238).
#
#   scripts/loadtest.sh [smoke|t1|login] [duration]     (login: RATE=…/s)
#
# For the run, the API trusts X-Forwarded-For from the compose network so
# every virtual player has its own client IP (rate limits stay on); the
# API is restored to the normal setting afterwards. Needs the dev stack up
# (scripts/dev-up.sh). Prints the API's peak memory at the end; Prometheus
# (--profile observability) shows the rest of the server side.
set -eu
tier=${1:-smoke}
duration=${2:-}
root=$(cd "$(dirname "$0")/.." && pwd)
export MSYS_NO_PATHCONV=1 # Git Bash: keep container paths as written
cd "$root/deploy"

script=mix.js
[ "$tier" = login ] && script=login.js

net=$(docker inspect -f '{{range $k, $v := .NetworkSettings.Networks}}{{$k}}{{end}}' otter_api)
subnet=$(docker network inspect -f '{{(index .IPAM.Config 0).Subnet}}' "$net")
mem=$(mktemp)

restore() {
  [ -n "${sampler:-}" ] && kill "$sampler" 2>/dev/null || true
  docker compose up -d --no-deps api >/dev/null 2>&1 || true
  rm -f "$mem"
}
trap restore EXIT
TRUSTED_PROXIES=$subnet docker compose up -d --no-deps api
until [ "$(docker inspect -f '{{.State.Health.Status}}' otter_api)" = healthy ]; do sleep 2; done

(while true; do docker stats --no-stream --format '{{.MemUsage}}' otter_api >>"$mem" 2>/dev/null; sleep 2; done) &
sampler=$!

docker run --rm --network "$net" -v "$root/perf/k6:/scripts:ro" \
  -e TIER="$tier" ${duration:+-e DURATION="$duration"} ${RATE:+-e RATE="$RATE"} -e BASE=http://api:8080 \
  grafana/k6:0.54.0 run --summary-trend-stats "avg,p(50),p(95),p(99),max" "/scripts/$script" || status=$?

awk '{ v=$1; u=v; gsub(/[0-9.]/,"",u); n=v+0; if (u=="GiB") n*=1024; else if (u=="KiB") n/=1024; if (n>m) m=n }
     END { printf "api peak memory: %.0f MiB\n", m }' "$mem"
exit "${status:-0}"
