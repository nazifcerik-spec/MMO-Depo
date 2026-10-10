# deploy/

Dev docker-compose, the production images and the single-host production
stack (DD-232).

## Dev

```bash
cd deploy
docker compose up -d          # postgres + redis + migrator + api + worker
docker compose logs -f api
docker compose down            # stop (volumes kept)
docker compose down -v         # also delete volumes
```

## Production (single host)

Requirements: a Linux server with Docker, a DNS record pointing the domain
at it, ports 80/443 open.

```bash
git clone <repo> && cd <repo>/deploy/prod
cp env.example .env            # fill in every value; generate secrets with openssl
../../scripts/deploy.sh        # build (tag = commit), migrate, start, wait healthy
```

- `caddy` obtains TLS certificates automatically and is the only public
  service; postgres and redis stay on the internal network.
- Updates: `git pull && scripts/deploy.sh`. Undo: `scripts/rollback.sh`
  (previous build) or `scripts/rollback.sh <tag>` (see `deployed.log`).
- Backups: the `backup` service writes a dump to `deploy/prod/backups` every
  24 h and keeps 35 days. **Copy that folder off the server** (object
  storage, another machine) — a backup on the same disk does not survive a
  disk loss.
- Restore drill (monthly): `scripts/restore-drill.sh deploy/prod/backups/<file>.dump`
  restores into a throwaway container and checks schema + ledger; it must
  print `DRILL PASSED`.
- Real restore: stop api/worker, `pg_restore --clean --if-exists` the dump
  into the live database, start them again.

## Images

`api.Dockerfile` / `worker.Dockerfile`: distroless nonroot, dependencies
verified against `go.sum`, the content pack baked into the API image, and a
`HEALTHCHECK` that runs the binary itself (`/app healthcheck`). The build
context is the repo root; the root `.dockerignore` keeps it to `server/` and
`content/`.
