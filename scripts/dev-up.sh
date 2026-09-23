#!/usr/bin/env bash
# Otter dev environment start.
set -euo pipefail

cd "$(dirname "$0")/../deploy"
docker compose up -d

echo ""
echo "→ Waiting for API…"
for i in {1..20}; do
  if curl -sf http://localhost:8080/live >/dev/null; then
    echo "✔ API live at http://localhost:8080"
    exit 0
  fi
  sleep 1
done

echo "✗ API did not become ready in 20s"
docker compose logs api | tail -50
exit 1
