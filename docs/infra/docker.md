# Docker (Production Images)

## 1. Amaç
Prod-grade container image: küçük, güvenli, hızlı start, immutable, imzalı.

## 2. Prensipler

- **Distroless nonroot** — `gcr.io/distroless/base-debian12:nonroot` (Go static binary).
- **Multi-stage build** — build stage'de compile, final stage'de binary + minimum runtime.
- **Reproducible** — same commit → same image digest (build cache + pin dep).
- **Immutable tag** — `sha:<git_commit>` primary tag; `v1.4.0` semver secondary.
- **Small size** — Go binary ~15-30 MB; final image < 50 MB target.
- **Signed** — cosign imzası; registry admission controller verify.
- **Scan** — Trivy/Grype CI'de her build; critical CVE fail.
- **No secrets in image** — env var runtime; multi-stage ile build secret leak yok.

## 3. Örnek Dockerfile (Go Service)

```dockerfile
# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG COMMIT_SHA
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/app ./cmd/api

# ---- final stage ----
FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app"]
```

Not: `-s -w` symbol strip; binary küçük. CGO_ENABLED=0 static; distroless glibc-free çalışır.

## 4. Registry

- **Production**: GHCR (GitHub Container Registry) veya ECR/GCR — ops karar.
- **Naming**: `<registry>/<org>/<service>:<tag>`.
- **Retention**: staging tags 30 gün; prod tags 1 yıl (rollback için).
- **Access**: OIDC federation CI → registry (no long-lived credentials).

## 5. Vulnerability Scanning

CI adımları:
1. **Build image**.
2. **Trivy scan** (image + fs + config).
3. **Critical**: fail; **High**: warn + issue create.
4. **Cosign sign**.
5. **Push registry**.

Runtime (Phase 02):
- Admission controller (Kyverno) → sadece signed + non-vuln image deploy.
- Nightly scan tüm çalışan image'lar; drift alert.

## 6. Image Boyutları Hedef

| Service | Base | Final Size |
|---|---|---|
| API (Go) | distroless base | < 50 MB |
| Worker (Go) | distroless base | < 50 MB |
| Frontend admin (React) | nginx alpine | < 30 MB |
| Migration tool (Go) | distroless base | < 30 MB |

## 7. Runtime Config

- Config env var + secret mount.
- **Health probe**:
  - `/live` (200 always if process up).
  - `/ready` (200 if DB + Redis reachable + warmup complete).
- **Graceful shutdown**: SIGTERM handler + drain WS + wait 30s max.
- **Resource limits**: CPU + memory k8s manifest (bkz. [performance/capacity-planning.md](../performance/capacity-planning.md)).

## 8. Layer Caching

Multi-stage cache:
- `go.mod` + `go.sum` copy önce → dep cache.
- Source copy sonra → incremental build fast.
- CI cache mount `.cache/go-build` + `.cache/go-mod`.

Build time hedef: < 2 dk full; < 30 sn incremental.

## 9. Local Development

`docker-compose.yml` staging-like:
- App + Postgres + Redis + prometheus + grafana + jaeger.
- Volume mount source (hot reload with `air` Go tool).
- Env `.env.dev` non-secret defaults.

## 10. Multi-Arch (Phase 02)

- arm64 support (M1/Mac dev + AWS Graviton).
- `docker buildx` cross-compile.
- Registry manifest list.

## 11. Metrikler
- `docker_image_build_duration_seconds{service}`
- `docker_image_scan_cve_total{severity}`
- `docker_image_size_bytes{service}`

## 12. Cross-Reference
- [architecture/docker.md](../architecture/docker.md) — v1 temel karar
- [ci-cd-prod.md](ci-cd-prod.md), [security/](../security/)

## 13. Bekleyen Kararlar
- Registry provider (GHCR/ECR/GCR) ops karar.
- Cosign key management (KMS integration) Phase 02.
- Kyverno admission controller Phase 02.
- Multi-arch arm64 rollout Phase 02.
