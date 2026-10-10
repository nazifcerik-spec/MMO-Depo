# syntax=docker/dockerfile:1.7
#
# Otter API — production multi-stage build.
# DD-203: distroless nonroot + layer caching + minimum attack surface.
# DD-232: dependencies pinned by go.sum, content pack baked in, self-probing
# HEALTHCHECK (distroless has no shell or wget).
# Build context: project root (.dockerignore keeps it to server/ + content/).

# ---- dependencies ----
FROM golang:1.25-alpine AS deps

WORKDIR /src/server

# Module files first so the download layer is cached; go.sum pins every
# dependency and `verify` checks the module cache against it.
COPY server/go.mod server/go.sum ./
RUN go mod download && go mod verify

# ---- build ----
FROM deps AS build

COPY server/ ./

ARG COMMIT_SHA=dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/api ./cmd/api

# ---- runtime ----
FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /out/api /app
# The content pack (items, mobs, names, legal pages) ships with the image so
# a deploy is one immutable artifact; dev mounts ../content over it.
COPY content/ /content/
ENV CONTENT_PATH=/content

USER nonroot:nonroot
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=4s --retries=5 CMD ["/app", "healthcheck"]
ENTRYPOINT ["/app"]
