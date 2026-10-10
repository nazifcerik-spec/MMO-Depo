#!/bin/sh
# Return api + worker to an earlier build (DD-232, docs/infra/rollback.md).
#
#   scripts/rollback.sh            # the build before the current one
#   scripts/rollback.sh <tag>      # a specific commit tag from deployed.log
#
# Only the application rolls back. Migrations are forward-compatible (expand
# first, contract in a later release), so the previous build runs on the
# newer schema; a migration that must be undone is reverted on purpose with
# migrate ... down 1, or the database restored from a backup.
set -eu

cd "$(dirname "$0")/../deploy/prod"
[ -f deployed.log ] || { echo "no deployed.log: nothing to roll back to" >&2; exit 1; }

tag=${1:-$(awk '{print $2}' deployed.log | tail -2 | head -1)}
docker image inspect "otter/api:$tag" >/dev/null 2>&1 || { echo "image otter/api:$tag not on this host" >&2; exit 1; }

IMAGE_TAG=$tag docker compose up -d --no-build api worker
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) $tag rollback" >>deployed.log
echo "rolled back to $tag"
