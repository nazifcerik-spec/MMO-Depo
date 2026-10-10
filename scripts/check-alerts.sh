#!/bin/sh
# Validates the alert rules and runs their unit tests (DD-236); with an
# argument also checks an Alertmanager config.
#
#   scripts/check-alerts.sh [deploy/prod/alertmanager.yml]
#
# Uses the pinned Prometheus/Alertmanager images, so only Docker is needed.
set -eu
root=$(cd "$(dirname "$0")/.." && pwd)
prom=prom/prometheus:v2.54.1
am=prom/alertmanager:v0.27.0
export MSYS_NO_PATHCONV=1 # Git Bash: keep container paths as written

docker run --rm --entrypoint promtool -v "$root/deploy:/rules:ro" -w /rules "$prom" check rules alerts.yml
docker run --rm --entrypoint promtool -v "$root/deploy:/rules:ro" -w /rules "$prom" test rules alerts_test.yml
docker run --rm --entrypoint promtool -v "$root/deploy:/rules:ro" -v "$root/deploy/alerts.yml:/etc/prometheus/alerts.yml:ro" \
  "$prom" check config /rules/prod/prometheus.yml

if [ $# -gt 0 ]; then
  docker run --rm --entrypoint amtool -v "$(cd "$(dirname "$1")" && pwd)/$(basename "$1"):/am.yml:ro" "$am" check-config /am.yml
fi
