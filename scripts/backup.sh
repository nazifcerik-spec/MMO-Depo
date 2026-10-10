#!/bin/sh
# Logical backup of the game database (DD-232).
#
#   backup.sh <out_dir>
#
# Connection from libpq env (PGHOST, PGUSER, PGPASSWORD, PGDATABASE). Writes
# otter-<UTC timestamp>.dump (pg_dump custom format), checks the archive is
# readable, deletes dumps older than KEEP_DAYS (default 35, the privacy
# policy's promise) and records LAST_SUCCESS for monitoring.
set -eu

out=${1:?usage: backup.sh <out_dir>}
keep=${KEEP_DAYS:-35}
mkdir -p "$out"

file="$out/otter-$(date -u +%Y%m%dT%H%M%SZ).dump"
pg_dump --format=custom --no-owner --no-privileges --file="$file.part"
pg_restore --list "$file.part" >/dev/null # unreadable archive = failed backup
mv "$file.part" "$file"

find "$out" -name 'otter-*.dump' -type f -mtime +"$keep" -delete
date -u +%Y-%m-%dT%H:%M:%SZ >"$out/LAST_SUCCESS"
echo "backup ok: $file ($(du -h "$file" | cut -f1))"
