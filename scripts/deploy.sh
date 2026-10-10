#!/bin/sh
# Deploy the current git checkout to the single-host production stack
# (DD-232): build images tagged with the commit, migrate, roll api + worker,
# wait until the API is healthy. Run on the server from the repo checkout.
#
#   git pull && scripts/deploy.sh
set -eu

cd "$(dirname "$0")/../deploy/prod"
[ -f .env ] || { echo "deploy/prod/.env missing (copy env.example)" >&2; exit 1; }
[ -z "$(git status --porcelain)" ] || { echo "uncommitted changes: refusing to deploy" >&2; exit 1; }

IMAGE_TAG=$(git rev-parse --short=12 HEAD)
export IMAGE_TAG

docker compose build api worker
docker compose run --rm migrate
docker compose up -d --remove-orphans

i=0
until [ "$(docker inspect --format '{{.State.Health.Status}}' "$(docker compose ps -q api)")" = healthy ]; do
  i=$((i + 1)); [ "$i" -lt 60 ] || { echo "api not healthy after 2 minutes — consider scripts/rollback.sh" >&2; exit 1; }
  sleep 2
done
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) $IMAGE_TAG" >>deployed.log
echo "deployed $IMAGE_TAG"
