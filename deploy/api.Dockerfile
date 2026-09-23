# syntax=docker/dockerfile:1.7
#
# Otter API — production-optimized multi-stage build.
# DD-203: distroless nonroot + layer caching + minimum attack surface.
# Build context: project root (docker compose handles this)

# ---- dependencies ----
FROM golang:1.25-alpine AS deps

RUN apk add --no-cache git

WORKDIR /src/server

# Copy only go.mod first
COPY server/go.mod ./

# Tidy and download dependencies (creates go.sum if missing)
RUN go mod tidy && go mod download && go mod verify

# ---- build ----
FROM deps AS build

# Copy full source
COPY server/ ./

ARG COMMIT_SHA=dev

# Build with optimizations: CGO disabled, stripped, minimal deps
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/api ./cmd/api

# ---- runtime ----
FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /out/api /app

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app"]
