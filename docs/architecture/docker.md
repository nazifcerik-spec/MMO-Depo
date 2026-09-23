# Docker (Containerization)

## 1. Amaç
Geliştirici ortamlarını tekrar edilebilir kılmak; CI/CD'de repeatable build; production'da immutable artifact.

## 2. Kapsam
- **Server binaries**: `cmd/server`, `cmd/worker`, `cmd/migrate` (her biri kendi image).
- **Flutter web**: static build → nginx image (opsiyonel; CDN de gidebilir).
- **Local infra**: `docker-compose.yml` ile Postgres + Redis + backend + workers.
- **Test infra**: `testcontainers-go` ile ephemeral Postgres/Redis.

## 3. Image Prensipleri
- **Multi-stage build**: build layer + minimal runtime.
- **Base**: `gcr.io/distroless/base-debian12` (Go static binary).
- **Non-root user**: `USER 65532:65532`.
- **Read-only rootfs**: Kubernetes deployment'ta enforce.
- **Small**: hedef < 40 MB (Go binary + certs).
- **No shell** in runtime image (distroless).
- Explicit `ENTRYPOINT`; `CMD` opsiyonel argümanlar.
- Signal handling: binary graceful shutdown (SIGTERM).

## 4. Örnek Dockerfile (server)
```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
  -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
  -o /out/server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /out/server /usr/local/bin/server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
```

## 5. Tagging & Versioning
- Tag pattern: `ghcr.io/<org>/mmo-server:<git-sha>`; ayrıca `:main`, `:staging`, `:v1.2.3`.
- Semver release tag'i `v*.*.*`.
- Image signing: cosign (Phase 02 hedef).
- SBOM: `syft` çıktısı CI artifact.

## 6. Local Development
```yaml
# docker-compose.yml (özet)
services:
  postgres:
    image: postgres:16-alpine
    environment: {POSTGRES_PASSWORD: dev, POSTGRES_DB: mmo}
    ports: ["5432:5432"]
    volumes: [pg-data:/var/lib/postgresql/data]
  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
  migrate:
    build: {context: ., dockerfile: deploy/docker/migrate.Dockerfile}
    depends_on: [postgres]
    env_file: [.env.local]
  server:
    build: {context: ., dockerfile: deploy/docker/server.Dockerfile}
    depends_on: [postgres, redis, migrate]
    env_file: [.env.local]
    ports: ["8080:8080"]
  worker:
    build: {context: ., dockerfile: deploy/docker/worker.Dockerfile}
    depends_on: [postgres, redis, migrate]
    env_file: [.env.local]
volumes:
  pg-data:
```

## 7. Build Cache
- CI'da BuildKit cache mount (`--cache-from`, `--cache-to`).
- `go mod download` ayrı layer.
- Reproducible: `-trimpath`, `-buildvcs=false`.

## 8. Security
- Base image regularly refreshed (weekly rebuild policy).
- `trivy` / `grype` scan CI'da; critical vulnerability fail-build.
- Only pinned base images (digest optional for critical builds).
- No secrets in image; env at runtime.

## 9. Runtime Standards
- Health endpoint: `/healthz/live`, `/healthz/ready`.
- Metrics endpoint: `/metrics` (Prometheus).
- Graceful shutdown: SIGTERM → drain (max 30s) → exit.
- Log to stdout/stderr (structured JSON) — collector picks up.

## 10. Resource Guidance (initial)
- **server**: 0.5–2 CPU, 512 MB–2 GB RAM (per node; horizontal scale).
- **worker**: 0.5–1 CPU, 512 MB.
- **postgres (self-hosted v1)**: 2–4 CPU, 8–16 GB RAM.
- **redis**: 1 CPU, 1–2 GB RAM.
- Prod: managed Postgres/Redis öneri (bkz. bekleyen kararlar).

## 11. Anti-Patterns
- `latest` tag production'da.
- Secret'ı Dockerfile'a bake.
- `root` user runtime.
- Full-OS image (ubuntu:latest gibi) 500 MB+ boyut.
- Build içinde `apt-get install` sonrası cache silinmemesi.

## 12. İlişkiler
- [System Architecture](system-architecture.md), [Go Server](go.md), [CI/CD](cicd.md), [Environment](environment.md), [Config](config.md), [Logging](logging.md)

## 13. Bekleyen Kararlar
- Base image: distroless mi Alpine mı (öneri: distroless nonroot).
- Managed vs self-hosted Postgres/Redis v1.
- Orchestrator seçimi (Kubernetes vs Nomad vs simpler) → farklı deployment manifest'i.
