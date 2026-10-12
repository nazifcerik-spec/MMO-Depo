#!/bin/sh
# Restore drill (DD-232, docs/infra/backup.md): restores a dump into a
# throwaway Postgres container and checks it is a usable game database.
# Run monthly and after any backup change; it never touches the live DB.
#
#   restore-drill.sh <dump file>
#
# Checks: the restore succeeds, the schema is at the newest migration in
# server/migrations, accounts exist, and every ledger transaction balances
# per currency (double-entry, DD-135). Exit 0 = drill passed.
set -eu

dump=${1:?usage: restore-drill.sh <dump file>}
[ -f "$dump" ] || { echo "no such file: $dump" >&2; exit 2; }
root=$(cd "$(dirname "$0")/.." && pwd)
expected=$(ls "$root/server/migrations" | sed -n 's/^\([0-9]*\)_.*\.up\.sql$/\1/p' | sort -n | tail -1 | sed 's/^0*//')

name="otter-restore-drill-$$"
cleanup() { docker rm -f "$name" >/dev/null 2>&1 || true; }
trap cleanup EXIT

start=$(date +%s)
docker run -d --name "$name" -e POSTGRES_PASSWORD=drill -e POSTGRES_DB=otter postgres:16.4-alpine >/dev/null
i=0
# TCP, not the socket: the image's first-boot init server listens on the socket only.
until docker exec "$name" pg_isready -h 127.0.0.1 -U postgres -d otter >/dev/null 2>&1; do
  i=$((i + 1)); [ "$i" -lt 60 ] || { echo "postgres did not start" >&2; exit 1; }
  sleep 1
done
docker cp "$dump" "$name:/tmp/drill.dump"
docker exec "$name" pg_restore --no-owner --no-privileges --exit-on-error -h 127.0.0.1 -U postgres -d otter /tmp/drill.dump

q() { docker exec "$name" psql -h 127.0.0.1 -U postgres -d otter -Atc "$1"; }
version=$(q "SELECT version FROM schema_migrations")
dirty=$(q "SELECT dirty FROM schema_migrations")
accounts=$(q "SELECT count(*) FROM identity.account")
characters=$(q "SELECT count(*) FROM character.character")
unbalanced=$(q "SELECT count(*) FROM (SELECT txn_id, currency_code FROM economy.ledger_entry
                                      GROUP BY 1, 2 HAVING sum(delta) <> 0) t")
took=$(( $(date +%s) - start ))

echo "restored:   $dump in ${took}s"
echo "schema:     version $version (expected $expected), dirty=$dirty"
echo "accounts:   $accounts, characters: $characters"
echo "ledger:     $unbalanced unbalanced transactions"

fail=0
[ "$version" = "$expected" ] || { echo "FAIL: schema version"; fail=1; }
[ "$dirty" = "f" ] || { echo "FAIL: dirty migration state"; fail=1; }
[ "$accounts" -gt 0 ] || { echo "FAIL: no accounts"; fail=1; }
[ "$unbalanced" -eq 0 ] || { echo "FAIL: ledger does not balance"; fail=1; }
[ "$fail" -eq 0 ] && echo "DRILL PASSED" || echo "DRILL FAILED"
exit "$fail"
